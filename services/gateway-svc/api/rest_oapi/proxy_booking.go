package rest_oapi

// S1-E-03 / BL-BOOK-001..006 / BL-GTW-003 — gateway REST handler for
// draft booking creation.
//
// POST /v1/bookings — public in S1 (bearer auth arrives with F4 later).
//
// The handler parses the REST request body into adapter params,
// forwards to booking-svc via booking_grpc_adapter, and maps the
// result back to a REST response.
//
// Per ADR-0009: gateway is the single REST entry-point; booking-svc is
// pure gRPC. No auth check in S1 (field channel already encodes the
// booking source; proper auth gates are a later F4 task).

import (
	"errors"

	"gateway-svc/adapter/booking_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ---------------------------------------------------------------------------
// Request / response body types (hand-written; would be oapi-generated)
// ---------------------------------------------------------------------------

// PilgrimBody is one jamaah in the booking request body.
type PilgrimBody struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email,omitempty"`
	Whatsapp    string `json:"whatsapp,omitempty"`
	Domicile    string `json:"domicile,omitempty"`
	IsLead      bool   `json:"is_lead,omitempty"`
	HasKTP      bool   `json:"has_ktp,omitempty"`
	HasPassport bool   `json:"has_passport,omitempty"`
}

// AddonBody is a selected add-on in the booking request body.
type AddonBody struct {
	AddonID            string `json:"addon_id"`
	AddonName          string `json:"addon_name,omitempty"`
	ListAmount         int64  `json:"list_amount,omitempty"`
	ListCurrency       string `json:"list_currency,omitempty"`
	SettlementCurrency string `json:"settlement_currency,omitempty"`
}

// LeadBody is the nested lead object in the booking request body (ISSUE-011 alias).
type LeadBody struct {
	FullName  string `json:"full_name"`
	Email     string `json:"email,omitempty"`
	Whatsapp  string `json:"whatsapp,omitempty"`
	Domicile  string `json:"domicile,omitempty"`
}

// CreateDraftBookingBody is the JSON body for POST /v1/bookings.
// Supports both flat fields (lead_full_name, pilgrims) and nested aliases
// (lead, jamaah) per ISSUE-011 — the UAT spec sends the nested form.
type CreateDraftBookingBody struct {
	Channel            string        `json:"channel"`
	AgentID            string        `json:"agent_id,omitempty"`
	StaffUserID        string        `json:"staff_user_id,omitempty"`
	PackageID          string        `json:"package_id"`
	DepartureID        string        `json:"departure_id"`
	RoomType           string        `json:"room_type"`
	// Flat form (legacy / internal)
	LeadFullName       string        `json:"lead_full_name"`
	LeadEmail          string        `json:"lead_email,omitempty"`
	LeadWhatsapp       string        `json:"lead_whatsapp"`
	LeadDomicile       string        `json:"lead_domicile"`
	Pilgrims           []PilgrimBody `json:"pilgrims,omitempty"`
	// Nested form (UAT spec / B2C form) — ISSUE-011 alias
	Lead               *LeadBody     `json:"lead,omitempty"`
	Jamaah             []PilgrimBody `json:"jamaah,omitempty"`
	Addons             []AddonBody   `json:"addons,omitempty"`
	ListAmount         int64         `json:"list_amount,omitempty"`
	ListCurrency       string        `json:"list_currency,omitempty"`
	SettlementCurrency string        `json:"settlement_currency,omitempty"`
	Notes              string        `json:"notes,omitempty"`
	MahramWarning      string        `json:"mahram_warning,omitempty"`
}

// BookingItemResponseBody is one pilgrim in the booking response.
type BookingItemResponseBody struct {
	ID              string `json:"id"`
	FullName        string `json:"full_name"`
	IsLead          bool   `json:"is_lead"`
	DocumentWarning string `json:"document_warning,omitempty"`
}

