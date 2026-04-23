package service

import (
	"context"

	"logistics-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for logistics-svc.
//
// Pilot scaffold surfaces the three standard scaffold endpoints plus the
// S3-E-02 fulfillment trigger:
//
//   - Liveness — process is up
//   - Readiness — process is up AND the database is reachable
//   - DbTxDiagnostic — writes + reads inside a WithTx, the canonical reference
//     for how services should use transactions (per docs/04-backend-conventions)
//   - OnBookingPaid — creates (or returns existing) fulfillment task for a
//     paid-in-full booking (S3-E-02 / BL-LOG-001)
//   - ShipFulfillmentTask — creates a shipment + tracking number (BL-LOG-002)
//   - GeneratePickupQR — generates a single-use pickup token with 7d TTL (BL-LOG-003)
//   - RedeemPickupQR — validates and marks a pickup token as used (BL-LOG-003)
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)
	DbTxDiagnostic(ctx context.Context, params *DbTxDiagnosticParams) (*DbTxDiagnosticResult, error)

	// OnBookingPaid creates a fulfillment task for a booking that has just
	// been fully paid. Idempotent: returns the existing task if one already
	// exists for this booking_id.
	OnBookingPaid(ctx context.Context, params *OnBookingPaidParams) (*OnBookingPaidResult, error)

	// ShipFulfillmentTask creates a shipment record + tracking number (BL-LOG-002).
	ShipFulfillmentTask(ctx context.Context, params *ShipFulfillmentTaskParams) (*ShipFulfillmentTaskResult, error)

	// GeneratePickupQR creates a single-use pickup QR token with 7d TTL (BL-LOG-003).
	GeneratePickupQR(ctx context.Context, params *GeneratePickupQRParams) (*GeneratePickupQRResult, error)

	// RedeemPickupQR validates and marks a pickup token as used (BL-LOG-003).
	RedeemPickupQR(ctx context.Context, params *RedeemPickupQRParams) (*RedeemPickupQRResult, error)
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
