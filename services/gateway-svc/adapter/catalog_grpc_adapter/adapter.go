// Package catalog_grpc_adapter is gateway-svc's consumer-side wrapper
// around catalog-svc's gRPC surface. It shields the rest of gateway-svc
// from proto types so the REST service layer + handlers stay transport-
// neutral.
//
// Per ADR 0009, gateway's public-read catalog routes (/v1/packages,
// /v1/packages/{id}, /v1/package-departures/{id}) proxy to
// catalog-svc.CatalogService via this adapter. The wire contract is in
// pb/catalog.proto, kept in sync by hand with
// services/catalog-svc/api/grpc_api/pb/catalog.proto.
//
// Landed with BL-GTW-002 / S1-E-10; catalog's REST surface is removed
// in the follow-up BL-REFACTOR-001 / S1-E-11.
package catalog_grpc_adapter

import (
	"gateway-svc/adapter/catalog_grpc_adapter/pb"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Adapter is a thin wrapper over a catalog.v1.CatalogService client.
type Adapter struct {
	logger *zerolog.Logger
	tracer trace.Tracer

	catalogClient          pb.CatalogServiceClient
	catalogMastersClient   pb.CatalogMastersClient
	catalogReadinessClient pb.CatalogReadinessClient
}

// NewAdapter creates a new catalog-svc gRPC adapter from an already-dialled
// conn. Ownership of the conn stays with the caller (shared pool lifetime).
func NewAdapter(logger *zerolog.Logger, tracer trace.Tracer, cc *grpc.ClientConn) *Adapter {
	return &Adapter{
		logger:                 logger,
		tracer:                 tracer,
		catalogClient:          pb.NewCatalogServiceClient(cc),
		catalogMastersClient:   pb.NewCatalogMastersClient(cc),
		catalogReadinessClient: pb.NewCatalogReadinessClient(cc),
	}
}
