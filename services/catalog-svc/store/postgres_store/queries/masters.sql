-- masters.sql — Catalog master data CRUD queries.
--
-- Covers hotels, airlines, muthawwif, addons, and package_pricing
-- (SetDeparturePricing / GetDeparturePricing).
--
-- ID generation is done in the service layer (ULID with type prefix).
-- All amount columns are cast to BIGINT on read for wire conformance.

-- ===========================================================================
-- Hotels
-- ===========================================================================

-- name: InsertHotel :one
INSERT INTO catalog.hotels (
    id,
    name,
    city,
    star_rating,
    walking_distance_m
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING
    id, name, city, star_rating, walking_distance_m, created_at, updated_at;

-- name: UpdateHotelFields :one
UPDATE catalog.hotels
SET
    name              = COALESCE($1::text,     name),
    city              = COALESCE($2::text,     city),
    star_rating       = COALESCE($3::smallint, star_rating),
    walking_distance_m = COALESCE($4::integer, walking_distance_m),
    updated_at        = now()
WHERE id = $5::text
RETURNING
    id, name, city, star_rating, walking_distance_m, created_at, updated_at;

-- name: DeleteHotel :exec
DELETE FROM catalog.hotels WHERE id = $1;

-- name: GetHotelByID :one
SELECT id, name, city, star_rating, walking_distance_m, created_at, updated_at
FROM catalog.hotels
WHERE id = $1;

-- name: ListHotels :many
SELECT id, name, city, star_rating, walking_distance_m, created_at, updated_at
FROM catalog.hotels
WHERE ($1::text = '' OR id > $1::text)
ORDER BY id
LIMIT $2;

-- name: CountHotelPackageRefs :one
-- Check if a hotel is still referenced by any package before hard-delete.
SELECT COUNT(*)::bigint AS ref_count
FROM catalog.package_hotels
WHERE hotel_id = $1;

-- ===========================================================================
-- Airlines
-- ===========================================================================

-- name: InsertAirline :one
INSERT INTO catalog.airlines (
    id,
    code,
    name,
    operator_kind
) VALUES (
    $1,
    $2,
    $3,
    $4::catalog.operator_kind
)
RETURNING
    id, code, name, operator_kind, created_at, updated_at;

-- name: UpdateAirlineFields :one
UPDATE catalog.airlines
SET
    code         = COALESCE($1::text, code),
    name         = COALESCE($2::text, name),
    updated_at   = now()
WHERE id = $3::text
RETURNING
    id, code, name, operator_kind, created_at, updated_at;

-- name: DeleteAirline :exec
DELETE FROM catalog.airlines WHERE id = $1;

-- name: GetAirlineByIDForStaff :one
SELECT id, code, name, operator_kind, created_at, updated_at
FROM catalog.airlines
WHERE id = $1;

-- name: ListAirlines :many
SELECT id, code, name, operator_kind, created_at, updated_at
FROM catalog.airlines
WHERE ($1::text = '' OR id > $1::text)
ORDER BY id
LIMIT $2;

-- ===========================================================================
-- Muthawwif
-- ===========================================================================

-- name: InsertMuthawwif :one
INSERT INTO catalog.muthawwif (
    id,
    name,
    portrait_url
) VALUES (
    $1,
    $2,
    $3
)
RETURNING
    id, name, portrait_url, created_at, updated_at;

-- name: UpdateMuthawwifFields :one
UPDATE catalog.muthawwif
SET
    name         = COALESCE($1::text, name),
    portrait_url = COALESCE($2::text, portrait_url),
    updated_at   = now()
WHERE id = $3::text
RETURNING
    id, name, portrait_url, created_at, updated_at;

-- name: DeleteMuthawwif :exec
DELETE FROM catalog.muthawwif WHERE id = $1;

-- name: GetMuthawwifByIDForStaff :one
SELECT id, name, portrait_url, created_at, updated_at
FROM catalog.muthawwif
WHERE id = $1;

-- name: ListMuthawwif :many
SELECT id, name, portrait_url, created_at, updated_at
FROM catalog.muthawwif
WHERE ($1::text = '' OR id > $1::text)
ORDER BY id
LIMIT $2;

-- ===========================================================================
-- Addons
-- ===========================================================================

-- name: InsertAddon :one
INSERT INTO catalog.addons (
    id,
    name,
    list_amount,
    list_currency,
    settlement_currency
) VALUES (
    $1,
    $2,
    $3,
    'IDR',
    'IDR'
)
RETURNING
    id,
    name,
    list_amount::bigint AS list_amount,
    list_currency,
    settlement_currency,
    created_at,
    updated_at;

-- name: UpdateAddonFields :one
UPDATE catalog.addons
SET
    name        = COALESCE($1::text,    name),
    list_amount = COALESCE($2::numeric, list_amount),
    updated_at  = now()
WHERE id = $3::text
RETURNING
    id,
    name,
    list_amount::bigint AS list_amount,
    list_currency,
    settlement_currency,
    created_at,
    updated_at;

-- name: DeleteAddon :exec
DELETE FROM catalog.addons WHERE id = $1;

-- name: GetAddonByID :one
SELECT
    id,
    name,
    list_amount::bigint AS list_amount,
    list_currency,
    settlement_currency,
    created_at,
    updated_at
FROM catalog.addons
WHERE id = $1;

-- name: ListAddons :many
SELECT
    id,
    name,
    list_amount::bigint AS list_amount,
    list_currency,
    settlement_currency,
    created_at,
    updated_at
FROM catalog.addons
WHERE ($1::text = '' OR id > $1::text)
ORDER BY id
LIMIT $2;

-- ===========================================================================
-- Package Departure Pricing (SetDeparturePricing / GetDeparturePricing)
-- ===========================================================================

-- name: UpsertDeparturePricing :one
-- Insert or update a single pricing row for a (departure_id, room_type) pair.
-- Service layer calls this in a loop for each room type in the request.
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
    COALESCE(NULLIF($5, ''), 'IDR'),
    COALESCE(NULLIF($6, ''), 'IDR')
)
ON CONFLICT (package_departure_id, room_type) DO UPDATE
    SET
        list_amount         = EXCLUDED.list_amount,
        list_currency       = EXCLUDED.list_currency,
        settlement_currency = EXCLUDED.settlement_currency,
        updated_at          = now()
RETURNING
    id,
    package_departure_id,
    room_type,
    list_amount::bigint AS list_amount,
    list_currency,
    settlement_currency,
    created_at,
    updated_at;

-- name: GetDeparturePricingRows :many
SELECT
    id,
    package_departure_id,
    room_type,
    list_amount::bigint AS list_amount,
    list_currency,
    settlement_currency,
    created_at,
    updated_at
FROM catalog.package_pricing
WHERE package_departure_id = $1
ORDER BY room_type;
