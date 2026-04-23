// pl_reports.sql.go — hand-written sqlc-style query implementations for
// P&L report and Balance Sheet (finance depth / Wave 1B).
//
// LEFT JOIN finance.account_codes is used for account names.
// COALESCE guards ensure queries still work if migration 000017 has not been
// applied yet (account_codes table absent → LEFT JOIN returns NULLs).

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// ---------------------------------------------------------------------------
// Row types
// ---------------------------------------------------------------------------

// PLReportLineRow holds one aggregated P&L line (revenue or expense account).
type PLReportLineRow struct {
	AccountCode  string         `json:"account_code"`
	AccountName  string         `json:"account_name"`
	AccountType  string         `json:"account_type"`
	TotalCredit  pgtype.Numeric `json:"total_credit"`
	TotalDebit   pgtype.Numeric `json:"total_debit"`
}

// BalanceSheetLineRow holds one aggregated balance-sheet line (asset/liability/equity).
type BalanceSheetLineRow struct {
	AccountCode   string         `json:"account_code"`
	AccountName   string         `json:"account_name"`
	AccountType   string         `json:"account_type"`
	NormalBalance string         `json:"normal_balance"`
	TotalDebit    pgtype.Numeric `json:"total_debit"`
	TotalCredit   pgtype.Numeric `json:"total_credit"`
}

// ---------------------------------------------------------------------------
// Params
// ---------------------------------------------------------------------------

// GetPLReportLinesParams holds the date-range filter for the P&L query.
type GetPLReportLinesParams struct {
	From pgtype.Timestamptz // optional; zero = no lower bound
	To   pgtype.Timestamptz // optional; zero = no upper bound
}

// GetBalanceSheetLinesParams holds the as-of-date filter for the balance sheet query.
type GetBalanceSheetLinesParams struct {
	AsOf pgtype.Timestamptz // optional; zero = no upper bound (all history)
}

// ---------------------------------------------------------------------------
// Query implementations
// ---------------------------------------------------------------------------

const getPLReportLines = `-- name: GetPLReportLines :many
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
ORDER BY jl.account_code`

// GetPLReportLines returns aggregated revenue and expense account lines for the
// given date range. Accounts 4xxx = revenue, 5xxx = expense.
func (q *Queries) GetPLReportLines(ctx context.Context, arg GetPLReportLinesParams) ([]PLReportLineRow, error) {
	rows, err := q.db.Query(ctx, getPLReportLines, arg.From, arg.To)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []PLReportLineRow
	for rows.Next() {
		var r PLReportLineRow
		if err := rows.Scan(
			&r.AccountCode,
			&r.AccountName,
			&r.AccountType,
			&r.TotalCredit,
			&r.TotalDebit,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

const getBalanceSheetLines = `-- name: GetBalanceSheetLines :many
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
ORDER BY jl.account_code`

// GetBalanceSheetLines returns aggregated balance-sheet account lines as of
// (and including) the given date. Accounts 1xxx = asset, 2xxx = liability,
// 3xxx = equity.
func (q *Queries) GetBalanceSheetLines(ctx context.Context, arg GetBalanceSheetLinesParams) ([]BalanceSheetLineRow, error) {
	rows, err := q.db.Query(ctx, getBalanceSheetLines, arg.AsOf)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []BalanceSheetLineRow
	for rows.Next() {
		var r BalanceSheetLineRow
		if err := rows.Scan(
			&r.AccountCode,
			&r.AccountName,
			&r.AccountType,
			&r.NormalBalance,
			&r.TotalDebit,
			&r.TotalCredit,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}
