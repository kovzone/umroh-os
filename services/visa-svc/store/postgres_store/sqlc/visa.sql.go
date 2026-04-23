// visa.sql.go — hand-written SQLC-style queries for visa pipeline (BL-VISA-001..003).

package sqlc

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// Models
// ---------------------------------------------------------------------------

type VisaApplication struct {
	ID          string             `json:"id"`
	JamaahID    string             `json:"jamaah_id"`
	BookingID   string             `json:"booking_id"`
	DepartureID string             `json:"departure_id"`
	Status      string             `json:"status"`
	ProviderID  pgtype.Text        `json:"provider_id"`
	ProviderRef pgtype.Text        `json:"provider_ref"`
	EVizaURL    pgtype.Text        `json:"e_visa_url"`
	IssuedDate  pgtype.Date        `json:"issued_date"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

type StatusHistoryEntry struct {
	ID            int64              `json:"id"`
	ApplicationID string             `json:"application_id"`
	FromStatus    pgtype.Text        `json:"from_status"`
	ToStatus      string             `json:"to_status"`
	Reason        pgtype.Text        `json:"reason"`
	ActorUserID   pgtype.Text        `json:"actor_user_id"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
}

// ---------------------------------------------------------------------------
// GetVisaApplication
// ---------------------------------------------------------------------------

const getVisaApplication = `
SELECT id, jamaah_id, booking_id, departure_id, status, provider_id, provider_ref, e_visa_url, issued_date, created_at, updated_at
FROM visa.visa_applications
WHERE id = $1
`

func (q *Queries) GetVisaApplication(ctx context.Context, id string) (VisaApplication, error) {
	row := q.db.QueryRow(ctx, getVisaApplication, id)
	var i VisaApplication
	err := row.Scan(
		&i.ID, &i.JamaahID, &i.BookingID, &i.DepartureID,
		&i.Status, &i.ProviderID, &i.ProviderRef, &i.EVizaURL,
		&i.IssuedDate, &i.CreatedAt, &i.UpdatedAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// UpdateVisaStatus
// ---------------------------------------------------------------------------

const updateVisaStatus = `
UPDATE visa.visa_applications
SET status = $1, updated_at = NOW()
WHERE id = $2
`

func (q *Queries) UpdateVisaStatus(ctx context.Context, newStatus, id string) error {
	_, err := q.db.Exec(ctx, updateVisaStatus, newStatus, id)
	return err
}

// ---------------------------------------------------------------------------
// UpdateVisaStatusAndProvider
// ---------------------------------------------------------------------------

const updateVisaStatusAndProvider = `
UPDATE visa.visa_applications
SET status = $1, provider_id = $2, updated_at = NOW()
WHERE id = $3
`

func (q *Queries) UpdateVisaStatusAndProvider(ctx context.Context, newStatus, providerID, id string) error {
	_, err := q.db.Exec(ctx, updateVisaStatusAndProvider, newStatus, providerID, id)
	return err
}

// ---------------------------------------------------------------------------
// InsertStatusHistory
// ---------------------------------------------------------------------------

type InsertStatusHistoryParams struct {
	ApplicationID string
	FromStatus    string
	ToStatus      string
	Reason        string
	ActorUserID   string
}

const insertStatusHistory = `
INSERT INTO visa.status_history (application_id, from_status, to_status, reason, actor_user_id)
VALUES ($1, $2, $3, $4, $5)
`

func (q *Queries) InsertStatusHistory(ctx context.Context, arg InsertStatusHistoryParams) error {
	_, err := q.db.Exec(ctx, insertStatusHistory,
		arg.ApplicationID,
		nullText(arg.FromStatus),
		arg.ToStatus,
		nullText(arg.Reason),
		nullText(arg.ActorUserID),
	)
	return err
}

// ---------------------------------------------------------------------------
// GetVisaApplicationsForDeparture
// ---------------------------------------------------------------------------

type GetVisaApplicationsForDepartureParams struct {
	DepartureID  string
	StatusFilter string // empty = no filter
}

const getVisaApplicationsAll = `
SELECT id, jamaah_id, booking_id, departure_id, status, provider_id, provider_ref, e_visa_url, issued_date, created_at, updated_at
FROM visa.visa_applications
WHERE departure_id = $1
ORDER BY created_at ASC
`

const getVisaApplicationsFiltered = `
SELECT id, jamaah_id, booking_id, departure_id, status, provider_id, provider_ref, e_visa_url, issued_date, created_at, updated_at
FROM visa.visa_applications
WHERE departure_id = $1 AND status = $2
ORDER BY created_at ASC
`

func (q *Queries) GetVisaApplicationsForDeparture(ctx context.Context, arg GetVisaApplicationsForDepartureParams) ([]VisaApplication, error) {
	var (
		rows interface{ Scan(...any) error }
	)
	if arg.StatusFilter == "" {
		r, err := q.db.Query(ctx, getVisaApplicationsAll, arg.DepartureID)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		var result []VisaApplication
		for r.Next() {
			var i VisaApplication
			if err := r.Scan(
				&i.ID, &i.JamaahID, &i.BookingID, &i.DepartureID,
				&i.Status, &i.ProviderID, &i.ProviderRef, &i.EVizaURL,
				&i.IssuedDate, &i.CreatedAt, &i.UpdatedAt,
			); err != nil {
				return nil, err
			}
			result = append(result, i)
		}
		return result, r.Err()
	}
	r, err := q.db.Query(ctx, getVisaApplicationsFiltered, arg.DepartureID, arg.StatusFilter)
	if err != nil {
		return nil, err
	}
	_ = rows
	defer r.Close()
	var result []VisaApplication
	for r.Next() {
		var i VisaApplication
		if err := r.Scan(
			&i.ID, &i.JamaahID, &i.BookingID, &i.DepartureID,
			&i.Status, &i.ProviderID, &i.ProviderRef, &i.EVizaURL,
			&i.IssuedDate, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, i)
	}
	return result, r.Err()
}

// ---------------------------------------------------------------------------
// GetStatusHistoryForApplication
// ---------------------------------------------------------------------------

const getStatusHistoryForApp = `
SELECT id, application_id, from_status, to_status, reason, actor_user_id, created_at
FROM visa.status_history
WHERE application_id = $1
ORDER BY id ASC
`

func (q *Queries) GetStatusHistoryForApplication(ctx context.Context, applicationID string) ([]StatusHistoryEntry, error) {
	rows, err := q.db.Query(ctx, getStatusHistoryForApp, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []StatusHistoryEntry
	for rows.Next() {
		var i StatusHistoryEntry
		if err := rows.Scan(
			&i.ID, &i.ApplicationID, &i.FromStatus, &i.ToStatus,
			&i.Reason, &i.ActorUserID, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, i)
	}
	return result, rows.Err()
}

// ---------------------------------------------------------------------------
// GetReadyVisaForJamaahDeparture
// ---------------------------------------------------------------------------

const getReadyVisaForJamaahDeparture = `
SELECT id, status FROM visa.visa_applications
WHERE jamaah_id = $1 AND departure_id = $2 AND status = 'READY'
`

type VisaIDStatus struct {
	ID     string
	Status string
}

func (q *Queries) GetReadyVisaForJamaahDeparture(ctx context.Context, jamaahID, departureID string) (VisaIDStatus, error) {
	row := q.db.QueryRow(ctx, getReadyVisaForJamaahDeparture, jamaahID, departureID)
	var i VisaIDStatus
	err := row.Scan(&i.ID, &i.Status)
	return i, err
}

// ---------------------------------------------------------------------------
// helper
// ---------------------------------------------------------------------------

func nullText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func nullTime(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}
