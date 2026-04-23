package iam_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/iam_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// EnrollTOTPParams — user_id from the validated bearer's payload.
type EnrollTOTPParams struct {
	UserID string
}

// EnrollTOTPResult carries the plaintext secret + otpauth URL exactly once.
// The gateway handler must not log or persist these anywhere downstream.
type EnrollTOTPResult struct {
	Secret     string
	OtpauthURL string
}

// EnrollTOTP delegates to iam-svc.EnrollTotp. Already-verified users get
// ErrConflict → HTTP 409. Admin-assisted reset is S1-E-06.
func (a *Adapter) EnrollTOTP(ctx context.Context, params *EnrollTOTPParams) (*EnrollTOTPResult, error) {
	const op = "iam_grpc_adapter.Adapter.EnrollTOTP"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "EnrollTotp"),
		attribute.String("user_id", params.UserID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.iamClient.EnrollTotp(ctx, &pb.EnrollTotpRequest{UserId: params.UserID})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &EnrollTOTPResult{
		Secret:     resp.GetSecret(),
		OtpauthURL: resp.GetOtpauthUrl(),
	}, nil
}
