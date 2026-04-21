-- name: CreateRole :one
INSERT INTO iam.roles (name, description)
VALUES ($1, $2)
RETURNING id, name, description, created_at, updated_at;

-- name: GetRoleByID :one
SELECT id, name, description, created_at, updated_at
FROM iam.roles
WHERE id = $1;

-- name: GetRoleByName :one
SELECT id, name, description, created_at, updated_at
FROM iam.roles
WHERE name = $1;

-- name: ListRoles :many
SELECT id, name, description, created_at, updated_at
FROM iam.roles
ORDER BY name;

-- name: UpdateRole :one
UPDATE iam.roles
SET name = $2,
    description = $3,
    updated_at = now()
WHERE id = $1
RETURNING id, name, description, created_at, updated_at;

-- name: DeleteRole :exec
DELETE FROM iam.roles
WHERE id = $1;
