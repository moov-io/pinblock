[![Moov Banner Logo](https://user-images.githubusercontent.com/20115216/104214617-885b3c80-53ec-11eb-8ce0-9fc745fb5bfc.png)](https://github.com/moov-io)

# moov-io/pinblock

The `pinblock` package implements Personal Identifiation (PIN) managment and security.

A block of data used to encapsulate a PIN during processing. The PIN block format defines the content of the PIN block and how it is processed to retrieve the PIN. The PIN block is composed of the PIN, the PIN length, and may contain subset of the PAN. Their are many formats for PIN blocks that are network and processor dependent. Regardless the data is shipped in [ISO8583](https://github.com/moov-io/iso8583) field 52. 


## Table of contents

- [Project status](#project-status)
- [Go module](#go-library)
    - [Installation](#installation)
    - [Pin blocks](#supported-pin-blocks)
- [Usage](#usage)
- [Docs](#docs)
- [Other examples](#other-examples)
- [Getting help](#getting-help)
- [Contributing](#contributing)
- [Related projects](#related-projects)
- [License](#license)

## Project status

Moov pinblock currently offers a Go package that support decode and encode features for a lot of pinblock formats. Feedback on this early version of the project is appreciated and vital to its success. Please let us know if you encounter any bugs/unclear documentation or have feature suggestions by opening up an issue. Thanks!

## Go library

This project uses [Go Modules](https://go.dev/blog/using-go-modules) and Go v1.18 or newer. See [Golang's install instructions](https://golang.org/doc/install) for help in setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/pinblock/releases/latest) as well. We highly recommend you use a tagged release for production.


### Installation

```
go get github.com/moov-io/pinblock@latest
```

### Supported pin blocks

- [x] ISO 9564-1 Format 0
- [x] ISO 9564-1:2003 Format 1
- [x] ISO 9564-3: 2003 Format 2
- [x] ISO 9564-1: 2002 Format 3
- [x] ANSI X9.8
- [x] OEM-1 (Diebold, Docutel, NCR)
- [x] ECI-1, ECI-2, ECI-3, ECI-4

## Usage

Every format implements one interface. PIN, PAN and PIN block are all `[]byte`:

```go
type Format interface {
	Encode(pin, account []byte) ([]byte, error)
	Decode(pinBlock, account []byte) ([]byte, error)
}
```

- `pin` is the PIN as ASCII decimal digits, for example `[]byte("1234")`.
- `account` is the PAN as ASCII decimal digits. Formats that do not use the PAN
  (ISO-1, ISO-2, OEM-1, ECI-2, ECI-3, ECI-4, VISA-2, VISA-3) ignore it; pass `nil`.
- A PIN block is binary: 8 bytes, or 16 bytes for ISO-4. Hex-encode it yourself
  if a message needs text.

Constructors:

```go
func NewFormatter(bType string) (Format, error) // "ISO-0", "ISO-1", ..., "VISA4"
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

```go
pin := []byte("1234")
account := []byte("5432101234567891")

iso0 := formats.NewISO0()
pinBlock, err := iso0.Encode(pin, account) // 04 12 15 FE DC BA 98 76
...
decoded, err := iso0.Decode(pinBlock, account)
...
// import "github.com/moov-io/pinblock"
pinblock.Clear(pin)
pinblock.Clear(account)
pinblock.Clear(pinBlock) // an ISO-0 block XORed with the PAN gives back the PIN
pinblock.Clear(decoded)
```

ISO-4 enciphers the PIN block with AES:

```go
cipher, err := encryption.NewAesECB(key)
iso4 := formats.NewISO4(cipher)
pinBlock, err := iso4.Encode([]byte("12344"), []byte("432198765432109870"))
pin, err := iso4.Decode(pinBlock, []byte("432198765432109870"))
```

### Handling PINs in memory

A Go `string` cannot be overwritten, so a PIN that becomes a string stays in
memory until the garbage collector reuses it. This package never converts a
PIN, PAN or PIN block to a string. Neither method modifies its arguments, and
every result is a new slice the caller owns and should clear with
`pinblock.Clear` once it is no longer needed. Prefer it to the `clear` builtin: a
`clear` with no later read is a dead store the compiler may remove, while
`pinblock.Clear` always writes. Clear the PIN block as well as the PIN: for ISO-0
and ISO-3 the block is the PIN field XORed with the PAN.

Every intermediate the package creates (the PIN field, the PAN field, XOR
results, random fill and cipher outputs) is cleared before the method returns,
including when it returns an error. Copies the package cannot reach are not
cleared: the AES key schedule inside `crypto/aes`, copies in CPU registers or on
goroutine stacks, and anything your own code or a `Cipher` implementation copies.
A host that needs registers and stacks erased can wrap its calls in
`runtime/secret.Do`, which is experimental and only erases on linux/amd64 and
linux/arm64.

Decode checks the block's structure (format nibble, PIN length, decimal digits
and, for ISO formats, the fill) and returns an error for a malformed block, which
is also what a block decrypted under the wrong key looks like.

## Docs

[ISO 9564 Wikipedia](https://en.wikipedia.org/wiki/ISO_9564)

[Complete list of Pin-blocks](https://www.eftlab.com/knowledge-base/complete-list-of-pin-blocks)

[MasterCard Pin block encryption](https://developer.mastercard.com/card-issuance/documentation/pin-block-encryption-process/)

[ICMA PIN block formats](http://icma.com/wp-content/uploads/2015/07/PinBlockFormats_SE1-15CM.pdf)

[PINBlock Explained](https://www.linkedin.com/pulse/pinblock-explained-iftekharul-haque/)

[Implementing ISO Format 4 PIN Blocks](https://listings.pcisecuritystandards.org/documents/Implementing_ISO_Format_4_PIN_Blocks_Information_Supplement.pdf)

## Other Examples 

[Newpay calculate PIN block for ISO8583](https://neapay.com/online-tools/calculate-pin-block.html)

[EMV Labs PIN Block calculator](https://emvlab.org/pinblock/)

[Payment Card Tools](https://paymentcardtools.com/pin-block-calculators/iso9564-format-0)


## Getting help

channel | info
 ------- | -------
[Project Documentation](https://github.com/moov-io/pinblock/tree/master/docs) | Our project documentation available online.
Twitter [@moov](https://twitter.com/moov)	| You can follow Moov.io's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io/pinblock/issues/new) | If you are able to reproduce a problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](https://slack.moov.io/) | Join our slack channel (`#pinblock`) to have an interactive discussion about the development of the project.

## Contributing

Please review our [Contributing guide](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) to get started! Check out our [issues for first time contributors](https://github.com/moov-io/pinblock/contribute) for something to help out with.

This project uses [Go Modules](https://go.dev/blog/using-go-modules) and Go v1.20 or newer. See [Golang's install instructions](https://golang.org/doc/install) for help setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/pinblock/releases/latest) as well. We highly recommend you use a tagged release for production.

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
