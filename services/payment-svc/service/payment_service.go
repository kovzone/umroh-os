// payment_service.go — PaymentService struct, interface extension, and constructor.
//
// The existing Service struct handles scaffold RPCs (liveness, readiness,
// db_tx_diagnostic). PaymentService extends the service layer with the F5
// payment domain (invoice, VA, webhook, refund, reconciliation).
//
// Both implement IService via composition: Service embeds the system methods;
// PaymentService wraps Service and adds payment methods. cmd/start.go constructs
// a PaymentService so all RPCs route through it.

package service

import (
	"context"

	"payment-svc/adapter/booking_grpc_adapter"
	"payment-svc/adapter/gateway"
	"payment-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// iamAuditAdapter is a minimal interface that payment-svc needs from the
// iam-svc adapter. Defined here to avoid a circular import; the concrete
// implementation lives in adapter/iam_grpc_adapter (to be created in S3).
// For MVP, inject nil and audit calls become no-ops (best-effort pattern).
type iamAuditAdapter struct {
	ActorUserID string
	Resource    string
	ResourceID  string
	Action      string
	OldValue    []byte
	NewValue    []byte
}

// IAMAuditEmitter is the narrow interface payment-svc needs for audit logging.
type IAMAuditEmitter interface {
	RecordAudit(ctx context.Context, params *iamAuditAdapter) (interface{}, error)
}

// IPaymentService extends IService with the payment-domain RPCs.
type IPaymentService interface {
	IService

	// IssueVirtualAccount creates an invoice and VA for a booking (F5-W1).
	IssueVirtualAccount(ctx context.Context, params *IssueVAParams) (*IssueVAResult, error)

	// ProcessWebhookEvent handles an incoming gateway webhook (F5-W2).
	ProcessWebhookEvent(ctx context.Context, params *WebhookEventParams) (*WebhookResult, error)

	// StartRefund initiates a refund flow (F5-W8).
	StartRefund(ctx context.Context, params *StartRefundParams) (*StartRefundResult, error)

	// ReconcileInvoices runs the reconciliation loop for a batch of invoices (F5-W5).
	// Called by the reconciliation cron in usecase/reconcile.
	ReconcileInvoices(ctx context.Context) (*ReconcileResult, error)

	// StartReconcileCron launches the hourly reconciliation ticker in a goroutine.
	// Call once from cmd/start.go after all dependencies are wired.
	// Blocks until ctx is cancelled — must be called as `go svc.StartReconcileCron(ctx)` or
	// via a method that itself starts the goroutine.
	StartReconcileCron(ctx context.Context)

	// ReissuePaymentLink retrieves the active VA for an existing booking's invoice,
	// or creates a new VA on the same invoice when the existing VA has expired.
	// CS-facing (BL-PAY-020).
	ReissuePaymentLink(ctx context.Context, params *ReissuePaymentLinkParams) (*ReissuePaymentLinkResult, error)

	// GetInvoiceByID fetches a single invoice by its UUID (BL-PAY-001 / ISSUE-005).
	// Used by gateway-svc for GET /v1/invoices/:id and the VA re-issuance flow.
	GetInvoiceByID(ctx context.Context, params *GetInvoiceByIDParams) (*GetInvoiceByIDResult, error)
}

// PaymentService is the concrete implementation of IPaymentService.
type PaymentService struct {
	// Embed the base Service so scaffold RPCs (Liveness, Readiness, DbTxDiagnostic)
	// are available without re-implementation.
	*Service

	// Gateway adapters.
	primaryGateway gateway.GatewayAdapter // Midtrans
	xenditGateway  gateway.GatewayAdapter // Xendit (hot-standby)
	mockGateway    gateway.GatewayAdapter // mock (MOCK_GATEWAY=true)

	// Downstream service adapters.
	bookingAdapter *booking_grpc_adapter.Adapter

	// IAM audit emitter (nil = no-op for MVP).
	iamAudit IAMAuditEmitter
}

// PaymentServiceConfig holds the dependencies for NewPaymentService.
type PaymentServiceConfig struct {
	Logger         *zerolog.Logger
	Tracer         trace.Tracer
	AppName        string
	Store          postgres_store.IStore
	PrimaryGateway gateway.GatewayAdapter
	XenditGateway  gateway.GatewayAdapter
	MockGateway    gateway.GatewayAdapter
	BookingAdapter *booking_grpc_adapter.Adapter
	IAMAudit       IAMAuditEmitter
}

// NewPaymentService constructs a PaymentService with all dependencies injected.
func NewPaymentService(cfg PaymentServiceConfig) IPaymentService {
	base := &Service{
		logger:  cfg.Logger,
		tracer:  cfg.Tracer,
		appName: cfg.AppName,
		store:   cfg.Store,
	}
	return &PaymentService{
		Service:        base,
		primaryGateway: cfg.PrimaryGateway,
		xenditGateway:  cfg.XenditGateway,
		mockGateway:    cfg.MockGateway,
		bookingAdapter: cfg.BookingAdapter,
		iamAudit:       cfg.IAMAudit,
	}
}
