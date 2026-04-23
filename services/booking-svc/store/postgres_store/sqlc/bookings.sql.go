// Hand-written sqlc-style query implementations for booking-svc draft creation.
// Run `make generate` (sqlc generate) to regenerate from bookings.sql once
// sqlc configuration supports the booking schema types.
//
// S1-E-03 / BL-BOOK-001..006 — booking draft (CreateBooking, InsertBookingItem,
// InsertBookingAddon, GetBookingByID, ListBookingItems, ListBookingAddons,
// GetBookingByIdempotencyKey).
//
// NOTE: This file intentionally mimics sqlc output style so the regenerated
// file will be a drop-in replacement. Do not add business logic here.

package sqlc

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// Enum types for the booking schema
// ---------------------------------------------------------------------------

// BookingStatus mirrors booking.status PostgreSQL enum.
type BookingStatus string

const (
	BookingStatusDraft          BookingStatus = "draft"
	BookingStatusPendingPayment BookingStatus = "pending_payment"
	BookingStatusPartiallyPaid  BookingStatus = "partially_paid"
	BookingStatusPaidInFull     BookingStatus = "paid_in_full"
	BookingStatusDeparted       BookingStatus = "departed"
	BookingStatusCompleted      BookingStatus = "completed"
	BookingStatusExpired        BookingStatus = "expired"
	BookingStatusCancelled      BookingStatus = "cancelled"
	BookingStatusFailed         BookingStatus = "failed"
)

// BookingChannel mirrors booking.channel PostgreSQL enum.
type BookingChannel string

const (
	BookingChannelB2cSelf  BookingChannel = "b2c_self"
	BookingChannelB2bAgent BookingChannel = "b2b_agent"
	BookingChannelCs       BookingChannel = "cs"
)

// BookingItemStatus mirrors booking.item_status PostgreSQL enum.
type BookingItemStatus string

const (
	BookingItemStatusActive    BookingItemStatus = "active"
	BookingItemStatusCancelled BookingItemStatus = "cancelled"
)

// ---------------------------------------------------------------------------
// Row types
// ---------------------------------------------------------------------------

// BookingRow maps to the columns returned by InsertBooking / GetBookingByID.
type BookingRow struct {
	ID                 string             `json:"id"`
	Status             BookingStatus      `json:"status"`
	Channel            BookingChannel     `json:"channel"`
	PackageID          string             `json:"package_id"`
	DepartureID        string             `json:"departure_id"`
	RoomType           string             `json:"room_type"`
	AgentID            pgtype.Text        `json:"agent_id"`
	StaffUserID        pgtype.Text        `json:"staff_user_id"`
	LeadFullName       string             `json:"lead_full_name"`
	LeadEmail          pgtype.Text        `json:"lead_email"`
	LeadWhatsapp       string             `json:"lead_whatsapp"`
	LeadDomicile       string             `json:"lead_domicile"`
	ListAmount         int64              `json:"list_amount"`
	ListCurrency       string             `json:"list_currency"`
	SettlementCurrency string             `json:"settlement_currency"`
	Notes              pgtype.Text        `json:"notes"`
	IdempotencyKey     pgtype.Text        `json:"idempotency_key"`
	IdempotencyBodyHash pgtype.Text       `json:"idempotency_body_hash"`
	CreatedAt          pgtype.Timestamptz `json:"created_at"`
	UpdatedAt          pgtype.Timestamptz `json:"updated_at"`
	ExpiresAt          pgtype.Timestamptz `json:"expires_at"`
}

