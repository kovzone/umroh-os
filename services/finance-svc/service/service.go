package service

import (
	"context"

	"finance-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for finance-svc.
//
// Pilot scaffold surfaces the three standard scaffold endpoints:
//
//   - Liveness — process is up
//   - Readiness — process is up AND the database is reachable
//   - DbTxDiagnostic — writes + reads inside a WithTx, the canonical reference
//     for how services should use transactions (per docs/04-backend-conventions)
//
// BL-IAM-002 adds FinancePing — the placeholder authenticated route that
// exercises the iam-svc permission gate so the "finance routes denied for
// non-finance roles" acceptance has a concrete surface. Real finance endpoints
// (journals, AR/AP, reports) land with S3-E-03 + S3-E-07.
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)
	DbTxDiagnostic(ctx context.Context, params *DbTxDiagnosticParams) (*DbTxDiagnosticResult, error)

	// Finance — BL-IAM-002 placeholder.
	FinancePing(ctx context.Context, params *FinancePingParams) (*FinancePingResult, error)
}

type Service struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	// appName is stamped on every row this service writes to the shared
	// public.diagnostics table so debugging can attribute rows to their origin.
	appName string

	store postgres_store.IStore

	// iamChecker is the consumer-side wrapper around iam-svc.CheckPermission.
	// Injected so tests can supply a double without spinning up a real gRPC server.
	iamChecker IamChecker
}

func NewService(
	logger *zerolog.Logger,
	tracer trace.Tracer,
	appName string,
	store postgres_store.IStore,
	iamChecker IamChecker,
) IService {
	return &Service{
		logger:     logger,
		tracer:     tracer,
		appName:    appName,
		store:      store,
		iamChecker: iamChecker,
	}
}
