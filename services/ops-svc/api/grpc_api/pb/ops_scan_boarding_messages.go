// ops_scan_boarding_messages.go — hand-written proto message types for scan
// events and bus boarding RPCs (BL-OPS-010/011).

package pb

// ---------------------------------------------------------------------------
// RecordScan (BL-OPS-010)
// ---------------------------------------------------------------------------

type RecordScanRequest struct {
	ScanType       string
	DepartureId    string
	JamaahId       string
	ScannedBy      string
	DeviceId       string
	Location       string
	IdempotencyKey string
	Metadata       []byte
}

func (x *RecordScanRequest) GetScanType() string {
	if x == nil {
		return ""
	}
	return x.ScanType
}
func (x *RecordScanRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *RecordScanRequest) GetJamaahId() string {
	if x == nil {
		return ""
	}
	return x.JamaahId
}
func (x *RecordScanRequest) GetScannedBy() string {
	if x == nil {
		return ""
	}
	return x.ScannedBy
}
func (x *RecordScanRequest) GetDeviceId() string {
	if x == nil {
		return ""
	}
	return x.DeviceId
}
func (x *RecordScanRequest) GetLocation() string {
	if x == nil {
		return ""
	}
	return x.Location
}
func (x *RecordScanRequest) GetIdempotencyKey() string {
	if x == nil {
		return ""
	}
	return x.IdempotencyKey
}
func (x *RecordScanRequest) GetMetadata() []byte {
	if x == nil {
		return nil
	}
	return x.Metadata
}

type RecordScanResponse struct {
	ScanId     string
	Idempotent bool
}

func (x *RecordScanResponse) GetScanId() string {
	if x == nil {
		return ""
	}
	return x.ScanId
}
func (x *RecordScanResponse) GetIdempotent() bool {
	if x == nil {
		return false
	}
	return x.Idempotent
}

// ---------------------------------------------------------------------------
// RecordBusBoarding (BL-OPS-011)
// ---------------------------------------------------------------------------

type RecordBusBoardingRequest struct {
	DepartureId string
	BusNumber   string
	JamaahId    string
	ScannedBy   string
	Status      string
}

func (x *RecordBusBoardingRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *RecordBusBoardingRequest) GetBusNumber() string {
	if x == nil {
		return ""
	}
	return x.BusNumber
}
func (x *RecordBusBoardingRequest) GetJamaahId() string {
	if x == nil {
		return ""
	}
	return x.JamaahId
}
func (x *RecordBusBoardingRequest) GetScannedBy() string {
	if x == nil {
		return ""
	}
	return x.ScannedBy
}
func (x *RecordBusBoardingRequest) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

type RecordBusBoardingResponse struct {
	BoardingId string
	Status     string
	Idempotent bool
}

func (x *RecordBusBoardingResponse) GetBoardingId() string {
	if x == nil {
		return ""
	}
	return x.BoardingId
}
func (x *RecordBusBoardingResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *RecordBusBoardingResponse) GetIdempotent() bool {
	if x == nil {
		return false
	}
	return x.Idempotent
}

// ---------------------------------------------------------------------------
// GetBoardingRoster (BL-OPS-011)
// ---------------------------------------------------------------------------

type GetBoardingRosterRequest struct {
	DepartureId string
	BusNumber   string // optional filter
}

func (x *GetBoardingRosterRequest) GetDepartureId() string {
	if x == nil {
		return ""
	}
	return x.DepartureId
}
func (x *GetBoardingRosterRequest) GetBusNumber() string {
	if x == nil {
		return ""
	}
	return x.BusNumber
}

type GetBoardingRosterResponse struct {
	Boardings    []*BoardingEntry
	TotalBoarded int32
	TotalAbsent  int32
}

func (x *GetBoardingRosterResponse) GetBoardings() []*BoardingEntry {
	if x == nil {
		return nil
	}
	return x.Boardings
}
func (x *GetBoardingRosterResponse) GetTotalBoarded() int32 {
	if x == nil {
		return 0
	}
	return x.TotalBoarded
}
func (x *GetBoardingRosterResponse) GetTotalAbsent() int32 {
	if x == nil {
		return 0
	}
	return x.TotalAbsent
}

type BoardingEntry struct {
	JamaahId  string
	BusNumber string
	Status    string
	BoardedAt string // RFC3339 or empty
}
