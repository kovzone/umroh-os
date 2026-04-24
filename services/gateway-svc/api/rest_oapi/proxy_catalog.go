package rest_oapi

import (
	"context"
	"errors"
	"time"

	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/api/rest_oapi/middleware"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

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
// Staff catalog write handlers (BL-CAT-014 / S1-E-07)
// ---------------------------------------------------------------------------
//
// Permission gate: every handler extracts the *middleware.Identity from
// c.Locals(middleware.IdentityKey) (set by RequireBearerToken middleware) and
// calls CheckPermission(catalog.package.manage) before delegating to the
// service layer. Failure → 403. catalog-svc also enforces the same check
// server-side as defense-in-depth.

// checkCatalogManagePermission is a shared helper that calls iam-svc.CheckPermission
// for the catalog.package.manage grant. Returns a gating error or nil on allow.
func (s *Server) checkCatalogManagePermission(ctx context.Context, userID string) error {
	perm, err := s.svc.CheckPermission(ctx, &iam_grpc_adapter.CheckPermissionParams{
		UserID:   userID,
		Resource: "catalog.package",
		Action:   "manage",
		Scope:    "global",
	})
	if err != nil {
		return err
	}
	if !perm.Allowed {
		return errors.Join(apperrors.ErrForbidden, errors.New("catalog.package.manage permission required"))
	}
	return nil
}

// CreatePackage implements POST /v1/packages (bearer + catalog.package.manage).
func (s *Server) CreatePackage(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreatePackage"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	id := c.Locals(middleware.IdentityKey).(*middleware.Identity)
	logger.Info().Str("op", op).Str("user_id", id.UserID).Msg("")

	// Permission gate.
	if err := s.checkCatalogManagePermission(ctx, id.UserID); err != nil {
		return writeCatalogError(c, span, err, "forbidden", "akses ditolak: catalog.package.manage diperlukan")
	}

	var body CreatePackageBody
	if err := c.BodyParser(&body); err != nil {
		return writeCatalogError(c, span, errors.Join(apperrors.ErrValidation, err), "invalid_body", "body permintaan tidak valid")
	}

	result, err := s.svc.CreatePackage(ctx, &catalog_grpc_adapter.CreatePackageParams{
		UserID:        id.UserID,
		BranchID:      id.BranchID,
		Kind:          body.Kind,
		Name:          body.Name,
		Description:   body.Description,
		CoverPhotoURL: body.CoverPhotoUrl,
		Highlights:    body.Highlights,
		ItineraryID:   body.ItineraryID,
		AirlineID:     body.AirlineID,
		MuthawwifID:   body.MuthawwifID,
		HotelIDs:      body.HotelIDs,
		AddonIDs:      body.AddonIDs,
		Status:        body.Status,
	})
	if err != nil {
		return writeCatalogWriteError(c, span, err)
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusCreated).JSON(GetPackageResponse{
		Package: packageDetailFromAdapter(result),
	})
}

// UpdatePackageById implements PUT /v1/packages/{id} (bearer + catalog.package.manage).
func (s *Server) UpdatePackageById(c *fiber.Ctx, packageID string) error {
	const op = "rest_oapi.Server.UpdatePackageById"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("id", packageID))

	logger := logging.LogWithTrace(ctx, s.logger)

	id := c.Locals(middleware.IdentityKey).(*middleware.Identity)
	logger.Info().Str("op", op).Str("user_id", id.UserID).Str("package_id", packageID).Msg("")

	if err := s.checkCatalogManagePermission(ctx, id.UserID); err != nil {
		return writeCatalogError(c, span, err, "forbidden", "akses ditolak: catalog.package.manage diperlukan")
	}

	var body UpdatePackageBody
	if err := c.BodyParser(&body); err != nil {
		return writeCatalogError(c, span, errors.Join(apperrors.ErrValidation, err), "invalid_body", "body permintaan tidak valid")
	}

	result, err := s.svc.UpdatePackage(ctx, &catalog_grpc_adapter.UpdatePackageParams{
		UserID:        id.UserID,
		BranchID:      id.BranchID,
		ID:            packageID,
		Name:          body.Name,
		Description:   body.Description,
		CoverPhotoURL: body.CoverPhotoUrl,
		Highlights:    body.Highlights,
		ItineraryID:   body.ItineraryID,
		AirlineID:     body.AirlineID,
		MuthawwifID:   body.MuthawwifID,
		HotelIDs:      body.HotelIDs,
		AddonIDs:      body.AddonIDs,
		Status:        body.Status,
	})
	if err != nil {
		return writeCatalogWriteError(c, span, err)
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(GetPackageResponse{
		Package: packageDetailFromAdapter(result),
	})
}

// DeletePackageById implements DELETE /v1/packages/{id} (bearer + catalog.package.manage).
func (s *Server) DeletePackageById(c *fiber.Ctx, packageID string) error {
	const op = "rest_oapi.Server.DeletePackageById"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("id", packageID))

	logger := logging.LogWithTrace(ctx, s.logger)

	id := c.Locals(middleware.IdentityKey).(*middleware.Identity)
	logger.Info().Str("op", op).Str("user_id", id.UserID).Str("package_id", packageID).Msg("")

	if err := s.checkCatalogManagePermission(ctx, id.UserID); err != nil {
		return writeCatalogError(c, span, err, "forbidden", "akses ditolak: catalog.package.manage diperlukan")
	}

	_, err := s.svc.DeletePackage(ctx, &catalog_grpc_adapter.DeletePackageParams{
		UserID:   id.UserID,
		BranchID: id.BranchID,
		ID:       packageID,
	})
	if err != nil {
		return writeCatalogWriteError(c, span, err)
	}

	span.SetStatus(codes.Ok, "success")
	return c.SendStatus(fiber.StatusNoContent)
}

