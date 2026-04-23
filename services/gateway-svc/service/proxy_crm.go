// proxy_crm.go — gateway service layer for CRM lead management (S4-E-02).
//
// Thin delegation layer: service methods map 1:1 to crm_grpc_adapter calls.
// No business logic here; all CRM logic lives in crm-svc.

package service

import (
	"context"

	"gateway-svc/adapter/crm_grpc_adapter"
)

func (s *Service) CreateLead(ctx context.Context, params *crm_grpc_adapter.CreateLeadParams) (*crm_grpc_adapter.LeadResult, error) {
	return s.adapters.crmGrpc.CreateLead(ctx, params)
}

func (s *Service) GetLead(ctx context.Context, id string) (*crm_grpc_adapter.LeadResult, error) {
	return s.adapters.crmGrpc.GetLead(ctx, id)
}

func (s *Service) UpdateLead(ctx context.Context, params *crm_grpc_adapter.UpdateLeadParams) (*crm_grpc_adapter.LeadResult, error) {
	return s.adapters.crmGrpc.UpdateLead(ctx, params)
}

func (s *Service) ListLeads(ctx context.Context, params *crm_grpc_adapter.ListLeadsParams) (*crm_grpc_adapter.ListLeadsResult, error) {
	return s.adapters.crmGrpc.ListLeads(ctx, params)
}
