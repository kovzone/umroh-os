package catalog_grpc_adapter

// Adapter-local types: proto types do not leak past this package. These
// mirror catalog-svc's service-layer shapes (and the response schemas in
// services/catalog-svc/api/rest_oapi/openapi.yaml) so the gateway REST
// handlers can marshal them straight into the generated oapi types.
// Dates are ISO YYYY-MM-DD strings; enums (kind, status, operator_kind,
// vendor_readiness_state) are strings that pass through unchanged.

// Money is the display-only currency triple.
type Money struct {
	ListAmount         int64
	ListCurrency       string
	SettlementCurrency string
}

// NextDeparture is the earliest open/closed departure attached to each
// list-row package (optional in the result).
type NextDeparture struct {
	ID             string
	DepartureDate  string
	ReturnDate     string
	RemainingSeats int
}

// Package is a list-row shape.
type Package struct {
	ID            string
	Kind          string
	Name          string
	Description   string
	CoverPhotoUrl string
	StartingPrice Money
	NextDeparture *NextDeparture // nil when no upcoming departure
}

// PageMeta is the pagination envelope for list responses.
type PageMeta struct {
	NextCursor string
	HasMore    bool
}

type HotelRef struct {
	ID               string
	Name             string
	City             string
	StarRating       int
	WalkingDistanceM int
}

type AirlineRef struct {
	ID           string
	Code         string
	Name         string
	OperatorKind string
}

type MuthawwifRef struct {
	ID          string
	Name        string
	PortraitUrl string
}

type AddonRef struct {
	ID                 string
	Name               string
	ListAmount         int64
	ListCurrency       string
	SettlementCurrency string
}

type ItineraryDay struct {
	Day         int
	Title       string
	Description string
	PhotoUrl    string
}

type Itinerary struct {
	ID        string
	Days      []ItineraryDay
	PublicUrl string
}

type DepartureSummary struct {
	ID             string
	DepartureDate  string
	ReturnDate     string
	RemainingSeats int
	Status         string
}

// PackageDetail is the full detail graph for GetPackage. Pointer-valued
// optional fields (Itinerary, Airline, Muthawwif) are nil when absent.
type PackageDetail struct {
	ID            string
	Kind          string
	Name          string
	Description   string
	Highlights    []string
	CoverPhotoUrl string
	Itinerary     *Itinerary
	Hotels        []HotelRef
	Airline       *AirlineRef
	Muthawwif     *MuthawwifRef
	Addons        []AddonRef
	Departures    []DepartureSummary
}

type PackagePricing struct {
	RoomType           string
	ListAmount         int64
	ListCurrency       string
	SettlementCurrency string
}

type VendorReadiness struct {
	Ticket string
	Hotel  string
	Visa   string
}

type DepartureDetail struct {
	ID              string
	PackageID       string
	DepartureDate   string
	ReturnDate      string
	TotalSeats      int
	RemainingSeats  int
	Status          string
	Pricing         []PackagePricing
	VendorReadiness VendorReadiness
}

// ---------------------------------------------------------------------------
// Adapter input / output envelopes.
// ---------------------------------------------------------------------------

type ListPackagesParams struct {
	Kind          string
	AirlineCode   string
	HotelID       string
	DepartureFrom string
	DepartureTo   string
	Cursor        string
	Limit         int
}

type ListPackagesResult struct {
	Packages []Package
	Page     PageMeta
}

type GetPackageParams struct {
	ID string
}

type GetPackageDepartureParams struct {
	ID string
}

// ---------------------------------------------------------------------------
// Write quartet input / output types (S1-E-07 / BL-CAT-014).
// ---------------------------------------------------------------------------

// PricingInputParam is the adapter-level (proto-free) shape for one pricing row.
type PricingInputParam struct {
	RoomType           string // "double" | "triple" | "quad"
	ListAmount         int64
	ListCurrency       string // "IDR" | "USD"
	SettlementCurrency string // always "IDR" in MVP
}

// CreatePackageParams is the adapter input for CreatePackage.
type CreatePackageParams struct {
	UserID        string
	BranchID      string
	Kind          string
	Name          string
	Description   string
	CoverPhotoURL string
	Highlights    []string
	ItineraryID   string
	AirlineID     string
	MuthawwifID   string
	HotelIDs      []string
	AddonIDs      []string
	Status        string // "draft" | "active" — defaults to "draft"
}

// UpdatePackageParams is the adapter input for UpdatePackage.
type UpdatePackageParams struct {
	UserID        string
	BranchID      string
	ID            string
	Name          string
	Description   string
	CoverPhotoURL string
	Highlights    []string // nil = no change
	ItineraryID   string
	AirlineID     string
	MuthawwifID   string
	HotelIDs      []string // nil = no change; non-nil = replace
	AddonIDs      []string // nil = no change; non-nil = replace
	Status        string
}

// DeletePackageParams is the adapter input for DeletePackage.
type DeletePackageParams struct {
	UserID   string
	BranchID string
	ID       string
}

// DeletePackageResult is the adapter output for DeletePackage.
type DeletePackageResult struct {
	OK bool
}

// CreateDepartureParams is the adapter input for CreateDeparture.
type CreateDepartureParams struct {
	UserID        string
	BranchID      string
	PackageID     string
	DepartureDate string // ISO YYYY-MM-DD
	ReturnDate    string // ISO YYYY-MM-DD
	TotalSeats    int32
	Status        string // "open" | "closed"
	Pricing       []PricingInputParam
}

// UpdateDepartureParams is the adapter input for UpdateDeparture.
type UpdateDepartureParams struct {
	UserID        string
	BranchID      string
	ID            string
	DepartureDate string // empty = no change
	ReturnDate    string // empty = no change
	TotalSeats    int32  // 0 = no change
	Status        string
	Pricing       []PricingInputParam // nil = no change; non-nil = replace all
}
