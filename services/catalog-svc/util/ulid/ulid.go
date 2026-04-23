// Package ulid provides a minimal ULID-like ID generator for catalog-svc.
//
// Generated IDs are time-sortable, URL-safe, and unique enough for a
// single-instance MVP. Format: prefix + 26-char Crockford base32 string
// (monotonic timestamp portion + random portion).
//
// The exact ULID spec (oklog/ulid) is not imported as a module dependency
// in this service; if the project later adds the canonical library, this
// package can be swapped to it with a one-line change.
package ulid

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// Crockford base32 alphabet (no I, L, O, U to avoid visual ambiguity).
const crockfordAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// New returns a new ULID-like ID with the given prefix (e.g. "pkg_").
// The returned string has the format: <prefix><10-char-timestamp><16-char-random>.
// Total non-prefix length is 26 chars (matching ULID spec).
func New(prefix string) (string, error) {
	// 10 chars of Crockford-encoded timestamp (ms precision).
	ts := uint64(time.Now().UnixMilli())
	tsPart := encodeCrockford(ts, 10)

	// 16 chars of random data.
	randPart, err := randomCrockford(16)
	if err != nil {
		return "", fmt.Errorf("ulid: generate random: %w", err)
	}

	return prefix + tsPart + randPart, nil
}

// MustNew is like New but panics on error. Use only in init/test code.
func MustNew(prefix string) string {
	id, err := New(prefix)
	if err != nil {
		panic(err)
	}
	return id
}

// encodeCrockford encodes val into n Crockford base32 characters (zero-padded).
func encodeCrockford(val uint64, n int) string {
	chars := make([]byte, n)
	for i := n - 1; i >= 0; i-- {
		chars[i] = crockfordAlphabet[val&0x1F]
		val >>= 5
	}
	return string(chars)
}

// randomCrockford returns n Crockford base32 characters of cryptographic randomness.
func randomCrockford(n int) (string, error) {
	var sb strings.Builder
	sb.Grow(n)

	// Each character needs 5 bits; generate enough random bytes.
	// We use crypto/rand.Int to generate uniformly in [0, 32^n).
	max := new(big.Int).Exp(big.NewInt(32), big.NewInt(int64(n)), nil)
	r, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Extract n base32 digits from r.
	digits := make([]byte, n)
	thirty2 := big.NewInt(32)
	mod := new(big.Int)
	for i := n - 1; i >= 0; i-- {
		r.DivMod(r, thirty2, mod)
		digits[i] = crockfordAlphabet[mod.Int64()]
	}
	sb.Write(digits)
	_ = binary.BigEndian // imported for potential future use

	return sb.String(), nil
}
