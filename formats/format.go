package formats

import (
	"fmt"
	"io"

	"github.com/moov-io/pinblock/encryption"
)

// Version identifiers for the ECI-style formats.
const (
	eci2Version = iota + 1
	eci3Version
	visa2Version
)

// Format encodes a PIN into a PIN block and decodes it back.
//
//   - pin is the PIN as ASCII decimal digits, for example "1234". Any other
//     character is an error.
//   - account is the PAN as ASCII decimal digits. Formats that do not use the
//     PAN (ISO-1, ISO-2, OEM-1, ECI-2, ECI-3, ECI-4, VISA-2, VISA-3) ignore it.
//   - A PIN block is upper-case hex text: 16 characters, or 32 for ISO-4.
//     Decode also accepts lower case.
//
// Every intermediate the methods build in memory (the PIN field, the PAN
// field, XOR results, random fill, cipher outputs and the text being
// converted) is held in a byte slice and cleared before the method returns,
// including when it returns an error.
//
// The arguments and results are Go strings, and a string cannot be
// overwritten: the PIN, the PAN and the PIN block passed in or returned stay in
// memory until the garbage collector reuses it, and neither this package nor
// the caller can clear them. Treat the PIN block as sensitive too: for ISO-0
// and ISO-3 it is the PIN field XORed with the PAN, so anyone holding the PAN
// can recover the PIN. An API that takes and returns byte slices the caller
// can clear is proposed as a separate enhancement.
type Format interface {
	// SetDebugWriter sets a writer that receives a short note after each
	// successful Encode or Decode, naming the operation and the format. It
	// never writes the PIN, the PAN, the PIN length, the fill or any PIN
	// block, clear or enciphered.
	SetDebugWriter(writer io.Writer)
	// Encode returns the PIN block for pin and account.
	Encode(pin, account string) (string, error)
	// Decode returns the PIN held in pinBlock. It checks the block's
	// structure (format nibble, PIN length, decimal digits and, where the
	// format defines it, the fill) and returns an error rather than a PIN if
	// it is malformed, which is also what a block decrypted under the wrong
	// key looks like.
	Decode(pinBlock, account string) (string, error)
}

// Cipher is the block cipher ISO-4 uses, with a 16-byte block.
//
// Encrypt and Decrypt must not modify their input and must return a new
// 16-byte slice. The ISO-4 format copies the result and clears the slice it
// was given, so an implementation that returns a buffer it keeps will find it
// cleared. A result of any other length is an error.
type Cipher interface {
	Encrypt(plainText []byte) ([]byte, error)
	Decrypt(cipherText []byte) ([]byte, error)
}

type constructorFunc func() Format

var (
	formatConstructor = map[string]constructorFunc{
		"ISO-0": func() Format { return NewISO0() },
		"ISO-1": func() Format { return NewISO1() },
		"ISO-2": func() Format { return NewISO2() },
		"ISO-3": func() Format { return NewISO3() },
		"ISO-4": func() Format { return NewISO4(encryption.NewNoOp()) },
		"ANSI":  func() Format { return NewANSIX98() },
		"ECI1":  func() Format { return NewECI1() },
		"ECI2":  func() Format { return NewECI2() },
		"ECI3":  func() Format { return NewECI3() },
		"ECI4":  func() Format { return NewECI4() },
		"VISA1": func() Format { return NewVISA1() },
		"VISA2": func() Format { return NewVISA2() },
		"VISA3": func() Format { return NewVISA3() },
		"VISA4": func() Format { return NewVISA4() },
	}
)

// NewFormatter returns the Format named bType: one of ISO-0, ISO-1, ISO-2,
// ISO-3, ISO-4, ANSI, ECI1, ECI2, ECI3, ECI4, VISA1, VISA2, VISA3 or VISA4.
func NewFormatter(bType string) (Format, error) {
	constructor := formatConstructor[bType]
	if constructor == nil {
		return nil, fmt.Errorf("unsupported pinblock type")
	}

	return constructor(), nil
}

// NewISO0 returns ISO 9564-1 Format 0: the PIN field (0, PIN length, PIN,
// fill F) XORed with the PAN field (0000 and the 12 rightmost PAN digits
// excluding the check digit).
func NewISO0() Format {
	return &iso0Object{format: "Format 0 (ISO-0)", version: 0x0}
}

// NewISO1 returns ISO 9564-1 Format 1: 1, PIN length, PIN and random fill
// 0 to F, with no PAN. It is recommended where no PAN is available.
func NewISO1() Format {
	return &iso1Object{format: "Format 1 (ISO-1)", version: 0x1}
}

// The ISO-2 PIN Block format is used for smart card offline authentication.
// It is similar to an ISO-1 PIN Block in that there is no PAN to associate with the PIN.
// It differs in that the fill is 0xF instead of random values.
func NewISO2() Format {
	return &iso1Object{format: "Format 2 (ISO-2)", version: 0x2}
}

// ISO 9564-1: 2002 Format 3.
//
//	Format 3 is the same as format 0, except that the “fill” digits are random values from 10 to 15,
//	and the first nibble (which identifies the block format) has the value 3.
func NewISO3() Format {
	return &iso0Object{format: "Format 3 (ISO-3)", version: 0x3}
}

// NewISO4 returns ISO 9564-1 Format 4, for AES. The PIN field is encrypted,
// XORed with the PAN field and encrypted again under cipher.
// NewFormatter("ISO-4") uses a no-op cipher, which yields the unencrypted
// structure and is meant for tests.
func NewISO4(cipher Cipher) Format {
	return &iso4Object{cipher: cipher, format: "Format 4 (ISO-4)"}
}

// ANSI X9.8:
//
//	Same as ISO-0.
func NewANSIX98() Format {
	return &iso0Object{format: "ANSI X9.8", version: 0x0}
}

// OEM-1 / Diebold / Docutel / NCR
//
//	The OEM-1 PIN block format is equivalent to the PIN block formats that Diebold, Docutel, and NCR define.
//	The OEM-1 PIN block format supports a PIN from 4 to 12 digits in length.
//	A PIN that is longer than 12 digits is truncated on the right.
func NewOEM1() Format {
	return &oemObject{
		format: "Diebold, Docutel, NCR",
	}
}

// ECI-1
//
//	Same as ISO-0.
func NewECI1() Format {
	return &iso0Object{format: "ECI-1", version: 0x0}
}

// ECI-2
func NewECI2() Format {
	return &eciObject{format: "ECI-2", version: eci2Version}
}

// ECI-3
func NewECI3() Format {
	return &eciObject{format: "ECI-3", version: eci3Version}
}

// ECI-4
//
//	Same as ISO-1.
func NewECI4() Format {
	return &iso1Object{format: "ECI-4", version: 0x1}
}

// VISA-1
//
//	Same as ISO-0.
func NewVISA1() Format {
	return &iso0Object{format: "VISA-1", version: 0x0}
}

// VISA-2
func NewVISA2() Format {
	return &eciObject{format: "VISA-2", version: visa2Version}
}

// VISA-3
func NewVISA3() Format {
	return &visa3Object{
		format: "VISA-3",
	}
}

// VISA-4
func NewVISA4() Format {
	return &iso0Object{format: "VISA-4", version: 0x0}
}
