// catalog_masters_stub.go — gateway-side gRPC client stub for catalog-svc
// master data RPCs (Wave 1A / Phase 6).
//
// These types mirror services/catalog-svc/api/grpc_api/pb/masters_messages.go
// and the method constants in masters_grpc_ext.go. Run `make genpb` to replace
// with generated code once protoc is available.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// Method name constants for master data RPCs
// ---------------------------------------------------------------------------

const (
	CatalogService_CreateHotel_FullMethodName         = "/pb.CatalogService/CreateHotel"
	CatalogService_UpdateHotel_FullMethodName         = "/pb.CatalogService/UpdateHotel"
	CatalogService_DeleteHotel_FullMethodName         = "/pb.CatalogService/DeleteHotel"
	CatalogService_ListHotels_FullMethodName          = "/pb.CatalogService/ListHotels"
	CatalogService_CreateAirline_FullMethodName       = "/pb.CatalogService/CreateAirline"
	CatalogService_UpdateAirline_FullMethodName       = "/pb.CatalogService/UpdateAirline"
	CatalogService_DeleteAirline_FullMethodName       = "/pb.CatalogService/DeleteAirline"
	CatalogService_ListAirlines_FullMethodName        = "/pb.CatalogService/ListAirlines"
	CatalogService_CreateMuthawwif_FullMethodName     = "/pb.CatalogService/CreateMuthawwif"
	CatalogService_UpdateMuthawwif_FullMethodName     = "/pb.CatalogService/UpdateMuthawwif"
	CatalogService_DeleteMuthawwif_FullMethodName     = "/pb.CatalogService/DeleteMuthawwif"
	CatalogService_ListMuthawwif_FullMethodName       = "/pb.CatalogService/ListMuthawwif"
	CatalogService_CreateAddon_FullMethodName         = "/pb.CatalogService/CreateAddon"
	CatalogService_UpdateAddon_FullMethodName         = "/pb.CatalogService/UpdateAddon"
	CatalogService_DeleteAddon_FullMethodName         = "/pb.CatalogService/DeleteAddon"
	CatalogService_ListAddons_FullMethodName          = "/pb.CatalogService/ListAddons"
	CatalogService_SetDeparturePricing_FullMethodName = "/pb.CatalogService/SetDeparturePricing"
	CatalogService_GetDeparturePricing_FullMethodName = "/pb.CatalogService/GetDeparturePricing"
)

// ---------------------------------------------------------------------------
// Domain objects
// ---------------------------------------------------------------------------

// MasterHotel mirrors catalog.hotels for wire use.
type MasterHotel struct {
	Id               string
	Name             string
	City             string
	StarRating       int32
	WalkingDistanceM int32
	CreatedAt        string
	UpdatedAt        string
}

func (x *MasterHotel) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *MasterHotel) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *MasterHotel) GetCity() string {
	if x == nil {
		return ""
	}
	return x.City
}
func (x *MasterHotel) GetStarRating() int32 {
	if x == nil {
		return 0
	}
	return x.StarRating
}
func (x *MasterHotel) GetWalkingDistanceM() int32 {
	if x == nil {
		return 0
	}
	return x.WalkingDistanceM
}
func (x *MasterHotel) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}
func (x *MasterHotel) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

// MasterAirline mirrors catalog.airlines for wire use.
type MasterAirline struct {
	Id           string
	Code         string
	Name         string
	OperatorKind string
	CreatedAt    string
	UpdatedAt    string
}

func (x *MasterAirline) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *MasterAirline) GetCode() string {
	if x == nil {
		return ""
	}
	return x.Code
}
func (x *MasterAirline) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *MasterAirline) GetOperatorKind() string {
	if x == nil {
		return ""
	}
	return x.OperatorKind
}
func (x *MasterAirline) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}
func (x *MasterAirline) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

// MasterMuthawwif mirrors catalog.muthawwif for wire use.
type MasterMuthawwif struct {
	Id          string
	Name        string
	PortraitUrl string
	CreatedAt   string
	UpdatedAt   string
}

func (x *MasterMuthawwif) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *MasterMuthawwif) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *MasterMuthawwif) GetPortraitUrl() string {
	if x == nil {
		return ""
	}
	return x.PortraitUrl
}
func (x *MasterMuthawwif) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}
func (x *MasterMuthawwif) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

