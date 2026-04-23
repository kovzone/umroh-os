-- S1-E-07 / BL-CAT-014 — catalog-svc staff write queries.
--
-- All mutating queries run through the service layer which gates calls via
-- iam-svc.CheckPermission (resource="catalog_package", action="manage",
-- scope="global") before reaching these store methods.
--
-- ID generation: service layer is responsible for minting ULIDs with the
-- correct prefix (pkg_, dep_, pkgpr_) using a ULID library and passing them
-- as $id parameters. The database CHECK constraints enforce prefix correctness
-- as a defense-in-depth guard.

-- name: InsertPackage :one
-- Insert a new package row. Highlights is a Go []string → TEXT[].
INSERT INTO catalog.packages (
    id,
    kind,
    name,
    description,
    highlights,
    cover_photo_url,
    itinerary_id,
    airline_id,
    muthawwif_id,
    status
) VALUES (
    $1,
    $2::catalog.package_kind,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10::catalog.package_status
)
RETURNING
    id,
    kind,
    name,
    description,
    highlights,
    cover_photo_url,
    itinerary_id,
    airline_id,
    muthawwif_id,
    status,
    created_at,
    updated_at;

-- name: UpdatePackageFields :one
-- Partial-update a package row. Only sets columns whose corresponding
-- narg is non-NULL. The status transition validity (e.g. draft→active) is
-- enforced in the service layer before this query runs.
UPDATE catalog.packages
SET
    name            = COALESCE(sqlc.narg('name')::text,            name),
    description     = COALESCE(sqlc.narg('description')::text,     description),
    cover_photo_url = COALESCE(sqlc.narg('cover_photo_url')::text, cover_photo_url),
    highlights      = COALESCE(sqlc.narg('highlights')::text[],    highlights),
    itinerary_id    = COALESCE(sqlc.narg('itinerary_id')::text,    itinerary_id),
    airline_id      = COALESCE(sqlc.narg('airline_id')::text,      airline_id),
    muthawwif_id    = COALESCE(sqlc.narg('muthawwif_id')::text,    muthawwif_id),
    status          = COALESCE(sqlc.narg('status')::catalog.package_status, status),
    updated_at      = now()
WHERE id = sqlc.arg('id')::text
  AND deleted_at IS NULL
RETURNING
    id,
    kind,
    name,
    description,
    highlights,
    cover_photo_url,
    itinerary_id,
    airline_id,
    muthawwif_id,
    status,
    created_at,
    updated_at;

-- name: SoftDeletePackage :one
-- Soft-delete: set deleted_at + status=archived atomically.
UPDATE catalog.packages
SET
    deleted_at = now(),
    status     = 'archived',
    updated_at = now()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, status, deleted_at;

-- name: GetPackageByIDForStaff :one
-- Staff-side read: returns a package regardless of status (draft/active/archived)
-- as long as it is not hard-deleted. Used by the service layer after a write to
-- hydrate the response (before lazy-loading relations).
SELECT
    id,
    kind,
    name,
    description,
    highlights,
    cover_photo_url,
    itinerary_id,
    airline_id,
    muthawwif_id,
    status,
    created_at,
    updated_at
FROM catalog.packages
WHERE id = $1
  AND deleted_at IS NULL;

-- name: DeletePackageHotels :exec
-- Remove all hotel associations for a package before replacing them.
DELETE FROM catalog.package_hotels
WHERE package_id = $1;

-- name: InsertPackageHotel :exec
-- Insert one hotel association. Called in a loop by the service layer.
INSERT INTO catalog.package_hotels (package_id, hotel_id, sort_order)
VALUES ($1, $2, $3)
ON CONFLICT (package_id, hotel_id) DO UPDATE SET sort_order = EXCLUDED.sort_order;

-- name: DeletePackageAddons :exec
-- Remove all add-on associations for a package before replacing them.
DELETE FROM catalog.package_addons
WHERE package_id = $1;

-- name: InsertPackageAddon :exec
-- Insert one add-on association.
INSERT INTO catalog.package_addons (package_id, addon_id)
VALUES ($1, $2)
ON CONFLICT (package_id, addon_id) DO NOTHING;

-- name: InsertDeparture :one
-- Insert a new departure row.
INSERT INTO catalog.package_departures (
    id,
    package_id,
    departure_date,
    return_date,
    total_seats,
    status
) VALUES (
    $1,
    $2,
    $3::date,
    $4::date,
    $5,
    $6::catalog.departure_status
)
RETURNING
    id,
    package_id,
    departure_date,
    return_date,
    total_seats,
    reserved_seats,
    (total_seats - reserved_seats)::integer AS remaining_seats,
    status,
    created_at,
    updated_at;

