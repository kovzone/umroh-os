// paid_fanout.go — FanOutBookingPaid service-layer implementation.
//
// S3-E-02 / S3-E-03 / S4-E-02 / ADR-0006.
//
// Called from the MarkBookingPaid gRPC handler after a booking transitions to
// paid_in_full. Fans out to:
//   - logistics-svc.OnBookingPaid      → creates fulfillment_task (S3-E-02)
//   - finance-svc.OnPaymentReceived    → posts double-entry journal (S3-E-03)
//   - crm-svc.OnBookingPaidInFull      → updates lead to 'converted' (S4-E-02)
//
// IMPORTANT (ADR-0006 / §S3-J-01): logistics + finance calls are SYNCHRONOUS.
// The CRM call is BEST-EFFORT: failure is logged but does not propagate.
// Both logistics and finance are idempotent, so retries are safe.

package service

import (
	"context"
	"fmt"

	"booking-svc/adapter/crm_grpc_adapter"
	"booking-svc/adapter/finance_grpc_adapter"
	"booking-svc/adapter/logistics_grpc_adapter"
	"booking-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// FanOutBookingPaidParams holds the inputs for the fan-out.
// Amount is int64 (integer IDR) per §S3-J-03 contract.
// JamaahIDs is required for the logistics fan-out per §S3-J-02 contract.
// LeadID is optional; used by the CRM fan-out to update lead status (S4-E-02).
type FanOutBookingPaidParams struct {
	BookingID   string
	DepartureID string
	JamaahIDs   []string // at least 1; ULIDs of jamaah on this booking
	InvoiceID   string
	Amount      int64  // integer IDR
	ReceivedAt  string // RFC3339; empty = server time
	LeadID      string // optional; CRM lead attribution (S4-E-02)
	PaidAt      string // RFC3339; empty = server time
}

// FanOutBookingPaidResult holds the combined result of both fan-out calls.
type FanOutBookingPaidResult struct {
	TaskID   string
	TaskStatus string
	EntryID  string
	Balanced bool
}

// FanOutBookingPaid triggers downstream services after a booking is confirmed
// as paid_in_full. Both calls are SYNCHRONOUS (ADR-0006 / §S3-J-01). Any
// error from either downstream is returned to the caller so payment-svc can
// retry via gateway.
func (svc *Service) FanOutBookingPaid(ctx context.Context, params *FanOutBookingPaidParams) (*FanOutBookingPaidResult, error) {
	const op = "service.Service.FanOutBookingPaid"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, svc.logger)

	logger.Info().
		Str("op", op).
		Str("booking_id", params.BookingID).
		Str("departure_id", params.DepartureID).
		Str("invoice_id", params.InvoiceID).
		Int64("amount", params.Amount).
		Msg("starting paid fan-out")

	var result FanOutBookingPaidResult

	// --- logistics-svc.OnBookingPaid (synchronous) ---
	if svc.logisticsClient != nil {
		logCtx, logSpan := svc.tracer.Start(ctx, op+".logistics")
		logSpan.SetAttributes(
			attribute.String("booking_id", params.BookingID),
			attribute.String("departure_id", params.DepartureID),
		)

		logResult, err := svc.logisticsClient.OnBookingPaid(logCtx, &logistics_grpc_adapter.OnBookingPaidParams{
			BookingID:   params.BookingID,
			DepartureID: params.DepartureID,
			JamaahIDs:   params.JamaahIDs,
		})
		if err != nil {
			logSpan.RecordError(err)
			logSpan.SetStatus(otelCodes.Error, err.Error())
			logSpan.End()
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, err.Error())
			logger.Error().
				Str("op", op+".logistics").
				Str("booking_id", params.BookingID).
				Err(err).
				Msg("logistics-svc.OnBookingPaid failed")
			return nil, fmt.Errorf("%s: logistics fan-out: %w", op, err)
		}

		logger.Info().
			Str("op", op+".logistics").
			Str("booking_id", params.BookingID).
			Str("task_id", logResult.TaskID).
			Str("task_status", logResult.Status).
			Msg("logistics-svc.OnBookingPaid succeeded")
		logSpan.SetStatus(otelCodes.Ok, "success")
		logSpan.End()
		result.TaskID = logResult.TaskID
		result.TaskStatus = logResult.Status
	} else {
		logger.Warn().
			Str("op", op).
			Msg("logisticsClient is nil — skipping logistics fan-out (service not configured)")
	}

	// --- finance-svc.OnPaymentReceived (synchronous) ---
	if svc.financeClient != nil {
		finCtx, finSpan := svc.tracer.Start(ctx, op+".finance")
		finSpan.SetAttributes(
			attribute.String("booking_id", params.BookingID),
			attribute.String("invoice_id", params.InvoiceID),
			attribute.Int64("amount", params.Amount),
		)

		finResult, err := svc.financeClient.OnPaymentReceived(finCtx, &finance_grpc_adapter.OnPaymentReceivedParams{
			BookingID:  params.BookingID,
			InvoiceID:  params.InvoiceID,
			Amount:     params.Amount,
			ReceivedAt: params.ReceivedAt,
		})
		if err != nil {
			finSpan.RecordError(err)
			finSpan.SetStatus(otelCodes.Error, err.Error())
			finSpan.End()
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, err.Error())
			logger.Error().
				Str("op", op+".finance").
				Str("booking_id", params.BookingID).
				Str("invoice_id", params.InvoiceID).
				Err(err).
				Msg("finance-svc.OnPaymentReceived failed")
			return nil, fmt.Errorf("%s: finance fan-out: %w", op, err)
		}

		logger.Info().
			Str("op", op+".finance").
			Str("booking_id", params.BookingID).
			Str("invoice_id", params.InvoiceID).
			Str("entry_id", finResult.EntryID).
			Bool("balanced", finResult.Balanced).
			Msg("finance-svc.OnPaymentReceived succeeded")
		finSpan.SetStatus(otelCodes.Ok, "success")
		finSpan.End()
		result.EntryID = finResult.EntryID
		result.Balanced = finResult.Balanced
	} else {
		logger.Warn().
			Str("op", op).
			Msg("financeClient is nil — skipping finance fan-out (service not configured)")
	}

	// --- crm-svc.OnBookingPaidInFull (best-effort, S4-E-02) ---
	//
	// CRM enrichment is NOT synchronous and NOT fatal. A failure here must
	// NOT cause the booking paid response to fail. Log + continue.
	if svc.crmClient != nil {
		crmCtx, crmSpan := svc.tracer.Start(ctx, op+".crm")
		crmSpan.SetAttributes(
			attribute.String("booking_id", params.BookingID),
			attribute.String("lead_id", params.LeadID),
		)

		crmResult, crmErr := svc.crmClient.OnBookingPaidInFull(crmCtx, &crm_grpc_adapter.OnBookingPaidInFullParams{
			BookingID: params.BookingID,
			LeadID:    params.LeadID,
			PaidAt:    params.PaidAt,
		})
		if crmErr != nil {
			// Best-effort: log but do NOT return error.
			crmSpan.RecordError(crmErr)
			crmSpan.SetStatus(otelCodes.Error, crmErr.Error())
			crmSpan.End()
			logger.Warn().
				Err(crmErr).
				Str("op", op+".crm").
				Str("booking_id", params.BookingID).
				Msg("crm-svc.OnBookingPaidInFull failed (best-effort, ignored)")
		} else {
			logger.Info().
				Str("op", op+".crm").
				Str("booking_id", params.BookingID).
				Bool("updated", crmResult.Updated).
				Str("lead_id", crmResult.LeadID).
				Msg("crm-svc.OnBookingPaidInFull succeeded")
			crmSpan.SetStatus(otelCodes.Ok, "success")
			crmSpan.End()
		}
	} else {
		logger.Warn().
			Str("op", op).
			Msg("crmClient is nil — skipping CRM paid fan-out (service not configured)")
	}

	span.SetStatus(otelCodes.Ok, "fan-out complete")
	return &result, nil
}
