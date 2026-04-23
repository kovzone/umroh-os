// logistics_phase6_stub.go — hand-written gRPC client stubs for logistics-svc
// Phase 6 RPCs (BL-LOG-010..012): CreatePurchaseRequest, ApprovePurchaseRequest,
// RecordGRNWithQC, CreateKitAssembly.

package pb

import (
	"context"

	"google.golang.org/grpc"
)

const (
	LogisticsService_CreatePurchaseRequest_FullMethodName  = "/pb.LogisticsService/CreatePurchaseRequest"
	LogisticsService_ApprovePurchaseRequest_FullMethodName = "/pb.LogisticsService/ApprovePurchaseRequest"
	LogisticsService_RecordGRNWithQC_FullMethodName        = "/pb.LogisticsService/RecordGRNWithQC"
	LogisticsService_CreateKitAssembly_FullMethodName      = "/pb.LogisticsService/CreateKitAssembly"
)

// ---------------------------------------------------------------------------
// Message types — mirror logistics-svc pb/logistics_procurement_messages.go
// ---------------------------------------------------------------------------

type CreatePurchaseRequestRequest struct {
	DepartureId    string
	RequestedBy    string
	ItemName       string
	Quantity       int32
	UnitPriceIdr   int64
	BudgetLimitIdr int64
}

func (x *CreatePurchaseRequestRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *CreatePurchaseRequestRequest) GetRequestedBy() string {
	if x == nil {
		return ""
	}
	return x.RequestedBy
}
func (x *CreatePurchaseRequestRequest) GetItemName() string {
	if x == nil {
		return ""
	}
	return x.ItemName
}
func (x *CreatePurchaseRequestRequest) GetQuantity() int32 {
	if x == nil {
		return 0
	}
	return x.Quantity
}
func (x *CreatePurchaseRequestRequest) GetUnitPriceIdr() int64 {
	if x == nil {
		return 0
	}
	return x.UnitPriceIdr
}
func (x *CreatePurchaseRequestRequest) GetBudgetLimitIdr() int64 {
	if x == nil {
		return 0
	}
	return x.BudgetLimitIdr
}

type CreatePurchaseRequestResponse struct {
	PrId          string
	Status        string
	TotalPriceIdr int64
}

func (x *CreatePurchaseRequestResponse) GetPrId() string {
	if x == nil {
		return ""
	}
	return x.PrId
}
func (x *CreatePurchaseRequestResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *CreatePurchaseRequestResponse) GetTotalPriceIdr() int64 {
	if x == nil {
		return 0
	}
	return x.TotalPriceIdr
}

type ApprovePurchaseRequestRequest struct {
	PrId       string
	ApprovedBy string
	Approved   bool
	Notes      string
}

