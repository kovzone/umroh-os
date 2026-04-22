-- 000010_seed_catalog_hidden_departure.up.sql
-- Adds a cancelled departure so the 404 departure_not_found path on
-- GET /v1/package-departures/{id} can be exercised by e2e tests.
-- Separate from 000009 because migrations landed on dev are immutable
-- — we extend by appending a new migration, not by editing.

INSERT INTO catalog.package_departures (
    id, package_id, departure_date, return_date,
    total_seats, reserved_seats, status
) VALUES (
    'dep_01JCDF00000000000000000003',
    'pkg_01JCDE00000000000000000001',
    CURRENT_DATE + INTERVAL '120 days',
    CURRENT_DATE + INTERVAL '132 days',
    45, 0, 'cancelled'
) ON CONFLICT (id) DO NOTHING;
