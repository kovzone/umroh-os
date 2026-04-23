// catalog_write_messages.go — hand-written proto message types for the
// S1-E-07 / BL-CAT-014 staff write quartet.
//
// NOTE: needs make genpb — these types would normally be generated into
// catalog.pb.go by protoc-gen-go. They are written by hand because protoc
// is not available in the current environment. Run `make genpb` to regenerate
// catalog.pb.go from catalog.proto, then delete this file.
//
// Types follow the same getter-method pattern as protoc-gen-go output so that
// the adapter code in catalog_write.go compiles unchanged once regenerated.
package pb

// ---------------------------------------------------------------------------
// PricingInput — shared by CreateDeparture and UpdateDeparture.
// ---------------------------------------------------------------------------

// PricingInput is the staff-side input for one room-type price row.
type PricingInput struct {
	RoomType           string
	ListAmount         int64
	ListCurrency       string
	SettlementCurrency string
}

func (x *PricingInput) GetRoomType() string {
	if x == nil {
		return ""
	}
	return x.RoomType
}
func (x *PricingInput) GetListAmount() int64 {
	if x == nil {
		return 0
	}
	return x.ListAmount
}
func (x *PricingInput) GetListCurrency() string {
	if x == nil {
		return ""
	}
	return x.ListCurrency
}
func (x *PricingInput) GetSettlementCurrency() string {
	if x == nil {
		return ""
	}
	return x.SettlementCurrency
}

// ---------------------------------------------------------------------------
// CreatePackage
// ---------------------------------------------------------------------------

// CreatePackageRequest mirrors the proto message.
type CreatePackageRequest struct {
	UserId        string
	BranchId      string
	Kind          string
	Name          string
	Description   string
	CoverPhotoUrl string
	Highlights    []string
	ItineraryId   string
	AirlineId     string
	MuthawwifId   string
	HotelIds      []string
	AddonIds      []string
	Status        string
}

func (x *CreatePackageRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *CreatePackageRequest) GetBranchId() string {
	if x == nil {
		return ""
	}
	return x.BranchId
}
func (x *CreatePackageRequest) GetKind() string {
	if x == nil {
		return ""
	}
	return x.Kind
}
func (x *CreatePackageRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *CreatePackageRequest) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *CreatePackageRequest) GetCoverPhotoUrl() string {
	if x == nil {
		return ""
	}
	return x.CoverPhotoUrl
}
func (x *CreatePackageRequest) GetHighlights() []string {
	if x == nil {
		return nil
	}
	return x.Highlights
}
func (x *CreatePackageRequest) GetItineraryId() string {
	if x == nil {
		return ""
	}
	return x.ItineraryId
}
func (x *CreatePackageRequest) GetAirlineId() string {
	if x == nil {
		return ""
	}
	return x.AirlineId
}
func (x *CreatePackageRequest) GetMuthawwifId() string {
	if x == nil {
		return ""
	}
	return x.MuthawwifId
}
func (x *CreatePackageRequest) GetHotelIds() []string {
	if x == nil {
		return nil
	}
	return x.HotelIds
}
func (x *CreatePackageRequest) GetAddonIds() []string {
	if x == nil {
		return nil
	}
	return x.AddonIds
}
func (x *CreatePackageRequest) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

// CreatePackageResponse wraps a PackageDetail.
type CreatePackageResponse struct {
	Package *PackageDetail
}

func (x *CreatePackageResponse) GetPackage() *PackageDetail {
	if x == nil {
		return nil
	}
	return x.Package
}

// ---------------------------------------------------------------------------
// UpdatePackage
// ---------------------------------------------------------------------------

// UpdatePackageRequest mirrors the proto message.
type UpdatePackageRequest struct {
	UserId        string
	BranchId      string
	Id            string
	Name          string
	Description   string
	CoverPhotoUrl string
	Highlights    []string
	ItineraryId   string
	AirlineId     string
	MuthawwifId   string
	HotelIds      []string
	AddonIds      []string
	Status        string
}