// MasterAddon mirrors catalog.addons for wire use.
type MasterAddon struct {
	Id            string
	Name          string
	ListAmountIdr int64
	CreatedAt     string
	UpdatedAt     string
}

func (x *MasterAddon) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *MasterAddon) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *MasterAddon) GetListAmountIdr() int64 {
	if x == nil {
		return 0
	}
	return x.ListAmountIdr
}
func (x *MasterAddon) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}
func (x *MasterAddon) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

// MasterDeparturePricing mirrors catalog.package_pricing for wire use.
type MasterDeparturePricing struct {
	Id            string
	DepartureId   string
	RoomType      string
	ListAmountIdr int64
	CreatedAt     string
	UpdatedAt     string
}

func (x *MasterDeparturePricing) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *MasterDeparturePricing) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *MasterDeparturePricing) GetRoomType() string {
	if x == nil {
		return ""
	}
	return x.RoomType
}
func (x *MasterDeparturePricing) GetListAmountIdr() int64 {
	if x == nil {
		return 0
	}
	return x.ListAmountIdr
}
func (x *MasterDeparturePricing) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}
func (x *MasterDeparturePricing) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

// ---------------------------------------------------------------------------
// Hotel request/response types
// ---------------------------------------------------------------------------

type CreateHotelRequest struct {
	UserId           string
	Name             string
	City             string
	StarRating       int32
	WalkingDistanceM int32
}

func (x *CreateHotelRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *CreateHotelRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *CreateHotelRequest) GetCity() string {
	if x == nil {
		return ""
	}
	return x.City
}
func (x *CreateHotelRequest) GetStarRating() int32 {
	if x == nil {
		return 0
	}
	return x.StarRating
}
func (x *CreateHotelRequest) GetWalkingDistanceM() int32 {
	if x == nil {
		return 0
	}
	return x.WalkingDistanceM
}

type CreateHotelResponse struct {
	Hotel *MasterHotel
}

func (x *CreateHotelResponse) GetHotel() *MasterHotel {
	if x == nil {
		return nil
	}
	return x.Hotel
}

type UpdateHotelRequest struct {
	UserId           string
	Id               string
	Name             string
	City             string
	StarRating       int32
	WalkingDistanceM int32
}

func (x *UpdateHotelRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *UpdateHotelRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *UpdateHotelRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *UpdateHotelRequest) GetCity() string {
	if x == nil {
		return ""
	}
	return x.City
}
func (x *UpdateHotelRequest) GetStarRating() int32 {
	if x == nil {
		return 0
	}
	return x.StarRating
}
func (x *UpdateHotelRequest) GetWalkingDistanceM() int32 {
	if x == nil {
		return 0
	}
	return x.WalkingDistanceM
}

type UpdateHotelResponse struct {
	Hotel *MasterHotel
}

func (x *UpdateHotelResponse) GetHotel() *MasterHotel {
	if x == nil {
		return nil
	}
	return x.Hotel
}

type DeleteHotelRequest struct {
	UserId string
	Id     string
}

func (x *DeleteHotelRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *DeleteHotelRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}

type DeleteHotelResponse struct {
	Ok bool
}

func (x *DeleteHotelResponse) GetOk() bool {
	if x == nil {
		return false
	}
	return x.Ok
}

type ListHotelsRequest struct {
	UserId string
	Cursor string
	Limit  int32
}

func (x *ListHotelsRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *ListHotelsRequest) GetCursor() string {
	if x == nil {
		return ""
	}
	return x.Cursor
}
func (x *ListHotelsRequest) GetLimit() int32 {
	if x == nil {
		return 0
	}
	return x.Limit
}

type ListHotelsResponse struct {
	Hotels  []*MasterHotel
	HasMore bool
	Cursor  string
}

func (x *ListHotelsResponse) GetHotels() []*MasterHotel {
	if x == nil {
		return nil
	}
	return x.Hotels
}
func (x *ListHotelsResponse) GetHasMore() bool {
	if x == nil {
		return false
	}
	return x.HasMore
}
func (x *ListHotelsResponse) GetCursor() string {
	if x == nil {
		return ""
	}
	return x.Cursor
}

// ---------------------------------------------------------------------------
// Airline request/response types
// ---------------------------------------------------------------------------

