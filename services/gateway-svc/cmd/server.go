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
//
// Route topology (per ADR 0009 / BL-GTW-001):
//
//   - /system/live, /system/ready      — public (own probes)
//   - /v1/<svc>/system/live            — public probe proxies (remaining REST backends)
//   - /v1/packages*, /v1/package-departures/* — public catalog read (security: [])
//   - /v1/auth/*                       — IAM auth routes (Login/Refresh public; Logout bearer)
//   - /v1/me, /v1/me/2fa/*             — bearer-auth IAM user routes
//   - /v1/users/:id/suspend            — bearer + permission (super_admin)
//   - POST /v1/packages                — staff create package (bearer + catalog.package.manage)
//   - PUT /v1/packages/:id             — staff update package (bearer + catalog.package.manage)
//   - DELETE /v1/packages/:id          — staff delete package (bearer + catalog.package.manage)
//   - POST /v1/packages/:id/departures — staff create departure (bearer + catalog.package.manage)
//   - PUT /v1/departures/:id           — staff update departure (bearer + catalog.package.manage)
//
// S1-E-11: /v1/catalog/system/live removed (catalog-svc is gRPC-only).
// S1-E-12: /v1/iam/system/live + /v1/iam/system/diagnostics/db-tx removed (iam-svc is gRPC-only).
//          IAM client-facing auth routes moved here from iam-svc REST.
//
// iamValidator is the *iam_grpc_adapter.Adapter produced in start.go; it is
// passed as the interface type so unit tests can substitute a stub.
func runRestServer(port int, api rest_oapi.ServerInterface, iamValidator middleware.IamValidator, metricsEnabled bool, serviceName string) {
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

	// ── Public routes ────────────────────────────────────────────────────────

	// System routes (own probes only — gateway has no DB so no /diagnostics/db-tx).
	system := app.Group("/system")
	{
		system.Get("/live", wrapper.Liveness)
		system.Get("/ready", wrapper.Readiness)
	}

	// /v1 proxy routes — one per backend, dispatched via its adapter.
	v1 := app.Group("/v1")
	{
		// System-probe proxies (REST adapters; retire with each BL-REFACTOR-* card).
		// Note: /v1/catalog/system/live removed in S1-E-11 (catalog-svc is gRPC-only).
		// Note: /v1/iam/system/live removed in S1-E-12 (iam-svc is gRPC-only).
		v1.Get("/booking/system/live", wrapper.GetBookingSystemLive)
		v1.Get("/jamaah/system/live", wrapper.GetJamaahSystemLive)
		v1.Get("/payment/system/live", wrapper.GetPaymentSystemLive)
		v1.Get("/visa/system/live", wrapper.GetVisaSystemLive)
		v1.Get("/ops/system/live", wrapper.GetOpsSystemLive)
		v1.Get("/logistics/system/live", wrapper.GetLogisticsSystemLive)
		v1.Get("/finance/system/live", wrapper.GetFinanceSystemLive)
		v1.Get("/crm/system/live", wrapper.GetCrmSystemLive)

		// Public catalog read (BL-GTW-002 / S1-E-10) — gRPC adapter.
		// security: [] in openapi.yaml — no bearer required.
		v1.Get("/packages", wrapper.ListPackages)
		v1.Get("/packages/:id", wrapper.GetPackageById)
		v1.Get("/package-departures/:id", wrapper.GetPackageDepartureById)
	}

	// IAM auth public routes (BL-IAM-018 / S1-E-12) — no bearer required.
	// Login + refresh are public; they issue / rotate the bearer.
	auth := v1.Group("/auth")
	{
		auth.Post("/login", wrapper.Login)
		auth.Post("/refresh", wrapper.RefreshSession)
	}

	// ── Protected routes (BL-GTW-001 / S1-E-09) ─────────────────────────────
	//
	// Every route in this group requires a valid Bearer token. The
	// RequireBearerToken middleware calls iam-svc.ValidateToken (gRPC) and
	// injects *middleware.Identity into c.Locals(middleware.IdentityKey) for
	// downstream handlers to consume. Failure modes per F1-W7:
	//
	//   - Missing / malformed Authorization header → 401 UNAUTHORIZED
	//   - iam-svc rejects (expired, revoked)       → 401 UNAUTHORIZED
	//   - iam-svc unreachable / timeout             → 502 SERVICE_UNAVAILABLE
	v1Protected := app.Group("/v1", middleware.RequireBearerToken(iamValidator))
	{
		// IAM auth bearer routes (BL-IAM-018 / S1-E-12).
		v1Protected.Delete("/auth/logout", wrapper.Logout)
		v1Protected.Get("/me", wrapper.GetMe)
		v1Protected.Post("/me/2fa/enroll", wrapper.EnrollTOTP)
		v1Protected.Post("/me/2fa/verify", wrapper.VerifyTOTP)
		v1Protected.Post("/users/:id/suspend", wrapper.SuspendUser)

		// Staff catalog write routes (BL-CAT-014 / S1-E-07).
		// All require bearer + catalog.package.manage permission.
		// Response shapes align with public read models (PackageDetail, DepartureDetail).
		v1Protected.Post("/packages", wrapper.CreatePackage)
		v1Protected.Put("/packages/:id", wrapper.UpdatePackageById)
		v1Protected.Delete("/packages/:id", wrapper.DeletePackageById)
		v1Protected.Post("/packages/:id/departures", wrapper.CreateDeparture)
		v1Protected.Put("/departures/:id", wrapper.UpdateDepartureById)
	}

	go func() {
		log.Printf("rest server started successfully on :%d", port)

		err := app.Listen(fmt.Sprintf(":%d", port))
		if err != nil {
			log.Printf("failed to listen at port: %v!\n", port)
			log.Printf("error: %v\n", err)
			os.Exit(1)
		}
	}()
}
