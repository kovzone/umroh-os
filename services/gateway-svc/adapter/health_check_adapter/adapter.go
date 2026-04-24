// Package health_check_adapter is the gateway's consumer-side wrapper around
// the standard grpc.health.v1.Health protocol (BL-MON-001) exposed by every
// downstream backend. It is NOT tied to a single backend — it holds one
// Health client per backend name and probes them all concurrently when the
// REST /v1/system/backends endpoint is hit by core-web's status page.
//
// Unlike the per-backend business adapters (catalog_grpc_adapter,
// finance_grpc_adapter, iam_grpc_adapter), this adapter does not care what
// RPCs each backend exposes — it talks only to `grpc.health.v1.Health` which
// is registered in the scaffold's cmd/server.go on every service. That's why
// it's also the right surface for the Svelte status page: one uniform
// protocol across every backend, independent of their business surface.
package health_check_adapter

import (
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
)

// Adapter holds one Health client per backend the gateway dials. Populated
// once at startup from cmd/start.go and queried on every /v1/system/backends
// request.
type Adapter struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	clients map[string]grpc_health_v1.HealthClient
}

// Backend is a (name, connection) pair accepted by NewAdapter. Lifted out so
// cmd/start.go can build the slice without depending on the internal client
// type.
type Backend struct {
	Name string
	Conn *grpc.ClientConn
}

// NewAdapter turns a slice of (name, conn) pairs into a Health-client map.
// Order is preserved via the input slice, but consumers that care about
// ordering should sort the result from CheckAll themselves.
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, backends []Backend) *Adapter {
	clients := make(map[string]grpc_health_v1.HealthClient, len(backends))
	for _, b := range backends {
		clients[b.Name] = grpc_health_v1.NewHealthClient(b.Conn)
	}
	return &Adapter{
		logger:  logger,
		tracer:  tracer,
		clients: clients,
	}
}

// Names returns the registered backend names in insertion order of the
// underlying map iteration. Callers that need deterministic ordering should
// sort the result or keep their own ordering.
func (a *Adapter) Names() []string {
	out := make([]string, 0, len(a.clients))
	for name := range a.clients {
		out = append(out, name)
	}
	return out
}
