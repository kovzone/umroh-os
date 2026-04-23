package grpc_api

import (
	"context"
	"errors"

	"iam-svc/api/grpc_api/pb"
	"iam-svc/service"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

// Permission tuple enforced before SuspendUser runs. Granted to super_admin
// by migration 000006_seed_iam_user_suspend_permission. Matches the tuple
// the legacy REST handler used to check before delegating to the service layer.
const (
	adminSuspendResource = "iam.users"
	adminSuspendAction   = "suspend"
	adminSuspendScope    = "global"
)

// SuspendUser — flip target user's status to `suspended` + revoke every active
// session in one tx. Mirrors POST /v1/users/{id}/suspend. The gateway bearer
// middleware has already validated the caller; this handler enforces the
// `iam.users/suspend/global` permission gate and self-suspend guard before
// delegating to service.SuspendUser (whose own WithTx writes the audit row).
func (s *Server) SuspendUser(ctx context.Context, req *pb.SuspendUserRequest) (*pb.SuspendUserResponse, error) {
	const op = "grpc_api.Server.SuspendUser"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "SuspendUser"),
		attribute.String("actor_user_id", req.GetActorUserId()),
		attribute.String("target_user_id", req.GetTargetUserId()),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetActorUserId() == "" || req.GetTargetUserId() == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("actor_user_id and target_user_id are required"))
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	// Permission gate — in-process service call (not a gRPC loopback). An
	// un-granted caller gets 403 before SuspendUser runs so the 404 branch
	// never leaks "target exists vs does not" to a caller without authority.
	perm, err := s.svc.CheckPermission(ctx, &service.CheckPermissionParams{
		UserID:   req.GetActorUserId(),
		Resource: adminSuspendResource,
		Action:   adminSuspendAction,
		Scope:    adminSuspendScope,
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}
	if !perm.Allowed {
		e := errors.Join(apperrors.ErrForbidden, errors.New("missing iam.users/suspend/global permission"))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, status.Error(apperrors.GRPCCode(e), apperrors.GRPCMessage(e))
	}

	result, err := s.svc.SuspendUser(ctx, &service.SuspendUserParams{
		ActorUserID:  req.GetActorUserId(),
		TargetUserID: req.GetTargetUserId(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.SuspendUserResponse{
		User: userProfileToProto(result.User),
	}, nil
}
