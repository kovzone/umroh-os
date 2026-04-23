// Hand-written proto stub for booking-svc's catalog_grpc_adapter.
// Run `make generate` (protoc) to regenerate from catalog.proto once protoc
// is available in the dev toolchain.
//
// Wire names match catalog-svc/api/grpc_api/pb exactly. Only the RPCs that
// booking-svc calls are included (consumer-stub rule, slice-S1 § IAM).

package pb

// ReserveSeatsRequest — matches CatalogService.ReserveSeats wire shape.
type ReserveSeatsRequest struct {
	ReservationId       string `protobuf:"bytes,1,opt,name=reservation_id,json=reservationId,proto3" json:"reservation_id,omitempty"`
	DepartureId         string `protobuf:"bytes,2,opt,name=departure_id,json=departureId,proto3" json:"departure_id,omitempty"`
	Seats               int32  `protobuf:"varint,3,opt,name=seats,proto3" json:"seats,omitempty"`
	IdempotencyTtlHours int32  `protobuf:"varint,4,opt,name=idempotency_ttl_hours,json=idempotencyTtlHours,proto3" json:"idempotency_ttl_hours,omitempty"`
}

func (r *ReserveSeatsRequest) Reset()         {}
func (r *ReserveSeatsRequest) String() string { return r.ReservationId }
func (r *ReserveSeatsRequest) ProtoMessage()  {}

func (r *ReserveSeatsRequest) GetReservationId() string       { return r.ReservationId }
func (r *ReserveSeatsRequest) GetDepartureId() string         { return r.DepartureId }
func (r *ReserveSeatsRequest) GetSeats() int32                { return r.Seats }
func (r *ReserveSeatsRequest) GetIdempotencyTtlHours() int32  { return r.IdempotencyTtlHours }

