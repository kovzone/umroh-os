package grpc_api

// S1-E-03 / BL-BOOK-001..006 — gRPC handler for draft booking creation.
//
// POST /v1/bookings (via gateway-svc → gRPC CreateDraftBooking).
//
// Auth: gateway-svc validates the Bearer token and populates the request
// payload. This handler trusts the incoming gRPC request (internal network
// isolation per ADR-0009 / BL-GTW-100 pending). The service layer calls
// catalog-svc gRPC to validate departure and reserve seats atomically.

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

// CreateDraftBooking implements BookingServiceServer.CreateDraftBooking.
// It delegates to the service layer, which calls catalog-svc to validate
// the departure and reserve seats atomically (BL-BOOK-004).
func (s *Server) CreateDraftBooking(ctx context.Context, req *pb.CreateDraftBookingRequest) (*pb.CreateDraftBookingResponse, error) {
	const op = "grpc_api.Server.CreateDraftBooking"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "CreateDraftBooking"),
		attribute.String("input.channel", req.GetChannel()),
		attribute.String("input.departure_id", req.GetDepartureId()),
	)

	logger.Info().
		Str("op", op).
		Str("channel", req.GetChannel()).
		Str("package_id", req.GetPackageId()).
		Str("departure_id", req.GetDepartureId()).
		Str("idempotency_key", req.GetIdempotencyKey()).
		Msg("")

	// --- Map proto request → service params ---
	params := &service.CreateDraftBookingParams{
		Channel:            req.GetChannel(),
		AgentID:            req.GetAgentId(),
		StaffUserID:        req.GetStaffUserId(),
		PackageID:          req.GetPackageId(),
		DepartureID:        req.GetDepartureId(),
		RoomType:           req.GetRoomType(),
		LeadFullName:       req.GetLeadFullName(),
		LeadEmail:          req.GetLeadEmail(),
		LeadWhatsapp:       req.GetLeadWhatsapp(),
		LeadDomicile:       req.GetLeadDomicile(),
		ListAmount:         req.GetListAmount(),
		ListCurrency:       req.GetListCurrency(),
		SettlementCurrency: req.GetSettlementCurrency(),
		Notes:              req.GetNotes(),
		IdempotencyKey:     req.GetIdempotencyKey(),
	}

	// Map pilgrims
	for _, p := range req.GetPilgrims() {
		params.Pilgrims = append(params.Pilgrims, service.PilgrimInput{
			FullName:    p.GetFullName(),
			Email:       p.GetEmail(),
			Whatsapp:    p.GetWhatsapp(),
			Domicile:    p.GetDomicile(),
			IsLead:      p.GetIsLead(),
			HasKTP:      p.GetHasKtp(),
			HasPassport: p.GetHasPassport(),
		})
	}

	// Map add-ons
	for _, a := range req.GetAddons() {
		params.Addons = append(params.Addons, service.AddonInput{
			AddonID:            a.GetAddonId(),
			AddonName:          a.GetAddonName(),
			ListAmount:         a.GetListAmount(),
			ListCurrency:       a.GetListCurrency(),
			SettlementCurrency: a.GetSettlementCurrency(),
		})
	}

	// Map mahram (BL-BOOK-006: non-blocking)
	if req.GetMahramWarning() != "" {
		params.Mahram = &service.MahramInfo{
			MahramWarning: req.GetMahramWarning(),
		}
	}

	// --- Call service layer ---
	result, err := s.svc.CreateDraftBooking(ctx, params)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(attribute.String("output.booking_id", result.ID))

	// --- Map service result → proto response ---
	items := make([]*pb.BookingItemResult, 0, len(result.Items))
	for _, it := range result.Items {
		items = append(items, &pb.BookingItemResult{
			Id:              it.ID,
			FullName:        it.FullName,
			IsLead:          it.IsLead,
			DocumentWarning: it.DocumentWarning,
		})
	}

	return &pb.CreateDraftBookingResponse{
		Id:                 result.ID,
		Status:             result.Status,
		Channel:            result.Channel,
		PackageId:          result.PackageID,
		DepartureId:        result.DepartureID,
		RoomType:           result.RoomType,
		AgentId:            result.AgentID,
		StaffUserId:        result.StaffUserID,
		LeadFullName:       result.LeadFullName,
		LeadEmail:          result.LeadEmail,
		LeadWhatsapp:       result.LeadWhatsapp,
		LeadDomicile:       result.LeadDomicile,
		ListAmount:         result.ListAmount,
		ListCurrency:       result.ListCurrency,
		SettlementCurrency: result.SettlementCurrency,
		Notes:              result.Notes,
		IdempotencyKey:     result.IdempotencyKey,
		CreatedAt:          result.CreatedAt,
		ExpiresAt:          result.ExpiresAt,
		Items:              items,
		MahramWarning:      result.MahramWarning,
		Replayed:           result.Replayed,
	}, nil
}
