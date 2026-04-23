// Package visa_grpc_adapter is gateway-svc's consumer-side wrapper around
// visa-svc's gRPC surface (Phase 6 / BL-VISA-001..003 / ADR-0009).
//
// Per ADR-0009, gateway's /v1/visas/* routes proxy to visa-svc.VisaService
// over gRPC via this adapter.
package visa_grpc_adapter

import (
	"errors"
	"fmt"

	"gateway-svc/adapter/visa_grpc_adapter/pb"
	"gateway-svc/util/apperrors"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Adapter wraps the visa-svc gRPC client.
type Adapter struct {
	logger      *zerolog.Logger
	tracer      trace.Tracer
	visaClient  pb.VisaServiceClient
}

// NewAdapter creates a new visa-svc gRPC adapter from an already-dialled conn.
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:     logger,
		tracer:     tracer,
		visaClient: pb.NewVisaServiceClient(cc),
	}
}

// mapVisaError converts a gRPC status error from visa-svc into an apperrors sentinel.
func mapVisaError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("visa call failed: %w", err))
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return errors.Join(apperrors.ErrNotFound, errors.New(st.Message()))
	case grpcCodes.InvalidArgument:
		return errors.Join(apperrors.ErrValidation, errors.New(st.Message()))
	case grpcCodes.AlreadyExists:
		return errors.Join(apperrors.ErrConflict, errors.New(st.Message()))
	case grpcCodes.PermissionDenied:
		return errors.Join(apperrors.ErrForbidden, errors.New(st.Message()))
	case grpcCodes.Unauthenticated:
		return errors.Join(apperrors.ErrUnauthorized, errors.New(st.Message()))
	case grpcCodes.Unavailable, grpcCodes.DeadlineExceeded, grpcCodes.Canceled, grpcCodes.Unknown:
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("visa-svc unreachable: %s", st.Message()))
	default:
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("visa call failed (%s): %s", st.Code(), st.Message()))
	}
}
