package formats_test

import (
	"testing"

	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
)

func TestISO1(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		pinBlock, err := formats.NewISO1().Encode(testPIN, nil)
		require.NoError(t, err)
		require.Equal(t, "141234", hexOf(pinBlock)[:6])

		pinBlock, err = formats.NewISO1().Encode([]byte("123456789012"), nil)
		require.NoError(t, err)
		require.Equal(t, "1C123456789012", hexOf(pinBlock)[:14])
	})

	t.Run("Decode", func(t *testing.T) {
		pin, err := formats.NewISO1().Decode(block(t, "141234FFFFFFFFFF"), nil)
		require.NoError(t, err)
		require.Equal(t, testPIN, pin)

		pin, err = formats.NewISO1().Decode(block(t, "1C123456789012FF"), nil)
		require.NoError(t, err)
		require.Equal(t, []byte("123456789012"), pin)

		pin, err = formats.NewECI4().Decode(block(t, "1C123456789012FF"), nil)
		require.NoError(t, err)
		require.Equal(t, []byte("123456789012"), pin)
	})

	// ISO 9564-1 allows 4 to 12 PIN digits.
	t.Run("pin length", func(t *testing.T) {
		_, err := formats.NewISO1().Encode([]byte("1234567890123"), nil)
		require.Error(t, err)
		_, err = formats.NewISO1().Encode([]byte("123"), nil)
		require.Error(t, err)
		_, err = formats.NewISO1().Decode(block(t, "13123FFFFFFFFFFF"), nil)
		require.ErrorContains(t, err, "invalid pin length")
		_, err = formats.NewISO1().Decode(block(t, "1E12345678901234"), nil)
		require.ErrorContains(t, err, "invalid pin length")
	})

	// Format 1 fill is random 0x0 to 0xF. Across 20 blocks (200 fill nibbles)
	// at least one is A to F unless the fill is limited to decimal digits.
	t.Run("fill is random 0-F", func(t *testing.T) {
		sawLetter := false
		for range 20 {
			pinBlock, err := formats.NewISO1().Encode(testPIN, nil)
			require.NoError(t, err)
			for _, c := range hexOf(pinBlock)[6:] {
				if c >= 'A' && c <= 'F' {
					sawLetter = true
				}
			}
		}
		require.True(t, sawLetter, "ISO-1 fill never used A-F")
	})

	t.Run("round trip", func(t *testing.T) {
		for _, ctor := range []func() formats.Format{formats.NewISO1, formats.NewECI4} {
			pinBlock, err := ctor().Encode([]byte("98765"), nil)
			require.NoError(t, err)
			pin, err := ctor().Decode(pinBlock, nil)
			require.NoError(t, err)
			require.Equal(t, []byte("98765"), pin)
		}
	})
}

func TestISO2(t *testing.T) {
	pinBlock, err := formats.NewISO2().Encode(testPIN, nil)
	require.NoError(t, err)
	require.Equal(t, "241234FFFFFFFFFF", hexOf(pinBlock))

	pinBlock, err = formats.NewISO2().Encode([]byte("123456789012"), nil)
	require.NoError(t, err)
	require.Equal(t, "2C123456789012FF", hexOf(pinBlock))

	pin, err := formats.NewISO2().Decode(block(t, "241234FFFFFFFFFF"), nil)
	require.NoError(t, err)
	require.Equal(t, testPIN, pin)

	pin, err = formats.NewISO2().Decode(block(t, "2C123456789012FF"), nil)
	require.NoError(t, err)
	require.Equal(t, []byte("123456789012"), pin)

	_, err = formats.NewISO2().Decode(block(t, "241234FFFFFFFFFA"), nil)
	require.ErrorContains(t, err, "invalid fill")
	_, err = formats.NewISO2().Decode(block(t, "141234FFFFFFFFFF"), nil)
	require.ErrorContains(t, err, "format is different")
}
