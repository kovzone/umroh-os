-- name: InsertAuditLog :one
INSERT INTO iam.audit_logs (
    user_id, branch_id, resource, resource_id, action,
    old_value, new_value, ip
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, user_id, branch_id, resource, resource_id, action,
          old_value, new_value, ip, created_at;

-- name: GetAuditLogByID :one
SELECT id, user_id, branch_id, resource, resource_id, action,
       old_value, new_value, ip, created_at
FROM iam.audit_logs
WHERE id = $1;

-- name: ListRecentAuditLogsByUser :many
SELECT id, user_id, branch_id, resource, resource_id, action,
       old_value, new_value, ip, created_at
FROM iam.audit_logs
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: ListRecentAuditLogsByResource :many
SELECT id, user_id, branch_id, resource, resource_id, action,
       old_value, new_value, ip, created_at
FROM iam.audit_logs
WHERE resource = $1
ORDER BY created_at DESC
LIMIT $2;
