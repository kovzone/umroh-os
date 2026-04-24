// finance_depth_wave4.go — gateway-svc adapter methods for Wave 4 finance depth
// RPCs (BL-FIN-020..041).
//
// All methods follow the pattern established in disbursement.go:
//   - Start OTel span
//   - Call wave4Client.Xxx(ctx, &pb.XxxRequest{...})
//   - Map errors via mapFinanceError
//   - Return domain result

package finance_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/finance_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// BL-FIN-020 ScheduleBilling
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

func (a *Adapter) ScheduleBilling(ctx context.Context, params *ScheduleBillingParams) (*ScheduleBillingResult, error) {
	const op = "finance_grpc_adapter.Adapter.ScheduleBilling"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.ScheduleBilling(ctx, &pb.ScheduleBillingRequest{
		DepartureID: params.DepartureID,
		DueDate:     params.DueDate,
		Notes:       params.Notes,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.ScheduleBilling failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &ScheduleBillingResult{
		InvoicesCreated: resp.GetInvoicesCreated(),
		TotalAmount:     resp.GetTotalAmount(),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-021 RecordBankTransaction / GetBankReconciliation
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

func (a *Adapter) RecordBankTransaction(ctx context.Context, params *RecordBankTransactionParams) (*RecordBankTransactionResult, error) {
	const op = "finance_grpc_adapter.Adapter.RecordBankTransaction"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.RecordBankTransaction(ctx, &pb.RecordBankTransactionRequest{
		AccountID:   params.AccountID,
		RefNo:       params.RefNo,
		Amount:      params.Amount,
		TxDate:      params.TxDate,
		Description: params.Description,
		Direction:   params.Direction,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.RecordBankTransaction failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RecordBankTransactionResult{TransactionID: resp.GetTransactionID()}, nil
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

func (a *Adapter) GetBankReconciliation(ctx context.Context, params *GetBankReconciliationParams) (*GetBankReconciliationResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetBankReconciliation"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.GetBankReconciliation(ctx, &pb.GetBankReconciliationRequest{
		AccountID: params.AccountID,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetBankReconciliation failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	rows := make([]*BankTxRowResult, 0, len(resp.GetRows()))
	for _, r := range resp.GetRows() {
		rows = append(rows, &BankTxRowResult{
			TxID:       r.GetTxID(),
			RefNo:      r.GetRefNo(),
			Amount:     r.GetAmount(),
			TxDate:     r.GetTxDate(),
			Direction:  r.GetDirection(),
			Reconciled: r.GetReconciled(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetBankReconciliationResult{
		AccountID:      resp.GetAccountID(),
		OpeningBalance: resp.GetOpeningBalance(),
		ClosingBalance: resp.GetClosingBalance(),
		Rows:           rows,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-022 GetARSubledger
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

func (a *Adapter) GetARSubledger(ctx context.Context, params *GetARSubledgerParams) (*GetARSubledgerResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetARSubledger"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.GetARSubledger(ctx, &pb.GetARSubledgerRequest{
		BookingID: params.BookingID,
		PilgrimID: params.PilgrimID,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetARSubledger failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	rows := make([]*SubledgerRowResult, 0, len(resp.GetRows()))
	for _, r := range resp.GetRows() {
		rows = append(rows, &SubledgerRowResult{
			EntryID:     r.GetEntryID(),
			Date:        r.GetDate(),
			Description: r.GetDescription(),
			Debit:       r.GetDebit(),
			Credit:      r.GetCredit(),
			Balance:     r.GetBalance(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetARSubledgerResult{
		BookingID: resp.GetBookingID(),
		PilgrimID: resp.GetPilgrimID(),
		Rows:      rows,
	}, nil
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

func (a *Adapter) IssueDigitalReceipt(ctx context.Context, params *IssueDigitalReceiptParams) (*IssueDigitalReceiptResult, error) {
	const op = "finance_grpc_adapter.Adapter.IssueDigitalReceipt"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.IssueDigitalReceipt(ctx, &pb.IssueDigitalReceiptRequest{
		BookingID: params.BookingID,
		PaymentID: params.PaymentID,
		Amount:    params.Amount,
		Notes:     params.Notes,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.IssueDigitalReceipt failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &IssueDigitalReceiptResult{
		ReceiptID:     resp.GetReceiptID(),
		ReceiptNumber: resp.GetReceiptNumber(),
		IssuedAt:      resp.GetIssuedAt(),
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

func (a *Adapter) GetDigitalReceipt(ctx context.Context, params *GetDigitalReceiptParams) (*DigitalReceiptResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetDigitalReceipt"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.GetDigitalReceipt(ctx, &pb.GetDigitalReceiptRequest{
		ReceiptID: params.ReceiptID,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetDigitalReceipt failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	r := resp.GetReceipt()
	if r == nil {
		span.SetStatus(codes.Ok, "ok")
		return &DigitalReceiptResult{}, nil
	}
	span.SetStatus(codes.Ok, "ok")
	return &DigitalReceiptResult{
		ReceiptID:     r.GetReceiptID(),
		ReceiptNumber: r.GetReceiptNumber(),
		BookingID:     r.GetBookingID(),
		PaymentID:     r.GetPaymentID(),
		Amount:        r.GetAmount(),
		IssuedAt:      r.GetIssuedAt(),
		Notes:         r.GetNotes(),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-024 RecordManualPayment
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

func (a *Adapter) RecordManualPayment(ctx context.Context, params *RecordManualPaymentParams) (*RecordManualPaymentResult, error) {
	const op = "finance_grpc_adapter.Adapter.RecordManualPayment"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.RecordManualPayment(ctx, &pb.RecordManualPaymentRequest{
		BookingID:   params.BookingID,
		Amount:      params.Amount,
		PaymentDate: params.PaymentDate,
		Method:      params.Method,
		EvidenceURL: params.EvidenceURL,
		Notes:       params.Notes,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.RecordManualPayment failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RecordManualPaymentResult{
		EntryID:   resp.GetEntryID(),
		JournalID: resp.GetJournalID(),
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

func (a *Adapter) CreateVendor(ctx context.Context, params *CreateVendorParams) (*CreateVendorResult, error) {
	const op = "finance_grpc_adapter.Adapter.CreateVendor"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.CreateVendor(ctx, &pb.CreateVendorRequest{
		Name:         params.Name,
		Category:     params.Category,
		BankAccount:  params.BankAccount,
		TaxID:        params.TaxID,
		ContactEmail: params.ContactEmail,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.CreateVendor failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &CreateVendorResult{VendorID: resp.GetVendorID()}, nil
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

func (a *Adapter) UpdateVendor(ctx context.Context, params *UpdateVendorParams) (*UpdateVendorResult, error) {
	const op = "finance_grpc_adapter.Adapter.UpdateVendor"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.UpdateVendor(ctx, &pb.UpdateVendorRequest{
		VendorID:     params.VendorID,
		Name:         params.Name,
		Category:     params.Category,
		BankAccount:  params.BankAccount,
		ContactEmail: params.ContactEmail,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.UpdateVendor failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &UpdateVendorResult{VendorID: resp.GetVendorID()}, nil
}

type ListVendorsParams struct {
	Category string
	PageSize int32
	Cursor   string
}

type VendorRowResult struct {
	VendorID     string
	Name         string
	Category     string
	BankAccount  string
	TaxID        string
	ContactEmail string
	CreatedAt    string
}

type ListVendorsResult struct {
	Vendors    []*VendorRowResult
	NextCursor string
}

func (a *Adapter) ListVendors(ctx context.Context, params *ListVendorsParams) (*ListVendorsResult, error) {
	const op = "finance_grpc_adapter.Adapter.ListVendors"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.ListVendors(ctx, &pb.ListVendorsRequest{
		Category: params.Category,
		PageSize: params.PageSize,
		Cursor:   params.Cursor,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.ListVendors failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	vendors := make([]*VendorRowResult, 0, len(resp.GetVendors()))
	for _, v := range resp.GetVendors() {
		vendors = append(vendors, &VendorRowResult{
			VendorID:     v.GetVendorID(),
			Name:         v.GetName(),
			Category:     v.GetCategory(),
			BankAccount:  v.GetBankAccount(),
			TaxID:        v.GetTaxID(),
			ContactEmail: v.GetContactEmail(),
			CreatedAt:    v.GetCreatedAt(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &ListVendorsResult{Vendors: vendors, NextCursor: resp.GetNextCursor()}, nil
}

type DeleteVendorParams struct {
	VendorID string
}

type DeleteVendorResult struct {
	Deleted bool
}

func (a *Adapter) DeleteVendor(ctx context.Context, params *DeleteVendorParams) (*DeleteVendorResult, error) {
	const op = "finance_grpc_adapter.Adapter.DeleteVendor"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.DeleteVendor(ctx, &pb.DeleteVendorRequest{VendorID: params.VendorID})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.DeleteVendor failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &DeleteVendorResult{Deleted: resp.GetDeleted()}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-026 GetAPSubledger
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

func (a *Adapter) GetAPSubledger(ctx context.Context, params *GetAPSubledgerParams) (*GetAPSubledgerResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetAPSubledger"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.GetAPSubledger(ctx, &pb.GetAPSubledgerRequest{
		VendorID:  params.VendorID,
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetAPSubledger failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	rows := make([]*SubledgerRowResult, 0, len(resp.GetRows()))
	for _, r := range resp.GetRows() {
		rows = append(rows, &SubledgerRowResult{
			EntryID:     r.GetEntryID(),
			Date:        r.GetDate(),
			Description: r.GetDescription(),
			Debit:       r.GetDebit(),
			Credit:      r.GetCredit(),
			Balance:     r.GetBalance(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetAPSubledgerResult{VendorID: resp.GetVendorID(), Rows: rows}, nil
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

func (a *Adapter) ListPendingAuthorizations(ctx context.Context, params *ListPendingAuthorizationsParams) (*ListPendingAuthorizationsResult, error) {
	const op = "finance_grpc_adapter.Adapter.ListPendingAuthorizations"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.ListPendingAuthorizations(ctx, &pb.ListPendingAuthorizationsRequest{Level: params.Level})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.ListPendingAuthorizations failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	items := make([]*AuthorizationRowResult, 0, len(resp.GetItems()))
	for _, it := range resp.GetItems() {
		items = append(items, &AuthorizationRowResult{
			AuthID:      it.GetAuthID(),
			BatchID:     it.GetBatchID(),
			Amount:      it.GetAmount(),
			RequestedBy: it.GetRequestedBy(),
			Level:       it.GetLevel(),
			CreatedAt:   it.GetCreatedAt(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &ListPendingAuthorizationsResult{Items: items}, nil
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

func (a *Adapter) DecidePaymentAuthorization(ctx context.Context, params *DecidePaymentAuthorizationParams) (*DecidePaymentAuthorizationResult, error) {
	const op = "finance_grpc_adapter.Adapter.DecidePaymentAuthorization"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.DecidePaymentAuthorization(ctx, &pb.DecidePaymentAuthorizationRequest{
		AuthID:   params.AuthID,
		Decision: params.Decision,
		Notes:    params.Notes,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.DecidePaymentAuthorization failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &DecidePaymentAuthorizationResult{AuthID: resp.GetAuthID(), Status: resp.GetStatus()}, nil
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

func (a *Adapter) RecordPettyCash(ctx context.Context, params *RecordPettyCashParams) (*RecordPettyCashResult, error) {
	const op = "finance_grpc_adapter.Adapter.RecordPettyCash"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.RecordPettyCash(ctx, &pb.RecordPettyCashRequest{
		Amount:      params.Amount,
		Direction:   params.Direction,
		Description: params.Description,
		Category:    params.Category,
		Date:        params.Date,
		EvidenceURL: params.EvidenceURL,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.RecordPettyCash failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RecordPettyCashResult{EntryID: resp.GetEntryID(), RunningBalance: resp.GetRunningBalance()}, nil
}

type ClosePettyCashPeriodParams struct {
	PeriodEnd string
	Notes     string
}

type ClosePettyCashPeriodResult struct {
	ClosingEntryID string
	ClosingBalance int64
}

func (a *Adapter) ClosePettyCashPeriod(ctx context.Context, params *ClosePettyCashPeriodParams) (*ClosePettyCashPeriodResult, error) {
	const op = "finance_grpc_adapter.Adapter.ClosePettyCashPeriod"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.ClosePettyCashPeriod(ctx, &pb.ClosePettyCashPeriodRequest{
		PeriodEnd: params.PeriodEnd,
		Notes:     params.Notes,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.ClosePettyCashPeriod failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &ClosePettyCashPeriodResult{
		ClosingEntryID: resp.GetClosingEntryID(),
		ClosingBalance: resp.GetClosingBalance(),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-029 GetProjectCosts
// ---------------------------------------------------------------------------

type GetProjectCostsParams struct {
	DepartureID string
}

type CostLineItemResult struct {
	Category    string
	Description string
	Amount      int64
}

type GetProjectCostsResult struct {
	DepartureID string
	TotalCost   int64
	Lines       []*CostLineItemResult
}

func (a *Adapter) GetProjectCosts(ctx context.Context, params *GetProjectCostsParams) (*GetProjectCostsResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetProjectCosts"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.GetProjectCosts(ctx, &pb.GetProjectCostsRequest{DepartureID: params.DepartureID})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetProjectCosts failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	lines := make([]*CostLineItemResult, 0, len(resp.GetLines()))
	for _, l := range resp.GetLines() {
		lines = append(lines, &CostLineItemResult{
			Category:    l.GetCategory(),
			Description: l.GetDescription(),
			Amount:      l.GetAmount(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetProjectCostsResult{
		DepartureID: resp.GetDepartureID(),
		TotalCost:   resp.GetTotalCost(),
		Lines:       lines,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-030 GetDeparturePL
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

func (a *Adapter) GetDeparturePL(ctx context.Context, params *GetDeparturePLParams) (*GetDeparturePLResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetDeparturePL"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.GetDeparturePL(ctx, &pb.GetDeparturePLRequest{DepartureID: params.DepartureID})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetDeparturePL failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetDeparturePLResult{
		DepartureID:   resp.GetDepartureID(),
		Revenue:       resp.GetRevenue(),
		Costs:         resp.GetCosts(),
		GrossProfit:   resp.GetGrossProfit(),
		BudgetRevenue: resp.GetBudgetRevenue(),
		BudgetCosts:   resp.GetBudgetCosts(),
		Variance:      resp.GetVariance(),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-031 GetBudgetVsActual
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

func (a *Adapter) GetBudgetVsActual(ctx context.Context, params *GetBudgetVsActualParams) (*GetBudgetVsActualResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetBudgetVsActual"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.GetBudgetVsActual(ctx, &pb.GetBudgetVsActualRequest{
		StartDate:   params.StartDate,
		EndDate:     params.EndDate,
		DepartureID: params.DepartureID,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetBudgetVsActual failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	lines := make([]*BudgetLineResult, 0, len(resp.GetLines()))
	for _, l := range resp.GetLines() {
		lines = append(lines, &BudgetLineResult{
			AccountCode: l.GetAccountCode(),
			AccountName: l.GetAccountName(),
			Budgeted:    l.GetBudgeted(),
			Actual:      l.GetActual(),
			Variance:    l.GetVariance(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetBudgetVsActualResult{StartDate: resp.GetStartDate(), EndDate: resp.GetEndDate(), Lines: lines}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-032 TriggerAutoJournal
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

func (a *Adapter) TriggerAutoJournal(ctx context.Context, params *TriggerAutoJournalParams) (*TriggerAutoJournalResult, error) {
	const op = "finance_grpc_adapter.Adapter.TriggerAutoJournal"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.TriggerAutoJournal(ctx, &pb.TriggerAutoJournalRequest{
		EventKind: params.EventKind,
		SourceID:  params.SourceID,
		Amount:    params.Amount,
		Notes:     params.Notes,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.TriggerAutoJournal failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &TriggerAutoJournalResult{JournalID: resp.GetJournalID(), Skipped: resp.GetSkipped()}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-033 Revenue recognition policy
// ---------------------------------------------------------------------------

type RevenueRecognitionPolicyResult struct {
	TriggerStatus      string
	DeferralAccount    string
	RecognitionAccount string
}

func (a *Adapter) GetRevenueRecognitionPolicy(ctx context.Context) (*RevenueRecognitionPolicyResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetRevenueRecognitionPolicy"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.GetRevenueRecognitionPolicy(ctx, &pb.GetRevenueRecognitionPolicyRequest{})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetRevenueRecognitionPolicy failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	p := resp.GetPolicy()
	if p == nil {
		span.SetStatus(codes.Ok, "ok")
		return &RevenueRecognitionPolicyResult{}, nil
	}
	span.SetStatus(codes.Ok, "ok")
	return &RevenueRecognitionPolicyResult{
		TriggerStatus:      p.GetTriggerStatus(),
		DeferralAccount:    p.GetDeferralAccount(),
		RecognitionAccount: p.GetRecognitionAccount(),
	}, nil
}

type SetRevenueRecognitionPolicyParams struct {
	TriggerStatus      string
	DeferralAccount    string
	RecognitionAccount string
}

func (a *Adapter) SetRevenueRecognitionPolicy(ctx context.Context, params *SetRevenueRecognitionPolicyParams) error {
	const op = "finance_grpc_adapter.Adapter.SetRevenueRecognitionPolicy"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	_, err := a.wave4Client.SetRevenueRecognitionPolicy(ctx, &pb.SetRevenueRecognitionPolicyRequest{
		TriggerStatus:      params.TriggerStatus,
		DeferralAccount:    params.DeferralAccount,
		RecognitionAccount: params.RecognitionAccount,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.SetRevenueRecognitionPolicy failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return nil
}

// ---------------------------------------------------------------------------
// BL-FIN-034 Exchange rates
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

func (a *Adapter) SetExchangeRate(ctx context.Context, params *SetExchangeRateParams) (*SetExchangeRateResult, error) {
	const op = "finance_grpc_adapter.Adapter.SetExchangeRate"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.SetExchangeRate(ctx, &pb.SetExchangeRateRequest{
		FromCurrency:  params.FromCurrency,
		ToCurrency:    params.ToCurrency,
		Rate:          params.Rate,
		EffectiveDate: params.EffectiveDate,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.SetExchangeRate failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &SetExchangeRateResult{RateID: resp.GetRateID()}, nil
}

type GetExchangeRateParams struct {
	FromCurrency string
	ToCurrency   string
	AsOf         string
}

type GetExchangeRateResult struct {
	RateID        string
	Rate          string
	EffectiveDate string
}

func (a *Adapter) GetExchangeRate(ctx context.Context, params *GetExchangeRateParams) (*GetExchangeRateResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetExchangeRate"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.GetExchangeRate(ctx, &pb.GetExchangeRateRequest{
		FromCurrency: params.FromCurrency,
		ToCurrency:   params.ToCurrency,
		AsOf:         params.AsOf,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetExchangeRate failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetExchangeRateResult{
		RateID:        resp.GetRateID(),
		Rate:          resp.GetRate(),
		EffectiveDate: resp.GetEffectiveDate(),
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

func (a *Adapter) CreateFixedAsset(ctx context.Context, params *CreateFixedAssetParams) (*CreateFixedAssetResult, error) {
	const op = "finance_grpc_adapter.Adapter.CreateFixedAsset"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.CreateFixedAsset(ctx, &pb.CreateFixedAssetRequest{
		Name:             params.Name,
		Category:         params.Category,
		PurchaseDate:     params.PurchaseDate,
		PurchaseCost:     params.PurchaseCost,
		UsefulLifeMonths: params.UsefulLifeMonths,
		ResidualValue:    params.ResidualValue,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.CreateFixedAsset failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &CreateFixedAssetResult{AssetID: resp.GetAssetID()}, nil
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

func (a *Adapter) ListFixedAssets(ctx context.Context, category string) (*ListFixedAssetsResult, error) {
	const op = "finance_grpc_adapter.Adapter.ListFixedAssets"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.ListFixedAssets(ctx, &pb.ListFixedAssetsRequest{Category: category})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.ListFixedAssets failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	assets := make([]*FixedAssetRowResult, 0, len(resp.GetAssets()))
	for _, a2 := range resp.GetAssets() {
		assets = append(assets, &FixedAssetRowResult{
			AssetID:                 a2.GetAssetID(),
			Name:                    a2.GetName(),
			Category:                a2.GetCategory(),
			PurchaseDate:            a2.GetPurchaseDate(),
			PurchaseCost:            a2.GetPurchaseCost(),
			AccumulatedDepreciation: a2.GetAccumulatedDepreciation(),
			BookValue:               a2.GetBookValue(),
			UsefulLifeMonths:        a2.GetUsefulLifeMonths(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &ListFixedAssetsResult{Assets: assets}, nil
}

type RunDepreciationParams struct {
	AsOf string
}

type RunDepreciationResult struct {
	AssetsProcessed   int32
	TotalDepreciation int64
	JournalID         string
}

func (a *Adapter) RunDepreciation(ctx context.Context, params *RunDepreciationParams) (*RunDepreciationResult, error) {
	const op = "finance_grpc_adapter.Adapter.RunDepreciation"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.RunDepreciation(ctx, &pb.RunDepreciationRequest{AsOf: params.AsOf})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.RunDepreciation failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RunDepreciationResult{
		AssetsProcessed:   resp.GetAssetsProcessed(),
		TotalDepreciation: resp.GetTotalDepreciation(),
		JournalID:         resp.GetJournalID(),
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

func (a *Adapter) CalculateTax(ctx context.Context, params *CalculateTaxParams) (*CalculateTaxResult, error) {
	const op = "finance_grpc_adapter.Adapter.CalculateTax"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.CalculateTax(ctx, &pb.CalculateTaxRequest{
		BaseAmount: params.BaseAmount,
		TaxType:    params.TaxType,
		Rate:       params.Rate,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.CalculateTax failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &CalculateTaxResult{
		BaseAmount: resp.GetBaseAmount(),
		TaxAmount:  resp.GetTaxAmount(),
		TaxType:    resp.GetTaxType(),
		Rate:       resp.GetRate(),
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

func (a *Adapter) GetTaxReport(ctx context.Context, params *GetTaxReportParams) (*GetTaxReportResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetTaxReport"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.GetTaxReport(ctx, &pb.GetTaxReportRequest{
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
		TaxType:   params.TaxType,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetTaxReport failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	rows := make([]*TaxReportRowResult, 0, len(resp.GetRows()))
	for _, r := range resp.GetRows() {
		rows = append(rows, &TaxReportRowResult{
			Date:        r.GetDate(),
			Description: r.GetDescription(),
			BaseAmount:  r.GetBaseAmount(),
			TaxAmount:   r.GetTaxAmount(),
			TaxType:     r.GetTaxType(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetTaxReportResult{
		StartDate: resp.GetStartDate(),
		EndDate:   resp.GetEndDate(),
		TotalTax:  resp.GetTotalTax(),
		Rows:      rows,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-037 Commission payouts
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

func (a *Adapter) CreateCommissionPayout(ctx context.Context, params *CreateCommissionPayoutParams) (*CreateCommissionPayoutResult, error) {
	const op = "finance_grpc_adapter.Adapter.CreateCommissionPayout"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.CreateCommissionPayout(ctx, &pb.CreateCommissionPayoutRequest{
		AgentID:     params.AgentID,
		DepartureID: params.DepartureID,
		Amount:      params.Amount,
		BasisAmount: params.BasisAmount,
		RatePercent: params.RatePercent,
		Notes:       params.Notes,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.CreateCommissionPayout failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &CreateCommissionPayoutResult{PayoutID: resp.GetPayoutID(), Status: resp.GetStatus()}, nil
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

func (a *Adapter) DecideCommissionPayout(ctx context.Context, params *DecideCommissionPayoutParams) (*DecideCommissionPayoutResult, error) {
	const op = "finance_grpc_adapter.Adapter.DecideCommissionPayout"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.DecideCommissionPayout(ctx, &pb.DecideCommissionPayoutRequest{
		PayoutID: params.PayoutID,
		Decision: params.Decision,
		Notes:    params.Notes,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.DecideCommissionPayout failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &DecideCommissionPayoutResult{PayoutID: resp.GetPayoutID(), Status: resp.GetStatus()}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-038 GetRealtimeFinancialSummary
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

func (a *Adapter) GetRealtimeFinancialSummary(ctx context.Context) (*GetRealtimeFinancialSummaryResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetRealtimeFinancialSummary"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.GetRealtimeFinancialSummary(ctx, &pb.GetRealtimeFinancialSummaryRequest{})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetRealtimeFinancialSummary failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	accounts := make([]*RealtimeSummaryAccountResult, 0, len(resp.GetAccounts()))
	for _, acc := range resp.GetAccounts() {
		accounts = append(accounts, &RealtimeSummaryAccountResult{
			AccountCode: acc.GetAccountCode(),
			AccountName: acc.GetAccountName(),
			Balance:     acc.GetBalance(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetRealtimeFinancialSummaryResult{
		AsOf:         resp.GetAsOf(),
		TotalRevenue: resp.GetTotalRevenue(),
		TotalExpense: resp.GetTotalExpense(),
		NetIncome:    resp.GetNetIncome(),
		CashBalance:  resp.GetCashBalance(),
		ARBalance:    resp.GetARBalance(),
		APBalance:    resp.GetAPBalance(),
		Accounts:     accounts,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-039 GetCashFlowDashboard
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

func (a *Adapter) GetCashFlowDashboard(ctx context.Context, params *GetCashFlowDashboardParams) (*GetCashFlowDashboardResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetCashFlowDashboard"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.GetCashFlowDashboard(ctx, &pb.GetCashFlowDashboardRequest{
		StartDate: params.StartDate,
		EndDate:   params.EndDate,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetCashFlowDashboard failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	lines := make([]*CashFlowLineResult, 0, len(resp.GetLines()))
	for _, l := range resp.GetLines() {
		lines = append(lines, &CashFlowLineResult{
			Date:           l.GetDate(),
			Description:    l.GetDescription(),
			Amount:         l.GetAmount(),
			RunningBalance: l.GetRunningBalance(),
			Category:       l.GetCategory(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetCashFlowDashboardResult{
		StartDate:      resp.GetStartDate(),
		EndDate:        resp.GetEndDate(),
		OpeningBalance: resp.GetOpeningBalance(),
		ClosingBalance: resp.GetClosingBalance(),
		OperatingNet:   resp.GetOperatingNet(),
		InvestingNet:   resp.GetInvestingNet(),
		FinancingNet:   resp.GetFinancingNet(),
		Lines:          lines,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-040 GetAgingAlerts
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

func (a *Adapter) GetAgingAlerts(ctx context.Context, params *GetAgingAlertsParams) (*GetAgingAlertsResult, error) {
	const op = "finance_grpc_adapter.Adapter.GetAgingAlerts"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.GetAgingAlerts(ctx, &pb.GetAgingAlertsRequest{
		LedgerType:    params.LedgerType,
		ThresholdDays: params.ThresholdDays,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.GetAgingAlerts failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	alerts := make([]*AgingAlertResult, 0, len(resp.GetAlerts()))
	for _, al := range resp.GetAlerts() {
		alerts = append(alerts, &AgingAlertResult{
			EntityID:    al.GetEntityID(),
			EntityName:  al.GetEntityName(),
			Amount:      al.GetAmount(),
			DaysOverdue: al.GetDaysOverdue(),
			LedgerType:  al.GetLedgerType(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetAgingAlertsResult{Alerts: alerts, TotalOverdue: resp.GetTotalOverdue()}, nil
}

// ---------------------------------------------------------------------------
// BL-FIN-041 SearchFinanceAuditLog
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

func (a *Adapter) SearchFinanceAuditLog(ctx context.Context, params *SearchFinanceAuditLogParams) (*SearchFinanceAuditLogResult, error) {
	const op = "finance_grpc_adapter.Adapter.SearchFinanceAuditLog"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, a.logger)
	resp, err := a.wave4Client.SearchFinanceAuditLog(ctx, &pb.SearchFinanceAuditLogRequest{
		UserID:     params.UserID,
		Action:     params.Action,
		EntityType: params.EntityType,
		EntityID:   params.EntityID,
		StartDate:  params.StartDate,
		EndDate:    params.EndDate,
		PageSize:   params.PageSize,
		Cursor:     params.Cursor,
	})
	if err != nil {
		wrapped := mapFinanceError(err)
		logger.Warn().Err(wrapped).Msg("finance-svc.SearchFinanceAuditLog failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	rows := make([]*AuditLogRowResult, 0, len(resp.GetRows()))
	for _, r := range resp.GetRows() {
		rows = append(rows, &AuditLogRowResult{
			LogID:      r.GetLogID(),
			UserID:     r.GetUserID(),
			Action:     r.GetAction(),
			EntityType: r.GetEntityType(),
			EntityID:   r.GetEntityID(),
			Diff:       r.GetDiff(),
			CreatedAt:  r.GetCreatedAt(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &SearchFinanceAuditLogResult{Rows: rows, NextCursor: resp.GetNextCursor()}, nil
}
