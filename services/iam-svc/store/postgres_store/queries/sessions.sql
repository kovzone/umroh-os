-- name: CreateSession :one
INSERT INTO iam.sessions (user_id, refresh_token_hash, user_agent, ip, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, refresh_token_hash, user_agent, ip,
          issued_at, expires_at, revoked_at;

-- name: GetSessionByID :one
SELECT id, user_id, refresh_token_hash, user_agent, ip,
       issued_at, expires_at, revoked_at
FROM iam.sessions
WHERE id = $1;

-- name: GetSessionByRefreshHash :one
SELECT id, user_id, refresh_token_hash, user_agent, ip,
       issued_at, expires_at, revoked_at
FROM iam.sessions
WHERE refresh_token_hash = $1;

-- name: ListActiveSessionsForUser :many
SELECT id, user_id, refresh_token_hash, user_agent, ip,
       issued_at, expires_at, revoked_at
FROM iam.sessions
WHERE user_id = $1
  AND revoked_at IS NULL
  AND expires_at > now()
ORDER BY issued_at DESC;

-- name: RevokeSession :exec
UPDATE iam.sessions
SET revoked_at = now()
WHERE id = $1 AND revoked_at IS NULL;

-- name: RevokeAllSessionsForUser :exec
UPDATE iam.sessions
SET revoked_at = now()
WHERE user_id = $1 AND revoked_at IS NULL;
