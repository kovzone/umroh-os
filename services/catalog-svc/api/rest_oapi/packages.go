package rest_oapi

import (
	"errors"
	"time"

	"catalog-svc/service"
	"catalog-svc/util/apperrors"
	"catalog-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ListPackages implements ServerInterface.
// GET /v1/packages — public, filterable, cursor-paginated. Returns only
// active packages per § Catalog. Error envelopes match the contract
// snake_case codes (`invalid_query_param`, `invalid_cursor`,
// `internal_error`) and include `trace_id`.
func (s *Server) ListPackages(c *fiber.Ctx, params ListPackagesParams) error {
	const op = "rest_oapi.Server.ListPackages"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/packages"),
		attribute.String("method", "GET"),
	)
	logger.Info().Str("op", op).Msg("")

	// Query-param validation. oapi-codegen does not enforce enum /
	// range constraints on query params by default, so we re-check
	// here before dispatching to the service layer. Failures map to
	// `invalid_query_param` 400 per § Catalog.
	svcParams := &service.GetPackagesParams{}
	if params.Kind != nil {
		if !isKnownPackageKind(string(*params.Kind)) {
			return writeInvalidQueryParam(c, span, logger, "kind")
		}
		svcParams.Kind = string(*params.Kind)
	}
	if params.AirlineCode != nil {
		svcParams.AirlineCode = *params.AirlineCode
	}
	if params.HotelId != nil {
		svcParams.HotelID = *params.HotelId
	}
	if params.DepartureFrom != nil {
		svcParams.DepartureFrom = params.DepartureFrom.Format("2006-01-02")
	}
	if params.DepartureTo != nil {
		svcParams.DepartureTo = params.DepartureTo.Format("2006-01-02")
	}
	if params.Cursor != nil {
		svcParams.Cursor = *params.Cursor
	}
	if params.Limit != nil {
		if *params.Limit < 1 || *params.Limit > 100 {
			return writeInvalidQueryParam(c, span, logger, "limit")
		}
		svcParams.Limit = *params.Limit
	}

	result, err := s.svc.GetPackages(ctx, svcParams)
	if err != nil {
		logger.Warn().Err(err).Str("op", op).Msg("")
		return s.writeCatalogError(c, span, err)
	}

	resp := ListPackagesResponse{
		Packages: make([]PackageListItem, 0, len(result.Packages)),
		Page: PageMeta{
			HasMore: result.HasMore,
		},
	}
	if result.NextCursor != "" {
		nc := result.NextCursor
		resp.Page.NextCursor = &nc
	}
	for _, p := range result.Packages {
		item := PackageListItem{
			Id:            p.ID,
			Kind:          PackageKind(p.Kind),
			Name:          p.Name,
			Description:   p.Description,
			CoverPhotoUrl: p.CoverPhotoUrl,
			StartingPrice: Money{
				ListAmount:         p.StartingPrice.ListAmount,
				ListCurrency:       p.StartingPrice.ListCurrency,
				SettlementCurrency: MoneySettlementCurrency(p.StartingPrice.SettlementCurrency),
			},
		}
		if p.NextDeparture != nil {
			item.NextDeparture = &NextDeparture{
				Id:             p.NextDeparture.ID,
				DepartureDate:  parseISODate(p.NextDeparture.DepartureDate),
				ReturnDate:     parseISODate(p.NextDeparture.ReturnDate),
				RemainingSeats: p.NextDeparture.RemainingSeats,
			}
		}
		resp.Packages = append(resp.Packages, item)
	}

	span.SetAttributes(attribute.Int("output.packages.count", len(resp.Packages)))
	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetPackageById implements ServerInterface.
// GET /v1/packages/{id} — public, single-row. `404 package_not_found`
// for any id that does not match an active package.
func (s *Server) GetPackageById(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.GetPackageById"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/packages/:id"),
		attribute.String("method", "GET"),
		attribute.String("package_id", id),
	)
	logger.Info().Str("op", op).Str("package_id", id).Msg("")

	detail, err := s.svc.GetPackageByID(ctx, &service.GetPackageByIDParams{ID: id})
	if err != nil {
		logger.Warn().Err(err).Str("op", op).Str("package_id", id).Msg("")
		return s.writeCatalogError(c, span, err)
	}

	resp := GetPackageResponse{
		Package: PackageDetail{
			Id:            detail.ID,
			Kind:          PackageKind(detail.Kind),
			Name:          detail.Name,
			Description:   detail.Description,
			Highlights:    detail.Highlights,
			CoverPhotoUrl: detail.CoverPhotoUrl,
			Hotels:        make([]HotelRef, 0, len(detail.Hotels)),
			AddOns:        make([]AddonRef, 0, len(detail.Addons)),
			Departures:    make([]DepartureSummary, 0, len(detail.Departures)),
		},
	}
	if detail.Itinerary != nil {
		days := make([]ItineraryDay, 0, len(detail.Itinerary.Days))
		for _, d := range detail.Itinerary.Days {
			day := ItineraryDay{
				Day:         d.Day,
				Title:       d.Title,
				Description: d.Description,
			}
			if d.PhotoUrl != "" {
				photo := d.PhotoUrl
				day.PhotoUrl = &photo
			}
			days = append(days, day)
		}
		resp.Package.Itinerary = &Itinerary{
			Id:        detail.Itinerary.ID,
			Days:      days,
			PublicUrl: detail.Itinerary.PublicUrl,
		}
	}
	if detail.Airline != nil {
		resp.Package.Airline = &AirlineRef{
			Id:           detail.Airline.ID,
			Code:         detail.Airline.Code,
			Name:         detail.Airline.Name,
			OperatorKind: OperatorKind(detail.Airline.OperatorKind),
		}
	}
	if detail.Muthawwif != nil {
		resp.Package.Muthawwif = &MuthawwifRef{
			Id:          detail.Muthawwif.ID,
			Name:        detail.Muthawwif.Name,
			PortraitUrl: detail.Muthawwif.PortraitUrl,
		}
	}
	for _, h := range detail.Hotels {
		resp.Package.Hotels = append(resp.Package.Hotels, HotelRef{
			Id:               h.ID,
			Name:             h.Name,
			City:             h.City,
			StarRating:       h.StarRating,
			WalkingDistanceM: h.WalkingDistanceM,
		})
	}
	for _, a := range detail.Addons {
		resp.Package.AddOns = append(resp.Package.AddOns, AddonRef{
			Id:                 a.ID,
			Name:               a.Name,
			ListAmount:         a.ListAmount,
			ListCurrency:       a.ListCurrency,
			SettlementCurrency: AddonRefSettlementCurrency(a.SettlementCurrency),
		})
	}
	for _, d := range detail.Departures {
		resp.Package.Departures = append(resp.Package.Departures, DepartureSummary{
			Id:             d.ID,
			DepartureDate:  parseISODate(d.DepartureDate),
			ReturnDate:     parseISODate(d.ReturnDate),
			RemainingSeats: d.RemainingSeats,
			Status:         DepartureStatus(d.Status),
		})
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}

