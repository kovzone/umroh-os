// visa.go — gateway-svc adapter methods for visa-svc Phase 6 RPCs
// (BL-VISA-001..003): TransitionStatus, BulkSubmit, GetApplications.

package visa_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/visa_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// TransitionStatus
// ---------------------------------------------------------------------------

// TransitionStatusParams holds inputs for PUT /v1/visas/:id/status.
type TransitionStatusParams struct {
	ApplicationID string
	ToStatus      string
	Reason        string
	ActorUserID   string
}

// TransitionStatusResult is the response.
type TransitionStatusResult struct {
	ApplicationID string
	FromStatus    string
	ToStatus      string
	Idempotent    bool
}

func (a *Adapter) TransitionStatus(ctx context.Context, params *TransitionStatusParams) (*TransitionStatusResult, error) {
	const op = "visa_grpc_adapter.Adapter.TransitionStatus"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("application_id", params.ApplicationID).Msg("")

	resp, err := a.visaClient.TransitionStatus(ctx, &pb.TransitionStatusRequest{
		ApplicationId: params.ApplicationID,
		ToStatus:      params.ToStatus,
		Reason:        params.Reason,
		ActorUserId:   params.ActorUserID,
	})
	if err != nil {
		wrapped := mapVisaError(err)
		logger.Warn().Err(wrapped).Msg("visa-svc.TransitionStatus failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &TransitionStatusResult{
		ApplicationID: resp.GetApplicationId(),
		FromStatus:    resp.GetFromStatus(),
		ToStatus:      resp.GetToStatus(),
		Idempotent:    resp.GetIdempotent(),
	}, nil
}

// ---------------------------------------------------------------------------
// BulkSubmit
// ---------------------------------------------------------------------------

// BulkSubmitParams holds inputs for POST /v1/visas/bulk-submit.
type BulkSubmitParams struct {
	DepartureID string
	JamaahIDs   []string
	ProviderID  string
}

// BulkSubmitResult is the response.
type BulkSubmitResult struct {
	SubmittedCount int32
	ApplicationIDs []string
}

func (a *Adapter) BulkSubmit(ctx context.Context, params *BulkSubmitParams) (*BulkSubmitResult, error) {
	const op = "visa_grpc_adapter.Adapter.BulkSubmit"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("departure_id", params.DepartureID).Msg("")

	resp, err := a.visaClient.BulkSubmit(ctx, &pb.BulkSubmitRequest{
		DepartureId: params.DepartureID,
		JamaahIds:   params.JamaahIDs,
		ProviderId:  params.ProviderID,
	})
	if err != nil {
		wrapped := mapVisaError(err)
		logger.Warn().Err(wrapped).Msg("visa-svc.BulkSubmit failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &BulkSubmitResult{
		SubmittedCount: resp.GetSubmittedCount(),
		ApplicationIDs: resp.GetApplicationIds(),
	}, nil
}

// ---------------------------------------------------------------------------
// GetApplications
// ---------------------------------------------------------------------------

// StatusHistoryEntry holds one status transition in application history.
type StatusHistoryEntry struct {
	FromStatus string
	ToStatus   string
	Reason     string
	CreatedAt  string
}

// ApplicationRecord holds one visa application with history.
type ApplicationRecord struct {
	ID          string
	JamaahID    string
	Status      string
	ProviderRef string
	IssuedDate  string
	History     []*StatusHistoryEntry
}

// GetApplicationsResult is the response for GET /v1/visas.
type GetApplicationsResult struct {
	Applications []*ApplicationRecord
}

func (a *Adapter) GetApplications(ctx context.Context, departureID, statusFilter string) (*GetApplicationsResult, error) {
	const op = "visa_grpc_adapter.Adapter.GetApplications"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("departure_id", departureID).Msg("")

	resp, err := a.visaClient.GetApplications(ctx, &pb.GetApplicationsRequest{
		DepartureId:  departureID,
		StatusFilter: statusFilter,
	})
	if err != nil {
		wrapped := mapVisaError(err)
		logger.Warn().Err(wrapped).Msg("visa-svc.GetApplications failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	apps := make([]*ApplicationRecord, 0, len(resp.GetApplications()))
	for _, a := range resp.GetApplications() {
		hist := make([]*StatusHistoryEntry, 0, len(a.GetHistory()))
		for _, h := range a.GetHistory() {
			hist = append(hist, &StatusHistoryEntry{
				FromStatus: h.GetFromStatus(),
				ToStatus:   h.GetToStatus(),
				Reason:     h.GetReason(),
				CreatedAt:  h.GetCreatedAt(),
			})
		}
		apps = append(apps, &ApplicationRecord{
			ID:          a.GetId(),
			JamaahID:    a.GetJamaahId(),
			Status:      a.GetStatus(),
			ProviderRef: a.GetProviderRef(),
			IssuedDate:  a.GetIssuedDate(),
			History:     hist,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &GetApplicationsResult{Applications: apps}, nil
}
