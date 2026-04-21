package grpc_api

import (
	"context"

	"iam-svc/api/grpc_api/pb"
	"iam-svc/service"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/status"
)

// Server implements pb.IamServiceServer by delegating to the service layer.
type Server struct {
	pb.UnimplementedIamServiceServer

	logger *zerolog.Logger
	tracer trace.Tracer

	svc service.IService
}

// NewServer constructs a gRPC server bound to the given service.
func NewServer(logger *zerolog.Logger, tracer trace.Tracer, svc service.IService) *Server {
	return &Server{
		logger: logger,
		tracer: tracer,
		svc:    svc,
	}
}

// Healthz is the pilot placeholder RPC kept from the scaffold.
// Real health checks go through the standard gRPC health protocol.
func (s *Server) Healthz(ctx context.Context, _ *pb.HealthzRequest) (*pb.HealthzResponse, error) {
	const op = "grpc_api.Server.Healthz"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	return &pb.HealthzResponse{Ok: true}, nil
}

// ValidateToken verifies a bearer access token and returns the caller's
// identity + current role-name snapshot. Token failures fail-closed as
// Unauthenticated per F1 "never default to allow".
func (s *Server) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	const op = "grpc_api.Server.ValidateToken"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.ValidateToken(ctx, &service.ValidateTokenParams{
		AccessToken: req.GetAccessToken(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.ValidateTokenResponse{
		UserId:        result.UserID,
		BranchId:      result.BranchID,
		SessionId:     result.SessionID,
		Roles:         result.Roles,
		ExpiresAtUnix: result.ExpiresAt.Unix(),
	}, nil
}

// CheckPermission resolves (resource, action, scope) against the user's grants.
// Returns allowed=false on the response (not PermissionDenied) when the grant
// is absent — consumers map that to HTTP 403 themselves. Malformed input maps
// to InvalidArgument; DB failures map to Internal.
func (s *Server) CheckPermission(ctx context.Context, req *pb.CheckPermissionRequest) (*pb.CheckPermissionResponse, error) {
	const op = "grpc_api.Server.CheckPermission"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("resource", req.GetResource()),
		attribute.String("action", req.GetAction()),
		attribute.String("scope", req.GetScope()),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.CheckPermission(ctx, &service.CheckPermissionParams{
		UserID:   req.GetUserId(),
		Resource: req.GetResource(),
		Action:   req.GetAction(),
		Scope:    req.GetScope(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetAttributes(attribute.Bool("allowed", result.Allowed))
	span.SetStatus(codes.Ok, "success")
	return &pb.CheckPermissionResponse{Allowed: result.Allowed}, nil
}
