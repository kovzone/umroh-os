package grpc_api

// roles.go — gRPC handler implementations for admin role management RPCs.
//
// Implements pb.IamAdminHandler (partial) for role-related methods:
// ListRoles, CreateRole, UpdateRole, DeleteRole, ListPermissions,
// AssignRoleToUser, RevokeRoleFromUser.
//
// S1-E-06 depth card (BL-IAM-010..017).

import (
	"context"
	"errors"

	"iam-svc/api/grpc_api/pb"
	"iam-svc/service"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

// ListRoles returns a cursor-paginated list of roles with their permissions.
func (s *Server) ListRoles(ctx context.Context, req *pb.ListRolesRequest) (*pb.ListRolesResponse, error) {
	const op = "grpc_api.Server.ListRoles"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.ListRolesAdmin(ctx, &service.ListRolesAdminParams{
		Cursor: req.Cursor,
		Limit:  req.Limit,
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	roles := make([]*pb.AdminRole, 0, len(result.Roles))
	for _, r := range result.Roles {
		roles = append(roles, adminRoleDetailToProto(r))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.ListRolesResponse{
		Roles:      roles,
		NextCursor: result.NextCursor,
	}, nil
}

// CreateRole creates a new role with optional permission grants.
func (s *Server) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.CreateRoleResponse, error) {
	const op = "grpc_api.Server.CreateRole"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("name", req.Name))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.Name == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("name is required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	result, err := s.svc.CreateRoleAdmin(ctx, &service.CreateRoleAdminParams{
		Name:          req.Name,
		Description:   req.Description,
		PermissionIDs: req.PermissionIds,
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.CreateRoleResponse{
		Role: adminRoleDetailToProto(result.Role),
	}, nil
}

// UpdateRole updates a role's name, description, and/or permission set.
func (s *Server) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.UpdateRoleResponse, error) {
	const op = "grpc_api.Server.UpdateRole"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("role_id", req.Id))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.Id == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("id is required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	params := &service.UpdateRoleAdminParams{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
	}
	if req.PermissionIds != nil {
		params.PermissionIDsProvided = true
		params.PermissionIDs = req.PermissionIds
	}

	result, err := s.svc.UpdateRoleAdmin(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.UpdateRoleResponse{
		Role: adminRoleDetailToProto(result.Role),
	}, nil
}

// DeleteRole deletes a role if no users currently hold it.
func (s *Server) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*pb.DeleteRoleResponse, error) {
	const op = "grpc_api.Server.DeleteRole"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("role_id", req.Id))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.Id == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("id is required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	_, err := s.svc.DeleteRoleAdmin(ctx, &service.DeleteRoleAdminParams{
		ID: req.Id,
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.DeleteRoleResponse{Ok: true}, nil
}

// ListPermissions returns all permissions defined in the system.
func (s *Server) ListPermissions(ctx context.Context, _ *pb.ListPermissionsRequest) (*pb.ListPermissionsResponse, error) {
	const op = "grpc_api.Server.ListPermissions"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.ListPermissionsAdmin(ctx, &service.ListPermissionsAdminParams{})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	perms := make([]*pb.AdminPermission, 0, len(result.Permissions))
	for _, p := range result.Permissions {
		perms = append(perms, &pb.AdminPermission{
			Id:       p.ID,
			Resource: p.Resource,
			Action:   p.Action,
			Scope:    p.Scope,
		})
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.ListPermissionsResponse{Permissions: perms}, nil
}

// AssignRoleToUser grants a role to a user.
func (s *Server) AssignRoleToUser(ctx context.Context, req *pb.AssignRoleToUserRequest) (*pb.AssignRoleToUserResponse, error) {
	const op = "grpc_api.Server.AssignRoleToUser"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("user_id", req.UserId),
		attribute.String("role_id", req.RoleId),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.UserId == "" || req.RoleId == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("user_id and role_id are required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	_, err := s.svc.AssignRoleToUserAdmin(ctx, &service.AssignRoleToUserAdminParams{
		UserID: req.UserId,
		RoleID: req.RoleId,
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.AssignRoleToUserResponse{Ok: true}, nil
}

// RevokeRoleFromUser removes a role from a user.
func (s *Server) RevokeRoleFromUser(ctx context.Context, req *pb.RevokeRoleFromUserRequest) (*pb.RevokeRoleFromUserResponse, error) {
	const op = "grpc_api.Server.RevokeRoleFromUser"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("user_id", req.UserId),
		attribute.String("role_id", req.RoleId),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.UserId == "" || req.RoleId == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("user_id and role_id are required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	_, err := s.svc.RevokeRoleFromUserAdmin(ctx, &service.RevokeRoleFromUserAdminParams{
		UserID: req.UserId,
		RoleID: req.RoleId,
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.RevokeRoleFromUserResponse{Ok: true}, nil
}

// ── Proto conversion helpers ──────────────────────────────────────────────────

func adminRoleDetailToProto(r service.AdminRoleDetail) *pb.AdminRole {
	perms := make([]*pb.AdminPermission, 0, len(r.Permissions))
	for _, p := range r.Permissions {
		perms = append(perms, &pb.AdminPermission{
			Id:       p.ID,
			Resource: p.Resource,
			Action:   p.Action,
			Scope:    p.Scope,
		})
	}
	return &pb.AdminRole{
		Id:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: perms,
		CreatedAt:   r.CreatedAt.Unix(),
	}
}
