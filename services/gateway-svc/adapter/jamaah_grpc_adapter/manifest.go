// manifest.go — gateway adapter methods for jamaah-svc manifest RPCs (Wave 1A / Phase 6).
//
// Each method translates gateway-local params → pb request, forwards via gRPC,
// and translates pb response → adapter-local types. Proto types do not leak
// past this package.
package jamaah_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/jamaah_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Adapter-local result types
// ---------------------------------------------------------------------------

// ManifestJamaahResult is the gateway-local representation of one pilgrim row.
type ManifestJamaahResult struct {
	BookingID     string
	Name          string
	NIK           string
	Phone         string
	RoomType      string
	BookingStatus string
	DocStatus     string
}

// GetDepartureManifestResult is the gateway-local result for GetDepartureManifest.
type GetDepartureManifestResult struct {
	DepartureID string
	TotalJamaah int32
	LunasPaid   int32
	DocComplete int32
	JamaahList  []*ManifestJamaahResult
}

// ---------------------------------------------------------------------------
// GetDepartureManifest
// ---------------------------------------------------------------------------

// GetDepartureManifest fetches the departure manifest from jamaah-svc.
func (a *Adapter) GetDepartureManifest(ctx context.Context, departureID string) (*GetDepartureManifestResult, error) {
	const op = "jamaah_grpc_adapter.Adapter.GetDepartureManifest"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "GetDepartureManifest"),
		attribute.String("departure_id", departureID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.manifestClient.GetDepartureManifest(ctx, &pb.GetDepartureManifestRequest{
		DepartureID: departureID,
	})
	if err != nil {
		wrapped := mapJamaahError(err)
		logger.Warn().Err(wrapped).Str("departure_id", departureID).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	list := make([]*ManifestJamaahResult, 0, len(resp.GetJamaahList()))
	for _, j := range resp.GetJamaahList() {
		list = append(list, &ManifestJamaahResult{
			BookingID:     j.GetBookingID(),
			Name:          j.GetName(),
			NIK:           j.GetNIK(),
			Phone:         j.GetPhone(),
			RoomType:      j.GetRoomType(),
			BookingStatus: j.GetBookingStatus(),
			DocStatus:     j.GetDocStatus(),
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &GetDepartureManifestResult{
		DepartureID: resp.GetDepartureID(),
		TotalJamaah: resp.GetTotalJamaah(),
		LunasPaid:   resp.GetLunasPaid(),
		DocComplete: resp.GetDocComplete(),
		JamaahList:  list,
	}, nil
}
