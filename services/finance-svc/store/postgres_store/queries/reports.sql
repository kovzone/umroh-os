-- reports.sql — queries for finance reports (S5-E-01).
-- GetFinanceSummary and ListJournalEntries for GET /v1/finance/summary
-- and GET /v1/finance/journals.

-- name: GetFinanceSummary :many
-- Returns aggregate debit and credit per account_code from finance.journal_lines.
-- net = debit_total - credit_total (positive = net debit, negative = net credit).
SELECT
    account_code,
    SUM(debit)  AS debit_total,
    SUM(credit) AS credit_total
FROM finance.journal_lines
GROUP BY account_code
ORDER BY account_code;

-- name: ListJournalEntries :many
-- Returns journal entries ordered by created_at DESC with optional date-range filter
-- and cursor-based pagination. Filter params are optional (pass NULL to skip).
-- Cursor is the posted_at timestamp of the last seen entry (exclusive).
SELECT
    je.id,
    je.idempotency_key,
    je.source_type,
    je.source_id,
    je.posted_at,
    je.description
FROM finance.journal_entries je
WHERE
    ($1::timestamptz IS NULL OR je.posted_at >= $1)   -- from
    AND ($2::timestamptz IS NULL OR je.posted_at <= $2) -- to
    AND ($3::timestamptz IS NULL OR je.posted_at < $3)  -- cursor (posted_at of last seen entry)
ORDER BY je.posted_at DESC
LIMIT $4; -- limit

-- name: GetJournalLines :many
-- Returns all lines for the given list of entry IDs.
-- Used to batch-load lines after fetching the entry list.
SELECT
    jl.id,
    jl.entry_id,
    jl.account_code,
    jl.debit,
    jl.credit
FROM finance.journal_lines jl
WHERE jl.entry_id = ANY($1::uuid[])
ORDER BY jl.entry_id, jl.id;
