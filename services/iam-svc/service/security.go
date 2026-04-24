// security.go — IAM Phase 6 security depth features.
// BL-IAM-010: Password policy + optional MFA enforcement
// BL-IAM-012: Login/action anomaly alerts
// BL-IAM-013: Session history + revoke
// BL-IAM-015: Communication template CRUD (WA/email)
// BL-IAM-017: DB backup schedule/status (procedure wrapper)

package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"iam-svc/store/postgres_store/sqlc"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// BL-IAM-010: Password policy
// ---------------------------------------------------------------------------

// PasswordPolicy defines complexity/length rules for user passwords.
// These are checked at create-user and password-reset time.
type PasswordPolicy struct {
	MinLength      int
	RequireUpper   bool
	RequireDigit   bool
	RequireSpecial bool
}

// DefaultPasswordPolicy is the production-grade password policy for UmrohOS.
var DefaultPasswordPolicy = PasswordPolicy{
	MinLength:      8,
	RequireUpper:   true,
	RequireDigit:   true,
	RequireSpecial: true,
}

// ValidatePassword checks the given plaintext password against the policy.
// Returns a joined validation error listing all violations.
func ValidatePassword(password string, policy PasswordPolicy) error {
	var violations []string

	if len(password) < policy.MinLength {
		violations = append(violations, fmt.Sprintf("at least %d characters required", policy.MinLength))
	}
	if policy.RequireUpper {
		hasUpper := false
		for _, r := range password {
			if unicode.IsUpper(r) {
				hasUpper = true
				break
			}
		}
		if !hasUpper {
			violations = append(violations, "at least one uppercase letter required")
		}
	}
	if policy.RequireDigit {
		hasDigit := false
		for _, r := range password {
			if unicode.IsDigit(r) {
				hasDigit = true
				break
			}
		}
		if !hasDigit {
			violations = append(violations, "at least one digit required")
		}
	}
	if policy.RequireSpecial {
		specialRe := regexp.MustCompile(`[^a-zA-Z0-9]`)
		if !specialRe.MatchString(password) {
			violations = append(violations, "at least one special character required")
		}
	}

	if len(violations) > 0 {
		return errors.Join(apperrors.ErrValidation,
			fmt.Errorf("password policy violations: %s", strings.Join(violations, "; ")))
	}
	return nil
}

// SetPasswordPolicyParams are the request fields for updating the global
// password policy stored in iam.global_config.
type SetPasswordPolicyParams struct {
	MinLength      int
	RequireUpper   bool
	RequireDigit   bool
	RequireSpecial bool
	RequireMFA     bool // BL-IAM-010: optional MFA enforcement flag
	UpdatedBy      string
}

type GetPasswordPolicyResult struct {
	MinLength      int
	RequireUpper   bool
	RequireDigit   bool
	RequireSpecial bool
	RequireMFA     bool
	UpdatedAt      string // RFC3339
}

// GetPasswordPolicy reads the current password policy from global config.
func (s *Service) GetPasswordPolicy(ctx context.Context) (*GetPasswordPolicyResult, error) {
	const op = "service.Service.GetPasswordPolicy"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	keys := []string{
		"iam.policy.password.min_length",
		"iam.policy.password.require_upper",
		"iam.policy.password.require_digit",
		"iam.policy.password.require_special",
		"iam.policy.mfa.enforce",
	}

	cfgResult, err := s.GetGlobalConfig(ctx, &GetGlobalConfigParams{Keys: keys})
	if err != nil {
		return nil, err
	}

	res := &GetPasswordPolicyResult{
		MinLength:      DefaultPasswordPolicy.MinLength,
		RequireUpper:   DefaultPasswordPolicy.RequireUpper,
		RequireDigit:   DefaultPasswordPolicy.RequireDigit,
		RequireSpecial: DefaultPasswordPolicy.RequireSpecial,
		RequireMFA:     false,
	}
	for _, cfg := range cfgResult.Configs {
		switch cfg.Key {
		case "iam.policy.password.min_length":
			var n int
			if _, err2 := fmt.Sscanf(cfg.Value, "%d", &n); err2 == nil && n > 0 {
				res.MinLength = n
			}
		case "iam.policy.password.require_upper":
			res.RequireUpper = cfg.Value == "true"
		case "iam.policy.password.require_digit":
			res.RequireDigit = cfg.Value == "true"
		case "iam.policy.password.require_special":
			res.RequireSpecial = cfg.Value == "true"
		case "iam.policy.mfa.enforce":
			res.RequireMFA = cfg.Value == "true"
		}
		if cfg.UpdatedAt != "" {
			res.UpdatedAt = cfg.UpdatedAt
		}
	}
	return res, nil
}

