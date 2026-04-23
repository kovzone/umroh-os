// proxy_iam_admin.go — gateway REST handlers for IAM admin user/role management
// (Wave 1C / Phase 6).
//
// Route topology (all bearer-protected):
//   GET    /v1/users                           → ListUsers
//   POST   /v1/users                           → CreateUser
//   GET    /v1/users/:id                       → GetUser
//   PUT    /v1/users/:id                       → UpdateUser
//   POST   /v1/users/:id/reset-password        → ResetUserPassword
//   GET    /v1/roles                           → ListRoles
//   POST   /v1/roles                           → CreateRole
//   PUT    /v1/roles/:id                       → UpdateRole
//   DELETE /v1/roles/:id                       → DeleteRole
//   GET    /v1/permissions                     → ListPermissions
//   POST   /v1/users/:id/roles/:role_id        → AssignRoleToUser
//   DELETE /v1/users/:id/roles/:role_id        → RevokeRoleFromUser
//
// Per ADR-0009: gateway is the single REST entry-point; iam-svc is pure gRPC.
package rest_oapi

import (
	"errors"

	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ---------------------------------------------------------------------------
// Request / response body types
// ---------------------------------------------------------------------------

// CreateUserAdminBody is the JSON body for POST /v1/users.
type CreateUserAdminBody struct {
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Password string   `json:"password"`
	BranchID string   `json:"branch_id,omitempty"`
	RoleIDs  []string `json:"role_ids,omitempty"`
}

// UpdateUserAdminBody is the JSON body for PUT /v1/users/:id.
type UpdateUserAdminBody struct {
	Name    string   `json:"name,omitempty"`
	Status  string   `json:"status,omitempty"`
	RoleIDs []string `json:"role_ids,omitempty"`
}

// ResetPasswordBody is the JSON body for POST /v1/users/:id/reset-password.
type ResetPasswordBody struct {
	NewPassword string `json:"new_password"`
}

// AdminUserResponseData is the JSON shape for a full user.
type AdminUserResponseData struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	BranchID  string   `json:"branch_id,omitempty"`
	Status    string   `json:"status"`
	Roles     []string `json:"roles"`
	CreatedAt int64    `json:"created_at"`
	UpdatedAt int64    `json:"updated_at"`
}

// AdminUserSummaryData is the JSON shape for a user summary.
type AdminUserSummaryData struct {
	ID          string   `json:"id"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	BranchID    string   `json:"branch_id,omitempty"`
	Status      string   `json:"status"`
	Roles       []string `json:"roles"`
	LastLoginAt int64    `json:"last_login_at"`
	CreatedAt   int64    `json:"created_at"`
}

// AdminRoleResponseData is the JSON shape for a role.
type AdminRoleResponseData struct {
	ID          string                        `json:"id"`
	Name        string                        `json:"name"`
	Description string                        `json:"description,omitempty"`
	Permissions []AdminPermissionResponseData `json:"permissions"`
	CreatedAt   int64                         `json:"created_at"`
}

// AdminPermissionResponseData is the JSON shape for a permission.
type AdminPermissionResponseData struct {
	ID       string `json:"id"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Scope    string `json:"scope"`
}

// CreateRoleAdminBody is the JSON body for POST /v1/roles.
type CreateRoleAdminBody struct {
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	PermissionIDs []string `json:"permission_ids,omitempty"`
}

// UpdateRoleAdminBody is the JSON body for PUT /v1/roles/:id.
type UpdateRoleAdminBody struct {
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	PermissionIDs []string `json:"permission_ids,omitempty"`
}

// ---------------------------------------------------------------------------
// Mapper helpers
// ---------------------------------------------------------------------------

