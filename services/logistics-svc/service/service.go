package service

import (
	"context"

	"logistics-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for logistics-svc.
//
// Pilot scaffold surfaces the three standard scaffold endpoints plus the
// S3-E-02 fulfillment trigger:
//
//   - Liveness — process is up
//   - Readiness — process is up AND the database is reachable
//   - DbTxDiagnostic — writes + reads inside a WithTx, the canonical reference
//     for how services should use transactions (per docs/04-backend-conventions)
//   - OnBookingPaid — creates (or returns existing) fulfillment task for a
//     paid-in-full booking (S3-E-02 / BL-LOG-001)
//   - ShipFulfillmentTask — creates a shipment + tracking number (BL-LOG-002)
//   - GeneratePickupQR — generates a single-use pickup token with 7d TTL (BL-LOG-003)
//   - RedeemPickupQR — validates and marks a pickup token as used (BL-LOG-003)
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)
	DbTxDiagnostic(ctx context.Context, params *DbTxDiagnosticParams) (*DbTxDiagnosticResult, error)

	// OnBookingPaid creates a fulfillment task for a booking that has just
	// been fully paid. Idempotent: returns the existing task if one already
	// exists for this booking_id.
	OnBookingPaid(ctx context.Context, params *OnBookingPaidParams) (*OnBookingPaidResult, error)

	// ShipFulfillmentTask creates a shipment record + tracking number (BL-LOG-002).
	ShipFulfillmentTask(ctx context.Context, params *ShipFulfillmentTaskParams) (*ShipFulfillmentTaskResult, error)

	// GeneratePickupQR creates a single-use pickup QR token with 7d TTL (BL-LOG-003).
	GeneratePickupQR(ctx context.Context, params *GeneratePickupQRParams) (*GeneratePickupQRResult, error)

	// RedeemPickupQR validates and marks a pickup token as used (BL-LOG-003).
	RedeemPickupQR(ctx context.Context, params *RedeemPickupQRParams) (*RedeemPickupQRResult, error)

	// CreatePurchaseRequest creates a PR with optional budget gate (BL-LOG-010).
	CreatePurchaseRequest(ctx context.Context, params *CreatePurchaseRequestParams) (*CreatePurchaseRequestResult, error)

	// ApprovePurchaseRequest approves or rejects a pending PR (BL-LOG-010).
	ApprovePurchaseRequest(ctx context.Context, params *ApprovePurchaseRequestParams) (*ApprovePurchaseRequestResult, error)

	// RecordGRNWithQC records a GRN with QC pass/fail; posts AP journal only when qc_passed=true (BL-LOG-011).
	RecordGRNWithQC(ctx context.Context, params *RecordGRNWithQCParams) (*RecordGRNWithQCResult, error)

	// CreateKitAssembly atomically creates an idempotent kit assembly (BL-LOG-012).
	CreateKitAssembly(ctx context.Context, params *CreateKitAssemblyParams) (*CreateKitAssemblyResult, error)

	// ListFulfillmentTasks returns a paginated, filtered list of fulfillment tasks (ISSUE-018).
	ListFulfillmentTasks(ctx context.Context, params *ListFulfillmentTasksParams) (*ListFulfillmentTasksResult, error)

	// Wave 5 depth (BL-LOG-013..029)

	// ListPurchaseRequests returns a paginated list of purchase requests (BL-LOG-013).
	ListPurchaseRequests(ctx context.Context, params *ListPurchaseRequestsParams) (*ListPurchaseRequestsResult, error)

	// GetBudgetSyncStatus returns budget vs committed vs actual for a departure (BL-LOG-014).
	GetBudgetSyncStatus(ctx context.Context, params *GetBudgetSyncStatusParams) (*GetBudgetSyncStatusResult, error)

	// GetTieredApprovals returns tiered approval rows with optional level/status filter (BL-LOG-015).
	GetTieredApprovals(ctx context.Context, params *GetTieredApprovalsParams) (*GetTieredApprovalsResult, error)

	// DecideTieredApproval approves or rejects a tiered approval record (BL-LOG-015).
	DecideTieredApproval(ctx context.Context, params *DecideTieredApprovalParams) (*DecideTieredApprovalResult, error)

	// AutoSelectVendor auto-selects the best matching vendor for a category (BL-LOG-016).
	AutoSelectVendor(ctx context.Context, params *AutoSelectVendorParams) (*AutoSelectVendorResult, error)

	// RecordPartialGRN records a partial goods receipt note with item details (BL-LOG-017).
	RecordPartialGRN(ctx context.Context, params *RecordPartialGRNParams) (*RecordPartialGRNResult, error)

	// ReverseGRN marks a GRN as reversed (BL-LOG-017).
	ReverseGRN(ctx context.Context, params *ReverseGRNParams) (*ReverseGRNResult, error)

	// GenerateBarcode generates barcode data and label URL for a SKU (BL-LOG-020).
	GenerateBarcode(ctx context.Context, params *GenerateBarcodeParams) (*GenerateBarcodeResult, error)

	// PrintSKULabel returns a label URL for printing SKU labels (BL-LOG-020).
	PrintSKULabel(ctx context.Context, params *PrintSKULabelParams) (*PrintSKULabelResult, error)

	// CreateWarehouse creates a new warehouse record (BL-LOG-021).
	CreateWarehouse(ctx context.Context, params *CreateWarehouseParams) (*CreateWarehouseResult, error)

	// TransferStock records a stock transfer between warehouses (BL-LOG-021).
	TransferStock(ctx context.Context, params *TransferStockParams) (*TransferStockResult, error)

	// GetStockAlerts returns inventory items below reorder level (BL-LOG-022).
	GetStockAlerts(ctx context.Context, params *GetStockAlertsParams) (*GetStockAlertsResult, error)

	// SetReorderLevel sets or updates the reorder level for a SKU/warehouse (BL-LOG-022).
	SetReorderLevel(ctx context.Context, params *SetReorderLevelParams) (*SetReorderLevelResult, error)

	// StartStocktake begins a stocktake session for a warehouse (BL-LOG-024).
	StartStocktake(ctx context.Context, params *StartStocktakeParams) (*StartStocktakeResult, error)

	// RecordStocktakeCount records a counted quantity for a SKU in a stocktake (BL-LOG-024).
	RecordStocktakeCount(ctx context.Context, params *RecordStocktakeCountParams) (*RecordStocktakeCountResult, error)

	// FinalizeStocktake completes a stocktake and returns variance lines (BL-LOG-024).
	FinalizeStocktake(ctx context.Context, params *FinalizeStocktakeParams) (*FinalizeStocktakeResult, error)

	// SyncFulfillmentSizes syncs pilgrim size data for a departure (BL-LOG-025).
	SyncFulfillmentSizes(ctx context.Context, params *SyncFulfillmentSizesParams) (*SyncFulfillmentSizesResult, error)

	// RecordCourierTracking records a courier tracking update for a fulfillment task (BL-LOG-027).
	RecordCourierTracking(ctx context.Context, params *RecordCourierTrackingParams) (*RecordCourierTrackingResult, error)

	// RecordReturn records a return for a fulfillment task (BL-LOG-029).
	RecordReturn(ctx context.Context, params *RecordReturnParams) (*RecordReturnResult, error)

	// ProcessExchange processes an exchange for a previously recorded return (BL-LOG-029).
	ProcessExchange(ctx context.Context, params *ProcessExchangeParams) (*ProcessExchangeResult, error)
}

type Service struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	// appName is stamped on every row this service writes to the shared
	// public.diagnostics table so debugging can attribute rows to their origin.
	appName string

	store postgres_store.IStore
}

func NewService(
	logger *zerolog.Logger,
	tracer trace.Tracer,
	appName string,
	store postgres_store.IStore,
) IService {
	return &Service{
		logger:  logger,
		tracer:  tracer,
		appName: appName,
		store:   store,
	}
}
