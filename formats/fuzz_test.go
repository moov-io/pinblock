package formats_test

import (
	"testing"

	"github.com/moov-io/pinblock/encryption"
	"github.com/moov-io/pinblock/formats"
)

// fuzzFormats is every format with the number of leading PIN digits it keeps.
var fuzzFormats = []struct {
	ctor func() formats.Format
	keep int
}{
	{formats.NewISO0, 12},
	{formats.NewISO1, 12},
	{formats.NewISO2, 12},
	{formats.NewISO3, 12},
	{func() formats.Format { return formats.NewISO4(encryption.NewNoOp()) }, 12},
	{formats.NewANSIX98, 12},
	{formats.NewOEM1, 12},
	{formats.NewECI1, 12},
	{formats.NewECI2, 4},
	{formats.NewECI3, 6},
	{formats.NewECI4, 12},
	{formats.NewVISA1, 12},
	{formats.NewVISA2, 6},
	{formats.NewVISA3, 12},
	{formats.NewVISA4, 12},
}

func FuzzISO0EncodeDecode(f *testing.F) {
	// Real test vectors from iso0_test.go
	seeds := []struct{ pin, pan string }{
		{"1234", "5432101234567891"},
		{"123456789012", "4000000000000002"},
		{"9", "4111111111111111"},
		{"", ""},
		{"1234", "123"},
		{"abcdef", "5432101234567891"},
	}
	for _, s := range seeds {
		f.Add(s.pin, s.pan)
	}

	f.Fuzz(func(t *testing.T, pin, pan string) {
		if len(pin) > 32 || len(pan) > 32 {
			t.Skip()
		}

		for _, tc := range fuzzFormats {
			fmtter := tc.ctor()
			block, err := fmtter.Encode(pin, pan)
			if err != nil {
				continue
			}
			requireUpperHex(t, block)
			// Whatever Encode accepts, Decode must give back.
			got, err := fmtter.Decode(block, pan)
			if err != nil {
				t.Fatalf("decode of an encoded block failed: %v", err)
			}
			if want := pin[:min(len(pin), tc.keep)]; got != want {
				t.Fatalf("round trip: got %q, want %q", got, want)
			}
		}
	})
}

func FuzzDecodeOnly(f *testing.F) {
	f.Add("041215FEDCBA9876", "5432101234567891")
	f.Add("", "")
	f.Add("0000000000000000", "4111111111111111")
	f.Add("FFFFFFFFFFFFFFFF", "4000000000000002")

	f.Fuzz(func(t *testing.T, block, pan string) {
		if len(block) > 64 || len(pan) > 32 {
			t.Skip()
		}

		for _, tc := range fuzzFormats {
			pin, err := tc.ctor().Decode(block, pan)
			if err != nil {
				continue
			}
			if len(pin) < 4 || len(pin) > 12 {
				t.Fatalf("decoded a PIN of %d digits", len(pin))
			}
			for i := 0; i < len(pin); i++ {
				if pin[i] < '0' || pin[i] > '9' {
					t.Fatalf("decoded a non-decimal PIN %q", pin)
				}
			}
		}
	})
}

// requireUpperHex fails unless block is upper-case hex text of a PIN block's
// length.
func requireUpperHex(t *testing.T, block string) {
	t.Helper()
	if len(block) != 16 && len(block) != 32 {
		t.Fatalf("pin block %q has length %d", block, len(block))
	}
	for i := 0; i < len(block); i++ {
		c := block[i]
		if (c < '0' || c > '9') && (c < 'A' || c > 'F') {
			t.Fatalf("pin block %q is not upper-case hex", block)
		}
	}
}