-- name: UpdateDepartureFields :one
-- Partial-update a departure row. Status transition validity enforced in service layer.
UPDATE catalog.package_departures
SET
    departure_date = COALESCE(sqlc.narg('departure_date')::date,              departure_date),
    return_date    = COALESCE(sqlc.narg('return_date')::date,                 return_date),
    total_seats    = COALESCE(sqlc.narg('total_seats')::integer,              total_seats),
    status         = COALESCE(sqlc.narg('status')::catalog.departure_status,  status),
    updated_at     = now()
WHERE id = sqlc.arg('id')::text
RETURNING
    id,
    package_id,
    departure_date,
    return_date,
    total_seats,
    reserved_seats,
    (total_seats - reserved_seats)::integer AS remaining_seats,
    status,
    created_at,
    updated_at;

-- name: GetDepartureByIDForStaff :one
-- Staff-side departure read: returns any departure regardless of status.
SELECT
    id,
    package_id,
    departure_date,
    return_date,
    total_seats,
    reserved_seats,
    (total_seats - reserved_seats)::integer AS remaining_seats,
    status
FROM catalog.package_departures
WHERE id = $1;

-- name: DeleteDeparturePricing :exec
-- Remove all pricing rows for a departure before replacing them.
DELETE FROM catalog.package_pricing
WHERE package_departure_id = $1;

-- name: InsertDeparturePricing :one
-- Insert one pricing row for a departure.
INSERT INTO catalog.package_pricing (
    id,
    package_departure_id,
    room_type,
    list_amount,
    list_currency,
    settlement_currency
) VALUES (
    $1,
    $2,
    $3::catalog.room_type,
    $4,
    $5,
    $6
)
RETURNING
    id,
    package_departure_id,
    room_type,
    list_amount::bigint AS list_amount,
    list_currency,
    settlement_currency;

-- ---------------------------------------------------------------------------
-- ReserveSeats / ReleaseSeats (§ Inventory, S1-J-03)
-- ---------------------------------------------------------------------------

-- name: ReserveSeatsAtomic :one
-- Atomic capacity decrement per § Inventory contract.
-- Returns zero rows (pgx.ErrNoRows) when reserved_seats + $seats > total_seats
-- at commit time → service layer maps to FAILED_PRECONDITION insufficient_capacity.
UPDATE catalog.package_departures
SET
    reserved_seats = reserved_seats + sqlc.arg('seats')::integer,
    updated_at     = now()
WHERE id = sqlc.arg('departure_id')::text
  AND status IN ('open', 'closed')
  AND reserved_seats + sqlc.arg('seats')::integer <= total_seats
RETURNING
    id,
    reserved_seats,
    total_seats,
    (total_seats - reserved_seats)::integer AS remaining_seats;

-- name: InsertSeatReservation :one
-- Write the dedup row for a ReserveSeats call in the same transaction as the
-- capacity decrement. TTL is derived in Go: expires_at = now() + interval.
INSERT INTO catalog.seat_reservations (
    reservation_id,
    departure_id,
    seats,
    expires_at
) VALUES ($1, $2, $3, $4)
RETURNING reservation_id, departure_id, seats, reserved_at, expires_at;

-- name: GetSeatReservation :one
-- Lookup a reservation for idempotency check.
SELECT reservation_id, departure_id, seats, reserved_at, expires_at, released_at
FROM catalog.seat_reservations
WHERE reservation_id = $1;

-- name: ReleaseSeatsAtomic :one
-- Atomic capacity increment. Returns the updated row.
UPDATE catalog.package_departures
SET
    reserved_seats = reserved_seats - sqlc.arg('seats_to_release')::integer,
    updated_at     = now()
WHERE id = sqlc.arg('departure_id')::text
  AND reserved_seats - sqlc.arg('seats_to_release')::integer >= 0
RETURNING
    id,
    reserved_seats,
    total_seats,
    (total_seats - reserved_seats)::integer AS remaining_seats;

-- name: MarkReservationReleased :one
-- Mark the dedup row as released so subsequent ReleaseSeats calls are no-ops.
UPDATE catalog.seat_reservations
SET released_at = now()
WHERE reservation_id = $1
  AND released_at IS NULL
RETURNING reservation_id, departure_id, seats, reserved_at, expires_at, released_at;
