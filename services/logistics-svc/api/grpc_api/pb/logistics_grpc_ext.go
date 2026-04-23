// logistics_grpc_ext.go — hand-written gRPC interface extension for S3-E-02.
//
// Extends LogisticsServiceServer / LogisticsServiceClient with OnBookingPaid.
// This file bridges the generated logistics_grpc.pb.go (Healthz only) and the
// hand-written message types in logistics_messages.go.
//
// Run `make genpb` after updating logistics.proto to replace these hand-written
// stubs with generated code.
//
// Context: booking-svc calls OnBookingPaid as a fire-and-forget gRPC call
// after a booking transitions to paid_in_full (ADR-0006: direct gRPC, no
// Temporal).

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	LogisticsService_OnBookingPaid_FullMethodName = "/pb.LogisticsService/OnBookingPaid"
)

// OnBookingPaidClient adds OnBookingPaid to the generated client interface.
// Consumers (e.g. booking-svc adapter) use this directly without casting.
type OnBookingPaidClient interface {
	OnBookingPaid(ctx context.Context, in *OnBookingPaidRequest, opts ...grpc.CallOption) (*OnBookingPaidResponse, error)
}

// logisticsOnBookingPaidClient is the concrete implementation used when
// dialling logistics-svc from booking-svc.
type logisticsOnBookingPaidClient struct {
	cc grpc.ClientConnInterface
}

// NewLogisticsOnBookingPaidClient wraps a conn so booking-svc can call
// OnBookingPaid without regenerating proto files.
func NewLogisticsOnBookingPaidClient(cc grpc.ClientConnInterface) OnBookingPaidClient {
	return &logisticsOnBookingPaidClient{cc}
}

func (c *logisticsOnBookingPaidClient) OnBookingPaid(ctx context.Context, in *OnBookingPaidRequest, opts ...grpc.CallOption) (*OnBookingPaidResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OnBookingPaidResponse)
	err := c.cc.Invoke(ctx, LogisticsService_OnBookingPaid_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Server side — UnimplementedLogisticsServiceServer already embeds Healthz.
// We extend LogisticsServiceServer (the generated interface) by registering
// the OnBookingPaid handler via a separate ServiceDesc so that both methods
// live on the same gRPC service name "pb.LogisticsService".
// ---------------------------------------------------------------------------

// OnBookingPaidHandler is the server-side interface for the new RPC.
type OnBookingPaidHandler interface {
	OnBookingPaid(context.Context, *OnBookingPaidRequest) (*OnBookingPaidResponse, error)
}

// UnimplementedOnBookingPaidHandler provides a safe default for services that
// have not yet implemented OnBookingPaid.
type UnimplementedOnBookingPaidHandler struct{}

func (UnimplementedOnBookingPaidHandler) OnBookingPaid(context.Context, *OnBookingPaidRequest) (*OnBookingPaidResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnBookingPaid not implemented")
}

func _LogisticsService_OnBookingPaid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OnBookingPaidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OnBookingPaidHandler).OnBookingPaid(ctx, req.(*OnBookingPaidRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LogisticsService_OnBookingPaid_FullMethodName,
	}
	return interceptor(ctx, in, info, handler)
}

// LogisticsService_OnBookingPaid_MethodDesc is appended to the service
// descriptor in RegisterLogisticsServiceServerWithExtensions.
var LogisticsService_OnBookingPaid_MethodDesc = grpc.MethodDesc{
	MethodName: "OnBookingPaid",
	Handler:    _LogisticsService_OnBookingPaid_Handler,
}

// RegisterLogisticsServiceServerWithExtensions registers the combined
// LogisticsService (generated Healthz + hand-written OnBookingPaid) on the
// given gRPC server.  Replace with generated RegisterLogisticsServiceServer
// once `make genpb` includes OnBookingPaid.
func RegisterLogisticsServiceServerWithExtensions(s grpc.ServiceRegistrar, srv interface {
	LogisticsServiceServer
	OnBookingPaidHandler
}) {
	desc := grpc.ServiceDesc{
		ServiceName: "pb.LogisticsService",
		HandlerType: (*LogisticsServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Healthz",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(HealthzRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(LogisticsServiceServer).Healthz(ctx, req.(*HealthzRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_Healthz_FullMethodName}
					return interceptor(ctx, in, info, handler)
				},
			},
			LogisticsService_OnBookingPaid_MethodDesc,
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "logistics.proto",
	}
	s.RegisterService(&desc, srv)
}
