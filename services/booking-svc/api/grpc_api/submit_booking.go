// submit_booking.go — gRPC handler for SubmitBooking (S2 / BL-BOOK-005).

package grpc_api

import (
	"context"
	"errors"

	"booking-svc/api/grpc_api/pb"
	"booking-svc/service"
	"booking-svc/util/apperrors"
	"booking-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// SubmitBooking implements pb.SubmitBookingHandler.
// Transitions a booking from 'draft' to 'pending_payment'.
func (s *Server) SubmitBooking(ctx context.Context, req *pb.SubmitBookingRequest) (*pb.SubmitBookingResponse, error) {
	const op = "grpc_api.Server.SubmitBooking"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("booking_id", req.GetBookingId()))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("booking_id", req.GetBookingId()).Msg("")

	result, err := s.svc.SubmitBooking(ctx, &service.SubmitBookingParams{
		BookingID: req.GetBookingId(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, mapBookingSubmitError(err)
	}

	span.SetStatus(codes.Ok, "ok")
	return &pb.SubmitBookingResponse{
		Booking: &pb.BookingSummary{
			Id:                 result.ID,
			Status:             result.Status,
			Channel:            result.Channel,
			PackageId:          result.PackageID,
			DepartureId:        result.DepartureID,
			RoomType:           result.RoomType,
			LeadFullName:       result.LeadFullName,
			LeadWhatsapp:       result.LeadWhatsapp,
			LeadDomicile:       result.LeadDomicile,
			ListAmount:         result.ListAmount,
			ListCurrency:       result.ListCurrency,
			SettlementCurrency: result.SettlementCurrency,
		},
	}, nil
}

func mapBookingSubmitError(err error) error {
	switch {
	case errors.Is(err, apperrors.ErrValidation):
		return grpcStatus.Error(grpcCodes.InvalidArgument, err.Error())
	case errors.Is(err, apperrors.ErrNotFound):
		return grpcStatus.Error(grpcCodes.NotFound, err.Error())
	case errors.Is(err, apperrors.ErrConflict):
		return grpcStatus.Error(grpcCodes.FailedPrecondition, err.Error())
	default:
		return grpcStatus.Error(grpcCodes.Internal, err.Error())
	}
}
