package service

import (
	"context"

	"visa-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for visa-svc.
//
// Covers scaffold lifecycle checks plus the three Phase 6 visa pipeline RPCs:
//   - TransitionStatus (BL-VISA-001) — state-machine-driven single application transition
//   - BulkSubmit       (BL-VISA-002) — atomic all-or-nothing batch READY→SUBMITTED
//   - GetApplications  (BL-VISA-003) — list applications for a departure with history
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)
	DbTxDiagnostic(ctx context.Context, params *DbTxDiagnosticParams) (*DbTxDiagnosticResult, error)

	// Visa pipeline (BL-VISA-001..003)
	TransitionStatus(ctx context.Context, params *TransitionStatusParams) (*TransitionStatusResult, error)
	BulkSubmit(ctx context.Context, params *BulkSubmitParams) (*BulkSubmitResult, error)
	GetApplications(ctx context.Context, params *GetApplicationsParams) (*GetApplicationsResult, error)
}

type Service struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	appName string
	store   postgres_store.IStore
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
