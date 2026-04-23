-- 000012 — rollback: drop payment schema and all tables.

DROP TABLE IF EXISTS payment.refunds;
DROP TABLE IF EXISTS payment.payment_events;
DROP TABLE IF EXISTS payment.virtual_accounts;
DROP TABLE IF EXISTS payment.invoices;
DROP SCHEMA IF EXISTS payment;
