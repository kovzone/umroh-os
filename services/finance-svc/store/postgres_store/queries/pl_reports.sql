-- pl_reports.sql — SQL queries for P&L report and Balance Sheet (finance depth / Wave 1B).
--
-- All queries LEFT JOIN finance.account_codes for human-readable account names.
-- COALESCE ensures graceful degradation if migration 000017 has not been run yet.

-- name: GetPLReportLines :many
-- Returns aggregated revenue (4xxx) and expense (5xxx) lines for a date range.
-- from_date and to_date are inclusive (posted_at >= from AND posted_at <= to).
SELECT
    jl.account_code,
    COALESCE(ac.name, jl.account_code)  AS account_name,
    COALESCE(ac.type, 'unknown')         AS account_type,
    SUM(jl.credit) AS total_credit,
    SUM(jl.debit)  AS total_debit
FROM finance.journal_lines jl
JOIN finance.journal_entries je ON je.id = jl.entry_id
LEFT JOIN finance.account_codes ac ON ac.code = jl.account_code
WHERE
    ($1::timestamptz IS NULL OR je.posted_at >= $1)
    AND ($2::timestamptz IS NULL OR je.posted_at <= $2)
    AND (jl.account_code LIKE '4%' OR jl.account_code LIKE '5%')
GROUP BY jl.account_code, ac.name, ac.type
ORDER BY jl.account_code;

-- name: GetBalanceSheetLines :many
-- Returns net balances for all balance-sheet accounts (1xxx, 2xxx, 3xxx)
-- as of (and including) as_of_date.
SELECT
    jl.account_code,
    COALESCE(ac.name, jl.account_code)           AS account_name,
    COALESCE(ac.type, 'unknown')                  AS account_type,
    COALESCE(ac.normal_balance, 'debit')          AS normal_balance,
    SUM(jl.debit)  AS total_debit,
    SUM(jl.credit) AS total_credit
FROM finance.journal_lines jl
JOIN finance.journal_entries je ON je.id = jl.entry_id
LEFT JOIN finance.account_codes ac ON ac.code = jl.account_code
WHERE
    ($1::timestamptz IS NULL OR je.posted_at <= $1)
    AND (jl.account_code LIKE '1%' OR jl.account_code LIKE '2%' OR jl.account_code LIKE '3%')
GROUP BY jl.account_code, ac.name, ac.type, ac.normal_balance
ORDER BY jl.account_code;
