package iam_grpc_adapter

import (
	"context"
	"time"

	"gateway-svc/adapter/iam_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// VerifyTOTPParams — six-digit code + user_id from the bearer payload.
type VerifyTOTPParams struct {
	UserID string
	Code   string
}

// VerifyTOTPResult — stamped verified_at on success.
type VerifyTOTPResult struct {
	VerifiedAt time.Time
}

// VerifyTOTP delegates to iam-svc.VerifyTotp. Invalid code → ErrUnauthorized
// (401). Not-yet-enrolled → ErrValidation (400 / 422 at the REST boundary).
func (a *Adapter) VerifyTOTP(ctx context.Context, params *VerifyTOTPParams) (*VerifyTOTPResult, error) {
	const op = "iam_grpc_adapter.Adapter.VerifyTOTP"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "VerifyTotp"),
		attribute.String("user_id", params.UserID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.iamClient.VerifyTotp(ctx, &pb.VerifyTotpRequest{
		UserId: params.UserID,
		Code:   params.Code,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &VerifyTOTPResult{
		VerifiedAt: time.Unix(resp.GetVerifiedAtUnix(), 0),
	}, nil
}
