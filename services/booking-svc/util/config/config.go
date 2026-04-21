package config

import (
	"time"
)

// Config holds all configuration for the application.
//
// Pilot scaffold scope: app name, REST + gRPC ports, Postgres pool, OTel
// tracer. Token (auth) config returns with F1.5 when real login/refresh/logout
// handlers land.
type Config struct {
	App        App        `mapstructure:"app"`
	Api        Api        `mapstructure:"api"`
	Store      Store      `mapstructure:"store"`
	OtelTracer OtelTracer `mapstructure:"otel_tracer"`
	Iam        Iam        `mapstructure:"iam"`
}

// Iam config — how booking-svc reaches iam-svc's internal gRPC surface.
//
// The target is a plain host:port (no scheme); the BL-IAM-004 adapter dials
// with insecure credentials because the traffic stays inside the docker-compose
// network. TLS-terminated ingress lands with the gateway-svc hardening card.
type Iam struct {
	GrpcTarget string `mapstructure:"grpc_target"`
}

// App config

type App struct {
	Name string `mapstructure:"name"`
}

// API config

type Rest struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Grpc struct {
	Address string `mapstructure:"address"`
}

type Api struct {
	Rest    Rest    `mapstructure:"rest"`
	Grpc    Grpc    `mapstructure:"grpc"`
	Metrics Metrics `mapstructure:"metrics"`
}

// Metrics config (Prometheus). Opt-in: set enabled to true to expose /metrics.
type Metrics struct {
	Enabled bool `mapstructure:"enabled"`
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