// SetPasswordPolicy persists the policy to global_config (5 keys, upsert each).
func (s *Service) SetPasswordPolicy(ctx context.Context, params *SetPasswordPolicyParams) (*GetPasswordPolicyResult, error) {
	const op = "service.Service.SetPasswordPolicy"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	boolStr := func(b bool) string {
		if b {
			return "true"
		}
		return "false"
	}

	entries := []struct {
		key   string
		value string
		desc  string
	}{
		{"iam.policy.password.min_length", fmt.Sprintf("%d", params.MinLength), "Minimum password length"},
		{"iam.policy.password.require_upper", boolStr(params.RequireUpper), "Require at least one uppercase letter"},
		{"iam.policy.password.require_digit", boolStr(params.RequireDigit), "Require at least one digit"},
		{"iam.policy.password.require_special", boolStr(params.RequireSpecial), "Require at least one special character"},
		{"iam.policy.mfa.enforce", boolStr(params.RequireMFA), "Enforce MFA for all users"},
	}

	for _, e := range entries {
		if _, err := s.SetGlobalConfig(ctx, &SetGlobalConfigParams{
			Key:       e.key,
			Value:     e.value,
			UpdatedBy: params.UpdatedBy,
		}); err != nil {
			return nil, fmt.Errorf("%s: upsert %s: %w", op, e.key, err)
		}
	}

	return s.GetPasswordPolicy(ctx)
}

// ---------------------------------------------------------------------------
// BL-IAM-012: Login/action anomaly alerts
// ---------------------------------------------------------------------------

type RecordLoginAnomalyParams struct {
	UserID    string
	IP        string
	UserAgent string
	AnomalyKind string // "new_ip" | "new_device" | "rapid_login" | "off_hours"
	Details   string  // human-readable explanation
}

type RecordLoginAnomalyResult struct {
	AlertID   string
	CreatedAt string // RFC3339
}

