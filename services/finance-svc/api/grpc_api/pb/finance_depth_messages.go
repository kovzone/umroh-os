// finance_depth_messages.go — hand-written proto message types for Wave 4
// Finance depth RPCs (BL-FIN-020..041).

package pb

// ---------------------------------------------------------------------------
// BL-FIN-020 ScheduleBilling
// ---------------------------------------------------------------------------

type ScheduleBillingRequest struct {
	DepartureID string
	DueDate     string
	Notes       string
}

func (x *ScheduleBillingRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *ScheduleBillingRequest) GetDueDate() string {
	if x == nil {
		return ""
	}
	return x.DueDate
}
func (x *ScheduleBillingRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type ScheduleBillingResponse struct {
	InvoicesCreated int32
	TotalAmount     int64
}

func (x *ScheduleBillingResponse) GetInvoicesCreated() int32 {
	if x == nil {
		return 0
	}
	return x.InvoicesCreated
}
func (x *ScheduleBillingResponse) GetTotalAmount() int64 {
	if x == nil {
		return 0
	}
	return x.TotalAmount
}

// ---------------------------------------------------------------------------
// BL-FIN-021 Bank integration
// ---------------------------------------------------------------------------

type RecordBankTransactionRequest struct {
	AccountID   string
	RefNo       string
	Amount      int64
	TxDate      string
	Description string
	Direction   string // "credit" | "debit"
}

func (x *RecordBankTransactionRequest) GetAccountID() string {
	if x == nil {
		return ""
	}
	return x.AccountID
}
func (x *RecordBankTransactionRequest) GetRefNo() string {
	if x == nil {
		return ""
	}
	return x.RefNo
}
func (x *RecordBankTransactionRequest) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *RecordBankTransactionRequest) GetTxDate() string {
	if x == nil {
		return ""
	}
	return x.TxDate
}
func (x *RecordBankTransactionRequest) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *RecordBankTransactionRequest) GetDirection() string {
	if x == nil {
		return ""
	}
	return x.Direction
}

type RecordBankTransactionResponse struct {
	TransactionID string
}

func (x *RecordBankTransactionResponse) GetTransactionID() string {
	if x == nil {
		return ""
	}
	return x.TransactionID
}

type GetBankReconciliationRequest struct {
	AccountID string
	StartDate string
	EndDate   string
}

func (x *GetBankReconciliationRequest) GetAccountID() string {
	if x == nil {
		return ""
	}
	return x.AccountID
}
func (x *GetBankReconciliationRequest) GetStartDate() string {
	if x == nil {
		return ""
	}
	return x.StartDate
}
func (x *GetBankReconciliationRequest) GetEndDate() string {
	if x == nil {
		return ""
	}
	return x.EndDate
}

type BankTxRow struct {
	TxID        string
	RefNo       string
	Amount      int64
	TxDate      string
	Direction   string
	Reconciled  bool
}

func (x *BankTxRow) GetTxID() string {
	if x == nil {
		return ""
	}
	return x.TxID
}
func (x *BankTxRow) GetRefNo() string {
	if x == nil {
		return ""
	}
	return x.RefNo
}
func (x *BankTxRow) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *BankTxRow) GetTxDate() string {
	if x == nil {
		return ""
	}
	return x.TxDate
}
func (x *BankTxRow) GetDirection() string {
	if x == nil {
		return ""
	}
	return x.Direction
}
func (x *BankTxRow) GetReconciled() bool {
	if x == nil {
		return false
	}
	return x.Reconciled
}

type GetBankReconciliationResponse struct {
	AccountID      string
	OpeningBalance int64
	ClosingBalance int64
	Rows           []*BankTxRow
}

func (x *GetBankReconciliationResponse) GetAccountID() string {
	if x == nil {
		return ""
	}
	return x.AccountID
}
func (x *GetBankReconciliationResponse) GetOpeningBalance() int64 {
	if x == nil {
		return 0
	}
	return x.OpeningBalance
}
func (x *GetBankReconciliationResponse) GetClosingBalance() int64 {
	if x == nil {
		return 0
	}
	return x.ClosingBalance
}
func (x *GetBankReconciliationResponse) GetRows() []*BankTxRow {
	if x == nil {
		return nil
	}
	return x.Rows
}

// ---------------------------------------------------------------------------
// BL-FIN-022 AR Subledger
// ---------------------------------------------------------------------------

type GetARSubledgerRequest struct {
	BookingID string
	PilgrimID string
	StartDate string
	EndDate   string
}

func (x *GetARSubledgerRequest) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *GetARSubledgerRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *GetARSubledgerRequest) GetStartDate() string {
	if x == nil {
		return ""
	}
	return x.StartDate
}
func (x *GetARSubledgerRequest) GetEndDate() string {
	if x == nil {
		return ""
	}
	return x.EndDate
}

type ARSubledgerRow struct {
	EntryID     string
	Date        string
	Description string
	Debit       int64
	Credit      int64
	Balance     int64
}

func (x *ARSubledgerRow) GetEntryID() string {
	if x == nil {
		return ""
	}
	return x.EntryID
}
func (x *ARSubledgerRow) GetDate() string {
	if x == nil {
		return ""
	}
	return x.Date
}
func (x *ARSubledgerRow) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *ARSubledgerRow) GetDebit() int64 {
	if x == nil {
		return 0
	}
	return x.Debit
}
func (x *ARSubledgerRow) GetCredit() int64 {
	if x == nil {
		return 0
	}
	return x.Credit
}
func (x *ARSubledgerRow) GetBalance() int64 {
	if x == nil {
		return 0
	}
	return x.Balance
}

type GetARSubledgerResponse struct {
	BookingID string
	PilgrimID string
	Rows      []*ARSubledgerRow
}

