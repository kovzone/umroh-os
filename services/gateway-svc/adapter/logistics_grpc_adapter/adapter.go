// Package logistics_grpc_adapter is gateway-svc's consumer-side wrapper around
// logistics-svc's gRPC surface (S3 Wave 2 / ADR-0009).
//
// Per ADR-0009, gateway's /v1/logistics/* routes proxy to
// logistics-svc.LogisticsService over gRPC via this adapter.
// The wire contract is in pb/logistics_stub.go, kept in sync by hand with
// services/logistics-svc/api/grpc_api/pb/logistics_new_messages.go.
package logistics_grpc_adapter

import (
	"errors"
	"fmt"

	"gateway-svc/adapter/logistics_grpc_adapter/pb"
	"gateway-svc/util/apperrors"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Adapter wraps the logistics-svc gRPC clients for logistics operations.
type Adapter struct {
	logger                *zerolog.Logger
	tracer                trace.Tracer
	logisticsClient       pb.LogisticsServiceClient
	logisticsPhase6Client pb.LogisticsPhase6Client
}

// NewAdapter creates a new logistics-svc gRPC adapter from an already-dialled conn.
// Ownership of the conn stays with the caller (shared pool lifetime).
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:                logger,
		tracer:                tracer,
		logisticsClient:       pb.NewLogisticsServiceClient(cc),
		logisticsPhase6Client: pb.NewLogisticsPhase6Client(cc),
	}
}

// mapLogisticsError converts a gRPC status error from logistics-svc into an
// apperrors sentinel the gateway error middleware can render.
func mapLogisticsError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("logistics call failed: %w", err))
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
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("logistics unreachable: %s", st.Message()))
	default:
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("logistics call failed (%s): %s", st.Code(), st.Message()))
	}
}
