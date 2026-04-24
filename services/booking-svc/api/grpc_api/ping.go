package grpc_api

import (
	"context"

	"booking-svc/api/grpc_api/pb"
	"booking-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Ping is the scaffold RPC used by the gateway's permission-gate smoke tests.
// Identity-agnostic per ADR 0009: the gateway validates the bearer and the
// permission tuple upstream; this handler returns a trivial {message: "ok"}
// proving the gateway→backend hop works. No service-layer call — there is
// no business logic to dispatch to.
func (s *Server) Ping(ctx context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	const op = "grpc_api.Server.Ping"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "Ping"))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	span.SetStatus(codes.Ok, "success")
	return &pb.PingResponse{Message: "ok"}, nil
}
