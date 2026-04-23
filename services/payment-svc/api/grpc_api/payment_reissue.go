// payment_reissue.go — gRPC handler for ReissuePaymentLink (BL-PAY-020).
//
// CS-facing RPC: retrieves or re-issues the VA link for an existing booking.

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

// ReissuePaymentLink retrieves the active VA for a booking or creates a fresh
// one on the same invoice when the existing VA has expired.
// Implements pb.PaymentServiceWithReissueServer.
func (s *Server) ReissuePaymentLink(ctx context.Context, req *pb.ReissuePaymentLinkRequest) (*pb.ReissuePaymentLinkResponse, error) {
	const op = "grpc_api.Server.ReissuePaymentLink"

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
		Msg("")

	if req.GetBookingId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "booking_id is required")
	}

	psvc, ok := s.svc.(service.IPaymentService)
	if !ok {
		span.SetStatus(codes.Error, "service does not implement IPaymentService")
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrInternal), "service misconfiguration")
	}

	result, err := psvc.ReissuePaymentLink(ctx, &service.ReissuePaymentLinkParams{
		BookingID:   req.GetBookingId(),
		BankCode:    req.GetBankCode(),
		GatewayPref: req.GetGatewayPref(),
		ActorUserID: req.GetActorUserId(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetAttributes(
		attribute.String("invoice_id", result.InvoiceID),
		attribute.Bool("is_new", result.IsNew),
	)
	span.SetStatus(codes.Ok, "success")

	return &pb.ReissuePaymentLinkResponse{
		InvoiceId:     result.InvoiceID,
		BookingId:     result.BookingID,
		AccountNumber: result.AccountNumber,
		BankCode:      result.BankCode,
		AmountTotal:   result.AmountTotal,
		ExpiresAt:     result.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"),
		Gateway:       result.Gateway,
		IsNew:         result.IsNew,
	}, nil
}
