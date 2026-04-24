package grpc_api

import (
	"ops-svc/api/grpc_api/pb"
	"ops-svc/service"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// Server implements pb.OpsServiceServer by delegating to the service layer.
// No business RPCs have been implemented yet — this scaffold exists so the
// service registers a gRPC listener and passes docker-compose healthchecks via
// the standard grpc.health.v1.Health protocol. Real RPCs land with the first
// operations feature slice.
type Server struct {
	pb.UnimplementedOpsServiceServer

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
