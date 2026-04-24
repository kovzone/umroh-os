// ops_depth_grpc_ext.go — hand-written gRPC service extension for ops-svc Wave 5 depth RPCs.
// BL-OPS-021..042: 22 new RPCs extending the OpsService.

package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Full method name constants for ops depth RPCs.
const (
	OpsService_StoreCollectiveDocument_FullMethodName      = "/pb.OpsService/StoreCollectiveDocument"
	OpsService_GetCollectiveDocuments_FullMethodName       = "/pb.OpsService/GetCollectiveDocuments"
	OpsService_SetDocumentACL_FullMethodName               = "/pb.OpsService/SetDocumentACL"
	OpsService_ExtractPassportOCR_FullMethodName           = "/pb.OpsService/ExtractPassportOCR"
	OpsService_SetMahramRelation_FullMethodName            = "/pb.OpsService/SetMahramRelation"
	OpsService_GetMahramRelations_FullMethodName           = "/pb.OpsService/GetMahramRelations"
	OpsService_GetDocumentProgress_FullMethodName          = "/pb.OpsService/GetDocumentProgress"
	OpsService_GetExpiryAlerts_FullMethodName              = "/pb.OpsService/GetExpiryAlerts"
	OpsService_GenerateOfficialLetter_FullMethodName       = "/pb.OpsService/GenerateOfficialLetter"
	OpsService_GenerateImmigrationManifest_FullMethodName  = "/pb.OpsService/GenerateImmigrationManifest"
	OpsService_RunSmartRooming_FullMethodName              = "/pb.OpsService/RunSmartRooming"
	OpsService_AssignTransport_FullMethodName              = "/pb.OpsService/AssignTransport"
	OpsService_GetTransportAssignments_FullMethodName      = "/pb.OpsService/GetTransportAssignments"
	OpsService_PublishManifestDelta_FullMethodName         = "/pb.OpsService/PublishManifestDelta"
	OpsService_AssignStaff_FullMethodName                  = "/pb.OpsService/AssignStaff"
	OpsService_RecordPassportHandover_FullMethodName       = "/pb.OpsService/RecordPassportHandover"
	OpsService_GetPassportLog_FullMethodName               = "/pb.OpsService/GetPassportLog"
	OpsService_GetVisaProgress_FullMethodName              = "/pb.OpsService/GetVisaProgress"
	OpsService_StoreEVisa_FullMethodName                   = "/pb.OpsService/StoreEVisa"
	OpsService_GetEVisa_FullMethodName                     = "/pb.OpsService/GetEVisa"
	OpsService_TriggerExternalProvider_FullMethodName      = "/pb.OpsService/TriggerExternalProvider"
	OpsService_CreateRefund_FullMethodName                 = "/pb.OpsService/CreateRefund"
	OpsService_ApproveRefund_FullMethodName                = "/pb.OpsService/ApproveRefund"
	OpsService_RecordPenalty_FullMethodName                = "/pb.OpsService/RecordPenalty"
	OpsService_RecordLuggageScan_FullMethodName            = "/pb.OpsService/RecordLuggageScan"
	OpsService_GetLuggageCount_FullMethodName              = "/pb.OpsService/GetLuggageCount"
	OpsService_BroadcastSchedule_FullMethodName            = "/pb.OpsService/BroadcastSchedule"
	OpsService_IssueDigitalTasreh_FullMethodName           = "/pb.OpsService/IssueDigitalTasreh"
	OpsService_RecordRaudhahEntry_FullMethodName           = "/pb.OpsService/RecordRaudhahEntry"
	OpsService_RegisterAudioDevice_FullMethodName          = "/pb.OpsService/RegisterAudioDevice"
	OpsService_UpdateAudioDeviceStatus_FullMethodName      = "/pb.OpsService/UpdateAudioDeviceStatus"
	OpsService_RecordZamzamDistribution_FullMethodName     = "/pb.OpsService/RecordZamzamDistribution"
	OpsService_RecordRoomCheckIn_FullMethodName            = "/pb.OpsService/RecordRoomCheckIn"
)

