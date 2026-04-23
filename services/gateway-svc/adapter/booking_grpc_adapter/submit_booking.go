// submit_booking.go — gateway-side adapter method for
// booking-svc.BookingService/SubmitBooking (S2 / BL-BOOK-005).
//
// Translates SubmitBookingParams → pb request, forwards via gRPC,
// translates pb response → SubmitBookingResult.
package booking_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/booking_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// SubmitBooking transitions a draft booking to pending_payment by calling
// booking-svc.BookingService/SubmitBooking over gRPC.
func (a *Adapter) SubmitBooking(ctx context.Context, params *SubmitBookingParams) (*SubmitBookingResult, error) {
	const op = "booking_grpc_adapter.Adapter.SubmitBooking"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "SubmitBooking"),
		attribute.String("booking_id", params.BookingID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.bookingClient.SubmitBooking(ctx, &pb.SubmitBookingRequest{
		BookingId: params.BookingID,
	})
	if err != nil {
		wrapped := mapBookingError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	b := resp.GetBooking()
	if b == nil {
		b = &pb.BookingSummary{}
	}
	return &SubmitBookingResult{
		ID:                 b.GetId(),
		Status:             b.GetStatus(),
		Channel:            b.GetChannel(),
		PackageID:          b.GetPackageId(),
		DepartureID:        b.GetDepartureId(),
		RoomType:           b.GetRoomType(),
		LeadFullName:       b.GetLeadFullName(),
		LeadWhatsapp:       b.GetLeadWhatsapp(),
		LeadDomicile:       b.GetLeadDomicile(),
		ListAmount:         b.GetListAmount(),
		ListCurrency:       b.GetListCurrency(),
		SettlementCurrency: b.GetSettlementCurrency(),
	}, nil
}
