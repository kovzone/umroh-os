package service

// S1-E-07 / BL-CAT-014 — Staff catalog write service methods.
//
// Permission gate: callers (gRPC handlers) validate user_id is non-empty and
// call iam-svc.CheckPermission before invoking these methods. Service layer
// trusts the gRPC handler has already gated the call.
//
// All database mutations run through the store layer (sqlc-generated types).
// ID generation uses the util/ulid package.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"catalog-svc/store/postgres_store"
	"catalog-svc/store/postgres_store/sqlc"
	"catalog-svc/util/apperrors"
	"catalog-svc/util/logging"
	"catalog-svc/util/ulid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Input / output types for write methods
// ---------------------------------------------------------------------------

// PricingInput is a single room-type price row for create/update.
type PricingInput struct {
	RoomType           string
	ListAmount         int64
	ListCurrency       string
	SettlementCurrency string
}

// CreatePackageParams is the input for CreatePackage.
type CreatePackageParams struct {
	UserID        string
	BranchID      string
	Kind          string
	Name          string
	Description   string
	CoverPhotoUrl string
	Highlights    []string
	ItineraryID   string
	AirlineID     string
	MuthawwifID   string
	HotelIDs      []string
	AddonIDs      []string
	Status        string // "draft" | "active" — defaults to "draft"
}

// UpdatePackageParams is the input for UpdatePackage.
type UpdatePackageParams struct {
	UserID        string
	BranchID      string
	ID            string
	Name          string   // empty = no change
	Description   string   // empty = no change
	CoverPhotoUrl string   // empty = no change
	Highlights    []string // nil = no change
	ItineraryID   string   // empty = no change
	AirlineID     string   // empty = no change
	MuthawwifID   string   // empty = no change
	HotelIDs      []string // nil = no change; non-nil = replace
	AddonIDs      []string // nil = no change; non-nil = replace
	Status        string   // empty = no change
}

// DeletePackageParams is the input for DeletePackage.
type DeletePackageParams struct {
	UserID   string
	BranchID string
	ID       string
}

// CreateDepartureParams is the input for CreateDeparture.
type CreateDepartureParams struct {
	UserID        string
	BranchID      string
	PackageID     string
	DepartureDate string // ISO YYYY-MM-DD
	ReturnDate    string // ISO YYYY-MM-DD
	TotalSeats    int    // must be > 0
	Status        string // "open" | "closed" — defaults to "open"
	Pricing       []PricingInput
}

// UpdateDepartureParams is the input for UpdateDeparture.
type UpdateDepartureParams struct {
	UserID        string
	BranchID      string
	ID            string
	DepartureDate string       // empty = no change
	ReturnDate    string       // empty = no change
	TotalSeats    int          // 0 = no change
	Status        string       // empty = no change
	Pricing       []PricingInput // nil = no change; non-nil = replace all
}

// ReserveSeatsParams mirrors the § Inventory contract.
type ReserveSeatsParams struct {
	ReservationID       string
	DepartureID         string
	Seats               int
	IdempotencyTTLHours int // 0 = default 24; clamped [1, 168]
}

// ReserveSeatsResult mirrors the § Inventory response.
type ReserveSeatsResult struct {
	ReservationID  string
	DepartureID    string
	Seats          int
	ReservedAt     string
	ExpiresAt      string
	RemainingSeats int
	Replayed       bool
}

// ReleaseSeatsParams mirrors the § Inventory ReleaseSeats contract.
type ReleaseSeatsParams struct {
	ReservationID string
	Seats         int    // 0 = full release
	Reason        string // audit note ≤ 256 chars
}

// ReleaseSeatsResult mirrors the § Inventory ReleaseSeats response.
type ReleaseSeatsResult struct {
	ReservationID  string
	DepartureID    string
	SeatsReleased  int
	ReleasedAt     string
	RemainingSeats int
	Replayed       bool
}

// ---------------------------------------------------------------------------
// BL-CAT-008: Package variant kind helpers
// ---------------------------------------------------------------------------

// travelKinds are the package kinds that require transport + accommodation.
// Financial and retail kinds are product-only and do not carry itinerary data.
var travelKinds = map[string]bool{
	"umrah_reguler": true,
	"umrah_plus":    true,
	"hajj_furoda":   true,
	"hajj_khusus":   true,
	"badal":         true,
}

