// disbursement.go — gateway-svc adapter methods for finance-svc AP disbursement
// and AR/AP aging RPCs (BL-FIN-010/011).

package finance_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/finance_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Extend Adapter to wire disbursementClient
// ---------------------------------------------------------------------------

// disbursementClient is added to Adapter in NewAdapterWithDisbursement.
// We extend via a separate unexported field to avoid changing NewAdapter's signature
// (which already wires 4 clients from the same conn).
// The actual client wiring is done by patching NewAdapter — see init below.

// ---------------------------------------------------------------------------
// CreateDisbursementBatch
// ---------------------------------------------------------------------------

// DisbursementItemInput is one line in a disbursement batch.
type DisbursementItemInput struct {
	VendorName  string
	Description string
	AmountIdr   int64
	Reference   string
}

// CreateDisbursementBatchParams holds inputs for POST /v1/finance/disbursements.
type CreateDisbursementBatchParams struct {
	Description string
	Items       []*DisbursementItemInput
	CreatedBy   string
}

// CreateDisbursementBatchResult is the response for a new disbursement batch.
type CreateDisbursementBatchResult struct {
	BatchID        string
	TotalAmountIdr int64
	ItemCount      int32
	Status         string
}

func (a *Adapter) CreateDisbursementBatch(ctx context.Context, params *CreateDisbursementBatchParams) (*CreateDisbursementBatchResult, error) {
	const op = "finance_grpc_adapter.Adapter.CreateDisbursementBatch"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	items := make([]*pb.DisbursementItemInputPb, 0, len(params.Items))
	for _, it := range params.Items {
		items = append(items, &pb.DisbursementItemInputPb{
			VendorName:  it.VendorName,
			Description: it.Description,
			AmountIdr:   it.AmountIdr,
			Reference:   it.Reference,
		})
	}

	resp, err := a.disbursementClient.CreateDisbursementBatch(ctx, &pb.CreateDisbursementBatchRequest{
		Description: params.Description,
		Items:       items,
		CreatedBy:   params.CreatedBy,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.CreateDisbursementBatch failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &CreateDisbursementBatchResult{
		BatchID:        resp.GetBatchId(),
		TotalAmountIdr: resp.GetTotalAmountIdr(),
		ItemCount:      resp.GetItemCount(),
		Status:         resp.GetStatus(),
	}, nil
}

// ---------------------------------------------------------------------------
// ApproveDisbursement
// ---------------------------------------------------------------------------

// ApproveDisbursementParams holds inputs for PUT /v1/finance/disbursements/:id/decision.
type ApproveDisbursementParams struct {
	BatchID    string
	ApprovedBy string
	Approved   bool
	Notes      string
}

// ApproveDisbursementResult is the response after approve/reject.
type ApproveDisbursementResult struct {
	BatchID         string
	Status          string
	JournalEntryIDs []string
}

func (a *Adapter) ApproveDisbursement(ctx context.Context, params *ApproveDisbursementParams) (*ApproveDisbursementResult, error) {
	const op = "finance_grpc_adapter.Adapter.ApproveDisbursement"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.disbursementClient.ApproveDisbursement(ctx, &pb.ApproveDisbursementRequest{
		BatchId:    params.BatchID,
		ApprovedBy: params.ApprovedBy,
		Approved:   params.Approved,
		Notes:      params.Notes,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.ApproveDisbursement failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &ApproveDisbursementResult{
		BatchID:         resp.GetBatchId(),
		Status:          resp.GetStatus(),
		JournalEntryIDs: resp.GetJournalEntryIds(),
	}, nil
}

// ---------------------------------------------------------------------------
// GetARAPAging
// ---------------------------------------------------------------------------

// GetARAPAgingParams holds inputs for GET /v1/finance/aging.
type GetARAPAgingParams struct {
	Type     string // "AR" | "AP" | "both"
	AsOfDate string // "YYYY-MM-DD"
}

// AgingBuckets holds one side's aging breakdown.
type AgingBuckets struct {
	Current int64
	Days30  int64
	Days60  int64
	Days90  int64
	Over90  int64
	Total   int64
}

// GetARAPAgingResult is the response for GET /v1/finance/aging.
type GetARAPAgingResult struct {
	AR          *AgingBuckets
	AP          *AgingBuckets
	GeneratedAt string
}

func (a *Adapter) GetARAPAging(ctx context.Context, params *GetARAPAgingParams) (*GetARAPAgingResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetARAPAging"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.disbursementClient.GetARAPAging(ctx, &pb.GetARAPAgingRequest{
		Type:     params.Type,
		AsOfDate: params.AsOfDate,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetARAPAging failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	mapBuckets := func(b *pb.AgingBucketsPb) *AgingBuckets {
		if b == nil {
			return nil
		}
		return &AgingBuckets{
			Current: b.GetCurrent(),
			Days30:  b.GetDays30(),
			Days60:  b.GetDays60(),
			Days90:  b.GetDays90(),
			Over90:  b.GetOver90(),
			Total:   b.GetTotal(),
		}
	}

	span.SetStatus(codes.Ok, "ok")
	return &GetARAPAgingResult{
		AR:          mapBuckets(resp.GetAr()),
		AP:          mapBuckets(resp.GetAp()),
		GeneratedAt: resp.GetGeneratedAt(),
	}, nil
}
