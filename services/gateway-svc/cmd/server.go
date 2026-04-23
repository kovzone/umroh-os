package main

import (
	"fmt"
	"log"
	"os"

	"gateway-svc/api/rest_oapi"
	"gateway-svc/api/rest_oapi/middleware"
	"gateway-svc/util/monitoring"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// runRestServer runs the REST server using the OpenAPI-generated routes and handler.
// gateway-svc currently exposes own probes, the finance liveness proxy (interim
// REST adapter, retires with BL-IAM-019 / S1-E-14), the BL-GTW-002 public
// catalog read, and — post BL-IAM-018 / S1-E-12 — the client-facing auth
// surface backed by iam-svc gRPC with the bearer middleware mounted on every
// protected route. As each remaining REST adapter graduates to gRPC-only, the
// corresponding proxy disappears.
//
// The api parameter is typed as *rest_oapi.Server (not the interface) so we
// can call RequireBearerToken() to mount the bearer middleware on protected
// route groups.
func runRestServer(port int, api *rest_oapi.Server, metricsEnabled bool, serviceName string) {
	app := fiber.New()

	// CORS — gateway is the edge, accept cross-origin from any browser client.
	corsConfig := cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}
	app.Use(cors.New(corsConfig))

	// OpenTelemetry — start an inbound span for every request so cross-service
	// traces initiated by an upstream caller (via W3C traceparent) continue in
	// this process's spans instead of starting a new trace. gateway-svc is the
	// trace origin for edge-initiated requests and a trace continuer when
	// something (e.g. a frontend with OTel instrumentation) already propagates.
	// Span names are prefixed with the service name so a multi-service trace
	// in Tempo is easy to scan (e.g. "gateway-svc GET /v1/iam/system/live").
	app.Use(otelfiber.Middleware(
		otelfiber.WithSpanNameFormatter(func(c *fiber.Ctx) string {
			return serviceName + " " + c.Method() + " " + c.Route().Path
		}),
	))

	if metricsEnabled {
		app.Use(monitoring.RecoveryMiddleware())
		app.Use(monitoring.Middleware())
		app.Get("/metrics", adaptor.HTTPHandler(monitoring.Handler()))
	}

	app.Use(middleware.ErrorHandler())

	wrapper := rest_oapi.ServerInterfaceWrapper{Handler: api}

	// System routes (own probes only — gateway has no DB so no /diagnostics/db-tx).
	system := app.Group("/system")
	{
		system.Get("/live", wrapper.Liveness)
		system.Get("/ready", wrapper.Readiness)
	}

	// Bearer middleware mounted per protected route below. Cached once so
	// each route registration reuses the same handler instance.
	bearer := api.RequireBearerToken()

	// /v1 proxy routes — one per backend, dispatched via its adapter.
	v1 := app.Group("/v1")
	{
		// System-probe proxies (REST adapters; retire with each BL-REFACTOR-* card).
		// iam-svc's proxies were removed in BL-IAM-018 / S1-E-12 — iam is now
		// gRPC-only; operators probe via grpc_health_probe. catalog-svc's
		// proxy was removed in G7. The seven pure-scaffold backends
		// (booking/crm/jamaah/logistics/ops/payment/visa) retired their
		// REST surfaces in BL-REFACTOR-002..008 / S1-E-13.
		v1.Get("/finance/system/live", wrapper.GetFinanceSystemLive)

		// Public catalog read (BL-GTW-002 / S1-E-10) — gRPC adapter.
		v1.Get("/packages", wrapper.ListPackages)
		v1.Get("/packages/:id", wrapper.GetPackageById)
		v1.Get("/package-departures/:id", wrapper.GetPackageDepartureById)

		// Client-facing auth (BL-IAM-018 / S1-E-12) — iam-svc gRPC adapter.
		// Public: login + refresh (no bearer). Protected: everything else.
		v1.Post("/sessions", wrapper.Login)
		v1.Post("/sessions/refresh", wrapper.RefreshSession)
		v1.Delete("/sessions", bearer, wrapper.Logout)
		v1.Get("/me", bearer, wrapper.GetMe)
		v1.Post("/me/2fa/enroll", bearer, wrapper.EnrollTotp)
		v1.Post("/me/2fa/verify", bearer, wrapper.VerifyTotp)
		v1.Post("/users/:id/suspend", bearer, wrapper.SuspendUser)
	}

	go func() {
		log.Printf("rest server started successfully 🚀")

		err := app.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			log.Printf("failed to listen at port: %v!\n", port)
			log.Printf("error: %v\n", err)
			os.Exit(1)
		}
	}()
}
