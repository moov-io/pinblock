package formats

import "io"

const (
	iso0Version  = "0"
	iso1Version  = "1"
	iso2Version  = "2"
	iso3Version  = "3"
	eci2Version  = "eci-2"
	eci3Version  = "eci-3"
	visa2Version = "visa-2"
)

var (
	hexLetters    = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'A', 'B', 'C', 'D', 'E', 'F'}
	hexDigits     = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	hexCharacters = []byte{'A', 'B', 'C', 'D', 'E', 'F'}
)

type FormatA interface {
	SetDebugWriter(writer io.Writer)
	Encode(pin, account string) (string, error)
	Decode(pinBlock, account string) (string, error)
}

type FormatB interface {
	SetDebugWriter(writer io.Writer)
	Encode(pin string) (string, error)
	Decode(pinBlock string) (string, error)
}

type Cipher interface {
	Encrypt(plainText []byte) ([]byte, error)
	Decrypt(cipherText []byte) ([]byte, error)
}

func NewISO0() FormatA {
	return &iso0Object{
		Filler: "F", // default to ISO0's Filler

		format:  "Format 0 (ISO-0)",
		version: iso0Version,
	}
}

func NewISO1() FormatB {
	return &iso1Object{
		version: iso1Version,
		format:  "Format 1 (ISO-1)",
	}
}

// The ISO-2 PIN Block format is used for smart card offline authentication.
// It is similar to an ISO-1 PIN Block in that there is no PAN to associate with the PIN.
// It differs in that the fill is 0xF instead of random digits
func NewISO2() FormatB {
	return &iso1Object{
		Filler: "F", // default to ISO1's Filler

		format:  "Format 2 (ISO-2)",
		version: iso2Version,
	}
}

// ISO 9564-1: 2002 Format 3.
//
//	Format 3 is the same as format 0, except that the “fill” digits are random values from 10 to 15,
//	and the first nibble (which identifies the block format) has the value 3.
func NewISO3() FormatA {
	return &iso0Object{
		format:  "Format 3 (ISO-3)",
		version: iso3Version,
	}
}

func NewISO4(cipher Cipher) FormatA {
	return &iso4Object{
		Filler: "A", // default to ISO-4

		cipher: cipher,
		format: "Format 4 (ISO-4)",
	}
}

// ANSI X9.8:
//
//	Same as ISO-0.
func NewANSIX98() FormatA {
	return &iso0Object{
		Filler: "F", // default to ISO0's Filler

		format:  "ANSI X9.8",
		version: iso0Version,
	}
}

// OEM-1 / Diebold / Docutel / NCR
//
//	The OEM-1 PIN block format is equivalent to the PIN block formats that Diebold, Docutel, and NCR define.
//	The OEM-1 PIN block format supports a PIN from 4 to 12 digits in length.
//	A PIN that is longer than 12 digits is truncated on the right.
func NewOEM1() FormatB {
	return &oemObject{
		format: "Diebold, Docutel, NCR",
	}
}

// ECI-1
//
//	Same as ISO-0.
func NewECI1() FormatA {
	return &iso0Object{
		Filler: "F", // default to ISO0's Filler

		version: iso0Version,
		format:  "ECI-1",
	}
}

// ECI-2
func NewECI2() FormatB {
	return &eciObject{
		format:  "ECI-2",
		version: eci2Version,
	}
}

// ECI-3
func NewECI3() FormatB {
	return &eciObject{
		format:  "ECI-3",
		version: eci3Version,
	}
}

// ECI-4
//
//	Same as ISO-1.
func NewECI4() FormatB {
	return &iso1Object{
		format:  "ECI-4",
		version: iso1Version,
	}
}

// VISA-1
//
//	Same as ISO-0.
func NewVISA1() FormatA {
	return &iso0Object{
		Filler: "F", // default to ISO0's Filler

		version: iso0Version,
		format:  "VISA-1",
	}
}

// VISA-2
func NewVISA2() FormatB {
	return &eciObject{
		version: visa2Version,
		format:  "VISA-2",
	}
}

// VISA-3
func NewVISA3() FormatB {
	return &visa3Object{
		format: "VISA-3",
	}
}

// VISA-4
func NewVISA4() FormatA {
	return &iso0Object{
		Filler: "F", // default to ISO0's Filler

		version: iso0Version,
		format:  "VISA-4",
	}
}
