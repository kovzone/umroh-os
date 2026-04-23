// Package pb — hand-written gRPC client stub for the gateway-svc → crm-svc
// connection (S4-E-02 / ADR-0009).
//
// gateway-svc is the edge proxy; it calls four RPCs on crm-svc:
//   - CreateLead  (POST /v1/leads  — public)
//   - GetLead     (GET  /v1/leads/:id — bearer)
//   - UpdateLead  (PUT  /v1/leads/:id — bearer)
//   - ListLeads   (GET  /v1/leads     — bearer cs/admin)
//
// Run `make genpb` with a shared proto path to replace with generated code.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

const (
	CrmService_CreateLead_FullMethodName = "/pb.CrmService/CreateLead"
	CrmService_GetLead_FullMethodName    = "/pb.CrmService/GetLead"
	CrmService_UpdateLead_FullMethodName = "/pb.CrmService/UpdateLead"
	CrmService_ListLeads_FullMethodName  = "/pb.CrmService/ListLeads"
)

// ---------------------------------------------------------------------------
// LeadResponse
// ---------------------------------------------------------------------------

type LeadResponse struct {
	Id                  string
	Source              string
	UtmSource           string
	UtmMedium           string
	UtmCampaign         string
	UtmContent          string
	UtmTerm             string
	Name                string
	Phone               string
	Email               string
	InterestPackageId   string
	InterestDepartureId string
	Status              string
	AssignedCsId        string
	Notes               string
	BookingId           string
	CreatedAt           string
	UpdatedAt           string
}

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
// CreateLead
// ---------------------------------------------------------------------------

type CreateLeadRequest struct {
	Name                string
	Phone               string
	Email               string
	Source              string
	UtmSource           string
	UtmMedium           string
	UtmCampaign         string
	UtmContent          string
	UtmTerm             string
	InterestPackageId   string
	InterestDepartureId string
	Notes               string
}

// ---------------------------------------------------------------------------
// GetLead
// ---------------------------------------------------------------------------

type GetLeadRequest struct {
	Id string
}

// ---------------------------------------------------------------------------
// UpdateLead
// ---------------------------------------------------------------------------

type UpdateLeadRequest struct {
	Id           string
	Status       string
	Notes        string
	AssignedCsId string
}

// ---------------------------------------------------------------------------
// ListLeads
// ---------------------------------------------------------------------------

type ListLeadsRequest struct {
	StatusFilter       string
	AssignedCsIdFilter string
	Page               int32
	PageSize           int32
}

type ListLeadsResponse struct {
	Leads    []*LeadResponse
	Total    int64
	Page     int32
	PageSize int32
}

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
// CrmServiceClient interface + implementation
// ---------------------------------------------------------------------------

type CrmServiceClient interface {
	CreateLead(ctx context.Context, in *CreateLeadRequest, opts ...grpc.CallOption) (*LeadResponse, error)
	GetLead(ctx context.Context, in *GetLeadRequest, opts ...grpc.CallOption) (*LeadResponse, error)
	UpdateLead(ctx context.Context, in *UpdateLeadRequest, opts ...grpc.CallOption) (*LeadResponse, error)
	ListLeads(ctx context.Context, in *ListLeadsRequest, opts ...grpc.CallOption) (*ListLeadsResponse, error)
}

type crmServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCrmServiceClient(cc grpc.ClientConnInterface) CrmServiceClient {
	return &crmServiceClient{cc}
}

func (c *crmServiceClient) CreateLead(ctx context.Context, in *CreateLeadRequest, opts ...grpc.CallOption) (*LeadResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LeadResponse)
	err := c.cc.Invoke(ctx, CrmService_CreateLead_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crmServiceClient) GetLead(ctx context.Context, in *GetLeadRequest, opts ...grpc.CallOption) (*LeadResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LeadResponse)
	err := c.cc.Invoke(ctx, CrmService_GetLead_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crmServiceClient) UpdateLead(ctx context.Context, in *UpdateLeadRequest, opts ...grpc.CallOption) (*LeadResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LeadResponse)
	err := c.cc.Invoke(ctx, CrmService_UpdateLead_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *crmServiceClient) ListLeads(ctx context.Context, in *ListLeadsRequest, opts ...grpc.CallOption) (*ListLeadsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListLeadsResponse)
	err := c.cc.Invoke(ctx, CrmService_ListLeads_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
