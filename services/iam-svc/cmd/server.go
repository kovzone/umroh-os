package main

import (
	"log"
	"net"
	"os"

	"iam-svc/api/grpc_api"
	"iam-svc/api/grpc_api/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	health_pb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// runGrpcServer runs the gRPC server exposing the IamService surface plus the
// standard gRPC health protocol. BL-IAM-018 / S1-E-12 added the client-facing
// auth RPCs that used to live on the retired REST port: Login / RefreshSession
// / Logout / GetMe / EnrollTotp / VerifyTotp / SuspendUser. gateway-svc
// validates bearers once at the edge (ADR 0009) and proxies here.
func runGrpcServer(address string, apiServer *grpc_api.Server) *grpc.Server {
	grpcServer := grpc.NewServer()

	pb.RegisterIamServiceServer(grpcServer, apiServer)

	healthServer := health.NewServer()
	health_pb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(
		"iam.v1.IamService",
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
