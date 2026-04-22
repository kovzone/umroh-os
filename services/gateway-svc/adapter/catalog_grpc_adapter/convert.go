package catalog_grpc_adapter

import "gateway-svc/adapter/catalog_grpc_adapter/pb"

// proto → adapter-local converters. Proto types are consumed only here
// and in the RPC wrappers; the rest of gateway-svc sees plain Go
// structs defined in types.go.

func fromProtoMoney(m *pb.Money) Money {
	if m == nil {
		return Money{}
	}
	return Money{
		ListAmount:         m.GetListAmount(),
		ListCurrency:       m.GetListCurrency(),
		SettlementCurrency: m.GetSettlementCurrency(),
	}
}

func fromProtoNextDeparture(n *pb.NextDeparture) *NextDeparture {
	if n == nil {
		return nil
	}
	return &NextDeparture{
		ID:             n.GetId(),
		DepartureDate:  n.GetDepartureDate(),
		ReturnDate:     n.GetReturnDate(),
		RemainingSeats: int(n.GetRemainingSeats()),
	}
}

func fromProtoPackageListItem(p *pb.PackageListItem) Package {
	if p == nil {
		return Package{}
	}
	return Package{
		ID:            p.GetId(),
		Kind:          p.GetKind(),
		Name:          p.GetName(),
		Description:   p.GetDescription(),
		CoverPhotoUrl: p.GetCoverPhotoUrl(),
		StartingPrice: fromProtoMoney(p.GetStartingPrice()),
		NextDeparture: fromProtoNextDeparture(p.GetNextDeparture()),
	}
}

func fromProtoHotelRef(h *pb.HotelRef) HotelRef {
	if h == nil {
		return HotelRef{}
	}
	return HotelRef{
		ID:               h.GetId(),
		Name:             h.GetName(),
		City:             h.GetCity(),
		StarRating:       int(h.GetStarRating()),
		WalkingDistanceM: int(h.GetWalkingDistanceM()),
	}
}

func fromProtoAirlineRef(a *pb.AirlineRef) *AirlineRef {
	if a == nil {
		return nil
	}
	return &AirlineRef{
		ID:           a.GetId(),
		Code:         a.GetCode(),
		Name:         a.GetName(),
		OperatorKind: a.GetOperatorKind(),
	}
}

func fromProtoMuthawwifRef(m *pb.MuthawwifRef) *MuthawwifRef {
	if m == nil {
		return nil
	}
	return &MuthawwifRef{
		ID:          m.GetId(),
		Name:        m.GetName(),
		PortraitUrl: m.GetPortraitUrl(),
	}
}

func fromProtoAddonRef(a *pb.AddonRef) AddonRef {
	if a == nil {
		return AddonRef{}
	}
	return AddonRef{
		ID:                 a.GetId(),
		Name:               a.GetName(),
		ListAmount:         a.GetListAmount(),
		ListCurrency:       a.GetListCurrency(),
		SettlementCurrency: a.GetSettlementCurrency(),
	}
}

func fromProtoItineraryDay(d *pb.ItineraryDay) ItineraryDay {
	if d == nil {
		return ItineraryDay{}
	}
	return ItineraryDay{
		Day:         int(d.GetDay()),
		Title:       d.GetTitle(),
		Description: d.GetDescription(),
		PhotoUrl:    d.GetPhotoUrl(),
	}
}

func fromProtoItinerary(it *pb.Itinerary) *Itinerary {
	if it == nil {
		return nil
	}
	days := make([]ItineraryDay, 0, len(it.GetDays()))
	for _, d := range it.GetDays() {
		days = append(days, fromProtoItineraryDay(d))
	}
	return &Itinerary{
		ID:        it.GetId(),
		Days:      days,
		PublicUrl: it.GetPublicUrl(),
	}
}

func fromProtoDepartureSummary(d *pb.DepartureSummary) DepartureSummary {
	if d == nil {
		return DepartureSummary{}
	}
	return DepartureSummary{
		ID:             d.GetId(),
		DepartureDate:  d.GetDepartureDate(),
		ReturnDate:     d.GetReturnDate(),
		RemainingSeats: int(d.GetRemainingSeats()),
		Status:         d.GetStatus(),
	}
}

func fromProtoPackageDetail(d *pb.PackageDetail) *PackageDetail {
	if d == nil {
		return nil
	}
	hotels := make([]HotelRef, 0, len(d.GetHotels()))
	for _, h := range d.GetHotels() {
		hotels = append(hotels, fromProtoHotelRef(h))
	}
	addons := make([]AddonRef, 0, len(d.GetAddOns()))
	for _, a := range d.GetAddOns() {
		addons = append(addons, fromProtoAddonRef(a))
	}
	departures := make([]DepartureSummary, 0, len(d.GetDepartures()))
	for _, s := range d.GetDepartures() {
		departures = append(departures, fromProtoDepartureSummary(s))
	}
	return &PackageDetail{
		ID:            d.GetId(),
		Kind:          d.GetKind(),
		Name:          d.GetName(),
		Description:   d.GetDescription(),
		Highlights:    d.GetHighlights(),
		CoverPhotoUrl: d.GetCoverPhotoUrl(),
		Itinerary:     fromProtoItinerary(d.GetItinerary()),
		Hotels:        hotels,
		Airline:       fromProtoAirlineRef(d.GetAirline()),
		Muthawwif:     fromProtoMuthawwifRef(d.GetMuthawwif()),
		Addons:        addons,
		Departures:    departures,
	}
}

func fromProtoPackagePricing(p *pb.PackagePricing) PackagePricing {
	if p == nil {
		return PackagePricing{}
	}
	return PackagePricing{
		RoomType:           p.GetRoomType(),
		ListAmount:         p.GetListAmount(),
		ListCurrency:       p.GetListCurrency(),
		SettlementCurrency: p.GetSettlementCurrency(),
	}
}

func fromProtoVendorReadiness(v *pb.VendorReadiness) VendorReadiness {
	if v == nil {
		return VendorReadiness{}
	}
	return VendorReadiness{
		Ticket: v.GetTicket(),
		Hotel:  v.GetHotel(),
		Visa:   v.GetVisa(),
	}
}

func fromProtoDepartureDetail(d *pb.DepartureDetail) *DepartureDetail {
	if d == nil {
		return nil
	}
	pricing := make([]PackagePricing, 0, len(d.GetPricing()))
	for _, p := range d.GetPricing() {
		pricing = append(pricing, fromProtoPackagePricing(p))
	}
	return &DepartureDetail{
		ID:              d.GetId(),
		PackageID:       d.GetPackageId(),
		DepartureDate:   d.GetDepartureDate(),
		ReturnDate:      d.GetReturnDate(),
		TotalSeats:      int(d.GetTotalSeats()),
		RemainingSeats:  int(d.GetRemainingSeats()),
		Status:          d.GetStatus(),
		Pricing:         pricing,
		VendorReadiness: fromProtoVendorReadiness(d.GetVendorReadiness()),
	}
}
