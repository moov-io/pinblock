package formats_test

import (
	"bytes"
	"testing"

	"github.com/moov-io/pinblock/formats"
)

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
		f.Add([]byte(s.pin), []byte(s.pan))
	}

	f.Fuzz(func(t *testing.T, pin, pan []byte) {
		if len(pin) > 32 || len(pan) > 32 {
			t.Skip()
		}
		pinIn, panIn := bytes.Clone(pin), bytes.Clone(pan)

		for _, tc := range []struct {
			ctor func() formats.Format
			// keep is how many leading PIN digits the format preserves
			keep int
		}{
			{formats.NewISO0, 12},
			{formats.NewISO1, 12},
			{formats.NewISO2, 12},
			{formats.NewISO3, 12},
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
		} {
			fmtter := tc.ctor()
			block, err := fmtter.Encode(pin, pan)
			if err != nil {
				continue
			}
			// Whatever Encode accepts, Decode must give back.
			got, err := fmtter.Decode(block, pan)
			if err != nil {
				t.Fatalf("decode of an encoded block failed: %v", err)
			}
			want := pin[:min(len(pin), tc.keep)]
			if !bytes.Equal(got, want) {
				t.Fatalf("round trip: got %q, want %q", got, want)
			}
		}
		if !bytes.Equal(pin, pinIn) || !bytes.Equal(pan, panIn) {
			t.Fatal("an input was modified")
		}
	})
}

func FuzzDecodeOnly(f *testing.F) {
	f.Add(block(f, "041215FEDCBA9876"), []byte("5432101234567891"))
	f.Add([]byte{}, []byte{})
	f.Add(block(f, "0000000000000000"), []byte("4111111111111111"))
	f.Add(block(f, "FFFFFFFFFFFFFFFF"), []byte("4000000000000002"))

	f.Fuzz(func(t *testing.T, pinBlock, pan []byte) {
		if len(pinBlock) > 64 || len(pan) > 32 {
			t.Skip()
		}
		for _, ctor := range []func() formats.Format{
			formats.NewISO0, formats.NewISO1, formats.NewISO2, formats.NewISO3,
			formats.NewANSIX98, formats.NewOEM1, formats.NewECI1, formats.NewECI2,
			formats.NewECI3, formats.NewECI4, formats.NewVISA1, formats.NewVISA2,
			formats.NewVISA3, formats.NewVISA4,
		} {
			pin, err := ctor().Decode(pinBlock, pan)
			if err != nil {
				continue
			}
			for _, c := range pin {
				if c < '0' || c > '9' {
					t.Fatalf("decoded a non-decimal PIN %q", pin)
				}
			}
		}
	})
}
