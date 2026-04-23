// ops_phase6.go — gateway-svc adapter methods for ops-svc Phase 6 RPCs
// (BL-OPS-010/011): RecordScan, RecordBusBoarding, GetBoardingRoster.

package ops_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/ops_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// RecordScan
// ---------------------------------------------------------------------------

// RecordScanParams holds inputs for POST /v1/ops/scans.
type RecordScanParams struct {
	ScanType       string
	DepartureID    string
	JamaahID       string
	ScannedBy      string
	DeviceID       string
	Location       string
	IdempotencyKey string
	Metadata       []byte
}

// RecordScanResult is the response.
type RecordScanResult struct {
	ScanID     string
	Idempotent bool
}

func (a *Adapter) RecordScan(ctx context.Context, params *RecordScanParams) (*RecordScanResult, error) {
	const op = "ops_grpc_adapter.Adapter.RecordScan"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("scan_type", params.ScanType).Msg("")

	resp, err := a.opsPhase6Client.RecordScan(ctx, &pb.RecordScanRequest{
		ScanType:       params.ScanType,
		DepartureId:    params.DepartureID,
		JamaahId:       params.JamaahID,
		ScannedBy:      params.ScannedBy,
		DeviceId:       params.DeviceID,
		Location:       params.Location,
		IdempotencyKey: params.IdempotencyKey,
		Metadata:       params.Metadata,
	})
	if err != nil {
		wrapped := mapOpsError(err)
		logger.Warn().Err(wrapped).Msg("ops-svc.RecordScan failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &RecordScanResult{
		ScanID:     resp.GetScanId(),
		Idempotent: resp.GetIdempotent(),
	}, nil
}

// ---------------------------------------------------------------------------
// RecordBusBoarding
// ---------------------------------------------------------------------------

// RecordBusBoardingParams holds inputs for POST /v1/ops/bus-boarding.
type RecordBusBoardingParams struct {
	DepartureID string
	BusNumber   string
	JamaahID    string
	ScannedBy   string
	Status      string
}

// RecordBusBoardingResult is the response.
type RecordBusBoardingResult struct {
	BoardingID string
	Status     string
	Idempotent bool
}

func (a *Adapter) RecordBusBoarding(ctx context.Context, params *RecordBusBoardingParams) (*RecordBusBoardingResult, error) {
	const op = "ops_grpc_adapter.Adapter.RecordBusBoarding"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("departure_id", params.DepartureID).Msg("")

	resp, err := a.opsPhase6Client.RecordBusBoarding(ctx, &pb.RecordBusBoardingRequest{
		DepartureId: params.DepartureID,
		BusNumber:   params.BusNumber,
		JamaahId:    params.JamaahID,
		ScannedBy:   params.ScannedBy,
		Status:      params.Status,
	})
	if err != nil {
		wrapped := mapOpsError(err)
		logger.Warn().Err(wrapped).Msg("ops-svc.RecordBusBoarding failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &RecordBusBoardingResult{
		BoardingID: resp.GetBoardingId(),
		Status:     resp.GetStatus(),
		Idempotent: resp.GetIdempotent(),
	}, nil
}

// ---------------------------------------------------------------------------
// GetBoardingRoster
// ---------------------------------------------------------------------------

// BoardingEntryResult holds one entry in the boarding roster.
type BoardingEntryResult struct {
	JamaahID  string
	BusNumber string
	Status    string
	BoardedAt string
}

// GetBoardingRosterResult is the response for GET /v1/ops/bus-boarding/:departure_id.
type GetBoardingRosterResult struct {
	Boardings    []*BoardingEntryResult
	TotalBoarded int32
	TotalAbsent  int32
}

func (a *Adapter) GetBoardingRoster(ctx context.Context, departureID, busNumber string) (*GetBoardingRosterResult, error) {
	const op = "ops_grpc_adapter.Adapter.GetBoardingRoster"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("departure_id", departureID).Msg("")

	resp, err := a.opsPhase6Client.GetBoardingRoster(ctx, &pb.GetBoardingRosterRequest{
		DepartureId: departureID,
		BusNumber:   busNumber,
	})
	if err != nil {
		wrapped := mapOpsError(err)
		logger.Warn().Err(wrapped).Msg("ops-svc.GetBoardingRoster failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	boardings := make([]*BoardingEntryResult, 0, len(resp.GetBoardings()))
	for _, b := range resp.GetBoardings() {
		boardings = append(boardings, &BoardingEntryResult{
			JamaahID:  b.GetJamaahId(),
			BusNumber: b.GetBusNumber(),
			Status:    b.GetStatus(),
			BoardedAt: b.GetBoardedAt(),
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &GetBoardingRosterResult{
		Boardings:    boardings,
		TotalBoarded: resp.GetTotalBoarded(),
		TotalAbsent:  resp.GetTotalAbsent(),
	}, nil
}
