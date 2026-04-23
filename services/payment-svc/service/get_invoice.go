// get_invoice.go — GetInvoiceByID service method (BL-PAY-001 / ISSUE-005).
//
// Fetches a single invoice by its UUID. Called by gateway-svc for:
//   GET /v1/invoices/:id
//   POST /v1/invoices/:id/virtual-accounts (to get booking_id + amount)

package service

import (
	"context"
	"errors"
	"fmt"

	"payment-svc/store/postgres_store"
	"payment-svc/util/apperrors"
	"payment-svc/util/logging"

	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/codes"
)

// GetInvoiceByIDParams is the input for GetInvoiceByID.
type GetInvoiceByIDParams struct {
	InvoiceID string // UUID string
}

// GetInvoiceByIDResult is the output for GetInvoiceByID.
type GetInvoiceByIDResult struct {
	ID          string
	BookingID   string
	Status      string
	AmountTotal float64
	PaidAmount  float64
	Currency    string
	CreatedAt   string // RFC 3339
	UpdatedAt   string // RFC 3339
}

// GetInvoiceByID fetches a single invoice by UUID.
// Returns ErrNotFound if the invoice does not exist.
// Returns ErrValidation if the invoice_id is not a valid UUID.
func (s *PaymentService) GetInvoiceByID(ctx context.Context, params *GetInvoiceByIDParams) (*GetInvoiceByIDResult, error) {
	const op = "service.PaymentService.GetInvoiceByID"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("invoice_id", params.InvoiceID).Msg("")

	if params.InvoiceID == "" {
		err := errors.Join(apperrors.ErrValidation, fmt.Errorf("invoice_id required"))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var invoiceUUID pgtype.UUID
	if err := invoiceUUID.Scan(params.InvoiceID); err != nil {
		valErr := errors.Join(apperrors.ErrValidation, fmt.Errorf("invalid invoice_id: %w", err))
		span.RecordError(valErr)
		span.SetStatus(codes.Error, valErr.Error())
		return nil, valErr
	}

	invoice, err := s.store.GetInvoiceByID(ctx, invoiceUUID)
	if err != nil {
		wrapped := postgres_store.WrapDBError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &GetInvoiceByIDResult{
		ID:          uuidToString(invoice.ID),
		BookingID:   uuidToString(invoice.BookingID),
		Status:      invoice.Status,
		AmountTotal: numericToFloat(invoice.AmountTotal),
		PaidAmount:  numericToFloat(invoice.PaidAmount),
		Currency:    invoice.Currency,
		CreatedAt:   invoice.CreatedAt.Time.UTC().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   invoice.UpdatedAt.Time.UTC().Format("2006-01-02T15:04:05Z"),
	}, nil
}
