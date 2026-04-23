// finance_disbursement_stub.go — hand-written gRPC client stubs for finance-svc
// AP disbursement and AR/AP aging RPCs (BL-FIN-010/011), gateway-svc consumer side.

package pb

import (
	"context"

	"google.golang.org/grpc"
)

const (
	FinanceService_CreateDisbursementBatch_FullMethodName = "/pb.FinanceService/CreateDisbursementBatch"
	FinanceService_ApproveDisbursement_FullMethodName     = "/pb.FinanceService/ApproveDisbursement"
	FinanceService_GetARAPAging_FullMethodName            = "/pb.FinanceService/GetARAPAging"
)

// ---------------------------------------------------------------------------
// Message types — mirror finance-svc pb/disbursement_messages.go
// ---------------------------------------------------------------------------

type DisbursementItemInputPb struct {
	VendorName  string
	Description string
	AmountIdr   int64
	Reference   string
}

func (x *DisbursementItemInputPb) GetVendorName() string {
	if x == nil {
		return ""
	}
	return x.VendorName
}
func (x *DisbursementItemInputPb) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *DisbursementItemInputPb) GetAmountIdr() int64 {
	if x == nil {
		return 0
	}
	return x.AmountIdr
}
func (x *DisbursementItemInputPb) GetReference() string {
	if x == nil {
		return ""
	}
	return x.Reference
}

type CreateDisbursementBatchRequest struct {
	Description string
	Items       []*DisbursementItemInputPb
	CreatedBy   string
}

func (x *CreateDisbursementBatchRequest) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *CreateDisbursementBatchRequest) GetItems() []*DisbursementItemInputPb {
	if x == nil {
		return nil
	}
	return x.Items
}
func (x *CreateDisbursementBatchRequest) GetCreatedBy() string {
	if x == nil {
		return ""
	}
	return x.CreatedBy
}

type CreateDisbursementBatchResponse struct {
	BatchId        string
	TotalAmountIdr int64
	ItemCount      int32
	Status         string
}

func (x *CreateDisbursementBatchResponse) GetBatchId() string {
	if x == nil {
		return ""
	}
	return x.BatchId
}
func (x *CreateDisbursementBatchResponse) GetTotalAmountIdr() int64 {
	if x == nil {
		return 0
	}
	return x.TotalAmountIdr
}
func (x *CreateDisbursementBatchResponse) GetItemCount() int32 {
	if x == nil {
		return 0
	}
	return x.ItemCount
}
func (x *CreateDisbursementBatchResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

type ApproveDisbursementRequest struct {
	BatchId    string
	ApprovedBy string
	Approved   bool
	Notes      string
}

func (x *ApproveDisbursementRequest) GetBatchId() string {
	if x == nil {
		return ""
	}
	return x.BatchId
}
func (x *ApproveDisbursementRequest) GetApprovedBy() string {
	if x == nil {
		return ""
	}
	return x.ApprovedBy
}
func (x *ApproveDisbursementRequest) GetApproved() bool {
	if x == nil {
		return false
	}
	return x.Approved
}
func (x *ApproveDisbursementRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type ApproveDisbursementResponse struct {
	BatchId         string
	Status          string
	JournalEntryIds []string
}

func (x *ApproveDisbursementResponse) GetBatchId() string {
	if x == nil {
		return ""
	}
	return x.BatchId
}
func (x *ApproveDisbursementResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *ApproveDisbursementResponse) GetJournalEntryIds() []string {
	if x == nil {
		return nil
	}
	return x.JournalEntryIds
}

type GetARAPAgingRequest struct {
	Type     string // "AR" | "AP" | "both"
	AsOfDate string // "YYYY-MM-DD"; empty = today
}

func (x *GetARAPAgingRequest) GetType() string {
	if x == nil {
		return ""
	}
	return x.Type
}
func (x *GetARAPAgingRequest) GetAsOfDate() string {
	if x == nil {
		return ""
	}
	return x.AsOfDate
}

type AgingBucketsPb struct {
	Current int64
	Days30  int64
	Days60  int64
	Days90  int64
	Over90  int64
	Total   int64
}

func (x *AgingBucketsPb) GetCurrent() int64 {
	if x == nil {
		return 0
	}
	return x.Current
}
func (x *AgingBucketsPb) GetDays30() int64 {
	if x == nil {
		return 0
	}
	return x.Days30
}
func (x *AgingBucketsPb) GetDays60() int64 {
	if x == nil {
		return 0
	}
	return x.Days60
}
func (x *AgingBucketsPb) GetDays90() int64 {
	if x == nil {
		return 0
	}
	return x.Days90
}
func (x *AgingBucketsPb) GetOver90() int64 {
	if x == nil {
		return 0
	}
	return x.Over90
}
func (x *AgingBucketsPb) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}

type GetARAPAgingResponse struct {
	Ar          *AgingBucketsPb
	Ap          *AgingBucketsPb
	GeneratedAt string
}

func (x *GetARAPAgingResponse) GetAr() *AgingBucketsPb {
	if x == nil {
		return nil
	}
	return x.Ar
}
func (x *GetARAPAgingResponse) GetAp() *AgingBucketsPb {
	if x == nil {
		return nil
	}
	return x.Ap
}
func (x *GetARAPAgingResponse) GetGeneratedAt() string {
	if x == nil {
		return ""
	}
	return x.GeneratedAt
}

// ---------------------------------------------------------------------------
// Client interface + implementation
// ---------------------------------------------------------------------------

// FinanceDisbursementClient is the narrow interface used by the gateway adapter.
type FinanceDisbursementClient interface {
	CreateDisbursementBatch(ctx context.Context, in *CreateDisbursementBatchRequest, opts ...grpc.CallOption) (*CreateDisbursementBatchResponse, error)
	ApproveDisbursement(ctx context.Context, in *ApproveDisbursementRequest, opts ...grpc.CallOption) (*ApproveDisbursementResponse, error)
	GetARAPAging(ctx context.Context, in *GetARAPAgingRequest, opts ...grpc.CallOption) (*GetARAPAgingResponse, error)
}

type financeDisbursementClient struct {
	cc grpc.ClientConnInterface
}

func NewFinanceDisbursementClient(cc grpc.ClientConnInterface) FinanceDisbursementClient {
	return &financeDisbursementClient{cc}
}

func (c *financeDisbursementClient) CreateDisbursementBatch(ctx context.Context, in *CreateDisbursementBatchRequest, opts ...grpc.CallOption) (*CreateDisbursementBatchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateDisbursementBatchResponse)
	if err := c.cc.Invoke(ctx, FinanceService_CreateDisbursementBatch_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *financeDisbursementClient) ApproveDisbursement(ctx context.Context, in *ApproveDisbursementRequest, opts ...grpc.CallOption) (*ApproveDisbursementResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApproveDisbursementResponse)
	if err := c.cc.Invoke(ctx, FinanceService_ApproveDisbursement_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *financeDisbursementClient) GetARAPAging(ctx context.Context, in *GetARAPAgingRequest, opts ...grpc.CallOption) (*GetARAPAgingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetARAPAgingResponse)
	if err := c.cc.Invoke(ctx, FinanceService_GetARAPAging_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
