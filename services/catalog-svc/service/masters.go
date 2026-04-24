package service

// Wave-1A master data CRUD service methods.
//
// Covers hotel, airline, muthawwif, addon CRUD and departure pricing upsert.
//
// Permission gate: callers (gRPC handlers) call checkCatalogManagePermission
// before delegating here. Service layer trusts the handler has gated the call.
//
// ID generation: ULID with type prefix via util/ulid.New.

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
// Input / output types
// ---------------------------------------------------------------------------

// --- Shared ---

// DeleteMasterParams is the input for any master delete operation.
type DeleteMasterParams struct {
	UserID string
	ID     string
}

// ListMasterParams is the input for cursor-paginated list operations.
type ListMasterParams struct {
	UserID string
	Cursor string
	Limit  int // 0 → default 50; max 200
}

// --- Hotel ---

type CreateHotelParams struct {
	UserID           string
	Name             string
	City             string
	StarRating       int
	WalkingDistanceM int
}

type UpdateHotelParams struct {
	UserID           string
	ID               string
	Name             string // empty = no change
	City             string // empty = no change
	StarRating       int    // -1 = no change; 0 is a valid value so use -1 sentinel
	WalkingDistanceM int    // -1 = no change
}

type HotelResult struct {
	ID               string
	Name             string
	City             string
	StarRating       int
	WalkingDistanceM int
	CreatedAt        string
	UpdatedAt        string
}

type ListHotelsResult struct {
	Hotels  []*HotelResult
	HasMore bool
	Cursor  string
}

// --- Airline ---

type CreateAirlineParams struct {
	UserID       string
	Code         string
	Name         string
	OperatorKind string // "airline" | "rail" | "bus"; defaults to "airline"
}

type UpdateAirlineParams struct {
	UserID string
	ID     string
	Code   string // empty = no change
	Name   string // empty = no change
}

type AirlineResult struct {
	ID           string
	Code         string
	Name         string
	OperatorKind string
	CreatedAt    string
	UpdatedAt    string
}

type ListAirlinesResult struct {
	Airlines []*AirlineResult
	HasMore  bool
	Cursor   string
}

// --- Muthawwif ---

type CreateMuthawwifParams struct {
	UserID      string
	Name        string
	PortraitUrl string
}

type UpdateMuthawwifParams struct {
	UserID      string
	ID          string
	Name        string // empty = no change
	PortraitUrl string // empty = no change
}

type MuthawwifResult struct {
	ID          string
	Name        string
	PortraitUrl string
	CreatedAt   string
	UpdatedAt   string
}

type ListMuthawwifResult struct {
	Muthawwif []*MuthawwifResult
	HasMore   bool
	Cursor    string
}

// --- Addon ---

type CreateAddonParams struct {
	UserID        string
	Name          string
	ListAmountIDR int64
}

type UpdateAddonParams struct {
	UserID        string
	ID            string
	Name          string // empty = no change
	ListAmountIDR int64  // 0 = no change; -1 = set to 0
}

type AddonResult struct {
	ID            string
	Name          string
	ListAmountIDR int64
	CreatedAt     string
	UpdatedAt     string
}

type ListAddonsResult struct {
	Addons  []*AddonResult
	HasMore bool
	Cursor  string
}

// --- Pricing ---

type PricingUpsertInput struct {
	RoomType      string
	ListAmountIDR int64
}

type SetDeparturePricingParams struct {
	UserID      string
	DepartureID string
	Pricings    []PricingUpsertInput
}

type GetDeparturePricingParams struct {
	UserID      string
	DepartureID string
}

type PricingResult struct {
	ID          string
	DepartureID string
	RoomType    string
	ListAmount  int64
	CreatedAt   string
	UpdatedAt   string
}

// ---------------------------------------------------------------------------
// Hotel CRUD
// ---------------------------------------------------------------------------

