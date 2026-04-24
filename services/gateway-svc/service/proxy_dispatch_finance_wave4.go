// proxy_dispatch_finance_wave4.go — gateway service dispatch for Wave 4 finance
// depth RPCs (BL-FIN-020..041).
//
// Each method is a thin delegation to finance_grpc_adapter.
// No business logic lives here.
package service

import (
	"context"

	"gateway-svc/adapter/finance_grpc_adapter"
)

func (s *Service) ScheduleBilling(ctx context.Context, params *finance_grpc_adapter.ScheduleBillingParams) (*finance_grpc_adapter.ScheduleBillingResult, error) {
	return s.adapters.financeGrpc.ScheduleBilling(ctx, params)
}

func (s *Service) RecordBankTransaction(ctx context.Context, params *finance_grpc_adapter.RecordBankTransactionParams) (*finance_grpc_adapter.RecordBankTransactionResult, error) {
	return s.adapters.financeGrpc.RecordBankTransaction(ctx, params)
}

func (s *Service) GetBankReconciliation(ctx context.Context, params *finance_grpc_adapter.GetBankReconciliationParams) (*finance_grpc_adapter.GetBankReconciliationResult, error) {
	return s.adapters.financeGrpc.GetBankReconciliation(ctx, params)
}

func (s *Service) GetARSubledger(ctx context.Context, params *finance_grpc_adapter.GetARSubledgerParams) (*finance_grpc_adapter.GetARSubledgerResult, error) {
	return s.adapters.financeGrpc.GetARSubledger(ctx, params)
}

func (s *Service) IssueDigitalReceipt(ctx context.Context, params *finance_grpc_adapter.IssueDigitalReceiptParams) (*finance_grpc_adapter.IssueDigitalReceiptResult, error) {
	return s.adapters.financeGrpc.IssueDigitalReceipt(ctx, params)
}

func (s *Service) GetDigitalReceipt(ctx context.Context, params *finance_grpc_adapter.GetDigitalReceiptParams) (*finance_grpc_adapter.DigitalReceiptResult, error) {
	return s.adapters.financeGrpc.GetDigitalReceipt(ctx, params)
}

func (s *Service) RecordManualPayment(ctx context.Context, params *finance_grpc_adapter.RecordManualPaymentParams) (*finance_grpc_adapter.RecordManualPaymentResult, error) {
	return s.adapters.financeGrpc.RecordManualPayment(ctx, params)
}

func (s *Service) CreateFinanceVendor(ctx context.Context, params *finance_grpc_adapter.CreateVendorParams) (*finance_grpc_adapter.CreateVendorResult, error) {
	return s.adapters.financeGrpc.CreateVendor(ctx, params)
}

func (s *Service) UpdateFinanceVendor(ctx context.Context, params *finance_grpc_adapter.UpdateVendorParams) (*finance_grpc_adapter.UpdateVendorResult, error) {
	return s.adapters.financeGrpc.UpdateVendor(ctx, params)
}

func (s *Service) ListFinanceVendors(ctx context.Context, params *finance_grpc_adapter.ListVendorsParams) (*finance_grpc_adapter.ListVendorsResult, error) {
	return s.adapters.financeGrpc.ListVendors(ctx, params)
}

func (s *Service) DeleteFinanceVendor(ctx context.Context, params *finance_grpc_adapter.DeleteVendorParams) (*finance_grpc_adapter.DeleteVendorResult, error) {
	return s.adapters.financeGrpc.DeleteVendor(ctx, params)
}

func (s *Service) GetAPSubledger(ctx context.Context, params *finance_grpc_adapter.GetAPSubledgerParams) (*finance_grpc_adapter.GetAPSubledgerResult, error) {
	return s.adapters.financeGrpc.GetAPSubledger(ctx, params)
}

func (s *Service) ListPendingAuthorizations(ctx context.Context, params *finance_grpc_adapter.ListPendingAuthorizationsParams) (*finance_grpc_adapter.ListPendingAuthorizationsResult, error) {
	return s.adapters.financeGrpc.ListPendingAuthorizations(ctx, params)
}

func (s *Service) DecidePaymentAuthorization(ctx context.Context, params *finance_grpc_adapter.DecidePaymentAuthorizationParams) (*finance_grpc_adapter.DecidePaymentAuthorizationResult, error) {
	return s.adapters.financeGrpc.DecidePaymentAuthorization(ctx, params)
}

func (s *Service) RecordPettyCash(ctx context.Context, params *finance_grpc_adapter.RecordPettyCashParams) (*finance_grpc_adapter.RecordPettyCashResult, error) {
	return s.adapters.financeGrpc.RecordPettyCash(ctx, params)
}

func (s *Service) ClosePettyCashPeriod(ctx context.Context, params *finance_grpc_adapter.ClosePettyCashPeriodParams) (*finance_grpc_adapter.ClosePettyCashPeriodResult, error) {
	return s.adapters.financeGrpc.ClosePettyCashPeriod(ctx, params)
}

