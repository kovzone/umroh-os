// Code generated manually (sqlc pattern) for users_admin.sql.
// S1-E-06 depth card — admin user management queries.

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// -------------------------------------------------------------------------
// AdminListUsers
// -------------------------------------------------------------------------

const adminListUsers = `-- name: AdminListUsers :many
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
LIMIT $5
`

type AdminListUsersParams struct {
	StatusFilter string             `json:"status_filter"`
	BranchFilter pgtype.UUID        `json:"branch_filter"`
	CursorTime   pgtype.Timestamptz `json:"cursor_time"`
	CursorID     pgtype.UUID        `json:"cursor_id"`
	Lim          int32              `json:"lim"`
}

type AdminListUsersRow struct {
	ID          pgtype.UUID        `json:"id"`
	Email       string             `json:"email"`
	Name        string             `json:"name"`
	BranchID    pgtype.UUID        `json:"branch_id"`
	Status      IamUserStatus      `json:"status"`
	LastLoginAt pgtype.Timestamptz `json:"last_login_at"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	RoleNames   []string           `json:"role_names"`
}

func (q *Queries) AdminListUsers(ctx context.Context, arg AdminListUsersParams) ([]AdminListUsersRow, error) {
	rows, err := q.db.Query(ctx, adminListUsers,
		arg.StatusFilter,
		arg.BranchFilter,
		arg.CursorTime,
		arg.CursorID,
		arg.Lim,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AdminListUsersRow{}
	for rows.Next() {
		var i AdminListUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.Name,
			&i.BranchID,
			&i.Status,
			&i.LastLoginAt,
			&i.CreatedAt,
			&i.RoleNames,
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
// AdminGetUser
// -------------------------------------------------------------------------

const adminGetUser = `-- name: AdminGetUser :one
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
GROUP BY u.id
`

type AdminGetUserRow struct {
	ID          pgtype.UUID        `json:"id"`
	Email       string             `json:"email"`
	Name        string             `json:"name"`
	BranchID    pgtype.UUID        `json:"branch_id"`
	Status      IamUserStatus      `json:"status"`
	LastLoginAt pgtype.Timestamptz `json:"last_login_at"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
	RoleNames   []string           `json:"role_names"`
}

func (q *Queries) AdminGetUser(ctx context.Context, id pgtype.UUID) (AdminGetUserRow, error) {
	row := q.db.QueryRow(ctx, adminGetUser, id)
	var i AdminGetUserRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.BranchID,
		&i.Status,
		&i.LastLoginAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.RoleNames,
	)
	return i, err
}

// -------------------------------------------------------------------------
// AdminUpdateUserNameStatus
// -------------------------------------------------------------------------

const adminUpdateUserNameStatus = `-- name: AdminUpdateUserNameStatus :one
UPDATE iam.users
SET name       = $2,
    status     = $3,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, email, name, branch_id, status, last_login_at, created_at, updated_at
`

type AdminUpdateUserNameStatusParams struct {
	ID     pgtype.UUID   `json:"id"`
	Name   string        `json:"name"`
	Status IamUserStatus `json:"status"`
}

type AdminUpdateUserNameStatusRow struct {
	ID          pgtype.UUID        `json:"id"`
	Email       string             `json:"email"`
	Name        string             `json:"name"`
	BranchID    pgtype.UUID        `json:"branch_id"`
	Status      IamUserStatus      `json:"status"`
	LastLoginAt pgtype.Timestamptz `json:"last_login_at"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) AdminUpdateUserNameStatus(ctx context.Context, arg AdminUpdateUserNameStatusParams) (AdminUpdateUserNameStatusRow, error) {
	row := q.db.QueryRow(ctx, adminUpdateUserNameStatus, arg.ID, arg.Name, arg.Status)
	var i AdminUpdateUserNameStatusRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.BranchID,
		&i.Status,
		&i.LastLoginAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

// -------------------------------------------------------------------------
// RevokeAllSessionsForUserAdmin
// -------------------------------------------------------------------------

const revokeAllSessionsForUserAdmin = `-- name: RevokeAllSessionsForUserAdmin :exec
UPDATE iam.sessions
SET revoked_at = now()
WHERE user_id = $1 AND revoked_at IS NULL
`

func (q *Queries) RevokeAllSessionsForUserAdmin(ctx context.Context, userID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, revokeAllSessionsForUserAdmin, userID)
	return err
}
