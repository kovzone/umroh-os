package service

// BL-BOOK-007 — Cross-channel seat tracking.
//
// GetSeatsByChannel queries the bookings table grouped by channel for a given
// departure, returning seats and booking_count per channel. Cancelled, failed,
// and expired bookings are excluded.

import (
	"context"
	"errors"
	"fmt"

	"booking-svc/store/postgres_store"
	"booking-svc/util/apperrors"
	"booking-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetSeatsByChannelParams is the input for GetSeatsByChannel.
type GetSeatsByChannelParams struct {
	DepartureID string
}

// ChannelSeatEntry carries seat metrics for one booking channel.
type ChannelSeatEntry struct {
	Channel      string
	Seats        int
	BookingCount int
}

// GetSeatsByChannelResult is the response for GetSeatsByChannel.
type GetSeatsByChannelResult struct {
	DepartureID string
	ByChannel   []ChannelSeatEntry
}

// GetSeatsByChannel returns per-channel seat statistics for a departure.
func (s *Service) GetSeatsByChannel(ctx context.Context, params *GetSeatsByChannelParams) (*GetSeatsByChannelResult, error) {
	const op = "service.Service.GetSeatsByChannel"

	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.departure_id", params.DepartureID),
	)
	logger.Info().Str("op", op).Str("departure_id", params.DepartureID).Msg("")

	if params.DepartureID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("departure_id is required"))
	}

	rows, err := s.store.GetSeatsByChannelForDeparture(ctx, params.DepartureID)
	if err != nil {
		return nil, fmt.Errorf("get seats by channel: %w", postgres_store.WrapDBError(err))
	}

	entries := make([]ChannelSeatEntry, 0, len(rows))
	for _, r := range rows {
		entries = append(entries, ChannelSeatEntry{
			Channel:      string(r.Channel),
			Seats:        int(r.Seats),
			BookingCount: int(r.BookingCount),
		})
	}

	span.SetStatus(codes.Ok, "success")
	return &GetSeatsByChannelResult{
		DepartureID: params.DepartureID,
		ByChannel:   entries,
	}, nil
}
