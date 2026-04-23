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
	viper.BindEnv("api.http.webhook_address", "API_HTTP_WEBHOOK_ADDRESS")

	// Store config

	viper.BindEnv("store.postgres.connection_string", "POSTGRES_CONNECTION_STRING")
	viper.BindEnv("store.postgres.pool.max_conns", "POSTGRES_MAX_CONNS")
	viper.BindEnv("store.postgres.pool.min_conns", "POSTGRES_MIN_CONNS")
	viper.BindEnv("store.postgres.pool.retry_max_attempts", "POSTGRES_POOL_RETRY_MAX_ATTEMPTS")
	viper.BindEnv("store.postgres.pool.retry_base_delay", "POSTGRES_POOL_RETRY_BASE_DELAY")
	viper.BindEnv("store.postgres.pool.retry_max_delay", "POSTGRES_POOL_RETRY_MAX_DELAY")

	// Gateway config
	// TODO(REAL_API_KEY): set MIDTRANS_SERVER_KEY and XENDIT_CALLBACK_TOKEN in production.

	viper.BindEnv("gateway.mock_gateway", "MOCK_GATEWAY")
	viper.BindEnv("gateway.midtrans_server_key", "MIDTRANS_SERVER_KEY")
	viper.BindEnv("gateway.midtrans_base_url", "MIDTRANS_BASE_URL")
	viper.BindEnv("gateway.xendit_callback_token", "XENDIT_CALLBACK_TOKEN")
	viper.BindEnv("gateway.xendit_secret_key", "XENDIT_SECRET_KEY")
	viper.BindEnv("gateway.xendit_base_url", "XENDIT_BASE_URL")

	// Downstream services config

	viper.BindEnv("services.booking_svc_addr", "BOOKING_SVC_ADDR")

	// Otel tracer config

	viper.BindEnv("otel_tracer.name", "OTEL_TRACER_NAME")
	viper.BindEnv("otel_tracer.endpoint", "OTEL_TRACER_ENDPOINT")
}
