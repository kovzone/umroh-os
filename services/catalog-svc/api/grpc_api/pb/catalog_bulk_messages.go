// catalog_bulk_messages.go — hand-written proto message types for bulk
// import/update RPCs (BL-CAT-010 / BL-CAT-011).
//
// Run `make genpb` to regenerate from catalog.proto once protoc is available,
// then delete this file.
//
// Types follow the same getter-method pattern as protoc-gen-go output so that
// handler code compiles unchanged once regenerated.
package pb

// ---------------------------------------------------------------------------
// BulkImportPackages
// ---------------------------------------------------------------------------

// BulkImportRow is a single row in a bulk import request.
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

// BulkImportRowResult is the per-row result for a bulk import.
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

// BulkImportPackagesRequest carries the rows to import.
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

// BulkImportPackagesResponse carries per-row results.
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
// BulkUpdatePackages
// ---------------------------------------------------------------------------

// BulkUpdateRow is a single row in a bulk update request.
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

// BulkUpdateRowResult is the per-row result for a bulk update.
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

// BulkUpdatePackagesRequest carries the rows to update.
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

// BulkUpdatePackagesResponse carries per-row results.
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
