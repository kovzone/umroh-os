// pl_grpc_ext.go — hand-written gRPC interface extension for finance depth
// RPCs: RecognizeRevenue, GetPLReport, GetBalanceSheet (Wave 1B).
//
// Follows the same pattern as finance_grpc_ext.go and reports_grpc_ext.go.
// Run `make genpb` after updating finance.proto to replace with generated code.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	FinanceService_RecognizeRevenue_FullMethodName = "/pb.FinanceService/RecognizeRevenue"
	FinanceService_GetPLReport_FullMethodName      = "/pb.FinanceService/GetPLReport"
	FinanceService_GetBalanceSheet_FullMethodName  = "/pb.FinanceService/GetBalanceSheet"
)

// ---------------------------------------------------------------------------
// Client interfaces
// ---------------------------------------------------------------------------

// FinanceDepthClient is the consumer-side interface for finance-depth RPCs.
// Used by gateway-svc → finance-svc.
type FinanceDepthClient interface {
	RecognizeRevenue(ctx context.Context, in *RecognizeRevenueRequest, opts ...grpc.CallOption) (*RecognizeRevenueResponse, error)
	GetPLReport(ctx context.Context, in *GetPLReportRequest, opts ...grpc.CallOption) (*PLReportProto, error)
	GetBalanceSheet(ctx context.Context, in *GetBalanceSheetRequest, opts ...grpc.CallOption) (*BalanceSheetProto, error)
}

type financeDepthClient struct {
	cc grpc.ClientConnInterface
}

// NewFinanceDepthClient wraps a conn so gateway-svc can call finance depth RPCs.
func NewFinanceDepthClient(cc grpc.ClientConnInterface) FinanceDepthClient {
	return &financeDepthClient{cc}
}

func (c *financeDepthClient) RecognizeRevenue(ctx context.Context, in *RecognizeRevenueRequest, opts ...grpc.CallOption) (*RecognizeRevenueResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecognizeRevenueResponse)
	err := c.cc.Invoke(ctx, FinanceService_RecognizeRevenue_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *financeDepthClient) GetPLReport(ctx context.Context, in *GetPLReportRequest, opts ...grpc.CallOption) (*PLReportProto, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PLReportProto)
	err := c.cc.Invoke(ctx, FinanceService_GetPLReport_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *financeDepthClient) GetBalanceSheet(ctx context.Context, in *GetBalanceSheetRequest, opts ...grpc.CallOption) (*BalanceSheetProto, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BalanceSheetProto)
	err := c.cc.Invoke(ctx, FinanceService_GetBalanceSheet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Server side handler interfaces
// ---------------------------------------------------------------------------

// FinanceDepthHandler is the server-side interface for finance-depth RPCs.
type FinanceDepthHandler interface {
	RecognizeRevenue(context.Context, *RecognizeRevenueRequest) (*RecognizeRevenueResponse, error)
	GetPLReport(context.Context, *GetPLReportRequest) (*PLReportProto, error)
	GetBalanceSheet(context.Context, *GetBalanceSheetRequest) (*BalanceSheetProto, error)
}

// UnimplementedFinanceDepthHandler provides safe defaults.
type UnimplementedFinanceDepthHandler struct{}

func (UnimplementedFinanceDepthHandler) RecognizeRevenue(context.Context, *RecognizeRevenueRequest) (*RecognizeRevenueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecognizeRevenue not implemented")
}

func (UnimplementedFinanceDepthHandler) GetPLReport(context.Context, *GetPLReportRequest) (*PLReportProto, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPLReport not implemented")
}

func (UnimplementedFinanceDepthHandler) GetBalanceSheet(context.Context, *GetBalanceSheetRequest) (*BalanceSheetProto, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBalanceSheet not implemented")
}

// ---------------------------------------------------------------------------
// Handler functions for gRPC service descriptor
// ---------------------------------------------------------------------------

func _FinanceService_RecognizeRevenue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecognizeRevenueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceDepthHandler).RecognizeRevenue(ctx, req.(*RecognizeRevenueRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_RecognizeRevenue_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetPLReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPLReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceDepthHandler).GetPLReport(ctx, req.(*GetPLReportRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetPLReport_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetBalanceSheet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBalanceSheetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceDepthHandler).GetBalanceSheet(ctx, req.(*GetBalanceSheetRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetBalanceSheet_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

// gRPC method descriptors.
var FinanceService_RecognizeRevenue_MethodDesc = grpc.MethodDesc{
	MethodName: "RecognizeRevenue",
	Handler:    _FinanceService_RecognizeRevenue_Handler,
}

var FinanceService_GetPLReport_MethodDesc = grpc.MethodDesc{
	MethodName: "GetPLReport",
	Handler:    _FinanceService_GetPLReport_Handler,
}

var FinanceService_GetBalanceSheet_MethodDesc = grpc.MethodDesc{
	MethodName: "GetBalanceSheet",
	Handler:    _FinanceService_GetBalanceSheet_Handler,
}

// RegisterFinanceServiceServerFull registers the complete FinanceService:
// Healthz + OnPaymentReceived + reports RPCs + finance-depth RPCs.
// This replaces RegisterFinanceServiceServerWithAllExtensions for new deployments.
func RegisterFinanceServiceServerFull(s grpc.ServiceRegistrar, srv interface {
	FinanceServiceServer
	OnPaymentReceivedHandler
	FinanceReportsHandler
	FinanceDepthHandler
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
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "finance.proto",
	}
	s.RegisterService(&desc, srv)
}
