package grpc_api

import (
	"context"
	"fmt"

	"catalog-svc/api/grpc_api/pb"
	"catalog-svc/service"
	"catalog-svc/util/apperrors"
	"catalog-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

// DiagnosticsDbTx is the state-mutating WithTx reference path for catalog-svc.
// Inserts one row into public.diagnostics, reads it back inside the same
// transaction, and echoes the message. Dev-only verification; invoke via
// grpcurl. Per ADR 0009 no gateway REST route proxies this — iam-svc's
// /v1/iam/system/diagnostics/db-tx already carries the cross-service trace
// demonstration (S0-J-05).
func (s *Server) DiagnosticsDbTx(ctx context.Context, req *pb.DiagnosticsDbTxRequest) (*pb.DiagnosticsDbTxResponse, error) {
	const op = "grpc_api.Server.DiagnosticsDbTx"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "DiagnosticsDbTx"))

	logger := logging.LogWithTrace(ctx, s.logger)

	message := req.GetMessage()
	if message == "" {
		message = "no message"
	}

	logger.Info().Str("op", op).Str("message", message).Msg("")

	result, err := s.svc.DbTxDiagnostic(ctx, &service.DbTxDiagnosticParams{Message: message})
	if err != nil {
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetAttributes(attribute.String("output.response", fmt.Sprintf("diagnostic_id=%d message=%s", result.ID, message)))
	span.SetStatus(codes.Ok, "success")

	return &pb.DiagnosticsDbTxResponse{
		DiagnosticId: result.ID,
		Message:      fmt.Sprintf("received '%s' from client", message),
	}, nil
}
