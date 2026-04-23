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
// IssueVirtualAccount
// ---------------------------------------------------------------------------

// IssueVirtualAccountRequest mirrors payment-svc's pb.IssueVirtualAccountRequest.
type IssueVirtualAccountRequest struct {
	BookingId   string
	AmountTotal float64
	GatewayPref string
	BankCode    string
	ActorUserId string
}

func (x *IssueVirtualAccountRequest) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}
func (x *IssueVirtualAccountRequest) GetAmountTotal() float64 {
	if x == nil {
		return 0
	}
	return x.AmountTotal
}
func (x *IssueVirtualAccountRequest) GetGatewayPref() string {
	if x == nil {
		return ""
	}
	return x.GatewayPref
}
func (x *IssueVirtualAccountRequest) GetBankCode() string {
	if x == nil {
		return ""
	}
	return x.BankCode
}
func (x *IssueVirtualAccountRequest) GetActorUserId() string {
	if x == nil {
		return ""
	}
	return x.ActorUserId
}

// IssueVirtualAccountResponse mirrors payment-svc's pb.IssueVirtualAccountResponse.
type IssueVirtualAccountResponse struct {
	InvoiceId     string
	BookingId     string
	AccountNumber string
	BankCode      string
	AmountTotal   float64
	ExpiresAt     string
	Gateway       string
	Replayed      bool
}

func (x *IssueVirtualAccountResponse) GetInvoiceId() string {
	if x == nil {
		return ""
	}
	return x.InvoiceId
}
func (x *IssueVirtualAccountResponse) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}
func (x *IssueVirtualAccountResponse) GetAccountNumber() string {
	if x == nil {
		return ""
	}
	return x.AccountNumber
}
func (x *IssueVirtualAccountResponse) GetBankCode() string {
	if x == nil {
		return ""
	}
	return x.BankCode
}
func (x *IssueVirtualAccountResponse) GetAmountTotal() float64 {
	if x == nil {
		return 0
	}
	return x.AmountTotal
}
func (x *IssueVirtualAccountResponse) GetExpiresAt() string {
	if x == nil {
		return ""
	}
	return x.ExpiresAt
}
func (x *IssueVirtualAccountResponse) GetGateway() string {
	if x == nil {
		return ""
	}
	return x.Gateway
}
func (x *IssueVirtualAccountResponse) GetReplayed() bool {
	if x == nil {
		return false
	}
	return x.Replayed
}

// ---------------------------------------------------------------------------
// GetInvoiceByID
// ---------------------------------------------------------------------------

// GetInvoiceByIDRequest mirrors payment-svc's pb.GetInvoiceByIDRequest.
type GetInvoiceByIDRequest struct {
	InvoiceId string
}

func (x *GetInvoiceByIDRequest) GetInvoiceId() string {
	if x == nil {
		return ""
	}
	return x.InvoiceId
}

// GetInvoiceByIDResponse mirrors payment-svc's pb.GetInvoiceByIDResponse.
type GetInvoiceByIDResponse struct {
	Id          string
	BookingId   string
	Status      string
	AmountTotal float64
	PaidAmount  float64
	Currency    string
	CreatedAt   string
	UpdatedAt   string
}

func (x *GetInvoiceByIDResponse) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *GetInvoiceByIDResponse) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}
func (x *GetInvoiceByIDResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *GetInvoiceByIDResponse) GetAmountTotal() float64 {
	if x == nil {
		return 0
	}
	return x.AmountTotal
}
func (x *GetInvoiceByIDResponse) GetPaidAmount() float64 {
	if x == nil {
		return 0
	}
	return x.PaidAmount
}
func (x *GetInvoiceByIDResponse) GetCurrency() string {
	if x == nil {
		return ""
	}
	return x.Currency
}
func (x *GetInvoiceByIDResponse) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}
func (x *GetInvoiceByIDResponse) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

// ---------------------------------------------------------------------------
// ProcessWebhook
// ---------------------------------------------------------------------------

// ProcessWebhookRequest mirrors payment-svc's pb.ProcessWebhookRequest.
type ProcessWebhookRequest struct {
	Gateway   string
	Payload   []byte
	Signature string
}

func (x *ProcessWebhookRequest) GetGateway() string {
	if x == nil {
		return ""
	}
	return x.Gateway
}
func (x *ProcessWebhookRequest) GetPayload() []byte {
	if x == nil {
		return nil
	}
	return x.Payload
}
func (x *ProcessWebhookRequest) GetSignature() string {
	if x == nil {
		return ""
	}
	return x.Signature
}

// ProcessWebhookResponse mirrors payment-svc's pb.ProcessWebhookResponse.
type ProcessWebhookResponse struct {
	Replayed  bool
	InvoiceId string
	NewStatus string
}

func (x *ProcessWebhookResponse) GetReplayed() bool {
	if x == nil {
		return false
	}
	return x.Replayed
}
func (x *ProcessWebhookResponse) GetInvoiceId() string {
	if x == nil {
		return ""
	}
	return x.InvoiceId
}
func (x *ProcessWebhookResponse) GetNewStatus() string {
	if x == nil {
		return ""
	}
	return x.NewStatus
}

// ---------------------------------------------------------------------------
// Client interface + implementation
// ---------------------------------------------------------------------------

const (
	PaymentService_ReissuePaymentLink_FullMethodName = "/pb.PaymentService/ReissuePaymentLink"
	PaymentService_IssueVA_FullMethodName            = "/pb.PaymentService/IssueVirtualAccount"
	PaymentService_GetInvoiceByID_FullMethodName     = "/pb.PaymentService/GetInvoiceByID"
	PaymentService_ProcessWebhook_FullMethodName     = "/pb.PaymentService/ProcessWebhook"
)

// PaymentServiceClient is the narrow gRPC client interface used by the gateway adapter.
type PaymentServiceClient interface {
	ReissuePaymentLink(ctx context.Context, in *ReissuePaymentLinkRequest, opts ...grpc.CallOption) (*ReissuePaymentLinkResponse, error)
	IssueVirtualAccount(ctx context.Context, in *IssueVirtualAccountRequest, opts ...grpc.CallOption) (*IssueVirtualAccountResponse, error)
	GetInvoiceByID(ctx context.Context, in *GetInvoiceByIDRequest, opts ...grpc.CallOption) (*GetInvoiceByIDResponse, error)
	ProcessWebhook(ctx context.Context, in *ProcessWebhookRequest, opts ...grpc.CallOption) (*ProcessWebhookResponse, error)
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

func (c *paymentServiceClient) IssueVirtualAccount(ctx context.Context, in *IssueVirtualAccountRequest, opts ...grpc.CallOption) (*IssueVirtualAccountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IssueVirtualAccountResponse)
	err := c.cc.Invoke(ctx, PaymentService_IssueVA_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) GetInvoiceByID(ctx context.Context, in *GetInvoiceByIDRequest, opts ...grpc.CallOption) (*GetInvoiceByIDResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetInvoiceByIDResponse)
	err := c.cc.Invoke(ctx, PaymentService_GetInvoiceByID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentServiceClient) ProcessWebhook(ctx context.Context, in *ProcessWebhookRequest, opts ...grpc.CallOption) (*ProcessWebhookResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProcessWebhookResponse)
	err := c.cc.Invoke(ctx, PaymentService_ProcessWebhook_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
