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
go get github.com/moov-io/pinblock
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

All the pin block types implemented based on two different interfaces
```
type FormatA interface {
	SetDebugWriter(writer io.Writer)
	Encode(pin, account string) (string, error)
	Decode(pinBlock, account string) (string, error)
}

type FormatB interface {
	SetDebugWriter(writer io.Writer)
	Encode(pin string) (string, error)
	Decode(pinBlock string) (string, error)
}
```
FormatA used pin and pan code for encrypt operation and decode operation, but FormatB used pin code only.
Generally user can not access code block instance (struct instance), should use via interface 

There are some global functions to get pin blocker instance

```
func NewISO0() FormatA
func NewISO1() FormatB
func NewISO2() FormatB
func NewISO3() FormatA
func NewISO4(cipher Cipher) FormatA
func NewANSIX98() FormatA
func NewOEM1() FormatB
func NewECI1() FormatA
func NewECI2() FormatB
func NewECI3() FormatB
func NewECI4() FormatB
func NewVISA1() FormatA
func NewVISA2() FormatB
func NewVISA3() FormatB
func NewVISA4() FormatA
```

To use a set of method signatures according to interface is very easy
```
		pin := "1234"
		account := "5432101234567891"
		
		iso0 := formats.NewISO0()
		pinBlock, err := iso0.Encode(pin, account)
		...
		pin, err := iso0.Decode(pinBlock, account)
		...
```

Especially pin block of iso-4 required that pin block is encrypted with AES key.
To support this type, encryption logic implemented in moov pinblock package
```
		cipher, err := encryption.NewAesECB([]byte("1234567890123456"))
		iso4 := NewISO4(cipher)
		pinBlock, err := iso4.Encode("12344", "432198765432109870")
		pin, err := iso4.Decode(pinBlock, "432198765432109870")
```

User can get debug messages that describe operation status intuitively with SetDebugWriter() function.
```
		pin := "1234"
		account := "5432101234567891"
		iso0 := formats.NewISO0()
		out := bytes.NewBuffer([]byte{})
		iso0.SetDebugWriter(out)
		pinBlock, err := iso0.Encode(pin, account)
		fmt.Println(out.String())
		
		out.Reset()
		
		pin, err := iso0.Decode(pinBlock, account)
		fmt.Println(out.String())
		
```

Expected output messages from above code.
```
PIN block encode operation finished
************************************
PAN     : 5432101234567891
PIN     : 1234
PAD     : FFFFFFFFFF
Format  : Format 0 (ISO-0)
------------------------------------
Formatted PIN block  : 041215FEDCBA9876
Formatted PAN block  : 0000210123456789

PIN block decode operation finished
************************************
Formatted PAN block  : 0000210123456789
Formatted PIN block  : 041234FFFFFFFFFF
PAD                  : FFFFFFFFFF
Format               : Format 0 (ISO-0)
------------------------------------
Decoded PIN  : 1234

```

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
