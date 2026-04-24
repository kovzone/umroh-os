// ops_depth_messages.go — hand-written proto message types for ops-svc Wave 5 depth RPCs.
// BL-OPS-021..042: collective documents, passport OCR, mahram, document progress,
// official letters, immigration manifest, smart rooming, transport, manifest delta,
// staff assignment, passport log, visa progress, e-visa, external provider,
// refunds, luggage, broadcast, tasreh, audio devices, zamzam, room check-in.

package pb

// ---------------------------------------------------------------------------
// BL-OPS-021: Collective document storage
// ---------------------------------------------------------------------------

type StoreCollectiveDocumentRequest struct {
	DepartureID  string
	DocumentType string
	URL          string
	PilgrimID    string
	Notes        string
}

func (x *StoreCollectiveDocumentRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *StoreCollectiveDocumentRequest) GetDocumentType() string {
	if x == nil {
		return ""
	}
	return x.DocumentType
}
func (x *StoreCollectiveDocumentRequest) GetURL() string {
	if x == nil {
		return ""
	}
	return x.URL
}
func (x *StoreCollectiveDocumentRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *StoreCollectiveDocumentRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type StoreCollectiveDocumentResponse struct {
	DocumentID string
}

func (x *StoreCollectiveDocumentResponse) GetDocumentID() string {
	if x == nil {
		return ""
	}
	return x.DocumentID
}

type GetCollectiveDocumentsRequest struct {
	DepartureID  string
	PilgrimID    string
	DocumentType string
}

func (x *GetCollectiveDocumentsRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *GetCollectiveDocumentsRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *GetCollectiveDocumentsRequest) GetDocumentType() string {
	if x == nil {
		return ""
	}
	return x.DocumentType
}

type CollectiveDocRow struct {
	DocumentID   string
	DocumentType string
	URL          string
	PilgrimID    string
	UploadedAt   string
}

func (x *CollectiveDocRow) GetDocumentID() string {
	if x == nil {
		return ""
	}
	return x.DocumentID
}
func (x *CollectiveDocRow) GetDocumentType() string {
	if x == nil {
		return ""
	}
	return x.DocumentType
}
func (x *CollectiveDocRow) GetURL() string {
	if x == nil {
		return ""
	}
	return x.URL
}
func (x *CollectiveDocRow) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *CollectiveDocRow) GetUploadedAt() string {
	if x == nil {
		return ""
	}
	return x.UploadedAt
}

type GetCollectiveDocumentsResponse struct {
	DepartureID string
	Documents   []*CollectiveDocRow
}

func (x *GetCollectiveDocumentsResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *GetCollectiveDocumentsResponse) GetDocuments() []*CollectiveDocRow {
	if x == nil {
		return nil
	}
	return x.Documents
}

type SetDocumentACLRequest struct {
	DocumentID     string
	AccessLevel    string
	AllowedUserIDs []string
}

func (x *SetDocumentACLRequest) GetDocumentID() string {
	if x == nil {
		return ""
	}
	return x.DocumentID
}
func (x *SetDocumentACLRequest) GetAccessLevel() string {
	if x == nil {
		return ""
	}
	return x.AccessLevel
}
func (x *SetDocumentACLRequest) GetAllowedUserIDs() []string {
	if x == nil {
		return nil
	}
	return x.AllowedUserIDs
}

type SetDocumentACLResponse struct {
	Updated bool
}

func (x *SetDocumentACLResponse) GetUpdated() bool {
	if x == nil {
		return false
	}
	return x.Updated
}

// ---------------------------------------------------------------------------
// BL-OPS-022: Passport OCR & mahram
// ---------------------------------------------------------------------------

type ExtractPassportOCRRequest struct {
	ImageURL  string
	PilgrimID string
}

func (x *ExtractPassportOCRRequest) GetImageURL() string {
	if x == nil {
		return ""
	}
	return x.ImageURL
}
func (x *ExtractPassportOCRRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}

type PassportOCRData struct {
	FullName       string
	PassportNumber string
	Nationality    string
	DateOfBirth    string
	ExpiryDate     string
	Gender         string
	Confidence     float64
}

