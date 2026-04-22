package monitoring

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/metric"
)

// DBPoolStats holds connection pool metrics. Use with RegisterDBPoolStats to
// expose them. Compatible with pgxpool.Stat (AcquiredConns, IdleConns,
// TotalConns).
type DBPoolStats struct {
	Acquired int32 // connections currently in use
	Idle     int32 // idle connections in the pool
	Total    int32 // total connections (acquired + idle + constructing, etc.)
}

// DBPoolStatsProvider returns current pool stats. The OTel SDK calls this once
// per collection cycle (driven by the PeriodicReader interval in InitMeter),
// so the callback stays in the SDK's goroutine — no ticker or cancellation
// needed here.
type DBPoolStatsProvider func() DBPoolStats

// RegisterDBPoolStats attaches three observable gauges (db_connections_acquired,
// db_connections_idle, db_connections_total) whose values are pulled from
// provider() on every SDK collection cycle. Call once from cmd/start.go
// after InitMeter has run.
func RegisterDBPoolStats(provider DBPoolStatsProvider) error {
	meter := Meter()

	acquired, err := meter.Int64ObservableGauge(
		"db_connections_acquired",
		metric.WithDescription("Connections currently in use (acquired from the pool). DB saturation shows up here early."),
	)
	if err != nil {
		return fmt.Errorf("new gauge db_connections_acquired: %w", err)
	}

	idle, err := meter.Int64ObservableGauge(
		"db_connections_idle",
		metric.WithDescription("Idle database connections in the pool."),
	)
	if err != nil {
		return fmt.Errorf("new gauge db_connections_idle: %w", err)
	}

	total, err := meter.Int64ObservableGauge(
		"db_connections_total",
		metric.WithDescription("Total database connections in the pool (open)."),
	)
	if err != nil {
		return fmt.Errorf("new gauge db_connections_total: %w", err)
	}

	_, err = meter.RegisterCallback(
		func(_ context.Context, o metric.Observer) error {
			s := provider()
			o.ObserveInt64(acquired, int64(s.Acquired))
			o.ObserveInt64(idle, int64(s.Idle))
			o.ObserveInt64(total, int64(s.Total))
			return nil
		},
		acquired, idle, total,
	)
	if err != nil {
		return fmt.Errorf("register dbpool callback: %w", err)
	}
	return nil
}
