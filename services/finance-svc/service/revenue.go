// revenue.go — RecognizeRevenue service-layer implementation for finance-svc.
//
// Wave 1B / BL-FIN-006.
//
// Called when a departure status transitions to `departed` or `completed`.
// Posts a balanced double-entry journal:
//
//	Dr 2001 (Pilgrim Liability) — total_amount_idr
//	Cr 4001 (Revenue)           — total_amount_idr
//
// Idempotency: idempotency_key = "revenue:" + departure_id. If the key already
// exists, returns the existing entry unchanged (no second posting).
//
// Balance check: debit_total == credit_total before any INSERT.

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"finance-svc/store/postgres_store"
	"finance-svc/store/postgres_store/sqlc"
	"finance-svc/util/logging"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// Account codes used for revenue recognition.
const (
	accountRevenue = "4001" // Pendapatan Paket (revenue, credit-normal)
	// accountPilgrimLiability = "2001" is already declared in journal.go
)

// RecognizeRevenueParams holds inputs for RecognizeRevenue.
type RecognizeRevenueParams struct {
	DepartureID     string
	TotalAmountIDR  int64 // integer IDR — no fractional amounts
}

// RecognizeRevenueResult holds the result of a RecognizeRevenue call.
type RecognizeRevenueResult struct {
	EntryID      string
	RecognizedAt time.Time
	Replayed     bool
}

// RecognizeRevenue posts a double-entry journal entry for revenue recognition:
//
//	Dr 2001 (Pilgrim Liability) — total_amount_idr
//	Cr 4001 (Revenue)           — total_amount_idr
//
// Idempotent on departure_id: if already posted, returns existing entry.
func (svc *Service) RecognizeRevenue(ctx context.Context, params *RecognizeRevenueParams) (*RecognizeRevenueResult, error) {
	const op = "service.Service.RecognizeRevenue"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("departure_id", params.DepartureID),
		attribute.Int64("total_amount_idr", params.TotalAmountIDR),
	)

	logger := logging.LogWithTrace(ctx, svc.logger)

	if params.DepartureID == "" {
		return nil, fmt.Errorf("%s: departure_id is required", op)
	}
	if params.TotalAmountIDR <= 0 {
		return nil, fmt.Errorf("%s: total_amount_idr must be positive, got %d", op, params.TotalAmountIDR)
	}

	idempotencyKey := "revenue:" + params.DepartureID

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
			Str("departure_id", params.DepartureID).
			Str("entry_id", existing.ID).
			Bool("replayed", true).
			Msg("revenue journal entry already exists, returning existing")

		recognizedAt := existing.PostedAt.Time
		if !existing.PostedAt.Valid {
			recognizedAt = time.Now()
		}

		span.SetStatus(otelCodes.Ok, "replayed")
		return &RecognizeRevenueResult{
			EntryID:      existing.ID,
			RecognizedAt: recognizedAt,
			Replayed:     true,
		}, nil
	}

	// --- Build and validate lines ---
	amountNumeric := int64ToNumeric(params.TotalAmountIDR)

	// Balance check: debit_total == credit_total.
	debitTotal := params.TotalAmountIDR
	creditTotal := params.TotalAmountIDR
	if debitTotal != creditTotal {
		err := fmt.Errorf("%s: unbalanced journal: debit=%d credit=%d", op, debitTotal, creditTotal)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	description := pgtype.Text{
		String: fmt.Sprintf("Revenue recognized for departure %s — amount %d IDR",
			params.DepartureID, params.TotalAmountIDR),
		Valid: true,
	}

	// --- Insert entry + lines in a transaction ---
	var result RecognizeRevenueResult

	_, err = svc.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(qtx *sqlc.Queries) error {
			entry, err := qtx.InsertJournalEntry(ctx, sqlc.InsertJournalEntryParams{
				IdempotencyKey: idempotencyKey,
				SourceType:     "departure",
				SourceID:       params.DepartureID,
				Description:    description,
			})
			if err != nil {
				return fmt.Errorf("insert journal entry: %w", err)
			}

			// Dr 2001 (Pilgrim Liability) — revenue recognized reduces the liability
			if _, err := qtx.InsertJournalLine(ctx, sqlc.InsertJournalLineParams{
				EntryID:     entry.ID,
				AccountCode: accountPilgrimLiability,
				Debit:       amountNumeric,
				Credit:      zeroNumeric(),
			}); err != nil {
				return fmt.Errorf("insert debit line (2001): %w", err)
			}

			// Cr 4001 (Revenue)
			if _, err := qtx.InsertJournalLine(ctx, sqlc.InsertJournalLineParams{
				EntryID:     entry.ID,
				AccountCode: accountRevenue,
				Debit:       zeroNumeric(),
				Credit:      amountNumeric,
			}); err != nil {
				return fmt.Errorf("insert credit line (4001): %w", err)
			}

			recognizedAt := entry.PostedAt.Time
			if !entry.PostedAt.Valid {
				recognizedAt = time.Now()
			}

			result = RecognizeRevenueResult{
				EntryID:      entry.ID,
				RecognizedAt: recognizedAt,
				Replayed:     false,
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
		Str("departure_id", params.DepartureID).
		Str("entry_id", result.EntryID).
		Int64("total_amount_idr", params.TotalAmountIDR).
		Bool("replayed", false).
		Msg("revenue recognition journal posted")

	span.SetStatus(otelCodes.Ok, "posted")
	return &result, nil
}
