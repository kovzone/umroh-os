package service

import (
	"context"

	"finance-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for finance-svc.
//
// Post BL-IAM-019 / S1-E-14 the service is gRPC-only (ADR 0009). The
// permission-gate smoke (`FinancePing`) lives on the gRPC handler directly —
// its only job is to return `{message:"ok"}` once the gateway has already
// approved the caller via its RequireBearerToken + RequirePermission chain.
// Real finance endpoints (journals, AR/AP, reports) land with S3-E-03 +
// S3-E-07 and plug into this interface.
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)
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
