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

// nibblesOf returns the hex text s as one nibble per byte.
func nibblesOf(t *testing.T, s string) []byte {
	t.Helper()
	return unpack(mustHex(t, s))
}

// requireDebugNote checks that the debug writer received exactly the note
// SetDebugWriter promises: the operation and the format, and nothing else.
func requireDebugNote(t *testing.T, got, operation, format string) {
	t.Helper()
	require.Equal(t, "PIN block "+operation+" operation finished\nFormat  : "+format+"\n\n", got)
}

const (
	wipePIN     = "1234"
	wipeAccount = "5432101234567891"
)

// The ISO-0 PAN field for wipeAccount, one nibble per byte.
var wipePANField = []byte{0, 0, 0, 0, 2, 1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

func TestISO0_ClearsFields(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		r := recordWipes(t)
		pinBlock, err := NewISO0().Encode(wipePIN, wipeAccount)
		require.NoError(t, err)
		require.Equal(t, "041215FEDCBA9876", pinBlock)
		require.Equal(t, 1, r.count(wipePANField), "PAN field")
		// The PIN field is XORed in place, so on return it holds the block.
		require.Equal(t, 1, r.count(nibblesOf(t, pinBlock)), "PIN field")
		require.Equal(t, 1, r.count([]byte(pinBlock)), "hex text of the block")
		r.requireAllZero(t)
	})
	t.Run("decode", func(t *testing.T) {
		r := recordWipes(t)
		pin, err := NewISO0().Decode("041215FEDCBA9876", wipeAccount)
		require.NoError(t, err)
		require.Equal(t, wipePIN, pin)
		require.Equal(t, 1, r.count(wipePANField), "PAN field")
		require.Equal(t, 1, r.count([]byte{0, 4, 1, 2, 3, 4, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15}), "clear PIN field")
		require.Equal(t, 1, r.count([]byte(wipePIN)), "PIN text")
		r.requireAllZero(t)
	})
	t.Run("encode error after the PIN field is started", func(t *testing.T) {
		r := recordWipes(t)
		_, err := NewISO0().Encode("123A", wipeAccount)
		require.Error(t, err)
		require.Equal(t, 1, r.count([]byte{0, 4, 1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), "partial PIN field")
		r.requireAllZero(t)
	})
	t.Run("decode error on a non-decimal PIN digit", func(t *testing.T) {
		r := recordWipes(t)
		// clear PIN field 04123AFFFFFFFFFF
		_, err := NewISO0().Decode("04121BFEDCBA9876", wipeAccount)
		require.ErrorContains(t, err, "non-decimal")
		require.Equal(t, 1, r.count([]byte{'1', '2', '3', 0}), "partial PIN")
		require.Equal(t, 1, r.count([]byte{0, 4, 1, 2, 3, 0xA, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15}), "clear PIN field")
		r.requireAllZero(t)
	})
	t.Run("decode error on a non-hex block", func(t *testing.T) {
		r := recordWipes(t)
		_, err := NewISO0().Decode("041215FEDCBA987G", wipeAccount)
		require.ErrorContains(t, err, "hex")
		require.Equal(t, 1, r.count([]byte{0, 4, 1, 2, 1, 5, 0xF, 0xE, 0xD, 0xC, 0xB, 0xA, 9, 8, 7, 0}), "partial block")
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
		pinBlock, err := NewISO1().Encode(wipePIN, "")
		require.NoError(t, err)
		require.Equal(t, 1, r.count(nibblesOf(t, pinBlock)), "clear PIN field")
		require.Equal(t, 1, r.count([]byte(pinBlock)), "hex text of the block")
		require.Equal(t, 1, r.countLen(10+16), "random fill buffer")
		r.requireAllZero(t)
	})
	t.Run("ISO-2 encode", func(t *testing.T) {
		r := recordWipes(t)
		pinBlock, err := NewISO2().Encode(wipePIN, "")
		require.NoError(t, err)
		require.Equal(t, "241234FFFFFFFFFF", pinBlock)
		require.Equal(t, 1, r.count(nibblesOf(t, pinBlock)), "clear PIN field")
		r.requireAllZero(t)
	})
	t.Run("decode", func(t *testing.T) {
		r := recordWipes(t)
		pinBlock := "241234FFFFFFFFFF"
		pin, err := NewISO2().Decode(pinBlock, "")
		require.NoError(t, err)
		require.Equal(t, wipePIN, pin)
		require.Equal(t, 1, r.count(nibblesOf(t, pinBlock)), "clear PIN field")
		require.Equal(t, 1, r.count([]byte(wipePIN)), "PIN text")
		r.requireAllZero(t)
	})
}

