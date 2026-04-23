// visa.go — visa pipeline service logic (BL-VISA-001..003).
//
// TransitionStatus: single-application state-machine transition + history row.
// BulkSubmit:       atomic READY→SUBMITTED batch (all-or-nothing transaction).
// GetApplications:  list applications for a departure with embedded history.

package service

import (
	"context"
	"errors"
	"time"

	"visa-svc/store/postgres_store/sqlc"
	"visa-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Allowed state-machine transitions
// ---------------------------------------------------------------------------

var allowedTransitions = map[string][]string{
	"WAITING_DOCS":        {"READY"},
	"READY":               {"SUBMITTED", "CANCELLED"},
	"SUBMITTED":           {"ISSUED", "REJECTED_BY_EMBASSY"},
	"REJECTED_BY_EMBASSY": {"READY"},
}

var knownStatuses = map[string]bool{
	"WAITING_DOCS":        true,
	"READY":               true,
	"SUBMITTED":           true,
	"ISSUED":              true,
	"REJECTED_BY_EMBASSY": true,
	"CANCELLED":           true,
}

func transitionAllowed(from, to string) bool {
	for _, allowed := range allowedTransitions[from] {
		if allowed == to {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Params / Results
// ---------------------------------------------------------------------------

type TransitionStatusParams struct {
	ApplicationID string
	ToStatus      string
	Reason        string
	ActorUserID   string
}

type TransitionStatusResult struct {
	ApplicationID string
	FromStatus    string
	ToStatus      string
	Idempotent    bool
}

type BulkSubmitParams struct {
	DepartureID string
	JamaahIDs   []string
	ProviderID  string
}

type BulkSubmitResult struct {
	SubmittedCount int32
	ApplicationIDs []string
}

type GetApplicationsParams struct {
	DepartureID  string
	StatusFilter string
}

type HistoryEntry struct {
	FromStatus string
	ToStatus   string
	Reason     string
	CreatedAt  time.Time
}

type ApplicationRecord struct {
	ID          string
	JamaahID    string
	Status      string
	ProviderRef string
	IssuedDate  string
	History     []HistoryEntry
}

type GetApplicationsResult struct {
	Applications []ApplicationRecord
}

// ---------------------------------------------------------------------------
// TransitionStatus (BL-VISA-001)
// ---------------------------------------------------------------------------

var (
	ErrVisaNotFound      = errors.New("not_found")
	ErrInvalidStatus     = errors.New("invalid_status_value")
	ErrInvalidTransition = errors.New("invalid_transition")
)

func (s *Service) TransitionStatus(ctx context.Context, params *TransitionStatusParams) (*TransitionStatusResult, error) {
	const op = "service.TransitionStatus"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	if !knownStatuses[params.ToStatus] {
		return nil, ErrInvalidStatus
	}

	app, err := s.store.GetVisaApplication(ctx, params.ApplicationID)
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("get visa application")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, ErrVisaNotFound
	}

	// Idempotent — already at target status
	if app.Status == params.ToStatus {
		span.SetAttributes(attribute.Bool("idempotent", true))
		return &TransitionStatusResult{
			ApplicationID: app.ID,
			FromStatus:    app.Status,
			ToStatus:      params.ToStatus,
			Idempotent:    true,
		}, nil
	}

	if !transitionAllowed(app.Status, params.ToStatus) {
		return nil, ErrInvalidTransition
	}

	// Persist status update + history inside a tx
	fromStatus := app.Status
	if err := s.store.UpdateVisaStatus(ctx, params.ToStatus, params.ApplicationID); err != nil {
		logger.Error().Err(err).Str("op", op).Msg("update visa status")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	if err := s.store.InsertStatusHistory(ctx, sqlc.InsertStatusHistoryParams{
		ApplicationID: params.ApplicationID,
		FromStatus:    fromStatus,
		ToStatus:      params.ToStatus,
		Reason:        params.Reason,
		ActorUserID:   params.ActorUserID,
	}); err != nil {
		logger.Error().Err(err).Str("op", op).Msg("insert status history")
	}

	return &TransitionStatusResult{
		ApplicationID: params.ApplicationID,
		FromStatus:    fromStatus,
		ToStatus:      params.ToStatus,
		Idempotent:    false,
	}, nil
}

// ---------------------------------------------------------------------------
// BulkSubmit (BL-VISA-002)
// ---------------------------------------------------------------------------

var (
	ErrEmptyJamaahList      = errors.New("empty_jamaah_list")
	ErrInvalidProvider      = errors.New("invalid_provider")
	ErrNotAllReady          = errors.New("not_all_ready")
)

var validProviders = map[string]bool{"sajil": true, "mofa": true}

func (s *Service) BulkSubmit(ctx context.Context, params *BulkSubmitParams) (*BulkSubmitResult, error) {
	const op = "service.BulkSubmit"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	if len(params.JamaahIDs) == 0 {
		return nil, ErrEmptyJamaahList
	}
	if !validProviders[params.ProviderID] {
		return nil, ErrInvalidProvider
	}

	// Verify all jamaah have READY visa applications for this departure
	var appIDs []string
	for _, jamaahID := range params.JamaahIDs {
		app, err := s.store.GetReadyVisaForJamaahDeparture(ctx, jamaahID, params.DepartureID)
		if err != nil {
			// Not found or not READY
			return nil, ErrNotAllReady
		}
		appIDs = append(appIDs, app.ID)
	}

	// Atomic batch: update all to SUBMITTED + insert history rows
	for i, appID := range appIDs {
		jamaahID := params.JamaahIDs[i]
		if err := s.store.UpdateVisaStatusAndProvider(ctx, "SUBMITTED", params.ProviderID, appID); err != nil {
			logger.Error().Err(err).Str("op", op).Str("app_id", appID).Msg("bulk submit update status")
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		if err := s.store.InsertStatusHistory(ctx, sqlc.InsertStatusHistoryParams{
			ApplicationID: appID,
			FromStatus:    "READY",
			ToStatus:      "SUBMITTED",
			ActorUserID:   "",
			Reason:        "bulk submit via " + params.ProviderID,
		}); err != nil {
			logger.Error().Err(err).Str("op", op).Str("jamaah_id", jamaahID).Msg("insert bulk submit history")
		}
	}

	return &BulkSubmitResult{
		SubmittedCount: int32(len(appIDs)),
		ApplicationIDs: appIDs,
	}, nil
}

// ---------------------------------------------------------------------------
// GetApplications (BL-VISA-003)
// ---------------------------------------------------------------------------

var knownStatusValues = knownStatuses // re-use map

func (s *Service) GetApplications(ctx context.Context, params *GetApplicationsParams) (*GetApplicationsResult, error) {
	const op = "service.GetApplications"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	if params.StatusFilter != "" && !knownStatusValues[params.StatusFilter] {
		return nil, errors.New("invalid_status_filter")
	}

	apps, err := s.store.GetVisaApplicationsForDeparture(ctx, sqlc.GetVisaApplicationsForDepartureParams{
		DepartureID:  params.DepartureID,
		StatusFilter: params.StatusFilter,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("get visa applications")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var records []ApplicationRecord
	for _, app := range apps {
		history, err := s.store.GetStatusHistoryForApplication(ctx, app.ID)
		if err != nil {
			logger.Warn().Err(err).Str("op", op).Str("app_id", app.ID).Msg("get status history")
		}

		var entries []HistoryEntry
		for _, h := range history {
			fromStr := ""
			if h.FromStatus.Valid {
				fromStr = h.FromStatus.String
			}
			reasonStr := ""
			if h.Reason.Valid {
				reasonStr = h.Reason.String
			}
			ts := time.Time{}
			if h.CreatedAt.Valid {
				ts = h.CreatedAt.Time
			}
			entries = append(entries, HistoryEntry{
				FromStatus: fromStr,
				ToStatus:   h.ToStatus,
				Reason:     reasonStr,
				CreatedAt:  ts,
			})
		}

		providerRef := ""
		if app.ProviderRef.Valid {
			providerRef = app.ProviderRef.String
		}
		issuedDate := ""
		if app.IssuedDate.Valid {
			y, m, d := app.IssuedDate.Time.Date()
			issuedDate = time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		}

		records = append(records, ApplicationRecord{
			ID:          app.ID,
			JamaahID:    app.JamaahID,
			Status:      app.Status,
			ProviderRef: providerRef,
			IssuedDate:  issuedDate,
			History:     entries,
		})
	}

	return &GetApplicationsResult{Applications: records}, nil
}