// OpsDepthHandler is the server-side interface for all Wave 5 ops depth RPCs.
type OpsDepthHandler interface {
	StoreCollectiveDocument(context.Context, *StoreCollectiveDocumentRequest) (*StoreCollectiveDocumentResponse, error)
	GetCollectiveDocuments(context.Context, *GetCollectiveDocumentsRequest) (*GetCollectiveDocumentsResponse, error)
	SetDocumentACL(context.Context, *SetDocumentACLRequest) (*SetDocumentACLResponse, error)
	ExtractPassportOCR(context.Context, *ExtractPassportOCRRequest) (*ExtractPassportOCRResponse, error)
	SetMahramRelation(context.Context, *SetMahramRelationRequest) (*SetMahramRelationResponse, error)
	GetMahramRelations(context.Context, *GetMahramRelationsRequest) (*GetMahramRelationsResponse, error)
	GetDocumentProgress(context.Context, *GetDocumentProgressRequest) (*GetDocumentProgressResponse, error)
	GetExpiryAlerts(context.Context, *GetExpiryAlertsRequest) (*GetExpiryAlertsResponse, error)
	GenerateOfficialLetter(context.Context, *GenerateOfficialLetterRequest) (*GenerateOfficialLetterResponse, error)
	GenerateImmigrationManifest(context.Context, *GenerateImmigrationManifestRequest) (*GenerateImmigrationManifestResponse, error)
	RunSmartRooming(context.Context, *RunSmartRoomingRequest) (*RunSmartRoomingResponse, error)
	AssignTransport(context.Context, *AssignTransportRequest) (*AssignTransportResponse, error)
	GetTransportAssignments(context.Context, *GetTransportAssignmentsRequest) (*GetTransportAssignmentsResponse, error)
	PublishManifestDelta(context.Context, *PublishManifestDeltaRequest) (*PublishManifestDeltaResponse, error)
	AssignStaff(context.Context, *AssignStaffRequest) (*AssignStaffResponse, error)
	RecordPassportHandover(context.Context, *RecordPassportHandoverRequest) (*RecordPassportHandoverResponse, error)
	GetPassportLog(context.Context, *GetPassportLogRequest) (*GetPassportLogResponse, error)
	GetVisaProgress(context.Context, *GetVisaProgressRequest) (*GetVisaProgressResponse, error)
	StoreEVisa(context.Context, *StoreEVisaRequest) (*StoreEVisaResponse, error)
	GetEVisa(context.Context, *GetEVisaRequest) (*GetEVisaResponse, error)
	TriggerExternalProvider(context.Context, *TriggerExternalProviderRequest) (*TriggerExternalProviderResponse, error)
	CreateRefund(context.Context, *CreateRefundRequest) (*CreateRefundResponse, error)
	ApproveRefund(context.Context, *ApproveRefundRequest) (*ApproveRefundResponse, error)
	RecordPenalty(context.Context, *RecordPenaltyRequest) (*RecordPenaltyResponse, error)
	RecordLuggageScan(context.Context, *RecordLuggageScanRequest) (*RecordLuggageScanResponse, error)
	GetLuggageCount(context.Context, *GetLuggageCountRequest) (*GetLuggageCountResponse, error)
	BroadcastSchedule(context.Context, *BroadcastScheduleRequest) (*BroadcastScheduleResponse, error)
	IssueDigitalTasreh(context.Context, *IssueDigitalTasrehRequest) (*IssueDigitalTasrehResponse, error)
	RecordRaudhahEntry(context.Context, *RecordRaudhahEntryRequest) (*RecordRaudhahEntryResponse, error)
	RegisterAudioDevice(context.Context, *RegisterAudioDeviceRequest) (*RegisterAudioDeviceResponse, error)
	UpdateAudioDeviceStatus(context.Context, *UpdateAudioDeviceStatusRequest) (*UpdateAudioDeviceStatusResponse, error)
	RecordZamzamDistribution(context.Context, *RecordZamzamDistributionRequest) (*RecordZamzamDistributionResponse, error)
	RecordRoomCheckIn(context.Context, *RecordRoomCheckInRequest) (*RecordRoomCheckInResponse, error)
}

