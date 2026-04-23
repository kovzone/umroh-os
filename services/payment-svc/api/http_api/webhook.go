// Package http_api provides a minimal HTTP server for webhook ingestion.
//
// # Architecture decision: direct HTTP for webhooks (MVP)
//
// Midtrans and Xendit POST signed payloads to a public URL. Two options exist:
//
//   A. gateway-svc proxies the webhook via gRPC to payment-svc.ProcessWebhook.
//   B. payment-svc exposes a second HTTP listener (this package) on a separate
//      port (50065); nginx/load-balancer forwards only POST /v1/webhooks/* from
//      the gateway's IP range to this port.
//
// Option B is chosen for MVP because:
//   - It avoids adding a gRPC round-trip on the 500ms webhook critical path.
//   - It does not require modifying gateway-svc (task constraint).
//   - The port is NOT exposed to the internet; it is behind nginx with IP-allow.
//   - If gateway-svc proxying is desired later, ProcessWebhook is also a gRPC
//     method — no business logic duplication needed.
//
// This HTTP server runs in a goroutine started by cmd/start.go alongside the
// existing gRPC server.

package http_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"payment-svc/service"
	"payment-svc/util/logging"

	"github.com/rs/zerolog"
)

// WebhookHandler handles POST /v1/webhooks/{gateway} requests.
type WebhookHandler struct {
	logger      *zerolog.Logger
	svc         service.IPaymentService
	mockEnabled bool
}

// NewWebhookHandler creates a WebhookHandler.
// mockEnabled should be true when MOCK_GATEWAY=true.
func NewWebhookHandler(logger *zerolog.Logger, svc service.IPaymentService, mockEnabled bool) *WebhookHandler {
	return &WebhookHandler{
		logger:      logger,
		svc:         svc,
		mockEnabled: mockEnabled,
	}
}

// RegisterRoutes mounts the webhook routes on the given mux.
// The mock trigger endpoint is only registered when mockEnabled=true (MOCK_GATEWAY=true).
func (h *WebhookHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/v1/webhooks/midtrans", h.handleWebhook("midtrans"))
	mux.HandleFunc("/v1/webhooks/xendit", h.handleWebhook("xendit"))
	if h.mockEnabled {
		mux.HandleFunc("/v1/webhooks/mock", h.handleWebhook("mock"))
		mux.HandleFunc("/v1/webhooks/mock/trigger", h.handleMockTrigger())
	}
	// Internal liveness probe for the HTTP listener.
	mux.HandleFunc("/internal/healthz", h.handleInternalHealthz())
}

// handleWebhook returns an http.HandlerFunc for the given gateway.
func (h *WebhookHandler) handleWebhook(gateway string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http_api.WebhookHandler.handleWebhook"

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()
		logger := logging.LogWithTrace(ctx, h.logger)

		// Read raw body (max 1 MB).
		body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if err != nil {
			logger.Error().Err(err).Str("op", op).Str("gateway", gateway).Msg("read body")
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Extract gateway-specific signature header.
		signature := webhookSignatureHeader(gateway, r)

		result, err := h.svc.ProcessWebhookEvent(ctx, &service.WebhookEventParams{
			Gateway:   gateway,
			Payload:   body,
			Signature: signature,
		})
		if err != nil {
			// Signature failure → 401 so gateway stops retrying bad-signature payloads.
			code := http.StatusInternalServerError
			if isUnauthorizedErr(err) {
				code = http.StatusUnauthorized
			}
			logger.Error().Err(err).Str("op", op).Str("gateway", gateway).Int("code", code).Msg("webhook processing failed")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"code":    "internal_error",
					"message": err.Error(),
				},
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":         true,
			"replayed":   result.Replayed,
			"invoice_id": result.InvoiceID,
			"status":     result.NewStatus,
		})
	}
}

// webhookSignatureHeader extracts the gateway-specific signature header value.
// Per S2-J-02 contract:
//   - Midtrans: X-Callback-Token (SHA512 HMAC)
//   - Xendit:   X-CALLBACK-TOKEN (static token comparison)
func webhookSignatureHeader(gateway string, r *http.Request) string {
	switch gateway {
	case "midtrans":
		return r.Header.Get("X-Callback-Token")
	case "xendit":
		return r.Header.Get("X-CALLBACK-TOKEN")
	default:
		return r.Header.Get("X-Webhook-Signature")
	}
}

