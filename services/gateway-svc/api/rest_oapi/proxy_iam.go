package rest_oapi

import (
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetIamSystemLive proxies the iam-svc liveness probe through the gateway.
// Scaffold-time proof of the REST adapter pattern.
//
// GetIamSystemLive implements ServerInterface.
func (s *Server) GetIamSystemLive(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetIamSystemLive"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/iam/system/live"),
		attribute.String("method", "GET"),
	)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetIamSystemLive(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "BAD_GATEWAY",
				"message": err.Error(),
			},
		})
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(LiveResponse{
		Data: struct {
			Ok bool `json:"ok"`
		}{Ok: result.OK},
	})
}

// GetIamSystemDbTxDiagnostic proxies iam-svc's state-mutating DB transaction
// diagnostic through the gateway. It's the traced cross-service path the
// S0-J-05 observability acceptance uses: one request produces one trace
// spanning gateway-svc + iam-svc, with matching trace_id log lines in both
// containers' Loki streams.
//
// GetIamSystemDbTxDiagnostic implements ServerInterface.
func (s *Server) GetIamSystemDbTxDiagnostic(c *fiber.Ctx, params GetIamSystemDbTxDiagnosticParams) error {
	const op = "rest_oapi.Server.GetIamSystemDbTxDiagnostic"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	message := "no message"
	if params.Message != nil {
		message = *params.Message
	}

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/iam/system/diagnostics/db-tx"),
		attribute.String("method", "GET"),
		attribute.String("message", message),
	)
	logger.Info().Str("op", op).Str("message", message).Msg("")

	result, err := s.svc.GetIamSystemDbTxDiagnostic(ctx, message)
	if err != nil {
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "BAD_GATEWAY",
				"message": err.Error(),
			},
		})
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(DbTxDiagnosticResponse{
		Data: struct {
			DiagnosticId int64  `json:"diagnostic_id"`
			Message      string `json:"message"`
		}{
			DiagnosticId: result.DiagnosticID,
			Message:      result.Message,
		},
	})
}
