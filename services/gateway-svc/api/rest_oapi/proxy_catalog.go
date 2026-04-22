package rest_oapi

import (
	"errors"
	"time"

	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// GetCatalogSystemLive proxies catalog-svc's liveness probe. Retires
// with BL-REFACTOR-001 / S1-E-11 when catalog-svc drops its REST port.
func (s *Server) GetCatalogSystemLive(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetCatalogSystemLive"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/catalog/system/live"),
		attribute.String("method", "GET"),
	)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetCatalogSystemLive(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "BAD_GATEWAY",
				"message": err.Error(),
			},
		})
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(LiveResponse{
		Data: struct {
			Ok bool `json:"ok"`
		}{Ok: result.OK},
	})
}

// ListPackages proxies GET /v1/packages to catalog-svc.ListPackages over
// gRPC via catalog_grpc_adapter. Handler-level validation mirrors
// catalog-svc's REST handler so 400 shapes are identical.
func (s *Server) ListPackages(c *fiber.Ctx, params ListPackagesParams) error {
	const op = "rest_oapi.Server.ListPackages"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("operation", op))

	// --- Handler-level query-param validation (mirrors catalog's). -----
	var kind string
	if params.Kind != nil {
		if !isKnownPackageKind(*params.Kind) {
			return writeInvalidCatalogQueryParam(c, span, logger, "kind")
		}
		kind = string(*params.Kind)
	}
	var airlineCode string
	if params.AirlineCode != nil {
		airlineCode = *params.AirlineCode
	}
	var hotelID string
	if params.HotelId != nil {
		hotelID = *params.HotelId
	}
	var departureFrom string
	if params.DepartureFrom != nil {
		departureFrom = dateToISO(*params.DepartureFrom)
	}
	var departureTo string
	if params.DepartureTo != nil {
		departureTo = dateToISO(*params.DepartureTo)
	}
	var cursor string
	if params.Cursor != nil {
		cursor = *params.Cursor
	}
	limit := 0
	if params.Limit != nil {
		if *params.Limit < 1 || *params.Limit > 100 {
			return writeInvalidCatalogQueryParam(c, span, logger, "limit")
		}
		limit = *params.Limit
	}

	// --- Call service layer -------------------------------------------
	result, err := s.svc.ListPackages(ctx, &catalog_grpc_adapter.ListPackagesParams{
		Kind:          kind,
		AirlineCode:   airlineCode,
		HotelID:       hotelID,
		DepartureFrom: departureFrom,
		DepartureTo:   departureTo,
		Cursor:        cursor,
		Limit:         limit,
	})
	if err != nil {
		return writeCatalogError(c, span, err, "package_not_found", "paket tidak ditemukan")
	}

	// --- Marshal adapter types → oapi response ------------------------
	items := make([]PackageListItem, 0, len(result.Packages))
	for i := range result.Packages {
		items = append(items, packageListItemFromAdapter(&result.Packages[i]))
	}
	page := PageMeta{HasMore: result.Page.HasMore}
	if result.Page.NextCursor != "" {
		nc := result.Page.NextCursor
		page.NextCursor = &nc
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(ListPackagesResponse{
		Packages: items,
		Page:     page,
	})
}

