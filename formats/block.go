package formats

import (
	"crypto/rand"
	"fmt"

	"github.com/moov-io/pinblock/internal/memzero"
)

// wipe clears a buffer that held a PIN, a PAN, or a clear or intermediate PIN
// block. It is a variable only so tests can observe which buffers the package
// clears.
var wipe = memzero.Bytes

// fieldNibbles is the number of 4-bit fields in a 64-bit PIN or PAN field.
const fieldNibbles = 16

// Tables of nibble values used for random fill.
var (
	anyNibble     = []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF}
	decimalNibble = []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9}
	letterNibble  = []byte{0xA, 0xB, 0xC, 0xD, 0xE, 0xF}
)

// digitNibbles checks that every byte of s is an ASCII decimal digit and writes
// its value into dst, which must be at least len(s) long. name is used in the
// error, which never includes the value.
func digitNibbles(dst, s []byte, name string) error {
	for i, c := range s {
		if c < '0' || c > '9' {
			return fmt.Errorf("%s must contain only decimal digits", name)
		}
		dst[i] = c - '0'
	}
	return nil
}

// checkDigits reports an error unless every byte of s is an ASCII decimal digit.
func checkDigits(s []byte, name string) error {
	for _, c := range s {
		if c < '0' || c > '9' {
			return fmt.Errorf("%s must contain only decimal digits", name)
		}
	}
	return nil
}

// pinFromNibbles returns the PIN held in nibbles as ASCII digits, or an error
// if any nibble is not a decimal digit. The result belongs to the caller.
func pinFromNibbles(nibbles []byte) ([]byte, error) {
	pin := make([]byte, len(nibbles))
	for i, n := range nibbles {
		if n > 9 {
			wipe(pin)
			return nil, fmt.Errorf("pin block holds a non-decimal PIN digit")
		}
		pin[i] = '0' + n
	}
	return pin, nil
}

// pack packs nibbles (one 4-bit value per byte) into bytes, two per byte, high
// nibble first. len(nibbles) must be even.
func pack(nibbles []byte) []byte {
	out := make([]byte, len(nibbles)/2)
	for i := range out {
		out[i] = nibbles[2*i]<<4 | nibbles[2*i+1]&0x0F
	}
	return out
}

// unpack splits each byte of b into two nibbles, high nibble first. The result
// is a new slice; b is not modified.
func unpack(b []byte) []byte {
	out := make([]byte, 2*len(b))
	for i, c := range b {
		out[2*i] = c >> 4
		out[2*i+1] = c & 0x0F
	}
	return out
}

// randomFill sets every element of dst to a value drawn uniformly at random
// from table, using crypto/rand. Draws that would bias the result toward the
// start of table are rejected rather than reduced modulo len(table).
func randomFill(dst, table []byte) error {
	limit := 256 - 256%len(table)
	buf := make([]byte, len(dst)+16)
	defer wipe(buf)

	filled := 0
	for filled < len(dst) {
		if _, err := rand.Read(buf); err != nil {
			return fmt.Errorf("reading random fill: %w", err)
		}
		for _, r := range buf {
			if int(r) >= limit {
				continue
			}
			dst[filled] = table[int(r)%len(table)]
			filled++
			if filled == len(dst) {
				break
			}
		}
	}
	return nil
}

// panFieldISO0 returns the 16-nibble PAN field used by formats 0 and 3: four
// zero nibbles followed by the 12 rightmost PAN digits excluding the check
// digit (ISO 9564-1). account must hold at least 13 decimal digits. The result
// belongs to the caller, who must clear it.
func panFieldISO0(account []byte) ([]byte, error) {
	if len(account) < 13 {
		return nil, fmt.Errorf("account length must be at least 13 digits")
	}
	field := make([]byte, fieldNibbles)
	if err := digitNibbles(field[4:], account[len(account)-13:len(account)-1], "account"); err != nil {
		wipe(field)
		return nil, err
	}
	return field, nil
}

// xorInto sets dst[i] = a[i] ^ b[i]. All three must be the same length.
func xorInto(dst, a, b []byte) {
	for i := range dst {
		dst[i] = a[i] ^ b[i]
	}
}
