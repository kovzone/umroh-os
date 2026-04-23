package rest_oapi

// S1-E-07 / BL-CAT-014 — Staff catalog write REST handlers.
//
// Per ADR-0009, catalog-svc is transitioning to gRPC-only. These REST handlers
// are kept for compile completeness and potential direct-access during the pilot;
// from S1-E-10 onwards, all external traffic flows through gateway-svc which
// calls catalog-svc's gRPC surface. The permission gate (bearer + CheckPermission)
// is enforced at the gRPC layer in api/grpc_api/catalog_write.go.
//
// Validation: 400 with `validation_error` code for missing required fields.
// Auth: 401 if Authorization header absent/invalid; 403 if permission denied.
// Per § Catalog error envelope: snake_case code + id-ID message + trace_id.

import (
	"errors"

	"catalog-svc/service"
	"catalog-svc/util/apperrors"
	"catalog-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// writeWriteError emits the catalog error envelope for write endpoints.
// Uses the same shape as the public read errors (§ Catalog contract).
func (s *Server) writeWriteError(c *fiber.Ctx, err error, notFoundCode, notFoundMsg string) error {
	code := "internal_error"
	msg := "terjadi kesalahan tidak terduga"
	httpStatus := fiber.StatusInternalServerError

	switch {
	case errors.Is(err, apperrors.ErrUnauthorized):
		httpStatus = fiber.StatusUnauthorized
		code = "unauthorized"
		msg = "autentikasi diperlukan"
	case errors.Is(err, apperrors.ErrForbidden):
		httpStatus = fiber.StatusForbidden
		code = "forbidden"
		msg = "izin tidak mencukupi"
	case errors.Is(err, apperrors.ErrValidation):
		httpStatus = fiber.StatusBadRequest
		code = "validation_error"
		msg = "data permintaan tidak valid"
	case errors.Is(err, apperrors.ErrNotFound):
		httpStatus = fiber.StatusNotFound
		code = notFoundCode
		msg = notFoundMsg
	case errors.Is(err, apperrors.ErrConflict):
		httpStatus = fiber.StatusConflict
		code = "conflict"
		msg = "konflik data"
	}

	resp := ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = msg
	return c.Status(httpStatus).JSON(resp)
}

// CreatePackage handles POST /v1/packages (staff only).
// Bearer token + catalog.package.manage permission enforced via iam-svc.
func (s *Server) CreatePackage(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreatePackage"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "POST /v1/packages"),
	)
	logger.Info().Str("op", op).Msg("")

	var body CreatePackageRequestBody
	if err := c.BodyParser(&body); err != nil {
		return s.writeWriteError(c,
			errors.Join(apperrors.ErrValidation, err),
			"", "")
	}

	if string(body.Kind) == "" {
		return s.writeWriteError(c,
			errors.Join(apperrors.ErrValidation, fiber.NewError(fiber.StatusBadRequest, "kind required")),
			"", "")
	}
	if body.Name == "" {
		return s.writeWriteError(c,
			errors.Join(apperrors.ErrValidation, fiber.NewError(fiber.StatusBadRequest, "name required")),
			"", "")
	}

	// Extract user_id from context — set by upstream auth middleware at gateway-svc.
	// In the pilot (direct service call) we read X-User-Id header as a fallback.
	userID := c.Get("X-User-Id")
	if userID == "" {
		return s.writeWriteError(c,
			errors.Join(apperrors.ErrUnauthorized, fiber.NewError(fiber.StatusUnauthorized, "X-User-Id header required")),
			"", "")
	}

	params := &service.CreatePackageParams{
		UserID:        userID,
		Kind:          string(body.Kind),
		Name:          body.Name,
		Description:   body.Description,
		CoverPhotoUrl: body.CoverPhotoUrl,
		Highlights:    body.Highlights,
		Status:        body.Status,
		ItineraryID:   body.ItineraryID,
		AirlineID:     body.AirlineID,
		MuthawwifID:   body.MuthawwifID,
		HotelIDs:      body.HotelIDs,
		AddonIDs:      body.AddonIDs,
	}

	detail, err := s.svc.CreatePackage(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return s.writeWriteError(c, err, "package_not_found", "paket tidak ditemukan")
	}

	resp := StaffPackageResponse{
		Package: staffPackageDetailFromSvc(detail),
	}
	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// UpdatePackage handles PUT /v1/packages/{id} (staff only).
func (s *Server) UpdatePackage(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.UpdatePackage"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "PUT /v1/packages/:id"),
		attribute.String("package_id", id),
	)
	logger.Info().Str("op", op).Str("package_id", id).Msg("")

	var body UpdatePackageRequestBody
	if err := c.BodyParser(&body); err != nil {
		return s.writeWriteError(c, errors.Join(apperrors.ErrValidation, err), "", "")
	}

	userID := c.Get("X-User-Id")
	if userID == "" {
		return s.writeWriteError(c,
			errors.Join(apperrors.ErrUnauthorized, fiber.NewError(fiber.StatusUnauthorized, "X-User-Id header required")),
			"", "")
	}

	params := &service.UpdatePackageParams{
		UserID:        userID,
		ID:            id,
		Name:          body.Name,
		Description:   body.Description,
		CoverPhotoUrl: body.CoverPhotoUrl,
		Highlights:    body.Highlights,
		Status:        body.Status,
		ItineraryID:   body.ItineraryID,
		AirlineID:     body.AirlineID,
		MuthawwifID:   body.MuthawwifID,
		HotelIDs:      body.HotelIDs,
		AddonIDs:      body.AddonIDs,
	}

	detail, err := s.svc.UpdatePackage(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Str("op", op).Str("package_id", id).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return s.writeWriteError(c, err, "package_not_found", "paket tidak ditemukan")
	}

	resp := StaffPackageResponse{
		Package: staffPackageDetailFromSvc(detail),
	}
	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}

