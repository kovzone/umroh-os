// readiness.go — gateway adapter methods for catalog-svc vendor readiness RPCs
// (BL-OPS-020).
//
// Each method translates gateway-local params → pb request, forwards via gRPC,
// and translates pb response → adapter-local types. Proto types do not leak
// past this package.
package catalog_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/catalog_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Adapter-local types
// ---------------------------------------------------------------------------

// ReadinessResult is the gateway-local representation of departure readiness.
type ReadinessResult struct {
	TicketState string
	HotelState  string
	VisaState   string
}

// UpdateVendorReadinessParams is the input for UpdateVendorReadiness.
type UpdateVendorReadinessParams struct {
	UserID        string
	DepartureID   string
	Kind          string // ticket | hotel | visa
	State         string // not_started | in_progress | done
	Notes         string
	AttachmentURL string
}

// GetDepartureReadinessParams is the input for GetDepartureReadiness.
type GetDepartureReadinessParams struct {
	UserID      string
	DepartureID string
}

// ---------------------------------------------------------------------------
// UpdateVendorReadiness
// ---------------------------------------------------------------------------

// UpdateVendorReadiness upserts one readiness kind for a departure and returns
// the full readiness summary.
func (a *Adapter) UpdateVendorReadiness(ctx context.Context, params *UpdateVendorReadinessParams) (*ReadinessResult, error) {
	const op = "catalog_grpc_adapter.Adapter.UpdateVendorReadiness"

	logger := logging.LogWithTrace(ctx, a.logger)
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("departure_id", params.DepartureID),
		attribute.String("kind", params.Kind),
	)
	logger.Info().Str("op", op).
		Str("departure_id", params.DepartureID).
		Str("kind", params.Kind).
		Msg("")

	resp, err := a.catalogReadinessClient.UpdateVendorReadiness(ctx, &pb.UpdateVendorReadinessRequest{
		UserId:        params.UserID,
		DepartureId:   params.DepartureID,
		Kind:          params.Kind,
		State:         params.State,
		Notes:         params.Notes,
		AttachmentUrl: params.AttachmentURL,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, mapCatalogError(err)
	}

	span.SetStatus(codes.Ok, "success")
	return readinessProtoToResult(resp.GetReadiness()), nil
}

// ---------------------------------------------------------------------------
// GetDepartureReadiness
// ---------------------------------------------------------------------------

// GetDepartureReadiness fetches the current readiness state for a departure.
func (a *Adapter) GetDepartureReadiness(ctx context.Context, params *GetDepartureReadinessParams) (*ReadinessResult, error) {
	const op = "catalog_grpc_adapter.Adapter.GetDepartureReadiness"

	logger := logging.LogWithTrace(ctx, a.logger)
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("departure_id", params.DepartureID),
	)
	logger.Info().Str("op", op).Str("departure_id", params.DepartureID).Msg("")

	resp, err := a.catalogReadinessClient.GetDepartureReadiness(ctx, &pb.GetDepartureReadinessRequest{
		UserId:      params.UserID,
		DepartureId: params.DepartureID,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, mapCatalogError(err)
	}

	span.SetStatus(codes.Ok, "success")
	return readinessProtoToResult(resp.GetReadiness()), nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func readinessProtoToResult(r *pb.ReadinessState) *ReadinessResult {
	if r == nil {
		return &ReadinessResult{
			TicketState: "not_started",
			HotelState:  "not_started",
			VisaState:   "not_started",
		}
	}
	ticket := r.GetTicketState()
	if ticket == "" {
		ticket = "not_started"
	}
	hotel := r.GetHotelState()
	if hotel == "" {
		hotel = "not_started"
	}
	visa := r.GetVisaState()
	if visa == "" {
		visa = "not_started"
	}
	return &ReadinessResult{
		TicketState: ticket,
		HotelState:  hotel,
		VisaState:   visa,
	}
}
