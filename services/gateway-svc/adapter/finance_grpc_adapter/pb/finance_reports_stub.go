// Package pb — hand-written gRPC client stub for the gateway-svc → finance-svc
// finance report RPCs (S5-E-01 / ADR-0009).
//
// gateway-svc calls two RPCs on finance-svc:
//   - GetFinanceSummary  (GET /v1/finance/summary — bearer)
//   - ListJournalEntries (GET /v1/finance/journals — bearer)
//
// Run `make genpb` with a shared proto path to replace with generated code.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

const (
	FinanceService_GetFinanceSummary_FullMethodName  = "/pb.FinanceService/GetFinanceSummary"
	FinanceService_ListJournalEntries_FullMethodName = "/pb.FinanceService/ListJournalEntries"
)

// ---------------------------------------------------------------------------
// GetFinanceSummary
// ---------------------------------------------------------------------------

type GetFinanceSummaryRequest struct{}

type AccountBalance struct {
	AccountCode string
	DebitTotal  int64
	CreditTotal int64
	Net         int64
}

func (x *AccountBalance) GetAccountCode() string {
	if x == nil {
		return ""
	}
	return x.AccountCode
}
func (x *AccountBalance) GetDebitTotal() int64 {
	if x == nil {
		return 0
	}
	return x.DebitTotal
}
func (x *AccountBalance) GetCreditTotal() int64 {
	if x == nil {
		return 0
	}
	return x.CreditTotal
}
func (x *AccountBalance) GetNet() int64 {
	if x == nil {
		return 0
	}
	return x.Net
}

type GetFinanceSummaryResponse struct {
	Accounts []*AccountBalance
}

func (x *GetFinanceSummaryResponse) GetAccounts() []*AccountBalance {
	if x == nil {
		return nil
	}
	return x.Accounts
}

// ---------------------------------------------------------------------------
// ListJournalEntries
// ---------------------------------------------------------------------------

type ListJournalEntriesRequest struct {
	From   string
	To     string
	Limit  int32
	Cursor string
}

type JournalLineProto struct {
	Id          string
	EntryId     string
	AccountCode string
	Debit       int64
	Credit      int64
}

func (x *JournalLineProto) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *JournalLineProto) GetEntryId() string {
	if x == nil {
		return ""
	}
	return x.EntryId
}
func (x *JournalLineProto) GetAccountCode() string {
	if x == nil {
		return ""
	}
	return x.AccountCode
}
func (x *JournalLineProto) GetDebit() int64 {
	if x == nil {
		return 0
	}
	return x.Debit
}
func (x *JournalLineProto) GetCredit() int64 {
	if x == nil {
		return 0
	}
	return x.Credit
}

type JournalEntryProto struct {
	Id             string
	IdempotencyKey string
	SourceType     string
	SourceId       string
	PostedAt       string
	Description    string
	Lines          []*JournalLineProto
}

func (x *JournalEntryProto) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *JournalEntryProto) GetIdempotencyKey() string {
	if x == nil {
		return ""
	}
	return x.IdempotencyKey
}
func (x *JournalEntryProto) GetSourceType() string {
	if x == nil {
		return ""
	}
	return x.SourceType
}
func (x *JournalEntryProto) GetSourceId() string {
	if x == nil {
		return ""
	}
	return x.SourceId
}
func (x *JournalEntryProto) GetPostedAt() string {
	if x == nil {
		return ""
	}
	return x.PostedAt
}
func (x *JournalEntryProto) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *JournalEntryProto) GetLines() []*JournalLineProto {
	if x == nil {
		return nil
	}
	return x.Lines
}

type ListJournalEntriesResponse struct {
	Entries    []*JournalEntryProto
	NextCursor string
}

func (x *ListJournalEntriesResponse) GetEntries() []*JournalEntryProto {
	if x == nil {
		return nil
	}
	return x.Entries
}
func (x *ListJournalEntriesResponse) GetNextCursor() string {
	if x == nil {
		return ""
	}
	return x.NextCursor
}

// ---------------------------------------------------------------------------
// FinanceReportsClient interface + implementation
// ---------------------------------------------------------------------------

type FinanceReportsClient interface {
	GetFinanceSummary(ctx context.Context, in *GetFinanceSummaryRequest, opts ...grpc.CallOption) (*GetFinanceSummaryResponse, error)
	ListJournalEntries(ctx context.Context, in *ListJournalEntriesRequest, opts ...grpc.CallOption) (*ListJournalEntriesResponse, error)
}

type financeReportsClient struct {
	cc grpc.ClientConnInterface
}

func NewFinanceReportsClient(cc grpc.ClientConnInterface) FinanceReportsClient {
	return &financeReportsClient{cc}
}

func (c *financeReportsClient) GetFinanceSummary(ctx context.Context, in *GetFinanceSummaryRequest, opts ...grpc.CallOption) (*GetFinanceSummaryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFinanceSummaryResponse)
	err := c.cc.Invoke(ctx, FinanceService_GetFinanceSummary_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *financeReportsClient) ListJournalEntries(ctx context.Context, in *ListJournalEntriesRequest, opts ...grpc.CallOption) (*ListJournalEntriesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListJournalEntriesResponse)
	err := c.cc.Invoke(ctx, FinanceService_ListJournalEntries_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