// UnimplementedOpsDepthHandler provides safe defaults for OpsDepthHandler.
type UnimplementedOpsDepthHandler struct{}

func (UnimplementedOpsDepthHandler) StoreCollectiveDocument(context.Context, *StoreCollectiveDocumentRequest) (*StoreCollectiveDocumentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreCollectiveDocument not implemented")
}
func (UnimplementedOpsDepthHandler) GetCollectiveDocuments(context.Context, *GetCollectiveDocumentsRequest) (*GetCollectiveDocumentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCollectiveDocuments not implemented")
}
func (UnimplementedOpsDepthHandler) SetDocumentACL(context.Context, *SetDocumentACLRequest) (*SetDocumentACLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetDocumentACL not implemented")
}
func (UnimplementedOpsDepthHandler) ExtractPassportOCR(context.Context, *ExtractPassportOCRRequest) (*ExtractPassportOCRResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExtractPassportOCR not implemented")
}
func (UnimplementedOpsDepthHandler) SetMahramRelation(context.Context, *SetMahramRelationRequest) (*SetMahramRelationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetMahramRelation not implemented")
}
func (UnimplementedOpsDepthHandler) GetMahramRelations(context.Context, *GetMahramRelationsRequest) (*GetMahramRelationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMahramRelations not implemented")
}
func (UnimplementedOpsDepthHandler) GetDocumentProgress(context.Context, *GetDocumentProgressRequest) (*GetDocumentProgressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDocumentProgress not implemented")
}
func (UnimplementedOpsDepthHandler) GetExpiryAlerts(context.Context, *GetExpiryAlertsRequest) (*GetExpiryAlertsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetExpiryAlerts not implemented")
}
func (UnimplementedOpsDepthHandler) GenerateOfficialLetter(context.Context, *GenerateOfficialLetterRequest) (*GenerateOfficialLetterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateOfficialLetter not implemented")
}
func (UnimplementedOpsDepthHandler) GenerateImmigrationManifest(context.Context, *GenerateImmigrationManifestRequest) (*GenerateImmigrationManifestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateImmigrationManifest not implemented")
}
func (UnimplementedOpsDepthHandler) RunSmartRooming(context.Context, *RunSmartRoomingRequest) (*RunSmartRoomingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunSmartRooming not implemented")
}
func (UnimplementedOpsDepthHandler) AssignTransport(context.Context, *AssignTransportRequest) (*AssignTransportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignTransport not implemented")
}
func (UnimplementedOpsDepthHandler) GetTransportAssignments(context.Context, *GetTransportAssignmentsRequest) (*GetTransportAssignmentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransportAssignments not implemented")
}
func (UnimplementedOpsDepthHandler) PublishManifestDelta(context.Context, *PublishManifestDeltaRequest) (*PublishManifestDeltaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishManifestDelta not implemented")
}
func (UnimplementedOpsDepthHandler) AssignStaff(context.Context, *AssignStaffRequest) (*AssignStaffResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignStaff not implemented")
}
func (UnimplementedOpsDepthHandler) RecordPassportHandover(context.Context, *RecordPassportHandoverRequest) (*RecordPassportHandoverResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordPassportHandover not implemented")
}
func (UnimplementedOpsDepthHandler) GetPassportLog(context.Context, *GetPassportLogRequest) (*GetPassportLogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPassportLog not implemented")
}
func (UnimplementedOpsDepthHandler) GetVisaProgress(context.Context, *GetVisaProgressRequest) (*GetVisaProgressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVisaProgress not implemented")
}
func (UnimplementedOpsDepthHandler) StoreEVisa(context.Context, *StoreEVisaRequest) (*StoreEVisaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreEVisa not implemented")
}
func (UnimplementedOpsDepthHandler) GetEVisa(context.Context, *GetEVisaRequest) (*GetEVisaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEVisa not implemented")
}
func (UnimplementedOpsDepthHandler) TriggerExternalProvider(context.Context, *TriggerExternalProviderRequest) (*TriggerExternalProviderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerExternalProvider not implemented")
}
func (UnimplementedOpsDepthHandler) CreateRefund(context.Context, *CreateRefundRequest) (*CreateRefundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRefund not implemented")
}
func (UnimplementedOpsDepthHandler) ApproveRefund(context.Context, *ApproveRefundRequest) (*ApproveRefundResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApproveRefund not implemented")
}
func (UnimplementedOpsDepthHandler) RecordPenalty(context.Context, *RecordPenaltyRequest) (*RecordPenaltyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordPenalty not implemented")
}
func (UnimplementedOpsDepthHandler) RecordLuggageScan(context.Context, *RecordLuggageScanRequest) (*RecordLuggageScanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordLuggageScan not implemented")
}
func (UnimplementedOpsDepthHandler) GetLuggageCount(context.Context, *GetLuggageCountRequest) (*GetLuggageCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLuggageCount not implemented")
}
func (UnimplementedOpsDepthHandler) BroadcastSchedule(context.Context, *BroadcastScheduleRequest) (*BroadcastScheduleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BroadcastSchedule not implemented")
}
func (UnimplementedOpsDepthHandler) IssueDigitalTasreh(context.Context, *IssueDigitalTasrehRequest) (*IssueDigitalTasrehResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IssueDigitalTasreh not implemented")
}
func (UnimplementedOpsDepthHandler) RecordRaudhahEntry(context.Context, *RecordRaudhahEntryRequest) (*RecordRaudhahEntryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordRaudhahEntry not implemented")
}
func (UnimplementedOpsDepthHandler) RegisterAudioDevice(context.Context, *RegisterAudioDeviceRequest) (*RegisterAudioDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterAudioDevice not implemented")
}
func (UnimplementedOpsDepthHandler) UpdateAudioDeviceStatus(context.Context, *UpdateAudioDeviceStatusRequest) (*UpdateAudioDeviceStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAudioDeviceStatus not implemented")
}
func (UnimplementedOpsDepthHandler) RecordZamzamDistribution(context.Context, *RecordZamzamDistributionRequest) (*RecordZamzamDistributionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordZamzamDistribution not implemented")
}
func (UnimplementedOpsDepthHandler) RecordRoomCheckIn(context.Context, *RecordRoomCheckInRequest) (*RecordRoomCheckInResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordRoomCheckIn not implemented")
}

