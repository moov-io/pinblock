package formats

import "fmt"

// iso0Object implements ISO 9564-1 Formats 0 and 3 (and ANSI X9.8, ECI-1,
// VISA-1 and VISA-4, which are Format 0).
//
// PIN field (16 nibbles): C N P1 .. P4 P5/F .. P12/F F13 F14
//
//	C = control field, 0 for Format 0 and 3 for Format 3
//	N = PIN length, 4 to 12
//	P = PIN digit
//	F = fill: 0xF for Format 0, random 0xA to 0xF for Format 3
//
// PAN field (16 nibbles): 0 0 0 0 A1 ... A12, the 12 rightmost PAN digits
// excluding the check digit. The PIN block is the two fields XORed.
type iso0Object struct {
	debugOutput

	version byte
	format  string
}

// fill writes the fill nibbles for this format into dst.
func (i *iso0Object) fill(dst []byte) error {
	if i.version == 0x3 {
		return randomFill(dst, letterNibble)
	}
	for n := range dst {
		dst[n] = 0xF
	}
	return nil
}

// validFill reports whether n is a fill nibble this format allows.
func (i *iso0Object) validFill(n byte) bool {
	if i.version == 0x3 {
		return n >= 0xA
	}
	return n == 0xF
}

// Encode returns the 16-character PIN block for pin and account. See Format.
func (i *iso0Object) Encode(pin, account string) (string, error) {
	if len(pin) < 4 || len(pin) > 12 {
		return "", fmt.Errorf("pin length must be between 4 and 12 digits")
	}

	pinField := make([]byte, fieldNibbles)
	defer wipe(pinField)
	pinField[0] = i.version
	pinField[1] = byte(len(pin) & 0x0F) // 4..12, checked above
	if err := digitNibbles(pinField[2:], pin, "pin"); err != nil {
		return "", err
	}
	if err := i.fill(pinField[2+len(pin):]); err != nil {
		return "", err
	}

	panField, err := panFieldISO0(account)
	if err != nil {
		return "", err
	}
	defer wipe(panField)

	xorInto(pinField, pinField, panField)
	pinBlock := nibblesHex(pinField)

	i.note("encode", i.format)
	return pinBlock, nil
}

// Decode returns the PIN from a 16-character PIN block. See Format.
func (i *iso0Object) Decode(pinBlock, account string) (string, error) {
	if len(pinBlock) != fieldNibbles {
		return "", fmt.Errorf("pin block must be 16 characters")
	}

	panField, err := panFieldISO0(account)
	if err != nil {
		return "", err
	}
	defer wipe(panField)

	pinField, err := hexNibbles(pinBlock)
	if err != nil {
		return "", err
	}
	defer wipe(pinField)
	xorInto(pinField, pinField, panField)

	if pinField[0] != i.version {
		return "", fmt.Errorf("format is different")
	}
	pinLength := int(pinField[1])
	if pinLength < 4 || pinLength > 12 {
		return "", fmt.Errorf("invalid pin length %d", pinLength)
	}
	for _, n := range pinField[2+pinLength:] {
		if !i.validFill(n) {
			return "", fmt.Errorf("invalid fill in pin block")
		}
	}
	pin, err := pinFromNibbles(pinField[2 : 2+pinLength])
	if err != nil {
		return "", err
	}

	i.note("decode", i.format)
	return pin, nil
}
