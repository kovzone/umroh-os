// security_adapter.go — gateway-svc adapter methods for iam-svc security RPCs.
// BL-IAM-010: GetPasswordPolicy, SetPasswordPolicy
// BL-IAM-012: RecordLoginAnomaly
// BL-IAM-013: ListSessions, RevokeSession
// BL-IAM-015: UpsertCommTemplate, ListCommTemplates
// BL-IAM-017: TriggerBackup, GetBackupHistory

package iam_grpc_adapter

import (
	"context"
	"fmt"

	"gateway-svc/adapter/iam_grpc_adapter/pb"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// Result types — returned to gateway service layer
// ---------------------------------------------------------------------------

type PasswordPolicyResult struct {
	MinLength      int32
	RequireUpper   bool
	RequireDigit   bool
	RequireSpecial bool
	RequireMfa     bool
	UpdatedAt      string
}

type SetPasswordPolicyParams struct {
	MinLength      int32
	RequireUpper   bool
	RequireDigit   bool
	RequireSpecial bool
	RequireMfa     bool
	UpdatedBy      string
}

type RecordLoginAnomalyParams struct {
	UserID      string
	IP          string
	UserAgent   string
	AnomalyKind string
	Details     string
}

type RecordLoginAnomalyResult struct {
	AlertID   string
	CreatedAt string
}

type ListSessionsParams struct {
	UserID     string
	IncludeAll bool
}

type SessionEntryResult struct {
	SessionID string
	UserAgent string
	IP        string
	IssuedAt  string
	ExpiresAt string
	RevokedAt string
	IsActive  bool
}

type ListSessionsResult struct {
	Sessions []*SessionEntryResult
}

type RevokeSessionParams struct {
	SessionID   string
	RequestorID string
}

type RevokeSessionResult struct {
	SessionID string
	RevokedAt string
}

type UpsertCommTemplateParams struct {
	Channel   string
	Name      string
	Subject   string
	Body      string
	Variables []string
	UpdatedBy string
}

type UpsertCommTemplateResult struct {
	Key       string
	UpdatedAt string
}

type ListCommTemplatesParams struct {
	Channel string
}

type CommTemplateResult struct {
	Key       string
	Channel   string
	Name      string
	Subject   string
	Body      string
	Variables []string
	UpdatedAt string
}

type ListCommTemplatesResult struct {
	Templates []*CommTemplateResult
}

type TriggerBackupParams struct {
	TriggeredBy string
	Label       string
}

type TriggerBackupResult struct {
	BackupID    string
	Status      string
	ScheduledAt string
}

type BackupEntryResult struct {
	BackupID    string
	Status      string
	TriggeredBy string
	ScheduledAt string
}

type GetBackupHistoryResult struct {
	Backups []*BackupEntryResult
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func mapIamSecurityError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return apperrors.ErrInternal
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return fmt.Errorf("%w: %s", apperrors.ErrNotFound, st.Message())
	case grpcCodes.InvalidArgument:
		return fmt.Errorf("%w: %s", apperrors.ErrValidation, st.Message())
	case grpcCodes.PermissionDenied:
		return fmt.Errorf("%w: %s", apperrors.ErrForbidden, st.Message())
	case grpcCodes.Unauthenticated:
		return fmt.Errorf("%w: %s", apperrors.ErrUnauthorized, st.Message())
	default:
		return apperrors.ErrInternal
	}
}

// ---------------------------------------------------------------------------
// BL-IAM-010: Password policy
// ---------------------------------------------------------------------------

