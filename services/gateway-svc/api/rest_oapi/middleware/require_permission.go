package middleware

import (
	"context"
	"errors"
	"fmt"

	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/util/apperrors"

	"github.com/gofiber/fiber/v2"
)

// PermissionChecker is the minimum surface the permission middleware needs. A
// real *iam_grpc_adapter.Adapter satisfies it; unit tests substitute a stub.
type PermissionChecker interface {
	CheckPermission(ctx context.Context, params *iam_grpc_adapter.CheckPermissionParams) (*iam_grpc_adapter.CheckPermissionResult, error)
}

// RequirePermission returns a Fiber middleware that enforces a specific
// (resource, action, scope) authorization tuple via iam-svc.CheckPermission.
// Compose it AFTER RequireBearerToken on a protected route — the identity
// envelope produced by bearer auth is the input to this check.
//
// Flow:
//  1. Read *Identity from c.Locals(IdentityKey). If absent, the middleware
//     chain was misconfigured (RequirePermission mounted without
//     RequireBearerToken upstream) → 500 ErrInternal. This is a bug guard,
//     not a user-facing condition; a correct deployment never hits it.
//  2. Call iam-svc.CheckPermission via the adapter.
//  3. On iam-svc transport failure (unreachable / timeout) → bubble the
//     wrapped error; the adapter already produced ErrServiceUnavailable → 502.
//     Fail-closed per F1-W7.
//  4. On allowed=false → 403 ErrForbidden.
//  5. On allowed=true → continue.
//
// The checker is passed by interface so unit tests can substitute a stub.
func RequirePermission(checker PermissionChecker, resource, action, scope string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		val := c.Locals(IdentityKey)
		identity, ok := val.(*Identity)
		if !ok || identity == nil {
			return errors.Join(apperrors.ErrInternal, errors.New("permission middleware invoked without identity in locals (mount RequireBearerToken first)"))
		}

		result, err := checker.CheckPermission(c.UserContext(), &iam_grpc_adapter.CheckPermissionParams{
			UserID:   identity.UserID,
			Resource: resource,
			Action:   action,
			Scope:    scope,
		})
		if err != nil {
			return fmt.Errorf("permission check: %w", err)
		}

		if !result.Allowed {
			return errors.Join(apperrors.ErrForbidden, fmt.Errorf("user %s lacks %s/%s/%s", identity.UserID, resource, action, scope))
		}

		return c.Next()
	}
}
