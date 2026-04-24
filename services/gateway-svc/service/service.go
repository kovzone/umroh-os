package service

import (
	"context"

	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/adapter/finance_grpc_adapter"
	"gateway-svc/adapter/health_check_adapter"
	"gateway-svc/adapter/iam_grpc_adapter"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for gateway-svc.
//
// gateway-svc fronts every backend. Methods either return process-local state
// (Liveness, Readiness) or dispatch to a backend through the corresponding
// adapter. Under ADR 0009 every downstream service is gRPC-only — after
// BL-IAM-019 / S1-E-14 finance-svc is the last backend to migrate and no
// REST adapter remains on the gateway side.
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)

	// Public catalog read (BL-GTW-002 / S1-E-10) — proxies to catalog-svc
	// gRPC via catalog_grpc_adapter. Mirrors GET /v1/packages,
	// /v1/packages/{id}, /v1/package-departures/{id}.
	ListPackages(ctx context.Context, params *catalog_grpc_adapter.ListPackagesParams) (*catalog_grpc_adapter.ListPackagesResult, error)
	GetPackage(ctx context.Context, params *catalog_grpc_adapter.GetPackageParams) (*catalog_grpc_adapter.PackageDetail, error)
	GetPackageDeparture(ctx context.Context, params *catalog_grpc_adapter.GetPackageDepartureParams) (*catalog_grpc_adapter.DepartureDetail, error)

	// Client-facing auth (BL-IAM-018 / S1-E-12) — proxies to iam-svc gRPC
	// via iam_grpc_adapter. The gateway bearer middleware enforces auth on
	// protected routes; iam-svc handlers trust the forwarded user_id /
	// session_id per ADR 0009 single-point-auth.
	Login(ctx context.Context, params *iam_grpc_adapter.LoginParams) (*iam_grpc_adapter.LoginResult, error)
	RefreshSession(ctx context.Context, params *iam_grpc_adapter.RefreshSessionParams) (*iam_grpc_adapter.RefreshSessionResult, error)
	Logout(ctx context.Context, params *iam_grpc_adapter.LogoutParams) (*iam_grpc_adapter.LogoutResult, error)
	GetMe(ctx context.Context, params *iam_grpc_adapter.GetMeParams) (*iam_grpc_adapter.GetMeResult, error)
	EnrollTOTP(ctx context.Context, params *iam_grpc_adapter.EnrollTOTPParams) (*iam_grpc_adapter.EnrollTOTPResult, error)
	VerifyTOTP(ctx context.Context, params *iam_grpc_adapter.VerifyTOTPParams) (*iam_grpc_adapter.VerifyTOTPResult, error)
	SuspendUser(ctx context.Context, params *iam_grpc_adapter.SuspendUserParams) (*iam_grpc_adapter.SuspendUserResult, error)

	// Permission-gate smoke (BL-IAM-019 / S1-E-14) — proxies to finance-svc
	// gRPC via finance_grpc_adapter. Authorization is enforced at the
	// gateway's RequireBearerToken + RequirePermission("journal_entry",
	// "read", "global") middleware chain; finance-svc is identity-agnostic.
	FinancePing(ctx context.Context) (*finance_grpc_adapter.PingResult, error)

	// Aggregate backend health (BL-IAM-019 / S1-E-14) — fans out
	// grpc.health.v1.Health.Check to every backend via health_check_adapter.
	// Unauthenticated; core-web's status page polls it.
	SystemBackends(ctx context.Context) []health_check_adapter.BackendStatus
}

// Adapters bundles the adapters this service can dispatch through.
// One field per backend; populated at construction time in cmd/start.go.
// Every adapter is gRPC-based — the REST era ended with BL-IAM-019 / S1-E-14.
type Adapters struct {
	iamGrpc     *iam_grpc_adapter.Adapter
	catalogGrpc *catalog_grpc_adapter.Adapter
	financeGrpc *finance_grpc_adapter.Adapter
	healthCheck *health_check_adapter.Adapter
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
	IamGrpc     *iam_grpc_adapter.Adapter
	CatalogGrpc *catalog_grpc_adapter.Adapter
	FinanceGrpc *finance_grpc_adapter.Adapter
	HealthCheck *health_check_adapter.Adapter
}

func NewService(p NewServiceParams) IService {
	return &Service{
		logger:  p.Logger,
		tracer:  p.Tracer,
		appName: p.AppName,
		adapters: Adapters{
			iamGrpc:     p.IamGrpc,
			catalogGrpc: p.CatalogGrpc,
			financeGrpc: p.FinanceGrpc,
			healthCheck: p.HealthCheck,
		},
	}
}
