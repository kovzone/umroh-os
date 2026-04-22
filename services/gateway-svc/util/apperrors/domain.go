// Package apperrors provides domain-level sentinel errors and transport adapters.
//
// Design:
//   - domain.go: Transport-agnostic sentinels. Used by service layer to wrap errors.
//   - http.go: REST adapter. Maps sentinels to HTTP status codes. Used by API layer.
//   - grpc.go: gRPC adapter (future). Maps sentinels to gRPC codes. Add when needed.
//
// This package is reusable and can be copied to other projects.
package apperrors

import "errors"

// Domain sentinel errors. Transport-agnostic; use with errors.Is/As.
var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrValidation   = errors.New("validation error")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrInternal     = errors.New("internal error")
	// ErrServiceUnavailable is raised when an upstream service gateway depends on
	// (notably iam-svc for bearer validation per F1-W7 / ADR 0009) is unreachable.
	// Distinct from ErrUnauthorized so the client can tell "your bearer is bad"
	// (401) apart from "the auth layer is down" (502).
	ErrServiceUnavailable = errors.New("service unavailable")
)

// ErrorCode returns a machine-readable code for the error. Used in REST/API response envelopes.
func ErrorCode(err error) string {
	switch {
	case err == nil:
		return ""
	case errors.Is(err, ErrNotFound):
		return "NOT_FOUND"
	case errors.Is(err, ErrConflict):
		return "CONFLICT"
	case errors.Is(err, ErrValidation):
		return "VALIDATION_ERROR"
	case errors.Is(err, ErrUnauthorized):
		return "UNAUTHORIZED"
	case errors.Is(err, ErrForbidden):
		return "FORBIDDEN"
	case errors.Is(err, ErrServiceUnavailable):
		return "SERVICE_UNAVAILABLE"
	case errors.Is(err, ErrInternal):
		return "INTERNAL_ERROR"
	default:
		return "INTERNAL_ERROR"
	}
}