func (s *Service) CreateHotel(ctx context.Context, params *CreateHotelParams) (*HotelResult, error) {
	const op = "service.Service.CreateHotel"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))
	logger.Info().Str("op", op).Str("user_id", params.UserID).Msg("")

	if params.UserID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.Name == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("name is required"))
	}
	if params.City == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("city is required"))
	}
	if params.StarRating < 0 || params.StarRating > 5 {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("star_rating must be between 0 and 5"))
	}

	id, err := ulid.New("htl_")
	if err != nil {
		return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("mint hotel id: %w", err))
	}

	row, err := s.store.InsertHotel(ctx, sqlc.InsertHotelParams{
		ID:               id,
		Name:             params.Name,
		City:             params.City,
		StarRating:       int16(params.StarRating),
		WalkingDistanceM: int32(params.WalkingDistanceM),
	})
	if err != nil {
		return nil, fmt.Errorf("insert hotel: %w", postgres_store.WrapDBError(err))
	}

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(attribute.String("output.hotel_id", id))
	return hotelRowToResult(row), nil
}

func (s *Service) UpdateHotel(ctx context.Context, params *UpdateHotelParams) (*HotelResult, error) {
	const op = "service.Service.UpdateHotel"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger.Info().Str("op", op).Str("id", params.ID).Msg("")

	if params.UserID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.ID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("id is required"))
	}

	arg := sqlc.UpdateHotelFieldsParams{ID: params.ID}
	if params.Name != "" {
		arg.Name = pgtype.Text{String: params.Name, Valid: true}
	}
	if params.City != "" {
		arg.City = pgtype.Text{String: params.City, Valid: true}
	}
	if params.StarRating >= 0 {
		arg.StarRating = pgtype.Int2{Int16: int16(params.StarRating), Valid: true}
	}
	if params.WalkingDistanceM >= 0 {
		arg.WalkingDistanceM = pgtype.Int4{Int32: int32(params.WalkingDistanceM), Valid: true}
	}

	row, err := s.store.UpdateHotelFields(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("hotel %q not found", params.ID))
		}
		return nil, fmt.Errorf("update hotel: %w", postgres_store.WrapDBError(err))
	}

	span.SetStatus(codes.Ok, "success")
	return hotelRowToResult(row), nil
}

func (s *Service) DeleteHotel(ctx context.Context, params *DeleteMasterParams) error {
	const op = "service.Service.DeleteHotel"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger.Info().Str("op", op).Str("id", params.ID).Msg("")

	if params.UserID == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.ID == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("id is required"))
	}

	// Existence check.
	if _, err := s.store.GetHotelByID(ctx, params.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.Join(apperrors.ErrNotFound, fmt.Errorf("hotel %q not found", params.ID))
		}
		return fmt.Errorf("get hotel: %w", postgres_store.WrapDBError(err))
	}

	// Ref check: block if still referenced by packages.
	refs, err := s.store.CountHotelPackageRefs(ctx, params.ID)
	if err != nil {
		return fmt.Errorf("count hotel refs: %w", postgres_store.WrapDBError(err))
	}
	if refs > 0 {
		return errors.Join(apperrors.ErrConflict, fmt.Errorf("hotel %q is referenced by %d package(s); remove those references first", params.ID, refs))
	}

	if err := s.store.DeleteHotel(ctx, params.ID); err != nil {
		return fmt.Errorf("delete hotel: %w", postgres_store.WrapDBError(err))
	}

	span.SetStatus(codes.Ok, "success")
	return nil
}

