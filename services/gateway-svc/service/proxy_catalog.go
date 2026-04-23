package service

import (
	"context"

	"gateway-svc/adapter/catalog_grpc_adapter"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

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

// ---------------------------------------------------------------------------
// Staff catalog write (BL-CAT-014 / S1-E-07)
// ---------------------------------------------------------------------------

// CreatePackage proxies POST /v1/packages to catalog-svc.CreatePackage (gRPC).
// Permission gate (catalog.package.manage) is enforced by the HTTP handler
// before this method is called; catalog-svc also enforces it server-side.
func (s *Service) CreatePackage(ctx context.Context, params *catalog_grpc_adapter.CreatePackageParams) (*catalog_grpc_adapter.PackageDetail, error) {
	const op = "service.Service.CreatePackage"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.catalogGrpc.CreatePackage(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// UpdatePackage proxies PUT /v1/packages/{id} to catalog-svc.UpdatePackage (gRPC).
func (s *Service) UpdatePackage(ctx context.Context, params *catalog_grpc_adapter.UpdatePackageParams) (*catalog_grpc_adapter.PackageDetail, error) {
	const op = "service.Service.UpdatePackage"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("id", params.ID))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.catalogGrpc.UpdatePackage(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// DeletePackage proxies DELETE /v1/packages/{id} to catalog-svc.DeletePackage (gRPC).
func (s *Service) DeletePackage(ctx context.Context, params *catalog_grpc_adapter.DeletePackageParams) (*catalog_grpc_adapter.DeletePackageResult, error) {
	const op = "service.Service.DeletePackage"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("id", params.ID))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.catalogGrpc.DeletePackage(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// CreateDeparture proxies POST /v1/packages/{id}/departures to
// catalog-svc.CreateDeparture (gRPC).
func (s *Service) CreateDeparture(ctx context.Context, params *catalog_grpc_adapter.CreateDepartureParams) (*catalog_grpc_adapter.DepartureDetail, error) {
	const op = "service.Service.CreateDeparture"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("package_id", params.PackageID))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.catalogGrpc.CreateDeparture(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}

// UpdateDeparture proxies PUT /v1/departures/{id} to
// catalog-svc.UpdateDeparture (gRPC).
func (s *Service) UpdateDeparture(ctx context.Context, params *catalog_grpc_adapter.UpdateDepartureParams) (*catalog_grpc_adapter.DepartureDetail, error) {
	const op = "service.Service.UpdateDeparture"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("id", params.ID))

	logger := logging.LogWithTrace(ctx, s.logger)

	result, err := s.adapters.catalogGrpc.UpdateDeparture(ctx, params)
	if err != nil {
		logger.Warn().Err(err).Msg("")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "success")
	return result, nil
}
