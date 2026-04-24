// proxy_ops_depth.go — gateway REST handlers for ops-svc Wave 5 depth RPCs.
// BL-OPS-021..042

package rest_oapi

import (
	"errors"
	"strconv"

	"gateway-svc/adapter/ops_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// BL-OPS-021: StoreCollectiveDocument — POST /v1/ops/collective-docs
// ---------------------------------------------------------------------------

func (s *Server) StoreCollectiveDocument(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.StoreCollectiveDocument"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		DepartureID  string `json:"departure_id"`
		DocumentType string `json:"document_type"`
		URL          string `json:"url"`
		PilgrimID    string `json:"pilgrim_id"`
		Notes        string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.StoreCollectiveDocument(ctx, &ops_grpc_adapter.StoreCollectiveDocumentParams{
		DepartureID:  body.DepartureID,
		DocumentType: body.DocumentType,
		URL:          body.URL,
		PilgrimID:    body.PilgrimID,
		Notes:        body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"document_id": result.DocumentID}})
}

// ---------------------------------------------------------------------------
// BL-OPS-021: GetCollectiveDocuments — GET /v1/ops/collective-docs/:departure_id
// ---------------------------------------------------------------------------

func (s *Server) GetCollectiveDocuments(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetCollectiveDocuments"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")
	departureID := c.Params("departure_id")

	result, err := s.svc.GetCollectiveDocuments(ctx, &ops_grpc_adapter.GetCollectiveDocumentsParams{
		DepartureID:  departureID,
		PilgrimID:    c.Query("pilgrim_id"),
		DocumentType: c.Query("document_type"),
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	type docJSON struct {
		DocumentID   string `json:"document_id"`
		DocumentType string `json:"document_type"`
		URL          string `json:"url"`
		PilgrimID    string `json:"pilgrim_id"`
		UploadedAt   string `json:"uploaded_at"`
	}
	docs := make([]docJSON, 0, len(result.Documents))
	for _, d := range result.Documents {
		docs = append(docs, docJSON{d.DocumentID, d.DocumentType, d.URL, d.PilgrimID, d.UploadedAt})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"departure_id": result.DepartureID, "documents": docs}})
}

// ---------------------------------------------------------------------------
// BL-OPS-021: SetDocumentACL — PUT /v1/ops/collective-docs/:id/acl
// ---------------------------------------------------------------------------

func (s *Server) SetDocumentACL(c *fiber.Ctx, documentID string) error {
	const op = "rest_oapi.Server.SetDocumentACL"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		AccessLevel    string   `json:"access_level"`
		AllowedUserIDs []string `json:"allowed_user_ids"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.SetDocumentACL(ctx, &ops_grpc_adapter.SetDocumentACLParams{
		DocumentID:     documentID,
		AccessLevel:    body.AccessLevel,
		AllowedUserIDs: body.AllowedUserIDs,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"updated": result.Updated}})
}

// ---------------------------------------------------------------------------
// BL-OPS-022: ExtractPassportOCR — POST /v1/ops/passport-ocr
// ---------------------------------------------------------------------------

func (s *Server) ExtractPassportOCR(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ExtractPassportOCR"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		ImageURL  string `json:"image_url"`
		PilgrimID string `json:"pilgrim_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.ExtractPassportOCR(ctx, &ops_grpc_adapter.ExtractPassportOCRParams{
		ImageURL:  body.ImageURL,
		PilgrimID: body.PilgrimID,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	var dataJSON interface{}
	if result.Data != nil {
		dataJSON = fiber.Map{
			"full_name":       result.Data.FullName,
			"passport_number": result.Data.PassportNumber,
			"nationality":     result.Data.Nationality,
			"date_of_birth":   result.Data.DateOfBirth,
			"expiry_date":     result.Data.ExpiryDate,
			"gender":          result.Data.Gender,
			"confidence":      result.Data.Confidence,
		}
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"pilgrim_id": result.PilgrimID, "data": dataJSON, "warnings": result.Warnings}})
}

