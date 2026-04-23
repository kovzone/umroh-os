package service

// S1-E-03 / BL-BOOK-001..006 — booking-svc draft booking service methods.
//
// BL-BOOK-001: Draft booking saved with minimum fields.
// BL-BOOK-002: Booking stores channel + agent_id when present.
// BL-BOOK-003: Valid status transitions — S1 only ever writes 'draft'.
// BL-BOOK-004: Hard fail if seats insufficient or departure invalid.
// BL-BOOK-005: Document gate validation — DRAFT assumption: KTP + passport
//              required per Q006. Missing docs → clear error per pilgrim/doc kind.
//              Marked as DRAFT assumption (not blocking in MVP, just validated).
// BL-BOOK-006: Mahram validation result stored, does NOT block submit.
//
// Architecture: service layer trusts the REST adapter to have validated the
// bearer token (via iam-svc.ValidateToken); the service calls catalog-svc
// (via gRPC) to validate the departure and reserve seats atomically.
//
// Idempotency: if Idempotency-Key header was provided and matches an existing
// draft booking with the same request body hash, the existing booking is returned.

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"booking-svc/adapter/catalog_grpc_adapter"
	"booking-svc/store/postgres_store"
	"booking-svc/store/postgres_store/sqlc"
	"booking-svc/util/apperrors"
	"booking-svc/util/logging"
	"booking-svc/util/ulid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Domain types for booking draft
// ---------------------------------------------------------------------------

// PilgrimInput is one jamaah (pilgrim) entry in the booking request.
type PilgrimInput struct {
	FullName string
	Email    string
	Whatsapp string
	Domicile string
	IsLead   bool
	// Document validation fields (BL-BOOK-005: DRAFT assumption)
	// KTP and passport are required per Q006 engineering assumption.
	// We store whether docs are present; validation warnings (not errors) for MVP.
	HasKTP      bool
	HasPassport bool
}

// AddonInput is a selected add-on for the booking.
type AddonInput struct {
	AddonID    string
	AddonName  string
	ListAmount int64
	ListCurrency string
	SettlementCurrency string
}

// MahramInfo is the mahram relationship info (BL-BOOK-006: stored, not blocking).
type MahramInfo struct {
	// MahramWarning is set when mahram validation found an issue (non-blocking).
	MahramWarning string
}

// CreateDraftBookingParams is the input for CreateDraftBooking.
type CreateDraftBookingParams struct {
	// Channel attribution (BL-BOOK-002)
	Channel string // "b2c_self" | "b2b_agent" | "cs"
	AgentID string // populated when Channel == "b2b_agent"
	StaffUserID string // populated when Channel == "cs"

	// Catalog references
	PackageID   string
	DepartureID string
	RoomType    string

	// Lead pilgrim contact
	LeadFullName string
	LeadEmail    string
	LeadWhatsapp string
	LeadDomicile string

	// All pilgrims (including the lead)
	Pilgrims []PilgrimInput

	// Optional add-ons
	Addons []AddonInput

	// Pricing snapshot
	ListAmount         int64
	ListCurrency       string
	SettlementCurrency string

	Notes string

	// Idempotency (optional)
	IdempotencyKey string

	// Mahram validation result (BL-BOOK-006) — populated by caller if applicable.
	Mahram *MahramInfo
}

// BookingItemResult is one pilgrim in the result.
type BookingItemResult struct {
	ID       string
	FullName string
	IsLead   bool
	// Document validation warning (BL-BOOK-005)
	// Empty = no issues; non-empty = advisory message (not an error).
	DocumentWarning string
}

// DraftBookingResult is the response from CreateDraftBooking.
type DraftBookingResult struct {
	ID                 string
	Status             string
	Channel            string
	PackageID          string
	DepartureID        string
	RoomType           string
	AgentID            string
	StaffUserID        string
	LeadFullName       string
	LeadEmail          string
	LeadWhatsapp       string
	LeadDomicile       string
	ListAmount         int64
	ListCurrency       string
	SettlementCurrency string
	Notes              string
	IdempotencyKey     string
	CreatedAt          string
	ExpiresAt          string
	Items              []BookingItemResult
	// Mahram warning (BL-BOOK-006) — non-empty = advisory only
	MahramWarning string
	// Replayed is true when the response is from an idempotency dedup hit.
	Replayed bool
}

// ---------------------------------------------------------------------------
// CreateDraftBooking
// ---------------------------------------------------------------------------

