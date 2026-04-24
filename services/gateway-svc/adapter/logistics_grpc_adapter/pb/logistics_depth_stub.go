// logistics_depth_stub.go — gateway-side gRPC client stubs for logistics-svc Wave 5 depth RPCs.
// BL-LOG-013..029: purchase request list, budget sync, tiered approvals, auto vendor selection,
// partial/reverse GRN, barcode/labels, warehouse creation, stock transfer, stock alerts,
// reorder levels, stocktake lifecycle, size sync, courier tracking, returns, exchanges.
//
// Run `make genpb` to replace with generated code once protoc is available.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// Method name constants
// ---------------------------------------------------------------------------

const (
	LogisticsService_ListPurchaseRequests_FullMethodName   = "/pb.LogisticsService/ListPurchaseRequests"
	LogisticsService_GetBudgetSyncStatus_FullMethodName    = "/pb.LogisticsService/GetBudgetSyncStatus"
	LogisticsService_GetTieredApprovals_FullMethodName     = "/pb.LogisticsService/GetTieredApprovals"
	LogisticsService_DecideTieredApproval_FullMethodName   = "/pb.LogisticsService/DecideTieredApproval"
	LogisticsService_AutoSelectVendor_FullMethodName       = "/pb.LogisticsService/AutoSelectVendor"
	LogisticsService_RecordPartialGRN_FullMethodName       = "/pb.LogisticsService/RecordPartialGRN"
	LogisticsService_ReverseGRN_FullMethodName             = "/pb.LogisticsService/ReverseGRN"
	LogisticsService_GenerateBarcode_FullMethodName        = "/pb.LogisticsService/GenerateBarcode"
	LogisticsService_PrintSKULabel_FullMethodName          = "/pb.LogisticsService/PrintSKULabel"
	LogisticsService_CreateWarehouse_FullMethodName        = "/pb.LogisticsService/CreateWarehouse"
	LogisticsService_TransferStock_FullMethodName          = "/pb.LogisticsService/TransferStock"
	LogisticsService_GetStockAlerts_FullMethodName         = "/pb.LogisticsService/GetStockAlerts"
	LogisticsService_SetReorderLevel_FullMethodName        = "/pb.LogisticsService/SetReorderLevel"
	LogisticsService_StartStocktake_FullMethodName         = "/pb.LogisticsService/StartStocktake"
	LogisticsService_RecordStocktakeCount_FullMethodName   = "/pb.LogisticsService/RecordStocktakeCount"
	LogisticsService_FinalizeStocktake_FullMethodName      = "/pb.LogisticsService/FinalizeStocktake"
	LogisticsService_SyncFulfillmentSizes_FullMethodName   = "/pb.LogisticsService/SyncFulfillmentSizes"
	LogisticsService_RecordCourierTracking_FullMethodName  = "/pb.LogisticsService/RecordCourierTracking"
	LogisticsService_RecordReturn_FullMethodName           = "/pb.LogisticsService/RecordReturn"
	LogisticsService_ProcessExchange_FullMethodName        = "/pb.LogisticsService/ProcessExchange"
)

// ---------------------------------------------------------------------------
// BL-LOG-013: List purchase requests
// ---------------------------------------------------------------------------

type LogDepthListPurchaseRequestsRequest struct {
	DepartureID string
	Status      string
	Page        int32
	PageSize    int32
}

func (x *LogDepthListPurchaseRequestsRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *LogDepthListPurchaseRequestsRequest) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *LogDepthListPurchaseRequestsRequest) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *LogDepthListPurchaseRequestsRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type LogDepthPurchaseRequestRow struct {
	ID          string
	DepartureID string
	Status      string
	TotalAmount int64
	CreatedAt   string
}

func (x *LogDepthPurchaseRequestRow) GetID() string {
	if x == nil {
		return ""
	}
	return x.ID
}
func (x *LogDepthPurchaseRequestRow) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *LogDepthPurchaseRequestRow) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *LogDepthPurchaseRequestRow) GetTotalAmount() int64 {
	if x == nil {
		return 0
	}
	return x.TotalAmount
}
func (x *LogDepthPurchaseRequestRow) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type LogDepthListPurchaseRequestsResponse struct {
	Rows  []*LogDepthPurchaseRequestRow
	Total int32
}