// CreateDraftBookingResponseData is the data envelope in the booking response.
type CreateDraftBookingResponseData struct {
	ID                 string                    `json:"id"`
	Status             string                    `json:"status"`
	Channel            string                    `json:"channel"`
	PackageID          string                    `json:"package_id"`
	DepartureID        string                    `json:"departure_id"`
	RoomType           string                    `json:"room_type"`
	AgentID            string                    `json:"agent_id,omitempty"`
	StaffUserID        string                    `json:"staff_user_id,omitempty"`
	LeadFullName       string                    `json:"lead_full_name"`
	LeadEmail          string                    `json:"lead_email,omitempty"`
	LeadWhatsapp       string                    `json:"lead_whatsapp"`
	LeadDomicile       string                    `json:"lead_domicile"`
	ListAmount         int64                     `json:"list_amount"`
	ListCurrency       string                    `json:"list_currency"`
	SettlementCurrency string                    `json:"settlement_currency"`
	Notes              string                    `json:"notes,omitempty"`
	IdempotencyKey     string                    `json:"idempotency_key,omitempty"`
	CreatedAt          string                    `json:"created_at"`
	ExpiresAt          string                    `json:"expires_at"`
	Items              []BookingItemResponseBody `json:"items"`
	MahramWarning      string                    `json:"mahram_warning,omitempty"`
	Replayed           bool                      `json:"replayed"`
}

// CreateDraftBookingResponse is the top-level response envelope.
type CreateDraftBookingResponse struct {
	Data CreateDraftBookingResponseData `json:"data"`
}

// ---------------------------------------------------------------------------
// Handler
// ---------------------------------------------------------------------------

// CreateDraftBooking implements POST /v1/bookings.
// Implements ServerInterface.CreateDraftBooking (hand-added in S1-E-03).
func (s *Server) CreateDraftBooking(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateDraftBooking"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "POST /v1/bookings"),
	)

	// --- Parse Idempotency-Key header (optional) ---
	idempotencyKey := c.Get("Idempotency-Key", "")

	// --- Parse request body ---
	var body CreateDraftBookingBody
	if err := c.BodyParser(&body); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "bad request body")
		return writeBookingError(c, span, errors.Join(apperrors.ErrValidation,
			err))
	}

	logger.Info().
		Str("op", op).
		Str("channel", body.Channel).
		Str("package_id", body.PackageID).
		Str("departure_id", body.DepartureID).
		Str("idempotency_key", idempotencyKey).
		Msg("")

	// --- Merge nested lead/jamaah into flat fields (ISSUE-011) ---
	// Prefer flat fields if set; fall back to nested object aliases.
	if body.Lead != nil {
		if body.LeadFullName == "" {
			body.LeadFullName = body.Lead.FullName
		}
		if body.LeadEmail == "" {
			body.LeadEmail = body.Lead.Email
		}
		if body.LeadWhatsapp == "" {
			body.LeadWhatsapp = body.Lead.Whatsapp
		}
		if body.LeadDomicile == "" {
			body.LeadDomicile = body.Lead.Domicile
		}
	}
	// Merge jamaah[] alias into pilgrims if pilgrims is empty
	if len(body.Pilgrims) == 0 && len(body.Jamaah) > 0 {
		body.Pilgrims = body.Jamaah
	}

	// --- Map body → adapter params ---
	params := &booking_grpc_adapter.CreateDraftBookingParams{
		Channel:            body.Channel,
		AgentID:            body.AgentID,
		StaffUserID:        body.StaffUserID,
		PackageID:          body.PackageID,
		DepartureID:        body.DepartureID,
		RoomType:           body.RoomType,
		LeadFullName:       body.LeadFullName,
		LeadEmail:          body.LeadEmail,
		LeadWhatsapp:       body.LeadWhatsapp,
		LeadDomicile:       body.LeadDomicile,
		ListAmount:         body.ListAmount,
		ListCurrency:       body.ListCurrency,
		SettlementCurrency: body.SettlementCurrency,
		Notes:              body.Notes,
		IdempotencyKey:     idempotencyKey,
		MahramWarning:      body.MahramWarning,
	}

	for _, p := range body.Pilgrims {
		params.Pilgrims = append(params.Pilgrims, booking_grpc_adapter.PilgrimInputParam{
			FullName:    p.FullName,
			Email:       p.Email,
			Whatsapp:    p.Whatsapp,
			Domicile:    p.Domicile,
			IsLead:      p.IsLead,
			HasKTP:      p.HasKTP,
			HasPassport: p.HasPassport,
		})
	}

	for _, a := range body.Addons {
		params.Addons = append(params.Addons, booking_grpc_adapter.AddonInputParam{
			AddonID:            a.AddonID,
			AddonName:          a.AddonName,
			ListAmount:         a.ListAmount,
			ListCurrency:       a.ListCurrency,
			SettlementCurrency: a.SettlementCurrency,
		})
	}

	// --- Call service layer ---
	result, err := s.svc.CreateDraftBooking(ctx, params)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return writeBookingError(c, span, err)
	}

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(attribute.String("output.booking_id", result.ID))

	// --- Map result → response ---
	items := make([]BookingItemResponseBody, 0, len(result.Items))
	for _, it := range result.Items {
		items = append(items, BookingItemResponseBody{
			ID:              it.ID,
			FullName:        it.FullName,
			IsLead:          it.IsLead,
			DocumentWarning: it.DocumentWarning,
		})
	}

	resp := CreateDraftBookingResponse{
		Data: CreateDraftBookingResponseData{
			ID:                 result.ID,
			Status:             result.Status,
			Channel:            result.Channel,
			PackageID:          result.PackageID,
			DepartureID:        result.DepartureID,
			RoomType:           result.RoomType,
			AgentID:            result.AgentID,
			StaffUserID:        result.StaffUserID,
			LeadFullName:       result.LeadFullName,
			LeadEmail:          result.LeadEmail,
			LeadWhatsapp:       result.LeadWhatsapp,
			LeadDomicile:       result.LeadDomicile,
			ListAmount:         result.ListAmount,
			ListCurrency:       result.ListCurrency,
			SettlementCurrency: result.SettlementCurrency,
			Notes:              result.Notes,
			IdempotencyKey:     result.IdempotencyKey,
			CreatedAt:          result.CreatedAt,
			ExpiresAt:          result.ExpiresAt,
			Items:              items,
			MahramWarning:      result.MahramWarning,
			Replayed:           result.Replayed,
		},
	}

	// 200 for idempotency replay, 201 for new booking.
	httpStatus := fiber.StatusCreated
	if result.Replayed {
		httpStatus = fiber.StatusOK
	}

	return c.Status(httpStatus).JSON(resp)
}

