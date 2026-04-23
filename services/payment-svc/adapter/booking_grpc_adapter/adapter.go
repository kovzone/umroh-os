// Package booking_grpc_adapter is payment-svc's consumer-side wrapper around
// booking-svc's internal gRPC surface.
//
// Only the MarkBookingPaid RPC is needed by payment-svc for S2-E-02.
// This is the only direction of F5→F4 signalling for the MVP payment flow
// (per F5 spec W2 step 7 and ADR-0006 in-process saga with direct gRPC calls).
//
// Architecture note — webhook transport choice:
//   For MVP, payment-svc exposes a SECOND HTTP listener (port 50065, internal
//   only) for webhook ingestion from Midtrans/Xendit. This is pragmatic:
//   - Midtrans/Xendit require a public URL they POST to; gateway-svc could
//     proxy this, but that adds a gRPC hop on the 500ms critical path.
//   - The webhook port is NOT exposed to the internet directly in production;
//     an nginx/load-balancer rule forwards only POST /v1/webhooks/* from the
//     gateway's IP range to this port.
//   - gateway-svc is NOT modified for S2 (per task constraint).
//   This decision is documented here so future S3 work can revisit.

package booking_grpc_adapter

import (
	"context"
	"errors"
	"fmt"

	"payment-svc/util/apperrors"
	"payment-svc/util/logging"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
	"google.golang.org/grpc"
)

// Adapter wraps booking-svc's gRPC connection.
type Adapter struct {
	logger *zerolog.Logger
	tracer trace.Tracer
	conn   *grpc.ClientConn
}

// NewAdapter creates a new booking-svc adapter from an already-dialled conn.
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, conn *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger: logger,
		tracer: tracer,
		conn:   conn,
	}
}

// MarkBookingPaidParams carries the parameters for the MarkBookingPaid call.
type MarkBookingPaidParams struct {
	BookingID     string
	AmountPaidIDR float64
	InvoiceStatus string // "partially_paid" | "paid"
	InvoiceID     string
}

// MarkBookingPaidResult is the successful response from booking-svc.
type MarkBookingPaidResult struct {
	BookingID string
	NewStatus string
}

// markBookingPaidRequest is the wire-compatible request struct.
// Fields match the hand-written proto in booking-svc/api/grpc_api/pb.
// Uses JSON field names matching proto snake_case convention.
type markBookingPaidRequest struct {
	BookingId     string  `json:"booking_id"`
	AmountPaidIdr float64 `json:"amount_paid_idr"`
	InvoiceStatus string  `json:"invoice_status"`
	InvoiceId     string  `json:"invoice_id"`
}

// markBookingPaidResponse matches the wire response.
type markBookingPaidResponse struct {
	BookingId string `json:"booking_id"`
	NewStatus string `json:"new_status"`
}

// ProtoMessage satisfies proto.Message interface minimally.
func (m *markBookingPaidRequest) ProtoMessage()  {}
func (m *markBookingPaidRequest) Reset()         {}
func (m *markBookingPaidRequest) String() string { return fmt.Sprintf("%+v", *m) }

func (m *markBookingPaidResponse) ProtoMessage()  {}
func (m *markBookingPaidResponse) Reset()         {}
func (m *markBookingPaidResponse) String() string { return fmt.Sprintf("%+v", *m) }

// MarkBookingPaid calls booking-svc.MarkBookingPaid to signal that a payment
// event has been processed and the invoice status has changed.
//
// The RPC is idempotent: booking-svc must handle duplicate signals gracefully
// (same invoice_status → no-op state transition).
func (a *Adapter) MarkBookingPaid(ctx context.Context, params *MarkBookingPaidParams) (*MarkBookingPaidResult, error) {
	const op = "booking_grpc_adapter.Adapter.MarkBookingPaid"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "MarkBookingPaid"),
		attribute.String("booking_id", params.BookingID),
		attribute.String("invoice_status", params.InvoiceStatus),
		attribute.String("invoice_id", params.InvoiceID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().
		Str("op", op).
		Str("booking_id", params.BookingID).
		Str("invoice_status", params.InvoiceStatus).
		Float64("amount_paid_idr", params.AmountPaidIDR).
		Msg("")

	req := &markBookingPaidRequest{
		BookingId:     params.BookingID,
		AmountPaidIdr: params.AmountPaidIDR,
		InvoiceStatus: params.InvoiceStatus,
		InvoiceId:     params.InvoiceID,
	}
	resp := &markBookingPaidResponse{}

	err := a.conn.Invoke(ctx,
		"/booking.v1.BookingService/MarkBookingPaid",
		req,
		resp,
	)
	if err != nil {
		wrapped := mapBookingError(err)
		logger.Warn().Err(wrapped).Str("op", op).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetAttributes(attribute.String("new_status", resp.NewStatus))
	span.SetStatus(codes.Ok, "success")
	return &MarkBookingPaidResult{
		BookingID: resp.BookingId,
		NewStatus: resp.NewStatus,
	}, nil
}

// ---------------------------------------------------------------------------
// Error mapping
// ---------------------------------------------------------------------------

func mapBookingError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := grpcStatus.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("booking-svc call failed: %w", err))
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return errors.Join(apperrors.ErrNotFound, errors.New(st.Message()))
	case grpcCodes.FailedPrecondition, grpcCodes.AlreadyExists:
		return errors.Join(apperrors.ErrConflict, errors.New(st.Message()))
	case grpcCodes.InvalidArgument:
		return errors.Join(apperrors.ErrValidation, errors.New(st.Message()))
	case grpcCodes.Unauthenticated:
		return errors.Join(apperrors.ErrUnauthorized, errors.New(st.Message()))
	case grpcCodes.PermissionDenied:
		return errors.Join(apperrors.ErrForbidden, errors.New(st.Message()))
	default:
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("booking-svc call failed: %s", st.Message()))
	}
}