type CreateAirlineRequest struct {
	UserId       string
	Code         string
	Name         string
	OperatorKind string
}

func (x *CreateAirlineRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *CreateAirlineRequest) GetCode() string {
	if x == nil {
		return ""
	}
	return x.Code
}
func (x *CreateAirlineRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *CreateAirlineRequest) GetOperatorKind() string {
	if x == nil {
		return ""
	}
	return x.OperatorKind
}

type CreateAirlineResponse struct {
	Airline *MasterAirline
}

func (x *CreateAirlineResponse) GetAirline() *MasterAirline {
	if x == nil {
		return nil
	}
	return x.Airline
}

type UpdateAirlineRequest struct {
	UserId string
	Id     string
	Code   string
	Name   string
}

func (x *UpdateAirlineRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *UpdateAirlineRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *UpdateAirlineRequest) GetCode() string {
	if x == nil {
		return ""
	}
	return x.Code
}
func (x *UpdateAirlineRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}

type UpdateAirlineResponse struct {
	Airline *MasterAirline
}

func (x *UpdateAirlineResponse) GetAirline() *MasterAirline {
	if x == nil {
		return nil
	}
	return x.Airline
}

type DeleteAirlineRequest struct {
	UserId string
	Id     string
}

func (x *DeleteAirlineRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *DeleteAirlineRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}

type DeleteAirlineResponse struct {
	Ok bool
}

func (x *DeleteAirlineResponse) GetOk() bool {
	if x == nil {
		return false
	}
	return x.Ok
}

type ListAirlinesRequest struct {
	UserId string
	Cursor string
	Limit  int32
}

func (x *ListAirlinesRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *ListAirlinesRequest) GetCursor() string {
	if x == nil {
		return ""
	}
	return x.Cursor
}
func (x *ListAirlinesRequest) GetLimit() int32 {
	if x == nil {
		return 0
	}
	return x.Limit
}

type ListAirlinesResponse struct {
	Airlines []*MasterAirline
	HasMore  bool
	Cursor   string
}

func (x *ListAirlinesResponse) GetAirlines() []*MasterAirline {
	if x == nil {
		return nil
	}
	return x.Airlines
}
func (x *ListAirlinesResponse) GetHasMore() bool {
	if x == nil {
		return false
	}
	return x.HasMore
}
func (x *ListAirlinesResponse) GetCursor() string {
	if x == nil {
		return ""
	}
	return x.Cursor
}

// ---------------------------------------------------------------------------
// Muthawwif request/response types
// ---------------------------------------------------------------------------

type CreateMuthawwifRequest struct {
	UserId      string
	Name        string
	PortraitUrl string
}

func (x *CreateMuthawwifRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *CreateMuthawwifRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *CreateMuthawwifRequest) GetPortraitUrl() string {
	if x == nil {
		return ""
	}
	return x.PortraitUrl
}

type CreateMuthawwifResponse struct {
	Muthawwif *MasterMuthawwif
}

func (x *CreateMuthawwifResponse) GetMuthawwif() *MasterMuthawwif {
	if x == nil {
		return nil
	}
	return x.Muthawwif
}

type UpdateMuthawwifRequest struct {
	UserId      string
	Id          string
	Name        string
	PortraitUrl string
}

func (x *UpdateMuthawwifRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *UpdateMuthawwifRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *UpdateMuthawwifRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *UpdateMuthawwifRequest) GetPortraitUrl() string {
	if x == nil {
		return ""
	}
	return x.PortraitUrl
}

type UpdateMuthawwifResponse struct {
	Muthawwif *MasterMuthawwif
}

func (x *UpdateMuthawwifResponse) GetMuthawwif() *MasterMuthawwif {
	if x == nil {
		return nil
	}
	return x.Muthawwif
}

type DeleteMuthawwifRequest struct {
	UserId string
	Id     string
}

func (x *DeleteMuthawwifRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *DeleteMuthawwifRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}

type DeleteMuthawwifResponse struct {
	Ok bool
}

func (x *DeleteMuthawwifResponse) GetOk() bool {
	if x == nil {
		return false
	}
	return x.Ok
}

type ListMuthawwifRequest struct {
	UserId string
	Cursor string
	Limit  int32
}

