package grpc_api

// S1-E-07 / BL-CAT-014 — gRPC handlers for staff catalog write RPCs.
//
// Auth contract: every handler here expects req.UserId to be non-empty
// (populated by gateway-svc from the validated PASETO/JWT claims).
// The handler calls iam-svc.CheckPermission before delegating to the
// service layer via the server's iamClient field. If iamClient is nil
// (legacy mode / tests) the gate is skipped and a WARN is logged.
//
// Permission required: catalog.package.manage (action: manage, resource: catalog.package)
// Roles that should hold this permission: catalog_manager, admin.
//
// Error mapping: service-layer apperrors.Err* → gRPC status codes via
// apperrors.GRPCCode (same helper used by read handlers).

import (
	"context"
	"errors"
	"fmt"

	"catalog-svc/adapter/iam_grpc_adapter"
	"catalog-svc/api/grpc_api/pb"
	"catalog-svc/service"
	"catalog-svc/util/apperrors"
	"catalog-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
)

// checkCatalogManagePermission validates that the user has the catalog.package.manage
// permission via iam-svc. Returns a gRPC-ready error on deny/failure, nil on allow.
// If iamClient is nil (test/legacy mode), the gate is skipped with a WARN log.
func (s *Server) checkCatalogManagePermission(ctx context.Context, userID string) error {
	if s.iamClient == nil {
		logger := logging.LogWithTrace(ctx, s.logger)
		logger.Warn().Str("user_id", userID).Msg("iamClient is nil — permission gate skipped (test/legacy mode)")
		return nil
	}
	result, err := s.iamClient.CheckPermission(ctx, &iam_grpc_adapter.CheckPermissionParams{
		UserID:   userID,
		Resource: "catalog.package",
		Action:   "manage",
		Scope:    "global",
	})
	if err != nil {
		return status.Error(apperrors.GRPCCode(errors.Join(apperrors.ErrUnauthorized, fmt.Errorf("permission check failed: %w", err))), "permission check failed")
	}
	if !result.Allowed {
		return status.Error(apperrors.GRPCCode(apperrors.ErrForbidden), "catalog.package.manage permission required")
	}
	return nil
}

// ---------------------------------------------------------------------------
// CreatePackage
// ---------------------------------------------------------------------------

func (s *Server) CreatePackage(ctx context.Context, req *pb.CreatePackageRequest) (*pb.CreatePackageResponse, error) {
	const op = "grpc_api.Server.CreatePackage"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "CreatePackage"))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	params := &service.CreatePackageParams{
		UserID:        req.GetUserId(),
		BranchID:      req.GetBranchId(),
		Kind:          req.GetKind(),
		Name:          req.GetName(),
		Description:   req.GetDescription(),
		CoverPhotoUrl: req.GetCoverPhotoUrl(),
		Highlights:    req.GetHighlights(),
		ItineraryID:   req.GetItineraryId(),
		AirlineID:     req.GetAirlineId(),
		MuthawwifID:   req.GetMuthawwifId(),
		HotelIDs:      req.GetHotelIds(),
		AddonIDs:      req.GetAddonIds(),
		Status:        req.GetStatus(),
	}

	detail, err := s.svc.CreatePackage(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.CreatePackageResponse{Package: packageDetailToProto(detail)}, nil
}

// ---------------------------------------------------------------------------
// UpdatePackage
// ---------------------------------------------------------------------------

func (s *Server) UpdatePackage(ctx context.Context, req *pb.UpdatePackageRequest) (*pb.UpdatePackageResponse, error) {
	const op = "grpc_api.Server.UpdatePackage"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "UpdatePackage"))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	params := &service.UpdatePackageParams{
		UserID:        req.GetUserId(),
		BranchID:      req.GetBranchId(),
		ID:            req.GetId(),
		Name:          req.GetName(),
		Description:   req.GetDescription(),
		CoverPhotoUrl: req.GetCoverPhotoUrl(),
		ItineraryID:   req.GetItineraryId(),
		AirlineID:     req.GetAirlineId(),
		MuthawwifID:   req.GetMuthawwifId(),
		Status:        req.GetStatus(),
	}
	// Highlights: proto repeated string; nil vs empty treated differently.
	// Empty slice in proto = "caller explicitly sent []" = replace with [].
	// We propagate it as-is; service treats non-nil as "replace".
	if len(req.GetHighlights()) > 0 {
		params.Highlights = req.GetHighlights()
	}
	if len(req.GetHotelIds()) > 0 {
		params.HotelIDs = req.GetHotelIds()
	}
	if len(req.GetAddonIds()) > 0 {
		params.AddonIDs = req.GetAddonIds()
	}

	detail, err := s.svc.UpdatePackage(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.UpdatePackageResponse{Package: packageDetailToProto(detail)}, nil
}

// ---------------------------------------------------------------------------
// DeletePackage
// ---------------------------------------------------------------------------

