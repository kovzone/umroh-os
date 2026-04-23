-- roles_admin.sql — admin-side role management queries (S1-E-06 depth card).
--
-- Backing: service/roles_admin.go + grpc_api/roles.go.

-- name: AdminListRoles :many
-- Cursor-paginated roles with permissions aggregated.
-- $1 cursor_time TIMESTAMPTZ (NULL = start from beginning)
-- $2 cursor_id   UUID
-- $3 lim         INT
SELECT
    r.id,
    r.name,
    r.description,
    r.created_at,
    r.updated_at,
    COALESCE(
        ARRAY_AGG(
            p.id::TEXT || ':' || p.resource || ':' || p.action || ':' || p.scope::TEXT
            ORDER BY p.resource, p.action
        ) FILTER (WHERE p.id IS NOT NULL),
        '{}'
    ) AS permission_tuples
FROM iam.roles r
LEFT JOIN iam.role_permissions rp ON rp.role_id    = r.id
LEFT JOIN iam.permissions p       ON p.id          = rp.permission_id
WHERE (
    $1::TIMESTAMPTZ IS NULL
    OR (r.created_at, r.id) < ($1::TIMESTAMPTZ, $2::UUID)
)
GROUP BY r.id
ORDER BY r.created_at DESC, r.id DESC
LIMIT $3;

-- name: AdminGetRoleWithPermissions :one
-- Fetch a single role with its permissions aggregated.
SELECT
    r.id,
    r.name,
    r.description,
    r.created_at,
    r.updated_at,
    COALESCE(
        ARRAY_AGG(
            p.id::TEXT || ':' || p.resource || ':' || p.action || ':' || p.scope::TEXT
            ORDER BY p.resource, p.action
        ) FILTER (WHERE p.id IS NOT NULL),
        '{}'
    ) AS permission_tuples
FROM iam.roles r
LEFT JOIN iam.role_permissions rp ON rp.role_id = r.id
LEFT JOIN iam.permissions p       ON p.id       = rp.permission_id
WHERE r.id = $1
GROUP BY r.id;

-- name: CountUserRolesForRole :one
-- Used by DeleteRole to refuse deletion if any users still hold this role.
SELECT COUNT(*)::BIGINT
FROM iam.user_roles
WHERE role_id = $1;
