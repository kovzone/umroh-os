package rest_oapi

import (
	"errors"
	"fmt"
	"net/netip"

	"iam-svc/api/rest_oapi/middleware"
	"iam-svc/service"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"
	"iam-svc/util/token"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Login implements ServerInterface.
// POST /v1/sessions — email + password → access/refresh token pair. Public endpoint.
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

	result, err := s.svc.Login(ctx, &service.LoginParams{
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
	resp.Data.User = userProfileFromService(result.User)

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}

// RefreshSession implements ServerInterface.
// POST /v1/sessions/refresh — rotate refresh token, return new pair. Public endpoint.
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

	result, err := s.svc.RefreshSession(ctx, &service.RefreshSessionParams{
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

// Logout implements ServerInterface.
// DELETE /v1/sessions — revoke the current session row. Requires bearer auth.
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

	payload, err := payloadFromLocals(c)
	if err != nil {
		return err
	}

	_, err = s.svc.Logout(ctx, &service.LogoutParams{SessionID: payload.ID.String()})
	if err != nil {
		logger.Warn().Err(err).Str("session_id", payload.ID.String()).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, "success")
	return c.SendStatus(fiber.StatusNoContent)
}

// ------------------------ helpers ------------------------

// userProfileFromService maps the service-layer UserProfile into the OpenAPI-generated shape.
func userProfileFromService(u service.UserProfile) UserProfile {
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

// parseUUIDOrZero converts the canonical UUID string to the oapi types type
// (which is an alias for google/uuid.UUID). Returns the zero value if the
// string is malformed — the service layer always emits canonical strings, so
// the fallback only masks logic bugs rather than user input.
func parseUUIDOrZero(s string) openapi_types.UUID {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return openapi_types.UUID{}
	}
	return parsed
}

// clientIP extracts the caller's IP as *netip.Addr. Prefers X-Forwarded-For leftmost when
// present, falls back to c.IP() (the remote address the proxy saw). Returns nil if the
// address cannot be parsed — iam.sessions.ip is nullable, so nil is safe.
func clientIP(c *fiber.Ctx) *netip.Addr {
	raw := c.Get(fiber.HeaderXForwardedFor)
	if raw == "" {
		raw = c.IP()
	} else {
		// leftmost entry is the original client
		for i := 0; i < len(raw); i++ {
			if raw[i] == ',' {
				raw = raw[:i]
				break
			}
		}
	}
	addr, err := netip.ParseAddr(trimSpace(raw))
	if err != nil {
		return nil
	}
	return &addr
}

// trimSpace trims leading/trailing ASCII space without pulling in strings. Fast path for IPs.
func trimSpace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

// payloadFromLocals extracts the verified PASETO payload that RequireBearerToken stashed.
// Returns ErrUnauthorized if absent (handler was wired without the middleware — a bug).
func payloadFromLocals(c *fiber.Ctx) (*token.Payload, error) {
	val := c.Locals(middleware.PayloadKey)
	payload, ok := val.(*token.Payload)
	if !ok || payload == nil {
		return nil, errors.Join(apperrors.ErrUnauthorized, errors.New("missing token payload"))
	}
	return payload, nil
}
