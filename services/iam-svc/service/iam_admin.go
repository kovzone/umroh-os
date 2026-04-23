// iam_admin.go — service implementations for IAM Phase 6 admin/security depth.
// BL-IAM-007: SetDataScope
// BL-IAM-014: CreateAPIKey, RevokeAPIKey
// BL-IAM-016: GetGlobalConfig, SetGlobalConfig
// BL-IAM-011: SearchActivityLog (keyset pagination over iam.audit_logs)

package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"iam-svc/store/postgres_store/sqlc"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Params / Results
// ---------------------------------------------------------------------------

type SetDataScopeParams struct {
	UserID    string // required
	ScopeType string // "global" | "branch" | "own_only"
	BranchID  string // required when ScopeType="branch"
}

type SetDataScopeResult struct {
	UserID    string
	ScopeType string
}

type CreateAPIKeyParams struct {
	Name      string
	Scopes    []string
	ExpiresAt string // RFC3339; "" = no expiry
	CreatedBy string
}

type CreateAPIKeyResult struct {
	KeyID        string
	PlaintextKey string // returned once; never stored
	KeyPrefix    string
	ExpiresAt    string // RFC3339 or ""
}

type RevokeAPIKeyParams struct {
	KeyID string
}

type RevokeAPIKeyResult struct {
	KeyID     string
	RevokedAt string // RFC3339
}

type GetGlobalConfigParams struct {
	Keys []string // empty = return all
}

type ConfigEntryResult struct {
	Key         string
	Value       string
	Description string
	UpdatedAt   string // RFC3339
}

type GetGlobalConfigResult struct {
	Configs []ConfigEntryResult
}

type SetGlobalConfigParams struct {
	Key         string
	Value       string
	Description string // optional — keeps existing if empty on update
	UpdatedBy   string
}

type SetGlobalConfigResult struct {
	Key       string
	Value     string
	UpdatedAt string // RFC3339
}

type SearchActivityLogParams struct {
	UserID   string
	Resource string
	Action   string
	From     string // RFC3339; "" = no lower bound
	To       string // RFC3339; "" = no upper bound
	Limit    int32  // default 50; max 200
	Cursor   string // opaque; "" = first page
}

type ActivityLogEntry struct {
	ID         string
	UserID     string
	Resource   string
	Action     string
	ResourceID string
	CreatedAt  string // RFC3339
}

type SearchActivityLogResult struct {
	Logs       []ActivityLogEntry
	NextCursor string // "" = last page
}

// ---------------------------------------------------------------------------
// Sentinel scope types
// ---------------------------------------------------------------------------

var validScopeTypes = map[string]bool{
	"global":   true,
	"branch":   true,
	"own_only": true,
}

// ---------------------------------------------------------------------------
// SetDataScope (BL-IAM-007)
// ---------------------------------------------------------------------------

func (s *Service) SetDataScope(ctx context.Context, params *SetDataScopeParams) (*SetDataScopeResult, error) {
	const op = "service.Service.SetDataScope"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("user_id", params.UserID))
	logger := logging.LogWithTrace(ctx, s.logger)

	if params.UserID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if !validScopeTypes[params.ScopeType] {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("invalid_scope_type: %q", params.ScopeType))
	}
	if params.ScopeType == "branch" && params.BranchID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("missing_branch_id"))
	}

	id := "dsc_" + strings.ReplaceAll(uuid.New().String(), "-", "")

	row, err := s.store.UpsertDataScope(ctx, sqlc.UpsertDataScopeParams{
		ID:        id,
		UserID:    params.UserID,
		ScopeType: params.ScopeType,
		BranchID:  params.BranchID,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("upsert data scope")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, errors.Join(apperrors.ErrInternal, err)
	}

	return &SetDataScopeResult{
		UserID:    row.UserID,
		ScopeType: row.ScopeType,
	}, nil
}

// ---------------------------------------------------------------------------
// CreateAPIKey (BL-IAM-014)
// ---------------------------------------------------------------------------

// generateAPIKey returns a cryptographically random API key.
// Format: "umroh_k1_<base64url(32 random bytes)>"
// Key prefix (for UI display) is the first 8 chars: "umroh_k1"
func generateAPIKey() (plaintext, keyPrefix, keyHash string, err error) {
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", "", fmt.Errorf("rand: %w", err)
	}
	encoded := base64.RawURLEncoding.EncodeToString(buf)
	plaintext = "umroh_k1_" + encoded
	keyPrefix = "umroh_k1"
	// sha256 hash — sufficient for opaque keys (no password stretching needed)
	sum := sha256.Sum256([]byte(plaintext))
	keyHash = fmt.Sprintf("%x", sum)
	return plaintext, keyPrefix, keyHash, nil
}

