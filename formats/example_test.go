package formats_test

import (
	"crypto/des"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/moov-io/pinblock/encryption"
	"github.com/moov-io/pinblock/formats"
)

// ISO-0 Encode returns a clear PIN block. It is not enciphered.
func ExampleNewISO0() {
	pin := "1234"
	account := "43219876543210987"

	iso0 := formats.NewISO0()

	pinBlock, err := iso0.Encode(pin, account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pinBlock)

	decoded, err := iso0.Decode(pinBlock, account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(decoded)

	// Output:
	// 0412AC89ABCDEF67
	// 1234
}

// ISO-4 Encode returns a PIN block that is enciphered with AES. The block
// contains random fill, so it is different each time.
func ExampleNewISO4() {
	key, err := hex.DecodeString("00112233445566778899AABBCCDDEEFF")
	if err != nil {
		log.Fatal(err)
	}
	cipher, err := encryption.NewAesECB(key)
	if err != nil {
		log.Fatal(err)
	}

	pin := "1234"
	account := "432198765432109870"

	iso4 := formats.NewISO4(cipher)

	pinBlock, err := iso4.Encode(pin, account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(pinBlock))

	decoded, err := iso4.Decode(pinBlock, account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(decoded)

	// Output:
	// 32
	// 1234
}

// PIN translation: decipher an ISO-0 PIN block that is enciphered under a
// TDES PIN key, then encode the PIN again as an ISO-4 block under an AES key.
// This package does not include a TDES cipher, so the example uses
// crypto/des.
func Example_translatePIN() {
	account := "43219876543210987"

	// The double-length TDES key K1 K2 is used as K1 K2 K1.
	tdesKey, err := hex.DecodeString("0123456789ABCDEFFEDCBA9876543210")
	if err != nil {
		log.Fatal(err)
	}
	tdes, err := des.NewTripleDESCipher(append(tdesKey, tdesKey[:8]...))
	if err != nil {
		log.Fatal(err)
	}

	// The terminal side: encode a clear ISO-0 block and encipher it.
	clearBlock, err := formats.NewISO0().Encode("1234", account)
	if err != nil {
		log.Fatal(err)
	}
	clearBytes, err := hex.DecodeString(clearBlock)
	if err != nil {
		log.Fatal(err)
	}
	enciphered := make([]byte, des.BlockSize)
	tdes.Encrypt(enciphered, clearBytes)
	fmt.Println(strings.ToUpper(hex.EncodeToString(enciphered)))

	// The host side: decipher the block and decode the PIN.
	deciphered := make([]byte, des.BlockSize)
	tdes.Decrypt(deciphered, enciphered)
	pin, err := formats.NewISO0().Decode(strings.ToUpper(hex.EncodeToString(deciphered)), account)
	if err != nil {
		log.Fatal(err)
	}

	// Encode the PIN as an ISO-4 block under the AES key.
	aesKey, err := hex.DecodeString("00112233445566778899AABBCCDDEEFF")
	if err != nil {
		log.Fatal(err)
	}
	cipher, err := encryption.NewAesECB(aesKey)
	if err != nil {
		log.Fatal(err)
	}
	iso4 := formats.NewISO4(cipher)
	iso4Block, err := iso4.Encode(pin, account)
	if err != nil {
		log.Fatal(err)
	}

	decoded, err := iso4.Decode(iso4Block, account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(decoded)

	// Output:
	// C967C8198151A458
	// 1234
}
