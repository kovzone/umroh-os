// grn_messages.go — hand-written proto message types for the OnGRNReceived RPC
// (BL-FIN-002).
//
// OnGRNReceived is called when a Goods Receipt Note (GRN) is received,
// posting Dr 5001 (COGS/Inventory Expense) / Cr 2001 (AP/Pilgrim Liability).
// Idempotent on idempotency_key = "grn:" + grn_id.

package pb

// OnGRNReceivedRequest carries the GRN details needed to create the AP journal.
type OnGRNReceivedRequest struct {
	GrnId       string // UUID of the goods receipt note
	DepartureId string // UUID of the associated departure
	AmountIdr   int64  // integer IDR — no fractional amounts
}

func (x *OnGRNReceivedRequest) GetGrnId() string {
	if x == nil {
		return ""
	}
	return x.GrnId
}

func (x *OnGRNReceivedRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}

func (x *OnGRNReceivedRequest) GetAmountIdr() int64 {
	if x == nil {
		return 0
	}
	return x.AmountIdr
}

// OnGRNReceivedResponse returns the journal entry that was created (or the
// existing entry if already posted for this GRN — idempotent).
type OnGRNReceivedResponse struct {
	EntryId    string // UUID of the journal entry
	Idempotent bool   // true if an existing entry was returned (replayed)
}

func (x *OnGRNReceivedResponse) GetEntryId() string {
	if x == nil {
		return ""
	}
	return x.EntryId
}

func (x *OnGRNReceivedResponse) GetIdempotent() bool {
	if x == nil {
		return false
	}
	return x.Idempotent
}