// isTravelKind reports whether kind is a travel (itinerary-bearing) package.
func isTravelKind(kind string) bool { return travelKinds[kind] }

// validKinds is the full set of accepted package kinds.
var validKinds = map[string]bool{
	"umrah_reguler": true,
	"umrah_plus":    true,
	"hajj_furoda":   true,
	"hajj_khusus":   true,
	"badal":         true,
	"financial":     true,
	"retail":        true,
}

// validateKindOnCreate enforces BL-CAT-008 constraints at creation time.
//
// Rules:
//   - kind must be a valid CatalogPackageKind.
//   - Travel-kind packages being published (status = "active") must have
//     airline_id set — a departure cannot be created without a carrier.
//   - Non-travel kinds (financial, retail) never require airline_id.
func validateKindOnCreate(params *CreatePackageParams) error {
	if !validKinds[params.Kind] {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf(
			"kind %q is not valid; accepted: umrah_reguler, umrah_plus, hajj_furoda, hajj_khusus, badal, financial, retail",
			params.Kind,
		))
	}
	if isTravelKind(params.Kind) && params.Status == "active" && params.AirlineID == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf(
			"travel-kind package cannot be published (status=active) without an airline_id; "+
				"set airline_id or create as draft first",
		))
	}
	return nil
}

// validateKindOnPublish enforces BL-CAT-008 constraints when a package is
// transitioned to active during an update. Called when params.Status == "active".
func validateKindOnPublish(pkg sqlc.GetActivePackageByIDRow, newAirlineID string) error {
	if !isTravelKind(string(pkg.Kind)) {
		return nil
	}
	effectiveAirline := newAirlineID
	if effectiveAirline == "" && pkg.AirlineID.Valid {
		effectiveAirline = pkg.AirlineID.String
	}
	if effectiveAirline == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf(
			"travel-kind package cannot be published without an airline_id; "+
				"set airline_id before activating",
		))
	}
	return nil
}

// ---------------------------------------------------------------------------
// CreatePackage
// ---------------------------------------------------------------------------

func (s *Service) CreatePackage(ctx context.Context, params *CreatePackageParams) (*PackageDetail, error) {
	const op = "service.Service.CreatePackage"

	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.user_id", params.UserID),
		attribute.String("input.kind", params.Kind),
	)
	logger.Info().Str("op", op).Str("user_id", params.UserID).Msg("")

	if params.UserID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.Kind == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("kind is required"))
	}
	if params.Name == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("name is required"))
	}

	// BL-CAT-008: validate kind constraints (variant type + publish gate).
	if err := validateKindOnCreate(params); err != nil {
		return nil, err
	}

	status := sqlc.CatalogPackageStatusDraft
	switch params.Status {
	case "", "draft":
		status = sqlc.CatalogPackageStatusDraft
	case "active":
		status = sqlc.CatalogPackageStatusActive
	default:
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("status must be 'draft' or 'active' on create; got %q", params.Status))
	}

	id, err := ulid.New("pkg_")
	if err != nil {
		return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("mint package id: %w", err))
	}

	highlights := params.Highlights
	if highlights == nil {
		highlights = []string{}
	}

	_, err = s.store.InsertPackage(ctx, sqlc.InsertPackageParams{
		ID:            id,
		Kind:          sqlc.CatalogPackageKind(params.Kind),
		Name:          params.Name,
		Description:   params.Description,
		Highlights:    highlights,
		CoverPhotoUrl: params.CoverPhotoUrl,
		ItineraryID:   optText(params.ItineraryID),
		AirlineID:     optText(params.AirlineID),
		MuthawwifID:   optText(params.MuthawwifID),
		Status:        status,
	})
	if err != nil {
		return nil, fmt.Errorf("insert package: %w", postgres_store.WrapDBError(err))
	}

	for i, hid := range params.HotelIDs {
		if err := s.store.InsertPackageHotel(ctx, sqlc.InsertPackageHotelParams{
			PackageID: id, HotelID: hid, SortOrder: int16(i),
		}); err != nil {
			return nil, fmt.Errorf("link hotel %s: %w", hid, postgres_store.WrapDBError(err))
		}
	}
	for _, aid := range params.AddonIDs {
		// BL-CAT-009: validate addon existence before linking.
		if _, err := s.store.GetAddonByID(ctx, aid); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("addon_id %q does not exist in catalog", aid))
			}
			return nil, fmt.Errorf("lookup addon %s: %w", aid, postgres_store.WrapDBError(err))
		}
		if err := s.store.InsertPackageAddon(ctx, sqlc.InsertPackageAddonParams{
			PackageID: id, AddonID: aid,
		}); err != nil {
			return nil, fmt.Errorf("link addon %s: %w", aid, postgres_store.WrapDBError(err))
		}
	}

	detail, err := s.hydratePackageDetailForStaff(ctx, id)
	if err != nil {
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(attribute.String("output.package_id", id))

	// BL-CAT-012: emit catalog.updated webhook (best-effort).
	s.emitCatalogUpdated(ctx, "package.created", id, "package")

	return detail, nil
}

