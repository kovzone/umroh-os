-- name: AssignRoleToUser :exec
INSERT INTO iam.user_roles (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT (user_id, role_id) DO NOTHING;

-- name: RevokeRoleFromUser :exec
DELETE FROM iam.user_roles
WHERE user_id = $1 AND role_id = $2;

-- name: ListRoleIDsForUser :many
SELECT role_id
FROM iam.user_roles
WHERE user_id = $1;

-- name: ListUserIDsForRole :many
SELECT user_id
FROM iam.user_roles
WHERE role_id = $1;

-- name: RevokeAllRolesForUser :exec
DELETE FROM iam.user_roles
WHERE user_id = $1;
