package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"catalog-svc/store/postgres_store"
	"catalog-svc/store/postgres_store/sqlc"
	"catalog-svc/util/apperrors"
	"catalog-svc/util/logging"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Domain types — wire-adjacent but transport-agnostic.
//
// These types mirror the § Catalog JSON shapes so the REST adapter is a
// 1:1 mapping. They are NOT the sqlc generated rows: the service layer
// translates pgtype.* fields into plain Go types and derives the
// public-read semantics (e.g. suppressing `next_departure` when no
// upcoming departure exists).
// ---------------------------------------------------------------------------

// Money is the display-only currency triple from § Catalog.
//
// - `ListAmount` is the integer figure in whole units of `ListCurrency`.
// - `SettlementCurrency` is always "IDR" in MVP per Q001.
type Money struct {
	ListAmount         int64  `json:"list_amount"`
	ListCurrency       string `json:"list_currency"`
	SettlementCurrency string `json:"settlement_currency"`
}

// NextDeparture is the single "earliest upcoming open/closed" departure
// attached to each list-row package. Omitted when absent.
type NextDeparture struct {
	ID             string `json:"id"`
	DepartureDate  string `json:"departure_date"` // ISO date YYYY-MM-DD
	ReturnDate     string `json:"return_date"`
	RemainingSeats int    `json:"remaining_seats"`
}

// Package is a list-item shape: lean fields + optional starting price +
// optional next departure.
type Package struct {
	ID            string         `json:"id"`
	Kind          string         `json:"kind"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	CoverPhotoUrl string         `json:"cover_photo_url"`
	StartingPrice Money          `json:"starting_price"`
	NextDeparture *NextDeparture `json:"next_departure,omitempty"`
}

// GetPackagesParams is the service input for the list endpoint. All
// filter fields are optional; empty/zero values are treated as "no
// filter" by the store layer.
type GetPackagesParams struct {
	Kind          string // one of the catalog.package_kind enum values; "" = all
	AirlineCode   string
	HotelID       string
	DepartureFrom string // ISO date; "" = no lower bound
	DepartureTo   string // ISO date; "" = no upper bound
	Cursor        string // opaque; may be ""
	Limit         int    // 1..100; 0 = default 20
}

// GetPackagesResult is the shape returned to the REST layer.
type GetPackagesResult struct {
	Packages   []Package
	NextCursor string
	HasMore    bool
}

// --- Detail shapes -------------------------------------------------------

type HotelRef struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	City             string `json:"city"`
	StarRating       int    `json:"star_rating"`
	WalkingDistanceM int    `json:"walking_distance_m"`
}

type AirlineRef struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	OperatorKind string `json:"operator_kind"`
}

type MuthawwifRef struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PortraitUrl string `json:"portrait_url"`
}

type AddonRef struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	ListAmount         int64  `json:"list_amount"`
	ListCurrency       string `json:"list_currency"`
	SettlementCurrency string `json:"settlement_currency"`
}

type ItineraryDay struct {
	Day         int    `json:"day"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PhotoUrl    string `json:"photo_url,omitempty"`
}

type Itinerary struct {
	ID        string         `json:"id"`
	Days      []ItineraryDay `json:"days"`
	PublicUrl string         `json:"public_url"`
}

type DepartureSummary struct {
	ID             string `json:"id"`
	DepartureDate  string `json:"departure_date"`
	ReturnDate     string `json:"return_date"`
	RemainingSeats int    `json:"remaining_seats"`
	Status         string `json:"status"`
}

