package iam_grpc_adapter

import (
	"context"
	"time"

	"gateway-svc/adapter/iam_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// RecordAuditParams mirrors iam-svc's RecordAuditRequest. Optional string
// fields use empty string for "absent" — the server maps them to NULL in
// iam.audit_logs. Callers marshal their own JSON for OldValue / NewValue.
type RecordAuditParams struct {
	ActorUserID string // UUID string; "" → NULL
	BranchID    string // UUID string; "" → NULL
	Resource    string // required (e.g. "booking")
	ResourceID  string // optional; "" permitted
	Action      string // required (e.g. "create", "update")
	OldValue    []byte // JSONB bytes; len 0 → NULL
	NewValue    []byte // JSONB bytes; len 0 → NULL
	IP          string // IP literal; "" → NULL
}

// RecordAuditResult echoes the inserted row's id + timestamp so callers can
// cite the audit entry without re-querying.
type RecordAuditResult struct {
	AuditLogID string
	CreatedAt  time.Time
}

// RecordAudit delegates to iam-svc.RecordAudit over gRPC.
//
// The RPC is synchronous: the wire call returns only after the row is
// durably inserted into iam.audit_logs. Callers that prefer best-effort
// audit emit (e.g. not blocking a business action on the audit write)
// wrap this call in a goroutine on their side — the adapter intentionally
// does not take that decision on their behalf.
//
// Error mapping uses the shared mapIamError helper: validation failures
// land as ErrValidation, transport / unclassified failures as ErrUnauthorized
// (fail-closed, matching ValidateToken / CheckPermission policy).
//
// Not currently called by any gateway handler. Scaffolded for future staff
// routes (e.g. mutating catalog writes) where gateway-side audit emission
// may land before the backend returns.
func (a *Adapter) RecordAudit(ctx context.Context, params *RecordAuditParams) (*RecordAuditResult, error) {
	const op = "iam_grpc_adapter.Adapter.RecordAudit"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "RecordAudit"),
		attribute.String("resource", params.Resource),
		attribute.String("action", params.Action),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.iamClient.RecordAudit(ctx, &pb.RecordAuditRequest{
		UserId:     params.ActorUserID,
		BranchId:   params.BranchID,
		Resource:   params.Resource,
		ResourceId: params.ResourceID,
		Action:     params.Action,
		OldValue:   params.OldValue,
		NewValue:   params.NewValue,
		Ip:         params.IP,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetAttributes(attribute.String("audit_log_id", resp.GetAuditLogId()))
	span.SetStatus(codes.Ok, "success")
	return &RecordAuditResult{
		AuditLogID: resp.GetAuditLogId(),
		CreatedAt:  time.Unix(resp.GetCreatedAtUnix(), 0),
	}, nil
}
