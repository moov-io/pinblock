// Package pinblock formats PINs into PIN blocks and back. The formats live in
// the formats subpackage; this package holds helpers shared by its callers.
package pinblock

import "github.com/moov-io/pinblock/internal/memzero"

// Clear overwrites every byte of b with zero, for a caller clearing a PIN, a
// PAN or a PIN block once it is no longer needed.
//
// Prefer it to the clear builtin. A clear(b) with no later read of b is a dead
// store the compiler may remove; Clear is not inlined and keeps b reachable
// until the write is done, so the write happens.
//
// Clear zeroes len(b) bytes. To clear a slice's full backing array, pass
// b[:cap(b)]. It is best-effort heap hygiene: it does not reach copies in CPU
// registers or on the stack, the AES key schedule inside crypto/aes, or
// strings.
func Clear(b []byte) {
	memzero.Bytes(b)
}
