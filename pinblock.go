package pinblock

import (
	"errors"

	"github.com/moov-io/pinblock/formats"
)

// PINBlock is an interface to represent a PIN block format.
type PINBlock interface {
	Create(pin string, account string) (string, error)
}

// GetPINBlockFormat returns a PINBlock of the specified format.
func GetPINBlockFormat(format string) (PINBlock, error) {
	switch format {
	case "ISO0":
		return formats.NewISO0(), nil
	// More cases here for other formats...
	default:
		return nil, errors.New("unsupported format")
	}
}

// CreatePINBlock creates a PIN block in the specified format.
func CreatePINBlock(pin string, account string, format string) (string, error) {
	pinBlockFormat, err := GetPINBlockFormat(format)
	if err != nil {
		return "", err
	}
	return pinBlockFormat.Create(pin, account)
}
