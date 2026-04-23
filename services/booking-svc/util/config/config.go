package config

import (
	"time"
)

// Config holds all configuration for the application.
//
// Post-ADR 0009: booking-svc is gRPC-only. Metrics push via OTLP to the
// OpenTelemetry Collector using the OtelTracer.Endpoint; there is no local
// /metrics HTTP endpoint and no opt-in knob.
// S1-E-03 adds `catalog` section for the gRPC connection to catalog-svc
// used by draft booking creation (ReserveSeats / GetDeparture calls).
// S3-E-02 adds `logistics` section for the OnBookingPaid fan-out.
// S3-E-03 adds `finance` section for the OnPaymentReceived fan-out.
// S4-E-02 adds `crm` section for the CRM lead lifecycle fan-out.
type Config struct {
	App        App        `mapstructure:"app"`
	Api        Api        `mapstructure:"api"`
	Store      Store      `mapstructure:"store"`
	OtelTracer OtelTracer `mapstructure:"otel_tracer"`
	Iam        Iam        `mapstructure:"iam"`
	Catalog    Catalog    `mapstructure:"catalog"`
	Logistics  Logistics  `mapstructure:"logistics"`
	Finance    Finance    `mapstructure:"finance"`
	Crm        Crm        `mapstructure:"crm"`
}

// Catalog config — gRPC connection to catalog-svc for ReserveSeats / GetDeparture.
type Catalog struct {
	GrpcTarget string `mapstructure:"grpc_target"`
}

// Logistics config — gRPC connection to logistics-svc for OnBookingPaid fan-out (S3-E-02).
type Logistics struct {
	GrpcTarget string `mapstructure:"grpc_target"`
}

// Finance config — gRPC connection to finance-svc for OnPaymentReceived fan-out (S3-E-03).
type Finance struct {
	GrpcTarget string `mapstructure:"grpc_target"`
}

// Crm config — gRPC connection to crm-svc for lead lifecycle events (S4-E-02).
// Optional: empty GrpcTarget disables the dial; CRM fan-out is skipped.
type Crm struct {
	GrpcTarget string `mapstructure:"grpc_target"`
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
