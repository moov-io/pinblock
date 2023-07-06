package formats_test

import (
	"bytes"
	"testing"

	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
)

func TestECI2(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		pin := "1234"

		iso1 := formats.NewECI2()
		pinBlock, err := iso1.Encode(pin)

		require.NoError(t, err)
		require.Contains(t, pinBlock, "1234")

		pin = "123456789012"
		pinBlock, err = iso1.Encode(pin)

		require.NoError(t, err)
		require.Contains(t, pinBlock, "1234")
	})

	t.Run("Decode", func(t *testing.T) {
		iso1 := formats.NewECI2()
		pin, err := iso1.Decode("1234555555555555")

		require.NoError(t, err)
		require.Equal(t, "1234", pin)

		pin, err = iso1.Decode("123456789012AAAA")

		require.NoError(t, err)
		require.Equal(t, "1234", pin)
	})

	t.Run("encode logs", func(t *testing.T) {

		iso1 := formats.NewECI2()
		out := bytes.NewBuffer([]byte{})
		iso1.SetDebugWriter(out)

		pin := "1234"
		pinBlock, err := iso1.Encode(pin)

		require.NoError(t, err)
		require.Contains(t, pinBlock, "1234")

		expectedOutput := `PIN block encode operation finished
************************************
PIN     : 1234`
		require.Contains(t, out.String(), expectedOutput)
	})

	t.Run("decode logs", func(t *testing.T) {
		iso1 := formats.NewECI2()
		out := bytes.NewBuffer([]byte{})
		iso1.SetDebugWriter(out)

		pin, err := iso1.Decode("123456789012AAAA")

		require.NoError(t, err)
		require.Equal(t, "1234", pin)

		expectedOutput := `PIN block decode operation finished
************************************
Formatted PIN block  : 123456789012AAAA
PAD                  : 56789012AAAA
Format               : ECI-2
------------------------------------
Decoded PIN  : 1234

`
		require.Equal(t, out.String(), expectedOutput)
	})
}

func TestECI3(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		pin := "1234"

		iso1 := formats.NewECI3()
		pinBlock, err := iso1.Encode(pin)

		require.NoError(t, err)
		require.Contains(t, pinBlock, "4123400")

		pin = "123456789012"
		pinBlock, err = iso1.Encode(pin)

		require.NoError(t, err)
		require.Contains(t, pinBlock, "61234")
	})

	t.Run("Decode", func(t *testing.T) {
		iso1 := formats.NewECI3()
		pin, err := iso1.Decode("4123400555555555")
		require.NoError(t, err)
		require.Equal(t, "1234", pin)

		pin, err = iso1.Decode("6123456789012AAA")
		require.NoError(t, err)
		require.Equal(t, "123456", pin)
	})

	t.Run("encode logs", func(t *testing.T) {

		iso1 := formats.NewECI3()
		out := bytes.NewBuffer([]byte{})
		iso1.SetDebugWriter(out)

		pin := "1234"
		pinBlock, err := iso1.Encode(pin)

		require.NoError(t, err)
		require.Contains(t, pinBlock, "1234")

		expectedOutput := `PIN block encode operation finished
************************************
PIN     : 1234`
		require.Contains(t, out.String(), expectedOutput)
	})

	t.Run("decode logs", func(t *testing.T) {
		iso1 := formats.NewECI3()
		out := bytes.NewBuffer([]byte{})
		iso1.SetDebugWriter(out)

		pin, err := iso1.Decode("6123456789012AAA")

		require.NoError(t, err)
		require.Equal(t, "123456", pin)

		expectedOutput := `PIN block decode operation finished
************************************
Formatted PIN block  : 123456789012AAA
PAD                  : 789012AAA
Format               : ECI-3
------------------------------------
Decoded PIN  : 123456

`
		require.Equal(t, out.String(), expectedOutput)
	})
}

func TestVISA2(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		pin := "1234"

		iso1 := formats.NewVISA2()
		pinBlock, err := iso1.Encode(pin)

		require.NoError(t, err)
		require.Contains(t, pinBlock, "4123400")

		pin = "123456789012"
		pinBlock, err = iso1.Encode(pin)

		require.NoError(t, err)
		require.Contains(t, pinBlock, "61234")
	})

	t.Run("Decode", func(t *testing.T) {
		iso1 := formats.NewVISA2()
		pin, err := iso1.Decode("4123400555555555")
		require.NoError(t, err)
		require.Equal(t, "1234", pin)

		pin, err = iso1.Decode("6123456789012AAA")
		require.NoError(t, err)
		require.Equal(t, "123456", pin)
	})

	t.Run("encode logs", func(t *testing.T) {

		iso1 := formats.NewVISA2()
		out := bytes.NewBuffer([]byte{})
		iso1.SetDebugWriter(out)

		pin := "1234"
		pinBlock, err := iso1.Encode(pin)

		require.NoError(t, err)
		require.Contains(t, pinBlock, "1234")

		expectedOutput := `PIN block encode operation finished
************************************
PIN     : 1234`
		require.Contains(t, out.String(), expectedOutput)
	})

	t.Run("decode logs", func(t *testing.T) {
		iso1 := formats.NewVISA2()
		out := bytes.NewBuffer([]byte{})
		iso1.SetDebugWriter(out)

		pin, err := iso1.Decode("6123456789012222")

		require.NoError(t, err)
		require.Equal(t, "123456", pin)

		expectedOutput := `PIN block decode operation finished
************************************
Formatted PIN block  : 123456789012222
PAD                  : 789012222
Format               : VISA-2
------------------------------------
Decoded PIN  : 123456

`
		require.Equal(t, out.String(), expectedOutput)
	})
}
