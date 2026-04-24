// finance_depth_grpc_ext.go — hand-written gRPC interface extension for Wave 4
// Finance depth RPCs (BL-FIN-020..041).
//
// Follows the same pattern as correction_grpc_ext.go and pl_grpc_ext.go.
// Run `make genpb` after updating finance.proto to replace with generated code.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const (
	FinanceService_ScheduleBilling_FullMethodName               = "/pb.FinanceService/ScheduleBilling"
	FinanceService_RecordBankTransaction_FullMethodName         = "/pb.FinanceService/RecordBankTransaction"
	FinanceService_GetBankReconciliation_FullMethodName         = "/pb.FinanceService/GetBankReconciliation"
	FinanceService_GetARSubledger_FullMethodName                = "/pb.FinanceService/GetARSubledger"
	FinanceService_IssueDigitalReceipt_FullMethodName           = "/pb.FinanceService/IssueDigitalReceipt"
	FinanceService_GetDigitalReceipt_FullMethodName             = "/pb.FinanceService/GetDigitalReceipt"
	FinanceService_RecordManualPayment_FullMethodName           = "/pb.FinanceService/RecordManualPayment"
	FinanceService_CreateVendor_FullMethodName                  = "/pb.FinanceService/CreateVendor"
	FinanceService_UpdateVendor_FullMethodName                  = "/pb.FinanceService/UpdateVendor"
	FinanceService_ListVendors_FullMethodName                   = "/pb.FinanceService/ListVendors"
	FinanceService_DeleteVendor_FullMethodName                  = "/pb.FinanceService/DeleteVendor"
	FinanceService_GetAPSubledger_FullMethodName                = "/pb.FinanceService/GetAPSubledger"
	FinanceService_ListPendingAuthorizations_FullMethodName     = "/pb.FinanceService/ListPendingAuthorizations"
	FinanceService_DecidePaymentAuthorization_FullMethodName    = "/pb.FinanceService/DecidePaymentAuthorization"
	FinanceService_RecordPettyCash_FullMethodName               = "/pb.FinanceService/RecordPettyCash"
	FinanceService_ClosePettyCashPeriod_FullMethodName          = "/pb.FinanceService/ClosePettyCashPeriod"
	FinanceService_GetProjectCosts_FullMethodName               = "/pb.FinanceService/GetProjectCosts"
	FinanceService_GetDeparturePL_FullMethodName                = "/pb.FinanceService/GetDeparturePL"
	FinanceService_GetBudgetVsActual_FullMethodName             = "/pb.FinanceService/GetBudgetVsActual"
	FinanceService_TriggerAutoJournal_FullMethodName            = "/pb.FinanceService/TriggerAutoJournal"
	FinanceService_GetRevenueRecognitionPolicy_FullMethodName   = "/pb.FinanceService/GetRevenueRecognitionPolicy"
	FinanceService_SetRevenueRecognitionPolicy_FullMethodName   = "/pb.FinanceService/SetRevenueRecognitionPolicy"
	FinanceService_SetExchangeRate_FullMethodName               = "/pb.FinanceService/SetExchangeRate"
	FinanceService_GetExchangeRate_FullMethodName               = "/pb.FinanceService/GetExchangeRate"
	FinanceService_CreateFixedAsset_FullMethodName              = "/pb.FinanceService/CreateFixedAsset"
	FinanceService_ListFixedAssets_FullMethodName               = "/pb.FinanceService/ListFixedAssets"
	FinanceService_RunDepreciation_FullMethodName               = "/pb.FinanceService/RunDepreciation"
	FinanceService_CalculateTax_FullMethodName                  = "/pb.FinanceService/CalculateTax"
	FinanceService_GetTaxReport_FullMethodName                  = "/pb.FinanceService/GetTaxReport"
	FinanceService_CreateCommissionPayout_FullMethodName        = "/pb.FinanceService/CreateCommissionPayout"
	FinanceService_DecideCommissionPayout_FullMethodName        = "/pb.FinanceService/DecideCommissionPayout"
	FinanceService_GetRealtimeFinancialSummary_FullMethodName   = "/pb.FinanceService/GetRealtimeFinancialSummary"
	FinanceService_GetCashFlowDashboard_FullMethodName          = "/pb.FinanceService/GetCashFlowDashboard"
	FinanceService_GetAgingAlerts_FullMethodName                = "/pb.FinanceService/GetAgingAlerts"
	FinanceService_SearchFinanceAuditLog_FullMethodName         = "/pb.FinanceService/SearchFinanceAuditLog"
)

