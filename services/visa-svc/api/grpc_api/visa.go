// visa.go — gRPC handlers for visa pipeline RPCs (BL-VISA-001..003).
//
// TransitionStatus: enforce state machine, record history, idempotency.
// BulkSubmit:       atomic READY→SUBMITTED batch with provider assignment.
// GetApplications:  list applications for departure with embedded history.

package grpc_api

import (
	"context"
	"errors"
	"time"

	"visa-svc/api/grpc_api/pb"
	"visa-svc/service"
	"visa-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
	grpcCodes "google.golang.org/grpc/codes"
)

// ---------------------------------------------------------------------------
// TransitionStatus (BL-VISA-001)
// ---------------------------------------------------------------------------

func (s *Server) TransitionStatus(ctx context.Context, req *pb.TransitionStatusRequest) (*pb.TransitionStatusResponse, error) {
	const op = "grpc_api.Server.TransitionStatus"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("application_id", req.GetApplicationId()),
		attribute.String("to_status", req.GetToStatus()),
	)
	logger.Info().Str("op", op).Str("application_id", req.GetApplicationId()).Msg("")

	if req.GetApplicationId() == "" {
		return nil, status.Error(grpcCodes.InvalidArgument, "application_id is required")
	}
	if req.GetToStatus() == "" {
		return nil, status.Error(grpcCodes.InvalidArgument, "to_status is required")
	}

	result, err := s.svc.TransitionStatus(ctx, &service.TransitionStatusParams{
		ApplicationID: req.GetApplicationId(),
		ToStatus:      req.GetToStatus(),
		Reason:        req.GetReason(),
		ActorUserID:   req.GetActorUserId(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		switch {
		case errors.Is(err, service.ErrVisaNotFound):
			return nil, status.Error(grpcCodes.NotFound, "not_found")
		case errors.Is(err, service.ErrInvalidStatus):
			return nil, status.Error(grpcCodes.InvalidArgument, "invalid_status_value")
		case errors.Is(err, service.ErrInvalidTransition):
			return nil, status.Error(grpcCodes.FailedPrecondition, "invalid_transition")
		default:
			return nil, status.Error(grpcCodes.Internal, "internal_error")
		}
	}

	return &pb.TransitionStatusResponse{
		ApplicationId: result.ApplicationID,
		FromStatus:    result.FromStatus,
		ToStatus:      result.ToStatus,
		Idempotent:    result.Idempotent,
	}, nil
}

// ---------------------------------------------------------------------------
// BulkSubmit (BL-VISA-002)
// ---------------------------------------------------------------------------

func (s *Server) BulkSubmit(ctx context.Context, req *pb.BulkSubmitRequest) (*pb.BulkSubmitResponse, error) {
	const op = "grpc_api.Server.BulkSubmit"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("departure_id", req.GetDepartureId()),
		attribute.Int("jamaah_count", len(req.GetJamaahIds())),
	)
	logger.Info().Str("op", op).Str("departure_id", req.GetDepartureId()).Msg("")

	if req.GetDepartureId() == "" {
		return nil, status.Error(grpcCodes.InvalidArgument, "departure_id is required")
	}

	result, err := s.svc.BulkSubmit(ctx, &service.BulkSubmitParams{
		DepartureID: req.GetDepartureId(),
		JamaahIDs:   req.GetJamaahIds(),
		ProviderID:  req.GetProviderId(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		switch {
		case errors.Is(err, service.ErrEmptyJamaahList):
			return nil, status.Error(grpcCodes.InvalidArgument, "empty_jamaah_list")
		case errors.Is(err, service.ErrInvalidProvider):
			return nil, status.Error(grpcCodes.InvalidArgument, "invalid_provider")
		case errors.Is(err, service.ErrNotAllReady):
			return nil, status.Error(grpcCodes.FailedPrecondition, "not_all_ready")
		default:
			return nil, status.Error(grpcCodes.Internal, "internal_error")
		}
	}

	return &pb.BulkSubmitResponse{
		SubmittedCount: result.SubmittedCount,
		ApplicationIds: result.ApplicationIDs,
	}, nil
}

// ---------------------------------------------------------------------------
// GetApplications (BL-VISA-003)
// ---------------------------------------------------------------------------

func (s *Server) GetApplications(ctx context.Context, req *pb.GetApplicationsRequest) (*pb.GetApplicationsResponse, error) {
	const op = "grpc_api.Server.GetApplications"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureId()))
	logger.Info().Str("op", op).Str("departure_id", req.GetDepartureId()).Msg("")

	if req.GetDepartureId() == "" {
		return nil, status.Error(grpcCodes.InvalidArgument, "missing_departure_id")
	}

	result, err := s.svc.GetApplications(ctx, &service.GetApplicationsParams{
		DepartureID:  req.GetDepartureId(),
		StatusFilter: req.GetStatusFilter(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		if err.Error() == "invalid_status_filter" {
			return nil, status.Error(grpcCodes.InvalidArgument, "invalid_status_filter")
		}
		return nil, status.Error(grpcCodes.Internal, "internal_error")
	}

	var apps []*pb.ApplicationRecord
	for _, r := range result.Applications {
		var history []*pb.StatusHistoryEntry
		for _, h := range r.History {
			ts := ""
			if !h.CreatedAt.IsZero() {
				ts = h.CreatedAt.UTC().Format(time.RFC3339)
			}
			history = append(history, &pb.StatusHistoryEntry{
				FromStatus: h.FromStatus,
				ToStatus:   h.ToStatus,
				Reason:     h.Reason,
				CreatedAt:  ts,
			})
		}
		apps = append(apps, &pb.ApplicationRecord{
			Id:          r.ID,
			JamaahId:    r.JamaahID,
			Status:      r.Status,
			ProviderRef: r.ProviderRef,
			IssuedDate:  r.IssuedDate,
			History:     history,
		})
	}

	return &pb.GetApplicationsResponse{Applications: apps}, nil
}
