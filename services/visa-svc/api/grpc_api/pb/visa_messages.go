// visa_messages.go — hand-written proto message types for the visa pipeline RPCs
// (BL-VISA-001..003).
//
// TransitionStatus — state-machine-driven single application status transition.
// BulkSubmit       — atomic all-or-nothing batch READY→SUBMITTED.
// GetApplications  — list applications for a departure with embedded history.

package pb

// ---------------------------------------------------------------------------
// TransitionStatus (BL-VISA-001)
// ---------------------------------------------------------------------------

type TransitionStatusRequest struct {
	ApplicationId string
	ToStatus      string
	Reason        string
	ActorUserId   string
}

func (x *TransitionStatusRequest) GetApplicationId() string {
	if x == nil {
		return ""
	}
	return x.ApplicationId
}
func (x *TransitionStatusRequest) GetToStatus() string {
	if x == nil {
		return ""
	}
	return x.ToStatus
}
func (x *TransitionStatusRequest) GetReason() string {
	if x == nil {
		return ""
	}
	return x.Reason
}
func (x *TransitionStatusRequest) GetActorUserId() string {
	if x == nil {
		return ""
	}
	return x.ActorUserId
}

type TransitionStatusResponse struct {
	ApplicationId string
	FromStatus    string
	ToStatus      string
	Idempotent    bool
}

func (x *TransitionStatusResponse) GetApplicationId() string {
	if x == nil {
		return ""
	}
	return x.ApplicationId
}
func (x *TransitionStatusResponse) GetFromStatus() string {
	if x == nil {
		return ""
	}
	return x.FromStatus
}
func (x *TransitionStatusResponse) GetToStatus() string {
	if x == nil {
		return ""
	}
	return x.ToStatus
}
func (x *TransitionStatusResponse) GetIdempotent() bool {
	if x == nil {
		return false
	}
	return x.Idempotent
}

// ---------------------------------------------------------------------------
// BulkSubmit (BL-VISA-002)
// ---------------------------------------------------------------------------

type BulkSubmitRequest struct {
	DepartureId string
	JamaahIds   []string
	ProviderId  string
}

func (x *BulkSubmitRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *BulkSubmitRequest) GetJamaahIds() []string {
	if x == nil {
		return nil
	}
	return x.JamaahIds
}
func (x *BulkSubmitRequest) GetProviderId() string {
	if x == nil {
		return ""
	}
	return x.ProviderId
}

type BulkSubmitResponse struct {
	SubmittedCount int32
	ApplicationIds []string
}

func (x *BulkSubmitResponse) GetSubmittedCount() int32 {
	if x == nil {
		return 0
	}
	return x.SubmittedCount
}
func (x *BulkSubmitResponse) GetApplicationIds() []string {
	if x == nil {
		return nil
	}
	return x.ApplicationIds
}

// ---------------------------------------------------------------------------
// GetApplications (BL-VISA-003)
// ---------------------------------------------------------------------------

type GetApplicationsRequest struct {
	DepartureId  string
	StatusFilter string
}

func (x *GetApplicationsRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *GetApplicationsRequest) GetStatusFilter() string {
	if x == nil {
		return ""
	}
	return x.StatusFilter
}

type GetApplicationsResponse struct {
	Applications []*ApplicationRecord
}

func (x *GetApplicationsResponse) GetApplications() []*ApplicationRecord {
	if x == nil {
		return nil
	}
	return x.Applications
}

type ApplicationRecord struct {
	Id          string
	JamaahId    string
	Status      string
	ProviderRef string
	IssuedDate  string
	History     []*StatusHistoryEntry
}

type StatusHistoryEntry struct {
	FromStatus string
	ToStatus   string
	Reason     string
	CreatedAt  string // RFC3339
}