func (x *UpdatePackageRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *UpdatePackageRequest) GetBranchId() string {
	if x == nil {
		return ""
	}
	return x.BranchId
}
func (x *UpdatePackageRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *UpdatePackageRequest) GetName() string {
	if x == nil {
		return ""
	}
	return x.Name
}
func (x *UpdatePackageRequest) GetDescription() string {
	if x == nil {
		return ""
	}
	return x.Description
}
func (x *UpdatePackageRequest) GetCoverPhotoUrl() string {
	if x == nil {
		return ""
	}
	return x.CoverPhotoUrl
}
func (x *UpdatePackageRequest) GetHighlights() []string {
	if x == nil {
		return nil
	}
	return x.Highlights
}
func (x *UpdatePackageRequest) GetItineraryId() string {
	if x == nil {
		return ""
	}
	return x.ItineraryId
}
func (x *UpdatePackageRequest) GetAirlineId() string {
	if x == nil {
		return ""
	}
	return x.AirlineId
}
func (x *UpdatePackageRequest) GetMuthawwifId() string {
	if x == nil {
		return ""
	}
	return x.MuthawwifId
}
func (x *UpdatePackageRequest) GetHotelIds() []string {
	if x == nil {
		return nil
	}
	return x.HotelIds
}
func (x *UpdatePackageRequest) GetAddonIds() []string {
	if x == nil {
		return nil
	}
	return x.AddonIds
}
func (x *UpdatePackageRequest) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

// UpdatePackageResponse wraps a PackageDetail.
type UpdatePackageResponse struct {
	Package *PackageDetail
}

func (x *UpdatePackageResponse) GetPackage() *PackageDetail {
	if x == nil {
		return nil
	}
	return x.Package
}

// ---------------------------------------------------------------------------
// DeletePackage
// ---------------------------------------------------------------------------

// DeletePackageRequest mirrors the proto message.
type DeletePackageRequest struct {
	UserId   string
	BranchId string
	Id       string
}

func (x *DeletePackageRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *DeletePackageRequest) GetBranchId() string {
	if x == nil {
		return ""
	}
	return x.BranchId
}
func (x *DeletePackageRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}

// DeletePackageResponse carries the ok flag.
type DeletePackageResponse struct {
	Ok bool
}

func (x *DeletePackageResponse) GetOk() bool {
	if x == nil {
		return false
	}
	return x.Ok
}

// ---------------------------------------------------------------------------
// CreateDeparture
// ---------------------------------------------------------------------------

// CreateDepartureRequest mirrors the proto message.
type CreateDepartureRequest struct {
	UserId        string
	BranchId      string
	PackageId     string
	DepartureDate string
	ReturnDate    string
	TotalSeats    int32
	Status        string
	Pricing       []*PricingInput
}

func (x *CreateDepartureRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *CreateDepartureRequest) GetBranchId() string {
	if x == nil {
		return ""
	}
	return x.BranchId
}
func (x *CreateDepartureRequest) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}
func (x *CreateDepartureRequest) GetDepartureDate() string {
	if x == nil {
		return ""
	}
	return x.DepartureDate
}
func (x *CreateDepartureRequest) GetReturnDate() string {
	if x == nil {
		return ""
	}
	return x.ReturnDate
}
func (x *CreateDepartureRequest) GetTotalSeats() int32 {
	if x == nil {
		return 0
	}
	return x.TotalSeats
}
func (x *CreateDepartureRequest) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *CreateDepartureRequest) GetPricing() []*PricingInput {
	if x == nil {
		return nil
	}
	return x.Pricing
}

// CreateDepartureResponse wraps a DepartureDetail.
type CreateDepartureResponse struct {
	Departure *DepartureDetail
}

func (x *CreateDepartureResponse) GetDeparture() *DepartureDetail {
	if x == nil {
		return nil
	}
	return x.Departure
}

// ---------------------------------------------------------------------------
// UpdateDeparture
// ---------------------------------------------------------------------------

// UpdateDepartureRequest mirrors the proto message.
type UpdateDepartureRequest struct {
	UserId        string
	BranchId      string
	Id            string
	DepartureDate string
	ReturnDate    string
	TotalSeats    int32
	Status        string
	Pricing       []*PricingInput
}

func (x *UpdateDepartureRequest) GetUserId() string {
	if x == nil {
		return ""
	}
	return x.UserId
}
func (x *UpdateDepartureRequest) GetBranchId() string {
	if x == nil {
		return ""
	}
	return x.BranchId
}
func (x *UpdateDepartureRequest) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *UpdateDepartureRequest) GetDepartureDate() string {
	if x == nil {
		return ""
	}
	return x.DepartureDate
}
func (x *UpdateDepartureRequest) GetReturnDate() string {
	if x == nil {
		return ""
	}
	return x.ReturnDate
}
func (x *UpdateDepartureRequest) GetTotalSeats() int32 {
	if x == nil {
		return 0
	}
	return x.TotalSeats
}
func (x *UpdateDepartureRequest) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *UpdateDepartureRequest) GetPricing() []*PricingInput {
	if x == nil {
		return nil
	}
	return x.Pricing
}

// UpdateDepartureResponse wraps a DepartureDetail.
type UpdateDepartureResponse struct {
	Departure *DepartureDetail
}

func (x *UpdateDepartureResponse) GetDeparture() *DepartureDetail {
	if x == nil {
		return nil
	}
	return x.Departure
}
