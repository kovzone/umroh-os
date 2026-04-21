package middleware_test

import (
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"finance-svc/adapter/iam_grpc_adapter"
	"finance-svc/api/rest_oapi/middleware"
	"finance-svc/util/apperrors"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// stubValidator lets each test set the expected result without importing
// testify/mock just for a one-method surface.
type stubValidator struct {
	result *iam_grpc_adapter.ValidateTokenResult
	err    error
}

func (v *stubValidator) ValidateToken(c *fiber.Ctx, params *iam_grpc_adapter.ValidateTokenParams) (*iam_grpc_adapter.ValidateTokenResult, error) {
	return v.result, v.err
}

// newAppWithMiddleware spins up a Fiber app that mounts the bearer middleware
// and a sentinel handler that reports the payload back via c.Locals.
func newAppWithMiddleware(validator middleware.IamValidator) *fiber.App {
	app := fiber.New()
	app.Use(middleware.ErrorHandler())
	app.Get("/ping",
		middleware.RequireBearerTokenWithValidator(validator),
		func(c *fiber.Ctx) error {
			payload, _ := c.Locals(middleware.PayloadKey).(*iam_grpc_adapter.ValidateTokenResult)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"user_id": payload.UserID})
		},
	)
	return app
}

func Test_RequireBearerToken_happyPath(t *testing.T) {
	validator := &stubValidator{
		result: &iam_grpc_adapter.ValidateTokenResult{
			UserID:    "user-1",
			BranchID:  "branch-1",
			SessionID: "session-1",
			Roles:     []string{"finance_admin"},
			ExpiresAt: time.Now().Add(time.Hour),
		},
	}
	app := newAppWithMiddleware(validator)

	req := httptest.NewRequest("GET", "/ping", nil)
	req.Header.Set("Authorization", "Bearer valid.token.here")

	resp, err := app.Test(req, int(5*time.Second/time.Millisecond))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(body), `"user_id":"user-1"`)
}

func Test_RequireBearerToken_rejectsMissingHeader(t *testing.T) {
	validator := &stubValidator{result: nil}
	app := newAppWithMiddleware(validator)

	req := httptest.NewRequest("GET", "/ping", nil)

	resp, err := app.Test(req, int(5*time.Second/time.Millisecond))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func Test_RequireBearerToken_rejectsMalformedHeader(t *testing.T) {
	validator := &stubValidator{result: nil}
	app := newAppWithMiddleware(validator)

	req := httptest.NewRequest("GET", "/ping", nil)
	req.Header.Set("Authorization", "Basic notbearer")

	resp, err := app.Test(req, int(5*time.Second/time.Millisecond))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func Test_RequireBearerToken_failsClosedOnIamError(t *testing.T) {
	// Simulate iam-svc returning ErrUnauthorized (or being unreachable —
	// the adapter already maps both to ErrUnauthorized by design). The
	// middleware must surface 401, never 500 and never allow.
	iamErr := apperrors.ErrUnauthorized
	validator := &stubValidator{err: iamErr}
	app := newAppWithMiddleware(validator)

	req := httptest.NewRequest("GET", "/ping", nil)
	req.Header.Set("Authorization", "Bearer ignored.token")

	resp, err := app.Test(req, int(5*time.Second/time.Millisecond))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}
