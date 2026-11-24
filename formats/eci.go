package formats

import "fmt"

// eciObject implements ECI-2, ECI-3 and VISA-2. None uses the PAN.
//
//	ECI-2:  P1 P2 P3 P4 R .. R    the first 4 PIN digits, then 12 random nibbles
//	ECI-3:  N P ... P 0 ... 0 R   PIN length (4 to 6), the PIN truncated to 6
//	                              digits, zeros to nibble 7, then 9 random nibbles
//	VISA-2: as ECI-3, with random decimal digits as the fill
type eciObject struct {
	debugOutput

	version int
	format  string
}

// Encode returns the 16-character PIN block for pin. account is not used. A
// PIN longer than the format holds (4 digits for ECI-2, 6 for ECI-3 and
// VISA-2) is truncated on the right, as the format defines. See Format.
func (i *eciObject) Encode(pin, _ string) (string, error) {
	if len(pin) < 4 || len(pin) > 12 {
		return "", fmt.Errorf("pin length must be between 4 and 12 digits")
	}
	if err := checkDigits(pin, "pin"); err != nil {
		return "", err
	}

	nibbles := make([]byte, fieldNibbles)
	defer wipe(nibbles)

	var used int
	switch i.version {
	case eci2Version:
		used = 4
		_ = digitNibbles(nibbles, pin[:4], "pin")
	default:
		n := min(len(pin), 6)
		nibbles[0] = byte(n & 0x0F) // 4..6
		_ = digitNibbles(nibbles[1:], pin[:n], "pin")
		used = 7 // length, up to 6 digits, zero-padded
	}

	table := anyNibble
	if i.version == visa2Version {
		table = decimalNibble
	}
	if err := randomFill(nibbles[used:], table); err != nil {
		return "", err
	}
	pinBlock := nibblesHex(nibbles)

	i.note("encode", i.format)
	return pinBlock, nil
}

// Decode returns the PIN from a 16-character PIN block. account is not used.
// See Format.
func (i *eciObject) Decode(pinBlock, _ string) (string, error) {
	if len(pinBlock) != fieldNibbles {
		return "", errBlockLength16
	}

	nibbles, err := hexNibbles(pinBlock)
	if err != nil {
		return "", err
	}
	defer wipe(nibbles)

	var digits []byte
	if i.version == eci2Version {
		digits = nibbles[:4]
	} else {
		pinLength := int(nibbles[0])
		if pinLength < 4 || pinLength > 6 {
			return "", errBlockUnparsable
		}
		digits = nibbles[1 : 1+pinLength]
	}
	pin, err := pinFromNibbles(digits)
	if err != nil {
		return "", err
	}

	i.note("decode", i.format)
	return pin, nil
}