func (s *Service) CreateDraftBooking(ctx context.Context, params *CreateDraftBookingParams) (*DraftBookingResult, error) {
	const op = "service.Service.CreateDraftBooking"

	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.channel", params.Channel),
		attribute.String("input.departure_id", params.DepartureID),
	)
	logger.Info().Str("op", op).Str("channel", params.Channel).Str("departure_id", params.DepartureID).Msg("")

	// --- Input validation ---
	if err := validateCreateBookingParams(params); err != nil {
		return nil, err
	}

	// --- Idempotency dedup (BL-BOOK-001: if key provided, check for existing) ---
	if params.IdempotencyKey != "" {
		existing, err := s.store.GetBookingByIdempotencyKey(ctx, sqlc.GetBookingByIdempotencyKeyParams{
			Channel:        sqlc.BookingChannel(params.Channel),
			IdempotencyKey: params.IdempotencyKey,
		})
		if err == nil {
			// Found existing — check body hash match.
			bodyHash := computeBodyHash(params)
			if existing.IdempotencyBodyHash.Valid && existing.IdempotencyBodyHash.String == bodyHash {
				logger.Info().Str("op", op).Str("booking_id", existing.ID).Msg("idempotency replay")
				span.SetStatus(codes.Ok, "replayed")
				result, hydErr := s.hydrateBookingResult(ctx, existing, true)
				return result, hydErr
			}
			// Body hash mismatch — different request with same key.
			return nil, errors.Join(apperrors.ErrConflict, fmt.Errorf(
				"idempotency_key_conflict: key %q was previously used with a different request body", params.IdempotencyKey))
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("check idempotency key: %w", postgres_store.WrapDBError(err))
		}
	}

	// --- BL-BOOK-004: Validate departure via catalog-svc gRPC ---
	dep, err := s.catalogClient.GetDeparture(ctx, params.DepartureID)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("departure %q not found or not bookable", params.DepartureID))
		}
		return nil, fmt.Errorf("validate departure: %w", err)
	}
	if dep.Status != "open" && dep.Status != "closed" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf(
			"departure %q is not accepting bookings (status: %s)", params.DepartureID, dep.Status))
	}
	// Verify departure belongs to the specified package.
	if dep.PackageID != params.PackageID {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf(
			"departure %q does not belong to package %q", params.DepartureID, params.PackageID))
	}

	numPilgrims := len(params.Pilgrims)
	if numPilgrims < 1 {
		numPilgrims = 1 // at minimum the lead
	}

	// --- BL-BOOK-004: Hard fail if seats insufficient ---
	if dep.RemainingSeats < numPilgrims {
		return nil, errors.Join(apperrors.ErrConflict, fmt.Errorf(
			"insufficient_capacity: %d seats requested but only %d remaining on departure %q",
			numPilgrims, dep.RemainingSeats, params.DepartureID))
	}

	// --- BL-BOOK-005: Document gate (DRAFT engineering assumption per Q006) ---
	// KTP + passport required. Missing docs → warning per pilgrim, not hard fail in MVP.
	// This is marked DRAFT — the actual doc validation engine lands in a later slice.
	var docWarnings []string
	for i, p := range params.Pilgrims {
		name := p.FullName
		if name == "" {
			name = fmt.Sprintf("pilgrim %d", i+1)
		}
		if !p.HasKTP {
			docWarnings = append(docWarnings, fmt.Sprintf("%s: KTP diperlukan", name))
		}
		if !p.HasPassport {
			docWarnings = append(docWarnings, fmt.Sprintf("%s: paspor diperlukan", name))
		}
	}
	// For MVP, doc warnings are advisory (not blocking). Log them but proceed.
	if len(docWarnings) > 0 {
		logger.Warn().Strs("doc_warnings", docWarnings).Str("op", op).Msg("DRAFT: document gate warnings (not blocking in S1)")
	}

	// --- BL-BOOK-006: Mahram validation (store result, do not block) ---
	mahramWarning := ""
	if params.Mahram != nil {
		mahramWarning = params.Mahram.MahramWarning
		if mahramWarning != "" {
			logger.Warn().Str("mahram_warning", mahramWarning).Msg("mahram validation warning (non-blocking)")
		}
	}

	// --- Reserve seats via catalog-svc gRPC (BL-BOOK-004) ---
	reservationID, err := ulid.New("rsv_")
	if err != nil {
		return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("mint reservation id: %w", err))
	}
	_, err = s.catalogClient.ReserveSeats(ctx, &catalog_grpc_adapter.ReserveSeatsParams{
		ReservationID:       reservationID,
		DepartureID:         params.DepartureID,
		Seats:               numPilgrims,
		IdempotencyTTLHours: 1, // draft bookings hold seats for 1 hour
	})
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, apperrors.ErrConflict) {
			return nil, errors.Join(apperrors.ErrConflict, fmt.Errorf("insufficient_capacity: seat reservation failed: %w", err))
		}
		return nil, fmt.Errorf("reserve seats: %w", err)
	}

	// --- Generate booking ID ---
	bookingID, err := ulid.New("bkg_")
	if err != nil {
		// Attempt to release seats before returning error.
		_, _ = s.catalogClient.ReleaseSeats(ctx, &catalog_grpc_adapter.ReleaseSeatsParams{
			ReservationID: reservationID,
			Reason:        "booking id mint failed",
		})
		return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("mint booking id: %w", err))
	}

	// --- Persist booking row ---
	currency := params.ListCurrency
	if currency == "" {
		currency = "IDR"
	}
	settlement := params.SettlementCurrency
	if settlement == "" {
		settlement = "IDR"
	}
	expiresAt := time.Now().UTC().Add(30 * time.Minute)
	bodyHash := computeBodyHash(params)

	insertArg := sqlc.InsertBookingParams{
		ID:                 bookingID,
		Channel:            sqlc.BookingChannel(params.Channel),
		PackageID:          params.PackageID,
		DepartureID:        params.DepartureID,
		RoomType:           params.RoomType,
		AgentID:            optText(params.AgentID),
		StaffUserID:        optText(params.StaffUserID),
		LeadFullName:       params.LeadFullName,
		LeadEmail:          optText(params.LeadEmail),
		LeadWhatsapp:       params.LeadWhatsapp,
		LeadDomicile:       params.LeadDomicile,
		ListAmount:         params.ListAmount,
		ListCurrency:       currency,
		SettlementCurrency: settlement,
		Notes:              optText(params.Notes),
		IdempotencyKey:     optText(params.IdempotencyKey),
		IdempotencyBodyHash: optText(bodyHash),
		ExpiresAt:          expiresAt,
	}

	booking, err := s.store.InsertBooking(ctx, insertArg)
	if err != nil {
		// Attempt to release seats before returning error.
		_, _ = s.catalogClient.ReleaseSeats(ctx, &catalog_grpc_adapter.ReleaseSeatsParams{
			ReservationID: reservationID,
			Reason:        "booking insert failed",
		})
		return nil, fmt.Errorf("insert booking: %w", postgres_store.WrapDBError(err))
	}

	// --- Persist booking items (one per pilgrim) ---
	var items []BookingItemResult
	for i, p := range params.Pilgrims {
		itemID, err := ulid.New("bkgitem_")
		if err != nil {
			return nil, errors.Join(apperrors.ErrInternal, fmt.Errorf("mint item id: %w", err))
		}
		item, err := s.store.InsertBookingItem(ctx, sqlc.InsertBookingItemParams{
			ID:        itemID,
			BookingID: bookingID,
			FullName:  p.FullName,
			Email:     optText(p.Email),
			Whatsapp:  optText(p.Whatsapp),
			Domicile:  p.Domicile,
			IsLead:    p.IsLead || i == 0,
		})
		if err != nil {
			return nil, fmt.Errorf("insert booking item %d: %w", i, postgres_store.WrapDBError(err))
		}

		// BL-BOOK-005: per-pilgrim document warning (advisory)
		docWarn := ""
		if !p.HasKTP {
			docWarn += "KTP diperlukan; "
		}
		if !p.HasPassport {
			docWarn += "paspor diperlukan"
		}
		items = append(items, BookingItemResult{
			ID:              item.ID,
			FullName:        item.FullName,
			IsLead:          item.IsLead,
			DocumentWarning: docWarn,
		})
	}

	// --- Persist add-ons (if any) ---
	for _, a := range params.Addons {
		addonCurrency := a.ListCurrency
		if addonCurrency == "" {
			addonCurrency = "IDR"
		}
		addonSettlement := a.SettlementCurrency
		if addonSettlement == "" {
			addonSettlement = "IDR"
		}
		if err := s.store.InsertBookingAddon(ctx, sqlc.InsertBookingAddonParams{
			BookingID:          bookingID,
			AddonID:            a.AddonID,
			AddonName:          a.AddonName,
			ListAmount:         a.ListAmount,
			ListCurrency:       addonCurrency,
			SettlementCurrency: addonSettlement,
		}); err != nil {
			return nil, fmt.Errorf("insert booking addon %s: %w", a.AddonID, postgres_store.WrapDBError(err))
		}
	}

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(attribute.String("output.booking_id", bookingID))

	result := &DraftBookingResult{
		ID:                 booking.ID,
		Status:             string(booking.Status),
		Channel:            string(booking.Channel),
		PackageID:          booking.PackageID,
		DepartureID:        booking.DepartureID,
		RoomType:           booking.RoomType,
		AgentID:            booking.AgentID.String,
		StaffUserID:        booking.StaffUserID.String,
		LeadFullName:       booking.LeadFullName,
		LeadEmail:          booking.LeadEmail.String,
		LeadWhatsapp:       booking.LeadWhatsapp,
		LeadDomicile:       booking.LeadDomicile,
		ListAmount:         booking.ListAmount,
		ListCurrency:       booking.ListCurrency,
		SettlementCurrency: booking.SettlementCurrency,
		Notes:              booking.Notes.String,
		IdempotencyKey:     booking.IdempotencyKey.String,
		CreatedAt:          tsToISO(booking.CreatedAt),
		ExpiresAt:          tsToISO(booking.ExpiresAt),
		Items:              items,
		MahramWarning:      mahramWarning,
		Replayed:           false,
	}

	// --- S4-E-02: CRM fan-out (best-effort, non-blocking) ---
	// Notify crm-svc that a booking was created. If a lead_id is embedded in
	// the notes or passed from the caller context it can be passed here.
	// For now we pass empty lead_id — attribution is established when the
	// gateway embeds lead context in a later slice.
	s.FanOutBookingCreated(ctx, &FanOutBookingCreatedParams{
		BookingID:   bookingID,
		LeadID:      "", // populated in a later slice (F4/F10 attribution)
		PackageID:   params.PackageID,
		DepartureID: params.DepartureID,
		JamaahCount: int32(numPilgrims),
		CreatedAt:   tsToISO(booking.CreatedAt),
	})

	return result, nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

