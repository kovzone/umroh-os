package service

import (
	"context"
	"errors"
	"fmt"

	"finance-svc/adapter/iam_grpc_adapter"
	"finance-svc/util/apperrors"
	"finance-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// IamChecker is the slice of iam-svc the finance-svc service layer needs.
// Defined here (not in adapter/) so tests can inject a testify/mock double
// without depending on the gRPC proto types. The adapter type satisfies it.
type IamChecker interface {
	CheckPermission(ctx context.Context, params *iam_grpc_adapter.CheckPermissionParams) (*iam_grpc_adapter.CheckPermissionResult, error)
}

// FinancePingParams is the input for the placeholder authenticated ping route.
// The UserID + Roles come from the bearer middleware (populated from
// iam-svc.ValidateToken) — handlers must not accept them from the request body.
type FinancePingParams struct {
	UserID string
	Roles  []string
}

// FinancePingResult is the success envelope returned to the caller.
type FinancePingResult struct {
	Message string
	UserID  string
	Roles   []string
}

// FinancePing exercises the finance-svc permission gate against iam-svc. It is
// the smallest realistic surface that proves "finance routes denied for
// non-finance roles" per F1 acceptance — intentionally a no-op handler; real
// finance endpoints (journals, AR/AP, reports) land with S3-E-03 + S3-E-07.
//
// Permission tuple (journal_entry / read / global) is the entry point every
// finance endpoint will reuse as its base gate in deeper cards.
func (s *Service) FinancePing(ctx context.Context, params *FinancePingParams) (*FinancePingResult, error) {
	const op = "service.Service.FinancePing"
	logger := logging.LogWithTrace(ctx, s.logger)
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("user_id", params.UserID),
	)

	if params == nil || params.UserID == "" {
		e := errors.Join(apperrors.ErrUnauthorized, errors.New("missing user_id on authenticated request"))
		logger.Warn().Err(e).Msg("")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	decision, err := s.iamChecker.CheckPermission(ctx, &iam_grpc_adapter.CheckPermissionParams{
		UserID:   params.UserID,
		Resource: "journal_entry",
		Action:   "read",
		Scope:    "global",
	})
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	if !decision.Allowed {
		e := errors.Join(apperrors.ErrForbidden, fmt.Errorf("user %s lacks journal_entry/read/global", params.UserID))
		logger.Info().Str("user_id", params.UserID).Msg("finance access denied")
		span.RecordError(e)
		span.SetStatus(codes.Error, e.Error())
		return nil, e
	}

	span.SetStatus(codes.Ok, "success")
	logger.Info().Str("user_id", params.UserID).Msg("finance ping allowed")
	return &FinancePingResult{
		Message: "ok",
		UserID:  params.UserID,
		Roles:   params.Roles,
	}, nil
}
