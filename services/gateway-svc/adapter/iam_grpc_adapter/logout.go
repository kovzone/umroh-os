package iam_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/iam_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// LogoutParams — session_id is the id the gateway middleware resolved from the
// validated bearer's payload. iam-svc revokes the row.
type LogoutParams struct {
	SessionID string
}

// LogoutResult is intentionally empty — gRPC success is the only signal. The
// gateway handler maps a nil-err return to HTTP 204.
type LogoutResult struct{}

// Logout delegates to iam-svc.Logout. Idempotent per iam-svc contract — a
// second call on the same session_id still returns ok.
func (a *Adapter) Logout(ctx context.Context, params *LogoutParams) (*LogoutResult, error) {
	const op = "iam_grpc_adapter.Adapter.Logout"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "Logout"),
		attribute.String("session_id", params.SessionID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	if _, err := a.iamClient.Logout(ctx, &pb.LogoutRequest{SessionId: params.SessionID}); err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &LogoutResult{}, nil
}
