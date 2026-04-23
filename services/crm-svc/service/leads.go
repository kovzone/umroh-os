// leads.go — CRM lead business logic (S4-E-02 / BL-CRM-001..003).
//
// Implements CreateLead (round-robin CS assignment), GetLead,
// UpdateLead (partial update + status transition validation),
// ListLeads (filtered pagination), OnBookingCreated (qualified),
// and OnBookingPaidInFull (converted) on the service layer.
//
// Per ADR-0006: all calls are synchronous. Per ADR-0007: schema lives in
// migration 000014. Per Elda conventions: all mutations have audit-log entries.

package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"crm-svc/store/postgres_store/sqlc"
	"crm-svc/util/logging"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/attribute"
)

// ---------------------------------------------------------------------------
// Params / results
// ---------------------------------------------------------------------------

// CreateLeadParams is the input for CreateLead.
type CreateLeadParams struct {
	Name                string
	Phone               string
	Email               string
	Source              string
	UtmSource           string
	UtmMedium           string
	UtmCampaign         string
	UtmContent          string
	UtmTerm             string
	InterestPackageID   string
	InterestDepartureID string
	Notes               string
}

// LeadResult is returned by all lead read/write methods.
type LeadResult struct {
	ID                  string
	Source              string
	UtmSource           string
	UtmMedium           string
	UtmCampaign         string
	UtmContent          string
	UtmTerm             string
	Name                string
	Phone               string
	Email               string
	InterestPackageID   string
	InterestDepartureID string
	Status              string
	AssignedCsID        string
	Notes               string
	BookingID           string
	CreatedAt           string
	UpdatedAt           string
}

// GetLeadParams is the input for GetLead.
type GetLeadParams struct {
	ID string
}

// UpdateLeadParams is the input for UpdateLead — all optional except ID.
type UpdateLeadParams struct {
	ID           string
	Status       string // empty = keep current
	Notes        string // empty = keep current
	AssignedCsID string // empty = keep current
}

// ListLeadsParams is the input for ListLeads.
type ListLeadsParams struct {
	StatusFilter     string
	AssignedCsFilter string
	Page             int32
	PageSize         int32
}

// ListLeadsResult holds a paginated list and the total count.
type ListLeadsResult struct {
	Leads    []*LeadResult
	Total    int64
	Page     int32
	PageSize int32
}

// OnBookingCreatedParams is the input for OnBookingCreated.
type OnBookingCreatedParams struct {
	BookingID   string
	LeadID      string
	PackageID   string
	DepartureID string
	JamaahCount int32
	CreatedAt   string
}

// OnBookingCreatedResult is the output of OnBookingCreated.
type OnBookingCreatedResult struct {
	Updated bool
	LeadID  string
}

// OnBookingPaidInFullParams is the input for OnBookingPaidInFull.
type OnBookingPaidInFullParams struct {
	BookingID string
	LeadID    string
	PaidAt    string
}

// OnBookingPaidInFullResult is the output of OnBookingPaidInFull.
type OnBookingPaidInFullResult struct {
	Updated bool
	LeadID  string
}

// ---------------------------------------------------------------------------
// Allowed status transitions
// ---------------------------------------------------------------------------

// validTransition returns true if moving from `from` to `to` is allowed.
// Per S4 contract: 'converted' and 'lost' are both terminal — no further transitions.
func validTransition(from, to string) bool {
	if from == to {
		return true // idempotent
	}
	// Both terminal states reject any outbound transition.
	if from == "converted" || from == "lost" {
		return false // terminal
	}
	// Only valid target statuses are accepted.
	valid := map[string]bool{
		"new":       true,
		"contacted": true,
		"qualified": true,
		"converted": true,
		"lost":      true,
	}
	return valid[to]
}

// ---------------------------------------------------------------------------
// CreateLead
// ---------------------------------------------------------------------------

