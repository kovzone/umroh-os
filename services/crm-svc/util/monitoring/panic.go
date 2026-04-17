package monitoring

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// panicsTotal counts recovered panics. Panics in banking = pager; you want visibility even if the process survives.
	panicsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "panics_total",
			Help: "Total number of panics recovered by the panic recovery middleware. Use for alerting on unexpected panics.",
		},
	)
)

func init() {
	DefaultRegistry.MustRegister(panicsTotal)
}

// RecoveryMiddleware returns a Fiber middleware that recovers from panics, increments panics_total, then re-panics
// so the request still fails and error handlers can run. Mount it as early as possible (e.g. first after CORS).
func RecoveryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				panicsTotal.Inc()
				panic(r)
			}
		}()
		return c.Next()
	}
}
