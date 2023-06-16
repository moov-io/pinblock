package formats_test

import (
	"bytes"
	"testing"

	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
)

func TestISO0(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		pin := "1234"
		account := "5432101234567891"

		iso0 := formats.NewISO0()
		pinBlock, err := iso0.Encode(pin, account)

		require.NoError(t, err)

		// can be checked here:
		// https://paymentcardtools.com/pin-block-calculators/iso9564-format-0
		require.Equal(t, "041215FEDCBA9876", pinBlock)
	})

	t.Run("encode logs", func(t *testing.T) {
		pin := "1234"
		account := "5432101234567891"

		iso0 := formats.NewISO0()

		out := bytes.NewBuffer([]byte{})
		iso0.SetWriter(out)

		pinBlock, err := iso0.Encode(pin, account)
		require.NoError(t, err)
		require.Equal(t, "041215FEDCBA9876", pinBlock)

		expectedOutput := `PIN block encode operation finished
************************************
PAN     : 5432101234567891
PIN     : 1234
PAD     : FFFFFFFFFF
Format  : Format 0 (ISO-0)
------------------------------------
Encoded PIN block  : 041215FEDCBA9876
Encoded PAN block  : 0000210123456789

`
		require.Equal(t, expectedOutput, out.String())
	})

	t.Run("bad pin length", func(t *testing.T) {
		iso0 := formats.NewISO0()

		account := "5432101234567891"

		// test short pin
		shortPin := "123"

		_, err := iso0.Encode(shortPin, account)

		require.Error(t, err)

		// test long pin
		longPin := "1234567890123"

		_, err = iso0.Encode(longPin, account)

		require.Error(t, err)
	})

	t.Run("bad account length", func(t *testing.T) {
		iso0 := formats.NewISO0()
		pin := "1234"

		// test short account
		shortAccount := "456789"

		_, err := iso0.Encode(pin, shortAccount)

		require.Error(t, err)
	})
}

func TestISO0_Decode(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		account := "5432101234567891"
		pinBlock := "041215FEDCBA9876"

		// test Decode
		iso0 := formats.NewISO0()
		pin, err := iso0.Decode(pinBlock, account)
		require.NoError(t, err)

		require.Equal(t, "1234", pin)
	})

	t.Run("decode logs", func(t *testing.T) {
		account := "5432101234567891"
		pinBlock := "041215FEDCBA9876"

		// test Decode
		iso0 := formats.NewISO0()

		out := bytes.NewBuffer([]byte{})
		iso0.SetWriter(out)

		pin, err := iso0.Decode(pinBlock, account)
		require.NoError(t, err)
		require.Equal(t, "1234", pin)

		expectedOutput := `PIN block decode operation finished
************************************
PAN block  : 0000210123456789
PIN block  : 041215FEDCBA9876
PAD        : FFFFFFFFFF
Format     : Format 0 (ISO-0)
------------------------------------
Decoded PIN  : 1234
Decoded PAN  : 5432101234567891

`
		require.Equal(t, expectedOutput, out.String())
	})

	t.Run("bad pin block length", func(t *testing.T) {
		iso0 := formats.NewISO0()
		account := "5432101234567891"

		// test short pin block
		shortPinBlock := "041215FEDCBA987"

		_, err := iso0.Decode(shortPinBlock, account)

		require.Error(t, err)

		// test long pin block
		longPinBlock := "041215FEDCBA98765"

		_, err = iso0.Decode(longPinBlock, account)

		require.Error(t, err)
	})

	t.Run("bad account length", func(t *testing.T) {
		iso0 := formats.NewISO0()

		// test short account
		shortAccount := "456789"

		_, err := iso0.Decode("041215FEDCBA9876", shortAccount)

		require.Error(t, err)
	})
}
