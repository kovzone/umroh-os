// S4-E-02 — hand-written sqlc equivalent for crm.leads queries.
// Regenerate with sqlc after `make migrate-up` applies migration 000014.
package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// InsertLead
// ---------------------------------------------------------------------------

const insertLead = `-- name: InsertLead :one
INSERT INTO crm.leads (
    source, utm_source, utm_medium, utm_campaign, utm_content, utm_term,
    name, phone, email, interest_package_id, interest_departure_id,
    status, assigned_cs_id, notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
)
RETURNING
    id, source, utm_source, utm_medium, utm_campaign, utm_content, utm_term,
    name, phone, email, interest_package_id, interest_departure_id,
    status, assigned_cs_id, notes, booking_id, created_at, updated_at`

type InsertLeadParams struct {
	Source              string      `json:"source"`
	UtmSource           pgtype.Text `json:"utm_source"`
	UtmMedium           pgtype.Text `json:"utm_medium"`
	UtmCampaign         pgtype.Text `json:"utm_campaign"`
	UtmContent          pgtype.Text `json:"utm_content"`
	UtmTerm             pgtype.Text `json:"utm_term"`
	Name                string      `json:"name"`
	Phone               string      `json:"phone"`
	Email               pgtype.Text `json:"email"`
	InterestPackageID   pgtype.UUID `json:"interest_package_id"`
	InterestDepartureID pgtype.UUID `json:"interest_departure_id"`
	Status              string      `json:"status"`
	AssignedCsID        pgtype.UUID `json:"assigned_cs_id"`
	Notes               pgtype.Text `json:"notes"`
}

func (q *Queries) InsertLead(ctx context.Context, arg InsertLeadParams) (Lead, error) {
	row := q.db.QueryRow(ctx, insertLead,
		arg.Source,
		arg.UtmSource,
		arg.UtmMedium,
		arg.UtmCampaign,
		arg.UtmContent,
		arg.UtmTerm,
		arg.Name,
		arg.Phone,
		arg.Email,
		arg.InterestPackageID,
		arg.InterestDepartureID,
		arg.Status,
		arg.AssignedCsID,
		arg.Notes,
	)
	return scanLead(row)
}

// ---------------------------------------------------------------------------
// GetLeadByID
// ---------------------------------------------------------------------------

const getLeadByID = `-- name: GetLeadByID :one
SELECT
    id, source, utm_source, utm_medium, utm_campaign, utm_content, utm_term,
    name, phone, email, interest_package_id, interest_departure_id,
    status, assigned_cs_id, notes, booking_id, created_at, updated_at
FROM crm.leads
WHERE id = $1`

func (q *Queries) GetLeadByID(ctx context.Context, id pgtype.UUID) (Lead, error) {
	row := q.db.QueryRow(ctx, getLeadByID, id)
	return scanLead(row)
}

// ---------------------------------------------------------------------------
// GetLeadByBookingID
// ---------------------------------------------------------------------------

const getLeadByBookingID = `-- name: GetLeadByBookingID :one
SELECT
    id, source, utm_source, utm_medium, utm_campaign, utm_content, utm_term,
    name, phone, email, interest_package_id, interest_departure_id,
    status, assigned_cs_id, notes, booking_id, created_at, updated_at
FROM crm.leads
WHERE booking_id = $1
LIMIT 1`

func (q *Queries) GetLeadByBookingID(ctx context.Context, bookingID pgtype.Text) (Lead, error) {
	row := q.db.QueryRow(ctx, getLeadByBookingID, bookingID)
	return scanLead(row)
}

// ---------------------------------------------------------------------------
// UpdateLeadStatus
// ---------------------------------------------------------------------------

const updateLeadStatus = `-- name: UpdateLeadStatus :one
UPDATE crm.leads
SET
    status         = COALESCE($2, status),
    notes          = COALESCE($3, notes),
    assigned_cs_id = COALESCE($4, assigned_cs_id),
    booking_id     = COALESCE($5, booking_id),
    updated_at     = now()
WHERE id = $1
RETURNING
    id, source, utm_source, utm_medium, utm_campaign, utm_content, utm_term,
    name, phone, email, interest_package_id, interest_departure_id,
    status, assigned_cs_id, notes, booking_id, created_at, updated_at`

type UpdateLeadStatusParams struct {
	ID           pgtype.UUID `json:"id"`
	Status       pgtype.Text `json:"status"`        // nullable → COALESCE keeps current
	Notes        pgtype.Text `json:"notes"`         // nullable → COALESCE keeps current
	AssignedCsID pgtype.UUID `json:"assigned_cs_id"` // nullable → COALESCE keeps current
	BookingID    pgtype.Text `json:"booking_id"`    // nullable → COALESCE keeps current
}

