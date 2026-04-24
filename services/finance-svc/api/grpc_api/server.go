package grpc_api

import (
	"finance-svc/api/grpc_api/pb"
	"finance-svc/service"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// Server implements pb.FinanceServiceServer by delegating to the service
// layer. Container-level liveness is served by the standard
// grpc.health.v1.Health protocol registered in cmd/server.go; the one
// business RPC (FinancePing) lives in finance.go.
type Server struct {
	pb.UnimplementedFinanceServiceServer

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
