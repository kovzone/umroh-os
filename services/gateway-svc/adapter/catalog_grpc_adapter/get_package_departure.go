package catalog_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/catalog_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetPackageDeparture delegates to catalog-svc.GetPackageDeparture over
// gRPC. Returns ErrNotFound (mapped from gRPC NotFound) for
// departed / completed / cancelled / unknown.
func (a *Adapter) GetPackageDeparture(ctx context.Context, params *GetPackageDepartureParams) (*DepartureDetail, error) {
	const op = "catalog_grpc_adapter.Adapter.GetPackageDeparture"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "GetPackageDeparture"),
		attribute.String("id", params.ID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.catalogClient.GetPackageDeparture(ctx, &pb.GetPackageDepartureRequest{Id: params.ID})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return fromProtoDepartureDetail(resp.GetDeparture()), nil
}
