-- Reverses 000006_seed_iam_user_suspend_permission.up.sql. Strictly deletes by
-- fixed UUID so later role-grant inserts (S1-E-06 admin CRUD) are never touched.
--
-- Order:
--   sessions rows for the sacrifice user → sacrifice user row →
--   role_permissions grant → permission row.
-- (Sessions rolled back first to respect the iam.sessions.user_id FK.)

BEGIN;

DELETE FROM iam.sessions
WHERE user_id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

DELETE FROM iam.users
WHERE id = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';

DELETE FROM iam.role_permissions
WHERE permission_id = '99999999-9999-9999-9999-999999999999';

DELETE FROM iam.permissions
WHERE id = '99999999-9999-9999-9999-999999999999';

COMMIT;
