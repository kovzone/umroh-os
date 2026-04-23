// Package crm_grpc_adapter is booking-svc's consumer-side wrapper around
// crm-svc's internal gRPC surface (S4-E-02 / ADR-0006).
//
// booking-svc calls two RPCs:
//   - OnBookingCreated  → sent after a draft booking is created (booking_created_fanout)
//   - OnBookingPaidInFull → sent inside paid_fanout after invoice_status == "paid"
//
// Both calls are best-effort: crm-svc is a downstream enrichment; failures are
// logged but do NOT return an error to the caller (non-blocking per ADR-0006 §S4).

package crm_grpc_adapter

import (
	"context"
	"errors"
	"fmt"

	"booking-svc/adapter/crm_grpc_adapter/pb"
	"booking-svc/util/apperrors"
	"booking-svc/util/logging"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Adapter wraps crm-svc's CrmServiceClient.
type Adapter struct {
	logger    *zerolog.Logger
	tracer    trace.Tracer
	crmClient pb.CrmServiceClient
}

// NewAdapter creates a crm adapter from an already-dialled conn.
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:    logger,
		tracer:    tracer,
		crmClient: pb.NewCrmServiceClient(cc),
	}
}

// ---------------------------------------------------------------------------
// OnBookingCreated
// ---------------------------------------------------------------------------

type OnBookingCreatedParams struct {
	BookingID   string
	LeadID      string // optional; empty = no lead attribution
	PackageID   string
	DepartureID string
	JamaahCount int32
	CreatedAt   string // RFC3339
}

type OnBookingCreatedResult struct {
	Updated bool
	LeadID  string
}

func (a *Adapter) OnBookingCreated(ctx context.Context, params *OnBookingCreatedParams) (*OnBookingCreatedResult, error) {
	const op = "crm_grpc_adapter.Adapter.OnBookingCreated"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("booking_id", params.BookingID),
		attribute.String("lead_id", params.LeadID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.crmClient.OnBookingCreated(ctx, &pb.OnBookingCreatedRequest{
		BookingId:   params.BookingID,
		LeadId:      params.LeadID,
		PackageId:   params.PackageID,
		DepartureId: params.DepartureID,
		JamaahCount: params.JamaahCount,
		CreatedAt:   params.CreatedAt,
	})
	if err != nil {
		wrapped := mapCrmError(err)
		logger.Warn().
			Err(wrapped).
			Str("booking_id", params.BookingID).
			Msg("crm-svc.OnBookingCreated failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &OnBookingCreatedResult{
		Updated: resp.GetUpdated(),
		LeadID:  resp.GetLeadId(),
	}, nil
}

// ---------------------------------------------------------------------------
// OnBookingPaidInFull
// ---------------------------------------------------------------------------

type OnBookingPaidInFullParams struct {
	BookingID string
	LeadID    string // optional; empty = no lead attribution
	PaidAt    string // RFC3339
}

type OnBookingPaidInFullResult struct {
	Updated bool
	LeadID  string
}

func (a *Adapter) OnBookingPaidInFull(ctx context.Context, params *OnBookingPaidInFullParams) (*OnBookingPaidInFullResult, error) {
	const op = "crm_grpc_adapter.Adapter.OnBookingPaidInFull"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("booking_id", params.BookingID),
		attribute.String("lead_id", params.LeadID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.crmClient.OnBookingPaidInFull(ctx, &pb.OnBookingPaidInFullRequest{
		BookingId: params.BookingID,
		LeadId:    params.LeadID,
		PaidAt:    params.PaidAt,
	})
	if err != nil {
		wrapped := mapCrmError(err)
		logger.Warn().
			Err(wrapped).
			Str("booking_id", params.BookingID).
			Msg("crm-svc.OnBookingPaidInFull failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &OnBookingPaidInFullResult{
		Updated: resp.GetUpdated(),
		LeadID:  resp.GetLeadId(),
	}, nil
}

// ---------------------------------------------------------------------------
// Error mapping
// ---------------------------------------------------------------------------

func mapCrmError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("crm call failed: %w", err))
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return errors.Join(apperrors.ErrNotFound, errors.New(st.Message()))
	case grpcCodes.InvalidArgument:
		return errors.Join(apperrors.ErrValidation, errors.New(st.Message()))
	case grpcCodes.AlreadyExists:
		return errors.Join(apperrors.ErrConflict, errors.New(st.Message()))
	case grpcCodes.Unauthenticated:
		return errors.Join(apperrors.ErrUnauthorized, errors.New(st.Message()))
	case grpcCodes.PermissionDenied:
		return errors.Join(apperrors.ErrForbidden, errors.New(st.Message()))
	default:
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("crm call failed: %s", st.Message()))
	}
}