// BookingItemRow maps to booking.booking_items columns.
type BookingItemRow struct {
	ID        string             `json:"id"`
	BookingID string             `json:"booking_id"`
	FullName  string             `json:"full_name"`
	Email     pgtype.Text        `json:"email"`
	Whatsapp  pgtype.Text        `json:"whatsapp"`
	Domicile  string             `json:"domicile"`
	IsLead    bool               `json:"is_lead"`
	Status    BookingItemStatus  `json:"status"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

// BookingAddonRow maps to booking.booking_addons columns.
type BookingAddonRow struct {
	BookingID          string `json:"booking_id"`
	AddonID            string `json:"addon_id"`
	AddonName          string `json:"addon_name"`
	ListAmount         int64  `json:"list_amount"`
	ListCurrency       string `json:"list_currency"`
	SettlementCurrency string `json:"settlement_currency"`
}

// ---------------------------------------------------------------------------
// InsertBooking
// ---------------------------------------------------------------------------

const insertBooking = `-- name: InsertBooking :one
INSERT INTO booking.bookings (
    id,
    status,
    channel,
    package_id,
    departure_id,
    room_type,
    agent_id,
    staff_user_id,
    lead_full_name,
    lead_email,
    lead_whatsapp,
    lead_domicile,
    list_amount,
    list_currency,
    settlement_currency,
    notes,
    idempotency_key,
    idempotency_body_hash,
    expires_at
) VALUES (
    $1,
    'draft'::booking.status,
    $2::booking.channel,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15,
    $16,
    $17,
    $18
)
RETURNING
    id, status, channel,
    package_id, departure_id, room_type,
    agent_id, staff_user_id,
    lead_full_name, lead_email, lead_whatsapp, lead_domicile,
    list_amount, list_currency, settlement_currency,
    notes, idempotency_key, idempotency_body_hash,
    created_at, updated_at, expires_at`

type InsertBookingParams struct {
	ID                  string         `json:"id"`
	Channel             BookingChannel `json:"channel"`
	PackageID           string         `json:"package_id"`
	DepartureID         string         `json:"departure_id"`
	RoomType            string         `json:"room_type"`
	AgentID             pgtype.Text    `json:"agent_id"`
	StaffUserID         pgtype.Text    `json:"staff_user_id"`
	LeadFullName        string         `json:"lead_full_name"`
	LeadEmail           pgtype.Text    `json:"lead_email"`
	LeadWhatsapp        string         `json:"lead_whatsapp"`
	LeadDomicile        string         `json:"lead_domicile"`
	ListAmount          int64          `json:"list_amount"`
	ListCurrency        string         `json:"list_currency"`
	SettlementCurrency  string         `json:"settlement_currency"`
	Notes               pgtype.Text    `json:"notes"`
	IdempotencyKey      pgtype.Text    `json:"idempotency_key"`
	IdempotencyBodyHash pgtype.Text    `json:"idempotency_body_hash"`
	ExpiresAt           time.Time      `json:"expires_at"`
}

func (q *Queries) InsertBooking(ctx context.Context, arg InsertBookingParams) (BookingRow, error) {
	row := q.db.QueryRow(ctx, insertBooking,
		arg.ID,
		arg.Channel,
		arg.PackageID,
		arg.DepartureID,
		arg.RoomType,
		arg.AgentID,
		arg.StaffUserID,
		arg.LeadFullName,
		arg.LeadEmail,
		arg.LeadWhatsapp,
		arg.LeadDomicile,
		arg.ListAmount,
		arg.ListCurrency,
		arg.SettlementCurrency,
		arg.Notes,
		arg.IdempotencyKey,
		arg.IdempotencyBodyHash,
		arg.ExpiresAt,
	)
	var i BookingRow
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Channel,
		&i.PackageID,
		&i.DepartureID,
		&i.RoomType,
		&i.AgentID,
		&i.StaffUserID,
		&i.LeadFullName,
		&i.LeadEmail,
		&i.LeadWhatsapp,
		&i.LeadDomicile,
		&i.ListAmount,
		&i.ListCurrency,
		&i.SettlementCurrency,
		&i.Notes,
		&i.IdempotencyKey,
		&i.IdempotencyBodyHash,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// InsertBookingItem
// ---------------------------------------------------------------------------

const insertBookingItem = `-- name: InsertBookingItem :one
INSERT INTO booking.booking_items (
    id,
    booking_id,
    full_name,
    email,
    whatsapp,
    domicile,
    is_lead,
    status
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    'active'::booking.item_status
)
RETURNING
    id, booking_id, full_name, email, whatsapp, domicile, is_lead, status, created_at, updated_at`

type InsertBookingItemParams struct {
	ID        string      `json:"id"`
	BookingID string      `json:"booking_id"`
	FullName  string      `json:"full_name"`
	Email     pgtype.Text `json:"email"`
	Whatsapp  pgtype.Text `json:"whatsapp"`
	Domicile  string      `json:"domicile"`
	IsLead    bool        `json:"is_lead"`
}

func (q *Queries) InsertBookingItem(ctx context.Context, arg InsertBookingItemParams) (BookingItemRow, error) {
	row := q.db.QueryRow(ctx, insertBookingItem,
		arg.ID,
		arg.BookingID,
		arg.FullName,
		arg.Email,
		arg.Whatsapp,
		arg.Domicile,
		arg.IsLead,
	)
	var i BookingItemRow
	err := row.Scan(
		&i.ID,
		&i.BookingID,
		&i.FullName,
		&i.Email,
		&i.Whatsapp,
		&i.Domicile,
		&i.IsLead,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// InsertBookingAddon
// ---------------------------------------------------------------------------

const insertBookingAddon = `-- name: InsertBookingAddon :exec
INSERT INTO booking.booking_addons (
    booking_id,
    addon_id,
    addon_name,
    list_amount,
    list_currency,
    settlement_currency
) VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (booking_id, addon_id) DO NOTHING`

type InsertBookingAddonParams struct {
	BookingID          string `json:"booking_id"`
	AddonID            string `json:"addon_id"`
	AddonName          string `json:"addon_name"`
	ListAmount         int64  `json:"list_amount"`
	ListCurrency       string `json:"list_currency"`
	SettlementCurrency string `json:"settlement_currency"`
}

func (q *Queries) InsertBookingAddon(ctx context.Context, arg InsertBookingAddonParams) error {
	_, err := q.db.Exec(ctx, insertBookingAddon,
		arg.BookingID,
		arg.AddonID,
		arg.AddonName,
		arg.ListAmount,
		arg.ListCurrency,
		arg.SettlementCurrency,
	)
	return err
}

// ---------------------------------------------------------------------------
// GetBookingByID
// ---------------------------------------------------------------------------

const getBookingByID = `-- name: GetBookingByID :one
SELECT
    id, status, channel,
    package_id, departure_id, room_type,
    agent_id, staff_user_id,
    lead_full_name, lead_email, lead_whatsapp, lead_domicile,
    list_amount, list_currency, settlement_currency,
    notes, idempotency_key,
    created_at, updated_at, expires_at
FROM booking.bookings
WHERE id = $1
  AND status = 'draft'`

func (q *Queries) GetBookingByID(ctx context.Context, id string) (BookingRow, error) {
	row := q.db.QueryRow(ctx, getBookingByID, id)
	var i BookingRow
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Channel,
		&i.PackageID,
		&i.DepartureID,
		&i.RoomType,
		&i.AgentID,
		&i.StaffUserID,
		&i.LeadFullName,
		&i.LeadEmail,
		&i.LeadWhatsapp,
		&i.LeadDomicile,
		&i.ListAmount,
		&i.ListCurrency,
		&i.SettlementCurrency,
		&i.Notes,
		&i.IdempotencyKey,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// ListBookingItems
// ---------------------------------------------------------------------------

const listBookingItems = `-- name: ListBookingItems :many
SELECT
    id, booking_id, full_name, email, whatsapp, domicile, is_lead, status, created_at, updated_at
FROM booking.booking_items
WHERE booking_id = $1
ORDER BY is_lead DESC, created_at ASC`

func (q *Queries) ListBookingItems(ctx context.Context, bookingID string) ([]BookingItemRow, error) {
	rows, err := q.db.Query(ctx, listBookingItems, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []BookingItemRow
	for rows.Next() {
		var i BookingItemRow
		if err := rows.Scan(
			&i.ID,
			&i.BookingID,
			&i.FullName,
			&i.Email,
			&i.Whatsapp,
			&i.Domicile,
			&i.IsLead,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// ---------------------------------------------------------------------------
// ListBookingAddons
// ---------------------------------------------------------------------------

const listBookingAddons = `-- name: ListBookingAddons :many
SELECT
    booking_id, addon_id, addon_name, list_amount, list_currency, settlement_currency
FROM booking.booking_addons
WHERE booking_id = $1`

func (q *Queries) ListBookingAddons(ctx context.Context, bookingID string) ([]BookingAddonRow, error) {
	rows, err := q.db.Query(ctx, listBookingAddons, bookingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addons []BookingAddonRow
	for rows.Next() {
		var a BookingAddonRow
		if err := rows.Scan(
			&a.BookingID,
			&a.AddonID,
			&a.AddonName,
			&a.ListAmount,
			&a.ListCurrency,
			&a.SettlementCurrency,
		); err != nil {
			return nil, err
		}
		addons = append(addons, a)
	}
	return addons, rows.Err()
}

// ---------------------------------------------------------------------------
// GetBookingByIdempotencyKey
// ---------------------------------------------------------------------------

const getBookingByIdempotencyKey = `-- name: GetBookingByIdempotencyKey :one
SELECT
    id, status, channel,
    package_id, departure_id, room_type,
    agent_id, staff_user_id,
    lead_full_name, lead_email, lead_whatsapp, lead_domicile,
    list_amount, list_currency, settlement_currency,
    notes, idempotency_key, idempotency_body_hash,
    created_at, updated_at, expires_at
FROM booking.bookings
WHERE channel = $1::booking.channel
  AND idempotency_key = $2`

type GetBookingByIdempotencyKeyParams struct {
	Channel        BookingChannel `json:"channel"`
	IdempotencyKey string         `json:"idempotency_key"`
}

func (q *Queries) GetBookingByIdempotencyKey(ctx context.Context, arg GetBookingByIdempotencyKeyParams) (BookingRow, error) {
	row := q.db.QueryRow(ctx, getBookingByIdempotencyKey, arg.Channel, arg.IdempotencyKey)
	var i BookingRow
	err := row.Scan(
		&i.ID,
		&i.Status,
		&i.Channel,
		&i.PackageID,
		&i.DepartureID,
		&i.RoomType,
		&i.AgentID,
		&i.StaffUserID,
		&i.LeadFullName,
		&i.LeadEmail,
		&i.LeadWhatsapp,
		&i.LeadDomicile,
		&i.ListAmount,
		&i.ListCurrency,
		&i.SettlementCurrency,
		&i.Notes,
		&i.IdempotencyKey,
		&i.IdempotencyBodyHash,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
	)
	return i, err
}
