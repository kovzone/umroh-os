// Hand-written sqlc-style query implementations for catalog master data CRUD.
// Run `make generate` (sqlc generate) to regenerate from masters.sql once
// sqlc is available.
//
// Covers: hotels, airlines, muthawwif, addons, package_pricing (upsert + list).

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// Hotels
// ---------------------------------------------------------------------------

const insertHotel = `-- name: InsertHotel :one
INSERT INTO catalog.hotels (id, name, city, star_rating, walking_distance_m)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, city, star_rating, walking_distance_m, created_at, updated_at`

type InsertHotelParams struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	City             string `json:"city"`
	StarRating       int16  `json:"star_rating"`
	WalkingDistanceM int32  `json:"walking_distance_m"`
}

type HotelRow struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	City             string             `json:"city"`
	StarRating       int16              `json:"star_rating"`
	WalkingDistanceM int32              `json:"walking_distance_m"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) InsertHotel(ctx context.Context, arg InsertHotelParams) (HotelRow, error) {
	row := q.db.QueryRow(ctx, insertHotel,
		arg.ID, arg.Name, arg.City, arg.StarRating, arg.WalkingDistanceM)
	var i HotelRow
	err := row.Scan(&i.ID, &i.Name, &i.City, &i.StarRating, &i.WalkingDistanceM,
		&i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const updateHotelFields = `-- name: UpdateHotelFields :one
UPDATE catalog.hotels
SET
    name               = COALESCE($1::text,     name),
    city               = COALESCE($2::text,     city),
    star_rating        = COALESCE($3::smallint, star_rating),
    walking_distance_m = COALESCE($4::integer,  walking_distance_m),
    updated_at         = now()
WHERE id = $5::text
RETURNING id, name, city, star_rating, walking_distance_m, created_at, updated_at`

type UpdateHotelFieldsParams struct {
	Name             pgtype.Text  `json:"name"`
	City             pgtype.Text  `json:"city"`
	StarRating       pgtype.Int2  `json:"star_rating"`
	WalkingDistanceM pgtype.Int4  `json:"walking_distance_m"`
	ID               string       `json:"id"`
}

func (q *Queries) UpdateHotelFields(ctx context.Context, arg UpdateHotelFieldsParams) (HotelRow, error) {
	row := q.db.QueryRow(ctx, updateHotelFields,
		arg.Name, arg.City, arg.StarRating, arg.WalkingDistanceM, arg.ID)
	var i HotelRow
	err := row.Scan(&i.ID, &i.Name, &i.City, &i.StarRating, &i.WalkingDistanceM,
		&i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const deleteHotel = `-- name: DeleteHotel :exec
DELETE FROM catalog.hotels WHERE id = $1`

func (q *Queries) DeleteHotel(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteHotel, id)
	return err
}

const getHotelByID = `-- name: GetHotelByID :one
SELECT id, name, city, star_rating, walking_distance_m, created_at, updated_at
FROM catalog.hotels WHERE id = $1`

func (q *Queries) GetHotelByID(ctx context.Context, id string) (HotelRow, error) {
	row := q.db.QueryRow(ctx, getHotelByID, id)
	var i HotelRow
	err := row.Scan(&i.ID, &i.Name, &i.City, &i.StarRating, &i.WalkingDistanceM,
		&i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const listHotels = `-- name: ListHotels :many
SELECT id, name, city, star_rating, walking_distance_m, created_at, updated_at
FROM catalog.hotels
WHERE ($1::text = '' OR id > $1::text)
ORDER BY id
LIMIT $2`

func (q *Queries) ListHotels(ctx context.Context, cursor string, limit int32) ([]HotelRow, error) {
	rows, err := q.db.Query(ctx, listHotels, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []HotelRow
	for rows.Next() {
		var i HotelRow
		if err := rows.Scan(&i.ID, &i.Name, &i.City, &i.StarRating, &i.WalkingDistanceM,
			&i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const countHotelPackageRefs = `-- name: CountHotelPackageRefs :one
SELECT COUNT(*)::bigint AS ref_count FROM catalog.package_hotels WHERE hotel_id = $1`

func (q *Queries) CountHotelPackageRefs(ctx context.Context, hotelID string) (int64, error) {
	row := q.db.QueryRow(ctx, countHotelPackageRefs, hotelID)
	var n int64
	err := row.Scan(&n)
	return n, err
}

// ---------------------------------------------------------------------------
// Airlines
// ---------------------------------------------------------------------------

const insertAirline = `-- name: InsertAirline :one
INSERT INTO catalog.airlines (id, code, name, operator_kind)
VALUES ($1, $2, $3, $4::catalog.operator_kind)
RETURNING id, code, name, operator_kind, created_at, updated_at`

type InsertAirlineParams struct {
	ID           string              `json:"id"`
	Code         string              `json:"code"`
	Name         string              `json:"name"`
	OperatorKind CatalogOperatorKind `json:"operator_kind"`
}

type AirlineRow struct {
	ID           string              `json:"id"`
	Code         string              `json:"code"`
	Name         string              `json:"name"`
	OperatorKind CatalogOperatorKind `json:"operator_kind"`
	CreatedAt    pgtype.Timestamptz  `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz  `json:"updated_at"`
}

