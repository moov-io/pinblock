package formats_test

import (
	"testing"

	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
)

func TestDecode_InvalidPinLengthNoPanic(t *testing.T) {
	// Previously panicked on ISO-0 when the pin length nibble exceeded the block
	require.NotPanics(t, func() {
		_, _ = formats.NewISO0().Decode(block(t, "0A00000000000000"), []byte("0000000000000"))
	})
	require.NotPanics(t, func() {
		for _, ctor := range []func() formats.Format{
			formats.NewISO0, formats.NewISO1, formats.NewISO3,
			formats.NewECI1, formats.NewECI2, formats.NewECI3, formats.NewECI4,
			formats.NewVISA1, formats.NewVISA2, formats.NewVISA3, formats.NewVISA4,
			formats.NewOEM1, formats.NewANSIX98,
		} {
			pinBlock, err := ctor().Encode([]byte("00 0"), []byte("0"))
			require.Error(t, err)
			require.Nil(t, pinBlock)
			_, _ = ctor().Decode(block(t, "0A00000000000000"), []byte("0000000000000"))
			// Every length is checked, including empty and nil.
			for _, b := range [][]byte{nil, {}, make([]byte, 7), make([]byte, 9), make([]byte, 16)} {
				_, _ = ctor().Decode(b, []byte("0"))
			}
		}
	})
}
