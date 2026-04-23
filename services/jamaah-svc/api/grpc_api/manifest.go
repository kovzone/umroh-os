package grpc_api

// manifest.go — gRPC handler for GetDepartureManifest RPC.
//
// Returns a departure manifest: one row per active pilgrim with booking status
// and document completion summary.
//
// Wave-1A manifest API.

import (
	"context"

	"jamaah-svc/api/grpc_api/pb"
	"jamaah-svc/service"
	"jamaah-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// GetDepartureManifest returns the manifest for a departure.
func (s *Server) GetDepartureManifest(ctx context.Context, req *pb.GetDepartureManifestRequest) (*pb.GetDepartureManifestResponse, error) {
	const op = "grpc_api.Server.GetDepartureManifest"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("departure_id", req.DepartureID),
	)

	if req.DepartureID == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "departure_id is required")
	}

	result, err := s.svc.GetDepartureManifest(ctx, &service.GetDepartureManifestParams{
		DepartureID: req.DepartureID,
	})
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("departure_id", req.DepartureID).
			Err(err).
			Msg("GetDepartureManifest failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "failed to get departure manifest: %v", err)
	}

	jamaahList := make([]*pb.ManifestJamaah, 0, len(result.JamaahList))
	for _, j := range result.JamaahList {
		jamaahList = append(jamaahList, &pb.ManifestJamaah{
			BookingID:     j.BookingID,
			Name:          j.Name,
			NIK:           j.NIK,
			Phone:         j.Phone,
			RoomType:      j.RoomType,
			BookingStatus: j.BookingStatus,
			DocStatus:     j.DocStatus,
		})
	}

	span.SetStatus(otelCodes.Ok, "success")
	return &pb.GetDepartureManifestResponse{
		DepartureID: result.DepartureID,
		TotalJamaah: result.TotalJamaah,
		LunasPaid:   result.LunasPaid,
		DocComplete: result.DocComplete,
		JamaahList:  jamaahList,
	}, nil
}
