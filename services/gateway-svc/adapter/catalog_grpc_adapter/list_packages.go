package catalog_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/catalog_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ListPackages delegates to catalog-svc.ListPackages over gRPC and
// translates proto messages into adapter-local Go types.
func (a *Adapter) ListPackages(ctx context.Context, params *ListPackagesParams) (*ListPackagesResult, error) {
	const op = "catalog_grpc_adapter.Adapter.ListPackages"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "ListPackages"))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.catalogClient.ListPackages(ctx, &pb.ListPackagesRequest{
		Kind:          params.Kind,
		AirlineCode:   params.AirlineCode,
		HotelId:       params.HotelID,
		DepartureFrom: params.DepartureFrom,
		DepartureTo:   params.DepartureTo,
		Cursor:        params.Cursor,
		Limit:         int32(params.Limit),
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	packages := make([]Package, 0, len(resp.GetPackages()))
	for _, p := range resp.GetPackages() {
		packages = append(packages, fromProtoPackageListItem(p))
	}

	result := &ListPackagesResult{
		Packages: packages,
		Page: PageMeta{
			NextCursor: resp.GetPage().GetNextCursor(),
			HasMore:    resp.GetPage().GetHasMore(),
		},
	}
	span.SetStatus(codes.Ok, "success")
	return result, nil
}
