package catalog_grpc_adapter

import (
	"errors"
	"fmt"

	"gateway-svc/util/apperrors"

	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mapCatalogError converts a grpc-status error from catalog-svc into an
// apperrors sentinel the gateway error middleware can render.
//
// NotFound, InvalidArgument, AlreadyExists map 1:1. Transport-level
// failures (Unavailable, DeadlineExceeded, Canceled, Unknown, non-grpc
// errors) map to ErrServiceUnavailable → 502, matching the iam_grpc_adapter
// policy for "backend down" cases.
func mapCatalogError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("catalog call failed: %w", err))
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
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("catalog unreachable: %s", st.Message()))
	default:
		return errors.Join(apperrors.ErrInternal, fmt.Errorf("catalog call failed (%s): %s", st.Code(), st.Message()))
	}
}
