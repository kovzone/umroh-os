package main

import (
	"log"
	"net"
	"os"

	"finance-svc/api/grpc_api"
	"finance-svc/api/grpc_api/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	health_pb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// runGrpcServer runs the gRPC server exposing the FinanceService surface plus
// the standard gRPC health protocol. Post BL-IAM-019 / S1-E-14 finance-svc is
// gRPC-only (ADR 0009); the REST port was retired together with the migration
// of /v1/finance/ping onto gateway-svc.
//
// FinanceService currently carries the permission-gate smoke (`FinancePing`)
// invoked through the gateway after its RequireBearerToken + RequirePermission
// middleware has approved the caller. Real finance RPCs (journals, AR/AP,
// reports) land with S3-E-03 + S3-E-07.
func runGrpcServer(address string, apiServer *grpc_api.Server) *grpc.Server {
	grpcServer := grpc.NewServer()

	pb.RegisterFinanceServiceServer(grpcServer, apiServer)

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
