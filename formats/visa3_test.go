package formats_test

import (
	"bytes"
	"testing"

	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
)

func TestVISA3(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		pin := "1234"

		visa3 := formats.NewVISA3()
		pinBlock, err := visa3.Encode(pin)

		require.NoError(t, err)
		require.Contains(t, pinBlock, "1234F")

		pin = "123456789012"
		pinBlock, err = visa3.Encode(pin)

		require.NoError(t, err)
		require.Contains(t, pinBlock, "123456789012F")
	})

	t.Run("Decode", func(t *testing.T) {
		visa3 := formats.NewVISA3()
		pin, err := visa3.Decode("1234F55555555555")

		require.NoError(t, err)
		require.Equal(t, "1234", pin)

		pin, err = visa3.Decode("123456789012FAAA")

		require.NoError(t, err)
		require.Equal(t, "123456789012", pin)
	})

	t.Run("encode logs", func(t *testing.T) {

		visa3 := formats.NewVISA3()
		out := bytes.NewBuffer([]byte{})
		visa3.SetDebugWriter(out)

		pin := "1234"
		pinBlock, err := visa3.Encode(pin)

		require.NoError(t, err)
		require.Contains(t, pinBlock, "1234F")

		expectedOutput := `PIN block encode operation finished
************************************
PIN     : 1234
PAD     : `
		require.Contains(t, out.String(), expectedOutput)
	})

	t.Run("decode logs", func(t *testing.T) {
		visa3 := formats.NewVISA3()
		out := bytes.NewBuffer([]byte{})
		visa3.SetDebugWriter(out)

		pin, err := visa3.Decode("123456789012FAAA")

		require.NoError(t, err)
		require.Equal(t, "123456789012", pin)

		expectedOutput := `PIN block decode operation finished
************************************
Formatted PIN block  : 123456789012FAAA
PAD                  : AAA
Format               : VISA-3
------------------------------------
Decoded PIN  : 123456789012

`
		require.Equal(t, expectedOutput, out.String())
	})
}
