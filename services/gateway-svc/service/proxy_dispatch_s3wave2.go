// proxy_dispatch_s3wave2.go — gateway service dispatch for S3 Wave 2 new routes.
//
// Covers:
//   - ops-svc: RunRoomAllocation, GetRoomAllocation, GenerateIDCard, VerifyIDCard,
//     ExportManifest
//   - logistics-svc: ShipFulfillmentTask, GeneratePickupQR, RedeemPickupQR
//   - finance-svc: OnGRNReceived
//   - jamaah-svc: TriggerOCR, GetOCRStatus
//
// Each method is a thin delegation to the appropriate adapter.
// No business logic lives here; all logic lives in the backend services.
package service

import (
	"context"

	"gateway-svc/adapter/finance_grpc_adapter"
	"gateway-svc/adapter/jamaah_grpc_adapter"
	"gateway-svc/adapter/logistics_grpc_adapter"
	"gateway-svc/adapter/ops_grpc_adapter"
)

// ---------------------------------------------------------------------------
// ops-svc — Room allocation
// ---------------------------------------------------------------------------

func (s *Service) RunRoomAllocation(ctx context.Context, departureID string, jamaahIDs []string) (*ops_grpc_adapter.RunRoomAllocationResult, error) {
	return s.adapters.opsGrpc.RunRoomAllocation(ctx, departureID, jamaahIDs)
}

func (s *Service) GetRoomAllocation(ctx context.Context, departureID string) (*ops_grpc_adapter.GetRoomAllocationResult, error) {
	return s.adapters.opsGrpc.GetRoomAllocation(ctx, departureID)
}

// ---------------------------------------------------------------------------
// ops-svc — ID cards
// ---------------------------------------------------------------------------

func (s *Service) GenerateIDCard(ctx context.Context, jamaahID, departureID, cardType, jamaahName, departureName string) (*ops_grpc_adapter.GenerateIDCardResult, error) {
	return s.adapters.opsGrpc.GenerateIDCard(ctx, jamaahID, departureID, cardType, jamaahName, departureName)
}

func (s *Service) VerifyIDCard(ctx context.Context, token string) (*ops_grpc_adapter.VerifyIDCardResult, error) {
	return s.adapters.opsGrpc.VerifyIDCard(ctx, token)
}

// ---------------------------------------------------------------------------
// ops-svc — Manifest export
// ---------------------------------------------------------------------------

func (s *Service) ExportManifest(ctx context.Context, departureID string) (*ops_grpc_adapter.ExportManifestResult, error) {
	return s.adapters.opsGrpc.ExportManifest(ctx, departureID)
}

// ---------------------------------------------------------------------------
// logistics-svc — Shipment + pickup QR
// ---------------------------------------------------------------------------

func (s *Service) ShipFulfillmentTask(ctx context.Context, bookingID, carrier, notes string) (*logistics_grpc_adapter.ShipFulfillmentTaskResult, error) {
	return s.adapters.logisticsGrpc.ShipFulfillmentTask(ctx, bookingID, carrier, notes)
}

func (s *Service) GeneratePickupQR(ctx context.Context, bookingID string) (*logistics_grpc_adapter.GeneratePickupQRResult, error) {
	return s.adapters.logisticsGrpc.GeneratePickupQR(ctx, bookingID)
}

func (s *Service) RedeemPickupQR(ctx context.Context, token string) (*logistics_grpc_adapter.RedeemPickupQRResult, error) {
	return s.adapters.logisticsGrpc.RedeemPickupQR(ctx, token)
}

func (s *Service) ListFulfillmentTasks(ctx context.Context, params *logistics_grpc_adapter.ListFulfillmentTasksParams) (*logistics_grpc_adapter.ListFulfillmentTasksResult, error) {
	return s.adapters.logisticsGrpc.ListFulfillmentTasks(ctx, params)
}

// ---------------------------------------------------------------------------
// finance-svc — GRN
// ---------------------------------------------------------------------------

func (s *Service) OnGRNReceived(ctx context.Context, grnID, departureID string, amountIdr int64) (*finance_grpc_adapter.OnGRNReceivedResult, error) {
	return s.adapters.financeGrpc.OnGRNReceived(ctx, grnID, departureID, amountIdr)
}

// ---------------------------------------------------------------------------
// finance-svc — CorrectJournal (BL-FIN-006)
// ---------------------------------------------------------------------------

func (s *Service) CorrectJournal(ctx context.Context, params *finance_grpc_adapter.CorrectJournalParams) (*finance_grpc_adapter.CorrectJournalResult, error) {
	return s.adapters.financeGrpc.CorrectJournal(ctx, params)
}

// ---------------------------------------------------------------------------
// jamaah-svc — OCR
// ---------------------------------------------------------------------------

func (s *Service) TriggerOCR(ctx context.Context, documentID string) (*jamaah_grpc_adapter.TriggerOCRResult, error) {
	return s.adapters.jamaahGrpc.TriggerOCR(ctx, documentID)
}

func (s *Service) GetOCRStatus(ctx context.Context, documentID string) (*jamaah_grpc_adapter.GetOCRStatusResult, error) {
	return s.adapters.jamaahGrpc.GetOCRStatus(ctx, documentID)
}
