package config

import (
	"time"
)

// Config holds all configuration for the application.
//
// Post-ADR 0009: payment-svc is gRPC-only for inter-service calls.
// Exception: webhook endpoints are served via a separate HTTP listener (see Api.Http)
// because Midtrans/Xendit POST to public URLs directly.
// Metrics push via OTLP to the OTel Collector using OtelTracer.Endpoint.
type Config struct {
	App        App        `mapstructure:"app"`
	Api        Api        `mapstructure:"api"`
	Store      Store      `mapstructure:"store"`
	Gateway    Gateway    `mapstructure:"gateway"`
	Services   Services   `mapstructure:"services"`
	OtelTracer OtelTracer `mapstructure:"otel_tracer"`
}

// App config

type App struct {
	Name string `mapstructure:"name"`
}

// API config

type Grpc struct {
	Address string `mapstructure:"address"`
}

type Http struct {
	// WebhookAddress is the internal HTTP listener for Midtrans/Xendit webhook ingestion.
	// Not exposed to the internet — nginx/LB forwards only POST /v1/webhooks/* from gateway IP.
	// Default: 0.0.0.0:50065
	WebhookAddress string `mapstructure:"webhook_address"`
}

type Api struct {
	Grpc Grpc `mapstructure:"grpc"`
	Http Http `mapstructure:"http"`
}

// Gateway config — payment gateway credentials and behaviour flags.
// TODO(REAL_API_KEY): set MIDTRANS_SERVER_KEY and XENDIT_CALLBACK_TOKEN in production.

type Gateway struct {
	// MockGateway activates the in-process mock adapter (MOCK_GATEWAY=true).
	// Must be false in production. payment-svc refuses to start if ENV=production && MockGateway=true.
	MockGateway bool `mapstructure:"mock_gateway"`

	// MidtransServerKey is the Midtrans server key for VA issuance and webhook verification.
	// TODO(REAL_API_KEY): set via env var MIDTRANS_SERVER_KEY.
	MidtransServerKey string `mapstructure:"midtrans_server_key"`

	// MidtransBaseURL is the Midtrans API base URL.
	// Production: https://api.midtrans.com — Sandbox: https://api.sandbox.midtrans.com
	MidtransBaseURL string `mapstructure:"midtrans_base_url"`

	// XenditCallbackToken is the static Xendit callback token for webhook verification.
	// TODO(REAL_API_KEY): set via env var XENDIT_CALLBACK_TOKEN.
	XenditCallbackToken string `mapstructure:"xendit_callback_token"`

	// XenditSecretKey is the Xendit API secret key for VA issuance and status queries.
	// TODO(REAL_API_KEY): set via env var XENDIT_SECRET_KEY.
	XenditSecretKey string `mapstructure:"xendit_secret_key"`

	// XenditBaseURL is the Xendit API base URL.
	// Default: https://api.xendit.co
	XenditBaseURL string `mapstructure:"xendit_base_url"`
}

// Services config — addresses of downstream gRPC services.

type Services struct {
	// BookingSvcAddr is the gRPC address of booking-svc.
	// payment-svc calls MarkBookingPaid on this address after webhook processing.
	// Default: booking-svc:50051
	BookingSvcAddr string `mapstructure:"booking_svc_addr"`
}

// Store config

type PostgresPool struct {
	MaxConns         int           `mapstructure:"max_conns"`
	MinConns         int           `mapstructure:"min_conns"`
	RetryMaxAttempts int           `mapstructure:"retry_max_attempts"`
	RetryBaseDelay   time.Duration `mapstructure:"retry_base_delay"`
	RetryMaxDelay    time.Duration `mapstructure:"retry_max_delay"`
}

type Postgres struct {
	ConnectionString string       `mapstructure:"connection_string"`
	Pool             PostgresPool `mapstructure:"pool"`
}

type Store struct {
	Postgres Postgres `mapstructure:"postgres"`
}

// Otel tracer config

type OtelTracer struct {
	Name     string `mapstructure:"name"`
	Endpoint string `mapstructure:"endpoint"`
}