// ---------------------------------------------------------------------------
// UpdatePackage
// ---------------------------------------------------------------------------

func (s *Service) UpdatePackage(ctx context.Context, params *UpdatePackageParams) (*PackageDetail, error) {
	const op = "service.Service.UpdatePackage"

	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.user_id", params.UserID),
		attribute.String("input.package_id", params.ID),
	)
	logger.Info().Str("op", op).Str("user_id", params.UserID).Str("id", params.ID).Msg("")

	if params.UserID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.ID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("id is required"))
	}

	existing, err := s.store.GetPackageByIDForStaff(ctx, params.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("package %q not found", params.ID))
		}
		return nil, fmt.Errorf("get package: %w", postgres_store.WrapDBError(err))
	}

	// BL-CAT-008: publish-gate — travel-kind packages need airline_id before activating.
	if params.Status == "active" {
		if err := validateKindOnPublish(&existing, params.AirlineID); err != nil {
			return nil, err
		}
	}

	arg := sqlc.UpdatePackageFieldsParams{ID: params.ID}
	if params.Name != "" {
		arg.Name = pgtype.Text{String: params.Name, Valid: true}
	}
	if params.Description != "" {
		arg.Description = pgtype.Text{String: params.Description, Valid: true}
	}
	if params.CoverPhotoUrl != "" {
		arg.CoverPhotoUrl = pgtype.Text{String: params.CoverPhotoUrl, Valid: true}
	}
	if params.Highlights != nil {
		arg.Highlights = params.Highlights
	}
	if params.ItineraryID != "" {
		arg.ItineraryID = pgtype.Text{String: params.ItineraryID, Valid: true}
	}
	if params.AirlineID != "" {
		arg.AirlineID = pgtype.Text{String: params.AirlineID, Valid: true}
	}
	if params.MuthawwifID != "" {
		arg.MuthawwifID = pgtype.Text{String: params.MuthawwifID, Valid: true}
	}
	if params.Status != "" {
		arg.Status = sqlc.NullCatalogPackageStatus{
			CatalogPackageStatus: sqlc.CatalogPackageStatus(params.Status),
			Valid:                true,
		}
	}

	if _, err := s.store.UpdatePackageFields(ctx, arg); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("package %q not found or deleted", params.ID))
		}
		return nil, fmt.Errorf("update package: %w", postgres_store.WrapDBError(err))
	}

	if params.HotelIDs != nil {
		if err := s.store.DeletePackageHotels(ctx, params.ID); err != nil {
			return nil, fmt.Errorf("clear hotels: %w", postgres_store.WrapDBError(err))
		}
		for i, hid := range params.HotelIDs {
			if err := s.store.InsertPackageHotel(ctx, sqlc.InsertPackageHotelParams{
				PackageID: params.ID, HotelID: hid, SortOrder: int16(i),
			}); err != nil {
				return nil, fmt.Errorf("link hotel %s: %w", hid, postgres_store.WrapDBError(err))
			}
		}
	}
	if params.AddonIDs != nil {
		// BL-CAT-009: validate all addon IDs before performing any writes.
		for _, aid := range params.AddonIDs {
			if _, err := s.store.GetAddonByID(ctx, aid); err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("addon_id %q does not exist in catalog", aid))
				}
				return nil, fmt.Errorf("lookup addon %s: %w", aid, postgres_store.WrapDBError(err))
			}
		}
		if err := s.store.DeletePackageAddons(ctx, params.ID); err != nil {
			return nil, fmt.Errorf("clear addons: %w", postgres_store.WrapDBError(err))
		}
		for _, aid := range params.AddonIDs {
			if err := s.store.InsertPackageAddon(ctx, sqlc.InsertPackageAddonParams{
				PackageID: params.ID, AddonID: aid,
			}); err != nil {
				return nil, fmt.Errorf("link addon %s: %w", aid, postgres_store.WrapDBError(err))
			}
		}
	}

	detail, err := s.hydratePackageDetailForStaff(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")

	// BL-CAT-012: emit catalog.updated webhook (best-effort).
	s.emitCatalogUpdated(ctx, "package.updated", params.ID, "package")

	return detail, nil
}

