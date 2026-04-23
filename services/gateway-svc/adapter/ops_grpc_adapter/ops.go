// ops.go — gateway adapter methods for ops-svc RPCs (S3 Wave 2).
//
// Each method translates gateway-local params → pb request, forwards via gRPC,
// and translates pb response → adapter-local types. Proto types do not leak
// past this package.
package ops_grpc_adapter

import (
	"context"

	"gateway-svc/adapter/ops_grpc_adapter/pb"
	"gateway-svc/util/logging"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ---------------------------------------------------------------------------
// Adapter-local result types
// ---------------------------------------------------------------------------

// RoomAssignmentResult is the gateway-local representation of one room assignment.
type RoomAssignmentResult struct {
	RoomNumber string
	JamaahID   string
}

// RunRoomAllocationResult is the gateway-local result for RunRoomAllocation.
type RunRoomAllocationResult struct {
	AllocationID string
	RoomCount    int32
	Assignments  []*RoomAssignmentResult
}

// GetRoomAllocationResult is the gateway-local result for GetRoomAllocation.
type GetRoomAllocationResult struct {
	AllocationID string
	RoomCount    int32
	Assignments  []*RoomAssignmentResult
	Status       string
}

// GenerateIDCardResult is the gateway-local result for GenerateIDCard.
type GenerateIDCardResult struct {
	Token    string
	QRData   string
	IssuedAt string
}

// VerifyIDCardResult is the gateway-local result for VerifyIDCard.
type VerifyIDCardResult struct {
	Valid        bool
	JamaahID     string
	DepartureID  string
	CardType     string
	ErrorReason  string
}

// ManifestExportRowResult is one row in the manifest export.
type ManifestExportRowResult struct {
	No          int32
	JamaahName  string
	PassportNo  string
	DocStatus   string
	RoomNumber  string
}

// ExportManifestResult is the gateway-local result for ExportManifest.
type ExportManifestResult struct {
	DepartureID string
	Rows        []*ManifestExportRowResult
}

// ---------------------------------------------------------------------------
// RunRoomAllocation
// ---------------------------------------------------------------------------

// RunRoomAllocation triggers a room allocation run for a departure.
func (a *Adapter) RunRoomAllocation(ctx context.Context, departureID string, jamaahIDs []string) (*RunRoomAllocationResult, error) {
	const op = "ops_grpc_adapter.Adapter.RunRoomAllocation"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "RunRoomAllocation"),
		attribute.String("departure_id", departureID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.opsClient.RunRoomAllocation(ctx, &pb.RunRoomAllocationRequest{
		DepartureId: departureID,
		JamaahIds:   jamaahIDs,
	})
	if err != nil {
		wrapped := mapOpsError(err)
		logger.Warn().Err(wrapped).Str("departure_id", departureID).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	assignments := make([]*RoomAssignmentResult, 0, len(resp.GetAssignments()))
	for _, ra := range resp.GetAssignments() {
		assignments = append(assignments, &RoomAssignmentResult{
			RoomNumber: ra.GetRoomNumber(),
			JamaahID:   ra.GetJamaahId(),
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &RunRoomAllocationResult{
		AllocationID: resp.GetAllocationId(),
		RoomCount:    resp.GetRoomCount(),
		Assignments:  assignments,
	}, nil
}

// ---------------------------------------------------------------------------
// GetRoomAllocation
// ---------------------------------------------------------------------------

// GetRoomAllocation fetches the current room allocation for a departure.
func (a *Adapter) GetRoomAllocation(ctx context.Context, departureID string) (*GetRoomAllocationResult, error) {
	const op = "ops_grpc_adapter.Adapter.GetRoomAllocation"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "GetRoomAllocation"),
		attribute.String("departure_id", departureID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.opsClient.GetRoomAllocation(ctx, &pb.GetRoomAllocationRequest{
		DepartureId: departureID,
	})
	if err != nil {
		wrapped := mapOpsError(err)
		logger.Warn().Err(wrapped).Str("departure_id", departureID).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	assignments := make([]*RoomAssignmentResult, 0, len(resp.GetAssignments()))
	for _, ra := range resp.GetAssignments() {
		assignments = append(assignments, &RoomAssignmentResult{
			RoomNumber: ra.GetRoomNumber(),
			JamaahID:   ra.GetJamaahId(),
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &GetRoomAllocationResult{
		AllocationID: resp.GetAllocationId(),
		RoomCount:    resp.GetRoomCount(),
		Assignments:  assignments,
		Status:       resp.GetStatus(),
	}, nil
}

// ---------------------------------------------------------------------------
// GenerateIDCard
// ---------------------------------------------------------------------------

// GenerateIDCard creates (or returns existing) HMAC-signed ID card token.
func (a *Adapter) GenerateIDCard(ctx context.Context, jamaahID, departureID, cardType, jamaahName, departureName string) (*GenerateIDCardResult, error) {
	const op = "ops_grpc_adapter.Adapter.GenerateIDCard"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "GenerateIDCard"),
		attribute.String("jamaah_id", jamaahID),
		attribute.String("departure_id", departureID),
		attribute.String("card_type", cardType),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.opsClient.GenerateIDCard(ctx, &pb.GenerateIDCardRequest{
		JamaahId:      jamaahID,
		DepartureId:   departureID,
		CardType:      cardType,
		JamaahName:    jamaahName,
		DepartureName: departureName,
	})
	if err != nil {
		wrapped := mapOpsError(err)
		logger.Warn().Err(wrapped).Str("jamaah_id", jamaahID).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &GenerateIDCardResult{
		Token:    resp.GetToken(),
		QRData:   resp.GetQrData(),
		IssuedAt: resp.GetIssuedAt(),
	}, nil
}

// ---------------------------------------------------------------------------
// VerifyIDCard
// ---------------------------------------------------------------------------

// VerifyIDCard verifies a previously issued ID card token.
func (a *Adapter) VerifyIDCard(ctx context.Context, token string) (*VerifyIDCardResult, error) {
	const op = "ops_grpc_adapter.Adapter.VerifyIDCard"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "VerifyIDCard"),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.opsClient.VerifyIDCard(ctx, &pb.VerifyIDCardRequest{
		Token: token,
	})
	if err != nil {
		wrapped := mapOpsError(err)
		logger.Warn().Err(wrapped).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	span.SetStatus(codes.Ok, "ok")
	return &VerifyIDCardResult{
		Valid:       resp.GetValid(),
		JamaahID:    resp.GetJamaahId(),
		DepartureID: resp.GetDepartureId(),
		CardType:    resp.GetCardType(),
		ErrorReason: resp.GetErrorReason(),
	}, nil
}

// ---------------------------------------------------------------------------
// ExportManifest
// ---------------------------------------------------------------------------

// ExportManifest requests a manifest export for a departure.
func (a *Adapter) ExportManifest(ctx context.Context, departureID string) (*ExportManifestResult, error) {
	const op = "ops_grpc_adapter.Adapter.ExportManifest"

	ctx, span := a.tracer.Start(ctx, op)
	defer span.End()
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("rpc", "ExportManifest"),
		attribute.String("departure_id", departureID),
	)

	logger := logging.LogWithTrace(ctx, a.logger)

	resp, err := a.opsClient.ExportManifest(ctx, &pb.ExportManifestRequest{
		DepartureId: departureID,
	})
	if err != nil {
		wrapped := mapOpsError(err)
		logger.Warn().Err(wrapped).Str("departure_id", departureID).Msg("")
		span.RecordError(wrapped)
		span.SetStatus(codes.Error, wrapped.Error())
		return nil, wrapped
	}

	rows := make([]*ManifestExportRowResult, 0, len(resp.GetRows()))
	for _, r := range resp.GetRows() {
		rows = append(rows, &ManifestExportRowResult{
			No:         r.GetNo(),
			JamaahName: r.GetJamaahName(),
			PassportNo: r.GetPassportNo(),
			DocStatus:  r.GetDocStatus(),
			RoomNumber: r.GetRoomNumber(),
		})
	}

	span.SetStatus(codes.Ok, "ok")
	return &ExportManifestResult{
		DepartureID: resp.GetDepartureId(),
		Rows:        rows,
	}, nil
}
