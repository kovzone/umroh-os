// ops_manifest.go — manifest export service-layer implementation (BL-OPS-001).
//
// ExportManifest returns structured rows for a departure manifest.
// NOTE: Real data will come from jamaah-svc via gRPC in a later wiring step.
// For now, returns empty rows (endpoint structure only).

package service

import (
	"context"
	"fmt"

	"ops-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// ManifestRow is a single row in the manifest export.
type ManifestRow struct {
	No         int32
	JamaahName string
	PassportNo string
	DocStatus  string
	RoomNumber string
}

// ExportManifestParams holds inputs for ExportManifest.
type ExportManifestParams struct {
	DepartureID string
}

// ExportManifestResult holds the result of ExportManifest.
type ExportManifestResult struct {
	Rows []ManifestRow
}

// ExportManifest returns manifest rows for a departure.
// NOTE: Currently returns empty rows. Real jamaah data is wired in a later
// sprint via gRPC call to jamaah-svc.
func (svc *Service) ExportManifest(ctx context.Context, params *ExportManifestParams) (*ExportManifestResult, error) {
	const op = "service.Service.ExportManifest"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("departure_id", params.DepartureID),
	)

	logger := logging.LogWithTrace(ctx, svc.logger)

	if params.DepartureID == "" {
		return nil, fmt.Errorf("%s: departure_id is required", op)
	}

	// TODO(S3-Wave2): fetch real data from jamaah-svc via gRPC and merge with
	// room assignment data. For now return empty rows as endpoint scaffold.
	logger.Info().
		Str("op", op).
		Str("departure_id", params.DepartureID).
		Msg("ExportManifest called (stub — empty rows)")

	span.SetStatus(otelCodes.Ok, "success")
	return &ExportManifestResult{
		Rows: []ManifestRow{},
	}, nil
}
