package iam_grpc_adapter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gateway-svc/adapter/iam_grpc_adapter/pb"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidateTokenParams is the adapter-level input.
type ValidateTokenParams struct {
	AccessToken string
}

// ValidateTokenResult is the adapter-level output. Keeps proto types out of
// the bearer-auth middleware and downstream handlers.
type ValidateTokenResult struct {
	UserID    string
	BranchID  string
	SessionID string
	Roles     []string
	ExpiresAt time.Time
}

// ValidateToken delegates to iam-svc.ValidateToken over gRPC.
//
// gRPC codes are translated back to apperrors sentinels so the gateway error
// middleware renders a consistent envelope. Unauthenticated → ErrUnauthorized,
// PermissionDenied → ErrForbidden, anything else → ErrUnauthorized so the
// caller *fails closed* per F1-W7 "never default to allow" and ADR 0009's
// single-point-auth trust boundary.
func (a *Adapter) ValidateToken(ctx context.Context, params *ValidateTokenParams) (*ValidateTokenResult, error) {
	const op = "iam_grpc_adapter.Adapter.ValidateToken"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "ValidateToken"))

	logger := logging.LogWithTrace(ctx, a.logger)

	// Do NOT log params.AccessToken — it is a bearer secret. Intentionally absent.
	resp, err := a.iamClient.ValidateToken(ctx, &pb.ValidateTokenRequest{
		AccessToken: params.AccessToken,
	})
	if err != nil {
		wrapped := mapIamError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	result := &ValidateTokenResult{
		UserID:    resp.GetUserId(),
		BranchID:  resp.GetBranchId(),
		SessionID: resp.GetSessionId(),
		Roles:     resp.GetRoles(),
		ExpiresAt: time.Unix(resp.GetExpiresAtUnix(), 0),
	}
	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// mapIamError converts a grpc-status error from iam-svc into an apperrors
// sentinel the gateway error middleware can render.
//
// Per F1-W7 / ADR 0009 the gateway distinguishes two failure classes:
//   - iam-svc is reachable but rejects the request (bearer bad, user suspended,
//     etc.) → ErrUnauthorized / ErrForbidden (401 / 403).
//   - iam-svc is NOT reachable (transport timeout, connection refused, RST,
//     unclassified error) → ErrServiceUnavailable (502 `SERVICE_UNAVAILABLE`).
//     Fail-closed: request never reaches the backend.
//
// Transport codes treated as "auth unavailable": Unavailable, DeadlineExceeded,
// Canceled, Unknown, and anything that isn't a grpc status at all.
func mapIamError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("iam call failed: %w", err))
	}
	switch st.Code() {
	case grpcCodes.Unauthenticated:
		return errors.Join(apperrors.ErrUnauthorized, errors.New(st.Message()))
	case grpcCodes.PermissionDenied:
		return errors.Join(apperrors.ErrForbidden, errors.New(st.Message()))
	case grpcCodes.InvalidArgument:
		return errors.Join(apperrors.ErrValidation, errors.New(st.Message()))
	case grpcCodes.NotFound:
		return errors.Join(apperrors.ErrNotFound, errors.New(st.Message()))
	case grpcCodes.Unavailable, grpcCodes.DeadlineExceeded, grpcCodes.Canceled, grpcCodes.Unknown:
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("iam unreachable: %s", st.Message()))
	default:
		// Unclassified iam failure → surface as service-unavailable (502),
		// never open the gate, but distinguishable from a bad bearer (401).
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("iam call failed (%s): %s", st.Code(), st.Message()))
	}
}
