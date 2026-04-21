-- Reverses 000005_seed_iam_test_roles_and_permissions.up.sql. Strictly deletes
-- by fixed UUID so subsequent feature-level inserts (e.g. staff users added
-- in S1-E-06) are never touched.
--
-- Order: user_roles → users → role_permissions → permission/roles (reverse of
-- insertion; respects FKs).

BEGIN;

DELETE FROM iam.user_roles
WHERE (user_id, role_id) IN (
    ('77777777-7777-7777-7777-777777777777', '55555555-5555-5555-5555-555555555555'),
    ('88888888-8888-8888-8888-888888888888', '66666666-6666-6666-6666-666666666666')
);

DELETE FROM iam.users
WHERE id IN (
    '77777777-7777-7777-7777-777777777777',
    '88888888-8888-8888-8888-888888888888'
);

DELETE FROM iam.role_permissions
WHERE permission_id = '44444444-4444-4444-4444-444444444444';

DELETE FROM iam.permissions
WHERE id = '44444444-4444-4444-4444-444444444444';

DELETE FROM iam.roles
WHERE id IN (
    '55555555-5555-5555-5555-555555555555',
    '66666666-6666-6666-6666-666666666666'
);

COMMIT;
