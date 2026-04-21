package rest_oapi

import (
	"errors"

	"finance-svc/adapter/iam_grpc_adapter"
	"finance-svc/api/rest_oapi/middleware"
	"finance-svc/service"
	"finance-svc/util/apperrors"
	"finance-svc/util/logging"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// FinancePing implements ServerInterface.
// GET /v1/finance/ping — BL-IAM-002 placeholder. Requires bearer auth.
func (s *Server) FinancePing(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.FinancePing"

	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("endpoint", "/v1/finance/ping"),
		attribute.String("method", "GET"),
	)

	payload, err := payloadFromLocals(c)
	if err != nil {
		return err
	}

	result, err := s.svc.FinancePing(ctx, &service.FinancePingParams{
		UserID: payload.UserID,
		Roles:  payload.Roles,
	})
	if err != nil {
		logger.Warn().Err(err).Str("user_id", payload.UserID).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	// Parse the user_id string into openapi_types.UUID for the generated schema.
	userUUID, err := uuid.Parse(result.UserID)
	if err != nil {
		e := errors.Join(apperrors.ErrInternal, errors.New("iam-returned user_id is not a valid uuid"))
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return e
	}

	resp := FinancePingResponse{}
	resp.Data.Message = result.Message
	resp.Data.UserId = openapi_types.UUID(userUUID)
	resp.Data.Roles = result.Roles

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}

// payloadFromLocals extracts the verified identity envelope that
// RequireBearerToken stashed. Absence signals the handler was wired without
// the middleware — that's a programmer bug, not an auth failure, but we still
// surface it as ErrUnauthorized so the caller sees a consistent 401.
func payloadFromLocals(c *fiber.Ctx) (*iam_grpc_adapter.ValidateTokenResult, error) {
	val := c.Locals(middleware.PayloadKey)
	payload, ok := val.(*iam_grpc_adapter.ValidateTokenResult)
	if !ok || payload == nil {
		return nil, errors.Join(apperrors.ErrUnauthorized, errors.New("missing token payload"))
	}
	return payload, nil
}
