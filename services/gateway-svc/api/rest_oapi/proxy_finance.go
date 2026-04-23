// proxy_finance.go — gateway REST handlers for finance reports (S5-E-01).
//
// Route topology:
//   GET /v1/finance/system/live — public probe proxy (finance_rest_adapter)
//   GET /v1/finance/summary    → GetFinanceSummary  (bearer, finance_grpc_adapter)
//   GET /v1/finance/journals   → ListJournalEntries (bearer, finance_grpc_adapter)
//
// Per ADR-0009: gateway is the single REST entry-point; finance-svc is pure gRPC
// for business routes.

package rest_oapi

import (
	"errors"

	"gateway-svc/adapter/finance_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// GetFinanceSystemLive proxies the iam-svc liveness probe through the gateway.
// Scaffold-time proof of the REST adapter pattern.
//
// GetFinanceSystemLive implements ServerInterface.
func (s *Server) GetFinanceSystemLive(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetFinanceSystemLive"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/finance/system/live"),
		attribute.String("method", "GET"),
	)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetFinanceSystemLive(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "BAD_GATEWAY",
				"message": err.Error(),
			},
		})
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(LiveResponse{
		Data: struct {
			Ok bool `json:"ok"`
		}{Ok: result.OK},
	})
}

// ---------------------------------------------------------------------------
// Response / request body types for finance reports
// ---------------------------------------------------------------------------

// AccountBalanceData is the JSON representation of one account balance.
type AccountBalanceData struct {
	AccountCode string `json:"account_code"`
	DebitTotal  int64  `json:"debit_total"`
	CreditTotal int64  `json:"credit_total"`
	Net         int64  `json:"net"`
}

// FinanceSummaryResponse is the response for GET /v1/finance/summary.
type FinanceSummaryResponse struct {
	Data []AccountBalanceData `json:"data"`
}

// JournalLineData is the JSON representation of one journal line.
type JournalLineData struct {
	ID          string `json:"id"`
	EntryID     string `json:"entry_id"`
	AccountCode string `json:"account_code"`
	Debit       int64  `json:"debit"`
	Credit      int64  `json:"credit"`
}

// JournalEntryData is the JSON representation of one journal entry with lines.
type JournalEntryData struct {
	ID             string            `json:"id"`
	IdempotencyKey string            `json:"idempotency_key"`
	SourceType     string            `json:"source_type"`
	SourceID       string            `json:"source_id"`
	PostedAt       string            `json:"posted_at"`
	Description    string            `json:"description,omitempty"`
	Lines          []JournalLineData `json:"lines"`
}

// JournalListResponse is the response for GET /v1/finance/journals.
type JournalListResponse struct {
	Data       []JournalEntryData `json:"data"`
	NextCursor string             `json:"next_cursor,omitempty"`
}

// ---------------------------------------------------------------------------
// GetFinanceSummary — GET /v1/finance/summary (bearer)
// ---------------------------------------------------------------------------

func (s *Server) GetFinanceSummary(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetFinanceSummary"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "GET /v1/finance/summary"),
	)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetFinanceSummary(ctx)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeFinanceError(c, span, err)
	}

	data := make([]AccountBalanceData, 0, len(result.Accounts))
	for _, a := range result.Accounts {
		data = append(data, AccountBalanceData{
			AccountCode: a.AccountCode,
			DebitTotal:  a.DebitTotal,
			CreditTotal: a.CreditTotal,
			Net:         a.Net,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(FinanceSummaryResponse{Data: data})
}

// ---------------------------------------------------------------------------
// ListJournals — GET /v1/finance/journals (bearer)
// ---------------------------------------------------------------------------

func (s *Server) ListJournals(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListJournals"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	limit := int32(c.QueryInt("limit", 50))
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	from := c.Query("from")
	to := c.Query("to")
	cursor := c.Query("cursor")

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "GET /v1/finance/journals"),
		attribute.Int("limit", int(limit)),
		attribute.String("from", from),
		attribute.String("to", to),
		attribute.String("cursor", cursor),
	)
	logger.Info().Str("op", op).Int32("limit", limit).Str("from", from).Str("to", to).Msg("")

	result, err := s.svc.ListJournalEntries(ctx, &finance_grpc_adapter.ListJournalEntriesParams{
		From:   from,
		To:     to,
		Limit:  limit,
		Cursor: cursor,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writeFinanceError(c, span, err)
	}

	data := make([]JournalEntryData, 0, len(result.Entries))
	for _, e := range result.Entries {
		lines := make([]JournalLineData, 0, len(e.Lines))
		for _, l := range e.Lines {
			lines = append(lines, JournalLineData{
				ID:          l.ID,
				EntryID:     l.EntryID,
				AccountCode: l.AccountCode,
				Debit:       l.Debit,
				Credit:      l.Credit,
			})
		}
		data = append(data, JournalEntryData{
			ID:             e.ID,
			IdempotencyKey: e.IdempotencyKey,
			SourceType:     e.SourceType,
			SourceID:       e.SourceID,
			PostedAt:       e.PostedAt,
			Description:    e.Description,
			Lines:          lines,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(JournalListResponse{
		Data:       data,
		NextCursor: result.NextCursor,
	})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeFinanceError(c *fiber.Ctx, span trace.Span, err error) error {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	var httpStatus int
	var code, message string

	switch {
	case errors.Is(err, apperrors.ErrValidation):
		httpStatus = fiber.StatusBadRequest
		code = "validation_error"
		message = err.Error()
	case errors.Is(err, apperrors.ErrNotFound):
		httpStatus = fiber.StatusNotFound
		code = "not_found"
		message = err.Error()
	case errors.Is(err, apperrors.ErrUnauthorized):
		httpStatus = fiber.StatusUnauthorized
		code = "unauthorized"
		message = "autentikasi diperlukan"
	case errors.Is(err, apperrors.ErrForbidden):
		httpStatus = fiber.StatusForbidden
		code = "forbidden"
		message = "akses ditolak"
	case errors.Is(err, apperrors.ErrServiceUnavailable):
		httpStatus = fiber.StatusBadGateway
		code = "service_unavailable"
		message = "layanan finance sementara tidak tersedia"
	default:
		httpStatus = fiber.StatusInternalServerError
		code = "internal_error"
		message = "terjadi kesalahan tidak terduga"
	}

	return c.Status(httpStatus).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": message,
		},
	})
}
