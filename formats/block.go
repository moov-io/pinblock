package formats

import (
	"crypto/rand"
	"fmt"
	"io"

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

// upperHex maps a nibble to its upper-case hex character.
const upperHex = "0123456789ABCDEF"

// digitNibbles checks that every byte of s is an ASCII decimal digit and writes
// its value into dst, which must be at least len(s) long. name is used in the
// error, which never includes the value.
func digitNibbles(dst []byte, s, name string) error {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return fmt.Errorf("%s must contain only decimal digits", name)
		}
		dst[i] = c - '0'
	}
	return nil
}

// checkDigits reports an error unless every byte of s is an ASCII decimal digit.
func checkDigits(s, name string) error {
	for i := 0; i < len(s); i++ {
		if c := s[i]; c < '0' || c > '9' {
			return fmt.Errorf("%s must contain only decimal digits", name)
		}
	}
	return nil
}

// pinFromNibbles returns the PIN held in nibbles as a string of decimal
// digits, or an error if any nibble is not a decimal digit. The PIN is built in
// a buffer that is cleared before returning; the returned string cannot be.
func pinFromNibbles(nibbles []byte) (string, error) {
	pin := make([]byte, len(nibbles))
	defer wipe(pin)
	for i, n := range nibbles {
		if n > 9 {
			return "", fmt.Errorf("pin block holds a non-decimal PIN digit")
		}
		pin[i] = '0' + n
	}
	return string(pin), nil
}

// hexNibbles parses s, hex text in either case, into one nibble per byte. The
// result belongs to the caller, who must clear it. On error nothing is
// returned and the partial result is cleared.
func hexNibbles(s string) ([]byte, error) {
	out := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
			out[i] = c - '0'
		case c >= 'A' && c <= 'F':
			out[i] = c - 'A' + 10
		case c >= 'a' && c <= 'f':
			out[i] = c - 'a' + 10
		default:
			wipe(out)
			return nil, fmt.Errorf("pin block must contain only hex characters")
		}
	}
	return out, nil
}

// nibblesHex returns nibbles as upper-case hex text. The text is built in a
// buffer that is cleared before returning; the returned string cannot be.
func nibblesHex(nibbles []byte) string {
	buf := make([]byte, len(nibbles))
	defer wipe(buf)
	for i, n := range nibbles {
		buf[i] = upperHex[n&0x0F]
	}
	return string(buf)
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
	return randomFillFrom(rand.Reader, dst, table)
}

// randomFillFrom is randomFill reading from r, so tests can supply the source.
func randomFillFrom(r io.Reader, dst, table []byte) error {
	limit := 256 - 256%len(table)
	buf := make([]byte, len(dst)+16)
	defer wipe(buf)

	filled := 0
	for filled < len(dst) {
		if _, err := io.ReadFull(r, buf); err != nil {
			return fmt.Errorf("reading random fill: %w", err)
		}
		for _, b := range buf {
			if int(b) >= limit {
				continue
			}
			dst[filled] = table[int(b)%len(table)]
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
// digit (ISO 9564-1). account must hold at least 13 digits, and the digits used
// must be decimal. The result belongs to the caller, who must clear it.
func panFieldISO0(account string) ([]byte, error) {
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

// debugOutput holds the writer set by SetDebugWriter. Every format embeds it.
type debugOutput struct {
	debugWriter io.Writer
}

// SetDebugWriter sets a writer that receives a short note after each
// successful Encode or Decode, naming the operation and the format. It never
// writes the PIN, the PAN, the PIN length, the fill or any PIN block, clear or
// enciphered.
func (d *debugOutput) SetDebugWriter(writer io.Writer) {
	d.debugWriter = writer
}

// note writes the summary SetDebugWriter promises, and nothing else.
func (d *debugOutput) note(operation, format string) {
	if d.debugWriter == nil {
		return
	}
	fmt.Fprintf(d.debugWriter, "PIN block %s operation finished\nFormat  : %s\n\n", operation, format)
}
