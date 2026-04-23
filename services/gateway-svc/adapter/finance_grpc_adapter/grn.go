// grn.go — gateway adapter method for finance-svc OnGRNReceived RPC (S3 Wave 2 / BL-FIN-002).
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

// OnGRNReceivedResult is the gateway-local result for OnGRNReceived.
type OnGRNReceivedResult struct {
	EntryID    string
	Idempotent bool
}

// ---------------------------------------------------------------------------
// OnGRNReceived
// ---------------------------------------------------------------------------

// OnGRNReceived posts a Dr COGS / Cr AP journal entry for a goods receipt note.
// Idempotent on grn_id.
func (a *Adapter) OnGRNReceived(ctx context.Context, grnID, departureID string, amountIdr int64) (*OnGRNReceivedResult, error) {
	const op = "finance_grpc_adapter.Adapter.OnGRNReceived"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("grn_id", grnID),
		attribute.String("departure_id", departureID),
		attribute.Int64("amount_idr", amountIdr),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.grnClient.OnGRNReceived(ctx, &pb.OnGRNReceivedRequest{
		GrnId:       grnID,
		DepartureId: departureID,
		AmountIdr:   amountIdr,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Str("grn_id", grnID).Msg("finance-svc.OnGRNReceived failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &OnGRNReceivedResult{
		EntryID:    resp.GetEntryId(),
		Idempotent: resp.GetIdempotent(),
	}, nil
}
