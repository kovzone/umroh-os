// errors.go — gRPC → apperrors mapping for booking_grpc_adapter.
//
// Mirrors the pattern used in catalog_grpc_adapter/errors.go.
// Booking-svc gRPC errors are mapped to the gateway's apperrors sentinels
// so that rest_oapi handlers can emit the correct HTTP status.
package booking_grpc_adapter

import (
	"errors"
	"fmt"

	"gateway-svc/util/apperrors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mapBookingError converts a gRPC status error from booking-svc into a
// gateway-local apperrors sentinel error. Non-gRPC errors pass through.
func mapBookingError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("booking-svc: %w", err)
	}
	switch st.Code() {
	case codes.NotFound:
		return errors.Join(apperrors.ErrNotFound, fmt.Errorf("booking-svc: %s", st.Message()))
	case codes.AlreadyExists:
		return errors.Join(apperrors.ErrConflict, fmt.Errorf("booking-svc: %s", st.Message()))
	case codes.InvalidArgument:
		return errors.Join(apperrors.ErrValidation, fmt.Errorf("booking-svc: %s", st.Message()))
	case codes.Unauthenticated:
		return errors.Join(apperrors.ErrUnauthorized, fmt.Errorf("booking-svc: %s", st.Message()))
	case codes.PermissionDenied:
		return errors.Join(apperrors.ErrForbidden, fmt.Errorf("booking-svc: %s", st.Message()))
	case codes.Unavailable:
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("booking-svc unavailable: %s", st.Message()))
	default:
		return fmt.Errorf("booking-svc internal error: %s", st.Message())
	}
}
