// ops_depth.go — gRPC handlers for Wave 5 depth RPCs (BL-OPS-021..042).
// Thin delegation layer: pb request → service params → pb response.

package grpc_api

import (
	"context"

	"ops-svc/api/grpc_api/pb"
	"ops-svc/service"
	"ops-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// ---------------------------------------------------------------------------
// BL-OPS-021: Collective documents
// ---------------------------------------------------------------------------

func (s *Server) StoreCollectiveDocument(ctx context.Context, req *pb.StoreCollectiveDocumentRequest) (*pb.StoreCollectiveDocumentResponse, error) {
	const op = "grpc_api.Server.StoreCollectiveDocument"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Str("departure_id", req.GetDepartureID()).Msg("")

	result, err := s.svc.StoreCollectiveDocument(ctx, &service.StoreCollectiveDocumentParams{
		DepartureID:  req.GetDepartureID(),
		DocumentType: req.GetDocumentType(),
		URL:          req.GetURL(),
		PilgrimID:    req.GetPilgrimID(),
		Notes:        req.GetNotes(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.StoreCollectiveDocumentResponse{DocumentID: result.DocumentID}, nil
}

func (s *Server) GetCollectiveDocuments(ctx context.Context, req *pb.GetCollectiveDocumentsRequest) (*pb.GetCollectiveDocumentsResponse, error) {
	const op = "grpc_api.Server.GetCollectiveDocuments"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetCollectiveDocuments(ctx, &service.GetCollectiveDocumentsParams{
		DepartureID:  req.GetDepartureID(),
		PilgrimID:    req.GetPilgrimID(),
		DocumentType: req.GetDocumentType(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	var docs []*pb.CollectiveDocRow
	for _, d := range result.Documents {
		docs = append(docs, &pb.CollectiveDocRow{
			DocumentID:   d.DocumentID,
			DocumentType: d.DocumentType,
			URL:          d.URL,
			PilgrimID:    d.PilgrimID,
			UploadedAt:   d.UploadedAt,
		})
	}
	return &pb.GetCollectiveDocumentsResponse{DepartureID: result.DepartureID, Documents: docs}, nil
}

func (s *Server) SetDocumentACL(ctx context.Context, req *pb.SetDocumentACLRequest) (*pb.SetDocumentACLResponse, error) {
	const op = "grpc_api.Server.SetDocumentACL"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("document_id", req.GetDocumentID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.SetDocumentACL(ctx, &service.SetDocumentACLParams{
		DocumentID:     req.GetDocumentID(),
		AccessLevel:    req.GetAccessLevel(),
		AllowedUserIDs: req.GetAllowedUserIDs(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.SetDocumentACLResponse{Updated: result.Updated}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-022: Passport OCR & mahram
// ---------------------------------------------------------------------------

func (s *Server) ExtractPassportOCR(ctx context.Context, req *pb.ExtractPassportOCRRequest) (*pb.ExtractPassportOCRResponse, error) {
	const op = "grpc_api.Server.ExtractPassportOCR"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("pilgrim_id", req.GetPilgrimID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.ExtractPassportOCR(ctx, &service.ExtractPassportOCRParams{
		ImageURL:  req.GetImageURL(),
		PilgrimID: req.GetPilgrimID(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	var ocrData *pb.PassportOCRData
	if result.Data != nil {
		ocrData = &pb.PassportOCRData{
			FullName:       result.Data.FullName,
			PassportNumber: result.Data.PassportNumber,
			Nationality:    result.Data.Nationality,
			DateOfBirth:    result.Data.DateOfBirth,
			ExpiryDate:     result.Data.ExpiryDate,
			Gender:         result.Data.Gender,
			Confidence:     result.Data.Confidence,
		}
	}
	return &pb.ExtractPassportOCRResponse{
		PilgrimID: result.PilgrimID,
		Data:      ocrData,
		Warnings:  result.Warnings,
	}, nil
}

func (s *Server) SetMahramRelation(ctx context.Context, req *pb.SetMahramRelationRequest) (*pb.SetMahramRelationResponse, error) {
	const op = "grpc_api.Server.SetMahramRelation"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("pilgrim_id", req.GetPilgrimID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.SetMahramRelation(ctx, &service.SetMahramRelationParams{
		PilgrimID:       req.GetPilgrimID(),
		MahramPilgrimID: req.GetMahramPilgrimID(),
		Relation:        req.GetRelation(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.SetMahramRelationResponse{RelationID: result.RelationID}, nil
}

func (s *Server) GetMahramRelations(ctx context.Context, req *pb.GetMahramRelationsRequest) (*pb.GetMahramRelationsResponse, error) {
	const op = "grpc_api.Server.GetMahramRelations"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("booking_id", req.GetBookingID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetMahramRelations(ctx, &service.GetMahramRelationsParams{
		BookingID: req.GetBookingID(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	var rels []*pb.MahramRelationRow
	for _, r := range result.Relations {
		rels = append(rels, &pb.MahramRelationRow{
			RelationID:      r.RelationID,
			PilgrimID:       r.PilgrimID,
			MahramPilgrimID: r.MahramPilgrimID,
			Relation:        r.Relation,
		})
	}
	return &pb.GetMahramRelationsResponse{BookingID: result.BookingID, Relations: rels}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-023: Document progress & expiry
// ---------------------------------------------------------------------------

func (s *Server) GetDocumentProgress(ctx context.Context, req *pb.GetDocumentProgressRequest) (*pb.GetDocumentProgressResponse, error) {
	const op = "grpc_api.Server.GetDocumentProgress"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetDocumentProgress(ctx, &service.GetDocumentProgressParams{
		DepartureID: req.GetDepartureID(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	var rows []*pb.DocProgressRow
	for _, r := range result.Rows {
		rows = append(rows, &pb.DocProgressRow{
			PilgrimID:       r.PilgrimID,
			DocumentType:    r.DocumentType,
			Status:          r.Status,
			ExpiryDate:      r.ExpiryDate,
			DaysUntilExpiry: r.DaysUntilExpiry,
		})
	}
	return &pb.GetDocumentProgressResponse{
		DepartureID:       result.DepartureID,
		Rows:              rows,
		TotalPilgrims:     result.TotalPilgrims,
		DocumentsComplete: result.DocumentsComplete,
		DocumentsExpiring: result.DocumentsExpiring,
	}, nil
}

func (s *Server) GetExpiryAlerts(ctx context.Context, req *pb.GetExpiryAlertsRequest) (*pb.GetExpiryAlertsResponse, error) {
	const op = "grpc_api.Server.GetExpiryAlerts"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetExpiryAlerts(ctx, &service.GetExpiryAlertsParams{
		ThresholdDays: req.GetThresholdDays(),
		DepartureID:   req.GetDepartureID(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	var alerts []*pb.ExpiryAlertRow
	for _, a := range result.Alerts {
		alerts = append(alerts, &pb.ExpiryAlertRow{
			PilgrimID:       a.PilgrimID,
			DocumentType:    a.DocumentType,
			ExpiryDate:      a.ExpiryDate,
			DaysUntilExpiry: a.DaysUntilExpiry,
		})
	}
	return &pb.GetExpiryAlertsResponse{Alerts: alerts}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-024: Official letter
// ---------------------------------------------------------------------------

func (s *Server) GenerateOfficialLetter(ctx context.Context, req *pb.GenerateOfficialLetterRequest) (*pb.GenerateOfficialLetterResponse, error) {
	const op = "grpc_api.Server.GenerateOfficialLetter"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GenerateOfficialLetter(ctx, &service.GenerateOfficialLetterParams{
		TemplateName: req.GetTemplateName(),
		DepartureID:  req.GetDepartureID(),
		PilgrimID:    req.GetPilgrimID(),
		IssuedTo:     req.GetIssuedTo(),
		Notes:        req.GetNotes(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.GenerateOfficialLetterResponse{
		LetterID:  result.LetterID,
		LetterURL: result.LetterURL,
		IssuedAt:  result.IssuedAt,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-025: Immigration manifest
// ---------------------------------------------------------------------------

func (s *Server) GenerateImmigrationManifest(ctx context.Context, req *pb.GenerateImmigrationManifestRequest) (*pb.GenerateImmigrationManifestResponse, error) {
	const op = "grpc_api.Server.GenerateImmigrationManifest"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GenerateImmigrationManifest(ctx, &service.GenerateImmigrationManifestParams{
		DepartureID: req.GetDepartureID(),
		Format:      req.GetFormat(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.GenerateImmigrationManifestResponse{
		ManifestID:  result.ManifestID,
		ManifestURL: result.ManifestURL,
		Version:     result.Version,
		RowCount:    result.RowCount,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-026: Smart rooming
// ---------------------------------------------------------------------------

func (s *Server) RunSmartRooming(ctx context.Context, req *pb.RunSmartRoomingRequest) (*pb.RunSmartRoomingResponse, error) {
	const op = "grpc_api.Server.RunSmartRooming"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	var opts *service.SmartRoomingOptionsParams
	if o := req.GetOptions(); o != nil {
		opts = &service.SmartRoomingOptionsParams{
			MaxPerRoom:       o.GetMaxPerRoom(),
			SeparateMahram:   o.GetSeparateMahram(),
			PreferSameFamily: o.GetPreferSameFamily(),
		}
	}

	result, err := s.svc.RunSmartRooming(ctx, &service.RunSmartRoomingParams{
		DepartureID: req.GetDepartureID(),
		Options:     opts,
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.RunSmartRoomingResponse{
		DepartureID:   result.DepartureID,
		RoomsAssigned: result.RoomsAssigned,
		TotalPilgrims: result.TotalPilgrims,
		Warnings:      result.Warnings,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-027: Transport arrangement
// ---------------------------------------------------------------------------

func (s *Server) AssignTransport(ctx context.Context, req *pb.AssignTransportRequest) (*pb.AssignTransportResponse, error) {
	const op = "grpc_api.Server.AssignTransport"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.AssignTransport(ctx, &service.AssignTransportParams{
		DepartureID: req.GetDepartureID(),
		VehicleType: req.GetVehicleType(),
		VehicleID:   req.GetVehicleID(),
		PilgrimIDs:  req.GetPilgrimIDs(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.AssignTransportResponse{
		AssignmentID:  result.AssignmentID,
		AssignedCount: result.AssignedCount,
	}, nil
}

func (s *Server) GetTransportAssignments(ctx context.Context, req *pb.GetTransportAssignmentsRequest) (*pb.GetTransportAssignmentsResponse, error) {
	const op = "grpc_api.Server.GetTransportAssignments"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetTransportAssignments(ctx, &service.GetTransportAssignmentsParams{
		DepartureID: req.GetDepartureID(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	var assignments []*pb.TransportAssignmentRow
	for _, a := range result.Assignments {
		assignments = append(assignments, &pb.TransportAssignmentRow{
			AssignmentID: a.AssignmentID,
			VehicleType:  a.VehicleType,
			VehicleID:    a.VehicleID,
			PilgrimIDs:   a.PilgrimIDs,
		})
	}
	return &pb.GetTransportAssignmentsResponse{DepartureID: result.DepartureID, Assignments: assignments}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-028: Manifest delta
// ---------------------------------------------------------------------------

func (s *Server) PublishManifestDelta(ctx context.Context, req *pb.PublishManifestDeltaRequest) (*pb.PublishManifestDeltaResponse, error) {
	const op = "grpc_api.Server.PublishManifestDelta"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.PublishManifestDelta(ctx, &service.PublishManifestDeltaParams{
		DepartureID: req.GetDepartureID(),
		ChangeType:  req.GetChangeType(),
		EntityID:    req.GetEntityID(),
		Notes:       req.GetNotes(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.PublishManifestDeltaResponse{DeltaID: result.DeltaID, PublishedAt: result.PublishedAt}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-029: Staff assignment
// ---------------------------------------------------------------------------

func (s *Server) AssignStaff(ctx context.Context, req *pb.AssignStaffRequest) (*pb.AssignStaffResponse, error) {
	const op = "grpc_api.Server.AssignStaff"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.AssignStaff(ctx, &service.AssignStaffParams{
		DepartureID: req.GetDepartureID(),
		StaffUserID: req.GetStaffUserID(),
		Role:        req.GetRole(),
		PilgrimIDs:  req.GetPilgrimIDs(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.AssignStaffResponse{AssignmentID: result.AssignmentID}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-030: Passport log
// ---------------------------------------------------------------------------

func (s *Server) RecordPassportHandover(ctx context.Context, req *pb.RecordPassportHandoverRequest) (*pb.RecordPassportHandoverResponse, error) {
	const op = "grpc_api.Server.RecordPassportHandover"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.RecordPassportHandover(ctx, &service.RecordPassportHandoverParams{
		DepartureID: req.GetDepartureID(),
		PilgrimID:   req.GetPilgrimID(),
		FromUserID:  req.GetFromUserID(),
		ToUserID:    req.GetToUserID(),
		Notes:       req.GetNotes(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.RecordPassportHandoverResponse{HandoverID: result.HandoverID, RecordedAt: result.RecordedAt}, nil
}

func (s *Server) GetPassportLog(ctx context.Context, req *pb.GetPassportLogRequest) (*pb.GetPassportLogResponse, error) {
	const op = "grpc_api.Server.GetPassportLog"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetPassportLog(ctx, &service.GetPassportLogParams{
		DepartureID: req.GetDepartureID(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	var rows []*pb.PassportHandoverRow
	for _, r := range result.Rows {
		rows = append(rows, &pb.PassportHandoverRow{
			HandoverID: r.HandoverID,
			PilgrimID:  r.PilgrimID,
			FromUserID: r.FromUserID,
			ToUserID:   r.ToUserID,
			RecordedAt: r.RecordedAt,
			Notes:      r.Notes,
		})
	}
	return &pb.GetPassportLogResponse{DepartureID: result.DepartureID, Rows: rows}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-031: Visa progress
// ---------------------------------------------------------------------------

func (s *Server) GetVisaProgress(ctx context.Context, req *pb.GetVisaProgressRequest) (*pb.GetVisaProgressResponse, error) {
	const op = "grpc_api.Server.GetVisaProgress"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetVisaProgress(ctx, &service.GetVisaProgressParams{
		DepartureID: req.GetDepartureID(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	var rows []*pb.VisaProgressRow
	for _, r := range result.Rows {
		rows = append(rows, &pb.VisaProgressRow{
			PilgrimID:     r.PilgrimID,
			VisaStatus:    r.VisaStatus,
			SubmittedAt:   r.SubmittedAt,
			ExpectedBy:    r.ExpectedBy,
			DaysRemaining: r.DaysRemaining,
		})
	}
	return &pb.GetVisaProgressResponse{
		DepartureID: result.DepartureID,
		Rows:        rows,
		Submitted:   result.Submitted,
		Approved:    result.Approved,
		Rejected:    result.Rejected,
		Pending:     result.Pending,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-032: E-visa repository
// ---------------------------------------------------------------------------

func (s *Server) StoreEVisa(ctx context.Context, req *pb.StoreEVisaRequest) (*pb.StoreEVisaResponse, error) {
	const op = "grpc_api.Server.StoreEVisa"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("pilgrim_id", req.GetPilgrimID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.StoreEVisa(ctx, &service.StoreEVisaParams{
		PilgrimID:   req.GetPilgrimID(),
		DepartureID: req.GetDepartureID(),
		VisaNumber:  req.GetVisaNumber(),
		VisaURL:     req.GetVisaURL(),
		IssuedDate:  req.GetIssuedDate(),
		ExpiryDate:  req.GetExpiryDate(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.StoreEVisaResponse{EVisaID: result.EVisaID}, nil
}

func (s *Server) GetEVisa(ctx context.Context, req *pb.GetEVisaRequest) (*pb.GetEVisaResponse, error) {
	const op = "grpc_api.Server.GetEVisa"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("pilgrim_id", req.GetPilgrimID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetEVisa(ctx, &service.GetEVisaParams{
		PilgrimID:   req.GetPilgrimID(),
		DepartureID: req.GetDepartureID(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	var evisa *pb.EVisa
	if result.EVisa != nil {
		evisa = &pb.EVisa{
			EVisaID:    result.EVisa.EVisaID,
			PilgrimID:  result.EVisa.PilgrimID,
			VisaNumber: result.EVisa.VisaNumber,
			VisaURL:    result.EVisa.VisaURL,
			IssuedDate: result.EVisa.IssuedDate,
			ExpiryDate: result.EVisa.ExpiryDate,
		}
	}
	return &pb.GetEVisaResponse{EVisa: evisa}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-033: External provider (stub)
// ---------------------------------------------------------------------------

func (s *Server) TriggerExternalProvider(ctx context.Context, req *pb.TriggerExternalProviderRequest) (*pb.TriggerExternalProviderResponse, error) {
	const op = "grpc_api.Server.TriggerExternalProvider"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("provider", req.GetProvider()))
	logger.Info().Str("op", op).Str("provider", req.GetProvider()).Msg("")

	result, err := s.svc.TriggerExternalProvider(ctx, &service.TriggerExternalProviderParams{
		Provider:    req.GetProvider(),
		Action:      req.GetAction(),
		ReferenceID: req.GetReferenceID(),
		Payload:     req.GetPayload(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.TriggerExternalProviderResponse{RequestID: result.RequestID, Status: result.Status}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-034: Refund & penalty
// ---------------------------------------------------------------------------

func (s *Server) CreateRefund(ctx context.Context, req *pb.CreateRefundRequest) (*pb.CreateRefundResponse, error) {
	const op = "grpc_api.Server.CreateRefund"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("booking_id", req.GetBookingID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.CreateRefund(ctx, &service.CreateRefundParams{
		BookingID: req.GetBookingID(),
		Reason:    req.GetReason(),
		Amount:    req.GetAmount(),
		Notes:     req.GetNotes(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.CreateRefundResponse{RefundID: result.RefundID, Status: result.Status}, nil
}

func (s *Server) ApproveRefund(ctx context.Context, req *pb.ApproveRefundRequest) (*pb.ApproveRefundResponse, error) {
	const op = "grpc_api.Server.ApproveRefund"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("refund_id", req.GetRefundID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.ApproveRefund(ctx, &service.ApproveRefundParams{
		RefundID: req.GetRefundID(),
		Decision: req.GetDecision(),
		Notes:    req.GetNotes(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.ApproveRefundResponse{RefundID: result.RefundID, Status: result.Status}, nil
}

func (s *Server) RecordPenalty(ctx context.Context, req *pb.RecordPenaltyRequest) (*pb.RecordPenaltyResponse, error) {
	const op = "grpc_api.Server.RecordPenalty"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("booking_id", req.GetBookingID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.RecordPenalty(ctx, &service.RecordPenaltyParams{
		BookingID:   req.GetBookingID(),
		PenaltyType: req.GetPenaltyType(),
		Amount:      req.GetAmount(),
		Notes:       req.GetNotes(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.RecordPenaltyResponse{PenaltyID: result.PenaltyID}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-036: Luggage counter
// ---------------------------------------------------------------------------

func (s *Server) RecordLuggageScan(ctx context.Context, req *pb.RecordLuggageScanRequest) (*pb.RecordLuggageScanResponse, error) {
	const op = "grpc_api.Server.RecordLuggageScan"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.RecordLuggageScan(ctx, &service.RecordLuggageScanParams{
		DepartureID: req.GetDepartureID(),
		PilgrimID:   req.GetPilgrimID(),
		TagID:       req.GetTagID(),
		ScanPoint:   req.GetScanPoint(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.RecordLuggageScanResponse{ScanID: result.ScanID, TotalBags: result.TotalBags}, nil
}

func (s *Server) GetLuggageCount(ctx context.Context, req *pb.GetLuggageCountRequest) (*pb.GetLuggageCountResponse, error) {
	const op = "grpc_api.Server.GetLuggageCount"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.GetLuggageCount(ctx, &service.GetLuggageCountParams{
		DepartureID: req.GetDepartureID(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}

	var rows []*pb.LuggageCountRow
	for _, r := range result.Rows {
		rows = append(rows, &pb.LuggageCountRow{
			PilgrimID:     r.PilgrimID,
			BagCount:      r.BagCount,
			LastScannedAt: r.LastScannedAt,
		})
	}
	return &pb.GetLuggageCountResponse{
		DepartureID: result.DepartureID,
		TotalBags:   result.TotalBags,
		Rows:        rows,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-037: Departure/arrival broadcast
// ---------------------------------------------------------------------------

func (s *Server) BroadcastSchedule(ctx context.Context, req *pb.BroadcastScheduleRequest) (*pb.BroadcastScheduleResponse, error) {
	const op = "grpc_api.Server.BroadcastSchedule"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.BroadcastSchedule(ctx, &service.BroadcastScheduleParams{
		DepartureID:   req.GetDepartureID(),
		BroadcastType: req.GetBroadcastType(),
		Message:       req.GetMessage(),
		Channel:       req.GetChannel(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.BroadcastScheduleResponse{BroadcastID: result.BroadcastID, RecipientCount: result.RecipientCount}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-039: Raudhah shield & tasreh
// ---------------------------------------------------------------------------

func (s *Server) IssueDigitalTasreh(ctx context.Context, req *pb.IssueDigitalTasrehRequest) (*pb.IssueDigitalTasrehResponse, error) {
	const op = "grpc_api.Server.IssueDigitalTasreh"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("pilgrim_id", req.GetPilgrimID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.IssueDigitalTasreh(ctx, &service.IssueDigitalTasrehParams{
		PilgrimID:   req.GetPilgrimID(),
		DepartureID: req.GetDepartureID(),
		VisitDate:   req.GetVisitDate(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.IssueDigitalTasrehResponse{TasrehID: result.TasrehID, QRCode: result.QRCode}, nil
}

func (s *Server) RecordRaudhahEntry(ctx context.Context, req *pb.RecordRaudhahEntryRequest) (*pb.RecordRaudhahEntryResponse, error) {
	const op = "grpc_api.Server.RecordRaudhahEntry"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("tasreh_id", req.GetTasrehID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.RecordRaudhahEntry(ctx, &service.RecordRaudhahEntryParams{
		TasrehID:  req.GetTasrehID(),
		PilgrimID: req.GetPilgrimID(),
		EntryTime: req.GetEntryTime(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.RecordRaudhahEntryResponse{RecordID: result.RecordID, Valid: result.Valid}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-040: Audio devices
// ---------------------------------------------------------------------------

func (s *Server) RegisterAudioDevice(ctx context.Context, req *pb.RegisterAudioDeviceRequest) (*pb.RegisterAudioDeviceResponse, error) {
	const op = "grpc_api.Server.RegisterAudioDevice"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.RegisterAudioDevice(ctx, &service.RegisterAudioDeviceParams{
		DepartureID:  req.GetDepartureID(),
		DeviceType:   req.GetDeviceType(),
		SerialNumber: req.GetSerialNumber(),
		AssignedTo:   req.GetAssignedTo(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.RegisterAudioDeviceResponse{DeviceID: result.DeviceID}, nil
}

func (s *Server) UpdateAudioDeviceStatus(ctx context.Context, req *pb.UpdateAudioDeviceStatusRequest) (*pb.UpdateAudioDeviceStatusResponse, error) {
	const op = "grpc_api.Server.UpdateAudioDeviceStatus"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("device_id", req.GetDeviceID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.UpdateAudioDeviceStatus(ctx, &service.UpdateAudioDeviceStatusParams{
		DeviceID: req.GetDeviceID(),
		Status:   req.GetStatus(),
		Notes:    req.GetNotes(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.UpdateAudioDeviceStatusResponse{DeviceID: result.DeviceID, Status: result.Status}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-041: Zamzam distribution
// ---------------------------------------------------------------------------

func (s *Server) RecordZamzamDistribution(ctx context.Context, req *pb.RecordZamzamDistributionRequest) (*pb.RecordZamzamDistributionResponse, error) {
	const op = "grpc_api.Server.RecordZamzamDistribution"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.RecordZamzamDistribution(ctx, &service.RecordZamzamDistributionParams{
		DepartureID: req.GetDepartureID(),
		PilgrimID:   req.GetPilgrimID(),
		LitersGiven: req.GetLitersGiven(),
		ReceivedBy:  req.GetReceivedBy(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.RecordZamzamDistributionResponse{DistributionID: result.DistributionID, RecordedAt: result.RecordedAt}, nil
}

// ---------------------------------------------------------------------------
// BL-OPS-042: Room check-in
// ---------------------------------------------------------------------------

func (s *Server) RecordRoomCheckIn(ctx context.Context, req *pb.RecordRoomCheckInRequest) (*pb.RecordRoomCheckInResponse, error) {
	const op = "grpc_api.Server.RecordRoomCheckIn"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(attribute.String("departure_id", req.GetDepartureID()))
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.RecordRoomCheckIn(ctx, &service.RecordRoomCheckInParams{
		DepartureID: req.GetDepartureID(),
		PilgrimID:   req.GetPilgrimID(),
		RoomNumber:  req.GetRoomNumber(),
		HotelID:     req.GetHotelID(),
	})
	if err != nil {
		logger.Error().Err(err).Str("op", op).Msg("")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Error(grpcCodes.Internal, "internal_error")
	}
	return &pb.RecordRoomCheckInResponse{CheckInID: result.CheckInID, CheckedInAt: result.CheckedInAt}, nil
}
