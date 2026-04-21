-- F1.2 — iam-svc domain schema.
--
-- Scaffolds the 8 tables + 2 enums that back F1 (Identity, Access, Audit).
-- Mirrors `docs/03-services/00-iam-svc/02-data-model.md` row-for-row.
--
-- All objects live in a dedicated `iam` schema per ADR 0007's per-service
-- schema-namespace convention. sqlc reads this file via its shared
-- `schema: "../../../../migration"` pointer in each service's sqlc.yaml.
--
-- The `audit_logs` table is append-only: an INSTEAD-OF trigger blocks
-- UPDATE / DELETE at the DB layer, enforcing Q008 (UU PDP) + F1
-- Acceptance criteria "Audit log rows cannot be mutated."

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS iam;

-- ---------------------------------------------------------------------------
-- Enum types
-- ---------------------------------------------------------------------------

CREATE TYPE iam.user_status AS ENUM (
    'active',
    'suspended',
    'pending'
);

CREATE TYPE iam.permission_scope AS ENUM (
    'global',
    'branch',
    'personal'
);

-- ---------------------------------------------------------------------------
-- branches
-- ---------------------------------------------------------------------------

CREATE TABLE iam.branches (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL,
    code       TEXT NOT NULL,
    parent_id  UUID REFERENCES iam.branches (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT branches_code_uk UNIQUE (code)
);

CREATE INDEX branches_parent_id_idx ON iam.branches (parent_id);

-- ---------------------------------------------------------------------------
-- users
-- ---------------------------------------------------------------------------

CREATE TABLE iam.users (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email            TEXT NOT NULL,
    password_hash    TEXT NOT NULL,
    name             TEXT NOT NULL,
    branch_id        UUID NOT NULL REFERENCES iam.branches (id) ON DELETE RESTRICT,
    status           iam.user_status NOT NULL DEFAULT 'pending',
    totp_secret      TEXT,
    totp_verified_at TIMESTAMPTZ,
    last_login_at    TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at       TIMESTAMPTZ
);

-- Email is unique per live user (soft-deleted rows can re-use the address).
CREATE UNIQUE INDEX users_email_live_uk
    ON iam.users (email)
    WHERE deleted_at IS NULL;

CREATE INDEX users_branch_id_idx ON iam.users (branch_id);
CREATE INDEX users_status_idx ON iam.users (status) WHERE deleted_at IS NULL;

-- ---------------------------------------------------------------------------
-- roles
-- ---------------------------------------------------------------------------

CREATE TABLE iam.roles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT roles_name_uk UNIQUE (name)
);

-- ---------------------------------------------------------------------------
-- permissions
-- ---------------------------------------------------------------------------

CREATE TABLE iam.permissions (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource   TEXT NOT NULL,
    action     TEXT NOT NULL,
    scope      iam.permission_scope NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT permissions_resource_action_scope_uk UNIQUE (resource, action, scope)
);

CREATE INDEX permissions_resource_idx ON iam.permissions (resource);

-- ---------------------------------------------------------------------------
-- role_permissions  (join)
-- ---------------------------------------------------------------------------

CREATE TABLE iam.role_permissions (
    role_id       UUID NOT NULL REFERENCES iam.roles (id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES iam.permissions (id) ON DELETE CASCADE,
    granted_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX role_permissions_permission_id_idx ON iam.role_permissions (permission_id);

-- ---------------------------------------------------------------------------
-- user_roles  (join)
-- ---------------------------------------------------------------------------

CREATE TABLE iam.user_roles (
    user_id    UUID NOT NULL REFERENCES iam.users (id) ON DELETE CASCADE,
    role_id    UUID NOT NULL REFERENCES iam.roles (id) ON DELETE CASCADE,
    granted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX user_roles_role_id_idx ON iam.user_roles (role_id);

-- ---------------------------------------------------------------------------
-- sessions
-- ---------------------------------------------------------------------------

CREATE TABLE iam.sessions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES iam.users (id) ON DELETE CASCADE,
    refresh_token_hash  TEXT NOT NULL,
    user_agent          TEXT NOT NULL DEFAULT '',
    ip                  INET,
    issued_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at          TIMESTAMPTZ NOT NULL,
    revoked_at          TIMESTAMPTZ,
    CONSTRAINT sessions_refresh_token_hash_uk UNIQUE (refresh_token_hash)
);

CREATE INDEX sessions_user_id_active_idx
    ON iam.sessions (user_id)
    WHERE revoked_at IS NULL;

CREATE INDEX sessions_expires_at_idx
    ON iam.sessions (expires_at)
    WHERE revoked_at IS NULL;

-- ---------------------------------------------------------------------------
-- audit_logs  (append-only)
-- ---------------------------------------------------------------------------

CREATE TABLE iam.audit_logs (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID REFERENCES iam.users (id) ON DELETE SET NULL,
    branch_id   UUID REFERENCES iam.branches (id) ON DELETE SET NULL,
    resource    TEXT NOT NULL,
    resource_id TEXT NOT NULL DEFAULT '',
    action      TEXT NOT NULL,
    old_value   JSONB,
    new_value   JSONB,
    ip          INET,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX audit_logs_user_id_created_at_idx
    ON iam.audit_logs (user_id, created_at DESC);

CREATE INDEX audit_logs_resource_created_at_idx
    ON iam.audit_logs (resource, created_at DESC);

-- Append-only enforcement: reject UPDATE and DELETE on audit_logs at the
-- database layer so a compromised service role cannot quietly rewrite
-- history. Per F1 Acceptance: "Audit log rows cannot be mutated."

CREATE OR REPLACE FUNCTION iam.audit_logs_reject_mutation()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    RAISE EXCEPTION 'iam.audit_logs is append-only; % is not permitted',
        TG_OP
        USING ERRCODE = 'insufficient_privilege';
END;
$$;

CREATE TRIGGER audit_logs_reject_update
    BEFORE UPDATE ON iam.audit_logs
    FOR EACH ROW EXECUTE FUNCTION iam.audit_logs_reject_mutation();

CREATE TRIGGER audit_logs_reject_delete
    BEFORE DELETE ON iam.audit_logs
    FOR EACH ROW EXECUTE FUNCTION iam.audit_logs_reject_mutation();
