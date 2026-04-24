package grpc_api

import (
	"payment-svc/api/grpc_api/pb"
	"payment-svc/service"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// Server implements pb.PaymentServiceServer by delegating to the service layer.
// No business RPCs have been implemented yet — this scaffold exists so the
// service registers a gRPC listener and passes docker-compose healthchecks via
// the standard grpc.health.v1.Health protocol. Real RPCs land with the first
// payment feature slice.
type Server struct {
	pb.UnimplementedPaymentServiceServer

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
