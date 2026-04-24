package service

import (
	"context"

	"ops-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for ops-svc.
//
//   - Liveness — process is up
//   - Readiness — process is up AND the database is reachable
//   - DbTxDiagnostic — writes + reads inside a WithTx, canonical transaction reference
//   - RunRoomAllocation — run room grouping algorithm (BL-OPS-002)
//   - GetRoomAllocation — retrieve current allocation for a departure (BL-OPS-002)
//   - GenerateIDCard — generate HMAC-signed QR token (BL-OPS-003)
//   - VerifyIDCard — verify HMAC-signed QR token (BL-OPS-003)
//   - ExportManifest — return manifest rows for a departure (BL-OPS-001)
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)
	DbTxDiagnostic(ctx context.Context, params *DbTxDiagnosticParams) (*DbTxDiagnosticResult, error)

	// RunRoomAllocation groups jamaah_ids into rooms for a departure.
	// Idempotent: re-runs if draft, errors if committed.
	RunRoomAllocation(ctx context.Context, params *RunRoomAllocationParams) (*RunRoomAllocationResult, error)

	// GetRoomAllocation returns the current room allocation for a departure.
	GetRoomAllocation(ctx context.Context, params *GetRoomAllocationParams) (*GetRoomAllocationResult, error)

	// GenerateIDCard creates (or refreshes) an HMAC-signed token for an ID card or
	// luggage tag and stores it in ops.id_card_issuances.
	GenerateIDCard(ctx context.Context, params *GenerateIDCardParams) (*GenerateIDCardResult, error)

	// VerifyIDCard verifies an HMAC-signed ID card token.
	// Returns valid=false with ErrorReason on tamper; infrastructure errors are returned as error.
	VerifyIDCard(ctx context.Context, params *VerifyIDCardParams) (*VerifyIDCardResult, error)

	// ExportManifest returns structured manifest rows for a departure.
	// Currently returns empty rows; real data wired via jamaah-svc in a later sprint.
	ExportManifest(ctx context.Context, params *ExportManifestParams) (*ExportManifestResult, error)

	// RecordScan records a scan event idempotently (BL-OPS-010).
	RecordScan(ctx context.Context, params *RecordScanParams) (*RecordScanResult, error)

	// RecordBusBoarding records a bus boarding event atomically (BL-OPS-011).
	RecordBusBoarding(ctx context.Context, params *RecordBusBoardingParams) (*RecordBusBoardingResult, error)

	// GetBoardingRoster retrieves the boarding status roster for a departure (BL-OPS-011).
	GetBoardingRoster(ctx context.Context, params *GetBoardingRosterParams) (*GetBoardingRosterResult, error)

	// Wave 5 depth RPCs (BL-OPS-021..042).
	StoreCollectiveDocument(ctx context.Context, params *StoreCollectiveDocumentParams) (*StoreCollectiveDocumentResult, error)
	GetCollectiveDocuments(ctx context.Context, params *GetCollectiveDocumentsParams) (*GetCollectiveDocumentsResult, error)
	SetDocumentACL(ctx context.Context, params *SetDocumentACLParams) (*SetDocumentACLResult, error)
	ExtractPassportOCR(ctx context.Context, params *ExtractPassportOCRParams) (*ExtractPassportOCRResult, error)
	SetMahramRelation(ctx context.Context, params *SetMahramRelationParams) (*SetMahramRelationResult, error)
	GetMahramRelations(ctx context.Context, params *GetMahramRelationsParams) (*GetMahramRelationsResult, error)
	GetDocumentProgress(ctx context.Context, params *GetDocumentProgressParams) (*GetDocumentProgressResult, error)
	GetExpiryAlerts(ctx context.Context, params *GetExpiryAlertsParams) (*GetExpiryAlertsResult, error)
	GenerateOfficialLetter(ctx context.Context, params *GenerateOfficialLetterParams) (*GenerateOfficialLetterResult, error)
	GenerateImmigrationManifest(ctx context.Context, params *GenerateImmigrationManifestParams) (*GenerateImmigrationManifestResult, error)
	RunSmartRooming(ctx context.Context, params *RunSmartRoomingParams) (*RunSmartRoomingResult, error)
	AssignTransport(ctx context.Context, params *AssignTransportParams) (*AssignTransportResult, error)
	GetTransportAssignments(ctx context.Context, params *GetTransportAssignmentsParams) (*GetTransportAssignmentsResult, error)
	PublishManifestDelta(ctx context.Context, params *PublishManifestDeltaParams) (*PublishManifestDeltaResult, error)
	AssignStaff(ctx context.Context, params *AssignStaffParams) (*AssignStaffResult, error)
	RecordPassportHandover(ctx context.Context, params *RecordPassportHandoverParams) (*RecordPassportHandoverResult, error)
	GetPassportLog(ctx context.Context, params *GetPassportLogParams) (*GetPassportLogResult, error)
	GetVisaProgress(ctx context.Context, params *GetVisaProgressParams) (*GetVisaProgressResult, error)
	StoreEVisa(ctx context.Context, params *StoreEVisaParams) (*StoreEVisaResult, error)
	GetEVisa(ctx context.Context, params *GetEVisaParams) (*GetEVisaResult, error)
	TriggerExternalProvider(ctx context.Context, params *TriggerExternalProviderParams) (*TriggerExternalProviderResult, error)
	CreateRefund(ctx context.Context, params *CreateRefundParams) (*CreateRefundResult, error)
	ApproveRefund(ctx context.Context, params *ApproveRefundParams) (*ApproveRefundResult, error)
	RecordPenalty(ctx context.Context, params *RecordPenaltyParams) (*RecordPenaltyResult, error)
	RecordLuggageScan(ctx context.Context, params *RecordLuggageScanParams) (*RecordLuggageScanResult, error)
	GetLuggageCount(ctx context.Context, params *GetLuggageCountParams) (*GetLuggageCountResult, error)
	BroadcastSchedule(ctx context.Context, params *BroadcastScheduleParams) (*BroadcastScheduleResult, error)
	IssueDigitalTasreh(ctx context.Context, params *IssueDigitalTasrehParams) (*IssueDigitalTasrehResult, error)
	RecordRaudhahEntry(ctx context.Context, params *RecordRaudhahEntryParams) (*RecordRaudhahEntryResult, error)
	RegisterAudioDevice(ctx context.Context, params *RegisterAudioDeviceParams) (*RegisterAudioDeviceResult, error)
	UpdateAudioDeviceStatus(ctx context.Context, params *UpdateAudioDeviceStatusParams) (*UpdateAudioDeviceStatusResult, error)
	RecordZamzamDistribution(ctx context.Context, params *RecordZamzamDistributionParams) (*RecordZamzamDistributionResult, error)
	RecordRoomCheckIn(ctx context.Context, params *RecordRoomCheckInParams) (*RecordRoomCheckInResult, error)
}

type Service struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	// appName is stamped on every row this service writes to the shared
	// public.diagnostics table so debugging can attribute rows to their origin.
	appName string

	store postgres_store.IStore
}

func NewService(
	logger *zerolog.Logger,
	tracer trace.Tracer,
	appName string,
	store postgres_store.IStore,
) IService {
	return &Service{
		logger:  logger,
		tracer:  tracer,
		appName: appName,
		store:   store,
	}
}
