package apperrors

import "errors"

// HTTPStatus maps domain sentinels to HTTP status codes for REST.
// Use this in the REST API layer when formatting error responses.
func HTTPStatus(err error) int {
	switch {
	case err == nil:
		return 200
	case errors.Is(err, ErrNotFound):
		return 404
	case errors.Is(err, ErrConflict):
		return 409
	case errors.Is(err, ErrValidation):
		return 400
	case errors.Is(err, ErrUnauthorized):
		return 401
	case errors.Is(err, ErrForbidden):
		return 403
	case errors.Is(err, ErrServiceUnavailable):
		return 502
	case errors.Is(err, ErrInternal):
		return 500
	default:
		return 500
	}
}
