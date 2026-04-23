// Package pb — hand-written proto message types for CrmService (S4-E-02).
//
// These supplement the protoc-generated crm.pb.go / crm_grpc.pb.go.
// All new message types for CreateLead, GetLead, UpdateLead, ListLeads,
// OnBookingCreated, OnBookingPaidInFull are defined here.
// Run `make genpb` to regenerate from crm.proto and delete this file.
package pb

// ---------------------------------------------------------------------------
// LeadResponse — returned by all lead-mutating and read RPCs.
// ---------------------------------------------------------------------------

type LeadResponse struct {
	Id                  string `json:"id"`
	Source              string `json:"source"`
	UtmSource           string `json:"utm_source"`
	UtmMedium           string `json:"utm_medium"`
	UtmCampaign         string `json:"utm_campaign"`
	UtmContent          string `json:"utm_content"`
	UtmTerm             string `json:"utm_term"`
	Name                string `json:"name"`
	Phone               string `json:"phone"`
	Email               string `json:"email"`
	InterestPackageId   string `json:"interest_package_id"`
	InterestDepartureId string `json:"interest_departure_id"`
	Status              string `json:"status"`
	AssignedCsId        string `json:"assigned_cs_id"`
	Notes               string `json:"notes"`
	BookingId           string `json:"booking_id"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

func (x *LeadResponse) Reset()         {}
func (x *LeadResponse) String() string { return x.Id }
func (x *LeadResponse) ProtoMessage()  {}

func (x *LeadResponse) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *LeadResponse) GetSource() string {
	if x == nil {
		return ""
	}
	return x.Source
}
func (x *LeadResponse) GetUtmSource() string {
	if x == nil {
		return ""
	}
	return x.UtmSource
}
func (x *LeadResponse) GetUtmMedium() string {
	if x == nil {
		return ""
	}
	return x.UtmMedium
}
func (x *LeadResponse) GetUtmCampaign() string {
	if x == nil {
		return ""
	}
	return x.UtmCampaign
}
func (x *LeadResponse) GetUtmContent() string {
	if x == nil {
		return ""
	}
	return x.UtmContent
}
func (x *LeadResponse) GetUtmTerm() string {
	if x == nil {
		return ""
	}
	return x.UtmTerm
}
func (x *LeadResponse) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *LeadResponse) GetPhone() string {
	if x == nil {
		return ""
	}
	return x.Phone
}
func (x *LeadResponse) GetEmail() string {
	if x == nil {
		return ""
	}
	return x.Email
}
func (x *LeadResponse) GetInterestPackageId() string {
	if x == nil {
		return ""
	}
	return x.InterestPackageId
}
func (x *LeadResponse) GetInterestDepartureId() string {
	if x == nil {
		return ""
	}
	return x.InterestDepartureId
}
func (x *LeadResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *LeadResponse) GetAssignedCsId() string {
	if x == nil {
		return ""
	}
	return x.AssignedCsId
}
func (x *LeadResponse) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}
func (x *LeadResponse) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}
func (x *LeadResponse) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}
func (x *LeadResponse) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

// ---------------------------------------------------------------------------
// CreateLeadRequest
// ---------------------------------------------------------------------------

type CreateLeadRequest struct {
	Name                string `json:"name"`
	Phone               string `json:"phone"`
	Email               string `json:"email"`
	Source              string `json:"source"`
	UtmSource           string `json:"utm_source"`
	UtmMedium           string `json:"utm_medium"`
	UtmCampaign         string `json:"utm_campaign"`
	UtmContent          string `json:"utm_content"`
	UtmTerm             string `json:"utm_term"`
	InterestPackageId   string `json:"interest_package_id"`
	InterestDepartureId string `json:"interest_departure_id"`
	Notes               string `json:"notes"`
}

func (x *CreateLeadRequest) Reset()         {}
func (x *CreateLeadRequest) String() string { return "" }
func (x *CreateLeadRequest) ProtoMessage()  {}

func (x *CreateLeadRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *CreateLeadRequest) GetPhone() string {
	if x == nil {
		return ""
	}
	return x.Phone
}
func (x *CreateLeadRequest) GetEmail() string {
	if x == nil {
		return ""
	}
	return x.Email
}
func (x *CreateLeadRequest) GetSource() string {
	if x == nil {
		return ""
	}
	return x.Source
}
func (x *CreateLeadRequest) GetUtmSource() string {
	if x == nil {
		return ""
	}
	return x.UtmSource
}
func (x *CreateLeadRequest) GetUtmMedium() string {
	if x == nil {
		return ""
	}
	return x.UtmMedium
}
func (x *CreateLeadRequest) GetUtmCampaign() string {
	if x == nil {
		return ""
	}
	return x.UtmCampaign
}
func (x *CreateLeadRequest) GetUtmContent() string {
	if x == nil {
		return ""
	}
	return x.UtmContent
}
func (x *CreateLeadRequest) GetUtmTerm() string {
	if x == nil {
		return ""
	}
	return x.UtmTerm
}
func (x *CreateLeadRequest) GetInterestPackageId() string {
	if x == nil {
		return ""
	}
	return x.InterestPackageId
}
func (x *CreateLeadRequest) GetInterestDepartureId() string {
	if x == nil {
		return ""
	}
	return x.InterestDepartureId
}
func (x *CreateLeadRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

// ---------------------------------------------------------------------------
// GetLeadRequest
// ---------------------------------------------------------------------------

type GetLeadRequest struct {
	Id string `json:"id"`
}

func (x *GetLeadRequest) Reset()         {}
func (x *GetLeadRequest) String() string { return "" }
func (x *GetLeadRequest) ProtoMessage()  {}

func (x *GetLeadRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}

// ---------------------------------------------------------------------------
// UpdateLeadRequest
// ---------------------------------------------------------------------------

type UpdateLeadRequest struct {
	Id           string `json:"id"`
	Status       string `json:"status"`        // empty = no change
	Notes        string `json:"notes"`         // empty = no change
	AssignedCsId string `json:"assigned_cs_id"` // empty = no change
}

func (x *UpdateLeadRequest) Reset()         {}
func (x *UpdateLeadRequest) String() string { return "" }
func (x *UpdateLeadRequest) ProtoMessage()  {}

func (x *UpdateLeadRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *UpdateLeadRequest) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *UpdateLeadRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}
func (x *UpdateLeadRequest) GetAssignedCsId() string {
	if x == nil {
		return ""
	}
	return x.AssignedCsId
}

// ---------------------------------------------------------------------------
// ListLeadsRequest / ListLeadsResponse
// ---------------------------------------------------------------------------

type ListLeadsRequest struct {
	StatusFilter      string `json:"status_filter"`
	AssignedCsIdFilter string `json:"assigned_cs_id_filter"`
	Page              int32  `json:"page"`
	PageSize          int32  `json:"page_size"`
}

func (x *ListLeadsRequest) Reset()         {}
func (x *ListLeadsRequest) String() string { return "" }
func (x *ListLeadsRequest) ProtoMessage()  {}

func (x *ListLeadsRequest) GetStatusFilter() string {
	if x == nil {
		return ""
	}
	return x.StatusFilter
}
func (x *ListLeadsRequest) GetAssignedCsIdFilter() string {
	if x == nil {
		return ""
	}
	return x.AssignedCsIdFilter
}
func (x *ListLeadsRequest) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *ListLeadsRequest) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

type ListLeadsResponse struct {
	Leads    []*LeadResponse `json:"leads"`
	Total    int64           `json:"total"`
	Page     int32           `json:"page"`
	PageSize int32           `json:"page_size"`
}

func (x *ListLeadsResponse) Reset()         {}
func (x *ListLeadsResponse) String() string { return "" }
func (x *ListLeadsResponse) ProtoMessage()  {}

func (x *ListLeadsResponse) GetLeads() []*LeadResponse {
	if x == nil {
		return nil
	}
	return x.Leads
}
func (x *ListLeadsResponse) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}
func (x *ListLeadsResponse) GetPage() int32 {
	if x == nil {
		return 0
	}
	return x.Page
}
func (x *ListLeadsResponse) GetPageSize() int32 {
	if x == nil {
		return 0
	}
	return x.PageSize
}

// ---------------------------------------------------------------------------
// OnBookingCreatedRequest / OnBookingCreatedResponse
// ---------------------------------------------------------------------------

type OnBookingCreatedRequest struct {
	BookingId    string `json:"booking_id"`
	LeadId       string `json:"lead_id"`
	PackageId    string `json:"package_id"`
	DepartureId  string `json:"departure_id"`
	JamaahCount  int32  `json:"jamaah_count"`
	CreatedAt    string `json:"created_at"`
}

func (x *OnBookingCreatedRequest) Reset()         {}
func (x *OnBookingCreatedRequest) String() string { return "" }
func (x *OnBookingCreatedRequest) ProtoMessage()  {}

func (x *OnBookingCreatedRequest) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}
func (x *OnBookingCreatedRequest) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}
func (x *OnBookingCreatedRequest) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}
func (x *OnBookingCreatedRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *OnBookingCreatedRequest) GetJamaahCount() int32 {
	if x == nil {
		return 0
	}
	return x.JamaahCount
}
func (x *OnBookingCreatedRequest) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type OnBookingCreatedResponse struct {
	Updated bool   `json:"updated"`
	LeadId  string `json:"lead_id"`
}

func (x *OnBookingCreatedResponse) Reset()         {}
func (x *OnBookingCreatedResponse) String() string { return "" }
func (x *OnBookingCreatedResponse) ProtoMessage()  {}

func (x *OnBookingCreatedResponse) GetUpdated() bool {
	if x == nil {
		return false
	}
	return x.Updated
}
func (x *OnBookingCreatedResponse) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}

// ---------------------------------------------------------------------------
// OnBookingPaidInFullRequest / OnBookingPaidInFullResponse
// ---------------------------------------------------------------------------

type OnBookingPaidInFullRequest struct {
	BookingId string `json:"booking_id"`
	LeadId    string `json:"lead_id"`
	PaidAt    string `json:"paid_at"`
}

func (x *OnBookingPaidInFullRequest) Reset()         {}
func (x *OnBookingPaidInFullRequest) String() string { return "" }
func (x *OnBookingPaidInFullRequest) ProtoMessage()  {}

func (x *OnBookingPaidInFullRequest) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}
func (x *OnBookingPaidInFullRequest) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}
func (x *OnBookingPaidInFullRequest) GetPaidAt() string {
	if x == nil {
		return ""
	}
	return x.PaidAt
}

type OnBookingPaidInFullResponse struct {
	Updated bool   `json:"updated"`
	LeadId  string `json:"lead_id"`
}

func (x *OnBookingPaidInFullResponse) Reset()         {}
func (x *OnBookingPaidInFullResponse) String() string { return "" }
func (x *OnBookingPaidInFullResponse) ProtoMessage()  {}

func (x *OnBookingPaidInFullResponse) GetUpdated() bool {
	if x == nil {
		return false
	}
	return x.Updated
}
func (x *OnBookingPaidInFullResponse) GetLeadId() string {
	if x == nil {
		return ""
	}
	return x.LeadId
}