func (x *GetARSubledgerResponse) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *GetARSubledgerResponse) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *GetARSubledgerResponse) GetRows() []*ARSubledgerRow {
	if x == nil {
		return nil
	}
	return x.Rows
}

// ---------------------------------------------------------------------------
// BL-FIN-023 Digital receipts
// ---------------------------------------------------------------------------

type IssueDigitalReceiptRequest struct {
	BookingID string
	PaymentID string
	Amount    int64
	Notes     string
}

func (x *IssueDigitalReceiptRequest) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *IssueDigitalReceiptRequest) GetPaymentID() string {
	if x == nil {
		return ""
	}
	return x.PaymentID
}
func (x *IssueDigitalReceiptRequest) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *IssueDigitalReceiptRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type IssueDigitalReceiptResponse struct {
	ReceiptID     string
	ReceiptNumber string
	IssuedAt      string
}

func (x *IssueDigitalReceiptResponse) GetReceiptID() string {
	if x == nil {
		return ""
	}
	return x.ReceiptID
}
func (x *IssueDigitalReceiptResponse) GetReceiptNumber() string {
	if x == nil {
		return ""
	}
	return x.ReceiptNumber
}
func (x *IssueDigitalReceiptResponse) GetIssuedAt() string {
	if x == nil {
		return ""
	}
	return x.IssuedAt
}

type GetDigitalReceiptRequest struct {
	ReceiptID string
}

func (x *GetDigitalReceiptRequest) GetReceiptID() string {
	if x == nil {
		return ""
	}
	return x.ReceiptID
}

type DigitalReceiptProto struct {
	ReceiptID     string
	ReceiptNumber string
	BookingID     string
	PaymentID     string
	Amount        int64
	IssuedAt      string
	Notes         string
}

func (x *DigitalReceiptProto) GetReceiptID() string {
	if x == nil {
		return ""
	}
	return x.ReceiptID
}
func (x *DigitalReceiptProto) GetReceiptNumber() string {
	if x == nil {
		return ""
	}
	return x.ReceiptNumber
}
func (x *DigitalReceiptProto) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *DigitalReceiptProto) GetPaymentID() string {
	if x == nil {
		return ""
	}
	return x.PaymentID
}
func (x *DigitalReceiptProto) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *DigitalReceiptProto) GetIssuedAt() string {
	if x == nil {
		return ""
	}
	return x.IssuedAt
}
func (x *DigitalReceiptProto) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type GetDigitalReceiptResponse struct {
	Receipt *DigitalReceiptProto
}

func (x *GetDigitalReceiptResponse) GetReceipt() *DigitalReceiptProto {
	if x == nil {
		return nil
	}
	return x.Receipt
}

// ---------------------------------------------------------------------------
// BL-FIN-024 Manual payment
// ---------------------------------------------------------------------------

type RecordManualPaymentRequest struct {
	BookingID   string
	Amount      int64
	PaymentDate string
	Method      string
	EvidenceURL string
	Notes       string
}

func (x *RecordManualPaymentRequest) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *RecordManualPaymentRequest) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *RecordManualPaymentRequest) GetPaymentDate() string {
	if x == nil {
		return ""
	}
	return x.PaymentDate
}
func (x *RecordManualPaymentRequest) GetMethod() string {
	if x == nil {
		return ""
	}
	return x.Method
}
func (x *RecordManualPaymentRequest) GetEvidenceURL() string {
	if x == nil {
		return ""
	}
	return x.EvidenceURL
}
func (x *RecordManualPaymentRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type RecordManualPaymentResponse struct {
	EntryID   string
	JournalID string
}

func (x *RecordManualPaymentResponse) GetEntryID() string {
	if x == nil {
		return ""
	}
	return x.EntryID
}
func (x *RecordManualPaymentResponse) GetJournalID() string {
	if x == nil {
		return ""
	}
	return x.JournalID
}

// ---------------------------------------------------------------------------
// BL-FIN-025 Vendor master
// ---------------------------------------------------------------------------

type CreateVendorRequest struct {
	Name         string
	Category     string
	BankAccount  string
	TaxID        string
	ContactEmail string
}

func (x *CreateVendorRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *CreateVendorRequest) GetCategory() string {
	if x == nil {
		return ""
	}
	return x.Category
}
func (x *CreateVendorRequest) GetBankAccount() string {
	if x == nil {
		return ""
	}
	return x.BankAccount
}
func (x *CreateVendorRequest) GetTaxID() string {
	if x == nil {
		return ""
	}
	return x.TaxID
}
func (x *CreateVendorRequest) GetContactEmail() string {
	if x == nil {
		return ""
	}
	return x.ContactEmail
}

type CreateVendorResponse struct {
	VendorID string
}

func (x *CreateVendorResponse) GetVendorID() string {
	if x == nil {
		return ""
	}
	return x.VendorID
}

type UpdateVendorRequest struct {
	VendorID     string
	Name         string
	Category     string
	BankAccount  string
	ContactEmail string
}

func (x *UpdateVendorRequest) GetVendorID() string {
	if x == nil {
		return ""
	}
	return x.VendorID
}
func (x *UpdateVendorRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *UpdateVendorRequest) GetCategory() string {
	if x == nil {
		return ""
	}
	return x.Category
}
func (x *UpdateVendorRequest) GetBankAccount() string {
	if x == nil {
		return ""
	}
	return x.BankAccount
}
func (x *UpdateVendorRequest) GetContactEmail() string {
	if x == nil {
		return ""
	}
	return x.ContactEmail
}

type UpdateVendorResponse struct {
	VendorID string
}

func (x *UpdateVendorResponse) GetVendorID() string {
	if x == nil {
		return ""
	}
	return x.VendorID
}

type ListVendorsRequest struct {
	Category string
	PageSize int32
	Cursor   string
}

func (x *ListVendorsRequest) GetCategory() string {
	if x == nil {
		return ""
	}
	return x.Category
}
func (x *ListVendorsRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}
func (x *ListVendorsRequest) GetCursor() string {
	if x == nil {
		return ""
	}
	return x.Cursor
}

type VendorRow struct {
	VendorID     string
	Name         string
	Category     string
	BankAccount  string
	TaxID        string
	ContactEmail string
	CreatedAt    string
}

func (x *VendorRow) GetVendorID() string {
	if x == nil {
		return ""
	}
	return x.VendorID
}
func (x *VendorRow) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *VendorRow) GetCategory() string {
	if x == nil {
		return ""
	}
	return x.Category
}
func (x *VendorRow) GetBankAccount() string {
	if x == nil {
		return ""
	}
	return x.BankAccount
}
func (x *VendorRow) GetTaxID() string {
	if x == nil {
		return ""
	}
	return x.TaxID
}
func (x *VendorRow) GetContactEmail() string {
	if x == nil {
		return ""
	}
	return x.ContactEmail
}
func (x *VendorRow) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type ListVendorsResponse struct {
	Vendors    []*VendorRow
	NextCursor string
}

