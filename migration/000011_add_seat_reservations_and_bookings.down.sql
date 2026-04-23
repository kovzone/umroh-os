-- Reverse of 000011_add_seat_reservations_and_bookings.up.sql

DROP TABLE IF EXISTS booking.booking_addons;
DROP TABLE IF EXISTS booking.booking_items;
DROP TABLE IF EXISTS booking.bookings;
DROP TYPE  IF EXISTS booking.item_status;
DROP TYPE  IF EXISTS booking.status;
DROP TYPE  IF EXISTS booking.channel;
DROP SCHEMA IF EXISTS booking;

DROP TABLE IF EXISTS catalog.seat_reservations;
