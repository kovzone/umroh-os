// proxy_dispatch_logistics_depth.go — gateway service dispatch for logistics-svc Wave 5 depth RPCs.
// BL-LOG-013..029: purchase request list, budget sync, tiered approvals, auto vendor,
// partial/reverse GRN, barcode/labels, warehouse, stock transfer, stock alerts,
// reorder levels, stocktake lifecycle, size sync, courier tracking, returns, exchanges.
//
// Each method is a thin delegation to logistics_grpc_adapter. No business logic here.
package service

import (
	"context"

	"gateway-svc/adapter/logistics_grpc_adapter"
)

// ---------------------------------------------------------------------------
// Logistics depth — Wave 5 (BL-LOG-013..029)
// ---------------------------------------------------------------------------

func (s *Service) ListPurchaseRequests(ctx context.Context, params *logistics_grpc_adapter.ListPurchaseRequestsParams) (*logistics_grpc_adapter.ListPurchaseRequestsResult, error) {
	return s.adapters.logisticsGrpc.ListPurchaseRequests(ctx, params)
}

func (s *Service) GetBudgetSyncStatus(ctx context.Context, params *logistics_grpc_adapter.GetBudgetSyncStatusParams) (*logistics_grpc_adapter.GetBudgetSyncStatusResult, error) {
	return s.adapters.logisticsGrpc.GetBudgetSyncStatus(ctx, params)
}

func (s *Service) GetTieredApprovals(ctx context.Context, params *logistics_grpc_adapter.GetTieredApprovalsParams) (*logistics_grpc_adapter.GetTieredApprovalsResult, error) {
	return s.adapters.logisticsGrpc.GetTieredApprovals(ctx, params)
}

func (s *Service) DecideTieredApproval(ctx context.Context, params *logistics_grpc_adapter.DecideTieredApprovalParams) (*logistics_grpc_adapter.DecideTieredApprovalResult, error) {
	return s.adapters.logisticsGrpc.DecideTieredApproval(ctx, params)
}

func (s *Service) AutoSelectVendor(ctx context.Context, params *logistics_grpc_adapter.AutoSelectVendorParams) (*logistics_grpc_adapter.AutoSelectVendorResult, error) {
	return s.adapters.logisticsGrpc.AutoSelectVendor(ctx, params)
}

func (s *Service) RecordPartialGRN(ctx context.Context, params *logistics_grpc_adapter.RecordPartialGRNParams) (*logistics_grpc_adapter.RecordPartialGRNResult, error) {
	return s.adapters.logisticsGrpc.RecordPartialGRN(ctx, params)
}

func (s *Service) ReverseGRN(ctx context.Context, params *logistics_grpc_adapter.ReverseGRNParams) (*logistics_grpc_adapter.ReverseGRNResult, error) {
	return s.adapters.logisticsGrpc.ReverseGRN(ctx, params)
}

func (s *Service) GenerateBarcode(ctx context.Context, params *logistics_grpc_adapter.GenerateBarcodeParams) (*logistics_grpc_adapter.GenerateBarcodeResult, error) {
	return s.adapters.logisticsGrpc.GenerateBarcode(ctx, params)
}

func (s *Service) PrintSKULabel(ctx context.Context, params *logistics_grpc_adapter.PrintSKULabelParams) (*logistics_grpc_adapter.PrintSKULabelResult, error) {
	return s.adapters.logisticsGrpc.PrintSKULabel(ctx, params)
}

func (s *Service) CreateWarehouse(ctx context.Context, params *logistics_grpc_adapter.CreateWarehouseParams) (*logistics_grpc_adapter.CreateWarehouseResult, error) {
	return s.adapters.logisticsGrpc.CreateWarehouse(ctx, params)
}

func (s *Service) TransferStock(ctx context.Context, params *logistics_grpc_adapter.TransferStockParams) (*logistics_grpc_adapter.TransferStockResult, error) {
	return s.adapters.logisticsGrpc.TransferStock(ctx, params)
}

func (s *Service) GetStockAlerts(ctx context.Context, params *logistics_grpc_adapter.GetStockAlertsParams) (*logistics_grpc_adapter.GetStockAlertsResult, error) {
	return s.adapters.logisticsGrpc.GetStockAlerts(ctx, params)
}

func (s *Service) SetReorderLevel(ctx context.Context, params *logistics_grpc_adapter.SetReorderLevelParams) (*logistics_grpc_adapter.SetReorderLevelResult, error) {
	return s.adapters.logisticsGrpc.SetReorderLevel(ctx, params)
}

func (s *Service) StartStocktake(ctx context.Context, params *logistics_grpc_adapter.StartStocktakeParams) (*logistics_grpc_adapter.StartStocktakeResult, error) {
	return s.adapters.logisticsGrpc.StartStocktake(ctx, params)
}

func (s *Service) RecordStocktakeCount(ctx context.Context, params *logistics_grpc_adapter.RecordStocktakeCountParams) (*logistics_grpc_adapter.RecordStocktakeCountResult, error) {
	return s.adapters.logisticsGrpc.RecordStocktakeCount(ctx, params)
}

func (s *Service) FinalizeStocktake(ctx context.Context, params *logistics_grpc_adapter.FinalizeStocktakeParams) (*logistics_grpc_adapter.FinalizeStocktakeResult, error) {
	return s.adapters.logisticsGrpc.FinalizeStocktake(ctx, params)
}

func (s *Service) SyncFulfillmentSizes(ctx context.Context, params *logistics_grpc_adapter.SyncFulfillmentSizesParams) (*logistics_grpc_adapter.SyncFulfillmentSizesResult, error) {
	return s.adapters.logisticsGrpc.SyncFulfillmentSizes(ctx, params)
}

func (s *Service) RecordCourierTracking(ctx context.Context, params *logistics_grpc_adapter.RecordCourierTrackingParams) (*logistics_grpc_adapter.RecordCourierTrackingResult, error) {
	return s.adapters.logisticsGrpc.RecordCourierTracking(ctx, params)
}

func (s *Service) RecordReturn(ctx context.Context, params *logistics_grpc_adapter.RecordReturnParams) (*logistics_grpc_adapter.RecordReturnResult, error) {
	return s.adapters.logisticsGrpc.RecordReturn(ctx, params)
}

func (s *Service) ProcessExchange(ctx context.Context, params *logistics_grpc_adapter.ProcessExchangeParams) (*logistics_grpc_adapter.ProcessExchangeResult, error) {
	return s.adapters.logisticsGrpc.ProcessExchange(ctx, params)
}