func TestISO4_ClearsFields(t *testing.T) {
	account := "432198765432109870"
	panNibbles := append([]byte{6, 4, 3, 2, 1, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 9, 8, 7, 0}, make([]byte, 13)...)
	panField := pack(panNibbles)

	// With the no-op cipher every intermediate can be computed from the
	// result: block A is the plain PIN field P, and the result is P XOR PAN.
	t.Run("encode", func(t *testing.T) {
		r := recordWipes(t)
		pinBlock, err := NewISO4(encryption.NewNoOp()).Encode(wipePIN, account)
		require.NoError(t, err)
		result := mustHex(t, pinBlock)
		p := make([]byte, 16)
		xorInto(p, result, panField)

		require.Equal(t, 1, r.count(panNibbles), "PAN field nibbles")
		require.Equal(t, 1, r.count(panField), "PAN field")
		require.Equal(t, 1, r.count(unpack(p)), "clear PIN field nibbles")
		// 32 bytes: the PAN field and PIN field nibbles, the random fill
		// buffer (16 nibbles plus 16 spare draws), the enciphered block's
		// nibbles and its hex text.
		require.Equal(t, 5, r.countLen(32), "buffers of 32 bytes, random fill among them")
		require.Equal(t, 3, r.count(p), "clear PIN field, block A and the cipher's output")
		require.Equal(t, 3, r.count(result), "block B, the cipher's output and the enciphered block")
		require.Equal(t, 1, r.count(unpack(result)), "enciphered block nibbles")
		require.Equal(t, 1, r.count([]byte(pinBlock)), "hex text of the block")
		r.requireAllZero(t)
	})
	t.Run("decode", func(t *testing.T) {
		pinBlock, err := NewISO4(encryption.NewNoOp()).Encode(wipePIN, account)
		require.NoError(t, err)
		encrypted := mustHex(t, pinBlock)
		p := make([]byte, 16)
		xorInto(p, encrypted, panField)

		r := recordWipes(t)
		pin, err := NewISO4(encryption.NewNoOp()).Decode(pinBlock, account)
		require.NoError(t, err)
		require.Equal(t, wipePIN, pin)

		require.Equal(t, 1, r.count(panNibbles), "PAN field nibbles")
		require.Equal(t, 1, r.count(panField), "PAN field")
		require.Equal(t, 1, r.count(unpack(encrypted)), "nibbles of the block")
		require.Equal(t, 3, r.count(encrypted), "the block, the cipher's output and block B")
		require.Equal(t, 3, r.count(p), "block A, the cipher's output and the clear PIN field")
		require.Equal(t, 1, r.count(unpack(p)), "clear PIN field nibbles")
		require.Equal(t, 1, r.count([]byte(wipePIN)), "PIN text")
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
			pinBlock, err := tc.format.Encode(wipePIN, "")
			require.NoError(t, err)
			require.Equal(t, 1, r.count(nibblesOf(t, pinBlock)), "clear PIN field on encode")
			require.Equal(t, 1, r.count([]byte(pinBlock)), "hex text of the block")
			r.requireAllZero(t)

			r = recordWipes(t)
			pin, err := tc.format.Decode(pinBlock, "")
			require.NoError(t, err)
			require.Equal(t, wipePIN, pin)
			require.Equal(t, 1, r.count(nibblesOf(t, pinBlock)), "clear PIN field on decode")
			require.Equal(t, 1, r.count([]byte(wipePIN)), "PIN text")
			r.requireAllZero(t)
		})
	}
	t.Run("VISA-3 filler", func(t *testing.T) {
		r := recordWipes(t)
		pinBlock, err := NewVISA3().Encode(wipePIN, "")
		require.NoError(t, err)
		filler := nibblesOf(t, pinBlock)[15]
		require.Equal(t, 1, r.count([]byte{filler}), "filler value")
		r.requireAllZero(t)
	})
}

