// booking_channel_messages.go — hand-written proto message types for
// cross-channel seat tracking RPC (BL-BOOK-007).
//
// Run `make genpb` to regenerate from booking.proto once protoc is available,
// then delete this file.
package pb

// ChannelSeatRow carries seat metrics for one booking channel.
type ChannelSeatRow struct {
	Channel      string
	Seats        int32
	BookingCount int32
}

func (x *ChannelSeatRow) GetChannel() string {
	if x == nil {
		return ""
	}
	return x.Channel
}
func (x *ChannelSeatRow) GetSeats() int32 {
	if x == nil {
		return 0
	}
	return x.Seats
}
func (x *ChannelSeatRow) GetBookingCount() int32 {
	if x == nil {
		return 0
	}
	return x.BookingCount
}

// GetSeatsByChannelRequest carries the departure_id to query.
type GetSeatsByChannelRequest struct {
	DepartureId string
}

func (x *GetSeatsByChannelRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}

// GetSeatsByChannelResponse carries the per-channel breakdown.
type GetSeatsByChannelResponse struct {
	DepartureId string
	ByChannel   []*ChannelSeatRow
}

func (x *GetSeatsByChannelResponse) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *GetSeatsByChannelResponse) GetByChannel() []*ChannelSeatRow {
	if x == nil {
		return nil
	}
	return x.ByChannel
}
