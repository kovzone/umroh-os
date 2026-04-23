// manifest.sql.go — hand-written sqlc-style query implementation for the
// departure manifest query.
//
// Wave-1A manifest API.  Run `make generate` once sqlc is configured to
// regenerate this from manifest.sql.

package sqlc

import (
	"context"
)

// ---------------------------------------------------------------------------
// Row type
// ---------------------------------------------------------------------------

// GetDepartureManifestRow is one row returned by GetDepartureManifest.
// booking.bookings does not have a NIK column; the NIK field is populated
// as an empty string until a pilgrim profile schema adds it.
type GetDepartureManifestRow struct {
	BookingID     string `json:"booking_id"`
	Name          string `json:"name"`
	NIK           string `json:"nik"`
	Phone         string `json:"phone"`
	RoomType      string `json:"room_type"`
	BookingStatus string `json:"booking_status"`
	ApprovedDocs  int64  `json:"approved_docs"`
	TotalDocs     int64  `json:"total_docs"`
}

// ---------------------------------------------------------------------------
// Query
// ---------------------------------------------------------------------------

const getDepartureManifest = `-- name: GetDepartureManifest :many
SELECT
    b.id                                                            AS booking_id,
    bi.full_name                                                    AS name,
    ''::TEXT                                                        AS nik,
    bi.whatsapp                                                     AS phone,
    b.room_type,
    b.status                                                        AS booking_status,
    COUNT(pd.id) FILTER (WHERE pd.status = 'approved')             AS approved_docs,
    COUNT(pd.id)                                                    AS total_docs
FROM booking.bookings b
JOIN booking.booking_items bi ON bi.booking_id = b.id AND bi.status = 'active'
LEFT JOIN jamaah.pilgrim_documents pd ON pd.booking_id = b.id
WHERE b.departure_id = $1
  AND b.status NOT IN ('draft', 'cancelled')
GROUP BY b.id, bi.full_name, bi.whatsapp, b.room_type, b.status
ORDER BY bi.full_name ASC`

// GetDepartureManifest returns all active pilgrims for a departure with their
// document completion counts.
func (q *Queries) GetDepartureManifest(ctx context.Context, departureID string) ([]GetDepartureManifestRow, error) {
	rows, err := q.db.Query(ctx, getDepartureManifest, departureID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []GetDepartureManifestRow
	for rows.Next() {
		var r GetDepartureManifestRow
		if err := rows.Scan(
			&r.BookingID,
			&r.Name,
			&r.NIK,
			&r.Phone,
			&r.RoomType,
			&r.BookingStatus,
			&r.ApprovedDocs,
			&r.TotalDocs,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}
