package service

import (
	"context"
	"fmt"

	"catalog-svc/store/postgres_store"
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
