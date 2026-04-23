// ops_room.go — room allocation service-layer implementation (BL-OPS-002).
//
// Algorithm: sort jamaah_ids alphabetically, assign groups of 4 (room
// capacity) to room numbers "1", "2", etc.
//
// Idempotency:
//   - If allocation for departure_id is 'draft', delete existing assignments
//     and re-run.
//   - If 'committed', return error.

package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"ops-svc/store/postgres_store"
	"ops-svc/store/postgres_store/sqlc"
	"ops-svc/util/logging"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	otelCodes "go.opentelemetry.io/otel/codes"
)

const roomCapacity = 4

// RoomAssignment represents a single room assignment.
type RoomAssignment struct {
	RoomNumber string
	JamaahID   string
}

// RunRoomAllocationParams holds inputs for RunRoomAllocation.
type RunRoomAllocationParams struct {
	DepartureID string
	JamaahIDs   []string
}

// RunRoomAllocationResult holds the result of RunRoomAllocation.
type RunRoomAllocationResult struct {
	AllocationID string
	RoomCount    int32
	Assignments  []RoomAssignment
}

// RunRoomAllocation runs the room allocation algorithm for a departure.
// Idempotent on departure_id: re-runs if draft, errors if committed.
func (svc *Service) RunRoomAllocation(ctx context.Context, params *RunRoomAllocationParams) (*RunRoomAllocationResult, error) {
	const op = "service.Service.RunRoomAllocation"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("departure_id", params.DepartureID),
		attribute.Int("jamaah_count", len(params.JamaahIDs)),
	)

	logger := logging.LogWithTrace(ctx, svc.logger)

	if params.DepartureID == "" {
		return nil, fmt.Errorf("%s: departure_id is required", op)
	}
	if len(params.JamaahIDs) == 0 {
		return nil, fmt.Errorf("%s: jamaah_ids must not be empty", op)
	}

	// Sort jamaah_ids alphabetically for deterministic grouping.
	sorted := make([]string, len(params.JamaahIDs))
	copy(sorted, params.JamaahIDs)
	sort.Strings(sorted)

	// Build assignments: groups of roomCapacity.
	var assignments []RoomAssignment
	roomNum := 1
	for i, jid := range sorted {
		if i > 0 && i%roomCapacity == 0 {
			roomNum++
		}
		assignments = append(assignments, RoomAssignment{
			RoomNumber: strconv.Itoa(roomNum),
			JamaahID:   jid,
		})
	}
	roomCount := int32(roomNum)

	var result RunRoomAllocationResult

	_, err := svc.store.WithTx(ctx, &postgres_store.WithTxArgs{
		Fn: func(qtx *sqlc.Queries) error {
			// Check if allocation exists.
			existing, err := qtx.GetRoomAllocationByDepartureID(ctx, params.DepartureID)
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("get existing allocation: %w", err)
			}

			var allocationID string
			if err == nil {
				// Allocation exists.
				if existing.Status == "committed" {
					return fmt.Errorf("allocation for departure %s is already committed, cannot re-run", params.DepartureID)
				}
				// Delete existing assignments and re-run.
				if err := qtx.DeleteRoomAssignmentsByAllocationID(ctx, existing.ID); err != nil {
					return fmt.Errorf("delete old assignments: %w", err)
				}
				allocationID = existing.ID
			} else {
				// No allocation yet — create one.
				alloc, err := qtx.InsertRoomAllocation(ctx, params.DepartureID)
				if err != nil {
					return fmt.Errorf("insert room allocation: %w", err)
				}
				allocationID = alloc.ID
			}

			// Insert new assignments.
			for _, a := range assignments {
				if _, err := qtx.InsertRoomAssignment(ctx, sqlc.InsertRoomAssignmentParams{
					AllocationID: allocationID,
					RoomNumber:   a.RoomNumber,
					JamaahID:     a.JamaahID,
				}); err != nil {
					return fmt.Errorf("insert room assignment jamaah=%s: %w", a.JamaahID, err)
				}
			}

			result = RunRoomAllocationResult{
				AllocationID: allocationID,
				RoomCount:    roomCount,
				Assignments:  assignments,
			}
			return nil
		},
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: transaction failed: %w", op, err)
	}

	logger.Info().
		Str("op", op).
		Str("departure_id", params.DepartureID).
		Str("allocation_id", result.AllocationID).
		Int32("room_count", result.RoomCount).
		Int("jamaah_count", len(params.JamaahIDs)).
		Msg("room allocation run completed")

	span.SetStatus(otelCodes.Ok, "success")
	return &result, nil
}

// GetRoomAllocationParams holds inputs for GetRoomAllocation.
type GetRoomAllocationParams struct {
	DepartureID string
}

// GetRoomAllocationResult holds the result of GetRoomAllocation.
type GetRoomAllocationResult struct {
	AllocationID string
	RoomCount    int32
	Assignments  []RoomAssignment
	Status       string
}

// GetRoomAllocation returns the current room allocation for a departure.
func (svc *Service) GetRoomAllocation(ctx context.Context, params *GetRoomAllocationParams) (*GetRoomAllocationResult, error) {
	const op = "service.Service.GetRoomAllocation"

	ctx, span := svc.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("departure_id", params.DepartureID),
	)

	if params.DepartureID == "" {
		return nil, fmt.Errorf("%s: departure_id is required", op)
	}

	alloc, err := svc.store.GetRoomAllocationByDepartureID(ctx, params.DepartureID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: get allocation: %w", op, err)
	}

	rows, err := svc.store.ListRoomAssignmentsByAllocationID(ctx, alloc.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelCodes.Error, err.Error())
		return nil, fmt.Errorf("%s: list assignments: %w", op, err)
	}

	assignments := make([]RoomAssignment, 0, len(rows))
	roomSet := make(map[string]struct{})
	for _, r := range rows {
		assignments = append(assignments, RoomAssignment{
			RoomNumber: r.RoomNumber,
			JamaahID:   r.JamaahID,
		})
		roomSet[r.RoomNumber] = struct{}{}
	}

	span.SetStatus(otelCodes.Ok, "success")
	return &GetRoomAllocationResult{
		AllocationID: alloc.ID,
		RoomCount:    int32(len(roomSet)),
		Assignments:  assignments,
		Status:       alloc.Status,
	}, nil
}
