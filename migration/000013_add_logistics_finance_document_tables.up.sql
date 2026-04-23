-- 000013 — Logistics fulfillment tasks, Finance journal engine, Document upload (S3-E-02 / S3-E-03).
--
-- Introduces three schema additions:
--   logistics.fulfillment_tasks  — created by OnBookingPaid (S3-E-02)
--   finance.journal_entries      — double-entry journal header (S3-E-03)
--   finance.journal_lines        — per-account debit/credit lines (S3-E-03)
--   jamaah.pilgrim_documents     — document upload stub (S3-E-02 partial / F3-W1)
--
-- All IDs are UUID (gen_random_uuid()).  Status columns use TEXT with CHECK
-- constraints rather than ENUM types so new values can be added without
-- blocking migrations (no ALTER TYPE required).
--
-- Schemas: logistics, finance, jamaah (per ADR-0007 per-service namespacing).

-- ---------------------------------------------------------------------------
-- Schemas
-- ---------------------------------------------------------------------------

CREATE SCHEMA IF NOT EXISTS logistics;
CREATE SCHEMA IF NOT EXISTS finance;
CREATE SCHEMA IF NOT EXISTS jamaah;

-- ---------------------------------------------------------------------------
-- logistics.fulfillment_tasks
-- ---------------------------------------------------------------------------
-- Created by logistics-svc.OnBookingPaid.  One row per booking; idempotent
-- on booking_id (unique constraint).  Status flows:
--   queued → processing → shipped → delivered
--   queued → cancelled  (cancellation path)

CREATE TABLE logistics.fulfillment_tasks (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id       UUID        NOT NULL,
    departure_id     UUID        NOT NULL,
    status           TEXT        NOT NULL DEFAULT 'queued'
        CHECK (status IN ('queued', 'processing', 'shipped', 'delivered', 'cancelled')),
    tracking_number  TEXT,
    shipped_at       TIMESTAMPTZ,
    delivered_at     TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Idempotency: one active task per booking.
    CONSTRAINT fulfillment_tasks_booking_id_unique UNIQUE (booking_id)
);

CREATE INDEX ON logistics.fulfillment_tasks (booking_id);
CREATE INDEX ON logistics.fulfillment_tasks (departure_id);
CREATE INDEX ON logistics.fulfillment_tasks (status);

-- ---------------------------------------------------------------------------
-- finance.journal_entries
-- ---------------------------------------------------------------------------
-- Header record for each double-entry journal posting.
-- idempotency_key enforces one posting per source event (e.g. "payment:<invoice_id>").

CREATE TABLE finance.journal_entries (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key  TEXT        NOT NULL UNIQUE,
    source_type      TEXT        NOT NULL
        CHECK (source_type IN ('payment', 'departure', 'refund')),
    source_id        UUID        NOT NULL,
    posted_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    description      TEXT
);

CREATE INDEX ON finance.journal_entries (source_type, source_id);

-- ---------------------------------------------------------------------------
-- finance.journal_lines
-- ---------------------------------------------------------------------------
-- Each line represents one leg of a double-entry posting.
-- Constraint: exactly one of debit/credit is positive; the other must be zero.
-- Balance check (Σ debit = Σ credit per entry) is enforced in the service layer
-- before INSERT rather than via a DB trigger, consistent with the pattern used
-- across the codebase.

CREATE TABLE finance.journal_lines (
    id           UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_id     UUID            NOT NULL REFERENCES finance.journal_entries (id),
    account_code TEXT            NOT NULL,
    debit        NUMERIC(15, 2)  NOT NULL DEFAULT 0,
    credit       NUMERIC(15, 2)  NOT NULL DEFAULT 0,
    CONSTRAINT journal_lines_debit_or_credit CHECK (
        (debit > 0 AND credit = 0) OR (debit = 0 AND credit > 0)
    )
);

CREATE INDEX ON finance.journal_lines (entry_id);

-- ---------------------------------------------------------------------------
-- jamaah.pilgrim_documents
-- ---------------------------------------------------------------------------
-- Scaffold table for F3-W1 document upload (S3-E-02 partial).
-- Full OCR pipeline and verification queue (BL-DOC-001..003) build on top.

CREATE TABLE jamaah.pilgrim_documents (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    jamaah_id        UUID        NOT NULL,
    booking_id       UUID        NOT NULL,
    doc_type         TEXT        NOT NULL
        CHECK (doc_type IN ('ktp', 'passport', 'photo', 'other')),
    file_path        TEXT        NOT NULL,
    status           TEXT        NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'approved', 'rejected')),
    rejection_reason TEXT,
    uploaded_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    reviewed_at      TIMESTAMPTZ,
    reviewed_by      UUID
);

CREATE INDEX ON jamaah.pilgrim_documents (jamaah_id);
CREATE INDEX ON jamaah.pilgrim_documents (booking_id);
