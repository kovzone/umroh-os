package grpc_api

import (
	"context"

	"catalog-svc/adapter/iam_grpc_adapter"
	"catalog-svc/api/grpc_api/pb"
	"catalog-svc/service"
	"catalog-svc/util/logging"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IamClient is the slice of iam-svc that catalog-svc's gRPC handlers use.
// Defined as an interface so tests can inject a mock without pulling in gRPC.
type IamClient interface {
	ValidateToken(ctx context.Context, params *iam_grpc_adapter.ValidateTokenParams) (*iam_grpc_adapter.ValidateTokenResult, error)
	CheckPermission(ctx context.Context, params *iam_grpc_adapter.CheckPermissionParams) (*iam_grpc_adapter.CheckPermissionResult, error)
	RecordAudit(ctx context.Context, params *iam_grpc_adapter.RecordAuditParams) (*iam_grpc_adapter.RecordAuditResult, error)
}

// Server implements pb.CatalogServiceServer by delegating to the service layer.
type Server struct {
	pb.UnimplementedCatalogServiceServer

	logger *zerolog.Logger
	tracer trace.Tracer

	svc       service.IService
	iamClient IamClient // nil → permission gate skipped (legacy / test mode)
}

// NewServer constructs a gRPC server bound to the given service and IAM adapter.
// Pass nil for iamClient to skip permission gating (useful in tests).
func NewServer(logger *zerolog.Logger, tracer trace.Tracer, svc service.IService, iamClient IamClient) *Server {
	return &Server{
		logger:    logger,
		tracer:    tracer,
		svc:       svc,
		iamClient: iamClient,
	}
}

// Healthz is the pilot placeholder RPC.
// Real RPCs (ValidateToken, CheckPermission, GetUser, RecordAudit) land in F1.7.
func (s *Server) Healthz(ctx context.Context, _ *pb.HealthzRequest) (*pb.HealthzResponse, error) {
	const op = "grpc_api.Server.Healthz"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	return &pb.HealthzResponse{Ok: true}, nil
}
