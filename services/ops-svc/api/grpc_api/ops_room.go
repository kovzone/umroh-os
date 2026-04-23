// ops_room.go — gRPC handlers for RunRoomAllocation and GetRoomAllocation RPCs
// (BL-OPS-002).

package grpc_api

import (
	"context"

	"ops-svc/api/grpc_api/pb"
	"ops-svc/service"
	"ops-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// RunRoomAllocation handles the RunRoomAllocation RPC.
// Groups jamaah_ids into rooms for a departure using the configured algorithm.
func (s *Server) RunRoomAllocation(ctx context.Context, req *pb.RunRoomAllocationRequest) (*pb.RunRoomAllocationResponse, error) {
	const op = "grpc_api.Server.RunRoomAllocation"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("op", op),
		attribute.String("departure_id", req.GetDepartureId()),
		attribute.Int("jamaah_count", len(req.GetJamaahIds())),
	)

	logger.Info().
		Str("op", op).
		Str("departure_id", req.GetDepartureId()).
		Int("jamaah_count", len(req.GetJamaahIds())).
		Msg("")

	if req.GetDepartureId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "departure_id is required")
	}
	if len(req.GetJamaahIds()) == 0 {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "jamaah_ids must not be empty")
	}

	result, err := s.svc.RunRoomAllocation(ctx, &service.RunRoomAllocationParams{
		DepartureID: req.GetDepartureId(),
		JamaahIDs:   req.GetJamaahIds(),
	})
	if err != nil {
		logger.Error().Str("op", op).Err(err).Msg("RunRoomAllocation failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "room allocation failed: %v", err)
	}

	assignments := make([]*pb.RoomAssignmentProto, 0, len(result.Assignments))
	for _, a := range result.Assignments {
		assignments = append(assignments, &pb.RoomAssignmentProto{
			RoomNumber: a.RoomNumber,
			JamaahId:   a.JamaahID,
		})
	}

	span.SetStatus(otelCodes.Ok, "success")
	return &pb.RunRoomAllocationResponse{
		AllocationId: result.AllocationID,
		RoomCount:    result.RoomCount,
		Assignments:  assignments,
	}, nil
}

// GetRoomAllocation handles the GetRoomAllocation RPC.
// Returns the current allocation for a departure.
func (s *Server) GetRoomAllocation(ctx context.Context, req *pb.GetRoomAllocationRequest) (*pb.GetRoomAllocationResponse, error) {
	const op = "grpc_api.Server.GetRoomAllocation"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("op", op),
		attribute.String("departure_id", req.GetDepartureId()),
	)

	logger.Info().
		Str("op", op).
		Str("departure_id", req.GetDepartureId()).
		Msg("")

	if req.GetDepartureId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "departure_id is required")
	}

	result, err := s.svc.GetRoomAllocation(ctx, &service.GetRoomAllocationParams{
		DepartureID: req.GetDepartureId(),
	})
	if err != nil {
		logger.Error().Str("op", op).Err(err).Msg("GetRoomAllocation failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "get room allocation failed: %v", err)
	}

	assignments := make([]*pb.RoomAssignmentProto, 0, len(result.Assignments))
	for _, a := range result.Assignments {
		assignments = append(assignments, &pb.RoomAssignmentProto{
			RoomNumber: a.RoomNumber,
			JamaahId:   a.JamaahID,
		})
	}

	span.SetStatus(otelCodes.Ok, "success")
	return &pb.GetRoomAllocationResponse{
		AllocationId: result.AllocationID,
		RoomCount:    result.RoomCount,
		Assignments:  assignments,
		Status:       result.Status,
	}, nil
}
