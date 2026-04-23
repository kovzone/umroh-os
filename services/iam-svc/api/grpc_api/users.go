package grpc_api

// users.go — gRPC handler implementations for admin user management RPCs.
//
// Implements pb.IamAdminHandler (partial) for user-related methods.
// S1-E-06 depth card (BL-IAM-005..009).

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

// ListUsers returns a cursor-paginated list of users with optional filters.
func (s *Server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	const op = "grpc_api.Server.ListUsers"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.ListUsers(ctx, &service.ListUsersParams{
		Status:   req.GetStatus(),
		BranchID: req.GetBranchId(),
		Cursor:   req.GetCursor(),
		Limit:    req.GetLimit(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	users := make([]*pb.AdminUserSummary, 0, len(result.Users))
	for _, u := range result.Users {
		pu := &pb.AdminUserSummary{
			Id:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			BranchId:  u.BranchID,
			Status:    u.Status,
			Roles:     u.Roles,
			CreatedAt: u.CreatedAt.Unix(),
		}
		if u.LastLoginAt != nil {
			pu.LastLoginAt = u.LastLoginAt.Unix()
		}
		users = append(users, pu)
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.ListUsersResponse{
		Users:      users,
		NextCursor: result.NextCursor,
	}, nil
}

// CreateUser creates a new user with optional role assignments.
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	const op = "grpc_api.Server.CreateUser"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("email", req.GetEmail()))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetEmail() == "" || req.GetName() == "" || req.GetPassword() == "" || req.GetBranchId() == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("email, name, password, branch_id are required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	result, err := s.svc.CreateUserAdmin(ctx, &service.CreateUserAdminParams{
		Email:    req.GetEmail(),
		Name:     req.GetName(),
		Password: req.GetPassword(),
		BranchID: req.GetBranchId(),
		RoleIDs:  req.GetRoleIds(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.CreateUserResponse{
		User: adminUserDetailToProto(result.User),
	}, nil
}

// UpdateUser updates a user's name, status, and/or role assignments.
func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	const op = "grpc_api.Server.UpdateUser"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("user_id", req.GetId()))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetId() == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("id is required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	params := &service.UpdateUserParams{
		ID:     req.GetId(),
		Name:   req.GetName(),
		Status: req.GetStatus(),
	}
	if req.RoleIds != nil {
		params.RoleIDsProvided = true
		params.RoleIDs = req.GetRoleIds()
	}

	result, err := s.svc.UpdateUser(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.UpdateUserResponse{
		User: adminUserDetailToProto(result.User),
	}, nil
}

// GetUser returns a single user's full profile including their role names.
func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	const op = "grpc_api.Server.GetUser"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("user_id", req.GetId()))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetId() == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("id is required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	result, err := s.svc.GetUser(ctx, &service.GetUserParams{ID: req.GetId()})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.GetUserResponse{
		User: adminUserDetailToProto(result.User),
	}, nil
}

// ResetUserPassword sets a new password and revokes all active sessions.
func (s *Server) ResetUserPassword(ctx context.Context, req *pb.ResetUserPasswordRequest) (*pb.ResetUserPasswordResponse, error) {
	const op = "grpc_api.Server.ResetUserPassword"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("user_id", req.GetId()))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetId() == "" || req.GetNewPassword() == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("id and new_password are required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	_, err := s.svc.ResetUserPassword(ctx, &service.ResetUserPasswordParams{
		ID:          req.GetId(),
		NewPassword: req.GetNewPassword(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.ResetUserPasswordResponse{Ok: true}, nil
}

// ── Proto conversion helpers ──────────────────────────────────────────────────

func adminUserDetailToProto(u service.AdminUserDetail) *pb.AdminUser {
	return &pb.AdminUser{
		Id:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		BranchId:  u.BranchID,
		Status:    u.Status,
		Roles:     u.Roles,
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	}
}