func (x *PassportOCRData) GetFullName() string {
	if x == nil {
		return ""
	}
	return x.FullName
}
func (x *PassportOCRData) GetPassportNumber() string {
	if x == nil {
		return ""
	}
	return x.PassportNumber
}
func (x *PassportOCRData) GetNationality() string {
	if x == nil {
		return ""
	}
	return x.Nationality
}
func (x *PassportOCRData) GetDateOfBirth() string {
	if x == nil {
		return ""
	}
	return x.DateOfBirth
}
func (x *PassportOCRData) GetExpiryDate() string {
	if x == nil {
		return ""
	}
	return x.ExpiryDate
}
func (x *PassportOCRData) GetGender() string {
	if x == nil {
		return ""
	}
	return x.Gender
}
func (x *PassportOCRData) GetConfidence() float64 {
	if x == nil {
		return 0
	}
	return x.Confidence
}

type ExtractPassportOCRResponse struct {
	PilgrimID string
	Data      *PassportOCRData
	Warnings  []string
}

func (x *ExtractPassportOCRResponse) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *ExtractPassportOCRResponse) GetData() *PassportOCRData {
	if x == nil {
		return nil
	}
	return x.Data
}
func (x *ExtractPassportOCRResponse) GetWarnings() []string {
	if x == nil {
		return nil
	}
	return x.Warnings
}

type SetMahramRelationRequest struct {
	PilgrimID       string
	MahramPilgrimID string
	Relation        string // "suami"/"ayah"/"saudara"/"lainnya"
}

func (x *SetMahramRelationRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *SetMahramRelationRequest) GetMahramPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.MahramPilgrimID
}
func (x *SetMahramRelationRequest) GetRelation() string {
	if x == nil {
		return ""
	}
	return x.Relation
}

type SetMahramRelationResponse struct {
	RelationID string
}

func (x *SetMahramRelationResponse) GetRelationID() string {
	if x == nil {
		return ""
	}
	return x.RelationID
}

type GetMahramRelationsRequest struct {
	BookingID string
}

func (x *GetMahramRelationsRequest) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}

type MahramRelationRow struct {
	RelationID      string
	PilgrimID       string
	MahramPilgrimID string
	Relation        string
}

func (x *MahramRelationRow) GetRelationID() string {
	if x == nil {
		return ""
	}
	return x.RelationID
}
func (x *MahramRelationRow) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *MahramRelationRow) GetMahramPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.MahramPilgrimID
}
func (x *MahramRelationRow) GetRelation() string {
	if x == nil {
		return ""
	}
	return x.Relation
}

type GetMahramRelationsResponse struct {
	BookingID string
	Relations []*MahramRelationRow
}

func (x *GetMahramRelationsResponse) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *GetMahramRelationsResponse) GetRelations() []*MahramRelationRow {
	if x == nil {
		return nil
	}
	return x.Relations
}

// ---------------------------------------------------------------------------
// BL-OPS-023: Document progress & expiry
// ---------------------------------------------------------------------------

type GetDocumentProgressRequest struct {
	DepartureID string
}

func (x *GetDocumentProgressRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type DocProgressRow struct {
	PilgrimID       string
	DocumentType    string
	Status          string
	ExpiryDate      string
	DaysUntilExpiry int32
}

func (x *DocProgressRow) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *DocProgressRow) GetDocumentType() string {
	if x == nil {
		return ""
	}
	return x.DocumentType
}
func (x *DocProgressRow) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *DocProgressRow) GetExpiryDate() string {
	if x == nil {
		return ""
	}
	return x.ExpiryDate
}
func (x *DocProgressRow) GetDaysUntilExpiry() int32 {
	if x == nil {
		return 0
	}
	return x.DaysUntilExpiry
}

type GetDocumentProgressResponse struct {
	DepartureID       string
	Rows              []*DocProgressRow
	TotalPilgrims     int32
	DocumentsComplete int32
	DocumentsExpiring int32
}

func (x *GetDocumentProgressResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *GetDocumentProgressResponse) GetRows() []*DocProgressRow {
	if x == nil {
		return nil
	}
	return x.Rows
}
func (x *GetDocumentProgressResponse) GetTotalPilgrims() int32 {
	if x == nil {
		return 0
	}
	return x.TotalPilgrims
}
func (x *GetDocumentProgressResponse) GetDocumentsComplete() int32 {
	if x == nil {
		return 0
	}
	return x.DocumentsComplete
}
func (x *GetDocumentProgressResponse) GetDocumentsExpiring() int32 {
	if x == nil {
		return 0
	}
	return x.DocumentsExpiring
}

type GetExpiryAlertsRequest struct {
	ThresholdDays int32
	DepartureID   string
}