// ---------------------------------------------------------------------------
// Error shaping — § Catalog envelope (snake_case code + trace_id)
// ---------------------------------------------------------------------------

// writeCatalogError emits the contract-exact error envelope for catalog
// endpoints. Records the error on the span, maps the domain sentinel
// to a (status, snake_case code, id-ID message) triple, and includes
// the current span's trace_id for correlation with Tempo/Grafana.
// Unrecognised errors fall back to `internal_error` / 500.
func (s *Server) writeCatalogError(c *fiber.Ctx, span trace.Span, err error) error {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	status, code, message := mapCatalogError(err)

	resp := ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	if sc := span.SpanContext(); sc.HasTraceID() {
		tid := sc.TraceID().String()
		resp.Error.TraceId = &tid
	}
	return c.Status(status).JSON(resp)
}

func mapCatalogError(err error) (int, string, string) {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		return fiber.StatusNotFound, "package_not_found", "paket tidak ditemukan"
	case errors.Is(err, apperrors.ErrValidation):
		// Reached from the service layer for bad cursor payloads. Handler-
		// level enum / range validation uses writeInvalidQueryParam
		// directly so this branch is cursor-specific.
		return fiber.StatusBadRequest, "invalid_cursor", "cursor tidak valid"
	default:
		return fiber.StatusInternalServerError, "internal_error", "terjadi kesalahan tidak terduga"
	}
}

// writeInvalidQueryParam emits a 400 envelope for a named bad query
// param with the contracted snake_case `invalid_query_param` code.
func writeInvalidQueryParam(c *fiber.Ctx, span trace.Span, logger zerolog.Logger, paramName string) error {
	logger.Warn().Str("param", paramName).Msg("invalid query param")
	span.SetAttributes(attribute.String("invalid_query_param", paramName))
	span.SetStatus(codes.Error, "invalid_query_param")

	resp := ErrorResponse{}
	resp.Error.Code = "invalid_query_param"
	resp.Error.Message = "parameter kueri tidak valid: " + paramName
	if sc := span.SpanContext(); sc.HasTraceID() {
		tid := sc.TraceID().String()
		resp.Error.TraceId = &tid
	}
	return c.Status(fiber.StatusBadRequest).JSON(resp)
}

// isKnownPackageKind returns true when s matches one of the § Catalog
// enum values. The generated type ListPackagesParamsKind is a string
// alias, so unknown values reach the handler unchanged.
func isKnownPackageKind(s string) bool {
	switch ListPackagesParamsKind(s) {
	case ListPackagesParamsKindUmrahReguler,
		ListPackagesParamsKindUmrahPlus,
		ListPackagesParamsKindHajjFuroda,
		ListPackagesParamsKindHajjKhusus,
		ListPackagesParamsKindBadal,
		ListPackagesParamsKindFinancial,
		ListPackagesParamsKindRetail:
		return true
	}
	return false
}

// parseISODate converts a "YYYY-MM-DD" string into an oapi-codegen
// Date value. Empty/invalid input yields a zero Date — only fed with
// DB-produced values, so invalid is not expected at runtime.
func parseISODate(s string) openapi_types.Date {
	if s == "" {
		return openapi_types.Date{}
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return openapi_types.Date{}
	}
	return openapi_types.Date{Time: t}
}
