// Package finance_grpc_adapter is gateway-svc's consumer-side wrapper around
// finance-svc's gRPC report surface (S5-E-01 / ADR-0009).
//
// Per ADR-0009, all finance REST report routes proxy to finance-svc over gRPC
// via this adapter:
//   GET /v1/finance/summary   → GetFinanceSummary  (bearer)
//   GET /v1/finance/journals  → ListJournalEntries (bearer)
package finance_grpc_adapter

import (
	"context"
	"errors"
	"fmt"

	"gateway-svc/adapter/finance_grpc_adapter/pb"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Adapter wraps finance-svc's FinanceReportsClient, FinanceDepthClient, FinanceGRNClient,
// FinanceCorrectionClient, and FinanceDisbursementClient.
type Adapter struct {
	logger              *zerolog.Logger
	tracer              trace.Tracer
	financeClient       pb.FinanceReportsClient
	financeDepthClient  pb.FinanceDepthClient
	grnClient           pb.FinanceGRNClient
	correctionClient    pb.FinanceCorrectionClient
	disbursementClient  pb.FinanceDisbursementClient
}

// NewAdapter creates a finance gRPC adapter from an already-dialled conn.
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:             logger,
		tracer:             tracer,
		financeClient:      pb.NewFinanceReportsClient(cc),
		financeDepthClient: pb.NewFinanceDepthClient(cc),
		grnClient:          pb.NewFinanceGRNClient(cc),
		correctionClient:   pb.NewFinanceCorrectionClient(cc),
		disbursementClient: pb.NewFinanceDisbursementClient(cc),
	}
}

// ---------------------------------------------------------------------------
// GetFinanceSummary
// ---------------------------------------------------------------------------

// AccountBalanceResult holds one account's aggregated balance.
type AccountBalanceResult struct {
	AccountCode string
	DebitTotal  int64
	CreditTotal int64
	Net         int64
}

// GetFinanceSummaryResult holds the list of account balances.
type GetFinanceSummaryResult struct {
	Accounts []*AccountBalanceResult
}

func (a *Adapter) GetFinanceSummary(ctx context.Context) (*GetFinanceSummaryResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetFinanceSummary"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.financeClient.GetFinanceSummary(ctx, &pb.GetFinanceSummaryRequest{})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetFinanceSummary failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	accounts := make([]*AccountBalanceResult, 0, len(resp.GetAccounts()))
	for _, acct := range resp.GetAccounts() {
		accounts = append(accounts, &AccountBalanceResult{
			AccountCode: acct.GetAccountCode(),
			DebitTotal:  acct.GetDebitTotal(),
			CreditTotal: acct.GetCreditTotal(),
			Net:         acct.GetNet(),
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &GetFinanceSummaryResult{Accounts: accounts}, nil
}

// ---------------------------------------------------------------------------
// ListJournalEntries
// ---------------------------------------------------------------------------

// ListJournalEntriesParams holds filter + pagination inputs.
type ListJournalEntriesParams struct {
	From   string // optional RFC3339
	To     string // optional RFC3339
	Limit  int32
	Cursor string // optional RFC3339 cursor
}

// JournalLineResult holds one journal line.
type JournalLineResult struct {
	ID          string
	EntryID     string
	AccountCode string
	Debit       int64
	Credit      int64
}

// JournalEntryResult holds one journal entry with its lines.
type JournalEntryResult struct {
	ID             string
	IdempotencyKey string
	SourceType     string
	SourceID       string
	PostedAt       string
	Description    string
	Lines          []JournalLineResult
}

// ListJournalEntriesResult holds the paginated result.
type ListJournalEntriesResult struct {
	Entries    []JournalEntryResult
	NextCursor string
}

func (a *Adapter) ListJournalEntries(ctx context.Context, params *ListJournalEntriesParams) (*ListJournalEntriesResult, error) {
	const op = "finance_grpc_adapter.Adapter.ListJournalEntries"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.financeClient.ListJournalEntries(ctx, &pb.ListJournalEntriesRequest{
		From:   params.From,
		To:     params.To,
		Limit:  params.Limit,
		Cursor: params.Cursor,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.ListJournalEntries failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	entries := make([]JournalEntryResult, 0, len(resp.GetEntries()))
	for _, e := range resp.GetEntries() {
		lines := make([]JournalLineResult, 0, len(e.GetLines()))
		for _, l := range e.GetLines() {
			lines = append(lines, JournalLineResult{
				ID:          l.GetId(),
				EntryID:     l.GetEntryId(),
				AccountCode: l.GetAccountCode(),
				Debit:       l.GetDebit(),
				Credit:      l.GetCredit(),
			})
		}
		entries = append(entries, JournalEntryResult{
			ID:             e.GetId(),
			IdempotencyKey: e.GetIdempotencyKey(),
			SourceType:     e.GetSourceType(),
			SourceID:       e.GetSourceId(),
			PostedAt:       e.GetPostedAt(),
			Description:    e.GetDescription(),
			Lines:          lines,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &ListJournalEntriesResult{
		Entries:    entries,
		NextCursor: resp.GetNextCursor(),
	}, nil
}

// ---------------------------------------------------------------------------
// Error mapping helper
// ---------------------------------------------------------------------------

func mapFinanceError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("finance call failed: %w", err))
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return errors.Join(apperrors.ErrNotFound, errors.New(st.Message()))
	case grpcCodes.InvalidArgument:
		return errors.Join(apperrors.ErrValidation, errors.New(st.Message()))
	case grpcCodes.Unauthenticated:
		return errors.Join(apperrors.ErrUnauthorized, errors.New(st.Message()))
	case grpcCodes.PermissionDenied:
		return errors.Join(apperrors.ErrForbidden, errors.New(st.Message()))
	case grpcCodes.Unavailable:
		return errors.Join(apperrors.ErrServiceUnavailable, errors.New(st.Message()))
	default:
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("finance call failed: %s", st.Message()))
	}
}
