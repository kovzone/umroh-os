-- S1-E-02 / BL-CAT-001 — package read queries.
--
-- All queries are read-only. Only `status = 'active'` rows are returned
-- by every query; draft/archived packages must never leak through these
-- endpoints per `slice-S1.md § Catalog`.
--
-- The list page is built by three queries: `ListActivePackages` returns
-- the base package rows (filtered + cursor-paginated), and for each row
-- the service layer fetches `GetStartingPriceForPackage` and
-- `GetNextDepartureForPackage` to hydrate `starting_price` and
-- `next_departure`. A package with no upcoming open/closed departure
-- returns `pgx.ErrNoRows` from the latter two — the handler emits
-- `next_departure: null` and `starting_price: {0, IDR, IDR}` (an
-- intentional "no priced departure" signal for the UI).
--
-- Cursor pagination keys off `packages.id` ULID (time-sortable), so
-- `WHERE id > $cursor ORDER BY id` yields a stable chronological page
-- across the active set.

-- name: ListActivePackages :many
-- Filters are all optional; pass NULL / empty string to skip a filter.
-- Caller fetches `row_limit+1` rows; the extra row signals
-- `has_more=true` and is dropped before shaping the response.
SELECT
    p.id,
    p.kind,
    p.name,
    p.description,
    p.cover_photo_url
FROM catalog.packages p
WHERE p.status = 'active'
  AND p.deleted_at IS NULL
  AND (sqlc.narg('kind')::catalog.package_kind IS NULL OR p.kind = sqlc.narg('kind')::catalog.package_kind)
  AND (sqlc.narg('airline_code')::text IS NULL OR EXISTS (
        SELECT 1 FROM catalog.airlines a
        WHERE a.id = p.airline_id AND a.code = sqlc.narg('airline_code')::text
  ))
  AND (sqlc.narg('hotel_id')::text IS NULL OR EXISTS (
        SELECT 1 FROM catalog.package_hotels ph
        WHERE ph.package_id = p.id AND ph.hotel_id = sqlc.narg('hotel_id')::text
  ))
  AND (
        sqlc.narg('departure_from')::date IS NULL
        AND sqlc.narg('departure_to')::date IS NULL
        OR EXISTS (
            SELECT 1 FROM catalog.package_departures d
            WHERE d.package_id = p.id
              AND d.status IN ('open', 'closed')
              AND (sqlc.narg('departure_from')::date IS NULL OR d.departure_date >= sqlc.narg('departure_from')::date)
              AND (sqlc.narg('departure_to')::date IS NULL OR d.departure_date <= sqlc.narg('departure_to')::date)
        )
  )
  AND (sqlc.narg('cursor_id')::text IS NULL OR p.id > sqlc.narg('cursor_id')::text)
ORDER BY p.id ASC
LIMIT sqlc.arg('row_limit')::integer;

-- name: GetActivePackageByID :one
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
  AND status = 'active'
  AND deleted_at IS NULL;

-- name: GetNextDepartureForPackage :one
-- Returns pgx.ErrNoRows when the package has no upcoming open/closed
-- departure. Service layer treats that as "no next_departure" and
-- emits `null` on the wire.
SELECT
    id,
    departure_date,
    return_date,
    total_seats,
    reserved_seats,
    (total_seats - reserved_seats)::integer AS remaining_seats
FROM catalog.package_departures
WHERE package_id = $1
  AND status IN ('open', 'closed')
  AND departure_date >= CURRENT_DATE
ORDER BY departure_date ASC
LIMIT 1;

-- name: GetStartingPriceForPackage :one
-- Returns the lowest priced room type across the package's upcoming
-- open/closed departures. Returns pgx.ErrNoRows if the package has
-- no priced upcoming departure (unusual for an active package).
--
-- `list_amount` is cast to BIGINT for wire conformance (whole currency
-- units per § Catalog). MVP single-currency (Q001) means no sub-unit
-- loss; multi-currency handling revisits this when it lands.
SELECT
    pr.list_amount::bigint AS list_amount,
    pr.list_currency,
    pr.settlement_currency
FROM catalog.package_pricing pr
JOIN catalog.package_departures d ON d.id = pr.package_departure_id
WHERE d.package_id = $1
  AND d.status IN ('open', 'closed')
  AND d.departure_date >= CURRENT_DATE
ORDER BY pr.list_amount ASC
LIMIT 1;

-- name: GetItineraryByID :one
SELECT id, name, days, public_url, created_at, updated_at
FROM catalog.itinerary_templates
WHERE id = $1;

-- name: GetAirlineByID :one
SELECT id, code, name, operator_kind, created_at, updated_at
FROM catalog.airlines
WHERE id = $1;

-- name: GetMuthawwifByID :one
SELECT id, name, portrait_url, created_at, updated_at
FROM catalog.muthawwif
WHERE id = $1;

-- name: ListHotelsForPackage :many
SELECT h.id, h.name, h.city, h.star_rating, h.walking_distance_m
FROM catalog.hotels h
JOIN catalog.package_hotels ph ON ph.hotel_id = h.id
WHERE ph.package_id = $1
ORDER BY ph.sort_order ASC, h.name ASC;

-- name: ListAddonsForPackage :many
SELECT
    a.id,
    a.name,
    a.list_amount::bigint AS list_amount,
    a.list_currency,
    a.settlement_currency
FROM catalog.addons a
JOIN catalog.package_addons pa ON pa.addon_id = a.id
WHERE pa.package_id = $1
ORDER BY a.name ASC;

-- name: ListOpenDeparturesForPackage :many
-- Public-readable departures: open/closed only, upcoming only.
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
WHERE package_id = $1
  AND status IN ('open', 'closed')
  AND departure_date >= CURRENT_DATE
ORDER BY departure_date ASC;

-- name: GetActiveDeparture :one
-- Returns the departure row only if its status is publicly visible
-- (open or closed). Any other status returns pgx.ErrNoRows which the
-- service layer maps to apperrors.ErrNotFound → 404 departure_not_found.
-- Identical 404 shape for unknown-id vs hidden-status — no existence oracle.
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
WHERE id = $1
  AND status IN ('open', 'closed');

-- name: ListPricingForDeparture :many
-- Returns all price rows for a departure ordered by list_amount ascending
-- so the cheapest room type surfaces first. `list_amount` cast to bigint
-- for wire-integer conformance (whole currency units per § Catalog + Q001).
SELECT
    room_type,
    list_amount::bigint AS list_amount,
    list_currency,
    settlement_currency
FROM catalog.package_pricing
WHERE package_departure_id = $1
ORDER BY list_amount ASC;
