// proxy_catalog_masters.go — gateway REST handlers for catalog master data
// CRUD and departure pricing (Wave 1A / Phase 6).
//
// Route topology (all bearer-protected):
//   GET    /v1/catalog/hotels           → ListHotels
//   POST   /v1/catalog/hotels           → CreateHotel
//   PUT    /v1/catalog/hotels/:id       → UpdateHotel
//   DELETE /v1/catalog/hotels/:id       → DeleteHotel
//   GET    /v1/catalog/airlines         → ListAirlines
//   POST   /v1/catalog/airlines         → CreateAirline
//   PUT    /v1/catalog/airlines/:id     → UpdateAirline
//   DELETE /v1/catalog/airlines/:id     → DeleteAirline
//   GET    /v1/catalog/muthawwif        → ListMuthawwif
//   POST   /v1/catalog/muthawwif        → CreateMuthawwif
//   PUT    /v1/catalog/muthawwif/:id    → UpdateMuthawwif
//   DELETE /v1/catalog/muthawwif/:id    → DeleteMuthawwif
//   GET    /v1/catalog/addons           → ListAddons
//   POST   /v1/catalog/addons           → CreateAddon
//   PUT    /v1/catalog/addons/:id       → UpdateAddon
//   DELETE /v1/catalog/addons/:id       → DeleteAddon
//   GET    /v1/departures/:id/pricing   → GetDeparturePricing
//   PUT    /v1/departures/:id/pricing   → SetDeparturePricing
//
// Per ADR-0009: gateway is the single REST entry-point; catalog-svc is pure gRPC.
package rest_oapi

