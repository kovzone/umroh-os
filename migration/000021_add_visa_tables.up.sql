-- 000021_add_visa_tables.up.sql
-- Phase 6.E — Visa Pipeline (BL-VISA-001..003)

CREATE SCHEMA IF NOT EXISTS visa;

CREATE TABLE visa.visa_applications (
  id             TEXT PRIMARY KEY,
  jamaah_id      TEXT NOT NULL,
  booking_id     TEXT NOT NULL,
  departure_id   TEXT NOT NULL,
  status         TEXT NOT NULL DEFAULT 'WAITING_DOCS',
  provider_id    TEXT,
  provider_ref   TEXT,
  e_visa_url     TEXT,
  issued_date    DATE,
  created_at     TIMESTAMPTZ DEFAULT NOW(),
  updated_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_visa_apps_departure ON visa.visa_applications(departure_id);
CREATE INDEX idx_visa_apps_jamaah    ON visa.visa_applications(jamaah_id);
CREATE INDEX idx_visa_apps_status    ON visa.visa_applications(status);

CREATE TABLE visa.status_history (
  id             BIGSERIAL PRIMARY KEY,
  application_id TEXT NOT NULL REFERENCES visa.visa_applications(id),
  from_status    TEXT,
  to_status      TEXT NOT NULL,
  reason         TEXT,
  actor_user_id  TEXT,
  created_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_visa_history_app ON visa.status_history(application_id);
