// scan_boarding.go — ops-svc service logic for scan events and bus boarding
// (BL-OPS-010/011).
//
// RecordScan:        idempotent scan event recording.
// RecordBusBoarding: atomic scan + boarding insert with idempotency check.
// GetBoardingRoster: aggregate roster with totals.

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ops-svc/store/postgres_store/sqlc"
	"ops-svc/util/logging"
	"ops-svc/util/ulid"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Params / Results
// ---------------------------------------------------------------------------

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

type RecordScanResult struct {
	ScanID     string
	Idempotent bool
}

type RecordBusBoardingParams struct {
	DepartureID string
	BusNumber   string
	JamaahID    string
	ScannedBy   string
	Status      string // "boarded" | "absent" | "late"
}

type RecordBusBoardingResult struct {
	BoardingID string
	Status     string
	Idempotent bool
}

type GetBoardingRosterParams struct {
	DepartureID string
	BusNumber   string // optional filter
}

type BoardingEntry struct {
	JamaahID  string
	BusNumber string
	Status    string
	BoardedAt string // RFC3339 or empty
}

type GetBoardingRosterResult struct {
	Boardings    []BoardingEntry
	TotalBoarded int32
	TotalAbsent  int32
}

// ---------------------------------------------------------------------------
// Sentinel errors
// ---------------------------------------------------------------------------

var (
	ErrInvalidScanType  = errors.New("invalid_scan_type")
	ErrMissingField     = errors.New("missing_required_field")
	ErrInvalidBoardingStatus = errors.New("invalid_status")
)

var validScanTypes = map[string]bool{
	"ALL": true, "bus_boarding": true, "luggage": true, "raudhah": true,
}

var validBoardingStatuses = map[string]bool{
	"boarded": true, "absent": true, "late": true,
}

// ---------------------------------------------------------------------------
// RecordScan (BL-OPS-010)
// ---------------------------------------------------------------------------

func (s *Service) RecordScan(ctx context.Context, params *RecordScanParams) (*RecordScanResult, error) {
	const op = "service.RecordScan"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	if !validScanTypes[params.ScanType] {
		return nil, ErrInvalidScanType
	}
	if params.DepartureID == "" || params.JamaahID == "" || params.ScannedBy == "" || params.IdempotencyKey == "" {
		return nil, ErrMissingField
	}

	id, err := ulid.New("scan_")
	if err != nil {
		return nil, fmt.Errorf("ulid: %w", err)
	}

	metadata := params.Metadata
	if metadata == nil {
		metadata = []byte("{}")
	}

	existingID, wasInserted, err := s.store.InsertScanEventIdempotent(ctx, sqlc.InsertScanEventIdempotentParams{
		ID:             id,
		ScanType:       params.ScanType,
		DepartureID:    params.DepartureID,
		JamaahID:       params.JamaahID,
		ScannedBy:      params.ScannedBy,
		DeviceID:       params.DeviceID,
		Location:       params.Location,
		IdempotencyKey: params.IdempotencyKey,
		Metadata:       metadata,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("insert scan event")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	idempotent := !wasInserted
	span.SetAttributes(
		attribute.String("scan_id", existingID),
		attribute.Bool("idempotent", idempotent),
	)

	return &RecordScanResult{ScanID: existingID, Idempotent: idempotent}, nil
}

// ---------------------------------------------------------------------------
// RecordBusBoarding (BL-OPS-011)
// ---------------------------------------------------------------------------

func (s *Service) RecordBusBoarding(ctx context.Context, params *RecordBusBoardingParams) (*RecordBusBoardingResult, error) {
	const op = "service.RecordBusBoarding"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	if params.DepartureID == "" || params.BusNumber == "" || params.JamaahID == "" || params.ScannedBy == "" {
		return nil, ErrMissingField
	}
	boardingStatus := params.Status
	if boardingStatus == "" {
		boardingStatus = "boarded"
	}
	if !validBoardingStatuses[boardingStatus] {
		return nil, ErrInvalidBoardingStatus
	}

	// Idempotency check: boarding already exists for (departure, jamaah)?
	existing, err := s.store.GetBoardingByDepartureJamaah(ctx, params.DepartureID, params.JamaahID)
	if err == nil {
		// Row found — return existing, idempotent
		return &RecordBusBoardingResult{
			BoardingID: existing.ID,
			Status:     existing.Status,
			Idempotent: true,
		}, nil
	}

	// Insert scan event
	ikey := fmt.Sprintf("bus_boarding:%s:%s:%s", params.DepartureID, params.JamaahID, time.Now().UTC().Format("2006-01-02"))
	scanID, err := ulid.New("scan_")
	if err != nil {
		return nil, fmt.Errorf("ulid scan: %w", err)
	}
	existingScanID, _, err := s.store.InsertScanEventIdempotent(ctx, sqlc.InsertScanEventIdempotentParams{
		ID:             scanID,
		ScanType:       "bus_boarding",
		DepartureID:    params.DepartureID,
		JamaahID:       params.JamaahID,
		ScannedBy:      params.ScannedBy,
		IdempotencyKey: ikey,
		Metadata:       []byte("{}"),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("insert scan event for boarding")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	boardingID, err := ulid.New("bbd_")
	if err != nil {
		return nil, fmt.Errorf("ulid boarding: %w", err)
	}
	boarding, err := s.store.InsertBusBoarding(ctx, sqlc.InsertBusBoardingParams{
		ID:          boardingID,
		DepartureID: params.DepartureID,
		BusNumber:   params.BusNumber,
		JamaahID:    params.JamaahID,
		Status:      boardingStatus,
		ScanEventID: existingScanID,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("insert bus boarding")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	return &RecordBusBoardingResult{
		BoardingID: boarding.ID,
		Status:     boarding.Status,
		Idempotent: false,
	}, nil
}

// ---------------------------------------------------------------------------
// GetBoardingRoster (BL-OPS-011)
// ---------------------------------------------------------------------------

func (s *Service) GetBoardingRoster(ctx context.Context, params *GetBoardingRosterParams) (*GetBoardingRosterResult, error) {
	const op = "service.GetBoardingRoster"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	rows, err := s.store.GetBoardingRoster(ctx, sqlc.GetBoardingRosterParams{
		DepartureID: params.DepartureID,
		BusNumber:   params.BusNumber,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("get boarding roster")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	var entries []BoardingEntry
	var totalBoarded, totalAbsent int32
	for _, r := range rows {
		boardedAt := ""
		if r.BoardedAt.Valid {
			boardedAt = r.BoardedAt.Time.UTC().Format(time.RFC3339)
		}
		entries = append(entries, BoardingEntry{
			JamaahID:  r.JamaahID,
			BusNumber: r.BusNumber,
			Status:    r.Status,
			BoardedAt: boardedAt,
		})
		switch r.Status {
		case "boarded", "late":
			totalBoarded++
		case "absent":
			totalAbsent++
		}
	}

	return &GetBoardingRosterResult{
		Boardings:    entries,
		TotalBoarded: totalBoarded,
		TotalAbsent:  totalAbsent,
	}, nil
}
