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
// Container-level liveness is served by the standard grpc.health.v1.Health
// protocol registered in cmd/server.go — no placeholder Healthz RPC is
// needed here.
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
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
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
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetAttributes(attribute.Bool("allowed", result.Allowed))
	span.SetStatus(codes.Ok, "success")
	return &pb.CheckPermissionResponse{Allowed: result.Allowed}, nil
}

// RecordAudit inserts one row into iam.audit_logs for a state-changing action.
// Validation failures (missing resource/action, malformed UUID or IP) map to
// InvalidArgument; DB failures map to Internal. Success returns the inserted
// row's id + timestamp so callers can cite the entry without a second query.
func (s *Server) RecordAudit(ctx context.Context, req *pb.RecordAuditRequest) (*pb.RecordAuditResponse, error) {
	const op = "grpc_api.Server.RecordAudit"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("resource", req.GetResource()),
		attribute.String("action", req.GetAction()),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.RecordAudit(ctx, &service.RecordAuditParams{
		ActorUserID: req.GetUserId(),
		BranchID:    req.GetBranchId(),
		Resource:    req.GetResource(),
		ResourceID:  req.GetResourceId(),
		Action:      req.GetAction(),
		OldValue:    req.GetOldValue(),
		NewValue:    req.GetNewValue(),
		IP:          req.GetIp(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), apperrors.GRPCMessage(err))
	}

	span.SetAttributes(attribute.String("audit_log_id", result.AuditLogID))
	span.SetStatus(codes.Ok, "success")
	return &pb.RecordAuditResponse{
		AuditLogId:    result.AuditLogID,
		CreatedAtUnix: result.CreatedAt.Unix(),
	}, nil
}
