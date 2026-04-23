// submit_booking.go — SubmitBooking service method (S2 / BL-BOOK-005).
//
// Transitions a booking from 'draft' → 'pending_payment'.
// This is the CS/admin trigger that locks in the booking and signals readiness
// for invoice issuance. Only 'draft' bookings can be submitted.

package service

import (
	"context"
	"errors"
	"fmt"

	"booking-svc/store/postgres_store/sqlc"
	"booking-svc/util/apperrors"
	"booking-svc/util/logging"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/codes"
)

// SubmitBookingParams is the input for SubmitBooking.
type SubmitBookingParams struct {
	BookingID string
}

// SubmitBookingResult holds the updated booking summary.
type SubmitBookingResult struct {
	ID                 string
	Status             string
	Channel            string
	PackageID          string
	DepartureID        string
	RoomType           string
	LeadFullName       string
	LeadWhatsapp       string
	LeadDomicile       string
	ListAmount         int64
	ListCurrency       string
	SettlementCurrency string
}

// SubmitBooking transitions a draft booking to pending_payment.
// Returns ErrNotFound if the booking does not exist.
// Returns ErrConflict if the booking is not in 'draft' state.
func (s *Service) SubmitBooking(ctx context.Context, params *SubmitBookingParams) (*SubmitBookingResult, error) {
	const op = "service.Service.SubmitBooking"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("booking_id", params.BookingID).Msg("")

	if params.BookingID == "" {
		err := errors.Join(apperrors.ErrValidation, fmt.Errorf("booking_id required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Verify booking exists regardless of status first (to give correct 404 vs 409)
	booking, err := s.store.GetBookingByIDAny(ctx, params.BookingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			notFound := errors.Join(apperrors.ErrNotFound, fmt.Errorf("booking %s not found", params.BookingID))
			span.RecordError(notFound)
			span.SetStatus(codes.Error, notFound.Error())
			return nil, notFound
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("%s: get booking: %w", op, err)
	}

	// Only draft bookings can be submitted
	if booking.Status != sqlc.BookingStatusDraft {
		conflict := errors.Join(apperrors.ErrConflict, fmt.Errorf(
			"booking %s is in state '%s'; only 'draft' bookings can be submitted",
			params.BookingID, booking.Status,
		))
		span.RecordError(conflict)
		span.SetStatus(codes.Error, conflict.Error())
		return nil, conflict
	}

	// Transition draft → pending_payment (optimistic guard)
	err = s.store.UpdateBookingStatus(ctx, sqlc.UpdateBookingStatusParams{
		ID:         params.BookingID,
		NewStatus:  sqlc.BookingStatusPendingPayment,
		FromStatus: sqlc.BookingStatusDraft,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Concurrent update — another request already transitioned it
			conflict := errors.Join(apperrors.ErrConflict, fmt.Errorf(
				"booking %s status changed concurrently; cannot submit", params.BookingID,
			))
			span.RecordError(conflict)
			span.SetStatus(codes.Error, conflict.Error())
			return nil, conflict
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("%s: update status: %w", op, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return &SubmitBookingResult{
		ID:                 booking.ID,
		Status:             string(sqlc.BookingStatusPendingPayment),
		Channel:            string(booking.Channel),
		PackageID:          booking.PackageID,
		DepartureID:        booking.DepartureID,
		RoomType:           booking.RoomType,
		LeadFullName:       booking.LeadFullName,
		LeadWhatsapp:       booking.LeadWhatsapp,
		LeadDomicile:       booking.LeadDomicile,
		ListAmount:         booking.ListAmount,
		ListCurrency:       booking.ListCurrency,
		SettlementCurrency: booking.SettlementCurrency,
	}, nil
}
