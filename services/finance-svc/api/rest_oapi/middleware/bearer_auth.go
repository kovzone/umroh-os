package middleware

import (
	"errors"
	"fmt"
	"strings"

	"finance-svc/adapter/iam_grpc_adapter"
	"finance-svc/util/apperrors"

	"github.com/gofiber/fiber/v2"
)

const (
	// AuthorizationHeader is the HTTP header that carries the Bearer token.
	AuthorizationHeader = "Authorization"

	// AuthorizationBearer is the scheme prefix expected in the Authorization header.
	AuthorizationBearer = "bearer"

	// PayloadKey is the Fiber Locals key where the verified identity envelope
	// is stored. Handlers retrieve it with:
	//   c.Locals(middleware.PayloadKey).(*iam_grpc_adapter.ValidateTokenResult)
	PayloadKey = "auth_payload"
)

// IamValidator is the narrow slice of iam_grpc_adapter the middleware needs.
// Keeping it as a package-local interface lets tests inject a double without
// spinning up a real gRPC server.
type IamValidator interface {
	ValidateToken(ctx *fiber.Ctx, params *iam_grpc_adapter.ValidateTokenParams) (*iam_grpc_adapter.ValidateTokenResult, error)
}

// iamValidatorFn lets us adapt the real adapter method (which takes a
// context.Context) to the fiber-friendly signature above without declaring
// multiple interfaces.
type iamValidatorFn func(ctx *fiber.Ctx, params *iam_grpc_adapter.ValidateTokenParams) (*iam_grpc_adapter.ValidateTokenResult, error)

func (f iamValidatorFn) ValidateToken(ctx *fiber.Ctx, params *iam_grpc_adapter.ValidateTokenParams) (*iam_grpc_adapter.ValidateTokenResult, error) {
	return f(ctx, params)
}

// RequireBearerToken returns a Fiber middleware that:
//  1. Extracts the Bearer token from the Authorization header.
//  2. Calls iam-svc.ValidateToken via the adapter (fail-closed on any failure).
//  3. Stores the identity envelope in c.Locals(PayloadKey).
//
// Per F1 acceptance: missing / malformed / invalid / unreachable iam all map
// to 401 — never allow-by-default.
func RequireBearerToken(adapter *iam_grpc_adapter.Adapter) fiber.Handler {
	validator := iamValidatorFn(func(c *fiber.Ctx, params *iam_grpc_adapter.ValidateTokenParams) (*iam_grpc_adapter.ValidateTokenResult, error) {
		return adapter.ValidateToken(c.UserContext(), params)
	})
	return RequireBearerTokenWithValidator(validator)
}

// RequireBearerTokenWithValidator is the testable variant. Production code
// uses RequireBearerToken which wires the real gRPC adapter.
func RequireBearerTokenWithValidator(validator IamValidator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(AuthorizationHeader)
		if authHeader == "" {
			return errors.Join(apperrors.ErrUnauthorized, errors.New("authorization header is missing"))
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], AuthorizationBearer) {
			return errors.Join(apperrors.ErrUnauthorized, errors.New("invalid authorization header format"))
		}

		result, err := validator.ValidateToken(c, &iam_grpc_adapter.ValidateTokenParams{AccessToken: parts[1]})
		if err != nil {
			// Adapter already mapped the gRPC status to an apperrors sentinel.
			// Tag the context so the error middleware keeps the correct code
			// but don't leak the raw bearer anywhere.
			return fmt.Errorf("bearer validation failed: %w", err)
		}

		c.Locals(PayloadKey, result)
		return c.Next()
	}
}
