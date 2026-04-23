-- 000024_add_finance_disbursements.up.sql
-- Phase 6.B — Finance AP/AR Depth (BL-FIN-010/011)

CREATE TABLE finance.disbursement_batches (
  id                TEXT PRIMARY KEY,
  description       TEXT NOT NULL,
  total_amount_idr  BIGINT NOT NULL,
  status            TEXT NOT NULL DEFAULT 'pending_approval',
  approved_by       TEXT,
  approved_at       TIMESTAMPTZ,
  created_by        TEXT NOT NULL,
  created_at        TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_disbursement_status ON finance.disbursement_batches(status);

CREATE TABLE finance.disbursement_items (
  id                BIGSERIAL PRIMARY KEY,
  batch_id          TEXT NOT NULL REFERENCES finance.disbursement_batches(id),
  vendor_name       TEXT NOT NULL,
  description       TEXT NOT NULL,
  amount_idr        BIGINT NOT NULL,
  reference         TEXT,
  journal_entry_id  TEXT
);

CREATE INDEX idx_disbursement_items_batch ON finance.disbursement_items(batch_id);