// isUnauthorizedErr returns true if the error indicates a signature failure.
func isUnauthorizedErr(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return containsString(s, "unauthorized") ||
		containsString(s, "signature") ||
		containsString(s, "invalid_signature") ||
		containsString(s, "callback token")
}

func containsString(s, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Mock trigger endpoint (MOCK_GATEWAY=true only) — S2-J-04
// ---------------------------------------------------------------------------

// mockTriggerRequest is the body for POST /v1/webhooks/mock/trigger.
// Simulates a payment event without calling a real gateway (dev/test only).
type mockTriggerRequest struct {
	InvoiceID string  `json:"invoice_id"`
	Amount    float64 `json:"amount"`
	// Status: "payment_received" | "settlement_received"
	Status string `json:"status"`
}

// handleMockTrigger returns an HTTP handler that fabricates a synthetic
// webhook payload and runs it through the exact same webhook processing pipeline
// as a real Midtrans/Xendit webhook (per S2-J-04).
//
// Only registered when MOCK_GATEWAY=true. The endpoint is not reachable in
// production because gateway-svc returns 404 for /v1/webhooks/mock/trigger when
// MOCK_GATEWAY=false (defence-in-depth: payment-svc also checks mockEnabled).
func (h *WebhookHandler) handleMockTrigger() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http_api.WebhookHandler.handleMockTrigger"

		if !h.mockEnabled {
			// Defence-in-depth: should not be reachable, but guard anyway.
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"code":    "mock_gateway_disabled",
					"message": "MOCK_GATEWAY is not enabled on this instance",
				},
			})
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()
		logger := logging.LogWithTrace(ctx, h.logger)

		body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var req mockTriggerRequest
		if err := json.Unmarshal(body, &req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"code":    "validation_failed",
					"message": "invalid JSON body: " + err.Error(),
				},
			})
			return
		}

		if req.InvoiceID == "" || req.Amount <= 0 ||
			(req.Status != "payment_received" && req.Status != "settlement_received") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"code":    "validation_failed",
					"message": "invoice_id required; amount must be > 0; status must be payment_received or settlement_received",
				},
			})
			return
		}

		// Build a synthetic mock webhook payload that matches the "mock" gateway parser.
		// gateway_txn_id is unique per call so multiple triggers on the same invoice are
		// recorded as separate payment events (matching the real webhook dedup behaviour).
		txnID := fmt.Sprintf("MOCK-PMT-%s-%d", req.InvoiceID, time.Now().UnixNano())
		syntheticPayload, _ := json.Marshal(map[string]interface{}{
			"gateway_txn_id": txnID,
			"gateway_va_id":  req.InvoiceID, // idempotency_key = invoice.id (set at VA creation)
			"amount":         req.Amount,
			"status":         req.Status,
			"mock_trigger":   true,
		})

		// Run through the exact same pipeline as a real webhook (steps 3–8 of S2-J-02).
		// Mock adapter's VerifyWebhookSignature always returns nil.
		result, err := h.svc.ProcessWebhookEvent(ctx, &service.WebhookEventParams{
			Gateway:   "mock",
			Payload:   syntheticPayload,
			Signature: "", // mock adapter accepts any signature
		})
		if err != nil {
			code := http.StatusInternalServerError
			if isUnauthorizedErr(err) {
				code = http.StatusUnauthorized
			}
			logger.Error().Err(err).Str("op", op).Int("code", code).Msg("mock trigger failed")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(code)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error": map[string]string{
					"code":    "internal_error",
					"message": err.Error(),
				},
			})
			return
		}

		logger.Info().Str("op", op).
			Str("invoice_id", req.InvoiceID).
			Str("txn_id", txnID).
			Float64("amount", req.Amount).
			Bool("replayed", result.Replayed).
			Str("new_status", result.NewStatus).
			Msg("mock trigger processed")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"payment_event_txn_id": txnID,
			"invoice_status":       result.NewStatus,
			"replayed":             result.Replayed,
		})
	}
}

// handleInternalHealthz is a lightweight liveness probe for the HTTP webhook listener.
func (h *WebhookHandler) handleInternalHealthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}
}
