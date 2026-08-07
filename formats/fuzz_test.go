package formats_test

import (
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
		f.Add(s.pin, s.pan)
	}

	f.Fuzz(func(t *testing.T, pin, pan string) {
		if len(pin) > 32 || len(pan) > 32 {
			t.Skip()
		}

		for _, ctor := range []func() formats.Format{
			formats.NewISO0,
			formats.NewISO1,
			formats.NewISO2,
			formats.NewISO3,
			formats.NewANSIX98,
			formats.NewOEM1,
			formats.NewECI1,
			formats.NewECI2,
			formats.NewECI3,
			formats.NewECI4,
			formats.NewVISA1,
			formats.NewVISA2,
			formats.NewVISA3,
			formats.NewVISA4,
		} {
			fmtter := ctor()
			block, err := fmtter.Encode(pin, pan)
			if err != nil {
				continue
			}
			_, _ = fmtter.Decode(block, pan)
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

		for _, ctor := range []func() formats.Format{
			formats.NewISO0,
			formats.NewISO1,
			formats.NewISO3,
			formats.NewANSIX98,
			formats.NewVISA1,
		} {
			_, _ = ctor().Decode(block, pan)
		}
	})
}