// ---------------------------------------------------------------------------
// Handler dispatch helpers
// ---------------------------------------------------------------------------

func _OpsService_StoreCollectiveDocument_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreCollectiveDocumentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).StoreCollectiveDocument(ctx, req.(*StoreCollectiveDocumentRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_StoreCollectiveDocument_FullMethodName}, handler)
}

func _OpsService_GetCollectiveDocuments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCollectiveDocumentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).GetCollectiveDocuments(ctx, req.(*GetCollectiveDocumentsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_GetCollectiveDocuments_FullMethodName}, handler)
}

func _OpsService_SetDocumentACL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetDocumentACLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).SetDocumentACL(ctx, req.(*SetDocumentACLRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_SetDocumentACL_FullMethodName}, handler)
}

func _OpsService_ExtractPassportOCR_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExtractPassportOCRRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).ExtractPassportOCR(ctx, req.(*ExtractPassportOCRRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_ExtractPassportOCR_FullMethodName}, handler)
}

func _OpsService_SetMahramRelation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetMahramRelationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).SetMahramRelation(ctx, req.(*SetMahramRelationRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_SetMahramRelation_FullMethodName}, handler)
}

func _OpsService_GetMahramRelations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMahramRelationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).GetMahramRelations(ctx, req.(*GetMahramRelationsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_GetMahramRelations_FullMethodName}, handler)
}

func _OpsService_GetDocumentProgress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDocumentProgressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).GetDocumentProgress(ctx, req.(*GetDocumentProgressRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_GetDocumentProgress_FullMethodName}, handler)
}

func _OpsService_GetExpiryAlerts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetExpiryAlertsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).GetExpiryAlerts(ctx, req.(*GetExpiryAlertsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_GetExpiryAlerts_FullMethodName}, handler)
}

