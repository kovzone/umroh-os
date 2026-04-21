package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// 32 bytes keeps encryptAESGCM at AES-256.
var testAESKey = []byte("iam_svc_test_totp_aes_256_key_00")

func Test_encryptDecryptAESGCM_roundTrip(t *testing.T) {
	require.Len(t, testAESKey, 32, "test key must be 32 bytes")

	plaintext := []byte("JBSWY3DPEHPK3PXP") // typical base32 TOTP secret
	ct, err := encryptAESGCM(testAESKey, plaintext)
	require.NoError(t, err)
	require.NotContains(t, ct, string(plaintext), "ciphertext must not leak plaintext")

	pt, err := decryptAESGCM(testAESKey, ct)
	require.NoError(t, err)
	require.Equal(t, plaintext, pt)
}

func Test_encryptAESGCM_nonDeterministic(t *testing.T) {
	a, err := encryptAESGCM(testAESKey, []byte("same plaintext"))
	require.NoError(t, err)
	b, err := encryptAESGCM(testAESKey, []byte("same plaintext"))
	require.NoError(t, err)
	require.NotEqual(t, a, b, "fresh nonce per call must make ciphertext unique")
}

func Test_decryptAESGCM_rejectsWrongKey(t *testing.T) {
	ct, err := encryptAESGCM(testAESKey, []byte("hello"))
	require.NoError(t, err)

	badKey := []byte("iam_svc_wrong_aes_256_key_length")
	require.Len(t, badKey, 32)
	_, err = decryptAESGCM(badKey, ct)
	require.Error(t, err, "GCM auth tag must fail on wrong key")
}

func Test_decryptAESGCM_rejectsMalformed(t *testing.T) {
	for _, input := range []string{
		"no-colon-at-all",
		"deadbeef:", // empty ciphertext
		":deadbeef", // empty nonce
		"zz:zz",     // non-hex
	} {
		_, err := decryptAESGCM(testAESKey, input)
		require.Error(t, err, "input %q must fail", input)
	}
}
