// finance_depth.sql.go — hand-written SQLC-style queries for Wave 4
// Finance depth features (BL-FIN-020..041).

package sqlc

import (
	"context"
	"time"
)

// ---------------------------------------------------------------------------
// Models
// ---------------------------------------------------------------------------

type BankTransactionRow struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	RefNo       string    `json:"ref_no"`
	Amount      int64     `json:"amount"`
	TxDate      time.Time `json:"tx_date"`
	Description string    `json:"description"`
	Direction   string    `json:"direction"`
	Reconciled  bool      `json:"reconciled"`
	CreatedAt   time.Time `json:"created_at"`
}

type ReceiptRow struct {
	ID            string    `json:"id"`
	ReceiptNumber string    `json:"receipt_number"`
	BookingID     string    `json:"booking_id"`
	PaymentID     string    `json:"payment_id"`
	Amount        int64     `json:"amount"`
	Notes         string    `json:"notes"`
	IssuedAt      time.Time `json:"issued_at"`
}

type FinanceVendorRow struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	BankAccount  string    `json:"bank_account"`
	TaxID        string    `json:"tax_id"`
	ContactEmail string    `json:"contact_email"`
	CreatedAt    time.Time `json:"created_at"`
}

type PaymentAuthorizationRow struct {
	ID          string    `json:"id"`
	BatchID     string    `json:"batch_id"`
	Amount      int64     `json:"amount"`
	RequestedBy string    `json:"requested_by"`
	Level       int32     `json:"level"`
	Status      string    `json:"status"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
}

type PettyCashEntryRow struct {
	ID             string    `json:"id"`
	Amount         int64     `json:"amount"`
	Direction      string    `json:"direction"`
	Description    string    `json:"description"`
	Category       string    `json:"category"`
	Date           time.Time `json:"date"`
	EvidenceURL    string    `json:"evidence_url"`
	RunningBalance int64     `json:"running_balance"`
	CreatedAt      time.Time `json:"created_at"`
}

type ExchangeRateRow struct {
	ID            string    `json:"id"`
	FromCurrency  string    `json:"from_currency"`
	ToCurrency    string    `json:"to_currency"`
	Rate          string    `json:"rate"`
	EffectiveDate time.Time `json:"effective_date"`
	CreatedAt     time.Time `json:"created_at"`
}

type FixedAssetRow struct {
	ID                      string    `json:"id"`
	Name                    string    `json:"name"`
	Category                string    `json:"category"`
	PurchaseDate            time.Time `json:"purchase_date"`
	PurchaseCost            int64     `json:"purchase_cost"`
	UsefulLifeMonths        int32     `json:"useful_life_months"`
	ResidualValue           int64     `json:"residual_value"`
	AccumulatedDepreciation int64     `json:"accumulated_depreciation"`
	BookValue               int64     `json:"book_value"`
	CreatedAt               time.Time `json:"created_at"`
}

type CommissionPayoutRow struct {
	ID          string    `json:"id"`
	AgentID     string    `json:"agent_id"`
	DepartureID string    `json:"departure_id"`
	Amount      int64     `json:"amount"`
	BasisAmount int64     `json:"basis_amount"`
	RatePercent string    `json:"rate_percent"`
	Notes       string    `json:"notes"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type AuditLogEntryRow struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	Diff       string    `json:"diff"`
	CreatedAt  time.Time `json:"created_at"`
}

type FinanceConfigRow struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ---------------------------------------------------------------------------
// BankTransaction queries
// ---------------------------------------------------------------------------

type InsertBankTransactionParams struct {
	ID          string
	AccountID   string
	RefNo       string
	Amount      int64
	TxDate      time.Time
	Description string
	Direction   string
}

const insertBankTransaction = `
INSERT INTO finance.bank_transactions (id, account_id, ref_no, amount, tx_date, description, direction)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, account_id, ref_no, amount, tx_date, description, direction, reconciled, created_at
`

