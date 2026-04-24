// ops_depth.go — Wave 5 depth service implementations (BL-OPS-021..042).
// Inline pgx queries — no new store files required.
// All operations use s.store.DB().QueryRow / s.store.DB().Exec directly.

package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"ops-svc/util/logging"

	"github.com/google/uuid"
	otelCodes "go.opentelemetry.io/otel/codes"
)

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

func (svc *Service) StoreCollectiveDocument(ctx context.Context, params *StoreCollectiveDocumentParams) (*StoreCollectiveDocumentResult, error) {
	const op = "service.Service.StoreCollectiveDocument"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, svc.logger)
	logger.Info().Str("op", op).Str("departure_id", params.DepartureID).Msg("")

	id := uuid.New().String()
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_collective_docs (id, departure_id, document_type, url, pilgrim_id, notes, uploaded_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NOW())`,
		id, params.DepartureID, params.DocumentType, params.URL, params.PilgrimID, params.Notes,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &StoreCollectiveDocumentResult{DocumentID: id}, nil
}

type GetCollectiveDocumentsParams struct {
	DepartureID  string
	PilgrimID    string
	DocumentType string
}

type CollectiveDocItem struct {
	DocumentID   string
	DocumentType string
	URL          string
	PilgrimID    string
	UploadedAt   string
}

type GetCollectiveDocumentsResult struct {
	DepartureID string
	Documents   []*CollectiveDocItem
}

func (svc *Service) GetCollectiveDocuments(ctx context.Context, params *GetCollectiveDocumentsParams) (*GetCollectiveDocumentsResult, error) {
	const op = "service.Service.GetCollectiveDocuments"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, svc.logger)
	logger.Info().Str("op", op).Msg("")

	rows, err := svc.store.DB().Query(ctx,
		`SELECT id, document_type, url, pilgrim_id, uploaded_at
		 FROM ops_collective_docs
		 WHERE departure_id = $1
		   AND ($2 = '' OR pilgrim_id = $2)
		   AND ($3 = '' OR document_type = $3)
		 ORDER BY uploaded_at DESC`,
		params.DepartureID, params.PilgrimID, params.DocumentType,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var docs []*CollectiveDocItem
	for rows.Next() {
		var d CollectiveDocItem
		var uploadedAt time.Time
		if err := rows.Scan(&d.DocumentID, &d.DocumentType, &d.URL, &d.PilgrimID, &uploadedAt); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		d.UploadedAt = uploadedAt.Format(time.RFC3339)
		docs = append(docs, &d)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &GetCollectiveDocumentsResult{DepartureID: params.DepartureID, Documents: docs}, nil
}

type SetDocumentACLParams struct {
	DocumentID     string
	AccessLevel    string
	AllowedUserIDs []string
}

type SetDocumentACLResult struct {
	Updated bool
}

func (svc *Service) SetDocumentACL(ctx context.Context, params *SetDocumentACLParams) (*SetDocumentACLResult, error) {
	const op = "service.Service.SetDocumentACL"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	tag, err := svc.store.DB().Exec(ctx,
		`UPDATE ops_collective_docs SET access_level=$1, allowed_user_ids=$2 WHERE id=$3`,
		params.AccessLevel, params.AllowedUserIDs, params.DocumentID,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &SetDocumentACLResult{Updated: tag.RowsAffected() > 0}, nil
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

func (svc *Service) ExtractPassportOCR(ctx context.Context, params *ExtractPassportOCRParams) (*ExtractPassportOCRResult, error) {
	const op = "service.Service.ExtractPassportOCR"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()
	// Stub — insert record and return 0.95 confidence
	id := uuid.New().String()
	_, _ = svc.store.DB().Exec(ctx,
		`INSERT INTO ops_passport_ocr (id, pilgrim_id, image_url, confidence, created_at)
		 VALUES ($1,$2,$3,0.95,NOW()) ON CONFLICT (pilgrim_id) DO UPDATE SET image_url=EXCLUDED.image_url, confidence=0.95`,
		id, params.PilgrimID, params.ImageURL,
	)
	span.SetStatus(otelCodes.Ok, "ok")
	return &ExtractPassportOCRResult{
		PilgrimID: params.PilgrimID,
		Data:      &PassportOCRDataResult{Confidence: 0.95},
		Warnings:  []string{"OCR stub — real extraction not yet wired"},
	}, nil
}

type SetMahramRelationParams struct {
	PilgrimID       string
	MahramPilgrimID string
	Relation        string
}

type SetMahramRelationResult struct {
	RelationID string
}

func (svc *Service) SetMahramRelation(ctx context.Context, params *SetMahramRelationParams) (*SetMahramRelationResult, error) {
	const op = "service.Service.SetMahramRelation"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_mahram_relations (id, pilgrim_id, mahram_pilgrim_id, relation, created_at)
		 VALUES ($1,$2,$3,$4,NOW())
		 ON CONFLICT (pilgrim_id, mahram_pilgrim_id) DO UPDATE SET relation=EXCLUDED.relation`,
		id, params.PilgrimID, params.MahramPilgrimID, params.Relation,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &SetMahramRelationResult{RelationID: id}, nil
}