func TestPANField_ClearedOnError(t *testing.T) {
	r := recordWipes(t)
	_, err := NewISO0().Encode(wipePIN, "54321012345678A1")
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

// shortCipher returns the first half of its input, a slice of the buffer it
// was given, so its capacity is the whole PIN field.
type shortCipher struct{}

func (shortCipher) Encrypt(in []byte) ([]byte, error) { return bytes.Clone(in)[:8], nil }
func (shortCipher) Decrypt(in []byte) ([]byte, error) { return bytes.Clone(in)[:8], nil }

func TestISO4_ClearsCipherOutputOnError(t *testing.T) {
	account := "432198765432109870"

	r := recordWipes(t)
	_, err := NewISO4(failingCipher{}).Encode(wipePIN, account)
	require.ErrorIs(t, err, errTest)
	require.Equal(t, 1, r.count([]byte{1, 2, 3}), "output of a failed cipher call")
	r.requireAllZero(t)

	// The whole capacity behind a short result is cleared.
	r = recordWipes(t)
	_, err = NewISO4(shortCipher{}).Encode(wipePIN, account)
	require.ErrorContains(t, err, "cipher returned 8 bytes")
	require.Zero(t, r.countLen(8), "only len(out) was cleared")
	r.requireAllZero(t)

	r = recordWipes(t)
	_, err = NewISO4(failingCipher{}).Decode("00000000000000000000000000000000", account)
	require.ErrorIs(t, err, errTest)
	require.Equal(t, 1, r.count([]byte{1, 2, 3}), "output of a failed cipher call")
	r.requireAllZero(t)
}

// scratchCipher returns its result as the first block of a larger buffer it
// owns, with key-derived bytes after it.
type scratchCipher struct{ last []byte }

func (c *scratchCipher) op(in []byte) ([]byte, error) {
	buf := bytes.Repeat([]byte{0xA5}, 64)
	copy(buf, in)
	c.last = buf
	return buf[:16], nil
}

// apply clears the whole capacity a cipher returns, not just the first block.
func TestApply_ClearsFullCapacity(t *testing.T) {
	c := &scratchCipher{}
	r := recordWipes(t)
	out, err := apply(c.op, make([]byte, 16), "test")
	require.NoError(t, err)
	require.Len(t, out, 16)
	require.Equal(t, make([]byte, 64), c.last, "the cipher's whole buffer is cleared")
	require.Equal(t, 1, r.countLen(64))
}

// randomFill rejects draws at or above the largest multiple of the table
// size, so every table entry is equally likely.
func TestRandomFill_RejectsBiasedDraws(t *testing.T) {
	// For the 10-entry table, 250 to 255 are rejected; 0, 11 and 249 map to
	// 0, 1 and 9.
	src := append([]byte{250, 251, 252, 253, 254, 255, 0, 11, 249}, make([]byte, 64)...)
	dst := make([]byte, 3)
	require.NoError(t, randomFillFrom(bytes.NewReader(src), dst, decimalNibble))
	require.Equal(t, []byte{0, 1, 9}, dst)

	// A short read is an error, not a partial fill.
	err := randomFillFrom(bytes.NewReader([]byte{1, 2}), make([]byte, 3), decimalNibble)
	require.ErrorContains(t, err, "reading random fill")
}
