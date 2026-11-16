package memzero

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytes(t *testing.T) {
	b := bytes.Repeat([]byte{0xA5}, 24)
	Bytes(b)
	require.Equal(t, make([]byte, 24), b)

	// A sub-slice clears only its own range.
	b = bytes.Repeat([]byte{0xA5}, 16)
	Bytes(b[4:8])
	require.Equal(t, []byte{0xA5, 0xA5, 0xA5, 0xA5, 0, 0, 0, 0}, b[:8])
	require.Equal(t, bytes.Repeat([]byte{0xA5}, 8), b[8:])

	// nil and empty slices are accepted.
	Bytes(nil)
	Bytes([]byte{})
}
