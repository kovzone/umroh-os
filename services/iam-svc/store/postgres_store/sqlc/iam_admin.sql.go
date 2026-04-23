// iam_admin.sql.go — hand-written SQLC-style queries for IAM Phase 6 admin
// tables: data_scopes, api_keys, global_config, and audit_log search
// (BL-IAM-007/010/011/014/016).

package sqlc

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// Models
// ---------------------------------------------------------------------------

type IamDataScope struct {
	ID        string             `json:"id"`
	UserID    string             `json:"user_id"`
	ScopeType string             `json:"scope_type"`
	BranchID  pgtype.Text        `json:"branch_id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

type IamAPIKey struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	KeyHash    string             `json:"key_hash"`
	KeyPrefix  string             `json:"key_prefix"`
	CreatedBy  string             `json:"created_by"`
	ExpiresAt  pgtype.Timestamptz `json:"expires_at"`
	LastUsedAt pgtype.Timestamptz `json:"last_used_at"`
	RevokedAt  pgtype.Timestamptz `json:"revoked_at"`
	Scopes     []string           `json:"scopes"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
}

type IamGlobalConfig struct {
	Key         string             `json:"key"`
	Value       string             `json:"value"`
	Description pgtype.Text        `json:"description"`
	UpdatedBy   pgtype.Text        `json:"updated_by"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

// ---------------------------------------------------------------------------
// data_scopes — UpsertDataScope
// ---------------------------------------------------------------------------

const upsertDataScope = `-- name: UpsertDataScope :one
INSERT INTO iam.data_scopes (id, user_id, scope_type, branch_id)
VALUES ($1, $2, $3, NULLIF($4, ''))
ON CONFLICT (user_id) DO UPDATE
  SET scope_type = EXCLUDED.scope_type,
      branch_id  = EXCLUDED.branch_id
RETURNING id, user_id, scope_type, branch_id, created_at
`

type UpsertDataScopeParams struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ScopeType string `json:"scope_type"`
	BranchID  string `json:"branch_id"` // "" → NULL
}

func (q *Queries) UpsertDataScope(ctx context.Context, arg UpsertDataScopeParams) (IamDataScope, error) {
	row := q.db.QueryRow(ctx, upsertDataScope, arg.ID, arg.UserID, arg.ScopeType, arg.BranchID)
	var i IamDataScope
	err := row.Scan(&i.ID, &i.UserID, &i.ScopeType, &i.BranchID, &i.CreatedAt)
	return i, err
}

// ---------------------------------------------------------------------------
// api_keys
// ---------------------------------------------------------------------------

const insertAPIKey = `-- name: InsertAPIKey :one
INSERT INTO iam.api_keys (id, name, key_hash, key_prefix, created_by, expires_at, scopes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, key_hash, key_prefix, created_by,
          expires_at, last_used_at, revoked_at, scopes, created_at
`

type InsertAPIKeyParams struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	KeyHash   string             `json:"key_hash"`
	KeyPrefix string             `json:"key_prefix"`
	CreatedBy string             `json:"created_by"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"` // Valid=false → NULL (no expiry)
	Scopes    []string           `json:"scopes"`
}

func (q *Queries) InsertAPIKey(ctx context.Context, arg InsertAPIKeyParams) (IamAPIKey, error) {
	row := q.db.QueryRow(ctx, insertAPIKey,
		arg.ID, arg.Name, arg.KeyHash, arg.KeyPrefix,
		arg.CreatedBy, arg.ExpiresAt, arg.Scopes,
	)
	var i IamAPIKey
	err := row.Scan(
		&i.ID, &i.Name, &i.KeyHash, &i.KeyPrefix, &i.CreatedBy,
		&i.ExpiresAt, &i.LastUsedAt, &i.RevokedAt, &i.Scopes, &i.CreatedAt,
	)
	return i, err
}

const getAPIKeyByID = `-- name: GetAPIKeyByID :one
SELECT id, name, key_hash, key_prefix, created_by,
       expires_at, last_used_at, revoked_at, scopes, created_at
FROM iam.api_keys
WHERE id = $1
`

func (q *Queries) GetAPIKeyByID(ctx context.Context, id string) (IamAPIKey, error) {
	row := q.db.QueryRow(ctx, getAPIKeyByID, id)
	var i IamAPIKey
	err := row.Scan(
		&i.ID, &i.Name, &i.KeyHash, &i.KeyPrefix, &i.CreatedBy,
		&i.ExpiresAt, &i.LastUsedAt, &i.RevokedAt, &i.Scopes, &i.CreatedAt,
	)
	return i, err
}

const revokeAPIKeyByID = `-- name: RevokeAPIKeyByID :exec
UPDATE iam.api_keys
SET revoked_at = COALESCE(revoked_at, $2)
WHERE id = $1
`

// RevokeAPIKeyByID sets revoked_at to now if not already set (idempotent).
func (q *Queries) RevokeAPIKeyByID(ctx context.Context, id string, revokedAt time.Time) error {
	_, err := q.db.Exec(ctx, revokeAPIKeyByID, id, revokedAt)
	return err
}

// ---------------------------------------------------------------------------
// global_config
// ---------------------------------------------------------------------------

const getAllGlobalConfig = `-- name: GetAllGlobalConfig :many
SELECT key, value, description, updated_by, updated_at
FROM iam.global_config
ORDER BY key
`

func (q *Queries) GetAllGlobalConfig(ctx context.Context) ([]IamGlobalConfig, error) {
	rows, err := q.db.Query(ctx, getAllGlobalConfig)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []IamGlobalConfig
	for rows.Next() {
		var i IamGlobalConfig
		if err := rows.Scan(&i.Key, &i.Value, &i.Description, &i.UpdatedBy, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const getGlobalConfigByKeys = `-- name: GetGlobalConfigByKeys :many
SELECT key, value, description, updated_by, updated_at
FROM iam.global_config
WHERE key = ANY($1::text[])
ORDER BY key
`

func (q *Queries) GetGlobalConfigByKeys(ctx context.Context, keys []string) ([]IamGlobalConfig, error) {
	rows, err := q.db.Query(ctx, getGlobalConfigByKeys, keys)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []IamGlobalConfig
	for rows.Next() {
		var i IamGlobalConfig
		if err := rows.Scan(&i.Key, &i.Value, &i.Description, &i.UpdatedBy, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const upsertGlobalConfig = `-- name: UpsertGlobalConfig :one
INSERT INTO iam.global_config (key, value, description, updated_by, updated_at)
VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''), NOW())
ON CONFLICT (key) DO UPDATE
  SET value       = EXCLUDED.value,
      description = CASE
                      WHEN $3 = '' THEN iam.global_config.description
                      ELSE EXCLUDED.description
                    END,
      updated_by  = EXCLUDED.updated_by,
      updated_at  = NOW()
RETURNING key, value, description, updated_by, updated_at
`

type UpsertGlobalConfigParams struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"` // "" = keep existing on update
	UpdatedBy   string `json:"updated_by"`
}

func (q *Queries) UpsertGlobalConfig(ctx context.Context, arg UpsertGlobalConfigParams) (IamGlobalConfig, error) {
	row := q.db.QueryRow(ctx, upsertGlobalConfig,
		arg.Key, arg.Value, arg.Description, arg.UpdatedBy)
	var i IamGlobalConfig
	err := row.Scan(&i.Key, &i.Value, &i.Description, &i.UpdatedBy, &i.UpdatedAt)
	return i, err
}

// ---------------------------------------------------------------------------
// SearchActivityLog — keyset pagination over iam.audit_logs
// ---------------------------------------------------------------------------

const searchActivityLog = `-- name: SearchActivityLog :many
SELECT id, user_id, branch_id, resource, resource_id, action,
       old_value, new_value, ip, created_at
FROM iam.audit_logs
WHERE ($1::text = '' OR user_id::text = $1)
  AND ($2::text = '' OR resource = $2)
  AND ($3::text = '' OR action = $3)
  AND ($4::timestamptz IS NULL OR created_at >= $4)
  AND ($5::timestamptz IS NULL OR created_at <= $5)
  AND (
    $6::timestamptz IS NULL
    OR created_at < $6
    OR (created_at = $6 AND id::text < $7)
  )
ORDER BY created_at DESC, id DESC
LIMIT $8
`

type SearchActivityLogParams struct {
	UserID    string             `json:"user_id"`    // "" = no filter
	Resource  string             `json:"resource"`   // "" = no filter
	Action    string             `json:"action"`     // "" = no filter
	From      pgtype.Timestamptz `json:"from"`       // Valid=false = no lower bound
	To        pgtype.Timestamptz `json:"to"`         // Valid=false = no upper bound
	CursorAt  pgtype.Timestamptz `json:"cursor_at"`  // Valid=false = first page
	CursorID  string             `json:"cursor_id"`  // "" = first page
	Limit     int32              `json:"limit"`
}

func (q *Queries) SearchActivityLog(ctx context.Context, arg SearchActivityLogParams) ([]IamAuditLog, error) {
	rows, err := q.db.Query(ctx, searchActivityLog,
		arg.UserID, arg.Resource, arg.Action,
		arg.From, arg.To,
		arg.CursorAt, arg.CursorID,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []IamAuditLog
	for rows.Next() {
		var i IamAuditLog
		if err := rows.Scan(
			&i.ID, &i.UserID, &i.BranchID, &i.Resource, &i.ResourceID,
			&i.Action, &i.OldValue, &i.NewValue, &i.Ip, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}
