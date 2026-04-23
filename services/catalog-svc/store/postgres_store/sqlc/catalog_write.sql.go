// Hand-written sqlc-style query implementations for catalog write operations.
// Run `make generate` (sqlc generate) to regenerate from catalog_write.sql
// once sqlc v1.30+ picks up the new queries.
//
// S1-E-07 / BL-CAT-014 — staff catalog write (CreatePackage, UpdatePackage,
// DeletePackage, CreateDeparture, UpdateDeparture, ReserveSeats, ReleaseSeats).
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
// InsertPackage
// ---------------------------------------------------------------------------

const insertPackage = `-- name: InsertPackage :one
INSERT INTO catalog.packages (
    id,
    kind,
    name,
    description,
    highlights,
    cover_photo_url,
    itinerary_id,
    airline_id,
    muthawwif_id,
    status
) VALUES (
    $1,
    $2::catalog.package_kind,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10::catalog.package_status
)
RETURNING
    id,
    kind,
    name,
    description,
    highlights,
    cover_photo_url,
    itinerary_id,
    airline_id,
    muthawwif_id,
    status,
    created_at,
    updated_at`

type InsertPackageParams struct {
	ID            string               `json:"id"`
	Kind          CatalogPackageKind   `json:"kind"`
	Name          string               `json:"name"`
	Description   string               `json:"description"`
	Highlights    []string             `json:"highlights"`
	CoverPhotoUrl string               `json:"cover_photo_url"`
	ItineraryID   pgtype.Text          `json:"itinerary_id"`
	AirlineID     pgtype.Text          `json:"airline_id"`
	MuthawwifID   pgtype.Text          `json:"muthawwif_id"`
	Status        CatalogPackageStatus `json:"status"`
}

type InsertPackageRow struct {
	ID            string               `json:"id"`
	Kind          CatalogPackageKind   `json:"kind"`
	Name          string               `json:"name"`
	Description   string               `json:"description"`
	Highlights    []string             `json:"highlights"`
	CoverPhotoUrl string               `json:"cover_photo_url"`
	ItineraryID   pgtype.Text          `json:"itinerary_id"`
	AirlineID     pgtype.Text          `json:"airline_id"`
	MuthawwifID   pgtype.Text          `json:"muthawwif_id"`
	Status        CatalogPackageStatus `json:"status"`
	CreatedAt     pgtype.Timestamptz   `json:"created_at"`
	UpdatedAt     pgtype.Timestamptz   `json:"updated_at"`
}

