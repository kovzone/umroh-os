package rest_oapi

import (
	"fmt"

	"booking-svc/service"
	"booking-svc/util/apperrors"
	"booking-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Liveness verifies that the service process is alive and responsive.
// Must not access any external dependencies; trivial and deterministic.
//
// No logging or tracing — probe endpoints are hit at high frequency; kept cheap.
// Liveness implements ServerInterface.
func (s *Server) Liveness(c *fiber.Ctx) error {
	result, err := s.svc.Liveness(c.Context(), &service.LivenessParams{})
	if err != nil {
		return fiber.NewError(apperrors.HTTPStatus(err), err.Error())
	}

	response := LiveResponse{
		Data: struct {
			Ok bool `json:"ok"`
		}{Ok: result.OK},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Readiness verifies that the service instance can handle real traffic.
// Checks DB via read-only SELECT 1; no state-mutating operations.
//
// No logging or tracing — probe endpoints are hit at high frequency; kept cheap.
// Readiness implements ServerInterface.
func (s *Server) Readiness(c *fiber.Ctx) error {
	result, err := s.svc.Readiness(c.Context(), &service.ReadinessParams{})
	if err != nil {
		return fiber.NewError(apperrors.HTTPStatus(err), err.Error())
	}

	response := ReadyResponse{
		Data: struct {
			Ok bool `json:"ok"`
		}{Ok: result.OK},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// DbTxDiagnostic executes a state-mutating DB transaction to validate the WithTx pattern.
// Not for health/readiness probes — diagnostic and reference implementation.
//
// DbTxDiagnostic implements ServerInterface.
func (s *Server) DbTxDiagnostic(c *fiber.Ctx, params DbTxDiagnosticParams) error {
	const op = "rest_oapi.Server.DbTxDiagnostic"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/system/diagnostics/db-tx"),
		attribute.String("method", "GET"),
	)

	message := "no message"
	if params.Message != nil && *params.Message != "" {
		message = *params.Message
	}

	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.DbTxDiagnostic(ctx, &service.DbTxDiagnosticParams{Message: message})
	if err != nil {
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fiber.NewError(apperrors.HTTPStatus(err), err.Error())
	}

	span.SetAttributes(
		attribute.String("output.response", fmt.Sprintf("diagnostic_id=%d message=%s", result.ID, message)),
	)
	span.SetStatus(codes.Ok, "success")

	response := DbTxDiagnosticResponse{
		Data: struct {
			DiagnosticId int64  `json:"diagnostic_id"`
			Message      string `json:"message"`
		}{
			DiagnosticId: result.ID,
			Message:      fmt.Sprintf("received '%s' from client", message),
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
