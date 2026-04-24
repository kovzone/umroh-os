package service

import (
	"context"

	"booking-svc/adapter/iam_grpc_adapter"
	"booking-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for booking-svc.
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

// IamClient is the slice of iam-svc the booking-svc service layer plans to
// call once handler-side BL-BKG-* cards land. Defined as an interface (not the
// concrete *iam_grpc_adapter.Adapter) so tests can inject a testify/mock
// double without pulling in the gRPC proto types. The adapter satisfies it.
//
// Scaffolded with BL-IAM-004; no method on IService consumes it yet. The first
// caller is S1-E-03 (booking draft create → RecordAudit + CheckPermission for
// create_on_behalf).
type IamClient interface {
	ValidateToken(ctx context.Context, params *iam_grpc_adapter.ValidateTokenParams) (*iam_grpc_adapter.ValidateTokenResult, error)
	CheckPermission(ctx context.Context, params *iam_grpc_adapter.CheckPermissionParams) (*iam_grpc_adapter.CheckPermissionResult, error)
	RecordAudit(ctx context.Context, params *iam_grpc_adapter.RecordAuditParams) (*iam_grpc_adapter.RecordAuditResult, error)
}

type Service struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	// appName is stamped on every row this service writes to the shared
	// public.diagnostics table so debugging can attribute rows to their origin.
	appName string

	store postgres_store.IStore

	// iamClient is the consumer-side wrapper around iam-svc's gRPC surface.
	// Held here so future booking handlers can call ValidateToken / CheckPermission
	// / RecordAudit without another constructor refactor when S1-E-03 lands.
	iamClient IamClient
}

func NewService(
	logger *zerolog.Logger,
	tracer trace.Tracer,
	appName string,
	store postgres_store.IStore,
	iamClient IamClient,
) IService {
	return &Service{
		logger:    logger,
		tracer:    tracer,
		appName:   appName,
		store:     store,
		iamClient: iamClient,
	}
}
