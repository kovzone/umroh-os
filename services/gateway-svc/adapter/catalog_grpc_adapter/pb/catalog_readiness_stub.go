// catalog_readiness_stub.go — gateway-side gRPC client stub for catalog-svc
// vendor readiness RPCs (BL-OPS-020).
//
// Mirrors services/catalog-svc/api/grpc_api/pb/vendor_readiness_messages.go
// and vendor_readiness_grpc_ext.go. Run `make genpb` to replace with
// generated code once protoc is available.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// Method name constants
// ---------------------------------------------------------------------------

const (
	CatalogService_UpdateVendorReadiness_FullMethodName = "/pb.CatalogService/UpdateVendorReadiness"
	CatalogService_GetDepartureReadiness_FullMethodName = "/pb.CatalogService/GetDepartureReadiness"
)

// ---------------------------------------------------------------------------
// Domain object
// ---------------------------------------------------------------------------

// ReadinessState carries the ticket/hotel/visa state for a departure.
type ReadinessState struct {
	TicketState string
	HotelState  string
	VisaState   string
}

func (x *ReadinessState) GetTicketState() string {
	if x == nil {
		return ""
	}
	return x.TicketState
}
func (x *ReadinessState) GetHotelState() string {
	if x == nil {
		return ""
	}
	return x.HotelState
}
func (x *ReadinessState) GetVisaState() string {
	if x == nil {
		return ""
	}
	return x.VisaState
}

// ---------------------------------------------------------------------------
// UpdateVendorReadiness
// ---------------------------------------------------------------------------

type UpdateVendorReadinessRequest struct {
	UserId        string
	DepartureId   string
	Kind          string
	State         string
	Notes         string
	AttachmentUrl string
}

func (x *UpdateVendorReadinessRequest) GetUserId() string        { return x.UserId }
func (x *UpdateVendorReadinessRequest) GetDepartureId() string   { return x.DepartureId }
func (x *UpdateVendorReadinessRequest) GetKind() string          { return x.Kind }
func (x *UpdateVendorReadinessRequest) GetState() string         { return x.State }
func (x *UpdateVendorReadinessRequest) GetNotes() string         { return x.Notes }
func (x *UpdateVendorReadinessRequest) GetAttachmentUrl() string { return x.AttachmentUrl }

type UpdateVendorReadinessResponse struct {
	Readiness *ReadinessState
}

func (x *UpdateVendorReadinessResponse) GetReadiness() *ReadinessState {
	if x == nil {
		return nil
	}
	return x.Readiness
}

// ---------------------------------------------------------------------------
// GetDepartureReadiness
// ---------------------------------------------------------------------------

type GetDepartureReadinessRequest struct {
	UserId      string
	DepartureId string
}

func (x *GetDepartureReadinessRequest) GetUserId() string      { return x.UserId }
func (x *GetDepartureReadinessRequest) GetDepartureId() string { return x.DepartureId }

type GetDepartureReadinessResponse struct {
	Readiness *ReadinessState
}

func (x *GetDepartureReadinessResponse) GetReadiness() *ReadinessState {
	if x == nil {
		return nil
	}
	return x.Readiness
}

// ---------------------------------------------------------------------------
// CatalogReadinessClient (gRPC client stub)
// ---------------------------------------------------------------------------

// CatalogReadinessClient is the client API for the vendor readiness RPCs.
type CatalogReadinessClient interface {
	UpdateVendorReadiness(ctx context.Context, in *UpdateVendorReadinessRequest, opts ...grpc.CallOption) (*UpdateVendorReadinessResponse, error)
	GetDepartureReadiness(ctx context.Context, in *GetDepartureReadinessRequest, opts ...grpc.CallOption) (*GetDepartureReadinessResponse, error)
}

type catalogReadinessClient struct {
	cc grpc.ClientConnInterface
}

func NewCatalogReadinessClient(cc grpc.ClientConnInterface) CatalogReadinessClient {
	return &catalogReadinessClient{cc}
}

func (c *catalogReadinessClient) UpdateVendorReadiness(ctx context.Context, in *UpdateVendorReadinessRequest, opts ...grpc.CallOption) (*UpdateVendorReadinessResponse, error) {
	out := new(UpdateVendorReadinessResponse)
	err := c.cc.Invoke(ctx, CatalogService_UpdateVendorReadiness_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogReadinessClient) GetDepartureReadiness(ctx context.Context, in *GetDepartureReadinessRequest, opts ...grpc.CallOption) (*GetDepartureReadinessResponse, error) {
	out := new(GetDepartureReadinessResponse)
	err := c.cc.Invoke(ctx, CatalogService_GetDepartureReadiness_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
