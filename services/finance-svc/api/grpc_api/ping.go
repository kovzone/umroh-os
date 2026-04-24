package grpc_api

import (
	"context"

	"finance-svc/api/grpc_api/pb"
	"finance-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Ping is the scaffold RPC used by the gateway's permission-gate smoke tests
// (BL-IAM-002 / BL-IAM-019). Authorization is enforced at the gateway via
// RequireBearerToken + RequirePermission("journal_entry", "read", "global")
// before this RPC is invoked. finance-svc is identity-agnostic per ADR 0009:
// the gateway's REST handler assembles the client-visible envelope from its
// own locals, so this method returns only a trivial "ok" — proving
// reachability is the full value of the backend call. If finance is down
// the gRPC transport error maps to ErrServiceUnavailable at the gateway
// adapter.
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
