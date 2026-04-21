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

// GRPCMessage returns a constant, code-derived message for the gRPC status body.
// It exists so the wire response never leaks the wrapped Go error chain (which
// can become a state oracle — e.g. "session revoked" vs "load session: not
// found"). The full detail still goes to zerolog and the span via RecordError
// at the caller; only the outbound string is sanitised.
func GRPCMessage(err error) string {
	switch {
	case errors.Is(err, ErrNotFound):
		return "not found"
	case errors.Is(err, ErrConflict):
		return "conflict"
	case errors.Is(err, ErrValidation):
		return "validation error"
	case errors.Is(err, ErrUnauthorized):
		return "unauthorized"
	case errors.Is(err, ErrForbidden):
		return "forbidden"
	default:
		return "internal error"
	}
}
