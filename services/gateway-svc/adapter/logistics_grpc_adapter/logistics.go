// logistics.go — gateway adapter methods for logistics-svc RPCs (S3 Wave 2).
//
// Each method translates gateway-local params → pb request, forwards via gRPC,
// and translates pb response → adapter-local types. Proto types do not leak
// past this package.
package logistics_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/logistics_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Adapter-local result types
// ---------------------------------------------------------------------------

// ShipFulfillmentTaskResult is the gateway-local result for ShipFulfillmentTask.
type ShipFulfillmentTaskResult struct {
	ShipmentID     string
	TrackingNumber string
	Status         string
}

// GeneratePickupQRResult is the gateway-local result for GeneratePickupQR.
type GeneratePickupQRResult struct {
	PickupTokenID string
	Token         string
	ExpiresAt     string
}

// RedeemPickupQRResult is the gateway-local result for RedeemPickupQR.
type RedeemPickupQRResult struct {
	Redeemed    bool
	BookingID   string
	TaskID      string
	ErrorReason string
}

// ---------------------------------------------------------------------------
// ShipFulfillmentTask
// ---------------------------------------------------------------------------

// ShipFulfillmentTask creates a shipment record for a booking kit.
func (a *Adapter) ShipFulfillmentTask(ctx context.Context, bookingID, carrier, notes string) (*ShipFulfillmentTaskResult, error) {
	const op = "logistics_grpc_adapter.Adapter.ShipFulfillmentTask"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "ShipFulfillmentTask"),
		attribute.String("booking_id", bookingID),
		attribute.String("carrier", carrier),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.logisticsClient.ShipFulfillmentTask(ctx, &pb.ShipFulfillmentTaskRequest{
		BookingId: bookingID,
		Carrier:   carrier,
		Notes:     notes,
	})
	if err != nil {
		wrapped := mapLogisticsError(err)
		logger.Warn().Err(wrapped).Str("booking_id", bookingID).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &ShipFulfillmentTaskResult{
		ShipmentID:     resp.GetShipmentId(),
		TrackingNumber: resp.GetTrackingNumber(),
		Status:         resp.GetStatus(),
	}, nil
}

// ---------------------------------------------------------------------------
// GeneratePickupQR
// ---------------------------------------------------------------------------

// GeneratePickupQR generates a single-use pickup QR token for a booking.
func (a *Adapter) GeneratePickupQR(ctx context.Context, bookingID string) (*GeneratePickupQRResult, error) {
	const op = "logistics_grpc_adapter.Adapter.GeneratePickupQR"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "GeneratePickupQR"),
		attribute.String("booking_id", bookingID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.logisticsClient.GeneratePickupQR(ctx, &pb.GeneratePickupQRRequest{
		BookingId: bookingID,
	})
	if err != nil {
		wrapped := mapLogisticsError(err)
		logger.Warn().Err(wrapped).Str("booking_id", bookingID).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &GeneratePickupQRResult{
		PickupTokenID: resp.GetPickupTokenId(),
		Token:         resp.GetToken(),
		ExpiresAt:     resp.GetExpiresAt(),
	}, nil
}

// ---------------------------------------------------------------------------
// RedeemPickupQR
// ---------------------------------------------------------------------------

// RedeemPickupQR redeems a single-use pickup QR token.
func (a *Adapter) RedeemPickupQR(ctx context.Context, token string) (*RedeemPickupQRResult, error) {
	const op = "logistics_grpc_adapter.Adapter.RedeemPickupQR"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "RedeemPickupQR"),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.logisticsClient.RedeemPickupQR(ctx, &pb.RedeemPickupQRRequest{
		Token: token,
	})
	if err != nil {
		wrapped := mapLogisticsError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &RedeemPickupQRResult{
		Redeemed:    resp.GetRedeemed(),
		BookingID:   resp.GetBookingId(),
		TaskID:      resp.GetTaskId(),
		ErrorReason: resp.GetErrorReason(),
	}, nil
}
