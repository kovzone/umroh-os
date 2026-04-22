package grpc_api

import (
	"context"

	"catalog-svc/api/grpc_api/pb"
	"catalog-svc/service"
	"catalog-svc/util/apperrors"
	"catalog-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

// ListPackages — public paginated list of active packages. Mirrors
// GET /v1/packages on the REST surface (per slice-S1 § Catalog) by
// calling the exact same service-layer method the REST handler calls;
// no business-logic duplication.
func (s *Server) ListPackages(ctx context.Context, req *pb.ListPackagesRequest) (*pb.ListPackagesResponse, error) {
	const op = "grpc_api.Server.ListPackages"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "ListPackages"))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.GetPackages(ctx, &service.GetPackagesParams{
		Kind:          req.GetKind(),
		AirlineCode:   req.GetAirlineCode(),
		HotelID:       req.GetHotelId(),
		DepartureFrom: req.GetDepartureFrom(),
		DepartureTo:   req.GetDepartureTo(),
		Cursor:        req.GetCursor(),
		Limit:         int(req.GetLimit()),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	items := make([]*pb.PackageListItem, 0, len(result.Packages))
	for i := range result.Packages {
		items = append(items, packageListItemToProto(&result.Packages[i]))
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.ListPackagesResponse{
		Packages: items,
		Page: &pb.PageMeta{
			NextCursor: result.NextCursor,
			HasMore:    result.HasMore,
		},
	}, nil
}

// GetPackage — public active-package detail. Mirrors GET /v1/packages/{id}.
func (s *Server) GetPackage(ctx context.Context, req *pb.GetPackageRequest) (*pb.GetPackageResponse, error) {
	const op = "grpc_api.Server.GetPackage"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "GetPackage"))

	logger := logging.LogWithTrace(ctx, s.logger)

	detail, err := s.svc.GetPackageByID(ctx, &service.GetPackageByIDParams{ID: req.GetId()})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.GetPackageResponse{Package: packageDetailToProto(detail)}, nil
}

// GetPackageDeparture — public departure detail with live seat math.
// Mirrors GET /v1/package-departures/{id}.
func (s *Server) GetPackageDeparture(ctx context.Context, req *pb.GetPackageDepartureRequest) (*pb.GetPackageDepartureResponse, error) {
	const op = "grpc_api.Server.GetPackageDeparture"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "GetPackageDeparture"))

	logger := logging.LogWithTrace(ctx, s.logger)

	detail, err := s.svc.GetDepartureByID(ctx, &service.GetDepartureByIDParams{ID: req.GetId()})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.GetPackageDepartureResponse{Departure: departureDetailToProto(detail)}, nil
}

// ---------------------------------------------------------------------------
// service ↔ proto converters. Proto types do not leak past this file; the
// service layer only produces its own Go structs.
// ---------------------------------------------------------------------------

func moneyToProto(m service.Money) *pb.Money {
	return &pb.Money{
		ListAmount:         m.ListAmount,
		ListCurrency:       m.ListCurrency,
		SettlementCurrency: m.SettlementCurrency,
	}
}

func nextDepartureToProto(d *service.NextDeparture) *pb.NextDeparture {
	if d == nil {
		return nil
	}
	return &pb.NextDeparture{
		Id:             d.ID,
		DepartureDate:  d.DepartureDate,
		ReturnDate:     d.ReturnDate,
		RemainingSeats: int32(d.RemainingSeats),
	}
}

func packageListItemToProto(p *service.Package) *pb.PackageListItem {
	return &pb.PackageListItem{
		Id:            p.ID,
		Kind:          p.Kind,
		Name:          p.Name,
		Description:   p.Description,
		CoverPhotoUrl: p.CoverPhotoUrl,
		StartingPrice: moneyToProto(p.StartingPrice),
		NextDeparture: nextDepartureToProto(p.NextDeparture),
	}
}

func hotelRefToProto(h service.HotelRef) *pb.HotelRef {
	return &pb.HotelRef{
		Id:               h.ID,
		Name:             h.Name,
		City:             h.City,
		StarRating:       int32(h.StarRating),
		WalkingDistanceM: int32(h.WalkingDistanceM),
	}
}

