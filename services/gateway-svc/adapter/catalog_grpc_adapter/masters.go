// masters.go — gateway adapter methods for catalog-svc master data RPCs
// (Wave 1A / Phase 6).
//
// Covers: Hotel, Airline, Muthawwif, Addon CRUD and departure pricing.
// Each method translates gateway-local params → pb request, forwards via gRPC,
// and translates pb response → adapter-local types. Proto types do not leak
// past this package.
//
// Permission check: gateway handlers MUST call iam_grpc_adapter.CheckPermission
// with resource="catalog.masters", action="manage", scope="global" BEFORE
// calling these methods.
package catalog_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/catalog_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Adapter-local result types
// ---------------------------------------------------------------------------

// HotelResult is the gateway-local representation of a catalog hotel.
type HotelResult struct {
	ID               string
	Name             string
	City             string
	StarRating       int32
	WalkingDistanceM int32
	CreatedAt        string
	UpdatedAt        string
}

// ListHotelsResult is the paginated list result for hotels.
type ListHotelsResult struct {
	Hotels  []*HotelResult
	HasMore bool
	Cursor  string
}

// AirlineResult is the gateway-local representation of a catalog airline.
type AirlineResult struct {
	ID           string
	Code         string
	Name         string
	OperatorKind string
	CreatedAt    string
	UpdatedAt    string
}

// ListAirlinesResult is the paginated list result for airlines.
type ListAirlinesResult struct {
	Airlines []*AirlineResult
	HasMore  bool
	Cursor   string
}

// MuthawwifResult is the gateway-local representation of a muthawwif record.
type MuthawwifResult struct {
	ID          string
	Name        string
	PortraitURL string
	CreatedAt   string
	UpdatedAt   string
}

// ListMuthawwifResult is the paginated list result for muthawwif.
type ListMuthawwifResult struct {
	Muthawwif []*MuthawwifResult
	HasMore   bool
	Cursor    string
}

// AddonResult is the gateway-local representation of a catalog addon.
type AddonResult struct {
	ID            string
	Name          string
	ListAmountIdr int64
	CreatedAt     string
	UpdatedAt     string
}

// ListAddonsResult is the paginated list result for addons.
type ListAddonsResult struct {
	Addons  []*AddonResult
	HasMore bool
	Cursor  string
}

// DeparturePricingResult is the gateway-local representation of a pricing row.
type DeparturePricingResult struct {
	ID            string
	DepartureID   string
	RoomType      string
	ListAmountIdr int64
	CreatedAt     string
	UpdatedAt     string
}

// PricingResult holds the list of departure pricing rows.
type PricingResult struct {
	Pricings []*DeparturePricingResult
}

// ---------------------------------------------------------------------------
// Params types
// ---------------------------------------------------------------------------

// CreateHotelParams is the input for CreateHotel.
type CreateHotelParams struct {
	UserID           string
	Name             string
	City             string
	StarRating       int32
	WalkingDistanceM int32
}

// UpdateHotelParams is the input for UpdateHotel.
type UpdateHotelParams struct {
	UserID           string
	ID               string
	Name             string
	City             string
	StarRating       int32
	WalkingDistanceM int32
}

// ListMastersParams is generic pagination params for list endpoints.
type ListMastersParams struct {
	UserID string
	Cursor string
	Limit  int32
}

// CreateAirlineParams is the input for CreateAirline.
type CreateAirlineParams struct {
	UserID       string
	Code         string
	Name         string
	OperatorKind string
}

// UpdateAirlineParams is the input for UpdateAirline.
type UpdateAirlineParams struct {
	UserID string
	ID     string
	Code   string
	Name   string
}

// CreateMuthawwifParams is the input for CreateMuthawwif.
type CreateMuthawwifParams struct {
	UserID      string
	Name        string
	PortraitURL string
}

// UpdateMuthawwifParams is the input for UpdateMuthawwif.
type UpdateMuthawwifParams struct {
	UserID      string
	ID          string
	Name        string
	PortraitURL string
}

