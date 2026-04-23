// payment_grpc_ext.go — hand-written gRPC service extension for payment-svc.
//
// Adds the ReissuePaymentLink RPC (BL-PAY-020) on top of the generated
// PaymentServiceServer interface using the same *_grpc_ext.go pattern used by
// other services in this repo (ops-svc, finance-svc, logistics-svc, etc.).
//
// RegisterPaymentServiceServerFull replaces RegisterPaymentServiceServer in
// cmd/server.go. It builds a single grpc.ServiceDesc that contains:
//   - All original methods from PaymentService_ServiceDesc (Healthz,
//     IssueVirtualAccount, ProcessWebhook, StartRefund)
//   - ReissuePaymentLink (added here)

package pb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	PaymentService_ReissuePaymentLink_FullMethodName = "/pb.PaymentService/ReissuePaymentLink"
)

// PaymentServiceWithReissueServer is implemented by the *grpc_api.Server in this service.
// It extends PaymentServiceServer with the CS-facing ReissuePaymentLink RPC.
type PaymentServiceWithReissueServer interface {
	PaymentServiceServer
	ReissuePaymentLink(context.Context, *ReissuePaymentLinkRequest) (*ReissuePaymentLinkResponse, error)
}

// UnimplementedPaymentServiceWithReissueServer provides a default no-op
// implementation so embedding structs only need to override what they implement.
type UnimplementedPaymentServiceWithReissueServer struct {
	UnimplementedPaymentServiceServer
}

func (UnimplementedPaymentServiceWithReissueServer) ReissuePaymentLink(context.Context, *ReissuePaymentLinkRequest) (*ReissuePaymentLinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReissuePaymentLink not implemented")
}

// _PaymentService_ReissuePaymentLink_Handler is the gRPC handler function for
// the ReissuePaymentLink RPC.
func _PaymentService_ReissuePaymentLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReissuePaymentLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServiceWithReissueServer).ReissuePaymentLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PaymentService_ReissuePaymentLink_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServiceWithReissueServer).ReissuePaymentLink(ctx, req.(*ReissuePaymentLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// paymentServiceFullDesc combines the original PaymentService_ServiceDesc
// methods with the new ReissuePaymentLink method.
var paymentServiceFullDesc = grpc.ServiceDesc{
	ServiceName: "pb.PaymentService",
	HandlerType: (*PaymentServiceWithReissueServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Healthz",
			Handler:    _PaymentService_Healthz_Handler,
		},
		{
			MethodName: "IssueVirtualAccount",
			Handler:    _PaymentService_IssueVirtualAccount_Handler,
		},
		{
			MethodName: "ProcessWebhook",
			Handler:    _PaymentService_ProcessWebhook_Handler,
		},
		{
			MethodName: "StartRefund",
			Handler:    _PaymentService_StartRefund_Handler,
		},
		{
			MethodName: "ReissuePaymentLink",
			Handler:    _PaymentService_ReissuePaymentLink_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "payment.proto",
}

// RegisterPaymentServiceServerFull registers the extended PaymentService
// (with ReissuePaymentLink) on the given gRPC server.
// Use this instead of RegisterPaymentServiceServer in cmd/server.go.
func RegisterPaymentServiceServerFull(s grpc.ServiceRegistrar, srv PaymentServiceWithReissueServer) {
	s.RegisterService(&paymentServiceFullDesc, srv)
}