func (x *GetExpiryAlertsRequest) GetThresholdDays() int32 {
	if x == nil {
		return 0
	}
	return x.ThresholdDays
}
func (x *GetExpiryAlertsRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type ExpiryAlertRow struct {
	PilgrimID       string
	DocumentType    string
	ExpiryDate      string
	DaysUntilExpiry int32
}

func (x *ExpiryAlertRow) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *ExpiryAlertRow) GetDocumentType() string {
	if x == nil {
		return ""
	}
	return x.DocumentType
}
func (x *ExpiryAlertRow) GetExpiryDate() string {
	if x == nil {
		return ""
	}
	return x.ExpiryDate
}
func (x *ExpiryAlertRow) GetDaysUntilExpiry() int32 {
	if x == nil {
		return 0
	}
	return x.DaysUntilExpiry
}

type GetExpiryAlertsResponse struct {
	Alerts []*ExpiryAlertRow
}

func (x *GetExpiryAlertsResponse) GetAlerts() []*ExpiryAlertRow {
	if x == nil {
		return nil
	}
	return x.Alerts
}

// ---------------------------------------------------------------------------
// BL-OPS-024: Official letter
// ---------------------------------------------------------------------------

type GenerateOfficialLetterRequest struct {
	TemplateName string
	DepartureID  string
	PilgrimID    string
	IssuedTo     string
	Notes        string
}

