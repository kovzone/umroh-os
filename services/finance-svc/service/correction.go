// correction.go — CorrectJournal service logic (BL-FIN-006).
//
// Finance audit trail & anti-delete:
//   - CorrectJournal:    posts a reversing counter-entry for an existing journal
//                        entry. The original entry is never modified or deleted.
//   - DeleteJournalEntry: explicitly forbidden; returns ErrForbidden to callers.
//
// Accounting logic:
//   The reversing entry has the same amounts as the original but with Dr and Cr
//   swapped on every line. This zeroes out the net effect of the original entry
//   in the ledger while preserving both entries in the audit trail.
//
// Idempotency:
//   idempotency_key = "correction:" + original_entry_id
//   If this key already exists, the existing correction entry is returned.

package service

import (
	"context"
	"errors"
	"fmt"

	"finance-svc/store/postgres_store"
	"finance-svc/store/postgres_store/sqlc"
	"finance-svc/util/apperrors"
	"finance-svc/util/logging"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// CorrectJournalParams holds the input for CorrectJournal.
type CorrectJournalParams struct {
	EntryID     string // UUID of the journal entry to reverse
	Reason      string // human-readable audit reason
	ActorUserID string // IAM user triggering this
}

// CorrectJournalResult holds the result of CorrectJournal.
type CorrectJournalResult struct {
	CorrectionEntryID string
	OriginalEntryID   string
	// Idempotent is true when the correction already existed.
	Idempotent bool
}

// CorrectJournal posts a reversing counter-entry for an existing journal entry.
// The original entry is left intact; only the new reversal is inserted.
func (svc *Service) CorrectJournal(ctx context.Context, params *CorrectJournalParams) (*CorrectJournalResult, error) {
	const op = "service.Service.CorrectJournal"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("entry_id", params.EntryID),
		attribute.String("actor_user_id", params.ActorUserID),
	)
	logger := logging.LogWithTrace(ctx, svc.logger)

	if params.EntryID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("entry_id is required"))
	}

	correctionKey := "correction:" + params.EntryID

	// --- Idempotency check ---
	existing, err := svc.store.GetJournalEntryByIdempotencyKey(ctx, correctionKey)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: check existing correction: %w", op, err)
	}
	if err == nil {
		logger.Info().
			Str("op", op).
			Str("entry_id", params.EntryID).
			Str("correction_id", existing.ID).
			Bool("idempotent", true).
			Msg("correction already exists, returning existing")
		span.SetStatus(otelCodes.Ok, "idempotent")
		return &CorrectJournalResult{
			CorrectionEntryID: existing.ID,
			OriginalEntryID:   params.EntryID,
			Idempotent:        true,
		}, nil
	}

	// --- Load original entry ---
	original, err := svc.store.GetJournalEntryByID(ctx, params.EntryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("journal entry %s not found", params.EntryID))
		}
		return nil, fmt.Errorf("%s: get original entry: %w", op, err)
	}

	// --- Load original lines ---
	lines, err := svc.store.GetJournalLines(ctx, []string{params.EntryID})
	if err != nil {
		return nil, fmt.Errorf("%s: get journal lines: %w", op, err)
	}
	if len(lines) == 0 {
		return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("journal entry %s has no lines", params.EntryID))
	}

	// --- Build reversal description ---
	reason := params.Reason
	if reason == "" {
		reason = "correction"
	}
	descStr := fmt.Sprintf("Correction of entry %s (%s): %s — original source: %s/%s",
		params.EntryID, original.SourceType, reason, original.SourceType, original.SourceID)
	description := pgtype.Text{String: descStr, Valid: true}

	// --- Insert reversal entry + swapped lines in a transaction ---
	var correctionID string
	_, txErr := svc.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(qtx *sqlc.Queries) error {
			correctionEntry, err := qtx.InsertJournalEntry(ctx, sqlc.InsertJournalEntryParams{
				IdempotencyKey: correctionKey,
				SourceType:     "correction",
				SourceID:       params.EntryID,
				Description:    description,
			})
			if err != nil {
				return fmt.Errorf("insert correction entry: %w", err)
			}
			correctionID = correctionEntry.ID

			// Insert each line with Dr/Cr swapped.
			for _, line := range lines {
				if _, err := qtx.InsertJournalLine(ctx, sqlc.InsertJournalLineParams{
					EntryID:     correctionEntry.ID,
					AccountCode: line.AccountCode,
					Debit:       line.Credit, // swap
					Credit:      line.Debit,  // swap
				}); err != nil {
					return fmt.Errorf("insert reversed line for account %s: %w", line.AccountCode, err)
				}
			}
			return nil
		},
	})
	if txErr != nil {
		span.RecordError(txErr)
		span.SetStatus(otelCodes.Error, txErr.Error())
		return nil, fmt.Errorf("%s: transaction failed: %w", op, txErr)
	}

	logger.Info().
		Str("op", op).
		Str("original_entry_id", params.EntryID).
		Str("correction_entry_id", correctionID).
		Int("lines_reversed", len(lines)).
		Msg("correction entry posted")

	span.SetStatus(otelCodes.Ok, "corrected")
	return &CorrectJournalResult{
		CorrectionEntryID: correctionID,
		OriginalEntryID:   params.EntryID,
		Idempotent:        false,
	}, nil
}

// DeleteJournalEntry is explicitly forbidden. Journal entries are immutable;
// corrections must be made via CorrectJournal (reversing counter-entry).
// This method exists solely to return a clear error to callers that attempt
// a direct delete.
func (svc *Service) DeleteJournalEntry(_ context.Context, _ string) error {
	return errors.Join(apperrors.ErrForbidden,
		fmt.Errorf("journal entries cannot be deleted; use CorrectJournal to post a reversing counter-entry"))
}