func (x *ListVendorsResponse) GetVendors() []*VendorRow {
	if x == nil {
		return nil
	}
	return x.Vendors
}
func (x *ListVendorsResponse) GetNextCursor() string {
	if x == nil {
		return ""
	}
	return x.NextCursor
}

type DeleteVendorRequest struct {
	VendorID string
}

func (x *DeleteVendorRequest) GetVendorID() string {
	if x == nil {
		return ""
	}
	return x.VendorID
}

type DeleteVendorResponse struct {
	Deleted bool
}

func (x *DeleteVendorResponse) GetDeleted() bool {
	if x == nil {
		return false
	}
	return x.Deleted
}

// ---------------------------------------------------------------------------
// BL-FIN-026 AP Subledger
// ---------------------------------------------------------------------------

type GetAPSubledgerRequest struct {
	VendorID  string
	StartDate string
	EndDate   string
}

func (x *GetAPSubledgerRequest) GetVendorID() string {
	if x == nil {
		return ""
	}
	return x.VendorID
}
func (x *GetAPSubledgerRequest) GetStartDate() string {
	if x == nil {
		return ""
	}
	return x.StartDate
}
func (x *GetAPSubledgerRequest) GetEndDate() string {
	if x == nil {
		return ""
	}
	return x.EndDate
}

type APSubledgerRow struct {
	EntryID     string
	Date        string
	Description string
	Debit       int64
	Credit      int64
	Balance     int64
}

func (x *APSubledgerRow) GetEntryID() string {
	if x == nil {
		return ""
	}
	return x.EntryID
}
func (x *APSubledgerRow) GetDate() string {
	if x == nil {
		return ""
	}
	return x.Date
}
func (x *APSubledgerRow) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *APSubledgerRow) GetDebit() int64 {
	if x == nil {
		return 0
	}
	return x.Debit
}
func (x *APSubledgerRow) GetCredit() int64 {
	if x == nil {
		return 0
	}
	return x.Credit
}
func (x *APSubledgerRow) GetBalance() int64 {
	if x == nil {
		return 0
	}
	return x.Balance
}

type GetAPSubledgerResponse struct {
	VendorID string
	Rows     []*APSubledgerRow
}

func (x *GetAPSubledgerResponse) GetVendorID() string {
	if x == nil {
		return ""
	}
	return x.VendorID
}
func (x *GetAPSubledgerResponse) GetRows() []*APSubledgerRow {
	if x == nil {
		return nil
	}
	return x.Rows
}

// ---------------------------------------------------------------------------
// BL-FIN-027 Payment authorization
// ---------------------------------------------------------------------------

type ListPendingAuthorizationsRequest struct {
	Level int32
}

func (x *ListPendingAuthorizationsRequest) GetLevel() int32 {
	if x == nil {
		return 0
	}
	return x.Level
}

type AuthorizationRow struct {
	AuthID      string
	BatchID     string
	Amount      int64
	RequestedBy string
	Level       int32
	CreatedAt   string
}

func (x *AuthorizationRow) GetAuthID() string {
	if x == nil {
		return ""
	}
	return x.AuthID
}
func (x *AuthorizationRow) GetBatchID() string {
	if x == nil {
		return ""
	}
	return x.BatchID
}
func (x *AuthorizationRow) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *AuthorizationRow) GetRequestedBy() string {
	if x == nil {
		return ""
	}
	return x.RequestedBy
}
func (x *AuthorizationRow) GetLevel() int32 {
	if x == nil {
		return 0
	}
	return x.Level
}
func (x *AuthorizationRow) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type ListPendingAuthorizationsResponse struct {
	Items []*AuthorizationRow
}

func (x *ListPendingAuthorizationsResponse) GetItems() []*AuthorizationRow {
	if x == nil {
		return nil
	}
	return x.Items
}

type DecidePaymentAuthorizationRequest struct {
	AuthID   string
	Decision string // "approve" | "reject"
	Notes    string
}

