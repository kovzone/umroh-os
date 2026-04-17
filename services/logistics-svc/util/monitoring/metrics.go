package monitoring

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// DefaultRegistry is the registry used for all HTTP and custom metrics.
// Use it when registering additional metrics (e.g. business counters).
var DefaultRegistry = prometheus.NewRegistry()

var (
	// httpRequestsTotal counts API requests by method, path, and status.
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of API requests received, by method, path and status. Use for throughput and SLA monitoring.",
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestDurationSeconds records request duration in seconds by method and path.
	httpRequestDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of API request handling in seconds, by method and path. Use for latency SLOs and performance monitoring.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// httpRequestsInFlight is the number of requests currently being handled (gauge). Early signal for saturation.
	httpRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed. Use as an early signal for saturation.",
		},
	)

	// httpResponseSizeBytes records response body size in bytes by method and path.
	httpResponseSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP response bodies in bytes, by method and path. Use for payload and bandwidth monitoring.",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8), // 100, 1k, 10k, 100k, 1M, 10M, 100M, 1G
		},
		[]string{"method", "path"},
	)

	// http4xxResponsesTotal counts client-error responses (4xx). Cheap alert signal.
	http4xxResponsesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_4xx_responses_total",
			Help: "Total number of HTTP 4xx (client error) responses. Use for alerting on client errors.",
		},
	)

	// http5xxResponsesTotal counts server-error responses (5xx). Cheap alert signal.
	http5xxResponsesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_5xx_responses_total",
			Help: "Total number of HTTP 5xx (server error) responses. Use for alerting on server errors.",
		},
	)
)

func init() {
	DefaultRegistry.MustRegister(
		httpRequestsTotal,
		httpRequestDurationSeconds,
		httpRequestsInFlight,
		httpResponseSizeBytes,
		http4xxResponsesTotal,
		http5xxResponsesTotal,
	)
	// Go runtime and process metrics (goroutines, heap, GC, CPU, RSS, open FDs).
	DefaultRegistry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
}

// Middleware returns a Fiber middleware that records request count, duration, in-flight, response size, and 4xx/5xx.
// Mount it early (e.g. after CORS) so all API traffic is measured.
func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		method := strings.Clone(c.Method())

		// Use c.Path() for full detail. Use strings.Clone() to force a deep copy
		// and prevent memory corruption from Fiber's buffer reuse.
		path := strings.Clone(c.Path())

		httpRequestsInFlight.Inc()
		defer func() { httpRequestsInFlight.Dec() }()

		err := c.Next()

		status := c.Response().StatusCode()
		statusStr := strconv.Itoa(status)
		duration := time.Since(start).Seconds()
		responseSize := float64(len(c.Response().Body()))

		httpRequestsTotal.WithLabelValues(method, path, statusStr).Inc()
		httpRequestDurationSeconds.WithLabelValues(method, path).Observe(duration)
		httpResponseSizeBytes.WithLabelValues(method, path).Observe(responseSize)

		if status >= 500 {
			http5xxResponsesTotal.Inc()
		} else if status >= 400 {
			http4xxResponsesTotal.Inc()
		}

		return err
	}
}

// Handler returns an http.Handler that serves Prometheus metrics in text format
// for the default registry. Mount it at GET /metrics (e.g. via adaptor.HTTPHandler).
func Handler() http.Handler {
	return promhttp.HandlerFor(DefaultRegistry, promhttp.HandlerOpts{})
}
