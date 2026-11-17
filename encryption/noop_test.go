package encryption

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// NoOp returns a copy, so clearing its input does not clear its output.
func TestNoOp_ReturnsCopy(t *testing.T) {
	in := []byte("1234567890123456")
	for _, op := range []func([]byte) ([]byte, error){NewNoOp().Encrypt, NewNoOp().Decrypt} {
		out, err := op(in)
		require.NoError(t, err)
		require.Equal(t, in, out)
		out[0] = 'X'
		require.Equal(t, byte('1'), in[0])
	}
}