func (x *DecidePaymentAuthorizationRequest) GetAuthID() string {
	if x == nil {
		return ""
	}
	return x.AuthID
}
func (x *DecidePaymentAuthorizationRequest) GetDecision() string {
	if x == nil {
		return ""
	}
	return x.Decision
}
func (x *DecidePaymentAuthorizationRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type DecidePaymentAuthorizationResponse struct {
	AuthID string
	Status string
}

func (x *DecidePaymentAuthorizationResponse) GetAuthID() string {
	if x == nil {
		return ""
	}
	return x.AuthID
}
func (x *DecidePaymentAuthorizationResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

// ---------------------------------------------------------------------------
// BL-FIN-028 Petty cash
// ---------------------------------------------------------------------------

type RecordPettyCashRequest struct {
	Amount      int64
	Direction   string // "in" | "out"
	Description string
	Category    string
	Date        string
	EvidenceURL string
}

func (x *RecordPettyCashRequest) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *RecordPettyCashRequest) GetDirection() string {
	if x == nil {
		return ""
	}
	return x.Direction
}
func (x *RecordPettyCashRequest) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *RecordPettyCashRequest) GetCategory() string {
	if x == nil {
		return ""
	}
	return x.Category
}
func (x *RecordPettyCashRequest) GetDate() string {
	if x == nil {
		return ""
	}
	return x.Date
}
func (x *RecordPettyCashRequest) GetEvidenceURL() string {
	if x == nil {
		return ""
	}
	return x.EvidenceURL
}

type RecordPettyCashResponse struct {
	EntryID        string
	RunningBalance int64
}

func (x *RecordPettyCashResponse) GetEntryID() string {
	if x == nil {
		return ""
	}
	return x.EntryID
}
func (x *RecordPettyCashResponse) GetRunningBalance() int64 {
	if x == nil {
		return 0
	}
	return x.RunningBalance
}

type ClosePettyCashPeriodRequest struct {
	PeriodEnd string
	Notes     string
}

func (x *ClosePettyCashPeriodRequest) GetPeriodEnd() string {
	if x == nil {
		return ""
	}
	return x.PeriodEnd
}
func (x *ClosePettyCashPeriodRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type ClosePettyCashPeriodResponse struct {
	ClosingEntryID string
	ClosingBalance int64
}

func (x *ClosePettyCashPeriodResponse) GetClosingEntryID() string {
	if x == nil {
		return ""
	}
	return x.ClosingEntryID
}
func (x *ClosePettyCashPeriodResponse) GetClosingBalance() int64 {
	if x == nil {
		return 0
	}
	return x.ClosingBalance
}

// ---------------------------------------------------------------------------
// BL-FIN-029 Project-based costing
// ---------------------------------------------------------------------------

type GetProjectCostsRequest struct {
	DepartureID string
}

func (x *GetProjectCostsRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type CostLineItem struct {
	Category    string
	Description string
	Amount      int64
}

func (x *CostLineItem) GetCategory() string {
	if x == nil {
		return ""
	}
	return x.Category
}
func (x *CostLineItem) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *CostLineItem) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}

type GetProjectCostsResponse struct {
	DepartureID string
	TotalCost   int64
	Lines       []*CostLineItem
}

func (x *GetProjectCostsResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *GetProjectCostsResponse) GetTotalCost() int64 {
	if x == nil {
		return 0
	}
	return x.TotalCost
}
func (x *GetProjectCostsResponse) GetLines() []*CostLineItem {
	if x == nil {
		return nil
	}
	return x.Lines
}

// ---------------------------------------------------------------------------
// BL-FIN-030 Departure P&L
// ---------------------------------------------------------------------------

type GetDeparturePLRequest struct {
	DepartureID string
}

