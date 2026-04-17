package main

import (
	"fmt"
	"log"
	"os"

	"gateway-svc/api/rest_oapi"
	"gateway-svc/api/rest_oapi/middleware"
	"gateway-svc/util/monitoring"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// runRestServer runs the REST server using the OpenAPI-generated routes and handler.
// gateway-svc scaffold exposes 12 endpoints:
//
//   - GET /system/live, /system/ready (own probes)
//   - GET /v1/<svc>/system/live for each of the 10 backends (proxies via REST adapters)
//
// Real per-route proxies (auth, packages, bookings, ...) land alongside each
// backend's first feature work.
func runRestServer(port int, api rest_oapi.ServerInterface, metricsEnabled bool) {
	app := fiber.New()

	// CORS — gateway is the edge, accept cross-origin from any browser client.
	corsConfig := cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}
	app.Use(cors.New(corsConfig))

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

	// /v1 proxy routes — one per backend, dispatched via its REST adapter.
	v1 := app.Group("/v1")
	{
		v1.Get("/iam/system/live", wrapper.GetIamSystemLive)
		v1.Get("/catalog/system/live", wrapper.GetCatalogSystemLive)
		v1.Get("/booking/system/live", wrapper.GetBookingSystemLive)
		v1.Get("/jamaah/system/live", wrapper.GetJamaahSystemLive)
		v1.Get("/payment/system/live", wrapper.GetPaymentSystemLive)
		v1.Get("/visa/system/live", wrapper.GetVisaSystemLive)
		v1.Get("/ops/system/live", wrapper.GetOpsSystemLive)
		v1.Get("/logistics/system/live", wrapper.GetLogisticsSystemLive)
		v1.Get("/finance/system/live", wrapper.GetFinanceSystemLive)
		v1.Get("/crm/system/live", wrapper.GetCrmSystemLive)
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
