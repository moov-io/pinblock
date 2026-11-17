package pinblock

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClear(t *testing.T) {
	pin := []byte("1234")
	Clear(pin)
	require.Equal(t, make([]byte, 4), pin)

	buf := bytes.Repeat([]byte{0x5A}, 16)
	Clear(buf[:8][:16])
	require.Equal(t, make([]byte, 16), buf)

	Clear(nil)
}
