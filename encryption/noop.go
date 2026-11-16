package encryption

// NoOp is a Cipher that does not encrypt. It returns a copy of its input, so
// the ISO-4 structure can be inspected in tests. Do not use it for real PIN
// blocks.
type NoOp struct{}

// NewNoOp returns a NoOp cipher.
func NewNoOp() *NoOp {
	return &NoOp{}
}

// Encrypt returns a copy of plaintext. It never returns plaintext itself, so
// a caller that clears its input does not also clear the result.
func (n *NoOp) Encrypt(plaintext []byte) ([]byte, error) {
	return append([]byte(nil), plaintext...), nil
}

// Decrypt returns a copy of ciphertext.
func (n *NoOp) Decrypt(ciphertext []byte) ([]byte, error) {
	return append([]byte(nil), ciphertext...), nil
}
