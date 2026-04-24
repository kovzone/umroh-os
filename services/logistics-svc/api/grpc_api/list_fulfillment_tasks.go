// list_fulfillment_tasks.go — gRPC handler for ListFulfillmentTasks (ISSUE-018).
//
// Route: /pb.LogisticsService/ListFulfillmentTasks
// Auth:  bearer (enforced at gateway); logistics-svc does not re-validate.

package grpc_api

import (
	"context"

	"logistics-svc/api/grpc_api/pb"
	"logistics-svc/service"
	"logistics-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/status"
	grpcCodes "google.golang.org/grpc/codes"
)

// ListFulfillmentTasks returns a paginated, optionally filtered list of fulfillment tasks.
func (s *Server) ListFulfillmentTasks(ctx context.Context, req *pb.ListFulfillmentTasksRequest) (*pb.ListFulfillmentTasksResponse, error) {
	const op = "grpc_api.Server.ListFulfillmentTasks"

	ctx, span := s.tracer.Start(ctx, op)
	defer span.End()

	logger := logging.LogWithTrace(ctx, s.logger)
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("status_filter", req.GetStatusFilter()),
		attribute.String("departure_id_filter", req.GetDepartureIdFilter()),
	)
	logger.Info().Str("op", op).Msg("")

	result, err := s.svc.ListFulfillmentTasks(ctx, &service.ListFulfillmentTasksParams{
		StatusFilter:      req.GetStatusFilter(),
		DepartureIDFilter: req.GetDepartureIdFilter(),
		Limit:             req.GetLimit(),
		Offset:            req.GetOffset(),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		logger.Error().Err(err).Str("op", op).Msg("")
		return nil, status.Errorf(grpcCodes.Internal, "list fulfillment tasks failed: %s", err.Error())
	}

	tasks := make([]*pb.FulfillmentTaskProto, 0, len(result.Tasks))
	for _, t := range result.Tasks {
		tasks = append(tasks, &pb.FulfillmentTaskProto{
			Id:             t.ID,
			BookingId:      t.BookingID,
			DepartureId:    t.DepartureID,
			Status:         t.Status,
			TrackingNumber: t.TrackingNumber,
			ShippedAt:      t.ShippedAt,
			DeliveredAt:    t.DeliveredAt,
			CreatedAt:      t.CreatedAt,
			UpdatedAt:      t.UpdatedAt,
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &pb.ListFulfillmentTasksResponse{
		Tasks: tasks,
		Total: result.Total,
	}, nil
}
