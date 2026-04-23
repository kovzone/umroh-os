// proxy_finance_reports.go — gateway service layer for finance report routes
// (S5-E-01 / BL-FIN-004..005).
//
// Thin delegation layer: service methods map 1:1 to finance_grpc_adapter calls.
// No business logic here; all finance logic lives in finance-svc.

package service

import (
	"context"

	"gateway-svc/adapter/finance_grpc_adapter"
)

func (s *Service) GetFinanceSummary(ctx context.Context) (*finance_grpc_adapter.GetFinanceSummaryResult, error) {
	return s.adapters.financeGrpc.GetFinanceSummary(ctx)
}

func (s *Service) ListJournalEntries(ctx context.Context, params *finance_grpc_adapter.ListJournalEntriesParams) (*finance_grpc_adapter.ListJournalEntriesResult, error) {
	return s.adapters.financeGrpc.ListJournalEntries(ctx, params)
}
