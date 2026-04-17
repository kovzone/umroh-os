package finance_rest_adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// LivenessResult mirrors finance-svc's /system/live response envelope.
type LivenessResult struct {
	OK bool `json:"ok"`
}

// envelope is finance-svc's outer JSON shape: {"data": <T>}.
type livenessEnvelope struct {
	Data LivenessResult `json:"data"`
}

// GetSystemLive calls finance-svc's GET /system/live and returns the decoded result.
//
// Used by gateway-svc's /v1/iam/system/live proxy proof — demonstrates that
// the REST adapter pattern works end-to-end: span propagation via otelhttp,
// typed result, error wrapping consistent with the baseline gRPC adapters.
func (a *Adapter) GetSystemLive(ctx context.Context) (*LivenessResult, error) {
	const op = "finance_rest_adapter.GetSystemLive"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)

	url := a.baseURL + "/system/live"
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("http.method", http.MethodGet),
		attribute.String("http.url", url),
	)
	logger.Info().Str("op", op).Str("url", url).Msg("calling finance-svc")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		err = fmt.Errorf("build request: %w", err)
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		err = fmt.Errorf("call finance-svc: %w", err)
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("finance-svc returned status %d", resp.StatusCode)
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	var env livenessEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		err = fmt.Errorf("decode response: %w", err)
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return &env.Data, nil
}
