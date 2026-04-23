// on_booking_paid.go — gRPC handler for OnBookingPaid RPC (S3-E-02).
//
// Called by booking-svc as a fire-and-forget after a booking transitions to
// paid_in_full (ADR-0006 direct gRPC, no Temporal).
//
// Behaviour:
//   - Validates booking_id and departure_id are non-empty.
//   - Delegates to service.OnBookingPaid which is idempotent.
//   - Returns { task_id, status } on success.

package grpc_api

import (
	"context"

	"logistics-svc/api/grpc_api/pb"
	"logistics-svc/service"
	"logistics-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// OnBookingPaid handles the OnBookingPaid RPC from booking-svc.
// It creates a fulfillment task in status='queued' for the paid booking.
// Idempotent: if a task already exists for this booking, returns it unchanged.
func (s *Server) OnBookingPaid(ctx context.Context, req *pb.OnBookingPaidRequest) (*pb.OnBookingPaidResponse, error) {
	const op = "grpc_api.Server.OnBookingPaid"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "OnBookingPaid"),
		attribute.String("booking_id", req.GetBookingId()),
		attribute.String("departure_id", req.GetDepartureId()),
		attribute.Int("jamaah_count", len(req.GetJamaahIds())),
	)

	logger.Info().
		Str("op", op).
		Str("booking_id", req.GetBookingId()).
		Str("departure_id", req.GetDepartureId()).
		Msg("")

	if req.GetBookingId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "booking_id is required")
	}
	if req.GetDepartureId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "departure_id is required")
	}
	if len(req.GetJamaahIds()) == 0 {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "jamaah_ids is required (min 1)")
	}

	result, err := s.svc.OnBookingPaid(ctx, &service.OnBookingPaidParams{
		BookingID:   req.GetBookingId(),
		DepartureID: req.GetDepartureId(),
		JamaahIDs:   req.GetJamaahIds(),
	})
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("booking_id", req.GetBookingId()).
			Err(err).
			Msg("OnBookingPaid failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "failed to create fulfillment task: %v", err)
	}

	logger.Info().
		Str("op", op).
		Str("booking_id", req.GetBookingId()).
		Str("task_id", result.TaskID).
		Str("status", result.Status).
		Bool("replayed", result.Replayed).
		Msg("OnBookingPaid succeeded")

	span.SetStatus(otelCodes.Ok, "success")
	return &pb.OnBookingPaidResponse{
		TaskId: result.TaskID,
		Status: result.Status,
	}, nil
}
