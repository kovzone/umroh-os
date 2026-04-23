package main

import (
	"log"
	"net"
	"os"

	"catalog-svc/api/grpc_api"
	"catalog-svc/api/grpc_api/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	health_pb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// runGrpcServer runs the gRPC server exposing the CatalogService surface plus the
// standard gRPC health protocol.
//
// Per ADR 0009 / BL-REFACTOR-001 (S1-E-11), catalog-svc is now gRPC-only.
// The former REST surface (/v1/packages*, /system/diagnostics/db-tx, /system/live, etc.)
// has been retired; all public reads flow through gateway-svc which proxies to this
// gRPC server via catalog_grpc_adapter. The docker-compose healthcheck uses
// grpc_health_probe against port 50052 (not HTTP).
func runGrpcServer(address string, apiServer *grpc_api.Server) *grpc.Server {
	grpcServer := grpc.NewServer()

	pb.RegisterCatalogServiceServerWithAll(grpcServer, apiServer)

	healthServer := health.NewServer()
	health_pb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(
		"catalog.v1.CatalogService",
		health_pb.HealthCheckResponse_SERVING,
	)
	// Empty service name = whole-server health. Required for grpc_health_probe's
	// default probe (called from docker-compose healthchecks per BL-MON-001).
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