func (a *Adapter) GetPasswordPolicy(ctx context.Context) (*PasswordPolicyResult, error) {
	const op = "iam_grpc_adapter.Adapter.GetPasswordPolicy"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Msg("")

	resp, err := a.securityClient.GetPasswordPolicy(ctx, &pb.GwGetPasswordPolicyRequest{})
	if err != nil {
		wrapped := mapIamSecurityError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.GetPasswordPolicy failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &PasswordPolicyResult{
		MinLength:      resp.GetMinLength(),
		RequireUpper:   resp.GetRequireUpper(),
		RequireDigit:   resp.GetRequireDigit(),
		RequireSpecial: resp.GetRequireSpecial(),
		RequireMfa:     resp.GetRequireMfa(),
		UpdatedAt:      resp.GetUpdatedAt(),
	}, nil
}

func (a *Adapter) SetPasswordPolicy(ctx context.Context, params *SetPasswordPolicyParams) (*PasswordPolicyResult, error) {
	const op = "iam_grpc_adapter.Adapter.SetPasswordPolicy"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Msg("")

	resp, err := a.securityClient.SetPasswordPolicy(ctx, &pb.GwSetPasswordPolicyRequest{
		MinLength:      params.MinLength,
		RequireUpper:   params.RequireUpper,
		RequireDigit:   params.RequireDigit,
		RequireSpecial: params.RequireSpecial,
		RequireMfa:     params.RequireMfa,
		UpdatedBy:      params.UpdatedBy,
	})
	if err != nil {
		wrapped := mapIamSecurityError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.SetPasswordPolicy failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &PasswordPolicyResult{
		MinLength:      resp.GetMinLength(),
		RequireUpper:   resp.GetRequireUpper(),
		RequireDigit:   resp.GetRequireDigit(),
		RequireSpecial: resp.GetRequireSpecial(),
		RequireMfa:     resp.GetRequireMfa(),
		UpdatedAt:      resp.GetUpdatedAt(),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-IAM-012: Anomaly
// ---------------------------------------------------------------------------

func (a *Adapter) RecordLoginAnomaly(ctx context.Context, params *RecordLoginAnomalyParams) (*RecordLoginAnomalyResult, error) {
	const op = "iam_grpc_adapter.Adapter.RecordLoginAnomaly"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("user_id", params.UserID).Msg("")

	resp, err := a.securityClient.RecordLoginAnomaly(ctx, &pb.GwRecordLoginAnomalyRequest{
		UserId:      params.UserID,
		Ip:          params.IP,
		UserAgent:   params.UserAgent,
		AnomalyKind: params.AnomalyKind,
		Details:     params.Details,
	})
	if err != nil {
		wrapped := mapIamSecurityError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.RecordLoginAnomaly failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &RecordLoginAnomalyResult{
		AlertID:   resp.GetAlertId(),
		CreatedAt: resp.GetCreatedAt(),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-IAM-013: Sessions
// ---------------------------------------------------------------------------

func (a *Adapter) ListSessions(ctx context.Context, params *ListSessionsParams) (*ListSessionsResult, error) {
	const op = "iam_grpc_adapter.Adapter.ListSessions"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("user_id", params.UserID).Msg("")

	resp, err := a.securityClient.ListSessions(ctx, &pb.GwListSessionsRequest{
		UserId:     params.UserID,
		IncludeAll: params.IncludeAll,
	})
	if err != nil {
		wrapped := mapIamSecurityError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.ListSessions failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	sessions := make([]*SessionEntryResult, 0, len(resp.GetSessions()))
	for _, s := range resp.GetSessions() {
		sessions = append(sessions, &SessionEntryResult{
			SessionID: s.GetSessionId(),
			UserAgent: s.GetUserAgent(),
			IP:        s.GetIp(),
			IssuedAt:  s.GetIssuedAt(),
			ExpiresAt: s.GetExpiresAt(),
			RevokedAt: s.GetRevokedAt(),
			IsActive:  s.GetIsActive(),
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &ListSessionsResult{Sessions: sessions}, nil
}

func (a *Adapter) RevokeSession(ctx context.Context, params *RevokeSessionParams) (*RevokeSessionResult, error) {
	const op = "iam_grpc_adapter.Adapter.RevokeSession"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("session_id", params.SessionID).Msg("")

	resp, err := a.securityClient.RevokeSession(ctx, &pb.GwRevokeSessionRequest{
		SessionId:   params.SessionID,
		RequestorId: params.RequestorID,
	})
	if err != nil {
		wrapped := mapIamSecurityError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.RevokeSession failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &RevokeSessionResult{
		SessionID: resp.GetSessionId(),
		RevokedAt: resp.GetRevokedAt(),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-IAM-015: Comm templates
// ---------------------------------------------------------------------------

func (a *Adapter) UpsertCommTemplate(ctx context.Context, params *UpsertCommTemplateParams) (*UpsertCommTemplateResult, error) {
	const op = "iam_grpc_adapter.Adapter.UpsertCommTemplate"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("channel", params.Channel).Str("name", params.Name).Msg("")

	resp, err := a.securityClient.UpsertCommTemplate(ctx, &pb.GwUpsertCommTemplateRequest{
		Channel:   params.Channel,
		Name:      params.Name,
		Subject:   params.Subject,
		Body:      params.Body,
		Variables: params.Variables,
		UpdatedBy: params.UpdatedBy,
	})
	if err != nil {
		wrapped := mapIamSecurityError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.UpsertCommTemplate failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &UpsertCommTemplateResult{
		Key:       resp.GetKey(),
		UpdatedAt: resp.GetUpdatedAt(),
	}, nil
}

func (a *Adapter) ListCommTemplates(ctx context.Context, params *ListCommTemplatesParams) (*ListCommTemplatesResult, error) {
	const op = "iam_grpc_adapter.Adapter.ListCommTemplates"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Msg("")

	resp, err := a.securityClient.ListCommTemplates(ctx, &pb.GwListCommTemplatesRequest{
		Channel: params.Channel,
	})
	if err != nil {
		wrapped := mapIamSecurityError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.ListCommTemplates failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	templates := make([]*CommTemplateResult, 0, len(resp.GetTemplates()))
	for _, t := range resp.GetTemplates() {
		templates = append(templates, &CommTemplateResult{
			Key:       t.GetKey(),
			Channel:   t.GetChannel(),
			Name:      t.GetName(),
			Subject:   t.GetSubject(),
			Body:      t.GetBody(),
			Variables: t.GetVariables(),
			UpdatedAt: t.GetUpdatedAt(),
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &ListCommTemplatesResult{Templates: templates}, nil
}

// ---------------------------------------------------------------------------
// BL-IAM-017: Backup
// ---------------------------------------------------------------------------

func (a *Adapter) TriggerBackup(ctx context.Context, params *TriggerBackupParams) (*TriggerBackupResult, error) {
	const op = "iam_grpc_adapter.Adapter.TriggerBackup"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("triggered_by", params.TriggeredBy).Msg("")

	resp, err := a.securityClient.TriggerBackup(ctx, &pb.GwTriggerBackupRequest{
		TriggeredBy: params.TriggeredBy,
		Label:       params.Label,
	})
	if err != nil {
		wrapped := mapIamSecurityError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.TriggerBackup failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &TriggerBackupResult{
		BackupID:    resp.GetBackupId(),
		Status:      resp.GetStatus(),
		ScheduledAt: resp.GetScheduledAt(),
	}, nil
}

func (a *Adapter) GetBackupHistory(ctx context.Context, limit int32) (*GetBackupHistoryResult, error) {
	const op = "iam_grpc_adapter.Adapter.GetBackupHistory"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Msg("")

	resp, err := a.securityClient.GetBackupHistory(ctx, &pb.GwGetBackupHistoryRequest{Limit: limit})
	if err != nil {
		wrapped := mapIamSecurityError(err)
		logger.Warn().Err(wrapped).Msg("iam-svc.GetBackupHistory failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	backups := make([]*BackupEntryResult, 0, len(resp.GetBackups()))
	for _, b := range resp.GetBackups() {
		backups = append(backups, &BackupEntryResult{
			BackupID:    b.GetBackupId(),
			Status:      b.GetStatus(),
			TriggeredBy: b.GetTriggeredBy(),
			ScheduledAt: b.GetScheduledAt(),
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &GetBackupHistoryResult{Backups: backups}, nil
}

