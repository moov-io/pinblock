package formats_test

import (
	"testing"

	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
)

func TestOEM1(t *testing.T) {
	pinBlock, err := formats.NewOEM1().Encode(testPIN, nil)
	require.NoError(t, err)
	require.Equal(t, "1234555555555555", hexOf(pinBlock))

	pinBlock, err = formats.NewOEM1().Encode([]byte("123456789012"), nil)
	require.NoError(t, err)
	require.Equal(t, "123456789012AAAA", hexOf(pinBlock))

	pin, err := formats.NewOEM1().Decode(block(t, "1234555555555555"), nil)
	require.NoError(t, err)
	require.Equal(t, testPIN, pin)

	pin, err = formats.NewOEM1().Decode(block(t, "123456789012AAAA"), nil)
	require.NoError(t, err)
	require.Equal(t, []byte("123456789012"), pin)

	_, err = formats.NewOEM1().Decode(block(t, "123555555555555"+"5"), nil)
	require.ErrorContains(t, err, "invalid pin length")
}
