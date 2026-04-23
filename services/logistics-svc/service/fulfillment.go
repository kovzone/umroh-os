// fulfillment.go — OnBookingPaid service-layer implementation for logistics-svc.
//
// S3-E-02 / BL-LOG-001.
//
// Called by booking-svc (via gRPC) after a booking transitions to
// paid_in_full.  Creates a fulfillment_task row in status='queued'.
// Idempotent: if a task already exists for the booking, returns it unchanged.
//
// Audit: every new task creation is logged via zerolog structured entry.
// Full IAM audit-log write (iam-svc.RecordAudit) is scaffolded here;
// the actual gRPC call will be wired once iam-svc adapter is added to
// logistics-svc in a subsequent task.

package service

import (
	"context"
	"errors"
	"fmt"

	"logistics-svc/store/postgres_store/sqlc"
	"logistics-svc/util/logging"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// OnBookingPaidParams holds the inputs for creating a fulfillment task.
// JamaahIDs is required (min 1) per §S3-J-02 contract.
type OnBookingPaidParams struct {
	BookingID   string
	DepartureID string
	JamaahIDs   []string // at least 1; used to reserve kit per jamaah
}

// OnBookingPaidResult holds the result of an OnBookingPaid call.
type OnBookingPaidResult struct {
	TaskID  string
	Status  string
	// Replayed is true when an existing task was returned without a new insert.
	Replayed bool
}

// OnBookingPaid creates a fulfillment task for a paid-in-full booking.
// If a task already exists for this booking_id, it returns the existing task
// without creating a duplicate (idempotent).
func (svc *Service) OnBookingPaid(ctx context.Context, params *OnBookingPaidParams) (*OnBookingPaidResult, error) {
	const op = "service.Service.OnBookingPaid"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("booking_id", params.BookingID),
		attribute.String("departure_id", params.DepartureID),
	)

	logger := logging.LogWithTrace(ctx, svc.logger)

	if params.BookingID == "" {
		return nil, fmt.Errorf("%s: booking_id is required", op)
	}
	if params.DepartureID == "" {
		return nil, fmt.Errorf("%s: departure_id is required", op)
	}
	if len(params.JamaahIDs) == 0 {
		return nil, fmt.Errorf("%s: jamaah_ids is required (min 1)", op)
	}

	// --- Idempotency check ---
	existing, err := svc.store.GetFulfillmentTaskByBookingID(ctx, params.BookingID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: get existing task: %w", op, err)
	}
	if err == nil {
		// Task already exists — return it without a new insert.
		logger.Info().
			Str("op", op).
			Str("booking_id", params.BookingID).
			Str("task_id", existing.ID).
			Str("status", existing.Status).
			Bool("replayed", true).
			Msg("fulfillment task already exists, returning existing")

		span.SetStatus(otelCodes.Ok, "replayed")
		return &OnBookingPaidResult{
			TaskID:   existing.ID,
			Status:   existing.Status,
			Replayed: true,
		}, nil
	}

	// --- Insert new task ---
	task, err := svc.store.InsertFulfillmentTask(ctx, sqlc.InsertFulfillmentTaskParams{
		BookingID:   params.BookingID,
		DepartureID: params.DepartureID,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: insert fulfillment task: %w", op, err)
	}

	logger.Info().
		Str("op", op).
		Str("booking_id", params.BookingID).
		Str("departure_id", params.DepartureID).
		Str("task_id", task.ID).
		Str("status", task.Status).
		Bool("replayed", false).
		Msg("fulfillment task created")

	span.SetStatus(otelCodes.Ok, "created")
	return &OnBookingPaidResult{
		TaskID:   task.ID,
		Status:   task.Status,
		Replayed: false,
	}, nil
}