// CreateAddonParams is the input for CreateAddon.
type CreateAddonParams struct {
	UserID        string
	Name          string
	ListAmountIdr int64
}

// UpdateAddonParams is the input for UpdateAddon.
type UpdateAddonParams struct {
	UserID        string
	ID            string
	Name          string
	ListAmountIdr int64
}

// PricingUpsertInput is a single room-type price entry for SetDeparturePricing.
type PricingUpsertInput struct {
	RoomType      string
	ListAmountIdr int64
}

// ---------------------------------------------------------------------------
// Helper: mastersClient lazily created from catalogClient conn
// ---------------------------------------------------------------------------

// mastersClient returns the CatalogMastersClient held on the Adapter.
// The Adapter stores it in the mastersClient field set by NewAdapter.
func (a *Adapter) mastersClient() pb.CatalogMastersClient {
	return a.catalogMastersClient
}

// ---------------------------------------------------------------------------
// Helper mappers
// ---------------------------------------------------------------------------

func fromProtoHotel(h *pb.MasterHotel) *HotelResult {
	if h == nil {
		return nil
	}
	return &HotelResult{
		ID:               h.GetId(),
		Name:             h.GetName(),
		City:             h.GetCity(),
		StarRating:       h.GetStarRating(),
		WalkingDistanceM: h.GetWalkingDistanceM(),
		CreatedAt:        h.GetCreatedAt(),
		UpdatedAt:        h.GetUpdatedAt(),
	}
}

func fromProtoAirline(a *pb.MasterAirline) *AirlineResult {
	if a == nil {
		return nil
	}
	return &AirlineResult{
		ID:           a.GetId(),
		Code:         a.GetCode(),
		Name:         a.GetName(),
		OperatorKind: a.GetOperatorKind(),
		CreatedAt:    a.GetCreatedAt(),
		UpdatedAt:    a.GetUpdatedAt(),
	}
}

func fromProtoMuthawwif(m *pb.MasterMuthawwif) *MuthawwifResult {
	if m == nil {
		return nil
	}
	return &MuthawwifResult{
		ID:          m.GetId(),
		Name:        m.GetName(),
		PortraitURL: m.GetPortraitUrl(),
		CreatedAt:   m.GetCreatedAt(),
		UpdatedAt:   m.GetUpdatedAt(),
	}
}

func fromProtoAddon(a *pb.MasterAddon) *AddonResult {
	if a == nil {
		return nil
	}
	return &AddonResult{
		ID:            a.GetId(),
		Name:          a.GetName(),
		ListAmountIdr: a.GetListAmountIdr(),
		CreatedAt:     a.GetCreatedAt(),
		UpdatedAt:     a.GetUpdatedAt(),
	}
}

func fromProtoDeparturePricing(p *pb.MasterDeparturePricing) *DeparturePricingResult {
	if p == nil {
		return nil
	}
	return &DeparturePricingResult{
		ID:            p.GetId(),
		DepartureID:   p.GetDepartureId(),
		RoomType:      p.GetRoomType(),
		ListAmountIdr: p.GetListAmountIdr(),
		CreatedAt:     p.GetCreatedAt(),
		UpdatedAt:     p.GetUpdatedAt(),
	}
}

// ---------------------------------------------------------------------------
// Hotel methods
// ---------------------------------------------------------------------------

