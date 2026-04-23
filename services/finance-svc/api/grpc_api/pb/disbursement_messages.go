// disbursement_messages.go — hand-written proto message types for AP
// disbursement and AR/AP aging RPCs (BL-FIN-010/011).

package pb

// ---------------------------------------------------------------------------
// CreateDisbursementBatch (BL-FIN-010)
// ---------------------------------------------------------------------------

type DisbursementItemInput struct {
	VendorName  string
	Description string
	AmountIdr   int64
	Reference   string
}

type CreateDisbursementBatchRequest struct {
	Description string
	Items       []*DisbursementItemInput
	CreatedBy   string
}

func (x *CreateDisbursementBatchRequest) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *CreateDisbursementBatchRequest) GetItems() []*DisbursementItemInput {
	if x == nil {
		return nil
	}
	return x.Items
}
func (x *CreateDisbursementBatchRequest) GetCreatedBy() string {
	if x == nil {
		return ""
	}
	return x.CreatedBy
}

type CreateDisbursementBatchResponse struct {
	BatchId        string
	TotalAmountIdr int64
	ItemCount      int32
	Status         string
}

func (x *CreateDisbursementBatchResponse) GetBatchId() string {
	if x == nil {
		return ""
	}
	return x.BatchId
}
func (x *CreateDisbursementBatchResponse) GetTotalAmountIdr() int64 {
	if x == nil {
		return 0
	}
	return x.TotalAmountIdr
}
func (x *CreateDisbursementBatchResponse) GetItemCount() int32 {
	if x == nil {
		return 0
	}
	return x.ItemCount
}
func (x *CreateDisbursementBatchResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

// ---------------------------------------------------------------------------
// ApproveDisbursement (BL-FIN-010)
// ---------------------------------------------------------------------------

type ApproveDisbursementRequest struct {
	BatchId    string
	ApprovedBy string
	Approved   bool
	Notes      string
}

func (x *ApproveDisbursementRequest) GetBatchId() string {
	if x == nil {
		return ""
	}
	return x.BatchId
}
func (x *ApproveDisbursementRequest) GetApprovedBy() string {
	if x == nil {
		return ""
	}
	return x.ApprovedBy
}
func (x *ApproveDisbursementRequest) GetApproved() bool {
	if x == nil {
		return false
	}
	return x.Approved
}
func (x *ApproveDisbursementRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type ApproveDisbursementResponse struct {
	BatchId         string
	Status          string
	JournalEntryIds []string
}

func (x *ApproveDisbursementResponse) GetBatchId() string {
	if x == nil {
		return ""
	}
	return x.BatchId
}
func (x *ApproveDisbursementResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *ApproveDisbursementResponse) GetJournalEntryIds() []string {
	if x == nil {
		return nil
	}
	return x.JournalEntryIds
}

// ---------------------------------------------------------------------------
// GetARAPAging (BL-FIN-011)
// ---------------------------------------------------------------------------

type GetARAPAgingRequest struct {
	Type      string // "AR" | "AP" | "both"
	AsOfDate  string // "YYYY-MM-DD"; empty = today
}

func (x *GetARAPAgingRequest) GetType() string {
	if x == nil {
		return ""
	}
	return x.Type
}
func (x *GetARAPAgingRequest) GetAsOfDate() string {
	if x == nil {
		return ""
	}
	return x.AsOfDate
}

type AgingBucketsPb struct {
	Current int64
	Days30  int64
	Days60  int64
	Days90  int64
	Over90  int64
	Total   int64
}

type GetARAPAgingResponse struct {
	Ar          *AgingBucketsPb
	Ap          *AgingBucketsPb
	GeneratedAt string
}

func (x *GetARAPAgingResponse) GetAr() *AgingBucketsPb {
	if x == nil {
		return nil
	}
	return x.Ar
}
func (x *GetARAPAgingResponse) GetAp() *AgingBucketsPb {
	if x == nil {
		return nil
	}
	return x.Ap
}
func (x *GetARAPAgingResponse) GetGeneratedAt() string {
	if x == nil {
		return ""
	}
	return x.GeneratedAt
}