func (q *Queries) InsertAirline(ctx context.Context, arg InsertAirlineParams) (AirlineRow, error) {
	row := q.db.QueryRow(ctx, insertAirline,
		arg.ID, arg.Code, arg.Name, arg.OperatorKind)
	var i AirlineRow
	err := row.Scan(&i.ID, &i.Code, &i.Name, &i.OperatorKind, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const updateAirlineFields = `-- name: UpdateAirlineFields :one
UPDATE catalog.airlines
SET
    code       = COALESCE($1::text, code),
    name       = COALESCE($2::text, name),
    updated_at = now()
WHERE id = $3::text
RETURNING id, code, name, operator_kind, created_at, updated_at`

type UpdateAirlineFieldsParams struct {
	Code pgtype.Text `json:"code"`
	Name pgtype.Text `json:"name"`
	ID   string      `json:"id"`
}

func (q *Queries) UpdateAirlineFields(ctx context.Context, arg UpdateAirlineFieldsParams) (AirlineRow, error) {
	row := q.db.QueryRow(ctx, updateAirlineFields, arg.Code, arg.Name, arg.ID)
	var i AirlineRow
	err := row.Scan(&i.ID, &i.Code, &i.Name, &i.OperatorKind, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const countAirlinePackageRefs = `-- name: CountAirlinePackageRefs :one
SELECT COUNT(*)::bigint AS ref_count FROM catalog.package_airlines WHERE airline_id = $1`

func (q *Queries) CountAirlinePackageRefs(ctx context.Context, airlineID string) (int64, error) {
	row := q.db.QueryRow(ctx, countAirlinePackageRefs, airlineID)
	var n int64
	err := row.Scan(&n)
	return n, err
}

const deleteAirline = `-- name: DeleteAirline :exec
DELETE FROM catalog.airlines WHERE id = $1`

func (q *Queries) DeleteAirline(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteAirline, id)
	return err
}

const getAirlineByIDForStaff = `-- name: GetAirlineByIDForStaff :one
SELECT id, code, name, operator_kind, created_at, updated_at
FROM catalog.airlines WHERE id = $1`

func (q *Queries) GetAirlineByIDForStaff(ctx context.Context, id string) (AirlineRow, error) {
	row := q.db.QueryRow(ctx, getAirlineByIDForStaff, id)
	var i AirlineRow
	err := row.Scan(&i.ID, &i.Code, &i.Name, &i.OperatorKind, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const listAirlines = `-- name: ListAirlines :many
SELECT id, code, name, operator_kind, created_at, updated_at
FROM catalog.airlines
WHERE ($1::text = '' OR id > $1::text)
ORDER BY id
LIMIT $2`

func (q *Queries) ListAirlines(ctx context.Context, cursor string, limit int32) ([]AirlineRow, error) {
	rows, err := q.db.Query(ctx, listAirlines, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AirlineRow
	for rows.Next() {
		var i AirlineRow
		if err := rows.Scan(&i.ID, &i.Code, &i.Name, &i.OperatorKind, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// ---------------------------------------------------------------------------
// Muthawwif
// ---------------------------------------------------------------------------

const insertMuthawwif = `-- name: InsertMuthawwif :one
INSERT INTO catalog.muthawwif (id, name, portrait_url)
VALUES ($1, $2, $3)
RETURNING id, name, portrait_url, created_at, updated_at`

type InsertMuthawwifParams struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PortraitUrl string `json:"portrait_url"`
}

type MuthawwifRow struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	PortraitUrl string             `json:"portrait_url"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) InsertMuthawwif(ctx context.Context, arg InsertMuthawwifParams) (MuthawwifRow, error) {
	row := q.db.QueryRow(ctx, insertMuthawwif, arg.ID, arg.Name, arg.PortraitUrl)
	var i MuthawwifRow
	err := row.Scan(&i.ID, &i.Name, &i.PortraitUrl, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const updateMuthawwifFields = `-- name: UpdateMuthawwifFields :one
UPDATE catalog.muthawwif
SET
    name         = COALESCE($1::text, name),
    portrait_url = COALESCE($2::text, portrait_url),
    updated_at   = now()
WHERE id = $3::text
RETURNING id, name, portrait_url, created_at, updated_at`

type UpdateMuthawwifFieldsParams struct {
	Name        pgtype.Text `json:"name"`
	PortraitUrl pgtype.Text `json:"portrait_url"`
	ID          string      `json:"id"`
}

func (q *Queries) UpdateMuthawwifFields(ctx context.Context, arg UpdateMuthawwifFieldsParams) (MuthawwifRow, error) {
	row := q.db.QueryRow(ctx, updateMuthawwifFields, arg.Name, arg.PortraitUrl, arg.ID)
	var i MuthawwifRow
	err := row.Scan(&i.ID, &i.Name, &i.PortraitUrl, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const countMuthawwifPackageRefs = `-- name: CountMuthawwifPackageRefs :one
SELECT COUNT(*)::bigint AS ref_count FROM catalog.package_muthawwif WHERE muthawwif_id = $1`

func (q *Queries) CountMuthawwifPackageRefs(ctx context.Context, muthawwifID string) (int64, error) {
	row := q.db.QueryRow(ctx, countMuthawwifPackageRefs, muthawwifID)
	var n int64
	err := row.Scan(&n)
	return n, err
}

const deleteMuthawwif = `-- name: DeleteMuthawwif :exec
DELETE FROM catalog.muthawwif WHERE id = $1`

func (q *Queries) DeleteMuthawwif(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteMuthawwif, id)
	return err
}

const getMuthawwifByIDForStaff = `-- name: GetMuthawwifByIDForStaff :one
SELECT id, name, portrait_url, created_at, updated_at
FROM catalog.muthawwif WHERE id = $1`

func (q *Queries) GetMuthawwifByIDForStaff(ctx context.Context, id string) (MuthawwifRow, error) {
	row := q.db.QueryRow(ctx, getMuthawwifByIDForStaff, id)
	var i MuthawwifRow
	err := row.Scan(&i.ID, &i.Name, &i.PortraitUrl, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const listMuthawwif = `-- name: ListMuthawwif :many
SELECT id, name, portrait_url, created_at, updated_at
FROM catalog.muthawwif
WHERE ($1::text = '' OR id > $1::text)
ORDER BY id
LIMIT $2`

func (q *Queries) ListMuthawwif(ctx context.Context, cursor string, limit int32) ([]MuthawwifRow, error) {
	rows, err := q.db.Query(ctx, listMuthawwif, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MuthawwifRow
	for rows.Next() {
		var i MuthawwifRow
		if err := rows.Scan(&i.ID, &i.Name, &i.PortraitUrl, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// ---------------------------------------------------------------------------
// Addons
// ---------------------------------------------------------------------------

const insertAddon = `-- name: InsertAddon :one
INSERT INTO catalog.addons (id, name, list_amount, list_currency, settlement_currency)
VALUES ($1, $2, $3, 'IDR', 'IDR')
RETURNING id, name, list_amount::bigint AS list_amount, list_currency, settlement_currency, created_at, updated_at`

type InsertAddonParams struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	ListAmount int64  `json:"list_amount"`
}

type AddonRow struct {
	ID                 string             `json:"id"`
	Name               string             `json:"name"`
	ListAmount         int64              `json:"list_amount"`
	ListCurrency       string             `json:"list_currency"`
	SettlementCurrency string             `json:"settlement_currency"`
	CreatedAt          pgtype.Timestamptz `json:"created_at"`
	UpdatedAt          pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) InsertAddon(ctx context.Context, arg InsertAddonParams) (AddonRow, error) {
	row := q.db.QueryRow(ctx, insertAddon, arg.ID, arg.Name, arg.ListAmount)
	var i AddonRow
	err := row.Scan(&i.ID, &i.Name, &i.ListAmount, &i.ListCurrency, &i.SettlementCurrency,
		&i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const updateAddonFields = `-- name: UpdateAddonFields :one
UPDATE catalog.addons
SET
    name        = COALESCE($1::text,    name),
    list_amount = COALESCE($2::numeric, list_amount),
    updated_at  = now()
WHERE id = $3::text
RETURNING id, name, list_amount::bigint AS list_amount, list_currency, settlement_currency, created_at, updated_at`

type UpdateAddonFieldsParams struct {
	Name       pgtype.Text    `json:"name"`
	ListAmount pgtype.Numeric `json:"list_amount"`
	ID         string         `json:"id"`
}

func (q *Queries) UpdateAddonFields(ctx context.Context, arg UpdateAddonFieldsParams) (AddonRow, error) {
	row := q.db.QueryRow(ctx, updateAddonFields, arg.Name, arg.ListAmount, arg.ID)
	var i AddonRow
	err := row.Scan(&i.ID, &i.Name, &i.ListAmount, &i.ListCurrency, &i.SettlementCurrency,
		&i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const deleteAddon = `-- name: DeleteAddon :exec
DELETE FROM catalog.addons WHERE id = $1`

func (q *Queries) DeleteAddon(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteAddon, id)
	return err
}

const getAddonByID = `-- name: GetAddonByID :one
SELECT id, name, list_amount::bigint AS list_amount, list_currency, settlement_currency, created_at, updated_at
FROM catalog.addons WHERE id = $1`

func (q *Queries) GetAddonByID(ctx context.Context, id string) (AddonRow, error) {
	row := q.db.QueryRow(ctx, getAddonByID, id)
	var i AddonRow
	err := row.Scan(&i.ID, &i.Name, &i.ListAmount, &i.ListCurrency, &i.SettlementCurrency,
		&i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const listAddons = `-- name: ListAddons :many
SELECT id, name, list_amount::bigint AS list_amount, list_currency, settlement_currency, created_at, updated_at
FROM catalog.addons
WHERE ($1::text = '' OR id > $1::text)
ORDER BY id
LIMIT $2`

func (q *Queries) ListAddons(ctx context.Context, cursor string, limit int32) ([]AddonRow, error) {
	rows, err := q.db.Query(ctx, listAddons, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AddonRow
	for rows.Next() {
		var i AddonRow
		if err := rows.Scan(&i.ID, &i.Name, &i.ListAmount, &i.ListCurrency, &i.SettlementCurrency,
			&i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// ---------------------------------------------------------------------------
// Package Departure Pricing
// ---------------------------------------------------------------------------

const upsertDeparturePricing = `-- name: UpsertDeparturePricing :one
INSERT INTO catalog.package_pricing (
    id, package_departure_id, room_type, list_amount, list_currency, settlement_currency
) VALUES (
    $1, $2, $3::catalog.room_type, $4, COALESCE(NULLIF($5, ''), 'IDR'), COALESCE(NULLIF($6, ''), 'IDR')
)
ON CONFLICT (package_departure_id, room_type) DO UPDATE
    SET
        list_amount         = EXCLUDED.list_amount,
        list_currency       = EXCLUDED.list_currency,
        settlement_currency = EXCLUDED.settlement_currency,
        updated_at          = now()
RETURNING id, package_departure_id, room_type, list_amount::bigint AS list_amount, list_currency, settlement_currency, created_at, updated_at`

type UpsertDeparturePricingParams struct {
	ID                 string          `json:"id"`
	PackageDepartureID string          `json:"package_departure_id"`
	RoomType           CatalogRoomType `json:"room_type"`
	ListAmount         int64           `json:"list_amount"`
	ListCurrency       string          `json:"list_currency"`
	SettlementCurrency string          `json:"settlement_currency"`
}

type PricingRow struct {
	ID                 string             `json:"id"`
	PackageDepartureID string             `json:"package_departure_id"`
	RoomType           CatalogRoomType    `json:"room_type"`
	ListAmount         int64              `json:"list_amount"`
	ListCurrency       string             `json:"list_currency"`
	SettlementCurrency string             `json:"settlement_currency"`
	CreatedAt          pgtype.Timestamptz `json:"created_at"`
	UpdatedAt          pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) UpsertDeparturePricing(ctx context.Context, arg UpsertDeparturePricingParams) (PricingRow, error) {
	row := q.db.QueryRow(ctx, upsertDeparturePricing,
		arg.ID, arg.PackageDepartureID, arg.RoomType, arg.ListAmount,
		arg.ListCurrency, arg.SettlementCurrency)
	var i PricingRow
	err := row.Scan(&i.ID, &i.PackageDepartureID, &i.RoomType, &i.ListAmount,
		&i.ListCurrency, &i.SettlementCurrency, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const getDeparturePricingRows = `-- name: GetDeparturePricingRows :many
SELECT id, package_departure_id, room_type, list_amount::bigint AS list_amount,
       list_currency, settlement_currency, created_at, updated_at
FROM catalog.package_pricing
WHERE package_departure_id = $1
ORDER BY room_type`

func (q *Queries) GetDeparturePricingRows(ctx context.Context, departureID string) ([]PricingRow, error) {
	rows, err := q.db.Query(ctx, getDeparturePricingRows, departureID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PricingRow
	for rows.Next() {
		var i PricingRow
		if err := rows.Scan(&i.ID, &i.PackageDepartureID, &i.RoomType, &i.ListAmount,
			&i.ListCurrency, &i.SettlementCurrency, &i.CreatedAt, &i.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}
