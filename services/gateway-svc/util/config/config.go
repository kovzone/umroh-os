package config

// Config holds all configuration for gateway-svc.
//
// gateway-svc is the edge REST proxy. It carries no DB. The External block
// names every backend service the gateway can call; one entry per backend.
type Config struct {
	App        App        `mapstructure:"app"`
	Api        Api        `mapstructure:"api"`
	External   External   `mapstructure:"external"`
	OtelTracer OtelTracer `mapstructure:"otel_tracer"`
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

type Api struct {
	Rest    Rest    `mapstructure:"rest"`
	Metrics Metrics `mapstructure:"metrics"`
}

// Metrics config (Prometheus). Opt-in: set enabled to true to expose /metrics.
type Metrics struct {
	Enabled bool `mapstructure:"enabled"`
}

// External — one entry per backend service the gateway calls over REST.

type External struct {
	IamSvc       BackendSvc `mapstructure:"iam_svc"`
	CatalogSvc   BackendSvc `mapstructure:"catalog_svc"`
	BookingSvc   BackendSvc `mapstructure:"booking_svc"`
	JamaahSvc    BackendSvc `mapstructure:"jamaah_svc"`
	PaymentSvc   BackendSvc `mapstructure:"payment_svc"`
	VisaSvc      BackendSvc `mapstructure:"visa_svc"`
	OpsSvc       BackendSvc `mapstructure:"ops_svc"`
	LogisticsSvc BackendSvc `mapstructure:"logistics_svc"`
	FinanceSvc   BackendSvc `mapstructure:"finance_svc"`
	CrmSvc       BackendSvc `mapstructure:"crm_svc"`
}

// BackendSvc is the address of one backend.
//
// Post BL-IAM-019 / S1-E-14 every backend is gRPC-only per ADR 0009; the
// REST `Address` field is gone. GrpcTarget is a plain host:port like
// "iam-svc:50051" used by the service's gRPC adapter. GrpcTarget is optional
// on backends whose gateway-side gRPC adapter hasn't been introduced yet;
// it becomes required when a card introduces one.
//
// Name is used for logging and span attribution.
type BackendSvc struct {
	Name       string `mapstructure:"name"`
	GrpcTarget string `mapstructure:"grpc_target"`
}

// Otel tracer config

type OtelTracer struct {
	Name     string `mapstructure:"name"`
	Endpoint string `mapstructure:"endpoint"`
}