// RecordLoginAnomaly writes an anomaly marker to the audit log so security
// dashboards can surface it. It is intentionally lightweight — the heavy
// detection logic lives in a dedicated SIEM (future) or periodic cron.
func (s *Service) RecordLoginAnomaly(ctx context.Context, params *RecordLoginAnomalyParams) (*RecordLoginAnomalyResult, error) {
	const op = "service.Service.RecordLoginAnomaly"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("user_id", params.UserID),
		attribute.String("anomaly_kind", params.AnomalyKind),
	)
	logger := logging.LogWithTrace(ctx, s.logger)

	if params.UserID == "" || params.AnomalyKind == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id and anomaly_kind are required"))
	}

	// Write to audit log with resource="iam.security" action="anomaly.{kind}"
	auditResult, err := s.RecordAudit(ctx, &RecordAuditParams{
		ActorUserID: params.UserID,
		Resource:    "iam.security",
		Action:      "anomaly." + params.AnomalyKind,
		ResourceID:  params.IP,
		IP:          params.IP,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Warn().
		Str("user_id", params.UserID).
		Str("anomaly_kind", params.AnomalyKind).
		Str("ip", params.IP).
		Str("user_agent", params.UserAgent).
		Str("details", params.Details).
		Msg("login anomaly detected")

	span.SetStatus(otelCodes.Ok, "success")
	return &RecordLoginAnomalyResult{
		AlertID:   auditResult.AuditLogID,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-IAM-013: Session history + revoke
// ---------------------------------------------------------------------------

type ListSessionsParams struct {
	UserID       string // target user; required
	IncludeAll   bool   // if false, only active (non-expired, non-revoked) sessions
}

type SessionEntry struct {
	SessionID string
	UserAgent string
	IP        string
	IssuedAt  string // RFC3339
	ExpiresAt string // RFC3339
	RevokedAt string // RFC3339 or ""
	IsActive  bool
}

type ListSessionsResult struct {
	Sessions []SessionEntry
}

func (s *Service) ListSessions(ctx context.Context, params *ListSessionsParams) (*ListSessionsResult, error) {
	const op = "service.Service.ListSessions"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("user_id", params.UserID))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("user_id", params.UserID).Msg("")

	if params.UserID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("user_id is required"))
	}

	uid, err := uuid.Parse(params.UserID)
	if err != nil {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("invalid user_id: %w", err))
	}

	rows, err := s.store.ListActiveSessionsForUser(ctx, pgtype.UUID{Bytes: uid, Valid: true})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	sessions := make([]SessionEntry, 0, len(rows))
	now := time.Now().UTC()
	for _, r := range rows {
		revokedAt := ""
		if r.RevokedAt.Valid {
			revokedAt = r.RevokedAt.Time.UTC().Format(time.RFC3339)
		}
		ipStr := ""
		if r.Ip != nil {
			ipStr = r.Ip.String()
		}
		active := !r.RevokedAt.Valid && r.ExpiresAt.Time.After(now)
		sessions = append(sessions, SessionEntry{
			SessionID: uuidToString(r.ID),
			UserAgent: r.UserAgent,
			IP:        ipStr,
			IssuedAt:  r.IssuedAt.Time.UTC().Format(time.RFC3339),
			ExpiresAt: r.ExpiresAt.Time.UTC().Format(time.RFC3339),
			RevokedAt: revokedAt,
			IsActive:  active,
		})
	}

	span.SetStatus(otelCodes.Ok, "success")
	return &ListSessionsResult{Sessions: sessions}, nil
}

type RevokeSessionParams struct {
	SessionID   string
	RequestorID string // admin or owner performing the revoke
}

type RevokeSessionResult struct {
	SessionID string
	RevokedAt string // RFC3339
}

func (s *Service) RevokeSession(ctx context.Context, params *RevokeSessionParams) (*RevokeSessionResult, error) {
	const op = "service.Service.RevokeSession"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("session_id", params.SessionID))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("session_id", params.SessionID).Msg("")

	if params.SessionID == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("session_id is required"))
	}

	sid, err := uuid.Parse(params.SessionID)
	if err != nil {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("invalid session_id: %w", err))
	}

	if err := s.store.RevokeSession(ctx, pgtype.UUID{Bytes: sid, Valid: true}); err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Emit audit log
	_, _ = s.RecordAudit(ctx, &RecordAuditParams{
		ActorUserID: params.RequestorID,
		Resource:    "iam.sessions",
		Action:      "revoke",
		ResourceID:  params.SessionID,
	})

	revokedAt := time.Now().UTC().Format(time.RFC3339)
	span.SetStatus(otelCodes.Ok, "success")
	return &RevokeSessionResult{SessionID: params.SessionID, RevokedAt: revokedAt}, nil
}

// ---------------------------------------------------------------------------
// BL-IAM-015: Communication template CRUD
// ---------------------------------------------------------------------------

// CommTemplate stores WA/email message templates that can reference
// {{variable}} placeholders. Templates are stored in global_config under
// the key prefix "comm.template.<channel>.<name>".

type CommTemplateChannel string

