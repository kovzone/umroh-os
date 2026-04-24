// logistics_depth_grpc_ext.go — gRPC dispatch handlers and registration for
// Wave 5 depth RPCs (BL-LOG-013..029).
//
// RegisterLogisticsServiceServerWithDepth replaces
// RegisterLogisticsServiceServerWithExtensions and includes ALL methods:
//   - Healthz (generated)
//   - OnBookingPaid, ShipFulfillmentTask, GeneratePickupQR, RedeemPickupQR,
//     ListFulfillmentTasks (Wave 3)
//   - CreatePurchaseRequest, ApprovePurchaseRequest, RecordGRNWithQC,
//     CreateKitAssembly (Wave 6 / BL-LOG-010..012)
//   - Wave 5 depth: 20 new RPCs (BL-LOG-013..029)

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// Full method name constants
// ---------------------------------------------------------------------------

const (
	LogisticsService_CreatePurchaseRequest_FullMethodName  = "/pb.LogisticsService/CreatePurchaseRequest"
	LogisticsService_ApprovePurchaseRequest_FullMethodName = "/pb.LogisticsService/ApprovePurchaseRequest"
	LogisticsService_RecordGRNWithQC_FullMethodName        = "/pb.LogisticsService/RecordGRNWithQC"
	LogisticsService_CreateKitAssembly_FullMethodName      = "/pb.LogisticsService/CreateKitAssembly"

	LogisticsService_ListPurchaseRequests_FullMethodName   = "/pb.LogisticsService/ListPurchaseRequests"
	LogisticsService_GetBudgetSyncStatus_FullMethodName    = "/pb.LogisticsService/GetBudgetSyncStatus"
	LogisticsService_GetTieredApprovals_FullMethodName     = "/pb.LogisticsService/GetTieredApprovals"
	LogisticsService_DecideTieredApproval_FullMethodName   = "/pb.LogisticsService/DecideTieredApproval"
	LogisticsService_AutoSelectVendor_FullMethodName       = "/pb.LogisticsService/AutoSelectVendor"
	LogisticsService_RecordPartialGRN_FullMethodName       = "/pb.LogisticsService/RecordPartialGRN"
	LogisticsService_ReverseGRN_FullMethodName             = "/pb.LogisticsService/ReverseGRN"
	LogisticsService_GenerateBarcode_FullMethodName        = "/pb.LogisticsService/GenerateBarcode"
	LogisticsService_PrintSKULabel_FullMethodName          = "/pb.LogisticsService/PrintSKULabel"
	LogisticsService_CreateWarehouse_FullMethodName        = "/pb.LogisticsService/CreateWarehouse"
	LogisticsService_TransferStock_FullMethodName          = "/pb.LogisticsService/TransferStock"
	LogisticsService_GetStockAlerts_FullMethodName         = "/pb.LogisticsService/GetStockAlerts"
	LogisticsService_SetReorderLevel_FullMethodName        = "/pb.LogisticsService/SetReorderLevel"
	LogisticsService_StartStocktake_FullMethodName         = "/pb.LogisticsService/StartStocktake"
	LogisticsService_RecordStocktakeCount_FullMethodName   = "/pb.LogisticsService/RecordStocktakeCount"
	LogisticsService_FinalizeStocktake_FullMethodName      = "/pb.LogisticsService/FinalizeStocktake"
	LogisticsService_SyncFulfillmentSizes_FullMethodName   = "/pb.LogisticsService/SyncFulfillmentSizes"
	LogisticsService_RecordCourierTracking_FullMethodName  = "/pb.LogisticsService/RecordCourierTracking"
	LogisticsService_RecordReturn_FullMethodName           = "/pb.LogisticsService/RecordReturn"
	LogisticsService_ProcessExchange_FullMethodName        = "/pb.LogisticsService/ProcessExchange"
)

// ---------------------------------------------------------------------------
// Procurement handlers (BL-LOG-010..012) — previously unregistered
// ---------------------------------------------------------------------------

