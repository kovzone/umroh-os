package service

// BL-CAT-012 — catalog.updated webhook event stub.
//
// emitCatalogUpdated fires an HTTP POST to CATALOG_WEBHOOK_URL with a JSON
// body describing the mutation. Errors are logged but never propagated so
// that a failing webhook never blocks a catalog mutation.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"catalog-svc/util/logging"
)

// catalogWebhookPayload is the JSON body sent to the webhook endpoint.
type catalogWebhookPayload struct {
	Event      string `json:"event"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Timestamp  string `json:"timestamp"`
}

// emitCatalogUpdated fires a best-effort webhook notification after any
// package or departure mutation.
//
// eventKind  — e.g. "package.created", "package.updated", "package.deleted",
//              "departure.created", "departure.updated"
// entityID   — the mutated entity's ID
// entityType — "package" or "departure"
func (s *Service) emitCatalogUpdated(ctx context.Context, eventKind, entityID, entityType string) {
	logger := logging.LogWithTrace(ctx, s.logger)

	payload := catalogWebhookPayload{
		Event:      eventKind,
		EntityType: entityType,
		EntityID:   entityID,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}
	logger.Info().
		Str("event", eventKind).
		Str("entity_type", entityType).
		Str("entity_id", entityID).
		Msg("catalog.updated event")

	webhookURL := os.Getenv("CATALOG_WEBHOOK_URL")
	if webhookURL == "" {
		// No URL configured — log only, nothing to deliver.
		return
	}

	body, err := json.Marshal(payload)
	if err != nil {
		logger.Error().Err(err).Msg("catalog webhook: marshal payload failed")
		return
	}

	// Fire and forget: use a short deadline so a slow webhook does not
	// hold the request goroutine open indefinitely.
	reqCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		logger.Error().Err(err).Msg("catalog webhook: build request failed")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Warn().Err(err).Str("url", webhookURL).Msg("catalog webhook: delivery failed (ignored)")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		logger.Warn().
			Str("url", webhookURL).
			Str("status", fmt.Sprintf("%d", resp.StatusCode)).
			Msg("catalog webhook: non-2xx response (ignored)")
	}
}