// Reservation sub-message.
type Reservation struct {
	ReservationId string `protobuf:"bytes,1,opt,name=reservation_id,json=reservationId,proto3" json:"reservation_id,omitempty"`
	DepartureId   string `protobuf:"bytes,2,opt,name=departure_id,json=departureId,proto3" json:"departure_id,omitempty"`
	Seats         int32  `protobuf:"varint,3,opt,name=seats,proto3" json:"seats,omitempty"`
	ReservedAt    string `protobuf:"bytes,4,opt,name=reserved_at,json=reservedAt,proto3" json:"reserved_at,omitempty"`
	ExpiresAt     string `protobuf:"bytes,5,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
}

func (r *Reservation) Reset()         {}
func (r *Reservation) String() string { return r.ReservationId }
func (r *Reservation) ProtoMessage()  {}

func (r *Reservation) GetReservationId() string { return r.ReservationId }
func (r *Reservation) GetDepartureId() string   { return r.DepartureId }
func (r *Reservation) GetSeats() int32          { return r.Seats }
func (r *Reservation) GetReservedAt() string    { return r.ReservedAt }
func (r *Reservation) GetExpiresAt() string     { return r.ExpiresAt }

// ReserveSeatsResponse.
type ReserveSeatsResponse struct {
	Reservation    *Reservation `protobuf:"bytes,1,opt,name=reservation,proto3" json:"reservation,omitempty"`
	RemainingSeats int32        `protobuf:"varint,2,opt,name=remaining_seats,json=remainingSeats,proto3" json:"remaining_seats,omitempty"`
	Replayed       bool         `protobuf:"varint,3,opt,name=replayed,proto3" json:"replayed,omitempty"`
}

func (r *ReserveSeatsResponse) Reset()         {}
func (r *ReserveSeatsResponse) String() string { return "" }
func (r *ReserveSeatsResponse) ProtoMessage()  {}

func (r *ReserveSeatsResponse) GetReservation() *Reservation    { return r.Reservation }
func (r *ReserveSeatsResponse) GetRemainingSeats() int32        { return r.RemainingSeats }
func (r *ReserveSeatsResponse) GetReplayed() bool               { return r.Replayed }

// ReleaseSeatsRequest.
type ReleaseSeatsRequest struct {
	ReservationId string `protobuf:"bytes,1,opt,name=reservation_id,json=reservationId,proto3" json:"reservation_id,omitempty"`
	Seats         int32  `protobuf:"varint,2,opt,name=seats,proto3" json:"seats,omitempty"`
	Reason        string `protobuf:"bytes,3,opt,name=reason,proto3" json:"reason,omitempty"`
}

func (r *ReleaseSeatsRequest) Reset()         {}
func (r *ReleaseSeatsRequest) String() string { return r.ReservationId }
func (r *ReleaseSeatsRequest) ProtoMessage()  {}

func (r *ReleaseSeatsRequest) GetReservationId() string { return r.ReservationId }
func (r *ReleaseSeatsRequest) GetSeats() int32          { return r.Seats }
func (r *ReleaseSeatsRequest) GetReason() string        { return r.Reason }

// Released sub-message.
type Released struct {
	ReservationId string `protobuf:"bytes,1,opt,name=reservation_id,json=reservationId,proto3" json:"reservation_id,omitempty"`
	DepartureId   string `protobuf:"bytes,2,opt,name=departure_id,json=departureId,proto3" json:"departure_id,omitempty"`
	SeatsReleased int32  `protobuf:"varint,3,opt,name=seats_released,json=seatsReleased,proto3" json:"seats_released,omitempty"`
	ReleasedAt    string `protobuf:"bytes,4,opt,name=released_at,json=releasedAt,proto3" json:"released_at,omitempty"`
}

func (r *Released) Reset()         {}
func (r *Released) String() string { return r.ReservationId }
func (r *Released) ProtoMessage()  {}

func (r *Released) GetReservationId() string { return r.ReservationId }
func (r *Released) GetDepartureId() string   { return r.DepartureId }
func (r *Released) GetSeatsReleased() int32  { return r.SeatsReleased }
func (r *Released) GetReleasedAt() string    { return r.ReleasedAt }

// ReleaseSeatsResponse.
type ReleaseSeatsResponse struct {
	Released       *Released `protobuf:"bytes,1,opt,name=released,proto3" json:"released,omitempty"`
	RemainingSeats int32     `protobuf:"varint,2,opt,name=remaining_seats,json=remainingSeats,proto3" json:"remaining_seats,omitempty"`
	Replayed       bool      `protobuf:"varint,3,opt,name=replayed,proto3" json:"replayed,omitempty"`
}

func (r *ReleaseSeatsResponse) Reset()         {}
func (r *ReleaseSeatsResponse) String() string { return "" }
func (r *ReleaseSeatsResponse) ProtoMessage()  {}

func (r *ReleaseSeatsResponse) GetReleased() *Released      { return r.Released }
func (r *ReleaseSeatsResponse) GetRemainingSeats() int32    { return r.RemainingSeats }
func (r *ReleaseSeatsResponse) GetReplayed() bool           { return r.Replayed }

// GetPackageDepartureRequest.
type GetPackageDepartureRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (r *GetPackageDepartureRequest) Reset()         {}
func (r *GetPackageDepartureRequest) String() string { return r.Id }
func (r *GetPackageDepartureRequest) ProtoMessage()  {}
func (r *GetPackageDepartureRequest) GetId() string  { return r.Id }

// DepartureDetail (narrow booking-svc view).
type DepartureDetail struct {
	Id             string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	PackageId      string `protobuf:"bytes,2,opt,name=package_id,json=packageId,proto3" json:"package_id,omitempty"`
	DepartureDate  string `protobuf:"bytes,3,opt,name=departure_date,json=departureDate,proto3" json:"departure_date,omitempty"`
	ReturnDate     string `protobuf:"bytes,4,opt,name=return_date,json=returnDate,proto3" json:"return_date,omitempty"`
	TotalSeats     int32  `protobuf:"varint,5,opt,name=total_seats,json=totalSeats,proto3" json:"total_seats,omitempty"`
	RemainingSeats int32  `protobuf:"varint,6,opt,name=remaining_seats,json=remainingSeats,proto3" json:"remaining_seats,omitempty"`
	Status         string `protobuf:"bytes,7,opt,name=status,proto3" json:"status,omitempty"`
}

func (d *DepartureDetail) Reset()               {}
func (d *DepartureDetail) String() string        { return d.Id }
func (d *DepartureDetail) ProtoMessage()         {}
func (d *DepartureDetail) GetId() string         { return d.Id }
func (d *DepartureDetail) GetPackageId() string  { return d.PackageId }
func (d *DepartureDetail) GetDepartureDate() string { return d.DepartureDate }
func (d *DepartureDetail) GetReturnDate() string { return d.ReturnDate }
func (d *DepartureDetail) GetTotalSeats() int32  { return d.TotalSeats }
func (d *DepartureDetail) GetRemainingSeats() int32 { return d.RemainingSeats }
func (d *DepartureDetail) GetStatus() string     { return d.Status }

// GetPackageDepartureResponse.
type GetPackageDepartureResponse struct {
	Departure *DepartureDetail `protobuf:"bytes,1,opt,name=departure,proto3" json:"departure,omitempty"`
}

func (r *GetPackageDepartureResponse) Reset()                        {}
func (r *GetPackageDepartureResponse) String() string                { return "" }
func (r *GetPackageDepartureResponse) ProtoMessage()                 {}
func (r *GetPackageDepartureResponse) GetDeparture() *DepartureDetail { return r.Departure }
