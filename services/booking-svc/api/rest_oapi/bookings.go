package rest_oapi

// S1-E-03 / BL-BOOK-001..006 — REST handler for booking draft creation.
//
// POST /v1/bookings
//
// Auth: the gateway-svc validates the Bearer token and injects X-User-Id.
// This handler trusts that header (no direct iam-svc call at this layer).
// The service layer calls catalog-svc gRPC to validate departure and reserve seats.

import (
	"errors"
	"fmt"

	"booking-svc/service"
	"booking-svc/util/apperrors"
	"booking-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// CreateDraftBooking handles POST /v1/bookings.
// Implements ServerInterface.
func (s *Server) CreateDraftBooking(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateDraftBooking"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "POST /v1/bookings"),
	)

	// --- Parse request body ---
	var body CreateBookingRequestBody
	if err := c.BodyParser(&body); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "bad request body")
		return writeBookingError(c, errors.Join(apperrors.ErrValidation,
			fmt.Errorf("request body tidak valid: %w", err)))
	}

	// --- Read Idempotency-Key header (optional) ---
	idempotencyKey := c.Get("Idempotency-Key", "")

	// --- Map request body → service params ---
	params := &service.CreateDraftBookingParams{
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
	}

	// Map pilgrims
	for _, p := range body.Pilgrims {
		params.Pilgrims = append(params.Pilgrims, service.PilgrimInput{
			FullName:    p.FullName,
			Email:       p.Email,
			Whatsapp:    p.Whatsapp,
			Domicile:    p.Domicile,
			IsLead:      p.IsLead,
			HasKTP:      p.HasKTP,
			HasPassport: p.HasPassport,
		})
	}

	// Map add-ons
	for _, a := range body.Addons {
		params.Addons = append(params.Addons, service.AddonInput{
			AddonID:            a.AddonID,
			AddonName:          a.AddonName,
			ListAmount:         a.ListAmount,
			ListCurrency:       a.ListCurrency,
			SettlementCurrency: a.SettlementCurrency,
		})
	}

	// Map mahram (BL-BOOK-006)
	if body.Mahram != nil {
		params.Mahram = &service.MahramInfo{
			MahramWarning: body.Mahram.MahramWarning,
		}
	}

	logger.Info().
		Str("op", op).
		Str("channel", params.Channel).
		Str("package_id", params.PackageID).
		Str("departure_id", params.DepartureID).
		Str("idempotency_key", idempotencyKey).
		Msg("")

	// --- Call service layer ---
	result, err := s.svc.CreateDraftBooking(ctx, params)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return writeBookingError(c, err)
	}

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(attribute.String("output.booking_id", result.ID))

	// --- Map result → response ---
	resp := CreateBookingResponse{
		Data: draftResultToBody(result),
	}

	// 200 for idempotency replay, 201 for new booking.
	status := fiber.StatusCreated
	if result.Replayed {
		status = fiber.StatusOK
	}

	return c.Status(status).JSON(resp)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// draftResultToBody maps service.DraftBookingResult to BookingDataBody.
func draftResultToBody(r *service.DraftBookingResult) BookingDataBody {
	items := make([]BookingItemResponseBody, 0, len(r.Items))
	for _, it := range r.Items {
		items = append(items, BookingItemResponseBody{
			ID:              it.ID,
			FullName:        it.FullName,
			IsLead:          it.IsLead,
			DocumentWarning: it.DocumentWarning,
		})
	}
	return BookingDataBody{
		ID:                 r.ID,
		Status:             r.Status,
		Channel:            r.Channel,
		PackageID:          r.PackageID,
		DepartureID:        r.DepartureID,
		RoomType:           r.RoomType,
		AgentID:            r.AgentID,
		StaffUserID:        r.StaffUserID,
		LeadFullName:       r.LeadFullName,
		LeadEmail:          r.LeadEmail,
		LeadWhatsapp:       r.LeadWhatsapp,
		LeadDomicile:       r.LeadDomicile,
		ListAmount:         r.ListAmount,
		ListCurrency:       r.ListCurrency,
		SettlementCurrency: r.SettlementCurrency,
		Notes:              r.Notes,
		IdempotencyKey:     r.IdempotencyKey,
		CreatedAt:          r.CreatedAt,
		ExpiresAt:          r.ExpiresAt,
		Items:              items,
		MahramWarning:      r.MahramWarning,
		Replayed:           r.Replayed,
	}
}

// writeBookingError maps a domain/service error to a booking-svc error envelope.
func writeBookingError(c *fiber.Ctx, err error) error {
	status := apperrors.HTTPStatus(err)
	code := apperrors.ErrorCode(err)
	return c.Status(status).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": err.Error(),
		},
	})
}
