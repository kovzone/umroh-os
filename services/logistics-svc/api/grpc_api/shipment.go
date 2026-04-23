// shipment.go — gRPC handler for ShipFulfillmentTask RPC (BL-LOG-002).
//
// ShipFulfillmentTask: lookup task by booking_id, create shipment record,
// generate tracking number (format: "UOS-" + 8 upper hex chars of UUID),
// update fulfillment_tasks.status to 'shipped', log WA notification stub.
// Idempotent: if a shipment already exists for the task, return it.

package grpc_api

import (
	"context"

	"logistics-svc/api/grpc_api/pb"
	"logistics-svc/service"
	"logistics-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

// ShipFulfillmentTask handles the ShipFulfillmentTask RPC.
// Creates a shipment + tracking number for a paid booking.
// Idempotent: returns existing shipment if one already exists for this booking.
func (s *Server) ShipFulfillmentTask(ctx context.Context, req *pb.ShipFulfillmentTaskRequest) (*pb.ShipFulfillmentTaskResponse, error) {
	const op = "grpc_api.Server.ShipFulfillmentTask"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("booking_id", req.GetBookingId()),
		attribute.String("carrier", req.GetCarrier()),
	)

	logger.Info().
		Str("op", op).
		Str("booking_id", req.GetBookingId()).
		Msg("")

	if req.GetBookingId() == "" {
		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "booking_id is required")
	}

	result, err := s.svc.ShipFulfillmentTask(ctx, &service.ShipFulfillmentTaskParams{
		BookingID: req.GetBookingId(),
		Carrier:   req.GetCarrier(),
		Notes:     req.GetNotes(),
	})
	if err != nil {
		logger.Error().
			Str("op", op).
			Str("booking_id", req.GetBookingId()).
			Err(err).
			Msg("ShipFulfillmentTask failed")
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, grpcStatus.Errorf(grpcCodes.Internal, "failed to ship fulfillment task: %v", err)
	}

	logger.Info().
		Str("op", op).
		Str("booking_id", req.GetBookingId()).
		Str("shipment_id", result.ShipmentID).
		Str("tracking_number", result.TrackingNumber).
		Bool("replayed", result.Replayed).
		Msg("ShipFulfillmentTask succeeded")

	span.SetStatus(otelCodes.Ok, result.Status)
	return &pb.ShipFulfillmentTaskResponse{
		ShipmentId:     result.ShipmentID,
		TrackingNumber: result.TrackingNumber,
		Status:         result.Status,
	}, nil
}
