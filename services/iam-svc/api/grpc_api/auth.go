package grpc_api

// auth.go — gRPC handler implementations for the client-facing auth RPCs added
// in BL-IAM-018 / S1-E-12. These mirror the former REST surface (sessions.go,
// me.go, admin.go in api/rest_oapi/) so the gateway can proxy IAM REST routes
// to iam-svc over gRPC per ADR 0009.
//
// Each handler follows the same span + log + error-map pattern as the existing
// ValidateToken / CheckPermission / RecordAudit handlers in server.go.
//
// The auth REST surface (api/rest_oapi/) is retired in S1-E-12.

import (
	"context"
	"errors"
	"net/netip"

	"iam-svc/api/grpc_api/pb"
	"iam-svc/service"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

// Login exchanges email + password for a PASETO access + opaque refresh token pair.
// Public — no bearer required. gateway forwards the client IP + User-Agent.
func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	const op = "grpc_api.Server.Login"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("email", req.GetEmail()))

	logger := logging.LogWithTrace(ctx, s.logger)

	var ip *netip.Addr
	if raw := req.GetIp(); raw != "" {
		parsed, err := netip.ParseAddr(raw)
		if err == nil {
			ip = &parsed
		}
	}

	result, err := s.svc.Login(ctx, &service.LoginParams{
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		TOTPCode:  req.GetTotpCode(),
		UserAgent: req.GetUserAgent(),
		IP:        ip,
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.LoginResponse{
		AccessToken:      result.AccessToken,
		RefreshToken:     result.RefreshToken,
		AccessExpiresAt:  result.AccessExpiresAt.Unix(),
		RefreshExpiresAt: result.RefreshExpiresAt.Unix(),
		User:             userProfileToProto(result.User),
	}, nil
}

// RefreshSession rotates the refresh token and issues a new access + refresh pair.
// Public — no bearer required.
func (s *Server) RefreshSession(ctx context.Context, req *pb.RefreshSessionRequest) (*pb.RefreshSessionResponse, error) {
	const op = "grpc_api.Server.RefreshSession"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	var ip *netip.Addr
	if raw := req.GetIp(); raw != "" {
		parsed, err := netip.ParseAddr(raw)
		if err == nil {
			ip = &parsed
		}
	}

	result, err := s.svc.RefreshSession(ctx, &service.RefreshSessionParams{
		RefreshToken: req.GetRefreshToken(),
		UserAgent:    req.GetUserAgent(),
		IP:           ip,
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.RefreshSessionResponse{
		AccessToken:      result.AccessToken,
		RefreshToken:     result.RefreshToken,
		AccessExpiresAt:  result.AccessExpiresAt.Unix(),
		RefreshExpiresAt: result.RefreshExpiresAt.Unix(),
	}, nil
}

// Logout revokes the session identified by session_id. gateway validates the
// bearer via ValidateToken first and passes the session_id claim here.
// Idempotent — a second call on an already-revoked session succeeds silently.
func (s *Server) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	const op = "grpc_api.Server.Logout"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetSessionId() == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("session_id is required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	_, err := s.svc.Logout(ctx, &service.LogoutParams{SessionID: req.GetSessionId()})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.LogoutResponse{Ok: true}, nil
}

// GetMe returns the current user's profile for the given user_id (from the bearer token claim).
func (s *Server) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.GetMeResponse, error) {
	const op = "grpc_api.Server.GetMe"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("user_id", req.GetUserId()))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.GetMe(ctx, &service.GetMeParams{UserID: req.GetUserId()})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.GetMeResponse{
		User:         userProfileToProto(result.User),
		TotpEnrolled: result.TOTPEnrolled,
		TotpVerified: result.TOTPVerified,
	}, nil
}

// EnrollTOTP generates a TOTP secret + otpauth URL for the user.
// Returns AlreadyExists (mapped to 409) if TOTP is already verified.
func (s *Server) EnrollTOTP(ctx context.Context, req *pb.EnrollTOTPRequest) (*pb.EnrollTOTPResponse, error) {
	const op = "grpc_api.Server.EnrollTOTP"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("user_id", req.GetUserId()))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.EnrollTOTP(ctx, &service.EnrollTOTPParams{UserID: req.GetUserId()})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.EnrollTOTPResponse{
		Secret:     result.Secret,
		OtpauthUrl: result.OtpauthURL,
	}, nil
}

// VerifyTOTP validates a TOTP code and stamps users.totp_verified_at on success.
func (s *Server) VerifyTOTP(ctx context.Context, req *pb.VerifyTOTPRequest) (*pb.VerifyTOTPResponse, error) {
	const op = "grpc_api.Server.VerifyTOTP"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("user_id", req.GetUserId()))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.VerifyTOTP(ctx, &service.VerifyTOTPParams{
		UserID: req.GetUserId(),
		Code:   req.GetCode(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.VerifyTOTPResponse{
		VerifiedAtUnix: result.VerifiedAt.Unix(),
	}, nil
}

// SuspendUser flips the target user's status to suspended + revokes all sessions.
// gateway validates the bearer first (via ValidateToken) and then CheckPermission
// for iam.users/suspend/global before forwarding actor_user_id + target_user_id here.
func (s *Server) SuspendUser(ctx context.Context, req *pb.SuspendUserRequest) (*pb.SuspendUserResponse, error) {
	const op = "grpc_api.Server.SuspendUser"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("actor_user_id", req.GetActorUserId()),
		attribute.String("target_user_id", req.GetTargetUserId()),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.SuspendUser(ctx, &service.SuspendUserParams{
		ActorUserID:  req.GetActorUserId(),
		TargetUserID: req.GetTargetUserId(),
	})
	if err != nil {
		logger.Warn().Err(err).
			Str("actor_user_id", req.GetActorUserId()).
			Str("target_user_id", req.GetTargetUserId()).
			Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.SuspendUserResponse{
		User: userProfileToProto(result.User),
	}, nil
}

// userProfileToProto converts the service-layer UserProfile to the proto type.
func userProfileToProto(u service.UserProfile) *pb.UserProfile {
	return &pb.UserProfile{
		UserId:   u.UserID,
		Email:    u.Email,
		Name:     u.Name,
		BranchId: u.BranchID,
		Status:   u.Status,
	}
}