import (
	"errors"

	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ---------------------------------------------------------------------------
// Request / response body types
// ---------------------------------------------------------------------------

type CreateHotelBody struct {
	Name             string `json:"name"`
	City             string `json:"city"`
	StarRating       int32  `json:"star_rating"`
	WalkingDistanceM int32  `json:"walking_distance_m"`
}

type UpdateHotelBody struct {
	Name             string `json:"name,omitempty"`
	City             string `json:"city,omitempty"`
	StarRating       int32  `json:"star_rating,omitempty"`
	WalkingDistanceM int32  `json:"walking_distance_m,omitempty"`
}

type HotelResponseData struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	City             string `json:"city"`
	StarRating       int32  `json:"star_rating"`
	WalkingDistanceM int32  `json:"walking_distance_m"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type SingleHotelResponse struct {
	Data HotelResponseData `json:"data"`
}

type HotelListResponse struct {
	Data    []HotelResponseData `json:"data"`
	HasMore bool                `json:"has_more"`
	Cursor  string              `json:"cursor,omitempty"`
}

func mapHotelResult(h *catalog_grpc_adapter.HotelResult) HotelResponseData {
	return HotelResponseData{
		ID:               h.ID,
		Name:             h.Name,
		City:             h.City,
		StarRating:       h.StarRating,
		WalkingDistanceM: h.WalkingDistanceM,
		CreatedAt:        h.CreatedAt,
		UpdatedAt:        h.UpdatedAt,
	}
}

type CreateAirlineBody struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	OperatorKind string `json:"operator_kind"`
}

type UpdateAirlineBody struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
}

type AirlineResponseData struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	OperatorKind string `json:"operator_kind"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type SingleAirlineResponse struct {
	Data AirlineResponseData `json:"data"`
}

type AirlineListResponse struct {
	Data    []AirlineResponseData `json:"data"`
	HasMore bool                  `json:"has_more"`
	Cursor  string                `json:"cursor,omitempty"`
}

func mapAirlineResult(a *catalog_grpc_adapter.AirlineResult) AirlineResponseData {
	return AirlineResponseData{
		ID:           a.ID,
		Code:         a.Code,
		Name:         a.Name,
		OperatorKind: a.OperatorKind,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

type CreateMuthawwifBody struct {
	Name        string `json:"name"`
	PortraitURL string `json:"portrait_url,omitempty"`
}

type UpdateMuthawwifBody struct {
	Name        string `json:"name,omitempty"`
	PortraitURL string `json:"portrait_url,omitempty"`
}

type MuthawwifResponseData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PortraitURL string `json:"portrait_url,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type SingleMuthawwifResponse struct {
	Data MuthawwifResponseData `json:"data"`
}

type MuthawwifListResponse struct {
	Data    []MuthawwifResponseData `json:"data"`
	HasMore bool                    `json:"has_more"`
	Cursor  string                  `json:"cursor,omitempty"`
}

func mapMuthawwifResult(m *catalog_grpc_adapter.MuthawwifResult) MuthawwifResponseData {
	return MuthawwifResponseData{
		ID:          m.ID,
		Name:        m.Name,
		PortraitURL: m.PortraitURL,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

type CreateAddonBody struct {
	Name          string `json:"name"`
	ListAmountIdr int64  `json:"list_amount_idr"`
}

type UpdateAddonBody struct {
	Name          string `json:"name,omitempty"`
	ListAmountIdr int64  `json:"list_amount_idr,omitempty"`
}

type AddonResponseData struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ListAmountIdr int64  `json:"list_amount_idr"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type SingleAddonResponse struct {
	Data AddonResponseData `json:"data"`
}

type AddonListResponse struct {
	Data    []AddonResponseData `json:"data"`
	HasMore bool                `json:"has_more"`
	Cursor  string              `json:"cursor,omitempty"`
}

func mapAddonResult(a *catalog_grpc_adapter.AddonResult) AddonResponseData {
	return AddonResponseData{
		ID:            a.ID,
		Name:          a.Name,
		ListAmountIdr: a.ListAmountIdr,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

type PricingEntryResponseData struct {
	ID            string `json:"id"`
	DepartureID   string `json:"departure_id"`
	RoomType      string `json:"room_type"`
	ListAmountIdr int64  `json:"list_amount_idr"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type PricingListResponse struct {
	Data []PricingEntryResponseData `json:"data"`
}

type SetPricingBody struct {
	Pricings []struct {
		RoomType      string `json:"room_type"`
		ListAmountIdr int64  `json:"list_amount_idr"`
	} `json:"pricings"`
}

func mapPricingResult(p *catalog_grpc_adapter.DeparturePricingResult) PricingEntryResponseData {
	return PricingEntryResponseData{
		ID:            p.ID,
		DepartureID:   p.DepartureID,
		RoomType:      p.RoomType,
		ListAmountIdr: p.ListAmountIdr,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

// ---------------------------------------------------------------------------
// Hotel handlers
// ---------------------------------------------------------------------------

func (s *Server) ListHotels(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListHotels"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("endpoint", "GET /v1/catalog/hotels"))
	logger.Info().Str("op", op).Msg("")

	cursor := c.Query("cursor")
	limit := int32(c.QueryInt("limit", 50))

	result, err := s.svc.ListHotels(ctx, &catalog_grpc_adapter.ListMastersParams{
		Cursor: cursor,
		Limit:  limit,
	})
	if err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	data := make([]HotelResponseData, 0, len(result.Hotels))
	for _, h := range result.Hotels {
		data = append(data, mapHotelResult(h))
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(HotelListResponse{Data: data, HasMore: result.HasMore, Cursor: result.Cursor})
}

func (s *Server) CreateHotel(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateHotel"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	var body CreateHotelBody
	if err := c.BodyParser(&body); err != nil {
		return writeCatalogMastersError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.CreateHotel(ctx, &catalog_grpc_adapter.CreateHotelParams{
		Name:             body.Name,
		City:             body.City,
		StarRating:       body.StarRating,
		WalkingDistanceM: body.WalkingDistanceM,
	})
	if err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(SingleHotelResponse{Data: mapHotelResult(result)})
}

func (s *Server) UpdateHotel(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.UpdateHotel"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("hotel_id", id))

	var body UpdateHotelBody
	if err := c.BodyParser(&body); err != nil {
		return writeCatalogMastersError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.UpdateHotel(ctx, &catalog_grpc_adapter.UpdateHotelParams{
		ID:               id,
		Name:             body.Name,
		City:             body.City,
		StarRating:       body.StarRating,
		WalkingDistanceM: body.WalkingDistanceM,
	})
	if err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(SingleHotelResponse{Data: mapHotelResult(result)})
}

func (s *Server) DeleteHotel(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.DeleteHotel"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("hotel_id", id))

	if err := s.svc.DeleteHotel(ctx, id); err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	span.SetStatus(codes.Ok, "deleted")
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// ---------------------------------------------------------------------------
// Airline handlers
// ---------------------------------------------------------------------------

func (s *Server) ListAirlines(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListAirlines"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	cursor := c.Query("cursor")
	limit := int32(c.QueryInt("limit", 50))

	result, err := s.svc.ListAirlines(ctx, &catalog_grpc_adapter.ListMastersParams{
		Cursor: cursor,
		Limit:  limit,
	})
	if err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	data := make([]AirlineResponseData, 0, len(result.Airlines))
	for _, a := range result.Airlines {
		data = append(data, mapAirlineResult(a))
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(AirlineListResponse{Data: data, HasMore: result.HasMore, Cursor: result.Cursor})
}

func (s *Server) CreateAirline(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateAirline"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	var body CreateAirlineBody
	if err := c.BodyParser(&body); err != nil {
		return writeCatalogMastersError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.CreateAirline(ctx, &catalog_grpc_adapter.CreateAirlineParams{
		Code:         body.Code,
		Name:         body.Name,
		OperatorKind: body.OperatorKind,
	})
	if err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(SingleAirlineResponse{Data: mapAirlineResult(result)})
}

func (s *Server) UpdateAirline(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.UpdateAirline"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	var body UpdateAirlineBody
	if err := c.BodyParser(&body); err != nil {
		return writeCatalogMastersError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.UpdateAirline(ctx, &catalog_grpc_adapter.UpdateAirlineParams{
		ID:   id,
		Code: body.Code,
		Name: body.Name,
	})
	if err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(SingleAirlineResponse{Data: mapAirlineResult(result)})
}

func (s *Server) DeleteAirline(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.DeleteAirline"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	if err := s.svc.DeleteAirline(ctx, id); err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	span.SetStatus(codes.Ok, "deleted")
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// ---------------------------------------------------------------------------
// Muthawwif handlers
// ---------------------------------------------------------------------------

func (s *Server) ListMuthawwif(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListMuthawwif"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	cursor := c.Query("cursor")
	limit := int32(c.QueryInt("limit", 50))

	result, err := s.svc.ListMuthawwif(ctx, &catalog_grpc_adapter.ListMastersParams{
		Cursor: cursor,
		Limit:  limit,
	})
	if err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	data := make([]MuthawwifResponseData, 0, len(result.Muthawwif))
	for _, m := range result.Muthawwif {
		data = append(data, mapMuthawwifResult(m))
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(MuthawwifListResponse{Data: data, HasMore: result.HasMore, Cursor: result.Cursor})
}

func (s *Server) CreateMuthawwif(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateMuthawwif"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	var body CreateMuthawwifBody
	if err := c.BodyParser(&body); err != nil {
		return writeCatalogMastersError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.CreateMuthawwif(ctx, &catalog_grpc_adapter.CreateMuthawwifParams{
		Name:        body.Name,
		PortraitURL: body.PortraitURL,
	})
	if err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(SingleMuthawwifResponse{Data: mapMuthawwifResult(result)})
}

func (s *Server) UpdateMuthawwif(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.UpdateMuthawwif"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	var body UpdateMuthawwifBody
	if err := c.BodyParser(&body); err != nil {
		return writeCatalogMastersError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.UpdateMuthawwif(ctx, &catalog_grpc_adapter.UpdateMuthawwifParams{
		ID:          id,
		Name:        body.Name,
		PortraitURL: body.PortraitURL,
	})
	if err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(SingleMuthawwifResponse{Data: mapMuthawwifResult(result)})
}

func (s *Server) DeleteMuthawwif(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.DeleteMuthawwif"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	if err := s.svc.DeleteMuthawwif(ctx, id); err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	span.SetStatus(codes.Ok, "deleted")
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// ---------------------------------------------------------------------------
// Addon handlers
// ---------------------------------------------------------------------------

func (s *Server) ListAddons(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListAddons"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	cursor := c.Query("cursor")
	limit := int32(c.QueryInt("limit", 50))

	result, err := s.svc.ListAddons(ctx, &catalog_grpc_adapter.ListMastersParams{
		Cursor: cursor,
		Limit:  limit,
	})
	if err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	data := make([]AddonResponseData, 0, len(result.Addons))
	for _, a := range result.Addons {
		data = append(data, mapAddonResult(a))
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(AddonListResponse{Data: data, HasMore: result.HasMore, Cursor: result.Cursor})
}

func (s *Server) CreateAddon(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateAddon"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	var body CreateAddonBody
	if err := c.BodyParser(&body); err != nil {
		return writeCatalogMastersError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.CreateAddon(ctx, &catalog_grpc_adapter.CreateAddonParams{
		Name:          body.Name,
		ListAmountIdr: body.ListAmountIdr,
	})
	if err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(SingleAddonResponse{Data: mapAddonResult(result)})
}

func (s *Server) UpdateAddon(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.UpdateAddon"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	var body UpdateAddonBody
	if err := c.BodyParser(&body); err != nil {
		return writeCatalogMastersError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.UpdateAddon(ctx, &catalog_grpc_adapter.UpdateAddonParams{
		ID:            id,
		Name:          body.Name,
		ListAmountIdr: body.ListAmountIdr,
	})
	if err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(SingleAddonResponse{Data: mapAddonResult(result)})
}

func (s *Server) DeleteAddon(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.DeleteAddon"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	if err := s.svc.DeleteAddon(ctx, id); err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	span.SetStatus(codes.Ok, "deleted")
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// ---------------------------------------------------------------------------
// Departure Pricing handlers
// ---------------------------------------------------------------------------

func (s *Server) GetDeparturePricing(c *fiber.Ctx, departureID string) error {
	const op = "rest_oapi.Server.GetDeparturePricing"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("departure_id", departureID))

	result, err := s.svc.GetDeparturePricing(ctx, departureID)
	if err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	data := make([]PricingEntryResponseData, 0, len(result.Pricings))
	for _, p := range result.Pricings {
		data = append(data, mapPricingResult(p))
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(PricingListResponse{Data: data})
}

func (s *Server) SetDeparturePricing(c *fiber.Ctx, departureID string) error {
	const op = "rest_oapi.Server.SetDeparturePricing"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("departure_id", departureID))

	var body SetPricingBody
	if err := c.BodyParser(&body); err != nil {
		return writeCatalogMastersError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	pricings := make([]*catalog_grpc_adapter.PricingUpsertInput, 0, len(body.Pricings))
	for _, p := range body.Pricings {
		pricings = append(pricings, &catalog_grpc_adapter.PricingUpsertInput{
			RoomType:      p.RoomType,
			ListAmountIdr: p.ListAmountIdr,
		})
	}

	result, err := s.svc.SetDeparturePricing(ctx, departureID, pricings)
	if err != nil {
		return writeCatalogMastersError(c, span, err)
	}

	data := make([]PricingEntryResponseData, 0, len(result.Pricings))
	for _, p := range result.Pricings {
		data = append(data, mapPricingResult(p))
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(PricingListResponse{Data: data})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeCatalogMastersError(c *fiber.Ctx, span trace.Span, err error) error {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	var httpStatus int
	var code, message string

	switch {
	case errors.Is(err, apperrors.ErrValidation):
		httpStatus = fiber.StatusBadRequest
		code = "validation_error"
		message = err.Error()
	case errors.Is(err, apperrors.ErrNotFound):
		httpStatus = fiber.StatusNotFound
		code = "not_found"
		message = err.Error()
	case errors.Is(err, apperrors.ErrConflict):
		httpStatus = fiber.StatusConflict
		code = "conflict"
		message = err.Error()
	case errors.Is(err, apperrors.ErrUnauthorized):
		httpStatus = fiber.StatusUnauthorized
		code = "unauthorized"
		message = "autentikasi diperlukan"
	case errors.Is(err, apperrors.ErrForbidden):
		httpStatus = fiber.StatusForbidden
		code = "forbidden"
		message = "akses ditolak"
	case errors.Is(err, apperrors.ErrServiceUnavailable):
		httpStatus = fiber.StatusBadGateway
		code = "service_unavailable"
		message = "layanan catalog sementara tidak tersedia"
	default:
		httpStatus = fiber.StatusInternalServerError
		code = "internal_error"
		message = "terjadi kesalahan tidak terduga"
	}

	return c.Status(httpStatus).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": message,
		},
	})
}