// ---------------------------------------------------------------------------
// DeletePackage
// ---------------------------------------------------------------------------

func (s *Service) DeletePackage(ctx context.Context, params *DeletePackageParams) error {
	const op = "service.Service.DeletePackage"

	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.user_id", params.UserID),
		attribute.String("input.package_id", params.ID),
	)
	logger.Info().Str("op", op).Str("user_id", params.UserID).Str("id", params.ID).Msg("")

	if params.UserID == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.ID == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("id is required"))
	}

	row, err := s.store.SoftDeletePackage(ctx, params.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.Join(apperrors.ErrNotFound, fmt.Errorf("package %q not found or already deleted", params.ID))
		}
		return fmt.Errorf("soft delete package: %w", postgres_store.WrapDBError(err))
	}

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(attribute.String("output.status", string(row.Status)))

	// BL-CAT-012: emit catalog.updated webhook (best-effort).
	s.emitCatalogUpdated(ctx, "package.deleted", params.ID, "package")

	return nil
}

// ---------------------------------------------------------------------------
// CreateDeparture
// ---------------------------------------------------------------------------

func (s *Service) CreateDeparture(ctx context.Context, params *CreateDepartureParams) (*DepartureDetail, error) {
	const op = "service.Service.CreateDeparture"

	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.user_id", params.UserID),
		attribute.String("input.package_id", params.PackageID),
	)
	logger.Info().Str("op", op).Str("user_id", params.UserID).Str("pkg", params.PackageID).Msg("")

	if params.UserID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.PackageID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("package_id is required"))
	}
	if params.DepartureDate == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("departure_date is required"))
	}
	if params.ReturnDate == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("return_date is required"))
	}
	if params.TotalSeats <= 0 {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("total_seats must be > 0"))
	}

	depDate, err := parseDate(params.DepartureDate)
	if err != nil {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("departure_date: %w", err))
	}
	retDate, err := parseDate(params.ReturnDate)
	if err != nil {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("return_date: %w", err))
	}

	depStatus := sqlc.CatalogDepartureStatusOpen
	switch params.Status {
	case "", "open":
		depStatus = sqlc.CatalogDepartureStatusOpen
	case "closed":
		depStatus = sqlc.CatalogDepartureStatusClosed
	default:
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("status must be 'open' or 'closed' on create; got %q", params.Status))
	}

	if _, err := s.store.GetPackageByIDForStaff(ctx, params.PackageID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("package %q not found", params.PackageID))
		}
		return nil, fmt.Errorf("get package: %w", postgres_store.WrapDBError(err))
	}

	id, err := ulid.New("dep_")
	if err != nil {
		return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("mint departure id: %w", err))
	}

	row, err := s.store.InsertDeparture(ctx, sqlc.InsertDepartureParams{
		ID:            id,
		PackageID:     params.PackageID,
		DepartureDate: depDate,
		ReturnDate:    retDate,
		TotalSeats:    int32(params.TotalSeats),
		Status:        depStatus,
	})
	if err != nil {
		return nil, fmt.Errorf("insert departure: %w", postgres_store.WrapDBError(err))
	}

	pricing, err := s.replaceDeparturePricing(ctx, id, params.Pricing)
	if err != nil {
		return nil, err
	}

	depDateVal, _ := row.DepartureDate.Value()
	retDateVal, _ := row.ReturnDate.Value()
	detail := &DepartureDetail{
		ID:             row.ID,
		PackageID:      row.PackageID,
		DepartureDate:  dateToISO(depDateVal),
		ReturnDate:     dateToISO(retDateVal),
		TotalSeats:     int(row.TotalSeats),
		RemainingSeats: int(row.RemainingSeats),
		Status:         string(row.Status),
		Pricing:        pricing,
		VendorReadiness: VendorReadiness{
			Ticket: "not_started",
			Hotel:  "not_started",
			Visa:   "not_started",
		},
	}

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(attribute.String("output.departure_id", id))

	// BL-CAT-012: emit catalog.updated webhook (best-effort).
	s.emitCatalogUpdated(ctx, "departure.created", id, "departure")

	return detail, nil
}