func (svc *Service) CreateLead(ctx context.Context, params *CreateLeadParams) (*LeadResult, error) {
	const op = "service.Service.CreateLead"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, svc.logger)

	// Sanitize phone: strip leading/trailing whitespace before validation/storage.
	params.Phone = strings.TrimSpace(params.Phone)

	// Validate phone (minimal: non-empty, ≥8 chars)
	if len(params.Phone) < 8 {
		err := fmt.Errorf("%s: phone must be at least 8 characters", op)
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	// Default source to 'direct' if not set or invalid
	source := params.Source
	allowed := map[string]bool{
		"organic": true, "whatsapp": true, "instagram": true,
		"facebook": true, "referral": true, "agent": true, "direct": true,
	}
	if !allowed[source] {
		source = "direct"
	}

	// Round-robin CS assignment
	assignedCSID, err := svc.store.GetLeastLoadedCS(ctx)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: get least loaded cs: %w", op, err)
	}
	// If no rows, assignedCSID is a zero pgtype.UUID (Valid=false) → NULL in DB.

	lead, err := svc.store.InsertLead(ctx, sqlc.InsertLeadParams{
		Source:              source,
		UtmSource:           textOrNull(params.UtmSource),
		UtmMedium:           textOrNull(params.UtmMedium),
		UtmCampaign:         textOrNull(params.UtmCampaign),
		UtmContent:          textOrNull(params.UtmContent),
		UtmTerm:             textOrNull(params.UtmTerm),
		Name:                params.Name,
		Phone:               params.Phone,
		Email:               textOrNull(params.Email),
		InterestPackageID:   uuidOrNull(params.InterestPackageID),
		InterestDepartureID: uuidOrNull(params.InterestDepartureID),
		Status:              "new",
		AssignedCsID:        assignedCSID,
		Notes:               textOrNull(params.Notes),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		logger.Error().Err(err).Str("op", op).Msg("InsertLead failed")
		return nil, fmt.Errorf("%s: insert lead: %w", op, err)
	}

	logger.Info().
		Str("op", op).
		Str("lead_id", uuidStr(lead.ID)).
		Str("phone", params.Phone).
		Str("assigned_cs_id", uuidStr(lead.AssignedCsID)).
		Msg("lead created")

	span.SetAttributes(attribute.String("lead_id", uuidStr(lead.ID)))
	span.SetStatus(otelCodes.Ok, "created")
	return mapLead(lead), nil
}

// ---------------------------------------------------------------------------
// GetLead
// ---------------------------------------------------------------------------

func (svc *Service) GetLead(ctx context.Context, params *GetLeadParams) (*LeadResult, error) {
	const op = "service.Service.GetLead"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id, err := parseUUID(params.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid id: %w", op, err)
	}

	lead, err := svc.store.GetLeadByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: lead not found", op)
		}
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: get lead: %w", op, err)
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return mapLead(lead), nil
}

// ---------------------------------------------------------------------------
// UpdateLead
// ---------------------------------------------------------------------------

func (svc *Service) UpdateLead(ctx context.Context, params *UpdateLeadParams) (*LeadResult, error) {
	const op = "service.Service.UpdateLead"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, svc.logger)

	id, err := parseUUID(params.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid id: %w", op, err)
	}

	// If a new status is requested, validate the transition.
	if params.Status != "" {
		// Fetch current lead to check the existing status.
		current, err := svc.store.GetLeadByID(ctx, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("%s: lead not found", op)
			}
			return nil, fmt.Errorf("%s: fetch current lead: %w", op, err)
		}
		if !validTransition(current.Status, params.Status) {
			return nil, fmt.Errorf("%s: invalid transition %s → %s", op, current.Status, params.Status)
		}
	}

	lead, err := svc.store.UpdateLeadStatus(ctx, sqlc.UpdateLeadStatusParams{
		ID:           id,
		Status:       textOrNull(params.Status),
		Notes:        textOrNull(params.Notes),
		AssignedCsID: uuidOrNull(params.AssignedCsID),
		BookingID:    pgtype.Text{}, // not changed via UpdateLead
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: lead not found", op)
		}
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: update lead: %w", op, err)
	}

	logger.Info().
		Str("op", op).
		Str("lead_id", params.ID).
		Str("new_status", lead.Status).
		Msg("lead updated")

	span.SetStatus(otelCodes.Ok, "updated")
	return mapLead(lead), nil
}

