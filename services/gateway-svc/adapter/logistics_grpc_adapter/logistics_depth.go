// logistics_depth.go — gateway logistics adapter methods for Wave 5 depth RPCs.
// BL-LOG-013..029: purchase request list, budget sync, tiered approvals, auto vendor,
// partial/reverse GRN, barcode/labels, warehouse, stock transfer, stock alerts,
// reorder levels, stocktake lifecycle, size sync, courier tracking, returns, exchanges.

package logistics_grpc_adapter

import (
	"context"
	"errors"
	"fmt"

	"gateway-svc/adapter/logistics_grpc_adapter/pb"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// Error mapper
// ---------------------------------------------------------------------------

func mapLogisticsDepthError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("logistics-depth call failed: %w", err))
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return fmt.Errorf("%w: %s", apperrors.ErrNotFound, st.Message())
	case grpcCodes.InvalidArgument:
		return fmt.Errorf("%w: %s", apperrors.ErrValidation, st.Message())
	case grpcCodes.PermissionDenied:
		return fmt.Errorf("%w: %s", apperrors.ErrForbidden, st.Message())
	default:
		return fmt.Errorf("%w: logistics-depth %s: %s", apperrors.ErrInternal, st.Code(), st.Message())
	}
}

// ---------------------------------------------------------------------------
// BL-LOG-013: List purchase requests
// ---------------------------------------------------------------------------

type ListPurchaseRequestsParams struct {
	DepartureID string
	Status      string
	Page        int32
	PageSize    int32
}

type PurchaseRequestRowResult struct {
	ID          string
	DepartureID string
	Status      string
	TotalAmount int64
	CreatedAt   string
}

type ListPurchaseRequestsResult struct {
	Rows  []PurchaseRequestRowResult
	Total int32
}

