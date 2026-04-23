package service

// roles_admin.go — admin role management: ListRoles, CreateRole, UpdateRole,
// DeleteRole, ListPermissions, AssignRoleToUser, RevokeRoleFromUser.
//
// S1-E-06 depth card (BL-IAM-010..017).

import (
	"context"
	"errors"
	"fmt"
	"time"

	"iam-svc/store/postgres_store"
	"iam-svc/store/postgres_store/sqlc"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ── Shared value types ────────────────────────────────────────────────────────

// AdminPermissionDetail is the representation of one permission.
type AdminPermissionDetail struct {
	ID       string
	Resource string
	Action   string
	Scope    string
}

// AdminRoleDetail is a role with its full permission set.
type AdminRoleDetail struct {
	ID          string
	Name        string
	Description string
	Permissions []AdminPermissionDetail
	CreatedAt   time.Time
}

// ── ListRolesAdmin ────────────────────────────────────────────────────────────

type ListRolesAdminParams struct {
	Cursor string // opaque "unix:uuid"; "" = beginning
	Limit  int32
}

type ListRolesAdminResult struct {
	Roles      []AdminRoleDetail
	NextCursor string
}

func (s *Service) ListRolesAdmin(ctx context.Context, params *ListRolesAdminParams) (*ListRolesAdminResult, error) {
	const op = "service.Service.ListRolesAdmin"
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

	var arg sqlc.AdminListRolesParams
	arg.Lim = lim
	if params.Cursor != "" {
		ct, cid, err := decodeCursor(params.Cursor)
		if err != nil {
			e := errors.Join(apperrors.ErrValidation, fmt.Errorf("parse cursor: %w", err))
			span.RecordError(e)
			span.SetStatus(codes.Error, e.Error())
			return nil, e
		}
		arg.CursorTime.Time = ct
		arg.CursorTime.Valid = true
		arg.CursorID = cid
	}

	rows, err := s.store.AdminListRoles(ctx, arg)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Error().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	roles := make([]AdminRoleDetail, 0, len(rows))
	for _, r := range rows {
		roles = append(roles, adminRoleRowToDetail(r))
	}

	nextCursor := ""
	if int32(len(rows)) == lim && len(rows) > 0 {
		last := rows[len(rows)-1]
		nextCursor = encodeCursor(last.CreatedAt.Time, last.ID)
	}

	span.SetStatus(codes.Ok, "success")
	return &ListRolesAdminResult{Roles: roles, NextCursor: nextCursor}, nil
}

// ── CreateRoleAdmin ───────────────────────────────────────────────────────────

type CreateRoleAdminParams struct {
	Name          string
	Description   string
	PermissionIDs []string
}

type CreateRoleAdminResult struct {
	Role AdminRoleDetail
}

func (s *Service) CreateRoleAdmin(ctx context.Context, params *CreateRoleAdminParams) (*CreateRoleAdminResult, error) {
	const op = "service.Service.CreateRoleAdmin"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("name", params.Name))

	if params.Name == "" {
		e := errors.Join(apperrors.ErrValidation, errors.New("name is required"))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	permUUIDs, err := parseUUIDSlice(params.PermissionIDs)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var createdRoleID sqlc.IamRole
	_, err = s.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(q *sqlc.Queries) error {
			role, err := q.CreateRole(ctx, sqlc.CreateRoleParams{
				Name:        params.Name,
				Description: params.Description,
			})
			if err != nil {
				return fmt.Errorf("create role: %w", postgres_store.WrapDBError(err))
			}
			createdRoleID = role

			for _, pid := range permUUIDs {
				if err := q.GrantPermissionToRole(ctx, sqlc.GrantPermissionToRoleParams{
					RoleID:       role.ID,
					PermissionID: pid,
				}); err != nil {
					return fmt.Errorf("grant permission: %w", postgres_store.WrapDBError(err))
				}
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

	// Fetch with permissions for response.
	row, err := s.store.AdminGetRoleWithPermissions(ctx, createdRoleID.ID)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	logger.Info().Str("role_id", uuidToString(createdRoleID.ID)).Str("name", params.Name).Msg("admin: role created")
	return &CreateRoleAdminResult{Role: adminRoleRowToDetail(row)}, nil
}

// ── UpdateRoleAdmin ───────────────────────────────────────────────────────────

type UpdateRoleAdminParams struct {
	ID          string
	Name        string // optional
	Description string // optional
	// PermissionIDsProvided distinguishes nil (omitted) from [] (explicitly empty).
	PermissionIDsProvided bool
	PermissionIDs         []string
}

type UpdateRoleAdminResult struct {
	Role AdminRoleDetail
}

func (s *Service) UpdateRoleAdmin(ctx context.Context, params *UpdateRoleAdminParams) (*UpdateRoleAdminResult, error) {
	const op = "service.Service.UpdateRoleAdmin"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("role_id", params.ID))

	if params.ID == "" {
		e := errors.Join(apperrors.ErrValidation, errors.New("id is required"))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	roleUUID, err := stringToUUID(params.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Load current state.
	current, err := s.store.GetRoleByID(ctx, roleUUID)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	newName := current.Name
	if params.Name != "" {
		newName = params.Name
	}
	newDesc := current.Description
	if params.Description != "" {
		newDesc = params.Description
	}

	var permIDs []pgtype.UUID
	if params.PermissionIDsProvided {
		permIDs, err = parseUUIDSlice(params.PermissionIDs)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}

	_, err = s.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(q *sqlc.Queries) error {
			if _, err := q.UpdateRole(ctx, sqlc.UpdateRoleParams{
				ID:          roleUUID,
				Name:        newName,
				Description: newDesc,
			}); err != nil {
				return fmt.Errorf("update role: %w", postgres_store.WrapDBError(err))
			}

			if params.PermissionIDsProvided {
				if err := q.RevokeAllPermissionsForRole(ctx, roleUUID); err != nil {
					return fmt.Errorf("revoke permissions: %w", postgres_store.WrapDBError(err))
				}
				for _, pid := range permIDs {
					if err := q.GrantPermissionToRole(ctx, sqlc.GrantPermissionToRoleParams{
						RoleID:       roleUUID,
						PermissionID: pid,
					}); err != nil {
						return fmt.Errorf("grant permission: %w", postgres_store.WrapDBError(err))
					}
				}
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

	row, err := s.store.AdminGetRoleWithPermissions(ctx, roleUUID)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	logger.Info().Str("role_id", params.ID).Msg("admin: role updated")
	return &UpdateRoleAdminResult{Role: adminRoleRowToDetail(row)}, nil
}

// ── DeleteRoleAdmin ───────────────────────────────────────────────────────────

type DeleteRoleAdminParams struct {
	ID string
}

type DeleteRoleAdminResult struct{}

func (s *Service) DeleteRoleAdmin(ctx context.Context, params *DeleteRoleAdminParams) (*DeleteRoleAdminResult, error) {
	const op = "service.Service.DeleteRoleAdmin"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("role_id", params.ID))

	if params.ID == "" {
		e := errors.Join(apperrors.ErrValidation, errors.New("id is required"))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	roleUUID, err := stringToUUID(params.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Refuse deletion if any users hold this role.
	count, err := s.store.CountUserRolesForRole(ctx, roleUUID)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	if count > 0 {
		e := errors.Join(apperrors.ErrConflict, fmt.Errorf("role is assigned to %d user(s); revoke before deleting", count))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	if err := s.store.DeleteRole(ctx, roleUUID); err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Error().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	logger.Info().Str("role_id", params.ID).Msg("admin: role deleted")
	return &DeleteRoleAdminResult{}, nil
}

// ── ListPermissionsAdmin ──────────────────────────────────────────────────────

type ListPermissionsAdminParams struct{}

type ListPermissionsAdminResult struct {
	Permissions []AdminPermissionDetail
}

func (s *Service) ListPermissionsAdmin(ctx context.Context, _ *ListPermissionsAdminParams) (*ListPermissionsAdminResult, error) {
	const op = "service.Service.ListPermissionsAdmin"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	rows, err := s.store.ListPermissions(ctx)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		logger.Error().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	perms := make([]AdminPermissionDetail, 0, len(rows))
	for _, p := range rows {
		perms = append(perms, AdminPermissionDetail{
			ID:       uuidToString(p.ID),
			Resource: p.Resource,
			Action:   p.Action,
			Scope:    string(p.Scope),
		})
	}

	span.SetStatus(codes.Ok, "success")
	return &ListPermissionsAdminResult{Permissions: perms}, nil
}

// ── AssignRoleToUserAdmin ─────────────────────────────────────────────────────

type AssignRoleToUserAdminParams struct {
	UserID      string
	RoleID      string
	ActorUserID string // for audit
}

type AssignRoleToUserAdminResult struct{}

func (s *Service) AssignRoleToUserAdmin(ctx context.Context, params *AssignRoleToUserAdminParams) (*AssignRoleToUserAdminResult, error) {
	const op = "service.Service.AssignRoleToUserAdmin"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("user_id", params.UserID),
		attribute.String("role_id", params.RoleID),
	)

	if params.UserID == "" || params.RoleID == "" {
		e := errors.Join(apperrors.ErrValidation, errors.New("user_id and role_id are required"))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	uid, err := stringToUUID(params.UserID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	rid, err := stringToUUID(params.RoleID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Load user for branch_id in audit.
	user, err := s.store.GetUserByID(ctx, uid)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	actorUUID, _ := optionalStringToUUID(params.ActorUserID)

	_, err = s.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(q *sqlc.Queries) error {
			if err := q.AssignRoleToUser(ctx, sqlc.AssignRoleToUserParams{
				UserID: uid,
				RoleID: rid,
			}); err != nil {
				return fmt.Errorf("assign role: %w", postgres_store.WrapDBError(err))
			}
			if _, err := q.InsertAuditLog(ctx, sqlc.InsertAuditLogParams{
				UserID:     actorUUID,
				BranchID:   user.BranchID,
				Resource:   "user_role",
				ResourceID: params.UserID,
				Action:     "assign_role",
				NewValue:   []byte(fmt.Sprintf(`{"role_id":%q}`, params.RoleID)),
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
	logger.Info().Str("user_id", params.UserID).Str("role_id", params.RoleID).Msg("admin: role assigned to user")
	return &AssignRoleToUserAdminResult{}, nil
}

// ── RevokeRoleFromUserAdmin ───────────────────────────────────────────────────

type RevokeRoleFromUserAdminParams struct {
	UserID      string
	RoleID      string
	ActorUserID string
}

type RevokeRoleFromUserAdminResult struct{}

func (s *Service) RevokeRoleFromUserAdmin(ctx context.Context, params *RevokeRoleFromUserAdminParams) (*RevokeRoleFromUserAdminResult, error) {
	const op = "service.Service.RevokeRoleFromUserAdmin"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("user_id", params.UserID),
		attribute.String("role_id", params.RoleID),
	)

	if params.UserID == "" || params.RoleID == "" {
		e := errors.Join(apperrors.ErrValidation, errors.New("user_id and role_id are required"))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	uid, err := stringToUUID(params.UserID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	rid, err := stringToUUID(params.RoleID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	user, err := s.store.GetUserByID(ctx, uid)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	actorUUID, _ := optionalStringToUUID(params.ActorUserID)

	_, err = s.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(q *sqlc.Queries) error {
			if err := q.RevokeRoleFromUser(ctx, sqlc.RevokeRoleFromUserParams{
				UserID: uid,
				RoleID: rid,
			}); err != nil {
				return fmt.Errorf("revoke role: %w", postgres_store.WrapDBError(err))
			}
			if _, err := q.InsertAuditLog(ctx, sqlc.InsertAuditLogParams{
				UserID:     actorUUID,
				BranchID:   user.BranchID,
				Resource:   "user_role",
				ResourceID: params.UserID,
				Action:     "revoke_role",
				OldValue:   []byte(fmt.Sprintf(`{"role_id":%q}`, params.RoleID)),
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
	logger.Info().Str("user_id", params.UserID).Str("role_id", params.RoleID).Msg("admin: role revoked from user")
	return &RevokeRoleFromUserAdminResult{}, nil
}

// ── Internal helpers ──────────────────────────────────────────────────────────

// adminRoleRowToDetail converts a DB row (with packed permission_tuples) to the
// service-layer AdminRoleDetail struct.
func adminRoleRowToDetail(r sqlc.AdminRoleRow) AdminRoleDetail {
	perms := make([]AdminPermissionDetail, 0)
	for _, item := range r.DecodePermissions() {
		perms = append(perms, AdminPermissionDetail{
			ID:       item.ID,
			Resource: item.Resource,
			Action:   item.Action,
			Scope:    item.Scope,
		})
	}
	return AdminRoleDetail{
		ID:          uuidToString(r.ID),
		Name:        r.Name,
		Description: r.Description,
		Permissions: perms,
		CreatedAt:   r.CreatedAt.Time,
	}
}
