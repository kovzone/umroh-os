package iam_grpc_adapter

import (
	"context"

	"finance-svc/adapter/iam_grpc_adapter/pb"
	"finance-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// CheckPermissionParams carries the tuple to evaluate.
type CheckPermissionParams struct {
	UserID   string
	Resource string
	Action   string
	Scope    string // "global" | "branch" | "personal"
}

// CheckPermissionResult carries the decision as a bare boolean. allowed=false
// is NOT an error here — the caller translates it to HTTP 403.
type CheckPermissionResult struct {
	Allowed bool
}

// CheckPermission delegates to iam-svc.CheckPermission over gRPC.
//
// Transport errors are mapped through mapIamError (same policy as ValidateToken):
// an unreachable iam-svc becomes ErrUnauthorized so the consumer fails closed.
// A successful RPC with allowed=false is a valid "deny" outcome — no error.
func (a *Adapter) CheckPermission(ctx context.Context, params *CheckPermissionParams) (*CheckPermissionResult, error) {
	const op = "iam_grpc_adapter.Adapter.CheckPermission"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "CheckPermission"),
		attribute.String("resource", params.Resource),
		attribute.String("action", params.Action),
		attribute.String("scope", params.Scope),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.iamClient.CheckPermission(ctx, &pb.CheckPermissionRequest{
		UserId:   params.UserID,
		Resource: params.Resource,
		Action:   params.Action,
		Scope:    params.Scope,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetAttributes(attribute.Bool("allowed", resp.GetAllowed()))
	span.SetStatus(codes.Ok, "success")
	return &CheckPermissionResult{Allowed: resp.GetAllowed()}, nil
}