func (a *Adapter) ListPurchaseRequests(ctx context.Context, params *ListPurchaseRequestsParams) (*ListPurchaseRequestsResult, error) {
	const op = "logistics_grpc_adapter.Adapter.ListPurchaseRequests"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.ListPurchaseRequests(ctx, &pb.LogDepthListPurchaseRequestsRequest{
		DepartureID: params.DepartureID,
		Status:      params.Status,
		Page:        params.Page,
		PageSize:    params.PageSize,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	rows := make([]PurchaseRequestRowResult, 0, len(resp.GetRows()))
	for _, r := range resp.GetRows() {
		rows = append(rows, PurchaseRequestRowResult{
			ID:          r.GetID(),
			DepartureID: r.GetDepartureID(),
			Status:      r.GetStatus(),
			TotalAmount: r.GetTotalAmount(),
			CreatedAt:   r.GetCreatedAt(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &ListPurchaseRequestsResult{Rows: rows, Total: resp.GetTotal()}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-014: Budget sync status
// ---------------------------------------------------------------------------

type GetBudgetSyncStatusParams struct {
	DepartureID string
}

type GetBudgetSyncStatusResult struct {
	DepartureID  string
	TotalBudget  int64
	Committed    int64
	Remaining    int64
	LastSyncedAt string
}

func (a *Adapter) GetBudgetSyncStatus(ctx context.Context, params *GetBudgetSyncStatusParams) (*GetBudgetSyncStatusResult, error) {
	const op = "logistics_grpc_adapter.Adapter.GetBudgetSyncStatus"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.GetBudgetSyncStatus(ctx, &pb.LogDepthGetBudgetSyncStatusRequest{
		DepartureID: params.DepartureID,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetBudgetSyncStatusResult{
		DepartureID:  resp.GetDepartureID(),
		TotalBudget:  resp.GetTotalBudget(),
		Committed:    resp.GetCommitted(),
		Remaining:    resp.GetRemaining(),
		LastSyncedAt: resp.GetLastSyncedAt(),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-015: Tiered approvals
// ---------------------------------------------------------------------------

type GetTieredApprovalsParams struct {
	DepartureID string
	Status      string
}

type TieredApprovalRowResult struct {
	ApprovalID string
	PRID       string
	Tier       string
	Status     string
	ApproverID string
	CreatedAt  string
}

type GetTieredApprovalsResult struct {
	Rows []TieredApprovalRowResult
}

func (a *Adapter) GetTieredApprovals(ctx context.Context, params *GetTieredApprovalsParams) (*GetTieredApprovalsResult, error) {
	const op = "logistics_grpc_adapter.Adapter.GetTieredApprovals"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.GetTieredApprovals(ctx, &pb.LogDepthGetTieredApprovalsRequest{
		DepartureID: params.DepartureID,
		Status:      params.Status,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	rows := make([]TieredApprovalRowResult, 0, len(resp.GetRows()))
	for _, r := range resp.GetRows() {
		rows = append(rows, TieredApprovalRowResult{
			ApprovalID: r.GetApprovalID(),
			PRID:       r.GetPRID(),
			Tier:       r.GetTier(),
			Status:     r.GetStatus(),
			ApproverID: r.GetApproverID(),
			CreatedAt:  r.GetCreatedAt(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetTieredApprovalsResult{Rows: rows}, nil
}

type DecideTieredApprovalParams struct {
	ApprovalID string
	Decision   string
	Notes      string
}

type DecideTieredApprovalResult struct {
	ApprovalID string
	Status     string
}

func (a *Adapter) DecideTieredApproval(ctx context.Context, params *DecideTieredApprovalParams) (*DecideTieredApprovalResult, error) {
	const op = "logistics_grpc_adapter.Adapter.DecideTieredApproval"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.DecideTieredApproval(ctx, &pb.LogDepthDecideTieredApprovalRequest{
		ApprovalID: params.ApprovalID,
		Decision:   params.Decision,
		Notes:      params.Notes,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &DecideTieredApprovalResult{ApprovalID: resp.GetApprovalID(), Status: resp.GetStatus()}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-016: Auto vendor selection
// ---------------------------------------------------------------------------

type AutoSelectVendorParams struct {
	PRID        string
	CategoryID  string
	BudgetLimit int64
}

type AutoSelectVendorResult struct {
	VendorID   string
	VendorName string
	Score      float64
}

func (a *Adapter) AutoSelectVendor(ctx context.Context, params *AutoSelectVendorParams) (*AutoSelectVendorResult, error) {
	const op = "logistics_grpc_adapter.Adapter.AutoSelectVendor"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.AutoSelectVendor(ctx, &pb.LogDepthAutoSelectVendorRequest{
		PRID:        params.PRID,
		CategoryID:  params.CategoryID,
		BudgetLimit: params.BudgetLimit,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &AutoSelectVendorResult{
		VendorID:   resp.GetVendorID(),
		VendorName: resp.GetVendorName(),
		Score:      resp.GetScore(),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-017: Partial / reverse GRN
// ---------------------------------------------------------------------------

type ReceivedItemInput struct {
	SKUID    string
	Quantity int32
}

type RecordPartialGRNParams struct {
	POID          string
	ReceivedItems []ReceivedItemInput
	Notes         string
}

type RecordPartialGRNResult struct {
	GRNID      string
	RecordedAt string
	IsComplete bool
}

func (a *Adapter) RecordPartialGRN(ctx context.Context, params *RecordPartialGRNParams) (*RecordPartialGRNResult, error) {
	const op = "logistics_grpc_adapter.Adapter.RecordPartialGRN"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	items := make([]*pb.LogDepthReceivedItem, 0, len(params.ReceivedItems))
	for _, it := range params.ReceivedItems {
		items = append(items, &pb.LogDepthReceivedItem{SKUID: it.SKUID, Quantity: it.Quantity})
	}
	resp, err := a.depthClient.RecordPartialGRN(ctx, &pb.LogDepthRecordPartialGRNRequest{
		POID:          params.POID,
		ReceivedItems: items,
		Notes:         params.Notes,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RecordPartialGRNResult{
		GRNID:      resp.GetGRNID(),
		RecordedAt: resp.GetRecordedAt(),
		IsComplete: resp.GetIsComplete(),
	}, nil
}

type ReverseGRNParams struct {
	GRNID  string
	Reason string
}

type ReverseGRNResult struct {
	ReversalID string
	ReversedAt string
}

func (a *Adapter) ReverseGRN(ctx context.Context, params *ReverseGRNParams) (*ReverseGRNResult, error) {
	const op = "logistics_grpc_adapter.Adapter.ReverseGRN"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.ReverseGRN(ctx, &pb.LogDepthReverseGRNRequest{
		GRNID:  params.GRNID,
		Reason: params.Reason,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &ReverseGRNResult{ReversalID: resp.GetReversalID(), ReversedAt: resp.GetReversedAt()}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-018: Barcode / SKU label
// ---------------------------------------------------------------------------

type GenerateBarcodeParams struct {
	SKUID  string
	Format string
}

type GenerateBarcodeResult struct {
	BarcodeID  string
	BarcodeURL string
}

func (a *Adapter) GenerateBarcode(ctx context.Context, params *GenerateBarcodeParams) (*GenerateBarcodeResult, error) {
	const op = "logistics_grpc_adapter.Adapter.GenerateBarcode"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.GenerateBarcode(ctx, &pb.LogDepthGenerateBarcodeRequest{
		SKUID:  params.SKUID,
		Format: params.Format,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &GenerateBarcodeResult{BarcodeID: resp.GetBarcodeID(), BarcodeURL: resp.GetBarcodeURL()}, nil
}

type PrintSKULabelParams struct {
	SKUID     string
	Copies    int32
	PrinterID string
}

type PrintSKULabelResult struct {
	JobID  string
	Queued bool
}

func (a *Adapter) PrintSKULabel(ctx context.Context, params *PrintSKULabelParams) (*PrintSKULabelResult, error) {
	const op = "logistics_grpc_adapter.Adapter.PrintSKULabel"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.PrintSKULabel(ctx, &pb.LogDepthPrintSKULabelRequest{
		SKUID:     params.SKUID,
		Copies:    params.Copies,
		PrinterID: params.PrinterID,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &PrintSKULabelResult{JobID: resp.GetJobID(), Queued: resp.GetQueued()}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-019: Warehouse management
// ---------------------------------------------------------------------------

type CreateWarehouseParams struct {
	Name     string
	Location string
	Type     string
}

type CreateWarehouseResult struct {
	WarehouseID string
	CreatedAt   string
}

func (a *Adapter) CreateWarehouse(ctx context.Context, params *CreateWarehouseParams) (*CreateWarehouseResult, error) {
	const op = "logistics_grpc_adapter.Adapter.CreateWarehouse"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.CreateWarehouse(ctx, &pb.LogDepthCreateWarehouseRequest{
		Name:     params.Name,
		Location: params.Location,
		Type:     params.Type,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &CreateWarehouseResult{WarehouseID: resp.GetWarehouseID(), CreatedAt: resp.GetCreatedAt()}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-020: Stock transfer
// ---------------------------------------------------------------------------

type TransferStockParams struct {
	SKUID           string
	Quantity        int32
	FromWarehouseID string
	ToWarehouseID   string
	Notes           string
}

type TransferStockResult struct {
	TransferID    string
	TransferredAt string
}

func (a *Adapter) TransferStock(ctx context.Context, params *TransferStockParams) (*TransferStockResult, error) {
	const op = "logistics_grpc_adapter.Adapter.TransferStock"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.TransferStock(ctx, &pb.LogDepthTransferStockRequest{
		SKUID:           params.SKUID,
		Quantity:        params.Quantity,
		FromWarehouseID: params.FromWarehouseID,
		ToWarehouseID:   params.ToWarehouseID,
		Notes:           params.Notes,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &TransferStockResult{TransferID: resp.GetTransferID(), TransferredAt: resp.GetTransferredAt()}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-021: Stock alerts
// ---------------------------------------------------------------------------

type GetStockAlertsParams struct {
	WarehouseID string
	Severity    string
}

type StockAlertRowResult struct {
	SKUID      string
	SKUName    string
	CurrentQty int32
	ReorderQty int32
	Severity   string
}

type GetStockAlertsResult struct {
	Alerts []StockAlertRowResult
}

func (a *Adapter) GetStockAlerts(ctx context.Context, params *GetStockAlertsParams) (*GetStockAlertsResult, error) {
	const op = "logistics_grpc_adapter.Adapter.GetStockAlerts"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.GetStockAlerts(ctx, &pb.LogDepthGetStockAlertsRequest{
		WarehouseID: params.WarehouseID,
		Severity:    params.Severity,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	alerts := make([]StockAlertRowResult, 0, len(resp.GetAlerts()))
	for _, al := range resp.GetAlerts() {
		alerts = append(alerts, StockAlertRowResult{
			SKUID:      al.GetSKUID(),
			SKUName:    al.GetSKUName(),
			CurrentQty: al.GetCurrentQty(),
			ReorderQty: al.GetReorderQty(),
			Severity:   al.GetSeverity(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetStockAlertsResult{Alerts: alerts}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-022: Reorder levels
// ---------------------------------------------------------------------------

type SetReorderLevelParams struct {
	SKUID      string
	ReorderQty int32
	MinQty     int32
}

type SetReorderLevelResult struct {
	SKUID     string
	UpdatedAt string
}

func (a *Adapter) SetReorderLevel(ctx context.Context, params *SetReorderLevelParams) (*SetReorderLevelResult, error) {
	const op = "logistics_grpc_adapter.Adapter.SetReorderLevel"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.SetReorderLevel(ctx, &pb.LogDepthSetReorderLevelRequest{
		SKUID:      params.SKUID,
		ReorderQty: params.ReorderQty,
		MinQty:     params.MinQty,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &SetReorderLevelResult{SKUID: resp.GetSKUID(), UpdatedAt: resp.GetUpdatedAt()}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-023: Stocktake lifecycle
// ---------------------------------------------------------------------------

type StartStocktakeParams struct {
	WarehouseID string
	Notes       string
}

type StartStocktakeResult struct {
	StocktakeID string
	StartedAt   string
}

func (a *Adapter) StartStocktake(ctx context.Context, params *StartStocktakeParams) (*StartStocktakeResult, error) {
	const op = "logistics_grpc_adapter.Adapter.StartStocktake"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.StartStocktake(ctx, &pb.LogDepthStartStocktakeRequest{
		WarehouseID: params.WarehouseID,
		Notes:       params.Notes,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &StartStocktakeResult{StocktakeID: resp.GetStocktakeID(), StartedAt: resp.GetStartedAt()}, nil
}

type RecordStocktakeCountParams struct {
	StocktakeID string
	SKUID       string
	CountedQty  int32
}

type RecordStocktakeCountResult struct {
	EntryID    string
	RecordedAt string
}

func (a *Adapter) RecordStocktakeCount(ctx context.Context, params *RecordStocktakeCountParams) (*RecordStocktakeCountResult, error) {
	const op = "logistics_grpc_adapter.Adapter.RecordStocktakeCount"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.RecordStocktakeCount(ctx, &pb.LogDepthRecordStocktakeCountRequest{
		StocktakeID: params.StocktakeID,
		SKUID:       params.SKUID,
		CountedQty:  params.CountedQty,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RecordStocktakeCountResult{EntryID: resp.GetEntryID(), RecordedAt: resp.GetRecordedAt()}, nil
}

type FinalizeStocktakeParams struct {
	StocktakeID string
}

type FinalizeStocktakeResult struct {
	StocktakeID   string
	FinalizedAt   string
	Discrepancies int32
}

func (a *Adapter) FinalizeStocktake(ctx context.Context, params *FinalizeStocktakeParams) (*FinalizeStocktakeResult, error) {
	const op = "logistics_grpc_adapter.Adapter.FinalizeStocktake"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.FinalizeStocktake(ctx, &pb.LogDepthFinalizeStocktakeRequest{
		StocktakeID: params.StocktakeID,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &FinalizeStocktakeResult{
		StocktakeID:   resp.GetStocktakeID(),
		FinalizedAt:   resp.GetFinalizedAt(),
		Discrepancies: resp.GetDiscrepancies(),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-024: Fulfillment size sync
// ---------------------------------------------------------------------------

type SizeSyncItemInput struct {
	BookingID string
	Size      string
}

type SyncFulfillmentSizesParams struct {
	DepartureID string
	Items       []SizeSyncItemInput
}

type SyncFulfillmentSizesResult struct {
	Synced    int32
	UpdatedAt string
}

func (a *Adapter) SyncFulfillmentSizes(ctx context.Context, params *SyncFulfillmentSizesParams) (*SyncFulfillmentSizesResult, error) {
	const op = "logistics_grpc_adapter.Adapter.SyncFulfillmentSizes"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	items := make([]*pb.LogDepthSizeSyncItem, 0, len(params.Items))
	for _, it := range params.Items {
		items = append(items, &pb.LogDepthSizeSyncItem{BookingID: it.BookingID, Size: it.Size})
	}
	resp, err := a.depthClient.SyncFulfillmentSizes(ctx, &pb.LogDepthSyncFulfillmentSizesRequest{
		DepartureID: params.DepartureID,
		Items:       items,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &SyncFulfillmentSizesResult{Synced: resp.GetSynced(), UpdatedAt: resp.GetUpdatedAt()}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-025: Courier tracking
// ---------------------------------------------------------------------------

type RecordCourierTrackingParams struct {
	BookingID      string
	Courier        string
	TrackingNumber string
	Status         string
	UpdatedAt      string
}

type RecordCourierTrackingResult struct {
	TrackingID string
	RecordedAt string
}

func (a *Adapter) RecordCourierTracking(ctx context.Context, params *RecordCourierTrackingParams) (*RecordCourierTrackingResult, error) {
	const op = "logistics_grpc_adapter.Adapter.RecordCourierTracking"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.RecordCourierTracking(ctx, &pb.LogDepthRecordCourierTrackingRequest{
		BookingID:      params.BookingID,
		Courier:        params.Courier,
		TrackingNumber: params.TrackingNumber,
		Status:         params.Status,
		UpdatedAt:      params.UpdatedAt,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RecordCourierTrackingResult{TrackingID: resp.GetTrackingID(), RecordedAt: resp.GetRecordedAt()}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-026: Returns
// ---------------------------------------------------------------------------

type ReturnItemInput struct {
	SKUID    string
	Quantity int32
}

type RecordReturnParams struct {
	BookingID string
	Reason    string
	Items     []ReturnItemInput
}

type RecordReturnResult struct {
	ReturnID   string
	RecordedAt string
}

func (a *Adapter) RecordReturn(ctx context.Context, params *RecordReturnParams) (*RecordReturnResult, error) {
	const op = "logistics_grpc_adapter.Adapter.RecordReturn"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	items := make([]*pb.LogDepthReturnItem, 0, len(params.Items))
	for _, it := range params.Items {
		items = append(items, &pb.LogDepthReturnItem{SKUID: it.SKUID, Quantity: it.Quantity})
	}
	resp, err := a.depthClient.RecordReturn(ctx, &pb.LogDepthRecordReturnRequest{
		BookingID: params.BookingID,
		Reason:    params.Reason,
		Items:     items,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RecordReturnResult{ReturnID: resp.GetReturnID(), RecordedAt: resp.GetRecordedAt()}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-027: Exchanges
// ---------------------------------------------------------------------------

type ExchangeItemInput struct {
	OldSKUID string
	NewSKUID string
	Quantity int32
}

type ProcessExchangeParams struct {
	BookingID string
	Items     []ExchangeItemInput
	Notes     string
}

type ProcessExchangeResult struct {
	ExchangeID  string
	ProcessedAt string
}

func (a *Adapter) ProcessExchange(ctx context.Context, params *ProcessExchangeParams) (*ProcessExchangeResult, error) {
	const op = "logistics_grpc_adapter.Adapter.ProcessExchange"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	items := make([]*pb.LogDepthExchangeItem, 0, len(params.Items))
	for _, it := range params.Items {
		items = append(items, &pb.LogDepthExchangeItem{
			OldSKUID: it.OldSKUID,
			NewSKUID: it.NewSKUID,
			Quantity: it.Quantity,
		})
	}
	resp, err := a.depthClient.ProcessExchange(ctx, &pb.LogDepthProcessExchangeRequest{
		BookingID: params.BookingID,
		Items:     items,
		Notes:     params.Notes,
	})
	if err != nil {
		wrapped := mapLogisticsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &ProcessExchangeResult{ExchangeID: resp.GetExchangeID(), ProcessedAt: resp.GetProcessedAt()}, nil
}
