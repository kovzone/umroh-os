// ops_models.go — hand-written sqlc model types for ops schema tables.
// These complement the generated models.go for the new ops.* tables.

package sqlc

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// RoomAllocation mirrors ops.room_allocations.
type RoomAllocation struct {
	ID          string             `json:"id"`
	DepartureID string             `json:"departure_id"`
	Status      string             `json:"status"`
	CommittedAt pgtype.Timestamptz `json:"committed_at"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

// RoomAssignment mirrors ops.room_assignments.
type RoomAssignment struct {
	ID           string             `json:"id"`
	AllocationID string             `json:"allocation_id"`
	RoomNumber   string             `json:"room_number"`
	JamaahID     string             `json:"jamaah_id"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
}

// IDCardIssuance mirrors ops.id_card_issuances.
type IDCardIssuance struct {
	ID          string             `json:"id"`
	JamaahID    string             `json:"jamaah_id"`
	DepartureID string             `json:"departure_id"`
	CardType    string             `json:"card_type"`
	Token       string             `json:"token"`
	IssuedAt    pgtype.Timestamptz `json:"issued_at"`
}
