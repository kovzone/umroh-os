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

-- ListRoleNamesForUser returns the role names currently assigned to the user.
-- Backs iam.v1.IamService/ValidateToken so downstream services receive role
-- strings alongside identity claims (BL-IAM-002).
--
-- name: ListRoleNamesForUser :many
SELECT r.name
FROM iam.user_roles ur
JOIN iam.roles r ON r.id = ur.role_id
WHERE ur.user_id = $1
ORDER BY r.name;
