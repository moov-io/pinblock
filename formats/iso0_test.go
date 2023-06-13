package formats_test

import (
	"testing"

	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
)

func TestISO0Encode(t *testing.T) {
	pin := "1234"
	account := "5432101234567891"

	iso0 := formats.NewISO0()
	pinBlock, err := iso0.Encode(pin, account)

	require.NoError(t, err)

	// can be checked here:
	// https://paymentcardtools.com/pin-block-calculators/iso9564-format-0
	require.Equal(t, "041215FEDCBA9876", pinBlock)

	// test Decode
	pin, err = iso0.Decode(pinBlock, account)
	require.NoError(t, err)

	require.Equal(t, "1234", pin)
}
