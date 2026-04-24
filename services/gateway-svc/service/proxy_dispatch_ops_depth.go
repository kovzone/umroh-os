// proxy_dispatch_ops_depth.go — gateway service dispatch for ops-svc Wave 5 depth RPCs.
// BL-OPS-021..042: collective docs, passport OCR, mahram, document progress,
// official letters, immigration manifest, transport, manifest delta, staff assignment,
// passport log, visa progress, e-visa, external provider, refunds, luggage, broadcast,
// tasreh, raudhah, audio devices, zamzam, room check-in.
//
// Each method is a thin delegation to ops_grpc_adapter. No business logic here.
package service

import (
	"context"

	"gateway-svc/adapter/ops_grpc_adapter"
)

// ---------------------------------------------------------------------------
// Ops depth — Wave 5 (BL-OPS-021..042)
// ---------------------------------------------------------------------------

func (s *Service) StoreCollectiveDocument(ctx context.Context, params *ops_grpc_adapter.StoreCollectiveDocumentParams) (*ops_grpc_adapter.StoreCollectiveDocumentResult, error) {
	return s.adapters.opsGrpc.StoreCollectiveDocument(ctx, params)
}

func (s *Service) GetCollectiveDocuments(ctx context.Context, params *ops_grpc_adapter.GetCollectiveDocumentsParams) (*ops_grpc_adapter.GetCollectiveDocumentsResult, error) {
	return s.adapters.opsGrpc.GetCollectiveDocuments(ctx, params)
}

func (s *Service) SetDocumentACL(ctx context.Context, params *ops_grpc_adapter.SetDocumentACLParams) (*ops_grpc_adapter.SetDocumentACLResult, error) {
	return s.adapters.opsGrpc.SetDocumentACL(ctx, params)
}

func (s *Service) ExtractPassportOCR(ctx context.Context, params *ops_grpc_adapter.ExtractPassportOCRParams) (*ops_grpc_adapter.ExtractPassportOCRResult, error) {
	return s.adapters.opsGrpc.ExtractPassportOCR(ctx, params)
}

func (s *Service) SetMahramRelation(ctx context.Context, params *ops_grpc_adapter.SetMahramRelationParams) (*ops_grpc_adapter.SetMahramRelationResult, error) {
	return s.adapters.opsGrpc.SetMahramRelation(ctx, params)
}

func (s *Service) GetMahramRelations(ctx context.Context, params *ops_grpc_adapter.GetMahramRelationsParams) (*ops_grpc_adapter.GetMahramRelationsResult, error) {
	return s.adapters.opsGrpc.GetMahramRelations(ctx, params)
}

func (s *Service) GetDocumentProgress(ctx context.Context, params *ops_grpc_adapter.GetDocumentProgressParams) (*ops_grpc_adapter.GetDocumentProgressResult, error) {
	return s.adapters.opsGrpc.GetDocumentProgress(ctx, params)
}

func (s *Service) GetExpiryAlerts(ctx context.Context, params *ops_grpc_adapter.GetExpiryAlertsParams) (*ops_grpc_adapter.GetExpiryAlertsResult, error) {
	return s.adapters.opsGrpc.GetExpiryAlerts(ctx, params)
}

func (s *Service) GenerateOfficialLetter(ctx context.Context, params *ops_grpc_adapter.GenerateOfficialLetterParams) (*ops_grpc_adapter.GenerateOfficialLetterResult, error) {
	return s.adapters.opsGrpc.GenerateOfficialLetter(ctx, params)
}

func (s *Service) GenerateImmigrationManifest(ctx context.Context, params *ops_grpc_adapter.GenerateImmigrationManifestParams) (*ops_grpc_adapter.GenerateImmigrationManifestResult, error) {
	return s.adapters.opsGrpc.GenerateImmigrationManifest(ctx, params)
}

func (s *Service) AssignTransport(ctx context.Context, params *ops_grpc_adapter.AssignTransportParams) (*ops_grpc_adapter.AssignTransportResult, error) {
	return s.adapters.opsGrpc.AssignTransport(ctx, params)
}

func (s *Service) GetTransportAssignments(ctx context.Context, params *ops_grpc_adapter.GetTransportAssignmentsParams) (*ops_grpc_adapter.GetTransportAssignmentsResult, error) {
	return s.adapters.opsGrpc.GetTransportAssignments(ctx, params)
}

func (s *Service) PublishManifestDelta(ctx context.Context, params *ops_grpc_adapter.PublishManifestDeltaParams) (*ops_grpc_adapter.PublishManifestDeltaResult, error) {
	return s.adapters.opsGrpc.PublishManifestDelta(ctx, params)
}

func (s *Service) AssignStaff(ctx context.Context, params *ops_grpc_adapter.AssignStaffParams) (*ops_grpc_adapter.AssignStaffResult, error) {
	return s.adapters.opsGrpc.AssignStaff(ctx, params)
}

