// grn.go — gRPC handler for OnGRNReceived RPC (BL-FIN-002).
//
// Called when a Goods Receipt Note is received, triggering auto-AP posting:
//   Dr 5001 (COGS/Inventory Expense) / Cr 2001 (AP/Pilgrim Liability).
//
// Idempotent on grn_id. Returns existing entry if already posted.
// GRN processing fails if AP posting fails (atomic transaction).

package grpc_api

import (
	"context"

	"finance-svc/api/grpc_api/pb"
	"finance-svc/service"
	"finance-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// OnGRNReceived handles the OnGRNReceived RPC.
// Posts Dr 5001 (COGS/Inventory Expense) / Cr 2001 (AP/Pilgrim Liability).
// Idempotent on grn_id.
func (s *Server) OnGRNReceived(ctx context.Context, req *pb.OnGRNReceivedRequest) (*pb.OnGRNReceivedResponse, error) {
	const op = "grpc_api.Server.OnGRNReceived"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "OnGRNReceived"),
		attribute.String("grn_id", req.GetGrnId()),
		attribute.String("departure_id", req.GetDepartureId()),
		attribute.Int64("amount_idr", req.GetAmountIdr()),
	)

	logger.Info().
		Str("op", op).
		Str("grn_id", req.GetGrnId()).
		Str("departure_id", req.GetDepartureId()).
		Int64("amount_idr", req.GetAmountIdr()).
		Msg("")

	if req.GetGrnId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "grn_id is required")
	}
	if req.GetDepartureId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "departure_id is required")
	}
	if req.GetAmountIdr() <= 0 {
		return nil, grpcStatus.Errorf(grpcCodes.InvalidArgument, "amount_idr must be positive, got %d", req.GetAmountIdr())
	}

	result, err := s.svc.OnGRNReceived(ctx, &service.OnGRNReceivedParams{
		GrnID:       req.GetGrnId(),
		DepartureID: req.GetDepartureId(),
		AmountIDR:   req.GetAmountIdr(),
	})
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("grn_id", req.GetGrnId()).
			Err(err).
			Msg("OnGRNReceived failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "failed to post GRN journal entry: %v", err)
	}

	logger.Info().
		Str("op", op).
		Str("grn_id", req.GetGrnId()).
		Str("departure_id", req.GetDepartureId()).
		Str("entry_id", result.EntryID).
		Bool("idempotent", result.Idempotent).
		Msg("OnGRNReceived succeeded")

	span.SetStatus(otelCodes.Ok, "success")
	return &pb.OnGRNReceivedResponse{
		EntryId:    result.EntryID,
		Idempotent: result.Idempotent,
	}, nil
}
