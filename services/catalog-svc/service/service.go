package service

import (
	"context"

	"catalog-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for catalog-svc.
//
// S1-E-02 (read) + S1-E-07 (staff write + inventory) methods.
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)
	DbTxDiagnostic(ctx context.Context, params *DbTxDiagnosticParams) (*DbTxDiagnosticResult, error)

	// Public read (S1-E-02)
	GetPackages(ctx context.Context, params *GetPackagesParams) (*GetPackagesResult, error)
	GetPackageByID(ctx context.Context, params *GetPackageByIDParams) (*PackageDetail, error)
	GetDepartureByID(ctx context.Context, params *GetDepartureByIDParams) (*DepartureDetail, error)

	// Staff write (S1-E-07 / BL-CAT-014)
	CreatePackage(ctx context.Context, params *CreatePackageParams) (*PackageDetail, error)
	UpdatePackage(ctx context.Context, params *UpdatePackageParams) (*PackageDetail, error)
	DeletePackage(ctx context.Context, params *DeletePackageParams) error
	CreateDeparture(ctx context.Context, params *CreateDepartureParams) (*DepartureDetail, error)
	UpdateDeparture(ctx context.Context, params *UpdateDepartureParams) (*DepartureDetail, error)

	// Inventory (§ Inventory / S1-J-03)
	ReserveSeats(ctx context.Context, params *ReserveSeatsParams) (*ReserveSeatsResult, error)
	ReleaseSeats(ctx context.Context, params *ReleaseSeatsParams) (*ReleaseSeatsResult, error)
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