func (x *ListMuthawwifRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *ListMuthawwifRequest) GetCursor() string {
	if x == nil {
		return ""
	}
	return x.Cursor
}
func (x *ListMuthawwifRequest) GetLimit() int32 {
	if x == nil {
		return 0
	}
	return x.Limit
}

type ListMuthawwifResponse struct {
	MuthawwifList []*MasterMuthawwif
	HasMore       bool
	Cursor        string
}

func (x *ListMuthawwifResponse) GetMuthawwifList() []*MasterMuthawwif {
	if x == nil {
		return nil
	}
	return x.MuthawwifList
}
func (x *ListMuthawwifResponse) GetHasMore() bool {
	if x == nil {
		return false
	}
	return x.HasMore
}
func (x *ListMuthawwifResponse) GetCursor() string {
	if x == nil {
		return ""
	}
	return x.Cursor
}

// ---------------------------------------------------------------------------
// Addon request/response types
// ---------------------------------------------------------------------------

type CreateAddonRequest struct {
	UserId        string
	Name          string
	ListAmountIdr int64
}

func (x *CreateAddonRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *CreateAddonRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *CreateAddonRequest) GetListAmountIdr() int64 {
	if x == nil {
		return 0
	}
	return x.ListAmountIdr
}

type CreateAddonResponse struct {
	Addon *MasterAddon
}

func (x *CreateAddonResponse) GetAddon() *MasterAddon {
	if x == nil {
		return nil
	}
	return x.Addon
}

type UpdateAddonRequest struct {
	UserId        string
	Id            string
	Name          string
	ListAmountIdr int64
}

func (x *UpdateAddonRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *UpdateAddonRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *UpdateAddonRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *UpdateAddonRequest) GetListAmountIdr() int64 {
	if x == nil {
		return 0
	}
	return x.ListAmountIdr
}

type UpdateAddonResponse struct {
	Addon *MasterAddon
}

func (x *UpdateAddonResponse) GetAddon() *MasterAddon {
	if x == nil {
		return nil
	}
	return x.Addon
}

type DeleteAddonRequest struct {
	UserId string
	Id     string
}

func (x *DeleteAddonRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *DeleteAddonRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}

type DeleteAddonResponse struct {
	Ok bool
}

func (x *DeleteAddonResponse) GetOk() bool {
	if x == nil {
		return false
	}
	return x.Ok
}

type ListAddonsRequest struct {
	UserId string
	Cursor string
	Limit  int32
}

func (x *ListAddonsRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *ListAddonsRequest) GetCursor() string {
	if x == nil {
		return ""
	}
	return x.Cursor
}
func (x *ListAddonsRequest) GetLimit() int32 {
	if x == nil {
		return 0
	}
	return x.Limit
}

type ListAddonsResponse struct {
	Addons  []*MasterAddon
	HasMore bool
	Cursor  string
}

func (x *ListAddonsResponse) GetAddons() []*MasterAddon {
	if x == nil {
		return nil
	}
	return x.Addons
}
func (x *ListAddonsResponse) GetHasMore() bool {
	if x == nil {
		return false
	}
	return x.HasMore
}
func (x *ListAddonsResponse) GetCursor() string {
	if x == nil {
		return ""
	}
	return x.Cursor
}

// ---------------------------------------------------------------------------
// Departure Pricing request/response types
// ---------------------------------------------------------------------------

// MasterPricingUpsertInput is a single room-type price for SetDeparturePricing.
type MasterPricingUpsertInput struct {
	RoomType      string
	ListAmountIdr int64
}

func (x *MasterPricingUpsertInput) GetRoomType() string {
	if x == nil {
		return ""
	}
	return x.RoomType
}
func (x *MasterPricingUpsertInput) GetListAmountIdr() int64 {
	if x == nil {
		return 0
	}
	return x.ListAmountIdr
}

type SetDeparturePricingRequest struct {
	UserId      string
	DepartureId string
	Pricings    []*MasterPricingUpsertInput
}

func (x *SetDeparturePricingRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *SetDeparturePricingRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *SetDeparturePricingRequest) GetPricings() []*MasterPricingUpsertInput {
	if x == nil {
		return nil
	}
	return x.Pricings
}

type SetDeparturePricingResponse struct {
	Pricings []*MasterDeparturePricing
}