// ---------------------------------------------------------------------------
// Server-side handler interface
// ---------------------------------------------------------------------------

// FinanceWave4Handler is the server-side interface for Wave 4 finance depth RPCs.
type FinanceWave4Handler interface {
	ScheduleBilling(context.Context, *ScheduleBillingRequest) (*ScheduleBillingResponse, error)
	RecordBankTransaction(context.Context, *RecordBankTransactionRequest) (*RecordBankTransactionResponse, error)
	GetBankReconciliation(context.Context, *GetBankReconciliationRequest) (*GetBankReconciliationResponse, error)
	GetARSubledger(context.Context, *GetARSubledgerRequest) (*GetARSubledgerResponse, error)
	IssueDigitalReceipt(context.Context, *IssueDigitalReceiptRequest) (*IssueDigitalReceiptResponse, error)
	GetDigitalReceipt(context.Context, *GetDigitalReceiptRequest) (*GetDigitalReceiptResponse, error)
	RecordManualPayment(context.Context, *RecordManualPaymentRequest) (*RecordManualPaymentResponse, error)
	CreateVendor(context.Context, *CreateVendorRequest) (*CreateVendorResponse, error)
	UpdateVendor(context.Context, *UpdateVendorRequest) (*UpdateVendorResponse, error)
	ListVendors(context.Context, *ListVendorsRequest) (*ListVendorsResponse, error)
	DeleteVendor(context.Context, *DeleteVendorRequest) (*DeleteVendorResponse, error)
	GetAPSubledger(context.Context, *GetAPSubledgerRequest) (*GetAPSubledgerResponse, error)
	ListPendingAuthorizations(context.Context, *ListPendingAuthorizationsRequest) (*ListPendingAuthorizationsResponse, error)
	DecidePaymentAuthorization(context.Context, *DecidePaymentAuthorizationRequest) (*DecidePaymentAuthorizationResponse, error)
	RecordPettyCash(context.Context, *RecordPettyCashRequest) (*RecordPettyCashResponse, error)
	ClosePettyCashPeriod(context.Context, *ClosePettyCashPeriodRequest) (*ClosePettyCashPeriodResponse, error)
	GetProjectCosts(context.Context, *GetProjectCostsRequest) (*GetProjectCostsResponse, error)
	GetDeparturePL(context.Context, *GetDeparturePLRequest) (*GetDeparturePLResponse, error)
	GetBudgetVsActual(context.Context, *GetBudgetVsActualRequest) (*GetBudgetVsActualResponse, error)
	TriggerAutoJournal(context.Context, *TriggerAutoJournalRequest) (*TriggerAutoJournalResponse, error)
	GetRevenueRecognitionPolicy(context.Context, *GetRevenueRecognitionPolicyRequest) (*GetRevenueRecognitionPolicyResponse, error)
	SetRevenueRecognitionPolicy(context.Context, *SetRevenueRecognitionPolicyRequest) (*SetRevenueRecognitionPolicyResponse, error)
	SetExchangeRate(context.Context, *SetExchangeRateRequest) (*SetExchangeRateResponse, error)
	GetExchangeRate(context.Context, *GetExchangeRateRequest) (*GetExchangeRateResponse, error)
	CreateFixedAsset(context.Context, *CreateFixedAssetRequest) (*CreateFixedAssetResponse, error)
	ListFixedAssets(context.Context, *ListFixedAssetsRequest) (*ListFixedAssetsResponse, error)
	RunDepreciation(context.Context, *RunDepreciationRequest) (*RunDepreciationResponse, error)
	CalculateTax(context.Context, *CalculateTaxRequest) (*CalculateTaxResponse, error)
	GetTaxReport(context.Context, *GetTaxReportRequest) (*GetTaxReportResponse, error)
	CreateCommissionPayout(context.Context, *CreateCommissionPayoutRequest) (*CreateCommissionPayoutResponse, error)
	DecideCommissionPayout(context.Context, *DecideCommissionPayoutRequest) (*DecideCommissionPayoutResponse, error)
	GetRealtimeFinancialSummary(context.Context, *GetRealtimeFinancialSummaryRequest) (*GetRealtimeFinancialSummaryResponse, error)
	GetCashFlowDashboard(context.Context, *GetCashFlowDashboardRequest) (*GetCashFlowDashboardResponse, error)
	GetAgingAlerts(context.Context, *GetAgingAlertsRequest) (*GetAgingAlertsResponse, error)
	SearchFinanceAuditLog(context.Context, *SearchFinanceAuditLogRequest) (*SearchFinanceAuditLogResponse, error)
}

