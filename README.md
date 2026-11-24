[![Moov Banner Logo](https://user-images.githubusercontent.com/20115216/104214617-885b3c80-53ec-11eb-8ce0-9fc745fb5bfc.png)](https://github.com/moov-io)

[![GoDoc](https://pkg.go.dev/badge/github.com/moov-io/pinblock/formats.svg)](https://pkg.go.dev/github.com/moov-io/pinblock/formats)

# moov-io/pinblock

The `pinblock` package implements Personal Identification Number (PIN) management and security.

A block of data used to encapsulate a PIN during processing. The PIN block format defines the content of the PIN block and how it is processed to retrieve the PIN. The PIN block is composed of the PIN, the PIN length, and may contain subset of the PAN. There are many formats for PIN blocks that are network and processor dependent. Regardless the data is shipped in [ISO8583](https://github.com/moov-io/iso8583) field 52. 


## Table of contents

- [Project status](#project-status)
- [Go module](#go-library)
    - [Installation](#installation)
    - [Pin blocks](#supported-pin-blocks)
- [Usage](#usage)
    - [Clear and enciphered PIN blocks](#clear-and-enciphered-pin-blocks)
    - [Select a format by name](#select-a-format-by-name)
    - [Validation](#validation)
- [Docs](#docs)
- [Other examples](#other-examples)
- [Getting help](#getting-help)
- [Contributing](#contributing)
- [Related projects](#related-projects)
- [License](#license)

## Project status

Moov pinblock currently offers a Go package that supports decode and encode features for a lot of pinblock formats. Feedback on this early version of the project is appreciated and vital to its success. Please let us know if you encounter any bugs/unclear documentation or have feature suggestions by opening up an issue. Thanks!

## Go library

This project uses [Go Modules](https://go.dev/blog/using-go-modules) and Go v1.22.2 or newer. See [Golang's install instructions](https://golang.org/doc/install) for help in setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/pinblock/releases/latest) as well. We highly recommend you use a tagged release for production.


### Installation

```
go get github.com/moov-io/pinblock@latest
```

The PIN block formats are in the `github.com/moov-io/pinblock/formats` package. The AES cipher for ISO-4 is in the `github.com/moov-io/pinblock/encryption` package. This project is a Go library only. It does not have a CLI, a server or a Docker image.

### Supported pin blocks

- [x] ISO 9564-1 Format 0 (`NewISO0`)
- [x] ISO 9564-1 Format 1 (`NewISO1`)
- [x] ISO 9564-1 Format 2 (`NewISO2`)
- [x] ISO 9564-1 Format 3 (`NewISO3`)
- [x] ISO 9564-1 Format 4, AES (`NewISO4`)
- [x] ANSI X9.8 (`NewANSIX98`)
- [x] OEM-1 (Diebold, Docutel, NCR) (`NewOEM1`)
- [x] ECI-1, ECI-2, ECI-3, ECI-4 (`NewECI1` to `NewECI4`)
- [x] VISA-1, VISA-2, VISA-3, VISA-4 (`NewVISA1` to `NewVISA4`)

## Usage

All the pin block types implement one interface, `formats.Format`:
```
type Format interface {
	SetDebugWriter(writer io.Writer)
	Encode(pin, account string) (string, error)
	Decode(pinBlock, account string) (string, error)
}
```
`pin` is the PIN as decimal digits and `account` is the PAN as decimal digits. Formats that do not use the PAN (ISO-1, ISO-2, OEM-1, ECI-2, ECI-3, ECI-4, VISA-2 and VISA-3) ignore `account`, so give `""`. A PIN block is upper-case hex text: 16 characters, or 32 characters for ISO-4.

Every constructor returns a `formats.Format`:

```
func NewISO0() Format
func NewISO1() Format
func NewISO2() Format
func NewISO3() Format
func NewISO4(cipher Cipher) Format
func NewANSIX98() Format
func NewOEM1() Format
func NewECI1() Format
func NewECI2() Format
func NewECI3() Format
func NewECI4() Format
func NewVISA1() Format
func NewVISA2() Format
func NewVISA3() Format
func NewVISA4() Format
```

Encode and decode an ISO-0 PIN block:
```go
package main

import (
	"fmt"
	"log"

	"github.com/moov-io/pinblock/formats"
)

func main() {
	pin := "1234"
	account := "43219876543210987"

	iso0 := formats.NewISO0()

	pinBlock, err := iso0.Encode(pin, account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pinBlock) // 0412AC89ABCDEF67

	decoded, err := iso0.Decode(pinBlock, account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(decoded) // 1234
}
```

ISO-0 and ISO-3 use the 12 rightmost digits of the PAN, without the check digit, so the PAN must have 13 digits or more. ISO-4 accepts a PAN of 12 to 19 digits.

An ISO-4 pin block is enciphered with an AES key. `encryption.NewAesECB` gives a `formats.Cipher` for a 16, 24 or 32-byte key:
```go
key, err := hex.DecodeString("00112233445566778899AABBCCDDEEFF")
if err != nil {
	log.Fatal(err)
}
cipher, err := encryption.NewAesECB(key)
if err != nil {
	log.Fatal(err)
}

iso4 := formats.NewISO4(cipher)

pinBlock, err := iso4.Encode("1234", "432198765432109870")
if err != nil {
	log.Fatal(err)
}
pin, err := iso4.Decode(pinBlock, "432198765432109870")
if err != nil {
	log.Fatal(err)
}
```
The ISO-4 block contains random fill, so `Encode` gives a different block each time.

These examples are tested in [formats/example_test.go](formats/example_test.go).

### Clear and enciphered PIN blocks

Only ISO-4 enciphers the PIN block. All the other formats return a **clear** PIN block from `Encode`, and `Decode` expects a clear PIN block. A clear PIN block is not safe to send: for ISO-0 and ISO-3, anyone who has the PAN can get the PIN from it. Before a clear block leaves the terminal or the host, encipher it under a PIN key (usually TDES), and decipher it before you call `Decode`.

This package does not include a TDES cipher. The `formats.Cipher` interface has a 16-byte block and it is only for ISO-4. Use `crypto/des` (or your HSM) for TDES. [`Example_translatePIN`](formats/example_test.go) shows a PIN translation: it deciphers a TDES-enciphered ISO-0 block, decodes the PIN, and encodes it again as an ISO-4 block under an AES key.

### Select a format by name

`formats.NewFormatter(name)` returns a `Format` by name, for example `"ISO-0"`, `"ISO-3"`, `"ANSI"`, `"ECI1"` or `"VISA1"`. See [its godoc](https://pkg.go.dev/github.com/moov-io/pinblock/formats#NewFormatter) for all the names. Use `NewISO4` with a cipher for ISO-4. Do not use `NewFormatter` for ISO-4.

### Debug output

`SetDebugWriter()` writes a short note after each successful operation, naming the operation and the format. It never writes the PIN, the PAN, the PIN length, the fill or any PIN block.
```
		iso0 := formats.NewISO0()
		out := bytes.NewBuffer([]byte{})
		iso0.SetDebugWriter(out)
		pinBlock, err := iso0.Encode("1234", "5432101234567891")
		fmt.Println(out.String())
```

Expected output:
```
PIN block encode operation finished
Format  : Format 0 (ISO-0)

```

### Validation

`Encode` accepts only decimal digits in the PIN and the PAN. A PIN must have 4 to 12 digits; OEM-1 is different, because it truncates a PIN that is longer than 12 digits, as the format defines. `Decode` checks the block's structure (format nibble, PIN length 4 to 12, decimal PIN digits and, for the ISO formats, the fill each one defines) and returns an error for a malformed block, which is also what a block decrypted under the wrong key looks like. PIN blocks are returned as upper-case hex; `Decode` accepts either case.

`Decode` cannot tell a block that was enciphered under a different key from a malformed block, because a block deciphered under the wrong key is random data. Use a key check value or other information if you must know which key was used.

### Handling PINs in memory

Every intermediate the package builds (the PIN field, the PAN field, XOR results, random fill, cipher outputs and the text being converted) is held in a byte slice and cleared before `Encode` or `Decode` returns, including when it returns an error.

The PIN, the PAN and the PIN block passed in and returned are Go strings, and a string cannot be overwritten: they stay in memory until the garbage collector reuses it, and neither this package nor the caller can clear them. Treat the PIN block as sensitive too: for ISO-0 and ISO-3 it is the PIN field XORed with the PAN. An API that takes and returns byte slices the caller can clear is proposed as a separate enhancement. `pinblock.Clear` is available for clearing your own byte buffers.

Copies the package cannot reach are not cleared either: the AES key schedule inside `crypto/aes`, copies in CPU registers or on goroutine stacks, and anything a `Cipher` implementation keeps.

## Docs

[ISO 9564 Wikipedia](https://en.wikipedia.org/wiki/ISO_9564)

[Complete list of Pin-blocks](https://www.eftlab.com/knowledge-base/complete-list-of-pin-blocks)

[MasterCard Pin block encryption](https://developer.mastercard.com/card-issuance/documentation/pin-block-encryption-process/)

[PINBlock Explained](https://www.linkedin.com/pulse/pinblock-explained-iftekharul-haque/)

[Implementing ISO Format 4 PIN Blocks](https://listings.pcisecuritystandards.org/documents/Implementing_ISO_Format_4_PIN_Blocks_Information_Supplement.pdf)

## Other Examples 

[Newpay calculate PIN block for ISO8583](https://neapay.com/online-tools/calculate-pin-block.html)

[EMV Labs PIN Block calculator](https://emvlab.org/pinblock/)

[Payment Card Tools](https://paymentcardtools.com/pin-block-calculators/iso9564-format-0)


## Getting help

channel | info
 ------- | -------
[Project Documentation](https://pkg.go.dev/github.com/moov-io/pinblock/formats) | Our Go package documentation on pkg.go.dev.
Twitter [@moov](https://twitter.com/moov)	| You can follow Moov.io's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io/pinblock/issues/new) | If you are able to reproduce a problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](https://slack.moov.io/) | Join our slack channel (`#pinblock`) to have an interactive discussion about the development of the project.

## Contributing

Please review our [Contributing guide](https://github.com/moov-io/ach/blob/master/CONTRIBUTING.md) and [Code of Conduct](https://github.com/moov-io/ach/blob/master/CODE_OF_CONDUCT.md) to get started! Check out our [issues for first time contributors](https://github.com/moov-io/pinblock/contribute) for something to help out with.

This project uses [Go Modules](https://go.dev/blog/using-go-modules) and Go v1.22.2 or newer. See [Golang's install instructions](https://golang.org/doc/install) for help setting up Go. Run the tests with `go test ./...`. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/pinblock/releases/latest) as well. We highly recommend you use a tagged release for production.

## Related projects
As part of Moov's initiative to offer open source fintech infrastructure, we have a large collection of active projects you may find useful:

- [Moov ACH](https://github.com/moov-io/ach) provides ACH file generation and parsing, supporting all Standard Entry Codes for the primary method of money movement throughout the United States.

- [Moov Watchman](https://github.com/moov-io/watchman) offers search functions over numerous trade sanction lists from the United States and European Union.

- [Moov Fed](https://github.com/moov-io/fed) implements utility services for searching the United States Federal Reserve System such as ABA routing numbers, financial institution name lookup, and FedACH and Fedwire routing information.

- [Moov Wire](https://github.com/moov-io/wire) implements an interface to write files for the Fedwire Funds Service, a real-time gross settlement funds transfer system operated by the United States Federal Reserve Banks.

- [Moov ImageCashLetter](https://github.com/moov-io/imagecashletter) implements Image Cash Letter (ICL) files used for Check21, X.9 or check truncation files for exchange and remote deposit in the U.S.

- [Moov Metro2](https://github.com/moov-io/metro2) provides a way to easily read, create, and validate Metro 2 format, which is used for consumer credit history reporting by the United States credit bureaus.

- [Moov iso8583](https://github.com/moov-io/iso8583) implements an ISO 8583 message reader and writer in Go.

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
