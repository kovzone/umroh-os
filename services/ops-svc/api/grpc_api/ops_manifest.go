// ops_manifest.go — gRPC handler for ExportManifest RPC (BL-OPS-001).

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

// ExportManifest handles the ExportManifest RPC.
// Returns structured manifest rows for a departure, suitable for CSV/XLSX
// formatting by the client. Currently returns empty rows (wired to jamaah-svc
// in a later sprint).
func (s *Server) ExportManifest(ctx context.Context, req *pb.ExportManifestRequest) (*pb.ExportManifestResponse, error) {
	const op = "grpc_api.Server.ExportManifest"

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

	result, err := s.svc.ExportManifest(ctx, &service.ExportManifestParams{
		DepartureID: req.GetDepartureId(),
	})
	if err != nil {
		logger.Error().Str("op", op).Err(err).Msg("ExportManifest failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "export manifest failed: %v", err)
	}

	rows := make([]*pb.ManifestRowProto, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, &pb.ManifestRowProto{
			No:         r.No,
			JamaahName: r.JamaahName,
			PassportNo: r.PassportNo,
			DocStatus:  r.DocStatus,
			RoomNumber: r.RoomNumber,
		})
	}

	span.SetStatus(otelCodes.Ok, "success")
	return &pb.ExportManifestResponse{
		Rows: rows,
	}, nil
}
