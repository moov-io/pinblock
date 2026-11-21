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
		pinBlock, err := iso1.Encode(pin, "")

		require.NoError(t, err)
		require.Equal(t, pinBlock, "1234555555555555")

		pin = "123456789012"
		pinBlock, err = iso1.Encode(pin, "")

		require.NoError(t, err)
		require.Equal(t, pinBlock, "123456789012AAAA")
	})

	t.Run("Decode", func(t *testing.T) {
		iso1 := formats.NewOEM1()
		pin, err := iso1.Decode("1234555555555555", "")

		require.NoError(t, err)
		require.Equal(t, "1234", pin)

		pin, err = iso1.Decode("123456789012AAAA", "")

		require.NoError(t, err)
		require.Equal(t, "123456789012", pin)
	})

	t.Run("encode logs", func(t *testing.T) {

		iso1 := formats.NewOEM1()
		out := bytes.NewBuffer([]byte{})
		iso1.SetDebugWriter(out)

		pin := "1234"
		pinBlock, err := iso1.Encode(pin, "")

		require.NoError(t, err)
		require.Equal(t, pinBlock, "1234555555555555")

		requireDebugNote(t, out.String(), "encode", "Diebold, Docutel, NCR")
	})

	t.Run("decode logs", func(t *testing.T) {
		iso1 := formats.NewOEM1()
		out := bytes.NewBuffer([]byte{})
		iso1.SetDebugWriter(out)

		pin, err := iso1.Decode("123456789012AAAA", "")

		require.NoError(t, err)
		require.Equal(t, "123456789012", pin)

		requireDebugNote(t, out.String(), "decode", "Diebold, Docutel, NCR")
	})
}
