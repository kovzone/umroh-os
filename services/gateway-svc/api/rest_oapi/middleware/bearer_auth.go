package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/util/apperrors"

	"github.com/gofiber/fiber/v2"
)

const (
	// AuthorizationHeader is the HTTP header that carries the Bearer token.
	AuthorizationHeader = "Authorization"

	// AuthorizationBearer is the scheme prefix expected in the Authorization header.
	AuthorizationBearer = "bearer"

	// IdentityKey is the Fiber Locals key where the validated identity envelope
	// is stored. Downstream handlers read it with:
	//   id := c.Locals(middleware.IdentityKey).(*middleware.Identity)
	IdentityKey = "auth_identity"
)

// Identity is the validated principal attached to the Fiber context after the
// bearer middleware succeeds. Proto types stay in the adapter; handlers see
// plain Go fields only.
type Identity struct {
	UserID    string
	BranchID  string
	SessionID string
	Roles     []string
}

// IamValidator is the minimum surface the bearer middleware needs. A real
// *iam_grpc_adapter.Adapter satisfies it; unit tests substitute a stub.
type IamValidator interface {
	ValidateToken(ctx context.Context, params *iam_grpc_adapter.ValidateTokenParams) (*iam_grpc_adapter.ValidateTokenResult, error)
}

// RequireBearerToken returns a Fiber middleware that enforces ADR 0009 /
// F1-W7 edge auth:
//
//  1. Extract `Authorization: Bearer <token>` from the incoming HTTP request.
//     Missing or malformed header → 401 ErrUnauthorized. Does not call iam.
//  2. Call iam-svc.ValidateToken via the gRPC adapter.
//  3. On iam-svc transport failure (unreachable / timeout) → 502
//     ErrServiceUnavailable. Fail-closed: request never reaches the backend.
//  4. On iam-svc Unauthenticated → 401 ErrUnauthorized.
//  5. On success: attach *Identity to c.Locals(IdentityKey); continue.
//
// Backends do not re-validate the bearer (per ADR 0009); the identity
// envelope is what the gateway hands downstream. The adapter is passed by
// interface so unit tests can substitute a stub.
func RequireBearerToken(validator IamValidator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(AuthorizationHeader)
		if authHeader == "" {
			return errors.Join(apperrors.ErrUnauthorized, errors.New("authorization header is missing"))
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], AuthorizationBearer) {
			return errors.Join(apperrors.ErrUnauthorized, errors.New("invalid authorization header format"))
		}

		result, err := validator.ValidateToken(c.UserContext(), &iam_grpc_adapter.ValidateTokenParams{
			AccessToken: parts[1],
		})
		if err != nil {
			// Adapter has already mapped transport failures to
			// ErrServiceUnavailable (502) and iam-returned Unauthenticated to
			// ErrUnauthorized (401). Bubble the wrapped error unchanged so the
			// central error handler renders the right envelope.
			return fmt.Errorf("bearer validation: %w", err)
		}

		c.Locals(IdentityKey, &Identity{
			UserID:    result.UserID,
			BranchID:  result.BranchID,
			SessionID: result.SessionID,
			Roles:     result.Roles,
		})
		return c.Next()
	}
}
