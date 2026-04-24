package grpc_api

// BL-CAT-013 — gRPC handler for GetPackageVersion.

import (
	"context"

	"catalog-svc/api/grpc_api/pb"
	"catalog-svc/service"
	"catalog-svc/util/apperrors"
	"catalog-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

// GetPackageVersion returns the computed version hash for a package.
func (s *Server) GetPackageVersion(ctx context.Context, req *pb.GetPackageVersionRequest) (*pb.GetPackageVersionResponse, error) {
	const op = "grpc_api.Server.GetPackageVersion"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "GetPackageVersion"),
		attribute.String("input.package_id", req.GetPackageId()),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.GetPackageVersion(ctx, &service.GetPackageVersionParams{
		PackageID: req.GetPackageId(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.GetPackageVersionResponse{
		PackageId: result.PackageID,
		Version:   result.Version,
		UpdatedAt: result.UpdatedAt,
	}, nil
}
