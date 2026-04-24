package service

import (
	"context"
	"fmt"

	"gateway-svc/adapter/finance_grpc_adapter"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// FinancePing dispatches the permission-gate smoke call to finance-svc via
// gRPC. Authorization was already enforced by the gateway's RequireBearerToken
// + RequirePermission middleware before this method runs; finance-svc is
// identity-agnostic and the request is empty. Transport errors bubble up as
// wrapped apperrors.ErrServiceUnavailable (502) per the adapter's fail-closed
// mapping.
func (s *Service) FinancePing(ctx context.Context) (*finance_grpc_adapter.PingResult, error) {
	const op = "service.Service.FinancePing"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("operation", op))
	logger.Info().Str("op", op).Msg("")

	result, err := s.adapters.financeGrpc.Ping(ctx)
	if err != nil {
		err = fmt.Errorf("call finance-svc: %w", err)
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}
