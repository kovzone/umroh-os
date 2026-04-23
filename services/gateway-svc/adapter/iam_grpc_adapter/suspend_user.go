package iam_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/iam_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// SuspendUserParams — actor_user_id from the validated bearer payload,
// target_user_id from the REST path param.
type SuspendUserParams struct {
	ActorUserID  string
	TargetUserID string
}

// SuspendUserResult echoes the post-suspend profile for the REST response.
type SuspendUserResult struct {
	User UserProfile
}

// SuspendUser delegates to iam-svc.SuspendUser. iam-svc enforces the
// iam.users/suspend/global permission gate + self-suspend guard before
// flipping status; errors flow through mapIamError.
func (a *Adapter) SuspendUser(ctx context.Context, params *SuspendUserParams) (*SuspendUserResult, error) {
	const op = "iam_grpc_adapter.Adapter.SuspendUser"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "SuspendUser"),
		attribute.String("actor_user_id", params.ActorUserID),
		attribute.String("target_user_id", params.TargetUserID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.iamClient.SuspendUser(ctx, &pb.SuspendUserRequest{
		ActorUserId:  params.ActorUserID,
		TargetUserId: params.TargetUserID,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &SuspendUserResult{
		User: userProfileFromProto(resp.GetUser()),
	}, nil
}
