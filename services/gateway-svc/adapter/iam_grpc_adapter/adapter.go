// Package iam_grpc_adapter is gateway-svc's consumer-side wrapper around
// iam-svc's internal gRPC surface. It shields the rest of gateway-svc from
// the producer's proto types so the bearer-auth middleware stays transport-
// neutral.
//
// Per ADR 0009 + F1-W7, gateway validates every authenticated request once
// at the edge via iam-svc.ValidateToken (gRPC); backends trust the forwarded
// call. CheckPermission and RecordAudit are included at scaffold time for
// symmetry with the booking-svc / finance-svc adapters — gateway-level
// handlers may use them in future cards (e.g. audit on staff mutating
// routes). The wire contract is in pb/iam.proto, kept in sync by hand with
// services/iam-svc/api/grpc_api/pb/iam.proto.
package iam_grpc_adapter

import (
	"gateway-svc/adapter/iam_grpc_adapter/pb"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Adapter is a thin wrapper over an iam.v1.IamService client.
type Adapter struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	iamClient   pb.IamServiceClient
	adminClient pb.IamAdminClient
}

// NewAdapter creates a new iam-svc gRPC adapter from an already-dialled conn.
// Ownership of the conn stays with the caller (shared pool lifetime).
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:      logger,
		tracer:      tracer,
		iamClient:   pb.NewIamServiceClient(cc),
		adminClient: pb.NewIamAdminClient(cc),
	}
}
