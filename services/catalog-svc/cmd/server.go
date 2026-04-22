package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"catalog-svc/api/grpc_api"
	"catalog-svc/api/grpc_api/pb"
	"catalog-svc/api/rest_oapi"
	"catalog-svc/api/rest_oapi/middleware"
	"catalog-svc/util/monitoring"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	health_pb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// runRestServer runs the REST server using the OpenAPI-generated routes and handler.
// catalog-svc exposes:
//
//   System probes:
//   - GET /system/live
//   - GET /system/ready
//   - GET /system/diagnostics/db-tx
//
//   Public catalog (§ Catalog, S1-E-02 / BL-CAT-001):
//   - GET /v1/packages       — list active packages (filterable, cursor-paginated)
//   - GET /v1/packages/{id}  — package detail (active only; 404 for draft/archived)
//
// Departure detail (`GET /v1/package-departures/{id}`) and admin write
// endpoints land in later S1-E-02 / S1-E-05 cards.
func runRestServer(port int, api rest_oapi.ServerInterface, metricsEnabled bool, serviceName string) {
	app := fiber.New()

	// CORS
	corsConfig := cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}
	app.Use(cors.New(corsConfig))

	// OpenTelemetry — start an inbound span for every request so cross-service
	// traces initiated by an upstream caller (via W3C traceparent) continue in
	// this process's spans instead of starting a new trace. Span names are
	// prefixed with the service name so a multi-service trace in Tempo is easy
	// to scan (e.g. "catalog-svc GET /system/diagnostics/db-tx").
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

	// System routes (probes + WithTx diagnostic)
	system := app.Group("/system")
	{
		system.Get("/live", wrapper.Liveness)
		system.Get("/ready", wrapper.Readiness)
		system.Get("/diagnostics/db-tx", wrapper.DbTxDiagnostic)
	}

	// Public catalog routes (no bearer; § Catalog contract)
	v1 := app.Group("/v1")
	{
		v1.Get("/packages", wrapper.ListPackages)
		v1.Get("/packages/:id", wrapper.GetPackageById)
		v1.Get("/package-departures/:id", wrapper.GetPackageDepartureById)
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

// runGrpcServer runs the gRPC server exposing the CatalogService surface plus the
// standard gRPC health protocol. Pilot wires a placeholder Healthz RPC;
// ValidateToken, CheckPermission, GetUser, and RecordAudit land in F1.7.
func runGrpcServer(address string, apiServer *grpc_api.Server) *grpc.Server {
	grpcServer := grpc.NewServer()

	pb.RegisterCatalogServiceServer(grpcServer, apiServer)

	healthServer := health.NewServer()
	health_pb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(
		"catalog.v1.CatalogService",
		health_pb.HealthCheckResponse_SERVING,
	)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("failed to listen at: %s!\n", address)
		log.Printf("error: %v\n", err)
		os.Exit(1)
	}

	log.Printf("gRPC listening at: %s", address)

	go func() {
		log.Printf("gRPC server started successfully 🚀")

		if err := grpcServer.Serve(listener); err != nil {
			log.Printf("failed to serve gRPC: %v\n", err)
			os.Exit(1)
		}
	}()

	return grpcServer
}
