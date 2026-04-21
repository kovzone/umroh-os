package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"iam-svc/store/postgres_store"
	"iam-svc/store/postgres_store/sqlc"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pquerna/otp/totp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ------------------------ GetMe ------------------------

type GetMeParams struct {
	UserID string
}

type GetMeResult struct {
	User         UserProfile
	TOTPEnrolled bool
	TOTPVerified bool
}

func (s *Service) GetMe(ctx context.Context, params *GetMeParams) (*GetMeResult, error) {
	const op = "service.Service.GetMe"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("user_id", params.UserID),
	)

	uid, err := stringToUUID(params.UserID)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	user, err := s.store.GetUserByID(ctx, uid)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Error().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	result := &GetMeResult{
		User: UserProfile{
			UserID:   uuidToString(user.ID),
			Email:    user.Email,
			Name:     user.Name,
			BranchID: uuidToString(user.BranchID),
			Status:   string(user.Status),
		},
		TOTPEnrolled: user.TotpSecret.Valid,
		TOTPVerified: user.TotpVerifiedAt.Valid,
	}
	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// ------------------------ EnrollTOTP ------------------------

type EnrollTOTPParams struct {
	UserID      string
	AccountName string // label in authenticator app; falls back to user email
}

type EnrollTOTPResult struct {
	Secret     string // plaintext base32 — shown to the user exactly once, never logged
	OtpauthURL string // otpauth:// URL for QR rendering
}

func (s *Service) EnrollTOTP(ctx context.Context, params *EnrollTOTPParams) (*EnrollTOTPResult, error) {
	const op = "service.Service.EnrollTOTP"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("user_id", params.UserID),
	)

	uid, err := stringToUUID(params.UserID)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	user, err := s.store.GetUserByID(ctx, uid)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Error().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	// Already-verified users must not re-enroll via this endpoint. Device-lost
	// recovery is an admin-triggered reset (S1-E-06 BL-IAM-005..017 depth card).
	if user.TotpVerifiedAt.Valid {
		e := errors.Join(apperrors.ErrConflict, fmt.Errorf("totp already verified for user %s", params.UserID))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	accountName := params.AccountName
	if accountName == "" {
		accountName = user.Email
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.totpIssuer,
		AccountName: accountName,
	})
	if err != nil {
		e := errors.Join(apperrors.ErrInternal, fmt.Errorf("totp generate: %w", err))
		logger.Error().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	encrypted, err := encryptAESGCM(s.totpEncryptionKey, []byte(key.Secret()))
	if err != nil {
		e := errors.Join(apperrors.ErrInternal, fmt.Errorf("encrypt secret: %w", err))
		logger.Error().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	if err := s.store.UpdateUserTOTP(ctx, sqlc.UpdateUserTOTPParams{
		ID:             uid,
		TotpSecret:     pgtype.Text{String: encrypted, Valid: true},
		TotpVerifiedAt: pgtype.Timestamptz{Valid: false}, // unverified until VerifyTOTP
	}); err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Error().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	result := &EnrollTOTPResult{
		Secret:     key.Secret(),
		OtpauthURL: key.URL(),
	}
	span.SetStatus(codes.Ok, "success")
	logger.Info().Str("user_id", params.UserID).Msg("totp enrolled")
	return result, nil
}

// ------------------------ VerifyTOTP ------------------------

type VerifyTOTPParams struct {
	UserID string
	Code   string
}

type VerifyTOTPResult struct {
	VerifiedAt time.Time
}

func (s *Service) VerifyTOTP(ctx context.Context, params *VerifyTOTPParams) (*VerifyTOTPResult, error) {
	const op = "service.Service.VerifyTOTP"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("user_id", params.UserID),
	)

	uid, err := stringToUUID(params.UserID)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	user, err := s.store.GetUserByID(ctx, uid)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Error().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	if !user.TotpSecret.Valid {
		e := errors.Join(apperrors.ErrValidation, errors.New("totp not enrolled"))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	secret, err := decryptAESGCM(s.totpEncryptionKey, user.TotpSecret.String)
	if err != nil {
		e := errors.Join(apperrors.ErrInternal, fmt.Errorf("decrypt secret: %w", err))
		logger.Error().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	if !totp.Validate(params.Code, string(secret)) {
		e := errors.Join(apperrors.ErrUnauthorized, errors.New("invalid totp code"))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	now := time.Now()
	if err := s.store.UpdateUserTOTP(ctx, sqlc.UpdateUserTOTPParams{
		ID:             uid,
		TotpSecret:     user.TotpSecret, // ciphertext unchanged
		TotpVerifiedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}); err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Error().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	logger.Info().Str("user_id", params.UserID).Msg("totp verified")
	return &VerifyTOTPResult{VerifiedAt: now}, nil
}

// ------------------------ AES-256-GCM helpers ------------------------

// Ciphertext format: "<hex-nonce>:<hex-ciphertext-including-gcm-tag>".
// Key is required to be exactly 32 bytes (AES-256) — enforced at start-up by config validation,
// not here, to keep hot-path helpers allocation-lean.

func encryptAESGCM(key, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ct := gcm.Seal(nil, nonce, plaintext, nil)
	return hex.EncodeToString(nonce) + ":" + hex.EncodeToString(ct), nil
}

func decryptAESGCM(key []byte, payload string) ([]byte, error) {
	parts := strings.SplitN(payload, ":", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid ciphertext format")
	}
	nonce, err := hex.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	ct, err := hex.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// gcm.Open panics if the nonce length does not match its expected size; guard so
	// malformed ciphertext returns an error rather than crashing the handler.
	if len(nonce) != gcm.NonceSize() {
		return nil, errors.New("invalid ciphertext: nonce length mismatch")
	}
	return gcm.Open(nil, nonce, ct, nil)
}
