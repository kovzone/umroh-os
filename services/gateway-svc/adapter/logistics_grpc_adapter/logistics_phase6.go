// logistics_phase6.go — gateway-svc adapter methods for logistics-svc Phase 6 RPCs
// (BL-LOG-010..012): CreatePurchaseRequest, ApprovePurchaseRequest,
// RecordGRNWithQC, CreateKitAssembly.

package logistics_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/logistics_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// CreatePurchaseRequest
// ---------------------------------------------------------------------------

// CreatePurchaseRequestParams holds inputs for POST /v1/logistics/purchase-requests.
type CreatePurchaseRequestParams struct {
	DepartureID    string
	RequestedBy    string
	ItemName       string
	Quantity       int32
	UnitPriceIdr   int64
	BudgetLimitIdr int64
}

// CreatePurchaseRequestResult is the response.
type CreatePurchaseRequestResult struct {
	PrID          string
	Status        string
	TotalPriceIdr int64
}

func (a *Adapter) CreatePurchaseRequest(ctx context.Context, params *CreatePurchaseRequestParams) (*CreatePurchaseRequestResult, error) {
	const op = "logistics_grpc_adapter.Adapter.CreatePurchaseRequest"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("departure_id", params.DepartureID).Msg("")

	resp, err := a.logisticsPhase6Client.CreatePurchaseRequest(ctx, &pb.CreatePurchaseRequestRequest{
		DepartureId:    params.DepartureID,
		RequestedBy:    params.RequestedBy,
		ItemName:       params.ItemName,
		Quantity:       params.Quantity,
		UnitPriceIdr:   params.UnitPriceIdr,
		BudgetLimitIdr: params.BudgetLimitIdr,
	})
	if err != nil {
		wrapped := mapLogisticsError(err)
		logger.Warn().Err(wrapped).Msg("logistics-svc.CreatePurchaseRequest failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &CreatePurchaseRequestResult{
		PrID:          resp.GetPrId(),
		Status:        resp.GetStatus(),
		TotalPriceIdr: resp.GetTotalPriceIdr(),
	}, nil
}

// ---------------------------------------------------------------------------
// ApprovePurchaseRequest
// ---------------------------------------------------------------------------

// ApprovePurchaseRequestParams holds inputs for PUT /v1/logistics/purchase-requests/:id/decision.
type ApprovePurchaseRequestParams struct {
	PrID       string
	ApprovedBy string
	Approved   bool
	Notes      string
}

// ApprovePurchaseRequestResult is the response.
type ApprovePurchaseRequestResult struct {
	PrID      string
	NewStatus string
}

func (a *Adapter) ApprovePurchaseRequest(ctx context.Context, params *ApprovePurchaseRequestParams) (*ApprovePurchaseRequestResult, error) {
	const op = "logistics_grpc_adapter.Adapter.ApprovePurchaseRequest"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("pr_id", params.PrID).Msg("")

	resp, err := a.logisticsPhase6Client.ApprovePurchaseRequest(ctx, &pb.ApprovePurchaseRequestRequest{
		PrId:       params.PrID,
		ApprovedBy: params.ApprovedBy,
		Approved:   params.Approved,
		Notes:      params.Notes,
	})
	if err != nil {
		wrapped := mapLogisticsError(err)
		logger.Warn().Err(wrapped).Msg("logistics-svc.ApprovePurchaseRequest failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &ApprovePurchaseRequestResult{
		PrID:      resp.GetPrId(),
		NewStatus: resp.GetNewStatus(),
	}, nil
}

// ---------------------------------------------------------------------------
// RecordGRNWithQC
// ---------------------------------------------------------------------------

// RecordGRNWithQCParams holds inputs for POST /v1/logistics/grn (Phase 6 extension).
type RecordGRNWithQCParams struct {
	GrnID       string
	DepartureID string
	AmountIdr   int64
	QcPassed    bool
	QcNotes     string
}

// RecordGRNWithQCResult is the response.
type RecordGRNWithQCResult struct {
	GrnID         string
	QcStatus      string
	JournalPosted bool
}

func (a *Adapter) RecordGRNWithQC(ctx context.Context, params *RecordGRNWithQCParams) (*RecordGRNWithQCResult, error) {
	const op = "logistics_grpc_adapter.Adapter.RecordGRNWithQC"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("grn_id", params.GrnID).Msg("")

	resp, err := a.logisticsPhase6Client.RecordGRNWithQC(ctx, &pb.RecordGRNWithQCRequest{
		GrnId:       params.GrnID,
		DepartureId: params.DepartureID,
		AmountIdr:   params.AmountIdr,
		QcPassed:    params.QcPassed,
		QcNotes:     params.QcNotes,
	})
	if err != nil {
		wrapped := mapLogisticsError(err)
		logger.Warn().Err(wrapped).Msg("logistics-svc.RecordGRNWithQC failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &RecordGRNWithQCResult{
		GrnID:         resp.GetGrnId(),
		QcStatus:      resp.GetQcStatus(),
		JournalPosted: resp.GetJournalPosted(),
	}, nil
}

// ---------------------------------------------------------------------------
// CreateKitAssembly
// ---------------------------------------------------------------------------

// KitItem holds one kit item line.
type KitItem struct {
	ItemName string
	Quantity int32
}

// CreateKitAssemblyParams holds inputs for POST /v1/logistics/kit-assembly.
type CreateKitAssemblyParams struct {
	DepartureID    string
	AssembledBy    string
	Items          []*KitItem
	IdempotencyKey string
}

// CreateKitAssemblyResult is the response.
type CreateKitAssemblyResult struct {
	AssemblyID string
	Status     string
	Idempotent bool
}

func (a *Adapter) CreateKitAssembly(ctx context.Context, params *CreateKitAssemblyParams) (*CreateKitAssemblyResult, error) {
	const op = "logistics_grpc_adapter.Adapter.CreateKitAssembly"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("departure_id", params.DepartureID).Msg("")

	items := make([]*pb.KitItemPb, 0, len(params.Items))
	for _, it := range params.Items {
		items = append(items, &pb.KitItemPb{
			ItemName: it.ItemName,
			Quantity: it.Quantity,
		})
	}

	resp, err := a.logisticsPhase6Client.CreateKitAssembly(ctx, &pb.CreateKitAssemblyRequest{
		DepartureId:    params.DepartureID,
		AssembledBy:    params.AssembledBy,
		Items:          items,
		IdempotencyKey: params.IdempotencyKey,
	})
	if err != nil {
		wrapped := mapLogisticsError(err)
		logger.Warn().Err(wrapped).Msg("logistics-svc.CreateKitAssembly failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &CreateKitAssemblyResult{
		AssemblyID: resp.GetAssemblyId(),
		Status:     resp.GetStatus(),
		Idempotent: resp.GetIdempotent(),
	}, nil
}
