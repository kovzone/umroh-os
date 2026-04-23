package iam_grpc_adapter

import (
	"context"
	"time"

	"catalog-svc/adapter/iam_grpc_adapter/pb"
	"catalog-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// RecordAuditParams carries the audit log entry to write.
type RecordAuditParams struct {
	UserID     string
	BranchID   string
	Resource   string
	ResourceID string
	Action     string
	OldValue   []byte
	NewValue   []byte
	IP         string
}

// RecordAuditResult is the minimal response from iam-svc.RecordAudit.
type RecordAuditResult struct {
	AuditLogID  string
	CreatedAtUnix time.Time
}

// RecordAudit delegates to iam-svc.RecordAudit over gRPC.
//
// Transport errors are logged as WARN but NOT returned to the caller — audit
// writes are best-effort in the MVP. Callers SHOULD still check the error for
// observability but MAY proceed on failure per the "audit SHOULD" clause in
// § Catalog — internal write.
func (a *Adapter) RecordAudit(ctx context.Context, params *RecordAuditParams) (*RecordAuditResult, error) {
	const op = "iam_grpc_adapter.Adapter.RecordAudit"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "RecordAudit"),
		attribute.String("resource", params.Resource),
		attribute.String("resource_id", params.ResourceID),
		attribute.String("action", params.Action),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.iamClient.RecordAudit(ctx, &pb.RecordAuditRequest{
		UserId:     params.UserID,
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
		logger.Warn().Err(wrapped).Str("resource", params.Resource).Str("resource_id", params.ResourceID).Msg("RecordAudit failed (best-effort)")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		// Return the error so callers can observe — they decide whether to abort.
		return nil, wrapped
	}

	span.SetAttributes(attribute.String("audit_log_id", resp.GetAuditLogId()))
	span.SetStatus(codes.Ok, "success")
	return &RecordAuditResult{
		AuditLogID:    resp.GetAuditLogId(),
		CreatedAtUnix: time.Unix(resp.GetCreatedAtUnix(), 0),
	}, nil
}
