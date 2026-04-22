package service

import (
	"context"
	"fmt"

	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/adapter/catalog_rest_adapter"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetCatalogSystemLive proxies a liveness probe to catalog-svc through the REST
// adapter. Scaffold-time proof of the REST adapter pattern. Retires with
// BL-REFACTOR-001 / S1-E-11 when catalog-svc drops its REST port.
func (s *Service) GetCatalogSystemLive(ctx context.Context) (*catalog_rest_adapter.LivenessResult, error) {
	const op = "service.Service.GetCatalogSystemLive"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("operation", op))
	logger.Info().Str("op", op).Msg("")

	result, err := s.adapters.catalogRest.GetSystemLive(ctx)
	if err != nil {
		err = fmt.Errorf("call catalog-svc: %w", err)
		logger.Error().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// ListPackages proxies GET /v1/packages to catalog-svc.ListPackages (gRPC)
// via catalog_grpc_adapter. Thin — all business logic stays in catalog-svc.
func (s *Service) ListPackages(ctx context.Context, params *catalog_grpc_adapter.ListPackagesParams) (*catalog_grpc_adapter.ListPackagesResult, error) {
	const op = "service.Service.ListPackages"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.catalogGrpc.ListPackages(ctx, params)
	if err != nil {
		// Error already shaped by adapter (apperrors-wrapped). Log at the
		// level appropriate to the eventual HTTP status.
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// GetPackage proxies GET /v1/packages/{id} to catalog-svc.GetPackage (gRPC).
func (s *Service) GetPackage(ctx context.Context, params *catalog_grpc_adapter.GetPackageParams) (*catalog_grpc_adapter.PackageDetail, error) {
	const op = "service.Service.GetPackage"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("id", params.ID))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.catalogGrpc.GetPackage(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// GetPackageDeparture proxies GET /v1/package-departures/{id} to
// catalog-svc.GetPackageDeparture (gRPC).
func (s *Service) GetPackageDeparture(ctx context.Context, params *catalog_grpc_adapter.GetPackageDepartureParams) (*catalog_grpc_adapter.DepartureDetail, error) {
	const op = "service.Service.GetPackageDeparture"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("id", params.ID))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.catalogGrpc.GetPackageDeparture(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}
