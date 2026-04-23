// crm_stub.go — minimal hand-written gRPC client stub for booking-svc
// to call crm-svc (S4-E-02).
//
// Covers only the two RPCs booking-svc uses as a producer:
//   - OnBookingCreated  (called on draft creation)
//   - OnBookingPaidInFull (called on paid_in_full)
//
// Run `make genpb` with a shared proto path to replace with generated code.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

const (
	CrmService_OnBookingCreated_FullMethodName    = "/pb.CrmService/OnBookingCreated"
	CrmService_OnBookingPaidInFull_FullMethodName = "/pb.CrmService/OnBookingPaidInFull"
)

// ---------------------------------------------------------------------------
// OnBookingCreated
// ---------------------------------------------------------------------------

type OnBookingCreatedRequest struct {
	BookingId   string
	LeadId      string
	PackageId   string
	DepartureId string
	JamaahCount int32
	CreatedAt   string
}

type OnBookingCreatedResponse struct {
	Updated bool
	LeadId  string
}

func (x *OnBookingCreatedResponse) GetUpdated() bool {
	if x == nil {
		return false
	}
	return x.Updated
}
func (x *OnBookingCreatedResponse) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}

// ---------------------------------------------------------------------------
// OnBookingPaidInFull
// ---------------------------------------------------------------------------

type OnBookingPaidInFullRequest struct {
	BookingId string
	LeadId    string
	PaidAt    string
}

type OnBookingPaidInFullResponse struct {
	Updated bool
	LeadId  string
}

func (x *OnBookingPaidInFullResponse) GetUpdated() bool {
	if x == nil {
		return false
	}
	return x.Updated
}
func (x *OnBookingPaidInFullResponse) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}

// ---------------------------------------------------------------------------
// CrmServiceClient interface + implementation
// ---------------------------------------------------------------------------

type CrmServiceClient interface {
	OnBookingCreated(ctx context.Context, in *OnBookingCreatedRequest, opts ...grpc.CallOption) (*OnBookingCreatedResponse, error)
	OnBookingPaidInFull(ctx context.Context, in *OnBookingPaidInFullRequest, opts ...grpc.CallOption) (*OnBookingPaidInFullResponse, error)
}

type crmServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCrmServiceClient(cc grpc.ClientConnInterface) CrmServiceClient {
	return &crmServiceClient{cc}
}

func (c *crmServiceClient) OnBookingCreated(ctx context.Context, in *OnBookingCreatedRequest, opts ...grpc.CallOption) (*OnBookingCreatedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OnBookingCreatedResponse)
	err := c.cc.Invoke(ctx, CrmService_OnBookingCreated_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crmServiceClient) OnBookingPaidInFull(ctx context.Context, in *OnBookingPaidInFullRequest, opts ...grpc.CallOption) (*OnBookingPaidInFullResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OnBookingPaidInFullResponse)
	err := c.cc.Invoke(ctx, CrmService_OnBookingPaidInFull_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
