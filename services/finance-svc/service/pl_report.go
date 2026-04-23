// pl_report.go — GetPLReport and GetBalanceSheet service-layer implementations
// for finance-svc (Wave 1B / BL-FIN-007..008).
//
// GetPLReport:
//   - Aggregates revenue (4xxx) and expense (5xxx) journal lines for a period.
//   - Returns TotalRevenue, TotalCOGS, GrossProfit, OtherIncome, OtherExpense, NetProfit.
//
// GetBalanceSheet:
//   - Aggregates asset (1xxx), liability (2xxx), and equity (3xxx) lines as of a date.
//   - Returns categorised lines with net balances.
//
// Both queries LEFT JOIN finance.account_codes for names; COALESCE guards
// ensure graceful fallback if migration 000017 has not been run.

package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"finance-svc/store/postgres_store/sqlc"
	"finance-svc/util/logging"

	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// P&L types
// ---------------------------------------------------------------------------

// PLLineItem holds one aggregated P&L line.
type PLLineItem struct {
	AccountCode string
	AccountName string
	Amount      int64
	Direction   string // "revenue" or "expense"
}

// PLReport is the result of GetPLReport.
type PLReport struct {
	PeriodFrom   time.Time
	PeriodTo     time.Time
	GeneratedAt  time.Time
	TotalRevenue int64 // sum of Cr 4001
	TotalCOGS    int64 // sum of Dr 5001
	GrossProfit  int64 // TotalRevenue - TotalCOGS
	OtherIncome  int64 // Cr account 4xxx (non-4001)
	OtherExpense int64 // Dr account 5xxx (non-5001)
	NetProfit    int64 // GrossProfit + OtherIncome - OtherExpense
	Entries      []PLLineItem
}

// GetPLReportParams holds the date range for the P&L report.
type GetPLReportParams struct {
	From time.Time // inclusive lower bound
	To   time.Time // inclusive upper bound
}

// ---------------------------------------------------------------------------
// Balance Sheet types
// ---------------------------------------------------------------------------

// BalanceSheetLine holds one aggregated balance-sheet line.
type BalanceSheetLine struct {
	AccountCode string
	AccountName string
	Balance     int64
}

// BalanceSheet is the result of GetBalanceSheet.
type BalanceSheet struct {
	AsOfDate         time.Time
	GeneratedAt      time.Time
	Assets           []BalanceSheetLine // account_code 1xxx, net debit balance
	Liabilities      []BalanceSheetLine // account_code 2xxx, net credit balance
	Equity           []BalanceSheetLine // account_code 3xxx
	TotalAssets      int64
	TotalLiabilities int64
	TotalEquity      int64
}

// GetBalanceSheetParams holds the as-of date for the balance sheet.
type GetBalanceSheetParams struct {
	AsOfDate time.Time
}

// ---------------------------------------------------------------------------
// GetPLReport
// ---------------------------------------------------------------------------

