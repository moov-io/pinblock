package formats

import (
	"fmt"
	"strings"
)

type ISO0 struct {
	Filler string
}

func NewISO0() *ISO0 {
	return &ISO0{
		Filler: "F", // default to ISO0's Filler
	}
}

// Encode returns the ISO0 PIN block for the given PIN and account number
func (i *ISO0) Encode(pin, account string) (string, error) {
	if len(pin) < 4 || len(pin) > 12 {
		return "", fmt.Errorf("pin length must be between 4 and 12 digits")
	}

	// pin block should start with 0, then add length of pin, then add pin,
	// then add F until 16 characters
	padding := strings.Repeat(i.Filler, 14-len(pin))
	pinBlock := fmt.Sprintf("0%d%s%s", len(pin), pin, padding)

	// account number must be at least 13 digits, including the check digit
	if len(account) < 13 {
		return "", fmt.Errorf("account length must be at least 13 digits")
	}

	// take the last 12 digits of the account number excluding the check digit
	accountBlock := fmt.Sprintf("0000%s", account[len(account)-13:len(account)-1])

	xorBlock, err := xorHex(pinBlock, accountBlock)
	if err != nil {
		return "", err
	}

	return strings.ToUpper(xorBlock), nil
}

func (i *ISO0) Decode(pinBlock, account string) (string, error) {
	if len(pinBlock) != 16 {
		return "", fmt.Errorf("pin block must be 16 characters")
	}

	if len(account) < 13 {
		return "", fmt.Errorf("account length must be at least 13 digits")
	}

	// take the last 12 digits of the account number excluding the check digit
	accountBlock := fmt.Sprintf("0000%s", account[len(account)-13:len(account)-1])

	decodedBlock, err := xorHex(pinBlock, accountBlock)
	if err != nil {
		return "", err
	}

	// decodedBlock should start with 0, then has length of pin, then has pin, then has F until 16 characters
	pinLength := int(decodedBlock[1] - '0')
	pin := decodedBlock[2 : 2+pinLength]

	return pin, nil
}
