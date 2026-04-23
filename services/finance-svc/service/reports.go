// reports.go — GetFinanceSummary and ListJournalEntries service-layer
// implementations for finance-svc (S5-E-01 / BL-FIN-004..005).
//
// GetFinanceSummary: aggregate debit/credit per account_code from
// finance.journal_lines. net = debit - credit.
//
// ListJournalEntries: paginated list of journal entries with their lines,
// ordered by posted_at DESC. Supports optional from/to date range filter and
// cursor-based pagination.

package service

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"finance-svc/store/postgres_store/sqlc"
	"finance-svc/util/logging"

	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// GetFinanceSummary
// ---------------------------------------------------------------------------

// AccountBalanceResult holds the aggregated balance for one account.
type AccountBalanceResult struct {
	AccountCode string
	DebitTotal  int64 // integer IDR
	CreditTotal int64 // integer IDR
	Net         int64 // DebitTotal - CreditTotal
}

// GetFinanceSummaryResult holds the list of account balances.
type GetFinanceSummaryResult struct {
	Accounts []AccountBalanceResult
}

// GetFinanceSummary returns per-account aggregated debit/credit balances.
func (svc *Service) GetFinanceSummary(ctx context.Context) (*GetFinanceSummaryResult, error) {
	const op = "service.Service.GetFinanceSummary"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, svc.logger)
	logger.Info().Str("op", op).Msg("")

	rows, err := svc.store.GetFinanceSummary(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}

	accounts := make([]AccountBalanceResult, 0, len(rows))
	for _, row := range rows {
		debit := numericToInt64(row.DebitTotal)
		credit := numericToInt64(row.CreditTotal)
		accounts = append(accounts, AccountBalanceResult{
			AccountCode: row.AccountCode,
			DebitTotal:  debit,
			CreditTotal: credit,
			Net:         debit - credit,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &GetFinanceSummaryResult{Accounts: accounts}, nil
}

// ---------------------------------------------------------------------------
// ListJournalEntries
// ---------------------------------------------------------------------------

// JournalLineResult holds one line of a journal entry.
type JournalLineResult struct {
	ID          string
	EntryID     string
	AccountCode string
	Debit       int64 // integer IDR
	Credit      int64 // integer IDR
}

// JournalEntryResult holds one journal entry with its lines.
type JournalEntryResult struct {
	ID             string
	IdempotencyKey string
	SourceType     string
	SourceID       string
	PostedAt       time.Time
	Description    string
	Lines          []JournalLineResult
}

// ListJournalEntriesParams holds filter + pagination inputs.
// From and To are optional; zero value means no bound.
// Cursor is the PostedAt of the last seen entry (exclusive upper bound for
// next page, since we order DESC). Empty = first page.
type ListJournalEntriesParams struct {
	From   *time.Time
	To     *time.Time
	Limit  int32
	Cursor *time.Time
}

// ListJournalEntriesResult holds the paginated result.
type ListJournalEntriesResult struct {
	Entries    []JournalEntryResult
	NextCursor string // RFC3339 of last entry's PostedAt; empty = no more pages
}

const defaultJournalLimit = 50
const maxJournalLimit = 200

// ListJournalEntries returns a cursor-paginated list of journal entries with lines.
func (svc *Service) ListJournalEntries(ctx context.Context, params *ListJournalEntriesParams) (*ListJournalEntriesResult, error) {
	const op = "service.Service.ListJournalEntries"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, svc.logger)

	limit := params.Limit
	if limit <= 0 {
		limit = defaultJournalLimit
	}
	if limit > maxJournalLimit {
		limit = maxJournalLimit
	}

	span.SetAttributes(attribute.Int("limit", int(limit)))

	// Build pgtype filter params.
	var from, to, cursor pgtype.Timestamptz
	if params.From != nil {
		from = pgtype.Timestamptz{Time: *params.From, Valid: true}
	}
	if params.To != nil {
		to = pgtype.Timestamptz{Time: *params.To, Valid: true}
	}
	if params.Cursor != nil {
		cursor = pgtype.Timestamptz{Time: *params.Cursor, Valid: true}
	}

	entries, err := svc.store.ListJournalEntries(ctx, sqlc.ListJournalEntriesParams{
		From:   from,
		To:     to,
		Cursor: cursor,
		Limit:  limit,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: list entries failed: %w", op, err)
	}

	if len(entries) == 0 {
		span.SetStatus(otelCodes.Ok, "empty")
		return &ListJournalEntriesResult{Entries: []JournalEntryResult{}}, nil
	}

	// Collect entry IDs for batch line fetch.
	entryIDs := make([]string, len(entries))
	for i, e := range entries {
		entryIDs[i] = e.ID
	}

	lines, err := svc.store.GetJournalLines(ctx, entryIDs)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: get journal lines failed: %w", op, err)
	}

	// Group lines by entry ID.
	linesByEntry := make(map[string][]JournalLineResult, len(entries))
	for _, l := range lines {
		linesByEntry[l.EntryID] = append(linesByEntry[l.EntryID], JournalLineResult{
			ID:          l.ID,
			EntryID:     l.EntryID,
			AccountCode: l.AccountCode,
			Debit:       numericToInt64(l.Debit),
			Credit:      numericToInt64(l.Credit),
		})
	}

	results := make([]JournalEntryResult, 0, len(entries))
	for _, e := range entries {
		var postedAt time.Time
		if e.PostedAt.Valid {
			postedAt = e.PostedAt.Time
		}
		desc := ""
		if e.Description.Valid {
			desc = e.Description.String
		}
		results = append(results, JournalEntryResult{
			ID:             e.ID,
			IdempotencyKey: e.IdempotencyKey,
			SourceType:     e.SourceType,
			SourceID:       e.SourceID,
			PostedAt:       postedAt,
			Description:    desc,
			Lines:          linesByEntry[e.ID],
		})
	}

	// Determine next cursor: PostedAt of the last entry in this page.
	var nextCursor string
	if len(results) > 0 {
		last := results[len(results)-1]
		nextCursor = last.PostedAt.UTC().Format(time.RFC3339Nano)
	}

	logger.Info().
		Str("op", op).
		Int("count", len(results)).
		Str("next_cursor", nextCursor).
		Msg("")

	span.SetStatus(otelCodes.Ok, "ok")
	return &ListJournalEntriesResult{
		Entries:    results,
		NextCursor: nextCursor,
	}, nil
}

// ---------------------------------------------------------------------------
// Numeric helper
// ---------------------------------------------------------------------------

// numericToInt64 converts a pgtype.Numeric (NUMERIC(15,2) column stored as
// integer×100) back to int64 IDR. Mirrors the reverse of int64ToNumeric in
// journal.go.
func numericToInt64(n pgtype.Numeric) int64 {
	if !n.Valid || n.Int == nil {
		return 0
	}
	// The value is stored as Int × 10^Exp. For our columns Exp = -2 (cents).
	// To recover the IDR integer: value = Int × 10^Exp = Int / 100.
	result := new(big.Int).Set(n.Int)
	exp := int(n.Exp)
	if exp < 0 {
		divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-exp)), nil)
		result.Div(result, divisor)
	} else if exp > 0 {
		multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(exp)), nil)
		result.Mul(result, multiplier)
	}
	return result.Int64()
}
