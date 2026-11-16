package formats

import "fmt"

const (
	eci2Version = iota + 1
	eci3Version
	visa2Version
)

// eciObject implements ECI-2, ECI-3 and VISA-2. None uses the PAN.
//
//	ECI-2:  P1 P2 P3 P4 R .. R    the first 4 PIN digits, then 12 random nibbles
//	ECI-3:  N P ... P 0 ... 0 R   PIN length (4 to 6), the PIN truncated to 6
//	                              digits, zeros to nibble 7, then 9 random nibbles
//	VISA-2: as ECI-3, with random decimal digits as the fill
type eciObject struct {
	version int
	format  string
}

// Encode returns the 8-byte PIN block for pin. account is not used. A PIN
// longer than the format holds (4 digits for ECI-2, 6 for ECI-3 and VISA-2) is
// truncated on the right, as the format defines. See Format.
func (i *eciObject) Encode(pin, _ []byte) ([]byte, error) {
	if len(pin) < 4 || len(pin) > 12 {
		return nil, fmt.Errorf("pin length must be between 4 and 12 digits")
	}
	if err := checkDigits(pin, "pin"); err != nil {
		return nil, err
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
		return nil, err
	}
	return pack(nibbles), nil
}

// Decode returns the PIN from an 8-byte PIN block. account is not used. See
// Format.
func (i *eciObject) Decode(pinBlock, _ []byte) ([]byte, error) {
	if len(pinBlock) != 8 {
		return nil, fmt.Errorf("pin block must be 8 bytes")
	}

	nibbles := unpack(pinBlock)
	defer wipe(nibbles)

	if i.version == eci2Version {
		return pinFromNibbles(nibbles[:4])
	}
	pinLength := int(nibbles[0])
	if pinLength < 4 || pinLength > 6 {
		return nil, fmt.Errorf("unable to parse pin block")
	}
	return pinFromNibbles(nibbles[1 : 1+pinLength])
}
