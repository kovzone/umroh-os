// ops_phase6_stub.go — hand-written gRPC client stubs for ops-svc Phase 6 RPCs
// (BL-OPS-010/011): RecordScan, RecordBusBoarding, GetBoardingRoster.

package pb

import (
	"context"

	"google.golang.org/grpc"
)

const (
	OpsService_RecordScan_FullMethodName         = "/pb.OpsService/RecordScan"
	OpsService_RecordBusBoarding_FullMethodName  = "/pb.OpsService/RecordBusBoarding"
	OpsService_GetBoardingRoster_FullMethodName  = "/pb.OpsService/GetBoardingRoster"
)

// RecordScanRequest mirrors ops-svc pb.RecordScanRequest.
type RecordScanRequest struct {
	ScanType       string
	DepartureId    string
	JamaahId       string
	ScannedBy      string
	DeviceId       string
	Location       string
	IdempotencyKey string
	Metadata       []byte
}

func (x *RecordScanRequest) GetScanType() string {
	if x == nil {
		return ""
	}
	return x.ScanType
}
func (x *RecordScanRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *RecordScanRequest) GetJamaahId() string {
	if x == nil {
		return ""
	}
	return x.JamaahId
}
func (x *RecordScanRequest) GetScannedBy() string {
	if x == nil {
		return ""
	}
	return x.ScannedBy
}
func (x *RecordScanRequest) GetDeviceId() string {
	if x == nil {
		return ""
	}
	return x.DeviceId
}
func (x *RecordScanRequest) GetLocation() string {
	if x == nil {
		return ""
	}
	return x.Location
}
func (x *RecordScanRequest) GetIdempotencyKey() string {
	if x == nil {
		return ""
	}
	return x.IdempotencyKey
}
func (x *RecordScanRequest) GetMetadata() []byte {
	if x == nil {
		return nil
	}
	return x.Metadata
}

type RecordScanResponse struct {
	ScanId     string
	Idempotent bool
}

func (x *RecordScanResponse) GetScanId() string {
	if x == nil {
		return ""
	}
	return x.ScanId
}
func (x *RecordScanResponse) GetIdempotent() bool {
	if x == nil {
		return false
	}
	return x.Idempotent
}

// RecordBusBoardingRequest mirrors ops-svc pb.RecordBusBoardingRequest.
type RecordBusBoardingRequest struct {
	DepartureId string
	BusNumber   string
	JamaahId    string
	ScannedBy   string
	Status      string
}

func (x *RecordBusBoardingRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *RecordBusBoardingRequest) GetBusNumber() string {
	if x == nil {
		return ""
	}
	return x.BusNumber
}
func (x *RecordBusBoardingRequest) GetJamaahId() string {
	if x == nil {
		return ""
	}
	return x.JamaahId
}
func (x *RecordBusBoardingRequest) GetScannedBy() string {
	if x == nil {
		return ""
	}
	return x.ScannedBy
}
func (x *RecordBusBoardingRequest) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

type RecordBusBoardingResponse struct {
	BoardingId string
	Status     string
	Idempotent bool
}

func (x *RecordBusBoardingResponse) GetBoardingId() string {
	if x == nil {
		return ""
	}
	return x.BoardingId
}
func (x *RecordBusBoardingResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *RecordBusBoardingResponse) GetIdempotent() bool {
	if x == nil {
		return false
	}
	return x.Idempotent
}

// GetBoardingRosterRequest mirrors ops-svc pb.GetBoardingRosterRequest.
type GetBoardingRosterRequest struct {
	DepartureId string
	BusNumber   string
}

func (x *GetBoardingRosterRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *GetBoardingRosterRequest) GetBusNumber() string {
	if x == nil {
		return ""
	}
	return x.BusNumber
}

type BoardingEntryPb struct {
	JamaahId  string
	BusNumber string
	Status    string
	BoardedAt string
}

func (x *BoardingEntryPb) GetJamaahId() string {
	if x == nil {
		return ""
	}
	return x.JamaahId
}
func (x *BoardingEntryPb) GetBusNumber() string {
	if x == nil {
		return ""
	}
	return x.BusNumber
}
func (x *BoardingEntryPb) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *BoardingEntryPb) GetBoardedAt() string {
	if x == nil {
		return ""
	}
	return x.BoardedAt
}

type GetBoardingRosterResponse struct {
	Boardings    []*BoardingEntryPb
	TotalBoarded int32
	TotalAbsent  int32
}

func (x *GetBoardingRosterResponse) GetBoardings() []*BoardingEntryPb {
	if x == nil {
		return nil
	}
	return x.Boardings
}
func (x *GetBoardingRosterResponse) GetTotalBoarded() int32 {
	if x == nil {
		return 0
	}
	return x.TotalBoarded
}
func (x *GetBoardingRosterResponse) GetTotalAbsent() int32 {
	if x == nil {
		return 0
	}
	return x.TotalAbsent
}

// ---------------------------------------------------------------------------
// Client interface + implementation
// ---------------------------------------------------------------------------

// OpsPhase6Client is the narrow interface for Phase 6 ops RPCs.
type OpsPhase6Client interface {
	RecordScan(ctx context.Context, in *RecordScanRequest, opts ...grpc.CallOption) (*RecordScanResponse, error)
	RecordBusBoarding(ctx context.Context, in *RecordBusBoardingRequest, opts ...grpc.CallOption) (*RecordBusBoardingResponse, error)
	GetBoardingRoster(ctx context.Context, in *GetBoardingRosterRequest, opts ...grpc.CallOption) (*GetBoardingRosterResponse, error)
}

type opsPhase6Client struct {
	cc grpc.ClientConnInterface
}

func NewOpsPhase6Client(cc grpc.ClientConnInterface) OpsPhase6Client {
	return &opsPhase6Client{cc}
}

func (c *opsPhase6Client) RecordScan(ctx context.Context, in *RecordScanRequest, opts ...grpc.CallOption) (*RecordScanResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecordScanResponse)
	if err := c.cc.Invoke(ctx, OpsService_RecordScan_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *opsPhase6Client) RecordBusBoarding(ctx context.Context, in *RecordBusBoardingRequest, opts ...grpc.CallOption) (*RecordBusBoardingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecordBusBoardingResponse)
	if err := c.cc.Invoke(ctx, OpsService_RecordBusBoarding_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *opsPhase6Client) GetBoardingRoster(ctx context.Context, in *GetBoardingRosterRequest, opts ...grpc.CallOption) (*GetBoardingRosterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetBoardingRosterResponse)
	if err := c.cc.Invoke(ctx, OpsService_GetBoardingRoster_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
