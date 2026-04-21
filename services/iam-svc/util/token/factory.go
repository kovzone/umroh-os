package token

import (
	"fmt"
	"strings"
	"time"
)

const (
	TypePaseto = "paseto"
	TypeJWT    = "jwt"
)

// Maker creates and verifies tokens. Only payload.go is project-specific; implementors use *Payload.
type Maker interface {
	CreateToken(payload *Payload, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

// NewMaker creates a Maker based on the given token type ("paseto" or "jwt").
// If tokenType is empty, defaults to PASETO.
// Key requirements: exactly 32 characters for PASETO (Symmetric Key), at least 32 for JWT (Secret Key).
func NewMaker(tokenType, key string) (Maker, error) {
	t := strings.ToLower(strings.TrimSpace(tokenType))
	if t == "" {
		t = TypePaseto
	}

	switch t {
	case TypePaseto:
		return NewPasetoMaker(key)
	case TypeJWT:
		return NewJWTMaker(key)
	default:
		return nil, fmt.Errorf("unsupported token type %q: must be %q or %q", tokenType, TypePaseto, TypeJWT)
	}
}
