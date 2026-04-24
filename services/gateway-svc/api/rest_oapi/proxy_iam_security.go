// proxy_iam_security.go — gateway REST handlers for IAM security depth.
// BL-IAM-010: GetPasswordPolicy, SetPasswordPolicy
// BL-IAM-012: RecordLoginAnomaly
// BL-IAM-013: ListSessions, RevokeSession
// BL-IAM-015: UpsertCommTemplate, ListCommTemplates
// BL-IAM-017: TriggerBackup, GetBackupHistory
//
// Route topology (all bearer-protected):
//   GET  /v1/admin/security/password-policy      → GetPasswordPolicy
//   PUT  /v1/admin/security/password-policy      → SetPasswordPolicy
//   POST /v1/admin/security/anomalies            → RecordLoginAnomaly
//   GET  /v1/admin/security/sessions             → ListSessions
//   DELETE /v1/admin/security/sessions/:id       → RevokeSession
//   GET  /v1/admin/comm-templates                → ListCommTemplates
//   PUT  /v1/admin/comm-templates/:channel/:name → UpsertCommTemplate
//   GET  /v1/admin/backups                       → GetBackupHistory
//   POST /v1/admin/backups                       → TriggerBackup
package rest_oapi

import (
	"errors"

	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Request body types
// ---------------------------------------------------------------------------

type SetPasswordPolicyBody struct {
	MinLength      int32    `json:"min_length"`
	RequireUpper   bool     `json:"require_upper"`
	RequireDigit   bool     `json:"require_digit"`
	RequireSpecial bool     `json:"require_special"`
	RequireMfa     bool     `json:"require_mfa"`
	UpdatedBy      string   `json:"updated_by,omitempty"`
}

type RecordLoginAnomalyBody struct {
	UserID      string `json:"user_id"`
	IP          string `json:"ip,omitempty"`
	UserAgent   string `json:"user_agent,omitempty"`
	AnomalyKind string `json:"anomaly_kind"`
	Details     string `json:"details,omitempty"`
}

type RevokeSessionBody struct {
	RequestorID string `json:"requestor_id,omitempty"`
}

type UpsertCommTemplateBody struct {
	Subject   string   `json:"subject,omitempty"`
	Body      string   `json:"body"`
	Variables []string `json:"variables,omitempty"`
	UpdatedBy string   `json:"updated_by,omitempty"`
}

type TriggerBackupBody struct {
	TriggeredBy string `json:"triggered_by,omitempty"`
	Label       string `json:"label,omitempty"`
}

// ---------------------------------------------------------------------------
// BL-IAM-010: Password policy
// ---------------------------------------------------------------------------

func (s *Server) GetPasswordPolicy(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetPasswordPolicy"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetPasswordPolicy(ctx)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"min_length":      result.MinLength,
			"require_upper":   result.RequireUpper,
			"require_digit":   result.RequireDigit,
			"require_special": result.RequireSpecial,
			"require_mfa":     result.RequireMfa,
			"updated_at":      result.UpdatedAt,
		},
	})
}

