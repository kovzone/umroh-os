// finance_stub.go — minimal hand-written gRPC client stub for booking-svc
// to call finance-svc.OnPaymentReceived (S3-E-03).
//
// This mirrors the subset of finance-svc's pb types that booking-svc needs
// as a consumer. Run `make genpb` with a shared proto path to replace with
// generated code.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

const (
	FinanceService_OnPaymentReceived_FullMethodName = "/pb.FinanceService/OnPaymentReceived"
)

// OnPaymentReceivedRequest is the request message for FinanceService.OnPaymentReceived.
// Amount is int64 (integer IDR) per §S3-J-03 contract.
type OnPaymentReceivedRequest struct {
	BookingId  string
	InvoiceId  string
	Amount     int64  // integer IDR — no fractional amounts (§S3-J-03)
	ReceivedAt string // RFC3339; empty = server time
}

// OnPaymentReceivedResponse is the response message from FinanceService.OnPaymentReceived.
type OnPaymentReceivedResponse struct {
	EntryId  string
	Balanced bool
}

func (x *OnPaymentReceivedResponse) GetEntryId() string {
	if x == nil {
		return ""
	}
	return x.EntryId
}
func (x *OnPaymentReceivedResponse) GetBalanced() bool {
	if x == nil {
		return false
	}
	return x.Balanced
}

// FinanceServiceClient is the consumer-side interface for the subset of
// FinanceService used by booking-svc.
type FinanceServiceClient interface {
	OnPaymentReceived(ctx context.Context, in *OnPaymentReceivedRequest, opts ...grpc.CallOption) (*OnPaymentReceivedResponse, error)
}

type financeServiceClient struct {
	cc grpc.ClientConnInterface
}

// NewFinanceServiceClient creates a client for finance-svc from a dialled conn.
func NewFinanceServiceClient(cc grpc.ClientConnInterface) FinanceServiceClient {
	return &financeServiceClient{cc}
}

func (c *financeServiceClient) OnPaymentReceived(ctx context.Context, in *OnPaymentReceivedRequest, opts ...grpc.CallOption) (*OnPaymentReceivedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OnPaymentReceivedResponse)
	err := c.cc.Invoke(ctx, FinanceService_OnPaymentReceived_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