// ---------------------------------------------------------------------------
// UpdateDeparture
// ---------------------------------------------------------------------------

func (s *Service) UpdateDeparture(ctx context.Context, params *UpdateDepartureParams) (*DepartureDetail, error) {
	const op = "service.Service.UpdateDeparture"

	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.user_id", params.UserID),
		attribute.String("input.departure_id", params.ID),
	)
	logger.Info().Str("op", op).Str("user_id", params.UserID).Str("id", params.ID).Msg("")

	if params.UserID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.ID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("id is required"))
	}

	if params.Status != "" {
		allowed := map[string]bool{"open": true, "closed": true, "cancelled": true}
		if !allowed[params.Status] {
			return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf(
				"invalid status transition target %q; allowed: open, closed, cancelled", params.Status))
		}
	}

	arg := sqlc.UpdateDepartureFieldsParams{ID: params.ID}
	if params.DepartureDate != "" {
		d, err := parseDate(params.DepartureDate)
		if err != nil {
			return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("departure_date: %w", err))
		}
		arg.DepartureDate = d
	}
	if params.ReturnDate != "" {
		d, err := parseDate(params.ReturnDate)
		if err != nil {
			return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("return_date: %w", err))
		}
		arg.ReturnDate = d
	}
	if params.TotalSeats > 0 {
		arg.TotalSeats = pgtype.Int4{Int32: int32(params.TotalSeats), Valid: true}
	}
	if params.Status != "" {
		arg.Status = sqlc.NullCatalogDepartureStatus{
			CatalogDepartureStatus: sqlc.CatalogDepartureStatus(params.Status),
			Valid:                  true,
		}
	}

	row, err := s.store.UpdateDepartureFields(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("departure %q not found", params.ID))
		}
		return nil, fmt.Errorf("update departure: %w", postgres_store.WrapDBError(err))
	}

	var pricing []PackagePricing
	if params.Pricing != nil {
		pricing, err = s.replaceDeparturePricing(ctx, params.ID, params.Pricing)
		if err != nil {
			return nil, err
		}
	} else {
		pRows, err := s.store.ListPricingForDeparture(ctx, params.ID)
		if err != nil {
			return nil, fmt.Errorf("list pricing: %w", postgres_store.WrapDBError(err))
		}
		for _, p := range pRows {
			pricing = append(pricing, PackagePricing{
				RoomType:           string(p.RoomType),
				ListAmount:         p.ListAmount,
				ListCurrency:       strings.TrimSpace(p.ListCurrency),
				SettlementCurrency: strings.TrimSpace(p.SettlementCurrency),
			})
		}
	}

	depDateVal, _ := row.DepartureDate.Value()
	retDateVal, _ := row.ReturnDate.Value()
	detail := &DepartureDetail{
		ID:             row.ID,
		PackageID:      row.PackageID,
		DepartureDate:  dateToISO(depDateVal),
		ReturnDate:     dateToISO(retDateVal),
		TotalSeats:     int(row.TotalSeats),
		RemainingSeats: int(row.RemainingSeats),
		Status:         string(row.Status),
		Pricing:        pricing,
		VendorReadiness: VendorReadiness{
			Ticket: "not_started",
			Hotel:  "not_started",
			Visa:   "not_started",
		},
	}

	span.SetStatus(codes.Ok, "success")

	// BL-CAT-012: emit catalog.updated webhook (best-effort).
	s.emitCatalogUpdated(ctx, "departure.updated", params.ID, "departure")

	return detail, nil
}

// ---------------------------------------------------------------------------
// ReserveSeats (§ Inventory / S1-J-03)
// ---------------------------------------------------------------------------

