package service

// users_admin.go — admin user management: ListUsers, CreateUser, UpdateUser,
// GetUser, ResetUserPassword.
//
// S1-E-06 depth card (BL-IAM-005..009).

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"iam-svc/store/postgres_store"
	"iam-svc/store/postgres_store/sqlc"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/crypto/bcrypt"
)

// ── Shared value types ────────────────────────────────────────────────────────

// AdminUserSummary is the lightweight user record returned from ListUsers.
type AdminUserSummary struct {
	ID          string
	Email       string
	Name        string
	BranchID    string
	Status      string
	Roles       []string
	LastLoginAt *time.Time // nil if never logged in
	CreatedAt   time.Time
}

// AdminUserDetail is the full user record returned from GetUser / CreateUser /
// UpdateUser. Roles are role names.
type AdminUserDetail struct {
	ID        string
	Email     string
	Name      string
	BranchID  string
	Status    string
	Roles     []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ── ListUsers ─────────────────────────────────────────────────────────────────

type ListUsersParams struct {
	Status   string // optional: active|suspended|pending; "" = all
	BranchID string // optional UUID string; "" = all
	Cursor   string // opaque "unix:uuid"; "" = beginning
	Limit    int32  // default 20
}

type ListUsersResult struct {
	Users      []AdminUserSummary
	NextCursor string // "" = no more pages
}

func (s *Service) ListUsers(ctx context.Context, params *ListUsersParams) (*ListUsersResult, error) {
	const op = "service.Service.ListUsers"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	lim := params.Limit
	if lim <= 0 {
		lim = 20
	}
	if lim > 100 {
		lim = 100
	}

	// Parse optional branch_id filter.
	var branchFilter pgtype.UUID
	if params.BranchID != "" {
		parsed, err := stringToUUID(params.BranchID)
		if err != nil {
			logger.Warn().Err(err).Msg("")
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		branchFilter = parsed
	} else {
		// zero UUID sentinel = no filter (matched by SQL: $2 = '00000000-...' check)
		branchFilter = pgtype.UUID{Valid: true} // all zeros
	}

	// Validate optional status filter.
	statusFilter := ""
	if params.Status != "" {
		switch sqlc.IamUserStatus(params.Status) {
		case sqlc.IamUserStatusActive, sqlc.IamUserStatusSuspended, sqlc.IamUserStatusPending:
			statusFilter = params.Status
		default:
			err := errors.Join(apperrors.ErrValidation, fmt.Errorf("unknown status %q", params.Status))
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}

	// Decode cursor.
	var cursorTime pgtype.Timestamptz
	var cursorID pgtype.UUID
	if params.Cursor != "" {
		ct, cid, err := decodeAdminListCursor(params.Cursor)
		if err != nil {
			e := errors.Join(apperrors.ErrValidation, fmt.Errorf("parse cursor: %w", err))
			span.RecordError(e)
			span.SetStatus(codes.Error, e.Error())
			return nil, e
		}
		cursorTime = pgtype.Timestamptz{Time: ct, Valid: true}
		cursorID = cid
	}

	rows, err := s.store.AdminListUsers(ctx, sqlc.AdminListUsersParams{
		StatusFilter: statusFilter,
		BranchFilter: branchFilter,
		CursorTime:   cursorTime,
		CursorID:     cursorID,
		Lim:          lim,
	})
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Error().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	users := make([]AdminUserSummary, 0, len(rows))
	for _, r := range rows {
		u := AdminUserSummary{
			ID:        uuidToString(r.ID),
			Email:     r.Email,
			Name:      r.Name,
			BranchID:  uuidToString(r.BranchID),
			Status:    string(r.Status),
			Roles:     r.RoleNames,
			CreatedAt: r.CreatedAt.Time,
		}
		if r.LastLoginAt.Valid {
			t := r.LastLoginAt.Time
			u.LastLoginAt = &t
		}
		users = append(users, u)
	}

	nextCursor := ""
	if int32(len(rows)) == lim && len(rows) > 0 {
		last := rows[len(rows)-1]
		nextCursor = encodeAdminListCursor(last.CreatedAt.Time, last.ID)
	}

	span.SetStatus(codes.Ok, "success")
	return &ListUsersResult{Users: users, NextCursor: nextCursor}, nil
}

// ── CreateUserAdmin ───────────────────────────────────────────────────────────

type CreateUserAdminParams struct {
	Email    string
	Name     string
	Password string
	BranchID string
	RoleIDs  []string // optional
}

type CreateUserAdminResult struct {
	User AdminUserDetail
}

func (s *Service) CreateUserAdmin(ctx context.Context, params *CreateUserAdminParams) (*CreateUserAdminResult, error) {
	const op = "service.Service.CreateUserAdmin"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("email", params.Email))

	if params.Email == "" || params.Name == "" || params.Password == "" || params.BranchID == "" {
		e := errors.Join(apperrors.ErrValidation, errors.New("email, name, password, branch_id are required"))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	branchUUID, err := stringToUUID(params.BranchID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Validate role IDs upfront.
	roleUUIDs, err := parseUUIDSlice(params.RoleIDs)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Hash password.
	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		e := errors.Join(apperrors.ErrInternal, fmt.Errorf("hash password: %w", err))
		logger.Error().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	var createdUser sqlc.IamUser
	var roleNames []string

	_, err = s.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(q *sqlc.Queries) error {
			u, err := q.CreateUser(ctx, sqlc.CreateUserParams{
				Email:        params.Email,
				PasswordHash: string(hash),
				Name:         params.Name,
				BranchID:     branchUUID,
				Status:       sqlc.IamUserStatusActive,
			})
			if err != nil {
				return fmt.Errorf("create user: %w", postgres_store.WrapDBError(err))
			}
			createdUser = u

			// Assign roles.
			names := make([]string, 0, len(roleUUIDs))
			for _, rid := range roleUUIDs {
				if err := q.AssignRoleToUser(ctx, sqlc.AssignRoleToUserParams{
					UserID: u.ID,
					RoleID: rid,
				}); err != nil {
					return fmt.Errorf("assign role: %w", postgres_store.WrapDBError(err))
				}
				// Fetch role name for response.
				role, err := q.GetRoleByID(ctx, rid)
				if err != nil {
					return fmt.Errorf("get role: %w", postgres_store.WrapDBError(err))
				}
				names = append(names, role.Name)
			}

			// Audit log.
			if _, err := q.InsertAuditLog(ctx, sqlc.InsertAuditLogParams{
				BranchID:   branchUUID,
				Resource:   "user",
				ResourceID: uuidToString(u.ID),
				Action:     "create",
				NewValue:   []byte(fmt.Sprintf(`{"email":%q,"name":%q}`, params.Email, params.Name)),
			}); err != nil {
				return fmt.Errorf("audit: %w", postgres_store.WrapDBError(err))
			}

			roleNames = names
			return nil
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	logger.Info().Str("user_id", uuidToString(createdUser.ID)).Str("email", params.Email).Msg("admin: user created")

	return &CreateUserAdminResult{
		User: AdminUserDetail{
			ID:        uuidToString(createdUser.ID),
			Email:     createdUser.Email,
			Name:      createdUser.Name,
			BranchID:  uuidToString(createdUser.BranchID),
			Status:    string(createdUser.Status),
			Roles:     roleNames,
			CreatedAt: createdUser.CreatedAt.Time,
			UpdatedAt: createdUser.UpdatedAt.Time,
		},
	}, nil
}

// ── UpdateUser ────────────────────────────────────────────────────────────────

type UpdateUserParams struct {
	ID      string
	Name    string   // optional; empty = keep current
	Status  string   // optional; empty = keep current
	RoleIDs []string // optional; nil = keep current; len-0 = remove all
	// RoleIDsProvided distinguishes nil (omitted) from [] (explicitly empty).
	RoleIDsProvided bool
}

type UpdateUserResult struct {
	User AdminUserDetail
}

func (s *Service) UpdateUser(ctx context.Context, params *UpdateUserParams) (*UpdateUserResult, error) {
	const op = "service.Service.UpdateUser"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("user_id", params.ID))

	if params.ID == "" {
		e := errors.Join(apperrors.ErrValidation, errors.New("id is required"))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	uid, err := stringToUUID(params.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Load current user to fill in unchanged fields.
	current, err := s.store.GetUserByID(ctx, uid)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	// Resolve final name / status.
	newName := current.Name
	if params.Name != "" {
		newName = params.Name
	}
	newStatus := current.Status
	if params.Status != "" {
		switch sqlc.IamUserStatus(params.Status) {
		case sqlc.IamUserStatusActive, sqlc.IamUserStatusSuspended, sqlc.IamUserStatusPending:
			newStatus = sqlc.IamUserStatus(params.Status)
		default:
			e := errors.Join(apperrors.ErrValidation, fmt.Errorf("unknown status %q", params.Status))
			span.RecordError(e)
			span.SetStatus(codes.Error, e.Error())
			return nil, e
		}
	}

	var roleUUIDs []pgtype.UUID
	if params.RoleIDsProvided {
		roleUUIDs, err = parseUUIDSlice(params.RoleIDs)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}

	oldJSON := []byte(fmt.Sprintf(`{"name":%q,"status":%q}`, current.Name, string(current.Status)))
	newJSON := []byte(fmt.Sprintf(`{"name":%q,"status":%q}`, newName, string(newStatus)))

	var updatedRow sqlc.AdminUpdateUserNameStatusRow
	var roleNames []string

	_, err = s.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(q *sqlc.Queries) error {
			row, err := q.AdminUpdateUserNameStatus(ctx, sqlc.AdminUpdateUserNameStatusParams{
				ID:     uid,
				Name:   newName,
				Status: newStatus,
			})
			if err != nil {
				return fmt.Errorf("update user: %w", postgres_store.WrapDBError(err))
			}
			updatedRow = row

			if params.RoleIDsProvided {
				// Replace role set.
				if err := q.RevokeAllRolesForUser(ctx, uid); err != nil {
					return fmt.Errorf("revoke roles: %w", postgres_store.WrapDBError(err))
				}
				names := make([]string, 0, len(roleUUIDs))
				for _, rid := range roleUUIDs {
					if err := q.AssignRoleToUser(ctx, sqlc.AssignRoleToUserParams{
						UserID: uid,
						RoleID: rid,
					}); err != nil {
						return fmt.Errorf("assign role: %w", postgres_store.WrapDBError(err))
					}
					role, err := q.GetRoleByID(ctx, rid)
					if err != nil {
						return fmt.Errorf("get role: %w", postgres_store.WrapDBError(err))
					}
					names = append(names, role.Name)
				}
				roleNames = names
			} else {
				// Load current role names.
				names, err := q.ListRoleNamesForUser(ctx, uid)
				if err != nil {
					return fmt.Errorf("list roles: %w", postgres_store.WrapDBError(err))
				}
				roleNames = names
			}

			// Audit.
			if _, err := q.InsertAuditLog(ctx, sqlc.InsertAuditLogParams{
				BranchID:   row.BranchID,
				Resource:   "user",
				ResourceID: params.ID,
				Action:     "update",
				OldValue:   oldJSON,
				NewValue:   newJSON,
			}); err != nil {
				return fmt.Errorf("audit: %w", postgres_store.WrapDBError(err))
			}
			return nil
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	logger.Info().Str("user_id", params.ID).Msg("admin: user updated")

	return &UpdateUserResult{
		User: AdminUserDetail{
			ID:        uuidToString(updatedRow.ID),
			Email:     updatedRow.Email,
			Name:      updatedRow.Name,
			BranchID:  uuidToString(updatedRow.BranchID),
			Status:    string(updatedRow.Status),
			Roles:     roleNames,
			CreatedAt: updatedRow.CreatedAt.Time,
			UpdatedAt: updatedRow.UpdatedAt.Time,
		},
	}, nil
}

// ── GetUser ───────────────────────────────────────────────────────────────────

type GetUserParams struct {
	ID string
}

type GetUserResult struct {
	User AdminUserDetail
}

func (s *Service) GetUser(ctx context.Context, params *GetUserParams) (*GetUserResult, error) {
	const op = "service.Service.GetUser"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("user_id", params.ID))

	if params.ID == "" {
		e := errors.Join(apperrors.ErrValidation, errors.New("id is required"))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	uid, err := stringToUUID(params.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	row, err := s.store.AdminGetUser(ctx, uid)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &GetUserResult{
		User: AdminUserDetail{
			ID:        uuidToString(row.ID),
			Email:     row.Email,
			Name:      row.Name,
			BranchID:  uuidToString(row.BranchID),
			Status:    string(row.Status),
			Roles:     row.RoleNames,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		},
	}, nil
}

// ── ResetUserPassword ─────────────────────────────────────────────────────────

type ResetUserPasswordParams struct {
	ID          string
	NewPassword string
	ActorUserID string // for audit; optional
}

type ResetUserPasswordResult struct{}

func (s *Service) ResetUserPassword(ctx context.Context, params *ResetUserPasswordParams) (*ResetUserPasswordResult, error) {
	const op = "service.Service.ResetUserPassword"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("user_id", params.ID))

	if params.ID == "" || params.NewPassword == "" {
		e := errors.Join(apperrors.ErrValidation, errors.New("id and new_password are required"))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	uid, err := stringToUUID(params.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Confirm user exists first.
	user, err := s.store.GetUserByID(ctx, uid)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		e := errors.Join(apperrors.ErrInternal, fmt.Errorf("hash password: %w", err))
		logger.Error().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	actorUUID, _ := optionalStringToUUID(params.ActorUserID)

	_, err = s.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(q *sqlc.Queries) error {
			if err := q.UpdateUserPasswordHash(ctx, sqlc.UpdateUserPasswordHashParams{
				ID:           uid,
				PasswordHash: string(hash),
			}); err != nil {
				return fmt.Errorf("update password: %w", postgres_store.WrapDBError(err))
			}

			// Revoke all active sessions for this user.
			if err := q.RevokeAllSessionsForUserAdmin(ctx, uid); err != nil {
				return fmt.Errorf("revoke sessions: %w", postgres_store.WrapDBError(err))
			}

			// Audit.
			if _, err := q.InsertAuditLog(ctx, sqlc.InsertAuditLogParams{
				UserID:     actorUUID,
				BranchID:   user.BranchID,
				Resource:   "user",
				ResourceID: params.ID,
				Action:     "reset_password",
				NewValue:   []byte(`{"action":"password_reset"}`),
			}); err != nil {
				return fmt.Errorf("audit: %w", postgres_store.WrapDBError(err))
			}
			return nil
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	logger.Info().Str("user_id", params.ID).Msg("admin: user password reset; all sessions revoked")
	return &ResetUserPasswordResult{}, nil
}

// ── Cursor helpers ────────────────────────────────────────────────────────────

// encodeCursor produces an opaque "unix:uuid" string for pagination.
func encodeAdminListCursor(t time.Time, id pgtype.UUID) string {
	return fmt.Sprintf("%d:%s", t.UnixMicro(), uuidToString(id))
}

// decodeCursor splits an opaque cursor into its (time, uuid) components.
func decodeAdminListCursor(cursor string) (time.Time, pgtype.UUID, error) {
	idx := strings.Index(cursor, ":")
	if idx < 0 {
		return time.Time{}, pgtype.UUID{}, fmt.Errorf("invalid cursor format")
	}
	us, err := strconv.ParseInt(cursor[:idx], 10, 64)
	if err != nil {
		return time.Time{}, pgtype.UUID{}, fmt.Errorf("parse timestamp: %w", err)
	}
	t := time.UnixMicro(us)
	uid, err := stringToUUID(cursor[idx+1:])
	if err != nil {
		return time.Time{}, pgtype.UUID{}, fmt.Errorf("parse uuid: %w", err)
	}
	return t, uid, nil
}

// parseUUIDSlice converts a slice of UUID strings to pgtype.UUID values.
func parseUUIDSlice(ids []string) ([]pgtype.UUID, error) {
	out := make([]pgtype.UUID, 0, len(ids))
	for _, id := range ids {
		u, err := stringToUUID(id)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}