// UnimplementedFinanceWave4Handler provides safe defaults.
type UnimplementedFinanceWave4Handler struct{}

func (UnimplementedFinanceWave4Handler) ScheduleBilling(context.Context, *ScheduleBillingRequest) (*ScheduleBillingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ScheduleBilling not implemented")
}
func (UnimplementedFinanceWave4Handler) RecordBankTransaction(context.Context, *RecordBankTransactionRequest) (*RecordBankTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordBankTransaction not implemented")
}
func (UnimplementedFinanceWave4Handler) GetBankReconciliation(context.Context, *GetBankReconciliationRequest) (*GetBankReconciliationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBankReconciliation not implemented")
}
func (UnimplementedFinanceWave4Handler) GetARSubledger(context.Context, *GetARSubledgerRequest) (*GetARSubledgerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetARSubledger not implemented")
}
func (UnimplementedFinanceWave4Handler) IssueDigitalReceipt(context.Context, *IssueDigitalReceiptRequest) (*IssueDigitalReceiptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IssueDigitalReceipt not implemented")
}
func (UnimplementedFinanceWave4Handler) GetDigitalReceipt(context.Context, *GetDigitalReceiptRequest) (*GetDigitalReceiptResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDigitalReceipt not implemented")
}
func (UnimplementedFinanceWave4Handler) RecordManualPayment(context.Context, *RecordManualPaymentRequest) (*RecordManualPaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordManualPayment not implemented")
}
func (UnimplementedFinanceWave4Handler) CreateVendor(context.Context, *CreateVendorRequest) (*CreateVendorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateVendor not implemented")
}
func (UnimplementedFinanceWave4Handler) UpdateVendor(context.Context, *UpdateVendorRequest) (*UpdateVendorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVendor not implemented")
}
func (UnimplementedFinanceWave4Handler) ListVendors(context.Context, *ListVendorsRequest) (*ListVendorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVendors not implemented")
}
func (UnimplementedFinanceWave4Handler) DeleteVendor(context.Context, *DeleteVendorRequest) (*DeleteVendorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteVendor not implemented")
}
func (UnimplementedFinanceWave4Handler) GetAPSubledger(context.Context, *GetAPSubledgerRequest) (*GetAPSubledgerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAPSubledger not implemented")
}
func (UnimplementedFinanceWave4Handler) ListPendingAuthorizations(context.Context, *ListPendingAuthorizationsRequest) (*ListPendingAuthorizationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPendingAuthorizations not implemented")
}
func (UnimplementedFinanceWave4Handler) DecidePaymentAuthorization(context.Context, *DecidePaymentAuthorizationRequest) (*DecidePaymentAuthorizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DecidePaymentAuthorization not implemented")
}
func (UnimplementedFinanceWave4Handler) RecordPettyCash(context.Context, *RecordPettyCashRequest) (*RecordPettyCashResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordPettyCash not implemented")
}
func (UnimplementedFinanceWave4Handler) ClosePettyCashPeriod(context.Context, *ClosePettyCashPeriodRequest) (*ClosePettyCashPeriodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClosePettyCashPeriod not implemented")
}
func (UnimplementedFinanceWave4Handler) GetProjectCosts(context.Context, *GetProjectCostsRequest) (*GetProjectCostsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProjectCosts not implemented")
}
func (UnimplementedFinanceWave4Handler) GetDeparturePL(context.Context, *GetDeparturePLRequest) (*GetDeparturePLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDeparturePL not implemented")
}
func (UnimplementedFinanceWave4Handler) GetBudgetVsActual(context.Context, *GetBudgetVsActualRequest) (*GetBudgetVsActualResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBudgetVsActual not implemented")
}
func (UnimplementedFinanceWave4Handler) TriggerAutoJournal(context.Context, *TriggerAutoJournalRequest) (*TriggerAutoJournalResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerAutoJournal not implemented")
}
func (UnimplementedFinanceWave4Handler) GetRevenueRecognitionPolicy(context.Context, *GetRevenueRecognitionPolicyRequest) (*GetRevenueRecognitionPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRevenueRecognitionPolicy not implemented")
}
func (UnimplementedFinanceWave4Handler) SetRevenueRecognitionPolicy(context.Context, *SetRevenueRecognitionPolicyRequest) (*SetRevenueRecognitionPolicyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetRevenueRecognitionPolicy not implemented")
}
func (UnimplementedFinanceWave4Handler) SetExchangeRate(context.Context, *SetExchangeRateRequest) (*SetExchangeRateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetExchangeRate not implemented")
}
func (UnimplementedFinanceWave4Handler) GetExchangeRate(context.Context, *GetExchangeRateRequest) (*GetExchangeRateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetExchangeRate not implemented")
}
func (UnimplementedFinanceWave4Handler) CreateFixedAsset(context.Context, *CreateFixedAssetRequest) (*CreateFixedAssetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFixedAsset not implemented")
}
func (UnimplementedFinanceWave4Handler) ListFixedAssets(context.Context, *ListFixedAssetsRequest) (*ListFixedAssetsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListFixedAssets not implemented")
}
func (UnimplementedFinanceWave4Handler) RunDepreciation(context.Context, *RunDepreciationRequest) (*RunDepreciationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunDepreciation not implemented")
}
func (UnimplementedFinanceWave4Handler) CalculateTax(context.Context, *CalculateTaxRequest) (*CalculateTaxResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CalculateTax not implemented")
}
func (UnimplementedFinanceWave4Handler) GetTaxReport(context.Context, *GetTaxReportRequest) (*GetTaxReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTaxReport not implemented")
}
func (UnimplementedFinanceWave4Handler) CreateCommissionPayout(context.Context, *CreateCommissionPayoutRequest) (*CreateCommissionPayoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCommissionPayout not implemented")
}
func (UnimplementedFinanceWave4Handler) DecideCommissionPayout(context.Context, *DecideCommissionPayoutRequest) (*DecideCommissionPayoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DecideCommissionPayout not implemented")
}
func (UnimplementedFinanceWave4Handler) GetRealtimeFinancialSummary(context.Context, *GetRealtimeFinancialSummaryRequest) (*GetRealtimeFinancialSummaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRealtimeFinancialSummary not implemented")
}
func (UnimplementedFinanceWave4Handler) GetCashFlowDashboard(context.Context, *GetCashFlowDashboardRequest) (*GetCashFlowDashboardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCashFlowDashboard not implemented")
}
func (UnimplementedFinanceWave4Handler) GetAgingAlerts(context.Context, *GetAgingAlertsRequest) (*GetAgingAlertsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAgingAlerts not implemented")
}
func (UnimplementedFinanceWave4Handler) SearchFinanceAuditLog(context.Context, *SearchFinanceAuditLogRequest) (*SearchFinanceAuditLogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchFinanceAuditLog not implemented")
}

