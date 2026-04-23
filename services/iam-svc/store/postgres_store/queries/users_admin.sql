-- users_admin.sql — admin-side user management queries (S1-E-06 depth card).
--
-- These complement users.sql which backs the auth hot path.
-- Backing: service/users_admin.go + grpc_api/users.go.

-- name: AdminListUsers :many
-- Cursor-paginated list with optional status / branch_id filters.
-- Caller passes:
--   $1  status_filter  TEXT ('' = no filter)
--   $2  branch_filter  UUID (zero-UUID = no filter)
--   $3  cursor_time    TIMESTAMPTZ (NULL = start from beginning, i.e. newest first)
--   $4  cursor_id      UUID        (used together with cursor_time for tie-breaking)
--   $5  lim            INT
SELECT
    u.id,
    u.email,
    u.name,
    u.branch_id,
    u.status,
    u.last_login_at,
    u.created_at,
    COALESCE(
        ARRAY_AGG(r.name ORDER BY r.name) FILTER (WHERE r.name IS NOT NULL),
        '{}'
    ) AS role_names
FROM iam.users u
LEFT JOIN iam.user_roles ur ON ur.user_id = u.id
LEFT JOIN iam.roles r       ON r.id       = ur.role_id
WHERE u.deleted_at IS NULL
  AND ($1::TEXT  = '' OR u.status::TEXT = $1)
  AND ($2::UUID  = '00000000-0000-0000-0000-000000000000' OR u.branch_id = $2)
  AND (
      $3::TIMESTAMPTZ IS NULL
      OR (u.created_at, u.id) < ($3::TIMESTAMPTZ, $4::UUID)
  )
GROUP BY u.id
ORDER BY u.created_at DESC, u.id DESC
LIMIT $5;

-- name: AdminGetUser :one
-- Return a single user row; role_names aggregated from user_roles.
SELECT
    u.id,
    u.email,
    u.name,
    u.branch_id,
    u.status,
    u.last_login_at,
    u.created_at,
    u.updated_at,
    COALESCE(
        ARRAY_AGG(r.name ORDER BY r.name) FILTER (WHERE r.name IS NOT NULL),
        '{}'
    ) AS role_names
FROM iam.users u
LEFT JOIN iam.user_roles ur ON ur.user_id = u.id
LEFT JOIN iam.roles r       ON r.id       = ur.role_id
WHERE u.id = $1 AND u.deleted_at IS NULL
GROUP BY u.id;

-- name: AdminUpdateUserNameStatus :one
-- Update name and/or status; returns updated row (without roles — caller re-fetches).
UPDATE iam.users
SET name       = $2,
    status     = $3,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, email, name, branch_id, status, last_login_at, created_at, updated_at;

-- name: RevokeAllSessionsForUserAdmin :exec
-- Revoke all active sessions for a user (used by ResetUserPassword).
UPDATE iam.sessions
SET revoked_at = now()
WHERE user_id = $1 AND revoked_at IS NULL;
