// Package logistics_grpc_adapter is booking-svc's consumer-side wrapper around
// logistics-svc's internal gRPC surface. It exposes only OnBookingPaid, the
// single RPC needed by the paid-booking fan-out path (S3-E-02 / ADR-0006).

package logistics_grpc_adapter

import (
	"context"
	"errors"
	"fmt"

	"booking-svc/adapter/logistics_grpc_adapter/pb"
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

// Adapter wraps logistics-svc's LogisticsServiceClient.
type Adapter struct {
	logger           *zerolog.Logger
	tracer           trace.Tracer
	logisticsClient  pb.LogisticsServiceClient
}

// NewAdapter creates a logistics adapter from an already-dialled conn.
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:          logger,
		tracer:          tracer,
		logisticsClient: pb.NewLogisticsServiceClient(cc),
	}
}

// ---------------------------------------------------------------------------
// OnBookingPaid — trigger fulfillment task creation.
// ---------------------------------------------------------------------------

// OnBookingPaidParams maps to LogisticsService.OnBookingPaid request.
// JamaahIDs is required (min 1) per §S3-J-02 contract.
type OnBookingPaidParams struct {
	BookingID   string
	DepartureID string
	JamaahIDs   []string // at least 1; ULIDs of jamaah on this booking
}

// OnBookingPaidResult maps to LogisticsService.OnBookingPaid response.
type OnBookingPaidResult struct {
	TaskID string
	Status string
}

// OnBookingPaid calls logistics-svc.OnBookingPaid and returns the created (or
// existing) fulfillment task.
func (a *Adapter) OnBookingPaid(ctx context.Context, params *OnBookingPaidParams) (*OnBookingPaidResult, error) {
	const op = "logistics_grpc_adapter.Adapter.OnBookingPaid"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("booking_id", params.BookingID),
		attribute.String("departure_id", params.DepartureID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.logisticsClient.OnBookingPaid(ctx, &pb.OnBookingPaidRequest{
		BookingId:   params.BookingID,
		DepartureId: params.DepartureID,
		JamaahIds:   params.JamaahIDs,
	})
	if err != nil {
		wrapped := mapLogisticsError(err)
		logger.Warn().
			Err(wrapped).
			Str("booking_id", params.BookingID).
			Msg("logistics-svc.OnBookingPaid failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &OnBookingPaidResult{
		TaskID: resp.GetTaskId(),
		Status: resp.GetStatus(),
	}, nil
}

// ---------------------------------------------------------------------------
// Error mapping
// ---------------------------------------------------------------------------

func mapLogisticsError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("logistics call failed: %w", err))
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
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("logistics call failed: %s", st.Message()))
	}
}
