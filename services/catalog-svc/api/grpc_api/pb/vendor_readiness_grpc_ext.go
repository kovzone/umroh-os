// vendor_readiness_grpc_ext.go — hand-written gRPC service extension for
// vendor readiness RPCs (BL-OPS-020).
//
// Extends CatalogServiceServer with UpdateVendorReadiness and
// GetDepartureReadiness RPCs. Run `make genpb` to regenerate from
// catalog.proto once protoc is available.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Full method name constants for vendor readiness RPCs.
const (
	CatalogService_UpdateVendorReadiness_FullMethodName  = "/pb.CatalogService/UpdateVendorReadiness"
	CatalogService_GetDepartureReadiness_FullMethodName  = "/pb.CatalogService/GetDepartureReadiness"
)

// VendorReadinessHandler is the server-side interface for vendor readiness RPCs.
type VendorReadinessHandler interface {
	UpdateVendorReadiness(context.Context, *UpdateVendorReadinessRequest) (*UpdateVendorReadinessResponse, error)
	GetDepartureReadiness(context.Context, *GetDepartureReadinessRequest) (*GetDepartureReadinessResponse, error)
}

// UnimplementedVendorReadinessHandler provides safe defaults for VendorReadinessHandler.
type UnimplementedVendorReadinessHandler struct{}

func (UnimplementedVendorReadinessHandler) UpdateVendorReadiness(context.Context, *UpdateVendorReadinessRequest) (*UpdateVendorReadinessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVendorReadiness not implemented")
}
func (UnimplementedVendorReadinessHandler) GetDepartureReadiness(context.Context, *GetDepartureReadinessRequest) (*GetDepartureReadinessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDepartureReadiness not implemented")
}

// ---------------------------------------------------------------------------
// Handler dispatch helpers
// ---------------------------------------------------------------------------

func _CatalogService_UpdateVendorReadiness_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateVendorReadinessRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VendorReadinessHandler).UpdateVendorReadiness(ctx, req.(*UpdateVendorReadinessRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_UpdateVendorReadiness_FullMethodName}, handler)
}

func _CatalogService_GetDepartureReadiness_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDepartureReadinessRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VendorReadinessHandler).GetDepartureReadiness(ctx, req.(*GetDepartureReadinessRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_GetDepartureReadiness_FullMethodName}, handler)
}

// RegisterCatalogServiceServerWithAll registers the full CatalogService surface:
// original RPCs + master data RPCs + vendor readiness RPCs.
// Use this instead of RegisterCatalogServiceServerWithMasters when vendor
// readiness is wired.
func RegisterCatalogServiceServerWithAll(s grpc.ServiceRegistrar, srv interface {
	CatalogServiceServer
	MastersHandler
	VendorReadinessHandler
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
			// Vendor readiness RPCs (BL-OPS-020)
			{MethodName: "UpdateVendorReadiness", Handler: _CatalogService_UpdateVendorReadiness_Handler},
			{MethodName: "GetDepartureReadiness", Handler: _CatalogService_GetDepartureReadiness_Handler},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "catalog.proto",
	}
	s.RegisterService(&desc, srv)
}
