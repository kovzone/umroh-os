package service

import (
	"context"
	"fmt"

	"gateway-svc/adapter/iam_rest_adapter"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetIamSystemLive proxies a liveness probe to iam-svc through the REST
// adapter. It's the scaffold-time proof that the REST adapter pattern works
// end-to-end (typed call, span propagation via otelhttp, error handling).
//
// Future iam proxy methods land here as the gateway exposes them.
func (s *Service) GetIamSystemLive(ctx context.Context) (*iam_rest_adapter.LivenessResult, error) {
	const op = "service.Service.GetIamSystemLive"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("operation", op))
	logger.Info().Str("op", op).Msg("")

	result, err := s.adapters.iamRest.GetSystemLive(ctx)
	if err != nil {
		err = fmt.Errorf("call iam-svc: %w", err)
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// GetIamSystemDbTxDiagnostic proxies the iam-svc DB transaction diagnostic.
// Its purpose is verification: gateway → iam produces one trace spanning both
// services, with matching trace_id log lines in each container. See S0-J-05.
func (s *Service) GetIamSystemDbTxDiagnostic(ctx context.Context, message string) (*iam_rest_adapter.DbTxDiagnosticResult, error) {
	const op = "service.Service.GetIamSystemDbTxDiagnostic"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("message", message),
	)
	logger.Info().Str("op", op).Str("message", message).Msg("")

	result, err := s.adapters.iamRest.GetSystemDbTxDiagnostic(ctx, message)
	if err != nil {
		err = fmt.Errorf("call iam-svc: %w", err)
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}
