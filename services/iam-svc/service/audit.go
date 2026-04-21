package service

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"iam-svc/store/postgres_store"
	"iam-svc/store/postgres_store/sqlc"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// RecordAuditParams mirrors iam.audit_logs 1:1. Optional string fields are empty
// when absent — the server maps them to NULL at the storage layer. OldValue and
// NewValue are opaque JSONB byte slices (the caller marshals their own JSON).
type RecordAuditParams struct {
	ActorUserID string // UUID string, optional ("" → NULL in storage)
	BranchID    string // UUID string, optional ("" → NULL)
	Resource    string // required
	ResourceID  string // optional; column default is ''
	Action      string // required
	OldValue    []byte // JSONB bytes, optional (len 0 → NULL)
	NewValue    []byte // JSONB bytes, optional (len 0 → NULL)
	IP          string // IP literal, optional ("" → NULL; non-empty must parse)
}

// RecordAuditResult echoes the inserted row's id + timestamp so callers can
// cite the audit entry without issuing a second query. No PII is returned.
type RecordAuditResult struct {
	AuditLogID string
	CreatedAt  time.Time
}

// RecordAudit inserts one row into iam.audit_logs. The table is append-only at
// the DB layer (UPDATE/DELETE raise insufficient_privilege via trigger), so the
// single statement is durable the moment the tx commits.
//
// Use this for cross-service audit emission over gRPC. For iam-svc's own
// state-changing actions that already own a `WithTx` (e.g. SuspendUser), emit
// from inside that closure directly against `q.InsertAuditLog` — the business
// write and its audit row land atomically without a round-trip.
//
// Validation is strict: resource + action are required, any non-empty UUID /
// IP string must parse, otherwise ErrValidation. The server never silently
// drops or coerces input.
func (s *Service) RecordAudit(ctx context.Context, params *RecordAuditParams) (*RecordAuditResult, error) {
	const op = "service.Service.RecordAudit"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("resource", safeOrEmpty(params, func(p *RecordAuditParams) string { return p.Resource })),
		attribute.String("action", safeOrEmpty(params, func(p *RecordAuditParams) string { return p.Action })),
	)

	if params == nil || params.Resource == "" || params.Action == "" {
		e := errors.Join(apperrors.ErrValidation, errors.New("resource and action are required"))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	userUUID, err := optionalStringToUUID(params.ActorUserID)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	branchUUID, err := optionalStringToUUID(params.BranchID)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	ipPtr, err := optionalStringToIP(params.IP)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	row, err := s.store.InsertAuditLog(ctx, sqlc.InsertAuditLogParams{
		UserID:     userUUID,
		BranchID:   branchUUID,
		Resource:   params.Resource,
		ResourceID: params.ResourceID,
		Action:     params.Action,
		OldValue:   jsonbOrNil(params.OldValue),
		NewValue:   jsonbOrNil(params.NewValue),
		Ip:         ipPtr,
	})
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Error().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	result := &RecordAuditResult{
		AuditLogID: uuidToString(row.ID),
		CreatedAt:  row.CreatedAt.Time,
	}

	span.SetAttributes(attribute.String("audit_log_id", result.AuditLogID))
	span.SetStatus(codes.Ok, "success")
	logger.Info().
		Str("audit_log_id", result.AuditLogID).
		Str("actor_user_id", params.ActorUserID).
		Str("resource", params.Resource).
		Str("resource_id", params.ResourceID).
		Str("action", params.Action).
		Msg("audit recorded")
	return result, nil
}

// optionalStringToUUID returns the zero pgtype.UUID (Valid=false → SQL NULL) for
// empty input, a parsed UUID otherwise. Non-empty-but-invalid is ErrValidation.
func optionalStringToUUID(s string) (pgtype.UUID, error) {
	if s == "" {
		return pgtype.UUID{}, nil
	}
	return stringToUUID(s)
}

// optionalStringToIP parses an optional IP literal. Empty input → nil pointer
// (stored as NULL); non-empty-but-invalid → ErrValidation.
func optionalStringToIP(s string) (*netip.Addr, error) {
	if s == "" {
		return nil, nil
	}
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("parse ip %q: %w", s, err))
	}
	return &addr, nil
}

// jsonbOrNil passes through a non-empty byte slice and converts empty to nil
// so pgx writes a JSON NULL instead of the literal string "".
func jsonbOrNil(b []byte) []byte {
	if len(b) == 0 {
		return nil
	}
	return b
}

// safeOrEmpty is a small nil-guard for span attributes — the early validation
// branch needs the attribute set before we reject a nil params.
func safeOrEmpty(p *RecordAuditParams, f func(*RecordAuditParams) string) string {
	if p == nil {
		return ""
	}
	return f(p)
}
