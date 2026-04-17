package monitoring

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var dbPoolStatsOnce sync.Once

// DBPoolStats holds connection pool metrics. Use with RegisterDBPoolStats to expose them on /metrics.
// Compatible with pgxpool.Stat (AcquiredConns, IdleConns, TotalConns).
type DBPoolStats struct {
	Acquired int32 // connections currently in use
	Idle     int32 // idle connections in the pool
	Total    int32 // total connections (acquired + idle + constructing, etc.)
}

// DBPoolStatsProvider returns current pool stats. Called periodically by the DB pool stats collector.
// From main, pass a function that calls pool.Stat() and returns the relevant fields.
type DBPoolStatsProvider func() DBPoolStats

var (
	dbConnectionsAcquired = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_acquired",
			Help: "Number of database connections currently in use (acquired from the pool). DB saturation shows up here early.",
		},
	)
	dbConnectionsIdle = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_idle",
			Help: "Number of idle database connections in the pool.",
		},
	)
	dbConnectionsTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_total",
			Help: "Total number of database connections in the pool (open).",
		},
	)
)

// RegisterDBPoolStats registers the DB pool gauges with the default registry (once) and starts a goroutine that
// periodically calls provider and updates the gauges. The goroutine runs until ctx is cancelled.
// Call this from main after creating the pool, e.g.:
//
//	go monitoring.RegisterDBPoolStats(ctx, func() monitoring.DBPoolStats {
//	    s := postgresPool.Stat()
//	    return monitoring.DBPoolStats{Acquired: s.AcquiredConns(), Idle: s.IdleConns(), Total: s.TotalConns()}
//	}, 10*time.Second)
func RegisterDBPoolStats(ctx context.Context, provider DBPoolStatsProvider, interval time.Duration) {
	dbPoolStatsOnce.Do(func() {
		DefaultRegistry.MustRegister(dbConnectionsAcquired, dbConnectionsIdle, dbConnectionsTotal)
	})

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s := provider()
				dbConnectionsAcquired.Set(float64(s.Acquired))
				dbConnectionsIdle.Set(float64(s.Idle))
				dbConnectionsTotal.Set(float64(s.Total))
			}
		}
	}()
}
