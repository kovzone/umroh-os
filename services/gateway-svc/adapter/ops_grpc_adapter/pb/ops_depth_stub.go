// ops_depth_stub.go — gateway-side gRPC client stubs for ops-svc Wave 5 depth RPCs.
// BL-OPS-021..042: collective documents, passport OCR, mahram, document progress,
// official letters, immigration manifest, transport, manifest delta, staff assignment,
// passport log, visa progress, e-visa, external provider, refunds, luggage, broadcast,
// tasreh, raudhah, audio devices, zamzam, room check-in.
//
// Mirrors services/ops-svc/api/grpc_api/pb/ops_depth_messages.go.
// Run `make genpb` to replace with generated code once protoc is available.
package pb

import (
	context "context"

	grpc "google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// Method name constants
// ---------------------------------------------------------------------------

const (
	OpsService_StoreCollectiveDocument_FullMethodName      = "/pb.OpsService/StoreCollectiveDocument"
	OpsService_GetCollectiveDocuments_FullMethodName       = "/pb.OpsService/GetCollectiveDocuments"
	OpsService_SetDocumentACL_FullMethodName               = "/pb.OpsService/SetDocumentACL"
	OpsService_ExtractPassportOCR_FullMethodName           = "/pb.OpsService/ExtractPassportOCR"
	OpsService_SetMahramRelation_FullMethodName            = "/pb.OpsService/SetMahramRelation"
	OpsService_GetMahramRelations_FullMethodName           = "/pb.OpsService/GetMahramRelations"
	OpsService_GetDocumentProgress_FullMethodName          = "/pb.OpsService/GetDocumentProgress"
	OpsService_GetExpiryAlerts_FullMethodName              = "/pb.OpsService/GetExpiryAlerts"
	OpsService_GenerateOfficialLetter_FullMethodName       = "/pb.OpsService/GenerateOfficialLetter"
	OpsService_GenerateImmigrationManifest_FullMethodName  = "/pb.OpsService/GenerateImmigrationManifest"
	OpsService_AssignTransport_FullMethodName              = "/pb.OpsService/AssignTransport"
	OpsService_GetTransportAssignments_FullMethodName      = "/pb.OpsService/GetTransportAssignments"
	OpsService_PublishManifestDelta_FullMethodName         = "/pb.OpsService/PublishManifestDelta"
	OpsService_AssignStaff_FullMethodName                  = "/pb.OpsService/AssignStaff"
	OpsService_RecordPassportHandover_FullMethodName       = "/pb.OpsService/RecordPassportHandover"
	OpsService_GetPassportLog_FullMethodName               = "/pb.OpsService/GetPassportLog"
	OpsService_GetVisaProgress_FullMethodName              = "/pb.OpsService/GetVisaProgress"
	OpsService_StoreEVisa_FullMethodName                   = "/pb.OpsService/StoreEVisa"
	OpsService_GetEVisa_FullMethodName                     = "/pb.OpsService/GetEVisa"
	OpsService_TriggerExternalProvider_FullMethodName      = "/pb.OpsService/TriggerExternalProvider"
	OpsService_CreateRefund_FullMethodName                 = "/pb.OpsService/CreateRefund"
	OpsService_ApproveRefund_FullMethodName                = "/pb.OpsService/ApproveRefund"
	OpsService_RecordPenalty_FullMethodName                = "/pb.OpsService/RecordPenalty"
	OpsService_RecordLuggageScan_FullMethodName            = "/pb.OpsService/RecordLuggageScan"
	OpsService_GetLuggageCount_FullMethodName              = "/pb.OpsService/GetLuggageCount"
	OpsService_BroadcastSchedule_FullMethodName            = "/pb.OpsService/BroadcastSchedule"
	OpsService_IssueDigitalTasreh_FullMethodName           = "/pb.OpsService/IssueDigitalTasreh"
	OpsService_RecordRaudhahEntry_FullMethodName           = "/pb.OpsService/RecordRaudhahEntry"
	OpsService_RegisterAudioDevice_FullMethodName          = "/pb.OpsService/RegisterAudioDevice"
	OpsService_UpdateAudioDeviceStatus_FullMethodName      = "/pb.OpsService/UpdateAudioDeviceStatus"
	OpsService_RecordZamzamDistribution_FullMethodName     = "/pb.OpsService/RecordZamzamDistribution"
	OpsService_RecordRoomCheckIn_FullMethodName            = "/pb.OpsService/RecordRoomCheckIn"
)

// ---------------------------------------------------------------------------
// BL-OPS-021: Collective document storage
// ---------------------------------------------------------------------------

type OpsDepthStoreCollectiveDocumentRequest struct {
	DepartureID  string
	DocumentType string
	URL          string
	PilgrimID    string
	Notes        string
}

func (x *OpsDepthStoreCollectiveDocumentRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthStoreCollectiveDocumentRequest) GetDocumentType() string {
	if x == nil {
		return ""
	}
	return x.DocumentType
}
func (x *OpsDepthStoreCollectiveDocumentRequest) GetURL() string {
	if x == nil {
		return ""
	}
	return x.URL
}
func (x *OpsDepthStoreCollectiveDocumentRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthStoreCollectiveDocumentRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type OpsDepthStoreCollectiveDocumentResponse struct {
	DocumentID string
}

func (x *OpsDepthStoreCollectiveDocumentResponse) GetDocumentID() string {
	if x == nil {
		return ""
	}
	return x.DocumentID
}

type OpsDepthGetCollectiveDocumentsRequest struct {
	DepartureID  string
	PilgrimID    string
	DocumentType string
}

func (x *OpsDepthGetCollectiveDocumentsRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthGetCollectiveDocumentsRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthGetCollectiveDocumentsRequest) GetDocumentType() string {
	if x == nil {
		return ""
	}
	return x.DocumentType
}

type OpsDepthCollectiveDocRow struct {
	DocumentID   string
	DocumentType string
	URL          string
	PilgrimID    string
	UploadedAt   string
}

func (x *OpsDepthCollectiveDocRow) GetDocumentID() string {
	if x == nil {
		return ""
	}
	return x.DocumentID
}
func (x *OpsDepthCollectiveDocRow) GetDocumentType() string {
	if x == nil {
		return ""
	}
	return x.DocumentType
}
func (x *OpsDepthCollectiveDocRow) GetURL() string {
	if x == nil {
		return ""
	}
	return x.URL
}
func (x *OpsDepthCollectiveDocRow) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthCollectiveDocRow) GetUploadedAt() string {
	if x == nil {
		return ""
	}
	return x.UploadedAt
}

type OpsDepthGetCollectiveDocumentsResponse struct {
	DepartureID string
	Documents   []*OpsDepthCollectiveDocRow
}

func (x *OpsDepthGetCollectiveDocumentsResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthGetCollectiveDocumentsResponse) GetDocuments() []*OpsDepthCollectiveDocRow {
	if x == nil {
		return nil
	}
	return x.Documents
}

type OpsDepthSetDocumentACLRequest struct {
	DocumentID     string
	AccessLevel    string
	AllowedUserIDs []string
}

func (x *OpsDepthSetDocumentACLRequest) GetDocumentID() string {
	if x == nil {
		return ""
	}
	return x.DocumentID
}
func (x *OpsDepthSetDocumentACLRequest) GetAccessLevel() string {
	if x == nil {
		return ""
	}
	return x.AccessLevel
}
func (x *OpsDepthSetDocumentACLRequest) GetAllowedUserIDs() []string {
	if x == nil {
		return nil
	}
	return x.AllowedUserIDs
}

type OpsDepthSetDocumentACLResponse struct {
	Updated bool
}

func (x *OpsDepthSetDocumentACLResponse) GetUpdated() bool {
	if x == nil {
		return false
	}
	return x.Updated
}

// ---------------------------------------------------------------------------
// BL-OPS-022: Passport OCR & mahram
// ---------------------------------------------------------------------------

type OpsDepthExtractPassportOCRRequest struct {
	ImageURL  string
	PilgrimID string
}

func (x *OpsDepthExtractPassportOCRRequest) GetImageURL() string {
	if x == nil {
		return ""
	}
	return x.ImageURL
}
func (x *OpsDepthExtractPassportOCRRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}

type OpsDepthPassportOCRData struct {
	FullName       string
	PassportNumber string
	Nationality    string
	DateOfBirth    string
	ExpiryDate     string
	Gender         string
	Confidence     float64
}

func (x *OpsDepthPassportOCRData) GetFullName() string {
	if x == nil {
		return ""
	}
	return x.FullName
}
func (x *OpsDepthPassportOCRData) GetPassportNumber() string {
	if x == nil {
		return ""
	}
	return x.PassportNumber
}
func (x *OpsDepthPassportOCRData) GetNationality() string {
	if x == nil {
		return ""
	}
	return x.Nationality
}
func (x *OpsDepthPassportOCRData) GetDateOfBirth() string {
	if x == nil {
		return ""
	}
	return x.DateOfBirth
}
func (x *OpsDepthPassportOCRData) GetExpiryDate() string {
	if x == nil {
		return ""
	}
	return x.ExpiryDate
}
func (x *OpsDepthPassportOCRData) GetGender() string {
	if x == nil {
		return ""
	}
	return x.Gender
}
func (x *OpsDepthPassportOCRData) GetConfidence() float64 {
	if x == nil {
		return 0
	}
	return x.Confidence
}

type OpsDepthExtractPassportOCRResponse struct {
	PilgrimID string
	Data      *OpsDepthPassportOCRData
	Warnings  []string
}

func (x *OpsDepthExtractPassportOCRResponse) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthExtractPassportOCRResponse) GetData() *OpsDepthPassportOCRData {
	if x == nil {
		return nil
	}
	return x.Data
}
func (x *OpsDepthExtractPassportOCRResponse) GetWarnings() []string {
	if x == nil {
		return nil
	}
	return x.Warnings
}

