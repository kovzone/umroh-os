package grpc_api

import (
	"context"

	"visa-svc/api/grpc_api/pb"
	"visa-svc/service"
	"visa-svc/util/logging"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// Server implements pb.VisaServiceServer by delegating to the service layer.
type Server struct {
	pb.UnimplementedVisaServiceServer

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