func validateCreateBookingParams(p *CreateDraftBookingParams) error {
	switch p.Channel {
	case "b2c_self", "b2b_agent", "cs":
	default:
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("channel must be b2c_self, b2b_agent, or cs; got %q", p.Channel))
	}
	if p.Channel == "b2b_agent" && p.AgentID == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("agent_id required when channel = b2b_agent"))
	}
	if p.Channel == "cs" && p.StaffUserID == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("staff_user_id required when channel = cs"))
	}
	if p.PackageID == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("package_id is required"))
	}
	if p.DepartureID == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("departure_id is required"))
	}
	if p.RoomType == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("room_type is required"))
	}
	if p.LeadFullName == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("lead_full_name is required"))
	}
	if p.LeadWhatsapp == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("lead_whatsapp is required"))
	}
	if p.LeadDomicile == "" {
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("lead_domicile is required"))
	}
	return nil
}

// computeBodyHash computes a deterministic SHA-256 hash of the booking request
// body for idempotency body-match checking.
func computeBodyHash(p *CreateDraftBookingParams) string {
	// Marshal the canonical fields that define request identity.
	type canonical struct {
		Channel     string
		PackageID   string
		DepartureID string
		RoomType    string
		LeadName    string
		LeadPhone   string
	}
	b, _ := json.Marshal(canonical{
		Channel:     p.Channel,
		PackageID:   p.PackageID,
		DepartureID: p.DepartureID,
		RoomType:    p.RoomType,
		LeadName:    p.LeadFullName,
		LeadPhone:   p.LeadWhatsapp,
	})
	h := sha256.Sum256(b)
	return fmt.Sprintf("%x", h)
}