func _OpsService_GenerateOfficialLetter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateOfficialLetterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).GenerateOfficialLetter(ctx, req.(*GenerateOfficialLetterRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_GenerateOfficialLetter_FullMethodName}, handler)
}

func _OpsService_GenerateImmigrationManifest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateImmigrationManifestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).GenerateImmigrationManifest(ctx, req.(*GenerateImmigrationManifestRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_GenerateImmigrationManifest_FullMethodName}, handler)
}

func _OpsService_RunSmartRooming_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunSmartRoomingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).RunSmartRooming(ctx, req.(*RunSmartRoomingRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_RunSmartRooming_FullMethodName}, handler)
}

func _OpsService_AssignTransport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignTransportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).AssignTransport(ctx, req.(*AssignTransportRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_AssignTransport_FullMethodName}, handler)
}

func _OpsService_GetTransportAssignments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransportAssignmentsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).GetTransportAssignments(ctx, req.(*GetTransportAssignmentsRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_GetTransportAssignments_FullMethodName}, handler)
}

func _OpsService_PublishManifestDelta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishManifestDeltaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).PublishManifestDelta(ctx, req.(*PublishManifestDeltaRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_PublishManifestDelta_FullMethodName}, handler)
}

func _OpsService_AssignStaff_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignStaffRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).AssignStaff(ctx, req.(*AssignStaffRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_AssignStaff_FullMethodName}, handler)
}

func _OpsService_RecordPassportHandover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordPassportHandoverRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).RecordPassportHandover(ctx, req.(*RecordPassportHandoverRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_RecordPassportHandover_FullMethodName}, handler)
}

func _OpsService_GetPassportLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPassportLogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).GetPassportLog(ctx, req.(*GetPassportLogRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_GetPassportLog_FullMethodName}, handler)
}

func _OpsService_GetVisaProgress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVisaProgressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).GetVisaProgress(ctx, req.(*GetVisaProgressRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_GetVisaProgress_FullMethodName}, handler)
}

func _OpsService_StoreEVisa_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreEVisaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).StoreEVisa(ctx, req.(*StoreEVisaRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_StoreEVisa_FullMethodName}, handler)
}

func _OpsService_GetEVisa_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEVisaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).GetEVisa(ctx, req.(*GetEVisaRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_GetEVisa_FullMethodName}, handler)
}

func _OpsService_TriggerExternalProvider_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TriggerExternalProviderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).TriggerExternalProvider(ctx, req.(*TriggerExternalProviderRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_TriggerExternalProvider_FullMethodName}, handler)
}

func _OpsService_CreateRefund_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRefundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).CreateRefund(ctx, req.(*CreateRefundRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_CreateRefund_FullMethodName}, handler)
}

func _OpsService_ApproveRefund_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApproveRefundRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).ApproveRefund(ctx, req.(*ApproveRefundRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_ApproveRefund_FullMethodName}, handler)
}

func _OpsService_RecordPenalty_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordPenaltyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).RecordPenalty(ctx, req.(*RecordPenaltyRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_RecordPenalty_FullMethodName}, handler)
}

func _OpsService_RecordLuggageScan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordLuggageScanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).RecordLuggageScan(ctx, req.(*RecordLuggageScanRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_RecordLuggageScan_FullMethodName}, handler)
}

func _OpsService_GetLuggageCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLuggageCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).GetLuggageCount(ctx, req.(*GetLuggageCountRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_GetLuggageCount_FullMethodName}, handler)
}

func _OpsService_BroadcastSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BroadcastScheduleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).BroadcastSchedule(ctx, req.(*BroadcastScheduleRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_BroadcastSchedule_FullMethodName}, handler)
}

func _OpsService_IssueDigitalTasreh_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IssueDigitalTasrehRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).IssueDigitalTasreh(ctx, req.(*IssueDigitalTasrehRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_IssueDigitalTasreh_FullMethodName}, handler)
}

