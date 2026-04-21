package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload is the only type you need to customize when copying this package.
// BERES: user_id and current_branch_id (see user lifecycle doc); roles at that branch.
// Maker implementations set IssuedAt/ExpiredAt.
type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"`   // UUID string (iam.users.id)
	BranchID  string    `json:"branch_id"` // UUID string (current/default branch)
	Roles     []string  `json:"roles"`     // role names at BranchID
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload builds a payload with a new ID and the given identity fields.
// IssuedAt/ExpiredAt are set by Maker.CreateToken.
func NewPayload(userID, branchID string, roles []string) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Payload{
		ID:       tokenID,
		UserID:   userID,
		BranchID: branchID,
		Roles:    roles,
	}, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
