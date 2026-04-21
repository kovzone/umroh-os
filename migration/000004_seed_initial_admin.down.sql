-- Reverses 000004_seed_initial_admin.up.sql. Deletes strictly by fixed UUID so
-- this migration never touches rows created by subsequent feature work (e.g. a
-- manually-created second admin).
--
-- Order: join row → user → role → branch (reverse of insertion; respects FKs).

BEGIN;

DELETE FROM iam.user_roles
WHERE user_id = '33333333-3333-3333-3333-333333333333'
  AND role_id = '22222222-2222-2222-2222-222222222222';

DELETE FROM iam.users
WHERE id = '33333333-3333-3333-3333-333333333333';

DELETE FROM iam.roles
WHERE id = '22222222-2222-2222-2222-222222222222';

DELETE FROM iam.branches
WHERE id = '11111111-1111-1111-1111-111111111111';

COMMIT;
