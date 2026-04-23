package rest_oapi

import (
	"errors"
	"fmt"

	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/api/rest_oapi/middleware"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// BL-IAM-018 (S1-E-12) — gateway REST handlers that front iam-svc's gRPC
// auth surface. Each returns apperrors-wrapped errors so the central
// middleware.ErrorHandler() renders the UPPER_SNAKE envelope that matches
// iam-svc's legacy REST response shape 1:1 (existing clients + e2e bodies
// are unchanged). The bearer middleware (RequireBearerToken) is mounted per
// protected route in cmd/server.go and stashes the validated Identity in
// c.Locals(middleware.IdentityKey).

// Login proxies POST /v1/sessions. Public route (no bearer).
func (s *Server) Login(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.Login"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/sessions"),
		attribute.String("method", "POST"),
	)

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		e := errors.Join(apperrors.ErrValidation, fmt.Errorf("parse body: %w", err))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return e
	}

	totpCode := ""
	if req.TotpCode != nil {
		totpCode = *req.TotpCode
	}

	result, err := s.svc.Login(ctx, &iam_grpc_adapter.LoginParams{
		Email:     string(req.Email),
		Password:  req.Password,
		TOTPCode:  totpCode,
		UserAgent: c.Get(fiber.HeaderUserAgent),
		IP:        clientIP(c),
	})
	if err != nil {
		logger.Warn().Err(err).Str("email", string(req.Email)).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	resp := LoginResponse{}
	resp.Data.AccessToken = result.AccessToken
	resp.Data.RefreshToken = result.RefreshToken
	resp.Data.AccessExpiresAt = result.AccessExpiresAt
	resp.Data.RefreshExpiresAt = result.RefreshExpiresAt
	resp.Data.User = userProfileFromAdapter(result.User)

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}

// RefreshSession proxies POST /v1/sessions/refresh. Public route.
func (s *Server) RefreshSession(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RefreshSession"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/sessions/refresh"),
		attribute.String("method", "POST"),
	)

	var req RefreshSessionRequest
	if err := c.BodyParser(&req); err != nil {
		e := errors.Join(apperrors.ErrValidation, fmt.Errorf("parse body: %w", err))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return e
	}

	result, err := s.svc.RefreshSession(ctx, &iam_grpc_adapter.RefreshSessionParams{
		RefreshToken: req.RefreshToken,
		UserAgent:    c.Get(fiber.HeaderUserAgent),
		IP:           clientIP(c),
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	resp := RefreshSessionResponse{}
	resp.Data.AccessToken = result.AccessToken
	resp.Data.RefreshToken = result.RefreshToken
	resp.Data.AccessExpiresAt = result.AccessExpiresAt
	resp.Data.RefreshExpiresAt = result.RefreshExpiresAt

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}

// Logout proxies DELETE /v1/sessions. Protected route; session_id comes from
// the validated bearer's Identity envelope stashed by the middleware.
func (s *Server) Logout(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.Logout"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/sessions"),
		attribute.String("method", "DELETE"),
	)

	id, err := identityFromLocals(c)
	if err != nil {
		return err
	}

	if _, err := s.svc.Logout(ctx, &iam_grpc_adapter.LogoutParams{SessionID: id.SessionID}); err != nil {
		logger.Warn().Err(err).Str("session_id", id.SessionID).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "success")
	return c.SendStatus(fiber.StatusNoContent)
}

// GetMe proxies GET /v1/me. Protected.
func (s *Server) GetMe(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetMe"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/me"),
		attribute.String("method", "GET"),
	)

	id, err := identityFromLocals(c)
	if err != nil {
		return err
	}

	result, err := s.svc.GetMe(ctx, &iam_grpc_adapter.GetMeParams{UserID: id.UserID})
	if err != nil {
		logger.Warn().Err(err).Str("user_id", id.UserID).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	resp := GetMeResponse{}
	resp.Data.User = userProfileFromAdapter(result.User)
	resp.Data.TotpEnrolled = result.TOTPEnrolled
	resp.Data.TotpVerified = result.TOTPVerified

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}

// EnrollTotp proxies POST /v1/me/2fa/enroll. Protected.
func (s *Server) EnrollTotp(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.EnrollTotp"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/me/2fa/enroll"),
		attribute.String("method", "POST"),
	)

	id, err := identityFromLocals(c)
	if err != nil {
		return err
	}

	result, err := s.svc.EnrollTOTP(ctx, &iam_grpc_adapter.EnrollTOTPParams{UserID: id.UserID})
	if err != nil {
		logger.Warn().Err(err).Str("user_id", id.UserID).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	resp := EnrollTotpResponse{}
	resp.Data.Secret = result.Secret
	resp.Data.OtpauthUrl = result.OtpauthURL

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}