// CreateDeparture handles POST /v1/packages/{id}/departures (staff only).
func (s *Server) CreateDeparture(c *fiber.Ctx, packageId string) error {
	const op = "rest_oapi.Server.CreateDeparture"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "POST /v1/packages/:id/departures"),
		attribute.String("package_id", packageId),
	)
	logger.Info().Str("op", op).Str("package_id", packageId).Msg("")

	var body CreateDepartureRequestBody
	if err := c.BodyParser(&body); err != nil {
		return s.writeWriteError(c, errors.Join(apperrors.ErrValidation, err), "", "")
	}

	if body.DepartureDate == "" {
		return s.writeWriteError(c,
			errors.Join(apperrors.ErrValidation, fiber.NewError(fiber.StatusBadRequest, "departure_date required")),
			"", "")
	}
	if body.ReturnDate == "" {
		return s.writeWriteError(c,
			errors.Join(apperrors.ErrValidation, fiber.NewError(fiber.StatusBadRequest, "return_date required")),
			"", "")
	}
	if body.TotalSeats < 1 {
		return s.writeWriteError(c,
			errors.Join(apperrors.ErrValidation, fiber.NewError(fiber.StatusBadRequest, "total_seats must be >= 1")),
			"", "")
	}

	userID := c.Get("X-User-Id")
	if userID == "" {
		return s.writeWriteError(c,
			errors.Join(apperrors.ErrUnauthorized, fiber.NewError(fiber.StatusUnauthorized, "X-User-Id header required")),
			"", "")
	}

	pricing := make([]service.PricingInput, 0, len(body.Pricing))
	for _, p := range body.Pricing {
		pricing = append(pricing, service.PricingInput{
			RoomType:           string(p.RoomType),
			ListAmount:         p.ListAmount,
			ListCurrency:       p.ListCurrency,
			SettlementCurrency: p.SettlementCurrency,
		})
	}

	detail, err := s.svc.CreateDeparture(ctx, &service.CreateDepartureParams{
		UserID:        userID,
		PackageID:     packageId,
		DepartureDate: body.DepartureDate,
		ReturnDate:    body.ReturnDate,
		TotalSeats:    body.TotalSeats,
		Status:        body.Status,
		Pricing:       pricing,
	})
	if err != nil {
		logger.Warn().Err(err).Str("op", op).Str("package_id", packageId).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return s.writeWriteError(c, err, "package_not_found", "paket tidak ditemukan")
	}

	resp := StaffDepartureResponse{
		Departure: staffDepartureDetailFromSvc(detail),
	}
	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// UpdateDeparture handles PUT /v1/package-departures/{id} (staff only).
func (s *Server) UpdateDeparture(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.UpdateDeparture"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "PUT /v1/package-departures/:id"),
		attribute.String("departure_id", id),
	)
	logger.Info().Str("op", op).Str("departure_id", id).Msg("")

	var body UpdateDepartureRequestBody
	if err := c.BodyParser(&body); err != nil {
		return s.writeWriteError(c, errors.Join(apperrors.ErrValidation, err), "", "")
	}

	userID := c.Get("X-User-Id")
	if userID == "" {
		return s.writeWriteError(c,
			errors.Join(apperrors.ErrUnauthorized, fiber.NewError(fiber.StatusUnauthorized, "X-User-Id header required")),
			"", "")
	}

	var pricing []service.PricingInput
	if body.Pricing != nil {
		pricing = make([]service.PricingInput, 0, len(body.Pricing))
		for _, p := range body.Pricing {
			pricing = append(pricing, service.PricingInput{
				RoomType:           string(p.RoomType),
				ListAmount:         p.ListAmount,
				ListCurrency:       p.ListCurrency,
				SettlementCurrency: p.SettlementCurrency,
			})
		}
	}

	detail, err := s.svc.UpdateDeparture(ctx, &service.UpdateDepartureParams{
		UserID:        userID,
		ID:            id,
		DepartureDate: body.DepartureDate,
		ReturnDate:    body.ReturnDate,
		TotalSeats:    body.TotalSeats,
		Status:        body.Status,
		Pricing:       pricing,
	})
	if err != nil {
		logger.Warn().Err(err).Str("op", op).Str("departure_id", id).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return s.writeWriteError(c, err, "departure_not_found", "keberangkatan tidak ditemukan")
	}

	resp := StaffDepartureResponse{
		Departure: staffDepartureDetailFromSvc(detail),
	}
	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}

