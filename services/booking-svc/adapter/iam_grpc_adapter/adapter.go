// Package iam_grpc_adapter is booking-svc's consumer-side wrapper around
// iam-svc's internal gRPC surface. It shields the rest of booking-svc from the
// producer's proto types so the service layer + middleware can stay transport-
// neutral. Kept in sync by hand with services/iam-svc/api/grpc_api/pb/iam.proto
// (see pb/iam.proto in this package for the wire contract).
//
// Scaffolded as part of BL-IAM-004: ValidateToken + CheckPermission land the
// bearer-auth + permission-gate surface for future booking handlers, and
// RecordAudit is the write path every state-changing booking action will call
// as a one-line emit. No booking-svc handler consumes the adapter yet — first
// consumer lands with S1-E-03 (POST /v1/bookings draft).
package iam_grpc_adapter

import (
	"booking-svc/adapter/iam_grpc_adapter/pb"

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
