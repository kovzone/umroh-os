// ops_depth.go — gateway ops adapter methods for Wave 5 depth RPCs.
// BL-OPS-021..042: collective documents, passport OCR, mahram, document progress,
// official letters, immigration manifest, transport, manifest delta, staff assignment,
// passport log, visa progress, e-visa, external provider, refunds, luggage, broadcast,
// tasreh, raudhah, audio devices, zamzam, room check-in.

package ops_grpc_adapter

import (
	"context"
	"errors"
	"fmt"

	"gateway-svc/adapter/ops_grpc_adapter/pb"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// Error mapper
// ---------------------------------------------------------------------------

func mapOpsDepthError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return errors.Join(apperrors.ErrServiceUnavailable, fmt.Errorf("ops-depth call failed: %w", err))
	}
	switch st.Code() {
	case grpcCodes.NotFound:
		return fmt.Errorf("%w: %s", apperrors.ErrNotFound, st.Message())
	case grpcCodes.InvalidArgument:
		return fmt.Errorf("%w: %s", apperrors.ErrValidation, st.Message())
	case grpcCodes.PermissionDenied:
		return fmt.Errorf("%w: %s", apperrors.ErrForbidden, st.Message())
	default:
		return fmt.Errorf("%w: ops-depth %s: %s", apperrors.ErrInternal, st.Code(), st.Message())
	}
}

// ---------------------------------------------------------------------------
// BL-OPS-021: Collective document storage
// ---------------------------------------------------------------------------

type StoreCollectiveDocumentParams struct {
	DepartureID  string
	DocumentType string
	URL          string
	PilgrimID    string
	Notes        string
}

type StoreCollectiveDocumentResult struct {
	DocumentID string
}