type CreatePurchaseRequestHandler interface {
	CreatePurchaseRequest(context.Context, *CreatePurchaseRequestRequest) (*CreatePurchaseRequestResponse, error)
}

type UnimplementedCreatePurchaseRequestHandler struct{}

func (UnimplementedCreatePurchaseRequestHandler) CreatePurchaseRequest(context.Context, *CreatePurchaseRequestRequest) (*CreatePurchaseRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePurchaseRequest not implemented")
}

func _LogisticsService_CreatePurchaseRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePurchaseRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CreatePurchaseRequestHandler).CreatePurchaseRequest(ctx, req.(*CreatePurchaseRequestRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_CreatePurchaseRequest_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

type ApprovePurchaseRequestHandler interface {
	ApprovePurchaseRequest(context.Context, *ApprovePurchaseRequestRequest) (*ApprovePurchaseRequestResponse, error)
}

type UnimplementedApprovePurchaseRequestHandler struct{}

func (UnimplementedApprovePurchaseRequestHandler) ApprovePurchaseRequest(context.Context, *ApprovePurchaseRequestRequest) (*ApprovePurchaseRequestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApprovePurchaseRequest not implemented")
}

func _LogisticsService_ApprovePurchaseRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApprovePurchaseRequestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApprovePurchaseRequestHandler).ApprovePurchaseRequest(ctx, req.(*ApprovePurchaseRequestRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_ApprovePurchaseRequest_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

type RecordGRNWithQCHandler interface {
	RecordGRNWithQC(context.Context, *RecordGRNWithQCRequest) (*RecordGRNWithQCResponse, error)
}

type UnimplementedRecordGRNWithQCHandler struct{}

func (UnimplementedRecordGRNWithQCHandler) RecordGRNWithQC(context.Context, *RecordGRNWithQCRequest) (*RecordGRNWithQCResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordGRNWithQC not implemented")
}

func _LogisticsService_RecordGRNWithQC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordGRNWithQCRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordGRNWithQCHandler).RecordGRNWithQC(ctx, req.(*RecordGRNWithQCRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_RecordGRNWithQC_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

type CreateKitAssemblyHandler interface {
	CreateKitAssembly(context.Context, *CreateKitAssemblyRequest) (*CreateKitAssemblyResponse, error)
}

type UnimplementedCreateKitAssemblyHandler struct{}

func (UnimplementedCreateKitAssemblyHandler) CreateKitAssembly(context.Context, *CreateKitAssemblyRequest) (*CreateKitAssemblyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateKitAssembly not implemented")
}

func _LogisticsService_CreateKitAssembly_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateKitAssemblyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CreateKitAssemblyHandler).CreateKitAssembly(ctx, req.(*CreateKitAssemblyRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_CreateKitAssembly_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

// ---------------------------------------------------------------------------
// LogisticsDepthHandler — Wave 5 depth (BL-LOG-013..029), 20 new RPCs
// ---------------------------------------------------------------------------

// LogisticsDepthHandler is the server-side interface for all Wave 5 depth RPCs.
type LogisticsDepthHandler interface {
	ListPurchaseRequests(context.Context, *ListPurchaseRequestsRequest) (*ListPurchaseRequestsResponse, error)
	GetBudgetSyncStatus(context.Context, *GetBudgetSyncStatusRequest) (*GetBudgetSyncStatusResponse, error)
	GetTieredApprovals(context.Context, *GetTieredApprovalsRequest) (*GetTieredApprovalsResponse, error)
	DecideTieredApproval(context.Context, *DecideTieredApprovalRequest) (*DecideTieredApprovalResponse, error)
	AutoSelectVendor(context.Context, *AutoSelectVendorRequest) (*AutoSelectVendorResponse, error)
	RecordPartialGRN(context.Context, *RecordPartialGRNRequest) (*RecordPartialGRNResponse, error)
	ReverseGRN(context.Context, *ReverseGRNRequest) (*ReverseGRNResponse, error)
	GenerateBarcode(context.Context, *GenerateBarcodeRequest) (*GenerateBarcodeResponse, error)
	PrintSKULabel(context.Context, *PrintSKULabelRequest) (*PrintSKULabelResponse, error)
	CreateWarehouse(context.Context, *CreateWarehouseRequest) (*CreateWarehouseResponse, error)
	TransferStock(context.Context, *TransferStockRequest) (*TransferStockResponse, error)
	GetStockAlerts(context.Context, *GetStockAlertsRequest) (*GetStockAlertsResponse, error)
	SetReorderLevel(context.Context, *SetReorderLevelRequest) (*SetReorderLevelResponse, error)
	StartStocktake(context.Context, *StartStocktakeRequest) (*StartStocktakeResponse, error)
	RecordStocktakeCount(context.Context, *RecordStocktakeCountRequest) (*RecordStocktakeCountResponse, error)
	FinalizeStocktake(context.Context, *FinalizeStocktakeRequest) (*FinalizeStocktakeResponse, error)
	SyncFulfillmentSizes(context.Context, *SyncFulfillmentSizesRequest) (*SyncFulfillmentSizesResponse, error)
	RecordCourierTracking(context.Context, *RecordCourierTrackingRequest) (*RecordCourierTrackingResponse, error)
	RecordReturn(context.Context, *RecordReturnRequest) (*RecordReturnResponse, error)
	ProcessExchange(context.Context, *ProcessExchangeRequest) (*ProcessExchangeResponse, error)
}

// UnimplementedLogisticsDepthHandler provides safe defaults for all 20 Wave 5 RPCs.
type UnimplementedLogisticsDepthHandler struct{}

func (UnimplementedLogisticsDepthHandler) ListPurchaseRequests(context.Context, *ListPurchaseRequestsRequest) (*ListPurchaseRequestsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPurchaseRequests not implemented")
}
func (UnimplementedLogisticsDepthHandler) GetBudgetSyncStatus(context.Context, *GetBudgetSyncStatusRequest) (*GetBudgetSyncStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBudgetSyncStatus not implemented")
}
func (UnimplementedLogisticsDepthHandler) GetTieredApprovals(context.Context, *GetTieredApprovalsRequest) (*GetTieredApprovalsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTieredApprovals not implemented")
}
func (UnimplementedLogisticsDepthHandler) DecideTieredApproval(context.Context, *DecideTieredApprovalRequest) (*DecideTieredApprovalResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DecideTieredApproval not implemented")
}
func (UnimplementedLogisticsDepthHandler) AutoSelectVendor(context.Context, *AutoSelectVendorRequest) (*AutoSelectVendorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AutoSelectVendor not implemented")
}
func (UnimplementedLogisticsDepthHandler) RecordPartialGRN(context.Context, *RecordPartialGRNRequest) (*RecordPartialGRNResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordPartialGRN not implemented")
}
func (UnimplementedLogisticsDepthHandler) ReverseGRN(context.Context, *ReverseGRNRequest) (*ReverseGRNResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReverseGRN not implemented")
}
func (UnimplementedLogisticsDepthHandler) GenerateBarcode(context.Context, *GenerateBarcodeRequest) (*GenerateBarcodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateBarcode not implemented")
}
func (UnimplementedLogisticsDepthHandler) PrintSKULabel(context.Context, *PrintSKULabelRequest) (*PrintSKULabelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrintSKULabel not implemented")
}
func (UnimplementedLogisticsDepthHandler) CreateWarehouse(context.Context, *CreateWarehouseRequest) (*CreateWarehouseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateWarehouse not implemented")
}
func (UnimplementedLogisticsDepthHandler) TransferStock(context.Context, *TransferStockRequest) (*TransferStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransferStock not implemented")
}
func (UnimplementedLogisticsDepthHandler) GetStockAlerts(context.Context, *GetStockAlertsRequest) (*GetStockAlertsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStockAlerts not implemented")
}
func (UnimplementedLogisticsDepthHandler) SetReorderLevel(context.Context, *SetReorderLevelRequest) (*SetReorderLevelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetReorderLevel not implemented")
}
func (UnimplementedLogisticsDepthHandler) StartStocktake(context.Context, *StartStocktakeRequest) (*StartStocktakeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartStocktake not implemented")
}
func (UnimplementedLogisticsDepthHandler) RecordStocktakeCount(context.Context, *RecordStocktakeCountRequest) (*RecordStocktakeCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordStocktakeCount not implemented")
}
func (UnimplementedLogisticsDepthHandler) FinalizeStocktake(context.Context, *FinalizeStocktakeRequest) (*FinalizeStocktakeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FinalizeStocktake not implemented")
}
func (UnimplementedLogisticsDepthHandler) SyncFulfillmentSizes(context.Context, *SyncFulfillmentSizesRequest) (*SyncFulfillmentSizesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncFulfillmentSizes not implemented")
}
func (UnimplementedLogisticsDepthHandler) RecordCourierTracking(context.Context, *RecordCourierTrackingRequest) (*RecordCourierTrackingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordCourierTracking not implemented")
}
func (UnimplementedLogisticsDepthHandler) RecordReturn(context.Context, *RecordReturnRequest) (*RecordReturnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordReturn not implemented")
}
func (UnimplementedLogisticsDepthHandler) ProcessExchange(context.Context, *ProcessExchangeRequest) (*ProcessExchangeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcessExchange not implemented")
}

// ---------------------------------------------------------------------------
// Dispatch functions — Wave 5 depth
// ---------------------------------------------------------------------------

func _LogisticsService_ListPurchaseRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPurchaseRequestsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).ListPurchaseRequests(ctx, req.(*ListPurchaseRequestsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_ListPurchaseRequests_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_GetBudgetSyncStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBudgetSyncStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).GetBudgetSyncStatus(ctx, req.(*GetBudgetSyncStatusRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_GetBudgetSyncStatus_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_GetTieredApprovals_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTieredApprovalsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).GetTieredApprovals(ctx, req.(*GetTieredApprovalsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_GetTieredApprovals_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_DecideTieredApproval_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DecideTieredApprovalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).DecideTieredApproval(ctx, req.(*DecideTieredApprovalRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_DecideTieredApproval_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_AutoSelectVendor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AutoSelectVendorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).AutoSelectVendor(ctx, req.(*AutoSelectVendorRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_AutoSelectVendor_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_RecordPartialGRN_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordPartialGRNRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).RecordPartialGRN(ctx, req.(*RecordPartialGRNRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_RecordPartialGRN_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_ReverseGRN_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReverseGRNRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).ReverseGRN(ctx, req.(*ReverseGRNRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_ReverseGRN_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_GenerateBarcode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateBarcodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).GenerateBarcode(ctx, req.(*GenerateBarcodeRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_GenerateBarcode_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_PrintSKULabel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrintSKULabelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).PrintSKULabel(ctx, req.(*PrintSKULabelRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_PrintSKULabel_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_CreateWarehouse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateWarehouseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).CreateWarehouse(ctx, req.(*CreateWarehouseRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_CreateWarehouse_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_TransferStock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferStockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).TransferStock(ctx, req.(*TransferStockRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_TransferStock_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_GetStockAlerts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStockAlertsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).GetStockAlerts(ctx, req.(*GetStockAlertsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_GetStockAlerts_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_SetReorderLevel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetReorderLevelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).SetReorderLevel(ctx, req.(*SetReorderLevelRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_SetReorderLevel_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_StartStocktake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartStocktakeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).StartStocktake(ctx, req.(*StartStocktakeRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_StartStocktake_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_RecordStocktakeCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordStocktakeCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).RecordStocktakeCount(ctx, req.(*RecordStocktakeCountRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_RecordStocktakeCount_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_FinalizeStocktake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FinalizeStocktakeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).FinalizeStocktake(ctx, req.(*FinalizeStocktakeRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_FinalizeStocktake_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_SyncFulfillmentSizes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncFulfillmentSizesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).SyncFulfillmentSizes(ctx, req.(*SyncFulfillmentSizesRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_SyncFulfillmentSizes_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_RecordCourierTracking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordCourierTrackingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).RecordCourierTracking(ctx, req.(*RecordCourierTrackingRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_RecordCourierTracking_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_RecordReturn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordReturnRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).RecordReturn(ctx, req.(*RecordReturnRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_RecordReturn_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _LogisticsService_ProcessExchange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessExchangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticsDepthHandler).ProcessExchange(ctx, req.(*ProcessExchangeRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: LogisticsService_ProcessExchange_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

// ---------------------------------------------------------------------------
// Method descriptors — Wave 5 depth
// ---------------------------------------------------------------------------

var LogisticsService_ListPurchaseRequests_MethodDesc = grpc.MethodDesc{MethodName: "ListPurchaseRequests", Handler: _LogisticsService_ListPurchaseRequests_Handler}
var LogisticsService_GetBudgetSyncStatus_MethodDesc = grpc.MethodDesc{MethodName: "GetBudgetSyncStatus", Handler: _LogisticsService_GetBudgetSyncStatus_Handler}
var LogisticsService_GetTieredApprovals_MethodDesc = grpc.MethodDesc{MethodName: "GetTieredApprovals", Handler: _LogisticsService_GetTieredApprovals_Handler}
var LogisticsService_DecideTieredApproval_MethodDesc = grpc.MethodDesc{MethodName: "DecideTieredApproval", Handler: _LogisticsService_DecideTieredApproval_Handler}
var LogisticsService_AutoSelectVendor_MethodDesc = grpc.MethodDesc{MethodName: "AutoSelectVendor", Handler: _LogisticsService_AutoSelectVendor_Handler}
var LogisticsService_RecordPartialGRN_MethodDesc = grpc.MethodDesc{MethodName: "RecordPartialGRN", Handler: _LogisticsService_RecordPartialGRN_Handler}
var LogisticsService_ReverseGRN_MethodDesc = grpc.MethodDesc{MethodName: "ReverseGRN", Handler: _LogisticsService_ReverseGRN_Handler}
var LogisticsService_GenerateBarcode_MethodDesc = grpc.MethodDesc{MethodName: "GenerateBarcode", Handler: _LogisticsService_GenerateBarcode_Handler}
var LogisticsService_PrintSKULabel_MethodDesc = grpc.MethodDesc{MethodName: "PrintSKULabel", Handler: _LogisticsService_PrintSKULabel_Handler}
var LogisticsService_CreateWarehouse_MethodDesc = grpc.MethodDesc{MethodName: "CreateWarehouse", Handler: _LogisticsService_CreateWarehouse_Handler}
var LogisticsService_TransferStock_MethodDesc = grpc.MethodDesc{MethodName: "TransferStock", Handler: _LogisticsService_TransferStock_Handler}
var LogisticsService_GetStockAlerts_MethodDesc = grpc.MethodDesc{MethodName: "GetStockAlerts", Handler: _LogisticsService_GetStockAlerts_Handler}
var LogisticsService_SetReorderLevel_MethodDesc = grpc.MethodDesc{MethodName: "SetReorderLevel", Handler: _LogisticsService_SetReorderLevel_Handler}
var LogisticsService_StartStocktake_MethodDesc = grpc.MethodDesc{MethodName: "StartStocktake", Handler: _LogisticsService_StartStocktake_Handler}
var LogisticsService_RecordStocktakeCount_MethodDesc = grpc.MethodDesc{MethodName: "RecordStocktakeCount", Handler: _LogisticsService_RecordStocktakeCount_Handler}
var LogisticsService_FinalizeStocktake_MethodDesc = grpc.MethodDesc{MethodName: "FinalizeStocktake", Handler: _LogisticsService_FinalizeStocktake_Handler}
var LogisticsService_SyncFulfillmentSizes_MethodDesc = grpc.MethodDesc{MethodName: "SyncFulfillmentSizes", Handler: _LogisticsService_SyncFulfillmentSizes_Handler}
var LogisticsService_RecordCourierTracking_MethodDesc = grpc.MethodDesc{MethodName: "RecordCourierTracking", Handler: _LogisticsService_RecordCourierTracking_Handler}
var LogisticsService_RecordReturn_MethodDesc = grpc.MethodDesc{MethodName: "RecordReturn", Handler: _LogisticsService_RecordReturn_Handler}
var LogisticsService_ProcessExchange_MethodDesc = grpc.MethodDesc{MethodName: "ProcessExchange", Handler: _LogisticsService_ProcessExchange_Handler}

// ---------------------------------------------------------------------------
// RegisterLogisticsServiceServerWithDepth — registers ALL methods
// ---------------------------------------------------------------------------

// RegisterLogisticsServiceServerWithDepth registers the full LogisticsService
// surface: original RPCs + procurement (BL-LOG-010..012) + Wave 5 depth
// (BL-LOG-013..029). Replaces RegisterLogisticsServiceServerWithExtensions.
func RegisterLogisticsServiceServerWithDepth(s grpc.ServiceRegistrar, srv interface {
	LogisticsServiceServer
	OnBookingPaidHandler
	ShipFulfillmentTaskHandler
	GeneratePickupQRHandler
	RedeemPickupQRHandler
	ListFulfillmentTasksHandler
	CreatePurchaseRequestHandler
	ApprovePurchaseRequestHandler
	RecordGRNWithQCHandler
	CreateKitAssemblyHandler
	LogisticsDepthHandler
}) {
	desc := grpc.ServiceDesc{
		ServiceName: "pb.LogisticsService",
		HandlerType: (*LogisticsServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			// Healthz — generated
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
			// Wave 3 RPCs
			LogisticsService_OnBookingPaid_MethodDesc,
			{MethodName: "ShipFulfillmentTask", Handler: _LogisticsService_ShipFulfillmentTask_Handler},
			{MethodName: "GeneratePickupQR", Handler: _LogisticsService_GeneratePickupQR_Handler},
			{MethodName: "RedeemPickupQR", Handler: _LogisticsService_RedeemPickupQR_Handler},
			{MethodName: "ListFulfillmentTasks", Handler: _LogisticsService_ListFulfillmentTasks_Handler},
			// BL-LOG-010..012 procurement
			{MethodName: "CreatePurchaseRequest", Handler: _LogisticsService_CreatePurchaseRequest_Handler},
			{MethodName: "ApprovePurchaseRequest", Handler: _LogisticsService_ApprovePurchaseRequest_Handler},
			{MethodName: "RecordGRNWithQC", Handler: _LogisticsService_RecordGRNWithQC_Handler},
			{MethodName: "CreateKitAssembly", Handler: _LogisticsService_CreateKitAssembly_Handler},
			// Wave 5 depth (BL-LOG-013..029)
			LogisticsService_ListPurchaseRequests_MethodDesc,
			LogisticsService_GetBudgetSyncStatus_MethodDesc,
			LogisticsService_GetTieredApprovals_MethodDesc,
			LogisticsService_DecideTieredApproval_MethodDesc,
			LogisticsService_AutoSelectVendor_MethodDesc,
			LogisticsService_RecordPartialGRN_MethodDesc,
			LogisticsService_ReverseGRN_MethodDesc,
			LogisticsService_GenerateBarcode_MethodDesc,
			LogisticsService_PrintSKULabel_MethodDesc,
			LogisticsService_CreateWarehouse_MethodDesc,
			LogisticsService_TransferStock_MethodDesc,
			LogisticsService_GetStockAlerts_MethodDesc,
			LogisticsService_SetReorderLevel_MethodDesc,
			LogisticsService_StartStocktake_MethodDesc,
			LogisticsService_RecordStocktakeCount_MethodDesc,
			LogisticsService_FinalizeStocktake_MethodDesc,
			LogisticsService_SyncFulfillmentSizes_MethodDesc,
			LogisticsService_RecordCourierTracking_MethodDesc,
			LogisticsService_RecordReturn_MethodDesc,
			LogisticsService_ProcessExchange_MethodDesc,
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "logistics.proto",
	}
	s.RegisterService(&desc, srv)
}
