// admin.go — gateway adapter methods for iam-svc admin RPCs (Wave 1C / Phase 6).
//
// Covers: user management (ListUsers, CreateUser, UpdateUser, GetUser,
// ResetUserPassword) and role management (ListRoles, CreateRole, UpdateRole,
// DeleteRole, ListPermissions, AssignRoleToUser, RevokeRoleFromUser).
//
// Each method follows the span + log + mapIamError pattern from validate_token.go.
package iam_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/iam_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Adapter-local result types
// ---------------------------------------------------------------------------

// AdminPermissionResult is the gateway-local representation of a permission.
type AdminPermissionResult struct {
	ID       string
	Resource string
	Action   string
	Scope    string
}

// AdminRoleResult is the gateway-local representation of a role.
type AdminRoleResult struct {
	ID          string
	Name        string
	Description string
	Permissions []*AdminPermissionResult
	CreatedAt   int64
}

// AdminUserSummaryResult is the gateway-local summary for list responses.
type AdminUserSummaryResult struct {
	ID          string
	Email       string
	Name        string
	BranchID    string
	Status      string
	Roles       []string
	LastLoginAt int64
	CreatedAt   int64
}

// AdminUserResult is the gateway-local full user record.
type AdminUserResult struct {
	ID        string
	Email     string
	Name      string
	BranchID  string
	Status    string
	Roles     []string
	CreatedAt int64
	UpdatedAt int64
}

// ListUsersResult is the paginated list result for users.
type ListUsersResult struct {
	Users      []*AdminUserSummaryResult
	NextCursor string
}

// ListRolesResult is the paginated list result for roles.
type ListRolesResult struct {
	Roles      []*AdminRoleResult
	NextCursor string
}

// ListPermissionsResult holds all permissions.
type ListPermissionsResult struct {
	Permissions []*AdminPermissionResult
}

// ---------------------------------------------------------------------------
// Params types
// ---------------------------------------------------------------------------

// ListUsersParams is the input for ListUsers.
type ListUsersParams struct {
	Status   string
	BranchID string
	Cursor   string
	Limit    int32
}

// CreateUserParams is the input for CreateUser.
type CreateUserParams struct {
	Email    string
	Name     string
	Password string
	BranchID string
	RoleIDs  []string
}

// UpdateUserParams is the input for UpdateUser.
type UpdateUserParams struct {
	Name    string
	Status  string
	RoleIDs []string
}

// CreateRoleParams is the input for CreateRole.
type CreateRoleParams struct {
	Name          string
	Description   string
	PermissionIDs []string
}

// UpdateRoleParams is the input for UpdateRole.
type UpdateRoleParams struct {
	Name          string
	Description   string
	PermissionIDs []string
}

// ---------------------------------------------------------------------------
// Helper mappers
// ---------------------------------------------------------------------------

func fromProtoAdminPermission(p *pb.AdminPermission) *AdminPermissionResult {
	if p == nil {
		return nil
	}
	return &AdminPermissionResult{
		ID:       p.GetId(),
		Resource: p.GetResource(),
		Action:   p.GetAction(),
		Scope:    p.GetScope(),
	}
}

func fromProtoAdminRole(r *pb.AdminRole) *AdminRoleResult {
	if r == nil {
		return nil
	}
	perms := make([]*AdminPermissionResult, 0, len(r.GetPermissions()))
	for _, p := range r.GetPermissions() {
		perms = append(perms, fromProtoAdminPermission(p))
	}
	return &AdminRoleResult{
		ID:          r.GetId(),
		Name:        r.GetName(),
		Description: r.GetDescription(),
		Permissions: perms,
		CreatedAt:   r.GetCreatedAt(),
	}
}

func fromProtoAdminUserSummary(u *pb.AdminUserSummary) *AdminUserSummaryResult {
	if u == nil {
		return nil
	}
	return &AdminUserSummaryResult{
		ID:          u.GetId(),
		Email:       u.GetEmail(),
		Name:        u.GetName(),
		BranchID:    u.GetBranchId(),
		Status:      u.GetStatus(),
		Roles:       u.GetRoles(),
		LastLoginAt: u.GetLastLoginAt(),
		CreatedAt:   u.GetCreatedAt(),
	}
}

func fromProtoAdminUser(u *pb.AdminUser) *AdminUserResult {
	if u == nil {
		return nil
	}
	return &AdminUserResult{
		ID:        u.GetId(),
		Email:     u.GetEmail(),
		Name:      u.GetName(),
		BranchID:  u.GetBranchId(),
		Status:    u.GetStatus(),
		Roles:     u.GetRoles(),
		CreatedAt: u.GetCreatedAt(),
		UpdatedAt: u.GetUpdatedAt(),
	}
}

// ---------------------------------------------------------------------------
// User management methods
// ---------------------------------------------------------------------------

