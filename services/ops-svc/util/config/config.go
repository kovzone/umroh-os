package config

import (
	"time"
)

// Config holds all configuration for the application.
//
// Post-ADR 0009: ops-svc is gRPC-only. Metrics push via OTLP to the
// OpenTelemetry Collector using the OtelTracer.Endpoint; there is no local
// /metrics HTTP endpoint and no opt-in knob.
type Config struct {
	App        App        `mapstructure:"app"`
	Api        Api        `mapstructure:"api"`
	Store      Store      `mapstructure:"store"`
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

// Otel tracer config

type OtelTracer struct {
	Name     string `mapstructure:"name"`
	Endpoint string `mapstructure:"endpoint"`
}
