// mark_booking_paid.go — gRPC handler for MarkBookingPaid RPC (S2-E-02 / S3).
//
// payment-svc calls this RPC after processing a payment event (webhook or
// reconciliation backfill) to signal that a booking's payment status has changed.
//
// Per F5 spec W2 step 7 and ADR-0006: direct synchronous gRPC call from
// payment-svc; no Temporal in MVP.
//
// Booking state machine transitions:
//   pending_payment + invoice_status=partially_paid → partially_paid
//   pending_payment | partially_paid + invoice_status=paid → paid_in_full
//
// S3 addition: when invoice_status == "paid" (paid_in_full), performs synchronous fan-out:
//   - logistics-svc.OnBookingPaid (S3-E-02)
//   - finance-svc.OnPaymentReceived (S3-E-03)
// Per ADR-0006 / §S3-J-01: fan-out is SYNCHRONOUS — errors surface to payment-svc
// as INTERNAL so the gateway can retry. Both downstream services are idempotent.
//
// The call is idempotent: if the booking is already in the target status, the
// RPC returns success without a no-op error.
//
// TODO (future): implement full booking DB state machine transition.
// The stub returns the target status so payment-svc can proceed.

package grpc_api

import (
	"context"
	"fmt"

	"booking-svc/api/grpc_api/pb"
	"booking-svc/service"
	"booking-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
	grpcCodes "google.golang.org/grpc/codes"
)

// MarkBookingPaid handles the MarkBookingPaid RPC from payment-svc.
// It updates the booking status based on the invoice_status received.
// When invoice_status == "paid", it performs a synchronous fan-out to logistics
// and finance services (ADR-0006 / §S3-J-01). Errors from fan-out are returned
// as INTERNAL so the gateway can retry.
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
		attribute.String("invoice_id", req.GetInvoiceId()),
		attribute.Float64("amount_paid_idr", req.GetAmountPaidIdr()),
		attribute.String("departure_id", req.GetDepartureId()),
	)

	logger.Info().
		Str("op", op).
		Str("booking_id", req.GetBookingId()).
		Str("invoice_status", req.GetInvoiceStatus()).
		Str("invoice_id", req.GetInvoiceId()).
		Str("departure_id", req.GetDepartureId()).
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

	// TODO (future): call s.svc.MarkBookingPaid to persist the DB transition.
	// Logging the intent; DB update deferred to a future task.
	logger.Info().
		Str("op", op).
		Str("booking_id", req.GetBookingId()).
		Str("target_status", targetBookingStatus).
		Msg("MarkBookingPaid — target status determined")

	// --- S3 Fan-out: only on fully paid bookings (§S3-J-01) ---
	// Partially paid bookings do NOT trigger logistics fan-out; finance fan-out
	// is called on every payment (partial or full) per §S3-J-01 consumer rules.
	// For now the fan-out gate is: invoice_status == "paid" triggers both calls.
	// TODO (S3-E-02): retrieve departure_id from booking DB record once the full
	// booking state machine persists it. For now, departure_id is read from the
	// request (caller must populate it once the proto field is extended).
	if req.GetInvoiceStatus() == "paid" {
		// departure_id must come from the booking record (TODO: DB lookup).
		// Temporarily use the field from the request; logistics-svc will return
		// INVALID_ARGUMENT and MarkBookingPaid will surface INTERNAL to payment-svc
		// so the gateway retries — this is correct per ADR-0006.
		departureID := req.GetDepartureId()

		// TODO (S3-E-02): resolve jamaah_ids from booking DB record.
		// Provided by caller in JamaahIds field until DB lookup is wired.
		if _, err := s.svc.FanOutBookingPaid(ctx, &service.FanOutBookingPaidParams{
			BookingID:   req.GetBookingId(),
			DepartureID: departureID,
			JamaahIDs:   req.GetJamaahIds(),
			InvoiceID:   req.GetInvoiceId(),
			Amount:      int64(req.GetAmountPaidIdr()), // contract: amount is integer IDR
			ReceivedAt:  req.GetReceivedAt(),
		}); err != nil {
			logger.Error().
				Str("op", op).
				Str("booking_id", req.GetBookingId()).
				Err(err).
				Msg("MarkBookingPaid fan-out failed")
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, status.Errorf(grpcCodes.Internal, "fan-out failed: %v", err)
		}
	}

	span.SetStatus(codes.Ok, "success")

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
