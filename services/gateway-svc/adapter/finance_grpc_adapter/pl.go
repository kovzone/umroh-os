// pl.go — gateway adapter methods for finance-svc depth RPCs (Wave 1B / Phase 6).
//
// Covers: RecognizeRevenue, GetPLReport, GetBalanceSheet.
// Each method translates gateway-local params → pb request, forwards via gRPC,
// and translates pb response → adapter-local types. Proto types do not leak
// past this package.
package finance_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/finance_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Adapter-local result types
// ---------------------------------------------------------------------------

// RecognizeRevenueResult is the gateway-local result for RecognizeRevenue.
type RecognizeRevenueResult struct {
	EntryID      string
	RecognizedAt string
	Replayed     bool
}

// PLLineItemResult is one P&L line item.
type PLLineItemResult struct {
	AccountCode string
	AccountName string
	Amount      int64
	Direction   string
}

// PLReportResult is the gateway-local result for GetPLReport.
type PLReportResult struct {
	PeriodFrom   string
	PeriodTo     string
	GeneratedAt  string
	TotalRevenue int64
	TotalCogs    int64
	GrossProfit  int64
	OtherIncome  int64
	OtherExpense int64
	NetProfit    int64
	Entries      []*PLLineItemResult
}

// BalanceSheetLineResult is one balance-sheet line item.
type BalanceSheetLineResult struct {
	AccountCode string
	AccountName string
	Balance     int64
}

// BalanceSheetResult is the gateway-local result for GetBalanceSheet.
type BalanceSheetResult struct {
	AsOfDate         string
	GeneratedAt      string
	Assets           []*BalanceSheetLineResult
	Liabilities      []*BalanceSheetLineResult
	Equity           []*BalanceSheetLineResult
	TotalAssets      int64
	TotalLiabilities int64
	TotalEquity      int64
}

// ---------------------------------------------------------------------------
// RecognizeRevenue
// ---------------------------------------------------------------------------

// RecognizeRevenue posts revenue recognition journal entries via finance-svc.
// Idempotent on departure_id.
func (a *Adapter) RecognizeRevenue(ctx context.Context, departureID string, totalAmountIdr int64) (*RecognizeRevenueResult, error) {
	const op = "finance_grpc_adapter.Adapter.RecognizeRevenue"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("departure_id", departureID),
		attribute.Int64("total_amount_idr", totalAmountIdr),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.financeDepthClient.RecognizeRevenue(ctx, &pb.RecognizeRevenueRequest{
		DepartureId:    departureID,
		TotalAmountIdr: totalAmountIdr,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.RecognizeRevenue failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &RecognizeRevenueResult{
		EntryID:      resp.GetEntryId(),
		RecognizedAt: resp.GetRecognizedAt(),
		Replayed:     resp.GetReplayed(),
	}, nil
}

// ---------------------------------------------------------------------------
// GetPLReport
// ---------------------------------------------------------------------------

// GetPLReport fetches the P&L report for the given date range.
// from and to are YYYY-MM-DD strings; empty = unbounded.
func (a *Adapter) GetPLReport(ctx context.Context, from, to string) (*PLReportResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetPLReport"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("from", from),
		attribute.String("to", to),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.financeDepthClient.GetPLReport(ctx, &pb.GetPLReportRequest{
		From: from,
		To:   to,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetPLReport failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	entries := make([]*PLLineItemResult, 0, len(resp.GetEntries()))
	for _, e := range resp.GetEntries() {
		entries = append(entries, &PLLineItemResult{
			AccountCode: e.GetAccountCode(),
			AccountName: e.GetAccountName(),
			Amount:      e.GetAmount(),
			Direction:   e.GetDirection(),
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &PLReportResult{
		PeriodFrom:   resp.GetPeriodFrom(),
		PeriodTo:     resp.GetPeriodTo(),
		GeneratedAt:  resp.GetGeneratedAt(),
		TotalRevenue: resp.GetTotalRevenue(),
		TotalCogs:    resp.GetTotalCogs(),
		GrossProfit:  resp.GetGrossProfit(),
		OtherIncome:  resp.GetOtherIncome(),
		OtherExpense: resp.GetOtherExpense(),
		NetProfit:    resp.GetNetProfit(),
		Entries:      entries,
	}, nil
}

// ---------------------------------------------------------------------------
// GetBalanceSheet
// ---------------------------------------------------------------------------

// GetBalanceSheet fetches the balance sheet as-of a date.
// asOf is a YYYY-MM-DD string; empty = latest.
func (a *Adapter) GetBalanceSheet(ctx context.Context, asOf string) (*BalanceSheetResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetBalanceSheet"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("as_of", asOf),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.financeDepthClient.GetBalanceSheet(ctx, &pb.GetBalanceSheetRequest{
		AsOf: asOf,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetBalanceSheet failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	mapLines := func(lines []*pb.BalanceSheetLineProto) []*BalanceSheetLineResult {
		out := make([]*BalanceSheetLineResult, 0, len(lines))
		for _, l := range lines {
			out = append(out, &BalanceSheetLineResult{
				AccountCode: l.GetAccountCode(),
				AccountName: l.GetAccountName(),
				Balance:     l.GetBalance(),
			})
		}
		return out
	}

	span.SetStatus(codes.Ok, "ok")
	return &BalanceSheetResult{
		AsOfDate:         resp.GetAsOfDate(),
		GeneratedAt:      resp.GetGeneratedAt(),
		Assets:           mapLines(resp.GetAssets()),
		Liabilities:      mapLines(resp.GetLiabilities()),
		Equity:           mapLines(resp.GetEquity()),
		TotalAssets:      resp.GetTotalAssets(),
		TotalLiabilities: resp.GetTotalLiabilities(),
		TotalEquity:      resp.GetTotalEquity(),
	}, nil
}
