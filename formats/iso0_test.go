package formats_test

import (
	"testing"

	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
)

var (
	testPIN     = []byte("1234")
	testAccount = []byte("5432101234567891")
)

func TestISO0(t *testing.T) {
	// can be checked here:
	// https://paymentcardtools.com/pin-block-calculators/iso9564-format-0
	for _, ctor := range []func() formats.Format{
		formats.NewISO0, formats.NewANSIX98, formats.NewECI1, formats.NewVISA1, formats.NewVISA4,
	} {
		pinBlock, err := ctor().Encode(testPIN, testAccount)
		require.NoError(t, err)
		require.Equal(t, "041215FEDCBA9876", hexOf(pinBlock))
	}

	t.Run("bad pin length", func(t *testing.T) {
		_, err := formats.NewISO0().Encode([]byte("123"), testAccount)
		require.Error(t, err)
		_, err = formats.NewISO0().Encode([]byte("1234567890123"), testAccount)
		require.Error(t, err)
	})

	t.Run("bad account length", func(t *testing.T) {
		_, err := formats.NewISO0().Encode(testPIN, []byte("456789"))
		require.Error(t, err)
	})

	t.Run("non-decimal pin or account", func(t *testing.T) {
		_, err := formats.NewISO0().Encode([]byte("12A4"), testAccount)
		require.ErrorContains(t, err, "pin must contain only decimal digits")
		_, err = formats.NewISO0().Encode([]byte("abcd"), testAccount)
		require.Error(t, err)
		_, err = formats.NewISO0().Encode(testPIN, []byte("54321012345678A1"))
		require.ErrorContains(t, err, "account must contain only decimal digits")
	})
}

func TestISO3(t *testing.T) {
	pinBlock, err := formats.NewISO3().Encode(testPIN, testAccount)
	require.NoError(t, err)
	require.Equal(t, "341215", hexOf(pinBlock)[:6])

	pin, err := formats.NewISO3().Decode(pinBlock, testAccount)
	require.NoError(t, err)
	require.Equal(t, testPIN, pin)

	_, err = formats.NewISO3().Encode([]byte("123"), testAccount)
	require.Error(t, err)
	_, err = formats.NewISO3().Encode([]byte("1234567890123"), testAccount)
	require.Error(t, err)
	_, err = formats.NewISO3().Encode(testPIN, []byte("456789"))
	require.Error(t, err)
}

func TestISO0_Decode(t *testing.T) {
	pin, err := formats.NewISO0().Decode(block(t, "041215FEDCBA9876"), testAccount)
	require.NoError(t, err)
	require.Equal(t, testPIN, pin)

	t.Run("bad pin block length", func(t *testing.T) {
		_, err := formats.NewISO0().Decode(block(t, "041215FEDCBA98"), testAccount)
		require.Error(t, err)
		_, err = formats.NewISO0().Decode(block(t, "041215FEDCBA987654"), testAccount)
		require.Error(t, err)
	})

	t.Run("bad account length", func(t *testing.T) {
		_, err := formats.NewISO0().Decode(block(t, "041215FEDCBA9876"), []byte("456789"))
		require.Error(t, err)
	})

	t.Run("format 3 block", func(t *testing.T) {
		_, err := formats.NewISO0().Decode(block(t, "341215FEDCBA9876"), testAccount)
		require.ErrorContains(t, err, "format is different")
	})

	// A block whose fill is not all F (for example one decrypted under the
	// wrong key) is refused rather than decoded to a PIN.
	t.Run("bad fill", func(t *testing.T) {
		_, err := formats.NewISO0().Decode(block(t, "041215FEDCBA9870"), testAccount)
		require.ErrorContains(t, err, "invalid fill")
	})

	t.Run("non-decimal pin digit", func(t *testing.T) {
		// PIN field 04123AFFFFFFFFFF XOR PAN field 0000210123456789
		_, err := formats.NewISO0().Decode(block(t, "04121BFEDCBA9876"), testAccount)
		require.ErrorContains(t, err, "non-decimal")
	})
}

func TestISO3_Decode(t *testing.T) {
	pin, err := formats.NewISO3().Decode(block(t, "341215FEDCBA9876"), testAccount)
	require.NoError(t, err)
	require.Equal(t, testPIN, pin)

	_, err = formats.NewISO3().Decode(block(t, "341215FEDCBA98"), testAccount)
	require.Error(t, err)
	_, err = formats.NewISO3().Decode(block(t, "341215FEDCBA9876"), []byte("456789"))
	require.Error(t, err)
	_, err = formats.NewISO3().Decode(block(t, "041215FEDCBA9876"), testAccount)
	require.ErrorContains(t, err, "format is different")
}

func TestISO0RoundTripTenDigitPIN(t *testing.T) {
	pin := []byte("1234567890")
	iso0 := formats.NewISO0()
	pinBlock, err := iso0.Encode(pin, testAccount)
	require.NoError(t, err)
	decoded, err := iso0.Decode(pinBlock, testAccount)
	require.NoError(t, err)
	require.Equal(t, pin, decoded)
}
