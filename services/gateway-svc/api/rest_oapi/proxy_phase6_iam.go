// proxy_phase6_iam.go — gateway REST handlers for IAM Phase 6 security features
// (BL-IAM-007/011/014/016).
//
// Route topology (all bearer-protected):
//   PUT    /v1/admin/users/:id/data-scope  → SetDataScope
//   POST   /v1/admin/api-keys             → CreateAPIKey
//   DELETE /v1/admin/api-keys/:id         → RevokeAPIKey
//   GET    /v1/admin/config               → GetGlobalConfig
//   PUT    /v1/admin/config/:key          → SetGlobalConfig
//   GET    /v1/admin/activity-log         → SearchActivityLog
//
// Per ADR-0009: gateway is the single REST entry-point; iam-svc is pure gRPC.
package rest_oapi

import (
	"errors"
	"strings"

	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Request / response body types
// ---------------------------------------------------------------------------

type SetDataScopeBody struct {
	ScopeType string `json:"scope_type"`
	BranchID  string `json:"branch_id,omitempty"`
}

type CreateAPIKeyBody struct {
	Name      string   `json:"name"`
	Scopes    []string `json:"scopes,omitempty"`
	ExpiresAt string   `json:"expires_at,omitempty"`
	CreatedBy string   `json:"created_by,omitempty"`
}

type SetGlobalConfigBody struct {
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	UpdatedBy   string `json:"updated_by,omitempty"`
}

// ---------------------------------------------------------------------------
// SetDataScope — PUT /v1/admin/users/:id/data-scope (bearer)
// ---------------------------------------------------------------------------

func (s *Server) SetDataScope(c *fiber.Ctx, userID string) error {
	const op = "rest_oapi.Server.SetDataScope"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("user_id", userID))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("user_id", userID).Msg("")

	var body SetDataScopeBody
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.SetDataScope(ctx, &iam_grpc_adapter.SetDataScopeParams{
		UserID:    userID,
		ScopeType: body.ScopeType,
		BranchID:  body.BranchID,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"user_id":    result.UserID,
			"scope_type": result.ScopeType,
		},
	})
}

// ---------------------------------------------------------------------------
// CreateAPIKey — POST /v1/admin/api-keys (bearer)
// ---------------------------------------------------------------------------

func (s *Server) CreateAPIKey(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateAPIKey"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body CreateAPIKeyBody
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.CreateAPIKey(ctx, &iam_grpc_adapter.CreateAPIKeyParams{
		Name:      body.Name,
		Scopes:    body.Scopes,
		ExpiresAt: body.ExpiresAt,
		CreatedBy: body.CreatedBy,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "created")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": fiber.Map{
			"key_id":        result.KeyID,
			"plaintext_key": result.PlaintextKey,
			"key_prefix":    result.KeyPrefix,
			"expires_at":    result.ExpiresAt,
		},
	})
}

// ---------------------------------------------------------------------------
// RevokeAPIKey — DELETE /v1/admin/api-keys/:id (bearer)
// ---------------------------------------------------------------------------

func (s *Server) RevokeAPIKey(c *fiber.Ctx, keyID string) error {
	const op = "rest_oapi.Server.RevokeAPIKey"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("key_id", keyID))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("key_id", keyID).Msg("")

	result, err := s.svc.RevokeAPIKey(ctx, keyID)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "deleted")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"key_id":     result.KeyID,
			"revoked_at": result.RevokedAt,
		},
	})
}

// ---------------------------------------------------------------------------
// GetGlobalConfig — GET /v1/admin/config (bearer)
// ---------------------------------------------------------------------------

func (s *Server) GetGlobalConfig(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetGlobalConfig"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	// ?keys=a,b,c — optional comma-separated key filter
	keysParam := c.Query("keys")
	var keys []string
	if keysParam != "" {
		for _, k := range strings.Split(keysParam, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				keys = append(keys, k)
			}
		}
	}

	result, err := s.svc.GetGlobalConfig(ctx, keys)
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	type configEntryJSON struct {
		Key         string `json:"key"`
		Value       string `json:"value"`
		Description string `json:"description,omitempty"`
		UpdatedAt   string `json:"updated_at"`
	}
	data := make([]configEntryJSON, 0, len(result.Configs))
	for _, c := range result.Configs {
		data = append(data, configEntryJSON{
			Key:         c.Key,
			Value:       c.Value,
			Description: c.Description,
			UpdatedAt:   c.UpdatedAt,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": data})
}

// ---------------------------------------------------------------------------
// SetGlobalConfig — PUT /v1/admin/config/:key (bearer)
// ---------------------------------------------------------------------------

func (s *Server) SetGlobalConfig(c *fiber.Ctx, key string) error {
	const op = "rest_oapi.Server.SetGlobalConfig"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("config_key", key))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Str("key", key).Msg("")

	var body SetGlobalConfigBody
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	result, err := s.svc.SetGlobalConfig(ctx, &iam_grpc_adapter.SetGlobalConfigParams{
		Key:         key,
		Value:       body.Value,
		Description: body.Description,
		UpdatedBy:   body.UpdatedBy,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"key":        result.Key,
			"value":      result.Value,
			"updated_at": result.UpdatedAt,
		},
	})
}

// ---------------------------------------------------------------------------
// SearchActivityLog — GET /v1/admin/activity-log (bearer)
// ---------------------------------------------------------------------------

func (s *Server) SearchActivityLog(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.SearchActivityLog"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	userID := c.Query("user_id")
	resource := c.Query("resource")
	action := c.Query("action")
	from := c.Query("from")
	to := c.Query("to")
	cursor := c.Query("cursor")
	limit := int32(c.QueryInt("limit", 50))

	result, err := s.svc.SearchActivityLog(ctx, &iam_grpc_adapter.SearchActivityLogParams{
		UserID:   userID,
		Resource: resource,
		Action:   action,
		From:     from,
		To:       to,
		Limit:    limit,
		Cursor:   cursor,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}

	type logEntryJSON struct {
		ID         string `json:"id"`
		UserID     string `json:"user_id"`
		Resource   string `json:"resource"`
		Action     string `json:"action"`
		ResourceID string `json:"resource_id,omitempty"`
		CreatedAt  string `json:"created_at"`
	}
	data := make([]logEntryJSON, 0, len(result.Logs))
	for _, e := range result.Logs {
		data = append(data, logEntryJSON{
			ID:         e.ID,
			UserID:     e.UserID,
			Resource:   e.Resource,
			Action:     e.Action,
			ResourceID: e.ResourceID,
			CreatedAt:  e.CreatedAt,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":        data,
		"next_cursor": result.NextCursor,
	})
}
