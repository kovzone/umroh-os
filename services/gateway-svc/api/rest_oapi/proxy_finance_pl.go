// proxy_finance_pl.go — gateway REST handlers for finance depth RPCs
// (Wave 1B / Phase 6).
//
// Route topology (all bearer-protected):
//   POST /v1/finance/recognize-revenue  → RecognizeRevenue
//   GET  /v1/finance/pl                 → GetPLReport (query: from, to as YYYY-MM-DD)
//   GET  /v1/finance/balance-sheet      → GetBalanceSheet (query: as_of as YYYY-MM-DD)
//
// Per ADR-0009: gateway is the single REST entry-point; finance-svc is pure gRPC.
package rest_oapi

import (
	"errors"

	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Request / response body types
// ---------------------------------------------------------------------------

// RecognizeRevenueBody is the JSON body for POST /v1/finance/recognize-revenue.
type RecognizeRevenueBody struct {
	DepartureID    string `json:"departure_id"`
	TotalAmountIdr int64  `json:"total_amount_idr"`
}

// RecognizeRevenueResponseData is the response for RecognizeRevenue.
type RecognizeRevenueResponseData struct {
	EntryID      string `json:"entry_id"`
	RecognizedAt string `json:"recognized_at"`
	Replayed     bool   `json:"replayed"`
}

// PLLineItemData is the JSON representation of one P&L line item.
type PLLineItemData struct {
	AccountCode string `json:"account_code"`
	AccountName string `json:"account_name"`
	Amount      int64  `json:"amount"`
	Direction   string `json:"direction"`
}

// PLReportResponseData is the response for GET /v1/finance/pl.
type PLReportResponseData struct {
	PeriodFrom   string           `json:"period_from"`
	PeriodTo     string           `json:"period_to"`
	GeneratedAt  string           `json:"generated_at"`
	TotalRevenue int64            `json:"total_revenue"`
	TotalCogs    int64            `json:"total_cogs"`
	GrossProfit  int64            `json:"gross_profit"`
	OtherIncome  int64            `json:"other_income"`
	OtherExpense int64            `json:"other_expense"`
	NetProfit    int64            `json:"net_profit"`
	Entries      []PLLineItemData `json:"entries"`
}

// BalanceSheetLineData is the JSON representation of one balance sheet line.
type BalanceSheetLineData struct {
	AccountCode string `json:"account_code"`
	AccountName string `json:"account_name"`
	Balance     int64  `json:"balance"`
}

// BalanceSheetResponseData is the response for GET /v1/finance/balance-sheet.
type BalanceSheetResponseData struct {
	AsOfDate         string                 `json:"as_of_date"`
	GeneratedAt      string                 `json:"generated_at"`
	Assets           []BalanceSheetLineData `json:"assets"`
	Liabilities      []BalanceSheetLineData `json:"liabilities"`
	Equity           []BalanceSheetLineData `json:"equity"`
	TotalAssets      int64                  `json:"total_assets"`
	TotalLiabilities int64                  `json:"total_liabilities"`
	TotalEquity      int64                  `json:"total_equity"`
}

// ---------------------------------------------------------------------------
// RecognizeRevenue — POST /v1/finance/recognize-revenue (bearer)
// ---------------------------------------------------------------------------

func (s *Server) RecognizeRevenue(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecognizeRevenue"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("endpoint", "POST /v1/finance/recognize-revenue"))
	logger.Info().Str("op", op).Msg("")

	var body RecognizeRevenueBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.RecognizeRevenue(ctx, body.DepartureID, body.TotalAmountIdr)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeFinanceError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": RecognizeRevenueResponseData{
			EntryID:      result.EntryID,
			RecognizedAt: result.RecognizedAt,
			Replayed:     result.Replayed,
		},
	})
}

// ---------------------------------------------------------------------------
// GetPLReport — GET /v1/finance/pl (bearer)
// ---------------------------------------------------------------------------

func (s *Server) GetPLReport(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetPLReport"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	from := c.Query("from")
	to := c.Query("to")

	span.SetAttributes(
		attribute.String("endpoint", "GET /v1/finance/pl"),
		attribute.String("from", from),
		attribute.String("to", to),
	)
	logger.Info().Str("op", op).Str("from", from).Str("to", to).Msg("")

	result, err := s.svc.GetPLReport(ctx, from, to)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeFinanceError(c, span, err)
	}

	entries := make([]PLLineItemData, 0, len(result.Entries))
	for _, e := range result.Entries {
		entries = append(entries, PLLineItemData{
			AccountCode: e.AccountCode,
			AccountName: e.AccountName,
			Amount:      e.Amount,
			Direction:   e.Direction,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": PLReportResponseData{
			PeriodFrom:   result.PeriodFrom,
			PeriodTo:     result.PeriodTo,
			GeneratedAt:  result.GeneratedAt,
			TotalRevenue: result.TotalRevenue,
			TotalCogs:    result.TotalCogs,
			GrossProfit:  result.GrossProfit,
			OtherIncome:  result.OtherIncome,
			OtherExpense: result.OtherExpense,
			NetProfit:    result.NetProfit,
			Entries:      entries,
		},
	})
}

// ---------------------------------------------------------------------------
// GetBalanceSheet — GET /v1/finance/balance-sheet (bearer)
// ---------------------------------------------------------------------------

func (s *Server) GetBalanceSheet(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetBalanceSheet"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	asOf := c.Query("as_of")

	span.SetAttributes(
		attribute.String("endpoint", "GET /v1/finance/balance-sheet"),
		attribute.String("as_of", asOf),
	)
	logger.Info().Str("op", op).Str("as_of", asOf).Msg("")

	result, err := s.svc.GetBalanceSheet(ctx, asOf)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeFinanceError(c, span, err)
	}

	// Map each category.
	assets := make([]BalanceSheetLineData, 0, len(result.Assets))
	for _, l := range result.Assets {
		assets = append(assets, BalanceSheetLineData{
			AccountCode: l.AccountCode,
			AccountName: l.AccountName,
			Balance:     l.Balance,
		})
	}
	liabilities := make([]BalanceSheetLineData, 0, len(result.Liabilities))
	for _, l := range result.Liabilities {
		liabilities = append(liabilities, BalanceSheetLineData{
			AccountCode: l.AccountCode,
			AccountName: l.AccountName,
			Balance:     l.Balance,
		})
	}
	equity := make([]BalanceSheetLineData, 0, len(result.Equity))
	for _, l := range result.Equity {
		equity = append(equity, BalanceSheetLineData{
			AccountCode: l.AccountCode,
			AccountName: l.AccountName,
			Balance:     l.Balance,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": BalanceSheetResponseData{
			AsOfDate:         result.AsOfDate,
			GeneratedAt:      result.GeneratedAt,
			Assets:           assets,
			Liabilities:      liabilities,
			Equity:           equity,
			TotalAssets:      result.TotalAssets,
			TotalLiabilities: result.TotalLiabilities,
			TotalEquity:      result.TotalEquity,
		},
	})
}
