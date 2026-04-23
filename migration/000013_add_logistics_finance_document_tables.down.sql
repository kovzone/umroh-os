-- 000013 down — Drop tables added in 000013 up.

DROP TABLE IF EXISTS jamaah.pilgrim_documents;
DROP TABLE IF EXISTS finance.journal_lines;
DROP TABLE IF EXISTS finance.journal_entries;
DROP TABLE IF EXISTS logistics.fulfillment_tasks;

-- Drop schemas only if empty (safety guard).
DROP SCHEMA IF EXISTS jamaah;
DROP SCHEMA IF EXISTS finance;
DROP SCHEMA IF EXISTS logistics;
