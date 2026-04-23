// scan_boarding.sql.go — hand-written SQLC-style queries for scan events and
// bus boarding roster (BL-OPS-010/011).

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// Models
// ---------------------------------------------------------------------------

type ScanEvent struct {
	ID             string             `json:"id"`
	ScanType       string             `json:"scan_type"`
	DepartureID    string             `json:"departure_id"`
	JamaahID       string             `json:"jamaah_id"`
	ScannedBy      string             `json:"scanned_by"`
	DeviceID       pgtype.Text        `json:"device_id"`
	Location       pgtype.Text        `json:"location"`
	IdempotencyKey string             `json:"idempotency_key"`
	Metadata       []byte             `json:"metadata"`
	ScannedAt      pgtype.Timestamptz `json:"scanned_at"`
}

type BusBoardingRow struct {
	ID          string             `json:"id"`
	DepartureID string             `json:"departure_id"`
	BusNumber   string             `json:"bus_number"`
	JamaahID    string             `json:"jamaah_id"`
	Status      string             `json:"status"`
	ScanEventID pgtype.Text        `json:"scan_event_id"`
	BoardedAt   pgtype.Timestamptz `json:"boarded_at"`
}

// ---------------------------------------------------------------------------
// InsertScanEventIdempotent
// ---------------------------------------------------------------------------

type InsertScanEventIdempotentParams struct {
	ID             string
	ScanType       string
	DepartureID    string
	JamaahID       string
	ScannedBy      string
	DeviceID       string
	Location       string
	IdempotencyKey string
	Metadata       []byte
}

const insertScanEventIdempotent = `
INSERT INTO ops.scan_events (id, scan_type, departure_id, jamaah_id, scanned_by, device_id, location, idempotency_key, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (idempotency_key) DO NOTHING
RETURNING id
`

const getScanByIdempotencyKey = `
SELECT id FROM ops.scan_events WHERE idempotency_key = $1
`

// InsertScanEventIdempotent inserts a scan event, returning (id, wasInserted, err).
func (q *Queries) InsertScanEventIdempotent(ctx context.Context, arg InsertScanEventIdempotentParams) (string, bool, error) {
	row := q.db.QueryRow(ctx, insertScanEventIdempotent,
		arg.ID, arg.ScanType, arg.DepartureID, arg.JamaahID, arg.ScannedBy,
		nullText(arg.DeviceID), nullText(arg.Location), arg.IdempotencyKey, arg.Metadata,
	)
	var id string
	err := row.Scan(&id)
	if err != nil {
		// ON CONFLICT DO NOTHING returns no row — fetch the existing id
		if err.Error() == "no rows in result set" || err.Error() == "pgx: rows closed" {
			row2 := q.db.QueryRow(ctx, getScanByIdempotencyKey, arg.IdempotencyKey)
			var existingID string
			if err2 := row2.Scan(&existingID); err2 != nil {
				return "", false, err2
			}
			return existingID, false, nil
		}
		return "", false, err
	}
	return id, true, nil
}

// ---------------------------------------------------------------------------
// GetBoardingByDeparturJamaah (for idempotency check)
// ---------------------------------------------------------------------------

const getBoardingByDepartureJamaah = `
SELECT id, departure_id, bus_number, jamaah_id, status, scan_event_id, boarded_at
FROM ops.bus_boardings
WHERE departure_id = $1 AND jamaah_id = $2
LIMIT 1
`

func (q *Queries) GetBoardingByDepartureJamaah(ctx context.Context, departureID, jamaahID string) (BusBoardingRow, error) {
	row := q.db.QueryRow(ctx, getBoardingByDepartureJamaah, departureID, jamaahID)
	var i BusBoardingRow
	err := row.Scan(
		&i.ID, &i.DepartureID, &i.BusNumber, &i.JamaahID,
		&i.Status, &i.ScanEventID, &i.BoardedAt,
	)
	return i, err
}

// ---------------------------------------------------------------------------
// InsertBusBoarding
// ---------------------------------------------------------------------------

type InsertBusBoardingParams struct {
	ID          string
	DepartureID string
	BusNumber   string
	JamaahID    string
	Status      string
	ScanEventID string
}

const insertBusBoarding = `
INSERT INTO ops.bus_boardings (id, departure_id, bus_number, jamaah_id, status, scan_event_id, boarded_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW())
RETURNING id, status
`

const insertBusBoardingAbsent = `
INSERT INTO ops.bus_boardings (id, departure_id, bus_number, jamaah_id, status, scan_event_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, status
`

func (q *Queries) InsertBusBoarding(ctx context.Context, arg InsertBusBoardingParams) (BusBoardingRow, error) {
	var query string
	if arg.Status == "absent" {
		query = insertBusBoardingAbsent
	} else {
		query = insertBusBoarding
	}
	row := q.db.QueryRow(ctx, query,
		arg.ID, arg.DepartureID, arg.BusNumber, arg.JamaahID,
		arg.Status, nullText(arg.ScanEventID),
	)
	var i BusBoardingRow
	err := row.Scan(&i.ID, &i.Status)
	if err != nil {
		return BusBoardingRow{}, err
	}
	i.DepartureID = arg.DepartureID
	i.BusNumber = arg.BusNumber
	i.JamaahID = arg.JamaahID
	return i, nil
}

// ---------------------------------------------------------------------------
// GetBoardingRoster
// ---------------------------------------------------------------------------

type GetBoardingRosterParams struct {
	DepartureID string
	BusNumber   string // optional
}

const getBoardingRosterAll = `
SELECT id, departure_id, bus_number, jamaah_id, status, scan_event_id, boarded_at
FROM ops.bus_boardings
WHERE departure_id = $1
ORDER BY jamaah_id ASC
`

const getBoardingRosterByBus = `
SELECT id, departure_id, bus_number, jamaah_id, status, scan_event_id, boarded_at
FROM ops.bus_boardings
WHERE departure_id = $1 AND bus_number = $2
ORDER BY jamaah_id ASC
`

func (q *Queries) GetBoardingRoster(ctx context.Context, arg GetBoardingRosterParams) ([]BusBoardingRow, error) {
	var (
		query string
		args  []any
	)
	if arg.BusNumber == "" {
		query = getBoardingRosterAll
		args = []any{arg.DepartureID}
	} else {
		query = getBoardingRosterByBus
		args = []any{arg.DepartureID, arg.BusNumber}
	}

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []BusBoardingRow
	for rows.Next() {
		var i BusBoardingRow
		if err := rows.Scan(
			&i.ID, &i.DepartureID, &i.BusNumber, &i.JamaahID,
			&i.Status, &i.ScanEventID, &i.BoardedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, i)
	}
	return result, rows.Err()
}
