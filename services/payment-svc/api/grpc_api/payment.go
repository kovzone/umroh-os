// payment.go — gRPC handler methods for payment-svc (S2-E-02).
//
// Implements the payment RPCs on *Server (which embeds pb.UnimplementedPaymentServiceServer).
// All business logic is delegated to the service layer (service.IPaymentService).
//
// RPCs implemented here:
//   IssueVirtualAccount  — F5-W1: create invoice + VA for a booking
//   ProcessWebhook       — F5-W2: gateway webhook ingestion
//   StartRefund          — F5-W8: initiate refund flow
//
// Note: ProcessWebhook is ALSO exposed via an HTTP listener in api/http_api/webhook.go
// because Midtrans/Xendit POST to a public URL. The gRPC method exists so that
// internal callers (e.g. a future gateway-svc proxy) can call it over gRPC.
// Both paths call the same service layer method.

package grpc_api

import (
	"context"

	"payment-svc/api/grpc_api/pb"
	"payment-svc/service"
	"payment-svc/util/apperrors"
	"payment-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

// IssueVirtualAccount creates an invoice and a virtual account for a booking.
// Called by booking-svc's submit saga (F4→F5 call per ADR-0006).
func (s *Server) IssueVirtualAccount(ctx context.Context, req *pb.IssueVirtualAccountRequest) (*pb.IssueVirtualAccountResponse, error) {
	const op = "grpc_api.Server.IssueVirtualAccount"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("booking_id", req.GetBookingId()),
	)
	logger.Info().
		Str("op", op).
		Str("booking_id", req.GetBookingId()).
		Float64("amount_total", req.GetAmountTotal()).
		Msg("")

	// Type-assert to IPaymentService.
	psvc, ok := s.svc.(service.IPaymentService)
	if !ok {
		err := apperrors.ErrInternal
		span.SetStatus(codes.Error, "service does not implement IPaymentService")
		return nil, status.Error(apperrors.GRPCCode(err), "service misconfiguration")
	}

	result, err := psvc.IssueVirtualAccount(ctx, &service.IssueVAParams{
		BookingID:             req.GetBookingId(),
		AmountTotal:           req.GetAmountTotal(),
		RoundingAdjustmentIDR: req.GetRoundingAdjustmentIdr(),
		GatewayPref:           req.GetGatewayPref(),
		BankCode:              req.GetBankCode(),
		ActorUserID:           req.GetActorUserId(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetAttributes(
		attribute.String("invoice_id", result.InvoiceID),
		attribute.Bool("replayed", result.Replayed),
	)
	span.SetStatus(codes.Ok, "success")

	return &pb.IssueVirtualAccountResponse{
		InvoiceId:     result.InvoiceID,
		BookingId:     result.BookingID,
		AccountNumber: result.AccountNumber,
		BankCode:      result.BankCode,
		AmountTotal:   result.AmountTotal,
		ExpiresAt:     result.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"),
		Gateway:       result.Gateway,
		Replayed:      result.Replayed,
	}, nil
}

// ProcessWebhook handles a gateway webhook forwarded via gRPC.
// The HTTP webhook handler (api/http_api/webhook.go) calls this same service
// method; this gRPC method exists for internal callers.
func (s *Server) ProcessWebhook(ctx context.Context, req *pb.ProcessWebhookRequest) (*pb.ProcessWebhookResponse, error) {
	const op = "grpc_api.Server.ProcessWebhook"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("gateway", req.GetGateway()),
	)
	logger.Info().Str("op", op).Str("gateway", req.GetGateway()).Msg("")

	psvc, ok := s.svc.(service.IPaymentService)
	if !ok {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrInternal), "service misconfiguration")
	}

	result, err := psvc.ProcessWebhookEvent(ctx, &service.WebhookEventParams{
		Gateway:   req.GetGateway(),
		Payload:   req.GetPayload(),
		Signature: req.GetSignature(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.ProcessWebhookResponse{
		Replayed:  result.Replayed,
		InvoiceId: result.InvoiceID,
		NewStatus: result.NewStatus,
	}, nil
}

// StartRefund initiates the refund flow for a booking.
// Called by booking-svc's cancel saga (F4→F5 gRPC call per ADR-0006).
func (s *Server) StartRefund(ctx context.Context, req *pb.StartRefundRequest) (*pb.StartRefundResponse, error) {
	const op = "grpc_api.Server.StartRefund"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("booking_id", req.GetBookingId()),
		attribute.String("reason_code", req.GetReasonCode()),
	)
	logger.Info().
		Str("op", op).
		Str("booking_id", req.GetBookingId()).
		Str("reason_code", req.GetReasonCode()).
		Msg("")

	psvc, ok := s.svc.(service.IPaymentService)
	if !ok {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrInternal), "service misconfiguration")
	}

	result, err := psvc.StartRefund(ctx, &service.StartRefundParams{
		BookingID:   req.GetBookingId(),
		ReasonCode:  req.GetReasonCode(),
		ActorUserID: req.GetActorUserId(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.StartRefundResponse{
		RefundId:  result.RefundID,
		InvoiceId: result.InvoiceID,
		AmountIdr: result.AmountIDR,
		Status:    result.Status,
	}, nil
}
