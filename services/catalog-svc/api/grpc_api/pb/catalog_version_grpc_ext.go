// catalog_version_grpc_ext.go — hand-written gRPC extension for package
// versioning RPC (BL-CAT-013).
//
// Extends CatalogServiceServer with GetPackageVersion.
// Run `make genpb` to regenerate from catalog.proto once protoc is available.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Full method name constant for package versioning RPC.
const (
	CatalogService_GetPackageVersion_FullMethodName = "/pb.CatalogService/GetPackageVersion"
)

// PackageVersionHandler is the server-side interface for the version RPC.
type PackageVersionHandler interface {
	GetPackageVersion(context.Context, *GetPackageVersionRequest) (*GetPackageVersionResponse, error)
}

// UnimplementedPackageVersionHandler provides a safe default for PackageVersionHandler.
type UnimplementedPackageVersionHandler struct{}

func (UnimplementedPackageVersionHandler) GetPackageVersion(context.Context, *GetPackageVersionRequest) (*GetPackageVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPackageVersion not implemented")
}

// ---------------------------------------------------------------------------
// Handler dispatch helper
// ---------------------------------------------------------------------------

func _CatalogService_GetPackageVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPackageVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PackageVersionHandler).GetPackageVersion(ctx, req.(*GetPackageVersionRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: CatalogService_GetPackageVersion_FullMethodName}, handler)
}
