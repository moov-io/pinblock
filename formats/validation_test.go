package formats_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/moov-io/pinblock/encryption"
	"github.com/moov-io/pinblock/formats"
	"github.com/stretchr/testify/require"
)

const (
	testPIN     = "1234"
	testAccount = "5432101234567891"
)

// allFormats returns every format, ISO-4 with the no-op cipher.
func allFormats() map[string]formats.Format {
	return map[string]formats.Format{
		"ISO-0": formats.NewISO0(), "ISO-1": formats.NewISO1(), "ISO-2": formats.NewISO2(),
		"ISO-3": formats.NewISO3(), "ISO-4": formats.NewISO4(encryption.NewNoOp()),
		"ANSI": formats.NewANSIX98(), "OEM-1": formats.NewOEM1(),
		"ECI-1": formats.NewECI1(), "ECI-2": formats.NewECI2(), "ECI-3": formats.NewECI3(), "ECI-4": formats.NewECI4(),
		"VISA-1": formats.NewVISA1(), "VISA-2": formats.NewVISA2(), "VISA-3": formats.NewVISA3(), "VISA-4": formats.NewVISA4(),
	}
}

func TestEncode_RejectsNonDecimalPIN(t *testing.T) {
	for name, f := range allFormats() {
		for _, pin := range []string{"abcd", "12A4", "12 4", "123456789O"} {
			_, err := f.Encode(pin, testAccount)
			require.ErrorContains(t, err, "pin must contain only decimal digits", "%s %q", name, pin)
		}
	}
}

func TestEncode_RejectsNonDecimalPAN(t *testing.T) {
	for _, f := range []formats.Format{formats.NewISO0(), formats.NewISO3(), formats.NewISO4(encryption.NewNoOp())} {
		_, err := f.Encode(testPIN, "54321012345678A1")
		require.ErrorContains(t, err, "account must contain only decimal digits")
		_, err = f.Decode(strings.Repeat("0", 32)[:len(mustEncode(t, f))], "54321012345678A1")
		require.ErrorContains(t, err, "account must contain only decimal digits")
	}
}

func mustEncode(t *testing.T, f formats.Format) string {
	t.Helper()
	pinBlock, err := f.Encode(testPIN, testAccount)
	require.NoError(t, err)
	return pinBlock
}

func TestDecode_RejectsNonHexBlock(t *testing.T) {
	for name, f := range allFormats() {
		pinBlock := mustEncode(t, f)
		_, err := f.Decode(pinBlock[:len(pinBlock)-1]+"G", testAccount)
		require.ErrorContains(t, err, "hex", name)
	}
}

func TestDecode_AcceptsLowerCaseHex(t *testing.T) {
	for name, f := range allFormats() {
		pin, err := f.Decode(strings.ToLower(mustEncode(t, f)), testAccount)
		require.NoError(t, err, name)
		require.Equal(t, testPIN, pin, name)
	}
}

func TestDecode_ChecksStructure(t *testing.T) {
	for _, tc := range []struct {
		name    string
		format  formats.Format
		block   string
		account string
		err     string
	}{
		// The format nibble must match (#5).
		{"ISO-0 given an ISO-3 block", formats.NewISO0(), "341215FEDCBA9876", testAccount, "format is different"},
		{"ISO-3 given an ISO-0 block", formats.NewISO3(), "041215FEDCBA9876", testAccount, "format is different"},
		{"ISO-2 given an ISO-1 block", formats.NewISO2(), "141234FFFFFFFFFF", "", "format is different"},
		{"ISO-1 given an ISO-2 block", formats.NewISO1(), "241234FFFFFFFFFF", "", "format is different"},
		// The fill each ISO format defines.
		{"ISO-0 fill not F", formats.NewISO0(), "041215FEDCBA9870", testAccount, "invalid fill"},
		{"ISO-3 fill below A", formats.NewISO3(), "3412345555555555", "0000000000000", "invalid fill"},
		{"ISO-2 fill not F", formats.NewISO2(), "241234FFFFFFFFFA", "", "invalid fill"},
		// PIN length 4 to 12.
		{"ISO-0 length 3", formats.NewISO0(), "03123FFFFFFFFFFF", "0000000000000", "invalid pin length"},
		{"ISO-1 length 3", formats.NewISO1(), "13123FFFFFFFFFFF", "", "invalid pin length"},
		{"ISO-1 length 14", formats.NewISO1(), "1E12345678901234", "", "invalid pin length"},
		{"ISO-2 length 13", formats.NewISO2(), "2D1234567890123F", "", "invalid pin length"},
		{"ECI-3 length 7", formats.NewECI3(), "7123456789012AAA", "", "unable to parse pin block"},
		{"OEM-1 three digits", formats.NewOEM1(), "123AAAAAAAAAAAAA", "", "invalid pin length"},
		// PIN digits must be decimal.
		{"ISO-0 PIN digit A", formats.NewISO0(), "04121BFEDCBA9876", testAccount, "non-decimal"},
		{"ISO-1 PIN digit A", formats.NewISO1(), "14123AFFFFFFFFFF", "", "non-decimal"},
		{"ECI-2 PIN digit A", formats.NewECI2(), "123A555555555555", "", "non-decimal"},
		{"VISA-3 PIN digit A", formats.NewVISA3(), "12A4F55555555555", "", "non-decimal"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.format.Decode(tc.block, tc.account)
			require.ErrorContains(t, err, tc.err)
		})
	}
}

