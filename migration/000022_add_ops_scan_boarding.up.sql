-- 000022_add_ops_scan_boarding.up.sql
-- Phase 6.D — Field Operations Scan (BL-OPS-010/011)

CREATE TABLE ops.scan_events (
  id               TEXT PRIMARY KEY,
  scan_type        TEXT NOT NULL,
  departure_id     TEXT NOT NULL,
  jamaah_id        TEXT NOT NULL,
  scanned_by       TEXT NOT NULL,
  device_id        TEXT,
  location         TEXT,
  idempotency_key  TEXT NOT NULL UNIQUE,
  metadata         JSONB DEFAULT '{}',
  scanned_at       TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_scan_events_departure ON ops.scan_events(departure_id);
CREATE INDEX idx_scan_events_jamaah    ON ops.scan_events(jamaah_id);

CREATE TABLE ops.bus_boardings (
  id              TEXT PRIMARY KEY,
  departure_id    TEXT NOT NULL,
  bus_number      TEXT NOT NULL,
  jamaah_id       TEXT NOT NULL,
  status          TEXT NOT NULL DEFAULT 'boarded',
  scan_event_id   TEXT REFERENCES ops.scan_events(id),
  boarded_at      TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_bus_boardings_dep_jamaah ON ops.bus_boardings(departure_id, jamaah_id);
CREATE INDEX idx_bus_boardings_departure ON ops.bus_boardings(departure_id);
