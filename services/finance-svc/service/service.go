package service

import (
	"context"

	"finance-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for finance-svc.
//
//   - Liveness — process is up
//   - Readiness — process is up AND the database is reachable
//   - DbTxDiagnostic — writes + reads inside a WithTx, the canonical reference
//     for how services should use transactions (per docs/04-backend-conventions)
//   - FinancePing — BL-IAM-002 authenticated placeholder
//   - OnPaymentReceived — creates (or returns existing) double-entry journal
//     for a payment received event (S3-E-03 / BL-FIN-001..003)
//   - RecognizeRevenue — posts Dr 2001 / Cr 4001 journal when departure departs
//     or completes (Wave 1B / BL-FIN-006).
//   - GetPLReport — P&L report for a date range (Wave 1B / BL-FIN-007).
//   - GetBalanceSheet — balance sheet as of a date (Wave 1B / BL-FIN-008).
//   - OnGRNReceived — posts Dr 5001 (COGS) / Cr 2001 (AP) when GRN received
//     (BL-FIN-002). Idempotent on grn_id.
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)
	DbTxDiagnostic(ctx context.Context, params *DbTxDiagnosticParams) (*DbTxDiagnosticResult, error)

	// Finance — BL-IAM-002 placeholder.
	FinancePing(ctx context.Context, params *FinancePingParams) (*FinancePingResult, error)

	// OnPaymentReceived posts a double-entry journal for a payment received
	// event. Idempotent on idempotency_key = "payment:<invoice_id>".
	// Dr 1001 (Bank) / Cr 2001 (Pilgrim Liability).
	OnPaymentReceived(ctx context.Context, params *OnPaymentReceivedParams) (*OnPaymentReceivedResult, error)

	// GetFinanceSummary returns per-account aggregated debit/credit balances
	// from finance.journal_lines (S5-E-01 / BL-FIN-004).
	GetFinanceSummary(ctx context.Context) (*GetFinanceSummaryResult, error)

	// ListJournalEntries returns a cursor-paginated list of journal entries
	// with their lines, ordered by posted_at DESC (S5-E-01 / BL-FIN-005).
	ListJournalEntries(ctx context.Context, params *ListJournalEntriesParams) (*ListJournalEntriesResult, error)

	// RecognizeRevenue posts a double-entry journal for revenue recognition
	// (Dr 2001 Pilgrim Liability / Cr 4001 Revenue).
	// Idempotent on idempotency_key = "revenue:<departure_id>".
	// Called when departure status transitions to departed or completed.
	RecognizeRevenue(ctx context.Context, params *RecognizeRevenueParams) (*RecognizeRevenueResult, error)

	// GetPLReport returns a Profit & Loss report for the given date range.
	// Aggregates revenue (4xxx) and expense (5xxx) account lines.
	GetPLReport(ctx context.Context, params *GetPLReportParams) (*PLReport, error)

	// GetBalanceSheet returns a balance sheet as of the given date.
	// Aggregates asset (1xxx), liability (2xxx), and equity (3xxx) account lines.
	GetBalanceSheet(ctx context.Context, params *GetBalanceSheetParams) (*BalanceSheet, error)

	// OnGRNReceived posts Dr 5001 (COGS/Inventory Expense) / Cr 2001 (AP/Pilgrim
	// Liability) when a Goods Receipt Note is received (BL-FIN-002).
	// Idempotent on idempotency_key = "grn:<grn_id>".
	OnGRNReceived(ctx context.Context, params *OnGRNReceivedParams) (*OnGRNReceivedResult, error)

	// CorrectJournal posts a reversing counter-entry for an existing journal
	// entry (BL-FIN-006). The original entry is never deleted; only a new
	// reversal entry is inserted. Idempotent on "correction:<entry_id>".
	CorrectJournal(ctx context.Context, params *CorrectJournalParams) (*CorrectJournalResult, error)

	// DeleteJournalEntry always returns ErrForbidden (BL-FIN-006 anti-delete
	// guard). Corrections must use CorrectJournal instead.
	DeleteJournalEntry(ctx context.Context, entryID string) error

	// CreateDisbursementBatch creates an AP disbursement batch in pending_approval (BL-FIN-010).
	CreateDisbursementBatch(ctx context.Context, params *CreateDisbursementBatchParams) (*CreateDisbursementBatchResult, error)

	// ApproveDisbursement approves or rejects a batch; on approval posts Dr AP / Cr Cash (BL-FIN-010).
	ApproveDisbursement(ctx context.Context, params *ApproveDisbursementParams) (*ApproveDisbursementResult, error)

	// GetARAPAging returns aging buckets for AR, AP, or both as of a given date (BL-FIN-011).
	GetARAPAging(ctx context.Context, params *GetARAPAgingParams) (*GetARAPAgingResult, error)

	// Wave 4 Finance depth (BL-FIN-020..041)

	// ScheduleBilling creates invoices for all bookings in a departure (BL-FIN-020).
	ScheduleBilling(ctx context.Context, params *ScheduleBillingParams) (*ScheduleBillingResult, error)

	// RecordBankTransaction records a bank transaction (BL-FIN-021).
	RecordBankTransaction(ctx context.Context, params *RecordBankTransactionParams) (*RecordBankTransactionResult, error)

	// GetBankReconciliation returns bank transactions for reconciliation (BL-FIN-021).
	GetBankReconciliation(ctx context.Context, params *GetBankReconciliationParams) (*GetBankReconciliationResult, error)

	// GetARSubledger returns AR subledger for a booking/pilgrim (BL-FIN-022).
	GetARSubledger(ctx context.Context, params *GetARSubledgerParams) (*GetARSubledgerResult, error)

	// IssueDigitalReceipt creates a digital receipt for a payment (BL-FIN-023).
	IssueDigitalReceipt(ctx context.Context, params *IssueDigitalReceiptParams) (*IssueDigitalReceiptResult, error)

	// GetDigitalReceipt fetches a receipt by ID (BL-FIN-023).
	GetDigitalReceipt(ctx context.Context, params *GetDigitalReceiptParams) (*DigitalReceiptResult, error)

	// RecordManualPayment records a manual/down payment (BL-FIN-024).
	RecordManualPayment(ctx context.Context, params *RecordManualPaymentParams) (*RecordManualPaymentResult, error)

	// CreateVendor creates a vendor master record (BL-FIN-025).
	CreateVendor(ctx context.Context, params *CreateVendorParams) (*CreateVendorResult, error)

	// UpdateVendor updates a vendor master record (BL-FIN-025).
	UpdateVendor(ctx context.Context, params *UpdateVendorParams) (*UpdateVendorResult, error)

	// ListVendors lists vendors with optional category filter (BL-FIN-025).
	ListVendors(ctx context.Context, params *ListVendorsParams) (*ListVendorsResult, error)

	// DeleteVendor removes a vendor (BL-FIN-025).
	DeleteVendor(ctx context.Context, params *DeleteVendorParams) (*DeleteVendorResult, error)

	// GetAPSubledger returns AP subledger for a vendor (BL-FIN-026).
	GetAPSubledger(ctx context.Context, params *GetAPSubledgerParams) (*GetAPSubledgerResult, error)

	// ListPendingAuthorizations returns pending payment authorizations (BL-FIN-027).
	ListPendingAuthorizations(ctx context.Context, params *ListPendingAuthorizationsParams) (*ListPendingAuthorizationsResult, error)

	// DecidePaymentAuthorization approves or rejects a payment authorization (BL-FIN-027).
	DecidePaymentAuthorization(ctx context.Context, params *DecidePaymentAuthorizationParams) (*DecidePaymentAuthorizationResult, error)

	// RecordPettyCash records a petty cash entry (BL-FIN-028).
	RecordPettyCash(ctx context.Context, params *RecordPettyCashParams) (*RecordPettyCashResult, error)

	// ClosePettyCashPeriod closes the current petty cash period (BL-FIN-028).
	ClosePettyCashPeriod(ctx context.Context, params *ClosePettyCashPeriodParams) (*ClosePettyCashPeriodResult, error)

	// GetProjectCosts returns cost breakdown for a departure (BL-FIN-029).
	GetProjectCosts(ctx context.Context, params *GetProjectCostsParams) (*GetProjectCostsResult, error)

	// GetDeparturePL returns P&L for a departure (BL-FIN-030).
	GetDeparturePL(ctx context.Context, params *GetDeparturePLParams) (*GetDeparturePLResult, error)

	// GetBudgetVsActual returns budget vs actual comparison (BL-FIN-031).
	GetBudgetVsActual(ctx context.Context, params *GetBudgetVsActualParams) (*GetBudgetVsActualResult, error)

	// TriggerAutoJournal creates an auto journal entry (BL-FIN-032).
	TriggerAutoJournal(ctx context.Context, params *TriggerAutoJournalParams) (*TriggerAutoJournalResult, error)

	// GetRevenueRecognitionPolicy returns the revenue recognition policy (BL-FIN-033).
	GetRevenueRecognitionPolicy(ctx context.Context) (*GetRevenueRecognitionPolicyResult, error)

	// SetRevenueRecognitionPolicy sets the revenue recognition policy (BL-FIN-033).
	SetRevenueRecognitionPolicy(ctx context.Context, params *SetRevenueRecognitionPolicyParams) error

	// SetExchangeRate sets an exchange rate (BL-FIN-034).
	SetExchangeRate(ctx context.Context, params *SetExchangeRateParams) (*SetExchangeRateResult, error)

	// GetExchangeRate returns the exchange rate for a currency pair as of a date (BL-FIN-034).
	GetExchangeRate(ctx context.Context, params *GetExchangeRateParamsService) (*GetExchangeRateResult, error)

	// CreateFixedAsset creates a fixed asset record (BL-FIN-035).
	CreateFixedAsset(ctx context.Context, params *CreateFixedAssetParams) (*CreateFixedAssetResult, error)

	// ListFixedAssets returns all fixed assets optionally filtered by category (BL-FIN-035).
	ListFixedAssets(ctx context.Context, category string) (*ListFixedAssetsResult, error)

	// RunDepreciation runs straight-line monthly depreciation (BL-FIN-035).
	RunDepreciation(ctx context.Context, params *RunDepreciationParams) (*RunDepreciationResult, error)

	// CalculateTax performs a pure-math tax calculation (BL-FIN-036).
	CalculateTax(ctx context.Context, params *CalculateTaxParams) (*CalculateTaxResult, error)

	// GetTaxReport returns a tax report for a date range (BL-FIN-036).
	GetTaxReport(ctx context.Context, params *GetTaxReportParams) (*GetTaxReportResult, error)

	// CreateCommissionPayout creates a commission payout request (BL-FIN-037).
	CreateCommissionPayout(ctx context.Context, params *CreateCommissionPayoutParams) (*CreateCommissionPayoutResult, error)

	// DecideCommissionPayout approves or rejects a commission payout (BL-FIN-037).
	DecideCommissionPayout(ctx context.Context, params *DecideCommissionPayoutParams) (*DecideCommissionPayoutResult, error)

	// GetRealtimeFinancialSummary returns real-time aggregated financial data (BL-FIN-038).
	GetRealtimeFinancialSummary(ctx context.Context) (*GetRealtimeFinancialSummaryResult, error)

	// GetCashFlowDashboard returns the cash flow dashboard (BL-FIN-039).
	GetCashFlowDashboard(ctx context.Context, params *GetCashFlowDashboardParams) (*GetCashFlowDashboardResult, error)

	// GetAgingAlerts returns overdue AR/AP entries (BL-FIN-040).
	GetAgingAlerts(ctx context.Context, params *GetAgingAlertsParams) (*GetAgingAlertsResult, error)

	// SearchFinanceAuditLog returns finance audit log entries (BL-FIN-041).
	SearchFinanceAuditLog(ctx context.Context, params *SearchFinanceAuditLogParams) (*SearchFinanceAuditLogResult, error)
}

type Service struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	// appName is stamped on every row this service writes to the shared
	// public.diagnostics table so debugging can attribute rows to their origin.
	appName string

	store postgres_store.IStore

	// iamChecker is the consumer-side wrapper around iam-svc.CheckPermission.
	// Injected so tests can supply a double without spinning up a real gRPC server.
	iamChecker IamChecker
}

func NewService(
	logger *zerolog.Logger,
	tracer trace.Tracer,
	appName string,
	store postgres_store.IStore,
	iamChecker IamChecker,
) IService {
	return &Service{
		logger:     logger,
		tracer:     tracer,
		appName:    appName,
		store:      store,
		iamChecker: iamChecker,
	}
}
