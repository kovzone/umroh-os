// seats_by_channel.go — gateway booking adapter method for cross-channel seat
// tracking RPC (BL-BOOK-007).

package booking_grpc_adapter

import (
	"context"
	"fmt"

	"gateway-svc/adapter/booking_grpc_adapter/pb"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// Result types
// ---------------------------------------------------------------------------

// ChannelSeatRow is the per-channel breakdown in GetSeatsByChannel result.
type ChannelSeatRow struct {
	Channel      string
	Seats        int
	BookingCount int
}

// GetSeatsByChannelResult is the gateway-layer result for GetSeatsByChannel.
type GetSeatsByChannelResult struct {
	DepartureID string
	ByChannel   []ChannelSeatRow
}

// ---------------------------------------------------------------------------
// Error mapper
// ---------------------------------------------------------------------------

func mapBookingChannelError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return apperrors.ErrInternal
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return fmt.Errorf("%w: %s", apperrors.ErrNotFound, st.Message())
	case grpcCodes.InvalidArgument:
		return fmt.Errorf("%w: %s", apperrors.ErrValidation, st.Message())
	case grpcCodes.PermissionDenied:
		return fmt.Errorf("%w: %s", apperrors.ErrForbidden, st.Message())
	default:
		return apperrors.ErrInternal
	}
}

// ---------------------------------------------------------------------------
// BL-BOOK-007: GetSeatsByChannel
// ---------------------------------------------------------------------------

// GetSeatsByChannel returns the seat count per booking channel for a departure.
func (a *Adapter) GetSeatsByChannel(ctx context.Context, departureID string) (*GetSeatsByChannelResult, error) {
	const op = "booking_grpc_adapter.Adapter.GetSeatsByChannel"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("departure_id", departureID).Msg("")

	resp, err := a.channelClient.GetSeatsByChannel(ctx, &pb.GetSeatsByChannelRequest{
		DepartureId: departureID,
	})
	if err != nil {
		wrapped := mapBookingChannelError(err)
		logger.Warn().Err(wrapped).Msg("booking-svc.GetSeatsByChannel failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	rows := make([]ChannelSeatRow, 0, len(resp.GetByChannel()))
	for _, r := range resp.GetByChannel() {
		rows = append(rows, ChannelSeatRow{
			Channel:      r.GetChannel(),
			Seats:        int(r.GetSeats()),
			BookingCount: int(r.GetBookingCount()),
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &GetSeatsByChannelResult{
		DepartureID: resp.GetDepartureId(),
		ByChannel:   rows,
	}, nil
}
