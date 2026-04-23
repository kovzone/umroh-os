// reports.sql.go — hand-written sqlc-style query implementations for
// finance-svc report queries (S5-E-01).
//
// Run `make generate` (sqlc generate) to regenerate from reports.sql once
// sqlc is configured to target the finance schema.

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// Row types
// ---------------------------------------------------------------------------

// AccountSummaryRow holds aggregated debit/credit per account_code.
type AccountSummaryRow struct {
	AccountCode  string         `json:"account_code"`
	DebitTotal   pgtype.Numeric `json:"debit_total"`
	CreditTotal  pgtype.Numeric `json:"credit_total"`
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

// ListJournalEntriesParams holds filter + pagination inputs for ListJournalEntries.
type ListJournalEntriesParams struct {
	From   pgtype.Timestamptz // optional; zero = no lower bound
	To     pgtype.Timestamptz // optional; zero = no upper bound
	Cursor pgtype.Timestamptz // optional; zero = no cursor (first page)
	Limit  int32
}

// GetJournalLinesParams holds the entry IDs for GetJournalLines.
type GetJournalLinesParams struct {
	EntryIDs []string // list of UUID strings
}

// ---------------------------------------------------------------------------
// Query implementations
// ---------------------------------------------------------------------------

const getFinanceSummary = `-- name: GetFinanceSummary :many
SELECT
    account_code,
    SUM(debit)  AS debit_total,
    SUM(credit) AS credit_total
FROM finance.journal_lines
GROUP BY account_code
ORDER BY account_code`

// GetFinanceSummary returns per-account aggregated debit/credit totals.
func (q *Queries) GetFinanceSummary(ctx context.Context) ([]AccountSummaryRow, error) {
	rows, err := q.db.Query(ctx, getFinanceSummary)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []AccountSummaryRow
	for rows.Next() {
		var r AccountSummaryRow
		if err := rows.Scan(&r.AccountCode, &r.DebitTotal, &r.CreditTotal); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

const listJournalEntries = `-- name: ListJournalEntries :many
SELECT
    je.id,
    je.idempotency_key,
    je.source_type,
    je.source_id,
    je.posted_at,
    je.description
FROM finance.journal_entries je
WHERE
    ($1::timestamptz IS NULL OR je.posted_at >= $1)
    AND ($2::timestamptz IS NULL OR je.posted_at <= $2)
    AND ($3::timestamptz IS NULL OR je.posted_at < $3)
ORDER BY je.posted_at DESC
LIMIT $4`

// ListJournalEntries returns a paginated, optionally date-filtered list of journal entries.
func (q *Queries) ListJournalEntries(ctx context.Context, arg ListJournalEntriesParams) ([]JournalEntryRow, error) {
	rows, err := q.db.Query(ctx, listJournalEntries,
		arg.From,
		arg.To,
		arg.Cursor,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []JournalEntryRow
	for rows.Next() {
		var r JournalEntryRow
		if err := rows.Scan(
			&r.ID,
			&r.IdempotencyKey,
			&r.SourceType,
			&r.SourceID,
			&r.PostedAt,
			&r.Description,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

const getJournalLines = `-- name: GetJournalLines :many
SELECT
    jl.id,
    jl.entry_id,
    jl.account_code,
    jl.debit,
    jl.credit
FROM finance.journal_lines jl
WHERE jl.entry_id = ANY($1::uuid[])
ORDER BY jl.entry_id, jl.id`

// GetJournalLines returns all lines for the given list of entry IDs.
func (q *Queries) GetJournalLines(ctx context.Context, entryIDs []string) ([]JournalLineRow, error) {
	rows, err := q.db.Query(ctx, getJournalLines, entryIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []JournalLineRow
	for rows.Next() {
		var r JournalLineRow
		if err := rows.Scan(
			&r.ID,
			&r.EntryID,
			&r.AccountCode,
			&r.Debit,
			&r.Credit,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}
