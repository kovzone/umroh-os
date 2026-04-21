package service

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func Test_generateRefreshToken_distinctHashes(t *testing.T) {
	seen := map[string]bool{}
	for range 16 {
		plain, hashed, err := generateRefreshToken()
		require.NoError(t, err)

		// Plaintext is 2× hex chars of the 32-byte random = 64 chars.
		require.Len(t, plain, 64)
		require.True(t, strings.IndexFunc(plain, func(r rune) bool {
			return !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f')
		}) == -1, "plaintext must be lowercase hex")

		// Hash must round-trip via hashRefreshToken.
		require.Equal(t, hashed, hashRefreshToken(plain))

		// And match an independent sha256 of the plaintext.
		sum := sha256.Sum256([]byte(plain))
		require.Equal(t, hex.EncodeToString(sum[:]), hashed)

		// Every token in the loop must be unique (rand.Read drives uniqueness).
		require.False(t, seen[plain], "refresh token collision: %s", plain)
		seen[plain] = true
	}
}

func Test_hashRefreshToken_deterministic(t *testing.T) {
	plain := "abc123"
	require.Equal(t, hashRefreshToken(plain), hashRefreshToken(plain))
	require.NotEqual(t, hashRefreshToken(plain), hashRefreshToken("abc124"))
}

func Test_uuidToString_validAndInvalid(t *testing.T) {
	var zero pgtype.UUID
	require.Equal(t, "", uuidToString(zero), "invalid pgtype.UUID returns empty string")

	canonical := uuid.New()
	valid := pgtype.UUID{Bytes: canonical, Valid: true}
	require.Equal(t, canonical.String(), uuidToString(valid))
}

func Test_stringToUUID_parsesAndRejects(t *testing.T) {
	id := uuid.New()
	parsed, err := stringToUUID(id.String())
	require.NoError(t, err)
	require.True(t, parsed.Valid)
	require.Equal(t, id, uuid.UUID(parsed.Bytes))

	_, err = stringToUUID("not-a-uuid")
	require.Error(t, err, "garbage input must not silently return zero value")
}
