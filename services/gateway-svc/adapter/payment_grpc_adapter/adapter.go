// Package payment_grpc_adapter is gateway-svc's consumer-side wrapper around
// payment-svc's gRPC surface (BL-PAY-020 / ADR-0009).
//
// Per ADR-0009, the CS-facing payment link route proxies to payment-svc over
// gRPC via this adapter:
//   POST /v1/payments/link → ReissuePaymentLink (bearer)

package payment_grpc_adapter

import (
	"context"
	"errors"
	"fmt"

	"gateway-svc/adapter/payment_grpc_adapter/pb"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Adapter wraps payment-svc's PaymentServiceClient.
type Adapter struct {
	logger         *zerolog.Logger
	tracer         trace.Tracer
	paymentClient  pb.PaymentServiceClient
}

// NewAdapter creates a payment gRPC adapter from an already-dialled conn.
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:        logger,
		tracer:        tracer,
		paymentClient: pb.NewPaymentServiceClient(cc),
	}
}

// ---------------------------------------------------------------------------
// ReissuePaymentLink
// ---------------------------------------------------------------------------

// ReissuePaymentLinkParams holds the CS-facing input.
type ReissuePaymentLinkParams struct {
	BookingID   string
	BankCode    string
	GatewayPref string
	ActorUserID string
}

// ReissuePaymentLinkResult holds the adapter-local result type.
// No proto types leak beyond this package.
type ReissuePaymentLinkResult struct {
	InvoiceID     string
	BookingID     string
	AccountNumber string
	BankCode      string
	AmountTotal   float64
	ExpiresAt     string
	Gateway       string
	IsNew         bool
}

func (a *Adapter) ReissuePaymentLink(ctx context.Context, params *ReissuePaymentLinkParams) (*ReissuePaymentLinkResult, error) {
	const op = "payment_grpc_adapter.Adapter.ReissuePaymentLink"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.paymentClient.ReissuePaymentLink(ctx, &pb.ReissuePaymentLinkRequest{
		BookingId:   params.BookingID,
		BankCode:    params.BankCode,
		GatewayPref: params.GatewayPref,
		ActorUserId: params.ActorUserID,
	})
	if err != nil {
		wrapped := mapPaymentError(err)
		logger.Warn().Err(wrapped).Msg("payment-svc.ReissuePaymentLink failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &ReissuePaymentLinkResult{
		InvoiceID:     resp.GetInvoiceId(),
		BookingID:     resp.GetBookingId(),
		AccountNumber: resp.GetAccountNumber(),
		BankCode:      resp.GetBankCode(),
		AmountTotal:   resp.GetAmountTotal(),
		ExpiresAt:     resp.GetExpiresAt(),
		Gateway:       resp.GetGateway(),
		IsNew:         resp.GetIsNew(),
	}, nil
}

// ---------------------------------------------------------------------------
// Error mapping
// ---------------------------------------------------------------------------

func mapPaymentError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("payment call failed: %w", err))
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return errors.Join(apperrors.ErrNotFound, errors.New(st.Message()))
	case grpcCodes.InvalidArgument:
		return errors.Join(apperrors.ErrValidation, errors.New(st.Message()))
	case grpcCodes.Unauthenticated:
		return errors.Join(apperrors.ErrUnauthorized, errors.New(st.Message()))
	case grpcCodes.PermissionDenied:
		return errors.Join(apperrors.ErrForbidden, errors.New(st.Message()))
	case grpcCodes.Unavailable:
		return errors.Join(apperrors.ErrServiceUnavailable, errors.New(st.Message()))
	default:
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("payment call failed: %s", st.Message()))
	}
}
