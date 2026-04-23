// catalog_write.go — gateway-side adapter methods for catalog-svc write RPCs.
// S1-E-07 / BL-CAT-014 — StaffCatalogWrite MVP.
//
// Each method translates gateway-local params → pb request, forwards via gRPC,
// and translates pb response → adapter-local types. Proto types do not leak
// past this package; the rest of gateway-svc sees plain Go structs from types.go.
//
// Permission check: gateway handlers MUST call iam_grpc_adapter.CheckPermission
// with resource="catalog.package", action="manage", scope="global" BEFORE
// calling these methods. catalog-svc also enforces the same check server-side
// as defense-in-depth.
package catalog_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/catalog_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// CreatePackage
// ---------------------------------------------------------------------------

// CreatePackage forwards a staff create-package request to catalog-svc over
// gRPC and returns the created PackageDetail.
func (a *Adapter) CreatePackage(ctx context.Context, params *CreatePackageParams) (*PackageDetail, error) {
	const op = "catalog_grpc_adapter.Adapter.CreatePackage"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "CreatePackage"))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.catalogClient.CreatePackage(ctx, &pb.CreatePackageRequest{
		UserId:        params.UserID,
		BranchId:      params.BranchID,
		Kind:          params.Kind,
		Name:          params.Name,
		Description:   params.Description,
		CoverPhotoUrl: params.CoverPhotoURL,
		Highlights:    params.Highlights,
		ItineraryId:   params.ItineraryID,
		AirlineId:     params.AirlineID,
		MuthawwifId:   params.MuthawwifID,
		HotelIds:      params.HotelIDs,
		AddonIds:      params.AddonIDs,
		Status:        params.Status,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return fromProtoPackageDetail(resp.GetPackage()), nil
}

// ---------------------------------------------------------------------------
// UpdatePackage
// ---------------------------------------------------------------------------

// UpdatePackage forwards a staff update-package request to catalog-svc.
func (a *Adapter) UpdatePackage(ctx context.Context, params *UpdatePackageParams) (*PackageDetail, error) {
	const op = "catalog_grpc_adapter.Adapter.UpdatePackage"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "UpdatePackage"),
		attribute.String("id", params.ID))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.catalogClient.UpdatePackage(ctx, &pb.UpdatePackageRequest{
		UserId:        params.UserID,
		BranchId:      params.BranchID,
		Id:            params.ID,
		Name:          params.Name,
		Description:   params.Description,
		CoverPhotoUrl: params.CoverPhotoURL,
		Highlights:    params.Highlights,
		ItineraryId:   params.ItineraryID,
		AirlineId:     params.AirlineID,
		MuthawwifId:   params.MuthawwifID,
		HotelIds:      params.HotelIDs,
		AddonIds:      params.AddonIDs,
		Status:        params.Status,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return fromProtoPackageDetail(resp.GetPackage()), nil
}

// ---------------------------------------------------------------------------
// DeletePackage
// ---------------------------------------------------------------------------

// DeletePackage soft-deletes a catalog package via catalog-svc.
func (a *Adapter) DeletePackage(ctx context.Context, params *DeletePackageParams) (*DeletePackageResult, error) {
	const op = "catalog_grpc_adapter.Adapter.DeletePackage"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "DeletePackage"),
		attribute.String("id", params.ID))

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.catalogClient.DeletePackage(ctx, &pb.DeletePackageRequest{
		UserId:   params.UserID,
		BranchId: params.BranchID,
		Id:       params.ID,
	})
	if err != nil {
		wrapped := mapCatalogError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "success")
	return &DeletePackageResult{OK: resp.GetOk()}, nil
}

// ---------------------------------------------------------------------------
// CreateDeparture
// ---------------------------------------------------------------------------

// CreateDeparture forwards a staff create-departure request to catalog-svc.
func (a *Adapter) CreateDeparture(ctx context.Context, params *CreateDepartureParams) (*DepartureDetail, error) {
	const op = "catalog_grpc_adapter.Adapter.CreateDeparture"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "CreateDeparture"),
		attribute.String("package_id", params.PackageID))

	logger := logging.LogWithTrace(ctx, a.logger)

	pbPricing := make([]*pb.PricingInput, 0, len(params.Pricing))
	for _, p := range params.Pricing {
		pbPricing = append(pbPricing, &pb.PricingInput{
			RoomType:           p.RoomType,
			ListAmount:         p.ListAmount,
			ListCurrency:       p.ListCurrency,
			SettlementCurrency: p.SettlementCurrency,
		})
	}

	resp, err := a.catalogClient.CreateDeparture(ctx, &pb.CreateDepartureRequest{
		UserId:        params.UserID,
		BranchId:      params.BranchID,
		PackageId:     params.PackageID,
		DepartureDate: params.DepartureDate,
		ReturnDate:    params.ReturnDate,
		TotalSeats:    params.TotalSeats,
		Status:        params.Status,
		Pricing:       pbPricing,
	})
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

// ---------------------------------------------------------------------------
// UpdateDeparture
// ---------------------------------------------------------------------------

// UpdateDeparture forwards a staff update-departure request to catalog-svc.
func (a *Adapter) UpdateDeparture(ctx context.Context, params *UpdateDepartureParams) (*DepartureDetail, error) {
	const op = "catalog_grpc_adapter.Adapter.UpdateDeparture"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(attribute.String("operation", op), attribute.String("rpc", "UpdateDeparture"),
		attribute.String("id", params.ID))

	logger := logging.LogWithTrace(ctx, a.logger)

	var pbPricing []*pb.PricingInput
	if params.Pricing != nil {
		pbPricing = make([]*pb.PricingInput, 0, len(params.Pricing))
		for _, p := range params.Pricing {
			pbPricing = append(pbPricing, &pb.PricingInput{
				RoomType:           p.RoomType,
				ListAmount:         p.ListAmount,
				ListCurrency:       p.ListCurrency,
				SettlementCurrency: p.SettlementCurrency,
			})
		}
	}

	resp, err := a.catalogClient.UpdateDeparture(ctx, &pb.UpdateDepartureRequest{
		UserId:        params.UserID,
		BranchId:      params.BranchID,
		Id:            params.ID,
		DepartureDate: params.DepartureDate,
		ReturnDate:    params.ReturnDate,
		TotalSeats:    params.TotalSeats,
		Status:        params.Status,
		Pricing:       pbPricing,
	})
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
