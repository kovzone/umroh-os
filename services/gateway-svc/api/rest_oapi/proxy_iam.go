package rest_oapi

// proxy_iam.go — REST handler implementations for IAM routes on gateway-svc.
//
// Interim (retired in S1-E-12 final cleanup): GetIamSystemLive and
// GetIamSystemDbTxDiagnostic forward to iam-svc via REST adapter.
//
// New (BL-IAM-018 / S1-E-12): Login, RefreshSession, Logout, GetMe,
// EnrollTOTP, VerifyTOTP, SuspendUser forward to iam-svc via gRPC adapter.
//
// Identity extraction (bearer routes): the RequireBearerToken middleware
// injects *middleware.Identity into c.Locals(middleware.IdentityKey) before
// the handler runs. Handlers cast it for user_id / session_id claims.

import (
	"errors"

	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/api/rest_oapi/middleware"
	"gateway-svc/util/logging"
	"gateway-svc/util/apperrors"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetIamSystemLive proxies the iam-svc liveness probe through the gateway.
// Scaffold-time proof of the REST adapter pattern.
//
// GetIamSystemLive implements ServerInterface.
func (s *Server) GetIamSystemLive(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetIamSystemLive"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/iam/system/live"),
		attribute.String("method", "GET"),
	)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetIamSystemLive(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "BAD_GATEWAY",
				"message": err.Error(),
			},
		})
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(LiveResponse{
		Data: struct {
			Ok bool `json:"ok"`
		}{Ok: result.OK},
	})
}

// GetIamSystemDbTxDiagnostic proxies iam-svc's state-mutating DB transaction
// diagnostic through the gateway. It's the traced cross-service path the
// S0-J-05 observability acceptance uses: one request produces one trace
// spanning gateway-svc + iam-svc, with matching trace_id log lines in both
// containers' Loki streams.
//
// GetIamSystemDbTxDiagnostic implements ServerInterface.
func (s *Server) GetIamSystemDbTxDiagnostic(c *fiber.Ctx, params GetIamSystemDbTxDiagnosticParams) error {
	const op = "rest_oapi.Server.GetIamSystemDbTxDiagnostic"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)

	message := "no message"
	if params.Message != nil {
		message = *params.Message
	}

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/iam/system/diagnostics/db-tx"),
		attribute.String("method", "GET"),
		attribute.String("message", message),
	)
	logger.Info().Str("op", op).Str("message", message).Msg("")

	result, err := s.svc.GetIamSystemDbTxDiagnostic(ctx, message)
	if err != nil {
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    "BAD_GATEWAY",
				"message": err.Error(),
			},
		})
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(DbTxDiagnosticResponse{
		Data: struct {
			DiagnosticId int64  `json:"diagnostic_id"`
			Message      string `json:"message"`
		}{
			DiagnosticId: result.DiagnosticID,
			Message:      result.Message,
		},
	})
}

// ---------------------------------------------------------------------------
// IAM auth handlers (BL-IAM-018 / S1-E-12)
// ---------------------------------------------------------------------------

// Login handles POST /v1/auth/login (public — no bearer required).
// Forwards email + password (+ optional TOTP code) to iam-svc.Login over gRPC.
// Client IP and User-Agent are forwarded so iam-svc can record them in sessions.
func (s *Server) Login(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.Login"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		TOTPCode string `json:"totp_code"`
	}
	if err := c.BodyParser(&body); err != nil {
		span.SetStatus(codes.Error, "invalid body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{"code": "BAD_REQUEST", "message": "invalid request body"},
		})
	}

	result, err := s.svc.Login(ctx, &iam_grpc_adapter.LoginParams{
		Email:     body.Email,
		Password:  body.Password,
		TOTPCode:  body.TOTPCode,
		UserAgent: c.Get("User-Agent"),
		IP:        c.IP(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Warn().Err(err).Msg("")
		return writeIamError(c, err)
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"access_token":       result.AccessToken,
			"refresh_token":      result.RefreshToken,
			"access_expires_at":  result.AccessExpiresAt,
			"refresh_expires_at": result.RefreshExpiresAt,
			"user":               userProfileToJSON(result.User),
		},
	})
}

