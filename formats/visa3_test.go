package formats_test

import (
	"testing"

	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
)

func TestVISA3(t *testing.T) {
	pinBlock, err := formats.NewVISA3().Encode(testPIN, nil)
	require.NoError(t, err)
	require.Equal(t, "1234F", hexOf(pinBlock)[:5])

	pinBlock, err = formats.NewVISA3().Encode([]byte("123456789012"), nil)
	require.NoError(t, err)
	require.Equal(t, "123456789012F", hexOf(pinBlock)[:13])

	pin, err := formats.NewVISA3().Decode(block(t, "1234F55555555555"), nil)
	require.NoError(t, err)
	require.Equal(t, testPIN, pin)

	pin, err = formats.NewVISA3().Decode(block(t, "123456789012FAAA"), nil)
	require.NoError(t, err)
	require.Equal(t, []byte("123456789012"), pin)

	_, err = formats.NewVISA3().Decode(block(t, "123F555555555555"), nil)
	require.Error(t, err)
}
