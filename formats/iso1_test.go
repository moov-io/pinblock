package formats_test

import (
	"bytes"
	"testing"

	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
)

func TestISO1(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		pin := "1234"

		iso1 := formats.NewISO1()
		pinBlock, err := iso1.Encode(pin, "")

		require.NoError(t, err)
		require.Contains(t, pinBlock, "141234")

		pin = "123456789012"
		pinBlock, err = iso1.Encode(pin, "")

		require.NoError(t, err)
		require.Contains(t, pinBlock, "1C123456789012")
	})

	t.Run("Decode", func(t *testing.T) {
		iso1 := formats.NewISO1()
		pin, err := iso1.Decode("141234FFFFFFFFFF", "")

		require.NoError(t, err)
		require.Equal(t, "1234", pin)

		pin, err = iso1.Decode("1C123456789012FF", "")

		require.NoError(t, err)
		require.Equal(t, "123456789012", pin)
	})

	t.Run("encode logs", func(t *testing.T) {

		iso1 := formats.NewISO1()
		out := bytes.NewBuffer([]byte{})
		iso1.SetDebugWriter(out)

		pin := "1234"
		pinBlock, err := iso1.Encode(pin, "")

		require.NoError(t, err)
		require.Contains(t, pinBlock, "141234")

		requireDebugNote(t, out.String(), "encode", "Format 1 (ISO-1)")
	})

	t.Run("decode logs", func(t *testing.T) {
		iso1 := formats.NewISO1()
		out := bytes.NewBuffer([]byte{})
		iso1.SetDebugWriter(out)

		pin, err := iso1.Decode("1C123456789012FF", "")

		require.NoError(t, err)
		require.Equal(t, "123456789012", pin)

		requireDebugNote(t, out.String(), "decode", "Format 1 (ISO-1)")
	})
}

func TestISO2(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		pin := "1234"

		iso2 := formats.NewISO2()
		pinBlock, err := iso2.Encode(pin, "")

		require.NoError(t, err)
		require.Equal(t, "241234FFFFFFFFFF", pinBlock)

		pin = "123456789012"
		pinBlock, err = iso2.Encode(pin, "")

		require.NoError(t, err)
		require.Equal(t, "2C123456789012FF", pinBlock)
	})

	t.Run("Decode", func(t *testing.T) {
		iso2 := formats.NewISO2()
		pin, err := iso2.Decode("241234FFFFFFFFFF", "")

		require.NoError(t, err)
		require.Equal(t, "1234", pin)

		pin, err = iso2.Decode("2C123456789012FF", "")

		require.NoError(t, err)
		require.Equal(t, "123456789012", pin)
	})

	t.Run("encode logs", func(t *testing.T) {

		iso2 := formats.NewISO2()
		out := bytes.NewBuffer([]byte{})
		iso2.SetDebugWriter(out)

		pin := "1234"
		pinBlock, err := iso2.Encode(pin, "")

		require.NoError(t, err)
		require.Equal(t, "241234FFFFFFFFFF", pinBlock)

		requireDebugNote(t, out.String(), "encode", "Format 2 (ISO-2)")
	})

	t.Run("decode logs", func(t *testing.T) {
		iso2 := formats.NewISO2()
		out := bytes.NewBuffer([]byte{})
		iso2.SetDebugWriter(out)

		pin, err := iso2.Decode("2C123456789012FF", "")

		require.NoError(t, err)
		require.Equal(t, "123456789012", pin)

		requireDebugNote(t, out.String(), "decode", "Format 2 (ISO-2)")
	})
}

func TestECI4(t *testing.T) {
	t.Run("ECI4 logs", func(t *testing.T) {
		iso1 := formats.NewECI4()
		out := bytes.NewBuffer([]byte{})
		iso1.SetDebugWriter(out)

		pin, err := iso1.Decode("1C123456789012FF", "")

		require.NoError(t, err)
		require.Equal(t, "123456789012", pin)

		requireDebugNote(t, out.String(), "decode", "ECI-4")
	})
}