// GetPackageById proxies GET /v1/packages/{id}.
func (s *Server) GetPackageById(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.GetPackageById"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("operation", op), attribute.String("id", id))
	logger.Info().Str("op", op).Str("id", id).Msg("")

	detail, err := s.svc.GetPackage(ctx, &catalog_grpc_adapter.GetPackageParams{ID: id})
	if err != nil {
		return writeCatalogError(c, span, err, "package_not_found", "paket tidak ditemukan")
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(GetPackageResponse{
		Package: packageDetailFromAdapter(detail),
	})
}

// GetPackageDepartureById proxies GET /v1/package-departures/{id}.
func (s *Server) GetPackageDepartureById(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.GetPackageDepartureById"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("operation", op), attribute.String("id", id))
	logger.Info().Str("op", op).Str("id", id).Msg("")

	detail, err := s.svc.GetPackageDeparture(ctx, &catalog_grpc_adapter.GetPackageDepartureParams{ID: id})
	if err != nil {
		return writeCatalogError(c, span, err, "departure_not_found", "keberangkatan tidak ditemukan")
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(GetDepartureResponse{
		Departure: departureDetailFromAdapter(detail),
	})
}

// ---------------------------------------------------------------------------
// Catalog error envelope (snake_case + id-ID + trace_id) — see the matching
// writer in services/catalog-svc/api/rest_oapi/packages.go. Gateway's
// default UPPER_SNAKE error middleware is bypassed for catalog routes so
// the e2e spec reads the same body shape whether it hits catalog-svc:4002
// (transitional) or gateway-svc:4000.
// ---------------------------------------------------------------------------

func writeCatalogError(c *fiber.Ctx, span trace.Span, err error, notFoundCode, notFoundMessage string) error {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	status, code, message := mapGatewayCatalogError(err, notFoundCode, notFoundMessage)

	resp := CatalogErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	if sc := span.SpanContext(); sc.HasTraceID() {
		tid := sc.TraceID().String()
		resp.Error.TraceId = &tid
	}
	return c.Status(status).JSON(resp)
}

func writeInvalidCatalogQueryParam(c *fiber.Ctx, span trace.Span, logger zerolog.Logger, paramName string) error {
	logger.Warn().Str("param", paramName).Msg("invalid query param")
	span.SetAttributes(attribute.String("invalid_query_param", paramName))
	span.SetStatus(codes.Error, "invalid_query_param")

	resp := CatalogErrorResponse{}
	resp.Error.Code = "invalid_query_param"
	resp.Error.Message = "parameter kueri tidak valid: " + paramName
	if sc := span.SpanContext(); sc.HasTraceID() {
		tid := sc.TraceID().String()
		resp.Error.TraceId = &tid
	}
	return c.Status(fiber.StatusBadRequest).JSON(resp)
}

func mapGatewayCatalogError(err error, notFoundCode, notFoundMessage string) (int, string, string) {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		return fiber.StatusNotFound, notFoundCode, notFoundMessage
	case errors.Is(err, apperrors.ErrValidation):
		return fiber.StatusBadRequest, "invalid_cursor", "cursor tidak valid"
	case errors.Is(err, apperrors.ErrServiceUnavailable):
		return fiber.StatusBadGateway, "service_unavailable", "layanan katalog sementara tidak tersedia"
	default:
		return fiber.StatusInternalServerError, "internal_error", "terjadi kesalahan tidak terduga"
	}
}

// ---------------------------------------------------------------------------
// adapter → oapi-generated type converters
// ---------------------------------------------------------------------------

func moneyFromAdapter(m catalog_grpc_adapter.Money) Money {
	return Money{
		ListAmount:         m.ListAmount,
		ListCurrency:       m.ListCurrency,
		SettlementCurrency: MoneySettlementCurrency(m.SettlementCurrency),
	}
}

func nextDepartureFromAdapter(n *catalog_grpc_adapter.NextDeparture) *NextDeparture {
	if n == nil {
		return nil
	}
	return &NextDeparture{
		Id:             n.ID,
		DepartureDate:  mustDate(n.DepartureDate),
		ReturnDate:     mustDate(n.ReturnDate),
		RemainingSeats: n.RemainingSeats,
	}
}

func packageListItemFromAdapter(p *catalog_grpc_adapter.Package) PackageListItem {
	return PackageListItem{
		Id:             p.ID,
		Kind:           PackageKind(p.Kind),
		Name:           p.Name,
		Description:    p.Description,
		CoverPhotoUrl:  p.CoverPhotoUrl,
		StartingPrice:  moneyFromAdapter(p.StartingPrice),
		NextDeparture:  nextDepartureFromAdapter(p.NextDeparture),
	}
}

func hotelRefFromAdapter(h catalog_grpc_adapter.HotelRef) HotelRef {
	return HotelRef{
		Id:               h.ID,
		Name:             h.Name,
		City:             h.City,
		StarRating:       h.StarRating,
		WalkingDistanceM: h.WalkingDistanceM,
	}
}

func airlineRefFromAdapter(a *catalog_grpc_adapter.AirlineRef) *AirlineRef {
	if a == nil {
		return nil
	}
	return &AirlineRef{
		Id:           a.ID,
		Code:         a.Code,
		Name:         a.Name,
		OperatorKind: OperatorKind(a.OperatorKind),
	}
}

func muthawwifRefFromAdapter(m *catalog_grpc_adapter.MuthawwifRef) *MuthawwifRef {
	if m == nil {
		return nil
	}
	return &MuthawwifRef{
		Id:          m.ID,
		Name:        m.Name,
		PortraitUrl: m.PortraitUrl,
	}
}

func addonRefFromAdapter(a catalog_grpc_adapter.AddonRef) AddonRef {
	return AddonRef{
		Id:                 a.ID,
		Name:               a.Name,
		ListAmount:         a.ListAmount,
		ListCurrency:       a.ListCurrency,
		SettlementCurrency: AddonRefSettlementCurrency(a.SettlementCurrency),
	}
}

func itineraryDayFromAdapter(d catalog_grpc_adapter.ItineraryDay) ItineraryDay {
	day := ItineraryDay{
		Day:         d.Day,
		Title:       d.Title,
		Description: d.Description,
	}
	if d.PhotoUrl != "" {
		url := d.PhotoUrl
		day.PhotoUrl = &url
	}
	return day
}

func itineraryFromAdapter(it *catalog_grpc_adapter.Itinerary) *Itinerary {
	if it == nil {
		return nil
	}
	days := make([]ItineraryDay, 0, len(it.Days))
	for _, d := range it.Days {
		days = append(days, itineraryDayFromAdapter(d))
	}
	return &Itinerary{
		Id:        it.ID,
		Days:      days,
		PublicUrl: it.PublicUrl,
	}
}

func departureSummaryFromAdapter(d catalog_grpc_adapter.DepartureSummary) DepartureSummary {
	return DepartureSummary{
		Id:             d.ID,
		DepartureDate:  mustDate(d.DepartureDate),
		ReturnDate:     mustDate(d.ReturnDate),
		RemainingSeats: d.RemainingSeats,
		Status:         DepartureStatus(d.Status),
	}
}

func packageDetailFromAdapter(d *catalog_grpc_adapter.PackageDetail) PackageDetail {
	hotels := make([]HotelRef, 0, len(d.Hotels))
	for _, h := range d.Hotels {
		hotels = append(hotels, hotelRefFromAdapter(h))
	}
	addons := make([]AddonRef, 0, len(d.Addons))
	for _, a := range d.Addons {
		addons = append(addons, addonRefFromAdapter(a))
	}
	departures := make([]DepartureSummary, 0, len(d.Departures))
	for _, s := range d.Departures {
		departures = append(departures, departureSummaryFromAdapter(s))
	}
	return PackageDetail{
		Id:            d.ID,
		Kind:          PackageKind(d.Kind),
		Name:          d.Name,
		Description:   d.Description,
		Highlights:    d.Highlights,
		CoverPhotoUrl: d.CoverPhotoUrl,
		Itinerary:     itineraryFromAdapter(d.Itinerary),
		Hotels:        hotels,
		Airline:       airlineRefFromAdapter(d.Airline),
		Muthawwif:     muthawwifRefFromAdapter(d.Muthawwif),
		AddOns:        addons,
		Departures:    departures,
	}
}

func packagePricingFromAdapter(p catalog_grpc_adapter.PackagePricing) PackagePricing {
	return PackagePricing{
		RoomType:           RoomType(p.RoomType),
		ListAmount:         p.ListAmount,
		ListCurrency:       p.ListCurrency,
		SettlementCurrency: PackagePricingSettlementCurrency(p.SettlementCurrency),
	}
}

func vendorReadinessFromAdapter(v catalog_grpc_adapter.VendorReadiness) VendorReadiness {
	return VendorReadiness{
		Ticket: VendorReadinessState(v.Ticket),
		Hotel:  VendorReadinessState(v.Hotel),
		Visa:   VendorReadinessState(v.Visa),
	}
}

func departureDetailFromAdapter(d *catalog_grpc_adapter.DepartureDetail) DepartureDetail {
	pricing := make([]PackagePricing, 0, len(d.Pricing))
	for _, p := range d.Pricing {
		pricing = append(pricing, packagePricingFromAdapter(p))
	}
	return DepartureDetail{
		Id:              d.ID,
		PackageId:       d.PackageID,
		DepartureDate:   mustDate(d.DepartureDate),
		ReturnDate:      mustDate(d.ReturnDate),
		TotalSeats:      d.TotalSeats,
		RemainingSeats:  d.RemainingSeats,
		Status:          DepartureStatus(d.Status),
		Pricing:         pricing,
		VendorReadiness: vendorReadinessFromAdapter(d.VendorReadiness),
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// mustDate parses an ISO YYYY-MM-DD string into openapi_types.Date. The
// adapter only receives backend-produced dates; a parse failure means the
// backend broke its contract, so we zero-value rather than returning an
// error to the handler (the backend will have logged the malformed row).
func mustDate(s string) openapi_types.Date {
	if s == "" {
		return openapi_types.Date{}
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return openapi_types.Date{}
	}
	return openapi_types.Date{Time: t}
}

func dateToISO(d openapi_types.Date) string {
	if d.Time.IsZero() {
		return ""
	}
	return d.Time.Format("2006-01-02")
}

// isKnownPackageKind accepts any PackageKind enum value from the generated
// oapi types. Mirrors catalog-svc's whitelist.
func isKnownPackageKind(k PackageKind) bool {
	switch k {
	case UmrahReguler, UmrahPlus, HajjFuroda, HajjKhusus, Badal, Financial, Retail:
		return true
	default:
		return false
	}
}

