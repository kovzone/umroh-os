// proxy_crm.go — gateway REST handlers for CRM lead management (S4-E-02).
//
// Route topology:
//   POST   /v1/leads        → CreateLead  (public — no bearer required for lead capture)
//   GET    /v1/leads        → ListLeads   (bearer, cs/admin)
//   GET    /v1/leads/:id    → GetLead     (bearer)
//   PUT    /v1/leads/:id    → UpdateLead  (bearer)
//
// Per ADR-0009: gateway is the single REST entry-point; crm-svc is pure gRPC.
// Auth for protected routes is enforced by the RequireBearerToken middleware
// registered in cmd/server.go.

package rest_oapi

import (
	"errors"

	"gateway-svc/adapter/crm_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ---------------------------------------------------------------------------
// Request / response body types
// ---------------------------------------------------------------------------

// CreateLeadBody is the JSON body for POST /v1/leads.
// Also accepts field aliases sent by UAT spec / landing pages (ISSUE-017):
//   message     → merged into notes if notes is empty
//   source_note → used as source if source is empty
type CreateLeadBody struct {
	Name                string `json:"name"`
	Phone               string `json:"phone"`
	Email               string `json:"email,omitempty"`
	Source              string `json:"source,omitempty"`
	SourceNote          string `json:"source_note,omitempty"` // alias for source
	UtmSource           string `json:"utm_source,omitempty"`
	UtmMedium           string `json:"utm_medium,omitempty"`
	UtmCampaign         string `json:"utm_campaign,omitempty"`
	UtmContent          string `json:"utm_content,omitempty"`
	UtmTerm             string `json:"utm_term,omitempty"`
	InterestPackageID   string `json:"interest_package_id,omitempty"`
	InterestDepartureID string `json:"interest_departure_id,omitempty"`
	Notes               string `json:"notes,omitempty"`
	Message             string `json:"message,omitempty"` // alias for notes
}

// UpdateLeadBody is the JSON body for PUT /v1/leads/:id.
type UpdateLeadBody struct {
	Status       string `json:"status,omitempty"`
	Notes        string `json:"notes,omitempty"`
	AssignedCsID string `json:"assigned_cs_id,omitempty"`
}

// LeadResponseData is the JSON data envelope for a single lead.
type LeadResponseData struct {
	ID                  string `json:"id"`
	Source              string `json:"source"`
	UtmSource           string `json:"utm_source,omitempty"`
	UtmMedium           string `json:"utm_medium,omitempty"`
	UtmCampaign         string `json:"utm_campaign,omitempty"`
	UtmContent          string `json:"utm_content,omitempty"`
	UtmTerm             string `json:"utm_term,omitempty"`
	Name                string `json:"name"`
	Phone               string `json:"phone"`
	Email               string `json:"email,omitempty"`
	InterestPackageID   string `json:"interest_package_id,omitempty"`
	InterestDepartureID string `json:"interest_departure_id,omitempty"`
	Status              string `json:"status"`
	AssignedCsID        string `json:"assigned_cs_id,omitempty"`
	Notes               string `json:"notes,omitempty"`
	BookingID           string `json:"booking_id,omitempty"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

// SingleLeadResponse wraps LeadResponseData in the standard envelope.
type SingleLeadResponse struct {
	Data LeadResponseData `json:"data"`
}

// LeadListResponse is the paginated list response for GET /v1/leads.
type LeadListResponse struct {
	Data     []LeadResponseData `json:"data"`
	Total    int64              `json:"total"`
	Page     int32              `json:"page"`
	PageSize int32              `json:"page_size"`
}

// ---------------------------------------------------------------------------
// CreateLead — POST /v1/leads (public)
// ---------------------------------------------------------------------------

// TODO(security): POST /v1/leads is a public endpoint (no auth required).
// Rate limiting should be applied to prevent abuse (e.g., lead spam, phone number
// enumeration). Recommended: per-IP token-bucket limit (e.g., 10 req/min/IP) via
// a Fiber rate-limiter middleware or an upstream proxy (nginx/Caddy). Track as
// BL-SEC-001 or Phase 6 hardening card.
func (s *Server) CreateLead(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateLead"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "POST /v1/leads"),
	)

	var body CreateLeadBody
	if err := c.BodyParser(&body); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "bad request body")
		return writeCrmError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	logger.Info().Str("op", op).Str("phone", body.Phone).Msg("")

	// Merge field aliases (ISSUE-017)
	if body.Notes == "" && body.Message != "" {
		body.Notes = body.Message
	}
	if body.Source == "" && body.SourceNote != "" {
		body.Source = body.SourceNote
	}

	result, err := s.svc.CreateLead(ctx, &crm_grpc_adapter.CreateLeadParams{
		Name:                body.Name,
		Phone:               body.Phone,
		Email:               body.Email,
		Source:              body.Source,
		UtmSource:           body.UtmSource,
		UtmMedium:           body.UtmMedium,
		UtmCampaign:         body.UtmCampaign,
		UtmContent:          body.UtmContent,
		UtmTerm:             body.UtmTerm,
		InterestPackageID:   body.InterestPackageID,
		InterestDepartureID: body.InterestDepartureID,
		Notes:               body.Notes,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeCrmError(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	span.SetAttributes(attribute.String("lead_id", result.ID))

	return c.Status(fiber.StatusCreated).JSON(SingleLeadResponse{Data: mapLeadResult(result)})
}

// ---------------------------------------------------------------------------
// ListLeads — GET /v1/leads (bearer, cs/admin)
// ---------------------------------------------------------------------------

func (s *Server) ListLeads(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListLeads"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	page := int32(c.QueryInt("page", 1))
	pageSize := int32(c.QueryInt("page_size", 20))
	statusFilter := c.Query("status")
	assignedCsFilter := c.Query("assigned_cs_id")

	result, err := s.svc.ListLeads(ctx, &crm_grpc_adapter.ListLeadsParams{
		StatusFilter:     statusFilter,
		AssignedCsFilter: assignedCsFilter,
		Page:             page,
		PageSize:         pageSize,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return writeCrmError(c, span, err)
	}

	data := make([]LeadResponseData, 0, len(result.Leads))
	for _, l := range result.Leads {
		data = append(data, mapLeadResult(l))
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(LeadListResponse{
		Data:     data,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	})
}

// ---------------------------------------------------------------------------
// GetLeadByID — GET /v1/leads/:id (bearer)
// ---------------------------------------------------------------------------

func (s *Server) GetLeadByID(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.GetLeadByID"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	span.SetAttributes(attribute.String("lead_id", id))

	result, err := s.svc.GetLead(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return writeCrmError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(SingleLeadResponse{Data: mapLeadResult(result)})
}

// ---------------------------------------------------------------------------
// UpdateLeadByID — PUT /v1/leads/:id (bearer)
// ---------------------------------------------------------------------------

// TODO(security/known-limitation): assigned_cs_id in the request body is NOT cross-validated
// against IAM to confirm the target user has role=CS. crm-svc stores any UUID provided.
// Full validation (gRPC call to iam-svc.GetUser to confirm role=CS) is deferred to Phase 6
// BL-CRM-003. Current risk: an admin could assign a lead to a non-CS user UUID.
func (s *Server) UpdateLeadByID(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.UpdateLeadByID"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	span.SetAttributes(attribute.String("lead_id", id))

	var body UpdateLeadBody
	if err := c.BodyParser(&body); err != nil {
		return writeCrmError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.UpdateLead(ctx, &crm_grpc_adapter.UpdateLeadParams{
		ID:           id,
		Status:       body.Status,
		Notes:        body.Notes,
		AssignedCsID: body.AssignedCsID,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return writeCrmError(c, span, err)
	}

	span.SetStatus(codes.Ok, "updated")
	return c.Status(fiber.StatusOK).JSON(SingleLeadResponse{Data: mapLeadResult(result)})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeCrmError(c *fiber.Ctx, span trace.Span, err error) error {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	var httpStatus int
	var code, message string

	switch {
	case errors.Is(err, apperrors.ErrValidation):
		httpStatus = fiber.StatusBadRequest
		code = "validation_error"
		message = err.Error()
	case errors.Is(err, apperrors.ErrNotFound):
		httpStatus = fiber.StatusNotFound
		code = "not_found"
		message = err.Error()
	case errors.Is(err, apperrors.ErrConflict):
		httpStatus = fiber.StatusConflict
		code = "conflict"
		message = err.Error()
	case errors.Is(err, apperrors.ErrUnauthorized):
		httpStatus = fiber.StatusUnauthorized
		code = "unauthorized"
		message = "autentikasi diperlukan"
	case errors.Is(err, apperrors.ErrForbidden):
		httpStatus = fiber.StatusForbidden
		code = "forbidden"
		message = "akses ditolak"
	case errors.Is(err, apperrors.ErrServiceUnavailable):
		httpStatus = fiber.StatusBadGateway
		code = "service_unavailable"
		message = "layanan CRM sementara tidak tersedia"
	default:
		httpStatus = fiber.StatusInternalServerError
		code = "internal_error"
		message = "terjadi kesalahan tidak terduga"
	}

	return c.Status(httpStatus).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": message,
		},
	})
}

// ---------------------------------------------------------------------------
// mapLeadResult — adapter result → response data
// ---------------------------------------------------------------------------

func mapLeadResult(l *crm_grpc_adapter.LeadResult) LeadResponseData {
	return LeadResponseData{
		ID:                  l.ID,
		Source:              l.Source,
		UtmSource:           l.UtmSource,
		UtmMedium:           l.UtmMedium,
		UtmCampaign:         l.UtmCampaign,
		UtmContent:          l.UtmContent,
		UtmTerm:             l.UtmTerm,
		Name:                l.Name,
		Phone:               l.Phone,
		Email:               l.Email,
		InterestPackageID:   l.InterestPackageID,
		InterestDepartureID: l.InterestDepartureID,
		Status:              l.Status,
		AssignedCsID:        l.AssignedCsID,
		Notes:               l.Notes,
		BookingID:           l.BookingID,
		CreatedAt:           l.CreatedAt,
		UpdatedAt:           l.UpdatedAt,
	}
}
