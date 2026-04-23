// finance_pl_stub.go — gateway-side gRPC client stub for finance-svc depth
// RPCs: RecognizeRevenue, GetPLReport, GetBalanceSheet (Wave 1B / Phase 6).
//
// Mirrors services/finance-svc/api/grpc_api/pb/pl_messages.go and
// pl_grpc_ext.go. Run `make genpb` to replace with generated code.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

const (
	FinanceService_RecognizeRevenue_FullMethodName = "/pb.FinanceService/RecognizeRevenue"
	FinanceService_GetPLReport_FullMethodName      = "/pb.FinanceService/GetPLReport"
	FinanceService_GetBalanceSheet_FullMethodName  = "/pb.FinanceService/GetBalanceSheet"
)

// ---------------------------------------------------------------------------
// RecognizeRevenue
// ---------------------------------------------------------------------------

type RecognizeRevenueRequest struct {
	DepartureId    string
	TotalAmountIdr int64
}

func (x *RecognizeRevenueRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *RecognizeRevenueRequest) GetTotalAmountIdr() int64 {
	if x == nil {
		return 0
	}
	return x.TotalAmountIdr
}

type RecognizeRevenueResponse struct {
	EntryId      string
	RecognizedAt string
	Replayed     bool
}

func (x *RecognizeRevenueResponse) GetEntryId() string {
	if x == nil {
		return ""
	}
	return x.EntryId
}
func (x *RecognizeRevenueResponse) GetRecognizedAt() string {
	if x == nil {
		return ""
	}
	return x.RecognizedAt
}
func (x *RecognizeRevenueResponse) GetReplayed() bool {
	if x == nil {
		return false
	}
	return x.Replayed
}

// ---------------------------------------------------------------------------
// GetPLReport
// ---------------------------------------------------------------------------

type GetPLReportRequest struct {
	From string
	To   string
}

func (x *GetPLReportRequest) GetFrom() string {
	if x == nil {
		return ""
	}
	return x.From
}
func (x *GetPLReportRequest) GetTo() string {
	if x == nil {
		return ""
	}
	return x.To
}

type PLLineItemProto struct {
	AccountCode string
	AccountName string
	Amount      int64
	Direction   string
}

func (x *PLLineItemProto) GetAccountCode() string {
	if x == nil {
		return ""
	}
	return x.AccountCode
}
func (x *PLLineItemProto) GetAccountName() string {
	if x == nil {
		return ""
	}
	return x.AccountName
}
func (x *PLLineItemProto) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *PLLineItemProto) GetDirection() string {
	if x == nil {
		return ""
	}
	return x.Direction
}

type PLReportProto struct {
	PeriodFrom   string
	PeriodTo     string
	GeneratedAt  string
	TotalRevenue int64
	TotalCogs    int64
	GrossProfit  int64
	OtherIncome  int64
	OtherExpense int64
	NetProfit    int64
	Entries      []*PLLineItemProto
}

func (x *PLReportProto) GetPeriodFrom() string {
	if x == nil {
		return ""
	}
	return x.PeriodFrom
}
func (x *PLReportProto) GetPeriodTo() string {
	if x == nil {
		return ""
	}
	return x.PeriodTo
}
func (x *PLReportProto) GetGeneratedAt() string {
	if x == nil {
		return ""
	}
	return x.GeneratedAt
}
func (x *PLReportProto) GetTotalRevenue() int64 {
	if x == nil {
		return 0
	}
	return x.TotalRevenue
}
func (x *PLReportProto) GetTotalCogs() int64 {
	if x == nil {
		return 0
	}
	return x.TotalCogs
}
func (x *PLReportProto) GetGrossProfit() int64 {
	if x == nil {
		return 0
	}
	return x.GrossProfit
}
func (x *PLReportProto) GetOtherIncome() int64 {
	if x == nil {
		return 0
	}
	return x.OtherIncome
}
func (x *PLReportProto) GetOtherExpense() int64 {
	if x == nil {
		return 0
	}
	return x.OtherExpense
}
func (x *PLReportProto) GetNetProfit() int64 {
	if x == nil {
		return 0
	}
	return x.NetProfit
}
func (x *PLReportProto) GetEntries() []*PLLineItemProto {
	if x == nil {
		return nil
	}
	return x.Entries
}

