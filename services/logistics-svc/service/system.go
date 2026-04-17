package service

import (
	"context"
	"errors"
	"fmt"

	"logistics-svc/store/postgres_store"
	"logistics-svc/store/postgres_store/sqlc"
	"logistics-svc/util/apperrors"
	"logistics-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type LivenessParams struct{}

type LivenessResult struct {
	OK bool `json:"ok"`
}

// Liveness verifies that the service process is alive and responsive at the service layer.
// It must not access any external dependencies (database, cache, network, filesystem).
// Trivial and deterministic; returns success as long as the service can execute and respond.
// Any failure indicates a broken or wedged process.
//
// No logging or tracing — probe endpoints are hit at high frequency; kept cheap.
func (s *Service) Liveness(ctx context.Context, params *LivenessParams) (*LivenessResult, error) {
	return &LivenessResult{OK: true}, nil
}

type ReadinessParams struct{}

type ReadinessResult struct {
	OK bool `json:"ok"`
}

// Readiness verifies that the service instance is capable of handling real traffic at this moment.
// It checks required external dependencies (e.g. database connectivity via read-only SELECT 1)
// but must not perform any state-mutating operations. Temporary failures are expected and
// should cause readiness to fail without crashing or restarting the service.
//
// No logging or tracing — probe endpoints are hit at high frequency; kept cheap.
func (s *Service) Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error) {
	_, err := s.store.ReadyCheck(ctx)
	if err != nil {
		return nil, fmt.Errorf("ready check: %w", postgres_store.WrapDBError(err))
	}
	return &ReadinessResult{OK: true}, nil
}

type DbTxDiagnosticParams struct {
	Message string `json:"message"`
}

type DbTxDiagnosticResult struct {
	ID int64 `json:"id"`
}

// DbTxDiagnostic executes a deliberate, state-mutating database transaction to
// demonstrate and validate the service's transaction handling pattern (WithTx),
// including commit behavior. It is the canonical reference for how services
// should use WithTx — NOT a health/readiness probe.
//
// The pattern: wrap every multi-step DB interaction in a WithTx callback; wrap
// every error returned from a sqlc call through WrapDBError so the API layer
// receives the correct HTTP status; return the domain sentinel on logical
// failures (joined via errors.Join) so apperrors maps it correctly.
//
// All services write to the same `public.diagnostics` table, stamped with the
// service's app name (per the 000002_scaffold_services migration).
func (s *Service) DbTxDiagnostic(ctx context.Context, params *DbTxDiagnosticParams) (*DbTxDiagnosticResult, error) {
	const op = "service.Service.DbTxDiagnostic"

	logger := logging.LogWithTrace(ctx, s.logger)

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.params", fmt.Sprintf("%+v", params)),
	)
	logger.Info().
		Str("op", op).
		Str("params", fmt.Sprintf("%+v", params)).
		Msg("")

	result := &DbTxDiagnosticResult{}

	_, err := s.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(q *sqlc.Queries) error {
			// 1. Insert diagnostic row (stamped with this service's app name).
			inserted, err := q.InsertDbTxDiagnostic(ctx, sqlc.InsertDbTxDiagnosticParams{
				Service: s.appName,
				Message: params.Message,
			})
			if err != nil {
				return fmt.Errorf("insert diagnostic: %w", postgres_store.WrapDBError(err))
			}

			// 2. Read it back inside the same transaction.
			row, err := q.GetDbTxDiagnostic(ctx, inserted.ID)
			if err != nil {
				return fmt.Errorf("get diagnostic: %w", postgres_store.WrapDBError(err))
			}

			// 3. Sanity-check commit-visible state.
			if inserted.Message != row.Message {
				return errors.Join(apperrors.ErrInternal,
					fmt.Errorf("diagnostic message mismatch: %s != %s", inserted.Message, row.Message))
			}

			result.ID = row.ID
			return nil
		},
	})
	if err != nil {
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	logger.Info().Str("result", fmt.Sprintf("%+v", result)).Msg("")
	span.SetAttributes(attribute.String("output.result", fmt.Sprintf("%+v", result)))
	span.SetStatus(codes.Ok, "success")

	return result, nil
}