func (s *Service) ReserveSeats(ctx context.Context, params *ReserveSeatsParams) (*ReserveSeatsResult, error) {
	const op = "service.Service.ReserveSeats"

	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.reservation_id", params.ReservationID),
		attribute.String("input.departure_id", params.DepartureID),
		attribute.Int("input.seats", params.Seats),
	)

	if params.ReservationID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("reservation_id is required"))
	}
	if params.DepartureID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("departure_id is required"))
	}
	if params.Seats < 1 {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("seats must be >= 1"))
	}

	ttlHours := params.IdempotencyTTLHours
	if ttlHours < 1 {
		ttlHours = 24
	}
	if ttlHours > 168 {
		ttlHours = 168
	}

	// --- Idempotency check ---
	existing, err := s.store.GetSeatReservation(ctx, params.ReservationID)
	if err == nil {
		// Existing dedup row.
		if existing.DepartureID != params.DepartureID || int(existing.Seats) != params.Seats {
			return nil, errors.Join(apperrors.ErrConflict, fmt.Errorf(
				"reservation_id_conflict: reservation %q was previously used for departure=%s seats=%d",
				params.ReservationID, existing.DepartureID, existing.Seats))
		}
		dep, depErr := s.store.GetDepartureByIDForStaff(ctx, params.DepartureID)
		if depErr != nil {
			return nil, fmt.Errorf("get departure for replay: %w", postgres_store.WrapDBError(depErr))
		}
		reservedAt, expiresAt := tsToRFC3339(existing.ReservedAt), tsToRFC3339(existing.ExpiresAt)
		logger.Info().Str("op", op).Str("reservation_id", params.ReservationID).Msg("replay")
		span.SetStatus(codes.Ok, "replayed")
		return &ReserveSeatsResult{
			ReservationID:  params.ReservationID,
			DepartureID:    params.DepartureID,
			Seats:          params.Seats,
			ReservedAt:     reservedAt,
			ExpiresAt:      expiresAt,
			RemainingSeats: int(dep.RemainingSeats),
			Replayed:       true,
		}, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("check reservation: %w", postgres_store.WrapDBError(err))
	}

	// --- Atomic decrement + dedup write inside one transaction ---
	type txOut struct {
		remaining int32
		reserved  sqlc.SeatReservationRow
	}
	var txr txOut

	_, txErr := s.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(qtx *sqlc.Queries) error {
			updated, err := qtx.ReserveSeatsAtomic(ctx, sqlc.ReserveSeatsAtomicParams{
				Seats:       int32(params.Seats),
				DepartureID: params.DepartureID,
			})
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					dep, lookupErr := qtx.GetDepartureByIDForStaff(ctx, params.DepartureID)
					if lookupErr != nil {
						if errors.Is(lookupErr, pgx.ErrNoRows) {
							return errors.Join(apperrors.ErrNotFound, fmt.Errorf("departure %q not found", params.DepartureID))
						}
						return fmt.Errorf("lookup departure: %w", lookupErr)
					}
					_ = dep
					return errors.Join(apperrors.ErrConflict, fmt.Errorf("insufficient_capacity: not enough seats on departure %q", params.DepartureID))
				}
				return fmt.Errorf("reserve atomic: %w", err)
			}
			expiresAt := time.Now().UTC().Add(time.Duration(ttlHours) * time.Hour)
			res, err := qtx.InsertSeatReservation(ctx, sqlc.InsertSeatReservationParams{
				ReservationID: params.ReservationID,
				DepartureID:   params.DepartureID,
				Seats:         int32(params.Seats),
				ExpiresAt:     expiresAt,
			})
			if err != nil {
				return fmt.Errorf("insert reservation: %w", err)
			}
			txr.remaining = updated.RemainingSeats
			txr.reserved = res
			return nil
		},
	})
	if txErr != nil {
		return nil, txErr
	}

	reservedAt := tsToRFC3339(txr.reserved.ReservedAt)
	expiresAt := tsToRFC3339(txr.reserved.ExpiresAt)

	span.SetStatus(codes.Ok, "success")
	return &ReserveSeatsResult{
		ReservationID:  params.ReservationID,
		DepartureID:    params.DepartureID,
		Seats:          params.Seats,
		ReservedAt:     reservedAt,
		ExpiresAt:      expiresAt,
		RemainingSeats: int(txr.remaining),
		Replayed:       false,
	}, nil
}

