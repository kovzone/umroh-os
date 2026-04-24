package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gateway-svc/adapter/iam_grpc_adapter"
	"gateway-svc/util/apperrors"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// stubChecker is an in-memory PermissionChecker for testing the middleware
// without a real iam-svc gRPC connection. Configure result + err to simulate
// every adapter outcome; captured* fields assert the middleware forwarded
// the correct tuple + user id.
type stubChecker struct {
	result           *iam_grpc_adapter.CheckPermissionResult
	err              error
	capturedUserID   string
	capturedResource string
	capturedAction   string
	capturedScope    string
	capturedCtxOk    bool
	callCount        int
}

func (s *stubChecker) CheckPermission(ctx context.Context, params *iam_grpc_adapter.CheckPermissionParams) (*iam_grpc_adapter.CheckPermissionResult, error) {
	s.callCount++
	s.capturedUserID = params.UserID
	s.capturedResource = params.Resource
	s.capturedAction = params.Action
	s.capturedScope = params.Scope
	s.capturedCtxOk = ctx != nil
	return s.result, s.err
}

// newPermTestApp wires the gateway error-envelope middleware + a protected
// route pre-populated with an Identity in locals (simulating RequireBearerToken
// upstream). Matches the shape bearer_auth_test.go uses so envelope semantics
// stay consistent.
func newPermTestApp(checker PermissionChecker, identity *Identity, resource, action, scope string) *fiber.App {
	app := fiber.New()
	app.Use(ErrorHandler())
	app.Get("/protected",
		// Pre-mount identity-injector simulating what RequireBearerToken does in prod.
		func(c *fiber.Ctx) error {
			if identity != nil {
				c.Locals(IdentityKey, identity)
			}
			return c.Next()
		},
		RequirePermission(checker, resource, action, scope),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"ok": true})
		},
	)
	return app
}

// -------------------------------------------------------------------------
// Case 1 — allow. Checker returns {Allowed: true} → 200 + handler runs.
// -------------------------------------------------------------------------
func TestRequirePermission_Allow(t *testing.T) {
	c := &stubChecker{result: &iam_grpc_adapter.CheckPermissionResult{Allowed: true}}
	identity := &Identity{UserID: "user-1", BranchID: "branch-1", Roles: []string{"finance_admin"}}

	app := newPermTestApp(c, identity, "journal_entry", "read", "global")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 1, c.callCount)
	require.True(t, c.capturedCtxOk, "middleware must pass a non-nil context")
}

// -------------------------------------------------------------------------
// Case 2 — deny. Checker returns {Allowed: false} → 403 FORBIDDEN.
// The error envelope must leak the tuple so the ops loki search is usable;
// user-facing message stays sanitised.
// -------------------------------------------------------------------------
func TestRequirePermission_Deny(t *testing.T) {
	c := &stubChecker{result: &iam_grpc_adapter.CheckPermissionResult{Allowed: false}}
	identity := &Identity{UserID: "user-1", BranchID: "branch-1", Roles: []string{"cs_agent"}}

	app := newPermTestApp(c, identity, "journal_entry", "read", "global")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(body), "FORBIDDEN")
	require.Equal(t, "user-1", c.capturedUserID, "checker must receive identity.UserID")
}

// -------------------------------------------------------------------------
// Case 3 — missing identity in locals. Middleware chain was misconfigured
// (RequirePermission mounted without RequireBearerToken upstream). Bug
// guard: surface ErrInternal → 500, do not panic.
// -------------------------------------------------------------------------
func TestRequirePermission_MissingIdentity(t *testing.T) {
	c := &stubChecker{}

	app := newPermTestApp(c, nil, "journal_entry", "read", "global")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	require.Equal(t, 0, c.callCount, "checker must not be called when identity is absent")

	body, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(body), "INTERNAL_ERROR")
}