func _OpsService_RecordRaudhahEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordRaudhahEntryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).RecordRaudhahEntry(ctx, req.(*RecordRaudhahEntryRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_RecordRaudhahEntry_FullMethodName}, handler)
}

func _OpsService_RegisterAudioDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterAudioDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).RegisterAudioDevice(ctx, req.(*RegisterAudioDeviceRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_RegisterAudioDevice_FullMethodName}, handler)
}

func _OpsService_UpdateAudioDeviceStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAudioDeviceStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).UpdateAudioDeviceStatus(ctx, req.(*UpdateAudioDeviceStatusRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_UpdateAudioDeviceStatus_FullMethodName}, handler)
}

func _OpsService_RecordZamzamDistribution_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordZamzamDistributionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).RecordZamzamDistribution(ctx, req.(*RecordZamzamDistributionRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_RecordZamzamDistribution_FullMethodName}, handler)
}

func _OpsService_RecordRoomCheckIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordRoomCheckInRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OpsDepthHandler).RecordRoomCheckIn(ctx, req.(*RecordRoomCheckInRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_RecordRoomCheckIn_FullMethodName}, handler)
}

// ---------------------------------------------------------------------------
// ScanBoardingHandler interface (Phase 6 scan/boarding RPCs)
// ---------------------------------------------------------------------------

// ScanBoardingHandler is the server-side interface for Phase 6 scan + boarding RPCs.
type ScanBoardingHandler interface {
	RecordScan(context.Context, *RecordScanRequest) (*RecordScanResponse, error)
	RecordBusBoarding(context.Context, *RecordBusBoardingRequest) (*RecordBusBoardingResponse, error)
	GetBoardingRoster(context.Context, *GetBoardingRosterRequest) (*GetBoardingRosterResponse, error)
}

func _OpsService_RecordScan_DepthHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordScanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScanBoardingHandler).RecordScan(ctx, req.(*RecordScanRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_RecordScan_FullMethodName}, handler)
}

func _OpsService_RecordBusBoarding_DepthHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordBusBoardingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScanBoardingHandler).RecordBusBoarding(ctx, req.(*RecordBusBoardingRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_RecordBusBoarding_FullMethodName}, handler)
}

func _OpsService_GetBoardingRoster_DepthHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBoardingRosterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScanBoardingHandler).GetBoardingRoster(ctx, req.(*GetBoardingRosterRequest))
	}
	if interceptor == nil {
		return handler(ctx, in)
	}
	return interceptor(ctx, in, &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_GetBoardingRoster_FullMethodName}, handler)
}

// ---------------------------------------------------------------------------
// RegisterOpsServiceServerWithDepth
// Registers the full OpsService surface: Healthz + scan/boarding (Phase 6) + depth (Wave 5).
// ---------------------------------------------------------------------------