// ---------------------------------------------------------------------------
// BL-OPS-022: SetMahramRelation — POST /v1/ops/mahram-relations
// ---------------------------------------------------------------------------

func (s *Server) SetMahramRelation(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.SetMahramRelation"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		PilgrimID       string `json:"pilgrim_id"`
		MahramPilgrimID string `json:"mahram_pilgrim_id"`
		Relation        string `json:"relation"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.SetMahramRelation(ctx, &ops_grpc_adapter.SetMahramRelationParams{
		PilgrimID:       body.PilgrimID,
		MahramPilgrimID: body.MahramPilgrimID,
		Relation:        body.Relation,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"relation_id": result.RelationID}})
}

// ---------------------------------------------------------------------------
// BL-OPS-022: GetMahramRelations — GET /v1/ops/mahram-relations/:booking_id
// ---------------------------------------------------------------------------

func (s *Server) GetMahramRelations(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetMahramRelations"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")
	bookingID := c.Params("booking_id")

	result, err := s.svc.GetMahramRelations(ctx, &ops_grpc_adapter.GetMahramRelationsParams{
		BookingID: bookingID,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	type relJSON struct {
		RelationID      string `json:"relation_id"`
		PilgrimID       string `json:"pilgrim_id"`
		MahramPilgrimID string `json:"mahram_pilgrim_id"`
		Relation        string `json:"relation"`
	}
	rels := make([]relJSON, 0, len(result.Relations))
	for _, r := range result.Relations {
		rels = append(rels, relJSON{r.RelationID, r.PilgrimID, r.MahramPilgrimID, r.Relation})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"booking_id": result.BookingID, "relations": rels}})
}

// ---------------------------------------------------------------------------
// BL-OPS-023: GetDocumentProgress — GET /v1/ops/document-progress/:departure_id
// ---------------------------------------------------------------------------

func (s *Server) GetDocumentProgress(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetDocumentProgress"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")
	departureID := c.Params("departure_id")

	result, err := s.svc.GetDocumentProgress(ctx, &ops_grpc_adapter.GetDocumentProgressParams{
		DepartureID: departureID,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	type rowJSON struct {
		PilgrimID       string `json:"pilgrim_id"`
		DocumentType    string `json:"document_type"`
		Status          string `json:"status"`
		ExpiryDate      string `json:"expiry_date"`
		DaysUntilExpiry int32  `json:"days_until_expiry"`
	}
	rows := make([]rowJSON, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, rowJSON{r.PilgrimID, r.DocumentType, r.Status, r.ExpiryDate, r.DaysUntilExpiry})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{
		"departure_id":       result.DepartureID,
		"rows":               rows,
		"total_pilgrims":     result.TotalPilgrims,
		"documents_complete": result.DocumentsComplete,
		"documents_expiring": result.DocumentsExpiring,
	}})
}

// ---------------------------------------------------------------------------
// BL-OPS-023: GetExpiryAlerts — GET /v1/ops/expiry-alerts
// ---------------------------------------------------------------------------

func (s *Server) GetExpiryAlerts(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetExpiryAlerts"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var threshold int32
	if v := c.Query("threshold_days"); v != "" {
		t, err := strconv.Atoi(v)
		if err != nil || t < 0 {
			return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, errors.New("threshold_days must be a non-negative integer")))
		}
		threshold = int32(t)
	}
	result, err := s.svc.GetExpiryAlerts(ctx, &ops_grpc_adapter.GetExpiryAlertsParams{
		ThresholdDays: threshold,
		DepartureID:   c.Query("departure_id"),
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	type alertJSON struct {
		PilgrimID       string `json:"pilgrim_id"`
		DocumentType    string `json:"document_type"`
		ExpiryDate      string `json:"expiry_date"`
		DaysUntilExpiry int32  `json:"days_until_expiry"`
	}
	alerts := make([]alertJSON, 0, len(result.Alerts))
	for _, a := range result.Alerts {
		alerts = append(alerts, alertJSON{a.PilgrimID, a.DocumentType, a.ExpiryDate, a.DaysUntilExpiry})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"alerts": alerts}})
}

// ---------------------------------------------------------------------------
// BL-OPS-024: GenerateOfficialLetter — POST /v1/ops/official-letters
// ---------------------------------------------------------------------------

func (s *Server) GenerateOfficialLetter(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GenerateOfficialLetter"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		TemplateName string `json:"template_name"`
		DepartureID  string `json:"departure_id"`
		PilgrimID    string `json:"pilgrim_id"`
		IssuedTo     string `json:"issued_to"`
		Notes        string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.GenerateOfficialLetter(ctx, &ops_grpc_adapter.GenerateOfficialLetterParams{
		TemplateName: body.TemplateName,
		DepartureID:  body.DepartureID,
		PilgrimID:    body.PilgrimID,
		IssuedTo:     body.IssuedTo,
		Notes:        body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{
		"letter_id":  result.LetterID,
		"letter_url": result.LetterURL,
		"issued_at":  result.IssuedAt,
	}})
}

// ---------------------------------------------------------------------------
// BL-OPS-025: GenerateImmigrationManifest — POST /v1/ops/immigration-manifest/:departure_id
// ---------------------------------------------------------------------------

func (s *Server) GenerateImmigrationManifest(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GenerateImmigrationManifest"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")
	departureID := c.Params("departure_id")

	var body struct {
		Format string `json:"format"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.GenerateImmigrationManifest(ctx, &ops_grpc_adapter.GenerateImmigrationManifestParams{
		DepartureID: departureID,
		Format:      body.Format,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{
		"manifest_id":  result.ManifestID,
		"manifest_url": result.ManifestURL,
		"version":      result.Version,
		"row_count":    result.RowCount,
	}})
}

