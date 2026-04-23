package service

// vendor_readiness.go — BL-OPS-020: departure vendor readiness checklist.
//
// Staff can mark each departure's ticket / hotel / visa readiness state as
// not_started → in_progress → done. The state is stored in
// catalog.departure_vendor_readiness (one row per departure × kind).

import (
	"context"
	"errors"
	"fmt"

	"catalog-svc/store/postgres_store"
	"catalog-svc/store/postgres_store/sqlc"
	"catalog-svc/util/apperrors"
	"catalog-svc/util/logging"
	"catalog-svc/util/ulid"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Params / result types
// ---------------------------------------------------------------------------

// UpdateVendorReadinessParams is the input for UpdateVendorReadiness.
type UpdateVendorReadinessParams struct {
	DepartureID   string // required
	Kind          string // ticket | hotel | visa
	State         string // not_started | in_progress | done
	Notes         string // optional free-text
	AttachmentURL string // optional URL
	UpdatedBy     string // actor subject (user ID)
}

// GetDepartureReadinessParams is the input for GetDepartureReadiness.
type GetDepartureReadinessParams struct {
	DepartureID string
}

// validReadinessKinds is the allowed set of readiness kind values.
var validReadinessKinds = map[string]bool{
	"ticket": true,
	"hotel":  true,
	"visa":   true,
}

// validReadinessStates is the allowed set of readiness state values.
var validReadinessStates = map[string]bool{
	"not_started": true,
	"in_progress": true,
	"done":        true,
}

// ---------------------------------------------------------------------------
// UpdateVendorReadiness
// ---------------------------------------------------------------------------

// UpdateVendorReadiness upserts a single vendor readiness record for a
// departure. Returns the full readiness summary (all three kinds) after the
// update so the caller can return the current state in one round-trip.
func (s *Service) UpdateVendorReadiness(ctx context.Context, params *UpdateVendorReadinessParams) (*VendorReadiness, error) {
	const op = "service.Service.UpdateVendorReadiness"

	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.departure_id", params.DepartureID),
		attribute.String("input.kind", params.Kind),
		attribute.String("input.state", params.State),
	)
	logger.Info().Str("op", op).
		Str("departure_id", params.DepartureID).
		Str("kind", params.Kind).
		Str("state", params.State).
		Msg("")

	// Validate kind and state.
	if !validReadinessKinds[params.Kind] {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf(
			"kind %q is not valid; accepted: ticket, hotel, visa", params.Kind))
	}
	if !validReadinessStates[params.State] {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf(
			"state %q is not valid; accepted: not_started, in_progress, done", params.State))
	}

	// Verify the departure exists.
	_, err := s.store.GetDepartureByIDForStaff(ctx, params.DepartureID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("departure %q not found", params.DepartureID))
		}
		return nil, fmt.Errorf("get departure: %w", postgres_store.WrapDBError(err))
	}

	// Generate a new ULID for the id — only used on INSERT path of the upsert.
	id, err := ulid.New("vr_")
	if err != nil {
		return nil, fmt.Errorf("generate readiness id: %w", err)
	}

	_, err = s.store.UpsertVendorReadiness(ctx, sqlc.UpsertVendorReadinessParams{
		ID:            id,
		DepartureID:   params.DepartureID,
		Kind:          params.Kind,
		State:         params.State,
		Notes:         params.Notes,
		AttachmentURL: params.AttachmentURL,
		UpdatedBy:     params.UpdatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("upsert vendor readiness: %w", postgres_store.WrapDBError(err))
	}

	// Return the full readiness summary.
	result, err := s.getReadinessFromDB(ctx, params.DepartureID)
	if err != nil {
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// ---------------------------------------------------------------------------
// GetDepartureReadiness
// ---------------------------------------------------------------------------

// GetDepartureReadiness fetches the full vendor readiness summary for a
// departure. Returns a VendorReadiness where any un-initialised kind defaults
// to "not_started".
func (s *Service) GetDepartureReadiness(ctx context.Context, params *GetDepartureReadinessParams) (*VendorReadiness, error) {
	const op = "service.Service.GetDepartureReadiness"

	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.departure_id", params.DepartureID),
	)
	logger.Info().Str("op", op).Str("departure_id", params.DepartureID).Msg("")

	// Verify departure exists (staff endpoint — uses staff query).
	_, err := s.store.GetDepartureByIDForStaff(ctx, params.DepartureID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("departure %q not found", params.DepartureID))
		}
		return nil, fmt.Errorf("get departure: %w", postgres_store.WrapDBError(err))
	}

	result, err := s.getReadinessFromDB(ctx, params.DepartureID)
	if err != nil {
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// ---------------------------------------------------------------------------
// internal helper
// ---------------------------------------------------------------------------

// getReadinessFromDB lists rows from catalog.departure_vendor_readiness and
// maps them to a VendorReadiness struct. Missing kinds default to "not_started".
func (s *Service) getReadinessFromDB(ctx context.Context, departureID string) (*VendorReadiness, error) {
	rows, err := s.store.ListVendorReadiness(ctx, departureID)
	if err != nil {
		return nil, fmt.Errorf("list vendor readiness: %w", postgres_store.WrapDBError(err))
	}

	r := &VendorReadiness{
		Ticket: "not_started",
		Hotel:  "not_started",
		Visa:   "not_started",
	}
	for _, row := range rows {
		switch row.Kind {
		case "ticket":
			r.Ticket = row.State
		case "hotel":
			r.Hotel = row.State
		case "visa":
			r.Visa = row.State
		}
	}
	return r, nil
}
