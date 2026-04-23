// Code generated manually (sqlc pattern) for roles_admin.sql.
// S1-E-06 depth card — admin role management queries.

package sqlc

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

// -------------------------------------------------------------------------
// AdminListRoles
// -------------------------------------------------------------------------

const adminListRoles = `-- name: AdminListRoles :many
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
LIMIT $3
`

type AdminListRolesParams struct {
	CursorTime pgtype.Timestamptz `json:"cursor_time"`
	CursorID   pgtype.UUID        `json:"cursor_id"`
	Lim        int32              `json:"lim"`
}

// AdminRoleRow is a role row with its permissions encoded as
// "uuid:resource:action:scope" tuples. The service layer decodes these.
type AdminRoleRow struct {
	ID               pgtype.UUID        `json:"id"`
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `json:"updated_at"`
	PermissionTuples []string           `json:"permission_tuples"`
}

// DecodePermissions splits the "id:resource:action:scope" tuples into
// AdminPermissionItem slices. Called by the service layer.
func (r *AdminRoleRow) DecodePermissions() []AdminPermissionItem {
	items := make([]AdminPermissionItem, 0, len(r.PermissionTuples))
	for _, t := range r.PermissionTuples {
		parts := strings.SplitN(t, ":", 4)
		if len(parts) != 4 {
			continue
		}
		items = append(items, AdminPermissionItem{
			ID:       parts[0],
			Resource: parts[1],
			Action:   parts[2],
			Scope:    parts[3],
		})
	}
	return items
}

// AdminPermissionItem is the decoded form of a permission tuple.
type AdminPermissionItem struct {
	ID       string `json:"id"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Scope    string `json:"scope"`
}

func (q *Queries) AdminListRoles(ctx context.Context, arg AdminListRolesParams) ([]AdminRoleRow, error) {
	rows, err := q.db.Query(ctx, adminListRoles,
		arg.CursorTime,
		arg.CursorID,
		arg.Lim,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AdminRoleRow{}
	for rows.Next() {
		var i AdminRoleRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.PermissionTuples,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// -------------------------------------------------------------------------
// AdminGetRoleWithPermissions
// -------------------------------------------------------------------------

const adminGetRoleWithPermissions = `-- name: AdminGetRoleWithPermissions :one
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
GROUP BY r.id
`

func (q *Queries) AdminGetRoleWithPermissions(ctx context.Context, id pgtype.UUID) (AdminRoleRow, error) {
	row := q.db.QueryRow(ctx, adminGetRoleWithPermissions, id)
	var i AdminRoleRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PermissionTuples,
	)
	return i, err
}

// -------------------------------------------------------------------------
// CountUserRolesForRole
// -------------------------------------------------------------------------

const countUserRolesForRole = `-- name: CountUserRolesForRole :one
SELECT COUNT(*)::BIGINT
FROM iam.user_roles
WHERE role_id = $1
`

func (q *Queries) CountUserRolesForRole(ctx context.Context, roleID pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countUserRolesForRole, roleID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