func (s *Service) ListHotels(ctx context.Context, params *ListMasterParams) (*ListHotelsResult, error) {
	const op = "service.Service.ListHotels"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	limit := clampLimit(params.Limit)
	rows, err := s.store.ListHotels(ctx, params.Cursor, int32(limit+1))
	if err != nil {
		return nil, fmt.Errorf("list hotels: %w", postgres_store.WrapDBError(err))
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	result := &ListHotelsResult{HasMore: hasMore}
	for _, r := range rows {
		result.Hotels = append(result.Hotels, hotelRowToResult(r))
	}
	if len(rows) > 0 {
		result.Cursor = rows[len(rows)-1].ID
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// ---------------------------------------------------------------------------
// Airline CRUD
// ---------------------------------------------------------------------------

func (s *Service) CreateAirline(ctx context.Context, params *CreateAirlineParams) (*AirlineResult, error) {
	const op = "service.Service.CreateAirline"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger.Info().Str("op", op).Str("user_id", params.UserID).Msg("")

	if params.UserID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.Code == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("code is required"))
	}
	if params.Name == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("name is required"))
	}

	opKind := sqlc.CatalogOperatorKindAirline
	switch params.OperatorKind {
	case "", "airline":
		opKind = sqlc.CatalogOperatorKindAirline
	case "rail":
		opKind = sqlc.CatalogOperatorKindRail
	case "bus":
		opKind = sqlc.CatalogOperatorKindBus
	default:
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("operator_kind must be 'airline', 'rail', or 'bus'; got %q", params.OperatorKind))
	}

	id, err := ulid.New("arl_")
	if err != nil {
		return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("mint airline id: %w", err))
	}

	row, err := s.store.InsertAirline(ctx, sqlc.InsertAirlineParams{
		ID:           id,
		Code:         strings.ToUpper(params.Code),
		Name:         params.Name,
		OperatorKind: opKind,
	})
	if err != nil {
		return nil, fmt.Errorf("insert airline: %w", postgres_store.WrapDBError(err))
	}

	span.SetStatus(codes.Ok, "success")
	return airlineRowToResult(row), nil
}

func (s *Service) UpdateAirline(ctx context.Context, params *UpdateAirlineParams) (*AirlineResult, error) {
	const op = "service.Service.UpdateAirline"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger.Info().Str("op", op).Str("id", params.ID).Msg("")

	if params.UserID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.ID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("id is required"))
	}

	arg := sqlc.UpdateAirlineFieldsParams{ID: params.ID}
	if params.Code != "" {
		arg.Code = pgtype.Text{String: strings.ToUpper(params.Code), Valid: true}
	}
	if params.Name != "" {
		arg.Name = pgtype.Text{String: params.Name, Valid: true}
	}

	row, err := s.store.UpdateAirlineFields(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("airline %q not found", params.ID))
		}
		return nil, fmt.Errorf("update airline: %w", postgres_store.WrapDBError(err))
	}

	span.SetStatus(codes.Ok, "success")
	return airlineRowToResult(row), nil
}

func (s *Service) DeleteAirline(ctx context.Context, params *DeleteMasterParams) error {
	const op = "service.Service.DeleteAirline"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger.Info().Str("op", op).Str("id", params.ID).Msg("")

	if params.UserID == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.ID == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("id is required"))
	}

	if _, err := s.store.GetAirlineByIDForStaff(ctx, params.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.Join(apperrors.ErrNotFound, fmt.Errorf("airline %q not found", params.ID))
		}
		return fmt.Errorf("get airline: %w", postgres_store.WrapDBError(err))
	}

	// Ref check: block if still referenced by packages.
	refs, err := s.store.CountAirlinePackageRefs(ctx, params.ID)
	if err != nil {
		return fmt.Errorf("count airline refs: %w", postgres_store.WrapDBError(err))
	}
	if refs > 0 {
		return errors.Join(apperrors.ErrConflict, fmt.Errorf("airline %q is referenced by %d package(s); remove those references first", params.ID, refs))
	}

	if err := s.store.DeleteAirline(ctx, params.ID); err != nil {
		return fmt.Errorf("delete airline: %w", postgres_store.WrapDBError(err))
	}

	span.SetStatus(codes.Ok, "success")
	return nil
}

func (s *Service) ListAirlines(ctx context.Context, params *ListMasterParams) (*ListAirlinesResult, error) {
	const op = "service.Service.ListAirlines"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	limit := clampLimit(params.Limit)
	rows, err := s.store.ListAirlines(ctx, params.Cursor, int32(limit+1))
	if err != nil {
		return nil, fmt.Errorf("list airlines: %w", postgres_store.WrapDBError(err))
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	result := &ListAirlinesResult{HasMore: hasMore}
	for _, r := range rows {
		result.Airlines = append(result.Airlines, airlineRowToResult(r))
	}
	if len(rows) > 0 {
		result.Cursor = rows[len(rows)-1].ID
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// ---------------------------------------------------------------------------
// Muthawwif CRUD
// ---------------------------------------------------------------------------

func (s *Service) CreateMuthawwif(ctx context.Context, params *CreateMuthawwifParams) (*MuthawwifResult, error) {
	const op = "service.Service.CreateMuthawwif"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger.Info().Str("op", op).Str("user_id", params.UserID).Msg("")

	if params.UserID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.Name == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("name is required"))
	}

	id, err := ulid.New("mtw_")
	if err != nil {
		return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("mint muthawwif id: %w", err))
	}

	row, err := s.store.InsertMuthawwif(ctx, sqlc.InsertMuthawwifParams{
		ID:          id,
		Name:        params.Name,
		PortraitUrl: params.PortraitUrl,
	})
	if err != nil {
		return nil, fmt.Errorf("insert muthawwif: %w", postgres_store.WrapDBError(err))
	}

	span.SetStatus(codes.Ok, "success")
	return muthawwifRowToResult(row), nil
}

