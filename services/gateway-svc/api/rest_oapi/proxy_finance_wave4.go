// proxy_finance_wave4.go — gateway REST handlers for Wave 4 finance depth routes
// (BL-FIN-020..041).
//
// Route topology (all bearer-protected):
//   POST /v1/finance/billing/schedule                    → ScheduleBilling
//   POST /v1/finance/bank-transactions                   → RecordBankTransaction
//   GET  /v1/finance/bank-reconciliation                 → GetBankReconciliation
//   GET  /v1/finance/ar-subledger                        → GetARSubledger
//   POST /v1/finance/receipts                            → IssueDigitalReceipt
//   GET  /v1/finance/receipts/:id                        → GetDigitalReceipt
//   POST /v1/finance/manual-payments                     → RecordManualPayment
//   POST /v1/finance/vendors                             → CreateFinanceVendor
//   PUT  /v1/finance/vendors/:id                         → UpdateFinanceVendor
//   GET  /v1/finance/vendors                             → ListFinanceVendors
//   DELETE /v1/finance/vendors/:id                       → DeleteFinanceVendor
//   GET  /v1/finance/ap-subledger                        → GetAPSubledger
//   GET  /v1/finance/payment-authorizations              → ListPendingAuthorizations
//   PUT  /v1/finance/payment-authorizations/:id/decision → DecidePaymentAuthorization
//   POST /v1/finance/petty-cash                          → RecordPettyCash
//   POST /v1/finance/petty-cash/close-period             → ClosePettyCashPeriod
//   GET  /v1/finance/project-costs/:departure_id         → GetProjectCosts
//   GET  /v1/finance/departure-pl/:departure_id          → GetDeparturePL
//   GET  /v1/finance/budget-vs-actual                    → GetBudgetVsActual
//   POST /v1/finance/auto-journal                        → TriggerAutoJournal
//   GET  /v1/finance/revenue-recognition-policy          → GetRevenueRecognitionPolicy
//   PUT  /v1/finance/revenue-recognition-policy          → SetRevenueRecognitionPolicy
//   POST /v1/finance/exchange-rates                      → SetExchangeRate
//   GET  /v1/finance/exchange-rates                      → GetExchangeRate
//   POST /v1/finance/fixed-assets                        → CreateFixedAsset
//   GET  /v1/finance/fixed-assets                        → ListFixedAssets
//   POST /v1/finance/depreciation                        → RunDepreciation
//   POST /v1/finance/tax/calculate                       → CalculateTax
//   GET  /v1/finance/tax/report                          → GetTaxReport
//   POST /v1/finance/commission-payouts                  → CreateCommissionPayout
//   PUT  /v1/finance/commission-payouts/:id/decision     → DecideCommissionPayout
//   GET  /v1/finance/realtime-summary                    → GetRealtimeFinancialSummary
//   GET  /v1/finance/cashflow                            → GetCashFlowDashboard
//   GET  /v1/finance/aging-alerts                        → GetAgingAlerts
//   GET  /v1/finance/audit-log                           → SearchFinanceAuditLog
//
// Per ADR-0009: gateway is the single REST entry-point; finance-svc is pure gRPC.
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

// ---------------------------------------------------------------------------
// BL-FIN-020 ScheduleBilling — POST /v1/finance/billing/schedule
// ---------------------------------------------------------------------------

type ScheduleBillingBody struct {
	DepartureID string `json:"departure_id"`
	DueDate     string `json:"due_date,omitempty"`
	Notes       string `json:"notes,omitempty"`
}

