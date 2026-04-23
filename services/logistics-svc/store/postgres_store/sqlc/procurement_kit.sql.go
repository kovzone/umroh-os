// procurement_kit.sql.go — hand-written SQLC-style queries for purchase
// requests and kit assembly (BL-LOG-010..012).

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// Models
// ---------------------------------------------------------------------------

type PurchaseRequestRow struct {
	ID            string             `json:"id"`
	DepartureID   string             `json:"departure_id"`
	RequestedBy   string             `json:"requested_by"`
	ItemName      string             `json:"item_name"`
	Quantity      int32              `json:"quantity"`
	UnitPriceIdr  int64              `json:"unit_price_idr"`
	TotalPriceIdr int64              `json:"total_price_idr"`
	BudgetLimit   pgtype.Int8        `json:"budget_limit_idr"`
	Status        string             `json:"status"`
	ApprovedBy    pgtype.Text        `json:"approved_by"`
	Notes         pgtype.Text        `json:"notes"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
}

type KitAssemblyRow struct {
	ID             string             `json:"id"`
	DepartureID    string             `json:"departure_id"`
	AssembledBy    string             `json:"assembled_by"`
	Status         string             `json:"status"`
	IdempotencyKey string             `json:"idempotency_key"`
	AssembledAt    pgtype.Timestamptz `json:"assembled_at"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
}

// ---------------------------------------------------------------------------
// InsertPurchaseRequest
// ---------------------------------------------------------------------------

type InsertPurchaseRequestParams struct {
	ID            string
	DepartureID   string
	RequestedBy   string
	ItemName      string
	Quantity      int32
	UnitPriceIdr  int64
	TotalPriceIdr int64
	BudgetLimit   int64 // 0 = no limit
}

const insertPurchaseRequest = `
INSERT INTO logistics.purchase_requests
  (id, departure_id, requested_by, item_name, quantity, unit_price_idr, total_price_idr, budget_limit_idr, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, NULLIF($8, 0), 'pending')
RETURNING id, departure_id, requested_by, item_name, quantity, unit_price_idr, total_price_idr, budget_limit_idr, status, approved_by, notes, created_at
`

