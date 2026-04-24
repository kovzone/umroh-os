// Package ops_grpc_adapter is gateway-svc's consumer-side wrapper around
// ops-svc's gRPC surface (S3 Wave 2 / ADR-0009).
//
// Per ADR-0009, gateway's /v1/ops/* routes proxy to ops-svc.OpsService over
// gRPC via this adapter. The wire contract is in pb/ops_stub.go, kept in sync
// by hand with services/ops-svc/api/grpc_api/pb/ops_messages.go.
package ops_grpc_adapter

import (
	"errors"
	"fmt"

	"gateway-svc/adapter/ops_grpc_adapter/pb"
	"gateway-svc/util/apperrors"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Adapter wraps the ops-svc gRPC clients for ops operations.
type Adapter struct {
	logger          *zerolog.Logger
	tracer          trace.Tracer
	opsClient       pb.OpsServiceClient
	opsPhase6Client pb.OpsPhase6Client
	depthClient     pb.OpsDepthClient
}

// NewAdapter creates a new ops-svc gRPC adapter from an already-dialled conn.
// Ownership of the conn stays with the caller (shared pool lifetime).
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:          logger,
		tracer:          tracer,
		opsClient:       pb.NewOpsServiceClient(cc),
		opsPhase6Client: pb.NewOpsPhase6Client(cc),
		depthClient:     pb.NewOpsDepthClient(cc),
	}
}

// mapOpsError converts a gRPC status error from ops-svc into an
// apperrors sentinel the gateway error middleware can render.
func mapOpsError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("ops call failed: %w", err))
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return errors.Join(apperrors.ErrNotFound, errors.New(st.Message()))
	case grpcCodes.InvalidArgument:
		return errors.Join(apperrors.ErrValidation, errors.New(st.Message()))
	case grpcCodes.PermissionDenied:
		return errors.Join(apperrors.ErrForbidden, errors.New(st.Message()))
	case grpcCodes.Unauthenticated:
		return errors.Join(apperrors.ErrUnauthorized, errors.New(st.Message()))
	case grpcCodes.Unavailable, grpcCodes.DeadlineExceeded, grpcCodes.Canceled, grpcCodes.Unknown:
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("ops unreachable: %s", st.Message()))
	default:
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("ops call failed (%s): %s", st.Code(), st.Message()))
	}
}
