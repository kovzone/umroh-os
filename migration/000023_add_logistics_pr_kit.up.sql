-- 000023_add_logistics_pr_kit.up.sql
-- Phase 6.C — Warehouse & Procurement (BL-LOG-010..012)

CREATE TABLE logistics.purchase_requests (
  id               TEXT PRIMARY KEY,
  departure_id     TEXT NOT NULL,
  requested_by     TEXT NOT NULL,
  item_name        TEXT NOT NULL,
  quantity         INT NOT NULL,
  unit_price_idr   BIGINT NOT NULL,
  total_price_idr  BIGINT NOT NULL,
  budget_limit_idr BIGINT,
  status           TEXT NOT NULL DEFAULT 'pending',
  approved_by      TEXT,
  notes            TEXT,
  created_at       TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_pr_departure ON logistics.purchase_requests(departure_id);
CREATE INDEX idx_pr_status    ON logistics.purchase_requests(status);

CREATE TABLE logistics.kit_assemblies (
  id               TEXT PRIMARY KEY,
  departure_id     TEXT NOT NULL,
  assembled_by     TEXT NOT NULL,
  status           TEXT NOT NULL DEFAULT 'pending',
  idempotency_key  TEXT NOT NULL UNIQUE,
  assembled_at     TIMESTAMPTZ,
  created_at       TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_kit_departure ON logistics.kit_assemblies(departure_id);

CREATE TABLE logistics.kit_assembly_items (
  id           BIGSERIAL PRIMARY KEY,
  assembly_id  TEXT NOT NULL REFERENCES logistics.kit_assemblies(id),
  item_name    TEXT NOT NULL,
  quantity     INT NOT NULL,
  fulfilled    BOOLEAN NOT NULL DEFAULT FALSE
);
