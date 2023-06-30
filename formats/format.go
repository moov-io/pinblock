package formats

import "io"

const (
	ISO0Version = "0"
	ISO1Version = "1"
	ISO2Version = "2"
	ISO3Version = "3"
)

type ISO0 interface {
	SetDebugWriter(writer io.Writer)
	Encode(pin, account string) (string, error)
	Decode(pinBlock, account string) (string, error)
}

type ISO1 interface {
	SetDebugWriter(writer io.Writer)
	Encode(pin string) (string, error)
	Decode(pinBlock string) (string, error)
}

type Cipher interface {
	Encrypt(plainText []byte) ([]byte, error)
	Decrypt(cipherText []byte) ([]byte, error)
}

func NewISO0() ISO0 {
	return &iso0Object{
		Filler:  "F", // default to ISO0's Filler
		format:  "Format 0 (ISO-0)",
		version: ISO0Version,
	}
}

func NewISO1() ISO1 {
	return &iso1Object{
		version: ISO1Version,
		format:  "Format 1 (ISO-1)",
	}
}

// The ISO-2 PIN Block format is used for smart card offline authentication.
// It is similar to an ISO-1 PIN Block in that there is no PAN to associate with the PIN.
// It differs in that the fill is 0xF instead of random digits
func NewISO2() ISO1 {
	return &iso1Object{
		Filler:  "F", // default to ISO1's Filler
		format:  "Format 2 (ISO-2)",
		version: ISO2Version,
	}
}

// ISO 9564-1: 2002 Format 3.
// Format 3 is the same as format 0, except that the “fill” digits are random values from 10 to 15,
// and the first nibble (which identifies the block format) has the value 3.
func NewISO3() ISO0 {
	return &iso0Object{
		format:  "Format 3 (ISO-3)",
		version: ISO3Version,
	}
}

func NewISO4(cipher Cipher) ISO0 {
	return &iso4Object{
		Filler: "A", // default to ISO-4
		format: "Format 4 (ISO-4)",
		cipher: cipher,
	}
}
