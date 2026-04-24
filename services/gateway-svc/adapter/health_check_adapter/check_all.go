package health_check_adapter

import (
	"context"
	"sort"
	"sync"
	"time"

	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

// perBackendTimeout caps a single Health.Check call so one slow backend
// cannot stall the whole /v1/system/backends response. 2s is generous enough
// for a healthy probe + a little network jitter; unhealthy backends surface
// as a DeadlineExceeded error which the handler renders as StatusUnknown +
// a non-empty Error string.
const perBackendTimeout = 2 * time.Second

// CheckAll probes every registered backend's grpc.health.v1.Health.Check
// concurrently and returns the results sorted by backend name (so the Svelte
// status grid renders in a stable order across refreshes). The total wall
// time is bounded by perBackendTimeout, not by the number of backends,
// because probes fan out in parallel.
//
// A Health.Check that returns an error (backend unreachable, timeout,
// service not registered) becomes BackendStatus{Status: StatusUnknown,
// Error: <message>}. The handler treats that as "not ok" on the UI without
// conflating it with an explicit NOT_SERVING signal from a reachable
// backend — both render as "fail" today but keep the distinction on the
// wire for future dashboards.
func (a *Adapter) CheckAll(ctx context.Context) []BackendStatus {
	const op = "health_check_adapter.Adapter.CheckAll"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.Int("backend_count", len(a.clients)),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	out := make([]BackendStatus, 0, len(a.clients))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for name, client := range a.clients {
		wg.Add(1)
		go func(name string, client grpc_health_v1.HealthClient) {
			defer wg.Done()

			probeCtx, cancel := context.WithTimeout(ctx, perBackendTimeout)
			defer cancel()

			resp, err := client.Check(probeCtx, &grpc_health_v1.HealthCheckRequest{})
			status := BackendStatus{Name: name}
			if err != nil {
				status.Status = StatusUnknown
				status.Error = err.Error()
			} else {
				switch resp.GetStatus() {
				case grpc_health_v1.HealthCheckResponse_SERVING:
					status.Status = StatusServing
				case grpc_health_v1.HealthCheckResponse_NOT_SERVING:
					status.Status = StatusNotServing
				default:
					status.Status = StatusUnknown
				}
			}

			mu.Lock()
			out = append(out, status)
			mu.Unlock()
		}(name, client)
	}

	wg.Wait()

	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })

	logger.Info().Int("backend_count", len(out)).Msg("")
	span.SetStatus(codes.Ok, "success")
	return out
}
