// shipments.sql.go — hand-written sqlc-style query implementations for
// logistics.shipments table (migration 000019).
//
// BL-LOG-002 / ShipFulfillmentTask.

package sqlc

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// Row type
// ---------------------------------------------------------------------------

// ShipmentRow mirrors a row from logistics.shipments.
type ShipmentRow struct {
	ID             string             `json:"id"`
	TaskID         string             `json:"task_id"`
	TrackingNumber string             `json:"tracking_number"`
	Carrier        string             `json:"carrier"`
	Status         string             `json:"status"`
	ShippedAt      pgtype.Timestamptz `json:"shipped_at"`
	DeliveredAt    pgtype.Timestamptz `json:"delivered_at"`
	Notes          pgtype.Text        `json:"notes"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

// InsertShipmentParams holds inputs for InsertShipment.
type InsertShipmentParams struct {
	TaskID         string
	TrackingNumber string
	Carrier        string
	Notes          string
}

// UpdateFulfillmentTaskStatusParams holds inputs for UpdateFulfillmentTaskStatus.
type UpdateFulfillmentTaskStatusParams struct {
	ID     string
	Status string
}

// ---------------------------------------------------------------------------
// Query implementations
// ---------------------------------------------------------------------------

const insertShipment = `-- name: InsertShipment :one
INSERT INTO logistics.shipments (
    task_id,
    tracking_number,
    carrier,
    notes,
    status
) VALUES (
    $1,
    $2,
    $3,
    NULLIF($4, ''),
    'shipped'
)
RETURNING id, task_id, tracking_number, carrier, status,
          shipped_at, delivered_at, notes, created_at`

// InsertShipment inserts a new shipment record and returns the created row.
func (q *Queries) InsertShipment(ctx context.Context, arg InsertShipmentParams) (ShipmentRow, error) {
	row := q.db.QueryRow(ctx, insertShipment,
		arg.TaskID,
		arg.TrackingNumber,
		arg.Carrier,
		arg.Notes,
	)
	var r ShipmentRow
	err := row.Scan(
		&r.ID, &r.TaskID, &r.TrackingNumber, &r.Carrier, &r.Status,
		&r.ShippedAt, &r.DeliveredAt, &r.Notes, &r.CreatedAt,
	)
	return r, err
}

const getShipmentByTaskID = `-- name: GetShipmentByTaskID :one
SELECT id, task_id, tracking_number, carrier, status,
       shipped_at, delivered_at, notes, created_at
FROM logistics.shipments
WHERE task_id = $1
LIMIT 1`

// GetShipmentByTaskID returns the shipment for the given task_id,
// or pgx.ErrNoRows if none exists.
func (q *Queries) GetShipmentByTaskID(ctx context.Context, taskID string) (ShipmentRow, error) {
	row := q.db.QueryRow(ctx, getShipmentByTaskID, taskID)
	var r ShipmentRow
	err := row.Scan(
		&r.ID, &r.TaskID, &r.TrackingNumber, &r.Carrier, &r.Status,
		&r.ShippedAt, &r.DeliveredAt, &r.Notes, &r.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return ShipmentRow{}, pgx.ErrNoRows
	}
	return r, err
}

const updateFulfillmentTaskStatus = `-- name: UpdateFulfillmentTaskStatus :exec
UPDATE logistics.fulfillment_tasks
SET status = $1, updated_at = now()
WHERE id = $2`

// UpdateFulfillmentTaskStatus sets a new status on a fulfillment task.
func (q *Queries) UpdateFulfillmentTaskStatus(ctx context.Context, arg UpdateFulfillmentTaskStatusParams) error {
	_, err := q.db.Exec(ctx, updateFulfillmentTaskStatus, arg.Status, arg.ID)
	return err
}
