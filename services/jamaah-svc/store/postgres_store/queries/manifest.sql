-- manifest.sql — departure manifest query joining bookings + pilgrim_documents.
--
-- Wave-1A manifest API.  The query returns one row per active booking_item
-- (pilgrim) for a given departure, with document completion counts aggregated
-- from jamaah.pilgrim_documents.
--
-- name: GetDepartureManifest :many
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
ORDER BY bi.full_name ASC