func (s *Server) DeletePackage(ctx context.Context, req *pb.DeletePackageRequest) (*pb.DeletePackageResponse, error) {
	const op = "grpc_api.Server.DeletePackage"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "DeletePackage"))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	err := s.svc.DeletePackage(ctx, &service.DeletePackageParams{
		UserID:   req.GetUserId(),
		BranchID: req.GetBranchId(),
		ID:       req.GetId(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.DeletePackageResponse{Ok: true}, nil
}

// ---------------------------------------------------------------------------
// CreateDeparture
// ---------------------------------------------------------------------------

func (s *Server) CreateDeparture(ctx context.Context, req *pb.CreateDepartureRequest) (*pb.CreateDepartureResponse, error) {
	const op = "grpc_api.Server.CreateDeparture"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "CreateDeparture"))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	pricing := make([]service.PricingInput, 0, len(req.GetPricing()))
	for _, p := range req.GetPricing() {
		pricing = append(pricing, service.PricingInput{
			RoomType:           p.GetRoomType(),
			ListAmount:         p.GetListAmount(),
			ListCurrency:       p.GetListCurrency(),
			SettlementCurrency: p.GetSettlementCurrency(),
		})
	}

	detail, err := s.svc.CreateDeparture(ctx, &service.CreateDepartureParams{
		UserID:        req.GetUserId(),
		BranchID:      req.GetBranchId(),
		PackageID:     req.GetPackageId(),
		DepartureDate: req.GetDepartureDate(),
		ReturnDate:    req.GetReturnDate(),
		TotalSeats:    int(req.GetTotalSeats()),
		Status:        req.GetStatus(),
		Pricing:       pricing,
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.CreateDepartureResponse{Departure: departureDetailToProto(detail)}, nil
}

// ---------------------------------------------------------------------------
// UpdateDeparture
// ---------------------------------------------------------------------------

func (s *Server) UpdateDeparture(ctx context.Context, req *pb.UpdateDepartureRequest) (*pb.UpdateDepartureResponse, error) {
	const op = "grpc_api.Server.UpdateDeparture"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "UpdateDeparture"))

	logger := logging.LogWithTrace(ctx, s.logger)

	if req.GetUserId() == "" {
		return nil, status.Error(apperrors.GRPCCode(apperrors.ErrUnauthorized), "user_id is required")
	}
	if err := s.checkCatalogManagePermission(ctx, req.GetUserId()); err != nil {
		return nil, err
	}

	var pricing []service.PricingInput
	if len(req.GetPricing()) > 0 {
		pricing = make([]service.PricingInput, 0, len(req.GetPricing()))
		for _, p := range req.GetPricing() {
			pricing = append(pricing, service.PricingInput{
				RoomType:           p.GetRoomType(),
				ListAmount:         p.GetListAmount(),
				ListCurrency:       p.GetListCurrency(),
				SettlementCurrency: p.GetSettlementCurrency(),
			})
		}
	}

	detail, err := s.svc.UpdateDeparture(ctx, &service.UpdateDepartureParams{
		UserID:        req.GetUserId(),
		BranchID:      req.GetBranchId(),
		ID:            req.GetId(),
		DepartureDate: req.GetDepartureDate(),
		ReturnDate:    req.GetReturnDate(),
		TotalSeats:    int(req.GetTotalSeats()),
		Status:        req.GetStatus(),
		Pricing:       pricing,
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.UpdateDepartureResponse{Departure: departureDetailToProto(detail)}, nil
}

// ---------------------------------------------------------------------------
// ReserveSeats
// ---------------------------------------------------------------------------

func (s *Server) ReserveSeats(ctx context.Context, req *pb.ReserveSeatsRequest) (*pb.ReserveSeatsResponse, error) {
	const op = "grpc_api.Server.ReserveSeats"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "ReserveSeats"))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.ReserveSeats(ctx, &service.ReserveSeatsParams{
		ReservationID:       req.GetReservationId(),
		DepartureID:         req.GetDepartureId(),
		Seats:               int(req.GetSeats()),
		IdempotencyTTLHours: int(req.GetIdempotencyTtlHours()),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.ReserveSeatsResponse{
		Reservation: &pb.Reservation{
			ReservationId: result.ReservationID,
			DepartureId:   result.DepartureID,
			Seats:         int32(result.Seats),
			ReservedAt:    result.ReservedAt,
			ExpiresAt:     result.ExpiresAt,
		},
		RemainingSeats: int32(result.RemainingSeats),
		Replayed:       result.Replayed,
	}, nil
}

// ---------------------------------------------------------------------------
// ReleaseSeats
// ---------------------------------------------------------------------------

func (s *Server) ReleaseSeats(ctx context.Context, req *pb.ReleaseSeatsRequest) (*pb.ReleaseSeatsResponse, error) {
	const op = "grpc_api.Server.ReleaseSeats"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "ReleaseSeats"))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.svc.ReleaseSeats(ctx, &service.ReleaseSeatsParams{
		ReservationID: req.GetReservationId(),
		Seats:         int(req.GetSeats()),
		Reason:        req.GetReason(),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "success")
	return &pb.ReleaseSeatsResponse{
		Released: &pb.Released{
			ReservationId: result.ReservationID,
			DepartureId:   result.DepartureID,
			SeatsReleased: int32(result.SeatsReleased),
			ReleasedAt:    result.ReleasedAt,
		},
		RemainingSeats: int32(result.RemainingSeats),
		Replayed:       result.Replayed,
	}, nil
}