func (s *Service) UpdateMuthawwif(ctx context.Context, params *UpdateMuthawwifParams) (*MuthawwifResult, error) {
	const op = "service.Service.UpdateMuthawwif"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger.Info().Str("op", op).Str("id", params.ID).Msg("")

	if params.UserID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.ID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("id is required"))
	}

	arg := sqlc.UpdateMuthawwifFieldsParams{ID: params.ID}
	if params.Name != "" {
		arg.Name = pgtype.Text{String: params.Name, Valid: true}
	}
	if params.PortraitUrl != "" {
		arg.PortraitUrl = pgtype.Text{String: params.PortraitUrl, Valid: true}
	}

	row, err := s.store.UpdateMuthawwifFields(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("muthawwif %q not found", params.ID))
		}
		return nil, fmt.Errorf("update muthawwif: %w", postgres_store.WrapDBError(err))
	}

	span.SetStatus(codes.Ok, "success")
	return muthawwifRowToResult(row), nil
}

func (s *Service) DeleteMuthawwif(ctx context.Context, params *DeleteMasterParams) error {
	const op = "service.Service.DeleteMuthawwif"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger.Info().Str("op", op).Str("id", params.ID).Msg("")

	if params.UserID == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.ID == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("id is required"))
	}

	if _, err := s.store.GetMuthawwifByIDForStaff(ctx, params.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.Join(apperrors.ErrNotFound, fmt.Errorf("muthawwif %q not found", params.ID))
		}
		return fmt.Errorf("get muthawwif: %w", postgres_store.WrapDBError(err))
	}

	// Ref check: block if still referenced by packages.
	refs, err := s.store.CountMuthawwifPackageRefs(ctx, params.ID)
	if err != nil {
		return fmt.Errorf("count muthawwif refs: %w", postgres_store.WrapDBError(err))
	}
	if refs > 0 {
		return errors.Join(apperrors.ErrConflict, fmt.Errorf("muthawwif %q is referenced by %d package(s); remove those references first", params.ID, refs))
	}

	if err := s.store.DeleteMuthawwif(ctx, params.ID); err != nil {
		return fmt.Errorf("delete muthawwif: %w", postgres_store.WrapDBError(err))
	}

	span.SetStatus(codes.Ok, "success")
	return nil
}

func (s *Service) ListMuthawwif(ctx context.Context, params *ListMasterParams) (*ListMuthawwifResult, error) {
	const op = "service.Service.ListMuthawwif"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	limit := clampLimit(params.Limit)
	rows, err := s.store.ListMuthawwif(ctx, params.Cursor, int32(limit+1))
	if err != nil {
		return nil, fmt.Errorf("list muthawwif: %w", postgres_store.WrapDBError(err))
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	result := &ListMuthawwifResult{HasMore: hasMore}
	for _, r := range rows {
		result.Muthawwif = append(result.Muthawwif, muthawwifRowToResult(r))
	}
	if len(rows) > 0 {
		result.Cursor = rows[len(rows)-1].ID
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// ---------------------------------------------------------------------------
// Addon CRUD
// ---------------------------------------------------------------------------

func (s *Service) CreateAddon(ctx context.Context, params *CreateAddonParams) (*AddonResult, error) {
	const op = "service.Service.CreateAddon"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger.Info().Str("op", op).Str("user_id", params.UserID).Msg("")

	if params.UserID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.Name == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("name is required"))
	}
	if params.ListAmountIDR < 0 {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("list_amount_idr must be >= 0"))
	}

	id, err := ulid.New("addon_")
	if err != nil {
		return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("mint addon id: %w", err))
	}

	row, err := s.store.InsertAddon(ctx, sqlc.InsertAddonParams{
		ID:         id,
		Name:       params.Name,
		ListAmount: params.ListAmountIDR,
	})
	if err != nil {
		return nil, fmt.Errorf("insert addon: %w", postgres_store.WrapDBError(err))
	}

	span.SetStatus(codes.Ok, "success")
	return addonRowToResult(row), nil
}