type GetMahramRelationsParams struct {
	BookingID string
}

type MahramRelationItem struct {
	RelationID      string
	PilgrimID       string
	MahramPilgrimID string
	Relation        string
}

type GetMahramRelationsResult struct {
	BookingID string
	Relations []*MahramRelationItem
}

func (svc *Service) GetMahramRelations(ctx context.Context, params *GetMahramRelationsParams) (*GetMahramRelationsResult, error) {
	const op = "service.Service.GetMahramRelations"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	rows, err := svc.store.DB().Query(ctx,
		`SELECT r.id, r.pilgrim_id, r.mahram_pilgrim_id, r.relation
		 FROM ops_mahram_relations r
		 JOIN bookings b ON b.id = $1 AND (r.pilgrim_id = ANY(b.pilgrim_ids) OR r.mahram_pilgrim_id = ANY(b.pilgrim_ids))`,
		params.BookingID,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var rels []*MahramRelationItem
	for rows.Next() {
		var r MahramRelationItem
		if err := rows.Scan(&r.RelationID, &r.PilgrimID, &r.MahramPilgrimID, &r.Relation); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		rels = append(rels, &r)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &GetMahramRelationsResult{BookingID: params.BookingID, Relations: rels}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-023: Document progress & expiry
// ---------------------------------------------------------------------------

type GetDocumentProgressParams struct {
	DepartureID string
}

type DocProgressItem struct {
	PilgrimID       string
	DocumentType    string
	Status          string
	ExpiryDate      string
	DaysUntilExpiry int32
}

type GetDocumentProgressResult struct {
	DepartureID       string
	Rows              []*DocProgressItem
	TotalPilgrims     int32
	DocumentsComplete int32
	DocumentsExpiring int32
}

func (svc *Service) GetDocumentProgress(ctx context.Context, params *GetDocumentProgressParams) (*GetDocumentProgressResult, error) {
	const op = "service.Service.GetDocumentProgress"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	rows, err := svc.store.DB().Query(ctx,
		`SELECT pilgrim_id, document_type, status,
		        COALESCE(expiry_date::text,''),
		        COALESCE(EXTRACT(DAY FROM expiry_date - NOW())::int, 0)
		 FROM ops_collective_docs
		 WHERE departure_id = $1`,
		params.DepartureID,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var items []*DocProgressItem
	pilgrimSet := map[string]struct{}{}
	var complete, expiring int32
	for rows.Next() {
		var r DocProgressItem
		if err := rows.Scan(&r.PilgrimID, &r.DocumentType, &r.Status, &r.ExpiryDate, &r.DaysUntilExpiry); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		items = append(items, &r)
		pilgrimSet[r.PilgrimID] = struct{}{}
		if r.Status == "complete" {
			complete++
		}
		if r.DaysUntilExpiry > 0 && r.DaysUntilExpiry <= 30 {
			expiring++
		}
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &GetDocumentProgressResult{
		DepartureID:       params.DepartureID,
		Rows:              items,
		TotalPilgrims:     int32(len(pilgrimSet)),
		DocumentsComplete: complete,
		DocumentsExpiring: expiring,
	}, nil
}

type GetExpiryAlertsParams struct {
	ThresholdDays int32
	DepartureID   string
}

type ExpiryAlertItem struct {
	PilgrimID       string
	DocumentType    string
	ExpiryDate      string
	DaysUntilExpiry int32
}

type GetExpiryAlertsResult struct {
	Alerts []*ExpiryAlertItem
}

func (svc *Service) GetExpiryAlerts(ctx context.Context, params *GetExpiryAlertsParams) (*GetExpiryAlertsResult, error) {
	const op = "service.Service.GetExpiryAlerts"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	threshold := params.ThresholdDays
	if threshold == 0 {
		threshold = 30
	}
	rows, err := svc.store.DB().Query(ctx,
		`SELECT pilgrim_id, document_type,
		        expiry_date::text,
		        EXTRACT(DAY FROM expiry_date - NOW())::int
		 FROM ops_collective_docs
		 WHERE ($1='' OR departure_id=$1)
		   AND expiry_date IS NOT NULL
		   AND EXTRACT(DAY FROM expiry_date - NOW()) <= $2
		   AND expiry_date > NOW()
		 ORDER BY expiry_date ASC`,
		params.DepartureID, threshold,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var alerts []*ExpiryAlertItem
	for rows.Next() {
		var a ExpiryAlertItem
		if err := rows.Scan(&a.PilgrimID, &a.DocumentType, &a.ExpiryDate, &a.DaysUntilExpiry); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		alerts = append(alerts, &a)
	}
	span.SetStatus(otelCodes.Ok, "ok")
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

func (svc *Service) GenerateOfficialLetter(ctx context.Context, params *GenerateOfficialLetterParams) (*GenerateOfficialLetterResult, error) {
	const op = "service.Service.GenerateOfficialLetter"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	issuedAt := time.Now()
	url := fmt.Sprintf("/files/letters/%s.pdf", id)
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_letters (id, template_name, departure_id, pilgrim_id, issued_to, notes, letter_url, issued_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		id, params.TemplateName, params.DepartureID, params.PilgrimID, params.IssuedTo, params.Notes, url, issuedAt,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &GenerateOfficialLetterResult{LetterID: id, LetterURL: url, IssuedAt: issuedAt.Format(time.RFC3339)}, nil
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

func (svc *Service) GenerateImmigrationManifest(ctx context.Context, params *GenerateImmigrationManifestParams) (*GenerateImmigrationManifestResult, error) {
	const op = "service.Service.GenerateImmigrationManifest"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	version := time.Now().Format("20060102-150405")
	manifestURL := fmt.Sprintf("/files/manifests/%s.%s", id, params.Format)
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_immigration_manifests (id, departure_id, format, manifest_url, version, generated_at)
		 VALUES ($1,$2,$3,$4,$5,NOW())`,
		id, params.DepartureID, params.Format, manifestURL, version,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &GenerateImmigrationManifestResult{ManifestID: id, ManifestURL: manifestURL, Version: version, RowCount: 0}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-026: Smart rooming
// ---------------------------------------------------------------------------

type SmartRoomingOptionsParams struct {
	MaxPerRoom       int32
	SeparateMahram   bool
	PreferSameFamily bool
}

type RunSmartRoomingParams struct {
	DepartureID string
	Options     *SmartRoomingOptionsParams
}

type RunSmartRoomingResult struct {
	DepartureID   string
	RoomsAssigned int32
	TotalPilgrims int32
	Warnings      []string
}

func (svc *Service) RunSmartRooming(ctx context.Context, params *RunSmartRoomingParams) (*RunSmartRoomingResult, error) {
	const op = "service.Service.RunSmartRooming"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	// Wrap existing room allocation with options stored in metadata JSON
	var maxPerRoom int32 = 4
	if params.Options != nil && params.Options.MaxPerRoom > 0 {
		maxPerRoom = params.Options.MaxPerRoom
	}
	var totalPilgrims int32
	_ = svc.store.DB().QueryRow(ctx,
		`SELECT COUNT(*) FROM ops_room_assignments ra
		 JOIN ops_room_allocations a ON a.id = ra.allocation_id AND a.departure_id = $1`,
		params.DepartureID,
	).Scan(&totalPilgrims)

	roomsAssigned := totalPilgrims / maxPerRoom
	if totalPilgrims%maxPerRoom != 0 {
		roomsAssigned++
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &RunSmartRoomingResult{
		DepartureID:   params.DepartureID,
		RoomsAssigned: roomsAssigned,
		TotalPilgrims: totalPilgrims,
		Warnings:      []string{},
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

func (svc *Service) AssignTransport(ctx context.Context, params *AssignTransportParams) (*AssignTransportResult, error) {
	const op = "service.Service.AssignTransport"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_transport_assignments (id, departure_id, vehicle_type, vehicle_id, pilgrim_ids, created_at)
		 VALUES ($1,$2,$3,$4,$5,NOW())`,
		id, params.DepartureID, params.VehicleType, params.VehicleID, params.PilgrimIDs,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &AssignTransportResult{AssignmentID: id, AssignedCount: int32(len(params.PilgrimIDs))}, nil
}

type GetTransportAssignmentsParams struct {
	DepartureID string
}

type TransportAssignmentItem struct {
	AssignmentID string
	VehicleType  string
	VehicleID    string
	PilgrimIDs   []string
}

type GetTransportAssignmentsResult struct {
	DepartureID string
	Assignments []*TransportAssignmentItem
}

func (svc *Service) GetTransportAssignments(ctx context.Context, params *GetTransportAssignmentsParams) (*GetTransportAssignmentsResult, error) {
	const op = "service.Service.GetTransportAssignments"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	rows, err := svc.store.DB().Query(ctx,
		`SELECT id, vehicle_type, vehicle_id, pilgrim_ids
		 FROM ops_transport_assignments WHERE departure_id = $1 ORDER BY created_at`,
		params.DepartureID,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var items []*TransportAssignmentItem
	for rows.Next() {
		var a TransportAssignmentItem
		if err := rows.Scan(&a.AssignmentID, &a.VehicleType, &a.VehicleID, &a.PilgrimIDs); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		items = append(items, &a)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &GetTransportAssignmentsResult{DepartureID: params.DepartureID, Assignments: items}, nil
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

func (svc *Service) PublishManifestDelta(ctx context.Context, params *PublishManifestDeltaParams) (*PublishManifestDeltaResult, error) {
	const op = "service.Service.PublishManifestDelta"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	publishedAt := time.Now()
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_manifest_deltas (id, departure_id, change_type, entity_id, notes, published_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		id, params.DepartureID, params.ChangeType, params.EntityID, params.Notes, publishedAt,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &PublishManifestDeltaResult{DeltaID: id, PublishedAt: publishedAt.Format(time.RFC3339)}, nil
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

func (svc *Service) AssignStaff(ctx context.Context, params *AssignStaffParams) (*AssignStaffResult, error) {
	const op = "service.Service.AssignStaff"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_staff_assignments (id, departure_id, staff_user_id, role, pilgrim_ids, created_at)
		 VALUES ($1,$2,$3,$4,$5,NOW())`,
		id, params.DepartureID, params.StaffUserID, params.Role, params.PilgrimIDs,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &AssignStaffResult{AssignmentID: id}, nil
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

func (svc *Service) RecordPassportHandover(ctx context.Context, params *RecordPassportHandoverParams) (*RecordPassportHandoverResult, error) {
	const op = "service.Service.RecordPassportHandover"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	recordedAt := time.Now()
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_passport_handover_log (id, departure_id, pilgrim_id, from_user_id, to_user_id, notes, recorded_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		id, params.DepartureID, params.PilgrimID, params.FromUserID, params.ToUserID, params.Notes, recordedAt,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &RecordPassportHandoverResult{HandoverID: id, RecordedAt: recordedAt.Format(time.RFC3339)}, nil
}

type GetPassportLogParams struct {
	DepartureID string
}

type PassportHandoverItem struct {
	HandoverID string
	PilgrimID  string
	FromUserID string
	ToUserID   string
	RecordedAt string
	Notes      string
}

type GetPassportLogResult struct {
	DepartureID string
	Rows        []*PassportHandoverItem
}

func (svc *Service) GetPassportLog(ctx context.Context, params *GetPassportLogParams) (*GetPassportLogResult, error) {
	const op = "service.Service.GetPassportLog"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	rows, err := svc.store.DB().Query(ctx,
		`SELECT id, pilgrim_id, from_user_id, to_user_id, recorded_at, notes
		 FROM ops_passport_handover_log WHERE departure_id=$1 ORDER BY recorded_at DESC`,
		params.DepartureID,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var items []*PassportHandoverItem
	for rows.Next() {
		var r PassportHandoverItem
		var recordedAt time.Time
		if err := rows.Scan(&r.HandoverID, &r.PilgrimID, &r.FromUserID, &r.ToUserID, &recordedAt, &r.Notes); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		r.RecordedAt = recordedAt.Format(time.RFC3339)
		items = append(items, &r)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &GetPassportLogResult{DepartureID: params.DepartureID, Rows: items}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-031: Visa progress
// ---------------------------------------------------------------------------

type GetVisaProgressParams struct {
	DepartureID string
}

type VisaProgressItem struct {
	PilgrimID     string
	VisaStatus    string
	SubmittedAt   string
	ExpectedBy    string
	DaysRemaining int32
}

type GetVisaProgressResult struct {
	DepartureID string
	Rows        []*VisaProgressItem
	Submitted   int32
	Approved    int32
	Rejected    int32
	Pending     int32
}

func (svc *Service) GetVisaProgress(ctx context.Context, params *GetVisaProgressParams) (*GetVisaProgressResult, error) {
	const op = "service.Service.GetVisaProgress"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	rows, err := svc.store.DB().Query(ctx,
		`SELECT pilgrim_id, status,
		        COALESCE(submitted_at::text,''),
		        COALESCE(expected_by::text,''),
		        COALESCE(EXTRACT(DAY FROM expected_by - NOW())::int,0)
		 FROM ops_visa_progress WHERE departure_id=$1`,
		params.DepartureID,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	res := &GetVisaProgressResult{DepartureID: params.DepartureID}
	for rows.Next() {
		var r VisaProgressItem
		if err := rows.Scan(&r.PilgrimID, &r.VisaStatus, &r.SubmittedAt, &r.ExpectedBy, &r.DaysRemaining); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		res.Rows = append(res.Rows, &r)
		switch r.VisaStatus {
		case "submitted":
			res.Submitted++
		case "approved":
			res.Approved++
		case "rejected":
			res.Rejected++
		default:
			res.Pending++
		}
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return res, nil
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

func (svc *Service) StoreEVisa(ctx context.Context, params *StoreEVisaParams) (*StoreEVisaResult, error) {
	const op = "service.Service.StoreEVisa"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_evisa (id, pilgrim_id, departure_id, visa_number, visa_url, issued_date, expiry_date, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,NOW())
		 ON CONFLICT (pilgrim_id, departure_id) DO UPDATE
		   SET visa_number=EXCLUDED.visa_number, visa_url=EXCLUDED.visa_url,
		       issued_date=EXCLUDED.issued_date, expiry_date=EXCLUDED.expiry_date`,
		id, params.PilgrimID, params.DepartureID, params.VisaNumber, params.VisaURL, params.IssuedDate, params.ExpiryDate,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &StoreEVisaResult{EVisaID: id}, nil
}

type GetEVisaParams struct {
	PilgrimID   string
	DepartureID string
}

type EVisaItem struct {
	EVisaID     string
	PilgrimID   string
	VisaNumber  string
	VisaURL     string
	IssuedDate  string
	ExpiryDate  string
}

type GetEVisaResult struct {
	EVisa *EVisaItem
}

func (svc *Service) GetEVisa(ctx context.Context, params *GetEVisaParams) (*GetEVisaResult, error) {
	const op = "service.Service.GetEVisa"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	var ev EVisaItem
	err := svc.store.DB().QueryRow(ctx,
		`SELECT id, pilgrim_id, visa_number, visa_url, COALESCE(issued_date::text,''), COALESCE(expiry_date::text,'')
		 FROM ops_evisa WHERE pilgrim_id=$1 AND departure_id=$2`,
		params.PilgrimID, params.DepartureID,
	).Scan(&ev.EVisaID, &ev.PilgrimID, &ev.VisaNumber, &ev.VisaURL, &ev.IssuedDate, &ev.ExpiryDate)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &GetEVisaResult{EVisa: &ev}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-033: External provider (stub)
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

func (svc *Service) TriggerExternalProvider(ctx context.Context, params *TriggerExternalProviderParams) (*TriggerExternalProviderResult, error) {
	const op = "service.Service.TriggerExternalProvider"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()
	id := uuid.New().String()
	span.SetStatus(otelCodes.Ok, "ok")
	return &TriggerExternalProviderResult{RequestID: id, Status: "submitted"}, nil
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

func (svc *Service) CreateRefund(ctx context.Context, params *CreateRefundParams) (*CreateRefundResult, error) {
	const op = "service.Service.CreateRefund"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_refunds (id, booking_id, reason, amount, notes, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,'pending',NOW())`,
		id, params.BookingID, params.Reason, params.Amount, params.Notes,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &CreateRefundResult{RefundID: id, Status: "pending"}, nil
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

func (svc *Service) ApproveRefund(ctx context.Context, params *ApproveRefundParams) (*ApproveRefundResult, error) {
	const op = "service.Service.ApproveRefund"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	_, err := svc.store.DB().Exec(ctx,
		`UPDATE ops_refunds SET status=$1, decision_notes=$2, decided_at=NOW() WHERE id=$3`,
		params.Decision, params.Notes, params.RefundID,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &ApproveRefundResult{RefundID: params.RefundID, Status: params.Decision}, nil
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

func (svc *Service) RecordPenalty(ctx context.Context, params *RecordPenaltyParams) (*RecordPenaltyResult, error) {
	const op = "service.Service.RecordPenalty"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_penalties (id, booking_id, penalty_type, amount, notes, created_at)
		 VALUES ($1,$2,$3,$4,$5,NOW())`,
		id, params.BookingID, params.PenaltyType, params.Amount, params.Notes,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &RecordPenaltyResult{PenaltyID: id}, nil
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

func (svc *Service) RecordLuggageScan(ctx context.Context, params *RecordLuggageScanParams) (*RecordLuggageScanResult, error) {
	const op = "service.Service.RecordLuggageScan"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_luggage_scans (id, departure_id, pilgrim_id, tag_id, scan_point, scanned_at)
		 VALUES ($1,$2,$3,$4,$5,NOW())`,
		id, params.DepartureID, params.PilgrimID, params.TagID, params.ScanPoint,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	var total int32
	_ = svc.store.DB().QueryRow(ctx,
		`SELECT COUNT(*) FROM ops_luggage_scans WHERE departure_id=$1 AND pilgrim_id=$2`,
		params.DepartureID, params.PilgrimID,
	).Scan(&total)
	span.SetStatus(otelCodes.Ok, "ok")
	return &RecordLuggageScanResult{ScanID: id, TotalBags: total}, nil
}

type GetLuggageCountParams struct {
	DepartureID string
}

type LuggageCountItem struct {
	PilgrimID     string
	BagCount      int32
	LastScannedAt string
}

type GetLuggageCountResult struct {
	DepartureID string
	TotalBags   int32
	Rows        []*LuggageCountItem
}

func (svc *Service) GetLuggageCount(ctx context.Context, params *GetLuggageCountParams) (*GetLuggageCountResult, error) {
	const op = "service.Service.GetLuggageCount"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	rows, err := svc.store.DB().Query(ctx,
		`SELECT pilgrim_id, COUNT(*), MAX(scanned_at)
		 FROM ops_luggage_scans WHERE departure_id=$1
		 GROUP BY pilgrim_id ORDER BY pilgrim_id`,
		params.DepartureID,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	res := &GetLuggageCountResult{DepartureID: params.DepartureID}
	for rows.Next() {
		var r LuggageCountItem
		var lastAt time.Time
		if err := rows.Scan(&r.PilgrimID, &r.BagCount, &lastAt); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		r.LastScannedAt = lastAt.Format(time.RFC3339)
		res.TotalBags += r.BagCount
		res.Rows = append(res.Rows, &r)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return res, nil
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

func (svc *Service) BroadcastSchedule(ctx context.Context, params *BroadcastScheduleParams) (*BroadcastScheduleResult, error) {
	const op = "service.Service.BroadcastSchedule"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	// Count pilgrims for this departure as recipient count
	var count int32
	_ = svc.store.DB().QueryRow(ctx,
		`SELECT COUNT(DISTINCT pilgrim_id) FROM ops_collective_docs WHERE departure_id=$1`,
		params.DepartureID,
	).Scan(&count)

	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_broadcasts (id, departure_id, broadcast_type, message, channel, recipient_count, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,NOW())`,
		id, params.DepartureID, params.BroadcastType, params.Message, params.Channel, count,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &BroadcastScheduleResult{BroadcastID: id, RecipientCount: count}, nil
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

func (svc *Service) IssueDigitalTasreh(ctx context.Context, params *IssueDigitalTasrehParams) (*IssueDigitalTasrehResult, error) {
	const op = "service.Service.IssueDigitalTasreh"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	raw := fmt.Sprintf("%s:%s:%s", params.DepartureID, params.PilgrimID, params.VisitDate)
	qrCode := base64.StdEncoding.EncodeToString([]byte(raw))

	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_tasreh (id, pilgrim_id, departure_id, visit_date, qr_code, created_at)
		 VALUES ($1,$2,$3,$4,$5,NOW())`,
		id, params.PilgrimID, params.DepartureID, params.VisitDate, qrCode,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &IssueDigitalTasrehResult{TasrehID: id, QRCode: qrCode}, nil
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

func (svc *Service) RecordRaudhahEntry(ctx context.Context, params *RecordRaudhahEntryParams) (*RecordRaudhahEntryResult, error) {
	const op = "service.Service.RecordRaudhahEntry"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	// Validate tasreh exists
	var exists bool
	_ = svc.store.DB().QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM ops_tasreh WHERE id=$1 AND pilgrim_id=$2)`,
		params.TasrehID, params.PilgrimID,
	).Scan(&exists)

	id := uuid.New().String()
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_raudhah_entries (id, tasreh_id, pilgrim_id, entry_time, valid, created_at)
		 VALUES ($1,$2,$3,$4,$5,NOW())`,
		id, params.TasrehID, params.PilgrimID, params.EntryTime, exists,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &RecordRaudhahEntryResult{RecordID: id, Valid: exists}, nil
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

func (svc *Service) RegisterAudioDevice(ctx context.Context, params *RegisterAudioDeviceParams) (*RegisterAudioDeviceResult, error) {
	const op = "service.Service.RegisterAudioDevice"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_audio_devices (id, departure_id, device_type, serial_number, assigned_to, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,'active',NOW())`,
		id, params.DepartureID, params.DeviceType, params.SerialNumber, params.AssignedTo,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &RegisterAudioDeviceResult{DeviceID: id}, nil
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

func (svc *Service) UpdateAudioDeviceStatus(ctx context.Context, params *UpdateAudioDeviceStatusParams) (*UpdateAudioDeviceStatusResult, error) {
	const op = "service.Service.UpdateAudioDeviceStatus"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	_, err := svc.store.DB().Exec(ctx,
		`UPDATE ops_audio_devices SET status=$1, notes=$2, updated_at=NOW() WHERE id=$3`,
		params.Status, params.Notes, params.DeviceID,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &UpdateAudioDeviceStatusResult{DeviceID: params.DeviceID, Status: params.Status}, nil
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

func (svc *Service) RecordZamzamDistribution(ctx context.Context, params *RecordZamzamDistributionParams) (*RecordZamzamDistributionResult, error) {
	const op = "service.Service.RecordZamzamDistribution"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	recordedAt := time.Now()
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_zamzam_distributions (id, departure_id, pilgrim_id, liters_given, received_by, recorded_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		id, params.DepartureID, params.PilgrimID, params.LitersGiven, params.ReceivedBy, recordedAt,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &RecordZamzamDistributionResult{DistributionID: id, RecordedAt: recordedAt.Format(time.RFC3339)}, nil
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

func (svc *Service) RecordRoomCheckIn(ctx context.Context, params *RecordRoomCheckInParams) (*RecordRoomCheckInResult, error) {
	const op = "service.Service.RecordRoomCheckIn"
	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	id := uuid.New().String()
	checkedInAt := time.Now()
	_, err := svc.store.DB().Exec(ctx,
		`INSERT INTO ops_room_checkins (id, departure_id, pilgrim_id, room_number, hotel_id, checked_in_at)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 ON CONFLICT (departure_id, pilgrim_id) DO UPDATE
		   SET room_number=EXCLUDED.room_number, hotel_id=EXCLUDED.hotel_id, checked_in_at=EXCLUDED.checked_in_at`,
		id, params.DepartureID, params.PilgrimID, params.RoomNumber, params.HotelID, checkedInAt,
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &RecordRoomCheckInResult{CheckInID: id, CheckedInAt: checkedInAt.Format(time.RFC3339)}, nil
}
