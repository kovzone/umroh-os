// Package finance_grpc_adapter is gateway-svc's consumer-side wrapper around
// finance-svc's gRPC surface. It shields the rest of gateway-svc from proto
// types so the REST service layer + handlers stay transport-neutral.
//
// Per ADR 0009, gateway's /v1/finance/ping proxies to
// finance-svc.FinanceService.FinancePing via this adapter after the
// RequireBearerToken + RequirePermission middleware chain has approved the
// request. The wire contract is in pb/finance.proto, kept in sync by hand
// with services/finance-svc/api/grpc_api/pb/finance.proto.
//
// Landed with BL-IAM-019 / S1-E-14 (the last REST → gRPC migration in the
// ADR 0009 refactor arc).
package finance_grpc_adapter

import (
	"gateway-svc/adapter/finance_grpc_adapter/pb"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Adapter is a thin wrapper over a finance.v1.FinanceService client.
type Adapter struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	financeClient pb.FinanceServiceClient
}

// NewAdapter creates a new finance-svc gRPC adapter from an already-dialled
// conn. Ownership of the conn stays with the caller (shared pool lifetime).
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:        logger,
		tracer:        tracer,
		financeClient: pb.NewFinanceServiceClient(cc),
	}
}
