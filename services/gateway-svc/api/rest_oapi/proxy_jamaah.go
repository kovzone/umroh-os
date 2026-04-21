package rest_oapi

import (
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetJamaahSystemLive proxies the iam-svc liveness probe through the gateway.
// Scaffold-time proof of the REST adapter pattern.
//
// GetJamaahSystemLive implements ServerInterface.
func (s *Server) GetJamaahSystemLive(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetJamaahSystemLive"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/jamaah/system/live"),
		attribute.String("method", "GET"),
	)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetJamaahSystemLive(ctx)
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
