// correction.go — gateway adapter method for finance-svc CorrectJournal RPC
// (BL-FIN-006).
//
// Translates gateway-local params → pb request, forwards via gRPC,
// and translates pb response → adapter-local types. Proto types do not leak
// past this package.

package finance_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/finance_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Adapter-local result types
// ---------------------------------------------------------------------------

// CorrectJournalParams holds inputs for CorrectJournal.
type CorrectJournalParams struct {
	EntryID     string
	Reason      string
	ActorUserID string
}

// CorrectJournalResult is the gateway-local result for CorrectJournal.
type CorrectJournalResult struct {
	CorrectionEntryID string
	OriginalEntryID   string
	Idempotent        bool
}

// ---------------------------------------------------------------------------
// CorrectJournal
// ---------------------------------------------------------------------------

// CorrectJournal posts a reversing counter-entry for an existing journal entry.
// Idempotent on entry_id.
func (a *Adapter) CorrectJournal(ctx context.Context, params *CorrectJournalParams) (*CorrectJournalResult, error) {
	const op = "finance_grpc_adapter.Adapter.CorrectJournal"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("entry_id", params.EntryID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.correctionClient.CorrectJournal(ctx, &pb.CorrectJournalRequest{
		EntryId:     params.EntryID,
		Reason:      params.Reason,
		ActorUserId: params.ActorUserID,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Str("entry_id", params.EntryID).Msg("finance-svc.CorrectJournal failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &CorrectJournalResult{
		CorrectionEntryID: resp.GetCorrectionEntryId(),
		OriginalEntryID:   resp.GetOriginalEntryId(),
		Idempotent:        resp.GetIdempotent(),
	}, nil
}
