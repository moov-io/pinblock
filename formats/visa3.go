package formats

import "fmt"

// visa3Delimiter is the nibble that ends the PIN in a VISA-3 block.
const visa3Delimiter = 0xF

// visa3Object implements VISA-3. It does not use the PAN.
//
// PIN block (16 nibbles): the PIN (4 to 12 digits), the delimiter 0xF, then
// one random value from 0xA to 0xF repeated to the end.
type visa3Object struct {
	debugOutput

	format string
}

// Encode returns the 16-character PIN block for pin. account is not used. A
// PIN longer than 12 digits is truncated on the right, as the format defines.
// See Format.
func (i *visa3Object) Encode(pin, _ string) (string, error) {
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
	nibbles[len(pin)] = visa3Delimiter

	filler := make([]byte, 1)
	defer wipe(filler)
	if err := randomFill(filler, letterNibble); err != nil {
		return "", err
	}
	for n := len(pin) + 1; n < fieldNibbles; n++ {
		nibbles[n] = filler[0]
	}
	pinBlock := nibblesHex(nibbles)

	i.note("encode", i.format)
	return pinBlock, nil
}

// Decode returns the PIN from a 16-character PIN block. account is not used.
// See Format.
func (i *visa3Object) Decode(pinBlock, _ string) (string, error) {
	if len(pinBlock) != fieldNibbles {
		return "", errBlockLength16
	}

	nibbles, err := hexNibbles(pinBlock)
	if err != nil {
		return "", err
	}
	defer wipe(nibbles)

	index := -1
	for n, v := range nibbles {
		if v == visa3Delimiter {
			index = n
			break
		}
	}
	if index < 4 || index > 12 {
		return "", errBlockUnparsable
	}
	pin, err := pinFromNibbles(nibbles[:index])
	if err != nil {
		return "", err
	}

	i.note("decode", i.format)
	return pin, nil
}
