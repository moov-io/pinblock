package formats

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

type ISO4 struct {
	Filler string
	cipher Cipher
}

type Cipher interface {
	Encrypt(plainText []byte) ([]byte, error)
	Decrypt(cipherText []byte) ([]byte, error)
}

func NewISO4(cipher Cipher) *ISO4 {
	return &ISO4{
		Filler: "A", // default to ISO-4
		cipher: cipher,
	}
}

// Format returns iso type
func (i *ISO4) Format() string {
	return "Format 4 (ISO-4)"
}

// Padding returns padding pattern
func (i *ISO4) Padding(pin string) string {
	return strings.Repeat(i.Filler, 16-len(pin)-2)
}

// Encode returns an ISO-4 formatted and encrypted PIN block
func (i *ISO4) Encode(pin, account string) (string, error) {
	// both pinBlock and panBlock are 16 bytes (128 bits)
	if len(pin) < 4 || len(pin) > 12 {
		return "", fmt.Errorf("pin length must be between 4 and 12 digits")
	}

	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("generating random bytes: %w", err)
	}

	pinBlock := fmt.Sprintf("4%X%s%s%X", len(pin), pin, i.Padding(pin), randomBytes)

	rawPinBlock, err := hex.DecodeString(pinBlock)
	if err != nil {
		return "", fmt.Errorf("decoding pinBlock: %w", err)
	}

	blockA, err := i.cipher.Encrypt(rawPinBlock)
	if err != nil {
		return "", fmt.Errorf("encrypting pinBlock: %w", err)
	}

	if len(account) < 12 || len(account) > 19 {
		return "", fmt.Errorf("account length must be between 12 and 19 digits")
	}

	// 4-bit field with permissible values 0000 (zero) to 0111 (7) indicate
	// a PAN length of 12 plus the value of the field (ranging then from 12
	// to 19). If the PAN is less than 12 digits, the digits are right
	// justified and padded to the left with zeros, and M is set to 0;
	controlField := fmt.Sprintf("%X", len(account)-12)

	panBlock := fmt.Sprintf("%s%s%s", controlField, account, strings.Repeat("0", 32-len(account)-1))

	rawPanBlock, err := hex.DecodeString(panBlock)
	if err != nil {
		return "", fmt.Errorf("decoding panBlock: %w", err)
	}

	blockB, err := xor(rawPanBlock, blockA)
	if err != nil {
		return "", fmt.Errorf("xor-ing block A and pan block: %w", err)
	}

	encryptedPinBlock, err := i.cipher.Encrypt(blockB)
	if err != nil {
		return "", fmt.Errorf("encrypting block B: %w", err)
	}

	return fmt.Sprintf("%X", encryptedPinBlock), nil
}

// Decode returns the PIN from an ISO-4 encrypted PIN block
func (i *ISO4) Decode(pinBlock, account string) (string, error) {
	if len(pinBlock) != 32 {
		return "", fmt.Errorf("pinBlock must be 32 hex characters (16 bytes)")
	}

	encryptedPinBlock, err := hex.DecodeString(pinBlock)
	if err != nil {
		return "", fmt.Errorf("decoding pinBlock: %w", err)
	}

	blockB, err := i.cipher.Decrypt(encryptedPinBlock)
	if err != nil {
		return "", fmt.Errorf("decrypting pinBlock: %w", err)
	}

	if len(account) < 12 || len(account) > 19 {
		return "", fmt.Errorf("account length must be between 12 and 19 digits")
	}

	controlField := fmt.Sprintf("%X", len(account)-12)

	panBlock := fmt.Sprintf("%s%s%s", controlField, account, strings.Repeat("0", 32-len(account)-1))

	rawPanBlock, err := hex.DecodeString(panBlock)
	if err != nil {
		return "", fmt.Errorf("decoding panBlock: %w", err)
	}

	blockA, err := xor(rawPanBlock, blockB)
	if err != nil {
		return "", fmt.Errorf("xor-ing block B and pan block: %w", err)
	}

	rawPinBlock, err := i.cipher.Decrypt(blockA)
	if err != nil {
		return "", fmt.Errorf("decrypting block A: %w", err)
	}

	plainPinBlock := hex.EncodeToString(rawPinBlock)

	// rawPinBlock should now be the original pinBlock, we'll parse it to get the PIN.
	pinLength, err := strconv.ParseInt(string(plainPinBlock[1]), 16, 64)
	if err != nil {
		return "", fmt.Errorf("parsing pin length: %w", err)
	}

	pin := string(plainPinBlock[2 : 2+pinLength])

	return pin, nil
}
