// Hand-written gRPC client stub for booking-svc's catalog_grpc_adapter.
// Run `make generate` (protoc + protoc-gen-go-grpc) to regenerate.

package pb

import (
	"context"

	"google.golang.org/grpc"
)

const (
	CatalogService_ReserveSeats_FullMethodName        = "/pb.CatalogService/ReserveSeats"
	CatalogService_ReleaseSeats_FullMethodName        = "/pb.CatalogService/ReleaseSeats"
	CatalogService_GetPackageDeparture_FullMethodName = "/pb.CatalogService/GetPackageDeparture"
)

// CatalogServiceClient is the client API for the narrow catalog stub.
type CatalogServiceClient interface {
	ReserveSeats(ctx context.Context, in *ReserveSeatsRequest, opts ...grpc.CallOption) (*ReserveSeatsResponse, error)
	ReleaseSeats(ctx context.Context, in *ReleaseSeatsRequest, opts ...grpc.CallOption) (*ReleaseSeatsResponse, error)
	GetPackageDeparture(ctx context.Context, in *GetPackageDepartureRequest, opts ...grpc.CallOption) (*GetPackageDepartureResponse, error)
}

type catalogServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCatalogServiceClient(cc grpc.ClientConnInterface) CatalogServiceClient {
	return &catalogServiceClient{cc}
}

func (c *catalogServiceClient) ReserveSeats(ctx context.Context, in *ReserveSeatsRequest, opts ...grpc.CallOption) (*ReserveSeatsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReserveSeatsResponse)
	err := c.cc.Invoke(ctx, CatalogService_ReserveSeats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogServiceClient) ReleaseSeats(ctx context.Context, in *ReleaseSeatsRequest, opts ...grpc.CallOption) (*ReleaseSeatsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReleaseSeatsResponse)
	err := c.cc.Invoke(ctx, CatalogService_ReleaseSeats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogServiceClient) GetPackageDeparture(ctx context.Context, in *GetPackageDepartureRequest, opts ...grpc.CallOption) (*GetPackageDepartureResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPackageDepartureResponse)
	err := c.cc.Invoke(ctx, CatalogService_GetPackageDeparture_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
