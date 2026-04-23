// scan_boarding.go — gRPC handlers for scan events and bus boarding RPCs
// (BL-OPS-010/011).
//
// RecordScan:        idempotent scan event recording.
// RecordBusBoarding: atomic bus boarding + scan insert.
// GetBoardingRoster: aggregate roster query with optional bus filter.

package grpc_api

import (
	"context"
	"errors"
	"time"

	"ops-svc/api/grpc_api/pb"
	"ops-svc/service"
	"ops-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// RecordScan (BL-OPS-010)
// ---------------------------------------------------------------------------

func (s *Server) RecordScan(ctx context.Context, req *pb.RecordScanRequest) (*pb.RecordScanResponse, error) {
	const op = "grpc_api.Server.RecordScan"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("scan_type", req.GetScanType()),
		attribute.String("departure_id", req.GetDepartureId()),
	)
	logger.Info().Str("op", op).Str("scan_type", req.GetScanType()).Msg("")

	result, err := s.svc.RecordScan(ctx, &service.RecordScanParams{
		ScanType:       req.GetScanType(),
		DepartureID:    req.GetDepartureId(),
		JamaahID:       req.GetJamaahId(),
		ScannedBy:      req.GetScannedBy(),
		DeviceID:       req.GetDeviceId(),
		Location:       req.GetLocation(),
		IdempotencyKey: req.GetIdempotencyKey(),
		Metadata:       req.GetMetadata(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		switch {
		case errors.Is(err, service.ErrInvalidScanType):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "invalid_scan_type")
		case errors.Is(err, service.ErrMissingField):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "missing_required_field")
		default:
			return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
		}
	}

	return &pb.RecordScanResponse{
		ScanId:     result.ScanID,
		Idempotent: result.Idempotent,
	}, nil
}

// ---------------------------------------------------------------------------
// RecordBusBoarding (BL-OPS-011)
// ---------------------------------------------------------------------------

func (s *Server) RecordBusBoarding(ctx context.Context, req *pb.RecordBusBoardingRequest) (*pb.RecordBusBoardingResponse, error) {
	const op = "grpc_api.Server.RecordBusBoarding"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("departure_id", req.GetDepartureId()),
		attribute.String("jamaah_id", req.GetJamaahId()),
	)
	logger.Info().Str("op", op).Str("departure_id", req.GetDepartureId()).Msg("")

	result, err := s.svc.RecordBusBoarding(ctx, &service.RecordBusBoardingParams{
		DepartureID: req.GetDepartureId(),
		BusNumber:   req.GetBusNumber(),
		JamaahID:    req.GetJamaahId(),
		ScannedBy:   req.GetScannedBy(),
		Status:      req.GetStatus(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())

		switch {
		case errors.Is(err, service.ErrInvalidBoardingStatus):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "invalid_status")
		case errors.Is(err, service.ErrMissingField):
			return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "missing_required_field")
		default:
			return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
		}
	}

	return &pb.RecordBusBoardingResponse{
		BoardingId: result.BoardingID,
		Status:     result.Status,
		Idempotent: result.Idempotent,
	}, nil
}

// ---------------------------------------------------------------------------
// GetBoardingRoster (BL-OPS-011)
// ---------------------------------------------------------------------------

func (s *Server) GetBoardingRoster(ctx context.Context, req *pb.GetBoardingRosterRequest) (*pb.GetBoardingRosterResponse, error) {
	const op = "grpc_api.Server.GetBoardingRoster"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureId()))
	logger.Info().Str("op", op).Str("departure_id", req.GetDepartureId()).Msg("")

	if req.GetDepartureId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "departure_id is required")
	}

	result, err := s.svc.GetBoardingRoster(ctx, &service.GetBoardingRosterParams{
		DepartureID: req.GetDepartureId(),
		BusNumber:   req.GetBusNumber(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	var boardings []*pb.BoardingEntry
	for _, b := range result.Boardings {
		ts := b.BoardedAt
		if ts == "" {
			ts = ""
		}
		_ = time.RFC3339 // import used
		boardings = append(boardings, &pb.BoardingEntry{
			JamaahId:  b.JamaahID,
			BusNumber: b.BusNumber,
			Status:    b.Status,
			BoardedAt: ts,
		})
	}

	return &pb.GetBoardingRosterResponse{
		Boardings:    boardings,
		TotalBoarded: result.TotalBoarded,
		TotalAbsent:  result.TotalAbsent,
	}, nil
}
