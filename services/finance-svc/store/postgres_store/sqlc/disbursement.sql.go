// disbursement.sql.go — hand-written SQLC-style queries for AP disbursement
// batches and AR/AP aging (BL-FIN-010/011).

package sqlc

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// Models
// ---------------------------------------------------------------------------

type DisbursementBatchRow struct {
	ID             string             `json:"id"`
	Description    string             `json:"description"`
	TotalAmountIdr int64              `json:"total_amount_idr"`
	Status         string             `json:"status"`
	ApprovedBy     pgtype.Text        `json:"approved_by"`
	ApprovedAt     pgtype.Timestamptz `json:"approved_at"`
	CreatedBy      string             `json:"created_by"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
}

type DisbursementItemRow struct {
	ID             int64       `json:"id"`
	BatchID        string      `json:"batch_id"`
	VendorName     string      `json:"vendor_name"`
	Description    string      `json:"description"`
	AmountIdr      int64       `json:"amount_idr"`
	Reference      pgtype.Text `json:"reference"`
	JournalEntryID pgtype.Text `json:"journal_entry_id"`
}

// ---------------------------------------------------------------------------
// InsertDisbursementBatch
// ---------------------------------------------------------------------------

type InsertDisbursementBatchParams struct {
	ID             string
	Description    string
	TotalAmountIdr int64
	CreatedBy      string
}

const insertDisbursementBatch = `
INSERT INTO finance.disbursement_batches (id, description, total_amount_idr, status, created_by)
VALUES ($1, $2, $3, 'pending_approval', $4)
RETURNING id, description, total_amount_idr, status, approved_by, approved_at, created_by, created_at
`

func (q *Queries) InsertDisbursementBatch(ctx context.Context, arg InsertDisbursementBatchParams) (DisbursementBatchRow, error) {
	row := q.db.QueryRow(ctx, insertDisbursementBatch,
		arg.ID, arg.Description, arg.TotalAmountIdr, arg.CreatedBy,
	)
	var i DisbursementBatchRow
	err := row.Scan(
		&i.ID, &i.Description, &i.TotalAmountIdr, &i.Status,
		&i.ApprovedBy, &i.ApprovedAt, &i.CreatedBy, &i.CreatedAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// InsertDisbursementItem
// ---------------------------------------------------------------------------

type InsertDisbursementItemParams struct {
	BatchID     string
	VendorName  string
	Description string
	AmountIdr   int64
	Reference   string
}

const insertDisbursementItem = `
INSERT INTO finance.disbursement_items (batch_id, vendor_name, description, amount_idr, reference)
VALUES ($1, $2, $3, $4, NULLIF($5, ''))
RETURNING id, batch_id, vendor_name, description, amount_idr, reference, journal_entry_id
`

func (q *Queries) InsertDisbursementItem(ctx context.Context, arg InsertDisbursementItemParams) (DisbursementItemRow, error) {
	row := q.db.QueryRow(ctx, insertDisbursementItem,
		arg.BatchID, arg.VendorName, arg.Description, arg.AmountIdr, arg.Reference,
	)
	var i DisbursementItemRow
	err := row.Scan(&i.ID, &i.BatchID, &i.VendorName, &i.Description, &i.AmountIdr, &i.Reference, &i.JournalEntryID)
	return i, err
}

// ---------------------------------------------------------------------------
// GetDisbursementBatchByID
// ---------------------------------------------------------------------------

const getDisbursementBatchByID = `
SELECT id, description, total_amount_idr, status, approved_by, approved_at, created_by, created_at
FROM finance.disbursement_batches
WHERE id = $1
`

func (q *Queries) GetDisbursementBatchByID(ctx context.Context, id string) (DisbursementBatchRow, error) {
	row := q.db.QueryRow(ctx, getDisbursementBatchByID, id)
	var i DisbursementBatchRow
	err := row.Scan(
		&i.ID, &i.Description, &i.TotalAmountIdr, &i.Status,
		&i.ApprovedBy, &i.ApprovedAt, &i.CreatedBy, &i.CreatedAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// GetDisbursementItemsByBatchID
// ---------------------------------------------------------------------------

const getDisbursementItemsByBatch = `
SELECT id, batch_id, vendor_name, description, amount_idr, reference, journal_entry_id
FROM finance.disbursement_items
WHERE batch_id = $1
ORDER BY id ASC
`

func (q *Queries) GetDisbursementItemsByBatchID(ctx context.Context, batchID string) ([]DisbursementItemRow, error) {
	rows, err := q.db.Query(ctx, getDisbursementItemsByBatch, batchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []DisbursementItemRow
	for rows.Next() {
		var i DisbursementItemRow
		if err := rows.Scan(&i.ID, &i.BatchID, &i.VendorName, &i.Description, &i.AmountIdr, &i.Reference, &i.JournalEntryID); err != nil {
			return nil, err
		}
		result = append(result, i)
	}
	return result, rows.Err()
}

// ---------------------------------------------------------------------------
// UpdateDisbursementBatchDecision
// ---------------------------------------------------------------------------

type UpdateDisbursementBatchDecisionParams struct {
	ID         string
	NewStatus  string
	ApprovedBy string
}

const updateDisbursementBatchDecision = `
UPDATE finance.disbursement_batches
SET status = $1, approved_by = NULLIF($2,''), approved_at = CASE WHEN $1 = 'approved' THEN NOW() ELSE NULL END
WHERE id = $3
`

func (q *Queries) UpdateDisbursementBatchDecision(ctx context.Context, arg UpdateDisbursementBatchDecisionParams) error {
	_, err := q.db.Exec(ctx, updateDisbursementBatchDecision, arg.NewStatus, arg.ApprovedBy, arg.ID)
	return err
}

// ---------------------------------------------------------------------------
// UpdateDisbursementItemJournal
// ---------------------------------------------------------------------------

const updateDisbursementItemJournal = `
UPDATE finance.disbursement_items SET journal_entry_id = $1 WHERE id = $2
`

func (q *Queries) UpdateDisbursementItemJournal(ctx context.Context, journalEntryID string, itemID int64) error {
	_, err := q.db.Exec(ctx, updateDisbursementItemJournal, journalEntryID, itemID)
	return err
}

// ---------------------------------------------------------------------------
// GetARAPAging
// ---------------------------------------------------------------------------

// AgingBuckets holds aggregated IDR amounts per aging bucket.
type AgingBuckets struct {
	Current int64
	Days30  int64
	Days60  int64
	Days90  int64
	Over90  int64
	Total   int64
}

// GetDisbursementAgingBuckets computes AP aging from disbursement_batches
// that are still in pending_approval status, using created_at as aging anchor.
const getDisbursementAgingBuckets = `
SELECT
  COALESCE(SUM(total_amount_idr) FILTER (WHERE ($1::date - created_at::date) <= 0), 0)   AS current,
  COALESCE(SUM(total_amount_idr) FILTER (WHERE ($1::date - created_at::date) BETWEEN 1 AND 30), 0)   AS days_30,
  COALESCE(SUM(total_amount_idr) FILTER (WHERE ($1::date - created_at::date) BETWEEN 31 AND 60), 0)  AS days_60,
  COALESCE(SUM(total_amount_idr) FILTER (WHERE ($1::date - created_at::date) BETWEEN 61 AND 90), 0)  AS days_90,
  COALESCE(SUM(total_amount_idr) FILTER (WHERE ($1::date - created_at::date) > 90), 0)               AS over_90
FROM finance.disbursement_batches
WHERE status = 'pending_approval'
`

func (q *Queries) GetDisbursementAgingBuckets(ctx context.Context, asOfDate time.Time) (AgingBuckets, error) {
	row := q.db.QueryRow(ctx, getDisbursementAgingBuckets, pgtype.Date{Time: asOfDate, Valid: true})
	var b AgingBuckets
	err := row.Scan(&b.Current, &b.Days30, &b.Days60, &b.Days90, &b.Over90)
	if err != nil {
		return AgingBuckets{}, err
	}
	b.Total = b.Current + b.Days30 + b.Days60 + b.Days90 + b.Over90
	return b, nil
}

// GetARAgingBuckets computes AR aging from journal entries credit lines (AR accounts = 1200-1299).
const getARAgingBuckets = `
SELECT
  COALESCE(SUM(jl.credit_idr) FILTER (WHERE ($1::date - je.created_at::date) <= 0), 0)   AS current,
  COALESCE(SUM(jl.credit_idr) FILTER (WHERE ($1::date - je.created_at::date) BETWEEN 1 AND 30), 0)  AS days_30,
  COALESCE(SUM(jl.credit_idr) FILTER (WHERE ($1::date - je.created_at::date) BETWEEN 31 AND 60), 0) AS days_60,
  COALESCE(SUM(jl.credit_idr) FILTER (WHERE ($1::date - je.created_at::date) BETWEEN 61 AND 90), 0) AS days_90,
  COALESCE(SUM(jl.credit_idr) FILTER (WHERE ($1::date - je.created_at::date) > 90), 0)              AS over_90
FROM finance.journal_lines jl
JOIN finance.journal_entries je ON je.id = jl.entry_id
WHERE jl.account_code LIKE '12%'
`

func (q *Queries) GetARAgingBuckets(ctx context.Context, asOfDate time.Time) (AgingBuckets, error) {
	row := q.db.QueryRow(ctx, getARAgingBuckets, pgtype.Date{Time: asOfDate, Valid: true})
	var b AgingBuckets
	err := row.Scan(&b.Current, &b.Days30, &b.Days60, &b.Days90, &b.Over90)
	if err != nil {
		return AgingBuckets{}, err
	}
	b.Total = b.Current + b.Days30 + b.Days60 + b.Days90 + b.Over90
	return b, nil
}
