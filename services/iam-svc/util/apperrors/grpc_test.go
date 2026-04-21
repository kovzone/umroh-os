package apperrors_test

import (
	"errors"
	"testing"

	"iam-svc/util/apperrors"

	"google.golang.org/grpc/codes"
)

func Test_GRPCCode(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want codes.Code
	}{
		{"nil", nil, codes.Internal},
		{"not found", apperrors.ErrNotFound, codes.NotFound},
		{"conflict", apperrors.ErrConflict, codes.AlreadyExists},
		{"validation", apperrors.ErrValidation, codes.InvalidArgument},
		{"unauthorized", apperrors.ErrUnauthorized, codes.Unauthenticated},
		{"forbidden", apperrors.ErrForbidden, codes.PermissionDenied},
		{"internal", apperrors.ErrInternal, codes.Internal},
		{"unknown", errors.New("bare error"), codes.Internal},
		{"wrapped unauthorized", errors.Join(apperrors.ErrUnauthorized, errors.New("bad password")), codes.Unauthenticated},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := apperrors.GRPCCode(tc.err)
			if got != tc.want {
				t.Fatalf("GRPCCode(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

// GRPCMessage returns constants only — never the wrapped chain. This guards
// against a state-oracle leak (e.g. "session revoked" vs "load session: not
// found") on the gRPC wire; the original chain still goes to logs/spans.
func Test_GRPCMessage(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"not found", apperrors.ErrNotFound, "not found"},
		{"conflict", apperrors.ErrConflict, "conflict"},
		{"validation", apperrors.ErrValidation, "validation error"},
		{"unauthorized", apperrors.ErrUnauthorized, "unauthorized"},
		{"forbidden", apperrors.ErrForbidden, "forbidden"},
		{"internal", apperrors.ErrInternal, "internal error"},
		{"unknown defaults to internal", errors.New("bare error"), "internal error"},
		{
			"wrapped unauthorized does NOT leak the inner chain",
			errors.Join(apperrors.ErrUnauthorized, errors.New("load session: not found")),
			"unauthorized",
		},
		{
			"wrapped unauthorized with revoked-session inner does NOT leak either",
			errors.Join(apperrors.ErrUnauthorized, errors.New("session revoked")),
			"unauthorized",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := apperrors.GRPCMessage(tc.err)
			if got != tc.want {
				t.Fatalf("GRPCMessage(%v) = %q, want %q", tc.err, got, tc.want)
			}
		})
	}
}
