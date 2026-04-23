-- revenue.sql — queries for revenue recognition (finance depth / Wave 1B).
--
-- RecognizeRevenue: Dr 2001 (Pilgrim Liability) / Cr 4001 (Revenue).
-- Idempotency key: "revenue:<departure_id>".
-- Lookup is handled by the existing GetJournalEntryByIdempotencyKey query.

-- name: GetJournalEntryByIdempotencyKeyForRevenue :one
-- Reuses the same idempotency_key lookup as payment entries.
-- Kept here for documentation; implementation uses the shared query.
SELECT id, idempotency_key, source_type, source_id, posted_at, description
FROM finance.journal_entries
WHERE idempotency_key = $1
LIMIT 1;
