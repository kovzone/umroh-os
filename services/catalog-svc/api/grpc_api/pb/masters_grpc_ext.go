// masters_grpc_ext.go — hand-written gRPC service extension for master data CRUD.
//
// Extends CatalogServiceServer with hotel, airline, muthawwif, addon and
// departure-pricing RPCs. Run `make genpb` to regenerate cleanly from
// catalog.proto once protoc is available.
//
// Pattern mirrors catalog_grpc.pb.go (hand-extended in S1-E-07).

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Full method name constants for master data RPCs.
const (
	CatalogService_CreateHotel_FullMethodName          = "/pb.CatalogService/CreateHotel"
	CatalogService_UpdateHotel_FullMethodName          = "/pb.CatalogService/UpdateHotel"
	CatalogService_DeleteHotel_FullMethodName          = "/pb.CatalogService/DeleteHotel"
	CatalogService_ListHotels_FullMethodName           = "/pb.CatalogService/ListHotels"
	CatalogService_CreateAirline_FullMethodName        = "/pb.CatalogService/CreateAirline"
	CatalogService_UpdateAirline_FullMethodName        = "/pb.CatalogService/UpdateAirline"
	CatalogService_DeleteAirline_FullMethodName        = "/pb.CatalogService/DeleteAirline"
	CatalogService_ListAirlines_FullMethodName         = "/pb.CatalogService/ListAirlines"
	CatalogService_CreateMuthawwif_FullMethodName      = "/pb.CatalogService/CreateMuthawwif"
	CatalogService_UpdateMuthawwif_FullMethodName      = "/pb.CatalogService/UpdateMuthawwif"
	CatalogService_DeleteMuthawwif_FullMethodName      = "/pb.CatalogService/DeleteMuthawwif"
	CatalogService_ListMuthawwif_FullMethodName        = "/pb.CatalogService/ListMuthawwif"
	CatalogService_CreateAddon_FullMethodName          = "/pb.CatalogService/CreateAddon"
	CatalogService_UpdateAddon_FullMethodName          = "/pb.CatalogService/UpdateAddon"
	CatalogService_DeleteAddon_FullMethodName          = "/pb.CatalogService/DeleteAddon"
	CatalogService_ListAddons_FullMethodName           = "/pb.CatalogService/ListAddons"
	CatalogService_SetDeparturePricing_FullMethodName  = "/pb.CatalogService/SetDeparturePricing"
	CatalogService_GetDeparturePricing_FullMethodName  = "/pb.CatalogService/GetDeparturePricing"
)

// MastersHandler is the server-side interface for all master data RPCs.
type MastersHandler interface {
	// Hotel
	CreateHotel(context.Context, *CreateHotelRequest) (*CreateHotelResponse, error)
	UpdateHotel(context.Context, *UpdateHotelRequest) (*UpdateHotelResponse, error)
	DeleteHotel(context.Context, *DeleteHotelRequest) (*DeleteHotelResponse, error)
	ListHotels(context.Context, *ListHotelsRequest) (*ListHotelsResponse, error)
	// Airline
	CreateAirline(context.Context, *CreateAirlineRequest) (*CreateAirlineResponse, error)
	UpdateAirline(context.Context, *UpdateAirlineRequest) (*UpdateAirlineResponse, error)
	DeleteAirline(context.Context, *DeleteAirlineRequest) (*DeleteAirlineResponse, error)
	ListAirlines(context.Context, *ListAirlinesRequest) (*ListAirlinesResponse, error)
	// Muthawwif
	CreateMuthawwif(context.Context, *CreateMuthawwifRequest) (*CreateMuthawwifResponse, error)
	UpdateMuthawwif(context.Context, *UpdateMuthawwifRequest) (*UpdateMuthawwifResponse, error)
	DeleteMuthawwif(context.Context, *DeleteMuthawwifRequest) (*DeleteMuthawwifResponse, error)
	ListMuthawwif(context.Context, *ListMuthawwifRequest) (*ListMuthawwifResponse, error)
	// Addon
	CreateAddon(context.Context, *CreateAddonRequest) (*CreateAddonResponse, error)
	UpdateAddon(context.Context, *UpdateAddonRequest) (*UpdateAddonResponse, error)
	DeleteAddon(context.Context, *DeleteAddonRequest) (*DeleteAddonResponse, error)
	ListAddons(context.Context, *ListAddonsRequest) (*ListAddonsResponse, error)
	// Departure pricing
	SetDeparturePricing(context.Context, *SetDeparturePricingRequest) (*SetDeparturePricingResponse, error)
	GetDeparturePricing(context.Context, *GetDeparturePricingRequest) (*GetDeparturePricingResponse, error)
}