func (x *GenerateOfficialLetterRequest) GetTemplateName() string {
	if x == nil {
		return ""
	}
	return x.TemplateName
}
func (x *GenerateOfficialLetterRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *GenerateOfficialLetterRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *GenerateOfficialLetterRequest) GetIssuedTo() string {
	if x == nil {
		return ""
	}
	return x.IssuedTo
}
func (x *GenerateOfficialLetterRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type GenerateOfficialLetterResponse struct {
	LetterID  string
	LetterURL string
	IssuedAt  string
}

func (x *GenerateOfficialLetterResponse) GetLetterID() string {
	if x == nil {
		return ""
	}
	return x.LetterID
}
func (x *GenerateOfficialLetterResponse) GetLetterURL() string {
	if x == nil {
		return ""
	}
	return x.LetterURL
}
func (x *GenerateOfficialLetterResponse) GetIssuedAt() string {
	if x == nil {
		return ""
	}
	return x.IssuedAt
}

// ---------------------------------------------------------------------------
// BL-OPS-025: Immigration manifest
// ---------------------------------------------------------------------------

type GenerateImmigrationManifestRequest struct {
	DepartureID string
	Format      string // "pdf"/"csv"
}

func (x *GenerateImmigrationManifestRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *GenerateImmigrationManifestRequest) GetFormat() string {
	if x == nil {
		return ""
	}
	return x.Format
}

type GenerateImmigrationManifestResponse struct {
	ManifestID  string
	ManifestURL string
	Version     string
	RowCount    int32
}

func (x *GenerateImmigrationManifestResponse) GetManifestID() string {
	if x == nil {
		return ""
	}
	return x.ManifestID
}
func (x *GenerateImmigrationManifestResponse) GetManifestURL() string {
	if x == nil {
		return ""
	}
	return x.ManifestURL
}
func (x *GenerateImmigrationManifestResponse) GetVersion() string {
	if x == nil {
		return ""
	}
	return x.Version
}
func (x *GenerateImmigrationManifestResponse) GetRowCount() int32 {
	if x == nil {
		return 0
	}
	return x.RowCount
}

// ---------------------------------------------------------------------------
// BL-OPS-026: Smart rooming
// ---------------------------------------------------------------------------

type SmartRoomingOptions struct {
	MaxPerRoom      int32
	SeparateMahram  bool
	PreferSameFamily bool
}

func (x *SmartRoomingOptions) GetMaxPerRoom() int32 {
	if x == nil {
		return 0
	}
	return x.MaxPerRoom
}
func (x *SmartRoomingOptions) GetSeparateMahram() bool {
	if x == nil {
		return false
	}
	return x.SeparateMahram
}
func (x *SmartRoomingOptions) GetPreferSameFamily() bool {
	if x == nil {
		return false
	}
	return x.PreferSameFamily
}

type RunSmartRoomingRequest struct {
	DepartureID string
	Options     *SmartRoomingOptions
}

func (x *RunSmartRoomingRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *RunSmartRoomingRequest) GetOptions() *SmartRoomingOptions {
	if x == nil {
		return nil
	}
	return x.Options
}

type RunSmartRoomingResponse struct {
	DepartureID   string
	RoomsAssigned int32
	TotalPilgrims int32
	Warnings      []string
}

func (x *RunSmartRoomingResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *RunSmartRoomingResponse) GetRoomsAssigned() int32 {
	if x == nil {
		return 0
	}
	return x.RoomsAssigned
}
func (x *RunSmartRoomingResponse) GetTotalPilgrims() int32 {
	if x == nil {
		return 0
	}
	return x.TotalPilgrims
}
func (x *RunSmartRoomingResponse) GetWarnings() []string {
	if x == nil {
		return nil
	}
	return x.Warnings
}

// ---------------------------------------------------------------------------
// BL-OPS-027: Transport arrangement
// ---------------------------------------------------------------------------

type AssignTransportRequest struct {
	DepartureID string
	VehicleType string
	VehicleID   string
	PilgrimIDs  []string
}

func (x *AssignTransportRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *AssignTransportRequest) GetVehicleType() string {
	if x == nil {
		return ""
	}
	return x.VehicleType
}
func (x *AssignTransportRequest) GetVehicleID() string {
	if x == nil {
		return ""
	}
	return x.VehicleID
}
func (x *AssignTransportRequest) GetPilgrimIDs() []string {
	if x == nil {
		return nil
	}
	return x.PilgrimIDs
}

type AssignTransportResponse struct {
	AssignmentID  string
	AssignedCount int32
}

func (x *AssignTransportResponse) GetAssignmentID() string {
	if x == nil {
		return ""
	}
	return x.AssignmentID
}
func (x *AssignTransportResponse) GetAssignedCount() int32 {
	if x == nil {
		return 0
	}
	return x.AssignedCount
}

type GetTransportAssignmentsRequest struct {
	DepartureID string
}

func (x *GetTransportAssignmentsRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type TransportAssignmentRow struct {
	AssignmentID string
	VehicleType  string
	VehicleID    string
	PilgrimIDs   []string
}

func (x *TransportAssignmentRow) GetAssignmentID() string {
	if x == nil {
		return ""
	}
	return x.AssignmentID
}
func (x *TransportAssignmentRow) GetVehicleType() string {
	if x == nil {
		return ""
	}
	return x.VehicleType
}
func (x *TransportAssignmentRow) GetVehicleID() string {
	if x == nil {
		return ""
	}
	return x.VehicleID
}
func (x *TransportAssignmentRow) GetPilgrimIDs() []string {
	if x == nil {
		return nil
	}
	return x.PilgrimIDs
}

type GetTransportAssignmentsResponse struct {
	DepartureID string
	Assignments []*TransportAssignmentRow
}

func (x *GetTransportAssignmentsResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *GetTransportAssignmentsResponse) GetAssignments() []*TransportAssignmentRow {
	if x == nil {
		return nil
	}
	return x.Assignments
}

// ---------------------------------------------------------------------------
// BL-OPS-028: Manifest delta
// ---------------------------------------------------------------------------

type PublishManifestDeltaRequest struct {
	DepartureID string
	ChangeType  string
	EntityID    string
	Notes       string
}

func (x *PublishManifestDeltaRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *PublishManifestDeltaRequest) GetChangeType() string {
	if x == nil {
		return ""
	}
	return x.ChangeType
}
func (x *PublishManifestDeltaRequest) GetEntityID() string {
	if x == nil {
		return ""
	}
	return x.EntityID
}
func (x *PublishManifestDeltaRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type PublishManifestDeltaResponse struct {
	DeltaID     string
	PublishedAt string
}

func (x *PublishManifestDeltaResponse) GetDeltaID() string {
	if x == nil {
		return ""
	}
	return x.DeltaID
}
func (x *PublishManifestDeltaResponse) GetPublishedAt() string {
	if x == nil {
		return ""
	}
	return x.PublishedAt
}

// ---------------------------------------------------------------------------
// BL-OPS-029: Staff assignment
// ---------------------------------------------------------------------------

type AssignStaffRequest struct {
	DepartureID string
	StaffUserID string
	Role        string
	PilgrimIDs  []string
}

func (x *AssignStaffRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *AssignStaffRequest) GetStaffUserID() string {
	if x == nil {
		return ""
	}
	return x.StaffUserID
}
func (x *AssignStaffRequest) GetRole() string {
	if x == nil {
		return ""
	}
	return x.Role
}
func (x *AssignStaffRequest) GetPilgrimIDs() []string {
	if x == nil {
		return nil
	}
	return x.PilgrimIDs
}

type AssignStaffResponse struct {
	AssignmentID string
}

func (x *AssignStaffResponse) GetAssignmentID() string {
	if x == nil {
		return ""
	}
	return x.AssignmentID
}

// ---------------------------------------------------------------------------
// BL-OPS-030: Passport log
// ---------------------------------------------------------------------------

type RecordPassportHandoverRequest struct {
	DepartureID string
	PilgrimID   string
	FromUserID  string
	ToUserID    string
	Notes       string
}

func (x *RecordPassportHandoverRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *RecordPassportHandoverRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *RecordPassportHandoverRequest) GetFromUserID() string {
	if x == nil {
		return ""
	}
	return x.FromUserID
}
func (x *RecordPassportHandoverRequest) GetToUserID() string {
	if x == nil {
		return ""
	}
	return x.ToUserID
}
func (x *RecordPassportHandoverRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type RecordPassportHandoverResponse struct {
	HandoverID string
	RecordedAt string
}

func (x *RecordPassportHandoverResponse) GetHandoverID() string {
	if x == nil {
		return ""
	}
	return x.HandoverID
}
func (x *RecordPassportHandoverResponse) GetRecordedAt() string {
	if x == nil {
		return ""
	}
	return x.RecordedAt
}

type GetPassportLogRequest struct {
	DepartureID string
}

func (x *GetPassportLogRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type PassportHandoverRow struct {
	HandoverID string
	PilgrimID  string
	FromUserID string
	ToUserID   string
	RecordedAt string
	Notes      string
}

func (x *PassportHandoverRow) GetHandoverID() string {
	if x == nil {
		return ""
	}
	return x.HandoverID
}
func (x *PassportHandoverRow) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *PassportHandoverRow) GetFromUserID() string {
	if x == nil {
		return ""
	}
	return x.FromUserID
}
func (x *PassportHandoverRow) GetToUserID() string {
	if x == nil {
		return ""
	}
	return x.ToUserID
}
func (x *PassportHandoverRow) GetRecordedAt() string {
	if x == nil {
		return ""
	}
	return x.RecordedAt
}
func (x *PassportHandoverRow) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type GetPassportLogResponse struct {
	DepartureID string
	Rows        []*PassportHandoverRow
}

func (x *GetPassportLogResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *GetPassportLogResponse) GetRows() []*PassportHandoverRow {
	if x == nil {
		return nil
	}
	return x.Rows
}

// ---------------------------------------------------------------------------
// BL-OPS-031: Visa progress
// ---------------------------------------------------------------------------

type GetVisaProgressRequest struct {
	DepartureID string
}

func (x *GetVisaProgressRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type VisaProgressRow struct {
	PilgrimID     string
	VisaStatus    string
	SubmittedAt   string
	ExpectedBy    string
	DaysRemaining int32
}

func (x *VisaProgressRow) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *VisaProgressRow) GetVisaStatus() string {
	if x == nil {
		return ""
	}
	return x.VisaStatus
}
func (x *VisaProgressRow) GetSubmittedAt() string {
	if x == nil {
		return ""
	}
	return x.SubmittedAt
}
func (x *VisaProgressRow) GetExpectedBy() string {
	if x == nil {
		return ""
	}
	return x.ExpectedBy
}
func (x *VisaProgressRow) GetDaysRemaining() int32 {
	if x == nil {
		return 0
	}
	return x.DaysRemaining
}

type GetVisaProgressResponse struct {
	DepartureID string
	Rows        []*VisaProgressRow
	Submitted   int32
	Approved    int32
	Rejected    int32
	Pending     int32
}

func (x *GetVisaProgressResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *GetVisaProgressResponse) GetRows() []*VisaProgressRow {
	if x == nil {
		return nil
	}
	return x.Rows
}
func (x *GetVisaProgressResponse) GetSubmitted() int32 {
	if x == nil {
		return 0
	}
	return x.Submitted
}
func (x *GetVisaProgressResponse) GetApproved() int32 {
	if x == nil {
		return 0
	}
	return x.Approved
}
func (x *GetVisaProgressResponse) GetRejected() int32 {
	if x == nil {
		return 0
	}
	return x.Rejected
}
func (x *GetVisaProgressResponse) GetPending() int32 {
	if x == nil {
		return 0
	}
	return x.Pending
}

// ---------------------------------------------------------------------------
// BL-OPS-032: E-visa repository
// ---------------------------------------------------------------------------

type StoreEVisaRequest struct {
	PilgrimID   string
	DepartureID string
	VisaNumber  string
	VisaURL     string
	IssuedDate  string
	ExpiryDate  string
}

func (x *StoreEVisaRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *StoreEVisaRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *StoreEVisaRequest) GetVisaNumber() string {
	if x == nil {
		return ""
	}
	return x.VisaNumber
}
func (x *StoreEVisaRequest) GetVisaURL() string {
	if x == nil {
		return ""
	}
	return x.VisaURL
}
func (x *StoreEVisaRequest) GetIssuedDate() string {
	if x == nil {
		return ""
	}
	return x.IssuedDate
}
func (x *StoreEVisaRequest) GetExpiryDate() string {
	if x == nil {
		return ""
	}
	return x.ExpiryDate
}

type StoreEVisaResponse struct {
	EVisaID string
}

func (x *StoreEVisaResponse) GetEVisaID() string {
	if x == nil {
		return ""
	}
	return x.EVisaID
}

type GetEVisaRequest struct {
	PilgrimID   string
	DepartureID string
}

func (x *GetEVisaRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *GetEVisaRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type EVisa struct {
	EVisaID     string
	PilgrimID   string
	VisaNumber  string
	VisaURL     string
	IssuedDate  string
	ExpiryDate  string
}

func (x *EVisa) GetEVisaID() string {
	if x == nil {
		return ""
	}
	return x.EVisaID
}
func (x *EVisa) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *EVisa) GetVisaNumber() string {
	if x == nil {
		return ""
	}
	return x.VisaNumber
}
func (x *EVisa) GetVisaURL() string {
	if x == nil {
		return ""
	}
	return x.VisaURL
}
func (x *EVisa) GetIssuedDate() string {
	if x == nil {
		return ""
	}
	return x.IssuedDate
}
func (x *EVisa) GetExpiryDate() string {
	if x == nil {
		return ""
	}
	return x.ExpiryDate
}

type GetEVisaResponse struct {
	EVisa *EVisa
}

func (x *GetEVisaResponse) GetEVisa() *EVisa {
	if x == nil {
		return nil
	}
	return x.EVisa
}

// ---------------------------------------------------------------------------
// BL-OPS-033: External provider (stub)
// ---------------------------------------------------------------------------

type TriggerExternalProviderRequest struct {
	Provider    string
	Action      string
	ReferenceID string
	Payload     string
}

func (x *TriggerExternalProviderRequest) GetProvider() string {
	if x == nil {
		return ""
	}
	return x.Provider
}
func (x *TriggerExternalProviderRequest) GetAction() string {
	if x == nil {
		return ""
	}
	return x.Action
}
func (x *TriggerExternalProviderRequest) GetReferenceID() string {
	if x == nil {
		return ""
	}
	return x.ReferenceID
}
func (x *TriggerExternalProviderRequest) GetPayload() string {
	if x == nil {
		return ""
	}
	return x.Payload
}

type TriggerExternalProviderResponse struct {
	RequestID string
	Status    string
}

func (x *TriggerExternalProviderResponse) GetRequestID() string {
	if x == nil {
		return ""
	}
	return x.RequestID
}
func (x *TriggerExternalProviderResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

// ---------------------------------------------------------------------------
// BL-OPS-034: Refund & penalty
// ---------------------------------------------------------------------------

type CreateRefundRequest struct {
	BookingID string
	Reason    string
	Amount    int64
	Notes     string
}

func (x *CreateRefundRequest) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *CreateRefundRequest) GetReason() string {
	if x == nil {
		return ""
	}
	return x.Reason
}
func (x *CreateRefundRequest) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *CreateRefundRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type CreateRefundResponse struct {
	RefundID string
	Status   string
}

func (x *CreateRefundResponse) GetRefundID() string {
	if x == nil {
		return ""
	}
	return x.RefundID
}
func (x *CreateRefundResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

type ApproveRefundRequest struct {
	RefundID string
	Decision string
	Notes    string
}

func (x *ApproveRefundRequest) GetRefundID() string {
	if x == nil {
		return ""
	}
	return x.RefundID
}
func (x *ApproveRefundRequest) GetDecision() string {
	if x == nil {
		return ""
	}
	return x.Decision
}
func (x *ApproveRefundRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type ApproveRefundResponse struct {
	RefundID string
	Status   string
}

func (x *ApproveRefundResponse) GetRefundID() string {
	if x == nil {
		return ""
	}
	return x.RefundID
}
func (x *ApproveRefundResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

type RecordPenaltyRequest struct {
	BookingID   string
	PenaltyType string
	Amount      int64
	Notes       string
}

func (x *RecordPenaltyRequest) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *RecordPenaltyRequest) GetPenaltyType() string {
	if x == nil {
		return ""
	}
	return x.PenaltyType
}
func (x *RecordPenaltyRequest) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *RecordPenaltyRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type RecordPenaltyResponse struct {
	PenaltyID string
}

func (x *RecordPenaltyResponse) GetPenaltyID() string {
	if x == nil {
		return ""
	}
	return x.PenaltyID
}

// ---------------------------------------------------------------------------
// BL-OPS-036: Luggage counter
// ---------------------------------------------------------------------------

type RecordLuggageScanRequest struct {
	DepartureID string
	PilgrimID   string
	TagID       string
	ScanPoint   string
}

func (x *RecordLuggageScanRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *RecordLuggageScanRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *RecordLuggageScanRequest) GetTagID() string {
	if x == nil {
		return ""
	}
	return x.TagID
}
func (x *RecordLuggageScanRequest) GetScanPoint() string {
	if x == nil {
		return ""
	}
	return x.ScanPoint
}

type RecordLuggageScanResponse struct {
	ScanID    string
	TotalBags int32
}

func (x *RecordLuggageScanResponse) GetScanID() string {
	if x == nil {
		return ""
	}
	return x.ScanID
}
func (x *RecordLuggageScanResponse) GetTotalBags() int32 {
	if x == nil {
		return 0
	}
	return x.TotalBags
}

type GetLuggageCountRequest struct {
	DepartureID string
}

func (x *GetLuggageCountRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type LuggageCountRow struct {
	PilgrimID     string
	BagCount      int32
	LastScannedAt string
}

func (x *LuggageCountRow) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *LuggageCountRow) GetBagCount() int32 {
	if x == nil {
		return 0
	}
	return x.BagCount
}
func (x *LuggageCountRow) GetLastScannedAt() string {
	if x == nil {
		return ""
	}
	return x.LastScannedAt
}

type GetLuggageCountResponse struct {
	DepartureID string
	TotalBags   int32
	Rows        []*LuggageCountRow
}

func (x *GetLuggageCountResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *GetLuggageCountResponse) GetTotalBags() int32 {
	if x == nil {
		return 0
	}
	return x.TotalBags
}
func (x *GetLuggageCountResponse) GetRows() []*LuggageCountRow {
	if x == nil {
		return nil
	}
	return x.Rows
}

// ---------------------------------------------------------------------------
// BL-OPS-037: Departure/arrival broadcast
// ---------------------------------------------------------------------------

type BroadcastScheduleRequest struct {
	DepartureID   string
	BroadcastType string // "departure"/"arrival"
	Message       string
	Channel       string // "whatsapp"/"sms"/"email"
}

func (x *BroadcastScheduleRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *BroadcastScheduleRequest) GetBroadcastType() string {
	if x == nil {
		return ""
	}
	return x.BroadcastType
}
func (x *BroadcastScheduleRequest) GetMessage() string {
	if x == nil {
		return ""
	}
	return x.Message
}
func (x *BroadcastScheduleRequest) GetChannel() string {
	if x == nil {
		return ""
	}
	return x.Channel
}

type BroadcastScheduleResponse struct {
	BroadcastID    string
	RecipientCount int32
}

func (x *BroadcastScheduleResponse) GetBroadcastID() string {
	if x == nil {
		return ""
	}
	return x.BroadcastID
}
func (x *BroadcastScheduleResponse) GetRecipientCount() int32 {
	if x == nil {
		return 0
	}
	return x.RecipientCount
}

// ---------------------------------------------------------------------------
// BL-OPS-039: Raudhah shield & tasreh
// ---------------------------------------------------------------------------

type IssueDigitalTasrehRequest struct {
	PilgrimID   string
	DepartureID string
	VisitDate   string
}

func (x *IssueDigitalTasrehRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *IssueDigitalTasrehRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *IssueDigitalTasrehRequest) GetVisitDate() string {
	if x == nil {
		return ""
	}
	return x.VisitDate
}

type IssueDigitalTasrehResponse struct {
	TasrehID string
	QRCode   string
}

func (x *IssueDigitalTasrehResponse) GetTasrehID() string {
	if x == nil {
		return ""
	}
	return x.TasrehID
}
func (x *IssueDigitalTasrehResponse) GetQRCode() string {
	if x == nil {
		return ""
	}
	return x.QRCode
}

type RecordRaudhahEntryRequest struct {
	TasrehID  string
	PilgrimID string
	EntryTime string
}

func (x *RecordRaudhahEntryRequest) GetTasrehID() string {
	if x == nil {
		return ""
	}
	return x.TasrehID
}
func (x *RecordRaudhahEntryRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *RecordRaudhahEntryRequest) GetEntryTime() string {
	if x == nil {
		return ""
	}
	return x.EntryTime
}

type RecordRaudhahEntryResponse struct {
	RecordID string
	Valid    bool
}

func (x *RecordRaudhahEntryResponse) GetRecordID() string {
	if x == nil {
		return ""
	}
	return x.RecordID
}
func (x *RecordRaudhahEntryResponse) GetValid() bool {
	if x == nil {
		return false
	}
	return x.Valid
}

// ---------------------------------------------------------------------------
// BL-OPS-040: Audio devices
// ---------------------------------------------------------------------------

type RegisterAudioDeviceRequest struct {
	DepartureID  string
	DeviceType   string
	SerialNumber string
	AssignedTo   string
}

func (x *RegisterAudioDeviceRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *RegisterAudioDeviceRequest) GetDeviceType() string {
	if x == nil {
		return ""
	}
	return x.DeviceType
}
func (x *RegisterAudioDeviceRequest) GetSerialNumber() string {
	if x == nil {
		return ""
	}
	return x.SerialNumber
}
func (x *RegisterAudioDeviceRequest) GetAssignedTo() string {
	if x == nil {
		return ""
	}
	return x.AssignedTo
}

type RegisterAudioDeviceResponse struct {
	DeviceID string
}

func (x *RegisterAudioDeviceResponse) GetDeviceID() string {
	if x == nil {
		return ""
	}
	return x.DeviceID
}

type UpdateAudioDeviceStatusRequest struct {
	DeviceID string
	Status   string // "active"/"lost"/"returned"/"damaged"
	Notes    string
}

func (x *UpdateAudioDeviceStatusRequest) GetDeviceID() string {
	if x == nil {
		return ""
	}
	return x.DeviceID
}
func (x *UpdateAudioDeviceStatusRequest) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *UpdateAudioDeviceStatusRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type UpdateAudioDeviceStatusResponse struct {
	DeviceID string
	Status   string
}

func (x *UpdateAudioDeviceStatusResponse) GetDeviceID() string {
	if x == nil {
		return ""
	}
	return x.DeviceID
}
func (x *UpdateAudioDeviceStatusResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

// ---------------------------------------------------------------------------
// BL-OPS-041: Zamzam distribution
// ---------------------------------------------------------------------------

type RecordZamzamDistributionRequest struct {
	DepartureID string
	PilgrimID   string
	LitersGiven float64
	ReceivedBy  string
}

func (x *RecordZamzamDistributionRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *RecordZamzamDistributionRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *RecordZamzamDistributionRequest) GetLitersGiven() float64 {
	if x == nil {
		return 0
	}
	return x.LitersGiven
}
func (x *RecordZamzamDistributionRequest) GetReceivedBy() string {
	if x == nil {
		return ""
	}
	return x.ReceivedBy
}

type RecordZamzamDistributionResponse struct {
	DistributionID string
	RecordedAt     string
}

func (x *RecordZamzamDistributionResponse) GetDistributionID() string {
	if x == nil {
		return ""
	}
	return x.DistributionID
}
func (x *RecordZamzamDistributionResponse) GetRecordedAt() string {
	if x == nil {
		return ""
	}
	return x.RecordedAt
}

// ---------------------------------------------------------------------------
// BL-OPS-042: Room check-in
// ---------------------------------------------------------------------------

type RecordRoomCheckInRequest struct {
	DepartureID string
	PilgrimID   string
	RoomNumber  string
	HotelID     string
}

func (x *RecordRoomCheckInRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *RecordRoomCheckInRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *RecordRoomCheckInRequest) GetRoomNumber() string {
	if x == nil {
		return ""
	}
	return x.RoomNumber
}
func (x *RecordRoomCheckInRequest) GetHotelID() string {
	if x == nil {
		return ""
	}
	return x.HotelID
}

type RecordRoomCheckInResponse struct {
	CheckInID   string
	CheckedInAt string
}

func (x *RecordRoomCheckInResponse) GetCheckInID() string {
	if x == nil {
		return ""
	}
	return x.CheckInID
}
func (x *RecordRoomCheckInResponse) GetCheckedInAt() string {
	if x == nil {
		return ""
	}
	return x.CheckedInAt
}
