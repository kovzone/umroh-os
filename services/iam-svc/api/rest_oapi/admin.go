package rest_oapi

import (
	"errors"

	"iam-svc/service"
	"iam-svc/util/apperrors"
	"iam-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Permission tuple enforced on POST /v1/users/{id}/suspend. Granted to
// super_admin by migration 000006_seed_iam_user_suspend_permission.
const (
	suspendResource = "iam.users"
	suspendAction   = "suspend"
	suspendScope    = "global"
)

// SuspendUser implements ServerInterface.
// POST /v1/users/{id}/suspend — flip iam.users.status=suspended + revoke every
// active session for the target in one transaction. Bearer + super_admin
// permission required. See service.SuspendUser for the underlying logic and
// the openapi description for the full F1-W5 rationale.
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

	payload, err := payloadFromLocals(c)
	if err != nil {
		return err
	}

	// Permission gate — call the in-process service (not a gRPC loopback). An
	// un-granted bearer returns 403 before the SuspendUser code path runs so
	// the 404 branch never leaks "target exists vs target does not" to a
	// caller who lacked authority to ask in the first place.
	perm, err := s.svc.CheckPermission(ctx, &service.CheckPermissionParams{
		UserID:   payload.UserID,
		Resource: suspendResource,
		Action:   suspendAction,
		Scope:    suspendScope,
	})
	if err != nil {
		logger.Warn().Err(err).Str("actor_user_id", payload.UserID).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	if !perm.Allowed {
		e := errors.Join(apperrors.ErrForbidden, errors.New("missing iam.users/suspend/global permission"))
		logger.Warn().Err(e).Str("actor_user_id", payload.UserID).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return e
	}

	result, err := s.svc.SuspendUser(ctx, &service.SuspendUserParams{
		ActorUserID:  payload.UserID,
		TargetUserID: id.String(),
	})
	if err != nil {
		logger.Warn().Err(err).
			Str("actor_user_id", payload.UserID).
			Str("target_user_id", id.String()).
			Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	resp := SuspendUserResponse{}
	resp.Data.User = userProfileFromService(result.User)

	span.SetStatus(codes.Ok, "success")
	return c.Status(fiber.StatusOK).JSON(resp)
}
