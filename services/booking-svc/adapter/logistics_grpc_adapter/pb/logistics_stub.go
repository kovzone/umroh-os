// logistics_stub.go — minimal hand-written gRPC client stub for booking-svc
// to call logistics-svc.OnBookingPaid (S3-E-02).
//
// This mirrors the subset of logistics-svc's pb types that booking-svc needs
// as a consumer. Run `make genpb` with a shared proto path to replace with
// generated code.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

const (
	LogisticsService_OnBookingPaid_FullMethodName = "/pb.LogisticsService/OnBookingPaid"
)

// OnBookingPaidRequest is the request message for LogisticsService.OnBookingPaid.
// JamaahIds is required (min 1) per §S3-J-02 contract.
type OnBookingPaidRequest struct {
	BookingId   string
	DepartureId string
	JamaahIds   []string // at least 1; ULIDs of jamaah on this booking
}

// OnBookingPaidResponse is the response message from LogisticsService.OnBookingPaid.
type OnBookingPaidResponse struct {
	TaskId string
	Status string
}

func (x *OnBookingPaidResponse) GetTaskId() string {
	if x == nil {
		return ""
	}
	return x.TaskId
}
func (x *OnBookingPaidResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

// LogisticsServiceClient is the consumer-side interface for the subset of
// LogisticsService used by booking-svc.
type LogisticsServiceClient interface {
	OnBookingPaid(ctx context.Context, in *OnBookingPaidRequest, opts ...grpc.CallOption) (*OnBookingPaidResponse, error)
}

type logisticsServiceClient struct {
	cc grpc.ClientConnInterface
}

// NewLogisticsServiceClient creates a client for logistics-svc from a dialled conn.
func NewLogisticsServiceClient(cc grpc.ClientConnInterface) LogisticsServiceClient {
	return &logisticsServiceClient{cc}
}

func (c *logisticsServiceClient) OnBookingPaid(ctx context.Context, in *OnBookingPaidRequest, opts ...grpc.CallOption) (*OnBookingPaidResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OnBookingPaidResponse)
	err := c.cc.Invoke(ctx, LogisticsService_OnBookingPaid_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
