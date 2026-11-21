package formats_test

import (
	"testing"

	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
)

func TestDecode_InvalidPinLengthNoPanic(t *testing.T) {
	// Previously paniced on ISO-0 when pin length nibble exceeded block
	require.NotPanics(t, func() {
		_, _ = formats.NewISO0().Decode("0A00000000000000", "0000000000000")
	})
	// ECI/VISA2 with spaces / short remainder after length digit
	require.NotPanics(t, func() {
		for _, ctor := range []func() formats.Format{
			formats.NewISO0, formats.NewISO1, formats.NewISO3,
			formats.NewECI1, formats.NewECI2, formats.NewECI3, formats.NewECI4,
			formats.NewVISA1, formats.NewVISA2, formats.NewVISA3, formats.NewVISA4,
			formats.NewOEM1, formats.NewANSIX98,
		} {
			// A PIN containing a space is not decimal, so every format refuses it.
			block, err := ctor().Encode("00 0", "0")
			require.Error(t, err)
			require.Empty(t, block)
			_, _ = ctor().Decode("0A00000000000000", "0000000000000")
			// Every length is checked, including empty.
			for _, b := range []string{"", "0A0000000000000", "0A000000000000000", "0A000000000000000000000000000000"} {
				_, _ = ctor().Decode(b, "0")
			}
		}
	})
}
