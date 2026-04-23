// ops_stub.go — gateway-side gRPC client stub for ops-svc RPCs (S3 Wave 2).
//
// Mirrors services/ops-svc/api/grpc_api/pb/ops_messages.go.
// Run `make genpb` to replace with generated code once protoc is available.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

const (
	OpsService_RunRoomAllocation_FullMethodName = "/pb.OpsService/RunRoomAllocation"
	OpsService_GetRoomAllocation_FullMethodName = "/pb.OpsService/GetRoomAllocation"
	OpsService_GenerateIDCard_FullMethodName    = "/pb.OpsService/GenerateIDCard"
	OpsService_VerifyIDCard_FullMethodName      = "/pb.OpsService/VerifyIDCard"
	OpsService_ExportManifest_FullMethodName    = "/pb.OpsService/ExportManifest"
)

// ---------------------------------------------------------------------------
// Domain objects
// ---------------------------------------------------------------------------

// RoomAssignmentProto represents a single room assignment.
type RoomAssignmentProto struct {
	RoomNumber string
	JamaahId   string
}

func (x *RoomAssignmentProto) GetRoomNumber() string {
	if x == nil {
		return ""
	}
	return x.RoomNumber
}
func (x *RoomAssignmentProto) GetJamaahId() string {
	if x == nil {
		return ""
	}
	return x.JamaahId
}

// ManifestRowProto is a single row in the manifest export.
type ManifestRowProto struct {
	No         int32
	JamaahName string
	PassportNo string
	DocStatus  string
	RoomNumber string
}

func (x *ManifestRowProto) GetNo() int32 {
	if x == nil {
		return 0
	}
	return x.No
}
func (x *ManifestRowProto) GetJamaahName() string {
	if x == nil {
		return ""
	}
	return x.JamaahName
}
func (x *ManifestRowProto) GetPassportNo() string {
	if x == nil {
		return ""
	}
	return x.PassportNo
}
func (x *ManifestRowProto) GetDocStatus() string {
	if x == nil {
		return ""
	}
	return x.DocStatus
}
func (x *ManifestRowProto) GetRoomNumber() string {
	if x == nil {
		return ""
	}
	return x.RoomNumber
}

// ---------------------------------------------------------------------------
// Request / response types
// ---------------------------------------------------------------------------

type RunRoomAllocationRequest struct {
	DepartureId string
	JamaahIds   []string
}

func (x *RunRoomAllocationRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *RunRoomAllocationRequest) GetJamaahIds() []string {
	if x == nil {
		return nil
	}
	return x.JamaahIds
}

type RunRoomAllocationResponse struct {
	AllocationId string
	RoomCount    int32
	Assignments  []*RoomAssignmentProto
}

func (x *RunRoomAllocationResponse) GetAllocationId() string {
	if x == nil {
		return ""
	}
	return x.AllocationId
}
func (x *RunRoomAllocationResponse) GetRoomCount() int32 {
	if x == nil {
		return 0
	}
	return x.RoomCount
}
func (x *RunRoomAllocationResponse) GetAssignments() []*RoomAssignmentProto {
	if x == nil {
		return nil
	}
	return x.Assignments
}

type GetRoomAllocationRequest struct {
	DepartureId string
}

func (x *GetRoomAllocationRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}

type GetRoomAllocationResponse struct {
	AllocationId string
	RoomCount    int32
	Assignments  []*RoomAssignmentProto
	Status       string
}

