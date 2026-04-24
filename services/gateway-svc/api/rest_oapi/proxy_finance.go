package rest_oapi

import (
	"errors"

	"gateway-svc/api/rest_oapi/middleware"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetFinancePing is the permission-gate smoke route (BL-IAM-002 /
// BL-IAM-019). The middleware chain on this endpoint (RequireBearerToken +
// RequirePermission("journal_entry", "read", "global")) has already approved
// the caller by the time this handler runs; finance-svc is called only to
// prove the backend is reachable. The client-visible envelope is assembled
// here from the gateway's auth locals — finance-svc is identity-agnostic.
//
// GetFinancePing implements ServerInterface.
func (s *Server) GetFinancePing(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetFinancePing"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/finance/ping"),
		attribute.String("method", "GET"),
	)

	identity, ok := c.Locals(middleware.IdentityKey).(*middleware.Identity)
	if !ok || identity == nil {
		// Bug guard: this handler is always mounted behind RequireBearerToken.
		// A nil/missing identity means the middleware chain was misconfigured.
		err := errors.Join(apperrors.ErrInternal, errors.New("finance ping invoked without identity in locals"))
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	logger.Info().Str("op", op).Str("user_id", identity.UserID).Msg("")

	result, err := s.svc.FinancePing(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	var resp FinancePingResponse
	resp.Data.Message = result.Message
	resp.Data.UserId = parseUUIDOrZero(identity.UserID)
	resp.Data.BranchId = parseUUIDOrZero(identity.BranchID)
	resp.Data.Roles = append([]string{}, identity.Roles...)

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}
