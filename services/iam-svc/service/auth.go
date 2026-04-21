package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"iam-svc/store/postgres_store"
	"iam-svc/store/postgres_store/sqlc"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"
	"iam-svc/util/token"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/crypto/bcrypt"
)

// refreshTokenBytes is the byte length of the cryptographically random
// refresh-token secret the caller receives. 32 bytes = 256 bits.
const refreshTokenBytes = 32

// UserProfile is the shared user envelope returned from Login / RefreshSession / GetMe.
type UserProfile struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	BranchID string `json:"branch_id,omitempty"`
	Status   string `json:"status"`
}

// ------------------------ Login ------------------------

type LoginParams struct {
	Email     string
	Password  string
	TOTPCode  string // optional; login-time TOTP enforcement deferred to S1-E-06
	UserAgent string
	IP        *netip.Addr
}

type LoginResult struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	User             UserProfile
}

func (s *Service) Login(ctx context.Context, params *LoginParams) (*LoginResult, error) {
	const op = "service.Service.Login"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("email", params.Email),
	)
	logger.Info().Str("op", op).Str("email", params.Email).Msg("")

	// 1. Lookup user. Map NOT_FOUND to UNAUTHORIZED — don't leak existence.
	user, err := s.store.GetUserByEmail(ctx, params.Email)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		if errors.Is(wrapped, apperrors.ErrNotFound) {
			wrapped = errors.Join(apperrors.ErrUnauthorized, errors.New("invalid credentials"))
		}
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	// 2. bcrypt compare.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password)); err != nil {
		e := errors.Join(apperrors.ErrUnauthorized, fmt.Errorf("password mismatch for %s", params.Email))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	// 3. Status gate.
	if user.Status != sqlc.IamUserStatusActive {
		e := errors.Join(apperrors.ErrForbidden, fmt.Errorf("user status=%s", user.Status))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	// 4. Mint refresh token + persist session + stamp last_login (atomic).
	refreshPlain, refreshHash, err := generateRefreshToken()
	if err != nil {
		e := errors.Join(apperrors.ErrInternal, fmt.Errorf("generate refresh: %w", err))
		logger.Error().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	refreshExpiresAt := time.Now().Add(s.refreshTokenTTL)
	var sessionRow sqlc.IamSession

	_, err = s.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(q *sqlc.Queries) error {
			created, err := q.CreateSession(ctx, sqlc.CreateSessionParams{
				UserID:           user.ID,
				RefreshTokenHash: refreshHash,
				UserAgent:        params.UserAgent,
				Ip:               params.IP,
				ExpiresAt:        pgtype.Timestamptz{Time: refreshExpiresAt, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("create session: %w", postgres_store.WrapDBError(err))
			}
			sessionRow = created
			if err := q.UpdateUserLastLoginAt(ctx, user.ID); err != nil {
				return fmt.Errorf("update last_login_at: %w", postgres_store.WrapDBError(err))
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

	// 5. Mint access token. Session UUID (not a random token id) is the payload.ID — logout uses it.
	userID := uuidToString(user.ID)
	branchID := uuidToString(user.BranchID)
	payload := &token.Payload{
		ID:       uuid.UUID(sessionRow.ID.Bytes),
		UserID:   userID,
		BranchID: branchID,
		Roles:    []string{}, // populated by BL-IAM-002 via user_roles × roles join
	}
	signed, err := s.tokenMaker.CreateToken(payload, s.accessTokenTTL)
	if err != nil {
		e := errors.Join(apperrors.ErrInternal, fmt.Errorf("sign access token: %w", err))
		logger.Error().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	result := &LoginResult{
		AccessToken:      signed,
		RefreshToken:     refreshPlain,
		AccessExpiresAt:  payload.ExpiredAt,
		RefreshExpiresAt: refreshExpiresAt,
		User: UserProfile{
			UserID:   userID,
			Email:    user.Email,
			Name:     user.Name,
			BranchID: branchID,
			Status:   string(user.Status),
		},
	}
	span.SetStatus(codes.Ok, "success")
	logger.Info().Str("user_id", userID).Msg("login success")
	return result, nil
}

// ------------------------ RefreshSession ------------------------

type RefreshSessionParams struct {
	RefreshToken string
	UserAgent    string
	IP           *netip.Addr
}

type RefreshSessionResult struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

func (s *Service) RefreshSession(ctx context.Context, params *RefreshSessionParams) (*RefreshSessionResult, error) {
	const op = "service.Service.RefreshSession"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	hashed := hashRefreshToken(params.RefreshToken)

	existing, err := s.store.GetSessionByRefreshHash(ctx, hashed)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		if errors.Is(wrapped, apperrors.ErrNotFound) {
			wrapped = errors.Join(apperrors.ErrUnauthorized, errors.New("unknown refresh token"))
		}
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	// Replay detection: if someone presents a token matching a revoked row,
	// assume the session has been compromised and wipe every active session for the user.
	if existing.RevokedAt.Valid {
		if err := s.store.RevokeAllSessionsForUser(ctx, existing.UserID); err != nil {
			logger.Error().Err(postgres_store.WrapDBError(err)).Msg("revoke-all after replay failed")
		}
		e := errors.Join(apperrors.ErrUnauthorized, errors.New("refresh-token replay; all sessions revoked"))
		logger.Error().Err(e).Msg("refresh replay detected")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}
	if !existing.ExpiresAt.Valid || existing.ExpiresAt.Time.Before(time.Now()) {
		e := errors.Join(apperrors.ErrUnauthorized, errors.New("refresh-token expired"))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	user, err := s.store.GetUserByID(ctx, existing.UserID)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Error().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	if user.Status != sqlc.IamUserStatusActive {
		e := errors.Join(apperrors.ErrForbidden, fmt.Errorf("user status=%s", user.Status))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	newPlain, newHash, err := generateRefreshToken()
	if err != nil {
		e := errors.Join(apperrors.ErrInternal, fmt.Errorf("generate refresh: %w", err))
		logger.Error().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	newExpires := time.Now().Add(s.refreshTokenTTL)
	var newSession sqlc.IamSession
	_, err = s.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(q *sqlc.Queries) error {
			if err := q.RevokeSession(ctx, existing.ID); err != nil {
				return fmt.Errorf("revoke old session: %w", postgres_store.WrapDBError(err))
			}
			created, err := q.CreateSession(ctx, sqlc.CreateSessionParams{
				UserID:           existing.UserID,
				RefreshTokenHash: newHash,
				UserAgent:        params.UserAgent,
				Ip:               params.IP,
				ExpiresAt:        pgtype.Timestamptz{Time: newExpires, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("create session: %w", postgres_store.WrapDBError(err))
			}
			newSession = created
			return nil
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	userID := uuidToString(user.ID)
	payload := &token.Payload{
		ID:       uuid.UUID(newSession.ID.Bytes),
		UserID:   userID,
		BranchID: uuidToString(user.BranchID),
		Roles:    []string{},
	}
	signed, err := s.tokenMaker.CreateToken(payload, s.accessTokenTTL)
	if err != nil {
		e := errors.Join(apperrors.ErrInternal, fmt.Errorf("sign access token: %w", err))
		logger.Error().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	result := &RefreshSessionResult{
		AccessToken:      signed,
		RefreshToken:     newPlain,
		AccessExpiresAt:  payload.ExpiredAt,
		RefreshExpiresAt: newExpires,
	}
	span.SetStatus(codes.Ok, "success")
	logger.Info().Str("user_id", userID).Msg("refresh success")
	return result, nil
}

// ------------------------ Logout ------------------------

type LogoutParams struct {
	// SessionID is the session UUID stored in the bearer-token payload's ID claim
	// (set at login / refresh). The handler reads it out of the verified payload.
	SessionID string
}

type LogoutResult struct{}

func (s *Service) Logout(ctx context.Context, params *LogoutParams) (*LogoutResult, error) {
	const op = "service.Service.Logout"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("session_id", params.SessionID),
	)

	sid, err := stringToUUID(params.SessionID)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if err := s.store.RevokeSession(ctx, sid); err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Error().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	logger.Info().Str("session_id", params.SessionID).Msg("logout success")
	return &LogoutResult{}, nil
}

// ------------------------ helpers ------------------------

// generateRefreshToken returns a plaintext hex-encoded random token and its SHA-256 hex digest.
// The plaintext goes to the client; only the digest is ever stored.
func generateRefreshToken() (plain string, hash string, err error) {
	buf := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	plain = hex.EncodeToString(buf)
	hash = hashRefreshToken(plain)
	return plain, hash, nil
}

// hashRefreshToken returns the canonical SHA-256 hex digest used to compare a caller-supplied
// refresh token against the row in iam.sessions.
func hashRefreshToken(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}

// uuidToString returns the canonical "xxxxxxxx-xxxx-..." form of a pgtype.UUID, or "" if invalid.
func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return uuid.UUID(u.Bytes).String()
}

// stringToUUID parses a UUID string into pgtype.UUID. Bad input is ErrValidation.
func stringToUUID(s string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, errors.Join(apperrors.ErrValidation, fmt.Errorf("parse uuid %q: %w", s, err))
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}
