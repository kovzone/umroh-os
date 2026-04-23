-- 000025_add_iam_admin_tables.up.sql
-- Phase 6.M — IAM Admin/Security Depth (BL-IAM-007/010/011/014/016)

CREATE TABLE iam.data_scopes (
  id          TEXT PRIMARY KEY,
  user_id     TEXT NOT NULL UNIQUE,
  scope_type  TEXT NOT NULL,
  branch_id   TEXT,
  created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE iam.api_keys (
  id           TEXT PRIMARY KEY,
  name         TEXT NOT NULL,
  key_hash     TEXT NOT NULL UNIQUE,
  key_prefix   TEXT NOT NULL,
  created_by   TEXT NOT NULL,
  expires_at   TIMESTAMPTZ,
  last_used_at TIMESTAMPTZ,
  revoked_at   TIMESTAMPTZ,
  scopes       TEXT[] NOT NULL DEFAULT '{}',
  created_at   TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_api_keys_prefix    ON iam.api_keys(key_prefix);
CREATE INDEX idx_api_keys_revoked   ON iam.api_keys(revoked_at);

CREATE TABLE iam.global_config (
  key         TEXT PRIMARY KEY,
  value       TEXT NOT NULL,
  description TEXT,
  updated_by  TEXT,
  updated_at  TIMESTAMPTZ DEFAULT NOW()
);
