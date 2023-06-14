package encryption

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAesECB(t *testing.T) {
	t.Run("Encode/Decode", func(t *testing.T) {
		cipher, err := NewAesECB([]byte("1234567890123456"))
		require.NoError(t, err)

		cipherText, err := cipher.Encrypt([]byte("1234567890123456"))

		require.NoError(t, err)
		require.Len(t, cipherText, 16)

		plainText, err := cipher.Decrypt(cipherText)

		require.NoError(t, err)
		require.Equal(t, "1234567890123456", string(plainText))
	})

	t.Run("Encrypt/Decrypt with wrong value", func(t *testing.T) {
		cipher, err := NewAesECB([]byte("1234567890123456"))
		require.NoError(t, err)

		// encrypt

		// longer than 16 bytes
		_, err = cipher.Encrypt([]byte("12345678901234561234567890123456"))
		require.EqualError(t, err, "plain text length must be 16 bytes")

		// shorter than 16 bytes
		_, err = cipher.Encrypt([]byte("123456789012345"))
		require.EqualError(t, err, "plain text length must be 16 bytes")

		// decrypt

		// longer than 16 bytes
		_, err = cipher.Decrypt([]byte("12345678901234561234567890123456"))
		require.EqualError(t, err, "cipher text length must be 16 bytes")

		// shorter than 16 bytes
		_, err = cipher.Decrypt([]byte("123456789012345"))
		require.EqualError(t, err, "cipher text length must be 16 bytes")
	})
}