func (q *Queries) UpdateLeadStatus(ctx context.Context, arg UpdateLeadStatusParams) (Lead, error) {
	row := q.db.QueryRow(ctx, updateLeadStatus,
		arg.ID,
		arg.Status,
		arg.Notes,
		arg.AssignedCsID,
		arg.BookingID,
	)
	return scanLead(row)
}

// ---------------------------------------------------------------------------
// ListLeads
// ---------------------------------------------------------------------------

const listLeads = `-- name: ListLeads :many
SELECT
    id, source, utm_source, utm_medium, utm_campaign, utm_content, utm_term,
    name, phone, email, interest_package_id, interest_departure_id,
    status, assigned_cs_id, notes, booking_id, created_at, updated_at
FROM crm.leads
WHERE
    ($1::text IS NULL OR status = $1)
    AND ($2::uuid IS NULL OR assigned_cs_id = $2)
ORDER BY created_at DESC
LIMIT $3
OFFSET $4`

type ListLeadsParams struct {
	StatusFilter     pgtype.Text `json:"status_filter"`      // nullable
	AssignedCsFilter pgtype.UUID `json:"assigned_cs_filter"` // nullable
	Limit            int32       `json:"limit"`
	Offset           int32       `json:"offset"`
}

func (q *Queries) ListLeads(ctx context.Context, arg ListLeadsParams) ([]Lead, error) {
	rows, err := q.db.Query(ctx, listLeads,
		arg.StatusFilter,
		arg.AssignedCsFilter,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []Lead
	for rows.Next() {
		var l Lead
		if err := rows.Scan(
			&l.ID,
			&l.Source,
			&l.UtmSource,
			&l.UtmMedium,
			&l.UtmCampaign,
			&l.UtmContent,
			&l.UtmTerm,
			&l.Name,
			&l.Phone,
			&l.Email,
			&l.InterestPackageID,
			&l.InterestDepartureID,
			&l.Status,
			&l.AssignedCsID,
			&l.Notes,
			&l.BookingID,
			&l.CreatedAt,
			&l.UpdatedAt,
		); err != nil {
			return nil, err
		}
		leads = append(leads, l)
	}
	return leads, rows.Err()
}

// ---------------------------------------------------------------------------
// CountLeads
// ---------------------------------------------------------------------------

const countLeads = `-- name: CountLeads :one
SELECT COUNT(*) FROM crm.leads
WHERE
    ($1::text IS NULL OR status = $1)
    AND ($2::uuid IS NULL OR assigned_cs_id = $2)`

type CountLeadsParams struct {
	StatusFilter     pgtype.Text `json:"status_filter"`
	AssignedCsFilter pgtype.UUID `json:"assigned_cs_filter"`
}

func (q *Queries) CountLeads(ctx context.Context, arg CountLeadsParams) (int64, error) {
	row := q.db.QueryRow(ctx, countLeads, arg.StatusFilter, arg.AssignedCsFilter)
	var count int64
	err := row.Scan(&count)
	return count, err
}

// ---------------------------------------------------------------------------
// GetLeastLoadedCS
// ---------------------------------------------------------------------------

const getLeastLoadedCS = `-- name: GetLeastLoadedCS :one
SELECT assigned_cs_id
FROM crm.leads
WHERE assigned_cs_id IS NOT NULL
  AND status NOT IN ('converted', 'lost')
GROUP BY assigned_cs_id
ORDER BY COUNT(*) ASC, assigned_cs_id ASC
LIMIT 1`

// GetLeastLoadedCS returns the UUID of the CS user with the fewest active leads.
// Returns pgtype.UUID with Valid=false when no CS has any active leads.
//
// TODO(race-condition): this query uses a plain SELECT (no FOR UPDATE / advisory lock).
// Two concurrent CreateLead calls can both read the same "least loaded CS" and
// both assign to that CS, causing uneven distribution. The S4 contract recommends
// a cs_assignment_state auxiliary table keyed by cs_user_id with SELECT FOR UPDATE
// SKIP LOCKED to serialize concurrent assignments. Fix in Phase 6 BL-CRM-003 or
// a dedicated S4 hardening card.
func (q *Queries) GetLeastLoadedCS(ctx context.Context) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, getLeastLoadedCS)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

// ---------------------------------------------------------------------------
// scanLead — shared scan helper
// ---------------------------------------------------------------------------

type scannable interface {
	Scan(dest ...any) error
}

func scanLead(row scannable) (Lead, error) {
	var l Lead
	err := row.Scan(
		&l.ID,
		&l.Source,
		&l.UtmSource,
		&l.UtmMedium,
		&l.UtmCampaign,
		&l.UtmContent,
		&l.UtmTerm,
		&l.Name,
		&l.Phone,
		&l.Email,
		&l.InterestPackageID,
		&l.InterestDepartureID,
		&l.Status,
		&l.AssignedCsID,
		&l.Notes,
		&l.BookingID,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	return l, err
}
