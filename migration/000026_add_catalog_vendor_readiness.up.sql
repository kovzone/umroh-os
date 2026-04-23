-- 000026_add_catalog_vendor_readiness.up.sql
--
-- BL-OPS-020 — Vendor readiness checklist.
--
-- Tracks ticket / hotel / visa readiness state per departure.
-- One row per (departure_id, kind); state transitions drive the
-- departure readiness badge in the ops console.

CREATE TABLE IF NOT EXISTS catalog.departure_vendor_readiness (
    id             text        NOT NULL,
    departure_id   text        NOT NULL
                               REFERENCES catalog.package_departures(id)
                               ON DELETE CASCADE,
    kind           text        NOT NULL,
    state          text        NOT NULL DEFAULT 'not_started',
    notes          text        NOT NULL DEFAULT '',
    attachment_url text        NOT NULL DEFAULT '',
    updated_at     timestamptz NOT NULL DEFAULT now(),
    updated_by     text        NOT NULL DEFAULT '',

    CONSTRAINT departure_vendor_readiness_pkey
        PRIMARY KEY (id),
    CONSTRAINT departure_vendor_readiness_kind_check
        CHECK (kind IN ('ticket', 'hotel', 'visa')),
    CONSTRAINT departure_vendor_readiness_state_check
        CHECK (state IN ('not_started', 'in_progress', 'done')),
    CONSTRAINT departure_vendor_readiness_unique
        UNIQUE (departure_id, kind)
);

CREATE INDEX IF NOT EXISTS departure_vendor_readiness_departure_id_idx
    ON catalog.departure_vendor_readiness (departure_id);