// CreateDeparture implements POST /v1/packages/{id}/departures (bearer + catalog.package.manage).
func (s *Server) CreateDeparture(c *fiber.Ctx, packageID string) error {
	const op = "rest_oapi.Server.CreateDeparture"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("package_id", packageID))

	logger := logging.LogWithTrace(ctx, s.logger)

	id := c.Locals(middleware.IdentityKey).(*middleware.Identity)
	logger.Info().Str("op", op).Str("user_id", id.UserID).Str("package_id", packageID).Msg("")

	if err := s.checkCatalogManagePermission(ctx, id.UserID); err != nil {
		return writeCatalogError(c, span, err, "forbidden", "akses ditolak: catalog.package.manage diperlukan")
	}

	var body CreateDepartureBody
	if err := c.BodyParser(&body); err != nil {
		return writeCatalogError(c, span, errors.Join(apperrors.ErrValidation, err), "invalid_body", "body permintaan tidak valid")
	}

	pricing := make([]catalog_grpc_adapter.PricingInputParam, 0, len(body.Pricing))
	for _, p := range body.Pricing {
		pricing = append(pricing, catalog_grpc_adapter.PricingInputParam{
			RoomType:           p.RoomType,
			ListAmount:         p.ListAmount,
			ListCurrency:       p.ListCurrency,
			SettlementCurrency: p.SettlementCurrency,
		})
	}

	result, err := s.svc.CreateDeparture(ctx, &catalog_grpc_adapter.CreateDepartureParams{
		UserID:        id.UserID,
		BranchID:      id.BranchID,
		PackageID:     packageID,
		DepartureDate: body.DepartureDate,
		ReturnDate:    body.ReturnDate,
		TotalSeats:    body.TotalSeats,
		Status:        body.Status,
		Pricing:       pricing,
	})
	if err != nil {
		return writeCatalogWriteError(c, span, err)
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusCreated).JSON(GetDepartureResponse{
		Departure: departureDetailFromAdapter(result),
	})
}