// ---------------------------------------------------------------------------
// BL-OPS-027: AssignTransport — POST /v1/ops/transport-assignments
// ---------------------------------------------------------------------------

func (s *Server) AssignTransport(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.AssignTransport"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		DepartureID string   `json:"departure_id"`
		VehicleType string   `json:"vehicle_type"`
		VehicleID   string   `json:"vehicle_id"`
		PilgrimIDs  []string `json:"pilgrim_ids"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.AssignTransport(ctx, &ops_grpc_adapter.AssignTransportParams{
		DepartureID: body.DepartureID,
		VehicleType: body.VehicleType,
		VehicleID:   body.VehicleID,
		PilgrimIDs:  body.PilgrimIDs,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{
		"assignment_id":  result.AssignmentID,
		"assigned_count": result.AssignedCount,
	}})
}

// ---------------------------------------------------------------------------
// BL-OPS-027: GetTransportAssignments — GET /v1/ops/transport-assignments/:departure_id
// ---------------------------------------------------------------------------

func (s *Server) GetTransportAssignments(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetTransportAssignments"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")
	departureID := c.Params("departure_id")

	result, err := s.svc.GetTransportAssignments(ctx, &ops_grpc_adapter.GetTransportAssignmentsParams{
		DepartureID: departureID,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	type rowJSON struct {
		AssignmentID string   `json:"assignment_id"`
		VehicleType  string   `json:"vehicle_type"`
		VehicleID    string   `json:"vehicle_id"`
		PilgrimIDs   []string `json:"pilgrim_ids"`
	}
	rows := make([]rowJSON, 0, len(result.Assignments))
	for _, r := range result.Assignments {
		rows = append(rows, rowJSON{r.AssignmentID, r.VehicleType, r.VehicleID, r.PilgrimIDs})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"departure_id": result.DepartureID, "assignments": rows}})
}

// ---------------------------------------------------------------------------
// BL-OPS-028: PublishManifestDelta — POST /v1/ops/manifest-delta
// ---------------------------------------------------------------------------

func (s *Server) PublishManifestDelta(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.PublishManifestDelta"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		DepartureID string `json:"departure_id"`
		ChangeType  string `json:"change_type"`
		EntityID    string `json:"entity_id"`
		Notes       string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.PublishManifestDelta(ctx, &ops_grpc_adapter.PublishManifestDeltaParams{
		DepartureID: body.DepartureID,
		ChangeType:  body.ChangeType,
		EntityID:    body.EntityID,
		Notes:       body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"delta_id": result.DeltaID, "published_at": result.PublishedAt}})
}

// ---------------------------------------------------------------------------
// BL-OPS-029: AssignStaff — POST /v1/ops/staff-assignments
// ---------------------------------------------------------------------------

func (s *Server) AssignStaff(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.AssignStaff"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		DepartureID string   `json:"departure_id"`
		StaffUserID string   `json:"staff_user_id"`
		Role        string   `json:"role"`
		PilgrimIDs  []string `json:"pilgrim_ids"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.AssignStaff(ctx, &ops_grpc_adapter.AssignStaffParams{
		DepartureID: body.DepartureID,
		StaffUserID: body.StaffUserID,
		Role:        body.Role,
		PilgrimIDs:  body.PilgrimIDs,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"assignment_id": result.AssignmentID}})
}

