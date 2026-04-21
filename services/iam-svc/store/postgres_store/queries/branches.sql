-- name: CreateBranch :one
INSERT INTO iam.branches (name, code, parent_id)
VALUES ($1, $2, $3)
RETURNING id, name, code, parent_id, created_at, updated_at;

-- name: GetBranchByID :one
SELECT id, name, code, parent_id, created_at, updated_at
FROM iam.branches
WHERE id = $1;

-- name: GetBranchByCode :one
SELECT id, name, code, parent_id, created_at, updated_at
FROM iam.branches
WHERE code = $1;

-- name: ListBranches :many
SELECT id, name, code, parent_id, created_at, updated_at
FROM iam.branches
ORDER BY name;

-- name: UpdateBranch :one
UPDATE iam.branches
SET name = $2,
    code = $3,
    parent_id = $4,
    updated_at = now()
WHERE id = $1
RETURNING id, name, code, parent_id, created_at, updated_at;

-- name: DeleteBranch :exec
DELETE FROM iam.branches
WHERE id = $1;
