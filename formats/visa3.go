package formats

import "fmt"

// visa3Delimiter is the nibble that ends the PIN in a VISA-3 block.
const visa3Delimiter = 0xF

// visa3Object implements VISA-3. It does not use the PAN.
//
// PIN block (16 nibbles): the PIN (4 to 12 digits), the delimiter 0xF, then
// one random value from 0xA to 0xF repeated to the end.
type visa3Object struct {
	format string
}

// Encode returns the 8-byte PIN block for pin. account is not used. A PIN
// longer than 12 digits is truncated on the right, as the format defines. See
// Format.
func (i *visa3Object) Encode(pin, _ []byte) ([]byte, error) {
	pin = pin[:min(len(pin), 12)]
	if len(pin) < 4 {
		return nil, fmt.Errorf("pin length must be between 4 and 12 digits")
	}

	nibbles := make([]byte, fieldNibbles)
	defer wipe(nibbles)
	if err := digitNibbles(nibbles, pin, "pin"); err != nil {
		return nil, err
	}
	nibbles[len(pin)] = visa3Delimiter

	filler := make([]byte, 1)
	defer wipe(filler)
	if err := randomFill(filler, letterNibble); err != nil {
		return nil, err
	}
	for n := len(pin) + 1; n < fieldNibbles; n++ {
		nibbles[n] = filler[0]
	}
	return pack(nibbles), nil
}

// Decode returns the PIN from an 8-byte PIN block. account is not used. See
// Format.
func (i *visa3Object) Decode(pinBlock, _ []byte) ([]byte, error) {
	if len(pinBlock) != 8 {
		return nil, fmt.Errorf("pin block must be 8 bytes")
	}

	nibbles := unpack(pinBlock)
	defer wipe(nibbles)

	index := -1
	for n, v := range nibbles {
		if v == visa3Delimiter {
			index = n
			break
		}
	}
	if index < 4 || index > 12 {
		return nil, fmt.Errorf("unable to parse pin block")
	}
	return pinFromNibbles(nibbles[:index])
}
