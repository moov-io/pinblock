package formats

import (
	"bytes"
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

		expectedOutput := `PIN block encode operation finished
************************************
PAN     : 432198765432109870
PIN     : 12344
PAD     : AAAAAAAAA
Format  : Format 4 (ISO-4)
------------------------------------
PAN block            : 64321987654321098700000000000000`
		require.Equal(t, true, strings.Contains(out.String(), expectedOutput))

		// flash buffer
		out.Reset()

		// Decode
		pin, err := iso4.Decode(pinBlock, "432198765432109870")

		require.NoError(t, err)
		require.Equal(t, "12344", pin)

		expectedOutput = `PIN block decode operation finished`
		require.Equal(t, true, strings.Contains(out.String(), expectedOutput))
	})
}
