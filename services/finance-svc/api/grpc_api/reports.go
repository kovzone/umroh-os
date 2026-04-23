// reports.go — gRPC handlers for GetFinanceSummary and ListJournalEntries
// RPCs (S5-E-01 / BL-FIN-004..005).
//
// Called by gateway-svc to populate:
//   GET /v1/finance/summary   — aggregate per-account balances
//   GET /v1/finance/journals  — paginated journal entries + lines

package grpc_api

import (
	"context"
	"time"

	"finance-svc/api/grpc_api/pb"
	"finance-svc/service"
	"finance-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// GetFinanceSummary handles the GetFinanceSummary RPC.
// Returns aggregated debit/credit per account_code.
func (s *Server) GetFinanceSummary(ctx context.Context, _ *pb.GetFinanceSummaryRequest) (*pb.GetFinanceSummaryResponse, error) {
	const op = "grpc_api.Server.GetFinanceSummary"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetFinanceSummary(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "get finance summary failed: %v", err)
	}

	accounts := make([]*pb.AccountBalance, 0, len(result.Accounts))
	for _, a := range result.Accounts {
		accounts = append(accounts, &pb.AccountBalance{
			AccountCode: a.AccountCode,
			DebitTotal:  a.DebitTotal,
			CreditTotal: a.CreditTotal,
			Net:         a.Net,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetFinanceSummaryResponse{Accounts: accounts}, nil
}

// ListJournalEntries handles the ListJournalEntries RPC.
// Returns a cursor-paginated list of journal entries with their lines.
func (s *Server) ListJournalEntries(ctx context.Context, req *pb.ListJournalEntriesRequest) (*pb.ListJournalEntriesResponse, error) {
	const op = "grpc_api.Server.ListJournalEntries"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("from", req.GetFrom()),
		attribute.String("to", req.GetTo()),
		attribute.Int("limit", int(req.GetLimit())),
		attribute.String("cursor", req.GetCursor()),
	)

	params := &service.ListJournalEntriesParams{
		Limit: req.GetLimit(),
	}

	if s := req.GetFrom(); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			params.From = &t
		} else if t, err := time.Parse("2006-01-02", s); err == nil {
			// date-only input from HTML date picker (YYYY-MM-DD) — treat as start of day UTC
			params.From = &t
		}
	}
	if s := req.GetTo(); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			params.To = &t
		} else if t, err := time.Parse("2006-01-02", s); err == nil {
			// date-only input — treat as end of day UTC (23:59:59) for inclusive upper bound
			endOfDay := t.Add(24*time.Hour - time.Second)
			params.To = &endOfDay
		}
	}
	if s := req.GetCursor(); s != "" {
		if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
			params.Cursor = &t
		} else if t, err := time.Parse(time.RFC3339, s); err == nil {
			params.Cursor = &t
		}
	}

	logger.Info().
		Str("op", op).
		Str("from", req.GetFrom()).
		Str("to", req.GetTo()).
		Int32("limit", req.GetLimit()).
		Str("cursor", req.GetCursor()).
		Msg("")

	result, err := s.svc.ListJournalEntries(ctx, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "list journal entries failed: %v", err)
	}

	protoEntries := make([]*pb.JournalEntryProto, 0, len(result.Entries))
	for _, e := range result.Entries {
		lines := make([]*pb.JournalLineProto, 0, len(e.Lines))
		for _, l := range e.Lines {
			lines = append(lines, &pb.JournalLineProto{
				Id:          l.ID,
				EntryId:     l.EntryID,
				AccountCode: l.AccountCode,
				Debit:       l.Debit,
				Credit:      l.Credit,
			})
		}
		protoEntries = append(protoEntries, &pb.JournalEntryProto{
			Id:             e.ID,
			IdempotencyKey: e.IdempotencyKey,
			SourceType:     e.SourceType,
			SourceId:       e.SourceID,
			PostedAt:       e.PostedAt.UTC().Format(time.RFC3339),
			Description:    e.Description,
			Lines:          lines,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.ListJournalEntriesResponse{
		Entries:    protoEntries,
		NextCursor: result.NextCursor,
	}, nil
}
