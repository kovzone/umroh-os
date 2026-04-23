// proxy_invoice.go — gateway REST handlers for invoice routes (BL-PAY-001 / ISSUE-005).
//
// Route topology (bearer-protected):
//   POST /v1/invoices                        — CreateInvoice (create invoice + VA for a booking)
//   GET  /v1/invoices/:id                    — GetInvoiceByID
//   POST /v1/invoices/:id/virtual-accounts   — IssueVirtualAccountForInvoice (re-issue VA on existing invoice)
//
// Per ADR-0009: gateway is the single REST entry-point; payment-svc is pure gRPC.

package rest_oapi

import (
	"errors"

	"gateway-svc/adapter/payment_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Request / response body types
// ---------------------------------------------------------------------------

// CreateInvoiceBody is the JSON body for POST /v1/invoices.
type CreateInvoiceBody struct {
	BookingID   string  `json:"booking_id"`
	AmountTotal float64 `json:"amount_total"`
	Gateway     string  `json:"gateway"`
	BankCode    string  `json:"bank_code"`
}

// InvoiceVAResponseData is the JSON response for CreateInvoice and IssueVirtualAccountForInvoice.
type InvoiceVAResponseData struct {
	InvoiceID     string  `json:"invoice_id"`
	BookingID     string  `json:"booking_id"`
	AccountNumber string  `json:"account_number"`
	BankCode      string  `json:"bank_code"`
	AmountTotal   float64 `json:"amount_total"`
	ExpiresAt     string  `json:"expires_at"`
	Gateway       string  `json:"gateway"`
	Replayed      bool    `json:"replayed"`
}

// InvoiceResponseData is the JSON response for GetInvoiceByID.
type InvoiceResponseData struct {
	ID          string  `json:"id"`
	BookingID   string  `json:"booking_id"`
	Status      string  `json:"status"`
	AmountTotal float64 `json:"amount_total"`
	PaidAmount  float64 `json:"paid_amount"`
	Currency    string  `json:"currency"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// IssueVAForInvoiceBody is the JSON body for POST /v1/invoices/:id/virtual-accounts.
type IssueVAForInvoiceBody struct {
	Gateway  string `json:"gateway"`
	BankCode string `json:"bank_code"`
}

// ---------------------------------------------------------------------------
// CreateInvoice — POST /v1/invoices (bearer)
// ---------------------------------------------------------------------------

// CreateInvoice creates an invoice + virtual account for a booking.
func (s *Server) CreateInvoice(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateInvoice"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("endpoint", "POST /v1/invoices"))
	logger.Info().Str("op", op).Msg("")

	var body CreateInvoiceBody
	if err := c.BodyParser(&body); err != nil {
		return writePaymentError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	if body.BookingID == "" {
		return writePaymentError(c, span, errors.Join(apperrors.ErrValidation, errors.New("booking_id diperlukan")))
	}

	result, err := s.svc.IssueVirtualAccount(ctx, &payment_grpc_adapter.IssueVirtualAccountParams{
		BookingID:   body.BookingID,
		AmountTotal: body.AmountTotal,
		GatewayPref: body.Gateway,
		BankCode:    body.BankCode,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writePaymentError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	span.SetAttributes(attribute.String("invoice_id", result.InvoiceID))

	httpStatus := fiber.StatusCreated
	if result.Replayed {
		httpStatus = fiber.StatusOK
	}

	return c.Status(httpStatus).JSON(fiber.Map{
		"data": InvoiceVAResponseData{
			InvoiceID:     result.InvoiceID,
			BookingID:     result.BookingID,
			AccountNumber: result.AccountNumber,
			BankCode:      result.BankCode,
			AmountTotal:   result.AmountTotal,
			ExpiresAt:     result.ExpiresAt,
			Gateway:       result.Gateway,
			Replayed:      result.Replayed,
		},
	})
}

// ---------------------------------------------------------------------------
// GetInvoiceByID — GET /v1/invoices/:id (bearer)
// ---------------------------------------------------------------------------

// GetInvoiceByID fetches a single invoice by its UUID.
func (s *Server) GetInvoiceByID(c *fiber.Ctx, invoiceID string) error {
	const op = "rest_oapi.Server.GetInvoiceByID"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("endpoint", "GET /v1/invoices/:id"),
		attribute.String("invoice_id", invoiceID),
	)
	logger.Info().Str("op", op).Str("invoice_id", invoiceID).Msg("")

	if invoiceID == "" {
		return writePaymentError(c, span, errors.Join(apperrors.ErrValidation, errors.New("invoice_id diperlukan")))
	}

	result, err := s.svc.GetInvoiceByID(ctx, &payment_grpc_adapter.GetInvoiceByIDParams{
		InvoiceID: invoiceID,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		return writePaymentError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": InvoiceResponseData{
			ID:          result.ID,
			BookingID:   result.BookingID,
			Status:      result.Status,
			AmountTotal: result.AmountTotal,
			PaidAmount:  result.PaidAmount,
			Currency:    result.Currency,
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
		},
	})
}

// ---------------------------------------------------------------------------
// IssueVirtualAccountForInvoice — POST /v1/invoices/:id/virtual-accounts (bearer)
// ---------------------------------------------------------------------------

// IssueVirtualAccountForInvoice re-issues a virtual account for an existing invoice.
// First fetches the invoice to resolve booking_id + amount, then calls IssueVirtualAccount.
func (s *Server) IssueVirtualAccountForInvoice(c *fiber.Ctx, invoiceID string) error {
	const op = "rest_oapi.Server.IssueVirtualAccountForInvoice"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("endpoint", "POST /v1/invoices/:id/virtual-accounts"),
		attribute.String("invoice_id", invoiceID),
	)
	logger.Info().Str("op", op).Str("invoice_id", invoiceID).Msg("")

	if invoiceID == "" {
		return writePaymentError(c, span, errors.Join(apperrors.ErrValidation, errors.New("invoice_id diperlukan")))
	}

	var body IssueVAForInvoiceBody
	if err := c.BodyParser(&body); err != nil {
		return writePaymentError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	// Step 1: fetch invoice to get booking_id + amount_total.
	invoice, err := s.svc.GetInvoiceByID(ctx, &payment_grpc_adapter.GetInvoiceByIDParams{
		InvoiceID: invoiceID,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("get invoice failed")
		return writePaymentError(c, span, err)
	}

	// Step 2: issue / re-issue VA using booking_id + amount from invoice.
	result, err := s.svc.IssueVirtualAccount(ctx, &payment_grpc_adapter.IssueVirtualAccountParams{
		BookingID:   invoice.BookingID,
		AmountTotal: invoice.AmountTotal,
		GatewayPref: body.Gateway,
		BankCode:    body.BankCode,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("issue VA failed")
		return writePaymentError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	span.SetAttributes(attribute.String("invoice_id", result.InvoiceID))

	httpStatus := fiber.StatusCreated
	if result.Replayed {
		httpStatus = fiber.StatusOK
	}

	return c.Status(httpStatus).JSON(fiber.Map{
		"data": InvoiceVAResponseData{
			InvoiceID:     result.InvoiceID,
			BookingID:     result.BookingID,
			AccountNumber: result.AccountNumber,
			BankCode:      result.BankCode,
			AmountTotal:   result.AmountTotal,
			ExpiresAt:     result.ExpiresAt,
			Gateway:       result.Gateway,
			Replayed:      result.Replayed,
		},
	})
}
