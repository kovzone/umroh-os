package grpc_api

// vendor_readiness.go — gRPC handler implementations for vendor readiness RPCs
// (BL-OPS-020). Implements pb.VendorReadinessHandler on *Server.
//
// Both RPCs gate on catalog.package.manage permission so only staff with
// departure-management access can read or update readiness state.

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

// UpdateVendorReadiness upserts one readiness kind for a departure and returns
// the full (ticket/hotel/visa) readiness summary.
func (s *Server) UpdateVendorReadiness(ctx context.Context, req *pb.UpdateVendorReadinessRequest) (*pb.UpdateVendorReadinessResponse, error) {
	const op = "grpc_api.Server.UpdateVendorReadiness"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("departure_id", req.GetDepartureId()),
		attribute.String("kind", req.GetKind()),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}
	if req.GetDepartureId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "departure_id is required")
	}

	result, err := s.svc.UpdateVendorReadiness(ctx, &service.UpdateVendorReadinessParams{
		DepartureID:   req.GetDepartureId(),
		Kind:          req.GetKind(),
		State:         req.GetState(),
		Notes:         req.GetNotes(),
		AttachmentURL: req.GetAttachmentUrl(),
		UpdatedBy:     req.GetUserId(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.UpdateVendorReadinessResponse{
		Readiness: readinessToProto(result),
	}, nil
}

// GetDepartureReadiness returns the current vendor readiness summary for a
// departure. Missing kinds (no row yet) default to "not_started".
func (s *Server) GetDepartureReadiness(ctx context.Context, req *pb.GetDepartureReadinessRequest) (*pb.GetDepartureReadinessResponse, error) {
	const op = "grpc_api.Server.GetDepartureReadiness"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("departure_id", req.GetDepartureId()),
	)

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}
	if req.GetDepartureId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrValidation), "departure_id is required")
	}

	result, err := s.svc.GetDepartureReadiness(ctx, &service.GetDepartureReadinessParams{
		DepartureID: req.GetDepartureId(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.GetDepartureReadinessResponse{
		Readiness: readinessToProto(result),
	}, nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// readinessToProto converts the service-layer VendorReadiness to the pb type.
func readinessToProto(r *service.VendorReadiness) *pb.VendorReadiness {
	if r == nil {
		return &pb.VendorReadiness{
			TicketState: "not_started",
			HotelState:  "not_started",
			VisaState:   "not_started",
		}
	}
	return &pb.VendorReadiness{
		TicketState: r.Ticket,
		HotelState:  r.Hotel,
		VisaState:   r.Visa,
	}
}
