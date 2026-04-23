package iam_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/iam_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetMeParams — user_id from the validated bearer's payload.
type GetMeParams struct {
	UserID string
}

// GetMeResult mirrors the REST /v1/me response body 1:1.
type GetMeResult struct {
	User         UserProfile
	TOTPEnrolled bool
	TOTPVerified bool
}

// GetMe delegates to iam-svc.GetMe. NotFound surfaces when the user was
// soft-deleted after token issuance — the gateway maps that to HTTP 404.
func (a *Adapter) GetMe(ctx context.Context, params *GetMeParams) (*GetMeResult, error) {
	const op = "iam_grpc_adapter.Adapter.GetMe"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "GetMe"),
		attribute.String("user_id", params.UserID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.iamClient.GetMe(ctx, &pb.GetMeRequest{UserId: params.UserID})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &GetMeResult{
		User:         userProfileFromProto(resp.GetUser()),
		TOTPEnrolled: resp.GetTotpEnrolled(),
		TOTPVerified: resp.GetTotpVerified(),
	}, nil
}
