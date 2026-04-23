package config

import (
	"github.com/spf13/viper"
)

// bindEnvironmentVariables sets up environment variable mappings for all config fields.
func bindEnvironmentVariables() {
	// App config

	viper.BindEnv("app.name", "APP_NAME")

	// API config

	viper.BindEnv("api.rest.host", "API_REST_HOST")
	viper.BindEnv("api.rest.port", "API_REST_PORT")
	viper.BindEnv("api.metrics.enabled", "API_METRICS_ENABLED")

	// External backends (one block per backend the gateway calls).

	// iam-svc is gRPC-only (per ADR 0009 / BL-IAM-018) — no REST address.
	viper.BindEnv("external.iam_svc.name", "IAM_SVC_NAME")
	viper.BindEnv("external.iam_svc.grpc_target", "IAM_SVC_GRPC_TARGET")

	// catalog-svc is gRPC-only (per ADR 0009 / BL-REFACTOR-001) — no REST address.
	viper.BindEnv("external.catalog_svc.name", "CATALOG_SVC_NAME")
	viper.BindEnv("external.catalog_svc.grpc_target", "CATALOG_SVC_GRPC_TARGET")

	// The seven pure-scaffold backends are gRPC-only (per ADR 0009 /
	// BL-REFACTOR-002..008) — no REST address binding. grpc_target stays
	// available for future gRPC adapters.
	viper.BindEnv("external.booking_svc.name", "BOOKING_SVC_NAME")
	viper.BindEnv("external.booking_svc.grpc_target", "BOOKING_SVC_GRPC_TARGET")

	viper.BindEnv("external.jamaah_svc.name", "JAMAAH_SVC_NAME")
	viper.BindEnv("external.jamaah_svc.grpc_target", "JAMAAH_SVC_GRPC_TARGET")

	viper.BindEnv("external.payment_svc.name", "PAYMENT_SVC_NAME")
	viper.BindEnv("external.payment_svc.grpc_target", "PAYMENT_SVC_GRPC_TARGET")

	viper.BindEnv("external.visa_svc.name", "VISA_SVC_NAME")
	viper.BindEnv("external.visa_svc.grpc_target", "VISA_SVC_GRPC_TARGET")

	viper.BindEnv("external.ops_svc.name", "OPS_SVC_NAME")
	viper.BindEnv("external.ops_svc.grpc_target", "OPS_SVC_GRPC_TARGET")

	viper.BindEnv("external.logistics_svc.name", "LOGISTICS_SVC_NAME")
	viper.BindEnv("external.logistics_svc.grpc_target", "LOGISTICS_SVC_GRPC_TARGET")

	// finance-svc still exposes REST (BL-IAM-019 / S1-E-14 will remove it).
	viper.BindEnv("external.finance_svc.name", "FINANCE_SVC_NAME")
	viper.BindEnv("external.finance_svc.address", "FINANCE_SVC_ADDRESS")
	viper.BindEnv("external.finance_svc.grpc_target", "FINANCE_SVC_GRPC_TARGET")

	viper.BindEnv("external.crm_svc.name", "CRM_SVC_NAME")
	viper.BindEnv("external.crm_svc.grpc_target", "CRM_SVC_GRPC_TARGET")

	// Otel tracer config

	viper.BindEnv("otel_tracer.name", "OTEL_TRACER_NAME")
	viper.BindEnv("otel_tracer.endpoint", "OTEL_TRACER_ENDPOINT")
}
