-- 000011 — Seat reservation dedup table (§ Inventory / S1-J-03) +
--          Booking draft tables (§ Booking / S1-J-02 / S1-E-03).
--
-- seat_reservations: tracks atomic ReserveSeats/ReleaseSeats dedup rows so
-- that idempotent replays don't double-decrement capacity. Lives in the
-- catalog schema because it is catalog-svc's data.
--
-- bookings + booking_items + booking_addons: the booking draft data model
-- that backs POST /v1/bookings (S1-E-03). Lives in a dedicated `booking`
-- schema per ADR 0007 per-service namespacing.
--
-- ID convention: bookings use ULID+prefix (bkg_, bkgitem_) matching
-- slice-S1.md § Booking "ULID IDs with type prefixes".

-- ---------------------------------------------------------------------------
-- catalog.seat_reservations (dedup table for ReserveSeats idempotency)
-- ---------------------------------------------------------------------------

CREATE TABLE catalog.seat_reservations (
    reservation_id  TEXT        PRIMARY KEY,   -- caller-supplied ULID
    departure_id    TEXT        NOT NULL REFERENCES catalog.package_departures (id) ON DELETE CASCADE,
    seats           INTEGER     NOT NULL CHECK (seats > 0),
    reserved_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at      TIMESTAMPTZ NOT NULL,
    released_at     TIMESTAMPTZ,               -- NULL = still active

    CONSTRAINT seat_reservations_expires_after_reserved CHECK (expires_at > reserved_at)
);

CREATE INDEX seat_reservations_departure_id_idx ON catalog.seat_reservations (departure_id);
CREATE INDEX seat_reservations_expires_at_idx   ON catalog.seat_reservations (expires_at);

-- ---------------------------------------------------------------------------
-- booking schema
-- ---------------------------------------------------------------------------

CREATE SCHEMA IF NOT EXISTS booking;

-- booking_channel: matches slice-S1.md § Booking auth table.
CREATE TYPE booking.channel AS ENUM ('b2c_self', 'b2b_agent', 'cs');

-- booking_status: full state machine from § Booking States.
-- S1-E-03 only ever writes 'draft'; other values reserved for later slices.
CREATE TYPE booking.status AS ENUM (
    'draft',
    'pending_payment',
    'partially_paid',
    'paid_in_full',
    'departed',
    'completed',
    'expired',
    'cancelled',
    'failed'
);

-- ---------------------------------------------------------------------------
-- bookings
-- ---------------------------------------------------------------------------

CREATE TABLE booking.bookings (
    id               TEXT               PRIMARY KEY,
    status           booking.status     NOT NULL DEFAULT 'draft',
    channel          booking.channel    NOT NULL,

    -- catalog references (snapshotted IDs; not FK because booking-svc has no
    -- foreign-key access to catalog schema; integrity enforced at app layer).
    package_id       TEXT               NOT NULL,
    departure_id     TEXT               NOT NULL,
    room_type        TEXT               NOT NULL, -- catalog.room_type value

    -- channel attribution
    agent_id         TEXT,              -- populated when channel = b2b_agent
    staff_user_id    TEXT,              -- populated when channel = cs; UUID from IAM

    -- lead contact (denormalized for quick lookup)
    lead_full_name   TEXT               NOT NULL,
    lead_email       TEXT,
    lead_whatsapp    TEXT               NOT NULL,
    lead_domicile    TEXT               NOT NULL,

    -- pricing snapshot at draft time (list price only; per Q001)
    list_amount      BIGINT             NOT NULL DEFAULT 0,
    list_currency    CHAR(3)            NOT NULL DEFAULT 'IDR',
    settlement_currency CHAR(3)         NOT NULL DEFAULT 'IDR',

    notes            TEXT,

    -- idempotency support
    idempotency_key  TEXT,
    idempotency_body_hash TEXT,         -- SHA-256 hex of the canonical request body

    -- timestamps
    created_at       TIMESTAMPTZ        NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ        NOT NULL DEFAULT now(),
    expires_at       TIMESTAMPTZ,       -- non-binding draft expiry hint (30 min post-create)

    CONSTRAINT bookings_id_prefix CHECK (id LIKE 'bkg_%'),
    CONSTRAINT bookings_settlement_currency_idr CHECK (settlement_currency = 'IDR'),
    CONSTRAINT bookings_agent_id_required CHECK (
        channel != 'b2b_agent' OR agent_id IS NOT NULL
    ),
    CONSTRAINT bookings_staff_user_id_required CHECK (
        channel != 'cs' OR staff_user_id IS NOT NULL
    )
);

CREATE INDEX bookings_status_created_idx ON booking.bookings (status, created_at DESC);
CREATE INDEX bookings_departure_id_idx   ON booking.bookings (departure_id);
CREATE INDEX bookings_agent_id_idx       ON booking.bookings (agent_id) WHERE agent_id IS NOT NULL;
CREATE UNIQUE INDEX bookings_idempotency_key_channel_idx
    ON booking.bookings (channel, idempotency_key)
    WHERE idempotency_key IS NOT NULL;

-- ---------------------------------------------------------------------------
-- booking_items (one row per jamaah / pilgrim)
-- ---------------------------------------------------------------------------

CREATE TYPE booking.item_status AS ENUM ('active', 'cancelled');

CREATE TABLE booking.booking_items (
    id           TEXT                 PRIMARY KEY,
    booking_id   TEXT                 NOT NULL REFERENCES booking.bookings (id) ON DELETE CASCADE,
    full_name    TEXT                 NOT NULL,
    email        TEXT,
    whatsapp     TEXT,
    domicile     TEXT                 NOT NULL,
    is_lead      BOOLEAN              NOT NULL DEFAULT FALSE,
    status       booking.item_status  NOT NULL DEFAULT 'active',
    created_at   TIMESTAMPTZ          NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ          NOT NULL DEFAULT now(),

    CONSTRAINT booking_items_id_prefix CHECK (id LIKE 'bkgitem_%')
);

CREATE INDEX booking_items_booking_id_idx ON booking.booking_items (booking_id);

-- ---------------------------------------------------------------------------
-- booking_addons (add-ons linked to a booking)
-- ---------------------------------------------------------------------------

CREATE TABLE booking.booking_addons (
    booking_id  TEXT    NOT NULL REFERENCES booking.bookings (id) ON DELETE CASCADE,
    addon_id    TEXT    NOT NULL,   -- catalog.addons.id value; no FK (cross-schema)
    addon_name  TEXT    NOT NULL,   -- snapshotted at booking time
    list_amount BIGINT  NOT NULL DEFAULT 0,
    list_currency CHAR(3) NOT NULL DEFAULT 'IDR',
    settlement_currency CHAR(3) NOT NULL DEFAULT 'IDR',

    PRIMARY KEY (booking_id, addon_id),
    CONSTRAINT booking_addons_settlement_currency_idr CHECK (settlement_currency = 'IDR')
);
