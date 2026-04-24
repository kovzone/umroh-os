// logistics_stub.go — gateway-side gRPC client stub for logistics-svc RPCs (S3 Wave 2).
//
// Mirrors services/logistics-svc/api/grpc_api/pb/logistics_new_messages.go.
// Run `make genpb` to replace with generated code once protoc is available.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

const (
	LogisticsService_ShipFulfillmentTask_FullMethodName  = "/pb.LogisticsService/ShipFulfillmentTask"
	LogisticsService_GeneratePickupQR_FullMethodName     = "/pb.LogisticsService/GeneratePickupQR"
	LogisticsService_RedeemPickupQR_FullMethodName       = "/pb.LogisticsService/RedeemPickupQR"
	LogisticsService_ListFulfillmentTasks_FullMethodName = "/pb.LogisticsService/ListFulfillmentTasks"
)

// ---------------------------------------------------------------------------
// Request / response types
// ---------------------------------------------------------------------------

type ShipFulfillmentTaskRequest struct {
	BookingId string
	Carrier   string
	Notes     string
}

func (x *ShipFulfillmentTaskRequest) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}
func (x *ShipFulfillmentTaskRequest) GetCarrier() string {
	if x == nil {
		return ""
	}
	return x.Carrier
}
func (x *ShipFulfillmentTaskRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type ShipFulfillmentTaskResponse struct {
	ShipmentId     string
	TrackingNumber string
	Status         string
}

func (x *ShipFulfillmentTaskResponse) GetShipmentId() string {
	if x == nil {
		return ""
	}
	return x.ShipmentId
}
func (x *ShipFulfillmentTaskResponse) GetTrackingNumber() string {
	if x == nil {
		return ""
	}
	return x.TrackingNumber
}
func (x *ShipFulfillmentTaskResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

type GeneratePickupQRRequest struct {
	BookingId string
}

func (x *GeneratePickupQRRequest) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}

type GeneratePickupQRResponse struct {
	PickupTokenId string
	Token         string
	ExpiresAt     string
}

func (x *GeneratePickupQRResponse) GetPickupTokenId() string {
	if x == nil {
		return ""
	}
	return x.PickupTokenId
}
func (x *GeneratePickupQRResponse) GetToken() string {
	if x == nil {
		return ""
	}
	return x.Token
}
func (x *GeneratePickupQRResponse) GetExpiresAt() string {
	if x == nil {
		return ""
	}
	return x.ExpiresAt
}

type RedeemPickupQRRequest struct {
	Token string
}

func (x *RedeemPickupQRRequest) GetToken() string {
	if x == nil {
		return ""
	}
	return x.Token
}

type RedeemPickupQRResponse struct {
	Redeemed    bool
	BookingId   string
	TaskId      string
	ErrorReason string
}

func (x *RedeemPickupQRResponse) GetRedeemed() bool {
	if x == nil {
		return false
	}
	return x.Redeemed
}
func (x *RedeemPickupQRResponse) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}
func (x *RedeemPickupQRResponse) GetTaskId() string {
	if x == nil {
		return ""
	}
	return x.TaskId
}
func (x *RedeemPickupQRResponse) GetErrorReason() string {
	if x == nil {
		return ""
	}
	return x.ErrorReason
}

// ---------------------------------------------------------------------------
// ListFulfillmentTasks message types
// ---------------------------------------------------------------------------

// ListFulfillmentTasksRequest mirrors the logistics-svc pb type (ISSUE-018).
type ListFulfillmentTasksRequest struct {
	StatusFilter      string
	DepartureIdFilter string
	Limit             int32
	Offset            int32
}

func (x *ListFulfillmentTasksRequest) GetStatusFilter() string {
	if x == nil {
		return ""
	}
	return x.StatusFilter
}
func (x *ListFulfillmentTasksRequest) GetDepartureIdFilter() string {
	if x == nil {
		return ""
	}
	return x.DepartureIdFilter
}
func (x *ListFulfillmentTasksRequest) GetLimit() int32 {
	if x == nil {
		return 0
	}
	return x.Limit
}
func (x *ListFulfillmentTasksRequest) GetOffset() int32 {
	if x == nil {
		return 0
	}
	return x.Offset
}

