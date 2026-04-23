// fulfillment_tasks.sql.go — hand-written sqlc-style query implementations for
// logistics-svc fulfillment task queries.
//
// Run `make generate` (sqlc generate) to regenerate from fulfillment_tasks.sql
// once sqlc is configured to target the logistics schema.
//
// S3-E-02 / BL-LOG-001.

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

// FulfillmentTaskRow mirrors a row from logistics.fulfillment_tasks.
type FulfillmentTaskRow struct {
	ID             string             `json:"id"`
	BookingID      string             `json:"booking_id"`
	DepartureID    string             `json:"departure_id"`
	Status         string             `json:"status"`
	TrackingNumber pgtype.Text        `json:"tracking_number"`
	ShippedAt      pgtype.Timestamptz `json:"shipped_at"`
	DeliveredAt    pgtype.Timestamptz `json:"delivered_at"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
	UpdatedAt      pgtype.Timestamptz `json:"updated_at"`
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

// InsertFulfillmentTaskParams holds inputs for InsertFulfillmentTask.
type InsertFulfillmentTaskParams struct {
	BookingID   string
	DepartureID string
}

// ---------------------------------------------------------------------------
// Query implementations
// ---------------------------------------------------------------------------

const getFulfillmentTaskByBookingID = `-- name: GetFulfillmentTaskByBookingID :one
SELECT id, booking_id, departure_id, status, tracking_number,
       shipped_at, delivered_at, created_at, updated_at
FROM logistics.fulfillment_tasks
WHERE booking_id = $1
LIMIT 1`

// GetFulfillmentTaskByBookingID returns the fulfillment task for the given
// booking_id, or pgx.ErrNoRows if none exists.
func (q *Queries) GetFulfillmentTaskByBookingID(ctx context.Context, bookingID string) (FulfillmentTaskRow, error) {
	row := q.db.QueryRow(ctx, getFulfillmentTaskByBookingID, bookingID)
	var r FulfillmentTaskRow
	err := row.Scan(
		&r.ID,
		&r.BookingID,
		&r.DepartureID,
		&r.Status,
		&r.TrackingNumber,
		&r.ShippedAt,
		&r.DeliveredAt,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return FulfillmentTaskRow{}, pgx.ErrNoRows
	}
	return r, err
}

const getFulfillmentTaskByID = `-- name: GetFulfillmentTaskByID :one
SELECT id, booking_id, departure_id, status, tracking_number,
       shipped_at, delivered_at, created_at, updated_at
FROM logistics.fulfillment_tasks
WHERE id = $1
LIMIT 1`

// GetFulfillmentTaskByID returns the fulfillment task for the given id,
// or pgx.ErrNoRows if none exists.
func (q *Queries) GetFulfillmentTaskByID(ctx context.Context, id string) (FulfillmentTaskRow, error) {
	row := q.db.QueryRow(ctx, getFulfillmentTaskByID, id)
	var r FulfillmentTaskRow
	err := row.Scan(
		&r.ID,
		&r.BookingID,
		&r.DepartureID,
		&r.Status,
		&r.TrackingNumber,
		&r.ShippedAt,
		&r.DeliveredAt,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return FulfillmentTaskRow{}, pgx.ErrNoRows
	}
	return r, err
}

const insertFulfillmentTask = `-- name: InsertFulfillmentTask :one
INSERT INTO logistics.fulfillment_tasks (
    booking_id,
    departure_id,
    status
) VALUES (
    $1,
    $2,
    'queued'
)
RETURNING id, booking_id, departure_id, status, tracking_number,
          shipped_at, delivered_at, created_at, updated_at`

// InsertFulfillmentTask inserts a new fulfillment task in status='queued' and
// returns the created row.
func (q *Queries) InsertFulfillmentTask(ctx context.Context, arg InsertFulfillmentTaskParams) (FulfillmentTaskRow, error) {
	row := q.db.QueryRow(ctx, insertFulfillmentTask, arg.BookingID, arg.DepartureID)
	var r FulfillmentTaskRow
	err := row.Scan(
		&r.ID,
		&r.BookingID,
		&r.DepartureID,
		&r.Status,
		&r.TrackingNumber,
		&r.ShippedAt,
		&r.DeliveredAt,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	return r, err
}
