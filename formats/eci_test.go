package formats_test

import (
	"testing"

	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
)

func TestECI2(t *testing.T) {
	pinBlock, err := formats.NewECI2().Encode(testPIN, nil)
	require.NoError(t, err)
	require.Equal(t, "1234", hexOf(pinBlock)[:4])

	// ECI-2 holds 4 digits; a longer PIN is truncated.
	pinBlock, err = formats.NewECI2().Encode([]byte("123456789012"), nil)
	require.NoError(t, err)
	require.Equal(t, "1234", hexOf(pinBlock)[:4])

	pin, err := formats.NewECI2().Decode(block(t, "1234555555555555"), nil)
	require.NoError(t, err)
	require.Equal(t, testPIN, pin)

	pin, err = formats.NewECI2().Decode(block(t, "123456789012AAAA"), nil)
	require.NoError(t, err)
	require.Equal(t, testPIN, pin)
}

func TestECI3AndVISA2(t *testing.T) {
	for _, ctor := range []func() formats.Format{formats.NewECI3, formats.NewVISA2} {
		pinBlock, err := ctor().Encode(testPIN, nil)
		require.NoError(t, err)
		require.Equal(t, "4123400", hexOf(pinBlock)[:7])

		pinBlock, err = ctor().Encode([]byte("123456789012"), nil)
		require.NoError(t, err)
		require.Equal(t, "6123456", hexOf(pinBlock)[:7])

		pin, err := ctor().Decode(block(t, "4123400555555555"), nil)
		require.NoError(t, err)
		require.Equal(t, testPIN, pin)

		pin, err = ctor().Decode(block(t, "6123456789012AAA"), nil)
		require.NoError(t, err)
		require.Equal(t, []byte("123456"), pin)

		_, err = ctor().Decode(block(t, "7123456789012AAA"), nil)
		require.Error(t, err)
	}

	// VISA-2 fills with decimal digits only.
	for range 20 {
		pinBlock, err := formats.NewVISA2().Encode(testPIN, nil)
		require.NoError(t, err)
		for _, c := range hexOf(pinBlock)[7:] {
			require.True(t, c >= '0' && c <= '9', "VISA-2 fill %s", hexOf(pinBlock))
		}
	}
}
