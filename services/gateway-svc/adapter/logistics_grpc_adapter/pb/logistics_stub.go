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
	LogisticsService_ShipFulfillmentTask_FullMethodName = "/pb.LogisticsService/ShipFulfillmentTask"
	LogisticsService_GeneratePickupQR_FullMethodName    = "/pb.LogisticsService/GeneratePickupQR"
	LogisticsService_RedeemPickupQR_FullMethodName      = "/pb.LogisticsService/RedeemPickupQR"
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
// LogisticsServiceClient interface + implementation
// ---------------------------------------------------------------------------

// LogisticsServiceClient is the consumer-side interface for logistics-svc RPCs.
type LogisticsServiceClient interface {
	ShipFulfillmentTask(ctx context.Context, req *ShipFulfillmentTaskRequest, opts ...grpc.CallOption) (*ShipFulfillmentTaskResponse, error)
	GeneratePickupQR(ctx context.Context, req *GeneratePickupQRRequest, opts ...grpc.CallOption) (*GeneratePickupQRResponse, error)
	RedeemPickupQR(ctx context.Context, req *RedeemPickupQRRequest, opts ...grpc.CallOption) (*RedeemPickupQRResponse, error)
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
