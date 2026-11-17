package formats

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/moov-io/pinblock/encryption"
	"github.com/stretchr/testify/require"
)

func TestISO4(t *testing.T) {
	t.Run("Encode/Decode with AES encryption", func(t *testing.T) {
		cipher, err := encryption.NewAesECB([]byte("1234567890123456"))
		require.NoError(t, err)

		iso4 := NewISO4(cipher)

		// Encode
		pinBlock, err := iso4.Encode("12344", "432198765432109870")

		require.NoError(t, err)
		require.Len(t, pinBlock, 32)

		// Decode
		pin, err := iso4.Decode(pinBlock, "432198765432109870")

		require.NoError(t, err)
		require.Equal(t, "12344", pin)
	})

	t.Run("Encode/Decode with NoOp encryption", func(t *testing.T) {
		cipher := encryption.NewNoOp()
		iso4 := NewISO4(cipher)

		pinBlock, err := iso4.Encode("1234", "432198765432109870")

		require.NoError(t, err)
		require.Len(t, pinBlock, 32)

		// compare only part of the pin block because the rest is random
		require.Equal(t, "20202D2DCFE98BA3", string(pinBlock[0:16]))

		pin, err := iso4.Decode(pinBlock, "432198765432109870")

		require.NoError(t, err)
		require.Equal(t, "1234", pin)
	})

	t.Run("encode/decode logs", func(t *testing.T) {
		cipher, err := encryption.NewAesECB([]byte("1234567890123456"))
		require.NoError(t, err)

		iso4 := NewISO4(cipher)
		out := bytes.NewBuffer([]byte{})
		iso4.SetDebugWriter(out)

		// Encode
		pinBlock, err := iso4.Encode("12344", "432198765432109870")

		require.NoError(t, err)
		require.Len(t, pinBlock, 32)

		requireDebugNote(t, out.String(), "encode", "Format 4 (ISO-4)")

		// flash buffer
		out.Reset()

		// Decode
		pin, err := iso4.Decode(pinBlock, "432198765432109870")

		require.NoError(t, err)
		require.Equal(t, "12344", pin)

		requireDebugNote(t, out.String(), "decode", "Format 4 (ISO-4)")
	})
}

func TestISO4_Validation(t *testing.T) {
	account := "432198765432109870"

	t.Run("wrong key is refused", func(t *testing.T) {
		c1, err := encryption.NewAesECB([]byte("1234567890123456"))
		require.NoError(t, err)
		c2, err := encryption.NewAesECB([]byte("6543210987654321"))
		require.NoError(t, err)

		pinBlock, err := NewISO4(c1).Encode("1234", account)
		require.NoError(t, err)
		_, err = NewISO4(c2).Decode(pinBlock, account)
		require.Error(t, err)
	})

	t.Run("bad lengths", func(t *testing.T) {
		iso4 := NewISO4(encryption.NewNoOp())
		_, err := iso4.Encode("123", account)
		require.Error(t, err)
		_, err = iso4.Encode("1234567890123", account)
		require.Error(t, err)
		_, err = iso4.Encode("1234", "12345678901")
		require.Error(t, err)
		_, err = iso4.Encode("1234", "12345678901234567890")
		require.Error(t, err)
		_, err = iso4.Decode(strings.Repeat("0", 30), account)
		require.Error(t, err)
	})

	t.Run("control nibble and fill", func(t *testing.T) {
		iso4 := NewISO4(encryption.NewNoOp())
		panField, err := iso4.(*iso4Object).panField(account)
		require.NoError(t, err)
		// With the no-op cipher the block is the PIN field XOR the PAN field.
		blockFor := func(pinField string) string {
			b := mustHex(t, pinField)
			xorInto(b, b, panField)
			return strings.ToUpper(hex.EncodeToString(b))
		}

		pin, err := iso4.Decode(blockFor("441234AAAAAAAAAA0123456789ABCDEF"), account)
		require.NoError(t, err)
		require.Equal(t, "1234", pin)

		_, err = iso4.Decode(blockFor("341234AAAAAAAAAA0123456789ABCDEF"), account)
		require.ErrorContains(t, err, "format is different")
		_, err = iso4.Decode(blockFor("441234AAAAAAAAAF0123456789ABCDEF"), account)
		require.ErrorContains(t, err, "invalid fill")
		_, err = iso4.Decode(blockFor("43123AAAAAAAAAAA0123456789ABCDEF"), account)
		require.ErrorContains(t, err, "invalid pin length")
	})

	// The PAN is validated before anything is encrypted.
	t.Run("PAN checked before encrypting", func(t *testing.T) {
		c := &countingCipher{}
		_, err := NewISO4(c).Encode("1234", "43219876543210987A")
		require.ErrorContains(t, err, "account must contain only decimal digits")
		_, err = NewISO4(c).Encode("1234", "1234")
		require.Error(t, err)
		require.Zero(t, c.calls)
	})

	t.Run("nil cipher", func(t *testing.T) {
		_, err := NewISO4(nil).Encode("1234", account)
		require.ErrorContains(t, err, "requires a cipher")
		_, err = NewISO4(nil).Decode(strings.Repeat("0", 32), account)
		require.ErrorContains(t, err, "requires a cipher")
	})
}

// countingCipher is a no-op cipher that counts its calls.
type countingCipher struct{ calls int }

func (c *countingCipher) Encrypt(p []byte) ([]byte, error) { c.calls++; return bytes.Clone(p), nil }
func (c *countingCipher) Decrypt(p []byte) ([]byte, error) { c.calls++; return bytes.Clone(p), nil }

// A cipher that returns the wrong length is an error, not a panic or a
// short block.
func TestISO4_CipherReturnsWrongLength(t *testing.T) {
	account := "432198765432109870"
	_, err := NewISO4(shortCipher{}).Encode("1234", account)
	require.ErrorContains(t, err, "cipher returned 8 bytes")
	_, err = NewISO4(shortCipher{}).Decode(strings.Repeat("0", 32), account)
	require.ErrorContains(t, err, "cipher returned 8 bytes")
}
