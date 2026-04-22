package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/util/apperrors"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// stubValidator is an in-memory IamValidator for testing the middleware
// without a real iam-svc gRPC connection. Use result + err to simulate
// every failure mode the adapter can produce.
type stubValidator struct {
	result     *iam_grpc_adapter.ValidateTokenResult
	err        error
	capturedTk string // last token the middleware passed in
}

func (s *stubValidator) ValidateToken(_ context.Context, params *iam_grpc_adapter.ValidateTokenParams) (*iam_grpc_adapter.ValidateTokenResult, error) {
	s.capturedTk = params.AccessToken
	return s.result, s.err
}

// newTestApp wires the gateway error-envelope middleware + a single
// /protected route that demands Bearer and echoes the validated identity.
// Mirrors iam-svc's bearer_auth_test harness so the envelope shape stays
// consistent across services.
func newTestApp(v IamValidator) *fiber.App {
	app := fiber.New()
	app.Use(ErrorHandler())
	app.Get("/protected", RequireBearerToken(v), func(c *fiber.Ctx) error {
		id := c.Locals(IdentityKey).(*Identity)
		return c.JSON(fiber.Map{
			"user_id":    id.UserID,
			"branch_id":  id.BranchID,
			"session_id": id.SessionID,
			"roles":      id.Roles,
		})
	})
	return app
}

func TestBearerAuth_ValidToken(t *testing.T) {
	v := &stubValidator{
		result: &iam_grpc_adapter.ValidateTokenResult{
			UserID:    "user-1",
			BranchID:  "branch-1",
			SessionID: "sess-1",
			Roles:     []string{"admin"},
			ExpiresAt: time.Now().Add(15 * time.Minute),
		},
	}

	app := newTestApp(v)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer test-token-xyz")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, "test-token-xyz", v.capturedTk, "middleware must pass the raw token to the validator")

	body, _ := io.ReadAll(resp.Body)
	var result map[string]any
	require.NoError(t, json.Unmarshal(body, &result))
	require.Equal(t, "user-1", result["user_id"])
	require.Equal(t, "branch-1", result["branch_id"])
	require.Equal(t, "sess-1", result["session_id"])
}

func TestBearerAuth_MissingHeader(t *testing.T) {
	v := &stubValidator{}

	app := newTestApp(v)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(body), "UNAUTHORIZED")
	require.Contains(t, string(body), "authorization header is missing")
	require.Empty(t, v.capturedTk, "validator must not be called when header is absent")
}

func TestBearerAuth_MalformedHeader(t *testing.T) {
	v := &stubValidator{}
	app := newTestApp(v)

	tests := []struct {
		name   string
		header string
	}{
		{"no scheme prefix", "somerandomtoken"},
		{"wrong scheme", "Basic abc123"},
		{"too many parts", "Bearer aaa bbb"},
		{"empty header value", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusUnauthorized, resp.StatusCode, "header=%q should 401", tc.header)
		})
	}
	require.Empty(t, v.capturedTk, "validator must not be called for malformed headers")
}

// TestBearerAuth_IamRejects simulates the iam-svc Unauthenticated response
// (bearer expired, session revoked, etc.). The middleware must surface it as
// 401 UNAUTHORIZED.
func TestBearerAuth_IamRejects(t *testing.T) {
	v := &stubValidator{
		err: errors.Join(apperrors.ErrUnauthorized, errors.New("token expired")),
	}

	app := newTestApp(v)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer bad-token")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(body), "UNAUTHORIZED")
}

// TestBearerAuth_IamUnreachable simulates iam-svc being down (connection
// refused, timeout, etc.). The adapter maps this to ErrServiceUnavailable;
// the middleware bubbles it up; the error middleware renders 502 per F1-W7
// "fail closed, never default to allow."
func TestBearerAuth_IamUnreachable(t *testing.T) {
	v := &stubValidator{
		err: errors.Join(apperrors.ErrServiceUnavailable, errors.New("iam unreachable: connection refused")),
	}

	app := newTestApp(v)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer any-token")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadGateway, resp.StatusCode, "iam unreachable must produce 502 per F1-W7 fail-closed rule")
	body, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(body), "SERVICE_UNAVAILABLE")
}
