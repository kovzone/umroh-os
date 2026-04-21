package service

import (
	"context"
	"errors"
	"fmt"

	"iam-svc/store/postgres_store"
	"iam-svc/store/postgres_store/sqlc"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ------------------------ SuspendUser ------------------------

// SuspendUserParams identifies the caller (for self-suspend guard + future audit
// row in BL-IAM-004) and the target whose `iam.users.status` flips to `suspended`.
// Permission-gating happens at the REST handler via CheckPermission; this layer
// assumes the caller has already passed that check.
type SuspendUserParams struct {
	ActorUserID  string
	TargetUserID string
}

// SuspendUserResult echoes the target's post-suspend profile. Re-suspending a
// user already in `suspended` returns the same shape (idempotent) — the call
// still revokes any sessions that raced in between status flips.
type SuspendUserResult struct {
	User UserProfile
}

// SuspendUser flips the target user's status to `suspended` and revokes every
// active session for that user inside one transaction. F1-W5 acceptance is
// outcome-based ("Suspended user cannot access again"): the status column gates
// new logins (see auth.go Login + RefreshSession), and the revoke-all sweeps
// every live refresh row so no in-flight access token can be traded for a new
// one. ValidateToken additionally reloads `users.status` on every call so the
// next consumer-side CheckPermission fails closed even before the access token
// TTL elapses.
//
// The actor and target must be different users — suspending yourself would
// lock you out of the one seat that holds the suspend grant.
//
// Idempotent: re-suspending a `suspended` user is a no-op on status and still
// sweeps any sessions that appeared between the first flip and this call.
func (s *Service) SuspendUser(ctx context.Context, params *SuspendUserParams) (*SuspendUserResult, error) {
	const op = "service.Service.SuspendUser"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("actor_user_id", params.ActorUserID),
		attribute.String("target_user_id", params.TargetUserID),
	)

	if params == nil || params.ActorUserID == "" || params.TargetUserID == "" {
		e := errors.Join(apperrors.ErrValidation, errors.New("actor_user_id and target_user_id are required"))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}
	if params.ActorUserID == params.TargetUserID {
		e := errors.Join(apperrors.ErrValidation, errors.New("cannot suspend self"))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	actorUUID, err := stringToUUID(params.ActorUserID)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	_ = actorUUID // reserved for BL-IAM-004 audit write

	targetUUID, err := stringToUUID(params.TargetUserID)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Pre-tx lookup — surfaces ErrNotFound before we open a transaction, and
	// gives us the profile fields to echo in the response.
	target, err := s.store.GetUserByID(ctx, targetUUID)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	_, err = s.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(q *sqlc.Queries) error {
			if err := q.UpdateUserStatus(ctx, sqlc.UpdateUserStatusParams{
				ID:     targetUUID,
				Status: sqlc.IamUserStatusSuspended,
			}); err != nil {
				return fmt.Errorf("update user status: %w", postgres_store.WrapDBError(err))
			}
			if err := q.RevokeAllSessionsForUser(ctx, targetUUID); err != nil {
				return fmt.Errorf("revoke all sessions: %w", postgres_store.WrapDBError(err))
			}
			return nil
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	result := &SuspendUserResult{
		User: UserProfile{
			UserID:   uuidToString(target.ID),
			Email:    target.Email,
			Name:     target.Name,
			BranchID: uuidToString(target.BranchID),
			Status:   string(sqlc.IamUserStatusSuspended),
		},
	}

	span.SetStatus(codes.Ok, "success")
	logger.Info().
		Str("actor_user_id", params.ActorUserID).
		Str("target_user_id", params.TargetUserID).
		Msg("user suspended; all sessions revoked")
	return result, nil
}
