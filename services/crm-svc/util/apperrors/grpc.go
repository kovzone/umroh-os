package apperrors

import (
	"errors"

	"google.golang.org/grpc/codes"
)

func GRPCCode(err error) codes.Code {
	switch {
	case errors.Is(err, ErrNotFound):
		return codes.NotFound
	case errors.Is(err, ErrConflict):
		return codes.AlreadyExists
	case errors.Is(err, ErrValidation):
		return codes.InvalidArgument
	case errors.Is(err, ErrUnauthorized):
		return codes.Unauthenticated
	case errors.Is(err, ErrForbidden):
		return codes.PermissionDenied
	default:
		return codes.Internal
	}
}
