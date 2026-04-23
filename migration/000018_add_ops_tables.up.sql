CREATE SCHEMA IF NOT EXISTS ops;

CREATE TABLE ops.room_allocations (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    departure_id     UUID NOT NULL,
    status           TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','committed','cancelled')),
    committed_at     TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT room_allocations_departure_id_unique UNIQUE (departure_id)
);

CREATE TABLE ops.room_assignments (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    allocation_id    UUID NOT NULL REFERENCES ops.room_allocations(id),
    room_number      TEXT NOT NULL,
    jamaah_id        UUID NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT room_assignments_allocation_jamaah_unique UNIQUE (allocation_id, jamaah_id)
);

CREATE TABLE ops.id_card_issuances (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    jamaah_id        UUID NOT NULL,
    departure_id     UUID NOT NULL,
    card_type        TEXT NOT NULL CHECK (card_type IN ('id_card','luggage_tag')),
    token            TEXT NOT NULL UNIQUE,
    issued_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT id_card_issuances_jamaah_type_dep_unique UNIQUE (jamaah_id, departure_id, card_type)
);

CREATE INDEX ON ops.room_assignments (allocation_id);
CREATE INDEX ON ops.id_card_issuances (jamaah_id, departure_id);
