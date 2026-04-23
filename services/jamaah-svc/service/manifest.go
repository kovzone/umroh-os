package service

// manifest.go — departure manifest service implementation.
//
// GetDepartureManifest fetches all active pilgrims for a given departure_id,
// aggregates document completion, and computes summary counters.
//
// Wave-1A manifest API.

import (
	"context"
	"errors"
	"fmt"

	"jamaah-svc/util/apperrors"
	"jamaah-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Input / output types
// ---------------------------------------------------------------------------

// GetDepartureManifestParams is the input for GetDepartureManifest.
type GetDepartureManifestParams struct {
	DepartureID string
}

// ManifestJamaahItem is one pilgrim row in the manifest result.
type ManifestJamaahItem struct {
	BookingID     string
	Name          string
	NIK           string
	Phone         string
	RoomType      string
	BookingStatus string
	DocStatus     string // "complete" | "partial" | "none"
}

// GetDepartureManifestResult is the output for GetDepartureManifest.
type GetDepartureManifestResult struct {
	DepartureID string
	TotalJamaah int32
	LunasPaid   int32
	DocComplete int32
	JamaahList  []*ManifestJamaahItem
}

// ---------------------------------------------------------------------------
// Implementation
// ---------------------------------------------------------------------------

// GetDepartureManifest returns the departure manifest for the given departure_id.
func (s *Service) GetDepartureManifest(ctx context.Context, params *GetDepartureManifestParams) (*GetDepartureManifestResult, error) {
	const op = "service.Service.GetDepartureManifest"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("departure_id", params.DepartureID),
	)

	if params.DepartureID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("departure_id is required"))
	}

	rows, err := s.store.GetDepartureManifest(ctx, params.DepartureID)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("GetDepartureManifest failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("get departure manifest: %w", err)
	}

	result := &GetDepartureManifestResult{
		DepartureID: params.DepartureID,
		TotalJamaah: int32(len(rows)),
		JamaahList:  make([]*ManifestJamaahItem, 0, len(rows)),
	}

	for _, r := range rows {
		// Determine doc_status.
		docStatus := "none"
		if r.TotalDocs > 0 {
			if r.ApprovedDocs == r.TotalDocs {
				docStatus = "complete"
			} else if r.ApprovedDocs > 0 {
				docStatus = "partial"
			} else {
				docStatus = "none"
			}
		}

		// Count lunas_paid: bookings that are paid_in_full or partially_paid.
		if r.BookingStatus == "paid_in_full" || r.BookingStatus == "partially_paid" {
			result.LunasPaid++
		}

		// Count doc_complete.
		if docStatus == "complete" {
			result.DocComplete++
		}

		result.JamaahList = append(result.JamaahList, &ManifestJamaahItem{
			BookingID:     r.BookingID,
			Name:          r.Name,
			NIK:           r.NIK,
			Phone:         r.Phone,
			RoomType:      r.RoomType,
			BookingStatus: r.BookingStatus,
			DocStatus:     docStatus,
		})
	}

	span.SetStatus(otelCodes.Ok, "success")
	return result, nil
}