// ISO 9564-1 allows 4 to 12 PIN digits; ISO-1 and ISO-2 refuse a longer PIN
// rather than truncating it, as the other ISO formats already did.
func TestISO1AndISO2_RefuseLongPIN(t *testing.T) {
	for _, f := range []formats.Format{formats.NewISO1(), formats.NewISO2(), formats.NewECI4()} {
		_, err := f.Encode("1234567890123", "")
		require.ErrorContains(t, err, "pin length must be between 4 and 12 digits")
		_, err = f.Encode("123", "")
		require.Error(t, err)
	}
}

// Format 1 fill is random 0x0 to 0xF. Across 20 blocks (200 fill nibbles) at
// least one is A to F unless the fill is limited to decimal digits.
func TestISO1_FillIsRandomHex(t *testing.T) {
	sawLetter := false
	for range 20 {
		pinBlock, err := formats.NewISO1().Encode(testPIN, "")
		require.NoError(t, err)
		if strings.ContainsAny(pinBlock[6:], "ABCDEF") {
			sawLetter = true
		}
	}
	require.True(t, sawLetter, "ISO-1 fill never used A-F")
}

// VISA-2 fills with decimal digits only.
func TestVISA2_FillIsDecimal(t *testing.T) {
	for range 20 {
		pinBlock, err := formats.NewVISA2().Encode(testPIN, "")
		require.NoError(t, err)
		require.Equal(t, -1, strings.IndexAny(pinBlock[7:], "ABCDEF"), pinBlock)
	}
}

// Every format decodes what it encodes, keeping as many leading digits as
// it holds.
func TestRoundTrip_AllFormats(t *testing.T) {
	keep := map[string]int{"ECI-2": 4, "ECI-3": 6, "VISA-2": 6}
	for name, f := range allFormats() {
		for _, pin := range []string{"1234", "98765", "1234567890", "123456789012"} {
			pinBlock, err := f.Encode(pin, testAccount)
			require.NoError(t, err, name)
			got, err := f.Decode(pinBlock, testAccount)
			require.NoError(t, err, name)
			want := pin
			if n, ok := keep[name]; ok {
				want = pin[:min(len(pin), n)]
			}
			require.Equal(t, want, got, "%s %s", name, pin)
		}
	}
}

// The debug writer names the operation and the format, and never the PIN, the
// PAN or the PIN block.
func TestSetDebugWriter_WritesNothingSensitive(t *testing.T) {
	const pin, account = "987654", "4111111111111111"
	for name, f := range allFormats() {
		var out bytes.Buffer
		f.SetDebugWriter(&out)

		pinBlock, err := f.Encode(pin, account)
		require.NoError(t, err, name)
		got, err := f.Decode(pinBlock, account)
		require.NoError(t, err, name)

		log := out.String()
		require.Contains(t, log, "PIN block encode operation finished", name)
		require.Contains(t, log, "PIN block decode operation finished", name)
		for _, secret := range []string{pin, got, account, account[3:15], pinBlock, strings.ToLower(pinBlock)} {
			require.NotContains(t, log, secret, name)
		}
	}
}