// ---------------------------------------------------------------------------
// GetBalanceSheet
// ---------------------------------------------------------------------------

type GetBalanceSheetRequest struct {
	AsOf string
}

func (x *GetBalanceSheetRequest) GetAsOf() string {
	if x == nil {
		return ""
	}
	return x.AsOf
}

type BalanceSheetLineProto struct {
	AccountCode string
	AccountName string
	Balance     int64
}

func (x *BalanceSheetLineProto) GetAccountCode() string {
	if x == nil {
		return ""
	}
	return x.AccountCode
}
func (x *BalanceSheetLineProto) GetAccountName() string {
	if x == nil {
		return ""
	}
	return x.AccountName
}
func (x *BalanceSheetLineProto) GetBalance() int64 {
	if x == nil {
		return 0
	}
	return x.Balance
}

type BalanceSheetProto struct {
	AsOfDate         string
	GeneratedAt      string
	Assets           []*BalanceSheetLineProto
	Liabilities      []*BalanceSheetLineProto
	Equity           []*BalanceSheetLineProto
	TotalAssets      int64
	TotalLiabilities int64
	TotalEquity      int64
}

func (x *BalanceSheetProto) GetAsOfDate() string {
	if x == nil {
		return ""
	}
	return x.AsOfDate
}
func (x *BalanceSheetProto) GetGeneratedAt() string {
	if x == nil {
		return ""
	}
	return x.GeneratedAt
}
func (x *BalanceSheetProto) GetAssets() []*BalanceSheetLineProto {
	if x == nil {
		return nil
	}
	return x.Assets
}
func (x *BalanceSheetProto) GetLiabilities() []*BalanceSheetLineProto {
	if x == nil {
		return nil
	}
	return x.Liabilities
}
func (x *BalanceSheetProto) GetEquity() []*BalanceSheetLineProto {
	if x == nil {
		return nil
	}
	return x.Equity
}
func (x *BalanceSheetProto) GetTotalAssets() int64 {
	if x == nil {
		return 0
	}
	return x.TotalAssets
}
func (x *BalanceSheetProto) GetTotalLiabilities() int64 {
	if x == nil {
		return 0
	}
	return x.TotalLiabilities
}
func (x *BalanceSheetProto) GetTotalEquity() int64 {
	if x == nil {
		return 0
	}
	return x.TotalEquity
}

// ---------------------------------------------------------------------------
// FinanceDepthClient interface + implementation
// ---------------------------------------------------------------------------

// FinanceDepthClient is the consumer-side interface for finance-depth RPCs.
type FinanceDepthClient interface {
	RecognizeRevenue(ctx context.Context, in *RecognizeRevenueRequest, opts ...grpc.CallOption) (*RecognizeRevenueResponse, error)
	GetPLReport(ctx context.Context, in *GetPLReportRequest, opts ...grpc.CallOption) (*PLReportProto, error)
	GetBalanceSheet(ctx context.Context, in *GetBalanceSheetRequest, opts ...grpc.CallOption) (*BalanceSheetProto, error)
}

type financeDepthClient struct {
	cc grpc.ClientConnInterface
}

// NewFinanceDepthClient wraps a conn so gateway-svc can call finance depth RPCs.
func NewFinanceDepthClient(cc grpc.ClientConnInterface) FinanceDepthClient {
	return &financeDepthClient{cc}
}

func (c *financeDepthClient) RecognizeRevenue(ctx context.Context, in *RecognizeRevenueRequest, opts ...grpc.CallOption) (*RecognizeRevenueResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecognizeRevenueResponse)
	err := c.cc.Invoke(ctx, FinanceService_RecognizeRevenue_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *financeDepthClient) GetPLReport(ctx context.Context, in *GetPLReportRequest, opts ...grpc.CallOption) (*PLReportProto, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PLReportProto)
	err := c.cc.Invoke(ctx, FinanceService_GetPLReport_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *financeDepthClient) GetBalanceSheet(ctx context.Context, in *GetBalanceSheetRequest, opts ...grpc.CallOption) (*BalanceSheetProto, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BalanceSheetProto)
	err := c.cc.Invoke(ctx, FinanceService_GetBalanceSheet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
