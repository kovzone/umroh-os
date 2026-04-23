// correction_grpc_ext.go — gRPC extension for CorrectJournal + DeleteJournalEntry
// (BL-FIN-006).
//
// Adds two RPCs to the FinanceService surface:
//   - CorrectJournal         — posts reversing counter-entry
//   - DeleteJournalEntry     — anti-delete guard (always PermissionDenied)
//
// RegisterFinanceServiceServerWithCorrections must be used in cmd/server.go
// (replaces RegisterFinanceServiceServerWithGRN).

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	FinanceService_CorrectJournal_FullMethodName     = "/pb.FinanceService/CorrectJournal"
	FinanceService_DeleteJournalEntry_FullMethodName = "/pb.FinanceService/DeleteJournalEntry"
)

// ---------------------------------------------------------------------------
// Server-side handler interfaces
// ---------------------------------------------------------------------------

// CorrectionHandler is implemented by the finance-svc gRPC server.
type CorrectionHandler interface {
	CorrectJournal(context.Context, *CorrectJournalRequest) (*CorrectJournalResponse, error)
	DeleteJournalEntry(context.Context, *DeleteJournalEntryRequest) (*DeleteJournalEntryResponse, error)
}

// UnimplementedCorrectionHandler provides safe defaults.
type UnimplementedCorrectionHandler struct{}

func (UnimplementedCorrectionHandler) CorrectJournal(context.Context, *CorrectJournalRequest) (*CorrectJournalResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CorrectJournal not implemented")
}

func (UnimplementedCorrectionHandler) DeleteJournalEntry(context.Context, *DeleteJournalEntryRequest) (*DeleteJournalEntryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteJournalEntry not implemented")
}

// ---------------------------------------------------------------------------
// gRPC handler functions
// ---------------------------------------------------------------------------

func _FinanceService_CorrectJournal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CorrectJournalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CorrectionHandler).CorrectJournal(ctx, req.(*CorrectJournalRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_CorrectJournal_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_DeleteJournalEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteJournalEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CorrectionHandler).DeleteJournalEntry(ctx, req.(*DeleteJournalEntryRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_DeleteJournalEntry_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

var FinanceService_CorrectJournal_MethodDesc = grpc.MethodDesc{
	MethodName: "CorrectJournal",
	Handler:    _FinanceService_CorrectJournal_Handler,
}

var FinanceService_DeleteJournalEntry_MethodDesc = grpc.MethodDesc{
	MethodName: "DeleteJournalEntry",
	Handler:    _FinanceService_DeleteJournalEntry_Handler,
}

// ---------------------------------------------------------------------------
// Client-side stub
// ---------------------------------------------------------------------------

// FinanceCorrectionClient is the gateway-side consumer interface.
type FinanceCorrectionClient interface {
	CorrectJournal(ctx context.Context, in *CorrectJournalRequest, opts ...grpc.CallOption) (*CorrectJournalResponse, error)
}

type financeCorrectionClient struct {
	cc grpc.ClientConnInterface
}

func NewFinanceCorrectionClient(cc grpc.ClientConnInterface) FinanceCorrectionClient {
	return &financeCorrectionClient{cc}
}

func (c *financeCorrectionClient) CorrectJournal(ctx context.Context, in *CorrectJournalRequest, opts ...grpc.CallOption) (*CorrectJournalResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CorrectJournalResponse)
	err := c.cc.Invoke(ctx, FinanceService_CorrectJournal_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// RegisterFinanceServiceServerWithCorrections
// ---------------------------------------------------------------------------

// RegisterFinanceServiceServerWithCorrections registers the complete
// FinanceService including CorrectJournal and DeleteJournalEntry (BL-FIN-006).
// This replaces RegisterFinanceServiceServerWithGRN in cmd/server.go.
func RegisterFinanceServiceServerWithCorrections(s grpc.ServiceRegistrar, srv interface {
	FinanceServiceServer
	OnPaymentReceivedHandler
	FinanceReportsHandler
	FinanceDepthHandler
	OnGRNReceivedHandler
	CorrectionHandler
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
			FinanceService_CorrectJournal_MethodDesc,
			FinanceService_DeleteJournalEntry_MethodDesc,
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "finance.proto",
	}
	s.RegisterService(&desc, srv)
}
