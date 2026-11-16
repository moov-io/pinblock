package formats

import (
	"bytes"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/moov-io/pinblock/encryption"
	"github.com/stretchr/testify/require"
)

// wipeRecorder replaces the package's wipe hook for one test. It records what
// each buffer held just before it was cleared, and the buffer itself, so a test
// can check both that a given intermediate was cleared and that the clear took
// effect.
type wipeRecorder struct {
	before  [][]byte
	buffers [][]byte
}

func recordWipes(t *testing.T) *wipeRecorder {
	t.Helper()
	r := &wipeRecorder{}
	original := wipe
	wipe = func(b []byte) {
		r.before = append(r.before, bytes.Clone(b))
		r.buffers = append(r.buffers, b)
		original(b)
	}
	t.Cleanup(func() { wipe = original })
	return r
}

// count returns how many cleared buffers held exactly value.
func (r *wipeRecorder) count(value []byte) int {
	n := 0
	for _, b := range r.before {
		if bytes.Equal(b, value) {
			n++
		}
	}
	return n
}

// countLen returns how many cleared buffers were n bytes long.
func (r *wipeRecorder) countLen(n int) int {
	c := 0
	for _, b := range r.before {
		if len(b) == n {
			c++
		}
	}
	return c
}

// requireAllZero fails unless every buffer passed to wipe is still all zero
// when the function under test has returned.
func (r *wipeRecorder) requireAllZero(t *testing.T) {
	t.Helper()
	require.NotEmpty(t, r.buffers)
	for i, b := range r.buffers {
		require.Equal(t, make([]byte, len(b)), b, "buffer %d is not zero after return", i)
	}
}

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	require.NoError(t, err)
	return b
}

var (
	wipePIN     = []byte("1234")
	wipeAccount = []byte("5432101234567891")
	// The ISO-0 PAN field for wipeAccount, one nibble per byte.
	wipePANField = []byte{0, 0, 0, 0, 2, 1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
)

func TestISO0_ClearsFields(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		r := recordWipes(t)
		pinBlock, err := NewISO0().Encode(wipePIN, wipeAccount)
		require.NoError(t, err)
		require.Equal(t, 1, r.count(wipePANField), "PAN field")
		// The PIN field is XORed in place, so on return it holds the block.
		require.Equal(t, 1, r.count(unpack(pinBlock)), "PIN field")
		r.requireAllZero(t)
	})
	t.Run("decode", func(t *testing.T) {
		r := recordWipes(t)
		pin, err := NewISO0().Decode(mustHex(t, "041215FEDCBA9876"), wipeAccount)
		require.NoError(t, err)
		require.Equal(t, wipePIN, pin)
		require.Equal(t, 1, r.count(wipePANField), "PAN field")
		require.Equal(t, 1, r.count([]byte{0, 4, 1, 2, 3, 4, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15}), "clear PIN field")
		r.requireAllZero(t)
	})
	t.Run("encode error after the PIN field is started", func(t *testing.T) {
		r := recordWipes(t)
		_, err := NewISO0().Encode([]byte("123A"), wipeAccount)
		require.Error(t, err)
		require.Equal(t, 1, r.count([]byte{0, 4, 1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), "partial PIN field")
		r.requireAllZero(t)
	})
	t.Run("decode error on a non-decimal PIN digit", func(t *testing.T) {
		r := recordWipes(t)
		// clear PIN field 04123AFFFFFFFFFF
		_, err := NewISO0().Decode(mustHex(t, "04121BFEDCBA9876"), wipeAccount)
		require.Error(t, err)
		require.Equal(t, 1, r.count([]byte{'1', '2', '3', 0}), "partial PIN")
		require.Equal(t, 1, r.count([]byte{0, 4, 1, 2, 3, 0xA, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15}), "clear PIN field")
		r.requireAllZero(t)
	})
	t.Run("ISO-3 random fill", func(t *testing.T) {
		r := recordWipes(t)
		_, err := NewISO3().Encode(wipePIN, wipeAccount)
		require.NoError(t, err)
		require.Equal(t, 1, r.countLen(10+16), "random fill buffer")
		r.requireAllZero(t)
	})
}

func TestISO1_ClearsFields(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		r := recordWipes(t)
		pinBlock, err := NewISO2().Encode(wipePIN, nil)
		require.NoError(t, err)
		require.Equal(t, 1, r.count(unpack(pinBlock)), "clear PIN field")
		r.requireAllZero(t)
	})
	t.Run("decode", func(t *testing.T) {
		r := recordWipes(t)
		pinBlock := mustHex(t, "241234FFFFFFFFFF")
		_, err := NewISO2().Decode(pinBlock, nil)
		require.NoError(t, err)
		require.Equal(t, 1, r.count(unpack(pinBlock)), "clear PIN field")
		r.requireAllZero(t)
	})
}

