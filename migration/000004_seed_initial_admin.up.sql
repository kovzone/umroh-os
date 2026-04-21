-- Seeds the dev IAM starter set: one HQ branch, one super_admin role, one admin user.
-- Serves BL-IAM-001's e2e spec (tests/e2e/tests/04-iam-svc-sessions.spec.ts) + the
-- curl walkthrough in docs/92-testing/testing-guide.md so the reviewer has a real
-- login target without side-channel scripts.
--
-- Fixed UUIDs keep the seed idempotent and let 000004_seed_initial_admin.down.sql
-- delete by ID. Passwords are *dev-only* — prod deployments must re-seed via a
-- one-shot admin-reset flow (land with S1-E-06 BL-IAM-005..017).

BEGIN;

-- HQ branch — parent of all future branches.
INSERT INTO iam.branches (id, name, code, parent_id)
VALUES (
    '11111111-1111-1111-1111-111111111111',
    'Headquarters',
    'HQ',
    NULL
)
ON CONFLICT (id) DO NOTHING;

-- super_admin role — full access.
INSERT INTO iam.roles (id, name, description)
VALUES (
    '22222222-2222-2222-2222-222222222222',
    'super_admin',
    'Full access to every resource and scope. Bootstrap role; restrict in production.'
)
ON CONFLICT (id) DO NOTHING;

-- Admin user — email "admin@umrohos.dev", password "password123" (bcrypt cost 12).
-- Hash generated with `golang.org/x/crypto/bcrypt` locally; the plain password is
-- never stored anywhere in the repo. Rotate in production — this hash is dev-only.
INSERT INTO iam.users (
    id,
    email,
    password_hash,
    name,
    branch_id,
    status
)
VALUES (
    '33333333-3333-3333-3333-333333333333',
    'admin@umrohos.dev',
    '$2a$12$Vo8eTydlb0sz8PMVlIkiqehbv6oLMlk34jEEvWS..WAmd19i7aBxC',
    'Dev Admin',
    '11111111-1111-1111-1111-111111111111',
    'active'
)
ON CONFLICT (id) DO NOTHING;

-- Link admin user → super_admin role.
INSERT INTO iam.user_roles (user_id, role_id)
VALUES (
    '33333333-3333-3333-3333-333333333333',
    '22222222-2222-2222-2222-222222222222'
)
ON CONFLICT (user_id, role_id) DO NOTHING;

COMMIT;
