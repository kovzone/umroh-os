// Package iam_grpc_adapter is finance-svc's consumer-side wrapper around
// iam-svc's internal gRPC surface. It shields the rest of finance-svc from the
// producer's proto types so the service layer + middleware can stay transport-
// neutral. Kept in sync by hand with services/iam-svc/api/grpc_api/pb/iam.proto
// (see pb/iam.proto in this package for the wire contract).
package iam_grpc_adapter

import (
	"finance-svc/adapter/iam_grpc_adapter/pb"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Adapter is a thin wrapper over an iam.v1.IamService client.
type Adapter struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	iamClient pb.IamServiceClient
}

// NewAdapter creates a new iam-svc gRPC adapter from an already-dialled conn.
// Ownership of the conn stays with the caller (shared pool lifetime).
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:    logger,
		tracer:    tracer,
		iamClient: pb.NewIamServiceClient(cc),
	}
}