func (q *Queries) InsertBankTransaction(ctx context.Context, arg InsertBankTransactionParams) (BankTransactionRow, error) {
	row := q.db.QueryRow(ctx, insertBankTransaction,
		arg.ID, arg.AccountID, arg.RefNo, arg.Amount, arg.TxDate, arg.Description, arg.Direction,
	)
	var i BankTransactionRow
	err := row.Scan(
		&i.ID, &i.AccountID, &i.RefNo, &i.Amount, &i.TxDate, &i.Description, &i.Direction, &i.Reconciled, &i.CreatedAt,
	)
	return i, err
}

type GetBankTransactionsParams struct {
	AccountID string
	StartDate time.Time
	EndDate   time.Time
}

const getBankTransactions = `
SELECT id, account_id, ref_no, amount, tx_date, description, direction, reconciled, created_at
FROM finance.bank_transactions
WHERE account_id = $1 AND tx_date >= $2 AND tx_date <= $3
ORDER BY tx_date ASC
`

func (q *Queries) GetBankTransactions(ctx context.Context, arg GetBankTransactionsParams) ([]BankTransactionRow, error) {
	rows, err := q.db.Query(ctx, getBankTransactions, arg.AccountID, arg.StartDate, arg.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BankTransactionRow
	for rows.Next() {
		var i BankTransactionRow
		if err := rows.Scan(
			&i.ID, &i.AccountID, &i.RefNo, &i.Amount, &i.TxDate, &i.Description, &i.Direction, &i.Reconciled, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

// ---------------------------------------------------------------------------
// Receipt queries
// ---------------------------------------------------------------------------

type InsertReceiptParams struct {
	ID        string
	BookingID string
	PaymentID string
	Amount    int64
	Notes     string
}

const insertReceipt = `
INSERT INTO finance.receipts (id, receipt_number, booking_id, payment_id, amount, notes)
VALUES ($1, 'RCP-' || lpad(nextval('finance.receipt_number_seq')::text, 8, '0'), $2, $3, $4, $5)
RETURNING id, receipt_number, booking_id, payment_id, amount, notes, issued_at
`

func (q *Queries) InsertReceipt(ctx context.Context, arg InsertReceiptParams) (ReceiptRow, error) {
	row := q.db.QueryRow(ctx, insertReceipt, arg.ID, arg.BookingID, arg.PaymentID, arg.Amount, arg.Notes)
	var i ReceiptRow
	err := row.Scan(&i.ID, &i.ReceiptNumber, &i.BookingID, &i.PaymentID, &i.Amount, &i.Notes, &i.IssuedAt)
	return i, err
}

const getReceiptByID = `
SELECT id, receipt_number, booking_id, payment_id, amount, notes, issued_at
FROM finance.receipts
WHERE id = $1
`

func (q *Queries) GetReceiptByID(ctx context.Context, id string) (ReceiptRow, error) {
	row := q.db.QueryRow(ctx, getReceiptByID, id)
	var i ReceiptRow
	err := row.Scan(&i.ID, &i.ReceiptNumber, &i.BookingID, &i.PaymentID, &i.Amount, &i.Notes, &i.IssuedAt)
	return i, err
}

// ---------------------------------------------------------------------------
// Vendor queries
// ---------------------------------------------------------------------------

type InsertVendorParams struct {
	ID           string
	Name         string
	Category     string
	BankAccount  string
	TaxID        string
	ContactEmail string
}

const insertVendor = `
INSERT INTO finance.vendors (id, name, category, bank_account, tax_id, contact_email)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, category, bank_account, tax_id, contact_email, created_at
`

func (q *Queries) InsertVendor(ctx context.Context, arg InsertVendorParams) (FinanceVendorRow, error) {
	row := q.db.QueryRow(ctx, insertVendor, arg.ID, arg.Name, arg.Category, arg.BankAccount, arg.TaxID, arg.ContactEmail)
	var i FinanceVendorRow
	err := row.Scan(&i.ID, &i.Name, &i.Category, &i.BankAccount, &i.TaxID, &i.ContactEmail, &i.CreatedAt)
	return i, err
}

type UpdateVendorParams struct {
	ID           string
	Name         string
	Category     string
	BankAccount  string
	ContactEmail string
}

const updateVendor = `
UPDATE finance.vendors
SET name = $2, category = $3, bank_account = $4, contact_email = $5
WHERE id = $1
RETURNING id, name, category, bank_account, tax_id, contact_email, created_at
`

func (q *Queries) UpdateVendor(ctx context.Context, arg UpdateVendorParams) (FinanceVendorRow, error) {
	row := q.db.QueryRow(ctx, updateVendor, arg.ID, arg.Name, arg.Category, arg.BankAccount, arg.ContactEmail)
	var i FinanceVendorRow
	err := row.Scan(&i.ID, &i.Name, &i.Category, &i.BankAccount, &i.TaxID, &i.ContactEmail, &i.CreatedAt)
	return i, err
}

type ListVendorsParams struct {
	Category string
	Limit    int32
	Cursor   string
}

const listVendors = `
SELECT id, name, category, bank_account, tax_id, contact_email, created_at
FROM finance.vendors
WHERE ($1 = '' OR category = $1)
  AND ($3 = '' OR id > $3)
ORDER BY id ASC
LIMIT $2
`

func (q *Queries) ListVendors(ctx context.Context, arg ListVendorsParams) ([]FinanceVendorRow, error) {
	limit := arg.Limit
	if limit <= 0 {
		limit = 50
	}
	rows, err := q.db.Query(ctx, listVendors, arg.Category, limit, arg.Cursor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FinanceVendorRow
	for rows.Next() {
		var i FinanceVendorRow
		if err := rows.Scan(&i.ID, &i.Name, &i.Category, &i.BankAccount, &i.TaxID, &i.ContactEmail, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

const deleteVendor = `DELETE FROM finance.vendors WHERE id = $1`

func (q *Queries) DeleteVendor(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteVendor, id)
	return err
}

// ---------------------------------------------------------------------------
// Payment authorization queries
// ---------------------------------------------------------------------------

type InsertPaymentAuthorizationParams struct {
	ID          string
	BatchID     string
	Amount      int64
	RequestedBy string
	Level       int32
}

const insertPaymentAuthorization = `
INSERT INTO finance.payment_authorizations (id, batch_id, amount, requested_by, level, status)
VALUES ($1, $2, $3, $4, $5, 'pending')
RETURNING id, batch_id, amount, requested_by, level, status, notes, created_at
`

func (q *Queries) InsertPaymentAuthorization(ctx context.Context, arg InsertPaymentAuthorizationParams) (PaymentAuthorizationRow, error) {
	row := q.db.QueryRow(ctx, insertPaymentAuthorization, arg.ID, arg.BatchID, arg.Amount, arg.RequestedBy, arg.Level)
	var i PaymentAuthorizationRow
	err := row.Scan(&i.ID, &i.BatchID, &i.Amount, &i.RequestedBy, &i.Level, &i.Status, &i.Notes, &i.CreatedAt)
	return i, err
}

const listPendingAuthorizations = `
SELECT id, batch_id, amount, requested_by, level, status, notes, created_at
FROM finance.payment_authorizations
WHERE status = 'pending' AND ($1 = 0 OR level = $1)
ORDER BY created_at ASC
`

func (q *Queries) ListPendingAuthorizations(ctx context.Context, level int32) ([]PaymentAuthorizationRow, error) {
	rows, err := q.db.Query(ctx, listPendingAuthorizations, level)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PaymentAuthorizationRow
	for rows.Next() {
		var i PaymentAuthorizationRow
		if err := rows.Scan(&i.ID, &i.BatchID, &i.Amount, &i.RequestedBy, &i.Level, &i.Status, &i.Notes, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

type UpdateAuthorizationDecisionParams struct {
	ID     string
	Status string
	Notes  string
}

const updateAuthorizationDecision = `
UPDATE finance.payment_authorizations SET status = $2, notes = $3 WHERE id = $1
`

func (q *Queries) UpdateAuthorizationDecision(ctx context.Context, arg UpdateAuthorizationDecisionParams) error {
	_, err := q.db.Exec(ctx, updateAuthorizationDecision, arg.ID, arg.Status, arg.Notes)
	return err
}

// ---------------------------------------------------------------------------
// Petty cash queries
// ---------------------------------------------------------------------------

type InsertPettyCashEntryParams struct {
	ID          string
	Amount      int64
	Direction   string
	Description string
	Category    string
	Date        time.Time
	EvidenceURL string
}

const insertPettyCashEntry = `
INSERT INTO finance.petty_cash_entries (id, amount, direction, description, category, date, evidence_url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, amount, direction, description, category, date, evidence_url,
  (SELECT COALESCE(SUM(CASE WHEN direction='in' THEN amount ELSE -amount END), 0) FROM finance.petty_cash_entries WHERE closed_at IS NULL) as running_balance,
  created_at
`

func (q *Queries) InsertPettyCashEntry(ctx context.Context, arg InsertPettyCashEntryParams) (PettyCashEntryRow, error) {
	row := q.db.QueryRow(ctx, insertPettyCashEntry, arg.ID, arg.Amount, arg.Direction, arg.Description, arg.Category, arg.Date, arg.EvidenceURL)
	var i PettyCashEntryRow
	err := row.Scan(&i.ID, &i.Amount, &i.Direction, &i.Description, &i.Category, &i.Date, &i.EvidenceURL, &i.RunningBalance, &i.CreatedAt)
	return i, err
}

const getPettyCashBalance = `
SELECT COALESCE(SUM(CASE WHEN direction='in' THEN amount ELSE -amount END), 0) as balance
FROM finance.petty_cash_entries
WHERE closed_at IS NULL
`

func (q *Queries) GetPettyCashBalance(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getPettyCashBalance)
	var balance int64
	err := row.Scan(&balance)
	return balance, err
}

const closePettyCashEntries = `
UPDATE finance.petty_cash_entries SET closed_at = NOW() WHERE closed_at IS NULL
`

func (q *Queries) ClosePettyCashEntries(ctx context.Context) error {
	_, err := q.db.Exec(ctx, closePettyCashEntries)
	return err
}

// ---------------------------------------------------------------------------
// Exchange rate queries
// ---------------------------------------------------------------------------

type InsertExchangeRateParams struct {
	ID            string
	FromCurrency  string
	ToCurrency    string
	Rate          string
	EffectiveDate time.Time
}

const insertExchangeRate = `
INSERT INTO finance.exchange_rates (id, from_currency, to_currency, rate, effective_date)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (from_currency, to_currency, effective_date) DO UPDATE SET rate = EXCLUDED.rate
RETURNING id, from_currency, to_currency, rate, effective_date, created_at
`

func (q *Queries) InsertExchangeRate(ctx context.Context, arg InsertExchangeRateParams) (ExchangeRateRow, error) {
	row := q.db.QueryRow(ctx, insertExchangeRate, arg.ID, arg.FromCurrency, arg.ToCurrency, arg.Rate, arg.EffectiveDate)
	var i ExchangeRateRow
	err := row.Scan(&i.ID, &i.FromCurrency, &i.ToCurrency, &i.Rate, &i.EffectiveDate, &i.CreatedAt)
	return i, err
}

type GetExchangeRateParams struct {
	FromCurrency string
	ToCurrency   string
	AsOf         time.Time
}

const getExchangeRate = `
SELECT id, from_currency, to_currency, rate, effective_date, created_at
FROM finance.exchange_rates
WHERE from_currency = $1 AND to_currency = $2 AND effective_date <= $3
ORDER BY effective_date DESC
LIMIT 1
`

func (q *Queries) GetExchangeRate(ctx context.Context, arg GetExchangeRateParams) (ExchangeRateRow, error) {
	row := q.db.QueryRow(ctx, getExchangeRate, arg.FromCurrency, arg.ToCurrency, arg.AsOf)
	var i ExchangeRateRow
	err := row.Scan(&i.ID, &i.FromCurrency, &i.ToCurrency, &i.Rate, &i.EffectiveDate, &i.CreatedAt)
	return i, err
}

// ---------------------------------------------------------------------------
// Fixed asset queries
// ---------------------------------------------------------------------------

type InsertFixedAssetParams struct {
	ID               string
	Name             string
	Category         string
	PurchaseDate     time.Time
	PurchaseCost     int64
	UsefulLifeMonths int32
	ResidualValue    int64
}

const insertFixedAsset = `
INSERT INTO finance.fixed_assets (id, name, category, purchase_date, purchase_cost, useful_life_months, residual_value)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, category, purchase_date, purchase_cost, useful_life_months, residual_value,
  accumulated_depreciation, book_value, created_at
`

func (q *Queries) InsertFixedAsset(ctx context.Context, arg InsertFixedAssetParams) (FixedAssetRow, error) {
	row := q.db.QueryRow(ctx, insertFixedAsset,
		arg.ID, arg.Name, arg.Category, arg.PurchaseDate, arg.PurchaseCost, arg.UsefulLifeMonths, arg.ResidualValue,
	)
	var i FixedAssetRow
	err := row.Scan(
		&i.ID, &i.Name, &i.Category, &i.PurchaseDate, &i.PurchaseCost, &i.UsefulLifeMonths, &i.ResidualValue,
		&i.AccumulatedDepreciation, &i.BookValue, &i.CreatedAt,
	)
	return i, err
}

const listFixedAssets = `
SELECT id, name, category, purchase_date, purchase_cost, useful_life_months, residual_value,
  accumulated_depreciation, book_value, created_at
FROM finance.fixed_assets
WHERE ($1 = '' OR category = $1)
ORDER BY created_at ASC
`

func (q *Queries) ListFixedAssets(ctx context.Context, category string) ([]FixedAssetRow, error) {
	rows, err := q.db.Query(ctx, listFixedAssets, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FixedAssetRow
	for rows.Next() {
		var i FixedAssetRow
		if err := rows.Scan(
			&i.ID, &i.Name, &i.Category, &i.PurchaseDate, &i.PurchaseCost, &i.UsefulLifeMonths, &i.ResidualValue,
			&i.AccumulatedDepreciation, &i.BookValue, &i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

type UpdateFixedAssetDepreciationParams struct {
	ID                      string
	AccumulatedDepreciation int64
	BookValue               int64
}

const updateFixedAssetDepreciation = `
UPDATE finance.fixed_assets SET accumulated_depreciation = $2, book_value = $3 WHERE id = $1
`

func (q *Queries) UpdateFixedAssetDepreciation(ctx context.Context, arg UpdateFixedAssetDepreciationParams) error {
	_, err := q.db.Exec(ctx, updateFixedAssetDepreciation, arg.ID, arg.AccumulatedDepreciation, arg.BookValue)
	return err
}

// ---------------------------------------------------------------------------
// Commission payout queries
// ---------------------------------------------------------------------------

type InsertCommissionPayoutParams struct {
	ID          string
	AgentID     string
	DepartureID string
	Amount      int64
	BasisAmount int64
	RatePercent string
	Notes       string
}

const insertCommissionPayout = `
INSERT INTO finance.commission_payouts (id, agent_id, departure_id, amount, basis_amount, rate_percent, notes, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, 'pending_approval')
RETURNING id, agent_id, departure_id, amount, basis_amount, rate_percent, notes, status, created_at
`

func (q *Queries) InsertCommissionPayout(ctx context.Context, arg InsertCommissionPayoutParams) (CommissionPayoutRow, error) {
	row := q.db.QueryRow(ctx, insertCommissionPayout,
		arg.ID, arg.AgentID, arg.DepartureID, arg.Amount, arg.BasisAmount, arg.RatePercent, arg.Notes,
	)
	var i CommissionPayoutRow
	err := row.Scan(&i.ID, &i.AgentID, &i.DepartureID, &i.Amount, &i.BasisAmount, &i.RatePercent, &i.Notes, &i.Status, &i.CreatedAt)
	return i, err
}

const getCommissionPayoutByID = `
SELECT id, agent_id, departure_id, amount, basis_amount, rate_percent, notes, status, created_at
FROM finance.commission_payouts WHERE id = $1
`

func (q *Queries) GetCommissionPayoutByID(ctx context.Context, id string) (CommissionPayoutRow, error) {
	row := q.db.QueryRow(ctx, getCommissionPayoutByID, id)
	var i CommissionPayoutRow
	err := row.Scan(&i.ID, &i.AgentID, &i.DepartureID, &i.Amount, &i.BasisAmount, &i.RatePercent, &i.Notes, &i.Status, &i.CreatedAt)
	return i, err
}

const updateCommissionPayoutStatus = `
UPDATE finance.commission_payouts SET status = $2 WHERE id = $1
`

func (q *Queries) UpdateCommissionPayoutStatus(ctx context.Context, id, newStatus string) error {
	_, err := q.db.Exec(ctx, updateCommissionPayoutStatus, id, newStatus)
	return err
}

// ---------------------------------------------------------------------------
// Finance config (key-value store for policies)
// ---------------------------------------------------------------------------

const upsertFinanceConfig = `
INSERT INTO finance.config (key, value) VALUES ($1, $2)
ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
`

func (q *Queries) UpsertFinanceConfig(ctx context.Context, key, value string) error {
	_, err := q.db.Exec(ctx, upsertFinanceConfig, key, value)
	return err
}

const getFinanceConfig = `SELECT key, value, updated_at FROM finance.config WHERE key = $1`

func (q *Queries) GetFinanceConfig(ctx context.Context, key string) (FinanceConfigRow, error) {
	row := q.db.QueryRow(ctx, getFinanceConfig, key)
	var i FinanceConfigRow
	err := row.Scan(&i.Key, &i.Value, &i.UpdatedAt)
	return i, err
}

// ---------------------------------------------------------------------------
// Finance audit log queries
// ---------------------------------------------------------------------------

type InsertAuditLogParams struct {
	ID         string
	UserID     string
	Action     string
	EntityType string
	EntityID   string
	Diff       string
}

const insertAuditLog = `
INSERT INTO finance.audit_log (id, user_id, action, entity_type, entity_id, diff)
VALUES ($1, $2, $3, $4, $5, $6)
`

func (q *Queries) InsertFinanceAuditLog(ctx context.Context, arg InsertAuditLogParams) error {
	_, err := q.db.Exec(ctx, insertAuditLog, arg.ID, arg.UserID, arg.Action, arg.EntityType, arg.EntityID, arg.Diff)
	return err
}

type SearchAuditLogParams struct {
	UserID     string
	Action     string
	EntityType string
	EntityID   string
	StartDate  time.Time
	EndDate    time.Time
	Limit      int32
	Cursor     string
}

const searchAuditLog = `
SELECT id, user_id, action, entity_type, entity_id, diff, created_at
FROM finance.audit_log
WHERE ($1 = '' OR user_id = $1)
  AND ($2 = '' OR action = $2)
  AND ($3 = '' OR entity_type = $3)
  AND ($4 = '' OR entity_id = $4)
  AND ($5::timestamptz IS NULL OR created_at >= $5)
  AND ($6::timestamptz IS NULL OR created_at <= $6)
  AND ($7 = '' OR id > $7)
ORDER BY created_at DESC, id DESC
LIMIT $8
`

func (q *Queries) SearchAuditLog(ctx context.Context, arg SearchAuditLogParams) ([]AuditLogEntryRow, error) {
	limit := arg.Limit
	if limit <= 0 {
		limit = 50
	}
	rows, err := q.db.Query(ctx, searchAuditLog,
		arg.UserID, arg.Action, arg.EntityType, arg.EntityID,
		arg.StartDate, arg.EndDate, arg.Cursor, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AuditLogEntryRow
	for rows.Next() {
		var i AuditLogEntryRow
		if err := rows.Scan(&i.ID, &i.UserID, &i.Action, &i.EntityType, &i.EntityID, &i.Diff, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}
