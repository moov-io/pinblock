package pinblock

import (
	"bytes"
	"github.com/moov-io/pinblock/encryption"
	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestISO0_describe(t *testing.T) {

	t.Run("encode describe", func(t *testing.T) {
		pin := "1234"
		account := "5432101234567891"

		iso0 := formats.NewISO0()

		out := bytes.NewBuffer([]byte{})
		err := Describe(iso0, out, DescribeEncode, pin, account)
		require.NoError(t, err)

		expectedOutput := `PIN block encode operation finished
************************************
PAN     : 5432101234567891
PIN     : 1234
PAD     : FFFFFFFFFF
Format  : Format 0 (ISO-0)
------------------------------------
Clear PIN block  : 041215FEDCBA9876

`
		require.Equal(t, expectedOutput, out.String())

	})

	t.Run("decode describe", func(t *testing.T) {
		account := "5432101234567891"
		pinBlock := "041215FEDCBA9876"

		iso0 := formats.NewISO0()

		out := bytes.NewBuffer([]byte{})
		err := Describe(iso0, out, DescribeDecode, pinBlock, account)
		require.NoError(t, err)

		expectedOutput := `PIN block decode operation finished
************************************
PIN block  : 041215FEDCBA9876
PAN block  : 5432101234567891
PAD        : FFFFFFFFFF
Format     : Format 0 (ISO-0)
------------------------------------
Decoded PIN  : 1234

`
		require.Equal(t, expectedOutput, out.String())
	})

	t.Run("unknown describe", func(t *testing.T) {
		pin := "1234"
		account := "5432101234567891"

		iso0 := formats.NewISO0()

		err := Describe(iso0, nil, "unknown", pin, account)
		require.Error(t, err)
	})

}

func TestISO4_describe(t *testing.T) {
	t.Run("encode describe", func(t *testing.T) {
		cipher, err := encryption.NewAesECB([]byte("1234567890123456"))
		require.NoError(t, err)

		iso4 := formats.NewISO4(cipher)

		pin := "1234"
		account := "5432101234567891"

		out := bytes.NewBuffer([]byte{})
		err = Describe(iso4, out, DescribeEncode, pin, account)
		require.NoError(t, err)

		expectedOutput := `PIN block encode operation finished
************************************
PAN     : 5432101234567891
PIN     : 1234
PAD     : AAAAAAAAAA
Format  : Format 4 (ISO-4)
------------------------------------
Clear PIN block  :`
		require.Equal(t, true, strings.Contains(out.String(), expectedOutput))
	})

	t.Run("decode describe", func(t *testing.T) {
		cipher, err := encryption.NewAesECB([]byte("1234567890123456"))
		require.NoError(t, err)

		iso4 := formats.NewISO4(cipher)

		pinBlock := "5BD956F4D13732997B7FFAF357541D2B"
		account := "5432101234567891"

		out := bytes.NewBuffer([]byte{})
		err = Describe(iso4, out, DescribeDecode, pinBlock, account)
		require.NoError(t, err)

		expectedOutput := `PIN block decode operation finished
************************************
PIN block  : 5BD956F4D13732997B7FFAF357541D2B
PAN block  : 5432101234567891
PAD        : AAAAAAAAAA
Format     : Format 4 (ISO-4)
------------------------------------
Decoded PIN  : 1234

`
		require.Equal(t, expectedOutput, out.String())
	})

	t.Run("unknown describe", func(t *testing.T) {
		pin := "1234"
		account := "5432101234567891"

		cipher, err := encryption.NewAesECB([]byte("1234567890123456"))
		require.NoError(t, err)

		iso4 := formats.NewISO4(cipher)

		err = Describe(iso4, nil, "unknown", pin, account)
		require.Error(t, err)
	})
}
