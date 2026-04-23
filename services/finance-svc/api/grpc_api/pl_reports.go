// pl_reports.go — gRPC handlers for GetPLReport and GetBalanceSheet RPCs
// (finance depth / Wave 1B / BL-FIN-007..008).
//
// Called by gateway-svc to populate:
//   GET /v1/finance/pl-report    — P&L for a date range
//   GET /v1/finance/balance-sheet — balance sheet as of a date

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

// GetPLReport handles the GetPLReport RPC.
// Returns a P&L report for the given date range.
func (s *Server) GetPLReport(ctx context.Context, req *pb.GetPLReportRequest) (*pb.PLReportProto, error) {
	const op = "grpc_api.Server.GetPLReport"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	span.SetAttributes(
		attribute.String("from", req.GetFrom()),
		attribute.String("to", req.GetTo()),
	)

	logger.Info().
		Str("op", op).
		Str("from", req.GetFrom()).
		Str("to", req.GetTo()).
		Msg("")

	params := &service.GetPLReportParams{}

	if s := req.GetFrom(); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			params.From = t
		} else if t, err := time.Parse("2006-01-02", s); err == nil {
			params.From = t
		}
	}
	if s := req.GetTo(); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			params.To = t
		} else if t, err := time.Parse("2006-01-02", s); err == nil {
			// Treat date-only as end of day UTC for inclusive upper bound.
			params.To = t.Add(24*time.Hour - time.Second)
		}
	}

	result, err := s.svc.GetPLReport(ctx, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "get P&L report failed: %v", err)
	}

	entries := make([]*pb.PLLineItemProto, 0, len(result.Entries))
	for _, e := range result.Entries {
		entries = append(entries, &pb.PLLineItemProto{
			AccountCode: e.AccountCode,
			AccountName: e.AccountName,
			Amount:      e.Amount,
			Direction:   e.Direction,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.PLReportProto{
		PeriodFrom:   result.PeriodFrom.UTC().Format(time.RFC3339),
		PeriodTo:     result.PeriodTo.UTC().Format(time.RFC3339),
		GeneratedAt:  result.GeneratedAt.UTC().Format(time.RFC3339),
		TotalRevenue: result.TotalRevenue,
		TotalCogs:    result.TotalCOGS,
		GrossProfit:  result.GrossProfit,
		OtherIncome:  result.OtherIncome,
		OtherExpense: result.OtherExpense,
		NetProfit:    result.NetProfit,
		Entries:      entries,
	}, nil
}

// GetBalanceSheet handles the GetBalanceSheet RPC.
// Returns a balance sheet as of the given date.
func (s *Server) GetBalanceSheet(ctx context.Context, req *pb.GetBalanceSheetRequest) (*pb.BalanceSheetProto, error) {
	const op = "grpc_api.Server.GetBalanceSheet"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("as_of", req.GetAsOf()))

	logger.Info().
		Str("op", op).
		Str("as_of", req.GetAsOf()).
		Msg("")

	params := &service.GetBalanceSheetParams{}

	if s := req.GetAsOf(); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			params.AsOfDate = t
		} else if t, err := time.Parse("2006-01-02", s); err == nil {
			// End of day UTC for inclusive upper bound.
			params.AsOfDate = t.Add(24*time.Hour - time.Second)
		}
	}

	result, err := s.svc.GetBalanceSheet(ctx, params)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "get balance sheet failed: %v", err)
	}

	toProtoLines := func(lines []service.BalanceSheetLine) []*pb.BalanceSheetLineProto {
		out := make([]*pb.BalanceSheetLineProto, 0, len(lines))
		for _, l := range lines {
			out = append(out, &pb.BalanceSheetLineProto{
				AccountCode: l.AccountCode,
				AccountName: l.AccountName,
				Balance:     l.Balance,
			})
		}
		return out
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.BalanceSheetProto{
		AsOfDate:         result.AsOfDate.UTC().Format(time.RFC3339),
		GeneratedAt:      result.GeneratedAt.UTC().Format(time.RFC3339),
		Assets:           toProtoLines(result.Assets),
		Liabilities:      toProtoLines(result.Liabilities),
		Equity:           toProtoLines(result.Equity),
		TotalAssets:      result.TotalAssets,
		TotalLiabilities: result.TotalLiabilities,
		TotalEquity:      result.TotalEquity,
	}, nil
}
