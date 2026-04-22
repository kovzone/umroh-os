package service

import (
	"context"

	"gateway-svc/adapter/booking_rest_adapter"
	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/adapter/catalog_rest_adapter"
	"gateway-svc/adapter/crm_rest_adapter"
	"gateway-svc/adapter/finance_rest_adapter"
	"gateway-svc/adapter/iam_rest_adapter"
	"gateway-svc/adapter/jamaah_rest_adapter"
	"gateway-svc/adapter/logistics_rest_adapter"
	"gateway-svc/adapter/ops_rest_adapter"
	"gateway-svc/adapter/payment_rest_adapter"
	"gateway-svc/adapter/visa_rest_adapter"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for gateway-svc.
//
// gateway-svc fronts every backend over REST. Methods either return
// process-local state (Liveness, Readiness) or dispatch to a backend through
// the corresponding REST adapter. Each Get<Svc>SystemLive proves end-to-end
// that the adapter pattern works for that backend; richer per-route methods
// (auth, packages, bookings, ...) land alongside their owning feature work.
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)

	// Per-backend liveness proxies — called by the web app's status page.
	GetIamSystemLive(ctx context.Context) (*iam_rest_adapter.LivenessResult, error)

	// Traced cross-service path — the S0-J-05 observability acceptance check
	// flows through gateway-svc → iam-svc here.
	GetIamSystemDbTxDiagnostic(ctx context.Context, message string) (*iam_rest_adapter.DbTxDiagnosticResult, error)

	GetCatalogSystemLive(ctx context.Context) (*catalog_rest_adapter.LivenessResult, error)

	// Public catalog read (BL-GTW-002 / S1-E-10) — proxies to catalog-svc
	// gRPC via catalog_grpc_adapter. Mirrors GET /v1/packages,
	// /v1/packages/{id}, /v1/package-departures/{id}.
	ListPackages(ctx context.Context, params *catalog_grpc_adapter.ListPackagesParams) (*catalog_grpc_adapter.ListPackagesResult, error)
	GetPackage(ctx context.Context, params *catalog_grpc_adapter.GetPackageParams) (*catalog_grpc_adapter.PackageDetail, error)
	GetPackageDeparture(ctx context.Context, params *catalog_grpc_adapter.GetPackageDepartureParams) (*catalog_grpc_adapter.DepartureDetail, error)

	GetBookingSystemLive(ctx context.Context) (*booking_rest_adapter.LivenessResult, error)
	GetJamaahSystemLive(ctx context.Context) (*jamaah_rest_adapter.LivenessResult, error)
	GetPaymentSystemLive(ctx context.Context) (*payment_rest_adapter.LivenessResult, error)
	GetVisaSystemLive(ctx context.Context) (*visa_rest_adapter.LivenessResult, error)
	GetOpsSystemLive(ctx context.Context) (*ops_rest_adapter.LivenessResult, error)
	GetLogisticsSystemLive(ctx context.Context) (*logistics_rest_adapter.LivenessResult, error)
	GetFinanceSystemLive(ctx context.Context) (*finance_rest_adapter.LivenessResult, error)
	GetCrmSystemLive(ctx context.Context) (*crm_rest_adapter.LivenessResult, error)
}

// Adapters bundles the adapters this service can dispatch through.
// One field per backend; populated at construction time in cmd/start.go.
// Mixed REST + gRPC during the ADR-0009 transition — each backend
// graduates to gRPC-only as its BL-REFACTOR-* card lands.
type Adapters struct {
	iamRest       *iam_rest_adapter.Adapter
	catalogRest   *catalog_rest_adapter.Adapter
	catalogGrpc   *catalog_grpc_adapter.Adapter
	bookingRest   *booking_rest_adapter.Adapter
	jamaahRest    *jamaah_rest_adapter.Adapter
	paymentRest   *payment_rest_adapter.Adapter
	visaRest      *visa_rest_adapter.Adapter
	opsRest       *ops_rest_adapter.Adapter
	logisticsRest *logistics_rest_adapter.Adapter
	financeRest   *finance_rest_adapter.Adapter
	crmRest       *crm_rest_adapter.Adapter
}

type Service struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	appName  string
	adapters Adapters
}

// NewServiceParams keeps the constructor readable as the adapter list grows.
type NewServiceParams struct {
	Logger        *zerolog.Logger
	Tracer        trace.Tracer
	AppName       string
	IamRest       *iam_rest_adapter.Adapter
	CatalogRest   *catalog_rest_adapter.Adapter
	CatalogGrpc   *catalog_grpc_adapter.Adapter
	BookingRest   *booking_rest_adapter.Adapter
	JamaahRest    *jamaah_rest_adapter.Adapter
	PaymentRest   *payment_rest_adapter.Adapter
	VisaRest      *visa_rest_adapter.Adapter
	OpsRest       *ops_rest_adapter.Adapter
	LogisticsRest *logistics_rest_adapter.Adapter
	FinanceRest   *finance_rest_adapter.Adapter
	CrmRest       *crm_rest_adapter.Adapter
}

func NewService(p NewServiceParams) IService {
	return &Service{
		logger:  p.Logger,
		tracer:  p.Tracer,
		appName: p.AppName,
		adapters: Adapters{
			iamRest:       p.IamRest,
			catalogRest:   p.CatalogRest,
			catalogGrpc:   p.CatalogGrpc,
			bookingRest:   p.BookingRest,
			jamaahRest:    p.JamaahRest,
			paymentRest:   p.PaymentRest,
			visaRest:      p.VisaRest,
			opsRest:       p.OpsRest,
			logisticsRest: p.LogisticsRest,
			financeRest:   p.FinanceRest,
			crmRest:       p.CrmRest,
		},
	}
}
