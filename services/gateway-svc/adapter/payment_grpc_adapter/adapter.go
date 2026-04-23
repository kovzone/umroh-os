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
// IssueVirtualAccount
// ---------------------------------------------------------------------------

// IssueVirtualAccountParams holds the input for creating an invoice + VA.
type IssueVirtualAccountParams struct {
	BookingID   string
	AmountTotal float64
	GatewayPref string
	BankCode    string
	ActorUserID string
}

// IssueVirtualAccountResult holds the issued VA details.
type IssueVirtualAccountResult struct {
	InvoiceID     string
	BookingID     string
	AccountNumber string
	BankCode      string
	AmountTotal   float64
	ExpiresAt     string
	Gateway       string
	Replayed      bool
}

func (a *Adapter) IssueVirtualAccount(ctx context.Context, params *IssueVirtualAccountParams) (*IssueVirtualAccountResult, error) {
	const op = "payment_grpc_adapter.Adapter.IssueVirtualAccount"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.paymentClient.IssueVirtualAccount(ctx, &pb.IssueVirtualAccountRequest{
		BookingId:   params.BookingID,
		AmountTotal: params.AmountTotal,
		GatewayPref: params.GatewayPref,
		BankCode:    params.BankCode,
		ActorUserId: params.ActorUserID,
	})
	if err != nil {
		wrapped := mapPaymentError(err)
		logger.Warn().Err(wrapped).Msg("payment-svc.IssueVirtualAccount failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &IssueVirtualAccountResult{
		InvoiceID:     resp.GetInvoiceId(),
		BookingID:     resp.GetBookingId(),
		AccountNumber: resp.GetAccountNumber(),
		BankCode:      resp.GetBankCode(),
		AmountTotal:   resp.GetAmountTotal(),
		ExpiresAt:     resp.GetExpiresAt(),
		Gateway:       resp.GetGateway(),
		Replayed:      resp.GetReplayed(),
	}, nil
}

// ---------------------------------------------------------------------------
// GetInvoiceByID
// ---------------------------------------------------------------------------

// GetInvoiceByIDParams holds the input for fetching a single invoice.
type GetInvoiceByIDParams struct {
	InvoiceID string
}

// GetInvoiceByIDResult holds the invoice details.
type GetInvoiceByIDResult struct {
	ID          string
	BookingID   string
	Status      string
	AmountTotal float64
	PaidAmount  float64
	Currency    string
	CreatedAt   string
	UpdatedAt   string
}

func (a *Adapter) GetInvoiceByID(ctx context.Context, params *GetInvoiceByIDParams) (*GetInvoiceByIDResult, error) {
	const op = "payment_grpc_adapter.Adapter.GetInvoiceByID"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.paymentClient.GetInvoiceByID(ctx, &pb.GetInvoiceByIDRequest{
		InvoiceId: params.InvoiceID,
	})
	if err != nil {
		wrapped := mapPaymentError(err)
		logger.Warn().Err(wrapped).Msg("payment-svc.GetInvoiceByID failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &GetInvoiceByIDResult{
		ID:          resp.GetId(),
		BookingID:   resp.GetBookingId(),
		Status:      resp.GetStatus(),
		AmountTotal: resp.GetAmountTotal(),
		PaidAmount:  resp.GetPaidAmount(),
		Currency:    resp.GetCurrency(),
		CreatedAt:   resp.GetCreatedAt(),
		UpdatedAt:   resp.GetUpdatedAt(),
	}, nil
}

// ---------------------------------------------------------------------------
// ProcessWebhook
// ---------------------------------------------------------------------------

// ProcessWebhookParams holds the input for forwarding a webhook payload.
type ProcessWebhookParams struct {
	Gateway   string
	Payload   []byte
	Signature string
}

// ProcessWebhookResult holds the result from processing a webhook.
type ProcessWebhookResult struct {
	Replayed  bool
	InvoiceID string
	NewStatus string
}

func (a *Adapter) ProcessWebhook(ctx context.Context, params *ProcessWebhookParams) (*ProcessWebhookResult, error) {
	const op = "payment_grpc_adapter.Adapter.ProcessWebhook"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.paymentClient.ProcessWebhook(ctx, &pb.ProcessWebhookRequest{
		Gateway:   params.Gateway,
		Payload:   params.Payload,
		Signature: params.Signature,
	})
	if err != nil {
		wrapped := mapPaymentError(err)
		logger.Warn().Err(wrapped).Msg("payment-svc.ProcessWebhook failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &ProcessWebhookResult{
		Replayed:  resp.GetReplayed(),
		InvoiceID: resp.GetInvoiceId(),
		NewStatus: resp.GetNewStatus(),
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
