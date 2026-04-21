-- name: CreatePermission :one
INSERT INTO iam.permissions (resource, action, scope)
VALUES ($1, $2, $3)
RETURNING id, resource, action, scope, created_at;

-- name: GetPermissionByID :one
SELECT id, resource, action, scope, created_at
FROM iam.permissions
WHERE id = $1;

-- name: GetPermissionByTuple :one
SELECT id, resource, action, scope, created_at
FROM iam.permissions
WHERE resource = $1 AND action = $2 AND scope = $3;

-- name: ListPermissions :many
SELECT id, resource, action, scope, created_at
FROM iam.permissions
ORDER BY resource, action, scope;

-- name: ListPermissionsByResource :many
SELECT id, resource, action, scope, created_at
FROM iam.permissions
WHERE resource = $1
ORDER BY action, scope;

-- name: DeletePermission :exec
DELETE FROM iam.permissions
WHERE id = $1;

-- UserHasPermission resolves whether the given user currently holds the
-- (resource, action, scope) permission via any of their assigned roles.
-- Backs iam.v1.IamService/CheckPermission (BL-IAM-002); the hot path must
-- stay a single index-backed join.
--
-- name: UserHasPermission :one
SELECT EXISTS (
    SELECT 1
    FROM iam.user_roles ur
    JOIN iam.role_permissions rp ON rp.role_id = ur.role_id
    JOIN iam.permissions p       ON p.id      = rp.permission_id
    WHERE ur.user_id = $1
      AND p.resource = $2
      AND p.action   = $3
      AND p.scope    = $4
) AS allowed;