func (s *Server) SetPasswordPolicy(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.SetPasswordPolicy"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body SetPasswordPolicyBody
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.SetPasswordPolicy(ctx, &iam_grpc_adapter.SetPasswordPolicyParams{
		MinLength:      body.MinLength,
		RequireUpper:   body.RequireUpper,
		RequireDigit:   body.RequireDigit,
		RequireSpecial: body.RequireSpecial,
		RequireMfa:     body.RequireMfa,
		UpdatedBy:      body.UpdatedBy,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"min_length":      result.MinLength,
			"require_upper":   result.RequireUpper,
			"require_digit":   result.RequireDigit,
			"require_special": result.RequireSpecial,
			"require_mfa":     result.RequireMfa,
			"updated_at":      result.UpdatedAt,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-IAM-012: Login anomaly
// ---------------------------------------------------------------------------

func (s *Server) RecordLoginAnomaly(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordLoginAnomaly"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body RecordLoginAnomalyBody
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.RecordLoginAnomaly(ctx, &iam_grpc_adapter.RecordLoginAnomalyParams{
		UserID:      body.UserID,
		IP:          body.IP,
		UserAgent:   body.UserAgent,
		AnomalyKind: body.AnomalyKind,
		Details:     body.Details,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"alert_id":   result.AlertID,
			"created_at": result.CreatedAt,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-IAM-013: Sessions
// ---------------------------------------------------------------------------

func (s *Server) ListSessions(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListSessions"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	userID := c.Query("user_id")
	includeAll := c.QueryBool("include_all", false)
	span.SetAttributes(attribute.String("user_id", userID))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("user_id", userID).Msg("")

	result, err := s.svc.ListSessions(ctx, &iam_grpc_adapter.ListSessionsParams{
		UserID:     userID,
		IncludeAll: includeAll,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	type sessionJSON struct {
		SessionID string `json:"session_id"`
		UserAgent string `json:"user_agent,omitempty"`
		IP        string `json:"ip,omitempty"`
		IssuedAt  string `json:"issued_at"`
		ExpiresAt string `json:"expires_at"`
		RevokedAt string `json:"revoked_at,omitempty"`
		IsActive  bool   `json:"is_active"`
	}
	data := make([]sessionJSON, 0, len(result.Sessions))
	for _, s := range result.Sessions {
		data = append(data, sessionJSON{
			SessionID: s.SessionID,
			UserAgent: s.UserAgent,
			IP:        s.IP,
			IssuedAt:  s.IssuedAt,
			ExpiresAt: s.ExpiresAt,
			RevokedAt: s.RevokedAt,
			IsActive:  s.IsActive,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": data})
}

func (s *Server) RevokeSession(c *fiber.Ctx, sessionID string) error {
	const op = "rest_oapi.Server.RevokeSession"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("session_id", sessionID))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("session_id", sessionID).Msg("")

	var body RevokeSessionBody
	_ = c.BodyParser(&body) // optional body

	result, err := s.svc.RevokeSession(ctx, &iam_grpc_adapter.RevokeSessionParams{
		SessionID:   sessionID,
		RequestorID: body.RequestorID,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"session_id": result.SessionID,
			"revoked_at": result.RevokedAt,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-IAM-015: Comm templates
// ---------------------------------------------------------------------------

func (s *Server) ListCommTemplates(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListCommTemplates"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	channel := c.Query("channel")

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.ListCommTemplates(ctx, &iam_grpc_adapter.ListCommTemplatesParams{
		Channel: channel,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	type templateJSON struct {
		Key       string   `json:"key"`
		Channel   string   `json:"channel"`
		Name      string   `json:"name"`
		Subject   string   `json:"subject,omitempty"`
		Body      string   `json:"body"`
		Variables []string `json:"variables,omitempty"`
		UpdatedAt string   `json:"updated_at"`
	}
	data := make([]templateJSON, 0, len(result.Templates))
	for _, t := range result.Templates {
		data = append(data, templateJSON{
			Key:       t.Key,
			Channel:   t.Channel,
			Name:      t.Name,
			Subject:   t.Subject,
			Body:      t.Body,
			Variables: t.Variables,
			UpdatedAt: t.UpdatedAt,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": data})
}

func (s *Server) UpsertCommTemplate(c *fiber.Ctx, channel, name string) error {
	const op = "rest_oapi.Server.UpsertCommTemplate"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(
		attribute.String("channel", channel),
		attribute.String("name", name),
	)

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("channel", channel).Str("name", name).Msg("")

	var body UpsertCommTemplateBody
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.UpsertCommTemplate(ctx, &iam_grpc_adapter.UpsertCommTemplateParams{
		Channel:   channel,
		Name:      name,
		Subject:   body.Subject,
		Body:      body.Body,
		Variables: body.Variables,
		UpdatedBy: body.UpdatedBy,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"key":        result.Key,
			"updated_at": result.UpdatedAt,
		},
	})
}

// ---------------------------------------------------------------------------
// BL-IAM-017: Backup
// ---------------------------------------------------------------------------

func (s *Server) GetBackupHistory(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetBackupHistory"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	limit := int32(c.QueryInt("limit", 20))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetBackupHistory(ctx, limit)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	type backupJSON struct {
		BackupID    string `json:"backup_id"`
		Status      string `json:"status"`
		TriggeredBy string `json:"triggered_by,omitempty"`
		ScheduledAt string `json:"scheduled_at"`
	}
	data := make([]backupJSON, 0, len(result.Backups))
	for _, b := range result.Backups {
		data = append(data, backupJSON{
			BackupID:    b.BackupID,
			Status:      b.Status,
			TriggeredBy: b.TriggeredBy,
			ScheduledAt: b.ScheduledAt,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": data})
}

func (s *Server) TriggerBackup(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.TriggerBackup"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body TriggerBackupBody
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.TriggerBackup(ctx, &iam_grpc_adapter.TriggerBackupParams{
		TriggeredBy: body.TriggeredBy,
		Label:       body.Label,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"backup_id":    result.BackupID,
			"status":       result.Status,
			"scheduled_at": result.ScheduledAt,
		},
	})
}
