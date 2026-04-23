// shipment.go — ShipFulfillmentTask service-layer implementation.
//
// BL-LOG-002 / S3-E-02.
//
// ShipFulfillmentTask:
//   1. Looks up fulfillment task by booking_id.
//   2. Idempotent: if a shipment already exists for this task, returns it.
//   3. Generates a tracking number: "UOS-" + 8 upper-hex chars from a UUID.
//   4. Inserts a shipment row and marks the fulfillment task status='shipped'.
//   5. Logs a notification stub (no real WA message sent).

package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"logistics-svc/store/postgres_store/sqlc"
	"logistics-svc/util/logging"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// ShipFulfillmentTask
// ---------------------------------------------------------------------------

// ShipFulfillmentTaskParams holds inputs for shipping a fulfillment task.
type ShipFulfillmentTaskParams struct {
	BookingID string
	Carrier   string // optional; defaults to "manual"
	Notes     string
}

// ShipFulfillmentTaskResult holds the result of a ShipFulfillmentTask call.
type ShipFulfillmentTaskResult struct {
	ShipmentID     string
	TrackingNumber string
	Status         string
	Replayed       bool // true if an existing shipment was returned
}

// ShipFulfillmentTask creates a shipment for a fulfillment task.
// Idempotent: if a shipment already exists for the task, it is returned unchanged.
func (svc *Service) ShipFulfillmentTask(ctx context.Context, params *ShipFulfillmentTaskParams) (*ShipFulfillmentTaskResult, error) {
	const op = "service.Service.ShipFulfillmentTask"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("booking_id", params.BookingID),
	)

	logger := logging.LogWithTrace(ctx, svc.logger)

	if params.BookingID == "" {
		return nil, fmt.Errorf("%s: booking_id is required", op)
	}

	carrier := params.Carrier
	if carrier == "" {
		carrier = "manual"
	}

	// --- Lookup task ---
	task, err := svc.store.GetFulfillmentTaskByBookingID(ctx, params.BookingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: no fulfillment task found for booking_id %s", op, params.BookingID)
		}
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: get task: %w", op, err)
	}

	// --- Idempotency: check if shipment already exists ---
	existing, err := svc.store.GetShipmentByTaskID(ctx, task.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: get existing shipment: %w", op, err)
	}
	if err == nil {
		logger.Info().
			Str("op", op).
			Str("booking_id", params.BookingID).
			Str("shipment_id", existing.ID).
			Str("tracking_number", existing.TrackingNumber).
			Bool("replayed", true).
			Msg("shipment already exists, returning existing")

		span.SetStatus(otelCodes.Ok, "replayed")
		return &ShipFulfillmentTaskResult{
			ShipmentID:     existing.ID,
			TrackingNumber: existing.TrackingNumber,
			Status:         existing.Status,
			Replayed:       true,
		}, nil
	}

	// --- Generate tracking number: "UOS-" + 8 upper-hex chars of a UUID ---
	rawUUID := uuid.New().String()
	hexPart := strings.ToUpper(strings.ReplaceAll(rawUUID, "-", ""))[:8]
	trackingNumber := "UOS-" + hexPart

	// --- Insert shipment ---
	shipment, err := svc.store.InsertShipment(ctx, sqlc.InsertShipmentParams{
		TaskID:         task.ID,
		TrackingNumber: trackingNumber,
		Carrier:        carrier,
		Notes:          params.Notes,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: insert shipment: %w", op, err)
	}

	// --- Update fulfillment task status to 'shipped' ---
	if err := svc.store.UpdateFulfillmentTaskStatus(ctx, sqlc.UpdateFulfillmentTaskStatusParams{
		ID:     task.ID,
		Status: "shipped",
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: update task status: %w", op, err)
	}

	// --- Notification stub (WA not wired yet) ---
	logger.Info().
		Str("op", op).
		Str("booking_id", params.BookingID).
		Str("task_id", task.ID).
		Str("shipment_id", shipment.ID).
		Str("tracking_number", trackingNumber).
		Str("carrier", carrier).
		Msgf("WA notify: booking %s shipped, tracking: %s", params.BookingID, trackingNumber)

	span.SetStatus(otelCodes.Ok, "shipped")
	return &ShipFulfillmentTaskResult{
		ShipmentID:     shipment.ID,
		TrackingNumber: trackingNumber,
		Status:         shipment.Status,
		Replayed:       false,
	}, nil
}
