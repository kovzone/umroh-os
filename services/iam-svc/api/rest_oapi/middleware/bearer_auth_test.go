package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"iam-svc/util/token"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// newTestApp builds a Fiber app with the BL-IAM-001 error envelope so
// bearer-auth failures produce the same `{error:{code,message}}` shape that
// handlers downstream rely on.
func newTestApp(maker token.Maker) *fiber.App {
	app := fiber.New()
	app.Use(ErrorHandler())
	app.Get("/protected", RequireBearerToken(maker), func(c *fiber.Ctx) error {
		payload := c.Locals(PayloadKey).(*token.Payload)
		return c.JSON(fiber.Map{"user_id": payload.UserID})
	})
	return app
}

const testPasetoKey = "iam_svc_test_paseto_32_byte_key0"

func TestBearerAuth_ValidToken(t *testing.T) {
	maker, err := token.NewPasetoMaker(testPasetoKey)
	require.NoError(t, err)

	payload, err := token.NewPayload("user-1", "branch-1", []string{"admin"})
	require.NoError(t, err)
	signed, err := maker.CreateToken(payload, 15*time.Minute)
	require.NoError(t, err)

	app := newTestApp(maker)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signed)

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	var result map[string]string
	require.NoError(t, json.Unmarshal(body, &result))
	require.Equal(t, "user-1", result["user_id"])
}

func TestBearerAuth_MissingHeader(t *testing.T) {
	maker, err := token.NewPasetoMaker(testPasetoKey)
	require.NoError(t, err)

	app := newTestApp(maker)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(body), "UNAUTHORIZED")
	require.Contains(t, string(body), "authorization header is missing")
}

func TestBearerAuth_BadFormat(t *testing.T) {
	maker, err := token.NewPasetoMaker(testPasetoKey)
	require.NoError(t, err)

	app := newTestApp(maker)

	tests := []struct {
		name   string
		header string
	}{
		{"no scheme", "just-a-token"},
		{"wrong scheme", "Basic abc123"},
		{"empty token", "Bearer "},
		{"extra parts", "Bearer abc def"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", tt.header)

			resp, err := app.Test(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

			body, _ := io.ReadAll(resp.Body)
			require.Contains(t, string(body), "UNAUTHORIZED")
		})
	}
}

func TestBearerAuth_InvalidToken(t *testing.T) {
	maker, err := token.NewPasetoMaker(testPasetoKey)
	require.NoError(t, err)

	app := newTestApp(maker)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer totally-invalid-token")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(body), "UNAUTHORIZED")
}

func TestBearerAuth_ExpiredToken(t *testing.T) {
	maker, err := token.NewPasetoMaker(testPasetoKey)
	require.NoError(t, err)

	payload, err := token.NewPayload("user-1", "branch-1", []string{"admin"})
	require.NoError(t, err)
	signed, err := maker.CreateToken(payload, -time.Minute)
	require.NoError(t, err)

	app := newTestApp(maker)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signed)

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	require.Contains(t, string(body), "UNAUTHORIZED")
}
