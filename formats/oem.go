package formats

import "fmt"

// oemPadOrder is the order in which OEM-1 picks its pad value: the first of
// these that is not a PIN digit.
var oemPadOrder = []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0x0, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF}

// oemObject implements OEM-1 (Diebold, Docutel, NCR). It does not use the PAN.
//
// PIN block (16 nibbles): the PIN (4 to 12 digits), then a pad value repeated
// to the end. The pad value is a nibble that differs from every PIN digit.
type oemObject struct {
	debugOutput

	format string
}

// Encode returns the 16-character PIN block for pin. account is not used. A
// PIN longer than 12 digits is truncated on the right, as the format defines.
// See Format.
func (i *oemObject) Encode(pin, _ string) (string, error) {
	if len(pin) < 4 {
		return "", fmt.Errorf("pin length must be between 4 and 12 digits")
	}
	if err := checkDigits(pin, "pin"); err != nil {
		return "", err
	}
	pin = pin[:min(len(pin), 12)]

	nibbles := make([]byte, fieldNibbles)
	defer wipe(nibbles)
	_ = digitNibbles(nibbles, pin, "pin")

	var used [16]bool
	for _, n := range nibbles[:len(pin)] {
		used[n] = true
	}
	var pad byte
	for _, n := range oemPadOrder {
		if !used[n] {
			pad = n
			break
		}
	}
	for n := len(pin); n < fieldNibbles; n++ {
		nibbles[n] = pad
	}
	pinBlock := nibblesHex(nibbles)

	i.note("encode", i.format)
	return pinBlock, nil
}

// Decode returns the PIN from a 16-character PIN block. account is not used.
// See Format.
func (i *oemObject) Decode(pinBlock, _ string) (string, error) {
	if len(pinBlock) != fieldNibbles {
		return "", fmt.Errorf("pin block must be 16 characters")
	}

	nibbles, err := hexNibbles(pinBlock)
	if err != nil {
		return "", err
	}
	defer wipe(nibbles)

	// The PIN ends where the run of pad values at the end of the block starts.
	pad := nibbles[fieldNibbles-1]
	pinLength := fieldNibbles
	for pinLength > 0 && nibbles[pinLength-1] == pad {
		pinLength--
	}
	if pinLength < 4 || pinLength > 12 {
		return "", fmt.Errorf("invalid pin length %d", pinLength)
	}
	pin, err := pinFromNibbles(nibbles[:pinLength])
	if err != nil {
		return "", err
	}

	i.note("decode", i.format)
	return pin, nil
}
