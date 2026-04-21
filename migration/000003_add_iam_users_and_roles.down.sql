-- Roll back F1.2 iam-svc domain schema.
-- Reverse dependency order: child tables first, then parents, then enums,
-- then the schema itself.

DROP TRIGGER IF EXISTS audit_logs_reject_delete ON iam.audit_logs;
DROP TRIGGER IF EXISTS audit_logs_reject_update ON iam.audit_logs;
DROP FUNCTION IF EXISTS iam.audit_logs_reject_mutation();

DROP TABLE IF EXISTS iam.audit_logs;
DROP TABLE IF EXISTS iam.sessions;
DROP TABLE IF EXISTS iam.user_roles;
DROP TABLE IF EXISTS iam.role_permissions;
DROP TABLE IF EXISTS iam.permissions;
DROP TABLE IF EXISTS iam.roles;
DROP TABLE IF EXISTS iam.users;
DROP TABLE IF EXISTS iam.branches;

DROP TYPE IF EXISTS iam.permission_scope;
DROP TYPE IF EXISTS iam.user_status;

DROP SCHEMA IF EXISTS iam;

-- pgcrypto extension is intentionally NOT dropped — other migrations (and
-- future schemas) will rely on gen_random_uuid().
