-- S1-E-02 / BL-CAT-001 — catalog-svc read-model schema.
--
-- Scaffolds the 10 tables + 5 enums that back the public-read slice of F2
-- (Product Catalog & Master Data). The catalog data model is the set of
-- rows the three § Catalog endpoints in `docs/contracts/slice-S1.md`
-- surface to B2C / B2B / mobile consumers:
--
--   GET /v1/packages              — list (active only)
--   GET /v1/packages/{id}         — detail (eager master refs + departures)
--   GET /v1/package-departures/{id} — departure detail (future: BL-CAT-002)
--
-- All objects live in a dedicated `catalog` schema per ADR 0007's per-
-- service schema-namespace convention. sqlc reads this file via the shared
-- `schema: "../../../../migration"` pointer in each service's sqlc.yaml.
--
-- **ID convention** deviates from `docs/04-backend-conventions/04-database-
-- and-sqlc.md`: PKs are TEXT ULIDs with a type prefix (`pkg_`, `dep_`,
-- `htl_`, `arl_`, `mtw_`, `itn_`, `addon_`, `pkgpr_`) per the frozen
-- § Catalog "Response ID format" clause. The catalog contract treats IDs
-- as opaque strings; consumers must not parse them. The deviation is
-- deliberate — iam-svc sticks with UUID because its F1 contract uses
-- UUID; catalog's F2 contract chose ULID+prefix for time-sortable IDs
-- with an embedded type hint.
--
-- **Currency** columns (`list_currency`, `settlement_currency`) follow
-- Q001: list can be IDR or USD for display; settlement is always IDR in
-- MVP (CHECK enforced). Catalog never commits a payable amount — that
-- lock happens at payment-svc VA issuance with an FX snapshot.

CREATE SCHEMA IF NOT EXISTS catalog;

-- ---------------------------------------------------------------------------
-- Enum types
-- ---------------------------------------------------------------------------

-- Seven-value package kind per the § Catalog `kind` query-param enum.
CREATE TYPE catalog.package_kind AS ENUM (
    'umrah_reguler',
    'umrah_plus',
    'hajj_furoda',
    'hajj_khusus',
    'badal',
    'financial',
    'retail'
);

CREATE TYPE catalog.package_status AS ENUM (
    'draft',
    'active',
    'archived'
);

CREATE TYPE catalog.departure_status AS ENUM (
    'open',
    'closed',
    'departed',
    'completed',
    'cancelled'
);

CREATE TYPE catalog.room_type AS ENUM (
    'double',
    'triple',
    'quad'
);

-- operator_kind on airlines: PRD Section D notes Haramain high-speed rail
-- is modelled through the same table under a three-value enum.
CREATE TYPE catalog.operator_kind AS ENUM (
    'airline',
    'rail',
    'bus'
);

-- ---------------------------------------------------------------------------
-- itinerary_templates
-- ---------------------------------------------------------------------------
-- `days` is JSONB of shape [{ day int, title text, description text, photo_url text }, ...]
-- per the § Catalog detail response. public_url is the shareable micro-web
-- URL generated on package publish (F2 W8).

CREATE TABLE catalog.itinerary_templates (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL DEFAULT '',
    days       JSONB NOT NULL DEFAULT '[]'::jsonb,
    public_url TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT itinerary_templates_id_prefix CHECK (id LIKE 'itn_%')
);

-- ---------------------------------------------------------------------------
-- hotels
-- ---------------------------------------------------------------------------

CREATE TABLE catalog.hotels (
    id                   TEXT PRIMARY KEY,
    name                 TEXT NOT NULL,
    city                 TEXT NOT NULL,
    star_rating          SMALLINT NOT NULL DEFAULT 0,
    walking_distance_m   INTEGER NOT NULL DEFAULT 0,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT hotels_id_prefix CHECK (id LIKE 'htl_%'),
    CONSTRAINT hotels_star_rating_range CHECK (star_rating BETWEEN 0 AND 5)
);

-- ---------------------------------------------------------------------------
-- airlines
-- ---------------------------------------------------------------------------

CREATE TABLE catalog.airlines (
    id             TEXT PRIMARY KEY,
    code           TEXT NOT NULL,
    name           TEXT NOT NULL,
    operator_kind  catalog.operator_kind NOT NULL DEFAULT 'airline',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT airlines_id_prefix CHECK (id LIKE 'arl_%'),
    CONSTRAINT airlines_code_uk UNIQUE (code)
);

-- ---------------------------------------------------------------------------
-- muthawwif
-- ---------------------------------------------------------------------------

CREATE TABLE catalog.muthawwif (
    id            TEXT PRIMARY KEY,
    name          TEXT NOT NULL,
    portrait_url  TEXT NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT muthawwif_id_prefix CHECK (id LIKE 'mtw_%')
);

-- ---------------------------------------------------------------------------
-- addons
-- ---------------------------------------------------------------------------

CREATE TABLE catalog.addons (
    id                   TEXT PRIMARY KEY,
    name                 TEXT NOT NULL,
    list_amount          NUMERIC(18, 4) NOT NULL DEFAULT 0,
    list_currency        CHAR(3) NOT NULL DEFAULT 'IDR',
    settlement_currency  CHAR(3) NOT NULL DEFAULT 'IDR',
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT addons_id_prefix CHECK (id LIKE 'addon_%'),
    CONSTRAINT addons_list_amount_nonneg CHECK (list_amount >= 0),
    CONSTRAINT addons_settlement_currency_idr CHECK (settlement_currency = 'IDR'),
    CONSTRAINT addons_currency_iso CHECK (
        list_currency ~ '^[A-Z]{3}$' AND settlement_currency ~ '^[A-Z]{3}$'
    )
);

