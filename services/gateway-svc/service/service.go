package service

import (
	"context"

	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/adapter/finance_rest_adapter"
	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/adapter/iam_rest_adapter"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for gateway-svc.
//
// gateway-svc fronts every backend. Methods either return process-local state
// (Liveness, Readiness) or dispatch to a backend through the corresponding
// adapter. Under ADR 0009 the target shape is REST-in, gRPC-out:
//   - catalog-svc: gRPC-only from S1-E-11 (catalog_grpc_adapter)
//   - iam-svc: gRPC-only from S1-E-12 (iam_grpc_adapter); iam_rest_adapter
//     retained only for interim probes until HTTP surface is fully removed.
//   - finance_rest_adapter: retires with BL-IAM-019 / S1-E-14.
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)

	// Per-backend liveness proxies — IAM still serves via REST adapter until
	// S1-E-12 removes the IAM REST surface.
	GetIamSystemLive(ctx context.Context) (*iam_rest_adapter.LivenessResult, error)

	// Traced cross-service path — the S0-J-05 observability acceptance check
	// flows through gateway-svc → iam-svc here.
	GetIamSystemDbTxDiagnostic(ctx context.Context, message string) (*iam_rest_adapter.DbTxDiagnosticResult, error)

	// catalog-svc is gRPC-only from S1-E-11; no REST liveness proxy.
	// Public catalog read (BL-GTW-002 / S1-E-10) — proxies to catalog-svc
	// gRPC via catalog_grpc_adapter. Mirrors GET /v1/packages,
	// /v1/packages/{id}, /v1/package-departures/{id}.
	ListPackages(ctx context.Context, params *catalog_grpc_adapter.ListPackagesParams) (*catalog_grpc_adapter.ListPackagesResult, error)
	GetPackage(ctx context.Context, params *catalog_grpc_adapter.GetPackageParams) (*catalog_grpc_adapter.PackageDetail, error)
	GetPackageDeparture(ctx context.Context, params *catalog_grpc_adapter.GetPackageDepartureParams) (*catalog_grpc_adapter.DepartureDetail, error)

	// Staff catalog write (BL-CAT-014 / S1-E-07) — bearer-protected; gateway
	// handler calls CheckPermission(catalog.package.manage) before delegating.
	CreatePackage(ctx context.Context, params *catalog_grpc_adapter.CreatePackageParams) (*catalog_grpc_adapter.PackageDetail, error)
	UpdatePackage(ctx context.Context, params *catalog_grpc_adapter.UpdatePackageParams) (*catalog_grpc_adapter.PackageDetail, error)
	DeletePackage(ctx context.Context, params *catalog_grpc_adapter.DeletePackageParams) (*catalog_grpc_adapter.DeletePackageResult, error)
	CreateDeparture(ctx context.Context, params *catalog_grpc_adapter.CreateDepartureParams) (*catalog_grpc_adapter.DepartureDetail, error)
	UpdateDeparture(ctx context.Context, params *catalog_grpc_adapter.UpdateDepartureParams) (*catalog_grpc_adapter.DepartureDetail, error)

	// IAM auth routes (BL-IAM-018 / S1-E-12) — proxied via iam_grpc_adapter.
	// Public (no bearer): Login, RefreshSession.
	// Bearer-required (via middleware): Logout, GetMe, EnrollTOTP, VerifyTOTP.
	// Bearer + permission: SuspendUser.
	CheckPermission(ctx context.Context, params *iam_grpc_adapter.CheckPermissionParams) (*iam_grpc_adapter.CheckPermissionResult, error)
	Login(ctx context.Context, params *iam_grpc_adapter.LoginParams) (*iam_grpc_adapter.LoginResult, error)
	RefreshSession(ctx context.Context, params *iam_grpc_adapter.RefreshSessionParams) (*iam_grpc_adapter.RefreshSessionResult, error)
	Logout(ctx context.Context, params *iam_grpc_adapter.LogoutParams) error
	GetMe(ctx context.Context, params *iam_grpc_adapter.GetMeParams) (*iam_grpc_adapter.GetMeResult, error)
	EnrollTOTP(ctx context.Context, params *iam_grpc_adapter.EnrollTOTPParams) (*iam_grpc_adapter.EnrollTOTPResult, error)
	VerifyTOTP(ctx context.Context, params *iam_grpc_adapter.VerifyTOTPParams) (*iam_grpc_adapter.VerifyTOTPResult, error)
	SuspendUser(ctx context.Context, params *iam_grpc_adapter.SuspendUserParams) (*iam_grpc_adapter.SuspendUserResult, error)

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
	iamGrpc     *iam_grpc_adapter.Adapter
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
	IamGrpc     *iam_grpc_adapter.Adapter
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
			iamGrpc:     p.IamGrpc,
			catalogGrpc: p.CatalogGrpc,
			financeRest: p.FinanceRest,
		},
	}
}
