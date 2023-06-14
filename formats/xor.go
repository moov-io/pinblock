package formats

import (
	"encoding/hex"
	"fmt"
)

func xorHex(a, b string) (string, error) {
	bytesA, err := hex.DecodeString(a)
	if err != nil {
		return "", err
	}
	bytesB, err := hex.DecodeString(b)
	if err != nil {
		return "", err
	}

	if len(bytesA) != len(bytesB) {
		return "", fmt.Errorf("length mismatch: %d vs %d", len(bytesA), len(bytesB))
	}

	// XOR the bytes
	xorBytes := make([]byte, len(bytesA))
	for i := 0; i < len(bytesA); i++ {
		xorBytes[i] = bytesA[i] ^ bytesB[i]
	}

	return hex.EncodeToString(xorBytes), nil
}

func xor(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length mismatch: %d vs %d", len(a), len(b))
	}

	// XOR the bytes
	xorBytes := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		xorBytes[i] = a[i] ^ b[i]
	}

	return xorBytes, nil
}
