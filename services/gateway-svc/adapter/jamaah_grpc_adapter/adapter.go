// Package jamaah_grpc_adapter is gateway-svc's consumer-side wrapper around
// jamaah-svc's gRPC manifest surface (Wave 1A / Phase 6).
//
// Per ADR-0009, gateway's /v1/manifest/:departure_id route proxies to
// jamaah-svc.JamaahService/GetDepartureManifest via this adapter. The wire
// contract is in pb/jamaah_manifest_stub.go, kept in sync by hand with
// services/jamaah-svc/api/grpc_api/pb/manifest_messages.go.
package jamaah_grpc_adapter

import (
	"errors"
	"fmt"

	"gateway-svc/adapter/jamaah_grpc_adapter/pb"
	"gateway-svc/util/apperrors"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Adapter wraps the jamaah-svc gRPC client for manifest and OCR operations.
type Adapter struct {
	logger         *zerolog.Logger
	tracer         trace.Tracer
	manifestClient pb.ManifestClient
	ocrClient      pb.JamaahOCRClient
}

// NewAdapter creates a new jamaah-svc gRPC adapter from an already-dialled conn.
// Ownership of the conn stays with the caller (shared pool lifetime).
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:         logger,
		tracer:         tracer,
		manifestClient: pb.NewManifestClient(cc),
		ocrClient:      pb.NewJamaahOCRClient(cc),
	}
}

// mapJamaahError converts a gRPC status error from jamaah-svc into an
// apperrors sentinel the gateway error middleware can render.
func mapJamaahError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("jamaah call failed: %w", err))
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
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("jamaah unreachable: %s", st.Message()))
	default:
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("jamaah call failed (%s): %s", st.Code(), st.Message()))
	}
}