func (x *GetRoomAllocationResponse) GetAllocationId() string {
	if x == nil {
		return ""
	}
	return x.AllocationId
}
func (x *GetRoomAllocationResponse) GetRoomCount() int32 {
	if x == nil {
		return 0
	}
	return x.RoomCount
}
func (x *GetRoomAllocationResponse) GetAssignments() []*RoomAssignmentProto {
	if x == nil {
		return nil
	}
	return x.Assignments
}
func (x *GetRoomAllocationResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

type GenerateIDCardRequest struct {
	JamaahId      string
	DepartureId   string
	CardType      string
	JamaahName    string
	DepartureName string
}

func (x *GenerateIDCardRequest) GetJamaahId() string {
	if x == nil {
		return ""
	}
	return x.JamaahId
}
func (x *GenerateIDCardRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *GenerateIDCardRequest) GetCardType() string {
	if x == nil {
		return ""
	}
	return x.CardType
}
func (x *GenerateIDCardRequest) GetJamaahName() string {
	if x == nil {
		return ""
	}
	return x.JamaahName
}
func (x *GenerateIDCardRequest) GetDepartureName() string {
	if x == nil {
		return ""
	}
	return x.DepartureName
}

type GenerateIDCardResponse struct {
	Token    string
	QrData   string
	IssuedAt string
}

func (x *GenerateIDCardResponse) GetToken() string {
	if x == nil {
		return ""
	}
	return x.Token
}
func (x *GenerateIDCardResponse) GetQrData() string {
	if x == nil {
		return ""
	}
	return x.QrData
}
func (x *GenerateIDCardResponse) GetIssuedAt() string {
	if x == nil {
		return ""
	}
	return x.IssuedAt
}

type VerifyIDCardRequest struct {
	Token string
}

func (x *VerifyIDCardRequest) GetToken() string {
	if x == nil {
		return ""
	}
	return x.Token
}

type VerifyIDCardResponse struct {
	Valid       bool
	JamaahId    string
	DepartureId string
	CardType    string
	ErrorReason string
}

func (x *VerifyIDCardResponse) GetValid() bool {
	if x == nil {
		return false
	}
	return x.Valid
}
func (x *VerifyIDCardResponse) GetJamaahId() string {
	if x == nil {
		return ""
	}
	return x.JamaahId
}
func (x *VerifyIDCardResponse) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *VerifyIDCardResponse) GetCardType() string {
	if x == nil {
		return ""
	}
	return x.CardType
}
func (x *VerifyIDCardResponse) GetErrorReason() string {
	if x == nil {
		return ""
	}
	return x.ErrorReason
}

type ExportManifestRequest struct {
	DepartureId string
}

func (x *ExportManifestRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}

type ExportManifestResponse struct {
	DepartureId string
	Rows        []*ManifestRowProto
}

func (x *ExportManifestResponse) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *ExportManifestResponse) GetRows() []*ManifestRowProto {
	if x == nil {
		return nil
	}
	return x.Rows
}

// ---------------------------------------------------------------------------
// OpsServiceClient interface + implementation
// ---------------------------------------------------------------------------

// OpsServiceClient is the consumer-side interface for ops-svc RPCs.
type OpsServiceClient interface {
	RunRoomAllocation(ctx context.Context, req *RunRoomAllocationRequest, opts ...grpc.CallOption) (*RunRoomAllocationResponse, error)
	GetRoomAllocation(ctx context.Context, req *GetRoomAllocationRequest, opts ...grpc.CallOption) (*GetRoomAllocationResponse, error)
	GenerateIDCard(ctx context.Context, req *GenerateIDCardRequest, opts ...grpc.CallOption) (*GenerateIDCardResponse, error)
	VerifyIDCard(ctx context.Context, req *VerifyIDCardRequest, opts ...grpc.CallOption) (*VerifyIDCardResponse, error)
	ExportManifest(ctx context.Context, req *ExportManifestRequest, opts ...grpc.CallOption) (*ExportManifestResponse, error)
}

type opsServiceClient struct {
	cc grpc.ClientConnInterface
}

// NewOpsServiceClient wraps a conn so gateway-svc can call ops-svc RPCs.
func NewOpsServiceClient(cc grpc.ClientConnInterface) OpsServiceClient {
	return &opsServiceClient{cc}
}

func (c *opsServiceClient) RunRoomAllocation(ctx context.Context, in *RunRoomAllocationRequest, opts ...grpc.CallOption) (*RunRoomAllocationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RunRoomAllocationResponse)
	err := c.cc.Invoke(ctx, OpsService_RunRoomAllocation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *opsServiceClient) GetRoomAllocation(ctx context.Context, in *GetRoomAllocationRequest, opts ...grpc.CallOption) (*GetRoomAllocationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRoomAllocationResponse)
	err := c.cc.Invoke(ctx, OpsService_GetRoomAllocation_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *opsServiceClient) GenerateIDCard(ctx context.Context, in *GenerateIDCardRequest, opts ...grpc.CallOption) (*GenerateIDCardResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GenerateIDCardResponse)
	err := c.cc.Invoke(ctx, OpsService_GenerateIDCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *opsServiceClient) VerifyIDCard(ctx context.Context, in *VerifyIDCardRequest, opts ...grpc.CallOption) (*VerifyIDCardResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VerifyIDCardResponse)
	err := c.cc.Invoke(ctx, OpsService_VerifyIDCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *opsServiceClient) ExportManifest(ctx context.Context, in *ExportManifestRequest, opts ...grpc.CallOption) (*ExportManifestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ExportManifestResponse)
	err := c.cc.Invoke(ctx, OpsService_ExportManifest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
