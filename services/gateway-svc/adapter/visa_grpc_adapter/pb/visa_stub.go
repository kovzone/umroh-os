// visa_stub.go — hand-written gRPC client stubs for visa-svc Phase 6 RPCs
// (BL-VISA-001..003): TransitionStatus, BulkSubmit, GetApplications.

package pb

import (
	"context"

	"google.golang.org/grpc"
)

const (
	VisaService_TransitionStatus_FullMethodName  = "/pb.VisaService/TransitionStatus"
	VisaService_BulkSubmit_FullMethodName        = "/pb.VisaService/BulkSubmit"
	VisaService_GetApplications_FullMethodName   = "/pb.VisaService/GetApplications"
)

// ---------------------------------------------------------------------------
// Message types — mirror visa-svc pb/visa_messages.go
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

type StatusHistoryEntryPb struct {
	FromStatus string
	ToStatus   string
	Reason     string
	CreatedAt  string
}

func (x *StatusHistoryEntryPb) GetFromStatus() string {
	if x == nil {
		return ""
	}
	return x.FromStatus
}
func (x *StatusHistoryEntryPb) GetToStatus() string {
	if x == nil {
		return ""
	}
	return x.ToStatus
}
func (x *StatusHistoryEntryPb) GetReason() string {
	if x == nil {
		return ""
	}
	return x.Reason
}
func (x *StatusHistoryEntryPb) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}

type ApplicationRecordPb struct {
	Id          string
	JamaahId    string
	Status      string
	ProviderRef string
	IssuedDate  string
	History     []*StatusHistoryEntryPb
}

func (x *ApplicationRecordPb) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *ApplicationRecordPb) GetJamaahId() string {
	if x == nil {
		return ""
	}
	return x.JamaahId
}
func (x *ApplicationRecordPb) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *ApplicationRecordPb) GetProviderRef() string {
	if x == nil {
		return ""
	}
	return x.ProviderRef
}
func (x *ApplicationRecordPb) GetIssuedDate() string {
	if x == nil {
		return ""
	}
	return x.IssuedDate
}
func (x *ApplicationRecordPb) GetHistory() []*StatusHistoryEntryPb {
	if x == nil {
		return nil
	}
	return x.History
}

type GetApplicationsResponse struct {
	Applications []*ApplicationRecordPb
}

func (x *GetApplicationsResponse) GetApplications() []*ApplicationRecordPb {
	if x == nil {
		return nil
	}
	return x.Applications
}

// ---------------------------------------------------------------------------
// Client interface + implementation
// ---------------------------------------------------------------------------

// VisaServiceClient is the narrow interface for visa-svc Phase 6 RPCs.
type VisaServiceClient interface {
	TransitionStatus(ctx context.Context, in *TransitionStatusRequest, opts ...grpc.CallOption) (*TransitionStatusResponse, error)
	BulkSubmit(ctx context.Context, in *BulkSubmitRequest, opts ...grpc.CallOption) (*BulkSubmitResponse, error)
	GetApplications(ctx context.Context, in *GetApplicationsRequest, opts ...grpc.CallOption) (*GetApplicationsResponse, error)
}

type visaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVisaServiceClient(cc grpc.ClientConnInterface) VisaServiceClient {
	return &visaServiceClient{cc}
}

func (c *visaServiceClient) TransitionStatus(ctx context.Context, in *TransitionStatusRequest, opts ...grpc.CallOption) (*TransitionStatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TransitionStatusResponse)
	if err := c.cc.Invoke(ctx, VisaService_TransitionStatus_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *visaServiceClient) BulkSubmit(ctx context.Context, in *BulkSubmitRequest, opts ...grpc.CallOption) (*BulkSubmitResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BulkSubmitResponse)
	if err := c.cc.Invoke(ctx, VisaService_BulkSubmit_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *visaServiceClient) GetApplications(ctx context.Context, in *GetApplicationsRequest, opts ...grpc.CallOption) (*GetApplicationsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetApplicationsResponse)
	if err := c.cc.Invoke(ctx, VisaService_GetApplications_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
