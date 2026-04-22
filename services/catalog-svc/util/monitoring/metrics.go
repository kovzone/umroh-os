// Package monitoring wires catalog-svc's runtime metrics via the OpenTelemetry
// metrics SDK, pushed over OTLP gRPC to the OpenTelemetry Collector. The
// Collector re-exports them on its Prometheus endpoint (see
// monitoring/otel-collector-config.yaml), which Prometheus scrapes.
//
// Per ADR 0009 catalog-svc is gRPC-only — there is no local /metrics HTTP
// endpoint. The service pushes; the Collector exposes. Service-name filtering
// in Grafana uses the `service_name` resource attribute attached here.
package monitoring

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// meterName is the instrumentation scope for every metric this service records.
const meterName = "catalog-svc/monitoring"

// InitMeter initializes the global MeterProvider with an OTLP gRPC exporter
// pushing to the OpenTelemetry Collector every 15 seconds. Returns a shutdown
// func to flush in-flight batches on process exit.
//
// Symmetric with util/tracing.InitTracer — same resource shape, same OTLP
// endpoint, same insecure transport for the in-cluster hop.
func InitMeter(serviceName, otlpEndpoint string) (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(serviceName)),
	)
	if err != nil {
		return nil, fmt.Errorf("new resource: %w", err)
	}

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(otlpEndpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("new otlp metric exporter: %w", err)
	}

	reader := sdkmetric.NewPeriodicReader(exporter,
		sdkmetric.WithInterval(15*time.Second),
	)

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(provider)

	return provider.Shutdown, nil
}

// Meter returns the service-scoped Meter used for every custom instrument in
// catalog-svc. Call after InitMeter; until then this returns a no-op meter
// from the global provider.
func Meter() metric.Meter {
	return otel.Meter(meterName)
}
