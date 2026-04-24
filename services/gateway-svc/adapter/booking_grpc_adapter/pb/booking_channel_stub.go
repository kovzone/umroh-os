// booking_channel_stub.go — gateway-side gRPC client stub for booking-svc
// cross-channel seat tracking RPC (BL-BOOK-007).
//
// Mirrors services/booking-svc/api/grpc_api/pb/booking_channel_messages.go and
// booking_channel_grpc_ext.go. Run `make genpb` to replace with generated code
// once protoc is available.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// Method name constants
// ---------------------------------------------------------------------------

const (
	BookingService_GetSeatsByChannel_FullMethodName = "/pb.BookingService/GetSeatsByChannel"
)

// ---------------------------------------------------------------------------
// GetSeatsByChannel — domain types
// ---------------------------------------------------------------------------

// ChannelSeatRow carries the seat breakdown for a single channel.
type ChannelSeatRow struct {
	Channel      string
	Seats        int32
	BookingCount int32
}

func (x *ChannelSeatRow) GetChannel() string {
	if x == nil {
		return ""
	}
	return x.Channel
}
func (x *ChannelSeatRow) GetSeats() int32 {
	if x == nil {
		return 0
	}
	return x.Seats
}
func (x *ChannelSeatRow) GetBookingCount() int32 {
	if x == nil {
		return 0
	}
	return x.BookingCount
}

// GetSeatsByChannelRequest is the gRPC request for GetSeatsByChannel.
type GetSeatsByChannelRequest struct {
	DepartureId string
}

func (x *GetSeatsByChannelRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}

// GetSeatsByChannelResponse is the gRPC response for GetSeatsByChannel.
type GetSeatsByChannelResponse struct {
	DepartureId string
	ByChannel   []*ChannelSeatRow
}

func (x *GetSeatsByChannelResponse) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *GetSeatsByChannelResponse) GetByChannel() []*ChannelSeatRow {
	if x == nil {
		return nil
	}
	return x.ByChannel
}

// ---------------------------------------------------------------------------
// BookingChannelClient (gRPC client stub)
// ---------------------------------------------------------------------------

// BookingChannelClient is the client API for the cross-channel seat RPC.
type BookingChannelClient interface {
	GetSeatsByChannel(ctx context.Context, in *GetSeatsByChannelRequest, opts ...grpc.CallOption) (*GetSeatsByChannelResponse, error)
}

type bookingChannelClient struct {
	cc grpc.ClientConnInterface
}

// NewBookingChannelClient returns a BookingChannelClient backed by the given conn.
func NewBookingChannelClient(cc grpc.ClientConnInterface) BookingChannelClient {
	return &bookingChannelClient{cc}
}

func (c *bookingChannelClient) GetSeatsByChannel(ctx context.Context, in *GetSeatsByChannelRequest, opts ...grpc.CallOption) (*GetSeatsByChannelResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetSeatsByChannelResponse)
	if err := c.cc.Invoke(ctx, BookingService_GetSeatsByChannel_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