func (a *Adapter) StoreCollectiveDocument(ctx context.Context, params *StoreCollectiveDocumentParams) (*StoreCollectiveDocumentResult, error) {
	const op = "ops_grpc_adapter.Adapter.StoreCollectiveDocument"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.StoreCollectiveDocument(ctx, &pb.OpsDepthStoreCollectiveDocumentRequest{
		DepartureID:  params.DepartureID,
		DocumentType: params.DocumentType,
		URL:          params.URL,
		PilgrimID:    params.PilgrimID,
		Notes:        params.Notes,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &StoreCollectiveDocumentResult{DocumentID: resp.GetDocumentID()}, nil
}

type GetCollectiveDocumentsParams struct {
	DepartureID  string
	PilgrimID    string
	DocumentType string
}

type CollectiveDocRowResult struct {
	DocumentID   string
	DocumentType string
	URL          string
	PilgrimID    string
	UploadedAt   string
}

type GetCollectiveDocumentsResult struct {
	DepartureID string
	Documents   []CollectiveDocRowResult
}

func (a *Adapter) GetCollectiveDocuments(ctx context.Context, params *GetCollectiveDocumentsParams) (*GetCollectiveDocumentsResult, error) {
	const op = "ops_grpc_adapter.Adapter.GetCollectiveDocuments"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.GetCollectiveDocuments(ctx, &pb.OpsDepthGetCollectiveDocumentsRequest{
		DepartureID:  params.DepartureID,
		PilgrimID:    params.PilgrimID,
		DocumentType: params.DocumentType,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	docs := make([]CollectiveDocRowResult, 0, len(resp.GetDocuments()))
	for _, d := range resp.GetDocuments() {
		docs = append(docs, CollectiveDocRowResult{
			DocumentID:   d.GetDocumentID(),
			DocumentType: d.GetDocumentType(),
			URL:          d.GetURL(),
			PilgrimID:    d.GetPilgrimID(),
			UploadedAt:   d.GetUploadedAt(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetCollectiveDocumentsResult{DepartureID: resp.GetDepartureID(), Documents: docs}, nil
}

type SetDocumentACLParams struct {
	DocumentID     string
	AccessLevel    string
	AllowedUserIDs []string
}

type SetDocumentACLResult struct {
	Updated bool
}

func (a *Adapter) SetDocumentACL(ctx context.Context, params *SetDocumentACLParams) (*SetDocumentACLResult, error) {
	const op = "ops_grpc_adapter.Adapter.SetDocumentACL"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.SetDocumentACL(ctx, &pb.OpsDepthSetDocumentACLRequest{
		DocumentID:     params.DocumentID,
		AccessLevel:    params.AccessLevel,
		AllowedUserIDs: params.AllowedUserIDs,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &SetDocumentACLResult{Updated: resp.GetUpdated()}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-022: Passport OCR & mahram
// ---------------------------------------------------------------------------

type ExtractPassportOCRParams struct {
	ImageURL  string
	PilgrimID string
}

type PassportOCRDataResult struct {
	FullName       string
	PassportNumber string
	Nationality    string
	DateOfBirth    string
	ExpiryDate     string
	Gender         string
	Confidence     float64
}

type ExtractPassportOCRResult struct {
	PilgrimID string
	Data      *PassportOCRDataResult
	Warnings  []string
}

func (a *Adapter) ExtractPassportOCR(ctx context.Context, params *ExtractPassportOCRParams) (*ExtractPassportOCRResult, error) {
	const op = "ops_grpc_adapter.Adapter.ExtractPassportOCR"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.ExtractPassportOCR(ctx, &pb.OpsDepthExtractPassportOCRRequest{
		ImageURL:  params.ImageURL,
		PilgrimID: params.PilgrimID,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	var d *PassportOCRDataResult
	if resp.GetData() != nil {
		d = &PassportOCRDataResult{
			FullName:       resp.GetData().GetFullName(),
			PassportNumber: resp.GetData().GetPassportNumber(),
			Nationality:    resp.GetData().GetNationality(),
			DateOfBirth:    resp.GetData().GetDateOfBirth(),
			ExpiryDate:     resp.GetData().GetExpiryDate(),
			Gender:         resp.GetData().GetGender(),
			Confidence:     resp.GetData().GetConfidence(),
		}
	}
	span.SetStatus(codes.Ok, "ok")
	return &ExtractPassportOCRResult{PilgrimID: resp.GetPilgrimID(), Data: d, Warnings: resp.GetWarnings()}, nil
}

type SetMahramRelationParams struct {
	PilgrimID       string
	MahramPilgrimID string
	Relation        string
}

type SetMahramRelationResult struct {
	RelationID string
}

func (a *Adapter) SetMahramRelation(ctx context.Context, params *SetMahramRelationParams) (*SetMahramRelationResult, error) {
	const op = "ops_grpc_adapter.Adapter.SetMahramRelation"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.SetMahramRelation(ctx, &pb.OpsDepthSetMahramRelationRequest{
		PilgrimID:       params.PilgrimID,
		MahramPilgrimID: params.MahramPilgrimID,
		Relation:        params.Relation,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &SetMahramRelationResult{RelationID: resp.GetRelationID()}, nil
}

type GetMahramRelationsParams struct {
	BookingID string
}

type MahramRelationRowResult struct {
	RelationID      string
	PilgrimID       string
	MahramPilgrimID string
	Relation        string
}

type GetMahramRelationsResult struct {
	BookingID string
	Relations []MahramRelationRowResult
}

func (a *Adapter) GetMahramRelations(ctx context.Context, params *GetMahramRelationsParams) (*GetMahramRelationsResult, error) {
	const op = "ops_grpc_adapter.Adapter.GetMahramRelations"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.GetMahramRelations(ctx, &pb.OpsDepthGetMahramRelationsRequest{
		BookingID: params.BookingID,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	rows := make([]MahramRelationRowResult, 0, len(resp.GetRelations()))
	for _, r := range resp.GetRelations() {
		rows = append(rows, MahramRelationRowResult{
			RelationID:      r.GetRelationID(),
			PilgrimID:       r.GetPilgrimID(),
			MahramPilgrimID: r.GetMahramPilgrimID(),
			Relation:        r.GetRelation(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetMahramRelationsResult{BookingID: resp.GetBookingID(), Relations: rows}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-023: Document progress & expiry
// ---------------------------------------------------------------------------

type GetDocumentProgressParams struct {
	DepartureID string
}

type DocProgressRowResult struct {
	PilgrimID       string
	DocumentType    string
	Status          string
	ExpiryDate      string
	DaysUntilExpiry int32
}

type GetDocumentProgressResult struct {
	DepartureID       string
	Rows              []DocProgressRowResult
	TotalPilgrims     int32
	DocumentsComplete int32
	DocumentsExpiring int32
}

func (a *Adapter) GetDocumentProgress(ctx context.Context, params *GetDocumentProgressParams) (*GetDocumentProgressResult, error) {
	const op = "ops_grpc_adapter.Adapter.GetDocumentProgress"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.GetDocumentProgress(ctx, &pb.OpsDepthGetDocumentProgressRequest{
		DepartureID: params.DepartureID,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	rows := make([]DocProgressRowResult, 0, len(resp.GetRows()))
	for _, r := range resp.GetRows() {
		rows = append(rows, DocProgressRowResult{
			PilgrimID:       r.GetPilgrimID(),
			DocumentType:    r.GetDocumentType(),
			Status:          r.GetStatus(),
			ExpiryDate:      r.GetExpiryDate(),
			DaysUntilExpiry: r.GetDaysUntilExpiry(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetDocumentProgressResult{
		DepartureID:       resp.GetDepartureID(),
		Rows:              rows,
		TotalPilgrims:     resp.GetTotalPilgrims(),
		DocumentsComplete: resp.GetDocumentsComplete(),
		DocumentsExpiring: resp.GetDocumentsExpiring(),
	}, nil
}

type GetExpiryAlertsParams struct {
	ThresholdDays int32
	DepartureID   string
}

type ExpiryAlertRowResult struct {
	PilgrimID       string
	DocumentType    string
	ExpiryDate      string
	DaysUntilExpiry int32
}

type GetExpiryAlertsResult struct {
	Alerts []ExpiryAlertRowResult
}

func (a *Adapter) GetExpiryAlerts(ctx context.Context, params *GetExpiryAlertsParams) (*GetExpiryAlertsResult, error) {
	const op = "ops_grpc_adapter.Adapter.GetExpiryAlerts"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.GetExpiryAlerts(ctx, &pb.OpsDepthGetExpiryAlertsRequest{
		ThresholdDays: params.ThresholdDays,
		DepartureID:   params.DepartureID,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	alerts := make([]ExpiryAlertRowResult, 0, len(resp.GetAlerts()))
	for _, a := range resp.GetAlerts() {
		alerts = append(alerts, ExpiryAlertRowResult{
			PilgrimID:       a.GetPilgrimID(),
			DocumentType:    a.GetDocumentType(),
			ExpiryDate:      a.GetExpiryDate(),
			DaysUntilExpiry: a.GetDaysUntilExpiry(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetExpiryAlertsResult{Alerts: alerts}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-024: Official letter
// ---------------------------------------------------------------------------

type GenerateOfficialLetterParams struct {
	TemplateName string
	DepartureID  string
	PilgrimID    string
	IssuedTo     string
	Notes        string
}

type GenerateOfficialLetterResult struct {
	LetterID  string
	LetterURL string
	IssuedAt  string
}

func (a *Adapter) GenerateOfficialLetter(ctx context.Context, params *GenerateOfficialLetterParams) (*GenerateOfficialLetterResult, error) {
	const op = "ops_grpc_adapter.Adapter.GenerateOfficialLetter"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.GenerateOfficialLetter(ctx, &pb.OpsDepthGenerateOfficialLetterRequest{
		TemplateName: params.TemplateName,
		DepartureID:  params.DepartureID,
		PilgrimID:    params.PilgrimID,
		IssuedTo:     params.IssuedTo,
		Notes:        params.Notes,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &GenerateOfficialLetterResult{
		LetterID:  resp.GetLetterID(),
		LetterURL: resp.GetLetterURL(),
		IssuedAt:  resp.GetIssuedAt(),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-025: Immigration manifest
// ---------------------------------------------------------------------------

type GenerateImmigrationManifestParams struct {
	DepartureID string
	Format      string
}

type GenerateImmigrationManifestResult struct {
	ManifestID  string
	ManifestURL string
	Version     string
	RowCount    int32
}

func (a *Adapter) GenerateImmigrationManifest(ctx context.Context, params *GenerateImmigrationManifestParams) (*GenerateImmigrationManifestResult, error) {
	const op = "ops_grpc_adapter.Adapter.GenerateImmigrationManifest"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.GenerateImmigrationManifest(ctx, &pb.OpsDepthGenerateImmigrationManifestRequest{
		DepartureID: params.DepartureID,
		Format:      params.Format,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &GenerateImmigrationManifestResult{
		ManifestID:  resp.GetManifestID(),
		ManifestURL: resp.GetManifestURL(),
		Version:     resp.GetVersion(),
		RowCount:    resp.GetRowCount(),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-027: Transport arrangement
// ---------------------------------------------------------------------------

type AssignTransportParams struct {
	DepartureID string
	VehicleType string
	VehicleID   string
	PilgrimIDs  []string
}

type AssignTransportResult struct {
	AssignmentID  string
	AssignedCount int32
}

func (a *Adapter) AssignTransport(ctx context.Context, params *AssignTransportParams) (*AssignTransportResult, error) {
	const op = "ops_grpc_adapter.Adapter.AssignTransport"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.AssignTransport(ctx, &pb.OpsDepthAssignTransportRequest{
		DepartureID: params.DepartureID,
		VehicleType: params.VehicleType,
		VehicleID:   params.VehicleID,
		PilgrimIDs:  params.PilgrimIDs,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &AssignTransportResult{
		AssignmentID:  resp.GetAssignmentID(),
		AssignedCount: resp.GetAssignedCount(),
	}, nil
}

type GetTransportAssignmentsParams struct {
	DepartureID string
}

type TransportAssignmentRowResult struct {
	AssignmentID string
	VehicleType  string
	VehicleID    string
	PilgrimIDs   []string
}

type GetTransportAssignmentsResult struct {
	DepartureID string
	Assignments []TransportAssignmentRowResult
}

func (a *Adapter) GetTransportAssignments(ctx context.Context, params *GetTransportAssignmentsParams) (*GetTransportAssignmentsResult, error) {
	const op = "ops_grpc_adapter.Adapter.GetTransportAssignments"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.GetTransportAssignments(ctx, &pb.OpsDepthGetTransportAssignmentsRequest{
		DepartureID: params.DepartureID,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	rows := make([]TransportAssignmentRowResult, 0, len(resp.GetAssignments()))
	for _, r := range resp.GetAssignments() {
		rows = append(rows, TransportAssignmentRowResult{
			AssignmentID: r.GetAssignmentID(),
			VehicleType:  r.GetVehicleType(),
			VehicleID:    r.GetVehicleID(),
			PilgrimIDs:   r.GetPilgrimIDs(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetTransportAssignmentsResult{DepartureID: resp.GetDepartureID(), Assignments: rows}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-028: Manifest delta
// ---------------------------------------------------------------------------

type PublishManifestDeltaParams struct {
	DepartureID string
	ChangeType  string
	EntityID    string
	Notes       string
}

type PublishManifestDeltaResult struct {
	DeltaID     string
	PublishedAt string
}

func (a *Adapter) PublishManifestDelta(ctx context.Context, params *PublishManifestDeltaParams) (*PublishManifestDeltaResult, error) {
	const op = "ops_grpc_adapter.Adapter.PublishManifestDelta"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.PublishManifestDelta(ctx, &pb.OpsDepthPublishManifestDeltaRequest{
		DepartureID: params.DepartureID,
		ChangeType:  params.ChangeType,
		EntityID:    params.EntityID,
		Notes:       params.Notes,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &PublishManifestDeltaResult{DeltaID: resp.GetDeltaID(), PublishedAt: resp.GetPublishedAt()}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-029: Staff assignment
// ---------------------------------------------------------------------------

type AssignStaffParams struct {
	DepartureID string
	StaffUserID string
	Role        string
	PilgrimIDs  []string
}

type AssignStaffResult struct {
	AssignmentID string
}

func (a *Adapter) AssignStaff(ctx context.Context, params *AssignStaffParams) (*AssignStaffResult, error) {
	const op = "ops_grpc_adapter.Adapter.AssignStaff"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.AssignStaff(ctx, &pb.OpsDepthAssignStaffRequest{
		DepartureID: params.DepartureID,
		StaffUserID: params.StaffUserID,
		Role:        params.Role,
		PilgrimIDs:  params.PilgrimIDs,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &AssignStaffResult{AssignmentID: resp.GetAssignmentID()}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-030: Passport log
// ---------------------------------------------------------------------------

type RecordPassportHandoverParams struct {
	DepartureID string
	PilgrimID   string
	FromUserID  string
	ToUserID    string
	Notes       string
}

type RecordPassportHandoverResult struct {
	HandoverID string
	RecordedAt string
}

func (a *Adapter) RecordPassportHandover(ctx context.Context, params *RecordPassportHandoverParams) (*RecordPassportHandoverResult, error) {
	const op = "ops_grpc_adapter.Adapter.RecordPassportHandover"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.RecordPassportHandover(ctx, &pb.OpsDepthRecordPassportHandoverRequest{
		DepartureID: params.DepartureID,
		PilgrimID:   params.PilgrimID,
		FromUserID:  params.FromUserID,
		ToUserID:    params.ToUserID,
		Notes:       params.Notes,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RecordPassportHandoverResult{HandoverID: resp.GetHandoverID(), RecordedAt: resp.GetRecordedAt()}, nil
}

type GetPassportLogParams struct {
	DepartureID string
}

type PassportHandoverRowResult struct {
	HandoverID string
	PilgrimID  string
	FromUserID string
	ToUserID   string
	RecordedAt string
	Notes      string
}

type GetPassportLogResult struct {
	DepartureID string
	Rows        []PassportHandoverRowResult
}

func (a *Adapter) GetPassportLog(ctx context.Context, params *GetPassportLogParams) (*GetPassportLogResult, error) {
	const op = "ops_grpc_adapter.Adapter.GetPassportLog"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.GetPassportLog(ctx, &pb.OpsDepthGetPassportLogRequest{
		DepartureID: params.DepartureID,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	rows := make([]PassportHandoverRowResult, 0, len(resp.GetRows()))
	for _, r := range resp.GetRows() {
		rows = append(rows, PassportHandoverRowResult{
			HandoverID: r.GetHandoverID(),
			PilgrimID:  r.GetPilgrimID(),
			FromUserID: r.GetFromUserID(),
			ToUserID:   r.GetToUserID(),
			RecordedAt: r.GetRecordedAt(),
			Notes:      r.GetNotes(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetPassportLogResult{DepartureID: resp.GetDepartureID(), Rows: rows}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-031: Visa progress
// ---------------------------------------------------------------------------

type GetVisaProgressParams struct {
	DepartureID string
}

type VisaProgressRowResult struct {
	PilgrimID     string
	VisaStatus    string
	SubmittedAt   string
	ExpectedBy    string
	DaysRemaining int32
}

type GetVisaProgressResult struct {
	DepartureID string
	Rows        []VisaProgressRowResult
	Submitted   int32
	Approved    int32
	Rejected    int32
	Pending     int32
}

func (a *Adapter) GetVisaProgress(ctx context.Context, params *GetVisaProgressParams) (*GetVisaProgressResult, error) {
	const op = "ops_grpc_adapter.Adapter.GetVisaProgress"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.GetVisaProgress(ctx, &pb.OpsDepthGetVisaProgressRequest{
		DepartureID: params.DepartureID,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	rows := make([]VisaProgressRowResult, 0, len(resp.GetRows()))
	for _, r := range resp.GetRows() {
		rows = append(rows, VisaProgressRowResult{
			PilgrimID:     r.GetPilgrimID(),
			VisaStatus:    r.GetVisaStatus(),
			SubmittedAt:   r.GetSubmittedAt(),
			ExpectedBy:    r.GetExpectedBy(),
			DaysRemaining: r.GetDaysRemaining(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetVisaProgressResult{
		DepartureID: resp.GetDepartureID(),
		Rows:        rows,
		Submitted:   resp.GetSubmitted(),
		Approved:    resp.GetApproved(),
		Rejected:    resp.GetRejected(),
		Pending:     resp.GetPending(),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-032: E-visa repository
// ---------------------------------------------------------------------------

type StoreEVisaParams struct {
	PilgrimID   string
	DepartureID string
	VisaNumber  string
	VisaURL     string
	IssuedDate  string
	ExpiryDate  string
}

type StoreEVisaResult struct {
	EVisaID string
}

func (a *Adapter) StoreEVisa(ctx context.Context, params *StoreEVisaParams) (*StoreEVisaResult, error) {
	const op = "ops_grpc_adapter.Adapter.StoreEVisa"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.StoreEVisa(ctx, &pb.OpsDepthStoreEVisaRequest{
		PilgrimID:   params.PilgrimID,
		DepartureID: params.DepartureID,
		VisaNumber:  params.VisaNumber,
		VisaURL:     params.VisaURL,
		IssuedDate:  params.IssuedDate,
		ExpiryDate:  params.ExpiryDate,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &StoreEVisaResult{EVisaID: resp.GetEVisaID()}, nil
}

type GetEVisaParams struct {
	PilgrimID   string
	DepartureID string
}

type EVisaResult struct {
	EVisaID    string
	PilgrimID  string
	VisaNumber string
	VisaURL    string
	IssuedDate string
	ExpiryDate string
}

type GetEVisaResult struct {
	EVisa *EVisaResult
}

func (a *Adapter) GetEVisa(ctx context.Context, params *GetEVisaParams) (*GetEVisaResult, error) {
	const op = "ops_grpc_adapter.Adapter.GetEVisa"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.GetEVisa(ctx, &pb.OpsDepthGetEVisaRequest{
		PilgrimID:   params.PilgrimID,
		DepartureID: params.DepartureID,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	var ev *EVisaResult
	if resp.GetEVisa() != nil {
		ev = &EVisaResult{
			EVisaID:    resp.GetEVisa().GetEVisaID(),
			PilgrimID:  resp.GetEVisa().GetPilgrimID(),
			VisaNumber: resp.GetEVisa().GetVisaNumber(),
			VisaURL:    resp.GetEVisa().GetVisaURL(),
			IssuedDate: resp.GetEVisa().GetIssuedDate(),
			ExpiryDate: resp.GetEVisa().GetExpiryDate(),
		}
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetEVisaResult{EVisa: ev}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-033: External provider
// ---------------------------------------------------------------------------

type TriggerExternalProviderParams struct {
	Provider    string
	Action      string
	ReferenceID string
	Payload     string
}

type TriggerExternalProviderResult struct {
	RequestID string
	Status    string
}

func (a *Adapter) TriggerExternalProvider(ctx context.Context, params *TriggerExternalProviderParams) (*TriggerExternalProviderResult, error) {
	const op = "ops_grpc_adapter.Adapter.TriggerExternalProvider"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.TriggerExternalProvider(ctx, &pb.OpsDepthTriggerExternalProviderRequest{
		Provider:    params.Provider,
		Action:      params.Action,
		ReferenceID: params.ReferenceID,
		Payload:     params.Payload,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &TriggerExternalProviderResult{RequestID: resp.GetRequestID(), Status: resp.GetStatus()}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-034: Refund & penalty
// ---------------------------------------------------------------------------

type CreateRefundParams struct {
	BookingID string
	Reason    string
	Amount    int64
	Notes     string
}

type CreateRefundResult struct {
	RefundID string
	Status   string
}

func (a *Adapter) CreateRefund(ctx context.Context, params *CreateRefundParams) (*CreateRefundResult, error) {
	const op = "ops_grpc_adapter.Adapter.CreateRefund"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.CreateRefund(ctx, &pb.OpsDepthCreateRefundRequest{
		BookingID: params.BookingID,
		Reason:    params.Reason,
		Amount:    params.Amount,
		Notes:     params.Notes,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &CreateRefundResult{RefundID: resp.GetRefundID(), Status: resp.GetStatus()}, nil
}

type ApproveRefundParams struct {
	RefundID string
	Decision string
	Notes    string
}

type ApproveRefundResult struct {
	RefundID string
	Status   string
}

func (a *Adapter) ApproveRefund(ctx context.Context, params *ApproveRefundParams) (*ApproveRefundResult, error) {
	const op = "ops_grpc_adapter.Adapter.ApproveRefund"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.ApproveRefund(ctx, &pb.OpsDepthApproveRefundRequest{
		RefundID: params.RefundID,
		Decision: params.Decision,
		Notes:    params.Notes,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &ApproveRefundResult{RefundID: resp.GetRefundID(), Status: resp.GetStatus()}, nil
}

type RecordPenaltyParams struct {
	BookingID   string
	PenaltyType string
	Amount      int64
	Notes       string
}

type RecordPenaltyResult struct {
	PenaltyID string
}

func (a *Adapter) RecordPenalty(ctx context.Context, params *RecordPenaltyParams) (*RecordPenaltyResult, error) {
	const op = "ops_grpc_adapter.Adapter.RecordPenalty"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.RecordPenalty(ctx, &pb.OpsDepthRecordPenaltyRequest{
		BookingID:   params.BookingID,
		PenaltyType: params.PenaltyType,
		Amount:      params.Amount,
		Notes:       params.Notes,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RecordPenaltyResult{PenaltyID: resp.GetPenaltyID()}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-036: Luggage counter
// ---------------------------------------------------------------------------

type RecordLuggageScanParams struct {
	DepartureID string
	PilgrimID   string
	TagID       string
	ScanPoint   string
}

type RecordLuggageScanResult struct {
	ScanID    string
	TotalBags int32
}

func (a *Adapter) RecordLuggageScan(ctx context.Context, params *RecordLuggageScanParams) (*RecordLuggageScanResult, error) {
	const op = "ops_grpc_adapter.Adapter.RecordLuggageScan"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.RecordLuggageScan(ctx, &pb.OpsDepthRecordLuggageScanRequest{
		DepartureID: params.DepartureID,
		PilgrimID:   params.PilgrimID,
		TagID:       params.TagID,
		ScanPoint:   params.ScanPoint,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RecordLuggageScanResult{ScanID: resp.GetScanID(), TotalBags: resp.GetTotalBags()}, nil
}

type GetLuggageCountParams struct {
	DepartureID string
}

type LuggageCountRowResult struct {
	PilgrimID     string
	BagCount      int32
	LastScannedAt string
}

type GetLuggageCountResult struct {
	DepartureID string
	TotalBags   int32
	Rows        []LuggageCountRowResult
}

func (a *Adapter) GetLuggageCount(ctx context.Context, params *GetLuggageCountParams) (*GetLuggageCountResult, error) {
	const op = "ops_grpc_adapter.Adapter.GetLuggageCount"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.GetLuggageCount(ctx, &pb.OpsDepthGetLuggageCountRequest{
		DepartureID: params.DepartureID,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	rows := make([]LuggageCountRowResult, 0, len(resp.GetRows()))
	for _, r := range resp.GetRows() {
		rows = append(rows, LuggageCountRowResult{
			PilgrimID:     r.GetPilgrimID(),
			BagCount:      r.GetBagCount(),
			LastScannedAt: r.GetLastScannedAt(),
		})
	}
	span.SetStatus(codes.Ok, "ok")
	return &GetLuggageCountResult{
		DepartureID: resp.GetDepartureID(),
		TotalBags:   resp.GetTotalBags(),
		Rows:        rows,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-037: Departure/arrival broadcast
// ---------------------------------------------------------------------------

type BroadcastScheduleParams struct {
	DepartureID   string
	BroadcastType string
	Message       string
	Channel       string
}

type BroadcastScheduleResult struct {
	BroadcastID    string
	RecipientCount int32
}

func (a *Adapter) BroadcastSchedule(ctx context.Context, params *BroadcastScheduleParams) (*BroadcastScheduleResult, error) {
	const op = "ops_grpc_adapter.Adapter.BroadcastSchedule"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.BroadcastSchedule(ctx, &pb.OpsDepthBroadcastScheduleRequest{
		DepartureID:   params.DepartureID,
		BroadcastType: params.BroadcastType,
		Message:       params.Message,
		Channel:       params.Channel,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &BroadcastScheduleResult{
		BroadcastID:    resp.GetBroadcastID(),
		RecipientCount: resp.GetRecipientCount(),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-039: Raudhah shield & tasreh
// ---------------------------------------------------------------------------

type IssueDigitalTasrehParams struct {
	PilgrimID   string
	DepartureID string
	VisitDate   string
}

type IssueDigitalTasrehResult struct {
	TasrehID string
	QRCode   string
}

func (a *Adapter) IssueDigitalTasreh(ctx context.Context, params *IssueDigitalTasrehParams) (*IssueDigitalTasrehResult, error) {
	const op = "ops_grpc_adapter.Adapter.IssueDigitalTasreh"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.IssueDigitalTasreh(ctx, &pb.OpsDepthIssueDigitalTasrehRequest{
		PilgrimID:   params.PilgrimID,
		DepartureID: params.DepartureID,
		VisitDate:   params.VisitDate,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &IssueDigitalTasrehResult{TasrehID: resp.GetTasrehID(), QRCode: resp.GetQRCode()}, nil
}

type RecordRaudhahEntryParams struct {
	TasrehID  string
	PilgrimID string
	EntryTime string
}

type RecordRaudhahEntryResult struct {
	RecordID string
	Valid    bool
}

func (a *Adapter) RecordRaudhahEntry(ctx context.Context, params *RecordRaudhahEntryParams) (*RecordRaudhahEntryResult, error) {
	const op = "ops_grpc_adapter.Adapter.RecordRaudhahEntry"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.RecordRaudhahEntry(ctx, &pb.OpsDepthRecordRaudhahEntryRequest{
		TasrehID:  params.TasrehID,
		PilgrimID: params.PilgrimID,
		EntryTime: params.EntryTime,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RecordRaudhahEntryResult{RecordID: resp.GetRecordID(), Valid: resp.GetValid()}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-040: Audio devices
// ---------------------------------------------------------------------------

type RegisterAudioDeviceParams struct {
	DepartureID  string
	DeviceType   string
	SerialNumber string
	AssignedTo   string
}

type RegisterAudioDeviceResult struct {
	DeviceID string
}

func (a *Adapter) RegisterAudioDevice(ctx context.Context, params *RegisterAudioDeviceParams) (*RegisterAudioDeviceResult, error) {
	const op = "ops_grpc_adapter.Adapter.RegisterAudioDevice"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.RegisterAudioDevice(ctx, &pb.OpsDepthRegisterAudioDeviceRequest{
		DepartureID:  params.DepartureID,
		DeviceType:   params.DeviceType,
		SerialNumber: params.SerialNumber,
		AssignedTo:   params.AssignedTo,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RegisterAudioDeviceResult{DeviceID: resp.GetDeviceID()}, nil
}

type UpdateAudioDeviceStatusParams struct {
	DeviceID string
	Status   string
	Notes    string
}

type UpdateAudioDeviceStatusResult struct {
	DeviceID string
	Status   string
}

func (a *Adapter) UpdateAudioDeviceStatus(ctx context.Context, params *UpdateAudioDeviceStatusParams) (*UpdateAudioDeviceStatusResult, error) {
	const op = "ops_grpc_adapter.Adapter.UpdateAudioDeviceStatus"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.UpdateAudioDeviceStatus(ctx, &pb.OpsDepthUpdateAudioDeviceStatusRequest{
		DeviceID: params.DeviceID,
		Status:   params.Status,
		Notes:    params.Notes,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &UpdateAudioDeviceStatusResult{DeviceID: resp.GetDeviceID(), Status: resp.GetStatus()}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-041: Zamzam distribution
// ---------------------------------------------------------------------------

type RecordZamzamDistributionParams struct {
	DepartureID string
	PilgrimID   string
	LitersGiven float64
	ReceivedBy  string
}

type RecordZamzamDistributionResult struct {
	DistributionID string
	RecordedAt     string
}

func (a *Adapter) RecordZamzamDistribution(ctx context.Context, params *RecordZamzamDistributionParams) (*RecordZamzamDistributionResult, error) {
	const op = "ops_grpc_adapter.Adapter.RecordZamzamDistribution"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.RecordZamzamDistribution(ctx, &pb.OpsDepthRecordZamzamDistributionRequest{
		DepartureID: params.DepartureID,
		PilgrimID:   params.PilgrimID,
		LitersGiven: params.LitersGiven,
		ReceivedBy:  params.ReceivedBy,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RecordZamzamDistributionResult{
		DistributionID: resp.GetDistributionID(),
		RecordedAt:     resp.GetRecordedAt(),
	}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-042: Room check-in
// ---------------------------------------------------------------------------

type RecordRoomCheckInParams struct {
	DepartureID string
	PilgrimID   string
	RoomNumber  string
	HotelID     string
}

type RecordRoomCheckInResult struct {
	CheckInID   string
	CheckedInAt string
}

func (a *Adapter) RecordRoomCheckIn(ctx context.Context, params *RecordRoomCheckInParams) (*RecordRoomCheckInResult, error) {
	const op = "ops_grpc_adapter.Adapter.RecordRoomCheckIn"
	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, a.logger).Info().Str("op", op).Msg("")

	resp, err := a.depthClient.RecordRoomCheckIn(ctx, &pb.OpsDepthRecordRoomCheckInRequest{
		DepartureID: params.DepartureID,
		PilgrimID:   params.PilgrimID,
		RoomNumber:  params.RoomNumber,
		HotelID:     params.HotelID,
	})
	if err != nil {
		wrapped := mapOpsDepthError(err)
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}
	span.SetStatus(codes.Ok, "ok")
	return &RecordRoomCheckInResult{
		CheckInID:   resp.GetCheckInID(),
		CheckedInAt: resp.GetCheckedInAt(),
	}, nil
}
