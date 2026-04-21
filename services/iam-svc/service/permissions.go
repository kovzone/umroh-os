package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"iam-svc/store/postgres_store"
	"iam-svc/store/postgres_store/sqlc"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ------------------------ ValidateToken ------------------------

// ValidateTokenParams carries the raw bearer string to verify.
type ValidateTokenParams struct {
	AccessToken string
}

// ValidateTokenResult is the identity envelope a downstream consumer needs to
// authenticate + authorize a request: identity claims from the token plus the
// current role-name snapshot (fresh from the DB, not the token, so a revoked
// or re-roled user cannot keep acting on a cached claim).
type ValidateTokenResult struct {
	UserID    string
	BranchID  string
	SessionID string
	Roles     []string
	ExpiresAt time.Time
}

// ValidateToken verifies the bearer, confirms the session is still live,
// and returns the user's current role names. Any failure maps to
// ErrUnauthorized per F1 acceptance ("Gateway cannot reach iam-svc → fail closed").
func (s *Service) ValidateToken(ctx context.Context, params *ValidateTokenParams) (*ValidateTokenResult, error) {
	const op = "service.Service.ValidateToken"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	if params == nil || params.AccessToken == "" {
		e := errors.Join(apperrors.ErrUnauthorized, errors.New("empty access token"))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	payload, err := s.tokenMaker.VerifyToken(params.AccessToken)
	if err != nil {
		e := errors.Join(apperrors.ErrUnauthorized, fmt.Errorf("verify token: %w", err))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	// Reload the session row. Two invariants must hold: the row exists, it's
	// not revoked, and it hasn't expired. Anything else is fail-closed.
	sessionID, err := stringToUUID(payload.ID.String())
	if err != nil {
		e := errors.Join(apperrors.ErrUnauthorized, fmt.Errorf("parse session id: %w", err))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	session, err := s.store.GetSessionByID(ctx, sessionID)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		// Don't leak existence — any session miss is "unauthorized".
		e := errors.Join(apperrors.ErrUnauthorized, fmt.Errorf("load session: %w", wrapped))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}
	if session.RevokedAt.Valid {
		e := errors.Join(apperrors.ErrUnauthorized, errors.New("session revoked"))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}
	if !session.ExpiresAt.Valid || session.ExpiresAt.Time.Before(time.Now()) {
		e := errors.Join(apperrors.ErrUnauthorized, errors.New("session expired"))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	// Reload users.status as belt-and-suspenders (BL-IAM-003). SuspendUser
	// already flips status + revokes every session in one tx, so an active
	// session + non-active status should never coexist — but any future admin
	// path that mutates status without also revoking sessions would otherwise
	// let an in-flight access token keep working until its TTL elapses. One
	// extra DB hit per call; ValidateToken is the hot path but the query is
	// index-backed on the primary key.
	user, err := s.store.GetUserByID(ctx, session.UserID)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		e := errors.Join(apperrors.ErrUnauthorized, fmt.Errorf("load user: %w", wrapped))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}
	if user.Status != sqlc.IamUserStatusActive {
		e := errors.Join(apperrors.ErrUnauthorized, fmt.Errorf("user status=%s", user.Status))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	// Role-name snapshot from the DB. The token's Roles field is not trusted
	// for authorization decisions — role changes must propagate without waiting
	// for the token to roll over.
	roles, err := s.store.ListRoleNamesForUser(ctx, session.UserID)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Error().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	result := &ValidateTokenResult{
		UserID:    uuidToString(session.UserID),
		BranchID:  payload.BranchID,
		SessionID: payload.ID.String(),
		Roles:     roles,
		ExpiresAt: payload.ExpiredAt,
	}
	span.SetStatus(codes.Ok, "success")
	logger.Info().Str("user_id", result.UserID).Int("roles", len(roles)).Msg("token validated")
	return result, nil
}

// ------------------------ CheckPermission ------------------------

// CheckPermissionParams names the tuple to evaluate against the user's current grants.
type CheckPermissionParams struct {
	UserID   string
	Resource string
	Action   string
	Scope    string // must be one of "global" | "branch" | "personal"
}

// CheckPermissionResult carries the boolean decision. A false result is NOT
// an error at this layer — it's a valid outcome of a well-formed query. The
// gRPC server maps allowed=false to the response field; it does NOT return
// PermissionDenied. This keeps the consumer's decision logic explicit and
// easy to audit.
type CheckPermissionResult struct {
	Allowed bool
}

func (s *Service) CheckPermission(ctx context.Context, params *CheckPermissionParams) (*CheckPermissionResult, error) {
	const op = "service.Service.CheckPermission"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("resource", params.Resource),
		attribute.String("action", params.Action),
		attribute.String("scope", params.Scope),
	)

	if params == nil || params.UserID == "" || params.Resource == "" || params.Action == "" || params.Scope == "" {
		e := errors.Join(apperrors.ErrValidation, errors.New("user_id, resource, action, scope are all required"))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	scope, err := parseScope(params.Scope)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	uid, err := stringToUUID(params.UserID)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	allowed, err := s.store.UserHasPermission(ctx, sqlc.UserHasPermissionParams{
		UserID:   uid,
		Resource: params.Resource,
		Action:   params.Action,
		Scope:    scope,
	})
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Error().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetAttributes(attribute.Bool("allowed", allowed))
	span.SetStatus(codes.Ok, "success")
	logger.Info().
		Str("user_id", params.UserID).
		Str("resource", params.Resource).
		Str("action", params.Action).
		Str("scope", params.Scope).
		Bool("allowed", allowed).
		Msg("permission resolved")
	return &CheckPermissionResult{Allowed: allowed}, nil
}

// parseScope validates the inbound string against the iam.permission_scope enum.
// Unknown values return ErrValidation so the gRPC layer maps them to InvalidArgument.
func parseScope(s string) (sqlc.IamPermissionScope, error) {
	switch sqlc.IamPermissionScope(s) {
	case sqlc.IamPermissionScopeGlobal,
		sqlc.IamPermissionScopeBranch,
		sqlc.IamPermissionScopePersonal:
		return sqlc.IamPermissionScope(s), nil
	default:
		return "", errors.Join(apperrors.ErrValidation, fmt.Errorf("unknown scope %q: must be global/branch/personal", s))
	}
}
