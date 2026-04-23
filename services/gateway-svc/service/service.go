package service

import (
	"context"

	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/adapter/finance_rest_adapter"
	"gateway-svc/adapter/iam_rest_adapter"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for gateway-svc.
//
// gateway-svc fronts every backend. Methods either return process-local state
// (Liveness, Readiness) or dispatch to a backend through the corresponding
// adapter. Under ADR 0009 the target shape is REST-in, gRPC-out; the iam and
// finance REST adapters are the remaining interim surfaces (retired by
// BL-IAM-018 / S1-E-12 and BL-IAM-019 / S1-E-14 respectively).
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)

	// Per-backend liveness proxies — called by the web app's status page.
	GetIamSystemLive(ctx context.Context) (*iam_rest_adapter.LivenessResult, error)

	// Traced cross-service path — the S0-J-05 observability acceptance check
	// flows through gateway-svc → iam-svc here.
	GetIamSystemDbTxDiagnostic(ctx context.Context, message string) (*iam_rest_adapter.DbTxDiagnosticResult, error)

	// Public catalog read (BL-GTW-002 / S1-E-10) — proxies to catalog-svc
	// gRPC via catalog_grpc_adapter. Mirrors GET /v1/packages,
	// /v1/packages/{id}, /v1/package-departures/{id}.
	ListPackages(ctx context.Context, params *catalog_grpc_adapter.ListPackagesParams) (*catalog_grpc_adapter.ListPackagesResult, error)
	GetPackage(ctx context.Context, params *catalog_grpc_adapter.GetPackageParams) (*catalog_grpc_adapter.PackageDetail, error)
	GetPackageDeparture(ctx context.Context, params *catalog_grpc_adapter.GetPackageDepartureParams) (*catalog_grpc_adapter.DepartureDetail, error)

	// Finance liveness proxy — interim REST adapter; retires with
	// BL-IAM-019 / S1-E-14 when /v1/finance/* moves to gRPC.
	GetFinanceSystemLive(ctx context.Context) (*finance_rest_adapter.LivenessResult, error)
}

// Adapters bundles the adapters this service can dispatch through.
// One field per backend; populated at construction time in cmd/start.go.
// Mixed REST + gRPC during the ADR-0009 transition — each backend
// graduates to gRPC-only as its BL-REFACTOR-* card lands.
type Adapters struct {
	iamRest     *iam_rest_adapter.Adapter
	catalogGrpc *catalog_grpc_adapter.Adapter
	financeRest *finance_rest_adapter.Adapter
}

type Service struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	appName  string
	adapters Adapters
}

// NewServiceParams keeps the constructor readable as the adapter list grows.
type NewServiceParams struct {
	Logger      *zerolog.Logger
	Tracer      trace.Tracer
	AppName     string
	IamRest     *iam_rest_adapter.Adapter
	CatalogGrpc *catalog_grpc_adapter.Adapter
	FinanceRest *finance_rest_adapter.Adapter
}

func NewService(p NewServiceParams) IService {
	return &Service{
		logger:  p.Logger,
		tracer:  p.Tracer,
		appName: p.AppName,
		adapters: Adapters{
			iamRest:     p.IamRest,
			catalogGrpc: p.CatalogGrpc,
			financeRest: p.FinanceRest,
		},
	}
}
