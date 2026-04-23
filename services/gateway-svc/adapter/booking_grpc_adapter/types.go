// types.go — adapter-local Go types for BookingService.
//
// Proto types do not leak past this package; the rest of gateway-svc
// sees only these plain Go structs from the service and handler layers.
package booking_grpc_adapter

// ---------------------------------------------------------------------------
// Input types
// ---------------------------------------------------------------------------

// PilgrimInputParam is one jamaah in the CreateDraftBooking request.
type PilgrimInputParam struct {
	FullName    string
	Email       string
	Whatsapp    string
	Domicile    string
	IsLead      bool
	HasKTP      bool
	HasPassport bool
}

// AddonInputParam is a selected add-on in the CreateDraftBooking request.
type AddonInputParam struct {
	AddonID            string
	AddonName          string
	ListAmount         int64
	ListCurrency       string
	SettlementCurrency string
}

// CreateDraftBookingParams is the adapter-side input for CreateDraftBooking.
type CreateDraftBookingParams struct {
	// Channel attribution (BL-BOOK-002).
	Channel     string // "b2c_self" | "b2b_agent" | "cs"
	AgentID     string // required when Channel = "b2b_agent"
	StaffUserID string // required when Channel = "cs"

	// Catalog references.
	PackageID   string
	DepartureID string
	RoomType    string

	// Lead pilgrim contact.
	LeadFullName string
	LeadEmail    string
	LeadWhatsapp string
	LeadDomicile string

	// All pilgrims (including the lead).
	Pilgrims []PilgrimInputParam

	// Optional add-ons.
	Addons []AddonInputParam

	// Pricing snapshot.
	ListAmount         int64
	ListCurrency       string
	SettlementCurrency string

	Notes string

	// Idempotency (optional).
	IdempotencyKey string

	// Mahram validation result (BL-BOOK-006 — non-blocking).
	MahramWarning string
}

// ---------------------------------------------------------------------------
// Output types
// ---------------------------------------------------------------------------

// BookingItemResult is one pilgrim in the CreateDraftBookingResult.
type BookingItemResult struct {
	ID              string
	FullName        string
	IsLead          bool
	DocumentWarning string // advisory (BL-BOOK-005), empty = no warning
}

// CreateDraftBookingResult is the adapter-side output for CreateDraftBooking.
type CreateDraftBookingResult struct {
	ID                 string
	Status             string // "draft"
	Channel            string
	PackageID          string
	DepartureID        string
	RoomType           string
	AgentID            string
	StaffUserID        string
	LeadFullName       string
	LeadEmail          string
	LeadWhatsapp       string
	LeadDomicile       string
	ListAmount         int64
	ListCurrency       string
	SettlementCurrency string
	Notes              string
	IdempotencyKey     string
	CreatedAt          string // RFC 3339
	ExpiresAt          string // RFC 3339
	Items              []BookingItemResult
	MahramWarning      string // non-empty = advisory (BL-BOOK-006)
	Replayed           bool   // true on idempotency dedup hit
}
