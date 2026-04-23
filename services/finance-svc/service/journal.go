// journal.go — OnPaymentReceived service-layer implementation for finance-svc.
//
// S3-E-03 / BL-FIN-001..003.
//
// Called by booking-svc (via gRPC) after a booking transitions to
// paid_in_full.  Posts a balanced double-entry journal:
//
//	Dr 1001 (Bank)              — amount
//	Cr 2001 (Pilgrim Liability) — amount
//
// Idempotency: idempotency_key = "payment:" + invoice_id. If the key already
// exists, returns the existing entry unchanged (no second posting).
//
// Balance check: debit_total == credit_total before any INSERT.

package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"finance-svc/store/postgres_store"
	"finance-svc/store/postgres_store/sqlc"
	"finance-svc/util/logging"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// Account codes for the payment received entry.
const (
	accountBank             = "1001" // Bank / Cash asset account
	accountPilgrimLiability = "2001" // Hutang Jamaah (Pilgrim Liability)
)

// OnPaymentReceivedParams holds the inputs for posting a payment journal.
// Amount is int64 (integer IDR) per §S3-J-03 contract.
type OnPaymentReceivedParams struct {
	BookingID  string
	InvoiceID  string
	Amount     int64     // integer IDR — no fractional amounts
	ReceivedAt time.Time // zero value = use server time
}

// OnPaymentReceivedResult holds the result of an OnPaymentReceived call.
type OnPaymentReceivedResult struct {
	EntryID  string
	Balanced bool
	// Replayed is true when an existing entry was returned without a new insert.
	Replayed bool
}

// OnPaymentReceived posts a double-entry journal for a payment received event.
// Idempotent: if an entry already exists for this invoice_id, returns it.
func (svc *Service) OnPaymentReceived(ctx context.Context, params *OnPaymentReceivedParams) (*OnPaymentReceivedResult, error) {
	const op = "service.Service.OnPaymentReceived"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("booking_id", params.BookingID),
		attribute.String("invoice_id", params.InvoiceID),
		attribute.Int64("amount", params.Amount),
	)

	logger := logging.LogWithTrace(ctx, svc.logger)

	if params.BookingID == "" {
		return nil, fmt.Errorf("%s: booking_id is required", op)
	}
	if params.InvoiceID == "" {
		return nil, fmt.Errorf("%s: invoice_id is required", op)
	}
	if params.Amount <= 0 {
		return nil, fmt.Errorf("%s: amount must be positive, got %d", op, params.Amount)
	}

	idempotencyKey := "payment:" + params.InvoiceID

	// --- Idempotency check ---
	existing, err := svc.store.GetJournalEntryByIdempotencyKey(ctx, idempotencyKey)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: get existing entry: %w", op, err)
	}
	if err == nil {
		logger.Info().
			Str("op", op).
			Str("invoice_id", params.InvoiceID).
			Str("entry_id", existing.ID).
			Bool("replayed", true).
			Msg("journal entry already exists, returning existing")

		span.SetStatus(otelCodes.Ok, "replayed")
		return &OnPaymentReceivedResult{
			EntryID:  existing.ID,
			Balanced: true,
			Replayed: true,
		}, nil
	}

	// --- Build and validate lines ---
	amountNumeric := int64ToNumeric(params.Amount)

	// Balanced check: debit_total == credit_total.
	// Both lines use the same amount so balance is guaranteed, but we verify
	// explicitly per the spec to catch future multi-line extensions early.
	debitTotal := params.Amount
	creditTotal := params.Amount
	if debitTotal != creditTotal {
		err := fmt.Errorf("%s: unbalanced journal: debit=%d credit=%d", op, debitTotal, creditTotal)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	// --- Resolve description ---
	receivedAt := params.ReceivedAt
	if receivedAt.IsZero() {
		receivedAt = time.Now()
	}
	description := pgtype.Text{
		String: fmt.Sprintf("Payment received for booking %s (invoice %s) at %s",
			params.BookingID, params.InvoiceID, receivedAt.Format(time.RFC3339)),
		Valid: true,
	}

	// --- Insert entry + lines in a transaction ---
	// We use WithTx to ensure the entry and both lines are atomic.
	var result OnPaymentReceivedResult

	_, err = svc.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(qtx *sqlc.Queries) error {
			entry, err := qtx.InsertJournalEntry(ctx, sqlc.InsertJournalEntryParams{
				IdempotencyKey: idempotencyKey,
				SourceType:     "payment",
				SourceID:       params.InvoiceID,
				Description:    description,
			})
			if err != nil {
				return fmt.Errorf("insert journal entry: %w", err)
			}

			// Dr Bank (1001)
			if _, err := qtx.InsertJournalLine(ctx, sqlc.InsertJournalLineParams{
				EntryID:     entry.ID,
				AccountCode: accountBank,
				Debit:       amountNumeric,
				Credit:      zeroNumeric(),
			}); err != nil {
				return fmt.Errorf("insert debit line: %w", err)
			}

			// Cr Pilgrim Liability (2001)
			if _, err := qtx.InsertJournalLine(ctx, sqlc.InsertJournalLineParams{
				EntryID:     entry.ID,
				AccountCode: accountPilgrimLiability,
				Debit:       zeroNumeric(),
				Credit:      amountNumeric,
			}); err != nil {
				return fmt.Errorf("insert credit line: %w", err)
			}

			result = OnPaymentReceivedResult{
				EntryID:  entry.ID,
				Balanced: true,
				Replayed: false,
			}
			return nil
		},
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: transaction failed: %w", op, err)
	}

	logger.Info().
		Str("op", op).
		Str("booking_id", params.BookingID).
		Str("invoice_id", params.InvoiceID).
		Str("entry_id", result.EntryID).
		Int64("amount", params.Amount).
		Bool("replayed", false).
		Msg("journal entry posted")

	span.SetStatus(otelCodes.Ok, "posted")
	return &result, nil
}

// ---------------------------------------------------------------------------
// Numeric helpers (pgtype.Numeric from int64)
// ---------------------------------------------------------------------------

// int64ToNumeric converts an int64 IDR amount to a pgtype.Numeric value
// suitable for NUMERIC(15,2) columns. Since IDR has no fractional part,
// we store as integer rupiah (Exp = 0) but the column accepts it fine.
// e.g. 15400000 → Int = 1540000000, Exp = -2 (= 15400000.00)
func int64ToNumeric(v int64) pgtype.Numeric {
	// Scale by 100 to represent .00 cents (NUMERIC(15,2) convention).
	bi := big.NewInt(v * 100)
	return pgtype.Numeric{
		Int:   bi,
		Exp:   -2,
		Valid: true,
	}
}

// zeroNumeric returns a pgtype.Numeric representing 0.00.
func zeroNumeric() pgtype.Numeric {
	return pgtype.Numeric{
		Int:   big.NewInt(0),
		Exp:   -2,
		Valid: true,
	}
}
