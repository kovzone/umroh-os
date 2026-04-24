// booking_channel_grpc_ext.go — hand-written gRPC extension for cross-channel
// seat tracking RPC (BL-BOOK-007).
//
// Extends the BookingService with GetSeatsByChannel.
// Run `make genpb` to regenerate from booking.proto once protoc is available.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Full method name constant for cross-channel seat tracking RPC.
const (
	BookingService_GetSeatsByChannel_FullMethodName = "/pb.BookingService/GetSeatsByChannel"
)

// SeatsByChannelHandler is the server-side interface for GetSeatsByChannel.
type SeatsByChannelHandler interface {
	GetSeatsByChannel(ctx context.Context, req *GetSeatsByChannelRequest) (*GetSeatsByChannelResponse, error)
}

// UnimplementedSeatsByChannelHandler provides a safe default for SeatsByChannelHandler.
type UnimplementedSeatsByChannelHandler struct{}

func (UnimplementedSeatsByChannelHandler) GetSeatsByChannel(_ context.Context, _ *GetSeatsByChannelRequest) (*GetSeatsByChannelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSeatsByChannel not implemented")
}

func _BookingService_GetSeatsByChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSeatsByChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SeatsByChannelHandler).GetSeatsByChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookingService_GetSeatsByChannel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SeatsByChannelHandler).GetSeatsByChannel(ctx, req.(*GetSeatsByChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterBookingServiceServerWithAll registers the BookingService with all
// known RPCs: base + SubmitBooking + GetSeatsByChannel.
func RegisterBookingServiceServerWithAll(s grpc.ServiceRegistrar, srv interface {
	BookingServiceServer
	SubmitBookingHandler
	SeatsByChannelHandler
}) {
	s.RegisterService(&bookingServiceWithAllDesc, srv)
}

// bookingServiceWithAllDesc extends bookingServiceWithSubmitDesc with GetSeatsByChannel.
var bookingServiceWithAllDesc = grpc.ServiceDesc{
	ServiceName: "pb.BookingService",
	HandlerType: (*BookingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Healthz",
			Handler:    _BookingService_Healthz_Handler,
		},
		{
			MethodName: "CreateDraftBooking",
			Handler:    _BookingService_CreateDraftBooking_Handler,
		},
		{
			MethodName: "SubmitBooking",
			Handler:    _BookingService_SubmitBooking_Handler,
		},
		{
			MethodName: "GetSeatsByChannel",
			Handler:    _BookingService_GetSeatsByChannel_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "booking.proto",
}
