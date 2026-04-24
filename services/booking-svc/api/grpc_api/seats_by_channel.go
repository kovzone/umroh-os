package grpc_api

// BL-BOOK-007 — gRPC handler for GetSeatsByChannel.

import (
	"context"

	"booking-svc/api/grpc_api/pb"
	"booking-svc/service"
	"booking-svc/util/apperrors"
	"booking-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

// GetSeatsByChannel returns per-channel seat statistics for a departure.
func (s *Server) GetSeatsByChannel(ctx context.Context, req *pb.GetSeatsByChannelRequest) (*pb.GetSeatsByChannelResponse, error) {
	const op = "grpc_api.Server.GetSeatsByChannel"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "GetSeatsByChannel"),
		attribute.String("input.departure_id", req.GetDepartureId()),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.GetSeatsByChannel(ctx, &service.GetSeatsByChannelParams{
		DepartureID: req.GetDepartureId(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	byChannel := make([]*pb.ChannelSeatRow, 0, len(result.ByChannel))
	for _, entry := range result.ByChannel {
		byChannel = append(byChannel, &pb.ChannelSeatRow{
			Channel:      entry.Channel,
			Seats:        int32(entry.Seats),
			BookingCount: int32(entry.BookingCount),
		})
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.GetSeatsByChannelResponse{
		DepartureId: result.DepartureID,
		ByChannel:   byChannel,
	}, nil
}