func (q *Queries) InsertPurchaseRequest(ctx context.Context, arg InsertPurchaseRequestParams) (PurchaseRequestRow, error) {
	row := q.db.QueryRow(ctx, insertPurchaseRequest,
		arg.ID, arg.DepartureID, arg.RequestedBy, arg.ItemName,
		arg.Quantity, arg.UnitPriceIdr, arg.TotalPriceIdr, arg.BudgetLimit,
	)
	var i PurchaseRequestRow
	err := row.Scan(
		&i.ID, &i.DepartureID, &i.RequestedBy, &i.ItemName,
		&i.Quantity, &i.UnitPriceIdr, &i.TotalPriceIdr, &i.BudgetLimit,
		&i.Status, &i.ApprovedBy, &i.Notes, &i.CreatedAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// GetPurchaseRequestByID
// ---------------------------------------------------------------------------

const getPurchaseRequestByID = `
SELECT id, departure_id, requested_by, item_name, quantity, unit_price_idr, total_price_idr, budget_limit_idr, status, approved_by, notes, created_at
FROM logistics.purchase_requests
WHERE id = $1
`

func (q *Queries) GetPurchaseRequestByID(ctx context.Context, id string) (PurchaseRequestRow, error) {
	row := q.db.QueryRow(ctx, getPurchaseRequestByID, id)
	var i PurchaseRequestRow
	err := row.Scan(
		&i.ID, &i.DepartureID, &i.RequestedBy, &i.ItemName,
		&i.Quantity, &i.UnitPriceIdr, &i.TotalPriceIdr, &i.BudgetLimit,
		&i.Status, &i.ApprovedBy, &i.Notes, &i.CreatedAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// UpdatePurchaseRequestDecision
// ---------------------------------------------------------------------------

type UpdatePurchaseRequestDecisionParams struct {
	ID         string
	NewStatus  string
	ApprovedBy string
	Notes      string
}

const updatePurchaseRequestDecision = `
UPDATE logistics.purchase_requests
SET status = $1, approved_by = $2, notes = $3
WHERE id = $4
`

func (q *Queries) UpdatePurchaseRequestDecision(ctx context.Context, arg UpdatePurchaseRequestDecisionParams) error {
	_, err := q.db.Exec(ctx, updatePurchaseRequestDecision,
		arg.NewStatus, nullText(arg.ApprovedBy), nullText(arg.Notes), arg.ID,
	)
	return err
}

// ---------------------------------------------------------------------------
// InsertKitAssembly
// ---------------------------------------------------------------------------

type InsertKitAssemblyParams struct {
	ID             string
	DepartureID    string
	AssembledBy    string
	Status         string
	IdempotencyKey string
}

const insertKitAssembly = `
INSERT INTO logistics.kit_assemblies (id, departure_id, assembled_by, status, idempotency_key)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (idempotency_key) DO NOTHING
RETURNING id, departure_id, assembled_by, status, idempotency_key, assembled_at, created_at
`

const getKitAssemblyByIdempotencyKey = `
SELECT id, departure_id, assembled_by, status, idempotency_key, assembled_at, created_at
FROM logistics.kit_assemblies
WHERE idempotency_key = $1
`

func (q *Queries) InsertKitAssemblyIdempotent(ctx context.Context, arg InsertKitAssemblyParams) (KitAssemblyRow, bool, error) {
	row := q.db.QueryRow(ctx, insertKitAssembly,
		arg.ID, arg.DepartureID, arg.AssembledBy, arg.Status, arg.IdempotencyKey,
	)
	var i KitAssemblyRow
	err := row.Scan(&i.ID, &i.DepartureID, &i.AssembledBy, &i.Status, &i.IdempotencyKey, &i.AssembledAt, &i.CreatedAt)
	if err != nil {
		// ON CONFLICT DO NOTHING — fetch existing
		row2 := q.db.QueryRow(ctx, getKitAssemblyByIdempotencyKey, arg.IdempotencyKey)
		var existing KitAssemblyRow
		if err2 := row2.Scan(&existing.ID, &existing.DepartureID, &existing.AssembledBy, &existing.Status, &existing.IdempotencyKey, &existing.AssembledAt, &existing.CreatedAt); err2 != nil {
			return KitAssemblyRow{}, false, err2
		}
		return existing, false, nil
	}
	return i, true, nil
}

// ---------------------------------------------------------------------------
// InsertKitAssemblyItem
// ---------------------------------------------------------------------------

type InsertKitAssemblyItemParams struct {
	AssemblyID string
	ItemName   string
	Quantity   int32
	Fulfilled  bool
}

const insertKitAssemblyItem = `
INSERT INTO logistics.kit_assembly_items (assembly_id, item_name, quantity, fulfilled)
VALUES ($1, $2, $3, $4)
`

func (q *Queries) InsertKitAssemblyItem(ctx context.Context, arg InsertKitAssemblyItemParams) error {
	_, err := q.db.Exec(ctx, insertKitAssemblyItem,
		arg.AssemblyID, arg.ItemName, arg.Quantity, arg.Fulfilled,
	)
	return err
}

// ---------------------------------------------------------------------------
// UpdateKitAssemblyStatus
// ---------------------------------------------------------------------------

const updateKitAssemblyStatus = `
UPDATE logistics.kit_assemblies
SET status = $1, assembled_at = CASE WHEN $1 = 'completed' THEN NOW() ELSE assembled_at END
WHERE id = $2
`

func (q *Queries) UpdateKitAssemblyStatus(ctx context.Context, newStatus, id string) error {
	_, err := q.db.Exec(ctx, updateKitAssemblyStatus, newStatus, id)
	return err
}

// ---------------------------------------------------------------------------
// nullText helper
// ---------------------------------------------------------------------------

func nullText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}
