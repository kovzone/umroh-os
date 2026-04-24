package service

import (
	"context"

	"logistics-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for logistics-svc.
//
// Pilot scaffold surfaces only the three standard scaffold endpoints:
//
//   - Liveness — process is up
//   - Readiness — process is up AND the database is reachable
//
// Real iam responsibilities (user/role/branch CRUD, auth login/refresh/logout,
// permission checks, session lifecycle, audit writes) land in F1.5–F1.11 and
// are deliberately out of scaffold scope.
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
