-- BL-IAM-003 fixtures. Two pieces:
--
--   1. One permission tuple for the admin-suspend action, granted to super_admin.
--      Drives the REST handler `POST /v1/users/{id}/suspend`, which calls
--      `service.CheckPermission(actor, "iam.users", "suspend", "global")` before
--      delegating to `service.SuspendUser`.
--
--   2. One dedicated "sacrifice" user (`suspend-target@umrohos.dev`) the e2e spec
--      `02c-iam-svc-suspend.spec.ts` can safely suspend on every run. Using this
--      user instead of cs@/finance@ keeps the 02b permission-gate spec green —
--      cs@ must stay active to prove its deny-path against /v1/finance/ping.
--      Re-running 02c is a no-op on status (idempotent) and still sweeps any
--      sessions that race in; neither side effect breaks 02b.
--
-- 000004 seeded super_admin with zero permissions; 000005 attached the finance
-- fixture. This migration layers the admin-only suspend grant + the sacrifice
-- user. Fixed UUIDs keep the seed idempotent; `.down.sql` deletes by id.
-- Real prod grants and user provisioning land with S1-E-06 admin CRUD — this
-- file is dev fixtures only.

BEGIN;

-- ---------------------------------------------------------------------------
-- Permission: iam.users / suspend / global
-- ---------------------------------------------------------------------------

INSERT INTO iam.permissions (id, resource, action, scope)
VALUES (
    '99999999-9999-9999-9999-999999999999',
    'iam.users',
    'suspend',
    'global'
)
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- Grant: super_admin → iam.users/suspend/global.
-- finance_admin + cs_agent intentionally do NOT receive this grant — they drive
-- the deny path in the e2e.
-- ---------------------------------------------------------------------------

INSERT INTO iam.role_permissions (role_id, permission_id)
VALUES (
    '22222222-2222-2222-2222-222222222222',
    '99999999-9999-9999-9999-999999999999'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- Fixture user: suspend-target@umrohos.dev (UUID aaaa...). Status active, no
-- role — the e2e only needs them to log in, not to perform authorized actions.
-- Password "password123" reuses the dev-only hash from 000004/000005.
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
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'suspend-target@umrohos.dev',
    '$2a$12$Vo8eTydlb0sz8PMVlIkiqehbv6oLMlk34jEEvWS..WAmd19i7aBxC',
    'Dev Suspend Target',
    '11111111-1111-1111-1111-111111111111',
    'active'
)
ON CONFLICT (id) DO NOTHING;

COMMIT;