// -------------------------------------------------------------------------
// Case 4 — iam rejects with PermissionDenied. Adapter maps that to
// apperrors.ErrForbidden; middleware bubbles; envelope renders 403.
// -------------------------------------------------------------------------
func TestRequirePermission_IamReturnsForbidden(t *testing.T) {
	c := &stubChecker{
		err: errors.Join(apperrors.ErrForbidden, errors.New("permission denied by iam")),
	}
	identity := &Identity{UserID: "user-1", BranchID: "branch-1"}

	app := newPermTestApp(c, identity, "journal_entry", "read", "global")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(body), "FORBIDDEN")
}

// -------------------------------------------------------------------------
// Case 5 — iam unreachable. Adapter maps transport failure to
// ErrServiceUnavailable; middleware bubbles; envelope renders 502.
// This is the F1-W7 fail-closed rule: ambiguous upstream state never
// defaults to allow.
// -------------------------------------------------------------------------
func TestRequirePermission_IamUnreachable(t *testing.T) {
	c := &stubChecker{
		err: errors.Join(apperrors.ErrServiceUnavailable, errors.New("iam unreachable: connection refused")),
	}
	identity := &Identity{UserID: "user-1", BranchID: "branch-1"}

	app := newPermTestApp(c, identity, "journal_entry", "read", "global")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadGateway, resp.StatusCode, "iam unreachable must produce 502 per F1-W7")
	body, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(body), "SERVICE_UNAVAILABLE")
}

// -------------------------------------------------------------------------
// Case 6 — tuple forwarded correctly. The factory args (resource, action,
// scope) must reach the checker unchanged on each invocation.
// -------------------------------------------------------------------------
func TestRequirePermission_TupleForwarded(t *testing.T) {
	c := &stubChecker{result: &iam_grpc_adapter.CheckPermissionResult{Allowed: true}}
	identity := &Identity{UserID: "user-7"}

	app := newPermTestApp(c, identity, "journal_entry", "write", "branch")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	_, err := app.Test(req)
	require.NoError(t, err)

	require.Equal(t, "journal_entry", c.capturedResource)
	require.Equal(t, "write", c.capturedAction)
	require.Equal(t, "branch", c.capturedScope)
}

// -------------------------------------------------------------------------
// Case 7 — UserID comes from identity, NOT from a factory argument.
// Guards against a regression where a future refactor mistakenly plumbs
// UserID through a closure.
// -------------------------------------------------------------------------
func TestRequirePermission_UserIDFromIdentity(t *testing.T) {
	c := &stubChecker{result: &iam_grpc_adapter.CheckPermissionResult{Allowed: true}}
	identity := &Identity{UserID: "the-real-user-id-42", BranchID: "b-1"}

	app := newPermTestApp(c, identity, "journal_entry", "read", "global")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	_, err := app.Test(req)
	require.NoError(t, err)

	require.Equal(t, "the-real-user-id-42", c.capturedUserID)
}

// -------------------------------------------------------------------------
// Case 8 — iam returns a generic internal error. Neither the F1-W7
// fail-closed path nor an authz decision; surface as-is so the gateway
// error middleware renders 500. Confirms we do NOT treat unknown errors
// as "deny" (which would be a silent authz decision) or "allow" (which
// would be a security leak).
// -------------------------------------------------------------------------
func TestRequirePermission_IamReturnsUnknownError(t *testing.T) {
	c := &stubChecker{
		err: errors.Join(apperrors.ErrInternal, errors.New("iam internal error")),
	}
	identity := &Identity{UserID: "user-1", BranchID: "b-1"}

	app := newPermTestApp(c, identity, "journal_entry", "read", "global")
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	var envelope map[string]any
	require.NoError(t, json.Unmarshal(body, &envelope))
	errObj, ok := envelope["error"].(map[string]any)
	require.True(t, ok, "error envelope must be present")
	require.Equal(t, "INTERNAL_ERROR", errObj["code"])
}
