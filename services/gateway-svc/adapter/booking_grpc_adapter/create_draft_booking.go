// create_draft_booking.go — gateway-side adapter method for
// booking-svc.BookingService/CreateDraftBooking.
//
// S1-E-03 / BL-BOOK-001..006 / BL-GTW-003.
//
// Translates gateway-local params → pb request, forwards via gRPC,
// and translates pb response → adapter-local types. Proto types do not
// leak past this package; the rest of gateway-svc sees plain Go structs.
package booking_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/booking_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// CreateDraftBooking forwards a draft-booking request to booking-svc
// over gRPC and returns the created booking data.
func (a *Adapter) CreateDraftBooking(ctx context.Context, params *CreateDraftBookingParams) (*CreateDraftBookingResult, error) {
	const op = "booking_grpc_adapter.Adapter.CreateDraftBooking"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "CreateDraftBooking"),
		attribute.String("channel", params.Channel),
		attribute.String("departure_id", params.DepartureID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	// Map pilgrims
	pbPilgrims := make([]*pb.PilgrimInput, 0, len(params.Pilgrims))
	for _, p := range params.Pilgrims {
		pbPilgrims = append(pbPilgrims, &pb.PilgrimInput{
			FullName:    p.FullName,
			Email:       p.Email,
			Whatsapp:    p.Whatsapp,
			Domicile:    p.Domicile,
			IsLead:      p.IsLead,
			HasKtp:      p.HasKTP,
			HasPassport: p.HasPassport,
		})
	}

	// Map add-ons
	pbAddons := make([]*pb.AddonInput, 0, len(params.Addons))
	for _, a := range params.Addons {
		pbAddons = append(pbAddons, &pb.AddonInput{
			AddonId:            a.AddonID,
			AddonName:          a.AddonName,
			ListAmount:         a.ListAmount,
			ListCurrency:       a.ListCurrency,
			SettlementCurrency: a.SettlementCurrency,
		})
	}

	resp, err := a.bookingClient.CreateDraftBooking(ctx, &pb.CreateDraftBookingRequest{
		Channel:            params.Channel,
		AgentId:            params.AgentID,
		StaffUserId:        params.StaffUserID,
		PackageId:          params.PackageID,
		DepartureId:        params.DepartureID,
		RoomType:           params.RoomType,
		LeadFullName:       params.LeadFullName,
		LeadEmail:          params.LeadEmail,
		LeadWhatsapp:       params.LeadWhatsapp,
		LeadDomicile:       params.LeadDomicile,
		Pilgrims:           pbPilgrims,
		Addons:             pbAddons,
		ListAmount:         params.ListAmount,
		ListCurrency:       params.ListCurrency,
		SettlementCurrency: params.SettlementCurrency,
		Notes:              params.Notes,
		IdempotencyKey:     params.IdempotencyKey,
		MahramWarning:      params.MahramWarning,
	})
	if err != nil {
		wrapped := mapBookingError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(attribute.String("output.booking_id", resp.GetId()))

	// Map items
	items := make([]BookingItemResult, 0, len(resp.GetItems()))
	for _, it := range resp.GetItems() {
		items = append(items, BookingItemResult{
			ID:              it.GetId(),
			FullName:        it.GetFullName(),
			IsLead:          it.GetIsLead(),
			DocumentWarning: it.GetDocumentWarning(),
		})
	}

	return &CreateDraftBookingResult{
		ID:                 resp.GetId(),
		Status:             resp.GetStatus(),
		Channel:            resp.GetChannel(),
		PackageID:          resp.GetPackageId(),
		DepartureID:        resp.GetDepartureId(),
		RoomType:           resp.GetRoomType(),
		AgentID:            resp.GetAgentId(),
		StaffUserID:        resp.GetStaffUserId(),
		LeadFullName:       resp.GetLeadFullName(),
		LeadEmail:          resp.GetLeadEmail(),
		LeadWhatsapp:       resp.GetLeadWhatsapp(),
		LeadDomicile:       resp.GetLeadDomicile(),
		ListAmount:         resp.GetListAmount(),
		ListCurrency:       resp.GetListCurrency(),
		SettlementCurrency: resp.GetSettlementCurrency(),
		Notes:              resp.GetNotes(),
		IdempotencyKey:     resp.GetIdempotencyKey(),
		CreatedAt:          resp.GetCreatedAt(),
		ExpiresAt:          resp.GetExpiresAt(),
		Items:              items,
		MahramWarning:      resp.GetMahramWarning(),
		Replayed:           resp.GetReplayed(),
	}, nil
}