func (x *LogDepthListPurchaseRequestsResponse) GetRows() []*LogDepthPurchaseRequestRow {
	if x == nil {
		return nil
	}
	return x.Rows
}
func (x *LogDepthListPurchaseRequestsResponse) GetTotal() int32 {
	if x == nil {
		return 0
	}
	return x.Total
}

// ---------------------------------------------------------------------------
// BL-LOG-014: Budget sync status
// ---------------------------------------------------------------------------

type LogDepthGetBudgetSyncStatusRequest struct {
	DepartureID string
}

func (x *LogDepthGetBudgetSyncStatusRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type LogDepthGetBudgetSyncStatusResponse struct {
	DepartureID   string
	TotalBudget   int64
	Committed     int64
	Remaining     int64
	LastSyncedAt  string
}

func (x *LogDepthGetBudgetSyncStatusResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *LogDepthGetBudgetSyncStatusResponse) GetTotalBudget() int64 {
	if x == nil {
		return 0
	}
	return x.TotalBudget
}
func (x *LogDepthGetBudgetSyncStatusResponse) GetCommitted() int64 {
	if x == nil {
		return 0
	}
	return x.Committed
}
func (x *LogDepthGetBudgetSyncStatusResponse) GetRemaining() int64 {
	if x == nil {
		return 0
	}
	return x.Remaining
}
func (x *LogDepthGetBudgetSyncStatusResponse) GetLastSyncedAt() string {
	if x == nil {
		return ""
	}
	return x.LastSyncedAt
}

// ---------------------------------------------------------------------------
// BL-LOG-015: Tiered approvals
// ---------------------------------------------------------------------------

type LogDepthGetTieredApprovalsRequest struct {
	DepartureID string
	Status      string
}

func (x *LogDepthGetTieredApprovalsRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *LogDepthGetTieredApprovalsRequest) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

type LogDepthTieredApprovalRow struct {
	ApprovalID  string
	PRID        string
	Tier        string
	Status      string
	ApproverID  string
	CreatedAt   string
}

func (x *LogDepthTieredApprovalRow) GetApprovalID() string {
	if x == nil {
		return ""
	}
	return x.ApprovalID
}
func (x *LogDepthTieredApprovalRow) GetPRID() string {
	if x == nil {
		return ""
	}
	return x.PRID
}
func (x *LogDepthTieredApprovalRow) GetTier() string {
	if x == nil {
		return ""
	}
	return x.Tier
}
func (x *LogDepthTieredApprovalRow) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *LogDepthTieredApprovalRow) GetApproverID() string {
	if x == nil {
		return ""
	}
	return x.ApproverID
}
func (x *LogDepthTieredApprovalRow) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type LogDepthGetTieredApprovalsResponse struct {
	Rows []*LogDepthTieredApprovalRow
}

func (x *LogDepthGetTieredApprovalsResponse) GetRows() []*LogDepthTieredApprovalRow {
	if x == nil {
		return nil
	}
	return x.Rows
}

type LogDepthDecideTieredApprovalRequest struct {
	ApprovalID string
	Decision   string
	Notes      string
}

