package service

import (
	"context"

	"booking-svc/adapter/catalog_grpc_adapter"
	"booking-svc/adapter/crm_grpc_adapter"
	"booking-svc/adapter/finance_grpc_adapter"
	"booking-svc/adapter/iam_grpc_adapter"
	"booking-svc/adapter/logistics_grpc_adapter"
	"booking-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for booking-svc.
//
// S1-E-03 adds CreateDraftBooking for BL-BOOK-001..006.
// S3 adds FanOutBookingPaid for the logistics + finance callback chain.
// S4-E-02 adds FanOutBookingCreated for the CRM callback (best-effort).
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)
	DbTxDiagnostic(ctx context.Context, params *DbTxDiagnosticParams) (*DbTxDiagnosticResult, error)

	// Booking draft (S1-E-03 / BL-BOOK-001..006)
	CreateDraftBooking(ctx context.Context, params *CreateDraftBookingParams) (*DraftBookingResult, error)

	// FanOutBookingPaid triggers downstream services after a booking is
	// confirmed as paid_in_full (S3-E-02 / S3-E-03).
	// Calls logistics-svc.OnBookingPaid and finance-svc.OnPaymentReceived
	// SYNCHRONOUSLY (ADR-0006 / §S3-J-01). Errors are returned to the caller
	// so the webhook pipeline can surface a 500 and trigger gateway retry.
	FanOutBookingPaid(ctx context.Context, params *FanOutBookingPaidParams) (*FanOutBookingPaidResult, error)
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

// LogisticsClient is the slice of logistics-svc the booking-svc service layer
// calls (S3-E-02). Defined as an interface so tests can inject a mock.
type LogisticsClient interface {
	OnBookingPaid(ctx context.Context, params *logistics_grpc_adapter.OnBookingPaidParams) (*logistics_grpc_adapter.OnBookingPaidResult, error)
}

// FinanceClient is the slice of finance-svc the booking-svc service layer
// calls (S3-E-03). Defined as an interface so tests can inject a mock.
type FinanceClient interface {
	OnPaymentReceived(ctx context.Context, params *finance_grpc_adapter.OnPaymentReceivedParams) (*finance_grpc_adapter.OnPaymentReceivedResult, error)
}

// CrmClient is the slice of crm-svc the booking-svc service layer calls (S4-E-02).
// Defined as an interface so tests can inject a mock.
// CRM calls are best-effort: failure does not block booking operations.
type CrmClient interface {
	OnBookingCreated(ctx context.Context, params *crm_grpc_adapter.OnBookingCreatedParams) (*crm_grpc_adapter.OnBookingCreatedResult, error)
	OnBookingPaidInFull(ctx context.Context, params *crm_grpc_adapter.OnBookingPaidInFullParams) (*crm_grpc_adapter.OnBookingPaidInFullResult, error)
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

	// logisticsClient is the consumer-side wrapper around logistics-svc's gRPC surface.
	// Used by MarkBookingPaid to trigger fulfillment task creation (S3-E-02).
	logisticsClient LogisticsClient

	// financeClient is the consumer-side wrapper around finance-svc's gRPC surface.
	// Used by MarkBookingPaid to post journal entries (S3-E-03).
	financeClient FinanceClient

	// crmClient is the consumer-side wrapper around crm-svc's gRPC surface (S4-E-02).
	// Used by CreateDraftBooking and MarkBookingPaid for lead lifecycle updates.
	// Best-effort: nil is acceptable; calls are skipped when nil.
	crmClient CrmClient
}

func NewService(
	logger *zerolog.Logger,
	tracer trace.Tracer,
	appName string,
	store postgres_store.IStore,
	iamClient IamClient,
	catalogClient CatalogClient,
	logisticsClient LogisticsClient,
	financeClient FinanceClient,
	crmClient CrmClient,
) IService {
	return &Service{
		logger:          logger,
		tracer:          tracer,
		appName:         appName,
		store:           store,
		iamClient:       iamClient,
		catalogClient:   catalogClient,
		logisticsClient: logisticsClient,
		financeClient:   financeClient,
		crmClient:       crmClient,
	}
}
