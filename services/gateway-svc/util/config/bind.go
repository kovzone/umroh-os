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

	viper.BindEnv("external.iam_svc.name", "IAM_SVC_NAME")
	viper.BindEnv("external.iam_svc.address", "IAM_SVC_ADDRESS")

	viper.BindEnv("external.catalog_svc.name", "CATALOG_SVC_NAME")
	viper.BindEnv("external.catalog_svc.address", "CATALOG_SVC_ADDRESS")

	viper.BindEnv("external.booking_svc.name", "BOOKING_SVC_NAME")
	viper.BindEnv("external.booking_svc.address", "BOOKING_SVC_ADDRESS")

	viper.BindEnv("external.jamaah_svc.name", "JAMAAH_SVC_NAME")
	viper.BindEnv("external.jamaah_svc.address", "JAMAAH_SVC_ADDRESS")

	viper.BindEnv("external.payment_svc.name", "PAYMENT_SVC_NAME")
	viper.BindEnv("external.payment_svc.address", "PAYMENT_SVC_ADDRESS")

	viper.BindEnv("external.visa_svc.name", "VISA_SVC_NAME")
	viper.BindEnv("external.visa_svc.address", "VISA_SVC_ADDRESS")

	viper.BindEnv("external.ops_svc.name", "OPS_SVC_NAME")
	viper.BindEnv("external.ops_svc.address", "OPS_SVC_ADDRESS")

	viper.BindEnv("external.logistics_svc.name", "LOGISTICS_SVC_NAME")
	viper.BindEnv("external.logistics_svc.address", "LOGISTICS_SVC_ADDRESS")

	viper.BindEnv("external.finance_svc.name", "FINANCE_SVC_NAME")
	viper.BindEnv("external.finance_svc.address", "FINANCE_SVC_ADDRESS")

	viper.BindEnv("external.crm_svc.name", "CRM_SVC_NAME")
	viper.BindEnv("external.crm_svc.address", "CRM_SVC_ADDRESS")

	// Otel tracer config

	viper.BindEnv("otel_tracer.name", "OTEL_TRACER_NAME")
	viper.BindEnv("otel_tracer.endpoint", "OTEL_TRACER_ENDPOINT")
}