// optText converts a string to pgtype.Text (Valid=true when non-empty).
func optText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

// tsToISO formats a pgtype.Timestamptz to RFC-3339 or "".
func tsToISO(ts pgtype.Timestamptz) string {
	if !ts.Valid {
		return ""
	}
	return ts.Time.UTC().Format(time.RFC3339)
}

// hydrateBookingResult converts a BookingRow to DraftBookingResult (used for replays).
func (s *Service) hydrateBookingResult(ctx context.Context, row sqlc.BookingRow, replayed bool) (*DraftBookingResult, error) {
	items, err := s.store.ListBookingItems(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("list booking items: %w", postgres_store.WrapDBError(err))
	}
	var resultItems []BookingItemResult
	for _, it := range items {
		resultItems = append(resultItems, BookingItemResult{
			ID:       it.ID,
			FullName: it.FullName,
			IsLead:   it.IsLead,
		})
	}
	return &DraftBookingResult{
		ID:                 row.ID,
		Status:             string(row.Status),
		Channel:            string(row.Channel),
		PackageID:          row.PackageID,
		DepartureID:        row.DepartureID,
		RoomType:           row.RoomType,
		AgentID:            row.AgentID.String,
		StaffUserID:        row.StaffUserID.String,
		LeadFullName:       row.LeadFullName,
		LeadEmail:          row.LeadEmail.String,
		LeadWhatsapp:       row.LeadWhatsapp,
		LeadDomicile:       row.LeadDomicile,
		ListAmount:         row.ListAmount,
		ListCurrency:       row.ListCurrency,
		SettlementCurrency: row.SettlementCurrency,
		Notes:              row.Notes.String,
		IdempotencyKey:     row.IdempotencyKey.String,
		CreatedAt:          tsToISO(row.CreatedAt),
		ExpiresAt:          tsToISO(row.ExpiresAt),
		Items:              resultItems,
		Replayed:           replayed,
	}, nil
}
