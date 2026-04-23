// Hand-written sqlc-style query implementations for catalog vendor readiness.
// Run `make generate` (sqlc generate) to regenerate from vendor_readiness.sql
// once sqlc is available.
//
// BL-OPS-020 — departure vendor readiness checklist.
// Table: catalog.departure_vendor_readiness (id, departure_id, kind, state,
// notes, attachment_url, updated_at, updated_by)

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// Row type
// ---------------------------------------------------------------------------

// VendorReadinessRow is the RETURNING / SELECT row for
// catalog.departure_vendor_readiness.
type VendorReadinessRow struct {
	ID            string             `json:"id"`
	DepartureID   string             `json:"departure_id"`
	Kind          string             `json:"kind"`
	State         string             `json:"state"`
	Notes         string             `json:"notes"`
	AttachmentURL string             `json:"attachment_url"`
	UpdatedAt     pgtype.Timestamptz `json:"updated_at"`
	UpdatedBy     string             `json:"updated_by"`
}

// ---------------------------------------------------------------------------
// UpsertVendorReadiness
// ---------------------------------------------------------------------------

const upsertVendorReadiness = `-- name: UpsertVendorReadiness :one
INSERT INTO catalog.departure_vendor_readiness (
    id, departure_id, kind, state, notes, attachment_url, updated_at, updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, now(), $7
)
ON CONFLICT (departure_id, kind) DO UPDATE SET
    state          = EXCLUDED.state,
    notes          = EXCLUDED.notes,
    attachment_url = EXCLUDED.attachment_url,
    updated_at     = now(),
    updated_by     = EXCLUDED.updated_by
RETURNING id, departure_id, kind, state, notes, attachment_url, updated_at, updated_by`

// UpsertVendorReadinessParams holds the input for UpsertVendorReadiness.
type UpsertVendorReadinessParams struct {
	ID            string `json:"id"`
	DepartureID   string `json:"departure_id"`
	Kind          string `json:"kind"`   // ticket | hotel | visa
	State         string `json:"state"`  // not_started | in_progress | done
	Notes         string `json:"notes"`
	AttachmentURL string `json:"attachment_url"`
	UpdatedBy     string `json:"updated_by"`
}

func (q *Queries) UpsertVendorReadiness(ctx context.Context, arg UpsertVendorReadinessParams) (VendorReadinessRow, error) {
	row := q.db.QueryRow(ctx, upsertVendorReadiness,
		arg.ID, arg.DepartureID, arg.Kind, arg.State, arg.Notes, arg.AttachmentURL, arg.UpdatedBy,
	)
	var i VendorReadinessRow
	err := row.Scan(
		&i.ID, &i.DepartureID, &i.Kind, &i.State,
		&i.Notes, &i.AttachmentURL, &i.UpdatedAt, &i.UpdatedBy,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// ListVendorReadiness
// ---------------------------------------------------------------------------

const listVendorReadiness = `-- name: ListVendorReadiness :many
SELECT id, departure_id, kind, state, notes, attachment_url, updated_at, updated_by
FROM catalog.departure_vendor_readiness
WHERE departure_id = $1
ORDER BY kind`

func (q *Queries) ListVendorReadiness(ctx context.Context, departureID string) ([]VendorReadinessRow, error) {
	rows, err := q.db.Query(ctx, listVendorReadiness, departureID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []VendorReadinessRow
	for rows.Next() {
		var i VendorReadinessRow
		if err := rows.Scan(
			&i.ID, &i.DepartureID, &i.Kind, &i.State,
			&i.Notes, &i.AttachmentURL, &i.UpdatedAt, &i.UpdatedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}