func (s *Service) RecordPassportHandover(ctx context.Context, params *ops_grpc_adapter.RecordPassportHandoverParams) (*ops_grpc_adapter.RecordPassportHandoverResult, error) {
	return s.adapters.opsGrpc.RecordPassportHandover(ctx, params)
}

func (s *Service) GetPassportLog(ctx context.Context, params *ops_grpc_adapter.GetPassportLogParams) (*ops_grpc_adapter.GetPassportLogResult, error) {
	return s.adapters.opsGrpc.GetPassportLog(ctx, params)
}

func (s *Service) GetVisaProgress(ctx context.Context, params *ops_grpc_adapter.GetVisaProgressParams) (*ops_grpc_adapter.GetVisaProgressResult, error) {
	return s.adapters.opsGrpc.GetVisaProgress(ctx, params)
}

func (s *Service) StoreEVisa(ctx context.Context, params *ops_grpc_adapter.StoreEVisaParams) (*ops_grpc_adapter.StoreEVisaResult, error) {
	return s.adapters.opsGrpc.StoreEVisa(ctx, params)
}

func (s *Service) GetEVisa(ctx context.Context, params *ops_grpc_adapter.GetEVisaParams) (*ops_grpc_adapter.GetEVisaResult, error) {
	return s.adapters.opsGrpc.GetEVisa(ctx, params)
}

func (s *Service) TriggerExternalProvider(ctx context.Context, params *ops_grpc_adapter.TriggerExternalProviderParams) (*ops_grpc_adapter.TriggerExternalProviderResult, error) {
	return s.adapters.opsGrpc.TriggerExternalProvider(ctx, params)
}

func (s *Service) CreateRefund(ctx context.Context, params *ops_grpc_adapter.CreateRefundParams) (*ops_grpc_adapter.CreateRefundResult, error) {
	return s.adapters.opsGrpc.CreateRefund(ctx, params)
}

func (s *Service) ApproveRefund(ctx context.Context, params *ops_grpc_adapter.ApproveRefundParams) (*ops_grpc_adapter.ApproveRefundResult, error) {
	return s.adapters.opsGrpc.ApproveRefund(ctx, params)
}

func (s *Service) RecordPenalty(ctx context.Context, params *ops_grpc_adapter.RecordPenaltyParams) (*ops_grpc_adapter.RecordPenaltyResult, error) {
	return s.adapters.opsGrpc.RecordPenalty(ctx, params)
}

func (s *Service) RecordLuggageScan(ctx context.Context, params *ops_grpc_adapter.RecordLuggageScanParams) (*ops_grpc_adapter.RecordLuggageScanResult, error) {
	return s.adapters.opsGrpc.RecordLuggageScan(ctx, params)
}

func (s *Service) GetLuggageCount(ctx context.Context, params *ops_grpc_adapter.GetLuggageCountParams) (*ops_grpc_adapter.GetLuggageCountResult, error) {
	return s.adapters.opsGrpc.GetLuggageCount(ctx, params)
}

func (s *Service) BroadcastSchedule(ctx context.Context, params *ops_grpc_adapter.BroadcastScheduleParams) (*ops_grpc_adapter.BroadcastScheduleResult, error) {
	return s.adapters.opsGrpc.BroadcastSchedule(ctx, params)
}

func (s *Service) IssueDigitalTasreh(ctx context.Context, params *ops_grpc_adapter.IssueDigitalTasrehParams) (*ops_grpc_adapter.IssueDigitalTasrehResult, error) {
	return s.adapters.opsGrpc.IssueDigitalTasreh(ctx, params)
}

func (s *Service) RecordRaudhahEntry(ctx context.Context, params *ops_grpc_adapter.RecordRaudhahEntryParams) (*ops_grpc_adapter.RecordRaudhahEntryResult, error) {
	return s.adapters.opsGrpc.RecordRaudhahEntry(ctx, params)
}

func (s *Service) RegisterAudioDevice(ctx context.Context, params *ops_grpc_adapter.RegisterAudioDeviceParams) (*ops_grpc_adapter.RegisterAudioDeviceResult, error) {
	return s.adapters.opsGrpc.RegisterAudioDevice(ctx, params)
}

func (s *Service) UpdateAudioDeviceStatus(ctx context.Context, params *ops_grpc_adapter.UpdateAudioDeviceStatusParams) (*ops_grpc_adapter.UpdateAudioDeviceStatusResult, error) {
	return s.adapters.opsGrpc.UpdateAudioDeviceStatus(ctx, params)
}

func (s *Service) RecordZamzamDistribution(ctx context.Context, params *ops_grpc_adapter.RecordZamzamDistributionParams) (*ops_grpc_adapter.RecordZamzamDistributionResult, error) {
	return s.adapters.opsGrpc.RecordZamzamDistribution(ctx, params)
}

func (s *Service) RecordRoomCheckIn(ctx context.Context, params *ops_grpc_adapter.RecordRoomCheckInParams) (*ops_grpc_adapter.RecordRoomCheckInResult, error) {
	return s.adapters.opsGrpc.RecordRoomCheckIn(ctx, params)
}