// ---------------------------------------------------------------------------
// BL-OPS-030: RecordPassportHandover — POST /v1/ops/passport-log
// ---------------------------------------------------------------------------

func (s *Server) RecordPassportHandover(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordPassportHandover"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		DepartureID string `json:"departure_id"`
		PilgrimID   string `json:"pilgrim_id"`
		FromUserID  string `json:"from_user_id"`
		ToUserID    string `json:"to_user_id"`
		Notes       string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.RecordPassportHandover(ctx, &ops_grpc_adapter.RecordPassportHandoverParams{
		DepartureID: body.DepartureID,
		PilgrimID:   body.PilgrimID,
		FromUserID:  body.FromUserID,
		ToUserID:    body.ToUserID,
		Notes:       body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"handover_id": result.HandoverID, "recorded_at": result.RecordedAt}})
}

// ---------------------------------------------------------------------------
// BL-OPS-030: GetPassportLog — GET /v1/ops/passport-log/:departure_id
// ---------------------------------------------------------------------------

func (s *Server) GetPassportLog(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetPassportLog"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")
	departureID := c.Params("departure_id")

	result, err := s.svc.GetPassportLog(ctx, &ops_grpc_adapter.GetPassportLogParams{
		DepartureID: departureID,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	type rowJSON struct {
		HandoverID string `json:"handover_id"`
		PilgrimID  string `json:"pilgrim_id"`
		FromUserID string `json:"from_user_id"`
		ToUserID   string `json:"to_user_id"`
		RecordedAt string `json:"recorded_at"`
		Notes      string `json:"notes"`
	}
	rows := make([]rowJSON, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, rowJSON{r.HandoverID, r.PilgrimID, r.FromUserID, r.ToUserID, r.RecordedAt, r.Notes})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"departure_id": result.DepartureID, "rows": rows}})
}

// ---------------------------------------------------------------------------
// BL-OPS-031: GetVisaProgress — GET /v1/ops/visa-progress/:departure_id
// ---------------------------------------------------------------------------

func (s *Server) GetVisaProgress(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetVisaProgress"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")
	departureID := c.Params("departure_id")

	result, err := s.svc.GetVisaProgress(ctx, &ops_grpc_adapter.GetVisaProgressParams{
		DepartureID: departureID,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	type rowJSON struct {
		PilgrimID     string `json:"pilgrim_id"`
		VisaStatus    string `json:"visa_status"`
		SubmittedAt   string `json:"submitted_at"`
		ExpectedBy    string `json:"expected_by"`
		DaysRemaining int32  `json:"days_remaining"`
	}
	rows := make([]rowJSON, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, rowJSON{r.PilgrimID, r.VisaStatus, r.SubmittedAt, r.ExpectedBy, r.DaysRemaining})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{
		"departure_id": result.DepartureID,
		"rows":         rows,
		"submitted":    result.Submitted,
		"approved":     result.Approved,
		"rejected":     result.Rejected,
		"pending":      result.Pending,
	}})
}

