package formats

import "fmt"

// iso4BlockLen is the length in bytes of an ISO-4 PIN block: one AES block.
const iso4BlockLen = 16

// iso4Object implements ISO 9564-1 Format 4, for AES.
//
// Plain PIN field (32 nibbles): C N P1 .. P4 P5/A .. P12/A A13 A14 R1 .. R16
//
//	C = control field, 4
//	N = PIN length, 4 to 12
//	P = PIN digit
//	A = fill, 0xA, up to nibble 16
//	R = 16 random nibbles (8 random bytes)
//
// PAN field (32 nibbles): M A1 ... An 0 ... 0, where M is the PAN length minus
// 12 (0 to 7) and the PAN is 12 to 19 digits, left-justified and zero-padded.
//
// Encipherment: block A = E(plain PIN field); block B = A XOR PAN field; the
// PIN block is E(block B). Decipherment reverses the steps.
type iso4Object struct {
	cipher Cipher
	format string
}

// SetCipher replaces the cipher used to encipher and decipher PIN blocks. A
// nil cipher is ignored.
func (i *iso4Object) SetCipher(cipher Cipher) {
	if i == nil || cipher == nil {
		return
	}
	i.cipher = cipher
}

// panField returns the 32-nibble ISO-4 PAN field packed into 16 bytes. The
// result belongs to the caller, who must clear it.
func (i *iso4Object) panField(account []byte) ([]byte, error) {
	if len(account) < 12 || len(account) > 19 {
		return nil, fmt.Errorf("account length must be between 12 and 19 digits")
	}
	nibbles := make([]byte, 2*iso4BlockLen)
	defer wipe(nibbles)
	nibbles[0] = byte((len(account) - 12) & 0x0F) // 0..7, checked above
	if err := digitNibbles(nibbles[1:], account, "account"); err != nil {
		return nil, err
	}
	return pack(nibbles), nil
}

// apply runs one cipher operation and returns its output in a new slice of
// exactly one block. It clears the whole capacity of the slice the cipher
// returned, so a cipher that returns part of a larger buffer has all of it
// cleared. A cipher that returns the wrong length is an error.
func apply(op func([]byte) ([]byte, error), in []byte, step string) ([]byte, error) {
	out, err := op(in)
	defer wipe(out[:cap(out)])
	if err != nil {
		return nil, fmt.Errorf("%s: %w", step, err)
	}
	if len(out) != iso4BlockLen {
		return nil, fmt.Errorf("%s: cipher returned %d bytes, want %d", step, len(out), iso4BlockLen)
	}
	result := make([]byte, iso4BlockLen)
	copy(result, out)
	return result, nil
}

// Encode returns the 16-byte enciphered PIN block for pin and account. See
// Format.
func (i *iso4Object) Encode(pin, account []byte) ([]byte, error) {
	if i.cipher == nil {
		return nil, fmt.Errorf("iso-4 requires a cipher")
	}
	if len(pin) < 4 || len(pin) > 12 {
		return nil, fmt.Errorf("pin length must be between 4 and 12 digits")
	}

	panField, err := i.panField(account)
	if err != nil {
		return nil, err
	}
	defer wipe(panField)

	nibbles := make([]byte, 2*iso4BlockLen)
	defer wipe(nibbles)
	nibbles[0] = 0x4
	nibbles[1] = byte(len(pin) & 0x0F) // 4..12, checked above
	if err := digitNibbles(nibbles[2:], pin, "pin"); err != nil {
		return nil, err
	}
	for n := 2 + len(pin); n < fieldNibbles; n++ {
		nibbles[n] = 0xA
	}
	if err := randomFill(nibbles[fieldNibbles:], anyNibble); err != nil {
		return nil, err
	}

	pinField := pack(nibbles)
	defer wipe(pinField)

	blockA, err := apply(i.cipher.Encrypt, pinField, "encrypting pin field")
	if err != nil {
		return nil, err
	}
	defer wipe(blockA)

	blockB := make([]byte, iso4BlockLen)
	defer wipe(blockB)
	xorInto(blockB, blockA, panField)

	return apply(i.cipher.Encrypt, blockB, "encrypting block B")
}

// Decode returns the PIN from a 16-byte enciphered PIN block. See Format.
func (i *iso4Object) Decode(pinBlock, account []byte) ([]byte, error) {
	if i.cipher == nil {
		return nil, fmt.Errorf("iso-4 requires a cipher")
	}
	if len(pinBlock) != iso4BlockLen {
		return nil, fmt.Errorf("pin block must be %d bytes", iso4BlockLen)
	}

	panField, err := i.panField(account)
	if err != nil {
		return nil, err
	}
	defer wipe(panField)

	// A copy, so a cipher that works in place cannot modify the caller's block.
	encrypted := make([]byte, iso4BlockLen)
	copy(encrypted, pinBlock)
	defer wipe(encrypted)

	blockB, err := apply(i.cipher.Decrypt, encrypted, "decrypting pin block")
	if err != nil {
		return nil, err
	}
	defer wipe(blockB)

	blockA := make([]byte, iso4BlockLen)
	defer wipe(blockA)
	xorInto(blockA, blockB, panField)

	pinField, err := apply(i.cipher.Decrypt, blockA, "decrypting block A")
	if err != nil {
		return nil, err
	}
	defer wipe(pinField)

	nibbles := unpack(pinField)
	defer wipe(nibbles)

	if nibbles[0] != 0x4 {
		return nil, fmt.Errorf("format is different")
	}
	pinLength := int(nibbles[1])
	if pinLength < 4 || pinLength > 12 {
		return nil, fmt.Errorf("invalid pin length %d", pinLength)
	}
	for _, n := range nibbles[2+pinLength : fieldNibbles] {
		if n != 0xA {
			return nil, fmt.Errorf("invalid fill in pin block")
		}
	}
	return pinFromNibbles(nibbles[2 : 2+pinLength])
}