// UnimplementedMastersHandler provides safe defaults for MastersHandler.
type UnimplementedMastersHandler struct{}

func (UnimplementedMastersHandler) CreateHotel(context.Context, *CreateHotelRequest) (*CreateHotelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateHotel not implemented")
}
func (UnimplementedMastersHandler) UpdateHotel(context.Context, *UpdateHotelRequest) (*UpdateHotelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateHotel not implemented")
}
func (UnimplementedMastersHandler) DeleteHotel(context.Context, *DeleteHotelRequest) (*DeleteHotelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteHotel not implemented")
}
func (UnimplementedMastersHandler) ListHotels(context.Context, *ListHotelsRequest) (*ListHotelsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListHotels not implemented")
}
func (UnimplementedMastersHandler) CreateAirline(context.Context, *CreateAirlineRequest) (*CreateAirlineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAirline not implemented")
}
func (UnimplementedMastersHandler) UpdateAirline(context.Context, *UpdateAirlineRequest) (*UpdateAirlineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAirline not implemented")
}
func (UnimplementedMastersHandler) DeleteAirline(context.Context, *DeleteAirlineRequest) (*DeleteAirlineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAirline not implemented")
}
func (UnimplementedMastersHandler) ListAirlines(context.Context, *ListAirlinesRequest) (*ListAirlinesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAirlines not implemented")
}
func (UnimplementedMastersHandler) CreateMuthawwif(context.Context, *CreateMuthawwifRequest) (*CreateMuthawwifResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMuthawwif not implemented")
}
func (UnimplementedMastersHandler) UpdateMuthawwif(context.Context, *UpdateMuthawwifRequest) (*UpdateMuthawwifResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMuthawwif not implemented")
}
func (UnimplementedMastersHandler) DeleteMuthawwif(context.Context, *DeleteMuthawwifRequest) (*DeleteMuthawwifResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMuthawwif not implemented")
}
func (UnimplementedMastersHandler) ListMuthawwif(context.Context, *ListMuthawwifRequest) (*ListMuthawwifResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMuthawwif not implemented")
}
func (UnimplementedMastersHandler) CreateAddon(context.Context, *CreateAddonRequest) (*CreateAddonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAddon not implemented")
}
func (UnimplementedMastersHandler) UpdateAddon(context.Context, *UpdateAddonRequest) (*UpdateAddonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAddon not implemented")
}
func (UnimplementedMastersHandler) DeleteAddon(context.Context, *DeleteAddonRequest) (*DeleteAddonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAddon not implemented")
}
func (UnimplementedMastersHandler) ListAddons(context.Context, *ListAddonsRequest) (*ListAddonsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAddons not implemented")
}
func (UnimplementedMastersHandler) SetDeparturePricing(context.Context, *SetDeparturePricingRequest) (*SetDeparturePricingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetDeparturePricing not implemented")
}
func (UnimplementedMastersHandler) GetDeparturePricing(context.Context, *GetDeparturePricingRequest) (*GetDeparturePricingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDeparturePricing not implemented")
}

// ---------------------------------------------------------------------------
// Handler dispatch helpers (used by RegisterCatalogServiceServerWithMasters)
// ---------------------------------------------------------------------------

func _CatalogService_CreateHotel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateHotelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).CreateHotel(ctx, req.(*CreateHotelRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_CreateHotel_FullMethodName}, handler)
}