// VerifyTotp proxies POST /v1/me/2fa/verify. Protected.
func (s *Server) VerifyTotp(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.VerifyTotp"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/me/2fa/verify"),
		attribute.String("method", "POST"),
	)

	id, err := identityFromLocals(c)
	if err != nil {
		return err
	}

	var req VerifyTotpRequest
	if err := c.BodyParser(&req); err != nil {
		e := errors.Join(apperrors.ErrValidation, fmt.Errorf("parse body: %w", err))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return e
	}

	result, err := s.svc.VerifyTOTP(ctx, &iam_grpc_adapter.VerifyTOTPParams{
		UserID: id.UserID,
		Code:   req.Code,
	})
	if err != nil {
		logger.Warn().Err(err).Str("user_id", id.UserID).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	resp := VerifyTotpResponse{}
	resp.Data.VerifiedAt = result.VerifiedAt

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}

// SuspendUser proxies POST /v1/users/{id}/suspend. Protected; iam-svc
// enforces the iam.users/suspend/global permission gate + self-suspend guard.
func (s *Server) SuspendUser(c *fiber.Ctx, id openapi_types.UUID) error {
	const op = "rest_oapi.Server.SuspendUser"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/users/:id/suspend"),
		attribute.String("method", "POST"),
		attribute.String("target_user_id", id.String()),
	)

	actor, err := identityFromLocals(c)
	if err != nil {
		return err
	}

	result, err := s.svc.SuspendUser(ctx, &iam_grpc_adapter.SuspendUserParams{
		ActorUserID:  actor.UserID,
		TargetUserID: id.String(),
	})
	if err != nil {
		logger.Warn().Err(err).
			Str("actor_user_id", actor.UserID).
			Str("target_user_id", id.String()).
			Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	resp := SuspendUserResponse{}
	resp.Data.User = userProfileFromAdapter(result.User)

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// identityFromLocals extracts the validated *middleware.Identity stashed by
// RequireBearerToken. Absence means the handler was mounted without the
// middleware — a bug that surfaces as 401 rather than a panic.
func identityFromLocals(c *fiber.Ctx) (*middleware.Identity, error) {
	val := c.Locals(middleware.IdentityKey)
	id, ok := val.(*middleware.Identity)
	if !ok || id == nil {
		return nil, errors.Join(apperrors.ErrUnauthorized, errors.New("missing auth identity"))
	}
	return id, nil
}

// userProfileFromAdapter converts the adapter-local UserProfile into the
// oapi-generated one. BranchID is optional on the wire; empty string → omit.
func userProfileFromAdapter(u iam_grpc_adapter.UserProfile) UserProfile {
	out := UserProfile{
		UserId: parseUUIDOrZero(u.UserID),
		Email:  openapi_types.Email(u.Email),
		Name:   u.Name,
		Status: UserProfileStatus(u.Status),
	}
	if u.BranchID != "" {
		b := parseUUIDOrZero(u.BranchID)
		out.BranchId = &b
	}
	return out
}

// parseUUIDOrZero parses a canonical UUID string. Malformed input yields the
// zero value — iam-svc always emits canonical strings so the fallback only
// masks logic bugs.
func parseUUIDOrZero(s string) openapi_types.UUID {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return openapi_types.UUID{}
	}
	return parsed
}

// clientIP returns the caller's IP as a string (empty if unparseable).
// Mirrors the legacy iam-svc REST handler's X-Forwarded-For / c.IP()
// derivation so behavior is unchanged across the migration.
func clientIP(c *fiber.Ctx) string {
	raw := c.Get(fiber.HeaderXForwardedFor)
	if raw == "" {
		return c.IP()
	}
	// leftmost entry is the original client
	for i := 0; i < len(raw); i++ {
		if raw[i] == ',' {
			raw = raw[:i]
			break
		}
	}
	// trim ASCII whitespace
	start := 0
	for start < len(raw) && (raw[start] == ' ' || raw[start] == '\t') {
		start++
	}
	end := len(raw)
	for end > start && (raw[end-1] == ' ' || raw[end-1] == '\t') {
		end--
	}
	return raw[start:end]
}
