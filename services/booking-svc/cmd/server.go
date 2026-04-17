package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"booking-svc/api/grpc_api"
	"booking-svc/api/grpc_api/pb"
	"booking-svc/api/rest_oapi"
	"booking-svc/api/rest_oapi/middleware"
	"booking-svc/util/monitoring"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	health_pb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// runRestServer runs the REST server using the OpenAPI-generated routes and handler.
// booking-svc pilot scaffold exposes only the three standard scaffold endpoints:
//
//   - GET /system/live
//   - GET /system/ready
//   - GET /system/diagnostics/db-tx
//
// Real iam REST routes (user/role/branch/audit + auth login/refresh/logout)
// land in F1.5–F1.11.
func runRestServer(port int, api rest_oapi.ServerInterface, metricsEnabled bool) {
	app := fiber.New()

	// CORS
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

	// System routes (probes + WithTx diagnostic)
	system := app.Group("/system")
	{
		system.Get("/live", wrapper.Liveness)
		system.Get("/ready", wrapper.Readiness)
		system.Get("/diagnostics/db-tx", wrapper.DbTxDiagnostic)
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

// runGrpcServer runs the gRPC server exposing the BookingService surface plus the
// standard gRPC health protocol. Pilot wires a placeholder Healthz RPC;
// ValidateToken, CheckPermission, GetUser, and RecordAudit land in F1.7.
func runGrpcServer(address string, apiServer *grpc_api.Server) *grpc.Server {
	grpcServer := grpc.NewServer()

	pb.RegisterBookingServiceServer(grpcServer, apiServer)

	healthServer := health.NewServer()
	health_pb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(
		"booking.v1.BookingService",
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