// RefreshSession handles POST /v1/auth/refresh (public — no bearer required).
// Rotates the refresh token and issues a new access + refresh pair.
func (s *Server) RefreshSession(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RefreshSession"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&body); err != nil {
		span.SetStatus(codes.Error, "invalid body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{"code": "BAD_REQUEST", "message": "invalid request body"},
		})
	}

	result, err := s.svc.RefreshSession(ctx, &iam_grpc_adapter.RefreshSessionParams{
		RefreshToken: body.RefreshToken,
		UserAgent:    c.Get("User-Agent"),
		IP:           c.IP(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Warn().Err(err).Msg("")
		return writeIamError(c, err)
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"access_token":       result.AccessToken,
			"refresh_token":      result.RefreshToken,
			"access_expires_at":  result.AccessExpiresAt,
			"refresh_expires_at": result.RefreshExpiresAt,
		},
	})
}

// Logout handles DELETE /v1/auth/logout (bearer-required).
// Revokes the session identified by the session_id claim in the validated token.
// Idempotent — a second call on an already-revoked session succeeds silently.
func (s *Server) Logout(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.Logout"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	id := c.Locals(middleware.IdentityKey).(*middleware.Identity)
	span.SetAttributes(attribute.String("session_id", id.SessionID))
	logger.Info().Str("op", op).Str("session_id", id.SessionID).Msg("")

	if err := s.svc.Logout(ctx, &iam_grpc_adapter.LogoutParams{
		SessionID: id.SessionID,
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Warn().Err(err).Msg("")
		return writeIamError(c, err)
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// GetMe handles GET /v1/me (bearer-required).
// Returns the current user's profile + TOTP state.
func (s *Server) GetMe(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetMe"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	id := c.Locals(middleware.IdentityKey).(*middleware.Identity)
	span.SetAttributes(attribute.String("user_id", id.UserID))
	logger.Info().Str("op", op).Str("user_id", id.UserID).Msg("")

	result, err := s.svc.GetMe(ctx, &iam_grpc_adapter.GetMeParams{UserID: id.UserID})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Warn().Err(err).Msg("")
		return writeIamError(c, err)
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"user":          userProfileToJSON(result.User),
			"totp_enrolled": result.TOTPEnrolled,
			"totp_verified": result.TOTPVerified,
		},
	})
}

// EnrollTOTP handles POST /v1/me/2fa/enroll (bearer-required).
// Generates a TOTP secret + otpauth URL for the current user.
// Returns 409 if TOTP is already verified.
func (s *Server) EnrollTOTP(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.EnrollTOTP"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	id := c.Locals(middleware.IdentityKey).(*middleware.Identity)
	span.SetAttributes(attribute.String("user_id", id.UserID))
	logger.Info().Str("op", op).Str("user_id", id.UserID).Msg("")

	result, err := s.svc.EnrollTOTP(ctx, &iam_grpc_adapter.EnrollTOTPParams{UserID: id.UserID})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Warn().Err(err).Msg("")
		return writeIamError(c, err)
	}

	span.SetStatus(codes.Ok, "success")
	// NOTE: Do NOT log result.Secret — it is a plaintext secret shown once.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"secret":      result.Secret,
			"otpauth_url": result.OtpauthURL,
		},
	})
}

