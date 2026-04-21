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
