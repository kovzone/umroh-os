// Package crm_grpc_adapter is gateway-svc's consumer-side wrapper around
// crm-svc's gRPC surface (S4-E-02 / ADR-0009).
//
// Per ADR-0009, all CRM REST routes (POST /v1/leads, GET /v1/leads,
// GET /v1/leads/:id, PUT /v1/leads/:id) proxy to crm-svc over gRPC via this adapter.
package crm_grpc_adapter

import (
	"context"
	"errors"
	"fmt"

	"gateway-svc/adapter/crm_grpc_adapter/pb"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Adapter wraps crm-svc's CrmServiceClient.
type Adapter struct {
	logger    *zerolog.Logger
	tracer    trace.Tracer
	crmClient pb.CrmServiceClient
}

// NewAdapter creates a crm adapter from an already-dialled conn.
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:    logger,
		tracer:    tracer,
		crmClient: pb.NewCrmServiceClient(cc),
	}
}

// ---------------------------------------------------------------------------
// CreateLead
// ---------------------------------------------------------------------------

type CreateLeadParams struct {
	Name                string
	Phone               string
	Email               string
	Source              string
	UtmSource           string
	UtmMedium           string
	UtmCampaign         string
	UtmContent          string
	UtmTerm             string
	InterestPackageID   string
	InterestDepartureID string
	Notes               string
}

type LeadResult struct {
	ID                  string
	Source              string
	UtmSource           string
	UtmMedium           string
	UtmCampaign         string
	UtmContent          string
	UtmTerm             string
	Name                string
	Phone               string
	Email               string
	InterestPackageID   string
	InterestDepartureID string
	Status              string
	AssignedCsID        string
	Notes               string
	BookingID           string
	CreatedAt           string
	UpdatedAt           string
}

func (a *Adapter) CreateLead(ctx context.Context, params *CreateLeadParams) (*LeadResult, error) {
	const op = "crm_grpc_adapter.Adapter.CreateLead"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("phone", params.Phone))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.crmClient.CreateLead(ctx, &pb.CreateLeadRequest{
		Name:                params.Name,
		Phone:               params.Phone,
		Email:               params.Email,
		Source:              params.Source,
		UtmSource:           params.UtmSource,
		UtmMedium:           params.UtmMedium,
		UtmCampaign:         params.UtmCampaign,
		UtmContent:          params.UtmContent,
		UtmTerm:             params.UtmTerm,
		InterestPackageId:   params.InterestPackageID,
		InterestDepartureId: params.InterestDepartureID,
		Notes:               params.Notes,
	})
	if err != nil {
		wrapped := mapCrmError(err)
		logger.Warn().Err(wrapped).Msg("crm-svc.CreateLead failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "created")
	return mapLeadResponse(resp), nil
}

// ---------------------------------------------------------------------------
// GetLead
// ---------------------------------------------------------------------------

func (a *Adapter) GetLead(ctx context.Context, id string) (*LeadResult, error) {
	const op = "crm_grpc_adapter.Adapter.GetLead"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("lead_id", id))

	resp, err := a.crmClient.GetLead(ctx, &pb.GetLeadRequest{Id: id})
	if err != nil {
		wrapped := mapCrmError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return mapLeadResponse(resp), nil
}

// ---------------------------------------------------------------------------
// UpdateLead
// ---------------------------------------------------------------------------

type UpdateLeadParams struct {
	ID           string
	Status       string
	Notes        string
	AssignedCsID string
}

func (a *Adapter) UpdateLead(ctx context.Context, params *UpdateLeadParams) (*LeadResult, error) {
	const op = "crm_grpc_adapter.Adapter.UpdateLead"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("lead_id", params.ID))

	resp, err := a.crmClient.UpdateLead(ctx, &pb.UpdateLeadRequest{
		Id:           params.ID,
		Status:       params.Status,
		Notes:        params.Notes,
		AssignedCsId: params.AssignedCsID,
	})
	if err != nil {
		wrapped := mapCrmError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "updated")
	return mapLeadResponse(resp), nil
}

// ---------------------------------------------------------------------------
// ListLeads
// ---------------------------------------------------------------------------

type ListLeadsParams struct {
	StatusFilter     string
	AssignedCsFilter string
	Page             int32
	PageSize         int32
}

type ListLeadsResult struct {
	Leads    []*LeadResult
	Total    int64
	Page     int32
	PageSize int32
}

func (a *Adapter) ListLeads(ctx context.Context, params *ListLeadsParams) (*ListLeadsResult, error) {
	const op = "crm_grpc_adapter.Adapter.ListLeads"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	resp, err := a.crmClient.ListLeads(ctx, &pb.ListLeadsRequest{
		StatusFilter:       params.StatusFilter,
		AssignedCsIdFilter: params.AssignedCsFilter,
		Page:               params.Page,
		PageSize:           params.PageSize,
	})
	if err != nil {
		wrapped := mapCrmError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	leads := make([]*LeadResult, 0, len(resp.GetLeads()))
	for _, l := range resp.GetLeads() {
		leads = append(leads, mapLeadResponse(l))
	}

	span.SetStatus(codes.Ok, "ok")
	return &ListLeadsResult{
		Leads:    leads,
		Total:    resp.GetTotal(),
		Page:     resp.GetPage(),
		PageSize: resp.GetPageSize(),
	}, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func mapLeadResponse(r *pb.LeadResponse) *LeadResult {
	return &LeadResult{
		ID:                  r.GetId(),
		Source:              r.GetSource(),
		UtmSource:           r.GetUtmSource(),
		UtmMedium:           r.GetUtmMedium(),
		UtmCampaign:         r.GetUtmCampaign(),
		UtmContent:          r.GetUtmContent(),
		UtmTerm:             r.GetUtmTerm(),
		Name:                r.GetName(),
		Phone:               r.GetPhone(),
		Email:               r.GetEmail(),
		InterestPackageID:   r.GetInterestPackageId(),
		InterestDepartureID: r.GetInterestDepartureId(),
		Status:              r.GetStatus(),
		AssignedCsID:        r.GetAssignedCsId(),
		Notes:               r.GetNotes(),
		BookingID:           r.GetBookingId(),
		CreatedAt:           r.GetCreatedAt(),
		UpdatedAt:           r.GetUpdatedAt(),
	}
}

func mapCrmError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("crm call failed: %w", err))
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return errors.Join(apperrors.ErrNotFound, errors.New(st.Message()))
	case grpcCodes.InvalidArgument:
		return errors.Join(apperrors.ErrValidation, errors.New(st.Message()))
	case grpcCodes.AlreadyExists:
		return errors.Join(apperrors.ErrConflict, errors.New(st.Message()))
	case grpcCodes.Unauthenticated:
		return errors.Join(apperrors.ErrUnauthorized, errors.New(st.Message()))
	case grpcCodes.PermissionDenied:
		return errors.Join(apperrors.ErrForbidden, errors.New(st.Message()))
	case grpcCodes.Unavailable:
		return errors.Join(apperrors.ErrServiceUnavailable, errors.New(st.Message()))
	default:
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("crm call failed: %s", st.Message()))
	}
}