// VerifyTOTP handles POST /v1/me/2fa/verify (bearer-required).
// Validates the TOTP code and stamps totp_verified_at on success.
func (s *Server) VerifyTOTP(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.VerifyTOTP"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	id := c.Locals(middleware.IdentityKey).(*middleware.Identity)
	span.SetAttributes(attribute.String("user_id", id.UserID))
	logger.Info().Str("op", op).Str("user_id", id.UserID).Msg("")

	var body struct {
		Code string `json:"code"`
	}
	if err := c.BodyParser(&body); err != nil {
		span.SetStatus(codes.Error, "invalid body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{"code": "BAD_REQUEST", "message": "invalid request body"},
		})
	}

	result, err := s.svc.VerifyTOTP(ctx, &iam_grpc_adapter.VerifyTOTPParams{
		UserID: id.UserID,
		Code:   body.Code,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Warn().Err(err).Msg("")
		return writeIamError(c, err)
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"verified_at": result.VerifiedAt,
		},
	})
}

// SuspendUser handles POST /v1/users/{id}/suspend (bearer + permission).
// The bearer middleware has already validated the token; permission gating
// (iam.users/suspend/global) is done here before forwarding to iam-svc.
func (s *Server) SuspendUser(c *fiber.Ctx, targetUserID string) error {
	const op = "rest_oapi.Server.SuspendUser"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	id := c.Locals(middleware.IdentityKey).(*middleware.Identity)
	span.SetAttributes(
		attribute.String("actor_user_id", id.UserID),
		attribute.String("target_user_id", targetUserID),
	)
	logger.Info().
		Str("op", op).
		Str("actor_user_id", id.UserID).
		Str("target_user_id", targetUserID).
		Msg("")

	// Permission gate: actor must hold iam.users/suspend/global.
	// The bearer middleware validated the token; this enforces the grant.
	perm, err := s.svc.CheckPermission(ctx, &iam_grpc_adapter.CheckPermissionParams{
		UserID:   id.UserID,
		Resource: "iam.users",
		Action:   "suspend",
		Scope:    "global",
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Warn().Err(err).Msg("CheckPermission call failed")
		return writeIamError(c, err)
	}
	if !perm.Allowed {
		span.SetStatus(codes.Error, "forbidden")
		logger.Warn().Str("user_id", id.UserID).Msg("suspend denied: no iam.users/suspend/global grant")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": fiber.Map{"code": "FORBIDDEN", "message": "akses ditolak"},
		})
	}

	result, err := s.svc.SuspendUser(ctx, &iam_grpc_adapter.SuspendUserParams{
		ActorUserID:  id.UserID,
		TargetUserID: targetUserID,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Warn().Err(err).Msg("")
		return writeIamError(c, err)
	}

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"user": userProfileToJSON(result.User),
		},
	})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// writeIamError maps an apperrors sentinel to the appropriate HTTP status and
// writes a standard ErrorResponse JSON body. Span error recording is the
// caller's responsibility (each handler records before calling this).
func writeIamError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, apperrors.ErrUnauthorized):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fiber.Map{"code": "UNAUTHORIZED", "message": "tidak terautentikasi"},
		})
	case errors.Is(err, apperrors.ErrForbidden):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": fiber.Map{"code": "FORBIDDEN", "message": "akses ditolak"},
		})
	case errors.Is(err, apperrors.ErrValidation):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
	case errors.Is(err, apperrors.ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fiber.Map{"code": "NOT_FOUND", "message": "data tidak ditemukan"},
		})
	case errors.Is(err, apperrors.ErrConflict):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": fiber.Map{"code": "CONFLICT", "message": err.Error()},
		})
	case errors.Is(err, apperrors.ErrServiceUnavailable):
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": fiber.Map{"code": "BAD_GATEWAY", "message": "layanan IAM sementara tidak tersedia"},
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{"code": "INTERNAL_ERROR", "message": "terjadi kesalahan tidak terduga"},
		})
	}
}

// userProfileToJSON converts adapter UserProfile to a JSON-serializable map.
func userProfileToJSON(u iam_grpc_adapter.UserProfile) fiber.Map {
	return fiber.Map{
		"user_id":   u.UserID,
		"email":     u.Email,
		"name":      u.Name,
		"branch_id": u.BranchID,
		"status":    u.Status,
	}
}
