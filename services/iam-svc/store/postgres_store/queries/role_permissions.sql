-- name: GrantPermissionToRole :exec
INSERT INTO iam.role_permissions (role_id, permission_id)
VALUES ($1, $2)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- name: RevokePermissionFromRole :exec
DELETE FROM iam.role_permissions
WHERE role_id = $1 AND permission_id = $2;

-- name: ListPermissionIDsForRole :many
SELECT permission_id
FROM iam.role_permissions
WHERE role_id = $1;

-- name: ListRoleIDsForPermission :many
SELECT role_id
FROM iam.role_permissions
WHERE permission_id = $1;

-- name: RevokeAllPermissionsForRole :exec
DELETE FROM iam.role_permissions
WHERE role_id = $1;
