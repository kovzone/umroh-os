package service

import (
	"context"

	"booking-svc/adapter/catalog_grpc_adapter"
	"booking-svc/adapter/iam_grpc_adapter"
	"booking-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for booking-svc.
//
// S1-E-03 adds CreateDraftBooking for BL-BOOK-001..006.
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)
	DbTxDiagnostic(ctx context.Context, params *DbTxDiagnosticParams) (*DbTxDiagnosticResult, error)

	// Booking draft (S1-E-03 / BL-BOOK-001..006)
	CreateDraftBooking(ctx context.Context, params *CreateDraftBookingParams) (*DraftBookingResult, error)
}

// IamClient is the slice of iam-svc the booking-svc service layer calls.
// Defined as an interface so tests can inject a mock.
type IamClient interface {
	ValidateToken(ctx context.Context, params *iam_grpc_adapter.ValidateTokenParams) (*iam_grpc_adapter.ValidateTokenResult, error)
	CheckPermission(ctx context.Context, params *iam_grpc_adapter.CheckPermissionParams) (*iam_grpc_adapter.CheckPermissionResult, error)
	RecordAudit(ctx context.Context, params *iam_grpc_adapter.RecordAuditParams) (*iam_grpc_adapter.RecordAuditResult, error)
}

// CatalogClient is the slice of catalog-svc the booking-svc service layer calls.
// Defined as an interface so tests can inject a mock.
type CatalogClient interface {
	GetDeparture(ctx context.Context, departureID string) (*catalog_grpc_adapter.GetDepartureResult, error)
	ReserveSeats(ctx context.Context, params *catalog_grpc_adapter.ReserveSeatsParams) (*catalog_grpc_adapter.ReserveSeatsResult, error)
	ReleaseSeats(ctx context.Context, params *catalog_grpc_adapter.ReleaseSeatsParams) (*catalog_grpc_adapter.ReleaseSeatsResult, error)
}

type Service struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	// appName is stamped on every row this service writes to the shared
	// public.diagnostics table so debugging can attribute rows to their origin.
	appName string

	store postgres_store.IStore

	// iamClient is the consumer-side wrapper around iam-svc's gRPC surface.
	iamClient IamClient

	// catalogClient is the consumer-side wrapper around catalog-svc's gRPC surface.
	// Used by CreateDraftBooking to validate departure and reserve seats.
	catalogClient CatalogClient
}

func NewService(
	logger *zerolog.Logger,
	tracer trace.Tracer,
	appName string,
	store postgres_store.IStore,
	iamClient IamClient,
	catalogClient CatalogClient,
) IService {
	return &Service{
		logger:        logger,
		tracer:        tracer,
		appName:       appName,
		store:         store,
		iamClient:     iamClient,
		catalogClient: catalogClient,
	}
}
