// on_payment_received.go — gRPC handler for OnPaymentReceived RPC (S3-E-03).
//
// Called by booking-svc as a fire-and-forget after a booking transitions to
// paid_in_full (ADR-0006 direct gRPC, no Temporal).
//
// Behaviour:
//   - Validates booking_id, invoice_id, and amount.
//   - Delegates to service.OnPaymentReceived which is idempotent.
//   - Returns { entry_id, balanced: true } on success.

package grpc_api

import (
	"context"
	"time"

	"finance-svc/api/grpc_api/pb"
	"finance-svc/service"
	"finance-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// OnPaymentReceived handles the OnPaymentReceived RPC from booking-svc.
// It posts a double-entry journal (Dr Bank / Cr Pilgrim Liability) for the
// payment. Idempotent: returns existing entry if already posted for this
// invoice.
func (s *Server) OnPaymentReceived(ctx context.Context, req *pb.OnPaymentReceivedRequest) (*pb.OnPaymentReceivedResponse, error) {
	const op = "grpc_api.Server.OnPaymentReceived"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "OnPaymentReceived"),
		attribute.String("booking_id", req.GetBookingId()),
		attribute.String("invoice_id", req.GetInvoiceId()),
		attribute.Int64("amount", req.GetAmount()),
	)

	logger.Info().
		Str("op", op).
		Str("booking_id", req.GetBookingId()).
		Str("invoice_id", req.GetInvoiceId()).
		Int64("amount", req.GetAmount()).
		Msg("")

	if req.GetBookingId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "booking_id is required")
	}
	if req.GetInvoiceId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "invoice_id is required")
	}
	if req.GetAmount() <= 0 {
		return nil, grpcStatus.Errorf(grpcCodes.InvalidArgument, "amount must be positive, got %d", req.GetAmount())
	}

	// Parse optional received_at; fall back to zero (service uses time.Now()).
	var receivedAt time.Time
	if ts := req.GetReceivedAt(); ts != "" {
		if t, err := time.Parse(time.RFC3339, ts); err == nil {
			receivedAt = t
		}
	}

	result, err := s.svc.OnPaymentReceived(ctx, &service.OnPaymentReceivedParams{
		BookingID:  req.GetBookingId(),
		InvoiceID:  req.GetInvoiceId(),
		Amount:     req.GetAmount(), // int64 IDR
		ReceivedAt: receivedAt,
	})
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("invoice_id", req.GetInvoiceId()).
			Err(err).
			Msg("OnPaymentReceived failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "failed to post journal entry: %v", err)
	}

	logger.Info().
		Str("op", op).
		Str("booking_id", req.GetBookingId()).
		Str("invoice_id", req.GetInvoiceId()).
		Str("entry_id", result.EntryID).
		Bool("balanced", result.Balanced).
		Bool("replayed", result.Replayed).
		Msg("OnPaymentReceived succeeded")

	span.SetStatus(otelCodes.Ok, "success")
	return &pb.OnPaymentReceivedResponse{
		EntryId:  result.EntryID,
		Balanced: result.Balanced,
	}, nil
}
