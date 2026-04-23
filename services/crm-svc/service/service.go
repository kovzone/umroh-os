package service

import (
	"context"

	"crm-svc/store/postgres_store"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// IService is the business-layer interface for crm-svc.
//
// Scaffold methods:
//   - Liveness / Readiness / DbTxDiagnostic
//
// S4-E-02 adds CRM lead management:
//   - CreateLead, GetLead, UpdateLead, ListLeads
//   - OnBookingCreated, OnBookingPaidInFull (event callbacks)
type IService interface {
	Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error)
	Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error)
	DbTxDiagnostic(ctx context.Context, params *DbTxDiagnosticParams) (*DbTxDiagnosticResult, error)

	// S4-E-02 — CRM lead management (BL-CRM-001..003)
	CreateLead(ctx context.Context, params *CreateLeadParams) (*LeadResult, error)
	GetLead(ctx context.Context, params *GetLeadParams) (*LeadResult, error)
	UpdateLead(ctx context.Context, params *UpdateLeadParams) (*LeadResult, error)
	ListLeads(ctx context.Context, params *ListLeadsParams) (*ListLeadsResult, error)
	OnBookingCreated(ctx context.Context, params *OnBookingCreatedParams) (*OnBookingCreatedResult, error)
	OnBookingPaidInFull(ctx context.Context, params *OnBookingPaidInFullParams) (*OnBookingPaidInFullResult, error)
}

type Service struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	// appName is stamped on every row this service writes to the shared
	// public.diagnostics table so debugging can attribute rows to their origin.
	appName string

	store postgres_store.IStore
}

func NewService(
	logger *zerolog.Logger,
	tracer trace.Tracer,
	appName string,
	store postgres_store.IStore,
) IService {
	return &Service{
		logger:  logger,
		tracer:  tracer,
		appName: appName,
		store:   store,
	}
}
