// logistics_depth.go — gRPC handlers for Wave 5 depth RPCs (BL-LOG-013..029).
//
// Implements all 20 LogisticsDepthHandler methods on *Server by delegating
// to the service layer and mapping results to pb response types.

package grpc_api

import (
	"context"

	"logistics-svc/api/grpc_api/pb"
	"logistics-svc/service"
	"logistics-svc/util/logging"

	otelCodes "go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// BL-LOG-013: ListPurchaseRequests
// ---------------------------------------------------------------------------

func (s *Server) ListPurchaseRequests(ctx context.Context, req *pb.ListPurchaseRequestsRequest) (*pb.ListPurchaseRequestsResponse, error) {
	const op = "grpc_api.Server.ListPurchaseRequests"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.ListPurchaseRequests(ctx, &service.ListPurchaseRequestsParams{
		Status:      req.GetStatus(),
		DepartureID: req.GetDepartureID(),
		PageSize:    req.GetPageSize(),
		Cursor:      req.GetCursor(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.ListPurchaseRequestsResponse{}, nil
	}

	rows := make([]*pb.PurchaseRequestRow, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, &pb.PurchaseRequestRow{
			PRID:        r.PRID,
			DepartureID: r.DepartureID,
			Description: r.Description,
			Status:      r.Status,
			RequestedBy: r.RequestedBy,
			Amount:      r.Amount,
			CreatedAt:   r.CreatedAt,
		})
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.ListPurchaseRequestsResponse{Rows: rows, NextCursor: result.NextCursor}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-014: GetBudgetSyncStatus
// ---------------------------------------------------------------------------

func (s *Server) GetBudgetSyncStatus(ctx context.Context, req *pb.GetBudgetSyncStatusRequest) (*pb.GetBudgetSyncStatusResponse, error) {
	const op = "grpc_api.Server.GetBudgetSyncStatus"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetBudgetSyncStatus(ctx, &service.GetBudgetSyncStatusParams{
		DepartureID: req.GetDepartureID(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.GetBudgetSyncStatusResponse{DepartureID: req.GetDepartureID()}, nil
	}

	lines := make([]*pb.BudgetSyncLine, 0, len(result.Lines))
	for _, l := range result.Lines {
		lines = append(lines, &pb.BudgetSyncLine{
			Category:  l.Category,
			Budgeted:  l.Budgeted,
			Committed: l.Committed,
			Actual:    l.Actual,
		})
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetBudgetSyncStatusResponse{
		DepartureID:    result.DepartureID,
		Lines:          lines,
		TotalBudgeted:  result.TotalBudgeted,
		TotalCommitted: result.TotalCommitted,
		TotalActual:    result.TotalActual,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-015: GetTieredApprovals
// ---------------------------------------------------------------------------

func (s *Server) GetTieredApprovals(ctx context.Context, req *pb.GetTieredApprovalsRequest) (*pb.GetTieredApprovalsResponse, error) {
	const op = "grpc_api.Server.GetTieredApprovals"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetTieredApprovals(ctx, &service.GetTieredApprovalsParams{
		Level:  req.GetLevel(),
		Status: req.GetStatus(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.GetTieredApprovalsResponse{}, nil
	}

	rows := make([]*pb.TieredApprovalRow, 0, len(result.Rows))
	for _, r := range result.Rows {
		rows = append(rows, &pb.TieredApprovalRow{
			ApprovalID:  r.ApprovalID,
			PRID:        r.PRID,
			RequestedBy: r.RequestedBy,
			Amount:      r.Amount,
			Level:       r.Level,
			Status:      r.Status,
			CreatedAt:   r.CreatedAt,
		})
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetTieredApprovalsResponse{Rows: rows}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-015: DecideTieredApproval
// ---------------------------------------------------------------------------

func (s *Server) DecideTieredApproval(ctx context.Context, req *pb.DecideTieredApprovalRequest) (*pb.DecideTieredApprovalResponse, error) {
	const op = "grpc_api.Server.DecideTieredApproval"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.DecideTieredApproval(ctx, &service.DecideTieredApprovalParams{
		ApprovalID: req.GetApprovalID(),
		Decision:   req.GetDecision(),
		Notes:      req.GetNotes(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.DecideTieredApprovalResponse{ApprovalID: req.GetApprovalID()}, nil
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.DecideTieredApprovalResponse{ApprovalID: result.ApprovalID, Status: result.Status}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-016: AutoSelectVendor
// ---------------------------------------------------------------------------

func (s *Server) AutoSelectVendor(ctx context.Context, req *pb.AutoSelectVendorRequest) (*pb.AutoSelectVendorResponse, error) {
	const op = "grpc_api.Server.AutoSelectVendor"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.AutoSelectVendor(ctx, &service.AutoSelectVendorParams{
		Category:    req.GetCategory(),
		DepartureID: req.GetDepartureID(),
		RequiredBy:  req.GetRequiredBy(),
		MaxBudget:   req.GetMaxBudget(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.AutoSelectVendorResponse{}, nil
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.AutoSelectVendorResponse{
		VendorID:      result.VendorID,
		VendorName:    result.VendorName,
		Reason:        result.Reason,
		EstimatedCost: result.EstimatedCost,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-017: RecordPartialGRN
// ---------------------------------------------------------------------------

func (s *Server) RecordPartialGRN(ctx context.Context, req *pb.RecordPartialGRNRequest) (*pb.RecordPartialGRNResponse, error) {
	const op = "grpc_api.Server.RecordPartialGRN"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	items := make([]service.GRNItemParam, 0, len(req.GetItems()))
	for _, it := range req.GetItems() {
		items = append(items, service.GRNItemParam{
			SKU:      it.GetSKU(),
			Quantity: it.GetQuantity(),
			UnitCost: it.GetUnitCost(),
		})
	}

	result, err := s.svc.RecordPartialGRN(ctx, &service.RecordPartialGRNParams{
		PRID:  req.GetPRID(),
		Notes: req.GetNotes(),
		Items: items,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.RecordPartialGRNResponse{}, nil
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.RecordPartialGRNResponse{
		GRNID:           result.GRNID,
		ItemsReceived:   result.ItemsReceived,
		IsFullyReceived: result.IsFullyReceived,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-017: ReverseGRN
// ---------------------------------------------------------------------------

func (s *Server) ReverseGRN(ctx context.Context, req *pb.ReverseGRNRequest) (*pb.ReverseGRNResponse, error) {
	const op = "grpc_api.Server.ReverseGRN"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.ReverseGRN(ctx, &service.ReverseGRNParams{
		GRNID:  req.GetGRNID(),
		Reason: req.GetReason(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.ReverseGRNResponse{}, nil
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.ReverseGRNResponse{ReversalID: result.ReversalID, Reversed: result.Reversed}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-020: GenerateBarcode
// ---------------------------------------------------------------------------

func (s *Server) GenerateBarcode(ctx context.Context, req *pb.GenerateBarcodeRequest) (*pb.GenerateBarcodeResponse, error) {
	const op = "grpc_api.Server.GenerateBarcode"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GenerateBarcode(ctx, &service.GenerateBarcodeParams{
		SKU:         req.GetSKU(),
		ItemName:    req.GetItemName(),
		DepartureID: req.GetDepartureID(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.GenerateBarcodeResponse{}, nil
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GenerateBarcodeResponse{BarcodeData: result.BarcodeData, LabelURL: result.LabelURL}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-020: PrintSKULabel
// ---------------------------------------------------------------------------

func (s *Server) PrintSKULabel(ctx context.Context, req *pb.PrintSKULabelRequest) (*pb.PrintSKULabelResponse, error) {
	const op = "grpc_api.Server.PrintSKULabel"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.PrintSKULabel(ctx, &service.PrintSKULabelParams{
		SKU:      req.GetSKU(),
		Quantity: req.GetQuantity(),
		Format:   req.GetFormat(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.PrintSKULabelResponse{}, nil
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.PrintSKULabelResponse{LabelURL: result.LabelURL, LabelCount: result.LabelCount}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-021: CreateWarehouse
// ---------------------------------------------------------------------------

func (s *Server) CreateWarehouse(ctx context.Context, req *pb.CreateWarehouseRequest) (*pb.CreateWarehouseResponse, error) {
	const op = "grpc_api.Server.CreateWarehouse"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.CreateWarehouse(ctx, &service.CreateWarehouseParams{
		Name:     req.GetName(),
		Location: req.GetLocation(),
		Type:     req.GetType(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.CreateWarehouseResponse{}, nil
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.CreateWarehouseResponse{WarehouseID: result.WarehouseID}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-021: TransferStock
// ---------------------------------------------------------------------------

func (s *Server) TransferStock(ctx context.Context, req *pb.TransferStockRequest) (*pb.TransferStockResponse, error) {
	const op = "grpc_api.Server.TransferStock"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.TransferStock(ctx, &service.TransferStockParams{
		FromWarehouseID: req.GetFromWarehouseID(),
		ToWarehouseID:   req.GetToWarehouseID(),
		SKU:             req.GetSKU(),
		Quantity:        req.GetQuantity(),
		Notes:           req.GetNotes(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.TransferStockResponse{}, nil
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.TransferStockResponse{TransferID: result.TransferID}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-022: GetStockAlerts
// ---------------------------------------------------------------------------

func (s *Server) GetStockAlerts(ctx context.Context, req *pb.GetStockAlertsRequest) (*pb.GetStockAlertsResponse, error) {
	const op = "grpc_api.Server.GetStockAlerts"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetStockAlerts(ctx, &service.GetStockAlertsParams{
		WarehouseID: req.GetWarehouseID(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.GetStockAlertsResponse{}, nil
	}

	alerts := make([]*pb.StockAlertRow, 0, len(result.Alerts))
	for _, a := range result.Alerts {
		alerts = append(alerts, &pb.StockAlertRow{
			SKU:          a.SKU,
			ItemName:     a.ItemName,
			WarehouseID:  a.WarehouseID,
			CurrentQty:   a.CurrentQty,
			ReorderLevel: a.ReorderLevel,
			Severity:     a.Severity,
		})
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.GetStockAlertsResponse{
		Alerts:        alerts,
		TotalCritical: result.TotalCritical,
		TotalWarning:  result.TotalWarning,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-022: SetReorderLevel
// ---------------------------------------------------------------------------

func (s *Server) SetReorderLevel(ctx context.Context, req *pb.SetReorderLevelRequest) (*pb.SetReorderLevelResponse, error) {
	const op = "grpc_api.Server.SetReorderLevel"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.SetReorderLevel(ctx, &service.SetReorderLevelParams{
		SKU:          req.GetSKU(),
		WarehouseID:  req.GetWarehouseID(),
		ReorderLevel: req.GetReorderLevel(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.SetReorderLevelResponse{Updated: true}, nil
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.SetReorderLevelResponse{Updated: result.Updated}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-024: StartStocktake
// ---------------------------------------------------------------------------

func (s *Server) StartStocktake(ctx context.Context, req *pb.StartStocktakeRequest) (*pb.StartStocktakeResponse, error) {
	const op = "grpc_api.Server.StartStocktake"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.StartStocktake(ctx, &service.StartStocktakeParams{
		WarehouseID: req.GetWarehouseID(),
		Notes:       req.GetNotes(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.StartStocktakeResponse{}, nil
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.StartStocktakeResponse{StocktakeID: result.StocktakeID, StartedAt: result.StartedAt}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-024: RecordStocktakeCount
// ---------------------------------------------------------------------------

func (s *Server) RecordStocktakeCount(ctx context.Context, req *pb.RecordStocktakeCountRequest) (*pb.RecordStocktakeCountResponse, error) {
	const op = "grpc_api.Server.RecordStocktakeCount"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.RecordStocktakeCount(ctx, &service.RecordStocktakeCountParams{
		StocktakeID: req.GetStocktakeID(),
		SKU:         req.GetSKU(),
		CountedQty:  req.GetCountedQty(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.RecordStocktakeCountResponse{}, nil
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.RecordStocktakeCountResponse{
		LineID:      result.LineID,
		SystemQty:   result.SystemQty,
		CountedQty:  result.CountedQty,
		VarianceQty: result.VarianceQty,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-024: FinalizeStocktake
// ---------------------------------------------------------------------------

func (s *Server) FinalizeStocktake(ctx context.Context, req *pb.FinalizeStocktakeRequest) (*pb.FinalizeStocktakeResponse, error) {
	const op = "grpc_api.Server.FinalizeStocktake"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.FinalizeStocktake(ctx, &service.FinalizeStocktakeParams{
		StocktakeID: req.GetStocktakeID(),
		Notes:       req.GetNotes(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.FinalizeStocktakeResponse{StocktakeID: req.GetStocktakeID()}, nil
	}

	lines := make([]*pb.StocktakeVarianceLine, 0, len(result.Lines))
	for _, l := range result.Lines {
		lines = append(lines, &pb.StocktakeVarianceLine{
			SKU:         l.SKU,
			ItemName:    l.ItemName,
			SystemQty:   l.SystemQty,
			CountedQty:  l.CountedQty,
			VarianceQty: l.VarianceQty,
		})
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.FinalizeStocktakeResponse{
		StocktakeID:   result.StocktakeID,
		TotalItems:    result.TotalItems,
		VarianceItems: result.VarianceItems,
		Lines:         lines,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-025: SyncFulfillmentSizes
// ---------------------------------------------------------------------------

func (s *Server) SyncFulfillmentSizes(ctx context.Context, req *pb.SyncFulfillmentSizesRequest) (*pb.SyncFulfillmentSizesResponse, error) {
	const op = "grpc_api.Server.SyncFulfillmentSizes"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.SyncFulfillmentSizes(ctx, &service.SyncFulfillmentSizesParams{
		DepartureID: req.GetDepartureID(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.SyncFulfillmentSizesResponse{DepartureID: req.GetDepartureID()}, nil
	}

	sizes := make([]*pb.SizePreset, 0, len(result.Sizes))
	for _, sz := range result.Sizes {
		sizes = append(sizes, &pb.SizePreset{PilgrimID: sz.PilgrimID, Size: sz.Size})
	}
	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.SyncFulfillmentSizesResponse{
		DepartureID: result.DepartureID,
		SyncedCount: result.SyncedCount,
		Sizes:       sizes,
	}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-027: RecordCourierTracking
// ---------------------------------------------------------------------------

func (s *Server) RecordCourierTracking(ctx context.Context, req *pb.RecordCourierTrackingRequest) (*pb.RecordCourierTrackingResponse, error) {
	const op = "grpc_api.Server.RecordCourierTracking"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.RecordCourierTracking(ctx, &service.RecordCourierTrackingParams{
		FulfillmentTaskID: req.GetFulfillmentTaskID(),
		CourierName:       req.GetCourierName(),
		TrackingNumber:    req.GetTrackingNumber(),
		Status:            req.GetStatus(),
		Note:              req.GetNote(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.RecordCourierTrackingResponse{}, nil
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.RecordCourierTrackingResponse{TrackingID: result.TrackingID, UpdatedAt: result.UpdatedAt}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-029: RecordReturn
// ---------------------------------------------------------------------------

func (s *Server) RecordReturn(ctx context.Context, req *pb.RecordReturnRequest) (*pb.RecordReturnResponse, error) {
	const op = "grpc_api.Server.RecordReturn"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.RecordReturn(ctx, &service.RecordReturnParams{
		FulfillmentTaskID: req.GetFulfillmentTaskID(),
		Reason:            req.GetReason(),
		Condition:         req.GetCondition(),
		ItemsSKUs:         req.GetItemsSKUs(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.RecordReturnResponse{}, nil
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.RecordReturnResponse{ReturnID: result.ReturnID, Status: result.Status}, nil
}

// ---------------------------------------------------------------------------
// BL-LOG-029: ProcessExchange
// ---------------------------------------------------------------------------

func (s *Server) ProcessExchange(ctx context.Context, req *pb.ProcessExchangeRequest) (*pb.ProcessExchangeResponse, error) {
	const op = "grpc_api.Server.ProcessExchange"
	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.ProcessExchange(ctx, &service.ProcessExchangeParams{
		ReturnID: req.GetReturnID(),
		NewSKUs:  req.GetNewSKUs(),
		Notes:    req.GetNotes(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return &pb.ProcessExchangeResponse{}, nil
	}

	span.SetStatus(otelCodes.Ok, "ok")
	return &pb.ProcessExchangeResponse{ExchangeID: result.ExchangeID, Status: result.Status}, nil
}
