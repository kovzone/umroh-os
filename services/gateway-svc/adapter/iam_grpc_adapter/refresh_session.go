package iam_grpc_adapter

import (
	"context"
	"time"

	"gateway-svc/adapter/iam_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// RefreshSessionParams — gateway-provided inputs for the refresh rotation.
type RefreshSessionParams struct {
	RefreshToken string
	UserAgent    string
	IP           string
}

// RefreshSessionResult — new access + refresh pair with their expiries.
type RefreshSessionResult struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

// RefreshSession delegates to iam-svc.RefreshSession over gRPC. Replay of an
// already-revoked refresh token arrives as ErrUnauthorized (iam-svc has
// revoked every active session server-side before returning).
func (a *Adapter) RefreshSession(ctx context.Context, params *RefreshSessionParams) (*RefreshSessionResult, error) {
	const op = "iam_grpc_adapter.Adapter.RefreshSession"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "RefreshSession"))

	logger := logging.LogWithTrace(ctx, a.logger)

	// Do NOT log params.RefreshToken.
	resp, err := a.iamClient.RefreshSession(ctx, &pb.RefreshSessionRequest{
		RefreshToken: params.RefreshToken,
		UserAgent:    params.UserAgent,
		Ip:           params.IP,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &RefreshSessionResult{
		AccessToken:      resp.GetAccessToken(),
		RefreshToken:     resp.GetRefreshToken(),
		AccessExpiresAt:  time.Unix(resp.GetAccessExpiresAtUnix(), 0),
		RefreshExpiresAt: time.Unix(resp.GetRefreshExpiresAtUnix(), 0),
	}, nil
}