func (s *Service) UpdateAddon(ctx context.Context, params *UpdateAddonParams) (*AddonResult, error) {
	const op = "service.Service.UpdateAddon"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger.Info().Str("op", op).Str("id", params.ID).Msg("")

	if params.UserID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.ID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("id is required"))
	}

	arg := sqlc.UpdateAddonFieldsParams{ID: params.ID}
	if params.Name != "" {
		arg.Name = pgtype.Text{String: params.Name, Valid: true}
	}
	// ListAmountIDR == -1 means "set to 0"; 0 means no change.
	if params.ListAmountIDR > 0 {
		var amount pgtype.Numeric
		// Use scan to convert int64 -> pgtype.Numeric cleanly.
		if err := amount.Scan(params.ListAmountIDR); err != nil {
			return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("convert amount: %w", err))
		}
		arg.ListAmount = amount
	} else if params.ListAmountIDR == -1 {
		if err := arg.ListAmount.Scan(int64(0)); err != nil {
			return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("convert amount: %w", err))
		}
	}

	row, err := s.store.UpdateAddonFields(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("addon %q not found", params.ID))
		}
		return nil, fmt.Errorf("update addon: %w", postgres_store.WrapDBError(err))
	}

	span.SetStatus(codes.Ok, "success")
	return addonRowToResult(row), nil
}

func (s *Service) DeleteAddon(ctx context.Context, params *DeleteMasterParams) error {
	const op = "service.Service.DeleteAddon"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger.Info().Str("op", op).Str("id", params.ID).Msg("")

	if params.UserID == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.ID == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("id is required"))
	}

	if _, err := s.store.GetAddonByID(ctx, params.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.Join(apperrors.ErrNotFound, fmt.Errorf("addon %q not found", params.ID))
		}
		return fmt.Errorf("get addon: %w", postgres_store.WrapDBError(err))
	}

	if err := s.store.DeleteAddon(ctx, params.ID); err != nil {
		return fmt.Errorf("delete addon: %w", postgres_store.WrapDBError(err))
	}

	span.SetStatus(codes.Ok, "success")
	return nil
}

func (s *Service) ListAddons(ctx context.Context, params *ListMasterParams) (*ListAddonsResult, error) {
	const op = "service.Service.ListAddons"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	limit := clampLimit(params.Limit)
	rows, err := s.store.ListAddons(ctx, params.Cursor, int32(limit+1))
	if err != nil {
		return nil, fmt.Errorf("list addons: %w", postgres_store.WrapDBError(err))
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	result := &ListAddonsResult{HasMore: hasMore}
	for _, r := range rows {
		result.Addons = append(result.Addons, addonRowToResult(r))
	}
	if len(rows) > 0 {
		result.Cursor = rows[len(rows)-1].ID
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// ---------------------------------------------------------------------------
// Departure Pricing
// ---------------------------------------------------------------------------

func (s *Service) SetDeparturePricing(ctx context.Context, params *SetDeparturePricingParams) ([]*PricingResult, error) {
	const op = "service.Service.SetDeparturePricing"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger.Info().Str("op", op).Str("departure_id", params.DepartureID).Msg("")

	if params.UserID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}
	if params.DepartureID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("departure_id is required"))
	}
	if len(params.Pricings) == 0 {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("pricings must not be empty"))
	}

	// Verify departure exists.
	if _, err := s.store.GetDepartureByIDForStaff(ctx, params.DepartureID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("departure %q not found", params.DepartureID))
		}
		return nil, fmt.Errorf("get departure: %w", postgres_store.WrapDBError(err))
	}

	out := make([]*PricingResult, 0, len(params.Pricings))
	for _, p := range params.Pricings {
		if p.RoomType == "" {
			return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("room_type is required in each pricing entry"))
		}
		if p.ListAmountIDR < 0 {
			return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("list_amount_idr must be >= 0"))
		}

		// Mint a new ID only for inserts; ON CONFLICT keeps the existing ID.
		id, err := ulid.New("pkgpr_")
		if err != nil {
			return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("mint pricing id: %w", err))
		}

		row, err := s.store.UpsertDeparturePricing(ctx, sqlc.UpsertDeparturePricingParams{
			ID:                 id,
			PackageDepartureID: params.DepartureID,
			RoomType:           sqlc.CatalogRoomType(p.RoomType),
			ListAmount:         p.ListAmountIDR,
			ListCurrency:       "IDR",
			SettlementCurrency: "IDR",
		})
		if err != nil {
			return nil, fmt.Errorf("upsert pricing %q: %w", p.RoomType, postgres_store.WrapDBError(err))
		}
		out = append(out, pricingRowToResult(row))
	}

	span.SetStatus(codes.Ok, "success")
	return out, nil
}