func (x *LogDepthDecideTieredApprovalRequest) GetApprovalID() string {
	if x == nil {
		return ""
	}
	return x.ApprovalID
}
func (x *LogDepthDecideTieredApprovalRequest) GetDecision() string {
	if x == nil {
		return ""
	}
	return x.Decision
}
func (x *LogDepthDecideTieredApprovalRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type LogDepthDecideTieredApprovalResponse struct {
	ApprovalID string
	Status     string
}

func (x *LogDepthDecideTieredApprovalResponse) GetApprovalID() string {
	if x == nil {
		return ""
	}
	return x.ApprovalID
}
func (x *LogDepthDecideTieredApprovalResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

// ---------------------------------------------------------------------------
// BL-LOG-016: Auto vendor selection
// ---------------------------------------------------------------------------

type LogDepthAutoSelectVendorRequest struct {
	PRID        string
	CategoryID  string
	BudgetLimit int64
}

func (x *LogDepthAutoSelectVendorRequest) GetPRID() string {
	if x == nil {
		return ""
	}
	return x.PRID
}
func (x *LogDepthAutoSelectVendorRequest) GetCategoryID() string {
	if x == nil {
		return ""
	}
	return x.CategoryID
}
func (x *LogDepthAutoSelectVendorRequest) GetBudgetLimit() int64 {
	if x == nil {
		return 0
	}
	return x.BudgetLimit
}

type LogDepthAutoSelectVendorResponse struct {
	VendorID   string
	VendorName string
	Score      float64
}

func (x *LogDepthAutoSelectVendorResponse) GetVendorID() string {
	if x == nil {
		return ""
	}
	return x.VendorID
}
func (x *LogDepthAutoSelectVendorResponse) GetVendorName() string {
	if x == nil {
		return ""
	}
	return x.VendorName
}
func (x *LogDepthAutoSelectVendorResponse) GetScore() float64 {
	if x == nil {
		return 0
	}
	return x.Score
}

// ---------------------------------------------------------------------------
// BL-LOG-017: Partial / reverse GRN
// ---------------------------------------------------------------------------

type LogDepthRecordPartialGRNRequest struct {
	POID          string
	ReceivedItems []*LogDepthReceivedItem
	Notes         string
}

func (x *LogDepthRecordPartialGRNRequest) GetPOID() string {
	if x == nil {
		return ""
	}
	return x.POID
}
func (x *LogDepthRecordPartialGRNRequest) GetReceivedItems() []*LogDepthReceivedItem {
	if x == nil {
		return nil
	}
	return x.ReceivedItems
}
func (x *LogDepthRecordPartialGRNRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type LogDepthReceivedItem struct {
	SKUID    string
	Quantity int32
}

func (x *LogDepthReceivedItem) GetSKUID() string {
	if x == nil {
		return ""
	}
	return x.SKUID
}
func (x *LogDepthReceivedItem) GetQuantity() int32 {
	if x == nil {
		return 0
	}
	return x.Quantity
}

type LogDepthRecordPartialGRNResponse struct {
	GRNID       string
	RecordedAt  string
	IsComplete  bool
}

func (x *LogDepthRecordPartialGRNResponse) GetGRNID() string {
	if x == nil {
		return ""
	}
	return x.GRNID
}
func (x *LogDepthRecordPartialGRNResponse) GetRecordedAt() string {
	if x == nil {
		return ""
	}
	return x.RecordedAt
}
func (x *LogDepthRecordPartialGRNResponse) GetIsComplete() bool {
	if x == nil {
		return false
	}
	return x.IsComplete
}

type LogDepthReverseGRNRequest struct {
	GRNID  string
	Reason string
}

func (x *LogDepthReverseGRNRequest) GetGRNID() string {
	if x == nil {
		return ""
	}
	return x.GRNID
}
func (x *LogDepthReverseGRNRequest) GetReason() string {
	if x == nil {
		return ""
	}
	return x.Reason
}

type LogDepthReverseGRNResponse struct {
	ReversalID string
	ReversedAt string
}

func (x *LogDepthReverseGRNResponse) GetReversalID() string {
	if x == nil {
		return ""
	}
	return x.ReversalID
}
func (x *LogDepthReverseGRNResponse) GetReversedAt() string {
	if x == nil {
		return ""
	}
	return x.ReversedAt
}

// ---------------------------------------------------------------------------
// BL-LOG-018: Barcode / SKU label
// ---------------------------------------------------------------------------

type LogDepthGenerateBarcodeRequest struct {
	SKUID  string
	Format string
}

func (x *LogDepthGenerateBarcodeRequest) GetSKUID() string {
	if x == nil {
		return ""
	}
	return x.SKUID
}
func (x *LogDepthGenerateBarcodeRequest) GetFormat() string {
	if x == nil {
		return ""
	}
	return x.Format
}

type LogDepthGenerateBarcodeResponse struct {
	BarcodeID  string
	BarcodeURL string
}

func (x *LogDepthGenerateBarcodeResponse) GetBarcodeID() string {
	if x == nil {
		return ""
	}
	return x.BarcodeID
}
func (x *LogDepthGenerateBarcodeResponse) GetBarcodeURL() string {
	if x == nil {
		return ""
	}
	return x.BarcodeURL
}

type LogDepthPrintSKULabelRequest struct {
	SKUID       string
	Copies      int32
	PrinterID   string
}

func (x *LogDepthPrintSKULabelRequest) GetSKUID() string {
	if x == nil {
		return ""
	}
	return x.SKUID
}
func (x *LogDepthPrintSKULabelRequest) GetCopies() int32 {
	if x == nil {
		return 0
	}
	return x.Copies
}
func (x *LogDepthPrintSKULabelRequest) GetPrinterID() string {
	if x == nil {
		return ""
	}
	return x.PrinterID
}

type LogDepthPrintSKULabelResponse struct {
	JobID    string
	Queued   bool
}

func (x *LogDepthPrintSKULabelResponse) GetJobID() string {
	if x == nil {
		return ""
	}
	return x.JobID
}
func (x *LogDepthPrintSKULabelResponse) GetQueued() bool {
	if x == nil {
		return false
	}
	return x.Queued
}

// ---------------------------------------------------------------------------
// BL-LOG-019: Warehouse management
// ---------------------------------------------------------------------------

type LogDepthCreateWarehouseRequest struct {
	Name     string
	Location string
	Type     string
}

func (x *LogDepthCreateWarehouseRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *LogDepthCreateWarehouseRequest) GetLocation() string {
	if x == nil {
		return ""
	}
	return x.Location
}
func (x *LogDepthCreateWarehouseRequest) GetType() string {
	if x == nil {
		return ""
	}
	return x.Type
}

type LogDepthCreateWarehouseResponse struct {
	WarehouseID string
	CreatedAt   string
}

func (x *LogDepthCreateWarehouseResponse) GetWarehouseID() string {
	if x == nil {
		return ""
	}
	return x.WarehouseID
}
func (x *LogDepthCreateWarehouseResponse) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

// ---------------------------------------------------------------------------
// BL-LOG-020: Stock transfer
// ---------------------------------------------------------------------------

type LogDepthTransferStockRequest struct {
	SKUID           string
	Quantity        int32
	FromWarehouseID string
	ToWarehouseID   string
	Notes           string
}

func (x *LogDepthTransferStockRequest) GetSKUID() string {
	if x == nil {
		return ""
	}
	return x.SKUID
}
func (x *LogDepthTransferStockRequest) GetQuantity() int32 {
	if x == nil {
		return 0
	}
	return x.Quantity
}
func (x *LogDepthTransferStockRequest) GetFromWarehouseID() string {
	if x == nil {
		return ""
	}
	return x.FromWarehouseID
}
func (x *LogDepthTransferStockRequest) GetToWarehouseID() string {
	if x == nil {
		return ""
	}
	return x.ToWarehouseID
}
func (x *LogDepthTransferStockRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type LogDepthTransferStockResponse struct {
	TransferID  string
	TransferredAt string
}

func (x *LogDepthTransferStockResponse) GetTransferID() string {
	if x == nil {
		return ""
	}
	return x.TransferID
}
func (x *LogDepthTransferStockResponse) GetTransferredAt() string {
	if x == nil {
		return ""
	}
	return x.TransferredAt
}

// ---------------------------------------------------------------------------
// BL-LOG-021: Stock alerts
// ---------------------------------------------------------------------------

type LogDepthGetStockAlertsRequest struct {
	WarehouseID string
	Severity    string
}

func (x *LogDepthGetStockAlertsRequest) GetWarehouseID() string {
	if x == nil {
		return ""
	}
	return x.WarehouseID
}
func (x *LogDepthGetStockAlertsRequest) GetSeverity() string {
	if x == nil {
		return ""
	}
	return x.Severity
}

type LogDepthStockAlertRow struct {
	SKUID       string
	SKUName     string
	CurrentQty  int32
	ReorderQty  int32
	Severity    string
}

func (x *LogDepthStockAlertRow) GetSKUID() string {
	if x == nil {
		return ""
	}
	return x.SKUID
}
func (x *LogDepthStockAlertRow) GetSKUName() string {
	if x == nil {
		return ""
	}
	return x.SKUName
}
func (x *LogDepthStockAlertRow) GetCurrentQty() int32 {
	if x == nil {
		return 0
	}
	return x.CurrentQty
}
func (x *LogDepthStockAlertRow) GetReorderQty() int32 {
	if x == nil {
		return 0
	}
	return x.ReorderQty
}
func (x *LogDepthStockAlertRow) GetSeverity() string {
	if x == nil {
		return ""
	}
	return x.Severity
}

type LogDepthGetStockAlertsResponse struct {
	Alerts []*LogDepthStockAlertRow
}

func (x *LogDepthGetStockAlertsResponse) GetAlerts() []*LogDepthStockAlertRow {
	if x == nil {
		return nil
	}
	return x.Alerts
}

// ---------------------------------------------------------------------------
// BL-LOG-022: Reorder levels
// ---------------------------------------------------------------------------

type LogDepthSetReorderLevelRequest struct {
	SKUID      string
	ReorderQty int32
	MinQty     int32
}

func (x *LogDepthSetReorderLevelRequest) GetSKUID() string {
	if x == nil {
		return ""
	}
	return x.SKUID
}
func (x *LogDepthSetReorderLevelRequest) GetReorderQty() int32 {
	if x == nil {
		return 0
	}
	return x.ReorderQty
}
func (x *LogDepthSetReorderLevelRequest) GetMinQty() int32 {
	if x == nil {
		return 0
	}
	return x.MinQty
}

type LogDepthSetReorderLevelResponse struct {
	SKUID     string
	UpdatedAt string
}

func (x *LogDepthSetReorderLevelResponse) GetSKUID() string {
	if x == nil {
		return ""
	}
	return x.SKUID
}
func (x *LogDepthSetReorderLevelResponse) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

// ---------------------------------------------------------------------------
// BL-LOG-023: Stocktake lifecycle
// ---------------------------------------------------------------------------

type LogDepthStartStocktakeRequest struct {
	WarehouseID string
	Notes       string
}

func (x *LogDepthStartStocktakeRequest) GetWarehouseID() string {
	if x == nil {
		return ""
	}
	return x.WarehouseID
}
func (x *LogDepthStartStocktakeRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type LogDepthStartStocktakeResponse struct {
	StocktakeID string
	StartedAt   string
}

func (x *LogDepthStartStocktakeResponse) GetStocktakeID() string {
	if x == nil {
		return ""
	}
	return x.StocktakeID
}
func (x *LogDepthStartStocktakeResponse) GetStartedAt() string {
	if x == nil {
		return ""
	}
	return x.StartedAt
}

type LogDepthRecordStocktakeCountRequest struct {
	StocktakeID string
	SKUID       string
	CountedQty  int32
}

func (x *LogDepthRecordStocktakeCountRequest) GetStocktakeID() string {
	if x == nil {
		return ""
	}
	return x.StocktakeID
}
func (x *LogDepthRecordStocktakeCountRequest) GetSKUID() string {
	if x == nil {
		return ""
	}
	return x.SKUID
}
func (x *LogDepthRecordStocktakeCountRequest) GetCountedQty() int32 {
	if x == nil {
		return 0
	}
	return x.CountedQty
}

type LogDepthRecordStocktakeCountResponse struct {
	EntryID    string
	RecordedAt string
}

func (x *LogDepthRecordStocktakeCountResponse) GetEntryID() string {
	if x == nil {
		return ""
	}
	return x.EntryID
}
func (x *LogDepthRecordStocktakeCountResponse) GetRecordedAt() string {
	if x == nil {
		return ""
	}
	return x.RecordedAt
}

type LogDepthFinalizeStocktakeRequest struct {
	StocktakeID string
}

func (x *LogDepthFinalizeStocktakeRequest) GetStocktakeID() string {
	if x == nil {
		return ""
	}
	return x.StocktakeID
}

type LogDepthFinalizeStocktakeResponse struct {
	StocktakeID   string
	FinalizedAt   string
	Discrepancies int32
}

func (x *LogDepthFinalizeStocktakeResponse) GetStocktakeID() string {
	if x == nil {
		return ""
	}
	return x.StocktakeID
}
func (x *LogDepthFinalizeStocktakeResponse) GetFinalizedAt() string {
	if x == nil {
		return ""
	}
	return x.FinalizedAt
}
func (x *LogDepthFinalizeStocktakeResponse) GetDiscrepancies() int32 {
	if x == nil {
		return 0
	}
	return x.Discrepancies
}

// ---------------------------------------------------------------------------
// BL-LOG-024: Fulfillment size sync
// ---------------------------------------------------------------------------

type LogDepthSizeSyncItem struct {
	BookingID string
	Size      string
}

func (x *LogDepthSizeSyncItem) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *LogDepthSizeSyncItem) GetSize() string {
	if x == nil {
		return ""
	}
	return x.Size
}

type LogDepthSyncFulfillmentSizesRequest struct {
	DepartureID string
	Items       []*LogDepthSizeSyncItem
}

func (x *LogDepthSyncFulfillmentSizesRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *LogDepthSyncFulfillmentSizesRequest) GetItems() []*LogDepthSizeSyncItem {
	if x == nil {
		return nil
	}
	return x.Items
}

type LogDepthSyncFulfillmentSizesResponse struct {
	Synced    int32
	UpdatedAt string
}

func (x *LogDepthSyncFulfillmentSizesResponse) GetSynced() int32 {
	if x == nil {
		return 0
	}
	return x.Synced
}
func (x *LogDepthSyncFulfillmentSizesResponse) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

// ---------------------------------------------------------------------------
// BL-LOG-025: Courier tracking
// ---------------------------------------------------------------------------

type LogDepthRecordCourierTrackingRequest struct {
	BookingID      string
	Courier        string
	TrackingNumber string
	Status         string
	UpdatedAt      string
}

func (x *LogDepthRecordCourierTrackingRequest) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *LogDepthRecordCourierTrackingRequest) GetCourier() string {
	if x == nil {
		return ""
	}
	return x.Courier
}
func (x *LogDepthRecordCourierTrackingRequest) GetTrackingNumber() string {
	if x == nil {
		return ""
	}
	return x.TrackingNumber
}
func (x *LogDepthRecordCourierTrackingRequest) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *LogDepthRecordCourierTrackingRequest) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

type LogDepthRecordCourierTrackingResponse struct {
	TrackingID string
	RecordedAt string
}

func (x *LogDepthRecordCourierTrackingResponse) GetTrackingID() string {
	if x == nil {
		return ""
	}
	return x.TrackingID
}
func (x *LogDepthRecordCourierTrackingResponse) GetRecordedAt() string {
	if x == nil {
		return ""
	}
	return x.RecordedAt
}

// ---------------------------------------------------------------------------
// BL-LOG-026: Returns
// ---------------------------------------------------------------------------

type LogDepthRecordReturnRequest struct {
	BookingID string
	Reason    string
	Items     []*LogDepthReturnItem
}

func (x *LogDepthRecordReturnRequest) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *LogDepthRecordReturnRequest) GetReason() string {
	if x == nil {
		return ""
	}
	return x.Reason
}
func (x *LogDepthRecordReturnRequest) GetItems() []*LogDepthReturnItem {
	if x == nil {
		return nil
	}
	return x.Items
}

type LogDepthReturnItem struct {
	SKUID    string
	Quantity int32
}

func (x *LogDepthReturnItem) GetSKUID() string {
	if x == nil {
		return ""
	}
	return x.SKUID
}
func (x *LogDepthReturnItem) GetQuantity() int32 {
	if x == nil {
		return 0
	}
	return x.Quantity
}

type LogDepthRecordReturnResponse struct {
	ReturnID   string
	RecordedAt string
}

func (x *LogDepthRecordReturnResponse) GetReturnID() string {
	if x == nil {
		return ""
	}
	return x.ReturnID
}
func (x *LogDepthRecordReturnResponse) GetRecordedAt() string {
	if x == nil {
		return ""
	}
	return x.RecordedAt
}

// ---------------------------------------------------------------------------
// BL-LOG-027: Exchanges
// ---------------------------------------------------------------------------

type LogDepthExchangeItem struct {
	OldSKUID string
	NewSKUID string
	Quantity int32
}

func (x *LogDepthExchangeItem) GetOldSKUID() string {
	if x == nil {
		return ""
	}
	return x.OldSKUID
}
func (x *LogDepthExchangeItem) GetNewSKUID() string {
	if x == nil {
		return ""
	}
	return x.NewSKUID
}
func (x *LogDepthExchangeItem) GetQuantity() int32 {
	if x == nil {
		return 0
	}
	return x.Quantity
}

type LogDepthProcessExchangeRequest struct {
	BookingID string
	Items     []*LogDepthExchangeItem
	Notes     string
}

func (x *LogDepthProcessExchangeRequest) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *LogDepthProcessExchangeRequest) GetItems() []*LogDepthExchangeItem {
	if x == nil {
		return nil
	}
	return x.Items
}
func (x *LogDepthProcessExchangeRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type LogDepthProcessExchangeResponse struct {
	ExchangeID  string
	ProcessedAt string
}

func (x *LogDepthProcessExchangeResponse) GetExchangeID() string {
	if x == nil {
		return ""
	}
	return x.ExchangeID
}
func (x *LogDepthProcessExchangeResponse) GetProcessedAt() string {
	if x == nil {
		return ""
	}
	return x.ProcessedAt
}

// ---------------------------------------------------------------------------
// LogisticsDepthClient — gRPC client interface
// ---------------------------------------------------------------------------

// LogisticsDepthClient is the client API for logistics-svc Wave 5 depth RPCs.
type LogisticsDepthClient interface {
	ListPurchaseRequests(ctx context.Context, in *LogDepthListPurchaseRequestsRequest, opts ...grpc.CallOption) (*LogDepthListPurchaseRequestsResponse, error)
	GetBudgetSyncStatus(ctx context.Context, in *LogDepthGetBudgetSyncStatusRequest, opts ...grpc.CallOption) (*LogDepthGetBudgetSyncStatusResponse, error)
	GetTieredApprovals(ctx context.Context, in *LogDepthGetTieredApprovalsRequest, opts ...grpc.CallOption) (*LogDepthGetTieredApprovalsResponse, error)
	DecideTieredApproval(ctx context.Context, in *LogDepthDecideTieredApprovalRequest, opts ...grpc.CallOption) (*LogDepthDecideTieredApprovalResponse, error)
	AutoSelectVendor(ctx context.Context, in *LogDepthAutoSelectVendorRequest, opts ...grpc.CallOption) (*LogDepthAutoSelectVendorResponse, error)
	RecordPartialGRN(ctx context.Context, in *LogDepthRecordPartialGRNRequest, opts ...grpc.CallOption) (*LogDepthRecordPartialGRNResponse, error)
	ReverseGRN(ctx context.Context, in *LogDepthReverseGRNRequest, opts ...grpc.CallOption) (*LogDepthReverseGRNResponse, error)
	GenerateBarcode(ctx context.Context, in *LogDepthGenerateBarcodeRequest, opts ...grpc.CallOption) (*LogDepthGenerateBarcodeResponse, error)
	PrintSKULabel(ctx context.Context, in *LogDepthPrintSKULabelRequest, opts ...grpc.CallOption) (*LogDepthPrintSKULabelResponse, error)
	CreateWarehouse(ctx context.Context, in *LogDepthCreateWarehouseRequest, opts ...grpc.CallOption) (*LogDepthCreateWarehouseResponse, error)
	TransferStock(ctx context.Context, in *LogDepthTransferStockRequest, opts ...grpc.CallOption) (*LogDepthTransferStockResponse, error)
	GetStockAlerts(ctx context.Context, in *LogDepthGetStockAlertsRequest, opts ...grpc.CallOption) (*LogDepthGetStockAlertsResponse, error)
	SetReorderLevel(ctx context.Context, in *LogDepthSetReorderLevelRequest, opts ...grpc.CallOption) (*LogDepthSetReorderLevelResponse, error)
	StartStocktake(ctx context.Context, in *LogDepthStartStocktakeRequest, opts ...grpc.CallOption) (*LogDepthStartStocktakeResponse, error)
	RecordStocktakeCount(ctx context.Context, in *LogDepthRecordStocktakeCountRequest, opts ...grpc.CallOption) (*LogDepthRecordStocktakeCountResponse, error)
	FinalizeStocktake(ctx context.Context, in *LogDepthFinalizeStocktakeRequest, opts ...grpc.CallOption) (*LogDepthFinalizeStocktakeResponse, error)
	SyncFulfillmentSizes(ctx context.Context, in *LogDepthSyncFulfillmentSizesRequest, opts ...grpc.CallOption) (*LogDepthSyncFulfillmentSizesResponse, error)
	RecordCourierTracking(ctx context.Context, in *LogDepthRecordCourierTrackingRequest, opts ...grpc.CallOption) (*LogDepthRecordCourierTrackingResponse, error)
	RecordReturn(ctx context.Context, in *LogDepthRecordReturnRequest, opts ...grpc.CallOption) (*LogDepthRecordReturnResponse, error)
	ProcessExchange(ctx context.Context, in *LogDepthProcessExchangeRequest, opts ...grpc.CallOption) (*LogDepthProcessExchangeResponse, error)
}

type logisticsDepthClient struct {
	cc grpc.ClientConnInterface
}

// NewLogisticsDepthClient returns a LogisticsDepthClient backed by the given conn.
func NewLogisticsDepthClient(cc grpc.ClientConnInterface) LogisticsDepthClient {
	return &logisticsDepthClient{cc}
}

func (c *logisticsDepthClient) ListPurchaseRequests(ctx context.Context, in *LogDepthListPurchaseRequestsRequest, opts ...grpc.CallOption) (*LogDepthListPurchaseRequestsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthListPurchaseRequestsResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_ListPurchaseRequests_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) GetBudgetSyncStatus(ctx context.Context, in *LogDepthGetBudgetSyncStatusRequest, opts ...grpc.CallOption) (*LogDepthGetBudgetSyncStatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthGetBudgetSyncStatusResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_GetBudgetSyncStatus_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) GetTieredApprovals(ctx context.Context, in *LogDepthGetTieredApprovalsRequest, opts ...grpc.CallOption) (*LogDepthGetTieredApprovalsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthGetTieredApprovalsResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_GetTieredApprovals_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) DecideTieredApproval(ctx context.Context, in *LogDepthDecideTieredApprovalRequest, opts ...grpc.CallOption) (*LogDepthDecideTieredApprovalResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthDecideTieredApprovalResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_DecideTieredApproval_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) AutoSelectVendor(ctx context.Context, in *LogDepthAutoSelectVendorRequest, opts ...grpc.CallOption) (*LogDepthAutoSelectVendorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthAutoSelectVendorResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_AutoSelectVendor_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) RecordPartialGRN(ctx context.Context, in *LogDepthRecordPartialGRNRequest, opts ...grpc.CallOption) (*LogDepthRecordPartialGRNResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthRecordPartialGRNResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_RecordPartialGRN_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) ReverseGRN(ctx context.Context, in *LogDepthReverseGRNRequest, opts ...grpc.CallOption) (*LogDepthReverseGRNResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthReverseGRNResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_ReverseGRN_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) GenerateBarcode(ctx context.Context, in *LogDepthGenerateBarcodeRequest, opts ...grpc.CallOption) (*LogDepthGenerateBarcodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthGenerateBarcodeResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_GenerateBarcode_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) PrintSKULabel(ctx context.Context, in *LogDepthPrintSKULabelRequest, opts ...grpc.CallOption) (*LogDepthPrintSKULabelResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthPrintSKULabelResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_PrintSKULabel_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) CreateWarehouse(ctx context.Context, in *LogDepthCreateWarehouseRequest, opts ...grpc.CallOption) (*LogDepthCreateWarehouseResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthCreateWarehouseResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_CreateWarehouse_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) TransferStock(ctx context.Context, in *LogDepthTransferStockRequest, opts ...grpc.CallOption) (*LogDepthTransferStockResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthTransferStockResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_TransferStock_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) GetStockAlerts(ctx context.Context, in *LogDepthGetStockAlertsRequest, opts ...grpc.CallOption) (*LogDepthGetStockAlertsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthGetStockAlertsResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_GetStockAlerts_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) SetReorderLevel(ctx context.Context, in *LogDepthSetReorderLevelRequest, opts ...grpc.CallOption) (*LogDepthSetReorderLevelResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthSetReorderLevelResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_SetReorderLevel_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) StartStocktake(ctx context.Context, in *LogDepthStartStocktakeRequest, opts ...grpc.CallOption) (*LogDepthStartStocktakeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthStartStocktakeResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_StartStocktake_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) RecordStocktakeCount(ctx context.Context, in *LogDepthRecordStocktakeCountRequest, opts ...grpc.CallOption) (*LogDepthRecordStocktakeCountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthRecordStocktakeCountResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_RecordStocktakeCount_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) FinalizeStocktake(ctx context.Context, in *LogDepthFinalizeStocktakeRequest, opts ...grpc.CallOption) (*LogDepthFinalizeStocktakeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthFinalizeStocktakeResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_FinalizeStocktake_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) SyncFulfillmentSizes(ctx context.Context, in *LogDepthSyncFulfillmentSizesRequest, opts ...grpc.CallOption) (*LogDepthSyncFulfillmentSizesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthSyncFulfillmentSizesResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_SyncFulfillmentSizes_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) RecordCourierTracking(ctx context.Context, in *LogDepthRecordCourierTrackingRequest, opts ...grpc.CallOption) (*LogDepthRecordCourierTrackingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthRecordCourierTrackingResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_RecordCourierTracking_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) RecordReturn(ctx context.Context, in *LogDepthRecordReturnRequest, opts ...grpc.CallOption) (*LogDepthRecordReturnResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthRecordReturnResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_RecordReturn_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *logisticsDepthClient) ProcessExchange(ctx context.Context, in *LogDepthProcessExchangeRequest, opts ...grpc.CallOption) (*LogDepthProcessExchangeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogDepthProcessExchangeResponse)
	if err := c.cc.Invoke(ctx, LogisticsService_ProcessExchange_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
