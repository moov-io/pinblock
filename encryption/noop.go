package encryption

type NoOp struct{}

func NewNoOp() *NoOp {
	return &NoOp{}
}

func (n *NoOp) Encrypt(plaintext []byte) ([]byte, error) {
	return plaintext, nil
}

func (n *NoOp) Decrypt(ciphertext []byte) ([]byte, error) {
	return ciphertext, nil
}
