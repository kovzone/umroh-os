// finance_grpc_ext.go — hand-written gRPC interface extension for S3-E-03.
//
// Extends FinanceServiceServer / FinanceServiceClient with OnPaymentReceived.
// This file bridges the generated finance_grpc.pb.go (Healthz only) and the
// hand-written message types in finance_messages.go.
//
// Run `make genpb` after updating finance.proto to replace these hand-written
// stubs with generated code.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	FinanceService_OnPaymentReceived_FullMethodName = "/pb.FinanceService/OnPaymentReceived"
)

// OnPaymentReceivedClient adds OnPaymentReceived to the generated client interface.
type OnPaymentReceivedClient interface {
	OnPaymentReceived(ctx context.Context, in *OnPaymentReceivedRequest, opts ...grpc.CallOption) (*OnPaymentReceivedResponse, error)
}

// financeOnPaymentReceivedClient is the concrete implementation used when
// dialling finance-svc from booking-svc.
type financeOnPaymentReceivedClient struct {
	cc grpc.ClientConnInterface
}

// NewFinanceOnPaymentReceivedClient wraps a conn so booking-svc can call
// OnPaymentReceived without regenerating proto files.
func NewFinanceOnPaymentReceivedClient(cc grpc.ClientConnInterface) OnPaymentReceivedClient {
	return &financeOnPaymentReceivedClient{cc}
}

func (c *financeOnPaymentReceivedClient) OnPaymentReceived(ctx context.Context, in *OnPaymentReceivedRequest, opts ...grpc.CallOption) (*OnPaymentReceivedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OnPaymentReceivedResponse)
	err := c.cc.Invoke(ctx, FinanceService_OnPaymentReceived_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Server side
// ---------------------------------------------------------------------------

// OnPaymentReceivedHandler is the server-side interface for the new RPC.
type OnPaymentReceivedHandler interface {
	OnPaymentReceived(context.Context, *OnPaymentReceivedRequest) (*OnPaymentReceivedResponse, error)
}

// UnimplementedOnPaymentReceivedHandler provides a safe default for services
// that have not yet implemented OnPaymentReceived.
type UnimplementedOnPaymentReceivedHandler struct{}

func (UnimplementedOnPaymentReceivedHandler) OnPaymentReceived(context.Context, *OnPaymentReceivedRequest) (*OnPaymentReceivedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnPaymentReceived not implemented")
}

func _FinanceService_OnPaymentReceived_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OnPaymentReceivedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OnPaymentReceivedHandler).OnPaymentReceived(ctx, req.(*OnPaymentReceivedRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FinanceService_OnPaymentReceived_FullMethodName,
	}
	return interceptor(ctx, in, info, handler)
}

// FinanceService_OnPaymentReceived_MethodDesc is appended to the service
// descriptor in RegisterFinanceServiceServerWithExtensions.
var FinanceService_OnPaymentReceived_MethodDesc = grpc.MethodDesc{
	MethodName: "OnPaymentReceived",
	Handler:    _FinanceService_OnPaymentReceived_Handler,
}

// RegisterFinanceServiceServerWithExtensions registers the combined
// FinanceService (generated Healthz + hand-written OnPaymentReceived).
func RegisterFinanceServiceServerWithExtensions(s grpc.ServiceRegistrar, srv interface {
	FinanceServiceServer
	OnPaymentReceivedHandler
}) {
	desc := grpc.ServiceDesc{
		ServiceName: "pb.FinanceService",
		HandlerType: (*FinanceServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Healthz",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(HealthzRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(FinanceServiceServer).Healthz(ctx, req.(*HealthzRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_Healthz_FullMethodName}
					return interceptor(ctx, in, info, handler)
				},
			},
			FinanceService_OnPaymentReceived_MethodDesc,
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "finance.proto",
	}
	s.RegisterService(&desc, srv)
}