func (x *GetDeparturePLRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type GetDeparturePLResponse struct {
	DepartureID   string
	Revenue       int64
	Costs         int64
	GrossProfit   int64
	BudgetRevenue int64
	BudgetCosts   int64
	Variance      int64
}

func (x *GetDeparturePLResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *GetDeparturePLResponse) GetRevenue() int64 {
	if x == nil {
		return 0
	}
	return x.Revenue
}
func (x *GetDeparturePLResponse) GetCosts() int64 {
	if x == nil {
		return 0
	}
	return x.Costs
}
func (x *GetDeparturePLResponse) GetGrossProfit() int64 {
	if x == nil {
		return 0
	}
	return x.GrossProfit
}
func (x *GetDeparturePLResponse) GetBudgetRevenue() int64 {
	if x == nil {
		return 0
	}
	return x.BudgetRevenue
}
func (x *GetDeparturePLResponse) GetBudgetCosts() int64 {
	if x == nil {
		return 0
	}
	return x.BudgetCosts
}
func (x *GetDeparturePLResponse) GetVariance() int64 {
	if x == nil {
		return 0
	}
	return x.Variance
}

// ---------------------------------------------------------------------------
// BL-FIN-031 Budget vs actual
// ---------------------------------------------------------------------------

type GetBudgetVsActualRequest struct {
	StartDate   string
	EndDate     string
	DepartureID string
}

func (x *GetBudgetVsActualRequest) GetStartDate() string {
	if x == nil {
		return ""
	}
	return x.StartDate
}
func (x *GetBudgetVsActualRequest) GetEndDate() string {
	if x == nil {
		return ""
	}
	return x.EndDate
}
func (x *GetBudgetVsActualRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type BudgetLine struct {
	AccountCode string
	AccountName string
	Budgeted    int64
	Actual      int64
	Variance    int64
}

func (x *BudgetLine) GetAccountCode() string {
	if x == nil {
		return ""
	}
	return x.AccountCode
}
func (x *BudgetLine) GetAccountName() string {
	if x == nil {
		return ""
	}
	return x.AccountName
}
func (x *BudgetLine) GetBudgeted() int64 {
	if x == nil {
		return 0
	}
	return x.Budgeted
}
func (x *BudgetLine) GetActual() int64 {
	if x == nil {
		return 0
	}
	return x.Actual
}
func (x *BudgetLine) GetVariance() int64 {
	if x == nil {
		return 0
	}
	return x.Variance
}

type GetBudgetVsActualResponse struct {
	StartDate string
	EndDate   string
	Lines     []*BudgetLine
}

func (x *GetBudgetVsActualResponse) GetStartDate() string {
	if x == nil {
		return ""
	}
	return x.StartDate
}
func (x *GetBudgetVsActualResponse) GetEndDate() string {
	if x == nil {
		return ""
	}
	return x.EndDate
}
func (x *GetBudgetVsActualResponse) GetLines() []*BudgetLine {
	if x == nil {
		return nil
	}
	return x.Lines
}

// ---------------------------------------------------------------------------
// BL-FIN-032 Automated journals
// ---------------------------------------------------------------------------

type TriggerAutoJournalRequest struct {
	EventKind string
	SourceID  string
	Amount    int64
	Notes     string
}

func (x *TriggerAutoJournalRequest) GetEventKind() string {
	if x == nil {
		return ""
	}
	return x.EventKind
}
func (x *TriggerAutoJournalRequest) GetSourceID() string {
	if x == nil {
		return ""
	}
	return x.SourceID
}
func (x *TriggerAutoJournalRequest) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *TriggerAutoJournalRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type TriggerAutoJournalResponse struct {
	JournalID string
	Skipped   bool
}

func (x *TriggerAutoJournalResponse) GetJournalID() string {
	if x == nil {
		return ""
	}
	return x.JournalID
}
func (x *TriggerAutoJournalResponse) GetSkipped() bool {
	if x == nil {
		return false
	}
	return x.Skipped
}

// ---------------------------------------------------------------------------
// BL-FIN-033 Revenue recognition policy
// ---------------------------------------------------------------------------

type GetRevenueRecognitionPolicyRequest struct{}

type RevenueRecognitionPolicy struct {
	TriggerStatus      string
	DeferralAccount    string
	RecognitionAccount string
}

func (x *RevenueRecognitionPolicy) GetTriggerStatus() string {
	if x == nil {
		return ""
	}
	return x.TriggerStatus
}
func (x *RevenueRecognitionPolicy) GetDeferralAccount() string {
	if x == nil {
		return ""
	}
	return x.DeferralAccount
}
func (x *RevenueRecognitionPolicy) GetRecognitionAccount() string {
	if x == nil {
		return ""
	}
	return x.RecognitionAccount
}

type GetRevenueRecognitionPolicyResponse struct {
	Policy *RevenueRecognitionPolicy
}

func (x *GetRevenueRecognitionPolicyResponse) GetPolicy() *RevenueRecognitionPolicy {
	if x == nil {
		return nil
	}
	return x.Policy
}

type SetRevenueRecognitionPolicyRequest struct {
	TriggerStatus      string
	DeferralAccount    string
	RecognitionAccount string
}

func (x *SetRevenueRecognitionPolicyRequest) GetTriggerStatus() string {
	if x == nil {
		return ""
	}
	return x.TriggerStatus
}
func (x *SetRevenueRecognitionPolicyRequest) GetDeferralAccount() string {
	if x == nil {
		return ""
	}
	return x.DeferralAccount
}
func (x *SetRevenueRecognitionPolicyRequest) GetRecognitionAccount() string {
	if x == nil {
		return ""
	}
	return x.RecognitionAccount
}

type SetRevenueRecognitionPolicyResponse struct {
	Updated bool
}

func (x *SetRevenueRecognitionPolicyResponse) GetUpdated() bool {
	if x == nil {
		return false
	}
	return x.Updated
}

// ---------------------------------------------------------------------------
// BL-FIN-034 Multi-currency
// ---------------------------------------------------------------------------

type SetExchangeRateRequest struct {
	FromCurrency  string
	ToCurrency    string
	Rate          string // string to avoid float precision
	EffectiveDate string
}

func (x *SetExchangeRateRequest) GetFromCurrency() string {
	if x == nil {
		return ""
	}
	return x.FromCurrency
}
func (x *SetExchangeRateRequest) GetToCurrency() string {
	if x == nil {
		return ""
	}
	return x.ToCurrency
}
func (x *SetExchangeRateRequest) GetRate() string {
	if x == nil {
		return ""
	}
	return x.Rate
}
func (x *SetExchangeRateRequest) GetEffectiveDate() string {
	if x == nil {
		return ""
	}
	return x.EffectiveDate
}

type SetExchangeRateResponse struct {
	RateID string
}

func (x *SetExchangeRateResponse) GetRateID() string {
	if x == nil {
		return ""
	}
	return x.RateID
}

type GetExchangeRateRequest struct {
	FromCurrency string
	ToCurrency   string
	AsOf         string
}

func (x *GetExchangeRateRequest) GetFromCurrency() string {
	if x == nil {
		return ""
	}
	return x.FromCurrency
}
func (x *GetExchangeRateRequest) GetToCurrency() string {
	if x == nil {
		return ""
	}
	return x.ToCurrency
}
func (x *GetExchangeRateRequest) GetAsOf() string {
	if x == nil {
		return ""
	}
	return x.AsOf
}

type GetExchangeRateResponse struct {
	RateID        string
	Rate          string
	EffectiveDate string
}

func (x *GetExchangeRateResponse) GetRateID() string {
	if x == nil {
		return ""
	}
	return x.RateID
}
func (x *GetExchangeRateResponse) GetRate() string {
	if x == nil {
		return ""
	}
	return x.Rate
}
func (x *GetExchangeRateResponse) GetEffectiveDate() string {
	if x == nil {
		return ""
	}
	return x.EffectiveDate
}

// ---------------------------------------------------------------------------
// BL-FIN-035 Fixed assets
// ---------------------------------------------------------------------------

type CreateFixedAssetRequest struct {
	Name             string
	Category         string
	PurchaseDate     string
	PurchaseCost     int64
	UsefulLifeMonths int32
	ResidualValue    int64
}

func (x *CreateFixedAssetRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *CreateFixedAssetRequest) GetCategory() string {
	if x == nil {
		return ""
	}
	return x.Category
}
func (x *CreateFixedAssetRequest) GetPurchaseDate() string {
	if x == nil {
		return ""
	}
	return x.PurchaseDate
}
func (x *CreateFixedAssetRequest) GetPurchaseCost() int64 {
	if x == nil {
		return 0
	}
	return x.PurchaseCost
}
func (x *CreateFixedAssetRequest) GetUsefulLifeMonths() int32 {
	if x == nil {
		return 0
	}
	return x.UsefulLifeMonths
}
func (x *CreateFixedAssetRequest) GetResidualValue() int64 {
	if x == nil {
		return 0
	}
	return x.ResidualValue
}

type CreateFixedAssetResponse struct {
	AssetID string
}

func (x *CreateFixedAssetResponse) GetAssetID() string {
	if x == nil {
		return ""
	}
	return x.AssetID
}

type FixedAssetRow struct {
	AssetID                 string
	Name                    string
	Category                string
	PurchaseDate            string
	PurchaseCost            int64
	AccumulatedDepreciation int64
	BookValue               int64
	UsefulLifeMonths        int32
}

func (x *FixedAssetRow) GetAssetID() string {
	if x == nil {
		return ""
	}
	return x.AssetID
}
func (x *FixedAssetRow) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *FixedAssetRow) GetCategory() string {
	if x == nil {
		return ""
	}
	return x.Category
}
func (x *FixedAssetRow) GetPurchaseDate() string {
	if x == nil {
		return ""
	}
	return x.PurchaseDate
}
func (x *FixedAssetRow) GetPurchaseCost() int64 {
	if x == nil {
		return 0
	}
	return x.PurchaseCost
}
func (x *FixedAssetRow) GetAccumulatedDepreciation() int64 {
	if x == nil {
		return 0
	}
	return x.AccumulatedDepreciation
}
func (x *FixedAssetRow) GetBookValue() int64 {
	if x == nil {
		return 0
	}
	return x.BookValue
}
func (x *FixedAssetRow) GetUsefulLifeMonths() int32 {
	if x == nil {
		return 0
	}
	return x.UsefulLifeMonths
}

type ListFixedAssetsRequest struct {
	Category string
}

func (x *ListFixedAssetsRequest) GetCategory() string {
	if x == nil {
		return ""
	}
	return x.Category
}

type ListFixedAssetsResponse struct {
	Assets []*FixedAssetRow
}

func (x *ListFixedAssetsResponse) GetAssets() []*FixedAssetRow {
	if x == nil {
		return nil
	}
	return x.Assets
}

type RunDepreciationRequest struct {
	AsOf string
}

func (x *RunDepreciationRequest) GetAsOf() string {
	if x == nil {
		return ""
	}
	return x.AsOf
}

type RunDepreciationResponse struct {
	AssetsProcessed  int32
	TotalDepreciation int64
	JournalID        string
}

func (x *RunDepreciationResponse) GetAssetsProcessed() int32 {
	if x == nil {
		return 0
	}
	return x.AssetsProcessed
}
func (x *RunDepreciationResponse) GetTotalDepreciation() int64 {
	if x == nil {
		return 0
	}
	return x.TotalDepreciation
}
func (x *RunDepreciationResponse) GetJournalID() string {
	if x == nil {
		return ""
	}
	return x.JournalID
}

// ---------------------------------------------------------------------------
// BL-FIN-036 Tax integration
// ---------------------------------------------------------------------------

type CalculateTaxRequest struct {
	BaseAmount int64
	TaxType    string // "vat" | "wht"
	Rate       string // string to avoid float precision
}

func (x *CalculateTaxRequest) GetBaseAmount() int64 {
	if x == nil {
		return 0
	}
	return x.BaseAmount
}
func (x *CalculateTaxRequest) GetTaxType() string {
	if x == nil {
		return ""
	}
	return x.TaxType
}
func (x *CalculateTaxRequest) GetRate() string {
	if x == nil {
		return ""
	}
	return x.Rate
}

type CalculateTaxResponse struct {
	BaseAmount int64
	TaxAmount  int64
	TaxType    string
	Rate       string
}

func (x *CalculateTaxResponse) GetBaseAmount() int64 {
	if x == nil {
		return 0
	}
	return x.BaseAmount
}
func (x *CalculateTaxResponse) GetTaxAmount() int64 {
	if x == nil {
		return 0
	}
	return x.TaxAmount
}
func (x *CalculateTaxResponse) GetTaxType() string {
	if x == nil {
		return ""
	}
	return x.TaxType
}
func (x *CalculateTaxResponse) GetRate() string {
	if x == nil {
		return ""
	}
	return x.Rate
}

type GetTaxReportRequest struct {
	StartDate string
	EndDate   string
	TaxType   string
}

func (x *GetTaxReportRequest) GetStartDate() string {
	if x == nil {
		return ""
	}
	return x.StartDate
}
func (x *GetTaxReportRequest) GetEndDate() string {
	if x == nil {
		return ""
	}
	return x.EndDate
}
func (x *GetTaxReportRequest) GetTaxType() string {
	if x == nil {
		return ""
	}
	return x.TaxType
}

type TaxReportRow struct {
	Date        string
	Description string
	BaseAmount  int64
	TaxAmount   int64
	TaxType     string
}

func (x *TaxReportRow) GetDate() string {
	if x == nil {
		return ""
	}
	return x.Date
}
func (x *TaxReportRow) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *TaxReportRow) GetBaseAmount() int64 {
	if x == nil {
		return 0
	}
	return x.BaseAmount
}
func (x *TaxReportRow) GetTaxAmount() int64 {
	if x == nil {
		return 0
	}
	return x.TaxAmount
}
func (x *TaxReportRow) GetTaxType() string {
	if x == nil {
		return ""
	}
	return x.TaxType
}

