-- Reverse of 000008_add_catalog_packages_and_masters.
-- Drops children before parents; drops enums after the last dependent
-- table is gone. Schema stays (other future catalog tables may still
-- need it); drop it manually if the whole catalog surface is retired.

DROP TABLE IF EXISTS catalog.package_pricing;
DROP TABLE IF EXISTS catalog.package_departures;
DROP TABLE IF EXISTS catalog.package_addons;
DROP TABLE IF EXISTS catalog.package_hotels;
DROP TABLE IF EXISTS catalog.packages;
DROP TABLE IF EXISTS catalog.addons;
DROP TABLE IF EXISTS catalog.muthawwif;
DROP TABLE IF EXISTS catalog.airlines;
DROP TABLE IF EXISTS catalog.hotels;
DROP TABLE IF EXISTS catalog.itinerary_templates;

DROP TYPE IF EXISTS catalog.operator_kind;
DROP TYPE IF EXISTS catalog.room_type;
DROP TYPE IF EXISTS catalog.departure_status;
DROP TYPE IF EXISTS catalog.package_status;
DROP TYPE IF EXISTS catalog.package_kind;
