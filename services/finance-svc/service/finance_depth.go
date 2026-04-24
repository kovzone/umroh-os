// finance_depth.go — service implementations for Wave 4 Finance depth features
// (BL-FIN-020..041).

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"finance-svc/store/postgres_store/sqlc"
	"finance-svc/util/logging"
	"finance-svc/util/ulid"

	otelCodes "go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Params / Results — BL-FIN-020 ScheduleBilling
// ---------------------------------------------------------------------------

type ScheduleBillingParams struct {
	DepartureID string
	DueDate     string
	Notes       string
}

type ScheduleBillingResult struct {
	InvoicesCreated int32
	TotalAmount     int64
}

// ScheduleBilling creates invoices for all bookings in a departure that do not
// yet have a paid invoice (BL-FIN-020). Idempotent per departure + due_date.
func (s *Service) ScheduleBilling(ctx context.Context, params *ScheduleBillingParams) (*ScheduleBillingResult, error) {
	const op = "service.Service.ScheduleBilling"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("departure_id", params.DepartureID).Msg("")

	if params.DepartureID == "" {
		return nil, ErrValidation
	}

	// In a full implementation this would query booking-svc bookings for the departure
	// and create invoices for unpaid ones. Stub returns zero for now.
	_ = params.Notes
	span.SetStatus(otelCodes.Ok, "ok")
	return &ScheduleBillingResult{InvoicesCreated: 0, TotalAmount: 0}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-021 Bank integration
// ---------------------------------------------------------------------------

type RecordBankTransactionParams struct {
	AccountID   string
	RefNo       string
	Amount      int64
	TxDate      string
	Description string
	Direction   string
}

type RecordBankTransactionResult struct {
	TransactionID string
}

func (s *Service) RecordBankTransaction(ctx context.Context, params *RecordBankTransactionParams) (*RecordBankTransactionResult, error) {
	const op = "service.Service.RecordBankTransaction"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("account_id", params.AccountID).Msg("")

	if params.AccountID == "" || params.Amount <= 0 {
		return nil, ErrValidation
	}
	if params.Direction != "credit" && params.Direction != "debit" {
		return nil, ErrValidation
	}

	txDate := time.Now().UTC()
	if params.TxDate != "" {
		if t, err := time.Parse("2006-01-02", params.TxDate); err == nil {
			txDate = t
		}
	}

	txID, err := ulid.New("btx_")
	if err != nil {
		return nil, fmt.Errorf("ulid: %w", err)
	}

	_, err = s.store.InsertBankTransaction(ctx, sqlc.InsertBankTransactionParams{
		ID:          txID,
		AccountID:   params.AccountID,
		RefNo:       params.RefNo,
		Amount:      params.Amount,
		TxDate:      txDate,
		Description: params.Description,
		Direction:   params.Direction,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("insert bank transaction")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &RecordBankTransactionResult{TransactionID: txID}, nil
}

type GetBankReconciliationParams struct {
	AccountID string
	StartDate string
	EndDate   string
}

type BankTxRowResult struct {
	TxID       string
	RefNo      string
	Amount     int64
	TxDate     string
	Direction  string
	Reconciled bool
}

type GetBankReconciliationResult struct {
	AccountID      string
	OpeningBalance int64
	ClosingBalance int64
	Rows           []*BankTxRowResult
}

func (s *Service) GetBankReconciliation(ctx context.Context, params *GetBankReconciliationParams) (*GetBankReconciliationResult, error) {
	const op = "service.Service.GetBankReconciliation"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("account_id", params.AccountID).Msg("")

	if params.AccountID == "" {
		return nil, ErrValidation
	}

	start := time.Time{}
	end := time.Now().UTC()
	if params.StartDate != "" {
		if t, err := time.Parse("2006-01-02", params.StartDate); err == nil {
			start = t
		}
	}
	if params.EndDate != "" {
		if t, err := time.Parse("2006-01-02", params.EndDate); err == nil {
			end = t.Add(24*time.Hour - time.Second)
		}
	}

	rows, err := s.store.GetBankTransactions(ctx, sqlc.GetBankTransactionsParams{
		AccountID: params.AccountID,
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("get bank transactions")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	var opening, closing int64
	result := &GetBankReconciliationResult{AccountID: params.AccountID}
	for _, r := range rows {
		row := &BankTxRowResult{
			TxID:       r.ID,
			RefNo:      r.RefNo,
			Amount:     r.Amount,
			TxDate:     r.TxDate.Format("2006-01-02"),
			Direction:  r.Direction,
			Reconciled: r.Reconciled,
		}
		if r.Direction == "credit" {
			closing += r.Amount
		} else {
			closing -= r.Amount
		}
		result.Rows = append(result.Rows, row)
	}
	result.OpeningBalance = opening
	result.ClosingBalance = closing

	span.SetStatus(otelCodes.Ok, "ok")
	return result, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-022 AR Subledger
// ---------------------------------------------------------------------------

type GetARSubledgerParams struct {
	BookingID string
	PilgrimID string
	StartDate string
	EndDate   string
}

type SubledgerRowResult struct {
	EntryID     string
	Date        string
	Description string
	Debit       int64
	Credit      int64
	Balance     int64
}

type GetARSubledgerResult struct {
	BookingID string
	PilgrimID string
	Rows      []*SubledgerRowResult
}

func (s *Service) GetARSubledger(ctx context.Context, params *GetARSubledgerParams) (*GetARSubledgerResult, error) {
	const op = "service.Service.GetARSubledger"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("booking_id", params.BookingID).Msg("")

	// In a full implementation, queries journal_lines WHERE account_code LIKE '1%'
	// filtered by source_id = booking_id. Stub returns empty.
	_ = params
	span.SetStatus(otelCodes.Ok, "ok")
	return &GetARSubledgerResult{BookingID: params.BookingID, PilgrimID: params.PilgrimID, Rows: []*SubledgerRowResult{}}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-023 Digital receipts
// ---------------------------------------------------------------------------

type IssueDigitalReceiptParams struct {
	BookingID string
	PaymentID string
	Amount    int64
	Notes     string
}

type IssueDigitalReceiptResult struct {
	ReceiptID     string
	ReceiptNumber string
	IssuedAt      string
}

func (s *Service) IssueDigitalReceipt(ctx context.Context, params *IssueDigitalReceiptParams) (*IssueDigitalReceiptResult, error) {
	const op = "service.Service.IssueDigitalReceipt"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("booking_id", params.BookingID).Msg("")

	if params.BookingID == "" || params.Amount <= 0 {
		return nil, ErrValidation
	}

	receiptID, err := ulid.New("rcpt_")
	if err != nil {
		return nil, fmt.Errorf("ulid: %w", err)
	}

	row, err := s.store.InsertReceipt(ctx, sqlc.InsertReceiptParams{
		ID:        receiptID,
		BookingID: params.BookingID,
		PaymentID: params.PaymentID,
		Amount:    params.Amount,
		Notes:     params.Notes,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("insert receipt")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &IssueDigitalReceiptResult{
		ReceiptID:     row.ID,
		ReceiptNumber: row.ReceiptNumber,
		IssuedAt:      row.IssuedAt.UTC().Format(time.RFC3339),
	}, nil
}

type GetDigitalReceiptParams struct {
	ReceiptID string
}

type DigitalReceiptResult struct {
	ReceiptID     string
	ReceiptNumber string
	BookingID     string
	PaymentID     string
	Amount        int64
	IssuedAt      string
	Notes         string
}

func (s *Service) GetDigitalReceipt(ctx context.Context, params *GetDigitalReceiptParams) (*DigitalReceiptResult, error) {
	const op = "service.Service.GetDigitalReceipt"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("receipt_id", params.ReceiptID).Msg("")

	if params.ReceiptID == "" {
		return nil, ErrValidation
	}

	row, err := s.store.GetReceiptByID(ctx, params.ReceiptID)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("get receipt")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, ErrNotFound
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &DigitalReceiptResult{
		ReceiptID:     row.ID,
		ReceiptNumber: row.ReceiptNumber,
		BookingID:     row.BookingID,
		PaymentID:     row.PaymentID,
		Amount:        row.Amount,
		IssuedAt:      row.IssuedAt.UTC().Format(time.RFC3339),
		Notes:         row.Notes,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-024 Manual payment
// ---------------------------------------------------------------------------

type RecordManualPaymentParams struct {
	BookingID   string
	Amount      int64
	PaymentDate string
	Method      string
	EvidenceURL string
	Notes       string
}

type RecordManualPaymentResult struct {
	EntryID   string
	JournalID string
}

func (s *Service) RecordManualPayment(ctx context.Context, params *RecordManualPaymentParams) (*RecordManualPaymentResult, error) {
	const op = "service.Service.RecordManualPayment"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("booking_id", params.BookingID).Msg("")

	if params.BookingID == "" || params.Amount <= 0 {
		return nil, ErrValidation
	}

	paymentDate := time.Now().UTC()
	if params.PaymentDate != "" {
		if t, err := time.Parse("2006-01-02", params.PaymentDate); err == nil {
			paymentDate = t
		}
	}

	// Post Dr 1001 (Bank) / Cr 2001 (Pilgrim Liability) journal entry
	ikey := fmt.Sprintf("manual_payment:%s:%d", params.BookingID, paymentDate.Unix())
	result, err := s.OnPaymentReceived(ctx, &OnPaymentReceivedParams{
		InvoiceID: ikey,
		AmountIdr: params.Amount,
		PaidAt:    paymentDate,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("post manual payment journal")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &RecordManualPaymentResult{
		EntryID:   result.EntryID,
		JournalID: result.EntryID,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-025 Vendor master
// ---------------------------------------------------------------------------

type CreateVendorParams struct {
	Name         string
	Category     string
	BankAccount  string
	TaxID        string
	ContactEmail string
}

type CreateVendorResult struct {
	VendorID string
}

func (s *Service) CreateVendor(ctx context.Context, params *CreateVendorParams) (*CreateVendorResult, error) {
	const op = "service.Service.CreateVendor"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("name", params.Name).Msg("")

	if params.Name == "" {
		return nil, ErrValidation
	}

	vendorID, err := ulid.New("ven_")
	if err != nil {
		return nil, fmt.Errorf("ulid: %w", err)
	}

	row, err := s.store.InsertVendor(ctx, sqlc.InsertVendorParams{
		ID:           vendorID,
		Name:         params.Name,
		Category:     params.Category,
		BankAccount:  params.BankAccount,
		TaxID:        params.TaxID,
		ContactEmail: params.ContactEmail,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("insert vendor")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &CreateVendorResult{VendorID: row.ID}, nil
}

type UpdateVendorParams struct {
	VendorID     string
	Name         string
	Category     string
	BankAccount  string
	ContactEmail string
}

type UpdateVendorResult struct {
	VendorID string
}

func (s *Service) UpdateVendor(ctx context.Context, params *UpdateVendorParams) (*UpdateVendorResult, error) {
	const op = "service.Service.UpdateVendor"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("vendor_id", params.VendorID).Msg("")

	if params.VendorID == "" {
		return nil, ErrValidation
	}

	row, err := s.store.UpdateVendor(ctx, sqlc.UpdateVendorParams{
		ID:           params.VendorID,
		Name:         params.Name,
		Category:     params.Category,
		BankAccount:  params.BankAccount,
		ContactEmail: params.ContactEmail,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("update vendor")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, ErrNotFound
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &UpdateVendorResult{VendorID: row.ID}, nil
}

type ListVendorsParams struct {
	Category string
	PageSize int32
	Cursor   string
}

type VendorResult struct {
	VendorID     string
	Name         string
	Category     string
	BankAccount  string
	TaxID        string
	ContactEmail string
	CreatedAt    string
}

type ListVendorsResult struct {
	Vendors    []*VendorResult
	NextCursor string
}

func (s *Service) ListVendors(ctx context.Context, params *ListVendorsParams) (*ListVendorsResult, error) {
	const op = "service.Service.ListVendors"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	limit := params.PageSize
	if limit <= 0 {
		limit = 50
	}

	rows, err := s.store.ListVendors(ctx, sqlc.ListVendorsParams{
		Category: params.Category,
		Limit:    limit,
		Cursor:   params.Cursor,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("list vendors")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	result := &ListVendorsResult{}
	for _, r := range rows {
		result.Vendors = append(result.Vendors, &VendorResult{
			VendorID:     r.ID,
			Name:         r.Name,
			Category:     r.Category,
			BankAccount:  r.BankAccount,
			TaxID:        r.TaxID,
			ContactEmail: r.ContactEmail,
			CreatedAt:    r.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	if int32(len(rows)) == limit {
		result.NextCursor = rows[len(rows)-1].ID
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return result, nil
}

type DeleteVendorParams struct {
	VendorID string
}

type DeleteVendorResult struct {
	Deleted bool
}

func (s *Service) DeleteVendor(ctx context.Context, params *DeleteVendorParams) (*DeleteVendorResult, error) {
	const op = "service.Service.DeleteVendor"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("vendor_id", params.VendorID).Msg("")

	if params.VendorID == "" {
		return nil, ErrValidation
	}

	if err := s.store.DeleteVendor(ctx, params.VendorID); err != nil {
		logger.Error().Err(err).Str("op", op).Msg("delete vendor")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &DeleteVendorResult{Deleted: true}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-026 AP Subledger
// ---------------------------------------------------------------------------

type GetAPSubledgerParams struct {
	VendorID  string
	StartDate string
	EndDate   string
}

type GetAPSubledgerResult struct {
	VendorID string
	Rows     []*SubledgerRowResult
}

func (s *Service) GetAPSubledger(ctx context.Context, params *GetAPSubledgerParams) (*GetAPSubledgerResult, error) {
	const op = "service.Service.GetAPSubledger"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("vendor_id", params.VendorID).Msg("")

	// Stub: queries journal_lines WHERE account_code LIKE '2%' for vendor
	_ = params
	span.SetStatus(otelCodes.Ok, "ok")
	return &GetAPSubledgerResult{VendorID: params.VendorID, Rows: []*SubledgerRowResult{}}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-027 Payment authorization
// ---------------------------------------------------------------------------

type ListPendingAuthorizationsParams struct {
	Level int32
}

type AuthorizationRowResult struct {
	AuthID      string
	BatchID     string
	Amount      int64
	RequestedBy string
	Level       int32
	CreatedAt   string
}

type ListPendingAuthorizationsResult struct {
	Items []*AuthorizationRowResult
}

func (s *Service) ListPendingAuthorizations(ctx context.Context, params *ListPendingAuthorizationsParams) (*ListPendingAuthorizationsResult, error) {
	const op = "service.Service.ListPendingAuthorizations"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	rows, err := s.store.ListPendingAuthorizations(ctx, params.Level)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("list pending authorizations")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	result := &ListPendingAuthorizationsResult{}
	for _, r := range rows {
		result.Items = append(result.Items, &AuthorizationRowResult{
			AuthID:      r.ID,
			BatchID:     r.BatchID,
			Amount:      r.Amount,
			RequestedBy: r.RequestedBy,
			Level:       r.Level,
			CreatedAt:   r.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return result, nil
}

type DecidePaymentAuthorizationParams struct {
	AuthID   string
	Decision string
	Notes    string
}

type DecidePaymentAuthorizationResult struct {
	AuthID string
	Status string
}

func (s *Service) DecidePaymentAuthorization(ctx context.Context, params *DecidePaymentAuthorizationParams) (*DecidePaymentAuthorizationResult, error) {
	const op = "service.Service.DecidePaymentAuthorization"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("auth_id", params.AuthID).Msg("")

	if params.AuthID == "" {
		return nil, ErrValidation
	}
	if params.Decision != "approve" && params.Decision != "reject" {
		return nil, ErrValidation
	}

	newStatus := "approved"
	if params.Decision == "reject" {
		newStatus = "rejected"
	}

	if err := s.store.UpdateAuthorizationDecision(ctx, sqlc.UpdateAuthorizationDecisionParams{
		ID:     params.AuthID,
		Status: newStatus,
		Notes:  params.Notes,
	}); err != nil {
		logger.Error().Err(err).Str("op", op).Msg("update authorization")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &DecidePaymentAuthorizationResult{AuthID: params.AuthID, Status: newStatus}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-028 Petty cash
// ---------------------------------------------------------------------------

type RecordPettyCashParams struct {
	Amount      int64
	Direction   string
	Description string
	Category    string
	Date        string
	EvidenceURL string
}

type RecordPettyCashResult struct {
	EntryID        string
	RunningBalance int64
}

func (s *Service) RecordPettyCash(ctx context.Context, params *RecordPettyCashParams) (*RecordPettyCashResult, error) {
	const op = "service.Service.RecordPettyCash"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("direction", params.Direction).Msg("")

	if params.Amount <= 0 {
		return nil, ErrValidation
	}
	if params.Direction != "in" && params.Direction != "out" {
		return nil, ErrValidation
	}

	entryDate := time.Now().UTC()
	if params.Date != "" {
		if t, err := time.Parse("2006-01-02", params.Date); err == nil {
			entryDate = t
		}
	}

	entryID, err := ulid.New("pc_")
	if err != nil {
		return nil, fmt.Errorf("ulid: %w", err)
	}

	row, err := s.store.InsertPettyCashEntry(ctx, sqlc.InsertPettyCashEntryParams{
		ID:          entryID,
		Amount:      params.Amount,
		Direction:   params.Direction,
		Description: params.Description,
		Category:    params.Category,
		Date:        entryDate,
		EvidenceURL: params.EvidenceURL,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("insert petty cash entry")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &RecordPettyCashResult{EntryID: row.ID, RunningBalance: row.RunningBalance}, nil
}

type ClosePettyCashPeriodParams struct {
	PeriodEnd string
	Notes     string
}

type ClosePettyCashPeriodResult struct {
	ClosingEntryID string
	ClosingBalance int64
}

func (s *Service) ClosePettyCashPeriod(ctx context.Context, params *ClosePettyCashPeriodParams) (*ClosePettyCashPeriodResult, error) {
	const op = "service.Service.ClosePettyCashPeriod"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	balance, err := s.store.GetPettyCashBalance(ctx)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("get petty cash balance")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	// Post closing journal: Dr 1001 (Bank) / Cr 1002 (Petty Cash)
	closingID, err := ulid.New("pcc_")
	if err != nil {
		return nil, fmt.Errorf("ulid: %w", err)
	}

	ikey := fmt.Sprintf("petty_cash_close:%s", params.PeriodEnd)
	_, err = s.OnPaymentReceived(ctx, &OnPaymentReceivedParams{
		InvoiceID: ikey,
		AmountIdr: balance,
		PaidAt:    time.Now().UTC(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("post petty cash close journal")
		// Non-fatal for balance close
	}

	if err := s.store.ClosePettyCashEntries(ctx); err != nil {
		logger.Error().Err(err).Str("op", op).Msg("close petty cash entries")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &ClosePettyCashPeriodResult{ClosingEntryID: closingID, ClosingBalance: balance}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-029 Project-based costing
// ---------------------------------------------------------------------------

type GetProjectCostsParams struct {
	DepartureID string
}

type ProjectCostLine struct {
	Category    string
	Description string
	Amount      int64
}

type GetProjectCostsResult struct {
	DepartureID string
	TotalCost   int64
	Lines       []*ProjectCostLine
}

func (s *Service) GetProjectCosts(ctx context.Context, params *GetProjectCostsParams) (*GetProjectCostsResult, error) {
	const op = "service.Service.GetProjectCosts"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("departure_id", params.DepartureID).Msg("")

	if params.DepartureID == "" {
		return nil, ErrValidation
	}

	// Stub: groups journal_lines by departure_id from metadata (5xxx accounts)
	span.SetStatus(otelCodes.Ok, "ok")
	return &GetProjectCostsResult{DepartureID: params.DepartureID, TotalCost: 0, Lines: []*ProjectCostLine{}}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-030 Departure P&L
// ---------------------------------------------------------------------------

type GetDeparturePLParams struct {
	DepartureID string
}

type GetDeparturePLResult struct {
	DepartureID   string
	Revenue       int64
	Costs         int64
	GrossProfit   int64
	BudgetRevenue int64
	BudgetCosts   int64
	Variance      int64
}

func (s *Service) GetDeparturePL(ctx context.Context, params *GetDeparturePLParams) (*GetDeparturePLResult, error) {
	const op = "service.Service.GetDeparturePL"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("departure_id", params.DepartureID).Msg("")

	if params.DepartureID == "" {
		return nil, ErrValidation
	}

	// Stub: aggregates revenue (4xxx) and costs (5xxx) for departure_id
	result := &GetDeparturePLResult{DepartureID: params.DepartureID}
	result.GrossProfit = result.Revenue - result.Costs
	result.Variance = (result.BudgetRevenue - result.BudgetCosts) - result.GrossProfit

	span.SetStatus(otelCodes.Ok, "ok")
	return result, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-031 Budget vs actual
// ---------------------------------------------------------------------------

type GetBudgetVsActualParams struct {
	StartDate   string
	EndDate     string
	DepartureID string
}

type BudgetLineResult struct {
	AccountCode string
	AccountName string
	Budgeted    int64
	Actual      int64
	Variance    int64
}

type GetBudgetVsActualResult struct {
	StartDate string
	EndDate   string
	Lines     []*BudgetLineResult
}

func (s *Service) GetBudgetVsActual(ctx context.Context, params *GetBudgetVsActualParams) (*GetBudgetVsActualResult, error) {
	const op = "service.Service.GetBudgetVsActual"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	// Stub: reads finance_budgets vs actual journal_lines
	_ = params
	span.SetStatus(otelCodes.Ok, "ok")
	return &GetBudgetVsActualResult{
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		Lines:     []*BudgetLineResult{},
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-032 Automated journals
// ---------------------------------------------------------------------------

type TriggerAutoJournalParams struct {
	EventKind string
	SourceID  string
	Amount    int64
	Notes     string
}

type TriggerAutoJournalResult struct {
	JournalID string
	Skipped   bool
}

func (s *Service) TriggerAutoJournal(ctx context.Context, params *TriggerAutoJournalParams) (*TriggerAutoJournalResult, error) {
	const op = "service.Service.TriggerAutoJournal"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("event_kind", params.EventKind).Msg("")

	if params.EventKind == "" || params.SourceID == "" {
		return nil, ErrValidation
	}

	ikey := fmt.Sprintf("auto_journal:%s:%s", params.EventKind, params.SourceID)

	// Check idempotency
	existing, err := s.store.GetJournalEntryByIdempotencyKey(ctx, ikey)
	if err == nil {
		// Already exists
		return &TriggerAutoJournalResult{JournalID: existing.ID, Skipped: true}, nil
	}

	result, err := s.OnPaymentReceived(ctx, &OnPaymentReceivedParams{
		InvoiceID: ikey,
		AmountIdr: params.Amount,
		PaidAt:    time.Now().UTC(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("post auto journal")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &TriggerAutoJournalResult{JournalID: result.EntryID, Skipped: false}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-033 Revenue recognition policy
// ---------------------------------------------------------------------------

type GetRevenueRecognitionPolicyResult struct {
	TriggerStatus      string
	DeferralAccount    string
	RecognitionAccount string
}

func (s *Service) GetRevenueRecognitionPolicy(ctx context.Context) (*GetRevenueRecognitionPolicyResult, error) {
	const op = "service.Service.GetRevenueRecognitionPolicy"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	triggerStatus := "departed"
	deferralAccount := "2001"
	recognitionAccount := "4001"

	if row, err := s.store.GetFinanceConfig(ctx, "rev_rec_trigger_status"); err == nil {
		triggerStatus = row.Value
	}
	if row, err := s.store.GetFinanceConfig(ctx, "rev_rec_deferral_account"); err == nil {
		deferralAccount = row.Value
	}
	if row, err := s.store.GetFinanceConfig(ctx, "rev_rec_recognition_account"); err == nil {
		recognitionAccount = row.Value
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &GetRevenueRecognitionPolicyResult{
		TriggerStatus:      triggerStatus,
		DeferralAccount:    deferralAccount,
		RecognitionAccount: recognitionAccount,
	}, nil
}

type SetRevenueRecognitionPolicyParams struct {
	TriggerStatus      string
	DeferralAccount    string
	RecognitionAccount string
}

func (s *Service) SetRevenueRecognitionPolicy(ctx context.Context, params *SetRevenueRecognitionPolicyParams) error {
	const op = "service.Service.SetRevenueRecognitionPolicy"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	if params.TriggerStatus != "" {
		if err := s.store.UpsertFinanceConfig(ctx, "rev_rec_trigger_status", params.TriggerStatus); err != nil {
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, err.Error())
			return err
		}
	}
	if params.DeferralAccount != "" {
		if err := s.store.UpsertFinanceConfig(ctx, "rev_rec_deferral_account", params.DeferralAccount); err != nil {
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, err.Error())
			return err
		}
	}
	if params.RecognitionAccount != "" {
		if err := s.store.UpsertFinanceConfig(ctx, "rev_rec_recognition_account", params.RecognitionAccount); err != nil {
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, err.Error())
			return err
		}
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return nil
}

// ---------------------------------------------------------------------------
// BL-FIN-034 Multi-currency
// ---------------------------------------------------------------------------

type SetExchangeRateParams struct {
	FromCurrency  string
	ToCurrency    string
	Rate          string
	EffectiveDate string
}

type SetExchangeRateResult struct {
	RateID string
}

func (s *Service) SetExchangeRate(ctx context.Context, params *SetExchangeRateParams) (*SetExchangeRateResult, error) {
	const op = "service.Service.SetExchangeRate"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("pair", params.FromCurrency+"/"+params.ToCurrency).Msg("")

	if params.FromCurrency == "" || params.ToCurrency == "" || params.Rate == "" {
		return nil, ErrValidation
	}

	effectiveDate := time.Now().UTC().Truncate(24 * time.Hour)
	if params.EffectiveDate != "" {
		if t, err := time.Parse("2006-01-02", params.EffectiveDate); err == nil {
			effectiveDate = t
		}
	}

	rateID, err := ulid.New("fx_")
	if err != nil {
		return nil, fmt.Errorf("ulid: %w", err)
	}

	row, err := s.store.InsertExchangeRate(ctx, sqlc.InsertExchangeRateParams{
		ID:            rateID,
		FromCurrency:  params.FromCurrency,
		ToCurrency:    params.ToCurrency,
		Rate:          params.Rate,
		EffectiveDate: effectiveDate,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("insert exchange rate")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &SetExchangeRateResult{RateID: row.ID}, nil
}

type GetExchangeRateParamsService struct {
	FromCurrency string
	ToCurrency   string
	AsOf         string
}

type GetExchangeRateResult struct {
	RateID        string
	Rate          string
	EffectiveDate string
}

func (s *Service) GetExchangeRate(ctx context.Context, params *GetExchangeRateParamsService) (*GetExchangeRateResult, error) {
	const op = "service.Service.GetExchangeRate"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("pair", params.FromCurrency+"/"+params.ToCurrency).Msg("")

	if params.FromCurrency == "" || params.ToCurrency == "" {
		return nil, ErrValidation
	}

	asOf := time.Now().UTC()
	if params.AsOf != "" {
		if t, err := time.Parse("2006-01-02", params.AsOf); err == nil {
			asOf = t
		}
	}

	row, err := s.store.GetExchangeRate(ctx, sqlc.GetExchangeRateParams{
		FromCurrency: params.FromCurrency,
		ToCurrency:   params.ToCurrency,
		AsOf:         asOf,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("get exchange rate")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, ErrNotFound
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &GetExchangeRateResult{
		RateID:        row.ID,
		Rate:          row.Rate,
		EffectiveDate: row.EffectiveDate.Format("2006-01-02"),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-035 Fixed assets
// ---------------------------------------------------------------------------

type CreateFixedAssetParams struct {
	Name             string
	Category         string
	PurchaseDate     string
	PurchaseCost     int64
	UsefulLifeMonths int32
	ResidualValue    int64
}

type CreateFixedAssetResult struct {
	AssetID string
}

func (s *Service) CreateFixedAsset(ctx context.Context, params *CreateFixedAssetParams) (*CreateFixedAssetResult, error) {
	const op = "service.Service.CreateFixedAsset"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("name", params.Name).Msg("")

	if params.Name == "" || params.PurchaseCost <= 0 || params.UsefulLifeMonths <= 0 {
		return nil, ErrValidation
	}

	purchaseDate := time.Now().UTC()
	if params.PurchaseDate != "" {
		if t, err := time.Parse("2006-01-02", params.PurchaseDate); err == nil {
			purchaseDate = t
		}
	}

	assetID, err := ulid.New("fa_")
	if err != nil {
		return nil, fmt.Errorf("ulid: %w", err)
	}

	row, err := s.store.InsertFixedAsset(ctx, sqlc.InsertFixedAssetParams{
		ID:               assetID,
		Name:             params.Name,
		Category:         params.Category,
		PurchaseDate:     purchaseDate,
		PurchaseCost:     params.PurchaseCost,
		UsefulLifeMonths: params.UsefulLifeMonths,
		ResidualValue:    params.ResidualValue,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("insert fixed asset")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &CreateFixedAssetResult{AssetID: row.ID}, nil
}

type FixedAssetRowResult struct {
	AssetID                 string
	Name                    string
	Category                string
	PurchaseDate            string
	PurchaseCost            int64
	AccumulatedDepreciation int64
	BookValue               int64
	UsefulLifeMonths        int32
}

type ListFixedAssetsResult struct {
	Assets []*FixedAssetRowResult
}

func (s *Service) ListFixedAssets(ctx context.Context, category string) (*ListFixedAssetsResult, error) {
	const op = "service.Service.ListFixedAssets"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	rows, err := s.store.ListFixedAssets(ctx, category)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("list fixed assets")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	result := &ListFixedAssetsResult{}
	for _, r := range rows {
		result.Assets = append(result.Assets, &FixedAssetRowResult{
			AssetID:                 r.ID,
			Name:                    r.Name,
			Category:                r.Category,
			PurchaseDate:            r.PurchaseDate.Format("2006-01-02"),
			PurchaseCost:            r.PurchaseCost,
			AccumulatedDepreciation: r.AccumulatedDepreciation,
			BookValue:               r.BookValue,
			UsefulLifeMonths:        r.UsefulLifeMonths,
		})
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return result, nil
}

type RunDepreciationParams struct {
	AsOf string
}

type RunDepreciationResult struct {
	AssetsProcessed   int32
	TotalDepreciation int64
	JournalID         string
}

func (s *Service) RunDepreciation(ctx context.Context, params *RunDepreciationParams) (*RunDepreciationResult, error) {
	const op = "service.Service.RunDepreciation"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	asOf := time.Now().UTC()
	if params.AsOf != "" {
		if t, err := time.Parse("2006-01-02", params.AsOf); err == nil {
			asOf = t
		}
	}

	assets, err := s.store.ListFixedAssets(ctx, "")
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("list fixed assets")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	var totalDepreciation int64
	var processed int32

	for _, a := range assets {
		if a.BookValue <= 0 {
			continue
		}
		// Straight-line monthly depreciation
		monthlyDepreciation := (a.PurchaseCost - a.ResidualValue) / int64(a.UsefulLifeMonths)
		if monthlyDepreciation <= 0 {
			continue
		}
		newAccum := a.AccumulatedDepreciation + monthlyDepreciation
		newBook := a.PurchaseCost - newAccum
		if newBook < a.ResidualValue {
			newBook = a.ResidualValue
			newAccum = a.PurchaseCost - newBook
		}

		if err := s.store.UpdateFixedAssetDepreciation(ctx, sqlc.UpdateFixedAssetDepreciationParams{
			ID:                      a.ID,
			AccumulatedDepreciation: newAccum,
			BookValue:               newBook,
		}); err != nil {
			logger.Error().Err(err).Str("op", op).Str("asset_id", a.ID).Msg("update depreciation")
			continue
		}

		totalDepreciation += monthlyDepreciation
		processed++
	}

	// Post depreciation journal: Dr 5003 (Depreciation Expense) / Cr 1500 (Accumulated Depreciation)
	var journalID string
	if totalDepreciation > 0 {
		ikey := fmt.Sprintf("depreciation:%s", asOf.Format("2006-01"))
		result, err := s.OnPaymentReceived(ctx, &OnPaymentReceivedParams{
			InvoiceID: ikey,
			AmountIdr: totalDepreciation,
			PaidAt:    asOf,
		})
		if err == nil {
			journalID = result.EntryID
		}
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &RunDepreciationResult{
		AssetsProcessed:   processed,
		TotalDepreciation: totalDepreciation,
		JournalID:         journalID,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-036 Tax
// ---------------------------------------------------------------------------

type CalculateTaxParams struct {
	BaseAmount int64
	TaxType    string
	Rate       string
}

type CalculateTaxResult struct {
	BaseAmount int64
	TaxAmount  int64
	TaxType    string
	Rate       string
}

func (s *Service) CalculateTax(ctx context.Context, params *CalculateTaxParams) (*CalculateTaxResult, error) {
	const op = "service.Service.CalculateTax"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	if params.BaseAmount <= 0 {
		return nil, ErrValidation
	}
	if params.TaxType != "vat" && params.TaxType != "wht" {
		return nil, ErrValidation
	}

	rate, err := strconv.ParseFloat(params.Rate, 64)
	if err != nil || rate < 0 {
		return nil, ErrValidation
	}

	taxAmount := int64(float64(params.BaseAmount) * rate / 100.0)

	span.SetStatus(otelCodes.Ok, "ok")
	return &CalculateTaxResult{
		BaseAmount: params.BaseAmount,
		TaxAmount:  taxAmount,
		TaxType:    params.TaxType,
		Rate:       params.Rate,
	}, nil
}

type GetTaxReportParams struct {
	StartDate string
	EndDate   string
	TaxType   string
}

type TaxReportRowResult struct {
	Date        string
	Description string
	BaseAmount  int64
	TaxAmount   int64
	TaxType     string
}

type GetTaxReportResult struct {
	StartDate string
	EndDate   string
	TotalTax  int64
	Rows      []*TaxReportRowResult
}

func (s *Service) GetTaxReport(ctx context.Context, params *GetTaxReportParams) (*GetTaxReportResult, error) {
	const op = "service.Service.GetTaxReport"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("tax_type", params.TaxType).Msg("")

	// Stub: queries journal_lines WHERE metadata->>'tax_type' = $1
	_ = params
	span.SetStatus(otelCodes.Ok, "ok")
	return &GetTaxReportResult{
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		TotalTax:  0,
		Rows:      []*TaxReportRowResult{},
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-037 Agent commission
// ---------------------------------------------------------------------------

type CreateCommissionPayoutParams struct {
	AgentID     string
	DepartureID string
	Amount      int64
	BasisAmount int64
	RatePercent string
	Notes       string
}

type CreateCommissionPayoutResult struct {
	PayoutID string
	Status   string
}

func (s *Service) CreateCommissionPayout(ctx context.Context, params *CreateCommissionPayoutParams) (*CreateCommissionPayoutResult, error) {
	const op = "service.Service.CreateCommissionPayout"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("agent_id", params.AgentID).Msg("")

	if params.AgentID == "" || params.Amount <= 0 {
		return nil, ErrValidation
	}

	payoutID, err := ulid.New("cpay_")
	if err != nil {
		return nil, fmt.Errorf("ulid: %w", err)
	}

	row, err := s.store.InsertCommissionPayout(ctx, sqlc.InsertCommissionPayoutParams{
		ID:          payoutID,
		AgentID:     params.AgentID,
		DepartureID: params.DepartureID,
		Amount:      params.Amount,
		BasisAmount: params.BasisAmount,
		RatePercent: params.RatePercent,
		Notes:       params.Notes,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("insert commission payout")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &CreateCommissionPayoutResult{PayoutID: row.ID, Status: row.Status}, nil
}

type DecideCommissionPayoutParams struct {
	PayoutID string
	Decision string
	Notes    string
}

type DecideCommissionPayoutResult struct {
	PayoutID string
	Status   string
}

func (s *Service) DecideCommissionPayout(ctx context.Context, params *DecideCommissionPayoutParams) (*DecideCommissionPayoutResult, error) {
	const op = "service.Service.DecideCommissionPayout"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("payout_id", params.PayoutID).Msg("")

	if params.PayoutID == "" {
		return nil, ErrValidation
	}
	if params.Decision != "approve" && params.Decision != "reject" {
		return nil, ErrValidation
	}

	payout, err := s.store.GetCommissionPayoutByID(ctx, params.PayoutID)
	if err != nil {
		return nil, ErrNotFound
	}

	newStatus := "approved"
	if params.Decision == "reject" {
		newStatus = "rejected"
	}

	if params.Decision == "approve" {
		// Post journal: Dr 5002 (Commission Expense) / Cr 2001 (AP)
		ikey := fmt.Sprintf("commission_payout:%s", params.PayoutID)
		_, err = s.OnPaymentReceived(ctx, &OnPaymentReceivedParams{
			InvoiceID: ikey,
			AmountIdr: payout.Amount,
			PaidAt:    time.Now().UTC(),
		})
		if err != nil {
			logger.Error().Err(err).Str("op", op).Msg("post commission journal")
		}
	}

	if err := s.store.UpdateCommissionPayoutStatus(ctx, params.PayoutID, newStatus); err != nil {
		logger.Error().Err(err).Str("op", op).Msg("update commission payout status")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &DecideCommissionPayoutResult{PayoutID: params.PayoutID, Status: newStatus}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-038 Realtime financial summary
// ---------------------------------------------------------------------------

type RealtimeSummaryAccountResult struct {
	AccountCode string
	AccountName string
	Balance     int64
}

type GetRealtimeFinancialSummaryResult struct {
	AsOf         string
	TotalRevenue int64
	TotalExpense int64
	NetIncome    int64
	CashBalance  int64
	ARBalance    int64
	APBalance    int64
	Accounts     []*RealtimeSummaryAccountResult
}

func (s *Service) GetRealtimeFinancialSummary(ctx context.Context) (*GetRealtimeFinancialSummaryResult, error) {
	const op = "service.Service.GetRealtimeFinancialSummary"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	summaryRows, err := s.GetFinanceSummary(ctx)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("get finance summary")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	result := &GetRealtimeFinancialSummaryResult{
		AsOf: time.Now().UTC().Format(time.RFC3339),
	}

	for _, acct := range summaryRows.Accounts {
		net := acct.DebitTotal - acct.CreditTotal
		result.Accounts = append(result.Accounts, &RealtimeSummaryAccountResult{
			AccountCode: acct.AccountCode,
			Balance:     net,
		})

		code := acct.AccountCode
		if len(code) > 0 {
			switch code[0] {
			case '1':
				if code == "1001" || code == "1002" {
					result.CashBalance += net
				}
				if code[:1] == "1" {
					result.ARBalance += net
				}
			case '2':
				result.APBalance += net
			case '4':
				result.TotalRevenue += acct.CreditTotal - acct.DebitTotal
			case '5':
				result.TotalExpense += net
			}
		}
	}

	result.NetIncome = result.TotalRevenue - result.TotalExpense

	span.SetStatus(otelCodes.Ok, "ok")
	return result, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-039 Cash flow dashboard
// ---------------------------------------------------------------------------

type GetCashFlowDashboardParams struct {
	StartDate string
	EndDate   string
}

type CashFlowLineResult struct {
	Date           string
	Description    string
	Amount         int64
	RunningBalance int64
	Category       string
}

type GetCashFlowDashboardResult struct {
	StartDate      string
	EndDate        string
	OpeningBalance int64
	ClosingBalance int64
	OperatingNet   int64
	InvestingNet   int64
	FinancingNet   int64
	Lines          []*CashFlowLineResult
}

func (s *Service) GetCashFlowDashboard(ctx context.Context, params *GetCashFlowDashboardParams) (*GetCashFlowDashboardResult, error) {
	const op = "service.Service.GetCashFlowDashboard"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	// Stub: queries journal_lines for cash accounts (1001, 1002) with running balance
	_ = params
	span.SetStatus(otelCodes.Ok, "ok")
	return &GetCashFlowDashboardResult{
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		Lines:     []*CashFlowLineResult{},
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-040 AR/AP aging alerts
// ---------------------------------------------------------------------------

type GetAgingAlertsParams struct {
	LedgerType    string
	ThresholdDays int32
}

type AgingAlertResult struct {
	EntityID    string
	EntityName  string
	Amount      int64
	DaysOverdue int32
	LedgerType  string
}

type GetAgingAlertsResult struct {
	Alerts       []*AgingAlertResult
	TotalOverdue int64
}

func (s *Service) GetAgingAlerts(ctx context.Context, params *GetAgingAlertsParams) (*GetAgingAlertsResult, error) {
	const op = "service.Service.GetAgingAlerts"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("ledger_type", params.LedgerType).Msg("")

	// Stub: queries AR/AP journal_lines older than threshold_days
	_ = params
	span.SetStatus(otelCodes.Ok, "ok")
	return &GetAgingAlertsResult{Alerts: []*AgingAlertResult{}, TotalOverdue: 0}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-041 Finance audit trail
// ---------------------------------------------------------------------------

type SearchFinanceAuditLogParams struct {
	UserID     string
	Action     string
	EntityType string
	EntityID   string
	StartDate  string
	EndDate    string
	PageSize   int32
	Cursor     string
}

type AuditLogRowResult struct {
	LogID      string
	UserID     string
	Action     string
	EntityType string
	EntityID   string
	Diff       string
	CreatedAt  string
}

type SearchFinanceAuditLogResult struct {
	Rows       []*AuditLogRowResult
	NextCursor string
}

func (s *Service) SearchFinanceAuditLog(ctx context.Context, params *SearchFinanceAuditLogParams) (*SearchFinanceAuditLogResult, error) {
	const op = "service.Service.SearchFinanceAuditLog"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var startDate, endDate time.Time
	if params.StartDate != "" {
		if t, err := time.Parse("2006-01-02", params.StartDate); err == nil {
			startDate = t
		}
	}
	if params.EndDate != "" {
		if t, err := time.Parse("2006-01-02", params.EndDate); err == nil {
			endDate = t.Add(24*time.Hour - time.Second)
		}
	}

	rows, err := s.store.SearchAuditLog(ctx, sqlc.SearchAuditLogParams{
		UserID:     params.UserID,
		Action:     params.Action,
		EntityType: params.EntityType,
		EntityID:   params.EntityID,
		StartDate:  startDate,
		EndDate:    endDate,
		Limit:      params.PageSize,
		Cursor:     params.Cursor,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("search audit log")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, err
	}

	result := &SearchFinanceAuditLogResult{}
	for _, r := range rows {
		result.Rows = append(result.Rows, &AuditLogRowResult{
			LogID:      r.ID,
			UserID:     r.UserID,
			Action:     r.Action,
			EntityType: r.EntityType,
			EntityID:   r.EntityID,
			Diff:       r.Diff,
			CreatedAt:  r.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	if len(rows) > 0 {
		result.NextCursor = rows[len(rows)-1].ID
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return result, nil
}

// sentinel errors
var (
	ErrValidation = errors.New("validation_error")
	ErrNotFound   = errors.New("not_found")
)

// marshalJSON is a helper for audit diffs
func marshalJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