// ---------------------------------------------------------------------------
// BL-OPS-032: StoreEVisa — POST /v1/ops/evisa
// ---------------------------------------------------------------------------

func (s *Server) StoreEVisa(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.StoreEVisa"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		PilgrimID   string `json:"pilgrim_id"`
		DepartureID string `json:"departure_id"`
		VisaNumber  string `json:"visa_number"`
		VisaURL     string `json:"visa_url"`
		IssuedDate  string `json:"issued_date"`
		ExpiryDate  string `json:"expiry_date"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.StoreEVisa(ctx, &ops_grpc_adapter.StoreEVisaParams{
		PilgrimID:   body.PilgrimID,
		DepartureID: body.DepartureID,
		VisaNumber:  body.VisaNumber,
		VisaURL:     body.VisaURL,
		IssuedDate:  body.IssuedDate,
		ExpiryDate:  body.ExpiryDate,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"evisa_id": result.EVisaID}})
}

// ---------------------------------------------------------------------------
// BL-OPS-032: GetEVisa — GET /v1/ops/evisa/:pilgrim_id
// ---------------------------------------------------------------------------

func (s *Server) GetEVisa(c *fiber.Ctx, pilgrimID string) error {
	const op = "rest_oapi.Server.GetEVisa"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetEVisa(ctx, &ops_grpc_adapter.GetEVisaParams{
		PilgrimID:   pilgrimID,
		DepartureID: c.Query("departure_id"),
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	var evJSON interface{}
	if result.EVisa != nil {
		evJSON = fiber.Map{
			"evisa_id":    result.EVisa.EVisaID,
			"pilgrim_id":  result.EVisa.PilgrimID,
			"visa_number": result.EVisa.VisaNumber,
			"visa_url":    result.EVisa.VisaURL,
			"issued_date": result.EVisa.IssuedDate,
			"expiry_date": result.EVisa.ExpiryDate,
		}
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"evisa": evJSON}})
}

// ---------------------------------------------------------------------------
// BL-OPS-033: TriggerExternalProvider — POST /v1/ops/external-provider
// ---------------------------------------------------------------------------

func (s *Server) TriggerExternalProvider(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.TriggerExternalProvider"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		Provider    string `json:"provider"`
		Action      string `json:"action"`
		ReferenceID string `json:"reference_id"`
		Payload     string `json:"payload"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.TriggerExternalProvider(ctx, &ops_grpc_adapter.TriggerExternalProviderParams{
		Provider:    body.Provider,
		Action:      body.Action,
		ReferenceID: body.ReferenceID,
		Payload:     body.Payload,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"request_id": result.RequestID, "status": result.Status}})
}

// ---------------------------------------------------------------------------
// BL-OPS-034: CreateRefund — POST /v1/ops/refunds
// ---------------------------------------------------------------------------

func (s *Server) CreateRefund(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateRefund"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		BookingID string `json:"booking_id"`
		Reason    string `json:"reason"`
		Amount    int64  `json:"amount"`
		Notes     string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.CreateRefund(ctx, &ops_grpc_adapter.CreateRefundParams{
		BookingID: body.BookingID,
		Reason:    body.Reason,
		Amount:    body.Amount,
		Notes:     body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"refund_id": result.RefundID, "status": result.Status}})
}

// ---------------------------------------------------------------------------
// BL-OPS-034: ApproveRefund — PUT /v1/ops/refunds/:id/decision
// ---------------------------------------------------------------------------

func (s *Server) ApproveRefund(c *fiber.Ctx, refundID string) error {
	const op = "rest_oapi.Server.ApproveRefund"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		Decision string `json:"decision"`
		Notes    string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.ApproveRefund(ctx, &ops_grpc_adapter.ApproveRefundParams{
		RefundID: refundID,
		Decision: body.Decision,
		Notes:    body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"refund_id": result.RefundID, "status": result.Status}})
}

