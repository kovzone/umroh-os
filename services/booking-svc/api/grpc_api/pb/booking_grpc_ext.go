// booking_grpc_ext.go — extends the generated booking gRPC surface with
// SubmitBooking (S2 / BL-BOOK-005).
//
// Pattern: same as catalog_grpc_ext.go / payment_grpc_ext.go in sibling services.
// When protoc is available, merge this into booking_grpc.pb.go via `make genpb`.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	BookingService_SubmitBooking_FullMethodName = "/pb.BookingService/SubmitBooking"
)

// SubmitBookingHandler is the narrow interface for the SubmitBooking RPC.
// Embed this in the server struct alongside BookingServiceServer.
type SubmitBookingHandler interface {
	SubmitBooking(ctx context.Context, req *SubmitBookingRequest) (*SubmitBookingResponse, error)
}

// UnimplementedSubmitBookingHandler returns Unimplemented for SubmitBooking.
type UnimplementedSubmitBookingHandler struct{}

func (UnimplementedSubmitBookingHandler) SubmitBooking(_ context.Context, _ *SubmitBookingRequest) (*SubmitBookingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitBooking not implemented")
}

func _BookingService_SubmitBooking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitBookingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubmitBookingHandler).SubmitBooking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BookingService_SubmitBooking_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubmitBookingHandler).SubmitBooking(ctx, req.(*SubmitBookingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterBookingServiceServerWithSubmit registers the BookingService with
// both the base RPCs and the SubmitBooking extension.
func RegisterBookingServiceServerWithSubmit(s grpc.ServiceRegistrar, srv interface {
	BookingServiceServer
	SubmitBookingHandler
}) {
	s.RegisterService(&bookingServiceWithSubmitDesc, srv)
}

// bookingServiceWithSubmitDesc extends BookingService_ServiceDesc with SubmitBooking.
var bookingServiceWithSubmitDesc = grpc.ServiceDesc{
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
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "booking.proto",
}