func _CatalogService_UpdateHotel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateHotelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).UpdateHotel(ctx, req.(*UpdateHotelRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_UpdateHotel_FullMethodName}, handler)
}

func _CatalogService_DeleteHotel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteHotelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).DeleteHotel(ctx, req.(*DeleteHotelRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_DeleteHotel_FullMethodName}, handler)
}

func _CatalogService_ListHotels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListHotelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).ListHotels(ctx, req.(*ListHotelsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_ListHotels_FullMethodName}, handler)
}

func _CatalogService_CreateAirline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAirlineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).CreateAirline(ctx, req.(*CreateAirlineRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_CreateAirline_FullMethodName}, handler)
}

func _CatalogService_UpdateAirline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAirlineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).UpdateAirline(ctx, req.(*UpdateAirlineRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_UpdateAirline_FullMethodName}, handler)
}

func _CatalogService_DeleteAirline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAirlineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).DeleteAirline(ctx, req.(*DeleteAirlineRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_DeleteAirline_FullMethodName}, handler)
}

func _CatalogService_ListAirlines_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAirlinesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).ListAirlines(ctx, req.(*ListAirlinesRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_ListAirlines_FullMethodName}, handler)
}

func _CatalogService_CreateMuthawwif_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMuthawwifRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).CreateMuthawwif(ctx, req.(*CreateMuthawwifRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_CreateMuthawwif_FullMethodName}, handler)
}

func _CatalogService_UpdateMuthawwif_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMuthawwifRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).UpdateMuthawwif(ctx, req.(*UpdateMuthawwifRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_UpdateMuthawwif_FullMethodName}, handler)
}

func _CatalogService_DeleteMuthawwif_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMuthawwifRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).DeleteMuthawwif(ctx, req.(*DeleteMuthawwifRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_DeleteMuthawwif_FullMethodName}, handler)
}

func _CatalogService_ListMuthawwif_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMuthawwifRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).ListMuthawwif(ctx, req.(*ListMuthawwifRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_ListMuthawwif_FullMethodName}, handler)
}

func _CatalogService_CreateAddon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAddonRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).CreateAddon(ctx, req.(*CreateAddonRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_CreateAddon_FullMethodName}, handler)
}

func _CatalogService_UpdateAddon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAddonRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).UpdateAddon(ctx, req.(*UpdateAddonRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_UpdateAddon_FullMethodName}, handler)
}

func _CatalogService_DeleteAddon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAddonRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).DeleteAddon(ctx, req.(*DeleteAddonRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_DeleteAddon_FullMethodName}, handler)
}

func _CatalogService_ListAddons_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAddonsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).ListAddons(ctx, req.(*ListAddonsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_ListAddons_FullMethodName}, handler)
}

func _CatalogService_SetDeparturePricing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetDeparturePricingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).SetDeparturePricing(ctx, req.(*SetDeparturePricingRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_SetDeparturePricing_FullMethodName}, handler)
}

func _CatalogService_GetDeparturePricing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDeparturePricingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MastersHandler).GetDeparturePricing(ctx, req.(*GetDeparturePricingRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_GetDeparturePricing_FullMethodName}, handler)
}