func (s *Server) ScheduleBilling(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ScheduleBilling"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body ScheduleBillingBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.ScheduleBilling(ctx, &finance_grpc_adapter.ScheduleBillingParams{
		DepartureID: body.DepartureID,
		DueDate:     body.DueDate,
		Notes:       body.Notes,
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"invoices_created": result.InvoicesCreated,
			"total_amount":     result.TotalAmount,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-021 RecordBankTransaction — POST /v1/finance/bank-transactions
// ---------------------------------------------------------------------------

type RecordBankTransactionBody struct {
	AccountID   string `json:"account_id"`
	RefNo       string `json:"ref_no"`
	Amount      int64  `json:"amount"`
	TxDate      string `json:"tx_date"`
	Description string `json:"description,omitempty"`
	Direction   string `json:"direction"` // "credit" | "debit"
}

func (s *Server) RecordBankTransaction(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordBankTransaction"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body RecordBankTransactionBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.RecordBankTransaction(ctx, &finance_grpc_adapter.RecordBankTransactionParams{
		AccountID:   body.AccountID,
		RefNo:       body.RefNo,
		Amount:      body.Amount,
		TxDate:      body.TxDate,
		Description: body.Description,
		Direction:   body.Direction,
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{"transaction_id": result.TransactionID},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-021 GetBankReconciliation — GET /v1/finance/bank-reconciliation
// ---------------------------------------------------------------------------

func (s *Server) GetBankReconciliation(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetBankReconciliation"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetBankReconciliation(ctx, &finance_grpc_adapter.GetBankReconciliationParams{
		AccountID: c.Query("account_id"),
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	rows := make([]fiber.Map, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, fiber.Map{
			"tx_id":      r.TxID,
			"ref_no":     r.RefNo,
			"amount":     r.Amount,
			"tx_date":    r.TxDate,
			"direction":  r.Direction,
			"reconciled": r.Reconciled,
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"account_id":      result.AccountID,
			"opening_balance": result.OpeningBalance,
			"closing_balance": result.ClosingBalance,
			"rows":            rows,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-022 GetARSubledger — GET /v1/finance/ar-subledger
// ---------------------------------------------------------------------------

func (s *Server) GetARSubledger(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetARSubledger"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetARSubledger(ctx, &finance_grpc_adapter.GetARSubledgerParams{
		BookingID: c.Query("booking_id"),
		PilgrimID: c.Query("pilgrim_id"),
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	rows := make([]fiber.Map, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, fiber.Map{
			"entry_id":    r.EntryID,
			"date":        r.Date,
			"description": r.Description,
			"debit":       r.Debit,
			"credit":      r.Credit,
			"balance":     r.Balance,
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"booking_id": result.BookingID,
			"pilgrim_id": result.PilgrimID,
			"rows":       rows,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-023 IssueDigitalReceipt — POST /v1/finance/receipts
// ---------------------------------------------------------------------------

type IssueDigitalReceiptBody struct {
	BookingID string `json:"booking_id"`
	PaymentID string `json:"payment_id"`
	Amount    int64  `json:"amount"`
	Notes     string `json:"notes,omitempty"`
}

func (s *Server) IssueDigitalReceipt(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.IssueDigitalReceipt"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body IssueDigitalReceiptBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.IssueDigitalReceipt(ctx, &finance_grpc_adapter.IssueDigitalReceiptParams{
		BookingID: body.BookingID,
		PaymentID: body.PaymentID,
		Amount:    body.Amount,
		Notes:     body.Notes,
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"receipt_id":     result.ReceiptID,
			"receipt_number": result.ReceiptNumber,
			"issued_at":      result.IssuedAt,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-023 GetDigitalReceipt — GET /v1/finance/receipts/:id
// ---------------------------------------------------------------------------

func (s *Server) GetDigitalReceipt(c *fiber.Ctx, receiptID string) error {
	const op = "rest_oapi.Server.GetDigitalReceipt"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("receipt_id", receiptID))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("receipt_id", receiptID).Msg("")

	result, err := s.svc.GetDigitalReceipt(ctx, &finance_grpc_adapter.GetDigitalReceiptParams{ReceiptID: receiptID})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"receipt_id":     result.ReceiptID,
			"receipt_number": result.ReceiptNumber,
			"booking_id":     result.BookingID,
			"payment_id":     result.PaymentID,
			"amount":         result.Amount,
			"issued_at":      result.IssuedAt,
			"notes":          result.Notes,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-024 RecordManualPayment — POST /v1/finance/manual-payments
// ---------------------------------------------------------------------------

type RecordManualPaymentBody struct {
	BookingID   string `json:"booking_id"`
	Amount      int64  `json:"amount"`
	PaymentDate string `json:"payment_date"`
	Method      string `json:"method"`
	EvidenceURL string `json:"evidence_url,omitempty"`
	Notes       string `json:"notes,omitempty"`
}

func (s *Server) RecordManualPayment(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordManualPayment"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body RecordManualPaymentBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.RecordManualPayment(ctx, &finance_grpc_adapter.RecordManualPaymentParams{
		BookingID:   body.BookingID,
		Amount:      body.Amount,
		PaymentDate: body.PaymentDate,
		Method:      body.Method,
		EvidenceURL: body.EvidenceURL,
		Notes:       body.Notes,
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"entry_id":   result.EntryID,
			"journal_id": result.JournalID,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-025 Vendor master CRUD
// ---------------------------------------------------------------------------

type CreateFinanceVendorBody struct {
	Name         string `json:"name"`
	Category     string `json:"category"`
	BankAccount  string `json:"bank_account,omitempty"`
	TaxID        string `json:"tax_id,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
}

func (s *Server) CreateFinanceVendor(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateFinanceVendor"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body CreateFinanceVendorBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.CreateFinanceVendor(ctx, &finance_grpc_adapter.CreateVendorParams{
		Name:         body.Name,
		Category:     body.Category,
		BankAccount:  body.BankAccount,
		TaxID:        body.TaxID,
		ContactEmail: body.ContactEmail,
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{"vendor_id": result.VendorID},
	})
}

type UpdateFinanceVendorBody struct {
	Name         string `json:"name,omitempty"`
	Category     string `json:"category,omitempty"`
	BankAccount  string `json:"bank_account,omitempty"`
	ContactEmail string `json:"contact_email,omitempty"`
}

func (s *Server) UpdateFinanceVendor(c *fiber.Ctx, vendorID string) error {
	const op = "rest_oapi.Server.UpdateFinanceVendor"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("vendor_id", vendorID))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("vendor_id", vendorID).Msg("")

	var body UpdateFinanceVendorBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.UpdateFinanceVendor(ctx, &finance_grpc_adapter.UpdateVendorParams{
		VendorID:     vendorID,
		Name:         body.Name,
		Category:     body.Category,
		BankAccount:  body.BankAccount,
		ContactEmail: body.ContactEmail,
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{"vendor_id": result.VendorID},
	})
}

func (s *Server) ListFinanceVendors(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListFinanceVendors"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.ListFinanceVendors(ctx, &finance_grpc_adapter.ListVendorsParams{
		Category: c.Query("category"),
		PageSize: int32(c.QueryInt("page_size", 50)),
		Cursor:   c.Query("cursor"),
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	vendors := make([]fiber.Map, 0, len(result.Vendors))
	for _, v := range result.Vendors {
		vendors = append(vendors, fiber.Map{
			"vendor_id":     v.VendorID,
			"name":          v.Name,
			"category":      v.Category,
			"bank_account":  v.BankAccount,
			"tax_id":        v.TaxID,
			"contact_email": v.ContactEmail,
			"created_at":    v.CreatedAt,
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":        vendors,
		"next_cursor": result.NextCursor,
	})
}

func (s *Server) DeleteFinanceVendor(c *fiber.Ctx, vendorID string) error {
	const op = "rest_oapi.Server.DeleteFinanceVendor"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("vendor_id", vendorID))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("vendor_id", vendorID).Msg("")

	result, err := s.svc.DeleteFinanceVendor(ctx, &finance_grpc_adapter.DeleteVendorParams{VendorID: vendorID})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{"deleted": result.Deleted},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-026 GetAPSubledger — GET /v1/finance/ap-subledger
// ---------------------------------------------------------------------------

func (s *Server) GetAPSubledger(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetAPSubledger"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetAPSubledger(ctx, &finance_grpc_adapter.GetAPSubledgerParams{
		VendorID:  c.Query("vendor_id"),
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	rows := make([]fiber.Map, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, fiber.Map{
			"entry_id":    r.EntryID,
			"date":        r.Date,
			"description": r.Description,
			"debit":       r.Debit,
			"credit":      r.Credit,
			"balance":     r.Balance,
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"vendor_id": result.VendorID,
			"rows":      rows,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-027 Payment authorizations
// ---------------------------------------------------------------------------

func (s *Server) ListPendingAuthorizations(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListPendingAuthorizations"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.ListPendingAuthorizations(ctx, &finance_grpc_adapter.ListPendingAuthorizationsParams{
		Level: int32(c.QueryInt("level", 0)),
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	items := make([]fiber.Map, 0, len(result.Items))
	for _, it := range result.Items {
		items = append(items, fiber.Map{
			"auth_id":      it.AuthID,
			"batch_id":     it.BatchID,
			"amount":       it.Amount,
			"requested_by": it.RequestedBy,
			"level":        it.Level,
			"created_at":   it.CreatedAt,
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": items})
}

type DecidePaymentAuthorizationBody struct {
	Decision string `json:"decision"` // "approve" | "reject"
	Notes    string `json:"notes,omitempty"`
}

func (s *Server) DecidePaymentAuthorization(c *fiber.Ctx, authID string) error {
	const op = "rest_oapi.Server.DecidePaymentAuthorization"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("auth_id", authID))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("auth_id", authID).Msg("")

	var body DecidePaymentAuthorizationBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.DecidePaymentAuthorization(ctx, &finance_grpc_adapter.DecidePaymentAuthorizationParams{
		AuthID:   authID,
		Decision: body.Decision,
		Notes:    body.Notes,
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{"auth_id": result.AuthID, "status": result.Status},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-028 Petty cash
// ---------------------------------------------------------------------------

type RecordPettyCashBody struct {
	Amount      int64  `json:"amount"`
	Direction   string `json:"direction"` // "in" | "out"
	Description string `json:"description"`
	Category    string `json:"category,omitempty"`
	Date        string `json:"date"`
	EvidenceURL string `json:"evidence_url,omitempty"`
}

func (s *Server) RecordPettyCash(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordPettyCash"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body RecordPettyCashBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.RecordPettyCash(ctx, &finance_grpc_adapter.RecordPettyCashParams{
		Amount:      body.Amount,
		Direction:   body.Direction,
		Description: body.Description,
		Category:    body.Category,
		Date:        body.Date,
		EvidenceURL: body.EvidenceURL,
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"entry_id":        result.EntryID,
			"running_balance": result.RunningBalance,
		},
	})
}

type ClosePettyCashPeriodBody struct {
	PeriodEnd string `json:"period_end"`
	Notes     string `json:"notes,omitempty"`
}

func (s *Server) ClosePettyCashPeriod(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ClosePettyCashPeriod"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body ClosePettyCashPeriodBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.ClosePettyCashPeriod(ctx, &finance_grpc_adapter.ClosePettyCashPeriodParams{
		PeriodEnd: body.PeriodEnd,
		Notes:     body.Notes,
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"closing_entry_id": result.ClosingEntryID,
			"closing_balance":  result.ClosingBalance,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-029 GetProjectCosts — GET /v1/finance/project-costs/:departure_id
// ---------------------------------------------------------------------------

func (s *Server) GetProjectCosts(c *fiber.Ctx, departureID string) error {
	const op = "rest_oapi.Server.GetProjectCosts"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("departure_id", departureID))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("departure_id", departureID).Msg("")

	result, err := s.svc.GetProjectCosts(ctx, &finance_grpc_adapter.GetProjectCostsParams{DepartureID: departureID})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	lines := make([]fiber.Map, 0, len(result.Lines))
	for _, l := range result.Lines {
		lines = append(lines, fiber.Map{
			"category":    l.Category,
			"description": l.Description,
			"amount":      l.Amount,
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"departure_id": result.DepartureID,
			"total_cost":   result.TotalCost,
			"lines":        lines,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-030 GetDeparturePL — GET /v1/finance/departure-pl/:departure_id
// ---------------------------------------------------------------------------

func (s *Server) GetDeparturePL(c *fiber.Ctx, departureID string) error {
	const op = "rest_oapi.Server.GetDeparturePL"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("departure_id", departureID))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("departure_id", departureID).Msg("")

	result, err := s.svc.GetDeparturePL(ctx, &finance_grpc_adapter.GetDeparturePLParams{DepartureID: departureID})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"departure_id":   result.DepartureID,
			"revenue":        result.Revenue,
			"costs":          result.Costs,
			"gross_profit":   result.GrossProfit,
			"budget_revenue": result.BudgetRevenue,
			"budget_costs":   result.BudgetCosts,
			"variance":       result.Variance,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-031 GetBudgetVsActual — GET /v1/finance/budget-vs-actual
// ---------------------------------------------------------------------------

func (s *Server) GetBudgetVsActual(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetBudgetVsActual"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetBudgetVsActual(ctx, &finance_grpc_adapter.GetBudgetVsActualParams{
		StartDate:   c.Query("start_date"),
		EndDate:     c.Query("end_date"),
		DepartureID: c.Query("departure_id"),
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	lines := make([]fiber.Map, 0, len(result.Lines))
	for _, l := range result.Lines {
		lines = append(lines, fiber.Map{
			"account_code": l.AccountCode,
			"account_name": l.AccountName,
			"budgeted":     l.Budgeted,
			"actual":       l.Actual,
			"variance":     l.Variance,
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"start_date": result.StartDate,
			"end_date":   result.EndDate,
			"lines":      lines,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-032 TriggerAutoJournal — POST /v1/finance/auto-journal
// ---------------------------------------------------------------------------

type TriggerAutoJournalBody struct {
	EventKind string `json:"event_kind"`
	SourceID  string `json:"source_id"`
	Amount    int64  `json:"amount"`
	Notes     string `json:"notes,omitempty"`
}

func (s *Server) TriggerAutoJournal(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.TriggerAutoJournal"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body TriggerAutoJournalBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.TriggerAutoJournal(ctx, &finance_grpc_adapter.TriggerAutoJournalParams{
		EventKind: body.EventKind,
		SourceID:  body.SourceID,
		Amount:    body.Amount,
		Notes:     body.Notes,
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"journal_id": result.JournalID,
			"skipped":    result.Skipped,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-033 Revenue recognition policy
// ---------------------------------------------------------------------------

func (s *Server) GetRevenueRecognitionPolicy(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetRevenueRecognitionPolicy"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetRevenueRecognitionPolicy(ctx)
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"trigger_status":      result.TriggerStatus,
			"deferral_account":    result.DeferralAccount,
			"recognition_account": result.RecognitionAccount,
		},
	})
}

type SetRevenueRecognitionPolicyBody struct {
	TriggerStatus      string `json:"trigger_status"`
	DeferralAccount    string `json:"deferral_account,omitempty"`
	RecognitionAccount string `json:"recognition_account,omitempty"`
}

func (s *Server) SetRevenueRecognitionPolicy(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.SetRevenueRecognitionPolicy"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body SetRevenueRecognitionPolicyBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	if err := s.svc.SetRevenueRecognitionPolicy(ctx, &finance_grpc_adapter.SetRevenueRecognitionPolicyParams{
		TriggerStatus:      body.TriggerStatus,
		DeferralAccount:    body.DeferralAccount,
		RecognitionAccount: body.RecognitionAccount,
	}); err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"updated": true}})
}

// ---------------------------------------------------------------------------
// BL-FIN-034 Exchange rates
// ---------------------------------------------------------------------------

type SetExchangeRateBody struct {
	FromCurrency  string `json:"from_currency"`
	ToCurrency    string `json:"to_currency"`
	Rate          string `json:"rate"`
	EffectiveDate string `json:"effective_date"`
}

func (s *Server) SetExchangeRate(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.SetExchangeRate"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body SetExchangeRateBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.SetExchangeRate(ctx, &finance_grpc_adapter.SetExchangeRateParams{
		FromCurrency:  body.FromCurrency,
		ToCurrency:    body.ToCurrency,
		Rate:          body.Rate,
		EffectiveDate: body.EffectiveDate,
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{"rate_id": result.RateID},
	})
}

func (s *Server) GetExchangeRate(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetExchangeRate"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetExchangeRate(ctx, &finance_grpc_adapter.GetExchangeRateParams{
		FromCurrency: c.Query("from"),
		ToCurrency:   c.Query("to"),
		AsOf:         c.Query("as_of"),
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"rate_id":        result.RateID,
			"rate":           result.Rate,
			"effective_date": result.EffectiveDate,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-035 Fixed assets
// ---------------------------------------------------------------------------

type CreateFixedAssetBody struct {
	Name             string `json:"name"`
	Category         string `json:"category"`
	PurchaseDate     string `json:"purchase_date"`
	PurchaseCost     int64  `json:"purchase_cost"`
	UsefulLifeMonths int32  `json:"useful_life_months"`
	ResidualValue    int64  `json:"residual_value,omitempty"`
}

func (s *Server) CreateFixedAsset(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateFixedAsset"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body CreateFixedAssetBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.CreateFixedAsset(ctx, &finance_grpc_adapter.CreateFixedAssetParams{
		Name:             body.Name,
		Category:         body.Category,
		PurchaseDate:     body.PurchaseDate,
		PurchaseCost:     body.PurchaseCost,
		UsefulLifeMonths: body.UsefulLifeMonths,
		ResidualValue:    body.ResidualValue,
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{"asset_id": result.AssetID},
	})
}

func (s *Server) ListFixedAssets(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListFixedAssets"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.ListFixedAssets(ctx, c.Query("category"))
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	assets := make([]fiber.Map, 0, len(result.Assets))
	for _, a := range result.Assets {
		assets = append(assets, fiber.Map{
			"asset_id":                 a.AssetID,
			"name":                     a.Name,
			"category":                 a.Category,
			"purchase_date":            a.PurchaseDate,
			"purchase_cost":            a.PurchaseCost,
			"accumulated_depreciation": a.AccumulatedDepreciation,
			"book_value":               a.BookValue,
			"useful_life_months":       a.UsefulLifeMonths,
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": assets})
}

type RunDepreciationBody struct {
	AsOf string `json:"as_of"`
}

func (s *Server) RunDepreciation(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RunDepreciation"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body RunDepreciationBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.RunDepreciation(ctx, &finance_grpc_adapter.RunDepreciationParams{AsOf: body.AsOf})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"assets_processed":   result.AssetsProcessed,
			"total_depreciation": result.TotalDepreciation,
			"journal_id":         result.JournalID,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-036 Tax
// ---------------------------------------------------------------------------

type CalculateTaxBody struct {
	BaseAmount int64  `json:"base_amount"`
	TaxType    string `json:"tax_type"` // "vat" | "wht"
	Rate       string `json:"rate"`
}

func (s *Server) CalculateTax(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CalculateTax"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body CalculateTaxBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.CalculateTax(ctx, &finance_grpc_adapter.CalculateTaxParams{
		BaseAmount: body.BaseAmount,
		TaxType:    body.TaxType,
		Rate:       body.Rate,
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"base_amount": result.BaseAmount,
			"tax_amount":  result.TaxAmount,
			"tax_type":    result.TaxType,
			"rate":        result.Rate,
		},
	})
}

func (s *Server) GetTaxReport(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetTaxReport"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetTaxReport(ctx, &finance_grpc_adapter.GetTaxReportParams{
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
		TaxType:   c.Query("tax_type"),
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	rows := make([]fiber.Map, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, fiber.Map{
			"date":        r.Date,
			"description": r.Description,
			"base_amount": r.BaseAmount,
			"tax_amount":  r.TaxAmount,
			"tax_type":    r.TaxType,
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"start_date": result.StartDate,
			"end_date":   result.EndDate,
			"total_tax":  result.TotalTax,
			"rows":       rows,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-037 Commission payouts
// ---------------------------------------------------------------------------

type CreateCommissionPayoutBody struct {
	AgentID     string `json:"agent_id"`
	DepartureID string `json:"departure_id"`
	Amount      int64  `json:"amount"`
	BasisAmount int64  `json:"basis_amount,omitempty"`
	RatePercent string `json:"rate_percent,omitempty"`
	Notes       string `json:"notes,omitempty"`
}

func (s *Server) CreateCommissionPayout(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateCommissionPayout"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body CreateCommissionPayoutBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.CreateCommissionPayout(ctx, &finance_grpc_adapter.CreateCommissionPayoutParams{
		AgentID:     body.AgentID,
		DepartureID: body.DepartureID,
		Amount:      body.Amount,
		BasisAmount: body.BasisAmount,
		RatePercent: body.RatePercent,
		Notes:       body.Notes,
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{"payout_id": result.PayoutID, "status": result.Status},
	})
}

type DecideCommissionPayoutBody struct {
	Decision string `json:"decision"` // "approve" | "reject"
	Notes    string `json:"notes,omitempty"`
}

func (s *Server) DecideCommissionPayout(c *fiber.Ctx, payoutID string) error {
	const op = "rest_oapi.Server.DecideCommissionPayout"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("payout_id", payoutID))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("payout_id", payoutID).Msg("")

	var body DecideCommissionPayoutBody
	if err := c.BodyParser(&body); err != nil {
		return writeFinanceWave4Error(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.DecideCommissionPayout(ctx, &finance_grpc_adapter.DecideCommissionPayoutParams{
		PayoutID: payoutID,
		Decision: body.Decision,
		Notes:    body.Notes,
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{"payout_id": result.PayoutID, "status": result.Status},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-038 GetRealtimeFinancialSummary — GET /v1/finance/realtime-summary
// ---------------------------------------------------------------------------

func (s *Server) GetRealtimeFinancialSummary(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetRealtimeFinancialSummary"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetRealtimeFinancialSummary(ctx)
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	accounts := make([]fiber.Map, 0, len(result.Accounts))
	for _, acc := range result.Accounts {
		accounts = append(accounts, fiber.Map{
			"account_code": acc.AccountCode,
			"account_name": acc.AccountName,
			"balance":      acc.Balance,
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"as_of":         result.AsOf,
			"total_revenue": result.TotalRevenue,
			"total_expense": result.TotalExpense,
			"net_income":    result.NetIncome,
			"cash_balance":  result.CashBalance,
			"ar_balance":    result.ARBalance,
			"ap_balance":    result.APBalance,
			"accounts":      accounts,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-039 GetCashFlowDashboard — GET /v1/finance/cashflow
// ---------------------------------------------------------------------------

func (s *Server) GetCashFlowDashboard(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetCashFlowDashboard"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetCashFlowDashboard(ctx, &finance_grpc_adapter.GetCashFlowDashboardParams{
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	lines := make([]fiber.Map, 0, len(result.Lines))
	for _, l := range result.Lines {
		lines = append(lines, fiber.Map{
			"date":            l.Date,
			"description":     l.Description,
			"amount":          l.Amount,
			"running_balance": l.RunningBalance,
			"category":        l.Category,
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"start_date":      result.StartDate,
			"end_date":        result.EndDate,
			"opening_balance": result.OpeningBalance,
			"closing_balance": result.ClosingBalance,
			"operating_net":   result.OperatingNet,
			"investing_net":   result.InvestingNet,
			"financing_net":   result.FinancingNet,
			"lines":           lines,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-040 GetAgingAlerts — GET /v1/finance/aging-alerts
// ---------------------------------------------------------------------------

func (s *Server) GetAgingAlerts(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetAgingAlerts"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetAgingAlerts(ctx, &finance_grpc_adapter.GetAgingAlertsParams{
		LedgerType:    c.Query("ledger_type", "both"),
		ThresholdDays: int32(c.QueryInt("threshold_days", 30)),
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	alerts := make([]fiber.Map, 0, len(result.Alerts))
	for _, al := range result.Alerts {
		alerts = append(alerts, fiber.Map{
			"entity_id":    al.EntityID,
			"entity_name":  al.EntityName,
			"amount":       al.Amount,
			"days_overdue": al.DaysOverdue,
			"ledger_type":  al.LedgerType,
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"alerts":        alerts,
			"total_overdue": result.TotalOverdue,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-FIN-041 SearchFinanceAuditLog — GET /v1/finance/audit-log
// ---------------------------------------------------------------------------

func (s *Server) SearchFinanceAuditLog(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.SearchFinanceAuditLog"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.SearchFinanceAuditLog(ctx, &finance_grpc_adapter.SearchFinanceAuditLogParams{
		UserID:     c.Query("user_id"),
		Action:     c.Query("action"),
		EntityType: c.Query("entity_type"),
		EntityID:   c.Query("entity_id"),
		StartDate:  c.Query("start_date"),
		EndDate:    c.Query("end_date"),
		PageSize:   int32(c.QueryInt("page_size", 50)),
		Cursor:     c.Query("cursor"),
	})
	if err != nil {
		return writeFinanceWave4Error(c, span, err)
	}
	rows := make([]fiber.Map, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, fiber.Map{
			"log_id":      r.LogID,
			"user_id":     r.UserID,
			"action":      r.Action,
			"entity_type": r.EntityType,
			"entity_id":   r.EntityID,
			"diff":        r.Diff,
			"created_at":  r.CreatedAt,
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":        rows,
		"next_cursor": result.NextCursor,
	})
}

// ---------------------------------------------------------------------------
// Error helper
// ---------------------------------------------------------------------------

func writeFinanceWave4Error(c *fiber.Ctx, span trace.Span, err error) error {
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
	case errors.Is(err, apperrors.ErrConflict):
		httpStatus = fiber.StatusConflict
		code = "conflict"
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
