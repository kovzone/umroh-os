// ops_qr.go — gRPC handlers for GenerateIDCard and VerifyIDCard RPCs
// (BL-OPS-003).

package grpc_api

import (
	"context"
	"time"

	"ops-svc/api/grpc_api/pb"
	"ops-svc/service"
	"ops-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// GenerateIDCard handles the GenerateIDCard RPC.
// Creates or refreshes an HMAC-signed token for an ID card or luggage tag.
func (s *Server) GenerateIDCard(ctx context.Context, req *pb.GenerateIDCardRequest) (*pb.GenerateIDCardResponse, error) {
	const op = "grpc_api.Server.GenerateIDCard"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("op", op),
		attribute.String("jamaah_id", req.GetJamaahId()),
		attribute.String("departure_id", req.GetDepartureId()),
		attribute.String("card_type", req.GetCardType()),
	)

	logger.Info().
		Str("op", op).
		Str("jamaah_id", req.GetJamaahId()).
		Str("departure_id", req.GetDepartureId()).
		Str("card_type", req.GetCardType()).
		Msg("")

	if req.GetJamaahId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "jamaah_id is required")
	}
	if req.GetDepartureId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "departure_id is required")
	}
	if req.GetCardType() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "card_type is required")
	}

	result, err := s.svc.GenerateIDCard(ctx, &service.GenerateIDCardParams{
		JamaahID:      req.GetJamaahId(),
		DepartureID:   req.GetDepartureId(),
		CardType:      req.GetCardType(),
		JamaahName:    req.GetJamaahName(),
		DepartureName: req.GetDepartureName(),
	})
	if err != nil {
		logger.Error().Str("op", op).Err(err).Msg("GenerateIDCard failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "generate ID card failed: %v", err)
	}

	span.SetStatus(otelCodes.Ok, "success")
	return &pb.GenerateIDCardResponse{
		Token:    result.Token,
		QrData:   result.QrData,
		IssuedAt: result.IssuedAt.Format(time.RFC3339),
	}, nil
}

// VerifyIDCard handles the VerifyIDCard RPC.
// Verifies an HMAC-signed token; returns valid=false on tamper.
func (s *Server) VerifyIDCard(ctx context.Context, req *pb.VerifyIDCardRequest) (*pb.VerifyIDCardResponse, error) {
	const op = "grpc_api.Server.VerifyIDCard"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("op", op))

	logger.Info().Str("op", op).Msg("")

	if req.GetToken() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "token is required")
	}

	result, err := s.svc.VerifyIDCard(ctx, &service.VerifyIDCardParams{
		Token: req.GetToken(),
	})
	if err != nil {
		logger.Error().Str("op", op).Err(err).Msg("VerifyIDCard failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "verify ID card failed: %v", err)
	}

	span.SetStatus(otelCodes.Ok, "success")
	return &pb.VerifyIDCardResponse{
		Valid:       result.Valid,
		JamaahId:    result.JamaahID,
		DepartureId: result.DepartureID,
		CardType:    result.CardType,
		ErrorReason: result.ErrorReason,
	}, nil
}
