// submit_booking_messages.go — hand-written proto message types for
// SubmitBooking RPC (S2 / BL-BOOK-005).
//
// Mirrors what protoc-gen-go would generate from:
//   message SubmitBookingRequest  { string booking_id = 1; }
//   message SubmitBookingResponse { BookingSummary booking = 1; }
//   message BookingSummary        { id, status, channel, ... }

package pb

// SubmitBookingRequest is sent by gateway-svc to transition a booking
// from 'draft' to 'pending_payment'.
type SubmitBookingRequest struct {
	BookingId string
}

func (x *SubmitBookingRequest) GetBookingId() string {
	if x == nil {
		return ""
	}
	return x.BookingId
}

// SubmitBookingResponse carries the updated booking summary.
type SubmitBookingResponse struct {
	Booking *BookingSummary
}

func (x *SubmitBookingResponse) GetBooking() *BookingSummary {
	if x == nil {
		return nil
	}
	return x.Booking
}

// BookingSummary is a lightweight booking projection used in SubmitBooking
// and other status-change responses.
type BookingSummary struct {
	Id                 string
	Status             string
	Channel            string
	PackageId          string
	DepartureId        string
	RoomType           string
	LeadFullName       string
	LeadWhatsapp       string
	LeadDomicile       string
	ListAmount         int64
	ListCurrency       string
	SettlementCurrency string
}

func (x *BookingSummary) GetId() string {
	if x == nil {
		return ""
	}
	return x.Id
}
func (x *BookingSummary) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *BookingSummary) GetChannel() string {
	if x == nil {
		return ""
	}
	return x.Channel
}
func (x *BookingSummary) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}
func (x *BookingSummary) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *BookingSummary) GetRoomType() string {
	if x == nil {
		return ""
	}
	return x.RoomType
}
func (x *BookingSummary) GetLeadFullName() string {
	if x == nil {
		return ""
	}
	return x.LeadFullName
}
func (x *BookingSummary) GetLeadWhatsapp() string {
	if x == nil {
		return ""
	}
	return x.LeadWhatsapp
}
func (x *BookingSummary) GetLeadDomicile() string {
	if x == nil {
		return ""
	}
	return x.LeadDomicile
}
func (x *BookingSummary) GetListAmount() int64 {
	if x == nil {
		return 0
	}
	return x.ListAmount
}
func (x *BookingSummary) GetListCurrency() string {
	if x == nil {
		return ""
	}
	return x.ListCurrency
}
func (x *BookingSummary) GetSettlementCurrency() string {
	if x == nil {
		return ""
	}
	return x.SettlementCurrency
}
