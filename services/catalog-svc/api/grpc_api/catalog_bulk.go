package grpc_api

// BL-CAT-010 / BL-CAT-011 — gRPC handlers for bulk package import/update RPCs.

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

// ---------------------------------------------------------------------------
// BulkImportPackages
// ---------------------------------------------------------------------------

func (s *Server) BulkImportPackages(ctx context.Context, req *pb.BulkImportPackagesRequest) (*pb.BulkImportPackagesResponse, error) {
	const op = "grpc_api.Server.BulkImportPackages"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "BulkImportPackages"),
		attribute.Int("input.row_count", len(req.GetRows())),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	rows := make([]service.BulkImportRowInput, 0, len(req.GetRows()))
	for _, r := range req.GetRows() {
		rows = append(rows, service.BulkImportRowInput{
			Name:          r.GetName(),
			Kind:          r.GetKind(),
			Description:   r.GetDescription(),
			CoverPhotoUrl: r.GetCoverPhotoUrl(),
			Highlights:    r.GetHighlights(),
			AddonIDs:      r.GetAddonIds(),
			HotelIDs:      r.GetHotelIds(),
			ItineraryID:   r.GetItineraryId(),
			AirlineID:     r.GetAirlineId(),
			MuthawwifID:   r.GetMuthawwifId(),
			Status:        r.GetStatus(),
		})
	}

	result, err := s.svc.BulkImportPackages(ctx, &service.BulkImportPackagesParams{
		UserID:   req.GetUserId(),
		BranchID: req.GetBranchId(),
		Rows:     rows,
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	pbResults := make([]*pb.BulkImportRowResult, 0, len(result.Results))
	for _, r := range result.Results {
		pbResults = append(pbResults, &pb.BulkImportRowResult{
			Index:     int32(r.Index),
			PackageId: r.PackageID,
			Error:     r.Error,
		})
	}

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(
		attribute.Int("output.successful", result.Successful),
		attribute.Int("output.failed", result.Failed),
	)
	return &pb.BulkImportPackagesResponse{
		Results:    pbResults,
		TotalRows:  int32(result.TotalRows),
		Successful: int32(result.Successful),
		Failed:     int32(result.Failed),
	}, nil
}

// ---------------------------------------------------------------------------
// BulkUpdatePackages
// ---------------------------------------------------------------------------

func (s *Server) BulkUpdatePackages(ctx context.Context, req *pb.BulkUpdatePackagesRequest) (*pb.BulkUpdatePackagesResponse, error) {
	const op = "grpc_api.Server.BulkUpdatePackages"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "BulkUpdatePackages"),
		attribute.Int("input.row_count", len(req.GetRows())),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	rows := make([]service.BulkUpdateRowInput, 0, len(req.GetRows()))
	for _, r := range req.GetRows() {
		rows = append(rows, service.BulkUpdateRowInput{
			ID:          r.GetId(),
			Name:        r.GetName(),
			Description: r.GetDescription(),
			Status:      r.GetStatus(),
			Highlights:  r.GetHighlights(),
			AddonIDs:    r.GetAddonIds(),
			HotelIDs:    r.GetHotelIds(),
		})
	}

	result, err := s.svc.BulkUpdatePackages(ctx, &service.BulkUpdatePackagesParams{
		UserID:   req.GetUserId(),
		BranchID: req.GetBranchId(),
		Rows:     rows,
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	pbResults := make([]*pb.BulkUpdateRowResult, 0, len(result.Results))
	for _, r := range result.Results {
		pbResults = append(pbResults, &pb.BulkUpdateRowResult{
			Index:     int32(r.Index),
			PackageId: r.PackageID,
			Error:     r.Error,
		})
	}

	span.SetStatus(codes.Ok, "success")
	span.SetAttributes(
		attribute.Int("output.successful", result.Successful),
		attribute.Int("output.failed", result.Failed),
	)
	return &pb.BulkUpdatePackagesResponse{
		Results:    pbResults,
		TotalRows:  int32(result.TotalRows),
		Successful: int32(result.Successful),
		Failed:     int32(result.Failed),
	}, nil
}
