// reports_grpc_ext.go — hand-written gRPC interface extension for S5-E-01.
//
// Extends FinanceServiceServer / FinanceServiceClient with GetFinanceSummary
// and ListJournalEntries RPCs.
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
	FinanceService_GetFinanceSummary_FullMethodName    = "/pb.FinanceService/GetFinanceSummary"
	FinanceService_ListJournalEntries_FullMethodName   = "/pb.FinanceService/ListJournalEntries"
)

// ---------------------------------------------------------------------------
// Client interfaces
// ---------------------------------------------------------------------------

// FinanceReportsClient is the consumer-side interface for report RPCs.
// Used by gateway-svc → finance-svc.
type FinanceReportsClient interface {
	GetFinanceSummary(ctx context.Context, in *GetFinanceSummaryRequest, opts ...grpc.CallOption) (*GetFinanceSummaryResponse, error)
	ListJournalEntries(ctx context.Context, in *ListJournalEntriesRequest, opts ...grpc.CallOption) (*ListJournalEntriesResponse, error)
}

type financeReportsClient struct {
	cc grpc.ClientConnInterface
}

// NewFinanceReportsClient wraps a conn so gateway-svc can call finance report RPCs.
func NewFinanceReportsClient(cc grpc.ClientConnInterface) FinanceReportsClient {
	return &financeReportsClient{cc}
}

func (c *financeReportsClient) GetFinanceSummary(ctx context.Context, in *GetFinanceSummaryRequest, opts ...grpc.CallOption) (*GetFinanceSummaryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetFinanceSummaryResponse)
	err := c.cc.Invoke(ctx, FinanceService_GetFinanceSummary_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *financeReportsClient) ListJournalEntries(ctx context.Context, in *ListJournalEntriesRequest, opts ...grpc.CallOption) (*ListJournalEntriesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListJournalEntriesResponse)
	err := c.cc.Invoke(ctx, FinanceService_ListJournalEntries_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Server side handler interfaces
// ---------------------------------------------------------------------------

// FinanceReportsHandler is the server-side interface for report RPCs.
type FinanceReportsHandler interface {
	GetFinanceSummary(context.Context, *GetFinanceSummaryRequest) (*GetFinanceSummaryResponse, error)
	ListJournalEntries(context.Context, *ListJournalEntriesRequest) (*ListJournalEntriesResponse, error)
}

// UnimplementedFinanceReportsHandler provides safe defaults.
type UnimplementedFinanceReportsHandler struct{}

func (UnimplementedFinanceReportsHandler) GetFinanceSummary(context.Context, *GetFinanceSummaryRequest) (*GetFinanceSummaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFinanceSummary not implemented")
}

func (UnimplementedFinanceReportsHandler) ListJournalEntries(context.Context, *ListJournalEntriesRequest) (*ListJournalEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListJournalEntries not implemented")
}

// ---------------------------------------------------------------------------
// Handler functions for gRPC service descriptor
// ---------------------------------------------------------------------------

func _FinanceService_GetFinanceSummary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFinanceSummaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceReportsHandler).GetFinanceSummary(ctx, req.(*GetFinanceSummaryRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetFinanceSummary_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_ListJournalEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListJournalEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceReportsHandler).ListJournalEntries(ctx, req.(*ListJournalEntriesRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_ListJournalEntries_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

// FinanceService_GetFinanceSummary_MethodDesc is the gRPC method descriptor.
var FinanceService_GetFinanceSummary_MethodDesc = grpc.MethodDesc{
	MethodName: "GetFinanceSummary",
	Handler:    _FinanceService_GetFinanceSummary_Handler,
}

// FinanceService_ListJournalEntries_MethodDesc is the gRPC method descriptor.
var FinanceService_ListJournalEntries_MethodDesc = grpc.MethodDesc{
	MethodName: "ListJournalEntries",
	Handler:    _FinanceService_ListJournalEntries_Handler,
}

// RegisterFinanceServiceServerWithAllExtensions registers the combined
// FinanceService (generated Healthz + OnPaymentReceived + reports RPCs).
func RegisterFinanceServiceServerWithAllExtensions(s grpc.ServiceRegistrar, srv interface {
	FinanceServiceServer
	OnPaymentReceivedHandler
	FinanceReportsHandler
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
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "finance.proto",
	}
	s.RegisterService(&desc, srv)
}