type OpsDepthSetMahramRelationRequest struct {
	PilgrimID       string
	MahramPilgrimID string
	Relation        string
}

func (x *OpsDepthSetMahramRelationRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthSetMahramRelationRequest) GetMahramPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.MahramPilgrimID
}
func (x *OpsDepthSetMahramRelationRequest) GetRelation() string {
	if x == nil {
		return ""
	}
	return x.Relation
}

type OpsDepthSetMahramRelationResponse struct {
	RelationID string
}

func (x *OpsDepthSetMahramRelationResponse) GetRelationID() string {
	if x == nil {
		return ""
	}
	return x.RelationID
}

type OpsDepthGetMahramRelationsRequest struct {
	BookingID string
}

func (x *OpsDepthGetMahramRelationsRequest) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}

type OpsDepthMahramRelationRow struct {
	RelationID      string
	PilgrimID       string
	MahramPilgrimID string
	Relation        string
}

func (x *OpsDepthMahramRelationRow) GetRelationID() string {
	if x == nil {
		return ""
	}
	return x.RelationID
}
func (x *OpsDepthMahramRelationRow) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthMahramRelationRow) GetMahramPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.MahramPilgrimID
}
func (x *OpsDepthMahramRelationRow) GetRelation() string {
	if x == nil {
		return ""
	}
	return x.Relation
}

type OpsDepthGetMahramRelationsResponse struct {
	BookingID string
	Relations []*OpsDepthMahramRelationRow
}

func (x *OpsDepthGetMahramRelationsResponse) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *OpsDepthGetMahramRelationsResponse) GetRelations() []*OpsDepthMahramRelationRow {
	if x == nil {
		return nil
	}
	return x.Relations
}

// ---------------------------------------------------------------------------
// BL-OPS-023: Document progress & expiry
// ---------------------------------------------------------------------------

type OpsDepthGetDocumentProgressRequest struct {
	DepartureID string
}

func (x *OpsDepthGetDocumentProgressRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type OpsDepthDocProgressRow struct {
	PilgrimID       string
	DocumentType    string
	Status          string
	ExpiryDate      string
	DaysUntilExpiry int32
}

func (x *OpsDepthDocProgressRow) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthDocProgressRow) GetDocumentType() string {
	if x == nil {
		return ""
	}
	return x.DocumentType
}
func (x *OpsDepthDocProgressRow) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *OpsDepthDocProgressRow) GetExpiryDate() string {
	if x == nil {
		return ""
	}
	return x.ExpiryDate
}
func (x *OpsDepthDocProgressRow) GetDaysUntilExpiry() int32 {
	if x == nil {
		return 0
	}
	return x.DaysUntilExpiry
}

type OpsDepthGetDocumentProgressResponse struct {
	DepartureID       string
	Rows              []*OpsDepthDocProgressRow
	TotalPilgrims     int32
	DocumentsComplete int32
	DocumentsExpiring int32
}

func (x *OpsDepthGetDocumentProgressResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthGetDocumentProgressResponse) GetRows() []*OpsDepthDocProgressRow {
	if x == nil {
		return nil
	}
	return x.Rows
}
func (x *OpsDepthGetDocumentProgressResponse) GetTotalPilgrims() int32 {
	if x == nil {
		return 0
	}
	return x.TotalPilgrims
}
func (x *OpsDepthGetDocumentProgressResponse) GetDocumentsComplete() int32 {
	if x == nil {
		return 0
	}
	return x.DocumentsComplete
}
func (x *OpsDepthGetDocumentProgressResponse) GetDocumentsExpiring() int32 {
	if x == nil {
		return 0
	}
	return x.DocumentsExpiring
}

type OpsDepthGetExpiryAlertsRequest struct {
	ThresholdDays int32
	DepartureID   string
}

func (x *OpsDepthGetExpiryAlertsRequest) GetThresholdDays() int32 {
	if x == nil {
		return 0
	}
	return x.ThresholdDays
}
func (x *OpsDepthGetExpiryAlertsRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type OpsDepthExpiryAlertRow struct {
	PilgrimID       string
	DocumentType    string
	ExpiryDate      string
	DaysUntilExpiry int32
}

func (x *OpsDepthExpiryAlertRow) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthExpiryAlertRow) GetDocumentType() string {
	if x == nil {
		return ""
	}
	return x.DocumentType
}
func (x *OpsDepthExpiryAlertRow) GetExpiryDate() string {
	if x == nil {
		return ""
	}
	return x.ExpiryDate
}
func (x *OpsDepthExpiryAlertRow) GetDaysUntilExpiry() int32 {
	if x == nil {
		return 0
	}
	return x.DaysUntilExpiry
}

type OpsDepthGetExpiryAlertsResponse struct {
	Alerts []*OpsDepthExpiryAlertRow
}