// ---------------------------------------------------------------------------
// SubmitBooking handler
// ---------------------------------------------------------------------------

// SubmitBookingResponseData is the data envelope in the SubmitBooking response.
type SubmitBookingResponseData struct {
	ID                 string `json:"id"`
	Status             string `json:"status"`
	Channel            string `json:"channel"`
	PackageID          string `json:"package_id"`
	DepartureID        string `json:"departure_id"`
	RoomType           string `json:"room_type"`
	LeadFullName       string `json:"lead_full_name"`
	LeadWhatsapp       string `json:"lead_whatsapp"`
	LeadDomicile       string `json:"lead_domicile"`
	ListAmount         int64  `json:"list_amount"`
	ListCurrency       string `json:"list_currency"`
	SettlementCurrency string `json:"settlement_currency"`
}

// SubmitBookingResponse is the top-level response envelope.
type SubmitBookingResponseBody struct {
	Data SubmitBookingResponseData `json:"data"`
}

// SubmitBooking implements POST /v1/bookings/:id/submit (S2 / BL-BOOK-005).
// Transitions a draft booking to pending_payment.
func (s *Server) SubmitBooking(c *fiber.Ctx, bookingID string) error {
	const op = "rest_oapi.Server.SubmitBooking"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "POST /v1/bookings/:id/submit"),
		attribute.String("booking_id", bookingID),
	)

	logger.Info().Str("op", op).Str("booking_id", bookingID).Msg("")

	result, err := s.svc.SubmitBooking(ctx, &booking_grpc_adapter.SubmitBookingParams{
		BookingID: bookingID,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeBookingError(c, span, err)
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(SubmitBookingResponseBody{
		Data: SubmitBookingResponseData{
			ID:                 result.ID,
			Status:             result.Status,
			Channel:            result.Channel,
			PackageID:          result.PackageID,
			DepartureID:        result.DepartureID,
			RoomType:           result.RoomType,
			LeadFullName:       result.LeadFullName,
			LeadWhatsapp:       result.LeadWhatsapp,
			LeadDomicile:       result.LeadDomicile,
			ListAmount:         result.ListAmount,
			ListCurrency:       result.ListCurrency,
			SettlementCurrency: result.SettlementCurrency,
		},
	})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

// writeBookingError maps a domain/service error to the booking error envelope.
func writeBookingError(c *fiber.Ctx, span trace.Span, err error) error {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	var httpStatus int
	var code, message string

	switch {
	case errors.Is(err, apperrors.ErrValidation):
		httpStatus = fiber.StatusBadRequest
		code = "validation_error"
		message = err.Error()
	case errors.Is(err, apperrors.ErrConflict):
		httpStatus = fiber.StatusConflict
		code = "conflict"
		message = err.Error()
	case errors.Is(err, apperrors.ErrNotFound):
		httpStatus = fiber.StatusNotFound
		code = "not_found"
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
		message = "layanan booking sementara tidak tersedia"
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
