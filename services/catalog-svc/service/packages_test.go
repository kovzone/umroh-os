package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Cursor helpers are the pure-function boundary between the public wire
// contract ("opaque base64 blob") and the internal `last_id` used by
// the SQL `WHERE p.id > $cursor` clause. Exercising them in isolation
// guards against accidental wire-format drift — renaming the JSON key
// or changing the base64 alphabet would break every paginating caller.
// End-to-end cursor behavior through the DB is covered by the e2e
// spec `tests/e2e/tests/02e-catalog-svc-read.spec.ts`.

func TestCursorRoundTrip(t *testing.T) {
	cases := []string{
		"pkg_01JCDE00000000000000000001",
		"pkg_01JCDZZZZZZZZZZZZZZZZZZZZZ", // max-lex ULID-ish
		"pkg_with-_non-alphanumeric.chars",
	}
	for _, id := range cases {
		t.Run(id, func(t *testing.T) {
			encoded := encodeCursor(id)
			require.NotEmpty(t, encoded)
			require.NotContains(t, encoded, id, "cursor must not leak the id verbatim")

			decoded, err := decodeCursor(encoded)
			require.NoError(t, err)
			require.Equal(t, id, decoded)
		})
	}
}

func TestDecodeEmptyCursor(t *testing.T) {
	id, err := decodeCursor("")
	require.NoError(t, err)
	require.Equal(t, "", id, "empty cursor means first-page request")
}

func TestDecodeMalformedCursor(t *testing.T) {
	cases := []string{
		"not-base64!!",
		"Z" + strings.Repeat("A", 10), // base64 but not valid JSON payload
	}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			_, err := decodeCursor(c)
			require.Error(t, err, "malformed cursor must error so the handler returns invalid_cursor 400")
		})
	}
}

func TestDecodeCursorMissingLastID(t *testing.T) {
	// Valid base64 + valid JSON but missing the `last_id` field.
	encoded := encodeCursor("") // encodes {"last_id":""}
	_, err := decodeCursor(encoded)
	require.Error(t, err, "cursor payload missing last_id must fail")
}