// FulfillmentTaskProto is a single task in the list response.
type FulfillmentTaskProto struct {
	Id             string
	BookingId      string
	DepartureId    string
	Status         string
	TrackingNumber string
	ShippedAt      string
	DeliveredAt    string
	CreatedAt      string
	UpdatedAt      string
}

func (x *FulfillmentTaskProto) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *FulfillmentTaskProto) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}
func (x *FulfillmentTaskProto) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *FulfillmentTaskProto) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *FulfillmentTaskProto) GetTrackingNumber() string {
	if x == nil {
		return ""
	}
	return x.TrackingNumber
}
func (x *FulfillmentTaskProto) GetShippedAt() string {
	if x == nil {
		return ""
	}
	return x.ShippedAt
}
func (x *FulfillmentTaskProto) GetDeliveredAt() string {
	if x == nil {
		return ""
	}
	return x.DeliveredAt
}
func (x *FulfillmentTaskProto) GetCreatedAt() string {
	if x == nil {
		return ""
	}
	return x.CreatedAt
}
func (x *FulfillmentTaskProto) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

// ListFulfillmentTasksResponse carries the page of tasks + total count.
type ListFulfillmentTasksResponse struct {
	Tasks []*FulfillmentTaskProto
	Total int64
}

func (x *ListFulfillmentTasksResponse) GetTasks() []*FulfillmentTaskProto {
	if x == nil {
		return nil
	}
	return x.Tasks
}
func (x *ListFulfillmentTasksResponse) GetTotal() int64 {
	if x == nil {
		return 0
	}
	return x.Total
}

// ---------------------------------------------------------------------------
// LogisticsServiceClient interface + implementation
// ---------------------------------------------------------------------------

// LogisticsServiceClient is the consumer-side interface for logistics-svc RPCs.
type LogisticsServiceClient interface {
	ShipFulfillmentTask(ctx context.Context, req *ShipFulfillmentTaskRequest, opts ...grpc.CallOption) (*ShipFulfillmentTaskResponse, error)
	GeneratePickupQR(ctx context.Context, req *GeneratePickupQRRequest, opts ...grpc.CallOption) (*GeneratePickupQRResponse, error)
	RedeemPickupQR(ctx context.Context, req *RedeemPickupQRRequest, opts ...grpc.CallOption) (*RedeemPickupQRResponse, error)
	ListFulfillmentTasks(ctx context.Context, req *ListFulfillmentTasksRequest, opts ...grpc.CallOption) (*ListFulfillmentTasksResponse, error)
}

type logisticsServiceClient struct {
	cc grpc.ClientConnInterface
}

// NewLogisticsServiceClient wraps a conn so gateway-svc can call logistics-svc RPCs.
func NewLogisticsServiceClient(cc grpc.ClientConnInterface) LogisticsServiceClient {
	return &logisticsServiceClient{cc}
}

func (c *logisticsServiceClient) ShipFulfillmentTask(ctx context.Context, in *ShipFulfillmentTaskRequest, opts ...grpc.CallOption) (*ShipFulfillmentTaskResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShipFulfillmentTaskResponse)
	err := c.cc.Invoke(ctx, LogisticsService_ShipFulfillmentTask_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticsServiceClient) GeneratePickupQR(ctx context.Context, in *GeneratePickupQRRequest, opts ...grpc.CallOption) (*GeneratePickupQRResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GeneratePickupQRResponse)
	err := c.cc.Invoke(ctx, LogisticsService_GeneratePickupQR_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticsServiceClient) RedeemPickupQR(ctx context.Context, in *RedeemPickupQRRequest, opts ...grpc.CallOption) (*RedeemPickupQRResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RedeemPickupQRResponse)
	err := c.cc.Invoke(ctx, LogisticsService_RedeemPickupQR_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticsServiceClient) ListFulfillmentTasks(ctx context.Context, in *ListFulfillmentTasksRequest, opts ...grpc.CallOption) (*ListFulfillmentTasksResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListFulfillmentTasksResponse)
	err := c.cc.Invoke(ctx, LogisticsService_ListFulfillmentTasks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
