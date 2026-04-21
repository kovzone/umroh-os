-- name: CreateUser :one
INSERT INTO iam.users (email, password_hash, name, branch_id, status)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, email, password_hash, name, branch_id, status,
          totp_secret, totp_verified_at, last_login_at,
          created_at, updated_at, deleted_at;

-- name: GetUserByID :one
SELECT id, email, password_hash, name, branch_id, status,
       totp_secret, totp_verified_at, last_login_at,
       created_at, updated_at, deleted_at
FROM iam.users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, name, branch_id, status,
       totp_secret, totp_verified_at, last_login_at,
       created_at, updated_at, deleted_at
FROM iam.users
WHERE email = $1 AND deleted_at IS NULL;

-- name: ListUsersByBranch :many
SELECT id, email, password_hash, name, branch_id, status,
       totp_secret, totp_verified_at, last_login_at,
       created_at, updated_at, deleted_at
FROM iam.users
WHERE branch_id = $1 AND deleted_at IS NULL
ORDER BY name;

-- name: UpdateUserProfile :one
UPDATE iam.users
SET name = $2,
    branch_id = $3,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, email, password_hash, name, branch_id, status,
          totp_secret, totp_verified_at, last_login_at,
          created_at, updated_at, deleted_at;

-- name: UpdateUserPasswordHash :exec
UPDATE iam.users
SET password_hash = $2,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateUserStatus :exec
UPDATE iam.users
SET status = $2,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateUserTOTP :exec
UPDATE iam.users
SET totp_secret = $2,
    totp_verified_at = $3,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateUserLastLoginAt :exec
UPDATE iam.users
SET last_login_at = now(),
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteUser :exec
UPDATE iam.users
SET deleted_at = now(),
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;
