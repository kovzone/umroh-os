// leads.go — gRPC handlers for CRM lead management (S4-E-02 / BL-CRM-001..003).
//
// Implements CrmServiceServer:
//   - CreateLead
//   - GetLead
//   - UpdateLead
//   - ListLeads
//   - OnBookingCreated
//   - OnBookingPaidInFull
//
// Per ADR-0009: crm-svc is pure gRPC — no REST surface here.
// Per ADR-0006: all calls are synchronous.

package grpc_api

import (
	"context"
	"errors"
	"strings"

	"crm-svc/api/grpc_api/pb"
	"crm-svc/service"
	"crm-svc/util/apperrors"
	"crm-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// CreateLead
// ---------------------------------------------------------------------------

func (s *Server) CreateLead(ctx context.Context, req *pb.CreateLeadRequest) (*pb.LeadResponse, error) {
	const op = "grpc_api.Server.CreateLead"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("phone", req.GetPhone()).Msg("")

	result, err := s.svc.CreateLead(ctx, &service.CreateLeadParams{
		Name:                req.GetName(),
		Phone:               req.GetPhone(),
		Email:               req.GetEmail(),
		Source:              req.GetSource(),
		UtmSource:           req.GetUtmSource(),
		UtmMedium:           req.GetUtmMedium(),
		UtmCampaign:         req.GetUtmCampaign(),
		UtmContent:          req.GetUtmContent(),
		UtmTerm:             req.GetUtmTerm(),
		InterestPackageID:   req.GetInterestPackageId(),
		InterestDepartureID: req.GetInterestDepartureId(),
		Notes:               req.GetNotes(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "created")
	span.SetAttributes(attribute.String("lead_id", result.ID))
	return mapLeadResult(result), nil
}

// ---------------------------------------------------------------------------
// GetLead
// ---------------------------------------------------------------------------

func (s *Server) GetLead(ctx context.Context, req *pb.GetLeadRequest) (*pb.LeadResponse, error) {
	const op = "grpc_api.Server.GetLead"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	result, err := s.svc.GetLead(ctx, &service.GetLeadParams{ID: req.GetId()})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "ok")
	return mapLeadResult(result), nil
}

// ---------------------------------------------------------------------------
// UpdateLead
// ---------------------------------------------------------------------------

func (s *Server) UpdateLead(ctx context.Context, req *pb.UpdateLeadRequest) (*pb.LeadResponse, error) {
	const op = "grpc_api.Server.UpdateLead"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("id", req.GetId()).Str("status", req.GetStatus()).Msg("")

	result, err := s.svc.UpdateLead(ctx, &service.UpdateLeadParams{
		ID:           req.GetId(),
		Status:       req.GetStatus(),
		Notes:        req.GetNotes(),
		AssignedCsID: req.GetAssignedCsId(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		// Map "invalid transition" → InvalidArgument
		if strings.Contains(err.Error(), "invalid transition") {
			return nil, grpcStatus.Error(apperrors.GRPCCode(errors.Join(apperrors.ErrValidation, err)), err.Error())
		}
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "updated")
	return mapLeadResult(result), nil
}

// ---------------------------------------------------------------------------
// ListLeads
// ---------------------------------------------------------------------------

func (s *Server) ListLeads(ctx context.Context, req *pb.ListLeadsRequest) (*pb.ListLeadsResponse, error) {
	const op = "grpc_api.Server.ListLeads"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	result, err := s.svc.ListLeads(ctx, &service.ListLeadsParams{
		StatusFilter:     req.GetStatusFilter(),
		AssignedCsFilter: req.GetAssignedCsIdFilter(),
		Page:             req.GetPage(),
		PageSize:         req.GetPageSize(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), err.Error())
	}

	leads := make([]*pb.LeadResponse, 0, len(result.Leads))
	for _, l := range result.Leads {
		leads = append(leads, mapLeadResult(l))
	}

	span.SetStatus(codes.Ok, "ok")
	return &pb.ListLeadsResponse{
		Leads:    leads,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}

// ---------------------------------------------------------------------------
// OnBookingCreated
// ---------------------------------------------------------------------------

func (s *Server) OnBookingCreated(ctx context.Context, req *pb.OnBookingCreatedRequest) (*pb.OnBookingCreatedResponse, error) {
	const op = "grpc_api.Server.OnBookingCreated"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().
		Str("op", op).
		Str("booking_id", req.GetBookingId()).
		Str("lead_id", req.GetLeadId()).
		Msg("")

	result, err := s.svc.OnBookingCreated(ctx, &service.OnBookingCreatedParams{
		BookingID:   req.GetBookingId(),
		LeadID:      req.GetLeadId(),
		PackageID:   req.GetPackageId(),
		DepartureID: req.GetDepartureId(),
		JamaahCount: req.GetJamaahCount(),
		CreatedAt:   req.GetCreatedAt(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "ok")
	return &pb.OnBookingCreatedResponse{
		Updated: result.Updated,
		LeadId:  result.LeadID,
	}, nil
}

// ---------------------------------------------------------------------------
// OnBookingPaidInFull
// ---------------------------------------------------------------------------

func (s *Server) OnBookingPaidInFull(ctx context.Context, req *pb.OnBookingPaidInFullRequest) (*pb.OnBookingPaidInFullResponse, error) {
	const op = "grpc_api.Server.OnBookingPaidInFull"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().
		Str("op", op).
		Str("booking_id", req.GetBookingId()).
		Str("lead_id", req.GetLeadId()).
		Msg("")

	result, err := s.svc.OnBookingPaidInFull(ctx, &service.OnBookingPaidInFullParams{
		BookingID: req.GetBookingId(),
		LeadID:    req.GetLeadId(),
		PaidAt:    req.GetPaidAt(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, grpcStatus.Error(apperrors.GRPCCode(err), err.Error())
	}

	span.SetStatus(codes.Ok, "ok")
	return &pb.OnBookingPaidInFullResponse{
		Updated: result.Updated,
		LeadId:  result.LeadID,
	}, nil
}

// ---------------------------------------------------------------------------
// mapLeadResult — service.LeadResult → pb.LeadResponse
// ---------------------------------------------------------------------------

func mapLeadResult(r *service.LeadResult) *pb.LeadResponse {
	return &pb.LeadResponse{
		Id:                  r.ID,
		Source:              r.Source,
		UtmSource:           r.UtmSource,
		UtmMedium:           r.UtmMedium,
		UtmCampaign:         r.UtmCampaign,
		UtmContent:          r.UtmContent,
		UtmTerm:             r.UtmTerm,
		Name:                r.Name,
		Phone:               r.Phone,
		Email:               r.Email,
		InterestPackageId:   r.InterestPackageID,
		InterestDepartureId: r.InterestDepartureID,
		Status:              r.Status,
		AssignedCsId:        r.AssignedCsID,
		Notes:               r.Notes,
		BookingId:           r.BookingID,
		CreatedAt:           r.CreatedAt,
		UpdatedAt:           r.UpdatedAt,
	}
}
