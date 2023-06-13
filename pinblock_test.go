package pinblock

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPINBlock(t *testing.T) {
	pin := "1234"
	account := "5432101234567891"

	// Create a PIN block in ISO0 format.
	pinBlock, err := EncodePINBlock(pin, account, "ISO0")
	require.NoError(t, err)

	require.Equal(t, "041215FEDCBA9876", pinBlock)
}
