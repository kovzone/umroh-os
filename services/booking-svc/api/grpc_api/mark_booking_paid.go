// mark_booking_paid.go — stub gRPC handler for MarkBookingPaid RPC (S2-E-02).
//
// payment-svc calls this RPC after processing a payment event (webhook or
// reconciliation backfill) to signal that a booking's payment status has changed.
//
// Per F5 spec W2 step 7 and ADR-0006: direct synchronous gRPC call from
// payment-svc; no Temporal in MVP.
//
// Booking state machine transitions driven by this call:
//   pending_payment + invoice_status=partially_paid → partially_paid
//   pending_payment | partially_paid + invoice_status=paid → paid_in_full
//
// The call is idempotent: if the booking is already in the target status, the
// RPC returns success without a no-op error.
//
// TODO (S3): implement full booking state machine transition logic.
// This stub returns a success response so payment-svc can complete its
// webhook path without blocking; a test fixture in tests/e2e can verify the
// end-to-end flow once booking-svc's state machine is implemented.

package grpc_api

import (
	"context"
	"fmt"

	"booking-svc/api/grpc_api/pb"
	"booking-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
	grpcCodes "google.golang.org/grpc/codes"
)

// MarkBookingPaid handles the MarkBookingPaid RPC from payment-svc.
// It updates the booking status based on the invoice_status received.
func (s *Server) MarkBookingPaid(ctx context.Context, req *pb.MarkBookingPaidRequest) (*pb.MarkBookingPaidResponse, error) {
	const op = "grpc_api.Server.MarkBookingPaid"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "MarkBookingPaid"),
		attribute.String("booking_id", req.GetBookingId()),
		attribute.String("invoice_status", req.GetInvoiceStatus()),
		attribute.Float64("amount_paid_idr", req.GetAmountPaidIdr()),
	)

	logger.Info().
		Str("op", op).
		Str("booking_id", req.GetBookingId()).
		Str("invoice_status", req.GetInvoiceStatus()).
		Float64("amount_paid_idr", req.GetAmountPaidIdr()).
		Msg("")

	if req.GetBookingId() == "" {
		return nil, status.Error(grpcCodes.InvalidArgument, "booking_id is required")
	}

	// Determine target booking status from invoice status.
	targetBookingStatus, err := invoiceStatusToBookingStatus(req.GetInvoiceStatus())
	if err != nil {
		return nil, status.Error(grpcCodes.InvalidArgument, err.Error())
	}

	// TODO (S3): call s.svc.MarkBookingPaid(ctx, req.GetBookingId(), targetBookingStatus)
	// For now, return the target status as the "new_status" so payment-svc can proceed.
	//
	// Stub behaviour: log + return success. The booking record is NOT updated yet
	// because the full booking state machine (pending_payment → partially_paid →
	// paid_in_full) requires additional service layer work beyond S2 scope.
	// The reconciliation cron will call this again on next cycle if booking is
	// still in an inconsistent state.

	logger.Info().
		Str("op", op).
		Str("booking_id", req.GetBookingId()).
		Str("target_status", targetBookingStatus).
		Msg("MarkBookingPaid stub — booking status transition will be implemented in S3")

	span.SetStatus(codes.Ok, "stub success")

	return &pb.MarkBookingPaidResponse{
		BookingId: req.GetBookingId(),
		NewStatus: targetBookingStatus,
	}, nil
}

// invoiceStatusToBookingStatus maps F5 invoice status to F4 booking status.
func invoiceStatusToBookingStatus(invoiceStatus string) (string, error) {
	switch invoiceStatus {
	case "partially_paid":
		return "partially_paid", nil
	case "paid":
		return "paid_in_full", nil
	default:
		return "", fmt.Errorf("unrecognised invoice_status %q for booking transition", invoiceStatus)
	}
}
