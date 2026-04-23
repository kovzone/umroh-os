package service

// S1-E-03 / BL-BOOK-001..006 / BL-GTW-003 — gateway service methods for
// booking-svc proxy.
//
// CreateDraftBooking forwards the REST request from the gateway handler to
// booking-svc via the booking_grpc_adapter. Per ADR-0009, gateway is the
// single REST entry point; booking-svc exposes only gRPC.

import (
	"context"
	"errors"
	"fmt"

	"gateway-svc/adapter/booking_grpc_adapter"
	"gateway-svc/util/apperrors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// CreateDraftBooking proxies to booking-svc.BookingService/CreateDraftBooking.
func (s *Service) CreateDraftBooking(ctx context.Context, params *booking_grpc_adapter.CreateDraftBookingParams) (*booking_grpc_adapter.CreateDraftBookingResult, error) {
	const op = "service.Service.CreateDraftBooking"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	if s.adapters.bookingGrpc == nil {
		span.SetStatus(codes.Error, "booking adapter not configured")
		return nil, errors.Join(apperrors.ErrServiceUnavailable,
			fmt.Errorf("booking_grpc_adapter not configured"))
	}

	result, err := s.adapters.bookingGrpc.CreateDraftBooking(ctx, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// SubmitBooking proxies to booking-svc.BookingService/SubmitBooking (S2 / BL-BOOK-005).
// Transitions a draft booking to pending_payment.
func (s *Service) SubmitBooking(ctx context.Context, params *booking_grpc_adapter.SubmitBookingParams) (*booking_grpc_adapter.SubmitBookingResult, error) {
	const op = "service.Service.SubmitBooking"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	if s.adapters.bookingGrpc == nil {
		span.SetStatus(codes.Error, "booking adapter not configured")
		return nil, errors.Join(apperrors.ErrServiceUnavailable,
			fmt.Errorf("booking_grpc_adapter not configured"))
	}

	result, err := s.adapters.bookingGrpc.SubmitBooking(ctx, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}
