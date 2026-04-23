package iam_grpc_adapter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"catalog-svc/adapter/iam_grpc_adapter/pb"
	"catalog-svc/util/apperrors"
	"catalog-svc/util/logging"

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
// downstream middleware and handlers.
type ValidateTokenResult struct {
	UserID    string
	BranchID  string
	SessionID string
	Roles     []string
	ExpiresAt time.Time
}

// ValidateToken delegates to iam-svc.ValidateToken over gRPC.
//
// gRPC codes are translated back to apperrors sentinels so the REST error
// middleware renders the same envelope as native catalog-svc errors. Unauthenticated
// → ErrUnauthorized, PermissionDenied → ErrForbidden, anything else → ErrUnauthorized
// so the consumer *fails closed* per F1 "never default to allow".
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
// sentinel the catalog-svc error middleware can render. The default branch is
// ErrUnauthorized (not ErrInternal) for ValidateToken specifically: unreachable
// iam-svc must fail-closed as 401, never as a 500 that could be mistaken for
// "transient, retry later".
func mapIamError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrUnauthorized, fmt.Errorf("iam call failed: %w", err))
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
	default:
		// Fail-closed: any unclassified iam failure looks like "auth unavailable",
		// and the consumer middleware returns 401, never opens the gate.
		return errors.Join(apperrors.ErrUnauthorized, fmt.Errorf("iam call failed: %s", st.Message()))
	}
}
