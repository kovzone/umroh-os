-- 000012 — Payment tables (F5 / S2-E-01).
--
-- Introduces the payment schema and four core tables:
--   invoices          — the billing record for a booking
--   virtual_accounts  — gateway-issued VA rows linked to an invoice
--   payment_events    — append-only ledger of every money movement
--   refunds           — refund lifecycle tracking
--
-- All IDs are UUID (gen_random_uuid()).  Status columns use VARCHAR(20/30)
-- with CHECK constraints rather than ENUM types so new values can be added
-- without blocking migrations (no ALTER TYPE required).
--
-- Schema: payment (per ADR 0007 per-service namespacing; same DB per ADR 0003).

CREATE SCHEMA IF NOT EXISTS payment;

-- ---------------------------------------------------------------------------
-- invoices
-- ---------------------------------------------------------------------------

CREATE TABLE payment.invoices (
    id                      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id              UUID        NOT NULL,
    amount_total            NUMERIC(15,2) NOT NULL,       -- IDR after Rp 1,000 half-up rounding (Q001)
    rounding_adjustment_idr NUMERIC(15,2) NOT NULL DEFAULT 0,  -- signed delta vs pre-round sum
    currency                VARCHAR(3)  NOT NULL DEFAULT 'IDR',
    fx_snapshot             JSONB       NOT NULL DEFAULT '{}', -- FX rates at issuance; immutable after first payment
    status                  VARCHAR(20) NOT NULL DEFAULT 'unpaid'
        CHECK (status IN ('unpaid', 'partially_paid', 'paid', 'void', 'refunded')),
    paid_amount             NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX invoices_booking_id_idx     ON payment.invoices (booking_id);
CREATE INDEX invoices_status_created_idx ON payment.invoices (status, created_at DESC);

-- ---------------------------------------------------------------------------
-- virtual_accounts
-- ---------------------------------------------------------------------------

CREATE TABLE payment.virtual_accounts (
    id               UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id       UUID         NOT NULL REFERENCES payment.invoices (id),
    gateway          VARCHAR(20)  NOT NULL
        CHECK (gateway IN ('midtrans', 'xendit', 'mock')),
    gateway_va_id    VARCHAR(255) NOT NULL,
    account_number   VARCHAR(50)  NOT NULL,
    bank_code        VARCHAR(10)  NOT NULL,
    expires_at       TIMESTAMPTZ  NOT NULL,
    status           VARCHAR(20)  NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'expired', 'paid')),
    -- idempotency_key = invoice.id (string representation of UUID)
    idempotency_key  VARCHAR(255) NOT NULL UNIQUE,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX va_invoice_id_idx ON payment.virtual_accounts (invoice_id);
CREATE INDEX va_status_idx     ON payment.virtual_accounts (status) WHERE status = 'active';

-- ---------------------------------------------------------------------------
-- payment_events  (append-only; never UPDATE/DELETE)
-- ---------------------------------------------------------------------------

CREATE TABLE payment.payment_events (
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id      UUID         NOT NULL REFERENCES payment.invoices (id),
    gateway         VARCHAR(20)  NOT NULL,
    gateway_txn_id  VARCHAR(255),
    kind            VARCHAR(30)  NOT NULL
        CHECK (kind IN ('va_created', 'payment_received', 'settlement_received', 'refund_issued', 'manual')),
    amount          NUMERIC(15,2) NOT NULL DEFAULT 0,
    raw_payload     JSONB,
    approval_status VARCHAR(20)
        CHECK (approval_status IS NULL OR approval_status IN ('pending', 'approved', 'rejected')),
    received_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    -- Webhook idempotency: (gateway, gateway_txn_id) must be unique.
    -- gateway_txn_id is NULL for manual events; partial UNIQUE index excludes NULLs.
    CONSTRAINT payment_events_gateway_txn_id_unique UNIQUE (gateway, gateway_txn_id)
);

CREATE INDEX pe_invoice_id_idx      ON payment.payment_events (invoice_id);
CREATE INDEX pe_received_at_idx     ON payment.payment_events (received_at DESC);
CREATE INDEX pe_gateway_txn_id_idx  ON payment.payment_events (gateway, gateway_txn_id)
    WHERE gateway_txn_id IS NOT NULL;

-- ---------------------------------------------------------------------------
-- refunds
-- ---------------------------------------------------------------------------

CREATE TABLE payment.refunds (
    id                UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id        UUID         NOT NULL REFERENCES payment.invoices (id),
    booking_id        UUID         NOT NULL,
    amount            NUMERIC(15,2) NOT NULL,
    reason_code       VARCHAR(50),
    status            VARCHAR(20)  NOT NULL DEFAULT 'requested'
        CHECK (status IN ('requested', 'processing', 'completed', 'failed')),
    gateway_refund_id VARCHAR(255),
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX refunds_invoice_id_idx  ON payment.refunds (invoice_id);
CREATE INDEX refunds_booking_id_idx  ON payment.refunds (booking_id);
CREATE INDEX refunds_status_idx      ON payment.refunds (status);