// ---------------------------------------------------------------------------
// ReleaseSeats (§ Inventory / S1-J-03)
// ---------------------------------------------------------------------------

func (s *Service) ReleaseSeats(ctx context.Context, params *ReleaseSeatsParams) (*ReleaseSeatsResult, error) {
	const op = "service.Service.ReleaseSeats"

	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.reservation_id", params.ReservationID),
	)

	if params.ReservationID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("reservation_id is required"))
	}

	res, err := s.store.GetSeatReservation(ctx, params.ReservationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("reservation %q not found", params.ReservationID))
		}
		return nil, fmt.Errorf("get reservation: %w", postgres_store.WrapDBError(err))
	}

	// Already fully released — replay.
	if res.ReleasedAt.Valid {
		dep, err := s.store.GetDepartureByIDForStaff(ctx, res.DepartureID)
		if err != nil {
			return nil, fmt.Errorf("get departure: %w", postgres_store.WrapDBError(err))
		}
		releasedAt := tsToRFC3339(res.ReleasedAt)
		logger.Info().Str("op", op).Str("reservation_id", params.ReservationID).Msg("release replay")
		span.SetStatus(codes.Ok, "replayed")
		return &ReleaseSeatsResult{
			ReservationID:  params.ReservationID,
			DepartureID:    res.DepartureID,
			SeatsReleased:  int(res.Seats),
			ReleasedAt:     releasedAt,
			RemainingSeats: int(dep.RemainingSeats),
			Replayed:       true,
		}, nil
	}

	seatsToRelease := int(res.Seats)
	if params.Seats > 0 {
		if params.Seats > int(res.Seats) {
			return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf(
				"partial release seats %d exceeds original reservation seats %d", params.Seats, res.Seats))
		}
		seatsToRelease = params.Seats
	}

	type txOut struct {
		remaining  int32
		releasedAt string
	}
	var txr txOut

	_, txErr := s.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(qtx *sqlc.Queries) error {
			updated, err := qtx.ReleaseSeatsAtomic(ctx, sqlc.ReleaseSeatsAtomicParams{
				SeatsToRelease: int32(seatsToRelease),
				DepartureID:    res.DepartureID,
			})
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return errors.Join(apperrors.ErrConflict, fmt.Errorf(
						"reservation_not_active: could not release %d seats from departure %q",
						seatsToRelease, res.DepartureID))
				}
				return fmt.Errorf("release atomic: %w", err)
			}
			marked, markErr := qtx.MarkReservationReleased(ctx, params.ReservationID)
			if markErr != nil {
				logger.Warn().Err(markErr).Str("reservation_id", params.ReservationID).Msg("mark released failed")
			}
			txr.remaining = updated.RemainingSeats
			if marked.ReleasedAt.Valid {
				txr.releasedAt = tsToRFC3339(marked.ReleasedAt)
			}
			return nil
		},
	})
	if txErr != nil {
		return nil, txErr
	}

	span.SetStatus(codes.Ok, "success")
	return &ReleaseSeatsResult{
		ReservationID:  params.ReservationID,
		DepartureID:    res.DepartureID,
		SeatsReleased:  seatsToRelease,
		ReleasedAt:     txr.releasedAt,
		RemainingSeats: int(txr.remaining),
		Replayed:       false,
	}, nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// hydratePackageDetailForStaff loads a full PackageDetail for any non-deleted
