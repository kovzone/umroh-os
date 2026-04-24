// catalog_bulk.go — gateway catalog adapter methods for bulk + versioning RPCs.
// BL-CAT-010: BulkImportPackages
// BL-CAT-011: BulkUpdatePackages
// BL-CAT-013: GetPackageVersion

package catalog_grpc_adapter

import (
	"context"
	"fmt"

	"gateway-svc/adapter/catalog_grpc_adapter/pb"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// Result types
// ---------------------------------------------------------------------------

// BulkImportRowResult is the per-row outcome for a bulk import.
type BulkImportRowResult struct {
	Index     int
	PackageID string
	Error     string
}

// BulkImportPackagesResult is the gateway-layer result for BulkImportPackages.
type BulkImportPackagesResult struct {
	Results    []BulkImportRowResult
	TotalRows  int
	Successful int
	Failed     int
}

// BulkImportRowInput is one row in a bulk import request.
type BulkImportRowInput struct {
	Name          string
	Kind          string
	Description   string
	CoverPhotoURL string
	Highlights    []string
	AddonIDs      []string
	HotelIDs      []string
	ItineraryID   string
	AirlineID     string
	MuthawwifID   string
	Status        string
}

// BulkImportPackagesParams is the gateway-layer input for BulkImportPackages.
type BulkImportPackagesParams struct {
	UserID   string
	BranchID string
	Rows     []BulkImportRowInput
}

// BulkUpdateRowResult is the per-row outcome for a bulk update.
type BulkUpdateRowResult struct {
	Index     int
	PackageID string
	Error     string
}

// BulkUpdatePackagesResult is the gateway-layer result for BulkUpdatePackages.
type BulkUpdatePackagesResult struct {
	Results    []BulkUpdateRowResult
	TotalRows  int
	Successful int
	Failed     int
}

// BulkUpdateRowInput is one row in a bulk update request.
type BulkUpdateRowInput struct {
	ID          string
	Name        string
	Description string
	Status      string
	Highlights  []string
	AddonIDs    []string
	HotelIDs    []string
}

// BulkUpdatePackagesParams is the gateway-layer input for BulkUpdatePackages.
type BulkUpdatePackagesParams struct {
	UserID   string
	BranchID string
	Rows     []BulkUpdateRowInput
}

// PackageVersionResult is the gateway-layer result for GetPackageVersion.
type PackageVersionResult struct {
	PackageID string
	Version   string
	UpdatedAt string
}

// ---------------------------------------------------------------------------
// Error mapper (mirrors mapCatalogError from catalog_write.go)
// ---------------------------------------------------------------------------

func mapCatalogBulkError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return apperrors.ErrInternal
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return fmt.Errorf("%w: %s", apperrors.ErrNotFound, st.Message())
	case grpcCodes.InvalidArgument:
		return fmt.Errorf("%w: %s", apperrors.ErrValidation, st.Message())
	case grpcCodes.PermissionDenied:
		return fmt.Errorf("%w: %s", apperrors.ErrForbidden, st.Message())
	default:
		return apperrors.ErrInternal
	}
}

// ---------------------------------------------------------------------------
// BL-CAT-010: BulkImportPackages
// ---------------------------------------------------------------------------

