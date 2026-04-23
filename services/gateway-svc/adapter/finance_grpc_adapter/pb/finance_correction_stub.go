// finance_correction_stub.go — hand-written gRPC client stub for finance-svc
// CorrectJournal RPC (BL-FIN-006), gateway-svc consumer side.

package pb

import (
	"context"

	"google.golang.org/grpc"
)

const (
	FinanceService_CorrectJournal_FullMethodName = "/pb.FinanceService/CorrectJournal"
)

// CorrectJournalRequest mirrors finance-svc pb.CorrectJournalRequest.
type CorrectJournalRequest struct {
	EntryId     string
	Reason      string
	ActorUserId string
}

func (x *CorrectJournalRequest) GetEntryId() string {
	if x == nil {
		return ""
	}
	return x.EntryId
}
func (x *CorrectJournalRequest) GetReason() string {
	if x == nil {
		return ""
	}
	return x.Reason
}
func (x *CorrectJournalRequest) GetActorUserId() string {
	if x == nil {
		return ""
	}
	return x.ActorUserId
}

// CorrectJournalResponse mirrors finance-svc pb.CorrectJournalResponse.
type CorrectJournalResponse struct {
	CorrectionEntryId string
	OriginalEntryId   string
	Idempotent        bool
}

func (x *CorrectJournalResponse) GetCorrectionEntryId() string {
	if x == nil {
		return ""
	}
	return x.CorrectionEntryId
}
func (x *CorrectJournalResponse) GetOriginalEntryId() string {
	if x == nil {
		return ""
	}
	return x.OriginalEntryId
}
func (x *CorrectJournalResponse) GetIdempotent() bool {
	if x == nil {
		return false
	}
	return x.Idempotent
}

// FinanceCorrectionClient is the narrow interface used by the gateway adapter.
type FinanceCorrectionClient interface {
	CorrectJournal(ctx context.Context, in *CorrectJournalRequest, opts ...grpc.CallOption) (*CorrectJournalResponse, error)
}

type financeCorrectionClient struct {
	cc grpc.ClientConnInterface
}

func NewFinanceCorrectionClient(cc grpc.ClientConnInterface) FinanceCorrectionClient {
	return &financeCorrectionClient{cc}
}

func (c *financeCorrectionClient) CorrectJournal(ctx context.Context, in *CorrectJournalRequest, opts ...grpc.CallOption) (*CorrectJournalResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CorrectJournalResponse)
	err := c.cc.Invoke(ctx, FinanceService_CorrectJournal_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
