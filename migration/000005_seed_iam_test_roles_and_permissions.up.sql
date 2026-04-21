-- BL-IAM-002 test fixtures: one finance permission tuple + two auxiliary
-- roles (finance_admin, cs_agent) + two dev-only users so the `CheckPermission`
-- acceptance criterion "Finance routes denied for non-finance roles" has
-- concrete rows to drive integration tests against.
--
-- 000004 seeded only the super_admin role with NO permissions attached.
-- This migration keeps super_admin's grant explicit (no implicit wildcard) and
-- introduces a granular fixture set the e2e suite can log into.
--
-- All rows use fixed UUIDs for idempotency; rolling back via 000005.down.sql
-- deletes by id. Passwords are dev-only — prod deployments must re-seed via a
-- one-shot admin-reset flow (land with S1-E-06 BL-IAM-005..017). Password
-- hash is the same bcrypt($2a$12) digest of "password123" used in 000004 for
-- admin@umrohos.dev — one hash, multiple fixture users, never in production.

BEGIN;

-- ---------------------------------------------------------------------------
-- Permission: journal_entry / read / global
-- ---------------------------------------------------------------------------

INSERT INTO iam.permissions (id, resource, action, scope)
VALUES (
    '44444444-4444-4444-4444-444444444444',
    'journal_entry',
    'read',
    'global'
)
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- Roles: finance_admin + cs_agent
-- super_admin is already present from 000004.
-- ---------------------------------------------------------------------------

INSERT INTO iam.roles (id, name, description)
VALUES (
    '55555555-5555-5555-5555-555555555555',
    'finance_admin',
    'Operates the finance module: journals, AR/AP, period close, reports.'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO iam.roles (id, name, description)
VALUES (
    '66666666-6666-6666-6666-666666666666',
    'cs_agent',
    'Customer-service front-line: booking lookups, jamaah outreach. No finance access.'
)
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- Role → permission grants.
-- super_admin + finance_admin both receive journal_entry/read/global.
-- cs_agent receives nothing — drives the deny path.
-- ---------------------------------------------------------------------------

INSERT INTO iam.role_permissions (role_id, permission_id)
VALUES
    ('22222222-2222-2222-2222-222222222222', '44444444-4444-4444-4444-444444444444'),
    ('55555555-5555-5555-5555-555555555555', '44444444-4444-4444-4444-444444444444')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- Fixture users: finance@umrohos.dev + cs@umrohos.dev.
-- Password "password123" — same bcrypt hash as admin@umrohos.dev in 000004.
-- Branch: HQ (seeded in 000004).
-- ---------------------------------------------------------------------------

INSERT INTO iam.users (
    id,
    email,
    password_hash,
    name,
    branch_id,
    status
)
VALUES (
    '77777777-7777-7777-7777-777777777777',
    'finance@umrohos.dev',
    '$2a$12$Vo8eTydlb0sz8PMVlIkiqehbv6oLMlk34jEEvWS..WAmd19i7aBxC',
    'Dev Finance Admin',
    '11111111-1111-1111-1111-111111111111',
    'active'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO iam.users (
    id,
    email,
    password_hash,
    name,
    branch_id,
    status
)
VALUES (
    '88888888-8888-8888-8888-888888888888',
    'cs@umrohos.dev',
    '$2a$12$Vo8eTydlb0sz8PMVlIkiqehbv6oLMlk34jEEvWS..WAmd19i7aBxC',
    'Dev CS Agent',
    '11111111-1111-1111-1111-111111111111',
    'active'
)
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- User → role bindings.
-- ---------------------------------------------------------------------------

INSERT INTO iam.user_roles (user_id, role_id)
VALUES
    ('77777777-7777-7777-7777-777777777777', '55555555-5555-5555-5555-555555555555'),
    ('88888888-8888-8888-8888-888888888888', '66666666-6666-6666-6666-666666666666')
ON CONFLICT (user_id, role_id) DO NOTHING;

COMMIT;
