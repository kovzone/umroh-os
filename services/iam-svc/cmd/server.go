package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"iam-svc/api/grpc_api"
	"iam-svc/api/grpc_api/pb"
	"iam-svc/api/rest_oapi"
	"iam-svc/api/rest_oapi/middleware"
	"iam-svc/util/monitoring"
	"iam-svc/util/token"

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
// Routes:
//
//   - GET    /system/live                 (public)
//   - GET    /system/ready                (public)
//   - GET    /system/diagnostics/db-tx    (public; pilot scope)
//   - POST   /v1/sessions                 (public — login)
//   - POST   /v1/sessions/refresh         (public — rotate)
//   - DELETE /v1/sessions                 (bearer — logout)
//   - GET    /v1/me                       (bearer)
//   - POST   /v1/me/2fa/enroll            (bearer)
//   - POST   /v1/me/2fa/verify            (bearer)
//
// `bearer` routes go through middleware.RequireBearerToken(tokenMaker) which
// puts a verified *token.Payload into c.Locals(middleware.PayloadKey).
func runRestServer(port int, api rest_oapi.ServerInterface, tokenMaker token.Maker, metricsEnabled bool, serviceName string) {
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
	// to scan (e.g. "iam-svc GET /system/diagnostics/db-tx").
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
	requireBearer := middleware.RequireBearerToken(tokenMaker)

	// System routes (probes + WithTx diagnostic) — public.
	system := app.Group("/system")
	{
		system.Get("/live", wrapper.Liveness)
		system.Get("/ready", wrapper.Readiness)
		system.Get("/diagnostics/db-tx", wrapper.DbTxDiagnostic)
	}

	// Session routes — /v1/sessions. Login + refresh are public; logout requires bearer.
	sessions := app.Group("/v1/sessions")
	{
		sessions.Post("", wrapper.Login)
		sessions.Post("/refresh", wrapper.RefreshSession)
		sessions.Delete("", requireBearer, wrapper.Logout)
	}

	// Current-user routes — /v1/me. All require bearer.
	me := app.Group("/v1/me", requireBearer)
	{
		me.Get("", wrapper.GetMe)
		me.Post("/2fa/enroll", wrapper.EnrollTOTP)
		me.Post("/2fa/verify", wrapper.VerifyTOTP)
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

// runGrpcServer runs the gRPC server exposing the IamService surface plus the
// standard gRPC health protocol. Pilot wires a placeholder Healthz RPC;
// ValidateToken, CheckPermission, GetUser, and RecordAudit land in F1.7.
func runGrpcServer(address string, apiServer *grpc_api.Server) *grpc.Server {
	grpcServer := grpc.NewServer()

	pb.RegisterIamServiceServer(grpcServer, apiServer)

	healthServer := health.NewServer()
	health_pb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(
		"iam.v1.IamService",
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