func (s *Service) GetProjectCosts(ctx context.Context, params *finance_grpc_adapter.GetProjectCostsParams) (*finance_grpc_adapter.GetProjectCostsResult, error) {
	return s.adapters.financeGrpc.GetProjectCosts(ctx, params)
}

func (s *Service) GetDeparturePL(ctx context.Context, params *finance_grpc_adapter.GetDeparturePLParams) (*finance_grpc_adapter.GetDeparturePLResult, error) {
	return s.adapters.financeGrpc.GetDeparturePL(ctx, params)
}

func (s *Service) GetBudgetVsActual(ctx context.Context, params *finance_grpc_adapter.GetBudgetVsActualParams) (*finance_grpc_adapter.GetBudgetVsActualResult, error) {
	return s.adapters.financeGrpc.GetBudgetVsActual(ctx, params)
}

func (s *Service) TriggerAutoJournal(ctx context.Context, params *finance_grpc_adapter.TriggerAutoJournalParams) (*finance_grpc_adapter.TriggerAutoJournalResult, error) {
	return s.adapters.financeGrpc.TriggerAutoJournal(ctx, params)
}

func (s *Service) GetRevenueRecognitionPolicy(ctx context.Context) (*finance_grpc_adapter.RevenueRecognitionPolicyResult, error) {
	return s.adapters.financeGrpc.GetRevenueRecognitionPolicy(ctx)
}

func (s *Service) SetRevenueRecognitionPolicy(ctx context.Context, params *finance_grpc_adapter.SetRevenueRecognitionPolicyParams) error {
	return s.adapters.financeGrpc.SetRevenueRecognitionPolicy(ctx, params)
}

func (s *Service) SetExchangeRate(ctx context.Context, params *finance_grpc_adapter.SetExchangeRateParams) (*finance_grpc_adapter.SetExchangeRateResult, error) {
	return s.adapters.financeGrpc.SetExchangeRate(ctx, params)
}

func (s *Service) GetExchangeRate(ctx context.Context, params *finance_grpc_adapter.GetExchangeRateParams) (*finance_grpc_adapter.GetExchangeRateResult, error) {
	return s.adapters.financeGrpc.GetExchangeRate(ctx, params)
}

func (s *Service) CreateFixedAsset(ctx context.Context, params *finance_grpc_adapter.CreateFixedAssetParams) (*finance_grpc_adapter.CreateFixedAssetResult, error) {
	return s.adapters.financeGrpc.CreateFixedAsset(ctx, params)
}

func (s *Service) ListFixedAssets(ctx context.Context, category string) (*finance_grpc_adapter.ListFixedAssetsResult, error) {
	return s.adapters.financeGrpc.ListFixedAssets(ctx, category)
}

func (s *Service) RunDepreciation(ctx context.Context, params *finance_grpc_adapter.RunDepreciationParams) (*finance_grpc_adapter.RunDepreciationResult, error) {
	return s.adapters.financeGrpc.RunDepreciation(ctx, params)
}

func (s *Service) CalculateTax(ctx context.Context, params *finance_grpc_adapter.CalculateTaxParams) (*finance_grpc_adapter.CalculateTaxResult, error) {
	return s.adapters.financeGrpc.CalculateTax(ctx, params)
}

func (s *Service) GetTaxReport(ctx context.Context, params *finance_grpc_adapter.GetTaxReportParams) (*finance_grpc_adapter.GetTaxReportResult, error) {
	return s.adapters.financeGrpc.GetTaxReport(ctx, params)
}

func (s *Service) CreateCommissionPayout(ctx context.Context, params *finance_grpc_adapter.CreateCommissionPayoutParams) (*finance_grpc_adapter.CreateCommissionPayoutResult, error) {
	return s.adapters.financeGrpc.CreateCommissionPayout(ctx, params)
}

func (s *Service) DecideCommissionPayout(ctx context.Context, params *finance_grpc_adapter.DecideCommissionPayoutParams) (*finance_grpc_adapter.DecideCommissionPayoutResult, error) {
	return s.adapters.financeGrpc.DecideCommissionPayout(ctx, params)
}

func (s *Service) GetRealtimeFinancialSummary(ctx context.Context) (*finance_grpc_adapter.GetRealtimeFinancialSummaryResult, error) {
	return s.adapters.financeGrpc.GetRealtimeFinancialSummary(ctx)
}

func (s *Service) GetCashFlowDashboard(ctx context.Context, params *finance_grpc_adapter.GetCashFlowDashboardParams) (*finance_grpc_adapter.GetCashFlowDashboardResult, error) {
	return s.adapters.financeGrpc.GetCashFlowDashboard(ctx, params)
}

func (s *Service) GetAgingAlerts(ctx context.Context, params *finance_grpc_adapter.GetAgingAlertsParams) (*finance_grpc_adapter.GetAgingAlertsResult, error) {
	return s.adapters.financeGrpc.GetAgingAlerts(ctx, params)
}

func (s *Service) SearchFinanceAuditLog(ctx context.Context, params *finance_grpc_adapter.SearchFinanceAuditLogParams) (*finance_grpc_adapter.SearchFinanceAuditLogResult, error) {
	return s.adapters.financeGrpc.SearchFinanceAuditLog(ctx, params)
}
