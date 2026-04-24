// proxy_logistics_depth.go — gateway REST handlers for logistics-svc Wave 5 depth RPCs.
// BL-LOG-013..029

package rest_oapi

import (
	"errors"
	"strconv"

	"gateway-svc/adapter/logistics_grpc_adapter"
	"gateway-svc/util/apperrors"
	"gateway-svc/util/logging"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// BL-LOG-013: ListPurchaseRequests — GET /v1/logistics/purchase-requests
// ---------------------------------------------------------------------------

func (s *Server) ListPurchaseRequests(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ListPurchaseRequests"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	result, err := s.svc.ListPurchaseRequests(ctx, &logistics_grpc_adapter.ListPurchaseRequestsParams{
		DepartureID: c.Query("departure_id"),
		Status:      c.Query("status"),
		Page:        int32(page),
		PageSize:    int32(pageSize),
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

// ---------------------------------------------------------------------------
// BL-LOG-014: GetBudgetSyncStatus — GET /v1/logistics/budget-sync/:departure_id
// ---------------------------------------------------------------------------

func (s *Server) GetBudgetSyncStatus(c *fiber.Ctx, departureID string) error {
	const op = "rest_oapi.Server.GetBudgetSyncStatus"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetBudgetSyncStatus(ctx, &logistics_grpc_adapter.GetBudgetSyncStatusParams{
		DepartureID: departureID,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

// ---------------------------------------------------------------------------
// BL-LOG-015: GetTieredApprovals — GET /v1/logistics/tiered-approvals
// ---------------------------------------------------------------------------

func (s *Server) GetTieredApprovals(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetTieredApprovals"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetTieredApprovals(ctx, &logistics_grpc_adapter.GetTieredApprovalsParams{
		DepartureID: c.Query("departure_id"),
		Status:      c.Query("status"),
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

// ---------------------------------------------------------------------------
// BL-LOG-015: DecideTieredApproval — PUT /v1/logistics/tiered-approvals/:id/decision
// ---------------------------------------------------------------------------

func (s *Server) DecideTieredApproval(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.DecideTieredApproval"
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
	result, err := s.svc.DecideTieredApproval(ctx, &logistics_grpc_adapter.DecideTieredApprovalParams{
		ApprovalID: id,
		Decision:   body.Decision,
		Notes:      body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"approval_id": result.ApprovalID, "status": result.Status}})
}

// ---------------------------------------------------------------------------
// BL-LOG-016: AutoSelectVendor — POST /v1/logistics/auto-vendor
// ---------------------------------------------------------------------------

func (s *Server) AutoSelectVendor(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.AutoSelectVendor"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		PRID        string `json:"pr_id"`
		CategoryID  string `json:"category_id"`
		BudgetLimit int64  `json:"budget_limit"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.AutoSelectVendor(ctx, &logistics_grpc_adapter.AutoSelectVendorParams{
		PRID:        body.PRID,
		CategoryID:  body.CategoryID,
		BudgetLimit: body.BudgetLimit,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"vendor_id": result.VendorID, "vendor_name": result.VendorName, "score": result.Score}})
}

// ---------------------------------------------------------------------------
// BL-LOG-017: RecordPartialGRN — POST /v1/logistics/partial-grn
// ---------------------------------------------------------------------------

func (s *Server) RecordPartialGRN(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordPartialGRN"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		POID          string `json:"po_id"`
		ReceivedItems []struct {
			SKUID    string `json:"sku_id"`
			Quantity int32  `json:"quantity"`
		} `json:"received_items"`
		Notes string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	items := make([]logistics_grpc_adapter.ReceivedItemInput, 0, len(body.ReceivedItems))
	for _, it := range body.ReceivedItems {
		items = append(items, logistics_grpc_adapter.ReceivedItemInput{
			SKUID:    it.SKUID,
			Quantity: it.Quantity,
		})
	}

	result, err := s.svc.RecordPartialGRN(ctx, &logistics_grpc_adapter.RecordPartialGRNParams{
		POID:          body.POID,
		ReceivedItems: items,
		Notes:         body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"grn_id": result.GRNID, "is_complete": result.IsComplete}})
}

// ---------------------------------------------------------------------------
// BL-LOG-017: ReverseGRN — POST /v1/logistics/reverse-grn
// ---------------------------------------------------------------------------

func (s *Server) ReverseGRN(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ReverseGRN"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		GRNID  string `json:"grn_id"`
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.ReverseGRN(ctx, &logistics_grpc_adapter.ReverseGRNParams{
		GRNID:  body.GRNID,
		Reason: body.Reason,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"reversal_id": result.ReversalID}})
}

// ---------------------------------------------------------------------------
// BL-LOG-018: GenerateBarcode — POST /v1/logistics/barcode
// ---------------------------------------------------------------------------

func (s *Server) GenerateBarcode(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GenerateBarcode"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		SKUID  string `json:"sku_id"`
		Format string `json:"format"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.GenerateBarcode(ctx, &logistics_grpc_adapter.GenerateBarcodeParams{
		SKUID:  body.SKUID,
		Format: body.Format,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"barcode_id": result.BarcodeID, "barcode_url": result.BarcodeURL}})
}

// ---------------------------------------------------------------------------
// BL-LOG-018: PrintSKULabel — POST /v1/logistics/sku-labels
// ---------------------------------------------------------------------------

func (s *Server) PrintSKULabel(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.PrintSKULabel"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		SKUID     string `json:"sku_id"`
		Copies    int32  `json:"copies"`
		PrinterID string `json:"printer_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.PrintSKULabel(ctx, &logistics_grpc_adapter.PrintSKULabelParams{
		SKUID:     body.SKUID,
		Copies:    body.Copies,
		PrinterID: body.PrinterID,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"job_id": result.JobID, "queued": result.Queued}})
}

// ---------------------------------------------------------------------------
// BL-LOG-019: CreateWarehouse — POST /v1/logistics/warehouses
// ---------------------------------------------------------------------------

func (s *Server) CreateWarehouse(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.CreateWarehouse"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Type     string `json:"type"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.CreateWarehouse(ctx, &logistics_grpc_adapter.CreateWarehouseParams{
		Name:     body.Name,
		Location: body.Location,
		Type:     body.Type,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": fiber.Map{"warehouse_id": result.WarehouseID}})
}

// ---------------------------------------------------------------------------
// BL-LOG-020: TransferStock — POST /v1/logistics/stock-transfer
// ---------------------------------------------------------------------------

func (s *Server) TransferStock(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.TransferStock"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		SKUID           string `json:"sku_id"`
		Quantity        int32  `json:"quantity"`
		FromWarehouseID string `json:"from_warehouse_id"`
		ToWarehouseID   string `json:"to_warehouse_id"`
		Notes           string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.TransferStock(ctx, &logistics_grpc_adapter.TransferStockParams{
		SKUID:           body.SKUID,
		Quantity:        body.Quantity,
		FromWarehouseID: body.FromWarehouseID,
		ToWarehouseID:   body.ToWarehouseID,
		Notes:           body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"transfer_id": result.TransferID}})
}

// ---------------------------------------------------------------------------
// BL-LOG-021: GetStockAlerts — GET /v1/logistics/stock-alerts
// ---------------------------------------------------------------------------

func (s *Server) GetStockAlerts(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.GetStockAlerts"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.GetStockAlerts(ctx, &logistics_grpc_adapter.GetStockAlertsParams{
		WarehouseID: c.Query("warehouse_id"),
		Severity:    c.Query("severity"),
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": result})
}

// ---------------------------------------------------------------------------
// BL-LOG-022: SetReorderLevel — PUT /v1/logistics/reorder-levels
// ---------------------------------------------------------------------------

func (s *Server) SetReorderLevel(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.SetReorderLevel"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		SKUID      string `json:"sku_id"`
		ReorderQty int32  `json:"reorder_qty"`
		MinQty     int32  `json:"min_qty"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.SetReorderLevel(ctx, &logistics_grpc_adapter.SetReorderLevelParams{
		SKUID:      body.SKUID,
		ReorderQty: body.ReorderQty,
		MinQty:     body.MinQty,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"sku_id": result.SKUID, "updated_at": result.UpdatedAt}})
}

// ---------------------------------------------------------------------------
// BL-LOG-023: StartStocktake — POST /v1/logistics/stocktake
// ---------------------------------------------------------------------------

func (s *Server) StartStocktake(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.StartStocktake"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		WarehouseID string `json:"warehouse_id"`
		Notes       string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.StartStocktake(ctx, &logistics_grpc_adapter.StartStocktakeParams{
		WarehouseID: body.WarehouseID,
		Notes:       body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": fiber.Map{"stocktake_id": result.StocktakeID}})
}

// ---------------------------------------------------------------------------
// BL-LOG-023: RecordStocktakeCount — POST /v1/logistics/stocktake/:id/count
// ---------------------------------------------------------------------------

func (s *Server) RecordStocktakeCount(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.RecordStocktakeCount"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		SKUID      string `json:"sku_id"`
		CountedQty int32  `json:"counted_qty"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.RecordStocktakeCount(ctx, &logistics_grpc_adapter.RecordStocktakeCountParams{
		StocktakeID: id,
		SKUID:       body.SKUID,
		CountedQty:  body.CountedQty,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"entry_id": result.EntryID}})
}

// ---------------------------------------------------------------------------
// BL-LOG-023: FinalizeStocktake — PUT /v1/logistics/stocktake/:id/finalize
// ---------------------------------------------------------------------------

func (s *Server) FinalizeStocktake(c *fiber.Ctx, id string) error {
	const op = "rest_oapi.Server.FinalizeStocktake"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	result, err := s.svc.FinalizeStocktake(ctx, &logistics_grpc_adapter.FinalizeStocktakeParams{
		StocktakeID: id,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{
		"stocktake_id":  result.StocktakeID,
		"finalized_at":  result.FinalizedAt,
		"discrepancies": result.Discrepancies,
	}})
}

// ---------------------------------------------------------------------------
// BL-LOG-024: SyncFulfillmentSizes — POST /v1/logistics/size-sync
// ---------------------------------------------------------------------------

func (s *Server) SyncFulfillmentSizes(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.SyncFulfillmentSizes"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		DepartureID string `json:"departure_id"`
		Items       []struct {
			BookingID string `json:"booking_id"`
			Size      string `json:"size"`
		} `json:"items"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	items := make([]logistics_grpc_adapter.SizeSyncItemInput, 0, len(body.Items))
	for _, it := range body.Items {
		items = append(items, logistics_grpc_adapter.SizeSyncItemInput{
			BookingID: it.BookingID,
			Size:      it.Size,
		})
	}

	result, err := s.svc.SyncFulfillmentSizes(ctx, &logistics_grpc_adapter.SyncFulfillmentSizesParams{
		DepartureID: body.DepartureID,
		Items:       items,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"synced": result.Synced, "updated_at": result.UpdatedAt}})
}

// ---------------------------------------------------------------------------
// BL-LOG-025: RecordCourierTracking — POST /v1/logistics/courier-tracking
// ---------------------------------------------------------------------------

func (s *Server) RecordCourierTracking(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordCourierTracking"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		BookingID      string `json:"booking_id"`
		Courier        string `json:"courier"`
		TrackingNumber string `json:"tracking_number"`
		Status         string `json:"status"`
		UpdatedAt      string `json:"updated_at"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}
	result, err := s.svc.RecordCourierTracking(ctx, &logistics_grpc_adapter.RecordCourierTrackingParams{
		BookingID:      body.BookingID,
		Courier:        body.Courier,
		TrackingNumber: body.TrackingNumber,
		Status:         body.Status,
		UpdatedAt:      body.UpdatedAt,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"tracking_id": result.TrackingID}})
}

// ---------------------------------------------------------------------------
// BL-LOG-026: RecordReturn — POST /v1/logistics/returns
// ---------------------------------------------------------------------------

func (s *Server) RecordReturn(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.RecordReturn"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		BookingID string `json:"booking_id"`
		Reason    string `json:"reason"`
		Items     []struct {
			SKUID    string `json:"sku_id"`
			Quantity int32  `json:"quantity"`
		} `json:"items"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	items := make([]logistics_grpc_adapter.ReturnItemInput, 0, len(body.Items))
	for _, it := range body.Items {
		items = append(items, logistics_grpc_adapter.ReturnItemInput{
			SKUID:    it.SKUID,
			Quantity: it.Quantity,
		})
	}

	result, err := s.svc.RecordReturn(ctx, &logistics_grpc_adapter.RecordReturnParams{
		BookingID: body.BookingID,
		Reason:    body.Reason,
		Items:     items,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"return_id": result.ReturnID}})
}

// ---------------------------------------------------------------------------
// BL-LOG-027: ProcessExchange — POST /v1/logistics/exchanges
// ---------------------------------------------------------------------------

func (s *Server) ProcessExchange(c *fiber.Ctx) error {
	const op = "rest_oapi.Server.ProcessExchange"
	ctx, span := s.tracer.Start(c.UserContext(), op)
	defer span.End()
	logging.LogWithTrace(ctx, s.logger).Info().Str("op", op).Msg("")

	var body struct {
		BookingID string `json:"booking_id"`
		Items     []struct {
			OldSKUID string `json:"old_sku_id"`
			NewSKUID string `json:"new_sku_id"`
			Quantity int32  `json:"quantity"`
		} `json:"items"`
		Notes string `json:"notes"`
	}
	if err := c.BodyParser(&body); err != nil {
		return writeIamAdminError(c, span, errors.Join(apperrors.ErrValidation, err))
	}

	items := make([]logistics_grpc_adapter.ExchangeItemInput, 0, len(body.Items))
	for _, it := range body.Items {
		items = append(items, logistics_grpc_adapter.ExchangeItemInput{
			OldSKUID: it.OldSKUID,
			NewSKUID: it.NewSKUID,
			Quantity: it.Quantity,
		})
	}

	result, err := s.svc.ProcessExchange(ctx, &logistics_grpc_adapter.ProcessExchangeParams{
		BookingID: body.BookingID,
		Items:     items,
		Notes:     body.Notes,
	})
	if err != nil {
		return writeIamAdminError(c, span, err)
	}
	span.SetStatus(codes.Ok, "ok")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": fiber.Map{"exchange_id": result.ExchangeID}})
}