// ---------------------------------------------------------------------------
// BL-OPS-034: RecordPenalty — POST /v1/ops/penalties
// ---------------------------------------------------------------------------

func (s *Server) RecordPenalty(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordPenalty"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		BookingID   string `json:"booking_id"`
		PenaltyType string `json:"penalty_type"`
		Amount      int64  `json:"amount"`
		Notes       string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.RecordPenalty(ctx, &ops_grpc_adapter.RecordPenaltyParams{
		BookingID:   body.BookingID,
		PenaltyType: body.PenaltyType,
		Amount:      body.Amount,
		Notes:       body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"penalty_id": result.PenaltyID}})
}

// ---------------------------------------------------------------------------
// BL-OPS-036: RecordLuggageScan — POST /v1/ops/luggage-scans
// ---------------------------------------------------------------------------

func (s *Server) RecordLuggageScan(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordLuggageScan"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		DepartureID string `json:"departure_id"`
		PilgrimID   string `json:"pilgrim_id"`
		TagID       string `json:"tag_id"`
		ScanPoint   string `json:"scan_point"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.RecordLuggageScan(ctx, &ops_grpc_adapter.RecordLuggageScanParams{
		DepartureID: body.DepartureID,
		PilgrimID:   body.PilgrimID,
		TagID:       body.TagID,
		ScanPoint:   body.ScanPoint,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"scan_id": result.ScanID, "total_bags": result.TotalBags}})
}

// ---------------------------------------------------------------------------
// BL-OPS-036: GetLuggageCount — GET /v1/ops/luggage-count/:departure_id
// ---------------------------------------------------------------------------

func (s *Server) GetLuggageCount(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetLuggageCount"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")
	departureID := c.Params("departure_id")

	result, err := s.svc.GetLuggageCount(ctx, &ops_grpc_adapter.GetLuggageCountParams{
		DepartureID: departureID,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	type rowJSON struct {
		PilgrimID     string `json:"pilgrim_id"`
		BagCount      int32  `json:"bag_count"`
		LastScannedAt string `json:"last_scanned_at"`
	}
	rows := make([]rowJSON, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, rowJSON{r.PilgrimID, r.BagCount, r.LastScannedAt})
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{
		"departure_id": result.DepartureID,
		"total_bags":   result.TotalBags,
		"rows":         rows,
	}})
}

// ---------------------------------------------------------------------------
// BL-OPS-037: BroadcastSchedule — POST /v1/ops/departure-broadcast
// ---------------------------------------------------------------------------

func (s *Server) BroadcastSchedule(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.BroadcastSchedule"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		DepartureID   string `json:"departure_id"`
		BroadcastType string `json:"broadcast_type"`
		Message       string `json:"message"`
		Channel       string `json:"channel"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.BroadcastSchedule(ctx, &ops_grpc_adapter.BroadcastScheduleParams{
		DepartureID:   body.DepartureID,
		BroadcastType: body.BroadcastType,
		Message:       body.Message,
		Channel:       body.Channel,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{
		"broadcast_id":    result.BroadcastID,
		"recipient_count": result.RecipientCount,
	}})
}

// ---------------------------------------------------------------------------
// BL-OPS-039: IssueDigitalTasreh — POST /v1/ops/tasreh
// ---------------------------------------------------------------------------