type GetTaxReportResponse struct {
	StartDate string
	EndDate   string
	TotalTax  int64
	Rows      []*TaxReportRow
}

func (x *GetTaxReportResponse) GetStartDate() string {
	if x == nil {
		return ""
	}
	return x.StartDate
}
func (x *GetTaxReportResponse) GetEndDate() string {
	if x == nil {
		return ""
	}
	return x.EndDate
}
func (x *GetTaxReportResponse) GetTotalTax() int64 {
	if x == nil {
		return 0
	}
	return x.TotalTax
}
func (x *GetTaxReportResponse) GetRows() []*TaxReportRow {
	if x == nil {
		return nil
	}
	return x.Rows
}

// ---------------------------------------------------------------------------
// BL-FIN-037 Agent commission
// ---------------------------------------------------------------------------

type CreateCommissionPayoutRequest struct {
	AgentID      string
	DepartureID  string
	Amount       int64
	BasisAmount  int64
	RatePercent  string
	Notes        string
}

func (x *CreateCommissionPayoutRequest) GetAgentID() string {
	if x == nil {
		return ""
	}
	return x.AgentID
}
func (x *CreateCommissionPayoutRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *CreateCommissionPayoutRequest) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *CreateCommissionPayoutRequest) GetBasisAmount() int64 {
	if x == nil {
		return 0
	}
	return x.BasisAmount
}
func (x *CreateCommissionPayoutRequest) GetRatePercent() string {
	if x == nil {
		return ""
	}
	return x.RatePercent
}
func (x *CreateCommissionPayoutRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type CreateCommissionPayoutResponse struct {
	PayoutID string
	Status   string
}

func (x *CreateCommissionPayoutResponse) GetPayoutID() string {
	if x == nil {
		return ""
	}
	return x.PayoutID
}
func (x *CreateCommissionPayoutResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

type DecideCommissionPayoutRequest struct {
	PayoutID string
	Decision string // "approve" | "reject"
	Notes    string
}

func (x *DecideCommissionPayoutRequest) GetPayoutID() string {
	if x == nil {
		return ""
	}
	return x.PayoutID
}
func (x *DecideCommissionPayoutRequest) GetDecision() string {
	if x == nil {
		return ""
	}
	return x.Decision
}
func (x *DecideCommissionPayoutRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type DecideCommissionPayoutResponse struct {
	PayoutID string
	Status   string
}

func (x *DecideCommissionPayoutResponse) GetPayoutID() string {
	if x == nil {
		return ""
	}
	return x.PayoutID
}
func (x *DecideCommissionPayoutResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

// ---------------------------------------------------------------------------
// BL-FIN-038 Real-time summary
// ---------------------------------------------------------------------------

type GetRealtimeFinancialSummaryRequest struct{}

type RealtimeSummaryAccount struct {
	AccountCode string
	AccountName string
	Balance     int64
}

func (x *RealtimeSummaryAccount) GetAccountCode() string {
	if x == nil {
		return ""
	}
	return x.AccountCode
}
func (x *RealtimeSummaryAccount) GetAccountName() string {
	if x == nil {
		return ""
	}
	return x.AccountName
}
func (x *RealtimeSummaryAccount) GetBalance() int64 {
	if x == nil {
		return 0
	}
	return x.Balance
}

type GetRealtimeFinancialSummaryResponse struct {
	AsOf         string
	TotalRevenue int64
	TotalExpense int64
	NetIncome    int64
	CashBalance  int64
	ARBalance    int64
	APBalance    int64
	Accounts     []*RealtimeSummaryAccount
}

func (x *GetRealtimeFinancialSummaryResponse) GetAsOf() string {
	if x == nil {
		return ""
	}
	return x.AsOf
}
func (x *GetRealtimeFinancialSummaryResponse) GetTotalRevenue() int64 {
	if x == nil {
		return 0
	}
	return x.TotalRevenue
}
func (x *GetRealtimeFinancialSummaryResponse) GetTotalExpense() int64 {
	if x == nil {
		return 0
	}
	return x.TotalExpense
}
func (x *GetRealtimeFinancialSummaryResponse) GetNetIncome() int64 {
	if x == nil {
		return 0
	}
	return x.NetIncome
}
func (x *GetRealtimeFinancialSummaryResponse) GetCashBalance() int64 {
	if x == nil {
		return 0
	}
	return x.CashBalance
}
func (x *GetRealtimeFinancialSummaryResponse) GetARBalance() int64 {
	if x == nil {
		return 0
	}
	return x.ARBalance
}
func (x *GetRealtimeFinancialSummaryResponse) GetAPBalance() int64 {
	if x == nil {
		return 0
	}
	return x.APBalance
}
func (x *GetRealtimeFinancialSummaryResponse) GetAccounts() []*RealtimeSummaryAccount {
	if x == nil {
		return nil
	}
	return x.Accounts
}

// ---------------------------------------------------------------------------
// BL-FIN-039 Cash flow dashboard
// ---------------------------------------------------------------------------

type GetCashFlowDashboardRequest struct {
	StartDate string
	EndDate   string
}

func (x *GetCashFlowDashboardRequest) GetStartDate() string {
	if x == nil {
		return ""
	}
	return x.StartDate
}
func (x *GetCashFlowDashboardRequest) GetEndDate() string {
	if x == nil {
		return ""
	}
	return x.EndDate
}

type CashFlowLine struct {
	Date           string
	Description    string
	Amount         int64
	RunningBalance int64
	Category       string // "operating" | "investing" | "financing"
}

func (x *CashFlowLine) GetDate() string {
	if x == nil {
		return ""
	}
	return x.Date
}
func (x *CashFlowLine) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *CashFlowLine) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *CashFlowLine) GetRunningBalance() int64 {
	if x == nil {
		return 0
	}
	return x.RunningBalance
}
func (x *CashFlowLine) GetCategory() string {
	if x == nil {
		return ""
	}
	return x.Category
}

type GetCashFlowDashboardResponse struct {
	StartDate      string
	EndDate        string
	OpeningBalance int64
	ClosingBalance int64
	OperatingNet   int64
	InvestingNet   int64
	FinancingNet   int64
	Lines          []*CashFlowLine
}

func (x *GetCashFlowDashboardResponse) GetStartDate() string {
	if x == nil {
		return ""
	}
	return x.StartDate
}
func (x *GetCashFlowDashboardResponse) GetEndDate() string {
	if x == nil {
		return ""
	}
	return x.EndDate
}
func (x *GetCashFlowDashboardResponse) GetOpeningBalance() int64 {
	if x == nil {
		return 0
	}
	return x.OpeningBalance
}
func (x *GetCashFlowDashboardResponse) GetClosingBalance() int64 {
	if x == nil {
		return 0
	}
	return x.ClosingBalance
}
func (x *GetCashFlowDashboardResponse) GetOperatingNet() int64 {
	if x == nil {
		return 0
	}
	return x.OperatingNet
}
func (x *GetCashFlowDashboardResponse) GetInvestingNet() int64 {
	if x == nil {
		return 0
	}
	return x.InvestingNet
}
func (x *GetCashFlowDashboardResponse) GetFinancingNet() int64 {
	if x == nil {
		return 0
	}
	return x.FinancingNet
}
func (x *GetCashFlowDashboardResponse) GetLines() []*CashFlowLine {
	if x == nil {
		return nil
	}
	return x.Lines
}

// ---------------------------------------------------------------------------
// BL-FIN-040 AR/AP aging alerts
// ---------------------------------------------------------------------------

type GetAgingAlertsRequest struct {
	LedgerType    string // "ar" | "ap" | "both"
	ThresholdDays int32
}

func (x *GetAgingAlertsRequest) GetLedgerType() string {
	if x == nil {
		return ""
	}
	return x.LedgerType
}
func (x *GetAgingAlertsRequest) GetThresholdDays() int32 {
	if x == nil {
		return 0
	}
	return x.ThresholdDays
}

type AgingAlert struct {
	EntityID    string
	EntityName  string
	Amount      int64
	DaysOverdue int32
	LedgerType  string
}

func (x *AgingAlert) GetEntityID() string {
	if x == nil {
		return ""
	}
	return x.EntityID
}
func (x *AgingAlert) GetEntityName() string {
	if x == nil {
		return ""
	}
	return x.EntityName
}
func (x *AgingAlert) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *AgingAlert) GetDaysOverdue() int32 {
	if x == nil {
		return 0
	}
	return x.DaysOverdue
}
func (x *AgingAlert) GetLedgerType() string {
	if x == nil {
		return ""
	}
	return x.LedgerType
}

type GetAgingAlertsResponse struct {
	Alerts       []*AgingAlert
	TotalOverdue int64
}

func (x *GetAgingAlertsResponse) GetAlerts() []*AgingAlert {
	if x == nil {
		return nil
	}
	return x.Alerts
}
func (x *GetAgingAlertsResponse) GetTotalOverdue() int64 {
	if x == nil {
		return 0
	}
	return x.TotalOverdue
}

// ---------------------------------------------------------------------------
// BL-FIN-041 Finance audit trail
// ---------------------------------------------------------------------------

type SearchFinanceAuditLogRequest struct {
	UserID     string
	Action     string
	EntityType string
	EntityID   string
	StartDate  string
	EndDate    string
	PageSize   int32
	Cursor     string
}

func (x *SearchFinanceAuditLogRequest) GetUserID() string {
	if x == nil {
		return ""
	}
	return x.UserID
}
func (x *SearchFinanceAuditLogRequest) GetAction() string {
	if x == nil {
		return ""
	}
	return x.Action
}
func (x *SearchFinanceAuditLogRequest) GetEntityType() string {
	if x == nil {
		return ""
	}
	return x.EntityType
}
func (x *SearchFinanceAuditLogRequest) GetEntityID() string {
	if x == nil {
		return ""
	}
	return x.EntityID
}
func (x *SearchFinanceAuditLogRequest) GetStartDate() string {
	if x == nil {
		return ""
	}
	return x.StartDate
}
func (x *SearchFinanceAuditLogRequest) GetEndDate() string {
	if x == nil {
		return ""
	}
	return x.EndDate
}
func (x *SearchFinanceAuditLogRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}
func (x *SearchFinanceAuditLogRequest) GetCursor() string {
	if x == nil {
		return ""
	}
	return x.Cursor
}

type AuditLogRow struct {
	LogID      string
	UserID     string
	Action     string
	EntityType string
	EntityID   string
	Diff       string
	CreatedAt  string
}

func (x *AuditLogRow) GetLogID() string {
	if x == nil {
		return ""
	}
	return x.LogID
}
func (x *AuditLogRow) GetUserID() string {
	if x == nil {
		return ""
	}
	return x.UserID
}
func (x *AuditLogRow) GetAction() string {
	if x == nil {
		return ""
	}
	return x.Action
}
func (x *AuditLogRow) GetEntityType() string {
	if x == nil {
		return ""
	}
	return x.EntityType
}
func (x *AuditLogRow) GetEntityID() string {
	if x == nil {
		return ""
	}
	return x.EntityID
}
func (x *AuditLogRow) GetDiff() string {
	if x == nil {
		return ""
	}
	return x.Diff
}
func (x *AuditLogRow) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type SearchFinanceAuditLogResponse struct {
	Rows       []*AuditLogRow
	NextCursor string
}

func (x *SearchFinanceAuditLogResponse) GetRows() []*AuditLogRow {
	if x == nil {
		return nil
	}
	return x.Rows
}
func (x *SearchFinanceAuditLogResponse) GetNextCursor() string {
	if x == nil {
		return ""
	}
	return x.NextCursor
}
