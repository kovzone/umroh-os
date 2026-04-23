// booking_created_fanout.go — fan-out to crm-svc after a draft booking is created.
//
// S4-E-02 / ADR-0006.
//
// Called from CreateDraftBooking after the booking row is persisted.
// Fans out to:
//   - crm-svc.OnBookingCreated → updates lead status to 'qualified' if lead_id provided
//
// Per ADR-0006: the CRM call is BEST-EFFORT. Unlike the paid fan-out (logistics +
// finance), a CRM enrichment failure must NOT block or fail the booking creation
// response. Failures are logged and the span records the error, but the function
// always returns nil to its caller.

package service

import (
	"context"

	"booking-svc/adapter/crm_grpc_adapter"
	"booking-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// FanOutBookingCreatedParams holds the inputs for the CRM booking-created fan-out.
type FanOutBookingCreatedParams struct {
	BookingID   string
	LeadID      string // optional; empty = no lead attribution
	PackageID   string
	DepartureID string
	JamaahCount int32
	CreatedAt   string // RFC3339
}

// FanOutBookingCreated calls crm-svc.OnBookingCreated in a best-effort manner.
// Returns nil in all cases — CRM enrichment failure must not block booking creation.
func (svc *Service) FanOutBookingCreated(ctx context.Context, params *FanOutBookingCreatedParams) {
	const op = "service.Service.FanOutBookingCreated"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, svc.logger)

	logger.Info().
		Str("op", op).
		Str("booking_id", params.BookingID).
		Str("lead_id", params.LeadID).
		Msg("starting booking-created CRM fan-out")

	if svc.crmClient == nil {
		logger.Warn().
			Str("op", op).
			Msg("crmClient is nil — skipping CRM booking-created fan-out (service not configured)")
		span.SetStatus(otelCodes.Ok, "skipped — crmClient nil")
		return
	}

	crmCtx, crmSpan := svc.tracer.Start(ctx, op+".crm")
	crmSpan.SetAttributes(
		attribute.String("booking_id", params.BookingID),
		attribute.String("lead_id", params.LeadID),
	)

	result, err := svc.crmClient.OnBookingCreated(crmCtx, &crm_grpc_adapter.OnBookingCreatedParams{
		BookingID:   params.BookingID,
		LeadID:      params.LeadID,
		PackageID:   params.PackageID,
		DepartureID: params.DepartureID,
		JamaahCount: params.JamaahCount,
		CreatedAt:   params.CreatedAt,
	})
	if err != nil {
		// Best-effort: log the error but do NOT propagate.
		crmSpan.RecordError(err)
		crmSpan.SetStatus(otelCodes.Error, err.Error())
		crmSpan.End()
		span.RecordError(err)
		logger.Warn().
			Err(err).
			Str("op", op+".crm").
			Str("booking_id", params.BookingID).
			Msg("crm-svc.OnBookingCreated failed (best-effort, ignored)")
		span.SetStatus(otelCodes.Ok, "crm fan-out failed but ignored")
		return
	}

	logger.Info().
		Str("op", op+".crm").
		Str("booking_id", params.BookingID).
		Bool("updated", result.Updated).
		Str("lead_id", result.LeadID).
		Msg("crm-svc.OnBookingCreated succeeded")

	crmSpan.SetStatus(otelCodes.Ok, "success")
	crmSpan.End()
	span.SetStatus(otelCodes.Ok, "fan-out complete")
}
