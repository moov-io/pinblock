package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// AesECB is an AES block cipher for ISO-4 PIN blocks. Encrypt and Decrypt
// take exactly one 16-byte block, do not modify it, and return a new slice.
type AesECB struct {
	cipherBlock cipher.Block
}

// NewAesECB returns an AesECB for key (16, 24 or 32 bytes). crypto/aes copies
// the key into its own key schedule, so the caller may clear key afterwards;
// that schedule is held by crypto/aes and cannot be cleared.
func NewAesECB(key []byte) (*AesECB, error) {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("creating cipher: %w", err)
	}

	return &AesECB{
		cipherBlock: cipher,
	}, nil
}

func (a *AesECB) Encrypt(plainText []byte) ([]byte, error) {
	if len(plainText) != 16 {
		return nil, fmt.Errorf("plain text length must be 16 bytes")
	}

	cipherText := make([]byte, len(plainText))

	a.cipherBlock.Encrypt(cipherText, plainText)

	return cipherText, nil
}

func (a *AesECB) Decrypt(cipherText []byte) ([]byte, error) {
	if len(cipherText) != 16 {
		return nil, fmt.Errorf("cipher text length must be 16 bytes")
	}

	plainText := make([]byte, len(cipherText))

	a.cipherBlock.Decrypt(plainText, cipherText)

	return plainText, nil
}