func (s *Server) IssueDigitalTasreh(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.IssueDigitalTasreh"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		PilgrimID   string `json:"pilgrim_id"`
		DepartureID string `json:"departure_id"`
		VisitDate   string `json:"visit_date"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.IssueDigitalTasreh(ctx, &ops_grpc_adapter.IssueDigitalTasrehParams{
		PilgrimID:   body.PilgrimID,
		DepartureID: body.DepartureID,
		VisitDate:   body.VisitDate,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"tasreh_id": result.TasrehID, "qr_code": result.QRCode}})
}

// ---------------------------------------------------------------------------
// BL-OPS-039: RecordRaudhahEntry — POST /v1/ops/raudhah-entry
// ---------------------------------------------------------------------------

func (s *Server) RecordRaudhahEntry(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordRaudhahEntry"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		TasrehID  string `json:"tasreh_id"`
		PilgrimID string `json:"pilgrim_id"`
		EntryTime string `json:"entry_time"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.RecordRaudhahEntry(ctx, &ops_grpc_adapter.RecordRaudhahEntryParams{
		TasrehID:  body.TasrehID,
		PilgrimID: body.PilgrimID,
		EntryTime: body.EntryTime,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"record_id": result.RecordID, "valid": result.Valid}})
}

// ---------------------------------------------------------------------------
// BL-OPS-040: RegisterAudioDevice — POST /v1/ops/audio-devices
// ---------------------------------------------------------------------------

func (s *Server) RegisterAudioDevice(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RegisterAudioDevice"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		DepartureID  string `json:"departure_id"`
		DeviceType   string `json:"device_type"`
		SerialNumber string `json:"serial_number"`
		AssignedTo   string `json:"assigned_to"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.RegisterAudioDevice(ctx, &ops_grpc_adapter.RegisterAudioDeviceParams{
		DepartureID:  body.DepartureID,
		DeviceType:   body.DeviceType,
		SerialNumber: body.SerialNumber,
		AssignedTo:   body.AssignedTo,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"device_id": result.DeviceID}})
}

// ---------------------------------------------------------------------------
// BL-OPS-040: UpdateAudioDeviceStatus — PUT /v1/ops/audio-devices/:id/status
// ---------------------------------------------------------------------------

func (s *Server) UpdateAudioDeviceStatus(c *fiber.Ctx, deviceID string) error {
	const op = "rest_oapi.Server.UpdateAudioDeviceStatus"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.UpdateAudioDeviceStatus(ctx, &ops_grpc_adapter.UpdateAudioDeviceStatusParams{
		DeviceID: deviceID,
		Status:   body.Status,
		Notes:    body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"device_id": result.DeviceID, "status": result.Status}})
}

// ---------------------------------------------------------------------------
// BL-OPS-041: RecordZamzamDistribution — POST /v1/ops/zamzam
// ---------------------------------------------------------------------------

func (s *Server) RecordZamzamDistribution(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordZamzamDistribution"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		DepartureID string  `json:"departure_id"`
		PilgrimID   string  `json:"pilgrim_id"`
		LitersGiven float64 `json:"liters_given"`
		ReceivedBy  string  `json:"received_by"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.RecordZamzamDistribution(ctx, &ops_grpc_adapter.RecordZamzamDistributionParams{
		DepartureID: body.DepartureID,
		PilgrimID:   body.PilgrimID,
		LitersGiven: body.LitersGiven,
		ReceivedBy:  body.ReceivedBy,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{
		"distribution_id": result.DistributionID,
		"recorded_at":     result.RecordedAt,
	}})
}

// ---------------------------------------------------------------------------
// BL-OPS-042: RecordRoomCheckIn — POST /v1/ops/room-checkin
// ---------------------------------------------------------------------------

func (s *Server) RecordRoomCheckIn(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordRoomCheckIn"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		DepartureID string `json:"departure_id"`
		PilgrimID   string `json:"pilgrim_id"`
		RoomNumber  string `json:"room_number"`
		HotelID     string `json:"hotel_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.RecordRoomCheckIn(ctx, &ops_grpc_adapter.RecordRoomCheckInParams{
		DepartureID: body.DepartureID,
		PilgrimID:   body.PilgrimID,
		RoomNumber:  body.RoomNumber,
		HotelID:     body.HotelID,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{
		"checkin_id":   result.CheckInID,
		"checked_in_at": result.CheckedInAt,
	}})
}
