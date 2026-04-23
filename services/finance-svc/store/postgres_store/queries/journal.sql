-- journal.sql — queries for finance.journal_entries and finance.journal_lines.
-- S3-E-03 / BL-FIN-001..003.

-- name: GetJournalEntryByIdempotencyKey :one
-- Returns an existing journal entry for the given idempotency key, or no rows.
SELECT id, idempotency_key, source_type, source_id, posted_at, description
FROM finance.journal_entries
WHERE idempotency_key = $1
LIMIT 1;

-- name: InsertJournalEntry :one
-- Inserts a new journal entry header.
INSERT INTO finance.journal_entries (
    idempotency_key,
    source_type,
    source_id,
    description
) VALUES (
    $1, -- idempotency_key TEXT
    $2, -- source_type     TEXT
    $3, -- source_id       UUID
    $4  -- description     TEXT (nullable)
)
RETURNING id, idempotency_key, source_type, source_id, posted_at, description;

-- name: InsertJournalLine :one
-- Inserts one journal line.  Exactly one of debit/credit must be > 0
-- (enforced by DB CHECK constraint).
INSERT INTO finance.journal_lines (
    entry_id,
    account_code,
    debit,
    credit
) VALUES (
    $1, -- entry_id     UUID
    $2, -- account_code TEXT
    $3, -- debit        NUMERIC(15,2)
    $4  -- credit       NUMERIC(15,2)
)
RETURNING id, entry_id, account_code, debit, credit;
