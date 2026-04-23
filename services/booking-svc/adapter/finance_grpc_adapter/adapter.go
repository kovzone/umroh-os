// Package finance_grpc_adapter is booking-svc's consumer-side wrapper around
// finance-svc's internal gRPC surface. It exposes only OnPaymentReceived, the
// single RPC needed by the paid-booking fan-out path (S3-E-03 / ADR-0006).

package finance_grpc_adapter

import (
	"context"
	"errors"
	"fmt"

	"booking-svc/adapter/finance_grpc_adapter/pb"
	"booking-svc/util/apperrors"
	"booking-svc/util/logging"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Adapter wraps finance-svc's FinanceServiceClient.
type Adapter struct {
	logger        *zerolog.Logger
	tracer        trace.Tracer
	financeClient pb.FinanceServiceClient
}

// NewAdapter creates a finance adapter from an already-dialled conn.
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:        logger,
		tracer:        tracer,
		financeClient: pb.NewFinanceServiceClient(cc),
	}
}

// ---------------------------------------------------------------------------
// OnPaymentReceived — trigger double-entry journal posting.
// ---------------------------------------------------------------------------

// OnPaymentReceivedParams maps to FinanceService.OnPaymentReceived request.
// Amount is int64 (integer IDR) per §S3-J-03 contract.
type OnPaymentReceivedParams struct {
	BookingID  string
	InvoiceID  string
	Amount     int64  // integer IDR — no fractional amounts
	ReceivedAt string // RFC3339; empty = server time
}

// OnPaymentReceivedResult maps to FinanceService.OnPaymentReceived response.
type OnPaymentReceivedResult struct {
	EntryID  string
	Balanced bool
}

// OnPaymentReceived calls finance-svc.OnPaymentReceived and returns the created
// (or existing) journal entry details.
func (a *Adapter) OnPaymentReceived(ctx context.Context, params *OnPaymentReceivedParams) (*OnPaymentReceivedResult, error) {
	const op = "finance_grpc_adapter.Adapter.OnPaymentReceived"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("booking_id", params.BookingID),
		attribute.String("invoice_id", params.InvoiceID),
		attribute.Int64("amount", params.Amount),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.financeClient.OnPaymentReceived(ctx, &pb.OnPaymentReceivedRequest{
		BookingId:  params.BookingID,
		InvoiceId:  params.InvoiceID,
		Amount:     params.Amount,
		ReceivedAt: params.ReceivedAt,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().
			Err(wrapped).
			Str("booking_id", params.BookingID).
			Str("invoice_id", params.InvoiceID).
			Msg("finance-svc.OnPaymentReceived failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &OnPaymentReceivedResult{
		EntryID:  resp.GetEntryId(),
		Balanced: resp.GetBalanced(),
	}, nil
}

// ---------------------------------------------------------------------------
// Error mapping
// ---------------------------------------------------------------------------

func mapFinanceError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("finance call failed: %w", err))
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return errors.Join(apperrors.ErrNotFound, errors.New(st.Message()))
	case grpcCodes.InvalidArgument:
		return errors.Join(apperrors.ErrValidation, errors.New(st.Message()))
	case grpcCodes.AlreadyExists:
		return errors.Join(apperrors.ErrConflict, errors.New(st.Message()))
	case grpcCodes.Unauthenticated:
		return errors.Join(apperrors.ErrUnauthorized, errors.New(st.Message()))
	case grpcCodes.PermissionDenied:
		return errors.Join(apperrors.ErrForbidden, errors.New(st.Message()))
	default:
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("finance call failed: %s", st.Message()))
	}
}
