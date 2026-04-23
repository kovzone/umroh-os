// revenue.go — gRPC handler for RecognizeRevenue RPC (Wave 1B / BL-FIN-006).
//
// Called by catalog-svc (or manually by ops) after a departure transitions
// to `departed` or `completed`.
//
// Posts a double-entry journal:
//   Dr 2001 (Pilgrim Liability) — total_amount_idr
//   Cr 4001 (Revenue)           — total_amount_idr
//
// Idempotent on departure_id.

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

// RecognizeRevenue handles the RecognizeRevenue RPC.
// Posts Dr 2001 / Cr 4001 journal entry for the given departure.
func (s *Server) RecognizeRevenue(ctx context.Context, req *pb.RecognizeRevenueRequest) (*pb.RecognizeRevenueResponse, error) {
	const op = "grpc_api.Server.RecognizeRevenue"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("departure_id", req.GetDepartureId()),
		attribute.Int64("total_amount_idr", req.GetTotalAmountIdr()),
	)

	logger.Info().
		Str("op", op).
		Str("departure_id", req.GetDepartureId()).
		Int64("total_amount_idr", req.GetTotalAmountIdr()).
		Msg("")

	if req.GetDepartureId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "departure_id is required")
	}
	if req.GetTotalAmountIdr() <= 0 {
		return nil, grpcStatus.Errorf(grpcCodes.InvalidArgument, "total_amount_idr must be positive, got %d", req.GetTotalAmountIdr())
	}

	result, err := s.svc.RecognizeRevenue(ctx, &service.RecognizeRevenueParams{
		DepartureID:    req.GetDepartureId(),
		TotalAmountIDR: req.GetTotalAmountIdr(),
	})
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("departure_id", req.GetDepartureId()).
			Err(err).
			Msg("RecognizeRevenue failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "failed to post revenue journal: %v", err)
	}

	logger.Info().
		Str("op", op).
		Str("departure_id", req.GetDepartureId()).
		Str("entry_id", result.EntryID).
		Bool("replayed", result.Replayed).
		Msg("RecognizeRevenue succeeded")

	span.SetStatus(otelCodes.Ok, "success")
	return &pb.RecognizeRevenueResponse{
		EntryId:      result.EntryID,
		RecognizedAt: result.RecognizedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		Replayed:     result.Replayed,
	}, nil
}
