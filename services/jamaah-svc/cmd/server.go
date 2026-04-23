package main

import (
	"log"
	"net"
	"os"

	"jamaah-svc/api/grpc_api"
	"jamaah-svc/api/grpc_api/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	health_pb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// runGrpcServer runs the gRPC server exposing the JamaahService surface plus
// the standard gRPC health protocol. Per ADR 0009 this is jamaah-svc's only
// transport — REST was removed in BL-REFACTOR-004 / S1-E-13.
func runGrpcServer(address string, apiServer *grpc_api.Server) *grpc.Server {
	grpcServer := grpc.NewServer()

	// RegisterJamaahServiceServerFull registers the generated Healthz RPC and
	// all hand-written RPCs: UploadDocument, ReviewDocument, GetDepartureManifest,
	// TriggerOCR, GetOCRStatus (S3-E-02 / BL-DOC-001..003, BL-DOC-002).
	// Replace with pb.RegisterJamaahServiceServer once `make genpb` includes
	// those RPCs from jamaah.proto.
	pb.RegisterJamaahServiceServerFull(grpcServer, apiServer)

	healthServer := health.NewServer()
	health_pb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(
		"jamaah.v1.JamaahService",
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