// RegisterCatalogServiceServerWithMasters registers the combined CatalogService
// (original RPCs + master data RPCs). Called from main instead of the plain
// RegisterCatalogServiceServer when the masters handlers are wired up.
func RegisterCatalogServiceServerWithMasters(s grpc.ServiceRegistrar, srv interface {
	CatalogServiceServer
	MastersHandler
}) {
	desc := grpc.ServiceDesc{
		ServiceName: "pb.CatalogService",
		HandlerType: (*CatalogServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Healthz",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(HealthzRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(CatalogServiceServer).Healthz(ctx, req.(*HealthzRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_Healthz_FullMethodName}, handler)
				},
			},
			{
				MethodName: "ListPackages",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(ListPackagesRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(CatalogServiceServer).ListPackages(ctx, req.(*ListPackagesRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_ListPackages_FullMethodName}, handler)
				},
			},
			{
				MethodName: "GetPackage",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(GetPackageRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(CatalogServiceServer).GetPackage(ctx, req.(*GetPackageRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_GetPackage_FullMethodName}, handler)
				},
			},
			{
				MethodName: "GetPackageDeparture",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(GetPackageDepartureRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(CatalogServiceServer).GetPackageDeparture(ctx, req.(*GetPackageDepartureRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_GetPackageDeparture_FullMethodName}, handler)
				},
			},
			{
				MethodName: "CreatePackage",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(CreatePackageRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(CatalogServiceServer).CreatePackage(ctx, req.(*CreatePackageRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_CreatePackage_FullMethodName}, handler)
				},
			},
			{
				MethodName: "UpdatePackage",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(UpdatePackageRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(CatalogServiceServer).UpdatePackage(ctx, req.(*UpdatePackageRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_UpdatePackage_FullMethodName}, handler)
				},
			},
			{
				MethodName: "DeletePackage",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(DeletePackageRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(CatalogServiceServer).DeletePackage(ctx, req.(*DeletePackageRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_DeletePackage_FullMethodName}, handler)
				},
			},
			{
				MethodName: "CreateDeparture",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(CreateDepartureRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(CatalogServiceServer).CreateDeparture(ctx, req.(*CreateDepartureRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_CreateDeparture_FullMethodName}, handler)
				},
			},
			{
				MethodName: "UpdateDeparture",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(UpdateDepartureRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(CatalogServiceServer).UpdateDeparture(ctx, req.(*UpdateDepartureRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_UpdateDeparture_FullMethodName}, handler)
				},
			},
			{
				MethodName: "ReserveSeats",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(ReserveSeatsRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(CatalogServiceServer).ReserveSeats(ctx, req.(*ReserveSeatsRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_ReserveSeats_FullMethodName}, handler)
				},
			},
			{
				MethodName: "ReleaseSeats",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(ReleaseSeatsRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(CatalogServiceServer).ReleaseSeats(ctx, req.(*ReleaseSeatsRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_ReleaseSeats_FullMethodName}, handler)
				},
			},
			// Master data RPCs
			{MethodName: "CreateHotel", Handler: _CatalogService_CreateHotel_Handler},
			{MethodName: "UpdateHotel", Handler: _CatalogService_UpdateHotel_Handler},
			{MethodName: "DeleteHotel", Handler: _CatalogService_DeleteHotel_Handler},
			{MethodName: "ListHotels", Handler: _CatalogService_ListHotels_Handler},
			{MethodName: "CreateAirline", Handler: _CatalogService_CreateAirline_Handler},
			{MethodName: "UpdateAirline", Handler: _CatalogService_UpdateAirline_Handler},
			{MethodName: "DeleteAirline", Handler: _CatalogService_DeleteAirline_Handler},
			{MethodName: "ListAirlines", Handler: _CatalogService_ListAirlines_Handler},
			{MethodName: "CreateMuthawwif", Handler: _CatalogService_CreateMuthawwif_Handler},
			{MethodName: "UpdateMuthawwif", Handler: _CatalogService_UpdateMuthawwif_Handler},
			{MethodName: "DeleteMuthawwif", Handler: _CatalogService_DeleteMuthawwif_Handler},
			{MethodName: "ListMuthawwif", Handler: _CatalogService_ListMuthawwif_Handler},
			{MethodName: "CreateAddon", Handler: _CatalogService_CreateAddon_Handler},
			{MethodName: "UpdateAddon", Handler: _CatalogService_UpdateAddon_Handler},
			{MethodName: "DeleteAddon", Handler: _CatalogService_DeleteAddon_Handler},
			{MethodName: "ListAddons", Handler: _CatalogService_ListAddons_Handler},
			{MethodName: "SetDeparturePricing", Handler: _CatalogService_SetDeparturePricing_Handler},
			{MethodName: "GetDeparturePricing", Handler: _CatalogService_GetDeparturePricing_Handler},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "catalog.proto",
	}
	s.RegisterService(&desc, srv)
}