func (x *ApprovePurchaseRequestRequest) GetPrId() string {
	if x == nil {
		return ""
	}
	return x.PrId
}
func (x *ApprovePurchaseRequestRequest) GetApprovedBy() string {
	if x == nil {
		return ""
	}
	return x.ApprovedBy
}
func (x *ApprovePurchaseRequestRequest) GetApproved() bool {
	if x == nil {
		return false
	}
	return x.Approved
}
func (x *ApprovePurchaseRequestRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type ApprovePurchaseRequestResponse struct {
	PrId      string
	NewStatus string
}

func (x *ApprovePurchaseRequestResponse) GetPrId() string {
	if x == nil {
		return ""
	}
	return x.PrId
}
func (x *ApprovePurchaseRequestResponse) GetNewStatus() string {
	if x == nil {
		return ""
	}
	return x.NewStatus
}

type RecordGRNWithQCRequest struct {
	GrnId       string
	DepartureId string
	AmountIdr   int64
	QcPassed    bool
	QcNotes     string
}

func (x *RecordGRNWithQCRequest) GetGrnId() string {
	if x == nil {
		return ""
	}
	return x.GrnId
}
func (x *RecordGRNWithQCRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *RecordGRNWithQCRequest) GetAmountIdr() int64 {
	if x == nil {
		return 0
	}
	return x.AmountIdr
}
func (x *RecordGRNWithQCRequest) GetQcPassed() bool {
	if x == nil {
		return false
	}
	return x.QcPassed
}
func (x *RecordGRNWithQCRequest) GetQcNotes() string {
	if x == nil {
		return ""
	}
	return x.QcNotes
}

type RecordGRNWithQCResponse struct {
	GrnId         string
	QcStatus      string
	JournalPosted bool
}

func (x *RecordGRNWithQCResponse) GetGrnId() string {
	if x == nil {
		return ""
	}
	return x.GrnId
}
func (x *RecordGRNWithQCResponse) GetQcStatus() string {
	if x == nil {
		return ""
	}
	return x.QcStatus
}
func (x *RecordGRNWithQCResponse) GetJournalPosted() bool {
	if x == nil {
		return false
	}
	return x.JournalPosted
}

type KitItemPb struct {
	ItemName string
	Quantity int32
}

func (x *KitItemPb) GetItemName() string {
	if x == nil {
		return ""
	}
	return x.ItemName
}
func (x *KitItemPb) GetQuantity() int32 {
	if x == nil {
		return 0
	}
	return x.Quantity
}

type CreateKitAssemblyRequest struct {
	DepartureId    string
	AssembledBy    string
	Items          []*KitItemPb
	IdempotencyKey string
}

func (x *CreateKitAssemblyRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *CreateKitAssemblyRequest) GetAssembledBy() string {
	if x == nil {
		return ""
	}
	return x.AssembledBy
}
func (x *CreateKitAssemblyRequest) GetItems() []*KitItemPb {
	if x == nil {
		return nil
	}
	return x.Items
}
func (x *CreateKitAssemblyRequest) GetIdempotencyKey() string {
	if x == nil {
		return ""
	}
	return x.IdempotencyKey
}

type CreateKitAssemblyResponse struct {
	AssemblyId string
	Status     string
	Idempotent bool
}

func (x *CreateKitAssemblyResponse) GetAssemblyId() string {
	if x == nil {
		return ""
	}
	return x.AssemblyId
}
func (x *CreateKitAssemblyResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *CreateKitAssemblyResponse) GetIdempotent() bool {
	if x == nil {
		return false
	}
	return x.Idempotent
}

// ---------------------------------------------------------------------------
// Client interface + implementation
// ---------------------------------------------------------------------------

// LogisticsPhase6Client is the narrow interface for Phase 6 logistics RPCs.
type LogisticsPhase6Client interface {
	CreatePurchaseRequest(ctx context.Context, in *CreatePurchaseRequestRequest, opts ...grpc.CallOption) (*CreatePurchaseRequestResponse, error)
	ApprovePurchaseRequest(ctx context.Context, in *ApprovePurchaseRequestRequest, opts ...grpc.CallOption) (*ApprovePurchaseRequestResponse, error)
	RecordGRNWithQC(ctx context.Context, in *RecordGRNWithQCRequest, opts ...grpc.CallOption) (*RecordGRNWithQCResponse, error)
	CreateKitAssembly(ctx context.Context, in *CreateKitAssemblyRequest, opts ...grpc.CallOption) (*CreateKitAssemblyResponse, error)
}

type logisticsPhase6Client struct {
	cc grpc.ClientConnInterface
}

func NewLogisticsPhase6Client(cc grpc.ClientConnInterface) LogisticsPhase6Client {
	return &logisticsPhase6Client{cc}
}

func (c *logisticsPhase6Client) CreatePurchaseRequest(ctx context.Context, in *CreatePurchaseRequestRequest, opts ...grpc.CallOption) (*CreatePurchaseRequestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreatePurchaseRequestResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_CreatePurchaseRequest_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticsPhase6Client) ApprovePurchaseRequest(ctx context.Context, in *ApprovePurchaseRequestRequest, opts ...grpc.CallOption) (*ApprovePurchaseRequestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApprovePurchaseRequestResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_ApprovePurchaseRequest_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticsPhase6Client) RecordGRNWithQC(ctx context.Context, in *RecordGRNWithQCRequest, opts ...grpc.CallOption) (*RecordGRNWithQCResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecordGRNWithQCResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_RecordGRNWithQC_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticsPhase6Client) CreateKitAssembly(ctx context.Context, in *CreateKitAssemblyRequest, opts ...grpc.CallOption) (*CreateKitAssemblyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateKitAssemblyResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_CreateKitAssembly_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
