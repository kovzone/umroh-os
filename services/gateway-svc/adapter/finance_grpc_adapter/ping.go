package finance_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/finance_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Ping delegates to finance-svc.Ping over gRPC. Request is empty — per ADR
// 0009 single-point auth, authorization was decided at the gateway's
// RequirePermission middleware before this call; finance-svc does not need
// caller identity to respond.
//
// On transport failure (finance unreachable / timeout), mapFinanceError
// produces ErrServiceUnavailable → 502 so the permission-gate smoke detects
// a downed backend per F1-W7 fail-closed.
func (a *Adapter) Ping(ctx context.Context) (*PingResult, error) {
	const op = "finance_grpc_adapter.Adapter.Ping"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "Ping"),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.financeClient.Ping(ctx, &pb.PingRequest{})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &PingResult{Message: resp.GetMessage()}, nil
}
