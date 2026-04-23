// correction_messages.go — hand-written proto message types for journal
// correction and anti-delete guard RPCs (BL-FIN-006).
//
// CorrectJournal  — posts a reversing counter-entry for an existing journal
//                   entry, preserving the full audit trail.
// DeleteJournalEntry — always returns PermissionDenied; corrections must use
//                   CorrectJournal instead of deleting entries.

package pb

// ---------------------------------------------------------------------------
// CorrectJournal
// ---------------------------------------------------------------------------

// CorrectJournalRequest requests a reversing counter-entry for an existing
// journal entry identified by entry_id.
type CorrectJournalRequest struct {
	// EntryId is the UUID of the journal entry to reverse.
	EntryId string
	// Reason is a human-readable explanation for the correction (audit trail).
	Reason string
	// ActorUserId is the IAM user triggering this operation.
	ActorUserId string
}

func (x *CorrectJournalRequest) GetEntryId() string {
	if x == nil {
		return ""
	}
	return x.EntryId
}
func (x *CorrectJournalRequest) GetReason() string {
	if x == nil {
		return ""
	}
	return x.Reason
}
func (x *CorrectJournalRequest) GetActorUserId() string {
	if x == nil {
		return ""
	}
	return x.ActorUserId
}

// CorrectJournalResponse carries the correction entry details.
type CorrectJournalResponse struct {
	// CorrectionEntryId is the UUID of the new reversing journal entry.
	CorrectionEntryId string
	// OriginalEntryId is the UUID of the original journal entry that was reversed.
	OriginalEntryId string
	// Idempotent is true when the correction already existed (safe replay).
	Idempotent bool
}

func (x *CorrectJournalResponse) GetCorrectionEntryId() string {
	if x == nil {
		return ""
	}
	return x.CorrectionEntryId
}
func (x *CorrectJournalResponse) GetOriginalEntryId() string {
	if x == nil {
		return ""
	}
	return x.OriginalEntryId
}
func (x *CorrectJournalResponse) GetIdempotent() bool {
	if x == nil {
		return false
	}
	return x.Idempotent
}

// ---------------------------------------------------------------------------
// DeleteJournalEntry (anti-delete guard — always returns PermissionDenied)
// ---------------------------------------------------------------------------

// DeleteJournalEntryRequest is accepted but always rejected.
// This RPC exists solely to return a clear error message directing callers
// to use CorrectJournal instead.
type DeleteJournalEntryRequest struct {
	EntryId string
}

func (x *DeleteJournalEntryRequest) GetEntryId() string {
	if x == nil {
		return ""
	}
	return x.EntryId
}

// DeleteJournalEntryResponse is never actually returned (PermissionDenied
// is returned before this can be populated).
type DeleteJournalEntryResponse struct {
	// Deleted is always false (never returned; the guard always rejects).
	Deleted bool
}

func (x *DeleteJournalEntryResponse) GetDeleted() bool {
	return false
}
