package grpc_api

import (
	"catalog-svc/api/grpc_api/pb"
	"catalog-svc/service"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// Server implements pb.CatalogServiceServer by delegating to the service layer.
// Container-level liveness is served by the standard grpc.health.v1.Health
// protocol registered in cmd/server.go; business RPCs (ListPackages,
// GetPackage, GetPackageDeparture) live in packages.go.
type Server struct {
	pb.UnimplementedCatalogServiceServer

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