// ---------------------------------------------------------------------------
// ListLeads
// ---------------------------------------------------------------------------

func (svc *Service) ListLeads(ctx context.Context, params *ListLeadsParams) (*ListLeadsResult, error) {
	const op = "service.Service.ListLeads"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	// Normalise pagination defaults.
	page := params.Page
	if page < 1 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	countP := sqlc.CountLeadsParams{
		StatusFilter:     textOrNull(params.StatusFilter),
		AssignedCsFilter: uuidOrNull(params.AssignedCsFilter),
	}
	total, err := svc.store.CountLeads(ctx, countP)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: count leads: %w", op, err)
	}

	listP := sqlc.ListLeadsParams{
		StatusFilter:     textOrNull(params.StatusFilter),
		AssignedCsFilter: uuidOrNull(params.AssignedCsFilter),
		Limit:            pageSize,
		Offset:           offset,
	}
	rows, err := svc.store.ListLeads(ctx, listP)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: list leads: %w", op, err)
	}

	out := make([]*LeadResult, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapLead(r))
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &ListLeadsResult{
		Leads:    out,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// ---------------------------------------------------------------------------
// OnBookingCreated
// ---------------------------------------------------------------------------

func (svc *Service) OnBookingCreated(ctx context.Context, params *OnBookingCreatedParams) (*OnBookingCreatedResult, error) {
	const op = "service.Service.OnBookingCreated"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, svc.logger)

	// No lead_id → nothing to update
	if params.LeadID == "" {
		span.SetStatus(otelCodes.Ok, "no lead_id, skipped")
		return &OnBookingCreatedResult{Updated: false}, nil
	}

	id, err := parseUUID(params.LeadID)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid lead_id: %w", op, err)
	}

	// Idempotency: if this booking_id is already stored on a lead, no-op.
	existing, err := svc.store.GetLeadByBookingID(ctx, pgtype.Text{String: params.BookingID, Valid: true})
	if err == nil && uuidStr(existing.ID) != "" {
		logger.Info().
			Str("op", op).
			Str("booking_id", params.BookingID).
			Msg("OnBookingCreated: duplicate booking_id, no-op")
		span.SetStatus(otelCodes.Ok, "duplicate, no-op")
		return &OnBookingCreatedResult{Updated: false, LeadID: uuidStr(existing.ID)}, nil
	}

	// Validate transition: only update if current status allows moving to 'qualified'.
	current, err := svc.store.GetLeadByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn().Str("op", op).Str("lead_id", params.LeadID).Msg("lead not found, skipping")
			span.SetStatus(otelCodes.Ok, "lead not found, skipped")
			return &OnBookingCreatedResult{Updated: false}, nil
		}
		return nil, fmt.Errorf("%s: fetch lead: %w", op, err)
	}
	if !validTransition(current.Status, "qualified") {
		logger.Info().
			Str("op", op).
			Str("lead_id", params.LeadID).
			Str("current_status", current.Status).
			Msg("OnBookingCreated: status transition not allowed, skipping")
		span.SetStatus(otelCodes.Ok, "transition not allowed, skipped")
		return &OnBookingCreatedResult{Updated: false, LeadID: params.LeadID}, nil
	}

	_, err = svc.store.UpdateLeadStatus(ctx, sqlc.UpdateLeadStatusParams{
		ID:           id,
		Status:       pgtype.Text{String: "qualified", Valid: true},
		Notes:        pgtype.Text{},
		AssignedCsID: pgtype.UUID{},
		BookingID:    pgtype.Text{String: params.BookingID, Valid: true},
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: update lead: %w", op, err)
	}

	logger.Info().
		Str("op", op).
		Str("lead_id", params.LeadID).
		Str("booking_id", params.BookingID).
		Msg("lead status updated to qualified")

	span.SetStatus(otelCodes.Ok, "updated")
	return &OnBookingCreatedResult{Updated: true, LeadID: params.LeadID}, nil
}

// ---------------------------------------------------------------------------
// OnBookingPaidInFull
// ---------------------------------------------------------------------------