func (x *OpsDepthGetExpiryAlertsResponse) GetAlerts() []*OpsDepthExpiryAlertRow {
	if x == nil {
		return nil
	}
	return x.Alerts
}

// ---------------------------------------------------------------------------
// BL-OPS-024: Official letter
// ---------------------------------------------------------------------------

type OpsDepthGenerateOfficialLetterRequest struct {
	TemplateName string
	DepartureID  string
	PilgrimID    string
	IssuedTo     string
	Notes        string
}

func (x *OpsDepthGenerateOfficialLetterRequest) GetTemplateName() string {
	if x == nil {
		return ""
	}
	return x.TemplateName
}
func (x *OpsDepthGenerateOfficialLetterRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthGenerateOfficialLetterRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthGenerateOfficialLetterRequest) GetIssuedTo() string {
	if x == nil {
		return ""
	}
	return x.IssuedTo
}
func (x *OpsDepthGenerateOfficialLetterRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type OpsDepthGenerateOfficialLetterResponse struct {
	LetterID  string
	LetterURL string
	IssuedAt  string
}

func (x *OpsDepthGenerateOfficialLetterResponse) GetLetterID() string {
	if x == nil {
		return ""
	}
	return x.LetterID
}
func (x *OpsDepthGenerateOfficialLetterResponse) GetLetterURL() string {
	if x == nil {
		return ""
	}
	return x.LetterURL
}
func (x *OpsDepthGenerateOfficialLetterResponse) GetIssuedAt() string {
	if x == nil {
		return ""
	}
	return x.IssuedAt
}

// ---------------------------------------------------------------------------
// BL-OPS-025: Immigration manifest
// ---------------------------------------------------------------------------

type OpsDepthGenerateImmigrationManifestRequest struct {
	DepartureID string
	Format      string
}

func (x *OpsDepthGenerateImmigrationManifestRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthGenerateImmigrationManifestRequest) GetFormat() string {
	if x == nil {
		return ""
	}
	return x.Format
}

type OpsDepthGenerateImmigrationManifestResponse struct {
	ManifestID  string
	ManifestURL string
	Version     string
	RowCount    int32
}

func (x *OpsDepthGenerateImmigrationManifestResponse) GetManifestID() string {
	if x == nil {
		return ""
	}
	return x.ManifestID
}
func (x *OpsDepthGenerateImmigrationManifestResponse) GetManifestURL() string {
	if x == nil {
		return ""
	}
	return x.ManifestURL
}
func (x *OpsDepthGenerateImmigrationManifestResponse) GetVersion() string {
	if x == nil {
		return ""
	}
	return x.Version
}
func (x *OpsDepthGenerateImmigrationManifestResponse) GetRowCount() int32 {
	if x == nil {
		return 0
	}
	return x.RowCount
}

// ---------------------------------------------------------------------------
// BL-OPS-027: Transport arrangement
// ---------------------------------------------------------------------------

type OpsDepthAssignTransportRequest struct {
	DepartureID string
	VehicleType string
	VehicleID   string
	PilgrimIDs  []string
}

func (x *OpsDepthAssignTransportRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthAssignTransportRequest) GetVehicleType() string {
	if x == nil {
		return ""
	}
	return x.VehicleType
}
func (x *OpsDepthAssignTransportRequest) GetVehicleID() string {
	if x == nil {
		return ""
	}
	return x.VehicleID
}
func (x *OpsDepthAssignTransportRequest) GetPilgrimIDs() []string {
	if x == nil {
		return nil
	}
	return x.PilgrimIDs
}

type OpsDepthAssignTransportResponse struct {
	AssignmentID  string
	AssignedCount int32
}

func (x *OpsDepthAssignTransportResponse) GetAssignmentID() string {
	if x == nil {
		return ""
	}
	return x.AssignmentID
}
func (x *OpsDepthAssignTransportResponse) GetAssignedCount() int32 {
	if x == nil {
		return 0
	}
	return x.AssignedCount
}

type OpsDepthGetTransportAssignmentsRequest struct {
	DepartureID string
}

func (x *OpsDepthGetTransportAssignmentsRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type OpsDepthTransportAssignmentRow struct {
	AssignmentID string
	VehicleType  string
	VehicleID    string
	PilgrimIDs   []string
}

func (x *OpsDepthTransportAssignmentRow) GetAssignmentID() string {
	if x == nil {
		return ""
	}
	return x.AssignmentID
}
func (x *OpsDepthTransportAssignmentRow) GetVehicleType() string {
	if x == nil {
		return ""
	}
	return x.VehicleType
}
func (x *OpsDepthTransportAssignmentRow) GetVehicleID() string {
	if x == nil {
		return ""
	}
	return x.VehicleID
}
func (x *OpsDepthTransportAssignmentRow) GetPilgrimIDs() []string {
	if x == nil {
		return nil
	}
	return x.PilgrimIDs
}

type OpsDepthGetTransportAssignmentsResponse struct {
	DepartureID string
	Assignments []*OpsDepthTransportAssignmentRow
}

func (x *OpsDepthGetTransportAssignmentsResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthGetTransportAssignmentsResponse) GetAssignments() []*OpsDepthTransportAssignmentRow {
	if x == nil {
		return nil
	}
	return x.Assignments
}

// ---------------------------------------------------------------------------
// BL-OPS-028: Manifest delta
// ---------------------------------------------------------------------------

type OpsDepthPublishManifestDeltaRequest struct {
	DepartureID string
	ChangeType  string
	EntityID    string
	Notes       string
}

func (x *OpsDepthPublishManifestDeltaRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthPublishManifestDeltaRequest) GetChangeType() string {
	if x == nil {
		return ""
	}
	return x.ChangeType
}
func (x *OpsDepthPublishManifestDeltaRequest) GetEntityID() string {
	if x == nil {
		return ""
	}
	return x.EntityID
}
func (x *OpsDepthPublishManifestDeltaRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type OpsDepthPublishManifestDeltaResponse struct {
	DeltaID     string
	PublishedAt string
}

func (x *OpsDepthPublishManifestDeltaResponse) GetDeltaID() string {
	if x == nil {
		return ""
	}
	return x.DeltaID
}
func (x *OpsDepthPublishManifestDeltaResponse) GetPublishedAt() string {
	if x == nil {
		return ""
	}
	return x.PublishedAt
}

// ---------------------------------------------------------------------------
// BL-OPS-029: Staff assignment
// ---------------------------------------------------------------------------

type OpsDepthAssignStaffRequest struct {
	DepartureID string
	StaffUserID string
	Role        string
	PilgrimIDs  []string
}

func (x *OpsDepthAssignStaffRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthAssignStaffRequest) GetStaffUserID() string {
	if x == nil {
		return ""
	}
	return x.StaffUserID
}
func (x *OpsDepthAssignStaffRequest) GetRole() string {
	if x == nil {
		return ""
	}
	return x.Role
}
func (x *OpsDepthAssignStaffRequest) GetPilgrimIDs() []string {
	if x == nil {
		return nil
	}
	return x.PilgrimIDs
}

type OpsDepthAssignStaffResponse struct {
	AssignmentID string
}

func (x *OpsDepthAssignStaffResponse) GetAssignmentID() string {
	if x == nil {
		return ""
	}
	return x.AssignmentID
}

// ---------------------------------------------------------------------------
// BL-OPS-030: Passport log
// ---------------------------------------------------------------------------

type OpsDepthRecordPassportHandoverRequest struct {
	DepartureID string
	PilgrimID   string
	FromUserID  string
	ToUserID    string
	Notes       string
}

func (x *OpsDepthRecordPassportHandoverRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthRecordPassportHandoverRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthRecordPassportHandoverRequest) GetFromUserID() string {
	if x == nil {
		return ""
	}
	return x.FromUserID
}
func (x *OpsDepthRecordPassportHandoverRequest) GetToUserID() string {
	if x == nil {
		return ""
	}
	return x.ToUserID
}
func (x *OpsDepthRecordPassportHandoverRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type OpsDepthRecordPassportHandoverResponse struct {
	HandoverID string
	RecordedAt string
}

func (x *OpsDepthRecordPassportHandoverResponse) GetHandoverID() string {
	if x == nil {
		return ""
	}
	return x.HandoverID
}
func (x *OpsDepthRecordPassportHandoverResponse) GetRecordedAt() string {
	if x == nil {
		return ""
	}
	return x.RecordedAt
}

type OpsDepthGetPassportLogRequest struct {
	DepartureID string
}

func (x *OpsDepthGetPassportLogRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type OpsDepthPassportHandoverRow struct {
	HandoverID string
	PilgrimID  string
	FromUserID string
	ToUserID   string
	RecordedAt string
	Notes      string
}

func (x *OpsDepthPassportHandoverRow) GetHandoverID() string {
	if x == nil {
		return ""
	}
	return x.HandoverID
}
func (x *OpsDepthPassportHandoverRow) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthPassportHandoverRow) GetFromUserID() string {
	if x == nil {
		return ""
	}
	return x.FromUserID
}
func (x *OpsDepthPassportHandoverRow) GetToUserID() string {
	if x == nil {
		return ""
	}
	return x.ToUserID
}
func (x *OpsDepthPassportHandoverRow) GetRecordedAt() string {
	if x == nil {
		return ""
	}
	return x.RecordedAt
}
func (x *OpsDepthPassportHandoverRow) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type OpsDepthGetPassportLogResponse struct {
	DepartureID string
	Rows        []*OpsDepthPassportHandoverRow
}

func (x *OpsDepthGetPassportLogResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthGetPassportLogResponse) GetRows() []*OpsDepthPassportHandoverRow {
	if x == nil {
		return nil
	}
	return x.Rows
}

// ---------------------------------------------------------------------------
// BL-OPS-031: Visa progress
// ---------------------------------------------------------------------------

type OpsDepthGetVisaProgressRequest struct {
	DepartureID string
}

func (x *OpsDepthGetVisaProgressRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type OpsDepthVisaProgressRow struct {
	PilgrimID     string
	VisaStatus    string
	SubmittedAt   string
	ExpectedBy    string
	DaysRemaining int32
}

func (x *OpsDepthVisaProgressRow) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthVisaProgressRow) GetVisaStatus() string {
	if x == nil {
		return ""
	}
	return x.VisaStatus
}
func (x *OpsDepthVisaProgressRow) GetSubmittedAt() string {
	if x == nil {
		return ""
	}
	return x.SubmittedAt
}
func (x *OpsDepthVisaProgressRow) GetExpectedBy() string {
	if x == nil {
		return ""
	}
	return x.ExpectedBy
}
func (x *OpsDepthVisaProgressRow) GetDaysRemaining() int32 {
	if x == nil {
		return 0
	}
	return x.DaysRemaining
}

type OpsDepthGetVisaProgressResponse struct {
	DepartureID string
	Rows        []*OpsDepthVisaProgressRow
	Submitted   int32
	Approved    int32
	Rejected    int32
	Pending     int32
}

func (x *OpsDepthGetVisaProgressResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthGetVisaProgressResponse) GetRows() []*OpsDepthVisaProgressRow {
	if x == nil {
		return nil
	}
	return x.Rows
}
func (x *OpsDepthGetVisaProgressResponse) GetSubmitted() int32 {
	if x == nil {
		return 0
	}
	return x.Submitted
}
func (x *OpsDepthGetVisaProgressResponse) GetApproved() int32 {
	if x == nil {
		return 0
	}
	return x.Approved
}
func (x *OpsDepthGetVisaProgressResponse) GetRejected() int32 {
	if x == nil {
		return 0
	}
	return x.Rejected
}
func (x *OpsDepthGetVisaProgressResponse) GetPending() int32 {
	if x == nil {
		return 0
	}
	return x.Pending
}

// ---------------------------------------------------------------------------
// BL-OPS-032: E-visa repository
// ---------------------------------------------------------------------------

type OpsDepthStoreEVisaRequest struct {
	PilgrimID   string
	DepartureID string
	VisaNumber  string
	VisaURL     string
	IssuedDate  string
	ExpiryDate  string
}

func (x *OpsDepthStoreEVisaRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthStoreEVisaRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthStoreEVisaRequest) GetVisaNumber() string {
	if x == nil {
		return ""
	}
	return x.VisaNumber
}
func (x *OpsDepthStoreEVisaRequest) GetVisaURL() string {
	if x == nil {
		return ""
	}
	return x.VisaURL
}
func (x *OpsDepthStoreEVisaRequest) GetIssuedDate() string {
	if x == nil {
		return ""
	}
	return x.IssuedDate
}
func (x *OpsDepthStoreEVisaRequest) GetExpiryDate() string {
	if x == nil {
		return ""
	}
	return x.ExpiryDate
}

type OpsDepthStoreEVisaResponse struct {
	EVisaID string
}

func (x *OpsDepthStoreEVisaResponse) GetEVisaID() string {
	if x == nil {
		return ""
	}
	return x.EVisaID
}

type OpsDepthGetEVisaRequest struct {
	PilgrimID   string
	DepartureID string
}

func (x *OpsDepthGetEVisaRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthGetEVisaRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type OpsDepthEVisa struct {
	EVisaID    string
	PilgrimID  string
	VisaNumber string
	VisaURL    string
	IssuedDate string
	ExpiryDate string
}

func (x *OpsDepthEVisa) GetEVisaID() string {
	if x == nil {
		return ""
	}
	return x.EVisaID
}
func (x *OpsDepthEVisa) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthEVisa) GetVisaNumber() string {
	if x == nil {
		return ""
	}
	return x.VisaNumber
}
func (x *OpsDepthEVisa) GetVisaURL() string {
	if x == nil {
		return ""
	}
	return x.VisaURL
}
func (x *OpsDepthEVisa) GetIssuedDate() string {
	if x == nil {
		return ""
	}
	return x.IssuedDate
}
func (x *OpsDepthEVisa) GetExpiryDate() string {
	if x == nil {
		return ""
	}
	return x.ExpiryDate
}

type OpsDepthGetEVisaResponse struct {
	EVisa *OpsDepthEVisa
}

func (x *OpsDepthGetEVisaResponse) GetEVisa() *OpsDepthEVisa {
	if x == nil {
		return nil
	}
	return x.EVisa
}

// ---------------------------------------------------------------------------
// BL-OPS-033: External provider
// ---------------------------------------------------------------------------

type OpsDepthTriggerExternalProviderRequest struct {
	Provider    string
	Action      string
	ReferenceID string
	Payload     string
}

func (x *OpsDepthTriggerExternalProviderRequest) GetProvider() string {
	if x == nil {
		return ""
	}
	return x.Provider
}
func (x *OpsDepthTriggerExternalProviderRequest) GetAction() string {
	if x == nil {
		return ""
	}
	return x.Action
}
func (x *OpsDepthTriggerExternalProviderRequest) GetReferenceID() string {
	if x == nil {
		return ""
	}
	return x.ReferenceID
}
func (x *OpsDepthTriggerExternalProviderRequest) GetPayload() string {
	if x == nil {
		return ""
	}
	return x.Payload
}

type OpsDepthTriggerExternalProviderResponse struct {
	RequestID string
	Status    string
}

func (x *OpsDepthTriggerExternalProviderResponse) GetRequestID() string {
	if x == nil {
		return ""
	}
	return x.RequestID
}
func (x *OpsDepthTriggerExternalProviderResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

// ---------------------------------------------------------------------------
// BL-OPS-034: Refund & penalty
// ---------------------------------------------------------------------------

type OpsDepthCreateRefundRequest struct {
	BookingID string
	Reason    string
	Amount    int64
	Notes     string
}

func (x *OpsDepthCreateRefundRequest) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *OpsDepthCreateRefundRequest) GetReason() string {
	if x == nil {
		return ""
	}
	return x.Reason
}
func (x *OpsDepthCreateRefundRequest) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *OpsDepthCreateRefundRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type OpsDepthCreateRefundResponse struct {
	RefundID string
	Status   string
}

func (x *OpsDepthCreateRefundResponse) GetRefundID() string {
	if x == nil {
		return ""
	}
	return x.RefundID
}
func (x *OpsDepthCreateRefundResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

type OpsDepthApproveRefundRequest struct {
	RefundID string
	Decision string
	Notes    string
}

func (x *OpsDepthApproveRefundRequest) GetRefundID() string {
	if x == nil {
		return ""
	}
	return x.RefundID
}
func (x *OpsDepthApproveRefundRequest) GetDecision() string {
	if x == nil {
		return ""
	}
	return x.Decision
}
func (x *OpsDepthApproveRefundRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type OpsDepthApproveRefundResponse struct {
	RefundID string
	Status   string
}

func (x *OpsDepthApproveRefundResponse) GetRefundID() string {
	if x == nil {
		return ""
	}
	return x.RefundID
}
func (x *OpsDepthApproveRefundResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

type OpsDepthRecordPenaltyRequest struct {
	BookingID   string
	PenaltyType string
	Amount      int64
	Notes       string
}

func (x *OpsDepthRecordPenaltyRequest) GetBookingID() string {
	if x == nil {
		return ""
	}
	return x.BookingID
}
func (x *OpsDepthRecordPenaltyRequest) GetPenaltyType() string {
	if x == nil {
		return ""
	}
	return x.PenaltyType
}
func (x *OpsDepthRecordPenaltyRequest) GetAmount() int64 {
	if x == nil {
		return 0
	}
	return x.Amount
}
func (x *OpsDepthRecordPenaltyRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type OpsDepthRecordPenaltyResponse struct {
	PenaltyID string
}

func (x *OpsDepthRecordPenaltyResponse) GetPenaltyID() string {
	if x == nil {
		return ""
	}
	return x.PenaltyID
}

// ---------------------------------------------------------------------------
// BL-OPS-036: Luggage counter
// ---------------------------------------------------------------------------

type OpsDepthRecordLuggageScanRequest struct {
	DepartureID string
	PilgrimID   string
	TagID       string
	ScanPoint   string
}

func (x *OpsDepthRecordLuggageScanRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthRecordLuggageScanRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthRecordLuggageScanRequest) GetTagID() string {
	if x == nil {
		return ""
	}
	return x.TagID
}
func (x *OpsDepthRecordLuggageScanRequest) GetScanPoint() string {
	if x == nil {
		return ""
	}
	return x.ScanPoint
}

type OpsDepthRecordLuggageScanResponse struct {
	ScanID    string
	TotalBags int32
}

func (x *OpsDepthRecordLuggageScanResponse) GetScanID() string {
	if x == nil {
		return ""
	}
	return x.ScanID
}
func (x *OpsDepthRecordLuggageScanResponse) GetTotalBags() int32 {
	if x == nil {
		return 0
	}
	return x.TotalBags
}

type OpsDepthGetLuggageCountRequest struct {
	DepartureID string
}

func (x *OpsDepthGetLuggageCountRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}

type OpsDepthLuggageCountRow struct {
	PilgrimID     string
	BagCount      int32
	LastScannedAt string
}

func (x *OpsDepthLuggageCountRow) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthLuggageCountRow) GetBagCount() int32 {
	if x == nil {
		return 0
	}
	return x.BagCount
}
func (x *OpsDepthLuggageCountRow) GetLastScannedAt() string {
	if x == nil {
		return ""
	}
	return x.LastScannedAt
}

type OpsDepthGetLuggageCountResponse struct {
	DepartureID string
	TotalBags   int32
	Rows        []*OpsDepthLuggageCountRow
}

func (x *OpsDepthGetLuggageCountResponse) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthGetLuggageCountResponse) GetTotalBags() int32 {
	if x == nil {
		return 0
	}
	return x.TotalBags
}
func (x *OpsDepthGetLuggageCountResponse) GetRows() []*OpsDepthLuggageCountRow {
	if x == nil {
		return nil
	}
	return x.Rows
}

// ---------------------------------------------------------------------------
// BL-OPS-037: Departure/arrival broadcast
// ---------------------------------------------------------------------------

type OpsDepthBroadcastScheduleRequest struct {
	DepartureID   string
	BroadcastType string
	Message       string
	Channel       string
}

func (x *OpsDepthBroadcastScheduleRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthBroadcastScheduleRequest) GetBroadcastType() string {
	if x == nil {
		return ""
	}
	return x.BroadcastType
}
func (x *OpsDepthBroadcastScheduleRequest) GetMessage() string {
	if x == nil {
		return ""
	}
	return x.Message
}
func (x *OpsDepthBroadcastScheduleRequest) GetChannel() string {
	if x == nil {
		return ""
	}
	return x.Channel
}

type OpsDepthBroadcastScheduleResponse struct {
	BroadcastID    string
	RecipientCount int32
}

func (x *OpsDepthBroadcastScheduleResponse) GetBroadcastID() string {
	if x == nil {
		return ""
	}
	return x.BroadcastID
}
func (x *OpsDepthBroadcastScheduleResponse) GetRecipientCount() int32 {
	if x == nil {
		return 0
	}
	return x.RecipientCount
}

// ---------------------------------------------------------------------------
// BL-OPS-039: Raudhah shield & tasreh
// ---------------------------------------------------------------------------

type OpsDepthIssueDigitalTasrehRequest struct {
	PilgrimID   string
	DepartureID string
	VisitDate   string
}

func (x *OpsDepthIssueDigitalTasrehRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthIssueDigitalTasrehRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthIssueDigitalTasrehRequest) GetVisitDate() string {
	if x == nil {
		return ""
	}
	return x.VisitDate
}

type OpsDepthIssueDigitalTasrehResponse struct {
	TasrehID string
	QRCode   string
}

func (x *OpsDepthIssueDigitalTasrehResponse) GetTasrehID() string {
	if x == nil {
		return ""
	}
	return x.TasrehID
}
func (x *OpsDepthIssueDigitalTasrehResponse) GetQRCode() string {
	if x == nil {
		return ""
	}
	return x.QRCode
}

type OpsDepthRecordRaudhahEntryRequest struct {
	TasrehID  string
	PilgrimID string
	EntryTime string
}

func (x *OpsDepthRecordRaudhahEntryRequest) GetTasrehID() string {
	if x == nil {
		return ""
	}
	return x.TasrehID
}
func (x *OpsDepthRecordRaudhahEntryRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthRecordRaudhahEntryRequest) GetEntryTime() string {
	if x == nil {
		return ""
	}
	return x.EntryTime
}

type OpsDepthRecordRaudhahEntryResponse struct {
	RecordID string
	Valid    bool
}

func (x *OpsDepthRecordRaudhahEntryResponse) GetRecordID() string {
	if x == nil {
		return ""
	}
	return x.RecordID
}
func (x *OpsDepthRecordRaudhahEntryResponse) GetValid() bool {
	if x == nil {
		return false
	}
	return x.Valid
}

// ---------------------------------------------------------------------------
// BL-OPS-040: Audio devices
// ---------------------------------------------------------------------------

type OpsDepthRegisterAudioDeviceRequest struct {
	DepartureID  string
	DeviceType   string
	SerialNumber string
	AssignedTo   string
}

func (x *OpsDepthRegisterAudioDeviceRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthRegisterAudioDeviceRequest) GetDeviceType() string {
	if x == nil {
		return ""
	}
	return x.DeviceType
}
func (x *OpsDepthRegisterAudioDeviceRequest) GetSerialNumber() string {
	if x == nil {
		return ""
	}
	return x.SerialNumber
}
func (x *OpsDepthRegisterAudioDeviceRequest) GetAssignedTo() string {
	if x == nil {
		return ""
	}
	return x.AssignedTo
}

type OpsDepthRegisterAudioDeviceResponse struct {
	DeviceID string
}

func (x *OpsDepthRegisterAudioDeviceResponse) GetDeviceID() string {
	if x == nil {
		return ""
	}
	return x.DeviceID
}

type OpsDepthUpdateAudioDeviceStatusRequest struct {
	DeviceID string
	Status   string
	Notes    string
}

func (x *OpsDepthUpdateAudioDeviceStatusRequest) GetDeviceID() string {
	if x == nil {
		return ""
	}
	return x.DeviceID
}
func (x *OpsDepthUpdateAudioDeviceStatusRequest) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}
func (x *OpsDepthUpdateAudioDeviceStatusRequest) GetNotes() string {
	if x == nil {
		return ""
	}
	return x.Notes
}

type OpsDepthUpdateAudioDeviceStatusResponse struct {
	DeviceID string
	Status   string
}

func (x *OpsDepthUpdateAudioDeviceStatusResponse) GetDeviceID() string {
	if x == nil {
		return ""
	}
	return x.DeviceID
}
func (x *OpsDepthUpdateAudioDeviceStatusResponse) GetStatus() string {
	if x == nil {
		return ""
	}
	return x.Status
}

// ---------------------------------------------------------------------------
// BL-OPS-041: Zamzam distribution
// ---------------------------------------------------------------------------

type OpsDepthRecordZamzamDistributionRequest struct {
	DepartureID string
	PilgrimID   string
	LitersGiven float64
	ReceivedBy  string
}

func (x *OpsDepthRecordZamzamDistributionRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthRecordZamzamDistributionRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthRecordZamzamDistributionRequest) GetLitersGiven() float64 {
	if x == nil {
		return 0
	}
	return x.LitersGiven
}
func (x *OpsDepthRecordZamzamDistributionRequest) GetReceivedBy() string {
	if x == nil {
		return ""
	}
	return x.ReceivedBy
}

type OpsDepthRecordZamzamDistributionResponse struct {
	DistributionID string
	RecordedAt     string
}

func (x *OpsDepthRecordZamzamDistributionResponse) GetDistributionID() string {
	if x == nil {
		return ""
	}
	return x.DistributionID
}
func (x *OpsDepthRecordZamzamDistributionResponse) GetRecordedAt() string {
	if x == nil {
		return ""
	}
	return x.RecordedAt
}

// ---------------------------------------------------------------------------
// BL-OPS-042: Room check-in
// ---------------------------------------------------------------------------

type OpsDepthRecordRoomCheckInRequest struct {
	DepartureID string
	PilgrimID   string
	RoomNumber  string
	HotelID     string
}

func (x *OpsDepthRecordRoomCheckInRequest) GetDepartureID() string {
	if x == nil {
		return ""
	}
	return x.DepartureID
}
func (x *OpsDepthRecordRoomCheckInRequest) GetPilgrimID() string {
	if x == nil {
		return ""
	}
	return x.PilgrimID
}
func (x *OpsDepthRecordRoomCheckInRequest) GetRoomNumber() string {
	if x == nil {
		return ""
	}
	return x.RoomNumber
}
func (x *OpsDepthRecordRoomCheckInRequest) GetHotelID() string {
	if x == nil {
		return ""
	}
	return x.HotelID
}

type OpsDepthRecordRoomCheckInResponse struct {
	CheckInID   string
	CheckedInAt string
}

func (x *OpsDepthRecordRoomCheckInResponse) GetCheckInID() string {
	if x == nil {
		return ""
	}
	return x.CheckInID
}
func (x *OpsDepthRecordRoomCheckInResponse) GetCheckedInAt() string {
	if x == nil {
		return ""
	}
	return x.CheckedInAt
}

// ---------------------------------------------------------------------------
// OpsDepthClient — gRPC client interface
// ---------------------------------------------------------------------------

// OpsDepthClient is the client API for ops-svc Wave 5 depth RPCs.
type OpsDepthClient interface {
	StoreCollectiveDocument(ctx context.Context, in *OpsDepthStoreCollectiveDocumentRequest, opts ...grpc.CallOption) (*OpsDepthStoreCollectiveDocumentResponse, error)
	GetCollectiveDocuments(ctx context.Context, in *OpsDepthGetCollectiveDocumentsRequest, opts ...grpc.CallOption) (*OpsDepthGetCollectiveDocumentsResponse, error)
	SetDocumentACL(ctx context.Context, in *OpsDepthSetDocumentACLRequest, opts ...grpc.CallOption) (*OpsDepthSetDocumentACLResponse, error)
	ExtractPassportOCR(ctx context.Context, in *OpsDepthExtractPassportOCRRequest, opts ...grpc.CallOption) (*OpsDepthExtractPassportOCRResponse, error)
	SetMahramRelation(ctx context.Context, in *OpsDepthSetMahramRelationRequest, opts ...grpc.CallOption) (*OpsDepthSetMahramRelationResponse, error)
	GetMahramRelations(ctx context.Context, in *OpsDepthGetMahramRelationsRequest, opts ...grpc.CallOption) (*OpsDepthGetMahramRelationsResponse, error)
	GetDocumentProgress(ctx context.Context, in *OpsDepthGetDocumentProgressRequest, opts ...grpc.CallOption) (*OpsDepthGetDocumentProgressResponse, error)
	GetExpiryAlerts(ctx context.Context, in *OpsDepthGetExpiryAlertsRequest, opts ...grpc.CallOption) (*OpsDepthGetExpiryAlertsResponse, error)
	GenerateOfficialLetter(ctx context.Context, in *OpsDepthGenerateOfficialLetterRequest, opts ...grpc.CallOption) (*OpsDepthGenerateOfficialLetterResponse, error)
	GenerateImmigrationManifest(ctx context.Context, in *OpsDepthGenerateImmigrationManifestRequest, opts ...grpc.CallOption) (*OpsDepthGenerateImmigrationManifestResponse, error)
	AssignTransport(ctx context.Context, in *OpsDepthAssignTransportRequest, opts ...grpc.CallOption) (*OpsDepthAssignTransportResponse, error)
	GetTransportAssignments(ctx context.Context, in *OpsDepthGetTransportAssignmentsRequest, opts ...grpc.CallOption) (*OpsDepthGetTransportAssignmentsResponse, error)
	PublishManifestDelta(ctx context.Context, in *OpsDepthPublishManifestDeltaRequest, opts ...grpc.CallOption) (*OpsDepthPublishManifestDeltaResponse, error)
	AssignStaff(ctx context.Context, in *OpsDepthAssignStaffRequest, opts ...grpc.CallOption) (*OpsDepthAssignStaffResponse, error)
	RecordPassportHandover(ctx context.Context, in *OpsDepthRecordPassportHandoverRequest, opts ...grpc.CallOption) (*OpsDepthRecordPassportHandoverResponse, error)
	GetPassportLog(ctx context.Context, in *OpsDepthGetPassportLogRequest, opts ...grpc.CallOption) (*OpsDepthGetPassportLogResponse, error)
	GetVisaProgress(ctx context.Context, in *OpsDepthGetVisaProgressRequest, opts ...grpc.CallOption) (*OpsDepthGetVisaProgressResponse, error)
	StoreEVisa(ctx context.Context, in *OpsDepthStoreEVisaRequest, opts ...grpc.CallOption) (*OpsDepthStoreEVisaResponse, error)
	GetEVisa(ctx context.Context, in *OpsDepthGetEVisaRequest, opts ...grpc.CallOption) (*OpsDepthGetEVisaResponse, error)
	TriggerExternalProvider(ctx context.Context, in *OpsDepthTriggerExternalProviderRequest, opts ...grpc.CallOption) (*OpsDepthTriggerExternalProviderResponse, error)
	CreateRefund(ctx context.Context, in *OpsDepthCreateRefundRequest, opts ...grpc.CallOption) (*OpsDepthCreateRefundResponse, error)
	ApproveRefund(ctx context.Context, in *OpsDepthApproveRefundRequest, opts ...grpc.CallOption) (*OpsDepthApproveRefundResponse, error)
	RecordPenalty(ctx context.Context, in *OpsDepthRecordPenaltyRequest, opts ...grpc.CallOption) (*OpsDepthRecordPenaltyResponse, error)
	RecordLuggageScan(ctx context.Context, in *OpsDepthRecordLuggageScanRequest, opts ...grpc.CallOption) (*OpsDepthRecordLuggageScanResponse, error)
	GetLuggageCount(ctx context.Context, in *OpsDepthGetLuggageCountRequest, opts ...grpc.CallOption) (*OpsDepthGetLuggageCountResponse, error)
	BroadcastSchedule(ctx context.Context, in *OpsDepthBroadcastScheduleRequest, opts ...grpc.CallOption) (*OpsDepthBroadcastScheduleResponse, error)
	IssueDigitalTasreh(ctx context.Context, in *OpsDepthIssueDigitalTasrehRequest, opts ...grpc.CallOption) (*OpsDepthIssueDigitalTasrehResponse, error)
	RecordRaudhahEntry(ctx context.Context, in *OpsDepthRecordRaudhahEntryRequest, opts ...grpc.CallOption) (*OpsDepthRecordRaudhahEntryResponse, error)
	RegisterAudioDevice(ctx context.Context, in *OpsDepthRegisterAudioDeviceRequest, opts ...grpc.CallOption) (*OpsDepthRegisterAudioDeviceResponse, error)
	UpdateAudioDeviceStatus(ctx context.Context, in *OpsDepthUpdateAudioDeviceStatusRequest, opts ...grpc.CallOption) (*OpsDepthUpdateAudioDeviceStatusResponse, error)
	RecordZamzamDistribution(ctx context.Context, in *OpsDepthRecordZamzamDistributionRequest, opts ...grpc.CallOption) (*OpsDepthRecordZamzamDistributionResponse, error)
	RecordRoomCheckIn(ctx context.Context, in *OpsDepthRecordRoomCheckInRequest, opts ...grpc.CallOption) (*OpsDepthRecordRoomCheckInResponse, error)
}

type opsDepthClient struct {
	cc grpc.ClientConnInterface
}

// NewOpsDepthClient returns an OpsDepthClient backed by the given conn.
func NewOpsDepthClient(cc grpc.ClientConnInterface) OpsDepthClient {
	return &opsDepthClient{cc}
}

func (c *opsDepthClient) StoreCollectiveDocument(ctx context.Context, in *OpsDepthStoreCollectiveDocumentRequest, opts ...grpc.CallOption) (*OpsDepthStoreCollectiveDocumentResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthStoreCollectiveDocumentResponse)
	if err := c.cc.Invoke(ctx, OpsService_StoreCollectiveDocument_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) GetCollectiveDocuments(ctx context.Context, in *OpsDepthGetCollectiveDocumentsRequest, opts ...grpc.CallOption) (*OpsDepthGetCollectiveDocumentsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthGetCollectiveDocumentsResponse)
	if err := c.cc.Invoke(ctx, OpsService_GetCollectiveDocuments_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) SetDocumentACL(ctx context.Context, in *OpsDepthSetDocumentACLRequest, opts ...grpc.CallOption) (*OpsDepthSetDocumentACLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthSetDocumentACLResponse)
	if err := c.cc.Invoke(ctx, OpsService_SetDocumentACL_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) ExtractPassportOCR(ctx context.Context, in *OpsDepthExtractPassportOCRRequest, opts ...grpc.CallOption) (*OpsDepthExtractPassportOCRResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthExtractPassportOCRResponse)
	if err := c.cc.Invoke(ctx, OpsService_ExtractPassportOCR_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) SetMahramRelation(ctx context.Context, in *OpsDepthSetMahramRelationRequest, opts ...grpc.CallOption) (*OpsDepthSetMahramRelationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthSetMahramRelationResponse)
	if err := c.cc.Invoke(ctx, OpsService_SetMahramRelation_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) GetMahramRelations(ctx context.Context, in *OpsDepthGetMahramRelationsRequest, opts ...grpc.CallOption) (*OpsDepthGetMahramRelationsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthGetMahramRelationsResponse)
	if err := c.cc.Invoke(ctx, OpsService_GetMahramRelations_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) GetDocumentProgress(ctx context.Context, in *OpsDepthGetDocumentProgressRequest, opts ...grpc.CallOption) (*OpsDepthGetDocumentProgressResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthGetDocumentProgressResponse)
	if err := c.cc.Invoke(ctx, OpsService_GetDocumentProgress_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) GetExpiryAlerts(ctx context.Context, in *OpsDepthGetExpiryAlertsRequest, opts ...grpc.CallOption) (*OpsDepthGetExpiryAlertsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthGetExpiryAlertsResponse)
	if err := c.cc.Invoke(ctx, OpsService_GetExpiryAlerts_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) GenerateOfficialLetter(ctx context.Context, in *OpsDepthGenerateOfficialLetterRequest, opts ...grpc.CallOption) (*OpsDepthGenerateOfficialLetterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthGenerateOfficialLetterResponse)
	if err := c.cc.Invoke(ctx, OpsService_GenerateOfficialLetter_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) GenerateImmigrationManifest(ctx context.Context, in *OpsDepthGenerateImmigrationManifestRequest, opts ...grpc.CallOption) (*OpsDepthGenerateImmigrationManifestResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthGenerateImmigrationManifestResponse)
	if err := c.cc.Invoke(ctx, OpsService_GenerateImmigrationManifest_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) AssignTransport(ctx context.Context, in *OpsDepthAssignTransportRequest, opts ...grpc.CallOption) (*OpsDepthAssignTransportResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthAssignTransportResponse)
	if err := c.cc.Invoke(ctx, OpsService_AssignTransport_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) GetTransportAssignments(ctx context.Context, in *OpsDepthGetTransportAssignmentsRequest, opts ...grpc.CallOption) (*OpsDepthGetTransportAssignmentsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthGetTransportAssignmentsResponse)
	if err := c.cc.Invoke(ctx, OpsService_GetTransportAssignments_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) PublishManifestDelta(ctx context.Context, in *OpsDepthPublishManifestDeltaRequest, opts ...grpc.CallOption) (*OpsDepthPublishManifestDeltaResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthPublishManifestDeltaResponse)
	if err := c.cc.Invoke(ctx, OpsService_PublishManifestDelta_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) AssignStaff(ctx context.Context, in *OpsDepthAssignStaffRequest, opts ...grpc.CallOption) (*OpsDepthAssignStaffResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthAssignStaffResponse)
	if err := c.cc.Invoke(ctx, OpsService_AssignStaff_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) RecordPassportHandover(ctx context.Context, in *OpsDepthRecordPassportHandoverRequest, opts ...grpc.CallOption) (*OpsDepthRecordPassportHandoverResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthRecordPassportHandoverResponse)
	if err := c.cc.Invoke(ctx, OpsService_RecordPassportHandover_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) GetPassportLog(ctx context.Context, in *OpsDepthGetPassportLogRequest, opts ...grpc.CallOption) (*OpsDepthGetPassportLogResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthGetPassportLogResponse)
	if err := c.cc.Invoke(ctx, OpsService_GetPassportLog_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) GetVisaProgress(ctx context.Context, in *OpsDepthGetVisaProgressRequest, opts ...grpc.CallOption) (*OpsDepthGetVisaProgressResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthGetVisaProgressResponse)
	if err := c.cc.Invoke(ctx, OpsService_GetVisaProgress_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) StoreEVisa(ctx context.Context, in *OpsDepthStoreEVisaRequest, opts ...grpc.CallOption) (*OpsDepthStoreEVisaResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthStoreEVisaResponse)
	if err := c.cc.Invoke(ctx, OpsService_StoreEVisa_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) GetEVisa(ctx context.Context, in *OpsDepthGetEVisaRequest, opts ...grpc.CallOption) (*OpsDepthGetEVisaResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthGetEVisaResponse)
	if err := c.cc.Invoke(ctx, OpsService_GetEVisa_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) TriggerExternalProvider(ctx context.Context, in *OpsDepthTriggerExternalProviderRequest, opts ...grpc.CallOption) (*OpsDepthTriggerExternalProviderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthTriggerExternalProviderResponse)
	if err := c.cc.Invoke(ctx, OpsService_TriggerExternalProvider_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) CreateRefund(ctx context.Context, in *OpsDepthCreateRefundRequest, opts ...grpc.CallOption) (*OpsDepthCreateRefundResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthCreateRefundResponse)
	if err := c.cc.Invoke(ctx, OpsService_CreateRefund_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) ApproveRefund(ctx context.Context, in *OpsDepthApproveRefundRequest, opts ...grpc.CallOption) (*OpsDepthApproveRefundResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthApproveRefundResponse)
	if err := c.cc.Invoke(ctx, OpsService_ApproveRefund_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) RecordPenalty(ctx context.Context, in *OpsDepthRecordPenaltyRequest, opts ...grpc.CallOption) (*OpsDepthRecordPenaltyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthRecordPenaltyResponse)
	if err := c.cc.Invoke(ctx, OpsService_RecordPenalty_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) RecordLuggageScan(ctx context.Context, in *OpsDepthRecordLuggageScanRequest, opts ...grpc.CallOption) (*OpsDepthRecordLuggageScanResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthRecordLuggageScanResponse)
	if err := c.cc.Invoke(ctx, OpsService_RecordLuggageScan_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) GetLuggageCount(ctx context.Context, in *OpsDepthGetLuggageCountRequest, opts ...grpc.CallOption) (*OpsDepthGetLuggageCountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthGetLuggageCountResponse)
	if err := c.cc.Invoke(ctx, OpsService_GetLuggageCount_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) BroadcastSchedule(ctx context.Context, in *OpsDepthBroadcastScheduleRequest, opts ...grpc.CallOption) (*OpsDepthBroadcastScheduleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthBroadcastScheduleResponse)
	if err := c.cc.Invoke(ctx, OpsService_BroadcastSchedule_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) IssueDigitalTasreh(ctx context.Context, in *OpsDepthIssueDigitalTasrehRequest, opts ...grpc.CallOption) (*OpsDepthIssueDigitalTasrehResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthIssueDigitalTasrehResponse)
	if err := c.cc.Invoke(ctx, OpsService_IssueDigitalTasreh_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) RecordRaudhahEntry(ctx context.Context, in *OpsDepthRecordRaudhahEntryRequest, opts ...grpc.CallOption) (*OpsDepthRecordRaudhahEntryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthRecordRaudhahEntryResponse)
	if err := c.cc.Invoke(ctx, OpsService_RecordRaudhahEntry_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) RegisterAudioDevice(ctx context.Context, in *OpsDepthRegisterAudioDeviceRequest, opts ...grpc.CallOption) (*OpsDepthRegisterAudioDeviceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthRegisterAudioDeviceResponse)
	if err := c.cc.Invoke(ctx, OpsService_RegisterAudioDevice_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) UpdateAudioDeviceStatus(ctx context.Context, in *OpsDepthUpdateAudioDeviceStatusRequest, opts ...grpc.CallOption) (*OpsDepthUpdateAudioDeviceStatusResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthUpdateAudioDeviceStatusResponse)
	if err := c.cc.Invoke(ctx, OpsService_UpdateAudioDeviceStatus_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) RecordZamzamDistribution(ctx context.Context, in *OpsDepthRecordZamzamDistributionRequest, opts ...grpc.CallOption) (*OpsDepthRecordZamzamDistributionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthRecordZamzamDistributionResponse)
	if err := c.cc.Invoke(ctx, OpsService_RecordZamzamDistribution_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *opsDepthClient) RecordRoomCheckIn(ctx context.Context, in *OpsDepthRecordRoomCheckInRequest, opts ...grpc.CallOption) (*OpsDepthRecordRoomCheckInResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpsDepthRecordRoomCheckInResponse)
	if err := c.cc.Invoke(ctx, OpsService_RecordRoomCheckIn_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}
