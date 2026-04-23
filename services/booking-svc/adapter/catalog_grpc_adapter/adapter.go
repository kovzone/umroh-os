// Package catalog_grpc_adapter is booking-svc's consumer-side wrapper around
// catalog-svc's internal gRPC surface. It exposes only the two Inventory RPCs
// (ReserveSeats + ReleaseSeats) and a narrow departure validation call
// (GetPackageDeparture) per the consumer-stub rule in slice-S1.md § IAM.
//
// S1-E-03 — first consumer of this adapter: POST /v1/bookings draft creation.
package catalog_grpc_adapter

import (
	"context"
	"errors"
	"fmt"

	"booking-svc/adapter/catalog_grpc_adapter/pb"
	"booking-svc/util/apperrors"
	"booking-svc/util/logging"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc"
)

// Adapter wraps catalog-svc's CatalogServiceClient.
type Adapter struct {
	logger        *zerolog.Logger
	tracer        trace.Tracer
	catalogClient pb.CatalogServiceClient
}

// NewAdapter creates a catalog adapter from an already-dialled conn.
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:        logger,
		tracer:        tracer,
		catalogClient: pb.NewCatalogServiceClient(cc),
	}
}

// ---------------------------------------------------------------------------
// GetDeparture — validate departure for booking draft creation.
// ---------------------------------------------------------------------------

// GetDepartureResult is the subset of departure fields booking-svc needs.
type GetDepartureResult struct {
	ID             string
	PackageID      string
	DepartureDate  string
	ReturnDate     string
	TotalSeats     int
	RemainingSeats int
	Status         string
}

// GetDeparture calls catalog-svc.GetPackageDeparture and returns the departure
// detail. Returns apperrors.ErrNotFound when the departure is not open/closed.
func (a *Adapter) GetDeparture(ctx context.Context, departureID string) (*GetDepartureResult, error) {
	const op = "catalog_grpc_adapter.Adapter.GetDeparture"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("departure_id", departureID))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.catalogClient.GetPackageDeparture(ctx, &pb.GetPackageDepartureRequest{Id: departureID})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Str("departure_id", departureID).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	d := resp.GetDeparture()
	if d == nil {
		return nil, errors.Join(apperrors.ErrNotFound, fmt.Errorf("departure %q not found", departureID))
	}

	span.SetStatus(codes.Ok, "success")
	return &GetDepartureResult{
		ID:             d.GetId(),
		PackageID:      d.GetPackageId(),
		DepartureDate:  d.GetDepartureDate(),
		ReturnDate:     d.GetReturnDate(),
		TotalSeats:     int(d.GetTotalSeats()),
		RemainingSeats: int(d.GetRemainingSeats()),
		Status:         d.GetStatus(),
	}, nil
}

// ---------------------------------------------------------------------------
// ReserveSeats
// ---------------------------------------------------------------------------

// ReserveSeatsParams maps to CatalogService.ReserveSeats request.
type ReserveSeatsParams struct {
	ReservationID       string
	DepartureID         string
	Seats               int
	IdempotencyTTLHours int
}

// ReserveSeatsResult maps to CatalogService.ReserveSeats response.
type ReserveSeatsResult struct {
	ReservationID  string
	DepartureID    string
	Seats          int
	ReservedAt     string
	ExpiresAt      string
	RemainingSeats int
	Replayed       bool
}

// ReserveSeats calls catalog-svc.ReserveSeats.
func (a *Adapter) ReserveSeats(ctx context.Context, params *ReserveSeatsParams) (*ReserveSeatsResult, error) {
	const op = "catalog_grpc_adapter.Adapter.ReserveSeats"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("reservation_id", params.ReservationID),
		attribute.String("departure_id", params.DepartureID),
		attribute.Int("seats", params.Seats),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.catalogClient.ReserveSeats(ctx, &pb.ReserveSeatsRequest{
		ReservationId:       params.ReservationID,
		DepartureId:         params.DepartureID,
		Seats:               int32(params.Seats),
		IdempotencyTtlHours: int32(params.IdempotencyTTLHours),
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	res := resp.GetReservation()
	span.SetStatus(codes.Ok, "success")
	return &ReserveSeatsResult{
		ReservationID:  res.GetReservationId(),
		DepartureID:    res.GetDepartureId(),
		Seats:          int(res.GetSeats()),
		ReservedAt:     res.GetReservedAt(),
		ExpiresAt:      res.GetExpiresAt(),
		RemainingSeats: int(resp.GetRemainingSeats()),
		Replayed:       resp.GetReplayed(),
	}, nil
}

// ---------------------------------------------------------------------------
// ReleaseSeats
// ---------------------------------------------------------------------------

// ReleaseSeatsParams maps to CatalogService.ReleaseSeats request.
type ReleaseSeatsParams struct {
	ReservationID string
	Seats         int
	Reason        string
}

// ReleaseSeatsResult maps to CatalogService.ReleaseSeats response.
type ReleaseSeatsResult struct {
	ReservationID  string
	DepartureID    string
	SeatsReleased  int
	ReleasedAt     string
	RemainingSeats int
	Replayed       bool
}

// ReleaseSeats calls catalog-svc.ReleaseSeats.
func (a *Adapter) ReleaseSeats(ctx context.Context, params *ReleaseSeatsParams) (*ReleaseSeatsResult, error) {
	const op = "catalog_grpc_adapter.Adapter.ReleaseSeats"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("reservation_id", params.ReservationID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.catalogClient.ReleaseSeats(ctx, &pb.ReleaseSeatsRequest{
		ReservationId: params.ReservationID,
		Seats:         int32(params.Seats),
		Reason:        params.Reason,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	rel := resp.GetReleased()
	span.SetStatus(codes.Ok, "success")
	return &ReleaseSeatsResult{
		ReservationID:  rel.GetReservationId(),
		DepartureID:    rel.GetDepartureId(),
		SeatsReleased:  int(rel.GetSeatsReleased()),
		ReleasedAt:     rel.GetReleasedAt(),
		RemainingSeats: int(resp.GetRemainingSeats()),
		Replayed:       resp.GetReplayed(),
	}, nil
}

// ---------------------------------------------------------------------------
// Error mapping
// ---------------------------------------------------------------------------

func mapCatalogError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("catalog call failed: %w", err))
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return errors.Join(apperrors.ErrNotFound, errors.New(st.Message()))
	case grpcCodes.FailedPrecondition, grpcCodes.AlreadyExists:
		return errors.Join(apperrors.ErrConflict, errors.New(st.Message()))
	case grpcCodes.InvalidArgument:
		return errors.Join(apperrors.ErrValidation, errors.New(st.Message()))
	case grpcCodes.Unauthenticated:
		return errors.Join(apperrors.ErrUnauthorized, errors.New(st.Message()))
	case grpcCodes.PermissionDenied:
		return errors.Join(apperrors.ErrForbidden, errors.New(st.Message()))
	default:
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("catalog call failed: %s", st.Message()))
	}
}