func airlineRefToProto(a *service.AirlineRef) *pb.AirlineRef {
	if a == nil {
		return nil
	}
	return &pb.AirlineRef{
		Id:           a.ID,
		Code:         a.Code,
		Name:         a.Name,
		OperatorKind: a.OperatorKind,
	}
}

func muthawwifRefToProto(m *service.MuthawwifRef) *pb.MuthawwifRef {
	if m == nil {
		return nil
	}
	return &pb.MuthawwifRef{
		Id:          m.ID,
		Name:        m.Name,
		PortraitUrl: m.PortraitUrl,
	}
}

func addonRefToProto(a service.AddonRef) *pb.AddonRef {
	return &pb.AddonRef{
		Id:                 a.ID,
		Name:               a.Name,
		ListAmount:         a.ListAmount,
		ListCurrency:       a.ListCurrency,
		SettlementCurrency: a.SettlementCurrency,
	}
}

func itineraryDayToProto(d service.ItineraryDay) *pb.ItineraryDay {
	return &pb.ItineraryDay{
		Day:         int32(d.Day),
		Title:       d.Title,
		Description: d.Description,
		PhotoUrl:    d.PhotoUrl,
	}
}

func itineraryToProto(it *service.Itinerary) *pb.Itinerary {
	if it == nil {
		return nil
	}
	days := make([]*pb.ItineraryDay, 0, len(it.Days))
	for _, d := range it.Days {
		days = append(days, itineraryDayToProto(d))
	}
	return &pb.Itinerary{
		Id:        it.ID,
		Days:      days,
		PublicUrl: it.PublicUrl,
	}
}

func departureSummaryToProto(d service.DepartureSummary) *pb.DepartureSummary {
	return &pb.DepartureSummary{
		Id:             d.ID,
		DepartureDate:  d.DepartureDate,
		ReturnDate:     d.ReturnDate,
		RemainingSeats: int32(d.RemainingSeats),
		Status:         d.Status,
	}
}

func packageDetailToProto(d *service.PackageDetail) *pb.PackageDetail {
	hotels := make([]*pb.HotelRef, 0, len(d.Hotels))
	for _, h := range d.Hotels {
		hotels = append(hotels, hotelRefToProto(h))
	}
	addons := make([]*pb.AddonRef, 0, len(d.Addons))
	for _, a := range d.Addons {
		addons = append(addons, addonRefToProto(a))
	}
	departures := make([]*pb.DepartureSummary, 0, len(d.Departures))
	for _, s := range d.Departures {
		departures = append(departures, departureSummaryToProto(s))
	}
	return &pb.PackageDetail{
		Id:            d.ID,
		Kind:          d.Kind,
		Name:          d.Name,
		Description:   d.Description,
		Highlights:    d.Highlights,
		CoverPhotoUrl: d.CoverPhotoUrl,
		Itinerary:     itineraryToProto(d.Itinerary),
		Hotels:        hotels,
		Airline:       airlineRefToProto(d.Airline),
		Muthawwif:     muthawwifRefToProto(d.Muthawwif),
		AddOns:        addons,
		Departures:    departures,
	}
}

func packagePricingToProto(p service.PackagePricing) *pb.PackagePricing {
	return &pb.PackagePricing{
		RoomType:           p.RoomType,
		ListAmount:         p.ListAmount,
		ListCurrency:       p.ListCurrency,
		SettlementCurrency: p.SettlementCurrency,
	}
}

func vendorReadinessToProto(v service.VendorReadiness) *pb.VendorReadiness {
	return &pb.VendorReadiness{
		Ticket: v.Ticket,
		Hotel:  v.Hotel,
		Visa:   v.Visa,
	}
}

func departureDetailToProto(d *service.DepartureDetail) *pb.DepartureDetail {
	pricing := make([]*pb.PackagePricing, 0, len(d.Pricing))
	for _, p := range d.Pricing {
		pricing = append(pricing, packagePricingToProto(p))
	}
	return &pb.DepartureDetail{
		Id:              d.ID,
		PackageId:       d.PackageID,
		DepartureDate:   d.DepartureDate,
		ReturnDate:      d.ReturnDate,
		TotalSeats:      int32(d.TotalSeats),
		RemainingSeats:  int32(d.RemainingSeats),
		Status:          d.Status,
		Pricing:         pricing,
		VendorReadiness: vendorReadinessToProto(d.VendorReadiness),
	}
}
