// grn_grpc_ext.go — hand-written gRPC interface extension for OnGRNReceived RPC
// (BL-FIN-002).
//
// Extends FinanceServiceServer with OnGRNReceived.
// Run `make genpb` after updating finance.proto to replace with generated code.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	FinanceService_OnGRNReceived_FullMethodName = "/pb.FinanceService/OnGRNReceived"
)

// ---------------------------------------------------------------------------
// Client side
// ---------------------------------------------------------------------------

// OnGRNReceivedClient is the consumer-side interface for the OnGRNReceived RPC.
type OnGRNReceivedClient interface {
	OnGRNReceived(ctx context.Context, in *OnGRNReceivedRequest, opts ...grpc.CallOption) (*OnGRNReceivedResponse, error)
}

type onGRNReceivedClient struct {
	cc grpc.ClientConnInterface
}

// NewOnGRNReceivedClient wraps a conn so callers can invoke OnGRNReceived.
func NewOnGRNReceivedClient(cc grpc.ClientConnInterface) OnGRNReceivedClient {
	return &onGRNReceivedClient{cc}
}

func (c *onGRNReceivedClient) OnGRNReceived(ctx context.Context, in *OnGRNReceivedRequest, opts ...grpc.CallOption) (*OnGRNReceivedResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OnGRNReceivedResponse)
	err := c.cc.Invoke(ctx, FinanceService_OnGRNReceived_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Server side
// ---------------------------------------------------------------------------

// OnGRNReceivedHandler is the server-side interface for the OnGRNReceived RPC.
type OnGRNReceivedHandler interface {
	OnGRNReceived(context.Context, *OnGRNReceivedRequest) (*OnGRNReceivedResponse, error)
}

// UnimplementedOnGRNReceivedHandler provides safe defaults.
type UnimplementedOnGRNReceivedHandler struct{}

func (UnimplementedOnGRNReceivedHandler) OnGRNReceived(context.Context, *OnGRNReceivedRequest) (*OnGRNReceivedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OnGRNReceived not implemented")
}

func _FinanceService_OnGRNReceived_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OnGRNReceivedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OnGRNReceivedHandler).OnGRNReceived(ctx, req.(*OnGRNReceivedRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FinanceService_OnGRNReceived_FullMethodName,
	}
	return interceptor(ctx, in, info, handler)
}

// FinanceService_OnGRNReceived_MethodDesc is the gRPC method descriptor.
var FinanceService_OnGRNReceived_MethodDesc = grpc.MethodDesc{
	MethodName: "OnGRNReceived",
	Handler:    _FinanceService_OnGRNReceived_Handler,
}

// RegisterFinanceServiceServerWithGRN registers the complete FinanceService
// including OnGRNReceived.
// This replaces RegisterFinanceServiceServerFull for new deployments.
func RegisterFinanceServiceServerWithGRN(s grpc.ServiceRegistrar, srv interface {
	FinanceServiceServer
	OnPaymentReceivedHandler
	FinanceReportsHandler
	FinanceDepthHandler
	OnGRNReceivedHandler
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
			FinanceService_GetFinanceSummary_MethodDesc,
			FinanceService_ListJournalEntries_MethodDesc,
			FinanceService_RecognizeRevenue_MethodDesc,
			FinanceService_GetPLReport_MethodDesc,
			FinanceService_GetBalanceSheet_MethodDesc,
			FinanceService_OnGRNReceived_MethodDesc,
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "finance.proto",
	}
	s.RegisterService(&desc, srv)
}
