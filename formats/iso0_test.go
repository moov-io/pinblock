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
		iso0.SetDebugWriter(out)

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
Formatted PIN block  : 041215FEDCBA9876
Formatted PAN block  : 0000210123456789

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

func TestISO3(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		pin := "1234"
		account := "5432101234567891"

		iso3 := formats.NewISO3()
		pinBlock, err := iso3.Encode(pin, account)

		require.NoError(t, err)

		// can be checked here:
		// https://paymentcardtools.com/pin-block-calculators/iso9564-format-0
		require.Contains(t, pinBlock, "341215")
	})

	t.Run("encode logs", func(t *testing.T) {
		pin := "1234"
		account := "5432101234567891"

		iso3 := formats.NewISO3()

		out := bytes.NewBuffer([]byte{})
		iso3.SetDebugWriter(out)

		pinBlock, err := iso3.Encode(pin, account)
		require.NoError(t, err)
		require.Contains(t, pinBlock, "341215")

		expectedOutput := `PIN block encode operation finished
************************************
PAN     : 5432101234567891
PIN     : 1234`
		require.Contains(t, out.String(), expectedOutput)
	})

	t.Run("bad pin length", func(t *testing.T) {
		iso3 := formats.NewISO3()

		account := "5432101234567891"

		// test short pin
		shortPin := "123"

		_, err := iso3.Encode(shortPin, account)

		require.Error(t, err)

		// test long pin
		longPin := "1234567890123"

		_, err = iso3.Encode(longPin, account)

		require.Error(t, err)
	})

	t.Run("bad account length", func(t *testing.T) {
		iso3 := formats.NewISO3()
		pin := "1234"

		// test short account
		shortAccount := "456789"

		_, err := iso3.Encode(pin, shortAccount)

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
		iso0.SetDebugWriter(out)

		pin, err := iso0.Decode(pinBlock, account)
		require.NoError(t, err)
		require.Equal(t, "1234", pin)

		expectedOutput := `PIN block decode operation finished
************************************
Formatted PAN block  : 0000210123456789
Formatted PIN block  : 041234FFFFFFFFFF
PAD                  : FFFFFFFFFF
Format               : Format 0 (ISO-0)
------------------------------------
Decoded PIN  : 1234

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

func TestISO3_Decode(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		account := "5432101234567891"
		pinBlock := "341215FEDCBA9876"

		// test Decode
		iso3 := formats.NewISO3()
		pin, err := iso3.Decode(pinBlock, account)
		require.NoError(t, err)

		require.Equal(t, "1234", pin)
	})

	t.Run("decode logs", func(t *testing.T) {
		account := "5432101234567891"
		pinBlock := "341215FEDCBA9876"

		// test Decode
		iso3 := formats.NewISO3()

		out := bytes.NewBuffer([]byte{})
		iso3.SetDebugWriter(out)

		pin, err := iso3.Decode(pinBlock, account)
		require.NoError(t, err)
		require.Equal(t, "1234", pin)

		expectedOutput := `PIN block decode operation finished
************************************
Formatted PAN block  : 0000210123456789
Formatted PIN block  : 341234FFFFFFFFFF`
		require.Contains(t, out.String(), expectedOutput)
	})

	t.Run("bad pin block length", func(t *testing.T) {
		iso3 := formats.NewISO3()
		account := "5432101234567891"

		// test short pin block
		shortPinBlock := "041215FEDCBA987"

		_, err := iso3.Decode(shortPinBlock, account)

		require.Error(t, err)

		// test long pin block
		longPinBlock := "041215FEDCBA98765"

		_, err = iso3.Decode(longPinBlock, account)

		require.Error(t, err)
	})

	t.Run("bad account length", func(t *testing.T) {
		iso3 := formats.NewISO3()

		// test short account
		shortAccount := "456789"

		_, err := iso3.Decode("041215FEDCBA9876", shortAccount)

		require.Error(t, err)
	})
}

func TestANSIX98(t *testing.T) {
	t.Run("ANSI X9.8 logs", func(t *testing.T) {
		pin := "1234"
		account := "5432101234567891"

		iso0 := formats.NewANSIX98()

		out := bytes.NewBuffer([]byte{})
		iso0.SetDebugWriter(out)

		pinBlock, err := iso0.Encode(pin, account)
		require.NoError(t, err)
		require.Equal(t, "041215FEDCBA9876", pinBlock)

		expectedOutput := `PIN block encode operation finished
************************************
PAN     : 5432101234567891
PIN     : 1234
PAD     : FFFFFFFFFF
Format  : ANSI X9.8
------------------------------------
Formatted PIN block  : 041215FEDCBA9876
Formatted PAN block  : 0000210123456789

`
		require.Equal(t, expectedOutput, out.String())
	})
}

func TestECI1(t *testing.T) {
	t.Run("ECI-1 logs", func(t *testing.T) {
		pin := "1234"
		account := "5432101234567891"

		iso0 := formats.NewECI1()

		out := bytes.NewBuffer([]byte{})
		iso0.SetDebugWriter(out)

		pinBlock, err := iso0.Encode(pin, account)
		require.NoError(t, err)
		require.Equal(t, "041215FEDCBA9876", pinBlock)

		expectedOutput := `PIN block encode operation finished
************************************
PAN     : 5432101234567891
PIN     : 1234
PAD     : FFFFFFFFFF
Format  : ECI-1
------------------------------------
Formatted PIN block  : 041215FEDCBA9876
Formatted PAN block  : 0000210123456789

`
		require.Equal(t, expectedOutput, out.String())
	})
}

func TestVISA1(t *testing.T) {
	t.Run("VISA-1 logs", func(t *testing.T) {
		pin := "1234"
		account := "5432101234567891"

		iso0 := formats.NewVISA1()

		out := bytes.NewBuffer([]byte{})
		iso0.SetDebugWriter(out)

		pinBlock, err := iso0.Encode(pin, account)
		require.NoError(t, err)
		require.Equal(t, "041215FEDCBA9876", pinBlock)

		expectedOutput := `PIN block encode operation finished
************************************
PAN     : 5432101234567891
PIN     : 1234
PAD     : FFFFFFFFFF
Format  : VISA-1
------------------------------------
Formatted PIN block  : 041215FEDCBA9876
Formatted PAN block  : 0000210123456789

`
		require.Equal(t, expectedOutput, out.String())
	})
}

func TestVISA4(t *testing.T) {
	t.Run("VISA-4 logs", func(t *testing.T) {
		pin := "1234"
		account := "5432101234567891"

		iso0 := formats.NewVISA4()

		out := bytes.NewBuffer([]byte{})
		iso0.SetDebugWriter(out)

		pinBlock, err := iso0.Encode(pin, account)
		require.NoError(t, err)
		require.Equal(t, "041215FEDCBA9876", pinBlock)

		expectedOutput := `PIN block encode operation finished
************************************
PAN     : 5432101234567891
PIN     : 1234
PAD     : FFFFFFFFFF
Format  : VISA-4
------------------------------------
Formatted PIN block  : 041215FEDCBA9876
Formatted PAN block  : 0000210123456789

`
		require.Equal(t, expectedOutput, out.String())
	})
}
