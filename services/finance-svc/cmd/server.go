package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"finance-svc/adapter/iam_grpc_adapter"
	"finance-svc/api/grpc_api"
	"finance-svc/api/grpc_api/pb"
	"finance-svc/api/rest_oapi"
	"finance-svc/api/rest_oapi/middleware"
	"finance-svc/util/monitoring"

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
//
// Public routes:
//   - GET /system/live
//   - GET /system/ready
//   - GET /system/diagnostics/db-tx
//
// Bearer-protected routes (BL-IAM-002) — require valid iam-svc token:
//   - GET /v1/finance/ping
func runRestServer(port int, api rest_oapi.ServerInterface, iamAdapter *iam_grpc_adapter.Adapter, metricsEnabled bool, serviceName string) {
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
	// to scan (e.g. "finance-svc GET /v1/finance/ping").
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

	// System routes (probes + WithTx diagnostic) — public.
	system := app.Group("/system")
	{
		system.Get("/live", wrapper.Liveness)
		system.Get("/ready", wrapper.Readiness)
		system.Get("/diagnostics/db-tx", wrapper.DbTxDiagnostic)
	}

	// Authenticated routes (BL-IAM-002) — every handler under /v1/finance
	// runs behind the bearer middleware, which fails closed on missing /
	// invalid / unreachable iam-svc per F1 "never default to allow".
	requireBearer := middleware.RequireBearerToken(iamAdapter)
	finance := app.Group("/v1/finance", requireBearer)
	{
		finance.Get("/ping", wrapper.FinancePing)
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

// runGrpcServer runs the gRPC server exposing the FinanceService surface plus the
// standard gRPC health protocol. Finance-svc's own gRPC surface stays a
// placeholder until S3-E-03 — real cross-service RPCs (e.g. the booking-saga
// callback for `payment.received`) land with the S3 contracts.
func runGrpcServer(address string, apiServer *grpc_api.Server) *grpc.Server {
	grpcServer := grpc.NewServer()

	// RegisterFinanceServiceServerWithGRN registers all finance RPCs:
	//   Healthz + OnPaymentReceived + reports + finance-depth + OnGRNReceived (BL-FIN-002).
	// Replace with pb.RegisterFinanceServiceServer once `make genpb` includes
	// all RPCs from finance.proto.
	// RegisterFinanceServiceServerWithCorrections registers all finance RPCs:
	//   Healthz + OnPaymentReceived + reports + finance-depth + OnGRNReceived
	//   + CorrectJournal + DeleteJournalEntry (BL-FIN-006).
	pb.RegisterFinanceServiceServerWithCorrections(grpcServer, apiServer)

	healthServer := health.NewServer()
	health_pb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(
		"finance.v1.FinanceService",
		health_pb.HealthCheckResponse_SERVING,
	)
	// Empty service name = whole-server health. Required for grpc_health_probe's default
	// probe (called from docker-compose healthchecks per ADR 0009 / BL-MON-001).
	healthServer.SetServingStatus(
		"",
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
