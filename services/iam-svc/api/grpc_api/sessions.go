package grpc_api

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

// Login — email/password → access + refresh token pair. Mirrors the legacy
// REST POST /v1/sessions. The gateway forwards user_agent + ip in the request
// message; iam-svc trusts those fields verbatim (single-point auth per ADR 0009).
func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	const op = "grpc_api.Server.Login"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "Login"),
		attribute.String("email", req.GetEmail()),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.Login(ctx, &service.LoginParams{
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		TOTPCode:  req.GetTotpCode(),
		UserAgent: req.GetUserAgent(),
		IP:        parseIPOrNil(req.GetIp()),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.LoginResponse{
		AccessToken:           result.AccessToken,
		RefreshToken:          result.RefreshToken,
		AccessExpiresAtUnix:   result.AccessExpiresAt.Unix(),
		RefreshExpiresAtUnix:  result.RefreshExpiresAt.Unix(),
		User:                  userProfileToProto(result.User),
	}, nil
}

// RefreshSession — rotate a refresh token. Mirrors POST /v1/sessions/refresh.
func (s *Server) RefreshSession(ctx context.Context, req *pb.RefreshSessionRequest) (*pb.RefreshSessionResponse, error) {
	const op = "grpc_api.Server.RefreshSession"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "RefreshSession"))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.RefreshSession(ctx, &service.RefreshSessionParams{
		RefreshToken: req.GetRefreshToken(),
		UserAgent:    req.GetUserAgent(),
		IP:           parseIPOrNil(req.GetIp()),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.RefreshSessionResponse{
		AccessToken:          result.AccessToken,
		RefreshToken:         result.RefreshToken,
		AccessExpiresAtUnix:  result.AccessExpiresAt.Unix(),
		RefreshExpiresAtUnix: result.RefreshExpiresAt.Unix(),
	}, nil
}

// Logout — revoke the session row identified by session_id (the gateway
// extracts it from the validated bearer payload and forwards here). Mirrors
// DELETE /v1/sessions.
func (s *Server) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	const op = "grpc_api.Server.Logout"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "Logout"),
		attribute.String("session_id", req.GetSessionId()),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetSessionId() == "" {
		err := errors.Join(apperrors.ErrValidation, errors.New("session_id is required"))
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	if _, err := s.svc.Logout(ctx, &service.LogoutParams{SessionID: req.GetSessionId()}); err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.LogoutResponse{}, nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// parseIPOrNil mirrors the REST-side clientIP behavior: empty or malformed
// input yields a nil *netip.Addr so iam.sessions.ip stays NULL.
func parseIPOrNil(raw string) *netip.Addr {
	if raw == "" {
		return nil
	}
	addr, err := netip.ParseAddr(raw)
	if err != nil {
		return nil
	}
	return &addr
}

// userProfileToProto maps the service-layer UserProfile to the proto envelope.
func userProfileToProto(u service.UserProfile) *pb.UserProfile {
	return &pb.UserProfile{
		UserId:   u.UserID,
		Email:    u.Email,
		Name:     u.Name,
		BranchId: u.BranchID,
		Status:   u.Status,
	}
}
