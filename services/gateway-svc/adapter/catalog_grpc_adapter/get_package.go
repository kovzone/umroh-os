package catalog_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/catalog_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetPackage delegates to catalog-svc.GetPackage over gRPC. Returns
// ErrNotFound (mapped from gRPC NotFound) for draft / archived / unknown.
func (a *Adapter) GetPackage(ctx context.Context, params *GetPackageParams) (*PackageDetail, error) {
	const op = "catalog_grpc_adapter.Adapter.GetPackage"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "GetPackage"),
		attribute.String("id", params.ID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.catalogClient.GetPackage(ctx, &pb.GetPackageRequest{Id: params.ID})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	detail := fromProtoPackageDetail(resp.GetPackage())

	// Enrich each DepartureSummary with its minimum list price by calling
	// GetPackageDeparture for each departure. This is an intentional N+1:
	// packages typically have 2–4 departures, so the cost is negligible and
	// avoids a proto schema change on the hot path.
	for i, dep := range detail.Departures {
		depDetail, err := a.GetPackageDeparture(ctx, &GetPackageDepartureParams{ID: dep.ID})
		if err != nil || len(depDetail.Pricing) == 0 {
			continue
		}
		var minPrice int64 = depDetail.Pricing[0].ListAmount
		for _, p := range depDetail.Pricing[1:] {
			if p.ListAmount > 0 && p.ListAmount < minPrice {
				minPrice = p.ListAmount
			}
		}
		if minPrice > 0 {
			detail.Departures[i].PricePerPax = &minPrice
		}
	}

	return detail, nil
}