func (s *Service) GetDeparturePricing(ctx context.Context, params *GetDeparturePricingParams) ([]*PricingResult, error) {
	const op = "service.Service.GetDeparturePricing"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	if params.DepartureID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("departure_id is required"))
	}

	rows, err := s.store.GetDeparturePricingRows(ctx, params.DepartureID)
	if err != nil {
		return nil, fmt.Errorf("get pricing: %w", postgres_store.WrapDBError(err))
	}

	out := make([]*PricingResult, 0, len(rows))
	for _, r := range rows {
		out = append(out, pricingRowToResult(r))
	}

	span.SetStatus(codes.Ok, "success")
	return out, nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

func clampLimit(limit int) int {
	if limit <= 0 {
		return 50
	}
	if limit > 200 {
		return 200
	}
	return limit
}

func tsToISO(ts pgtype.Timestamptz) string {
	if !ts.Valid {
		return ""
	}
	return ts.Time.UTC().Format(time.RFC3339)
}

func hotelRowToResult(r sqlc.HotelRow) *HotelResult {
	return &HotelResult{
		ID:               r.ID,
		Name:             r.Name,
		City:             r.City,
		StarRating:       int(r.StarRating),
		WalkingDistanceM: int(r.WalkingDistanceM),
		CreatedAt:        tsToISO(r.CreatedAt),
		UpdatedAt:        tsToISO(r.UpdatedAt),
	}
}

func airlineRowToResult(r sqlc.AirlineRow) *AirlineResult {
	return &AirlineResult{
		ID:           r.ID,
		Code:         strings.TrimSpace(r.Code),
		Name:         r.Name,
		OperatorKind: string(r.OperatorKind),
		CreatedAt:    tsToISO(r.CreatedAt),
		UpdatedAt:    tsToISO(r.UpdatedAt),
	}
}

func muthawwifRowToResult(r sqlc.MuthawwifRow) *MuthawwifResult {
	return &MuthawwifResult{
		ID:          r.ID,
		Name:        r.Name,
		PortraitUrl: r.PortraitUrl,
		CreatedAt:   tsToISO(r.CreatedAt),
		UpdatedAt:   tsToISO(r.UpdatedAt),
	}
}

func addonRowToResult(r sqlc.AddonRow) *AddonResult {
	return &AddonResult{
		ID:            r.ID,
		Name:          r.Name,
		ListAmountIDR: r.ListAmount,
		CreatedAt:     tsToISO(r.CreatedAt),
		UpdatedAt:     tsToISO(r.UpdatedAt),
	}
}

func pricingRowToResult(r sqlc.PricingRow) *PricingResult {
	return &PricingResult{
		ID:          r.ID,
		DepartureID: r.PackageDepartureID,
		RoomType:    string(r.RoomType),
		ListAmount:  r.ListAmount,
		CreatedAt:   tsToISO(r.CreatedAt),
		UpdatedAt:   tsToISO(r.UpdatedAt),
	}
}
