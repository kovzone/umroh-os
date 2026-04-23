// pickup_qr.go — gRPC handlers for GeneratePickupQR and RedeemPickupQR RPCs (BL-LOG-003).
//
// GeneratePickupQR: lookup task, create pickup_token (token = UUID v4, expires 7 days),
// idempotent (returns existing non-expired, non-used token if one exists).
//
// RedeemPickupQR: find token, reject if expired or already used, mark used,
// return booking_id.

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

// GeneratePickupQR handles the GeneratePickupQR RPC.
// Creates a single-use pickup QR token with 7-day TTL for the given booking.
// Idempotent: returns existing active token if one exists.
func (s *Server) GeneratePickupQR(ctx context.Context, req *pb.GeneratePickupQRRequest) (*pb.GeneratePickupQRResponse, error) {
	const op = "grpc_api.Server.GeneratePickupQR"

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
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "booking_id is required")
	}

	result, err := s.svc.GeneratePickupQR(ctx, &service.GeneratePickupQRParams{
		BookingID: req.GetBookingId(),
	})
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("booking_id", req.GetBookingId()).
			Err(err).
			Msg("GeneratePickupQR failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "failed to generate pickup QR: %v", err)
	}

	logger.Info().
		Str("op", op).
		Str("booking_id", req.GetBookingId()).
		Str("token_id", result.PickupTokenID).
		Bool("replayed", result.Replayed).
		Msg("GeneratePickupQR succeeded")

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GeneratePickupQRResponse{
		PickupTokenId: result.PickupTokenID,
		Token:         result.Token,
		ExpiresAt:     result.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"),
	}, nil
}

// RedeemPickupQR handles the RedeemPickupQR RPC.
// Validates and marks a pickup token as used.
// Returns Redeemed=false with ErrorReason if the token is expired or already used.
func (s *Server) RedeemPickupQR(ctx context.Context, req *pb.RedeemPickupQRRequest) (*pb.RedeemPickupQRResponse, error) {
	const op = "grpc_api.Server.RedeemPickupQR"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("token", req.GetToken()),
	)

	logger.Info().
		Str("op", op).
		Str("token", req.GetToken()).
		Msg("")

	if req.GetToken() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "token is required")
	}

	result, err := s.svc.RedeemPickupQR(ctx, &service.RedeemPickupQRParams{
		Token: req.GetToken(),
	})
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("token", req.GetToken()).
			Err(err).
			Msg("RedeemPickupQR failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "failed to redeem pickup QR: %v", err)
	}

	logger.Info().
		Str("op", op).
		Bool("redeemed", result.Redeemed).
		Str("booking_id", result.BookingID).
		Str("task_id", result.TaskID).
		Str("error_reason", result.ErrorReason).
		Msg("RedeemPickupQR completed")

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.RedeemPickupQRResponse{
		Redeemed:    result.Redeemed,
		BookingId:   result.BookingID,
		TaskId:      result.TaskID,
		ErrorReason: result.ErrorReason,
	}, nil
}
