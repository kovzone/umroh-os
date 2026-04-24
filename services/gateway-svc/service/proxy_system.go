package service

import (
	"context"

	"gateway-svc/adapter/health_check_adapter"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// SystemBackends delegates to the health_check_adapter to probe every
// registered backend's grpc.health.v1.Health.Check concurrently. The
// adapter owns the per-probe timeout + concurrency; this method adds the
// service-layer span and logging so Tempo/Loki can correlate a status-page
// poll with each fan-out probe.
func (s *Service) SystemBackends(ctx context.Context) []health_check_adapter.BackendStatus {
	const op = "service.Service.SystemBackends"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("operation", op))
	logger.Info().Str("op", op).Msg("")

	result := s.adapters.healthCheck.CheckAll(ctx)

	span.SetStatus(codes.Ok, "success")
	return result
}
