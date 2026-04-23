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