func (svc *Service) OnBookingPaidInFull(ctx context.Context, params *OnBookingPaidInFullParams) (*OnBookingPaidInFullResult, error) {
	const op = "service.Service.OnBookingPaidInFull"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, svc.logger)

	// No lead_id → nothing to update
	if params.LeadID == "" {
		span.SetStatus(otelCodes.Ok, "no lead_id, skipped")
		return &OnBookingPaidInFullResult{Updated: false}, nil
	}

	id, err := parseUUID(params.LeadID)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid lead_id: %w", op, err)
	}

	// Idempotency: if lead is already 'converted', no-op.
	current, err := svc.store.GetLeadByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Warn().Str("op", op).Str("lead_id", params.LeadID).Msg("lead not found, skipping")
			span.SetStatus(otelCodes.Ok, "lead not found, skipped")
			return &OnBookingPaidInFullResult{Updated: false}, nil
		}
		return nil, fmt.Errorf("%s: fetch lead: %w", op, err)
	}
	if current.Status == "converted" {
		logger.Info().
			Str("op", op).
			Str("lead_id", params.LeadID).
			Msg("OnBookingPaidInFull: already converted, no-op")
		span.SetStatus(otelCodes.Ok, "already converted, no-op")
		return &OnBookingPaidInFullResult{Updated: false, LeadID: params.LeadID}, nil
	}
	// Per S4 contract: if lead is already 'lost', do NOT override — log warning and return success.
	if current.Status == "lost" {
		logger.Warn().
			Str("op", op).
			Str("lead_id", params.LeadID).
			Msg("OnBookingPaidInFull: lead is 'lost' — booking paid but lead status not overridden")
		span.SetStatus(otelCodes.Ok, "lead lost, not overridden")
		return &OnBookingPaidInFullResult{Updated: false, LeadID: params.LeadID}, nil
	}

	_, err = svc.store.UpdateLeadStatus(ctx, sqlc.UpdateLeadStatusParams{
		ID:           id,
		Status:       pgtype.Text{String: "converted", Valid: true},
		Notes:        pgtype.Text{},
		AssignedCsID: pgtype.UUID{},
		BookingID:    pgtype.Text{String: params.BookingID, Valid: params.BookingID != ""},
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: update lead: %w", op, err)
	}

	logger.Info().
		Str("op", op).
		Str("lead_id", params.LeadID).
		Str("booking_id", params.BookingID).
		Msg("lead status updated to converted")

	span.SetStatus(otelCodes.Ok, "updated")
	return &OnBookingPaidInFullResult{Updated: true, LeadID: params.LeadID}, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func textOrNull(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

func uuidOrNull(s string) pgtype.UUID {
	if s == "" {
		return pgtype.UUID{}
	}
	var u pgtype.UUID
	if err := u.Scan(s); err != nil {
		return pgtype.UUID{}
	}
	return u
}

func parseUUID(s string) (pgtype.UUID, error) {
	var u pgtype.UUID
	if err := u.Scan(s); err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid UUID %q: %w", s, err)
	}
	return u, nil
}

func uuidStr(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	// Format as standard UUID string.
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		u.Bytes[0:4], u.Bytes[4:6], u.Bytes[6:8], u.Bytes[8:10], u.Bytes[10:16])
}

func textStr(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func mapLead(l sqlc.Lead) *LeadResult {
	return &LeadResult{
		ID:                  uuidStr(l.ID),
		Source:              l.Source,
		UtmSource:           textStr(l.UtmSource),
		UtmMedium:           textStr(l.UtmMedium),
		UtmCampaign:         textStr(l.UtmCampaign),
		UtmContent:          textStr(l.UtmContent),
		UtmTerm:             textStr(l.UtmTerm),
		Name:                l.Name,
		Phone:               l.Phone,
		Email:               textStr(l.Email),
		InterestPackageID:   uuidStr(l.InterestPackageID),
		InterestDepartureID: uuidStr(l.InterestDepartureID),
		Status:              l.Status,
		AssignedCsID:        uuidStr(l.AssignedCsID),
		Notes:               textStr(l.Notes),
		BookingID:           textStr(l.BookingID),
		CreatedAt:           l.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:           l.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}
}
