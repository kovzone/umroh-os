// pickup_tokens.sql.go — hand-written sqlc-style query implementations for
// logistics.pickup_tokens table (migration 000019).
//
// BL-LOG-003 / GeneratePickupQR + RedeemPickupQR.

package sqlc

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// Row type
// ---------------------------------------------------------------------------

// PickupTokenRow mirrors a row from logistics.pickup_tokens.
type PickupTokenRow struct {
	ID        string             `json:"id"`
	TaskID    string             `json:"task_id"`
	Token     string             `json:"token"`
	Used      bool               `json:"used"`
	UsedAt    pgtype.Timestamptz `json:"used_at"`
	ExpiresAt time.Time          `json:"expires_at"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

// InsertPickupTokenParams holds inputs for InsertPickupToken.
type InsertPickupTokenParams struct {
	TaskID    string
	Token     string
	ExpiresAt time.Time
}

// ---------------------------------------------------------------------------
// Query implementations
// ---------------------------------------------------------------------------

const insertPickupToken = `-- name: InsertPickupToken :one
INSERT INTO logistics.pickup_tokens (
    task_id,
    token,
    expires_at
) VALUES (
    $1,
    $2,
    $3
)
RETURNING id, task_id, token, used, used_at, expires_at, created_at`

// InsertPickupToken inserts a new pickup token and returns the created row.
func (q *Queries) InsertPickupToken(ctx context.Context, arg InsertPickupTokenParams) (PickupTokenRow, error) {
	row := q.db.QueryRow(ctx, insertPickupToken,
		arg.TaskID,
		arg.Token,
		arg.ExpiresAt,
	)
	var r PickupTokenRow
	err := row.Scan(
		&r.ID, &r.TaskID, &r.Token, &r.Used, &r.UsedAt, &r.ExpiresAt, &r.CreatedAt,
	)
	return r, err
}

const getPickupTokenByToken = `-- name: GetPickupTokenByToken :one
SELECT id, task_id, token, used, used_at, expires_at, created_at
FROM logistics.pickup_tokens
WHERE token = $1
LIMIT 1`

// GetPickupTokenByToken returns the pickup token for the given token string,
// or pgx.ErrNoRows if none exists.
func (q *Queries) GetPickupTokenByToken(ctx context.Context, token string) (PickupTokenRow, error) {
	row := q.db.QueryRow(ctx, getPickupTokenByToken, token)
	var r PickupTokenRow
	err := row.Scan(
		&r.ID, &r.TaskID, &r.Token, &r.Used, &r.UsedAt, &r.ExpiresAt, &r.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return PickupTokenRow{}, pgx.ErrNoRows
	}
	return r, err
}

const getActivePickupTokenByTaskID = `-- name: GetActivePickupTokenByTaskID :one
SELECT id, task_id, token, used, used_at, expires_at, created_at
FROM logistics.pickup_tokens
WHERE task_id = $1
  AND used = false
  AND expires_at > now()
ORDER BY created_at DESC
LIMIT 1`

// GetActivePickupTokenByTaskID returns the most recent non-used, non-expired
// pickup token for the given task, or pgx.ErrNoRows if none exists.
func (q *Queries) GetActivePickupTokenByTaskID(ctx context.Context, taskID string) (PickupTokenRow, error) {
	row := q.db.QueryRow(ctx, getActivePickupTokenByTaskID, taskID)
	var r PickupTokenRow
	err := row.Scan(
		&r.ID, &r.TaskID, &r.Token, &r.Used, &r.UsedAt, &r.ExpiresAt, &r.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return PickupTokenRow{}, pgx.ErrNoRows
	}
	return r, err
}

const markPickupTokenUsed = `-- name: MarkPickupTokenUsed :one
UPDATE logistics.pickup_tokens
SET used = true, used_at = now()
WHERE id = $1
RETURNING id, task_id, token, used, used_at, expires_at, created_at`

// MarkPickupTokenUsed marks the given token as used and records used_at timestamp.
func (q *Queries) MarkPickupTokenUsed(ctx context.Context, id string) (PickupTokenRow, error) {
	row := q.db.QueryRow(ctx, markPickupTokenUsed, id)
	var r PickupTokenRow
	err := row.Scan(
		&r.ID, &r.TaskID, &r.Token, &r.Used, &r.UsedAt, &r.ExpiresAt, &r.CreatedAt,
	)
	return r, err
}
