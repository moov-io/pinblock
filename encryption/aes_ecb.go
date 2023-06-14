package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

type AesECB struct {
	cipherBlock cipher.Block
}

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
