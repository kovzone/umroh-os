// finance_grn_stub.go — gateway-side gRPC client stub for finance-svc GRN RPC
// (S3 Wave 2 / BL-FIN-002).
//
// Mirrors services/finance-svc/api/grpc_api/pb/grn_messages.go.
// Run `make genpb` to replace with generated code once protoc is available.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

const (
	FinanceService_OnGRNReceived_FullMethodName = "/pb.FinanceService/OnGRNReceived"
)

// ---------------------------------------------------------------------------
// Request / response types
// ---------------------------------------------------------------------------

type OnGRNReceivedRequest struct {
	GrnId       string
	DepartureId string
	AmountIdr   int64
}

func (x *OnGRNReceivedRequest) GetGrnId() string {
	if x == nil {
		return ""
	}
	return x.GrnId
}
func (x *OnGRNReceivedRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *OnGRNReceivedRequest) GetAmountIdr() int64 {
	if x == nil {
		return 0
	}
	return x.AmountIdr
}

type OnGRNReceivedResponse struct {
	EntryId    string
	Idempotent bool
}

func (x *OnGRNReceivedResponse) GetEntryId() string {
	if x == nil {
		return ""
	}
	return x.EntryId
}
func (x *OnGRNReceivedResponse) GetIdempotent() bool {
	if x == nil {
		return false
	}
	return x.Idempotent
}

// ---------------------------------------------------------------------------
// FinanceGRNClient interface + implementation
// ---------------------------------------------------------------------------

// FinanceGRNClient is the consumer-side interface for finance-svc GRN RPC.
type FinanceGRNClient interface {
	OnGRNReceived(ctx context.Context, req *OnGRNReceivedRequest, opts ...grpc.CallOption) (*OnGRNReceivedResponse, error)
}

type financeGRNClient struct {
	cc grpc.ClientConnInterface
}

// NewFinanceGRNClient wraps a conn so gateway-svc can call finance-svc GRN RPC.
func NewFinanceGRNClient(cc grpc.ClientConnInterface) FinanceGRNClient {
	return &financeGRNClient{cc}
}

func (c *financeGRNClient) OnGRNReceived(ctx context.Context, in *OnGRNReceivedRequest, opts ...grpc.CallOption) (*OnGRNReceivedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OnGRNReceivedResponse)
	err := c.cc.Invoke(ctx, FinanceService_OnGRNReceived_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
