package formats

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type ISO0 struct {
	Filler string
}

func NewISO0() *ISO0 {
	return &ISO0{
		Filler: "F",
	}
}

// Encode returns the ISO0 PIN block for the given PIN and account number
func (i *ISO0) Encode(pin string, account string) (string, error) {
	// pin block should start with 0, then add length of pin, then add pin,
	// then add F until 16 characters
	padding := strings.Repeat(i.Filler, 14-len(pin))
	pinBlock := fmt.Sprintf("0%d%s%s", len(pin), pin, padding)

	// take the last 12 digits of the account number excluding the check digit
	accountBlock := fmt.Sprintf("0000%s", account[len(account)-13:len(account)-1])

	xorBlock, err := xorHex(pinBlock, accountBlock)
	if err != nil {
		return "", err
	}

	return strings.ToUpper(xorBlock), nil
}

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
