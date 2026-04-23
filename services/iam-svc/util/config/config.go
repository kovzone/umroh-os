package config

import (
	"time"
)

// Config holds all configuration for the application.
type Config struct {
	App        App        `mapstructure:"app"`
	Api        Api        `mapstructure:"api"`
	Store      Store      `mapstructure:"store"`
	Token      Token      `mapstructure:"token"`
	Totp       Totp       `mapstructure:"totp"`
	OtelTracer OtelTracer `mapstructure:"otel_tracer"`
}

// App config

type App struct {
	Name string `mapstructure:"name"`
}

// API config — iam-svc is gRPC-only post BL-IAM-018 / S1-E-12. Metrics flow
// via OTLP push (see util/monitoring) — no local /metrics HTTP endpoint.

type Grpc struct {
	Address string `mapstructure:"address"`
}

type Api struct {
	Grpc Grpc `mapstructure:"grpc"`
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

// Token config
// Type selects the token scheme: "paseto" (default) or "jwt".
// Key must be exactly 32 bytes for PASETO (ChaCha20-Poly1305 symmetric key);
// at least 32 bytes for JWT (HS256 secret key).
// AccessDuration is the short-lived access-token lifetime (e.g. 15m).
// RefreshDuration is the opaque-refresh-token lifetime (e.g. 168h / 7d).

type Token struct {
	Type            string        `mapstructure:"type"`
	Key             string        `mapstructure:"key"`
	AccessDuration  time.Duration `mapstructure:"access_duration"`
	RefreshDuration time.Duration `mapstructure:"refresh_duration"`
}

// TOTP config
// Issuer is the label shown in the jamaah/staff authenticator app (e.g. "UmrohOS").
// EncryptionKey is exactly 32 bytes; used as the AES-256-GCM key that wraps
// iam.users.totp_secret at rest (F1 data-model requirement).

type Totp struct {
	Issuer        string `mapstructure:"issuer"`
	EncryptionKey string `mapstructure:"encryption_key"`
}

// Otel tracer config

type OtelTracer struct {
	Name     string `mapstructure:"name"`
	Endpoint string `mapstructure:"endpoint"`
}