func (s *Service) CreateAPIKey(ctx context.Context, params *CreateAPIKeyParams) (*CreateAPIKeyResult, error) {
	const op = "service.Service.CreateAPIKey"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	if params.Name == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("missing_name"))
	}

	var expiresAt pgtype.Timestamptz
	var expiresAtStr string
	if params.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, params.ExpiresAt)
		if err != nil {
			return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("invalid_expires_at: %w", err))
		}
		if t.Before(time.Now()) {
			return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("expires_in_past"))
		}
		expiresAt = pgtype.Timestamptz{Time: t, Valid: true}
		expiresAtStr = t.UTC().Format(time.RFC3339)
	}

	plaintext, keyPrefix, keyHash, err := generateAPIKey()
	if err != nil {
		return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("generate api key: %w", err))
	}

	id := "akey_" + strings.ReplaceAll(uuid.New().String(), "-", "")
	scopes := params.Scopes
	if scopes == nil {
		scopes = []string{}
	}

	row, err := s.store.InsertAPIKey(ctx, sqlc.InsertAPIKeyParams{
		ID:        id,
		Name:      params.Name,
		KeyHash:   keyHash,
		KeyPrefix: keyPrefix,
		CreatedBy: params.CreatedBy,
		ExpiresAt: expiresAt,
		Scopes:    scopes,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("insert api key")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, errors.Join(apperrors.ErrInternal, err)
	}

	return &CreateAPIKeyResult{
		KeyID:        row.ID,
		PlaintextKey: plaintext,
		KeyPrefix:    row.KeyPrefix,
		ExpiresAt:    expiresAtStr,
	}, nil
}

// ---------------------------------------------------------------------------
// RevokeAPIKey (BL-IAM-014)
// ---------------------------------------------------------------------------

func (s *Service) RevokeAPIKey(ctx context.Context, params *RevokeAPIKeyParams) (*RevokeAPIKeyResult, error) {
	const op = "service.Service.RevokeAPIKey"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("key_id", params.KeyID))
	logger := logging.LogWithTrace(ctx, s.logger)

	if params.KeyID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("key_id is required"))
	}

	// Verify the key exists before revoking
	existing, err := s.store.GetAPIKeyByID(ctx, params.KeyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("api key %s not found", params.KeyID))
		}
		logger.Error().Err(err).Str("op", op).Msg("get api key")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, errors.Join(apperrors.ErrInternal, err)
	}

	// Already revoked — idempotent; return existing revoked_at
	if existing.RevokedAt.Valid {
		return &RevokeAPIKeyResult{
			KeyID:     existing.ID,
			RevokedAt: existing.RevokedAt.Time.UTC().Format(time.RFC3339),
		}, nil
	}

	now := time.Now().UTC()
	if err := s.store.RevokeAPIKeyByID(ctx, params.KeyID, now); err != nil {
		logger.Error().Err(err).Str("op", op).Msg("revoke api key")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, errors.Join(apperrors.ErrInternal, err)
	}

	return &RevokeAPIKeyResult{
		KeyID:     params.KeyID,
		RevokedAt: now.Format(time.RFC3339),
	}, nil
}

// ---------------------------------------------------------------------------
// GetGlobalConfig (BL-IAM-016)
// ---------------------------------------------------------------------------

func (s *Service) GetGlobalConfig(ctx context.Context, params *GetGlobalConfigParams) (*GetGlobalConfigResult, error) {
	const op = "service.Service.GetGlobalConfig"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	var rows []sqlc.IamGlobalConfig
	var err error
	if len(params.Keys) == 0 {
		rows, err = s.store.GetAllGlobalConfig(ctx)
	} else {
		rows, err = s.store.GetGlobalConfigByKeys(ctx, params.Keys)
	}
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("get global config")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, errors.Join(apperrors.ErrInternal, err)
	}

	configs := make([]ConfigEntryResult, 0, len(rows))
	for _, r := range rows {
		desc := ""
		if r.Description.Valid {
			desc = r.Description.String
		}
		configs = append(configs, ConfigEntryResult{
			Key:         r.Key,
			Value:       r.Value,
			Description: desc,
			UpdatedAt:   r.UpdatedAt.Time.UTC().Format(time.RFC3339),
		})
	}
	return &GetGlobalConfigResult{Configs: configs}, nil
}

// ---------------------------------------------------------------------------
// SetGlobalConfig (BL-IAM-016)
// ---------------------------------------------------------------------------

