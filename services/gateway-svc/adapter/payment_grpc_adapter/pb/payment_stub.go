// payment_stub.go — hand-written gRPC client stub for payment-svc (gateway-svc side).
//
// Mirrors the server-side message types defined in payment-svc/api/grpc_api/pb/
// so that the gateway adapter compiles without importing payment-svc packages
// (per the gateway adapter pattern used throughout this repo).
//
// RPCs covered:
//   - ReissuePaymentLink (BL-PAY-020)

package pb

import (
	"context"

	"google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// Message types (mirror payment-svc/api/grpc_api/pb)
// ---------------------------------------------------------------------------

// ReissuePaymentLinkRequest mirrors payment-svc's pb.ReissuePaymentLinkRequest.
type ReissuePaymentLinkRequest struct {
	BookingId   string
	BankCode    string
	GatewayPref string
	ActorUserId string
}

func (x *ReissuePaymentLinkRequest) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}
func (x *ReissuePaymentLinkRequest) GetBankCode() string {
	if x == nil {
		return ""
	}
	return x.BankCode
}
func (x *ReissuePaymentLinkRequest) GetGatewayPref() string {
	if x == nil {
		return ""
	}
	return x.GatewayPref
}
func (x *ReissuePaymentLinkRequest) GetActorUserId() string {
	if x == nil {
		return ""
	}
	return x.ActorUserId
}

// ReissuePaymentLinkResponse mirrors payment-svc's pb.ReissuePaymentLinkResponse.
type ReissuePaymentLinkResponse struct {
	InvoiceId     string
	BookingId     string
	AccountNumber string
	BankCode      string
	AmountTotal   float64
	ExpiresAt     string
	Gateway       string
	IsNew         bool
}

func (x *ReissuePaymentLinkResponse) GetInvoiceId() string {
	if x == nil {
		return ""
	}
	return x.InvoiceId
}
func (x *ReissuePaymentLinkResponse) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}
func (x *ReissuePaymentLinkResponse) GetAccountNumber() string {
	if x == nil {
		return ""
	}
	return x.AccountNumber
}
func (x *ReissuePaymentLinkResponse) GetBankCode() string {
	if x == nil {
		return ""
	}
	return x.BankCode
}
func (x *ReissuePaymentLinkResponse) GetAmountTotal() float64 {
	if x == nil {
		return 0
	}
	return x.AmountTotal
}
func (x *ReissuePaymentLinkResponse) GetExpiresAt() string {
	if x == nil {
		return ""
	}
	return x.ExpiresAt
}
func (x *ReissuePaymentLinkResponse) GetGateway() string {
	if x == nil {
		return ""
	}
	return x.Gateway
}
func (x *ReissuePaymentLinkResponse) GetIsNew() bool {
	if x == nil {
		return false
	}
	return x.IsNew
}

// ---------------------------------------------------------------------------
// Client interface + implementation
// ---------------------------------------------------------------------------

const (
	PaymentService_ReissuePaymentLink_FullMethodName = "/pb.PaymentService/ReissuePaymentLink"
)

// PaymentServiceClient is the narrow gRPC client interface used by the gateway adapter.
type PaymentServiceClient interface {
	ReissuePaymentLink(ctx context.Context, in *ReissuePaymentLinkRequest, opts ...grpc.CallOption) (*ReissuePaymentLinkResponse, error)
}

type paymentServiceClient struct {
	cc grpc.ClientConnInterface
}

// NewPaymentServiceClient returns a client backed by the given connection.
func NewPaymentServiceClient(cc grpc.ClientConnInterface) PaymentServiceClient {
	return &paymentServiceClient{cc}
}

func (c *paymentServiceClient) ReissuePaymentLink(ctx context.Context, in *ReissuePaymentLinkRequest, opts ...grpc.CallOption) (*ReissuePaymentLinkResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReissuePaymentLinkResponse)
	err := c.cc.Invoke(ctx, PaymentService_ReissuePaymentLink_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