// RegisterOpsServiceServerWithDepth registers all ops RPCs including Wave 5 depth methods.
func RegisterOpsServiceServerWithDepth(s grpc.ServiceRegistrar, srv interface {
	OpsServiceServer
	OpsServiceHandler
	ScanBoardingHandler
	OpsDepthHandler
}) {
	desc := grpc.ServiceDesc{
		ServiceName: "pb.OpsService",
		HandlerType: (*OpsServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Healthz",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					in := new(HealthzRequest)
					if err := dec(in); err != nil {
						return nil, err
					}
					handler := func(ctx context.Context, req interface{}) (interface{}, error) {
						return srv.(OpsServiceServer).Healthz(ctx, req.(*HealthzRequest))
					}
					if interceptor == nil {
						return handler(ctx, in)
					}
					info := &grpc.UnaryServerInfo{Server: srv, FullMethod: OpsService_Healthz_FullMethodName}
					return interceptor(ctx, in, info, handler)
				},
			},
			// Phase 6 room/id-card/manifest
			OpsService_RunRoomAllocation_MethodDesc,
			OpsService_GetRoomAllocation_MethodDesc,
			OpsService_GenerateIDCard_MethodDesc,
			OpsService_VerifyIDCard_MethodDesc,
			OpsService_ExportManifest_MethodDesc,
			// Phase 6 scan/boarding
			{MethodName: "RecordScan", Handler: _OpsService_RecordScan_DepthHandler},
			{MethodName: "RecordBusBoarding", Handler: _OpsService_RecordBusBoarding_DepthHandler},
			{MethodName: "GetBoardingRoster", Handler: _OpsService_GetBoardingRoster_DepthHandler},
			// Wave 5 depth
			{MethodName: "StoreCollectiveDocument", Handler: _OpsService_StoreCollectiveDocument_Handler},
			{MethodName: "GetCollectiveDocuments", Handler: _OpsService_GetCollectiveDocuments_Handler},
			{MethodName: "SetDocumentACL", Handler: _OpsService_SetDocumentACL_Handler},
			{MethodName: "ExtractPassportOCR", Handler: _OpsService_ExtractPassportOCR_Handler},
			{MethodName: "SetMahramRelation", Handler: _OpsService_SetMahramRelation_Handler},
			{MethodName: "GetMahramRelations", Handler: _OpsService_GetMahramRelations_Handler},
			{MethodName: "GetDocumentProgress", Handler: _OpsService_GetDocumentProgress_Handler},
			{MethodName: "GetExpiryAlerts", Handler: _OpsService_GetExpiryAlerts_Handler},
			{MethodName: "GenerateOfficialLetter", Handler: _OpsService_GenerateOfficialLetter_Handler},
			{MethodName: "GenerateImmigrationManifest", Handler: _OpsService_GenerateImmigrationManifest_Handler},
			{MethodName: "RunSmartRooming", Handler: _OpsService_RunSmartRooming_Handler},
			{MethodName: "AssignTransport", Handler: _OpsService_AssignTransport_Handler},
			{MethodName: "GetTransportAssignments", Handler: _OpsService_GetTransportAssignments_Handler},
			{MethodName: "PublishManifestDelta", Handler: _OpsService_PublishManifestDelta_Handler},
			{MethodName: "AssignStaff", Handler: _OpsService_AssignStaff_Handler},
			{MethodName: "RecordPassportHandover", Handler: _OpsService_RecordPassportHandover_Handler},
			{MethodName: "GetPassportLog", Handler: _OpsService_GetPassportLog_Handler},
			{MethodName: "GetVisaProgress", Handler: _OpsService_GetVisaProgress_Handler},
			{MethodName: "StoreEVisa", Handler: _OpsService_StoreEVisa_Handler},
			{MethodName: "GetEVisa", Handler: _OpsService_GetEVisa_Handler},
			{MethodName: "TriggerExternalProvider", Handler: _OpsService_TriggerExternalProvider_Handler},
			{MethodName: "CreateRefund", Handler: _OpsService_CreateRefund_Handler},
			{MethodName: "ApproveRefund", Handler: _OpsService_ApproveRefund_Handler},
			{MethodName: "RecordPenalty", Handler: _OpsService_RecordPenalty_Handler},
			{MethodName: "RecordLuggageScan", Handler: _OpsService_RecordLuggageScan_Handler},
			{MethodName: "GetLuggageCount", Handler: _OpsService_GetLuggageCount_Handler},
			{MethodName: "BroadcastSchedule", Handler: _OpsService_BroadcastSchedule_Handler},
			{MethodName: "IssueDigitalTasreh", Handler: _OpsService_IssueDigitalTasreh_Handler},
			{MethodName: "RecordRaudhahEntry", Handler: _OpsService_RecordRaudhahEntry_Handler},
			{MethodName: "RegisterAudioDevice", Handler: _OpsService_RegisterAudioDevice_Handler},
			{MethodName: "UpdateAudioDeviceStatus", Handler: _OpsService_UpdateAudioDeviceStatus_Handler},
			{MethodName: "RecordZamzamDistribution", Handler: _OpsService_RecordZamzamDistribution_Handler},
			{MethodName: "RecordRoomCheckIn", Handler: _OpsService_RecordRoomCheckIn_Handler},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "ops.proto",
	}
	s.RegisterService(&desc, srv)
}
