// Package memzero clears buffers that held a PIN, a PAN or a PIN block.
//
// A clear PIN block, and anything a PIN can be recovered from, is a secret.
// PCI PIN Security Requirements expect it to be cleared from memory as soon as
// it is no longer needed, rather than left on the heap until the garbage
// collector reuses the memory.
//
// What this does and does not guarantee: it overwrites the buffer it is given.
// It cannot reach copies made elsewhere, such as a cipher's internal state or
// a string, because Go offers no way to clear those.
package memzero

import "runtime"

// Bytes overwrites every byte of b with zero.
//
// It is not inlined, and keeps b reachable until the write is done, so the
// compiler cannot treat the clear as a dead store to a buffer that is about to
// go out of scope.
//
//go:noinline
func Bytes(b []byte) {
	clear(b)
	runtime.KeepAlive(b)
}
