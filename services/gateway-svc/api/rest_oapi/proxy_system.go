package rest_oapi

import (
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetSystemBackends renders the aggregate backend-health payload consumed by
// core-web's status page. The REST response shape mirrors
// grpc.health.v1.Health.Check per-backend — SERVING / NOT_SERVING /
// UNKNOWN — so dashboards branch on a stable, protocol-native vocabulary.
// Unauthenticated: the status page is reachable pre-login.
//
// GetSystemBackends implements ServerInterface.
func (s *Server) GetSystemBackends(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetSystemBackends"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/system/backends"),
		attribute.String("method", "GET"),
	)
	logger.Info().Str("op", op).Msg("")

	backends := s.svc.SystemBackends(ctx)

	resp := SystemBackendsResponse{}
	resp.Data.Backends = make([]BackendStatus, 0, len(backends))
	for _, b := range backends {
		entry := BackendStatus{
			Name:   b.Name,
			Status: BackendStatusStatus(b.Status),
		}
		if b.Error != "" {
			e := b.Error
			entry.Error = &e
		}
		resp.Data.Backends = append(resp.Data.Backends, entry)
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}
