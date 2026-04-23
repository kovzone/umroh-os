-- fulfillment_tasks.sql — queries for logistics.fulfillment_tasks table.
-- S3-E-02 / BL-LOG-001.

-- name: GetFulfillmentTaskByBookingID :one
-- Returns an existing fulfillment task for the given booking, or no rows.
SELECT id, booking_id, departure_id, status, tracking_number,
       shipped_at, delivered_at, created_at, updated_at
FROM logistics.fulfillment_tasks
WHERE booking_id = $1
LIMIT 1;

-- name: InsertFulfillmentTask :one
-- Inserts a new fulfillment task in status='queued'.
-- The UNIQUE constraint on booking_id ensures idempotency; callers should
-- check for existing tasks before calling this query.
INSERT INTO logistics.fulfillment_tasks (
    booking_id,
    departure_id,
    status
) VALUES (
    $1, -- booking_id   UUID
    $2, -- departure_id UUID
    'queued'
)
RETURNING id, booking_id, departure_id, status, tracking_number,
          shipped_at, delivered_at, created_at, updated_at;
