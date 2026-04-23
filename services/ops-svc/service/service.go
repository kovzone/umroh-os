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
