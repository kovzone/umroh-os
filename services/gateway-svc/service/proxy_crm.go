package service

import (
	"context"
	"fmt"

	"gateway-svc/adapter/crm_rest_adapter"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetCrmSystemLive proxies a liveness probe to crm-svc through the REST
// adapter. It's the scaffold-time proof that the REST adapter pattern works
// end-to-end (typed call, span propagation via otelhttp, error handling).
//
// Future iam proxy methods land here as the gateway exposes them.
func (s *Service) GetCrmSystemLive(ctx context.Context) (*crm_rest_adapter.LivenessResult, error) {
	const op = "service.Service.GetCrmSystemLive"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("operation", op))
	logger.Info().Str("op", op).Msg("")

	result, err := s.adapters.crmRest.GetSystemLive(ctx)
	if err != nil {
		err = fmt.Errorf("call crm-svc: %w", err)
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}