func mapAdminUserResult(u *iam_grpc_adapter.AdminUserResult) AdminUserResponseData {
	return AdminUserResponseData{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		BranchID:  u.BranchID,
		Status:    u.Status,
		Roles:     u.Roles,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func mapAdminUserSummary(u *iam_grpc_adapter.AdminUserSummaryResult) AdminUserSummaryData {
	return AdminUserSummaryData{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		BranchID:    u.BranchID,
		Status:      u.Status,
		Roles:       u.Roles,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
	}
}

func mapAdminRoleResult(r *iam_grpc_adapter.AdminRoleResult) AdminRoleResponseData {
	perms := make([]AdminPermissionResponseData, 0, len(r.Permissions))
	for _, p := range r.Permissions {
		perms = append(perms, AdminPermissionResponseData{
			ID:       p.ID,
			Resource: p.Resource,
			Action:   p.Action,
			Scope:    p.Scope,
		})
	}
	return AdminRoleResponseData{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: perms,
		CreatedAt:   r.CreatedAt,
	}
}

// ---------------------------------------------------------------------------
// User handlers
// ---------------------------------------------------------------------------

// ListUsers — GET /v1/users (bearer)
func (s *Server) ListUsers(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListUsers"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	statusFilter := c.Query("status")
	branchID := c.Query("branch_id")
	cursor := c.Query("cursor")
	limit := int32(c.QueryInt("limit", 20))

	result, err := s.svc.ListUsers(ctx, &iam_grpc_adapter.ListUsersParams{
		Status:   statusFilter,
		BranchID: branchID,
		Cursor:   cursor,
		Limit:    limit,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	data := make([]AdminUserSummaryData, 0, len(result.Users))
	for _, u := range result.Users {
		data = append(data, mapAdminUserSummary(u))
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":        data,
		"next_cursor": result.NextCursor,
	})
}

// CreateUser — POST /v1/users (bearer)
func (s *Server) CreateUser(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateUser"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	var body CreateUserAdminBody
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.CreateUser(ctx, &iam_grpc_adapter.CreateUserParams{
		Email:    body.Email,
		Name:     body.Name,
		Password: body.Password,
		BranchID: body.BranchID,
		RoleIDs:  body.RoleIDs,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": mapAdminUserResult(result)})
}

// GetUser — GET /v1/users/:id (bearer)
func (s *Server) GetUser(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.GetUser"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("user_id", id))

	result, err := s.svc.GetUser(ctx, id)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": mapAdminUserResult(result)})
}

// UpdateUser — PUT /v1/users/:id (bearer)
func (s *Server) UpdateUser(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.UpdateUser"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("user_id", id))

	var body UpdateUserAdminBody
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.UpdateUser(ctx, id, &iam_grpc_adapter.UpdateUserParams{
		Name:    body.Name,
		Status:  body.Status,
		RoleIDs: body.RoleIDs,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": mapAdminUserResult(result)})
}

// ResetUserPassword — POST /v1/users/:id/reset-password (bearer)
func (s *Server) ResetUserPassword(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.ResetUserPassword"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("user_id", id))

	var body ResetPasswordBody
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	if err := s.svc.ResetUserPassword(ctx, id, body.NewPassword); err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"ok": true})
}

// ---------------------------------------------------------------------------
// Role handlers
// ---------------------------------------------------------------------------

// ListRoles — GET /v1/roles (bearer)
func (s *Server) ListRoles(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListRoles"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	result, err := s.svc.ListRoles(ctx)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	data := make([]AdminRoleResponseData, 0, len(result.Roles))
	for _, r := range result.Roles {
		data = append(data, mapAdminRoleResult(r))
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":        data,
		"next_cursor": result.NextCursor,
	})
}

// CreateRole — POST /v1/roles (bearer)
func (s *Server) CreateRole(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateRole"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	var body CreateRoleAdminBody
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.CreateRole(ctx, &iam_grpc_adapter.CreateRoleParams{
		Name:          body.Name,
		Description:   body.Description,
		PermissionIDs: body.PermissionIDs,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": mapAdminRoleResult(result)})
}

// UpdateRole — PUT /v1/roles/:id (bearer)
func (s *Server) UpdateRole(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.UpdateRole"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("role_id", id))

	var body UpdateRoleAdminBody
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.UpdateRole(ctx, id, &iam_grpc_adapter.UpdateRoleParams{
		Name:          body.Name,
		Description:   body.Description,
		PermissionIDs: body.PermissionIDs,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": mapAdminRoleResult(result)})
}

// DeleteRole — DELETE /v1/roles/:id (bearer)
func (s *Server) DeleteRole(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.DeleteRole"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("role_id", id))

	if err := s.svc.DeleteRole(ctx, id); err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "deleted")
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// ListPermissions — GET /v1/permissions (bearer)
func (s *Server) ListPermissions(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListPermissions"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	result, err := s.svc.ListPermissions(ctx)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	data := make([]AdminPermissionResponseData, 0, len(result.Permissions))
	for _, p := range result.Permissions {
		data = append(data, AdminPermissionResponseData{
			ID:       p.ID,
			Resource: p.Resource,
			Action:   p.Action,
			Scope:    p.Scope,
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": data})
}

// AssignRoleToUser — POST /v1/users/:id/roles/:role_id (bearer)
func (s *Server) AssignRoleToUser(c *fiber.Ctx, userID, roleID string) error {
	const op = "rest_oapi.Server.AssignRoleToUser"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(
		attribute.String("user_id", userID),
		attribute.String("role_id", roleID),
	)

	if err := s.svc.AssignRoleToUser(ctx, userID, roleID); err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"ok": true})
}

// RevokeRoleFromUser — DELETE /v1/users/:id/roles/:role_id (bearer)
func (s *Server) RevokeRoleFromUser(c *fiber.Ctx, userID, roleID string) error {
	const op = "rest_oapi.Server.RevokeRoleFromUser"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(
		attribute.String("user_id", userID),
		attribute.String("role_id", roleID),
	)

	if err := s.svc.RevokeRoleFromUser(ctx, userID, roleID); err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "deleted")
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeIamAdminError(c *fiber.Ctx, span trace.Span, err error) error {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	var httpStatus int
	var code, message string

	switch {
	case errors.Is(err, apperrors.ErrValidation):
		httpStatus = fiber.StatusBadRequest
		code = "validation_error"
		message = err.Error()
	case errors.Is(err, apperrors.ErrNotFound):
		httpStatus = fiber.StatusNotFound
		code = "not_found"
		message = err.Error()
	case errors.Is(err, apperrors.ErrConflict):
		httpStatus = fiber.StatusConflict
		code = "conflict"
		message = err.Error()
	case errors.Is(err, apperrors.ErrUnauthorized):
		httpStatus = fiber.StatusUnauthorized
		code = "unauthorized"
		message = "autentikasi diperlukan"
	case errors.Is(err, apperrors.ErrForbidden):
		httpStatus = fiber.StatusForbidden
		code = "forbidden"
		message = "akses ditolak"
	case errors.Is(err, apperrors.ErrServiceUnavailable):
		httpStatus = fiber.StatusBadGateway
		code = "service_unavailable"
		message = "layanan IAM sementara tidak tersedia"
	default:
		httpStatus = fiber.StatusInternalServerError
		code = "internal_error"
		message = "terjadi kesalahan tidak terduga"
	}

	return c.Status(httpStatus).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": message,
		},
	})
}
