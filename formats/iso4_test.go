package formats

import (
	"testing"

	"github.com/moov-io/pinblock/encryption"
	"github.com/stretchr/testify/require"
)

func TestISO4(t *testing.T) {
	account := []byte("432198765432109870")

	t.Run("Encode/Decode with AES encryption", func(t *testing.T) {
		cipher, err := encryption.NewAesECB([]byte("1234567890123456"))
		require.NoError(t, err)
		iso4 := NewISO4(cipher)

		pinBlock, err := iso4.Encode([]byte("12344"), account)
		require.NoError(t, err)
		require.Len(t, pinBlock, 16)

		pin, err := iso4.Decode(pinBlock, account)
		require.NoError(t, err)
		require.Equal(t, []byte("12344"), pin)
	})

	t.Run("Encode/Decode with NoOp encryption", func(t *testing.T) {
		iso4 := NewISO4(encryption.NewNoOp())

		pinBlock, err := iso4.Encode([]byte("1234"), account)
		require.NoError(t, err)
		require.Len(t, pinBlock, 16)

		// compare only part of the pin block because the rest is random
		require.Equal(t, []byte{0x20, 0x20, 0x2D, 0x2D, 0xCF, 0xE9, 0x8B, 0xA3}, pinBlock[:8])

		pin, err := iso4.Decode(pinBlock, account)
		require.NoError(t, err)
		require.Equal(t, []byte("1234"), pin)
	})

	t.Run("wrong key is refused", func(t *testing.T) {
		c1, err := encryption.NewAesECB([]byte("1234567890123456"))
		require.NoError(t, err)
		c2, err := encryption.NewAesECB([]byte("6543210987654321"))
		require.NoError(t, err)

		pinBlock, err := NewISO4(c1).Encode([]byte("1234"), account)
		require.NoError(t, err)
		_, err = NewISO4(c2).Decode(pinBlock, account)
		require.Error(t, err)
	})

	t.Run("bad lengths", func(t *testing.T) {
		iso4 := NewISO4(encryption.NewNoOp())
		_, err := iso4.Encode([]byte("123"), account)
		require.Error(t, err)
		_, err = iso4.Encode([]byte("1234"), []byte("12345678901"))
		require.Error(t, err)
		_, err = iso4.Encode([]byte("1234"), []byte("12345678901234567890"))
		require.Error(t, err)
		_, err = iso4.Decode(make([]byte, 15), account)
		require.Error(t, err)
	})
}

// shortCipher returns the wrong length, which used to reach a slice index or
// an XOR length mismatch; now it is a plain error.
type shortCipher struct{}

func (shortCipher) Encrypt(p []byte) ([]byte, error) { return p[:8], nil }
func (shortCipher) Decrypt(c []byte) ([]byte, error) { return c[:8], nil }

func TestISO4_CipherReturnsWrongLength(t *testing.T) {
	account := []byte("432198765432109870")
	_, err := NewISO4(shortCipher{}).Encode([]byte("1234"), account)
	require.ErrorContains(t, err, "cipher returned 8 bytes")
	_, err = NewISO4(shortCipher{}).Decode(make([]byte, 16), account)
	require.ErrorContains(t, err, "cipher returned 8 bytes")
}
