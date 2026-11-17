package formats

import "fmt"

// iso1Object implements ISO 9564-1 Formats 1 and 2 (and ECI-4, which is
// Format 1). Neither uses the PAN.
//
// PIN block (16 nibbles): C N P1 .. P4 P5/F .. P12/F F13 F14
//
//	C = control field, 1 for Format 1 and 2 for Format 2
//	N = PIN length, 4 to 12
//	P = PIN digit
//	F = fill: random 0x0 to 0xF for Format 1, 0xF for Format 2
type iso1Object struct {
	debugOutput

	version byte
	format  string
}

// fill writes the fill nibbles for this format into dst.
func (i *iso1Object) fill(dst []byte) error {
	if i.version == 0x2 {
		for n := range dst {
			dst[n] = 0xF
		}
		return nil
	}
	return randomFill(dst, anyNibble)
}

// Encode returns the 16-character PIN block for pin. account is not used. See
// Format.
//
//	The `ISO-1` PIN block format is equivalent to an `ECI-4` PIN block format
//	and is recommended for usage where no PAN data is available.
func (i *iso1Object) Encode(pin, _ string) (string, error) {
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
	pinBlock := nibblesHex(pinField)

	i.note("encode", i.format)
	return pinBlock, nil
}

// Decode returns the PIN from a 16-character PIN block. account is not used.
// See Format.
func (i *iso1Object) Decode(pinBlock, _ string) (string, error) {
	if len(pinBlock) != fieldNibbles {
		return "", fmt.Errorf("pin block must be 16 characters")
	}

	pinField, err := hexNibbles(pinBlock)
	if err != nil {
		return "", err
	}
	defer wipe(pinField)

	if pinField[0] != i.version {
		return "", fmt.Errorf("format is different")
	}
	pinLength := int(pinField[1])
	if pinLength < 4 || pinLength > 12 {
		return "", fmt.Errorf("invalid pin length %d", pinLength)
	}
	if i.version == 0x2 {
		for _, n := range pinField[2+pinLength:] {
			if n != 0xF {
				return "", fmt.Errorf("invalid fill in pin block")
			}
		}
	}
	pin, err := pinFromNibbles(pinField[2 : 2+pinLength])
	if err != nil {
		return "", err
	}

	i.note("decode", i.format)
	return pin, nil
}
