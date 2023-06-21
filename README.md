[![Moov Banner Logo](https://user-images.githubusercontent.com/20115216/104214617-885b3c80-53ec-11eb-8ce0-9fc745fb5bfc.png)](https://github.com/moov-io)

# moov-io/pinblock

The `pinblock` package implements Personal Identifiation (PIN) managment and security.

A block of data used to encapsulate a PIN during processing. The PIN block format defines the content of the PIN block and how it is processed to retrieve the PIN. The PIN block is composed of the PIN, the PIN length, and may contain subset of the PAN. Their are many formats for PIN blocks that are network and processor dependent. Regardless the data is shipped in [ISO8583](https://github.com/moov-io/iso8583) field 52. 

# todo 

- [x] PIN block ISO-0 (ISO 9564-1 Format 0, ANSI X9.8, VISA-1, ECI-0, Format 0)
- [x] PIN block ISO-1 (ISO 9564-1:2003 Format 1, ECI-4)
- [ ] PIN block ISO-2 (ISO 9564-3: 2003 Format 2)
- [ ] PIN block ISO-3 (ISO 9564-1: 2002 Format 3)
- [x] PIN block ISO-4 (ISO 9564-1: 2017 Format 4)
- [ ] PIN block Docutel 2 (VISA-2)


# further reading 

[ISO 9564 Wikipedia](https://en.wikipedia.org/wiki/ISO_9564)

[Complete list of Pin-blocks](https://www.eftlab.com/knowledge-base/complete-list-of-pin-blocks)

[MasterCard Pin block encryption](https://developer.mastercard.com/card-issuance/documentation/pin-block-encryption-process/)

[ICMA PIN block formats](http://icma.com/wp-content/uploads/2015/07/PinBlockFormats_SE1-15CM.pdf)

[PINBlock Explained](https://www.linkedin.com/pulse/pinblock-explained-iftekharul-haque/)

[Implementing ISO Format 4 PIN Blocks](https://listings.pcisecuritystandards.org/documents/Implementing_ISO_Format_4_PIN_Blocks_Information_Supplement.pdf)

# Example implementations 

[Newpay calculate PIN block for ISO8583](https://neapay.com/online-tools/calculate-pin-block.html)

[EMV Labs PIN Block calculator](https://emvlab.org/pinblock/)

[Payment Card Tools](https://paymentcardtools.com/pin-block-calculators/iso9564-format-0)


## Getting help

 channel | info
 ------- | -------
Twitter [@moov](https://twitter.com/moov)	| You can follow Moov.io's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io/pinblock/issues/new) | If you are able to reproduce a problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](https://slack.moov.io/) | Join our slack channel to have an interactive discussion about the development of the project.

## Supported and tested platforms

- 64-bit Linux (Ubuntu, Debian), macOS, and Windows

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
