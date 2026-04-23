// grn.go — OnGRNReceived service-layer implementation (BL-FIN-002).
//
// Called when a Goods Receipt Note is received. Posts a balanced double-entry
// journal:
//
//	Dr 5001 (COGS/Inventory Expense) — amount_idr
//	Cr 2001 (AP/Pilgrim Liability)   — amount_idr
//
// Idempotency: idempotency_key = "grn:" + grn_id. If the key already exists,
// returns the existing entry unchanged (no second posting).
//
// Balance check: debit_total == credit_total before any INSERT.
//
// Acceptance criterion: "GRN fails if AP posting fails" — the transaction
// wraps both lines; if either insert fails the entry is rolled back.

package service

import (
	"context"
	"errors"
	"fmt"

	"finance-svc/store/postgres_store"
	"finance-svc/store/postgres_store/sqlc"
	"finance-svc/util/logging"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// Account codes used for GRN auto-AP posting.
const (
	accountCOGS = "5001" // COGS / Inventory Expense (debit-normal)
	// accountPilgrimLiability = "2001" declared in journal.go
)

// OnGRNReceivedParams holds inputs for OnGRNReceived.
type OnGRNReceivedParams struct {
	GrnID       string
	DepartureID string
	AmountIDR   int64 // integer IDR — no fractional amounts
}

// OnGRNReceivedResult holds the result of OnGRNReceived.
type OnGRNReceivedResult struct {
	EntryID    string
	Idempotent bool // true if an existing entry was returned (replayed)
}

// OnGRNReceived posts Dr 5001 (COGS) / Cr 2001 (AP/Pilgrim Liability) for a
// Goods Receipt Note. Idempotent on idempotency_key = "grn:" + grn_id.
func (svc *Service) OnGRNReceived(ctx context.Context, params *OnGRNReceivedParams) (*OnGRNReceivedResult, error) {
	const op = "service.Service.OnGRNReceived"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("grn_id", params.GrnID),
		attribute.String("departure_id", params.DepartureID),
		attribute.Int64("amount_idr", params.AmountIDR),
	)

	logger := logging.LogWithTrace(ctx, svc.logger)

	if params.GrnID == "" {
		return nil, fmt.Errorf("%s: grn_id is required", op)
	}
	if params.DepartureID == "" {
		return nil, fmt.Errorf("%s: departure_id is required", op)
	}
	if params.AmountIDR <= 0 {
		return nil, fmt.Errorf("%s: amount_idr must be positive, got %d", op, params.AmountIDR)
	}

	idempotencyKey := "grn:" + params.GrnID

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
			Str("grn_id", params.GrnID).
			Str("entry_id", existing.ID).
			Bool("idempotent", true).
			Msg("GRN journal entry already exists, returning existing")

		span.SetStatus(otelCodes.Ok, "replayed")
		return &OnGRNReceivedResult{
			EntryID:    existing.ID,
			Idempotent: true,
		}, nil
	}

	// --- Build and validate lines ---
	amountNumeric := int64ToNumeric(params.AmountIDR)

	// Balance check: debit_total == credit_total.
	debitTotal := params.AmountIDR
	creditTotal := params.AmountIDR
	if debitTotal != creditTotal {
		err := fmt.Errorf("%s: unbalanced journal: debit=%d credit=%d", op, debitTotal, creditTotal)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	description := pgtype.Text{
		String: fmt.Sprintf("Auto-AP from GRN %s (departure %s) — amount %d IDR",
			params.GrnID, params.DepartureID, params.AmountIDR),
		Valid: true,
	}

	// --- Insert entry + lines in a transaction ---
	// If AP posting (Cr 2001) fails, the transaction rolls back, satisfying
	// the acceptance criterion "GRN fails if AP posting fails".
	var result OnGRNReceivedResult

	_, err = svc.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(qtx *sqlc.Queries) error {
			entry, err := qtx.InsertJournalEntry(ctx, sqlc.InsertJournalEntryParams{
				IdempotencyKey: idempotencyKey,
				SourceType:     "grn",
				SourceID:       params.GrnID,
				Description:    description,
			})
			if err != nil {
				return fmt.Errorf("insert journal entry: %w", err)
			}

			// Dr 5001 (COGS/Inventory Expense)
			if _, err := qtx.InsertJournalLine(ctx, sqlc.InsertJournalLineParams{
				EntryID:     entry.ID,
				AccountCode: accountCOGS,
				Debit:       amountNumeric,
				Credit:      zeroNumeric(),
			}); err != nil {
				return fmt.Errorf("insert debit line (5001): %w", err)
			}

			// Cr 2001 (AP/Pilgrim Liability) — this is the AP posting
			if _, err := qtx.InsertJournalLine(ctx, sqlc.InsertJournalLineParams{
				EntryID:     entry.ID,
				AccountCode: accountPilgrimLiability,
				Debit:       zeroNumeric(),
				Credit:      amountNumeric,
			}); err != nil {
				return fmt.Errorf("insert credit line AP (2001): %w", err)
			}

			result = OnGRNReceivedResult{
				EntryID:    entry.ID,
				Idempotent: false,
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
		Str("grn_id", params.GrnID).
		Str("departure_id", params.DepartureID).
		Str("entry_id", result.EntryID).
		Int64("amount_idr", params.AmountIDR).
		Bool("idempotent", false).
		Msg("GRN auto-AP journal posted")

	span.SetStatus(otelCodes.Ok, "posted")
	return &result, nil
}
