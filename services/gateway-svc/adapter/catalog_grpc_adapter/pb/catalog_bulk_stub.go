// catalog_bulk_stub.go — gateway-side gRPC client stubs for catalog-svc bulk
// package RPCs (BL-CAT-010/011) and package versioning RPC (BL-CAT-013).
//
// Mirrors services/catalog-svc/api/grpc_api/pb/catalog_bulk_messages.go,
// catalog_bulk_grpc_ext.go, catalog_version_messages.go, and
// catalog_version_grpc_ext.go. Run `make genpb` to replace with generated
// code once protoc is available.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// Method name constants
// ---------------------------------------------------------------------------

const (
	CatalogService_BulkImportPackages_FullMethodName = "/pb.CatalogService/BulkImportPackages"
	CatalogService_BulkUpdatePackages_FullMethodName = "/pb.CatalogService/BulkUpdatePackages"
	CatalogService_GetPackageVersion_FullMethodName  = "/pb.CatalogService/GetPackageVersion"
)

// ---------------------------------------------------------------------------
// BulkImportPackages — domain types
// ---------------------------------------------------------------------------

// BulkImportRow is a single package definition in a bulk import payload.
type BulkImportRow struct {
	Name          string
	Kind          string
	Description   string
	CoverPhotoUrl string
	Highlights    []string
	AddonIds      []string
	HotelIds      []string
	ItineraryId   string
	AirlineId     string
	MuthawwifId   string
	Status        string
}

func (x *BulkImportRow) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *BulkImportRow) GetKind() string {
	if x == nil {
		return ""
	}
	return x.Kind
}
func (x *BulkImportRow) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *BulkImportRow) GetCoverPhotoUrl() string {
	if x == nil {
		return ""
	}
	return x.CoverPhotoUrl
}
func (x *BulkImportRow) GetHighlights() []string {
	if x == nil {
		return nil
	}
	return x.Highlights
}
func (x *BulkImportRow) GetAddonIds() []string {
	if x == nil {
		return nil
	}
	return x.AddonIds
}
func (x *BulkImportRow) GetHotelIds() []string {
	if x == nil {
		return nil
	}
	return x.HotelIds
}
func (x *BulkImportRow) GetItineraryId() string {
	if x == nil {
		return ""
	}
	return x.ItineraryId
}
func (x *BulkImportRow) GetAirlineId() string {
	if x == nil {
		return ""
	}
	return x.AirlineId
}
func (x *BulkImportRow) GetMuthawwifId() string {
	if x == nil {
		return ""
	}
	return x.MuthawwifId
}
func (x *BulkImportRow) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

// BulkImportRowResult is the per-row outcome for a bulk import.
type BulkImportRowResult struct {
	Index     int32
	PackageId string
	Error     string
}

func (x *BulkImportRowResult) GetIndex() int32 {
	if x == nil {
		return 0
	}
	return x.Index
}
func (x *BulkImportRowResult) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}
func (x *BulkImportRowResult) GetError() string {
	if x == nil {
		return ""
	}
	return x.Error
}

// BulkImportPackagesRequest is the gRPC request for BulkImportPackages.
type BulkImportPackagesRequest struct {
	UserId   string
	BranchId string
	Rows     []*BulkImportRow
}

func (x *BulkImportPackagesRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *BulkImportPackagesRequest) GetBranchId() string {
	if x == nil {
		return ""
	}
	return x.BranchId
}
func (x *BulkImportPackagesRequest) GetRows() []*BulkImportRow {
	if x == nil {
		return nil
	}
	return x.Rows
}

// BulkImportPackagesResponse is the gRPC response for BulkImportPackages.
type BulkImportPackagesResponse struct {
	Results    []*BulkImportRowResult
	TotalRows  int32
	Successful int32
	Failed     int32
}

func (x *BulkImportPackagesResponse) GetResults() []*BulkImportRowResult {
	if x == nil {
		return nil
	}
	return x.Results
}
func (x *BulkImportPackagesResponse) GetTotalRows() int32 {
	if x == nil {
		return 0
	}
	return x.TotalRows
}
func (x *BulkImportPackagesResponse) GetSuccessful() int32 {
	if x == nil {
		return 0
	}
	return x.Successful
}
func (x *BulkImportPackagesResponse) GetFailed() int32 {
	if x == nil {
		return 0
	}
	return x.Failed
}

// ---------------------------------------------------------------------------
// BulkUpdatePackages — domain types
// ---------------------------------------------------------------------------

// BulkUpdateRow is a single update spec in a bulk update payload.
type BulkUpdateRow struct {
	Id          string
	Name        string
	Description string
	Status      string
	Highlights  []string
	AddonIds    []string
	HotelIds    []string
}