// GetPLReport returns a P&L report for the given date range.
func (svc *Service) GetPLReport(ctx context.Context, params *GetPLReportParams) (*PLReport, error) {
	const op = "service.Service.GetPLReport"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, svc.logger)
	span.SetAttributes(
		attribute.String("from", params.From.Format(time.RFC3339)),
		attribute.String("to", params.To.Format(time.RFC3339)),
	)
	logger.Info().Str("op", op).
		Str("from", params.From.Format(time.RFC3339)).
		Str("to", params.To.Format(time.RFC3339)).
		Msg("")

	var from, to pgtype.Timestamptz
	if !params.From.IsZero() {
		from = pgtype.Timestamptz{Time: params.From, Valid: true}
	}
	if !params.To.IsZero() {
		to = pgtype.Timestamptz{Time: params.To, Valid: true}
	}

	rows, err := svc.store.GetPLReportLines(ctx, sqlc.GetPLReportLinesParams{
		From: from,
		To:   to,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	report := &PLReport{
		PeriodFrom:  params.From,
		PeriodTo:    params.To,
		GeneratedAt: time.Now().UTC(),
		Entries:     make([]PLLineItem, 0, len(rows)),
	}

	for _, row := range rows {
		credit := numericToInt64(row.TotalCredit)
		debit := numericToInt64(row.TotalDebit)

		if strings.HasPrefix(row.AccountCode, "4") {
			// Revenue account — net credit is the revenue amount
			amount := credit
			report.Entries = append(report.Entries, PLLineItem{
				AccountCode: row.AccountCode,
				AccountName: row.AccountName,
				Amount:      amount,
				Direction:   "revenue",
			})
			if row.AccountCode == "4001" {
				report.TotalRevenue += amount
			} else {
				report.OtherIncome += amount
			}
		} else if strings.HasPrefix(row.AccountCode, "5") {
			// Expense account — net debit is the expense amount
			amount := debit
			report.Entries = append(report.Entries, PLLineItem{
				AccountCode: row.AccountCode,
				AccountName: row.AccountName,
				Amount:      amount,
				Direction:   "expense",
			})
			if row.AccountCode == "5001" {
				report.TotalCOGS += amount
			} else {
				report.OtherExpense += amount
			}
		}
	}

	report.GrossProfit = report.TotalRevenue - report.TotalCOGS
	report.NetProfit = report.GrossProfit + report.OtherIncome - report.OtherExpense

	span.SetStatus(otelCodes.Ok, "ok")
	return report, nil
}

// ---------------------------------------------------------------------------
// GetBalanceSheet
// ---------------------------------------------------------------------------

// GetBalanceSheet returns a balance sheet as of the given date.
func (svc *Service) GetBalanceSheet(ctx context.Context, params *GetBalanceSheetParams) (*BalanceSheet, error) {
	const op = "service.Service.GetBalanceSheet"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, svc.logger)
	span.SetAttributes(attribute.String("as_of", params.AsOfDate.Format(time.RFC3339)))
	logger.Info().Str("op", op).
		Str("as_of", params.AsOfDate.Format(time.RFC3339)).
		Msg("")

	var asOf pgtype.Timestamptz
	if !params.AsOfDate.IsZero() {
		asOf = pgtype.Timestamptz{Time: params.AsOfDate, Valid: true}
	}

	rows, err := svc.store.GetBalanceSheetLines(ctx, sqlc.GetBalanceSheetLinesParams{
		AsOf: asOf,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	bs := &BalanceSheet{
		AsOfDate:    params.AsOfDate,
		GeneratedAt: time.Now().UTC(),
		Assets:      make([]BalanceSheetLine, 0),
		Liabilities: make([]BalanceSheetLine, 0),
		Equity:      make([]BalanceSheetLine, 0),
	}

	for _, row := range rows {
		debit := numericToInt64(row.TotalDebit)
		credit := numericToInt64(row.TotalCredit)

		// Net balance: for normal-debit accounts (assets), net = debit - credit.
		// For normal-credit accounts (liabilities, equity), net = credit - debit.
		var balance int64
		if row.NormalBalance == "credit" {
			balance = credit - debit
		} else {
			balance = debit - credit
		}

		line := BalanceSheetLine{
			AccountCode: row.AccountCode,
			AccountName: row.AccountName,
			Balance:     balance,
		}

		switch {
		case strings.HasPrefix(row.AccountCode, "1"):
			bs.Assets = append(bs.Assets, line)
			bs.TotalAssets += balance
		case strings.HasPrefix(row.AccountCode, "2"):
			bs.Liabilities = append(bs.Liabilities, line)
			bs.TotalLiabilities += balance
		case strings.HasPrefix(row.AccountCode, "3"):
			bs.Equity = append(bs.Equity, line)
			bs.TotalEquity += balance
		}
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return bs, nil
}
