// journal.sql.go — hand-written sqlc-style query implementations for
// finance-svc journal entry queries.
//
// Run `make generate` (sqlc generate) to regenerate from journal.sql once
// sqlc is configured to target the finance schema.
//
// S3-E-03 / BL-FIN-001..003.

package sqlc

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// Row types
// ---------------------------------------------------------------------------

// JournalEntryRow mirrors a row from finance.journal_entries.
type JournalEntryRow struct {
	ID             string             `json:"id"`
	IdempotencyKey string             `json:"idempotency_key"`
	SourceType     string             `json:"source_type"`
	SourceID       string             `json:"source_id"`
	PostedAt       pgtype.Timestamptz `json:"posted_at"`
	Description    pgtype.Text        `json:"description"`
}

// JournalLineRow mirrors a row from finance.journal_lines.
type JournalLineRow struct {
	ID          string         `json:"id"`
	EntryID     string         `json:"entry_id"`
	AccountCode string         `json:"account_code"`
	Debit       pgtype.Numeric `json:"debit"`
	Credit      pgtype.Numeric `json:"credit"`
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

// InsertJournalEntryParams holds inputs for InsertJournalEntry.
type InsertJournalEntryParams struct {
	IdempotencyKey string
	SourceType     string
	SourceID       string
	Description    pgtype.Text
}

// InsertJournalLineParams holds inputs for InsertJournalLine.
type InsertJournalLineParams struct {
	EntryID     string
	AccountCode string
	Debit       pgtype.Numeric
	Credit      pgtype.Numeric
}

// ---------------------------------------------------------------------------
// Query implementations
// ---------------------------------------------------------------------------

const getJournalEntryByIdempotencyKey = `-- name: GetJournalEntryByIdempotencyKey :one
SELECT id, idempotency_key, source_type, source_id, posted_at, description
FROM finance.journal_entries
WHERE idempotency_key = $1
LIMIT 1`

// GetJournalEntryByIdempotencyKey returns a journal entry for the given key,
// or pgx.ErrNoRows if none exists.
func (q *Queries) GetJournalEntryByIdempotencyKey(ctx context.Context, idempotencyKey string) (JournalEntryRow, error) {
	row := q.db.QueryRow(ctx, getJournalEntryByIdempotencyKey, idempotencyKey)
	var r JournalEntryRow
	err := row.Scan(
		&r.ID,
		&r.IdempotencyKey,
		&r.SourceType,
		&r.SourceID,
		&r.PostedAt,
		&r.Description,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return JournalEntryRow{}, pgx.ErrNoRows
	}
	return r, err
}

const getJournalEntryByID = `-- name: GetJournalEntryByID :one
SELECT id, idempotency_key, source_type, source_id, posted_at, description
FROM finance.journal_entries
WHERE id = $1
LIMIT 1`

// GetJournalEntryByID returns a journal entry for the given UUID string ID,
// or pgx.ErrNoRows if none exists.
func (q *Queries) GetJournalEntryByID(ctx context.Context, id string) (JournalEntryRow, error) {
	row := q.db.QueryRow(ctx, getJournalEntryByID, id)
	var r JournalEntryRow
	err := row.Scan(
		&r.ID,
		&r.IdempotencyKey,
		&r.SourceType,
		&r.SourceID,
		&r.PostedAt,
		&r.Description,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return JournalEntryRow{}, pgx.ErrNoRows
	}
	return r, err
}

const insertJournalEntry = `-- name: InsertJournalEntry :one
INSERT INTO finance.journal_entries (
    idempotency_key,
    source_type,
    source_id,
    description
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, idempotency_key, source_type, source_id, posted_at, description`

// InsertJournalEntry inserts a new journal entry header.
func (q *Queries) InsertJournalEntry(ctx context.Context, arg InsertJournalEntryParams) (JournalEntryRow, error) {
	row := q.db.QueryRow(ctx, insertJournalEntry,
		arg.IdempotencyKey,
		arg.SourceType,
		arg.SourceID,
		arg.Description,
	)
	var r JournalEntryRow
	err := row.Scan(
		&r.ID,
		&r.IdempotencyKey,
		&r.SourceType,
		&r.SourceID,
		&r.PostedAt,
		&r.Description,
	)
	return r, err
}

const insertJournalLine = `-- name: InsertJournalLine :one
INSERT INTO finance.journal_lines (
    entry_id,
    account_code,
    debit,
    credit
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, entry_id, account_code, debit, credit`

// InsertJournalLine inserts one journal line for the given entry.
func (q *Queries) InsertJournalLine(ctx context.Context, arg InsertJournalLineParams) (JournalLineRow, error) {
	row := q.db.QueryRow(ctx, insertJournalLine,
		arg.EntryID,
		arg.AccountCode,
		arg.Debit,
		arg.Credit,
	)
	var r JournalLineRow
	err := row.Scan(
		&r.ID,
		&r.EntryID,
		&r.AccountCode,
		&r.Debit,
		&r.Credit,
	)
	return r, err
}
