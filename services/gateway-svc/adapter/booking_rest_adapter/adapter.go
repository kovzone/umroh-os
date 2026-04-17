// Package booking_rest_adapter is gateway-svc's REST client into booking-svc.
//
// It mirrors the baseline gRPC adapter shape (see
// baseline/go-backend-template/demo-svc/adapter/demo_grpc_adapter/) but
// transports over REST instead of gRPC: an *http.Client wrapped with otelhttp
// for trace propagation, and a baseURL string in place of *grpc.ClientConn.
//
// Each typed method (one per topic file — system.go, users.go later, etc.)
// follows the same span+log+errwrap pattern as the gRPC variant. Switching to
// gRPC later is a one-sided change inside this package — the gateway service
// layer keeps calling the same methods.
package booking_rest_adapter

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

// Adapter wraps an HTTP client targeting booking-svc.
type Adapter struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	client  *http.Client
	baseURL string
}

// NewAdapter creates a REST adapter for booking-svc.
//
// The HTTP client is wrapped with otelhttp.NewTransport so the trace context
// from the inbound gateway request propagates over the wire to booking-svc as a
// W3C Trace Context header. booking-svc's Fiber middleware picks it up and
// continues the trace, so a single Tempo trace covers the full call chain.
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, baseURL string) *Adapter {
	httpClient := &http.Client{
		Timeout:   5 * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	return &Adapter{
		logger:  logger,
		tracer:  tracer,
		client:  httpClient,
		baseURL: baseURL,
	}
}