-- ---------------------------------------------------------------------------
-- packages
-- ---------------------------------------------------------------------------

CREATE TABLE catalog.packages (
    id               TEXT PRIMARY KEY,
    kind             catalog.package_kind NOT NULL,
    name             TEXT NOT NULL,
    description      TEXT NOT NULL DEFAULT '',
    highlights       TEXT[] NOT NULL DEFAULT '{}',
    cover_photo_url  TEXT NOT NULL DEFAULT '',
    itinerary_id     TEXT REFERENCES catalog.itinerary_templates (id) ON DELETE RESTRICT,
    airline_id       TEXT REFERENCES catalog.airlines (id) ON DELETE RESTRICT,
    muthawwif_id     TEXT REFERENCES catalog.muthawwif (id) ON DELETE RESTRICT,
    status           catalog.package_status NOT NULL DEFAULT 'draft',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at       TIMESTAMPTZ,
    CONSTRAINT packages_id_prefix CHECK (id LIKE 'pkg_%')
);

-- Hot paths: list filters by status (= 'active') + optional kind, and
-- cursor paginates by id. Partial indexes match the public read pattern.
CREATE INDEX packages_status_kind_idx
    ON catalog.packages (status, kind)
    WHERE deleted_at IS NULL;

CREATE INDEX packages_active_id_idx
    ON catalog.packages (id)
    WHERE status = 'active' AND deleted_at IS NULL;

CREATE INDEX packages_airline_id_idx
    ON catalog.packages (airline_id)
    WHERE deleted_at IS NULL;

-- ---------------------------------------------------------------------------
-- package_hotels (join)
-- ---------------------------------------------------------------------------

CREATE TABLE catalog.package_hotels (
    package_id TEXT NOT NULL REFERENCES catalog.packages (id) ON DELETE CASCADE,
    hotel_id   TEXT NOT NULL REFERENCES catalog.hotels (id) ON DELETE RESTRICT,
    sort_order SMALLINT NOT NULL DEFAULT 0,
    PRIMARY KEY (package_id, hotel_id)
);

CREATE INDEX package_hotels_hotel_id_idx ON catalog.package_hotels (hotel_id);

-- ---------------------------------------------------------------------------
-- package_addons (join)
-- ---------------------------------------------------------------------------

CREATE TABLE catalog.package_addons (
    package_id TEXT NOT NULL REFERENCES catalog.packages (id) ON DELETE CASCADE,
    addon_id   TEXT NOT NULL REFERENCES catalog.addons (id) ON DELETE RESTRICT,
    PRIMARY KEY (package_id, addon_id)
);

CREATE INDEX package_addons_addon_id_idx ON catalog.package_addons (addon_id);

-- ---------------------------------------------------------------------------
-- package_departures
-- ---------------------------------------------------------------------------
-- `reserved_seats` is the atomic counter driven by the gRPC
-- `ReserveSeats` / `ReleaseSeats` pair (S1-J-03 contract). Capacity guard
-- enforced at the application layer via a single-statement UPDATE ...
-- WHERE reserved_seats + $n <= total_seats. The CHECK here is a defense-
-- in-depth bound; the app-side atomic guard is the primary safety.

CREATE TABLE catalog.package_departures (
    id              TEXT PRIMARY KEY,
    package_id      TEXT NOT NULL REFERENCES catalog.packages (id) ON DELETE CASCADE,
    departure_date  DATE NOT NULL,
    return_date     DATE NOT NULL,
    total_seats     INTEGER NOT NULL,
    reserved_seats  INTEGER NOT NULL DEFAULT 0,
    status          catalog.departure_status NOT NULL DEFAULT 'open',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT package_departures_id_prefix CHECK (id LIKE 'dep_%'),
    CONSTRAINT package_departures_total_seats_pos CHECK (total_seats > 0),
    CONSTRAINT package_departures_reserved_seats_nonneg CHECK (reserved_seats >= 0),
    CONSTRAINT package_departures_capacity CHECK (reserved_seats <= total_seats),
    CONSTRAINT package_departures_date_order CHECK (return_date >= departure_date)
);

CREATE INDEX package_departures_package_id_date_idx
    ON catalog.package_departures (package_id, departure_date);

CREATE INDEX package_departures_status_date_idx
    ON catalog.package_departures (status, departure_date);

-- ---------------------------------------------------------------------------
-- package_pricing
-- ---------------------------------------------------------------------------
-- One row per (departure, room_type). `starting_price` in the list response
-- is min(list_amount) across the departure's pricing rows for the cheapest
-- `next_departure` of each package.

CREATE TABLE catalog.package_pricing (
    id                   TEXT PRIMARY KEY,
    package_departure_id TEXT NOT NULL REFERENCES catalog.package_departures (id) ON DELETE CASCADE,
    room_type            catalog.room_type NOT NULL,
    list_amount          NUMERIC(18, 4) NOT NULL,
    list_currency        CHAR(3) NOT NULL DEFAULT 'IDR',
    settlement_currency  CHAR(3) NOT NULL DEFAULT 'IDR',
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT package_pricing_id_prefix CHECK (id LIKE 'pkgpr_%'),
    CONSTRAINT package_pricing_list_amount_nonneg CHECK (list_amount >= 0),
    CONSTRAINT package_pricing_settlement_currency_idr CHECK (settlement_currency = 'IDR'),
    CONSTRAINT package_pricing_currency_iso CHECK (
        list_currency ~ '^[A-Z]{3}$' AND settlement_currency ~ '^[A-Z]{3}$'
    ),
    CONSTRAINT package_pricing_departure_room_uk UNIQUE (package_departure_id, room_type)
);

CREATE INDEX package_pricing_departure_id_idx ON catalog.package_pricing (package_departure_id);
