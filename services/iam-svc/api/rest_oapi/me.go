package rest_oapi

import (
	"errors"
	"fmt"

	"iam-svc/service"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetMe implements ServerInterface.
// GET /v1/me — current user profile. Requires bearer auth.
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

	payload, err := payloadFromLocals(c)
	if err != nil {
		return err
	}

	result, err := s.svc.GetMe(ctx, &service.GetMeParams{UserID: payload.UserID})
	if err != nil {
		logger.Warn().Err(err).Str("user_id", payload.UserID).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	resp := GetMeResponse{}
	resp.Data.User = userProfileFromService(result.User)
	resp.Data.TotpEnrolled = result.TOTPEnrolled
	resp.Data.TotpVerified = result.TOTPVerified

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}

// EnrollTOTP implements ServerInterface.
// POST /v1/me/2fa/enroll — generate + persist TOTP secret, return otpauth URL. Requires bearer auth.
func (s *Server) EnrollTOTP(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.EnrollTOTP"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/me/2fa/enroll"),
		attribute.String("method", "POST"),
	)

	payload, err := payloadFromLocals(c)
	if err != nil {
		return err
	}

	result, err := s.svc.EnrollTOTP(ctx, &service.EnrollTOTPParams{UserID: payload.UserID})
	if err != nil {
		logger.Warn().Err(err).Str("user_id", payload.UserID).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	resp := EnrollTOTPResponse{}
	resp.Data.Secret = result.Secret
	resp.Data.OtpauthUrl = result.OtpauthURL

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}

// VerifyTOTP implements ServerInterface.
// POST /v1/me/2fa/verify — validate TOTP code, stamp verified_at. Requires bearer auth.
func (s *Server) VerifyTOTP(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.VerifyTOTP"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/me/2fa/verify"),
		attribute.String("method", "POST"),
	)

	payload, err := payloadFromLocals(c)
	if err != nil {
		return err
	}

	var req VerifyTOTPRequest
	if err := c.BodyParser(&req); err != nil {
		e := errors.Join(apperrors.ErrValidation, fmt.Errorf("parse body: %w", err))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return e
	}

	result, err := s.svc.VerifyTOTP(ctx, &service.VerifyTOTPParams{
		UserID: payload.UserID,
		Code:   req.Code,
	})
	if err != nil {
		logger.Warn().Err(err).Str("user_id", payload.UserID).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	resp := VerifyTOTPResponse{}
	resp.Data.VerifiedAt = result.VerifiedAt

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}
