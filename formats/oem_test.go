package formats_test

import (
	"bytes"
	"testing"

	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
)

func TestOEM1(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		pin := "1234"

		iso1 := formats.NewOEM1()
		pinBlock, err := iso1.Encode(pin)

		require.NoError(t, err)
		require.Equal(t, pinBlock, "1234555555555555")

		pin = "123456789012"
		pinBlock, err = iso1.Encode(pin)

		require.NoError(t, err)
		require.Equal(t, pinBlock, "123456789012AAAA")
	})

	t.Run("Decode", func(t *testing.T) {
		iso1 := formats.NewOEM1()
		pin, err := iso1.Decode("1234555555555555")

		require.NoError(t, err)
		require.Equal(t, "1234", pin)

		pin, err = iso1.Decode("123456789012AAAA")

		require.NoError(t, err)
		require.Equal(t, "123456789012", pin)
	})

	t.Run("encode logs", func(t *testing.T) {

		iso1 := formats.NewOEM1()
		out := bytes.NewBuffer([]byte{})
		iso1.SetDebugWriter(out)

		pin := "1234"
		pinBlock, err := iso1.Encode(pin)

		require.NoError(t, err)
		require.Equal(t, pinBlock, "1234555555555555")

		expectedOutput := `PIN block encode operation finished
************************************
PIN     : 1234
PAD     : 555555555555
Format  : Diebold, Docutel, NCR
------------------------------------
Formatted PIN block  : 1234555555555555

`
		require.Equal(t, out.String(), expectedOutput)
	})

	t.Run("decode logs", func(t *testing.T) {
		iso1 := formats.NewOEM1()
		out := bytes.NewBuffer([]byte{})
		iso1.SetDebugWriter(out)

		pin, err := iso1.Decode("123456789012AAAA")

		require.NoError(t, err)
		require.Equal(t, "123456789012", pin)

		expectedOutput := `PIN block decode operation finished
************************************
Formatted PIN block  : 123456789012AAAA
PAD                  : AAAA
Format               : Diebold, Docutel, NCR
------------------------------------
Decoded PIN  : 123456789012

`
		require.Equal(t, expectedOutput, out.String())
	})
}
