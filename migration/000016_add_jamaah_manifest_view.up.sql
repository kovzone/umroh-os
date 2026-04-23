-- 000016 — Jamaah departure manifest view.
--
-- Creates a denormalized view joining booking.bookings and
-- booking.booking_items to produce a flat manifest row per pilgrim,
-- keyed by departure_id.  Used by the GetDepartureManifest RPC.
--
-- The view is intentionally non-materialized; the manifest endpoint is
-- an ops tool (low-frequency read) so a plain view is sufficient for MVP.

CREATE OR REPLACE VIEW jamaah.departure_manifest AS
SELECT
    b.departure_id,
    b.id            AS booking_id,
    bi.full_name    AS name,
    bi.email,
    bi.whatsapp     AS phone,
    b.room_type,
    b.status
FROM booking.bookings  b
JOIN booking.booking_items bi ON bi.booking_id = b.id
WHERE bi.status = 'active';

COMMENT ON VIEW jamaah.departure_manifest IS
    'Flat manifest of active pilgrims per departure, joining bookings + booking_items.';
