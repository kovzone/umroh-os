package config

import (
	"github.com/spf13/viper"
)

// bindEnvironmentVariables sets up environment variable mappings for all config fields.
func bindEnvironmentVariables() {
	// App config

	viper.BindEnv("app.name", "APP_NAME")

	// API config

	viper.BindEnv("api.grpc.address", "API_GRPC_ADDRESS")

	// Store config

	viper.BindEnv("store.postgres.connection_string", "POSTGRES_CONNECTION_STRING")
	viper.BindEnv("store.postgres.pool.max_conns", "POSTGRES_MAX_CONNS")
	viper.BindEnv("store.postgres.pool.min_conns", "POSTGRES_MIN_CONNS")
	viper.BindEnv("store.postgres.pool.retry_max_attempts", "POSTGRES_POOL_RETRY_MAX_ATTEMPTS")
	viper.BindEnv("store.postgres.pool.retry_base_delay", "POSTGRES_POOL_RETRY_BASE_DELAY")
	viper.BindEnv("store.postgres.pool.retry_max_delay", "POSTGRES_POOL_RETRY_MAX_DELAY")

	// Otel tracer config

	viper.BindEnv("otel_tracer.name", "OTEL_TRACER_NAME")
	viper.BindEnv("otel_tracer.endpoint", "OTEL_TRACER_ENDPOINT")

	// Downstream service gRPC targets (S3-E-02 / S3-E-03)

	viper.BindEnv("logistics.grpc_target", "LOGISTICS_GRPC_TARGET")
	viper.BindEnv("finance.grpc_target", "FINANCE_GRPC_TARGET")
}