func (x *SetDeparturePricingResponse) GetPricings() []*MasterDeparturePricing {
	if x == nil {
		return nil
	}
	return x.Pricings
}

type GetDeparturePricingRequest struct {
	UserId      string
	DepartureId string
}

func (x *GetDeparturePricingRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *GetDeparturePricingRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}

type GetDeparturePricingResponse struct {
	Pricings []*MasterDeparturePricing
}

func (x *GetDeparturePricingResponse) GetPricings() []*MasterDeparturePricing {
	if x == nil {
		return nil
	}
	return x.Pricings
}

// ---------------------------------------------------------------------------
// CatalogMastersClient interface + implementation
// ---------------------------------------------------------------------------

// CatalogMastersClient is the consumer-side interface for catalog-svc master
// data RPCs. Used by gateway-svc → catalog-svc.
type CatalogMastersClient interface {
	// Hotel
	CreateHotel(ctx context.Context, in *CreateHotelRequest, opts ...grpc.CallOption) (*CreateHotelResponse, error)
	UpdateHotel(ctx context.Context, in *UpdateHotelRequest, opts ...grpc.CallOption) (*UpdateHotelResponse, error)
	DeleteHotel(ctx context.Context, in *DeleteHotelRequest, opts ...grpc.CallOption) (*DeleteHotelResponse, error)
	ListHotels(ctx context.Context, in *ListHotelsRequest, opts ...grpc.CallOption) (*ListHotelsResponse, error)
	// Airline
	CreateAirline(ctx context.Context, in *CreateAirlineRequest, opts ...grpc.CallOption) (*CreateAirlineResponse, error)
	UpdateAirline(ctx context.Context, in *UpdateAirlineRequest, opts ...grpc.CallOption) (*UpdateAirlineResponse, error)
	DeleteAirline(ctx context.Context, in *DeleteAirlineRequest, opts ...grpc.CallOption) (*DeleteAirlineResponse, error)
	ListAirlines(ctx context.Context, in *ListAirlinesRequest, opts ...grpc.CallOption) (*ListAirlinesResponse, error)
	// Muthawwif
	CreateMuthawwif(ctx context.Context, in *CreateMuthawwifRequest, opts ...grpc.CallOption) (*CreateMuthawwifResponse, error)
	UpdateMuthawwif(ctx context.Context, in *UpdateMuthawwifRequest, opts ...grpc.CallOption) (*UpdateMuthawwifResponse, error)
	DeleteMuthawwif(ctx context.Context, in *DeleteMuthawwifRequest, opts ...grpc.CallOption) (*DeleteMuthawwifResponse, error)
	ListMuthawwif(ctx context.Context, in *ListMuthawwifRequest, opts ...grpc.CallOption) (*ListMuthawwifResponse, error)
	// Addon
	CreateAddon(ctx context.Context, in *CreateAddonRequest, opts ...grpc.CallOption) (*CreateAddonResponse, error)
	UpdateAddon(ctx context.Context, in *UpdateAddonRequest, opts ...grpc.CallOption) (*UpdateAddonResponse, error)
	DeleteAddon(ctx context.Context, in *DeleteAddonRequest, opts ...grpc.CallOption) (*DeleteAddonResponse, error)
	ListAddons(ctx context.Context, in *ListAddonsRequest, opts ...grpc.CallOption) (*ListAddonsResponse, error)
	// Departure pricing
	SetDeparturePricing(ctx context.Context, in *SetDeparturePricingRequest, opts ...grpc.CallOption) (*SetDeparturePricingResponse, error)
	GetDeparturePricing(ctx context.Context, in *GetDeparturePricingRequest, opts ...grpc.CallOption) (*GetDeparturePricingResponse, error)
}

type catalogMastersClient struct {
	cc grpc.ClientConnInterface
}

// NewCatalogMastersClient wraps a conn so gateway-svc can call catalog master RPCs.
func NewCatalogMastersClient(cc grpc.ClientConnInterface) CatalogMastersClient {
	return &catalogMastersClient{cc}
}

