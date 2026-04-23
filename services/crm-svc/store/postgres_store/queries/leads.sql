-- S4-E-02 — crm-svc lead queries.
--
-- All mutations use sqlc's :one / :exec return modes.
-- Status transitions and business logic are enforced in the service layer;
-- the store layer is intentionally dumb (pure CRUD + filter).

-- name: InsertLead :one
INSERT INTO crm.leads (
    source,
    utm_source,
    utm_medium,
    utm_campaign,
    utm_content,
    utm_term,
    name,
    phone,
    email,
    interest_package_id,
    interest_departure_id,
    status,
    assigned_cs_id,
    notes
) VALUES (
    $1,  -- source
    $2,  -- utm_source
    $3,  -- utm_medium
    $4,  -- utm_campaign
    $5,  -- utm_content
    $6,  -- utm_term
    $7,  -- name
    $8,  -- phone
    $9,  -- email
    $10, -- interest_package_id
    $11, -- interest_departure_id
    $12, -- status
    $13, -- assigned_cs_id
    $14  -- notes
)
RETURNING
    id, source, utm_source, utm_medium, utm_campaign, utm_content, utm_term,
    name, phone, email, interest_package_id, interest_departure_id,
    status, assigned_cs_id, notes, booking_id, created_at, updated_at;

-- name: GetLeadByID :one
SELECT
    id, source, utm_source, utm_medium, utm_campaign, utm_content, utm_term,
    name, phone, email, interest_package_id, interest_departure_id,
    status, assigned_cs_id, notes, booking_id, created_at, updated_at
FROM crm.leads
WHERE id = $1;

-- name: GetLeadByBookingID :one
-- Used for idempotency checks on booking-event callbacks.
SELECT
    id, source, utm_source, utm_medium, utm_campaign, utm_content, utm_term,
    name, phone, email, interest_package_id, interest_departure_id,
    status, assigned_cs_id, notes, booking_id, created_at, updated_at
FROM crm.leads
WHERE booking_id = $1
LIMIT 1;

-- name: UpdateLeadStatus :one
-- Partial update: only touches status, notes, assigned_cs_id, booking_id.
-- Caller sets only the fields they want changed; others pass current values.
UPDATE crm.leads
SET
    status         = COALESCE($2, status),
    notes          = COALESCE($3, notes),
    assigned_cs_id = COALESCE($4, assigned_cs_id),
    booking_id     = COALESCE($5, booking_id),
    updated_at     = now()
WHERE id = $1
RETURNING
    id, source, utm_source, utm_medium, utm_campaign, utm_content, utm_term,
    name, phone, email, interest_package_id, interest_departure_id,
    status, assigned_cs_id, notes, booking_id, created_at, updated_at;

-- name: ListLeads :many
-- Filter by optional status + assigned_cs_id, sorted newest-first.
-- Pagination via LIMIT + OFFSET.
SELECT
    id, source, utm_source, utm_medium, utm_campaign, utm_content, utm_term,
    name, phone, email, interest_package_id, interest_departure_id,
    status, assigned_cs_id, notes, booking_id, created_at, updated_at
FROM crm.leads
WHERE
    ($1::text IS NULL OR status = $1)
    AND ($2::uuid IS NULL OR assigned_cs_id = $2)
ORDER BY created_at DESC
LIMIT $3
OFFSET $4;

-- name: CountLeads :one
SELECT COUNT(*) FROM crm.leads
WHERE
    ($1::text IS NULL OR status = $1)
    AND ($2::uuid IS NULL OR assigned_cs_id = $2);

-- name: GetLeastLoadedCS :one
-- Round-robin CS assignment: pick the assigned_cs_id with the fewest active
-- (non-converted, non-lost) leads. Ties broken by assigned_cs_id ASC for
-- determinism. Returns NULL if no CS has active leads yet.
SELECT assigned_cs_id
FROM crm.leads
WHERE assigned_cs_id IS NOT NULL
  AND status NOT IN ('converted', 'lost')
GROUP BY assigned_cs_id
ORDER BY COUNT(*) ASC, assigned_cs_id ASC
LIMIT 1;
