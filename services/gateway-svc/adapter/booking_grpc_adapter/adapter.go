// Package booking_grpc_adapter is gateway-svc's consumer-side wrapper
// around booking-svc's gRPC surface (BookingService).
//
// Per ADR-0009, gateway's booking route (POST /v1/bookings) proxies to
// booking-svc.BookingService via this adapter. The wire contract is in
// pb/booking.go, kept in sync by hand with
// services/booking-svc/api/grpc_api/pb/booking.proto.
//
// Landed with BL-GTW-003 / S1-E-03.
package booking_grpc_adapter

import (
	"gateway-svc/adapter/booking_grpc_adapter/pb"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Adapter is a thin wrapper over a booking.v1.BookingService client.
type Adapter struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	bookingClient pb.BookingServiceClient
}

// NewAdapter creates a new booking-svc gRPC adapter from an already-dialled
// conn. Ownership of the conn stays with the caller (shared pool lifetime).
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:        logger,
		tracer:        tracer,
		bookingClient: pb.NewBookingServiceClient(cc),
	}
}