func (c *catalogMastersClient) CreateHotel(ctx context.Context, in *CreateHotelRequest, opts ...grpc.CallOption) (*CreateHotelResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateHotelResponse)
	err := c.cc.Invoke(ctx, CatalogService_CreateHotel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) UpdateHotel(ctx context.Context, in *UpdateHotelRequest, opts ...grpc.CallOption) (*UpdateHotelResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateHotelResponse)
	err := c.cc.Invoke(ctx, CatalogService_UpdateHotel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) DeleteHotel(ctx context.Context, in *DeleteHotelRequest, opts ...grpc.CallOption) (*DeleteHotelResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteHotelResponse)
	err := c.cc.Invoke(ctx, CatalogService_DeleteHotel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) ListHotels(ctx context.Context, in *ListHotelsRequest, opts ...grpc.CallOption) (*ListHotelsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListHotelsResponse)
	err := c.cc.Invoke(ctx, CatalogService_ListHotels_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) CreateAirline(ctx context.Context, in *CreateAirlineRequest, opts ...grpc.CallOption) (*CreateAirlineResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateAirlineResponse)
	err := c.cc.Invoke(ctx, CatalogService_CreateAirline_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) UpdateAirline(ctx context.Context, in *UpdateAirlineRequest, opts ...grpc.CallOption) (*UpdateAirlineResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateAirlineResponse)
	err := c.cc.Invoke(ctx, CatalogService_UpdateAirline_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) DeleteAirline(ctx context.Context, in *DeleteAirlineRequest, opts ...grpc.CallOption) (*DeleteAirlineResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteAirlineResponse)
	err := c.cc.Invoke(ctx, CatalogService_DeleteAirline_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) ListAirlines(ctx context.Context, in *ListAirlinesRequest, opts ...grpc.CallOption) (*ListAirlinesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListAirlinesResponse)
	err := c.cc.Invoke(ctx, CatalogService_ListAirlines_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) CreateMuthawwif(ctx context.Context, in *CreateMuthawwifRequest, opts ...grpc.CallOption) (*CreateMuthawwifResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateMuthawwifResponse)
	err := c.cc.Invoke(ctx, CatalogService_CreateMuthawwif_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) UpdateMuthawwif(ctx context.Context, in *UpdateMuthawwifRequest, opts ...grpc.CallOption) (*UpdateMuthawwifResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMuthawwifResponse)
	err := c.cc.Invoke(ctx, CatalogService_UpdateMuthawwif_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) DeleteMuthawwif(ctx context.Context, in *DeleteMuthawwifRequest, opts ...grpc.CallOption) (*DeleteMuthawwifResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteMuthawwifResponse)
	err := c.cc.Invoke(ctx, CatalogService_DeleteMuthawwif_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) ListMuthawwif(ctx context.Context, in *ListMuthawwifRequest, opts ...grpc.CallOption) (*ListMuthawwifResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListMuthawwifResponse)
	err := c.cc.Invoke(ctx, CatalogService_ListMuthawwif_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) CreateAddon(ctx context.Context, in *CreateAddonRequest, opts ...grpc.CallOption) (*CreateAddonResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateAddonResponse)
	err := c.cc.Invoke(ctx, CatalogService_CreateAddon_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) UpdateAddon(ctx context.Context, in *UpdateAddonRequest, opts ...grpc.CallOption) (*UpdateAddonResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateAddonResponse)
	err := c.cc.Invoke(ctx, CatalogService_UpdateAddon_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) DeleteAddon(ctx context.Context, in *DeleteAddonRequest, opts ...grpc.CallOption) (*DeleteAddonResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteAddonResponse)
	err := c.cc.Invoke(ctx, CatalogService_DeleteAddon_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) ListAddons(ctx context.Context, in *ListAddonsRequest, opts ...grpc.CallOption) (*ListAddonsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListAddonsResponse)
	err := c.cc.Invoke(ctx, CatalogService_ListAddons_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) SetDeparturePricing(ctx context.Context, in *SetDeparturePricingRequest, opts ...grpc.CallOption) (*SetDeparturePricingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetDeparturePricingResponse)
	err := c.cc.Invoke(ctx, CatalogService_SetDeparturePricing_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogMastersClient) GetDeparturePricing(ctx context.Context, in *GetDeparturePricingRequest, opts ...grpc.CallOption) (*GetDeparturePricingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetDeparturePricingResponse)
	err := c.cc.Invoke(ctx, CatalogService_GetDeparturePricing_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
