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
// S1-E-03: POST /v1/bookings added (public in S1; proxied to booking-svc gRPC).
// S4-E-02: POST /v1/leads (public), GET+PUT /v1/leads[/:id] (bearer) added.
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
		// Note: /v1/iam/system/live + /v1/iam/system/diagnostics/db-tx removed in S1-E-12
		//       (iam-svc is gRPC-only). The seven pure-scaffold backends (booking/crm/
		//       jamaah/logistics/ops/payment/visa) retired their REST surfaces in S1-E-13.
		v1.Get("/finance/system/live", wrapper.GetFinanceSystemLive)

		// Public catalog read (BL-GTW-002 / S1-E-10) — gRPC adapter.
		// security: [] in openapi.yaml — no bearer required.
		v1.Get("/packages", wrapper.ListPackages)
		v1.Get("/packages/:id", wrapper.GetPackageById)
		v1.Get("/package-departures/:id", wrapper.GetPackageDepartureById)

		// Booking draft (BL-GTW-003 / S1-E-03) — public in S1; auth arrives with F4.
		// Proxied to booking-svc.BookingService/CreateDraftBooking via gRPC.
		v1.Post("/bookings", wrapper.CreateDraftBooking)

		// CRM lead capture (S4-E-02 / BL-CRM-001) — public.
		// POST /v1/leads — submitting a lead from a landing page or B2C form.
		v1.Post("/leads", wrapper.CreateLead)
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

		// CRM lead management (S4-E-02 / BL-CRM-002..003) — bearer required.
		// GET  /v1/leads     — list leads (cs/admin)
		// GET  /v1/leads/:id — get single lead
		// PUT  /v1/leads/:id — update lead (status, notes, assigned_cs_id)
		v1Protected.Get("/leads", wrapper.ListLeads)
		v1Protected.Get("/leads/:id", wrapper.GetLeadByID)
		v1Protected.Put("/leads/:id", wrapper.UpdateLeadByID)

		// Finance report routes (S5-E-01 / BL-FIN-004..005) — bearer required.
		// GET /v1/finance/summary  — aggregate per-account balances
		// GET /v1/finance/journals — paginated journal entries + lines
		v1Protected.Get("/finance/summary", wrapper.GetFinanceSummary)
		v1Protected.Get("/finance/journals", wrapper.ListJournals)
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
