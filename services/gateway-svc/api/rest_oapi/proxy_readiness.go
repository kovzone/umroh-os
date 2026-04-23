// proxy_readiness.go — gateway REST handlers for vendor readiness endpoints
// (BL-OPS-020).
//
// Route topology (all bearer-protected):
//   GET /v1/departures/:id/readiness — get current ticket/hotel/visa states
//   PUT /v1/departures/:id/readiness — update one readiness kind
//
// Per ADR-0009: gateway is the single REST entry-point; catalog-svc is pure gRPC.
package rest_oapi

import (
	"errors"

	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/api/rest_oapi/middleware"
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

type UpdateDepartureReadinessBody struct {
	Kind          string `json:"kind"`           // ticket | hotel | visa
	State         string `json:"state"`          // not_started | in_progress | done
	Notes         string `json:"notes"`
	AttachmentURL string `json:"attachment_url"`
}

// ---------------------------------------------------------------------------
// GetDepartureReadiness — GET /v1/departures/:id/readiness (bearer)
// ---------------------------------------------------------------------------

func (s *Server) GetDepartureReadiness(c *fiber.Ctx, departureID string) error {
	const op = "rest_oapi.Server.GetDepartureReadiness"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("departure_id", departureID),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	id := c.Locals(middleware.IdentityKey).(*middleware.Identity)
	logger.Info().Str("op", op).Str("user_id", id.UserID).Str("departure_id", departureID).Msg("")

	result, err := s.svc.GetDepartureReadiness(ctx, &catalog_grpc_adapter.GetDepartureReadinessParams{
		UserID:      id.UserID,
		DepartureID: departureID,
	})
	if err != nil {
		return writeReadinessError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"departure_id": departureID,
			"ticket":       result.TicketState,
			"hotel":        result.HotelState,
			"visa":         result.VisaState,
		},
	})
}

// ---------------------------------------------------------------------------
// UpdateDepartureReadiness — PUT /v1/departures/:id/readiness (bearer)
// ---------------------------------------------------------------------------

func (s *Server) UpdateDepartureReadiness(c *fiber.Ctx, departureID string) error {
	const op = "rest_oapi.Server.UpdateDepartureReadiness"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("departure_id", departureID),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	id := c.Locals(middleware.IdentityKey).(*middleware.Identity)
	logger.Info().Str("op", op).Str("user_id", id.UserID).Str("departure_id", departureID).Msg("")

	var body UpdateDepartureReadinessBody
	if err := c.BodyParser(&body); err != nil {
		return writeReadinessError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.UpdateVendorReadiness(ctx, &catalog_grpc_adapter.UpdateVendorReadinessParams{
		UserID:        id.UserID,
		DepartureID:   departureID,
		Kind:          body.Kind,
		State:         body.State,
		Notes:         body.Notes,
		AttachmentURL: body.AttachmentURL,
	})
	if err != nil {
		return writeReadinessError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"departure_id": departureID,
			"ticket":       result.TicketState,
			"hotel":        result.HotelState,
			"visa":         result.VisaState,
		},
	})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeReadinessError(c *fiber.Ctx, span trace.Span, err error) error {
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
		message = "layanan catalog sementara tidak tersedia"
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
