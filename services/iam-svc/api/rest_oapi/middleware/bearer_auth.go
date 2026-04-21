package middleware

import (
	"errors"
	"fmt"
	"strings"

	"iam-svc/util/apperrors"
	"iam-svc/util/token"

	"github.com/gofiber/fiber/v2"
)

const (
	// AuthorizationHeader is the HTTP header that carries the Bearer token.
	AuthorizationHeader = "Authorization"

	// AuthorizationBearer is the scheme prefix expected in the Authorization header.
	AuthorizationBearer = "bearer"

	// PayloadKey is the Fiber Locals key where the verified *token.Payload is stored.
	// Handlers retrieve it with: c.Locals(middleware.PayloadKey).(*token.Payload)
	PayloadKey = "auth_payload"
)

// RequireBearerToken returns a Fiber middleware that:
//  1. Extracts the Bearer token from the Authorization header.
//  2. Verifies the token using the provided Maker (PASETO or JWT per config.token.type).
//  3. Stores the *token.Payload in c.Locals(PayloadKey) for downstream handlers.
//
// On any failure (missing header, malformed header, invalid/expired token) the
// middleware returns an apperrors.ErrUnauthorized error so the central error
// middleware renders a consistent `{error: {code: "UNAUTHORIZED", ...}}` envelope.
func RequireBearerToken(maker token.Maker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(AuthorizationHeader)
		if authHeader == "" {
			return errors.Join(apperrors.ErrUnauthorized, errors.New("authorization header is missing"))
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], AuthorizationBearer) {
			return errors.Join(apperrors.ErrUnauthorized, errors.New("invalid authorization header format"))
		}

		payload, err := maker.VerifyToken(parts[1])
		if err != nil {
			return errors.Join(apperrors.ErrUnauthorized, fmt.Errorf("verify bearer token: %w", err))
		}

		c.Locals(PayloadKey, payload)
		return c.Next()
	}
}