func (x *BulkUpdateRow) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *BulkUpdateRow) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *BulkUpdateRow) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *BulkUpdateRow) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *BulkUpdateRow) GetHighlights() []string {
	if x == nil {
		return nil
	}
	return x.Highlights
}
func (x *BulkUpdateRow) GetAddonIds() []string {
	if x == nil {
		return nil
	}
	return x.AddonIds
}
func (x *BulkUpdateRow) GetHotelIds() []string {
	if x == nil {
		return nil
	}
	return x.HotelIds
}

// BulkUpdateRowResult is the per-row outcome for a bulk update.
type BulkUpdateRowResult struct {
	Index     int32
	PackageId string
	Error     string
}

func (x *BulkUpdateRowResult) GetIndex() int32 {
	if x == nil {
		return 0
	}
	return x.Index
}
func (x *BulkUpdateRowResult) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}
func (x *BulkUpdateRowResult) GetError() string {
	if x == nil {
		return ""
	}
	return x.Error
}

// BulkUpdatePackagesRequest is the gRPC request for BulkUpdatePackages.
type BulkUpdatePackagesRequest struct {
	UserId   string
	BranchId string
	Rows     []*BulkUpdateRow
}

func (x *BulkUpdatePackagesRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *BulkUpdatePackagesRequest) GetBranchId() string {
	if x == nil {
		return ""
	}
	return x.BranchId
}
func (x *BulkUpdatePackagesRequest) GetRows() []*BulkUpdateRow {
	if x == nil {
		return nil
	}
	return x.Rows
}

// BulkUpdatePackagesResponse is the gRPC response for BulkUpdatePackages.
type BulkUpdatePackagesResponse struct {
	Results    []*BulkUpdateRowResult
	TotalRows  int32
	Successful int32
	Failed     int32
}

func (x *BulkUpdatePackagesResponse) GetResults() []*BulkUpdateRowResult {
	if x == nil {
		return nil
	}
	return x.Results
}
func (x *BulkUpdatePackagesResponse) GetTotalRows() int32 {
	if x == nil {
		return 0
	}
	return x.TotalRows
}
func (x *BulkUpdatePackagesResponse) GetSuccessful() int32 {
	if x == nil {
		return 0
	}
	return x.Successful
}
func (x *BulkUpdatePackagesResponse) GetFailed() int32 {
	if x == nil {
		return 0
	}
	return x.Failed
}

// ---------------------------------------------------------------------------
// GetPackageVersion — domain types
// ---------------------------------------------------------------------------

// GetPackageVersionRequest is the gRPC request for GetPackageVersion.
type GetPackageVersionRequest struct {
	PackageId string
}

func (x *GetPackageVersionRequest) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}

// GetPackageVersionResponse is the gRPC response for GetPackageVersion.
type GetPackageVersionResponse struct {
	PackageId string
	Version   string
	UpdatedAt string
}

func (x *GetPackageVersionResponse) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}
func (x *GetPackageVersionResponse) GetVersion() string {
	if x == nil {
		return ""
	}
	return x.Version
}
func (x *GetPackageVersionResponse) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}

// ---------------------------------------------------------------------------
// CatalogBulkClient (gRPC client stub)
// ---------------------------------------------------------------------------

// CatalogBulkClient is the client API for bulk + versioning catalog RPCs.
type CatalogBulkClient interface {
	BulkImportPackages(ctx context.Context, in *BulkImportPackagesRequest, opts ...grpc.CallOption) (*BulkImportPackagesResponse, error)
	BulkUpdatePackages(ctx context.Context, in *BulkUpdatePackagesRequest, opts ...grpc.CallOption) (*BulkUpdatePackagesResponse, error)
	GetPackageVersion(ctx context.Context, in *GetPackageVersionRequest, opts ...grpc.CallOption) (*GetPackageVersionResponse, error)
}

type catalogBulkClient struct {
	cc grpc.ClientConnInterface
}

// NewCatalogBulkClient returns a CatalogBulkClient backed by the given conn.
func NewCatalogBulkClient(cc grpc.ClientConnInterface) CatalogBulkClient {
	return &catalogBulkClient{cc}
}

func (c *catalogBulkClient) BulkImportPackages(ctx context.Context, in *BulkImportPackagesRequest, opts ...grpc.CallOption) (*BulkImportPackagesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BulkImportPackagesResponse)
	if err := c.cc.Invoke(ctx, CatalogService_BulkImportPackages_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogBulkClient) BulkUpdatePackages(ctx context.Context, in *BulkUpdatePackagesRequest, opts ...grpc.CallOption) (*BulkUpdatePackagesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BulkUpdatePackagesResponse)
	if err := c.cc.Invoke(ctx, CatalogService_BulkUpdatePackages_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogBulkClient) GetPackageVersion(ctx context.Context, in *GetPackageVersionRequest, opts ...grpc.CallOption) (*GetPackageVersionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPackageVersionResponse)
	if err := c.cc.Invoke(ctx, CatalogService_GetPackageVersion_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
