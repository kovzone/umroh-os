package service

import (
	"context"

	"catalog-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for catalog-svc.
//
// Current scope (S1-E-02 / BL-CAT-001 + BL-CAT-002):
//
//   - Liveness — process is up
//   - Readiness — process is up AND the database is reachable
//   - GetPackages — § Catalog public list (active only)
//   - GetPackageByID — § Catalog public detail (active only)
//   - GetDepartureByID — § Catalog departure detail with live remaining_seats
//
// Admin write endpoints (create/update/archive), bulk import/export,
// and the gRPC Reserve/ReleaseSeats pair land in later S1 cards.
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)

	GetPackages(ctx context.Context, params *GetPackagesParams) (*GetPackagesResult, error)
	GetPackageByID(ctx context.Context, params *GetPackageByIDParams) (*PackageDetail, error)
	GetDepartureByID(ctx context.Context, params *GetDepartureByIDParams) (*DepartureDetail, error)
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