// CreateHotel creates a new hotel master record.
func (a *Adapter) CreateHotel(ctx context.Context, params *CreateHotelParams) (*HotelResult, error) {
	const op = "catalog_grpc_adapter.Adapter.CreateHotel"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "CreateHotel"))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.mastersClient().CreateHotel(ctx, &pb.CreateHotelRequest{
		UserId:           params.UserID,
		Name:             params.Name,
		City:             params.City,
		StarRating:       params.StarRating,
		WalkingDistanceM: params.WalkingDistanceM,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return fromProtoHotel(resp.GetHotel()), nil
}

// UpdateHotel updates an existing hotel master record.
func (a *Adapter) UpdateHotel(ctx context.Context, params *UpdateHotelParams) (*HotelResult, error) {
	const op = "catalog_grpc_adapter.Adapter.UpdateHotel"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("id", params.ID))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.mastersClient().UpdateHotel(ctx, &pb.UpdateHotelRequest{
		UserId:           params.UserID,
		Id:               params.ID,
		Name:             params.Name,
		City:             params.City,
		StarRating:       params.StarRating,
		WalkingDistanceM: params.WalkingDistanceM,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return fromProtoHotel(resp.GetHotel()), nil
}

// DeleteHotel soft-deletes a hotel master record.
func (a *Adapter) DeleteHotel(ctx context.Context, userID, id string) error {
	const op = "catalog_grpc_adapter.Adapter.DeleteHotel"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("id", id))

	logger := logging.LogWithTrace(ctx, a.logger)

	_, err := a.mastersClient().DeleteHotel(ctx, &pb.DeleteHotelRequest{
		UserId: userID,
		Id:     id,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return nil
}

// ListHotels returns a paginated list of hotel master records.
func (a *Adapter) ListHotels(ctx context.Context, params *ListMastersParams) (*ListHotelsResult, error) {
	const op = "catalog_grpc_adapter.Adapter.ListHotels"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.mastersClient().ListHotels(ctx, &pb.ListHotelsRequest{
		UserId: params.UserID,
		Cursor: params.Cursor,
		Limit:  params.Limit,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	hotels := make([]*HotelResult, 0, len(resp.GetHotels()))
	for _, h := range resp.GetHotels() {
		hotels = append(hotels, fromProtoHotel(h))
	}

	span.SetStatus(codes.Ok, "ok")
	return &ListHotelsResult{
		Hotels:  hotels,
		HasMore: resp.GetHasMore(),
		Cursor:  resp.GetCursor(),
	}, nil
}

// ---------------------------------------------------------------------------
// Airline methods
// ---------------------------------------------------------------------------

// CreateAirline creates a new airline master record.
func (a *Adapter) CreateAirline(ctx context.Context, params *CreateAirlineParams) (*AirlineResult, error) {
	const op = "catalog_grpc_adapter.Adapter.CreateAirline"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.mastersClient().CreateAirline(ctx, &pb.CreateAirlineRequest{
		UserId:       params.UserID,
		Code:         params.Code,
		Name:         params.Name,
		OperatorKind: params.OperatorKind,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return fromProtoAirline(resp.GetAirline()), nil
}

// UpdateAirline updates an existing airline master record.
func (a *Adapter) UpdateAirline(ctx context.Context, params *UpdateAirlineParams) (*AirlineResult, error) {
	const op = "catalog_grpc_adapter.Adapter.UpdateAirline"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("id", params.ID))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.mastersClient().UpdateAirline(ctx, &pb.UpdateAirlineRequest{
		UserId: params.UserID,
		Id:     params.ID,
		Code:   params.Code,
		Name:   params.Name,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return fromProtoAirline(resp.GetAirline()), nil
}

// DeleteAirline soft-deletes an airline master record.
func (a *Adapter) DeleteAirline(ctx context.Context, userID, id string) error {
	const op = "catalog_grpc_adapter.Adapter.DeleteAirline"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("id", id))

	logger := logging.LogWithTrace(ctx, a.logger)

	_, err := a.mastersClient().DeleteAirline(ctx, &pb.DeleteAirlineRequest{
		UserId: userID,
		Id:     id,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return nil
}

// ListAirlines returns a paginated list of airline master records.
func (a *Adapter) ListAirlines(ctx context.Context, params *ListMastersParams) (*ListAirlinesResult, error) {
	const op = "catalog_grpc_adapter.Adapter.ListAirlines"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.mastersClient().ListAirlines(ctx, &pb.ListAirlinesRequest{
		UserId: params.UserID,
		Cursor: params.Cursor,
		Limit:  params.Limit,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	airlines := make([]*AirlineResult, 0, len(resp.GetAirlines()))
	for _, al := range resp.GetAirlines() {
		airlines = append(airlines, fromProtoAirline(al))
	}

	span.SetStatus(codes.Ok, "ok")
	return &ListAirlinesResult{
		Airlines: airlines,
		HasMore:  resp.GetHasMore(),
		Cursor:   resp.GetCursor(),
	}, nil
}

// ---------------------------------------------------------------------------
// Muthawwif methods
// ---------------------------------------------------------------------------

// CreateMuthawwif creates a new muthawwif master record.
func (a *Adapter) CreateMuthawwif(ctx context.Context, params *CreateMuthawwifParams) (*MuthawwifResult, error) {
	const op = "catalog_grpc_adapter.Adapter.CreateMuthawwif"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.mastersClient().CreateMuthawwif(ctx, &pb.CreateMuthawwifRequest{
		UserId:      params.UserID,
		Name:        params.Name,
		PortraitUrl: params.PortraitURL,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return fromProtoMuthawwif(resp.GetMuthawwif()), nil
}

// UpdateMuthawwif updates an existing muthawwif master record.
func (a *Adapter) UpdateMuthawwif(ctx context.Context, params *UpdateMuthawwifParams) (*MuthawwifResult, error) {
	const op = "catalog_grpc_adapter.Adapter.UpdateMuthawwif"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("id", params.ID))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.mastersClient().UpdateMuthawwif(ctx, &pb.UpdateMuthawwifRequest{
		UserId:      params.UserID,
		Id:          params.ID,
		Name:        params.Name,
		PortraitUrl: params.PortraitURL,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return fromProtoMuthawwif(resp.GetMuthawwif()), nil
}

// DeleteMuthawwif soft-deletes a muthawwif master record.
func (a *Adapter) DeleteMuthawwif(ctx context.Context, userID, id string) error {
	const op = "catalog_grpc_adapter.Adapter.DeleteMuthawwif"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("id", id))

	logger := logging.LogWithTrace(ctx, a.logger)

	_, err := a.mastersClient().DeleteMuthawwif(ctx, &pb.DeleteMuthawwifRequest{
		UserId: userID,
		Id:     id,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return nil
}

// ListMuthawwif returns a paginated list of muthawwif master records.
func (a *Adapter) ListMuthawwif(ctx context.Context, params *ListMastersParams) (*ListMuthawwifResult, error) {
	const op = "catalog_grpc_adapter.Adapter.ListMuthawwif"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.mastersClient().ListMuthawwif(ctx, &pb.ListMuthawwifRequest{
		UserId: params.UserID,
		Cursor: params.Cursor,
		Limit:  params.Limit,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	muthawwifList := make([]*MuthawwifResult, 0, len(resp.GetMuthawwifList()))
	for _, m := range resp.GetMuthawwifList() {
		muthawwifList = append(muthawwifList, fromProtoMuthawwif(m))
	}

	span.SetStatus(codes.Ok, "ok")
	return &ListMuthawwifResult{
		Muthawwif: muthawwifList,
		HasMore:   resp.GetHasMore(),
		Cursor:    resp.GetCursor(),
	}, nil
}

// ---------------------------------------------------------------------------
// Addon methods
// ---------------------------------------------------------------------------

// CreateAddon creates a new addon master record.
func (a *Adapter) CreateAddon(ctx context.Context, params *CreateAddonParams) (*AddonResult, error) {
	const op = "catalog_grpc_adapter.Adapter.CreateAddon"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.mastersClient().CreateAddon(ctx, &pb.CreateAddonRequest{
		UserId:        params.UserID,
		Name:          params.Name,
		ListAmountIdr: params.ListAmountIdr,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return fromProtoAddon(resp.GetAddon()), nil
}

// UpdateAddon updates an existing addon master record.
func (a *Adapter) UpdateAddon(ctx context.Context, params *UpdateAddonParams) (*AddonResult, error) {
	const op = "catalog_grpc_adapter.Adapter.UpdateAddon"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("id", params.ID))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.mastersClient().UpdateAddon(ctx, &pb.UpdateAddonRequest{
		UserId:        params.UserID,
		Id:            params.ID,
		Name:          params.Name,
		ListAmountIdr: params.ListAmountIdr,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return fromProtoAddon(resp.GetAddon()), nil
}

// DeleteAddon soft-deletes an addon master record.
func (a *Adapter) DeleteAddon(ctx context.Context, userID, id string) error {
	const op = "catalog_grpc_adapter.Adapter.DeleteAddon"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("id", id))

	logger := logging.LogWithTrace(ctx, a.logger)

	_, err := a.mastersClient().DeleteAddon(ctx, &pb.DeleteAddonRequest{
		UserId: userID,
		Id:     id,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return nil
}

// ListAddons returns a paginated list of addon master records.
func (a *Adapter) ListAddons(ctx context.Context, params *ListMastersParams) (*ListAddonsResult, error) {
	const op = "catalog_grpc_adapter.Adapter.ListAddons"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.mastersClient().ListAddons(ctx, &pb.ListAddonsRequest{
		UserId: params.UserID,
		Cursor: params.Cursor,
		Limit:  params.Limit,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	addons := make([]*AddonResult, 0, len(resp.GetAddons()))
	for _, ad := range resp.GetAddons() {
		addons = append(addons, fromProtoAddon(ad))
	}

	span.SetStatus(codes.Ok, "ok")
	return &ListAddonsResult{
		Addons:  addons,
		HasMore: resp.GetHasMore(),
		Cursor:  resp.GetCursor(),
	}, nil
}

// ---------------------------------------------------------------------------
// Departure Pricing methods
// ---------------------------------------------------------------------------

// SetDeparturePricing upserts room-type pricing for a departure.
func (a *Adapter) SetDeparturePricing(ctx context.Context, userID, departureID string, pricings []*PricingUpsertInput) (*PricingResult, error) {
	const op = "catalog_grpc_adapter.Adapter.SetDeparturePricing"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("departure_id", departureID))

	logger := logging.LogWithTrace(ctx, a.logger)

	pbPricings := make([]*pb.MasterPricingUpsertInput, 0, len(pricings))
	for _, p := range pricings {
		pbPricings = append(pbPricings, &pb.MasterPricingUpsertInput{
			RoomType:      p.RoomType,
			ListAmountIdr: p.ListAmountIdr,
		})
	}

	resp, err := a.mastersClient().SetDeparturePricing(ctx, &pb.SetDeparturePricingRequest{
		UserId:      userID,
		DepartureId: departureID,
		Pricings:    pbPricings,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	result := make([]*DeparturePricingResult, 0, len(resp.GetPricings()))
	for _, p := range resp.GetPricings() {
		result = append(result, fromProtoDeparturePricing(p))
	}

	span.SetStatus(codes.Ok, "ok")
	return &PricingResult{Pricings: result}, nil
}

// GetDeparturePricing returns the current pricing rows for a departure.
func (a *Adapter) GetDeparturePricing(ctx context.Context, userID, departureID string) (*PricingResult, error) {
	const op = "catalog_grpc_adapter.Adapter.GetDeparturePricing"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("departure_id", departureID))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.mastersClient().GetDeparturePricing(ctx, &pb.GetDeparturePricingRequest{
		UserId:      userID,
		DepartureId: departureID,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	result := make([]*DeparturePricingResult, 0, len(resp.GetPricings()))
	for _, p := range resp.GetPricings() {
		result = append(result, fromProtoDeparturePricing(p))
	}

	span.SetStatus(codes.Ok, "ok")
	return &PricingResult{Pricings: result}, nil
}
