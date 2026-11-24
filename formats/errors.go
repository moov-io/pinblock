package formats

import "errors"

// Errors that callers can test with errors.Is.
//
// Decode cannot tell a PIN block enciphered under a different key from a
// malformed block: a block deciphered under the wrong key is random data, and
// Decode rejects it as malformed. Both cases match ErrInvalidPINBlock, so a
// caller that must tell them apart (for example, to choose between ISO 8583
// response codes 55 and 81) must use other information, such as a key check
// value.
var (
	// ErrUnsupportedFormat is returned by NewFormatter for a name it does
	// not know.
	ErrUnsupportedFormat = errors.New("unsupported pinblock type")

	// ErrCipherRequired is returned when ISO-4 has no cipher. NewFormatter
	// returns it for "ISO-4"; use NewISO4 with a cipher.
	ErrCipherRequired = errors.New("iso-4 requires a cipher")

	// ErrInvalidPINBlock matches every error that Decode returns for a PIN
	// block that does not have the structure of its format: the wrong
	// length, characters that are not hex, the wrong format nibble, a PIN
	// length out of range, a PIN digit that is not decimal or the wrong
	// fill.
	ErrInvalidPINBlock = errors.New("invalid pin block")

	// ErrFormatMismatch is returned when the format nibble of the block is
	// not the nibble of the format that decodes it. It matches
	// ErrInvalidPINBlock.
	ErrFormatMismatch error = &blockError{"format is different"}

	// ErrInvalidPINLength is returned when the PIN length in the block is
	// out of range. It matches ErrInvalidPINBlock.
	ErrInvalidPINLength error = &blockError{"invalid pin length"}

	// ErrInvalidFill is returned when the fill in the block is not the fill
	// the format defines. It matches ErrInvalidPINBlock.
	ErrInvalidFill error = &blockError{"invalid fill in pin block"}
)

// Other errors for a malformed block. They are not exported, but they match
// ErrInvalidPINBlock.
var (
	errBlockLength16   error = &blockError{"pin block must be 16 characters"}
	errBlockLength32   error = &blockError{"pinBlock must be 32 hex characters (16 bytes)"}
	errBlockNotHex     error = &blockError{"pin block must contain only hex characters"}
	errBlockPINDigit   error = &blockError{"pin block holds a non-decimal PIN digit"}
	errBlockUnparsable error = &blockError{"unable to parse pin block"}
)

// blockError is an error for a malformed PIN block. It keeps the error text
// that the package returned before these errors were exported.
type blockError struct {
	msg string
}

func (e *blockError) Error() string {
	return e.msg
}

// Is makes every blockError match ErrInvalidPINBlock.
func (e *blockError) Is(target error) bool {
	return target == ErrInvalidPINBlock
}