// package regardless of status. Used by write-path response hydration.
func (s *Service) hydratePackageDetailForStaff(ctx context.Context, id string) (*PackageDetail, error) {
	row, err := s.store.GetPackageByIDForStaff(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("package %q not found", id))
		}
		return nil, fmt.Errorf("get package: %w", postgres_store.WrapDBError(err))
	}

	detail := &PackageDetail{
		ID:            row.ID,
		Kind:          string(row.Kind),
		Name:          row.Name,
		Description:   row.Description,
		Highlights:    row.Highlights,
		CoverPhotoUrl: row.CoverPhotoUrl,
		Status:        string(row.Status),
		Hotels:        []HotelRef{},
		Addons:        []AddonRef{},
		Departures:    []DepartureSummary{},
	}
	if detail.Highlights == nil {
		detail.Highlights = []string{}
	}

	if row.ItineraryID.Valid {
		it, err := s.store.GetItineraryByID(ctx, row.ItineraryID.String)
		if err == nil {
			detail.Itinerary = &Itinerary{ID: it.ID, Days: []ItineraryDay{}, PublicUrl: it.PublicUrl}
		}
	}
	if row.AirlineID.Valid {
		a, err := s.store.GetAirlineByID(ctx, row.AirlineID.String)
		if err == nil {
			detail.Airline = &AirlineRef{ID: a.ID, Code: a.Code, Name: a.Name, OperatorKind: string(a.OperatorKind)}
		}
	}
	if row.MuthawwifID.Valid {
		m, err := s.store.GetMuthawwifByID(ctx, row.MuthawwifID.String)
		if err == nil {
			detail.Muthawwif = &MuthawwifRef{ID: m.ID, Name: m.Name, PortraitUrl: m.PortraitUrl}
		}
	}

	hotels, _ := s.store.ListHotelsForPackage(ctx, id)
	for _, h := range hotels {
		detail.Hotels = append(detail.Hotels, HotelRef{
			ID: h.ID, Name: h.Name, City: h.City,
			StarRating: int(h.StarRating), WalkingDistanceM: int(h.WalkingDistanceM),
		})
	}

	addons, _ := s.store.ListAddonsForPackage(ctx, id)
	for _, a := range addons {
		detail.Addons = append(detail.Addons, AddonRef{
			ID:                 a.ID,
			Name:               a.Name,
			ListAmount:         a.ListAmount,
			ListCurrency:       strings.TrimSpace(a.ListCurrency),
			SettlementCurrency: strings.TrimSpace(a.SettlementCurrency),
		})
	}

	deps, _ := s.store.ListOpenDeparturesForPackage(ctx, id)
	for _, d := range deps {
		depDate, _ := d.DepartureDate.Value()
		retDate, _ := d.ReturnDate.Value()
		detail.Departures = append(detail.Departures, DepartureSummary{
			ID:             d.ID,
			DepartureDate:  dateToISO(depDate),
			ReturnDate:     dateToISO(retDate),
			RemainingSeats: int(d.RemainingSeats),
			Status:         string(d.Status),
		})
	}

	return detail, nil
}

// replaceDeparturePricing removes all existing pricing rows and inserts fresh ones.
func (s *Service) replaceDeparturePricing(ctx context.Context, departureID string, inputs []PricingInput) ([]PackagePricing, error) {
	if len(inputs) == 0 {
		return []PackagePricing{}, nil
	}
	if err := s.store.DeleteDeparturePricing(ctx, departureID); err != nil {
		return nil, fmt.Errorf("delete pricing: %w", postgres_store.WrapDBError(err))
	}
	out := make([]PackagePricing, 0, len(inputs))
	for _, p := range inputs {
		currency := p.ListCurrency
		if currency == "" {
			currency = "IDR"
		}
		settlement := p.SettlementCurrency
		if settlement == "" {
			settlement = "IDR"
		}
		priceID, err := ulid.New("pkgpr_")
		if err != nil {
			return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("mint pricing id: %w", err))
		}
		row, err := s.store.InsertDeparturePricing(ctx, sqlc.InsertDeparturePricingParams{
			ID:                 priceID,
			PackageDepartureID: departureID,
			RoomType:           sqlc.CatalogRoomType(p.RoomType),
			ListAmount:         p.ListAmount,
			ListCurrency:       currency,
			SettlementCurrency: settlement,
		})
		if err != nil {
			return nil, fmt.Errorf("insert pricing: %w", postgres_store.WrapDBError(err))
		}
		out = append(out, PackagePricing{
			RoomType:           string(row.RoomType),
			ListAmount:         row.ListAmount,
			ListCurrency:       strings.TrimSpace(row.ListCurrency),
			SettlementCurrency: strings.TrimSpace(row.SettlementCurrency),
		})
	}
	return out, nil
}

// optText converts a string to pgtype.Text (Valid=true when non-empty).
func optText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

// tsToRFC3339 formats a pgtype.Timestamptz to RFC 3339 string or "".
func tsToRFC3339(ts pgtype.Timestamptz) string {
	if !ts.Valid {
		return ""
	}
	return ts.Time.UTC().Format(time.RFC3339)
}
