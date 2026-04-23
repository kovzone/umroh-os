package service

import (
	"context"

	"jamaah-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for jamaah-svc.
//
// Pilot scaffold surfaces the three standard scaffold endpoints plus the
// S3-E-02 document upload/review scaffold:
//
//   - Liveness — process is up
//   - Readiness — process is up AND the database is reachable
//   - DbTxDiagnostic — canonical WithTx reference
//   - UploadDocument — store a document record in status='pending' (F3-W1 / BL-DOC-001)
//   - ReviewDocument — approve or reject a document (F3-W3 / BL-DOC-003)
//   - TriggerOCR — run stub OCR on a document (BL-DOC-002)
//   - GetOCRStatus — retrieve OCR results for a document (BL-DOC-002)
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)
	DbTxDiagnostic(ctx context.Context, params *DbTxDiagnosticParams) (*DbTxDiagnosticResult, error)

	// Document upload (S3-E-02 / F3-W1 / BL-DOC-001)
	UploadDocument(ctx context.Context, params *UploadDocumentParams) (*UploadDocumentResult, error)
	// Document review (S3-E-02 / F3-W3 / BL-DOC-003)
	ReviewDocument(ctx context.Context, params *ReviewDocumentParams) (*ReviewDocumentResult, error)

	// Departure manifest (Wave-1A)
	GetDepartureManifest(ctx context.Context, params *GetDepartureManifestParams) (*GetDepartureManifestResult, error)

	// OCR stub (BL-DOC-002 / migration 000020)
	TriggerOCR(ctx context.Context, params *TriggerOCRParams) (*TriggerOCRResult, error)
	GetOCRStatus(ctx context.Context, params *GetOCRStatusParams) (*GetOCRStatusResult, error)
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