// ---------------------------------------------------------------------------
// gRPC handler functions
// ---------------------------------------------------------------------------

func _FinanceService_ScheduleBilling_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScheduleBillingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).ScheduleBilling(ctx, req.(*ScheduleBillingRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_ScheduleBilling_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_RecordBankTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordBankTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).RecordBankTransaction(ctx, req.(*RecordBankTransactionRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_RecordBankTransaction_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetBankReconciliation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBankReconciliationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).GetBankReconciliation(ctx, req.(*GetBankReconciliationRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetBankReconciliation_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetARSubledger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetARSubledgerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).GetARSubledger(ctx, req.(*GetARSubledgerRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetARSubledger_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_IssueDigitalReceipt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IssueDigitalReceiptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).IssueDigitalReceipt(ctx, req.(*IssueDigitalReceiptRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_IssueDigitalReceipt_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetDigitalReceipt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDigitalReceiptRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).GetDigitalReceipt(ctx, req.(*GetDigitalReceiptRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetDigitalReceipt_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_RecordManualPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordManualPaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).RecordManualPayment(ctx, req.(*RecordManualPaymentRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_RecordManualPayment_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_CreateVendor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateVendorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).CreateVendor(ctx, req.(*CreateVendorRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_CreateVendor_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_UpdateVendor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateVendorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).UpdateVendor(ctx, req.(*UpdateVendorRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_UpdateVendor_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_ListVendors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListVendorsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).ListVendors(ctx, req.(*ListVendorsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_ListVendors_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_DeleteVendor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteVendorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).DeleteVendor(ctx, req.(*DeleteVendorRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_DeleteVendor_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetAPSubledger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAPSubledgerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).GetAPSubledger(ctx, req.(*GetAPSubledgerRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetAPSubledger_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_ListPendingAuthorizations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPendingAuthorizationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).ListPendingAuthorizations(ctx, req.(*ListPendingAuthorizationsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_ListPendingAuthorizations_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_DecidePaymentAuthorization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DecidePaymentAuthorizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).DecidePaymentAuthorization(ctx, req.(*DecidePaymentAuthorizationRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_DecidePaymentAuthorization_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_RecordPettyCash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordPettyCashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).RecordPettyCash(ctx, req.(*RecordPettyCashRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_RecordPettyCash_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_ClosePettyCashPeriod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClosePettyCashPeriodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).ClosePettyCashPeriod(ctx, req.(*ClosePettyCashPeriodRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_ClosePettyCashPeriod_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetProjectCosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProjectCostsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).GetProjectCosts(ctx, req.(*GetProjectCostsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetProjectCosts_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetDeparturePL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDeparturePLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).GetDeparturePL(ctx, req.(*GetDeparturePLRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetDeparturePL_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetBudgetVsActual_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBudgetVsActualRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).GetBudgetVsActual(ctx, req.(*GetBudgetVsActualRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetBudgetVsActual_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_TriggerAutoJournal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TriggerAutoJournalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).TriggerAutoJournal(ctx, req.(*TriggerAutoJournalRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_TriggerAutoJournal_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetRevenueRecognitionPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRevenueRecognitionPolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).GetRevenueRecognitionPolicy(ctx, req.(*GetRevenueRecognitionPolicyRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetRevenueRecognitionPolicy_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_SetRevenueRecognitionPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRevenueRecognitionPolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).SetRevenueRecognitionPolicy(ctx, req.(*SetRevenueRecognitionPolicyRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_SetRevenueRecognitionPolicy_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_SetExchangeRate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetExchangeRateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).SetExchangeRate(ctx, req.(*SetExchangeRateRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_SetExchangeRate_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetExchangeRate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetExchangeRateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).GetExchangeRate(ctx, req.(*GetExchangeRateRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetExchangeRate_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_CreateFixedAsset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFixedAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).CreateFixedAsset(ctx, req.(*CreateFixedAssetRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_CreateFixedAsset_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_ListFixedAssets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListFixedAssetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).ListFixedAssets(ctx, req.(*ListFixedAssetsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_ListFixedAssets_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_RunDepreciation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunDepreciationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).RunDepreciation(ctx, req.(*RunDepreciationRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_RunDepreciation_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_CalculateTax_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CalculateTaxRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).CalculateTax(ctx, req.(*CalculateTaxRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_CalculateTax_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetTaxReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTaxReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).GetTaxReport(ctx, req.(*GetTaxReportRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetTaxReport_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_CreateCommissionPayout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCommissionPayoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).CreateCommissionPayout(ctx, req.(*CreateCommissionPayoutRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_CreateCommissionPayout_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_DecideCommissionPayout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DecideCommissionPayoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).DecideCommissionPayout(ctx, req.(*DecideCommissionPayoutRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_DecideCommissionPayout_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetRealtimeFinancialSummary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRealtimeFinancialSummaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).GetRealtimeFinancialSummary(ctx, req.(*GetRealtimeFinancialSummaryRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetRealtimeFinancialSummary_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetCashFlowDashboard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCashFlowDashboardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).GetCashFlowDashboard(ctx, req.(*GetCashFlowDashboardRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetCashFlowDashboard_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_GetAgingAlerts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAgingAlertsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).GetAgingAlerts(ctx, req.(*GetAgingAlertsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_GetAgingAlerts_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

func _FinanceService_SearchFinanceAuditLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchFinanceAuditLogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FinanceWave4Handler).SearchFinanceAuditLog(ctx, req.(*SearchFinanceAuditLogRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: FinanceService_SearchFinanceAuditLog_FullMethodName}
	return interceptor(ctx, in, info, handler)
}

// ---------------------------------------------------------------------------
// gRPC method descriptors
// ---------------------------------------------------------------------------

var FinanceService_ScheduleBilling_MethodDesc = grpc.MethodDesc{
	MethodName: "ScheduleBilling",
	Handler:    _FinanceService_ScheduleBilling_Handler,
}
var FinanceService_RecordBankTransaction_MethodDesc = grpc.MethodDesc{
	MethodName: "RecordBankTransaction",
	Handler:    _FinanceService_RecordBankTransaction_Handler,
}
var FinanceService_GetBankReconciliation_MethodDesc = grpc.MethodDesc{
	MethodName: "GetBankReconciliation",
	Handler:    _FinanceService_GetBankReconciliation_Handler,
}
var FinanceService_GetARSubledger_MethodDesc = grpc.MethodDesc{
	MethodName: "GetARSubledger",
	Handler:    _FinanceService_GetARSubledger_Handler,
}
var FinanceService_IssueDigitalReceipt_MethodDesc = grpc.MethodDesc{
	MethodName: "IssueDigitalReceipt",
	Handler:    _FinanceService_IssueDigitalReceipt_Handler,
}
var FinanceService_GetDigitalReceipt_MethodDesc = grpc.MethodDesc{
	MethodName: "GetDigitalReceipt",
	Handler:    _FinanceService_GetDigitalReceipt_Handler,
}
var FinanceService_RecordManualPayment_MethodDesc = grpc.MethodDesc{
	MethodName: "RecordManualPayment",
	Handler:    _FinanceService_RecordManualPayment_Handler,
}
var FinanceService_CreateVendor_MethodDesc = grpc.MethodDesc{
	MethodName: "CreateVendor",
	Handler:    _FinanceService_CreateVendor_Handler,
}
var FinanceService_UpdateVendor_MethodDesc = grpc.MethodDesc{
	MethodName: "UpdateVendor",
	Handler:    _FinanceService_UpdateVendor_Handler,
}
var FinanceService_ListVendors_MethodDesc = grpc.MethodDesc{
	MethodName: "ListVendors",
	Handler:    _FinanceService_ListVendors_Handler,
}
var FinanceService_DeleteVendor_MethodDesc = grpc.MethodDesc{
	MethodName: "DeleteVendor",
	Handler:    _FinanceService_DeleteVendor_Handler,
}
var FinanceService_GetAPSubledger_MethodDesc = grpc.MethodDesc{
	MethodName: "GetAPSubledger",
	Handler:    _FinanceService_GetAPSubledger_Handler,
}
var FinanceService_ListPendingAuthorizations_MethodDesc = grpc.MethodDesc{
	MethodName: "ListPendingAuthorizations",
	Handler:    _FinanceService_ListPendingAuthorizations_Handler,
}
var FinanceService_DecidePaymentAuthorization_MethodDesc = grpc.MethodDesc{
	MethodName: "DecidePaymentAuthorization",
	Handler:    _FinanceService_DecidePaymentAuthorization_Handler,
}
var FinanceService_RecordPettyCash_MethodDesc = grpc.MethodDesc{
	MethodName: "RecordPettyCash",
	Handler:    _FinanceService_RecordPettyCash_Handler,
}
var FinanceService_ClosePettyCashPeriod_MethodDesc = grpc.MethodDesc{
	MethodName: "ClosePettyCashPeriod",
	Handler:    _FinanceService_ClosePettyCashPeriod_Handler,
}
var FinanceService_GetProjectCosts_MethodDesc = grpc.MethodDesc{
	MethodName: "GetProjectCosts",
	Handler:    _FinanceService_GetProjectCosts_Handler,
}
var FinanceService_GetDeparturePL_MethodDesc = grpc.MethodDesc{
	MethodName: "GetDeparturePL",
	Handler:    _FinanceService_GetDeparturePL_Handler,
}
var FinanceService_GetBudgetVsActual_MethodDesc = grpc.MethodDesc{
	MethodName: "GetBudgetVsActual",
	Handler:    _FinanceService_GetBudgetVsActual_Handler,
}
var FinanceService_TriggerAutoJournal_MethodDesc = grpc.MethodDesc{
	MethodName: "TriggerAutoJournal",
	Handler:    _FinanceService_TriggerAutoJournal_Handler,
}
var FinanceService_GetRevenueRecognitionPolicy_MethodDesc = grpc.MethodDesc{
	MethodName: "GetRevenueRecognitionPolicy",
	Handler:    _FinanceService_GetRevenueRecognitionPolicy_Handler,
}
var FinanceService_SetRevenueRecognitionPolicy_MethodDesc = grpc.MethodDesc{
	MethodName: "SetRevenueRecognitionPolicy",
	Handler:    _FinanceService_SetRevenueRecognitionPolicy_Handler,
}
var FinanceService_SetExchangeRate_MethodDesc = grpc.MethodDesc{
	MethodName: "SetExchangeRate",
	Handler:    _FinanceService_SetExchangeRate_Handler,
}
var FinanceService_GetExchangeRate_MethodDesc = grpc.MethodDesc{
	MethodName: "GetExchangeRate",
	Handler:    _FinanceService_GetExchangeRate_Handler,
}
var FinanceService_CreateFixedAsset_MethodDesc = grpc.MethodDesc{
	MethodName: "CreateFixedAsset",
	Handler:    _FinanceService_CreateFixedAsset_Handler,
}
var FinanceService_ListFixedAssets_MethodDesc = grpc.MethodDesc{
	MethodName: "ListFixedAssets",
	Handler:    _FinanceService_ListFixedAssets_Handler,
}
var FinanceService_RunDepreciation_MethodDesc = grpc.MethodDesc{
	MethodName: "RunDepreciation",
	Handler:    _FinanceService_RunDepreciation_Handler,
}
var FinanceService_CalculateTax_MethodDesc = grpc.MethodDesc{
	MethodName: "CalculateTax",
	Handler:    _FinanceService_CalculateTax_Handler,
}
var FinanceService_GetTaxReport_MethodDesc = grpc.MethodDesc{
	MethodName: "GetTaxReport",
	Handler:    _FinanceService_GetTaxReport_Handler,
}
var FinanceService_CreateCommissionPayout_MethodDesc = grpc.MethodDesc{
	MethodName: "CreateCommissionPayout",
	Handler:    _FinanceService_CreateCommissionPayout_Handler,
}
var FinanceService_DecideCommissionPayout_MethodDesc = grpc.MethodDesc{
	MethodName: "DecideCommissionPayout",
	Handler:    _FinanceService_DecideCommissionPayout_Handler,
}
var FinanceService_GetRealtimeFinancialSummary_MethodDesc = grpc.MethodDesc{
	MethodName: "GetRealtimeFinancialSummary",
	Handler:    _FinanceService_GetRealtimeFinancialSummary_Handler,
}
var FinanceService_GetCashFlowDashboard_MethodDesc = grpc.MethodDesc{
	MethodName: "GetCashFlowDashboard",
	Handler:    _FinanceService_GetCashFlowDashboard_Handler,
}
var FinanceService_GetAgingAlerts_MethodDesc = grpc.MethodDesc{
	MethodName: "GetAgingAlerts",
	Handler:    _FinanceService_GetAgingAlerts_Handler,
}
var FinanceService_SearchFinanceAuditLog_MethodDesc = grpc.MethodDesc{
	MethodName: "SearchFinanceAuditLog",
	Handler:    _FinanceService_SearchFinanceAuditLog_Handler,
}

// ---------------------------------------------------------------------------
// RegisterFinanceServiceServerWithWave4 — the canonical registration function
// for Wave 4. Supersedes RegisterFinanceServiceServerWithCorrections.
// ---------------------------------------------------------------------------

// RegisterFinanceServiceServerWithWave4 registers the complete FinanceService
// including all Wave 4 depth RPCs (BL-FIN-020..041).
func RegisterFinanceServiceServerWithWave4(s grpc.ServiceRegistrar, srv interface {
	FinanceServiceServer
	OnPaymentReceivedHandler
	FinanceReportsHandler
	FinanceDepthHandler
	OnGRNReceivedHandler
	CorrectionHandler
	FinanceWave4Handler
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
			// Wave 4 methods
			FinanceService_ScheduleBilling_MethodDesc,
			FinanceService_RecordBankTransaction_MethodDesc,
			FinanceService_GetBankReconciliation_MethodDesc,
			FinanceService_GetARSubledger_MethodDesc,
			FinanceService_IssueDigitalReceipt_MethodDesc,
			FinanceService_GetDigitalReceipt_MethodDesc,
			FinanceService_RecordManualPayment_MethodDesc,
			FinanceService_CreateVendor_MethodDesc,
			FinanceService_UpdateVendor_MethodDesc,
			FinanceService_ListVendors_MethodDesc,
			FinanceService_DeleteVendor_MethodDesc,
			FinanceService_GetAPSubledger_MethodDesc,
			FinanceService_ListPendingAuthorizations_MethodDesc,
			FinanceService_DecidePaymentAuthorization_MethodDesc,
			FinanceService_RecordPettyCash_MethodDesc,
			FinanceService_ClosePettyCashPeriod_MethodDesc,
			FinanceService_GetProjectCosts_MethodDesc,
			FinanceService_GetDeparturePL_MethodDesc,
			FinanceService_GetBudgetVsActual_MethodDesc,
			FinanceService_TriggerAutoJournal_MethodDesc,
			FinanceService_GetRevenueRecognitionPolicy_MethodDesc,
			FinanceService_SetRevenueRecognitionPolicy_MethodDesc,
			FinanceService_SetExchangeRate_MethodDesc,
			FinanceService_GetExchangeRate_MethodDesc,
			FinanceService_CreateFixedAsset_MethodDesc,
			FinanceService_ListFixedAssets_MethodDesc,
			FinanceService_RunDepreciation_MethodDesc,
			FinanceService_CalculateTax_MethodDesc,
			FinanceService_GetTaxReport_MethodDesc,
			FinanceService_CreateCommissionPayout_MethodDesc,
			FinanceService_DecideCommissionPayout_MethodDesc,
			FinanceService_GetRealtimeFinancialSummary_MethodDesc,
			FinanceService_GetCashFlowDashboard_MethodDesc,
			FinanceService_GetAgingAlerts_MethodDesc,
			FinanceService_SearchFinanceAuditLog_MethodDesc,
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "finance.proto",
	}
	s.RegisterService(&desc, srv)
}