// BulkImportPackages sends a batch of package definitions to catalog-svc for
// creation. Per-row results are returned; no single-row failure aborts others.
func (a *Adapter) BulkImportPackages(ctx context.Context, params *BulkImportPackagesParams) (*BulkImportPackagesResult, error) {
	const op = "catalog_grpc_adapter.Adapter.BulkImportPackages"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Int("rows", len(params.Rows)).Msg("")

	rows := make([]*pb.BulkImportRow, 0, len(params.Rows))
	for _, r := range params.Rows {
		rows = append(rows, &pb.BulkImportRow{
			Name:          r.Name,
			Kind:          r.Kind,
			Description:   r.Description,
			CoverPhotoUrl: r.CoverPhotoURL,
			Highlights:    r.Highlights,
			AddonIds:      r.AddonIDs,
			HotelIds:      r.HotelIDs,
			ItineraryId:   r.ItineraryID,
			AirlineId:     r.AirlineID,
			MuthawwifId:   r.MuthawwifID,
			Status:        r.Status,
		})
	}

	resp, err := a.catalogBulkClient.BulkImportPackages(ctx, &pb.BulkImportPackagesRequest{
		UserId:   params.UserID,
		BranchId: params.BranchID,
		Rows:     rows,
	})
	if err != nil {
		wrapped := mapCatalogBulkError(err)
		logger.Warn().Err(wrapped).Msg("catalog-svc.BulkImportPackages failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	results := make([]BulkImportRowResult, 0, len(resp.GetResults()))
	for _, r := range resp.GetResults() {
		results = append(results, BulkImportRowResult{
			Index:     int(r.GetIndex()),
			PackageID: r.GetPackageId(),
			Error:     r.GetError(),
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &BulkImportPackagesResult{
		Results:    results,
		TotalRows:  int(resp.GetTotalRows()),
		Successful: int(resp.GetSuccessful()),
		Failed:     int(resp.GetFailed()),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-CAT-011: BulkUpdatePackages
// ---------------------------------------------------------------------------

// BulkUpdatePackages applies mass field updates to existing packages.
func (a *Adapter) BulkUpdatePackages(ctx context.Context, params *BulkUpdatePackagesParams) (*BulkUpdatePackagesResult, error) {
	const op = "catalog_grpc_adapter.Adapter.BulkUpdatePackages"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Int("rows", len(params.Rows)).Msg("")

	rows := make([]*pb.BulkUpdateRow, 0, len(params.Rows))
	for _, r := range params.Rows {
		rows = append(rows, &pb.BulkUpdateRow{
			Id:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Status:      r.Status,
			Highlights:  r.Highlights,
			AddonIds:    r.AddonIDs,
			HotelIds:    r.HotelIDs,
		})
	}

	resp, err := a.catalogBulkClient.BulkUpdatePackages(ctx, &pb.BulkUpdatePackagesRequest{
		UserId:   params.UserID,
		BranchId: params.BranchID,
		Rows:     rows,
	})
	if err != nil {
		wrapped := mapCatalogBulkError(err)
		logger.Warn().Err(wrapped).Msg("catalog-svc.BulkUpdatePackages failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	results := make([]BulkUpdateRowResult, 0, len(resp.GetResults()))
	for _, r := range resp.GetResults() {
		results = append(results, BulkUpdateRowResult{
			Index:     int(r.GetIndex()),
			PackageID: r.GetPackageId(),
			Error:     r.GetError(),
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &BulkUpdatePackagesResult{
		Results:    results,
		TotalRows:  int(resp.GetTotalRows()),
		Successful: int(resp.GetSuccessful()),
		Failed:     int(resp.GetFailed()),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-CAT-013: GetPackageVersion
// ---------------------------------------------------------------------------

// GetPackageVersion returns the computed content-hash version for a package.
func (a *Adapter) GetPackageVersion(ctx context.Context, packageID string) (*PackageVersionResult, error) {
	const op = "catalog_grpc_adapter.Adapter.GetPackageVersion"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, a.logger)
	logger.Info().Str("op", op).Str("package_id", packageID).Msg("")

	resp, err := a.catalogBulkClient.GetPackageVersion(ctx, &pb.GetPackageVersionRequest{
		PackageId: packageID,
	})
	if err != nil {
		wrapped := mapCatalogBulkError(err)
		logger.Warn().Err(wrapped).Msg("catalog-svc.GetPackageVersion failed")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &PackageVersionResult{
		PackageID: resp.GetPackageId(),
		Version:   resp.GetVersion(),
		UpdatedAt: resp.GetUpdatedAt(),
	}, nil
}