func (s *Service) SetGlobalConfig(ctx context.Context, params *SetGlobalConfigParams) (*SetGlobalConfigResult, error) {
	const op = "service.Service.SetGlobalConfig"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("key", params.Key))
	logger := logging.LogWithTrace(ctx, s.logger)

	if params.Key == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("missing_key"))
	}
	if params.Value == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("missing_value"))
	}

	row, err := s.store.UpsertGlobalConfig(ctx, sqlc.UpsertGlobalConfigParams{
		Key:         params.Key,
		Value:       params.Value,
		Description: params.Description,
		UpdatedBy:   params.UpdatedBy,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("upsert global config")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, errors.Join(apperrors.ErrInternal, err)
	}

	return &SetGlobalConfigResult{
		Key:       row.Key,
		Value:     row.Value,
		UpdatedAt: row.UpdatedAt.Time.UTC().Format(time.RFC3339),
	}, nil
}

// ---------------------------------------------------------------------------
// SearchActivityLog (BL-IAM-011)
// ---------------------------------------------------------------------------

// encodeCursor base64url-encodes (created_at_rfc3339 + "|" + id_uuid_text).
func encodeCursor(createdAt time.Time, idHex string) string {
	raw := createdAt.UTC().Format(time.RFC3339Nano) + "|" + idHex
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

// decodeCursor reverses encodeCursor.
func decodeCursor(cursor string) (cursorAt pgtype.Timestamptz, cursorID string, err error) {
	b, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return pgtype.Timestamptz{}, "", fmt.Errorf("base64 decode: %w", err)
	}
	parts := strings.SplitN(string(b), "|", 2)
	if len(parts) != 2 {
		return pgtype.Timestamptz{}, "", fmt.Errorf("malformed cursor")
	}
	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return pgtype.Timestamptz{}, "", fmt.Errorf("parse timestamp: %w", err)
	}
	return pgtype.Timestamptz{Time: t, Valid: true}, parts[1], nil
}

func (s *Service) SearchActivityLog(ctx context.Context, params *SearchActivityLogParams) (*SearchActivityLogResult, error) {
	const op = "service.Service.SearchActivityLog"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	lim := params.Limit
	if lim <= 0 {
		lim = 50
	}
	if lim > 200 {
		lim = 200
	}

	var fromTS, toTS pgtype.Timestamptz
	if params.From != "" {
		t, err := time.Parse(time.RFC3339, params.From)
		if err != nil {
			return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("invalid_date_format: from"))
		}
		fromTS = pgtype.Timestamptz{Time: t, Valid: true}
	}
	if params.To != "" {
		t, err := time.Parse(time.RFC3339, params.To)
		if err != nil {
			return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("invalid_date_format: to"))
		}
		toTS = pgtype.Timestamptz{Time: t, Valid: true}
	}
	if fromTS.Valid && toTS.Valid && fromTS.Time.After(toTS.Time) {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("invalid_date_range: from is after to"))
	}

	var cursorAt pgtype.Timestamptz
	var cursorID string
	if params.Cursor != "" {
		var err error
		cursorAt, cursorID, err = decodeCursor(params.Cursor)
		if err != nil {
			return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("invalid_cursor: %w", err))
		}
	}

	rows, err := s.store.SearchActivityLog(ctx, sqlc.SearchActivityLogParams{
		UserID:   params.UserID,
		Resource: params.Resource,
		Action:   params.Action,
		From:     fromTS,
		To:       toTS,
		CursorAt: cursorAt,
		CursorID: cursorID,
		Limit:    lim + 1, // fetch one extra to determine if there's a next page
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("search activity log")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, errors.Join(apperrors.ErrInternal, err)
	}

	hasNext := len(rows) > int(lim)
	if hasNext {
		rows = rows[:lim]
	}

	entries := make([]ActivityLogEntry, 0, len(rows))
	for _, r := range rows {
		entries = append(entries, ActivityLogEntry{
			ID:         uuidToString(r.ID),
			UserID:     uuidToString(r.UserID),
			Resource:   r.Resource,
			Action:     r.Action,
			ResourceID: r.ResourceID,
			CreatedAt:  r.CreatedAt.Time.UTC().Format(time.RFC3339),
		})
	}

	var nextCursor string
	if hasNext && len(rows) > 0 {
		last := rows[len(rows)-1]
		nextCursor = encodeCursor(last.CreatedAt.Time, uuidToString(last.ID))
	}

	return &SearchActivityLogResult{
		Logs:       entries,
		NextCursor: nextCursor,
	}, nil
}