func TestISO4_ClearsFields(t *testing.T) {
	account := []byte("432198765432109870")
	panNibbles := append([]byte{6, 4, 3, 2, 1, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 9, 8, 7, 0}, make([]byte, 13)...)
	panField := pack(panNibbles)

	// With the no-op cipher every intermediate can be computed from the
	// result: block A is the plain PIN field P, and the result is P XOR PAN.
	t.Run("encode", func(t *testing.T) {
		r := recordWipes(t)
		result, err := NewISO4(encryption.NewNoOp()).Encode(wipePIN, account)
		require.NoError(t, err)
		p := make([]byte, 16)
		xorInto(p, result, panField)

		require.Equal(t, 1, r.count(panNibbles), "PAN field nibbles")
		require.Equal(t, 1, r.count(panField), "PAN field")
		require.Equal(t, 1, r.count(unpack(p)), "clear PIN field nibbles")
		require.Equal(t, 3, r.count(p), "clear PIN field, block A and the cipher's copy")
		require.Equal(t, 2, r.count(result), "block B and the cipher's copy")
		r.requireAllZero(t)
	})
	t.Run("decode", func(t *testing.T) {
		pinBlock, err := NewISO4(encryption.NewNoOp()).Encode(wipePIN, account)
		require.NoError(t, err)
		p := make([]byte, 16)
		xorInto(p, pinBlock, panField)
		before := bytes.Clone(pinBlock)

		r := recordWipes(t)
		pin, err := NewISO4(encryption.NewNoOp()).Decode(pinBlock, account)
		require.NoError(t, err)
		require.Equal(t, wipePIN, pin)
		require.Equal(t, before, pinBlock, "the caller's block must not be modified")

		require.Equal(t, 1, r.count(panNibbles), "PAN field nibbles")
		require.Equal(t, 1, r.count(panField), "PAN field")
		require.Equal(t, 3, r.count(pinBlock), "copy of the block, the cipher's output and block B")
		require.Equal(t, 3, r.count(p), "block A, the cipher's output and the clear PIN field")
		require.Equal(t, 1, r.count(unpack(p)), "clear PIN field nibbles")
		r.requireAllZero(t)
	})
}

func TestProprietaryFormats_ClearFields(t *testing.T) {
	for _, tc := range []struct {
		name   string
		format Format
	}{
		{"OEM-1", NewOEM1()},
		{"ECI-2", NewECI2()},
		{"ECI-3", NewECI3()},
		{"VISA-2", NewVISA2()},
		{"VISA-3", NewVISA3()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := recordWipes(t)
			pinBlock, err := tc.format.Encode(wipePIN, nil)
			require.NoError(t, err)
			require.Equal(t, 1, r.count(unpack(pinBlock)), "clear PIN field on encode")
			r.requireAllZero(t)

			r = recordWipes(t)
			pin, err := tc.format.Decode(pinBlock, nil)
			require.NoError(t, err)
			require.Equal(t, wipePIN, pin)
			require.Equal(t, 1, r.count(unpack(pinBlock)), "clear PIN field on decode")
			r.requireAllZero(t)
		})
	}
	t.Run("VISA-3 filler", func(t *testing.T) {
		r := recordWipes(t)
		pinBlock, err := NewVISA3().Encode(wipePIN, nil)
		require.NoError(t, err)
		require.Equal(t, 1, r.count([]byte{pinBlock[7] & 0x0F}), "filler value")
		r.requireAllZero(t)
	})
}

// The functions never modify their inputs, and the result never shares the
// caller's arrays.
func TestInputsAreNotModified(t *testing.T) {
	for name, f := range map[string]Format{
		"ISO-0": NewISO0(), "ISO-1": NewISO1(), "ISO-2": NewISO2(), "ISO-3": NewISO3(),
		"ISO-4": NewISO4(encryption.NewNoOp()), "OEM-1": NewOEM1(), "ECI-2": NewECI2(),
		"ECI-3": NewECI3(), "VISA-2": NewVISA2(), "VISA-3": NewVISA3(),
	} {
		pin := []byte("123456")
		account := []byte("5432101234567891")
		pinBlock, err := f.Encode(pin, account)
		require.NoError(t, err, name)
		require.Equal(t, []byte("123456"), pin, name)
		require.Equal(t, []byte("5432101234567891"), account, name)

		blockCopy := bytes.Clone(pinBlock)
		got, err := f.Decode(pinBlock, account)
		require.NoError(t, err, name)
		require.Equal(t, blockCopy, pinBlock, name)

		clear(got)
		again, err := f.Decode(pinBlock, account)
		require.NoError(t, err, name)
		require.NotEqual(t, make([]byte, len(again)), again, "%s: the result must not share a buffer the library reuses", name)
	}
}

func TestPANField_ClearedOnError(t *testing.T) {
	r := recordWipes(t)
	_, err := NewISO0().Encode(wipePIN, []byte("54321012345678A1"))
	require.ErrorContains(t, err, "account must contain only decimal digits")
	require.Equal(t, 1, r.count([]byte{0, 0, 0, 0, 2, 1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 0}), "partial PAN field")
	r.requireAllZero(t)
}

// failingCipher returns a buffer along with an error, as a cipher might on a
// partial failure. The buffer is cleared.
type failingCipher struct{}

func (failingCipher) Encrypt([]byte) ([]byte, error) { return []byte{1, 2, 3}, errTest }
func (failingCipher) Decrypt([]byte) ([]byte, error) { return []byte{1, 2, 3}, errTest }

var errTest = errors.New("cipher failed")

func TestISO4_ClearsCipherOutputOnError(t *testing.T) {
	account := []byte("432198765432109870")

	r := recordWipes(t)
	_, err := NewISO4(failingCipher{}).Encode(wipePIN, account)
	require.ErrorIs(t, err, errTest)
	require.Equal(t, 1, r.count([]byte{1, 2, 3}), "output of a failed cipher call")
	r.requireAllZero(t)

	r = recordWipes(t)
	_, err = NewISO4(shortCipher{}).Encode(wipePIN, account)
	require.ErrorContains(t, err, "cipher returned 8 bytes")
	require.Equal(t, 1, r.countLen(8), "short cipher output")
	r.requireAllZero(t)
}