const (
	CommChannelWhatsApp CommTemplateChannel = "whatsapp"
	CommChannelEmail    CommTemplateChannel = "email"
	CommChannelSMS      CommTemplateChannel = "sms"
)

type CommTemplate struct {
	Key       string              // e.g. "comm.template.whatsapp.booking_confirmed"
	Channel   CommTemplateChannel
	Name      string
	Subject   string // email subject; "" for WA/SMS
	Body      string // template body with {{var}} placeholders
	Variables []string
	UpdatedAt string // RFC3339
	UpdatedBy string
}

type UpsertCommTemplateParams struct {
	Channel   CommTemplateChannel
	Name      string
	Subject   string
	Body      string
	Variables []string
	UpdatedBy string
}

type UpsertCommTemplateResult struct {
	Key       string
	UpdatedAt string // RFC3339
}

type ListCommTemplatesParams struct {
	Channel CommTemplateChannel // "" = all channels
}

type ListCommTemplatesResult struct {
	Templates []CommTemplate
}

type DeleteCommTemplateParams struct {
	Channel CommTemplateChannel
	Name    string
	DeletedBy string
}

// UpsertCommTemplate stores a communication template in global_config.
func (s *Service) UpsertCommTemplate(ctx context.Context, params *UpsertCommTemplateParams) (*UpsertCommTemplateResult, error) {
	const op = "service.Service.UpsertCommTemplate"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)

	if params.Channel == "" || params.Name == "" || params.Body == "" {
		return nil, errors.Join(apperrors.ErrValidation, fmt.Errorf("channel, name, and body are required"))
	}

	key := fmt.Sprintf("comm.template.%s.%s", params.Channel, params.Name)
	// Encode as JSON-ish value: subject|||body|||var1,var2
	vars := strings.Join(params.Variables, ",")
	value := fmt.Sprintf("%s|||%s|||%s", params.Subject, params.Body, vars)

	res, err := s.SetGlobalConfig(ctx, &SetGlobalConfigParams{
		Key:         key,
		Value:       value,
		Description: fmt.Sprintf("Communication template [%s/%s]", params.Channel, params.Name),
		UpdatedBy:   params.UpdatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info().Str("op", op).Str("key", key).Str("channel", string(params.Channel)).Msg("upserted comm template")
	span.SetStatus(otelCodes.Ok, "success")
	return &UpsertCommTemplateResult{Key: key, UpdatedAt: res.UpdatedAt}, nil
}

// ListCommTemplates returns all templates, optionally filtered by channel.
func (s *Service) ListCommTemplates(ctx context.Context, params *ListCommTemplatesParams) (*ListCommTemplatesResult, error) {
	const op = "service.Service.ListCommTemplates"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("channel", string(params.Channel)).Msg("")

	allCfg, err := s.GetGlobalConfig(ctx, &GetGlobalConfigParams{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	prefix := "comm.template."
	if params.Channel != "" {
		prefix = fmt.Sprintf("comm.template.%s.", params.Channel)
	}

	templates := []CommTemplate{}
	for _, cfg := range allCfg.Configs {
		if !strings.HasPrefix(cfg.Key, prefix) {
			continue
		}
		// Parse key: comm.template.<channel>.<name>
		parts := strings.SplitN(cfg.Key, ".", 4)
		if len(parts) < 4 {
			continue
		}
		channel := CommTemplateChannel(parts[2])
		name := parts[3]

		// Decode value: subject|||body|||vars
		valueParts := strings.SplitN(cfg.Value, "|||", 3)
		subject, body, varsStr := "", cfg.Value, ""
		if len(valueParts) >= 2 {
			subject = valueParts[0]
			body = valueParts[1]
		}
		if len(valueParts) >= 3 {
			varsStr = valueParts[2]
		}
		var vars []string
		if varsStr != "" {
			vars = strings.Split(varsStr, ",")
		}

		templates = append(templates, CommTemplate{
			Key:       cfg.Key,
			Channel:   channel,
			Name:      name,
			Subject:   subject,
			Body:      body,
			Variables: vars,
			UpdatedAt: cfg.UpdatedAt,
		})
	}

	span.SetStatus(otelCodes.Ok, "success")
	return &ListCommTemplatesResult{Templates: templates}, nil
}

// ---------------------------------------------------------------------------
// BL-IAM-017: DB backup schedule/status
// ---------------------------------------------------------------------------

type BackupStatus string

const (
	BackupStatusPending  BackupStatus = "pending"
	BackupStatusRunning  BackupStatus = "running"
	BackupStatusSuccess  BackupStatus = "success"
	BackupStatusFailed   BackupStatus = "failed"
)

type TriggerBackupParams struct {
	TriggeredBy string // user ID or "scheduler"
	Label       string // optional human-readable label
}

type TriggerBackupResult struct {
	BackupID   string       // stored as audit resource_id
	Status     BackupStatus
	ScheduledAt string // RFC3339
}

type GetBackupStatusResult struct {
	Backups []BackupEntry
}

type BackupEntry struct {
	BackupID    string
	Label       string
	Status      BackupStatus
	TriggeredBy string
	ScheduledAt string // RFC3339
	Details     string
}

// TriggerBackup records a backup initiation event to the audit log.
// The actual pg_dump is executed by the DevSecOps-managed backup-agent
// (sidecar or cron), which later calls SetBackupStatus. This RPC is
// the "schedule" half described in BL-IAM-017.
func (s *Service) TriggerBackup(ctx context.Context, params *TriggerBackupParams) (*TriggerBackupResult, error) {
	const op = "service.Service.TriggerBackup"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("triggered_by", params.TriggeredBy))
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("triggered_by", params.TriggeredBy).Msg("backup triggered")

	backupID := uuid.NewString()
	label := params.Label
	if label == "" {
		label = fmt.Sprintf("backup-%s", time.Now().UTC().Format("20060102-150405"))
	}

	// Persist the initiation event as an audit row
	_, err := s.RecordAudit(ctx, &RecordAuditParams{
		ActorUserID: params.TriggeredBy,
		Resource:    "iam.backup",
		Action:      "trigger",
		ResourceID:  backupID,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	scheduledAt := time.Now().UTC().Format(time.RFC3339)
	span.SetStatus(otelCodes.Ok, "success")
	return &TriggerBackupResult{
		BackupID:    backupID,
		Status:      BackupStatusPending,
		ScheduledAt: scheduledAt,
	}, nil
}

// GetBackupHistory returns the last N backup events from the audit log.
func (s *Service) GetBackupHistory(ctx context.Context, limit int32) (*GetBackupStatusResult, error) {
	const op = "service.Service.GetBackupHistory"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	if limit <= 0 {
		limit = 20
	}

	logResult, err := s.SearchActivityLog(ctx, &SearchActivityLogParams{
		Resource: "iam.backup",
		Limit:    limit,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	backups := make([]BackupEntry, 0, len(logResult.Logs))
	for _, entry := range logResult.Logs {
		status := BackupStatusPending
		if strings.Contains(entry.Action, "success") {
			status = BackupStatusSuccess
		} else if strings.Contains(entry.Action, "failed") {
			status = BackupStatusFailed
		} else if strings.Contains(entry.Action, "running") {
			status = BackupStatusRunning
		}
		backups = append(backups, BackupEntry{
			BackupID:    entry.ResourceID,
			Status:      status,
			TriggeredBy: entry.UserID,
			ScheduledAt: entry.CreatedAt,
		})
	}

	span.SetStatus(otelCodes.Ok, "success")
	return &GetBackupStatusResult{Backups: backups}, nil
}

// auditDetailsForSession formats consistent audit details for session events.
func auditDetailsForSession(userAgent, ip string) string {
	return fmt.Sprintf("ua=%s ip=%s", userAgent, ip)
}

// Ensure sqlc import is used (IamSession struct referenced via store queries).
var _ sqlc.IamSession