type PackageDetail struct {
	ID            string             `json:"id"`
	Kind          string             `json:"kind"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	Highlights    []string           `json:"highlights"`
	CoverPhotoUrl string             `json:"cover_photo_url"`
	// Status is populated for staff write responses; empty for public read
	// (public read only surfaces active packages so status is implicit).
	Status        string             `json:"status,omitempty"`
	Itinerary     *Itinerary         `json:"itinerary,omitempty"`
	Airline       *AirlineRef        `json:"airline,omitempty"`
	Muthawwif     *MuthawwifRef      `json:"muthawwif,omitempty"`
	Hotels        []HotelRef         `json:"hotels"`
	Addons        []AddonRef         `json:"add_ons"`
	Departures    []DepartureSummary `json:"departures"`
}

type GetPackageByIDParams struct {
	ID string
}

// --- Departure detail shapes -----------------------------------------------

type PackagePricing struct {
	RoomType           string `json:"room_type"`
	ListAmount         int64  `json:"list_amount"`
	ListCurrency       string `json:"list_currency"`
	SettlementCurrency string `json:"settlement_currency"`
}

type VendorReadiness struct {
	Ticket string `json:"ticket"`
	Hotel  string `json:"hotel"`
	Visa   string `json:"visa"`
}

type DepartureDetail struct {
	ID              string           `json:"id"`
	PackageID       string           `json:"package_id"`
	DepartureDate   string           `json:"departure_date"` // ISO date YYYY-MM-DD
	ReturnDate      string           `json:"return_date"`
	TotalSeats      int              `json:"total_seats"`
	RemainingSeats  int              `json:"remaining_seats"`
	Status          string           `json:"status"`
	Pricing         []PackagePricing `json:"pricing"`
	VendorReadiness VendorReadiness  `json:"vendor_readiness"`
}

type GetDepartureByIDParams struct {
	ID string
}

// ---------------------------------------------------------------------------
// Service-layer extensions
// ---------------------------------------------------------------------------

// PackagesService is the subset of IService implementing § Catalog
// endpoints. Declared separately for clarity; Service still satisfies
// the consolidated IService interface below via the same receiver.
type PackagesService interface {
	GetPackages(ctx context.Context, params *GetPackagesParams) (*GetPackagesResult, error)
	GetPackageByID(ctx context.Context, params *GetPackageByIDParams) (*PackageDetail, error)
	GetDepartureByID(ctx context.Context, params *GetDepartureByIDParams) (*DepartureDetail, error)
}

const (
	listDefaultLimit = 20
	listMaxLimit     = 100
)

// GetPackages returns an active-filtered, cursor-paginated page of
// catalog packages. Draft/archived packages are filtered out by the
// store query and never reach this method.
func (s *Service) GetPackages(ctx context.Context, params *GetPackagesParams) (*GetPackagesResult, error) {
	const op = "service.Service.GetPackages"

	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.params", fmt.Sprintf("%+v", params)),
	)
	logger.Info().Str("op", op).Str("params", fmt.Sprintf("%+v", params)).Msg("")

	limit := params.Limit
	if limit <= 0 {
		limit = listDefaultLimit
	}
	if limit > listMaxLimit {
		limit = listMaxLimit
	}

	cursorID, err := decodeCursor(params.Cursor)
	if err != nil {
		// Treat an unparseable cursor as a caller validation error.
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("decode cursor: %w", err))
	}

	storeArgs := sqlc.ListActivePackagesParams{
		RowLimit: int32(limit + 1), // one extra to detect has_more
	}
	if params.Kind != "" {
		storeArgs.Kind = sqlc.NullCatalogPackageKind{
			CatalogPackageKind: sqlc.CatalogPackageKind(params.Kind),
			Valid:              true,
		}
	}
	if params.AirlineCode != "" {
		storeArgs.AirlineCode = pgtype.Text{String: params.AirlineCode, Valid: true}
	}
	if params.HotelID != "" {
		storeArgs.HotelID = pgtype.Text{String: params.HotelID, Valid: true}
	}
	if params.DepartureFrom != "" {
		d, err := parseDate(params.DepartureFrom)
		if err != nil {
			return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("departure_from: %w", err))
		}
		storeArgs.DepartureFrom = d
	}
	if params.DepartureTo != "" {
		d, err := parseDate(params.DepartureTo)
		if err != nil {
			return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("departure_to: %w", err))
		}
		storeArgs.DepartureTo = d
	}
	if cursorID != "" {
		storeArgs.CursorID = pgtype.Text{String: cursorID, Valid: true}
	}

	rows, err := s.store.ListActivePackages(ctx, storeArgs)
	if err != nil {
		return nil, fmt.Errorf("list active packages: %w", postgres_store.WrapDBError(err))
	}

	hasMore := false
	if len(rows) > limit {
		hasMore = true
		rows = rows[:limit]
	}

	packages := make([]Package, 0, len(rows))
	for _, r := range rows {
		pkg := Package{
			ID:            r.ID,
			Kind:          string(r.Kind),
			Name:          r.Name,
			Description:   r.Description,
			CoverPhotoUrl: r.CoverPhotoUrl,
			StartingPrice: Money{
				ListAmount:         0,
				ListCurrency:       "IDR",
				SettlementCurrency: "IDR",
			},
		}

		// Starting price — optional (missing when the package has no
		// priced upcoming open/closed departure).
		price, err := s.store.GetStartingPriceForPackage(ctx, r.ID)
		switch {
		case err == nil:
			pkg.StartingPrice = Money{
				ListAmount:         price.ListAmount,
				ListCurrency:       strings.TrimSpace(price.ListCurrency),
				SettlementCurrency: strings.TrimSpace(price.SettlementCurrency),
			}
		case errors.Is(err, pgx.ErrNoRows):
			// keep the zero Money default
		default:
			return nil, fmt.Errorf("starting price for %s: %w", r.ID, postgres_store.WrapDBError(err))
		}

		// Next departure — optional.
		nd, err := s.store.GetNextDepartureForPackage(ctx, r.ID)
		switch {
		case err == nil:
			depDate, _ := nd.DepartureDate.Value()
			retDate, _ := nd.ReturnDate.Value()
			pkg.NextDeparture = &NextDeparture{
				ID:             nd.ID,
				DepartureDate:  dateToISO(depDate),
				ReturnDate:     dateToISO(retDate),
				RemainingSeats: int(nd.RemainingSeats),
			}
		case errors.Is(err, pgx.ErrNoRows):
			// keep pkg.NextDeparture as nil
		default:
			return nil, fmt.Errorf("next departure for %s: %w", r.ID, postgres_store.WrapDBError(err))
		}

		packages = append(packages, pkg)
	}

	result := &GetPackagesResult{
		Packages: packages,
		HasMore:  hasMore,
	}
	if hasMore && len(packages) > 0 {
		result.NextCursor = encodeCursor(packages[len(packages)-1].ID)
	}

	span.SetAttributes(
		attribute.Int("output.packages.count", len(packages)),
		attribute.Bool("output.has_more", hasMore),
	)
	span.SetStatus(codes.Ok, "success")

	return result, nil
}

// GetPackageByID fetches a single active package with eager-loaded
// master references and the list of upcoming open/closed departures.
// Returns apperrors.ErrNotFound for any non-active or unknown id —
// matching the § Catalog rule that the public endpoint does not leak
// the existence of draft/archived rows.
func (s *Service) GetPackageByID(ctx context.Context, params *GetPackageByIDParams) (*PackageDetail, error) {
	const op = "service.Service.GetPackageByID"

	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.id", params.ID),
	)
	logger.Info().Str("op", op).Str("id", params.ID).Msg("")

	row, err := s.store.GetActivePackageByID(ctx, params.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("package %q not active or not found", params.ID))
		}
		return nil, fmt.Errorf("get active package: %w", postgres_store.WrapDBError(err))
	}

	detail := &PackageDetail{
		ID:            row.ID,
		Kind:          string(row.Kind),
		Name:          row.Name,
		Description:   row.Description,
		Highlights:    row.Highlights,
		CoverPhotoUrl: row.CoverPhotoUrl,
		Hotels:        []HotelRef{},
		Addons:        []AddonRef{},
		Departures:    []DepartureSummary{},
	}
	if detail.Highlights == nil {
		detail.Highlights = []string{}
	}

	if row.ItineraryID.Valid {
		it, err := s.store.GetItineraryByID(ctx, row.ItineraryID.String)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get itinerary %s: %w", row.ItineraryID.String, postgres_store.WrapDBError(err))
		}
		if err == nil {
			days := []ItineraryDay{}
			if len(it.Days) > 0 {
				if err := json.Unmarshal(it.Days, &days); err != nil {
					return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("decode itinerary days: %w", err))
				}
			}
			detail.Itinerary = &Itinerary{
				ID:        it.ID,
				Days:      days,
				PublicUrl: it.PublicUrl,
			}
		}
	}

	if row.AirlineID.Valid {
		a, err := s.store.GetAirlineByID(ctx, row.AirlineID.String)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get airline %s: %w", row.AirlineID.String, postgres_store.WrapDBError(err))
		}
		if err == nil {
			detail.Airline = &AirlineRef{
				ID:           a.ID,
				Code:         a.Code,
				Name:         a.Name,
				OperatorKind: string(a.OperatorKind),
			}
		}
	}

	if row.MuthawwifID.Valid {
		m, err := s.store.GetMuthawwifByID(ctx, row.MuthawwifID.String)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get muthawwif %s: %w", row.MuthawwifID.String, postgres_store.WrapDBError(err))
		}
		if err == nil {
			detail.Muthawwif = &MuthawwifRef{
				ID:          m.ID,
				Name:        m.Name,
				PortraitUrl: m.PortraitUrl,
			}
		}
	}

	hotels, err := s.store.ListHotelsForPackage(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("list hotels: %w", postgres_store.WrapDBError(err))
	}
	for _, h := range hotels {
		detail.Hotels = append(detail.Hotels, HotelRef{
			ID:               h.ID,
			Name:             h.Name,
			City:             h.City,
			StarRating:       int(h.StarRating),
			WalkingDistanceM: int(h.WalkingDistanceM),
		})
	}

	addons, err := s.store.ListAddonsForPackage(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("list addons: %w", postgres_store.WrapDBError(err))
	}
	for _, a := range addons {
		detail.Addons = append(detail.Addons, AddonRef{
			ID:                 a.ID,
			Name:               a.Name,
			ListAmount:         a.ListAmount,
			ListCurrency:       strings.TrimSpace(a.ListCurrency),
			SettlementCurrency: strings.TrimSpace(a.SettlementCurrency),
		})
	}

	deps, err := s.store.ListOpenDeparturesForPackage(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("list departures: %w", postgres_store.WrapDBError(err))
	}
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

	span.SetStatus(codes.Ok, "success")
	return detail, nil
}

// GetDepartureByID fetches a single public-visible departure (status in
// {open, closed}) with its pricing rows and the real vendor readiness summary
// from catalog.departure_vendor_readiness. Returns apperrors.ErrNotFound with
// identical shape for unknown-id / hidden-status — no existence oracle per
// § Catalog.
//
// BL-OPS-020: vendor_readiness is now read from the DB; missing kinds default
// to "not_started".
func (s *Service) GetDepartureByID(ctx context.Context, params *GetDepartureByIDParams) (*DepartureDetail, error) {
	const op = "service.Service.GetDepartureByID"

	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.id", params.ID),
	)
	logger.Info().Str("op", op).Str("id", params.ID).Msg("")

	row, err := s.store.GetActiveDeparture(ctx, params.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("departure %q not visible or not found", params.ID))
		}
		return nil, fmt.Errorf("get active departure: %w", postgres_store.WrapDBError(err))
	}

	pricingRows, err := s.store.ListPricingForDeparture(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("list pricing: %w", postgres_store.WrapDBError(err))
	}
	pricing := make([]PackagePricing, 0, len(pricingRows))
	for _, p := range pricingRows {
		pricing = append(pricing, PackagePricing{
			RoomType:           string(p.RoomType),
			ListAmount:         p.ListAmount,
			ListCurrency:       strings.TrimSpace(p.ListCurrency),
			SettlementCurrency: strings.TrimSpace(p.SettlementCurrency),
		})
	}

	depDate, _ := row.DepartureDate.Value()
	retDate, _ := row.ReturnDate.Value()

	// Fetch real vendor readiness from DB (BL-OPS-020).
	readiness, err := s.getReadinessFromDB(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("get vendor readiness: %w", err)
	}

	detail := &DepartureDetail{
		ID:              row.ID,
		PackageID:       row.PackageID,
		DepartureDate:   dateToISO(depDate),
		ReturnDate:      dateToISO(retDate),
		TotalSeats:      int(row.TotalSeats),
		RemainingSeats:  int(row.RemainingSeats),
		Status:          string(row.Status),
		Pricing:         pricing,
		VendorReadiness: *readiness,
	}

	span.SetStatus(codes.Ok, "success")
	return detail, nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

type cursorPayload struct {
	LastID string `json:"last_id"`
}

func encodeCursor(lastID string) string {
	b, _ := json.Marshal(cursorPayload{LastID: lastID})
	return base64.URLEncoding.EncodeToString(b)
}

func decodeCursor(c string) (string, error) {
	if c == "" {
		return "", nil
	}
	raw, err := base64.URLEncoding.DecodeString(c)
	if err != nil {
		return "", fmt.Errorf("malformed cursor: %w", err)
	}
	var p cursorPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return "", fmt.Errorf("malformed cursor payload: %w", err)
	}
	if p.LastID == "" {
		return "", errors.New("cursor missing last_id")
	}
	return p.LastID, nil
}

func parseDate(s string) (pgtype.Date, error) {
	var d pgtype.Date
	if err := d.Scan(s); err != nil {
		return pgtype.Date{}, fmt.Errorf("parse date %q: %w", s, err)
	}
	return d, nil
}

// dateToISO converts whatever pgtype.Date.Value() returned (time.Time or
// nil) into a YYYY-MM-DD string. Empty string on nil.
func dateToISO(v any) string {
	if v == nil {
		return ""
	}
	// pgtype.Date.Value() returns time.Time; format as date-only.
	if t, ok := v.(interface{ Format(string) string }); ok {
		return t.Format("2006-01-02")
	}
	return fmt.Sprintf("%v", v)
}