func (q *Queries) InsertPackage(ctx context.Context, arg InsertPackageParams) (InsertPackageRow, error) {
	row := q.db.QueryRow(ctx, insertPackage,
		arg.ID,
		arg.Kind,
		arg.Name,
		arg.Description,
		arg.Highlights,
		arg.CoverPhotoUrl,
		arg.ItineraryID,
		arg.AirlineID,
		arg.MuthawwifID,
		arg.Status,
	)
	var i InsertPackageRow
	err := row.Scan(
		&i.ID,
		&i.Kind,
		&i.Name,
		&i.Description,
		&i.Highlights,
		&i.CoverPhotoUrl,
		&i.ItineraryID,
		&i.AirlineID,
		&i.MuthawwifID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// UpdatePackageFields
// ---------------------------------------------------------------------------

const updatePackageFields = `-- name: UpdatePackageFields :one
UPDATE catalog.packages
SET
    name            = COALESCE($1::text,                        name),
    description     = COALESCE($2::text,                        description),
    cover_photo_url = COALESCE($3::text,                        cover_photo_url),
    highlights      = COALESCE($4::text[],                      highlights),
    itinerary_id    = COALESCE($5::text,                        itinerary_id),
    airline_id      = COALESCE($6::text,                        airline_id),
    muthawwif_id    = COALESCE($7::text,                        muthawwif_id),
    status          = COALESCE($8::catalog.package_status,      status),
    updated_at      = now()
WHERE id = $9::text
  AND deleted_at IS NULL
RETURNING
    id,
    kind,
    name,
    description,
    highlights,
    cover_photo_url,
    itinerary_id,
    airline_id,
    muthawwif_id,
    status,
    created_at,
    updated_at`

type UpdatePackageFieldsParams struct {
	Name          pgtype.Text               `json:"name"`
	Description   pgtype.Text               `json:"description"`
	CoverPhotoUrl pgtype.Text               `json:"cover_photo_url"`
	Highlights    []string                  `json:"highlights"` // nil = no update
	ItineraryID   pgtype.Text               `json:"itinerary_id"`
	AirlineID     pgtype.Text               `json:"airline_id"`
	MuthawwifID   pgtype.Text               `json:"muthawwif_id"`
	Status        NullCatalogPackageStatus  `json:"status"`
	ID            string                    `json:"id"`
}

type UpdatePackageFieldsRow = InsertPackageRow

func (q *Queries) UpdatePackageFields(ctx context.Context, arg UpdatePackageFieldsParams) (UpdatePackageFieldsRow, error) {
	// For the highlights field: pass nil to postgres to trigger COALESCE no-op.
	var highlights interface{}
	if arg.Highlights != nil {
		highlights = arg.Highlights
	}

	row := q.db.QueryRow(ctx, updatePackageFields,
		arg.Name,
		arg.Description,
		arg.CoverPhotoUrl,
		highlights,
		arg.ItineraryID,
		arg.AirlineID,
		arg.MuthawwifID,
		arg.Status,
		arg.ID,
	)
	var i UpdatePackageFieldsRow
	err := row.Scan(
		&i.ID,
		&i.Kind,
		&i.Name,
		&i.Description,
		&i.Highlights,
		&i.CoverPhotoUrl,
		&i.ItineraryID,
		&i.AirlineID,
		&i.MuthawwifID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// SoftDeletePackage
// ---------------------------------------------------------------------------

const softDeletePackage = `-- name: SoftDeletePackage :one
UPDATE catalog.packages
SET
    deleted_at = now(),
    status     = 'archived',
    updated_at = now()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, status, deleted_at`

type SoftDeletePackageRow struct {
	ID        string               `json:"id"`
	Status    CatalogPackageStatus `json:"status"`
	DeletedAt pgtype.Timestamptz   `json:"deleted_at"`
}

func (q *Queries) SoftDeletePackage(ctx context.Context, id string) (SoftDeletePackageRow, error) {
	row := q.db.QueryRow(ctx, softDeletePackage, id)
	var i SoftDeletePackageRow
	err := row.Scan(&i.ID, &i.Status, &i.DeletedAt)
	return i, err
}

// ---------------------------------------------------------------------------
// GetPackageByIDForStaff
// ---------------------------------------------------------------------------

const getPackageByIDForStaff = `-- name: GetPackageByIDForStaff :one
SELECT
    id,
    kind,
    name,
    description,
    highlights,
    cover_photo_url,
    itinerary_id,
    airline_id,
    muthawwif_id,
    status,
    created_at,
    updated_at
FROM catalog.packages
WHERE id = $1
  AND deleted_at IS NULL`

type GetPackageByIDForStaffRow = GetActivePackageByIDRow

func (q *Queries) GetPackageByIDForStaff(ctx context.Context, id string) (GetPackageByIDForStaffRow, error) {
	row := q.db.QueryRow(ctx, getPackageByIDForStaff, id)
	var i GetPackageByIDForStaffRow
	err := row.Scan(
		&i.ID,
		&i.Kind,
		&i.Name,
		&i.Description,
		&i.Highlights,
		&i.CoverPhotoUrl,
		&i.ItineraryID,
		&i.AirlineID,
		&i.MuthawwifID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// DeletePackageHotels / InsertPackageHotel
// ---------------------------------------------------------------------------

const deletePackageHotels = `-- name: DeletePackageHotels :exec
DELETE FROM catalog.package_hotels WHERE package_id = $1`

func (q *Queries) DeletePackageHotels(ctx context.Context, packageID string) error {
	_, err := q.db.Exec(ctx, deletePackageHotels, packageID)
	return err
}

const insertPackageHotel = `-- name: InsertPackageHotel :exec
INSERT INTO catalog.package_hotels (package_id, hotel_id, sort_order)
VALUES ($1, $2, $3)
ON CONFLICT (package_id, hotel_id) DO UPDATE SET sort_order = EXCLUDED.sort_order`

type InsertPackageHotelParams struct {
	PackageID string `json:"package_id"`
	HotelID   string `json:"hotel_id"`
	SortOrder int16  `json:"sort_order"`
}

func (q *Queries) InsertPackageHotel(ctx context.Context, arg InsertPackageHotelParams) error {
	_, err := q.db.Exec(ctx, insertPackageHotel, arg.PackageID, arg.HotelID, arg.SortOrder)
	return err
}

// ---------------------------------------------------------------------------
// DeletePackageAddons / InsertPackageAddon
// ---------------------------------------------------------------------------

const deletePackageAddons = `-- name: DeletePackageAddons :exec
DELETE FROM catalog.package_addons WHERE package_id = $1`

func (q *Queries) DeletePackageAddons(ctx context.Context, packageID string) error {
	_, err := q.db.Exec(ctx, deletePackageAddons, packageID)
	return err
}

const insertPackageAddon = `-- name: InsertPackageAddon :exec
INSERT INTO catalog.package_addons (package_id, addon_id)
VALUES ($1, $2)
ON CONFLICT (package_id, addon_id) DO NOTHING`

type InsertPackageAddonParams struct {
	PackageID string `json:"package_id"`
	AddonID   string `json:"addon_id"`
}

func (q *Queries) InsertPackageAddon(ctx context.Context, arg InsertPackageAddonParams) error {
	_, err := q.db.Exec(ctx, insertPackageAddon, arg.PackageID, arg.AddonID)
	return err
}

// ---------------------------------------------------------------------------
// InsertDeparture
// ---------------------------------------------------------------------------

const insertDeparture = `-- name: InsertDeparture :one
INSERT INTO catalog.package_departures (
    id,
    package_id,
    departure_date,
    return_date,
    total_seats,
    status
) VALUES (
    $1,
    $2,
    $3::date,
    $4::date,
    $5,
    $6::catalog.departure_status
)
RETURNING
    id,
    package_id,
    departure_date,
    return_date,
    total_seats,
    reserved_seats,
    (total_seats - reserved_seats)::integer AS remaining_seats,
    status,
    created_at,
    updated_at`

type InsertDepartureParams struct {
	ID            string                 `json:"id"`
	PackageID     string                 `json:"package_id"`
	DepartureDate pgtype.Date            `json:"departure_date"`
	ReturnDate    pgtype.Date            `json:"return_date"`
	TotalSeats    int32                  `json:"total_seats"`
	Status        CatalogDepartureStatus `json:"status"`
}

type InsertDepartureRow struct {
	ID             string                 `json:"id"`
	PackageID      string                 `json:"package_id"`
	DepartureDate  pgtype.Date            `json:"departure_date"`
	ReturnDate     pgtype.Date            `json:"return_date"`
	TotalSeats     int32                  `json:"total_seats"`
	ReservedSeats  int32                  `json:"reserved_seats"`
	RemainingSeats int32                  `json:"remaining_seats"`
	Status         CatalogDepartureStatus `json:"status"`
	CreatedAt      pgtype.Timestamptz     `json:"created_at"`
	UpdatedAt      pgtype.Timestamptz     `json:"updated_at"`
}

func (q *Queries) InsertDeparture(ctx context.Context, arg InsertDepartureParams) (InsertDepartureRow, error) {
	row := q.db.QueryRow(ctx, insertDeparture,
		arg.ID,
		arg.PackageID,
		arg.DepartureDate,
		arg.ReturnDate,
		arg.TotalSeats,
		arg.Status,
	)
	var i InsertDepartureRow
	err := row.Scan(
		&i.ID,
		&i.PackageID,
		&i.DepartureDate,
		&i.ReturnDate,
		&i.TotalSeats,
		&i.ReservedSeats,
		&i.RemainingSeats,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// UpdateDepartureFields
// ---------------------------------------------------------------------------

const updateDepartureFields = `-- name: UpdateDepartureFields :one
UPDATE catalog.package_departures
SET
    departure_date = COALESCE($1::date,                         departure_date),
    return_date    = COALESCE($2::date,                         return_date),
    total_seats    = COALESCE($3::integer,                      total_seats),
    status         = COALESCE($4::catalog.departure_status,     status),
    updated_at     = now()
WHERE id = $5::text
RETURNING
    id,
    package_id,
    departure_date,
    return_date,
    total_seats,
    reserved_seats,
    (total_seats - reserved_seats)::integer AS remaining_seats,
    status,
    created_at,
    updated_at`

type UpdateDepartureFieldsParams struct {
	DepartureDate pgtype.Date                  `json:"departure_date"`
	ReturnDate    pgtype.Date                  `json:"return_date"`
	TotalSeats    pgtype.Int4                  `json:"total_seats"`
	Status        NullCatalogDepartureStatus   `json:"status"`
	ID            string                       `json:"id"`
}

type UpdateDepartureFieldsRow = InsertDepartureRow

func (q *Queries) UpdateDepartureFields(ctx context.Context, arg UpdateDepartureFieldsParams) (UpdateDepartureFieldsRow, error) {
	row := q.db.QueryRow(ctx, updateDepartureFields,
		arg.DepartureDate,
		arg.ReturnDate,
		arg.TotalSeats,
		arg.Status,
		arg.ID,
	)
	var i UpdateDepartureFieldsRow
	err := row.Scan(
		&i.ID,
		&i.PackageID,
		&i.DepartureDate,
		&i.ReturnDate,
		&i.TotalSeats,
		&i.ReservedSeats,
		&i.RemainingSeats,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// GetDepartureByIDForStaff
// ---------------------------------------------------------------------------

const getDepartureByIDForStaff = `-- name: GetDepartureByIDForStaff :one
SELECT
    id,
    package_id,
    departure_date,
    return_date,
    total_seats,
    reserved_seats,
    (total_seats - reserved_seats)::integer AS remaining_seats,
    status
FROM catalog.package_departures
WHERE id = $1`

type GetDepartureByIDForStaffRow = GetActiveDepartureRow

func (q *Queries) GetDepartureByIDForStaff(ctx context.Context, id string) (GetDepartureByIDForStaffRow, error) {
	row := q.db.QueryRow(ctx, getDepartureByIDForStaff, id)
	var i GetDepartureByIDForStaffRow
	err := row.Scan(
		&i.ID,
		&i.PackageID,
		&i.DepartureDate,
		&i.ReturnDate,
		&i.TotalSeats,
		&i.ReservedSeats,
		&i.RemainingSeats,
		&i.Status,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// DeleteDeparturePricing / InsertDeparturePricing
// ---------------------------------------------------------------------------

const deleteDeparturePricing = `-- name: DeleteDeparturePricing :exec
DELETE FROM catalog.package_pricing WHERE package_departure_id = $1`

func (q *Queries) DeleteDeparturePricing(ctx context.Context, packageDepartureID string) error {
	_, err := q.db.Exec(ctx, deleteDeparturePricing, packageDepartureID)
	return err
}

const insertDeparturePricing = `-- name: InsertDeparturePricing :one
INSERT INTO catalog.package_pricing (
    id,
    package_departure_id,
    room_type,
    list_amount,
    list_currency,
    settlement_currency
) VALUES (
    $1,
    $2,
    $3::catalog.room_type,
    $4,
    $5,
    $6
)
RETURNING
    id,
    package_departure_id,
    room_type,
    list_amount::bigint AS list_amount,
    list_currency,
    settlement_currency`

type InsertDeparturePricingParams struct {
	ID                 string          `json:"id"`
	PackageDepartureID string          `json:"package_departure_id"`
	RoomType           CatalogRoomType `json:"room_type"`
	ListAmount         int64           `json:"list_amount"`
	ListCurrency       string          `json:"list_currency"`
	SettlementCurrency string          `json:"settlement_currency"`
}

type InsertDeparturePricingRow struct {
	ID                 string          `json:"id"`
	PackageDepartureID string          `json:"package_departure_id"`
	RoomType           CatalogRoomType `json:"room_type"`
	ListAmount         int64           `json:"list_amount"`
	ListCurrency       string          `json:"list_currency"`
	SettlementCurrency string          `json:"settlement_currency"`
}

func (q *Queries) InsertDeparturePricing(ctx context.Context, arg InsertDeparturePricingParams) (InsertDeparturePricingRow, error) {
	row := q.db.QueryRow(ctx, insertDeparturePricing,
		arg.ID,
		arg.PackageDepartureID,
		arg.RoomType,
		arg.ListAmount,
		arg.ListCurrency,
		arg.SettlementCurrency,
	)
	var i InsertDeparturePricingRow
	err := row.Scan(
		&i.ID,
		&i.PackageDepartureID,
		&i.RoomType,
		&i.ListAmount,
		&i.ListCurrency,
		&i.SettlementCurrency,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// ReserveSeatsAtomic
// ---------------------------------------------------------------------------

const reserveSeatsAtomic = `-- name: ReserveSeatsAtomic :one
UPDATE catalog.package_departures
SET
    reserved_seats = reserved_seats + $1::integer,
    updated_at     = now()
WHERE id = $2::text
  AND status IN ('open', 'closed')
  AND reserved_seats + $1::integer <= total_seats
RETURNING
    id,
    reserved_seats,
    total_seats,
    (total_seats - reserved_seats)::integer AS remaining_seats`

type ReserveSeatsAtomicParams struct {
	Seats       int32  `json:"seats"`
	DepartureID string `json:"departure_id"`
}

type ReserveSeatsAtomicRow struct {
	ID             string `json:"id"`
	ReservedSeats  int32  `json:"reserved_seats"`
	TotalSeats     int32  `json:"total_seats"`
	RemainingSeats int32  `json:"remaining_seats"`
}

func (q *Queries) ReserveSeatsAtomic(ctx context.Context, arg ReserveSeatsAtomicParams) (ReserveSeatsAtomicRow, error) {
	row := q.db.QueryRow(ctx, reserveSeatsAtomic, arg.Seats, arg.DepartureID)
	var i ReserveSeatsAtomicRow
	err := row.Scan(&i.ID, &i.ReservedSeats, &i.TotalSeats, &i.RemainingSeats)
	return i, err
}

// ---------------------------------------------------------------------------
// InsertSeatReservation
// ---------------------------------------------------------------------------

const insertSeatReservation = `-- name: InsertSeatReservation :one
INSERT INTO catalog.seat_reservations (
    reservation_id,
    departure_id,
    seats,
    expires_at
) VALUES ($1, $2, $3, $4)
RETURNING reservation_id, departure_id, seats, reserved_at, expires_at`

type InsertSeatReservationParams struct {
	ReservationID string    `json:"reservation_id"`
	DepartureID   string    `json:"departure_id"`
	Seats         int32     `json:"seats"`
	ExpiresAt     time.Time `json:"expires_at"`
}

type SeatReservationRow struct {
	ReservationID string             `json:"reservation_id"`
	DepartureID   string             `json:"departure_id"`
	Seats         int32              `json:"seats"`
	ReservedAt    pgtype.Timestamptz `json:"reserved_at"`
	ExpiresAt     pgtype.Timestamptz `json:"expires_at"`
	ReleasedAt    pgtype.Timestamptz `json:"released_at"`
}

func (q *Queries) InsertSeatReservation(ctx context.Context, arg InsertSeatReservationParams) (SeatReservationRow, error) {
	row := q.db.QueryRow(ctx, insertSeatReservation,
		arg.ReservationID,
		arg.DepartureID,
		arg.Seats,
		arg.ExpiresAt,
	)
	var i SeatReservationRow
	err := row.Scan(
		&i.ReservationID,
		&i.DepartureID,
		&i.Seats,
		&i.ReservedAt,
		&i.ExpiresAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// GetSeatReservation
// ---------------------------------------------------------------------------

const getSeatReservation = `-- name: GetSeatReservation :one
SELECT reservation_id, departure_id, seats, reserved_at, expires_at, released_at
FROM catalog.seat_reservations
WHERE reservation_id = $1`

func (q *Queries) GetSeatReservation(ctx context.Context, reservationID string) (SeatReservationRow, error) {
	row := q.db.QueryRow(ctx, getSeatReservation, reservationID)
	var i SeatReservationRow
	err := row.Scan(
		&i.ReservationID,
		&i.DepartureID,
		&i.Seats,
		&i.ReservedAt,
		&i.ExpiresAt,
		&i.ReleasedAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// ReleaseSeatsAtomic
// ---------------------------------------------------------------------------

const releaseSeatsAtomic = `-- name: ReleaseSeatsAtomic :one
UPDATE catalog.package_departures
SET
    reserved_seats = reserved_seats - $1::integer,
    updated_at     = now()
WHERE id = $2::text
  AND reserved_seats - $1::integer >= 0
RETURNING
    id,
    reserved_seats,
    total_seats,
    (total_seats - reserved_seats)::integer AS remaining_seats`

type ReleaseSeatsAtomicParams struct {
	SeatsToRelease int32  `json:"seats_to_release"`
	DepartureID    string `json:"departure_id"`
}

type ReleaseSeatsAtomicRow = ReserveSeatsAtomicRow

func (q *Queries) ReleaseSeatsAtomic(ctx context.Context, arg ReleaseSeatsAtomicParams) (ReleaseSeatsAtomicRow, error) {
	row := q.db.QueryRow(ctx, releaseSeatsAtomic, arg.SeatsToRelease, arg.DepartureID)
	var i ReleaseSeatsAtomicRow
	err := row.Scan(&i.ID, &i.ReservedSeats, &i.TotalSeats, &i.RemainingSeats)
	return i, err
}

// ---------------------------------------------------------------------------
// MarkReservationReleased
// ---------------------------------------------------------------------------

const markReservationReleased = `-- name: MarkReservationReleased :one
UPDATE catalog.seat_reservations
SET released_at = now()
WHERE reservation_id = $1
  AND released_at IS NULL
RETURNING reservation_id, departure_id, seats, reserved_at, expires_at, released_at`

func (q *Queries) MarkReservationReleased(ctx context.Context, reservationID string) (SeatReservationRow, error) {
	row := q.db.QueryRow(ctx, markReservationReleased, reservationID)
	var i SeatReservationRow
	err := row.Scan(
		&i.ReservationID,
		&i.DepartureID,
		&i.Seats,
		&i.ReservedAt,
		&i.ExpiresAt,
		&i.ReleasedAt,
	)
	return i, err
}
