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
		Status:   req.Status,
		BranchID: req.BranchId,
		Cursor:   req.Cursor,
		Limit:    req.Limit,
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
	span.SetAttributes(attribute.String("operation", op), attribute.String("email", req.Email))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.Email == "" || req.Name == "" || req.Password == "" || req.BranchId == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("email, name, password, branch_id are required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	result, err := s.svc.CreateUserAdmin(ctx, &service.CreateUserAdminParams{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		BranchID: req.BranchId,
		RoleIDs:  req.RoleIds,
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
	span.SetAttributes(attribute.String("operation", op), attribute.String("user_id", req.Id))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.Id == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("id is required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	params := &service.UpdateUserParams{
		ID:     req.Id,
		Name:   req.Name,
		Status: req.Status,
	}
	if req.RoleIds != nil {
		params.RoleIDsProvided = true
		params.RoleIDs = req.RoleIds
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
	span.SetAttributes(attribute.String("operation", op), attribute.String("user_id", req.Id))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.Id == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("id is required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	result, err := s.svc.GetUser(ctx, &service.GetUserParams{ID: req.Id})
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
	span.SetAttributes(attribute.String("operation", op), attribute.String("user_id", req.Id))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.Id == "" || req.NewPassword == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("id and new_password are required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	_, err := s.svc.ResetUserPassword(ctx, &service.ResetUserPasswordParams{
		ID:          req.Id,
		NewPassword: req.NewPassword,
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
