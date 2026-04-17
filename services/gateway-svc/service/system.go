package service

import (
	"context"
)

type LivenessParams struct{}

type LivenessResult struct {
	OK bool `json:"ok"`
}

// Liveness verifies that the service process is alive and responsive at the service layer.
// Trivial and deterministic; returns success as long as the service can execute and respond.
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
// gateway-svc has no required external dependency at scaffold time (it's a stateless edge proxy);
// future iterations may extend this to ping each wired backend adapter and aggregate the result.
//
// No logging or tracing — probe endpoints are hit at high frequency; kept cheap.
func (s *Service) Readiness(ctx context.Context, params *ReadinessParams) (*ReadinessResult, error) {
	return &ReadinessResult{OK: true}, nil
}