// ListUsers returns a paginated list of IAM users.
func (a *Adapter) ListUsers(ctx context.Context, params *ListUsersParams) (*ListUsersResult, error) {
	const op = "iam_grpc_adapter.Adapter.ListUsers"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.adminClient.ListUsers(ctx, &pb.ListUsersRequest{
		Status:   params.Status,
		BranchId: params.BranchID,
		Cursor:   params.Cursor,
		Limit:    params.Limit,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.ListUsers failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	users := make([]*AdminUserSummaryResult, 0, len(resp.GetUsers()))
	for _, u := range resp.GetUsers() {
		users = append(users, fromProtoAdminUserSummary(u))
	}

	span.SetStatus(codes.Ok, "ok")
	return &ListUsersResult{
		Users:      users,
		NextCursor: resp.GetNextCursor(),
	}, nil
}

// CreateUser creates a new IAM user.
func (a *Adapter) CreateUser(ctx context.Context, params *CreateUserParams) (*AdminUserResult, error) {
	const op = "iam_grpc_adapter.Adapter.CreateUser"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("email", params.Email))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.adminClient.CreateUser(ctx, &pb.CreateUserRequest{
		Email:    params.Email,
		Name:     params.Name,
		Password: params.Password,
		BranchId: params.BranchID,
		RoleIds:  params.RoleIDs,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.CreateUser failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return fromProtoAdminUser(resp.GetUser()), nil
}

// UpdateUser updates an existing IAM user.
func (a *Adapter) UpdateUser(ctx context.Context, id string, params *UpdateUserParams) (*AdminUserResult, error) {
	const op = "iam_grpc_adapter.Adapter.UpdateUser"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("user_id", id))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.adminClient.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:      id,
		Name:    params.Name,
		Status:  params.Status,
		RoleIds: params.RoleIDs,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.UpdateUser failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return fromProtoAdminUser(resp.GetUser()), nil
}

// GetUser returns the full record of a single IAM user.
func (a *Adapter) GetUser(ctx context.Context, id string) (*AdminUserResult, error) {
	const op = "iam_grpc_adapter.Adapter.GetUser"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("user_id", id))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.adminClient.GetUser(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.GetUser failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return fromProtoAdminUser(resp.GetUser()), nil
}

// ResetUserPassword resets the password for an IAM user.
func (a *Adapter) ResetUserPassword(ctx context.Context, id, newPassword string) error {
	const op = "iam_grpc_adapter.Adapter.ResetUserPassword"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("user_id", id))

	logger := logging.LogWithTrace(ctx, a.logger)

	_, err := a.adminClient.ResetUserPassword(ctx, &pb.ResetUserPasswordRequest{
		Id:          id,
		NewPassword: newPassword,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.ResetUserPassword failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return nil
}

// ---------------------------------------------------------------------------
// Role management methods
// ---------------------------------------------------------------------------

// ListRoles returns a paginated list of IAM roles.
func (a *Adapter) ListRoles(ctx context.Context) (*ListRolesResult, error) {
	const op = "iam_grpc_adapter.Adapter.ListRoles"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.adminClient.ListRoles(ctx, &pb.ListRolesRequest{})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.ListRoles failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	roles := make([]*AdminRoleResult, 0, len(resp.GetRoles()))
	for _, r := range resp.GetRoles() {
		roles = append(roles, fromProtoAdminRole(r))
	}

	span.SetStatus(codes.Ok, "ok")
	return &ListRolesResult{
		Roles:      roles,
		NextCursor: resp.GetNextCursor(),
	}, nil
}

// CreateRole creates a new IAM role.
func (a *Adapter) CreateRole(ctx context.Context, params *CreateRoleParams) (*AdminRoleResult, error) {
	const op = "iam_grpc_adapter.Adapter.CreateRole"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.adminClient.CreateRole(ctx, &pb.CreateRoleRequest{
		Name:          params.Name,
		Description:   params.Description,
		PermissionIds: params.PermissionIDs,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.CreateRole failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return fromProtoAdminRole(resp.GetRole()), nil
}

// UpdateRole updates an existing IAM role.
func (a *Adapter) UpdateRole(ctx context.Context, id string, params *UpdateRoleParams) (*AdminRoleResult, error) {
	const op = "iam_grpc_adapter.Adapter.UpdateRole"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("role_id", id))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.adminClient.UpdateRole(ctx, &pb.UpdateRoleRequest{
		Id:            id,
		Name:          params.Name,
		Description:   params.Description,
		PermissionIds: params.PermissionIDs,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.UpdateRole failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return fromProtoAdminRole(resp.GetRole()), nil
}

// DeleteRole deletes an IAM role.
func (a *Adapter) DeleteRole(ctx context.Context, id string) error {
	const op = "iam_grpc_adapter.Adapter.DeleteRole"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("role_id", id))

	logger := logging.LogWithTrace(ctx, a.logger)

	_, err := a.adminClient.DeleteRole(ctx, &pb.DeleteRoleRequest{Id: id})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.DeleteRole failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return nil
}

// ListPermissions returns all available IAM permissions.
func (a *Adapter) ListPermissions(ctx context.Context) (*ListPermissionsResult, error) {
	const op = "iam_grpc_adapter.Adapter.ListPermissions"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.adminClient.ListPermissions(ctx, &pb.ListPermissionsRequest{})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.ListPermissions failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	perms := make([]*AdminPermissionResult, 0, len(resp.GetPermissions()))
	for _, p := range resp.GetPermissions() {
		perms = append(perms, fromProtoAdminPermission(p))
	}

	span.SetStatus(codes.Ok, "ok")
	return &ListPermissionsResult{Permissions: perms}, nil
}

// AssignRoleToUser assigns a role to a user.
func (a *Adapter) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	const op = "iam_grpc_adapter.Adapter.AssignRoleToUser"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("user_id", userID), attribute.String("role_id", roleID))

	logger := logging.LogWithTrace(ctx, a.logger)

	_, err := a.adminClient.AssignRoleToUser(ctx, &pb.AssignRoleToUserRequest{
		UserId: userID,
		RoleId: roleID,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.AssignRoleToUser failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return nil
}

// RevokeRoleFromUser revokes a role from a user.
func (a *Adapter) RevokeRoleFromUser(ctx context.Context, userID, roleID string) error {
	const op = "iam_grpc_adapter.Adapter.RevokeRoleFromUser"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("user_id", userID), attribute.String("role_id", roleID))

	logger := logging.LogWithTrace(ctx, a.logger)

	_, err := a.adminClient.RevokeRoleFromUser(ctx, &pb.RevokeRoleFromUserRequest{
		UserId: userID,
		RoleId: roleID,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.RevokeRoleFromUser failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return nil
}
