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
	viper.BindEnv("api.grpc.address", "API_GRPC_ADDRESS")
	viper.BindEnv("api.metrics.enabled", "API_METRICS_ENABLED")

	// Store config

	viper.BindEnv("store.postgres.connection_string", "POSTGRES_CONNECTION_STRING")
	viper.BindEnv("store.postgres.pool.max_conns", "POSTGRES_MAX_CONNS")
	viper.BindEnv("store.postgres.pool.min_conns", "POSTGRES_MIN_CONNS")
	viper.BindEnv("store.postgres.pool.retry_max_attempts", "POSTGRES_POOL_RETRY_MAX_ATTEMPTS")
	viper.BindEnv("store.postgres.pool.retry_base_delay", "POSTGRES_POOL_RETRY_BASE_DELAY")
	viper.BindEnv("store.postgres.pool.retry_max_delay", "POSTGRES_POOL_RETRY_MAX_DELAY")

	// Token config

	viper.BindEnv("token.type", "TOKEN_TYPE")
	viper.BindEnv("token.key", "TOKEN_KEY")
	viper.BindEnv("token.access_duration", "TOKEN_ACCESS_DURATION")
	viper.BindEnv("token.refresh_duration", "TOKEN_REFRESH_DURATION")

	// TOTP config

	viper.BindEnv("totp.issuer", "TOTP_ISSUER")
	viper.BindEnv("totp.encryption_key", "TOTP_ENCRYPTION_KEY")

	// Otel tracer config

	viper.BindEnv("otel_tracer.name", "OTEL_TRACER_NAME")
	viper.BindEnv("otel_tracer.endpoint", "OTEL_TRACER_ENDPOINT")
}
