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
			block, err := ctor().Encode("00 0", "0")
			if err != nil {
				_, _ = ctor().Decode("0A00000000000000", "0000000000000")
				continue
			}
			_, _ = ctor().Decode(block, "0")
		}
	})
}
