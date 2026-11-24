package formats_test

import (
	"errors"
	"testing"

	"github.com/moov-io/pinblock/encryption"
	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
)

func TestNewFormatter(t *testing.T) {
	names := []string{
		"ISO-0", "ISO-1", "ISO-2", "ISO-3", "ANSI", "OEM1",
		"ECI1", "ECI2", "ECI3", "ECI4", "VISA1", "VISA2", "VISA3", "VISA4",
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			f, err := formats.NewFormatter(name)
			require.NoError(t, err)
			require.NotNil(t, f)
		})
	}

	t.Run("OEM1 is OEM-1", func(t *testing.T) {
		f, err := formats.NewFormatter("OEM1")
		require.NoError(t, err)
		got, err := f.Encode("1234", "")
		require.NoError(t, err)
		want, err := formats.NewOEM1().Encode("1234", "")
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("ISO-4 needs a cipher", func(t *testing.T) {
		f, err := formats.NewFormatter("ISO-4")
		require.Nil(t, f)
		require.ErrorIs(t, err, formats.ErrCipherRequired)
		require.Contains(t, err.Error(), "NewISO4")
	})

	t.Run("unknown name", func(t *testing.T) {
		f, err := formats.NewFormatter("ISO-9")
		require.Nil(t, f)
		require.ErrorIs(t, err, formats.ErrUnsupportedFormat)
		for _, name := range names {
			require.Contains(t, err.Error(), name)
		}
	})
}

func TestDecodeErrors(t *testing.T) {
	const account = "43219876543210987"

	iso3Block, err := formats.NewISO3().Encode("1234", account)
	require.NoError(t, err)

	cases := []struct {
		name     string
		format   formats.Format
		pinBlock string
		want     error
	}{
		{"wrong format nibble", formats.NewISO0(), iso3Block, formats.ErrFormatMismatch},
		{"PIN length out of range", formats.NewISO0(), "0F12AC89ABCDEF67", formats.ErrInvalidPINLength},
		{"wrong fill", formats.NewISO0(), "0412AC89ABCDEF68", formats.ErrInvalidFill},
		{"wrong length", formats.NewISO0(), "0412", nil},
		{"not hex", formats.NewISO0(), "0412AC89ABCDEFZZ", nil},
		{"ECI-3 unparsable", formats.NewECI3(), "9123456789012345", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.format.Decode(tc.pinBlock, account)
			require.Error(t, err)
			require.ErrorIs(t, err, formats.ErrInvalidPINBlock)
			if tc.want != nil {
				require.ErrorIs(t, err, tc.want)
			}
		})
	}

	t.Run("ISO-4 wrong key", func(t *testing.T) {
		key := []byte{
			0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
			0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
		}
		other := append([]byte(nil), key...)
		other[0] ^= 0x01

		c, err := encryption.NewAesECB(key)
		require.NoError(t, err)
		block, err := formats.NewISO4(c).Encode("1234", account)
		require.NoError(t, err)

		wrong, err := encryption.NewAesECB(other)
		require.NoError(t, err)
		_, err = formats.NewISO4(wrong).Decode(block, account)
		require.ErrorIs(t, err, formats.ErrInvalidPINBlock)
	})

	t.Run("error text does not change", func(t *testing.T) {
		_, err := formats.NewISO0().Decode(iso3Block, account)
		require.EqualError(t, err, "format is different")
		_, err = formats.NewISO0().Decode("0F12AC89ABCDEF67", account)
		require.EqualError(t, err, "invalid pin length 15")
	})

	t.Run("input errors do not match", func(t *testing.T) {
		_, err := formats.NewISO0().Encode("12", account)
		require.Error(t, err)
		require.False(t, errors.Is(err, formats.ErrInvalidPINBlock))
	})
}
