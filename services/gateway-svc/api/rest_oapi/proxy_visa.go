// proxy_visa.go — gateway REST handlers for visa pipeline Phase 6 routes
// (BL-VISA-001..003).
//
// Route topology (all bearer-protected):
//   PUT  /v1/visas/:id/status  → TransitionStatus
//   POST /v1/visas/bulk-submit → BulkSubmit
//   GET  /v1/visas             → GetApplications
//
// Per ADR-0009: gateway is the single REST entry-point; visa-svc is pure gRPC.
package rest_oapi

import (
	"errors"

	"gateway-svc/adapter/visa_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ---------------------------------------------------------------------------
// Request body types
// ---------------------------------------------------------------------------

type TransitionVisaStatusBody struct {
	ToStatus    string `json:"to_status"`
	Reason      string `json:"reason,omitempty"`
	ActorUserID string `json:"actor_user_id,omitempty"`
}

type BulkSubmitVisaBody struct {
	DepartureID string   `json:"departure_id"`
	JamaahIDs   []string `json:"jamaah_ids"`
	ProviderID  string   `json:"provider_id,omitempty"`
}

// ---------------------------------------------------------------------------
// TransitionVisaStatus — PUT /v1/visas/:id/status (bearer)
// ---------------------------------------------------------------------------

func (s *Server) TransitionVisaStatus(c *fiber.Ctx, applicationID string) error {
	const op = "rest_oapi.Server.TransitionVisaStatus"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("application_id", applicationID))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("application_id", applicationID).Msg("")

	var body TransitionVisaStatusBody
	if err := c.BodyParser(&body); err != nil {
		return writeVisaError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.TransitionVisaStatus(ctx, &visa_grpc_adapter.TransitionStatusParams{
		ApplicationID: applicationID,
		ToStatus:      body.ToStatus,
		Reason:        body.Reason,
		ActorUserID:   body.ActorUserID,
	})
	if err != nil {
		return writeVisaError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"application_id": result.ApplicationID,
			"from_status":    result.FromStatus,
			"to_status":      result.ToStatus,
			"idempotent":     result.Idempotent,
		},
	})
}

// ---------------------------------------------------------------------------
// BulkSubmitVisa — POST /v1/visas/bulk-submit (bearer)
// ---------------------------------------------------------------------------

func (s *Server) BulkSubmitVisa(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.BulkSubmitVisa"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body BulkSubmitVisaBody
	if err := c.BodyParser(&body); err != nil {
		return writeVisaError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.BulkSubmitVisa(ctx, &visa_grpc_adapter.BulkSubmitParams{
		DepartureID: body.DepartureID,
		JamaahIDs:   body.JamaahIDs,
		ProviderID:  body.ProviderID,
	})
	if err != nil {
		return writeVisaError(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"submitted_count": result.SubmittedCount,
			"application_ids": result.ApplicationIDs,
		},
	})
}

// ---------------------------------------------------------------------------
// GetVisaApplications — GET /v1/visas (bearer)
// ---------------------------------------------------------------------------

func (s *Server) GetVisaApplications(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetVisaApplications"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	departureID := c.Query("departure_id")
	statusFilter := c.Query("status")

	result, err := s.svc.GetVisaApplications(ctx, departureID, statusFilter)
	if err != nil {
		return writeVisaError(c, span, err)
	}

	type histEntryJSON struct {
		FromStatus string `json:"from_status"`
		ToStatus   string `json:"to_status"`
		Reason     string `json:"reason,omitempty"`
		CreatedAt  string `json:"created_at"`
	}
	type appJSON struct {
		ID          string          `json:"id"`
		JamaahID    string          `json:"jamaah_id"`
		Status      string          `json:"status"`
		ProviderRef string          `json:"provider_ref,omitempty"`
		IssuedDate  string          `json:"issued_date,omitempty"`
		History     []histEntryJSON `json:"history,omitempty"`
	}

	data := make([]appJSON, 0, len(result.Applications))
	for _, a := range result.Applications {
		hist := make([]histEntryJSON, 0, len(a.History))
		for _, h := range a.History {
			hist = append(hist, histEntryJSON{
				FromStatus: h.FromStatus,
				ToStatus:   h.ToStatus,
				Reason:     h.Reason,
				CreatedAt:  h.CreatedAt,
			})
		}
		data = append(data, appJSON{
			ID:          a.ID,
			JamaahID:    a.JamaahID,
			Status:      a.Status,
			ProviderRef: a.ProviderRef,
			IssuedDate:  a.IssuedDate,
			History:     hist,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": data})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeVisaError(c *fiber.Ctx, span trace.Span, err error) error {
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
		message = "layanan visa sementara tidak tersedia"
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