// UpdateDepartureById implements PUT /v1/departures/{id} (bearer + catalog.package.manage).
func (s *Server) UpdateDepartureById(c *fiber.Ctx, departureID string) error {
	const op = "rest_oapi.Server.UpdateDepartureById"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("id", departureID))

	logger := logging.LogWithTrace(ctx, s.logger)

	id := c.Locals(middleware.IdentityKey).(*middleware.Identity)
	logger.Info().Str("op", op).Str("user_id", id.UserID).Str("departure_id", departureID).Msg("")

	if err := s.checkCatalogManagePermission(ctx, id.UserID); err != nil {
		return writeCatalogError(c, span, err, "forbidden", "akses ditolak: catalog.package.manage diperlukan")
	}

	var body UpdateDepartureBody
	if err := c.BodyParser(&body); err != nil {
		return writeCatalogError(c, span, errors.Join(apperrors.ErrValidation, err), "invalid_body", "body permintaan tidak valid")
	}

	var pricing []catalog_grpc_adapter.PricingInputParam
	if len(body.Pricing) > 0 {
		pricing = make([]catalog_grpc_adapter.PricingInputParam, 0, len(body.Pricing))
		for _, p := range body.Pricing {
			pricing = append(pricing, catalog_grpc_adapter.PricingInputParam{
				RoomType:           p.RoomType,
				ListAmount:         p.ListAmount,
				ListCurrency:       p.ListCurrency,
				SettlementCurrency: p.SettlementCurrency,
			})
		}
	}

	result, err := s.svc.UpdateDeparture(ctx, &catalog_grpc_adapter.UpdateDepartureParams{
		UserID:        id.UserID,
		BranchID:      id.BranchID,
		ID:            departureID,
		DepartureDate: body.DepartureDate,
		ReturnDate:    body.ReturnDate,
		TotalSeats:    body.TotalSeats,
		Status:        body.Status,
		Pricing:       pricing,
	})
	if err != nil {
		return writeCatalogWriteError(c, span, err)
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(GetDepartureResponse{
		Departure: departureDetailFromAdapter(result),
	})
}

// writeCatalogWriteError handles errors from write operations, extending
// mapGatewayCatalogError with an additional forbidden case.
func writeCatalogWriteError(c *fiber.Ctx, span trace.Span, err error) error {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	var status int
	var code, message string
	switch {
	case errors.Is(err, apperrors.ErrForbidden):
		status = fiber.StatusForbidden
		code = "forbidden"
		message = "akses ditolak: izin tidak mencukupi"
	case errors.Is(err, apperrors.ErrNotFound):
		status = fiber.StatusNotFound
		code = "not_found"
		message = "sumber daya tidak ditemukan"
	case errors.Is(err, apperrors.ErrValidation):
		status = fiber.StatusBadRequest
		code = "validation_error"
		message = "data tidak valid: " + err.Error()
	case errors.Is(err, apperrors.ErrConflict):
		status = fiber.StatusConflict
		code = "conflict"
		message = "konflik data: " + err.Error()
	case errors.Is(err, apperrors.ErrServiceUnavailable):
		status = fiber.StatusBadGateway
		code = "service_unavailable"
		message = "layanan katalog sementara tidak tersedia"
	default:
		status = fiber.StatusInternalServerError
		code = "internal_error"
		message = "terjadi kesalahan tidak terduga"
	}

	resp := CatalogErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	if sc := span.SpanContext(); sc.HasTraceID() {
		tid := sc.TraceID().String()
		resp.Error.TraceId = &tid
	}
	return c.Status(status).JSON(resp)
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

func writeInvalidCatalogQueryParam(c *fiber.Ctx, span trace.Span, logger *zerolog.Logger, paramName string) error {
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
	ds := DepartureSummary{
		Id:             d.ID,
		DepartureDate:  mustDate(d.DepartureDate),
		ReturnDate:     mustDate(d.ReturnDate),
		RemainingSeats: d.RemainingSeats,
		Status:         DepartureStatus(d.Status),
	}
	if d.PricePerPax != nil {
		ds.PricePerPax = d.PricePerPax
	}
	return ds
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

