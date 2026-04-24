// catalog_bulk_grpc_ext.go — hand-written gRPC extension for bulk package
// import/update RPCs (BL-CAT-010 / BL-CAT-011).
//
// Extends CatalogServiceServer with BulkImportPackages and BulkUpdatePackages.
// Run `make genpb` to regenerate from catalog.proto once protoc is available.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Full method name constants for bulk RPCs.
const (
	CatalogService_BulkImportPackages_FullMethodName = "/pb.CatalogService/BulkImportPackages"
	CatalogService_BulkUpdatePackages_FullMethodName = "/pb.CatalogService/BulkUpdatePackages"
)

// BulkPackagesHandler is the server-side interface for bulk package RPCs.
type BulkPackagesHandler interface {
	BulkImportPackages(context.Context, *BulkImportPackagesRequest) (*BulkImportPackagesResponse, error)
	BulkUpdatePackages(context.Context, *BulkUpdatePackagesRequest) (*BulkUpdatePackagesResponse, error)
}

// UnimplementedBulkPackagesHandler provides safe defaults for BulkPackagesHandler.
type UnimplementedBulkPackagesHandler struct{}

func (UnimplementedBulkPackagesHandler) BulkImportPackages(context.Context, *BulkImportPackagesRequest) (*BulkImportPackagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BulkImportPackages not implemented")
}
func (UnimplementedBulkPackagesHandler) BulkUpdatePackages(context.Context, *BulkUpdatePackagesRequest) (*BulkUpdatePackagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BulkUpdatePackages not implemented")
}

// ---------------------------------------------------------------------------
// Handler dispatch helpers
// ---------------------------------------------------------------------------

func _CatalogService_BulkImportPackages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BulkImportPackagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BulkPackagesHandler).BulkImportPackages(ctx, req.(*BulkImportPackagesRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_BulkImportPackages_FullMethodName}, handler)
}

func _CatalogService_BulkUpdatePackages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BulkUpdatePackagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BulkPackagesHandler).BulkUpdatePackages(ctx, req.(*BulkUpdatePackagesRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_BulkUpdatePackages_FullMethodName}, handler)
}
