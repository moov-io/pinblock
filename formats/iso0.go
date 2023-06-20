package formats

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type ISO0 struct {
	Filler string

	debugWriter io.Writer
}

func NewISO0() *ISO0 {
	return &ISO0{
		Filler: "F", // default to ISO0's Filler
	}
}

// Format returns iso type
func (i *ISO0) format() string {
	return "Format 0 (ISO-0)"
}

// Padding returns padding pattern
func (i *ISO0) padding(pin string) (string, error) {

	if len(pin) < 4 || len(pin) > 12 {
		return "", fmt.Errorf("pin length must be between 4 and 12 digits")
	}

	return strings.Repeat(i.Filler, 14-len(pin)), nil
}

// SetDebugWriter will set writer for getting output message of encoding and decoding logic
func (i *ISO0) SetDebugWriter(writer io.Writer) {
	i.debugWriter = tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
}

// Encode returns the ISO0 PIN block for the given PIN and account number
func (i *ISO0) Encode(pin, account string) (string, error) {

	pad, err := i.padding(pin)
	if err != nil {
		return "", err
	}

	// pin block should start with 0, then add length of pin, then add pin,
	// then add F until 16 characters
	pinBlock := fmt.Sprintf("0%d%s%s", len(pin), pin, pad)

	// account number must be at least 13 digits, including the check digit
	if len(account) < 13 {
		return "", fmt.Errorf("account length must be at least 13 digits")
	}

	// take the last 12 digits of the account number excluding the check digit
	accountBlock := fmt.Sprintf("0000%s", account[len(account)-13:len(account)-1])

	xorBlock, err := xorHex(pinBlock, accountBlock)
	if err != nil {
		return "", err
	}

	// write encode information
	if i.debugWriter != nil {
		tw := i.debugWriter
		fmt.Fprintf(tw, "PIN block encode operation finished\n")
		fmt.Fprintf(tw, "%s\n", strings.Repeat("*", 36))
		fmt.Fprintf(tw, "PAN\t: %s\n", account)
		fmt.Fprintf(tw, "PIN\t: %s\n", pin)
		if pad, _ := i.padding(pin); len(pad) == 0 {
			fmt.Fprintf(tw, "PAD\t: N/A\n")
		} else {
			fmt.Fprintf(tw, "PAD\t: %s\n", pad)
		}
		fmt.Fprintf(tw, "Format\t: %s\n", i.format())
		fmt.Fprintf(tw, "%s\n", strings.Repeat("-", 36))
		fmt.Fprintf(tw, "Formatted PIN block\t: %s\n", strings.ToUpper(xorBlock))
		fmt.Fprintf(tw, "Formatted PAN block\t: %s\n", strings.ToUpper(accountBlock))
		tw.Write([]byte("\n"))
	}

	return strings.ToUpper(xorBlock), nil
}

func (i *ISO0) Decode(pinBlock, account string) (string, error) {
	if len(pinBlock) != 16 {
		return "", fmt.Errorf("pin block must be 16 characters")
	}

	if len(account) < 13 {
		return "", fmt.Errorf("account length must be at least 13 digits")
	}

	// take the last 12 digits of the account number excluding the check digit
	accountBlock := fmt.Sprintf("0000%s", account[len(account)-13:len(account)-1])

	decodedBlock, err := xorHex(pinBlock, accountBlock)
	if err != nil {
		return "", err
	}

	// decodedBlock should start with 0, then has length of pin, then has pin, then has F until 16 characters
	pinLength := int(decodedBlock[1] - '0')
	pin := decodedBlock[2 : 2+pinLength]

	// write decode information
	if i.debugWriter != nil {
		tw := i.debugWriter
		fmt.Fprintf(tw, "PIN block decode operation finished\n")
		fmt.Fprintf(tw, "%s\n", strings.Repeat("*", 36))
		fmt.Fprintf(tw, "Formatted PAN block\t: %s\n", strings.ToUpper(accountBlock))
		fmt.Fprintf(tw, "Formatted PIN block\t: %s\n", strings.ToUpper(decodedBlock))
		if pad, _ := i.padding(pin); len(pad) == 0 {
			fmt.Fprintf(tw, "PAD\t: N/A\n")
		} else {
			fmt.Fprintf(tw, "PAD\t: %s\n", pad)
		}
		fmt.Fprintf(tw, "Format\t: %s\n", i.format())
		fmt.Fprintf(tw, "%s\n", strings.Repeat("-", 36))
		fmt.Fprintf(tw, "Decoded PIN\t: %s\n", pin)
		tw.Write([]byte("\n"))
	}

	return pin, nil
}