// ---------------------------------------------------------------------------
// Response shape helpers
// ---------------------------------------------------------------------------

// staffPackageDetailFromSvc converts service.PackageDetail to StaffPackageDetail
// (which surfaces the status field hidden from the public read model).
func staffPackageDetailFromSvc(d *service.PackageDetail) StaffPackageDetail {
	out := StaffPackageDetail{
		Id:            d.ID,
		Kind:          PackageKind(d.Kind),
		Name:          d.Name,
		Description:   d.Description,
		Highlights:    d.Highlights,
		CoverPhotoUrl: d.CoverPhotoUrl,
		Status:        d.Status,
		Hotels:        make([]HotelRef, 0, len(d.Hotels)),
		AddOns:        make([]AddonRef, 0, len(d.Addons)),
		Departures:    make([]DepartureSummary, 0, len(d.Departures)),
	}
	if d.Itinerary != nil {
		days := make([]ItineraryDay, 0, len(d.Itinerary.Days))
		for _, day := range d.Itinerary.Days {
			itDay := ItineraryDay{Day: day.Day, Title: day.Title, Description: day.Description}
			if day.PhotoUrl != "" {
				u := day.PhotoUrl
				itDay.PhotoUrl = &u
			}
			days = append(days, itDay)
		}
		out.Itinerary = &Itinerary{Id: d.Itinerary.ID, Days: days, PublicUrl: d.Itinerary.PublicUrl}
	}
	if d.Airline != nil {
		out.Airline = &AirlineRef{
			Id: d.Airline.ID, Code: d.Airline.Code,
			Name: d.Airline.Name, OperatorKind: OperatorKind(d.Airline.OperatorKind),
		}
	}
	if d.Muthawwif != nil {
		out.Muthawwif = &MuthawwifRef{Id: d.Muthawwif.ID, Name: d.Muthawwif.Name, PortraitUrl: d.Muthawwif.PortraitUrl}
	}
	for _, h := range d.Hotels {
		out.Hotels = append(out.Hotels, HotelRef{
			Id: h.ID, Name: h.Name, City: h.City,
			StarRating: h.StarRating, WalkingDistanceM: h.WalkingDistanceM,
		})
	}
	for _, a := range d.Addons {
		out.AddOns = append(out.AddOns, AddonRef{
			Id: a.ID, Name: a.Name, ListAmount: a.ListAmount,
			ListCurrency: a.ListCurrency, SettlementCurrency: AddonRefSettlementCurrency(a.SettlementCurrency),
		})
	}
	for _, dep := range d.Departures {
		out.Departures = append(out.Departures, DepartureSummary{
			Id:             dep.ID,
			DepartureDate:  parseISODate(dep.DepartureDate),
			ReturnDate:     parseISODate(dep.ReturnDate),
			RemainingSeats: dep.RemainingSeats,
			Status:         DepartureStatus(dep.Status),
		})
	}
	return out
}

// staffDepartureDetailFromSvc converts service.DepartureDetail to StaffDepartureDetail.
func staffDepartureDetailFromSvc(d *service.DepartureDetail) StaffDepartureDetail {
	pricing := make([]PackagePricing, 0, len(d.Pricing))
	for _, p := range d.Pricing {
		pricing = append(pricing, PackagePricing{
			RoomType:           RoomType(p.RoomType),
			ListAmount:         p.ListAmount,
			ListCurrency:       p.ListCurrency,
			SettlementCurrency: PackagePricingSettlementCurrency(p.SettlementCurrency),
		})
	}
	return StaffDepartureDetail{
		Id:             d.ID,
		PackageId:      d.PackageID,
		DepartureDate:  d.DepartureDate,
		ReturnDate:     d.ReturnDate,
		TotalSeats:     d.TotalSeats,
		RemainingSeats: d.RemainingSeats,
		Status:         d.Status,
		Pricing:        pricing,
		VendorReadiness: VendorReadiness{
			Ticket: VendorReadinessState(d.VendorReadiness.Ticket),
			Hotel:  VendorReadinessState(d.VendorReadiness.Hotel),
			Visa:   VendorReadinessState(d.VendorReadiness.Visa),
		},
	}
}
