package formats

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type visa3Object struct {
	format      string
	debugWriter io.Writer
}

const delimiter = "F"

// Padding returns padding pattern
func (i *visa3Object) padding(pin string) (string, error) {
	if len(pin) < 4 || len(pin) > 12 {
		return "", fmt.Errorf("pin length must be between 4 and 12 digits")
	}

	filler, err := randomLetters(1, hexCharacters)
	if err != nil {
		return "", err
	}

	return delimiter + strings.Repeat(filler, 16-len(pin)-1), nil
}

// SetDebugWriter will set writer for getting output message of encoding and decoding logic
func (i *visa3Object) SetDebugWriter(writer io.Writer) {
	i.debugWriter = tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
}

// Encode returns the OEM-1 PIN block for the given PIN
func (i *visa3Object) Encode(pin, account string) (string, error) {
	isTruncated := false

	// A PIN that is longer than 12 digits is truncated on the right.
	if len(pin) > 12 {
		pin = pin[:12]
		isTruncated = true
	}

	pad, err := i.padding(pin)
	if err != nil {
		return "", err
	}

	// Pad value has a 4-bit value from X’0′ to X’F’ and must be different from any PIN digit.
	// The number of pad values for this format is in the range from 4 to 12, and all the pad values must have the same value.
	pinBlock := fmt.Sprintf("%s%s", pin, pad)

	// write encode information
	if i.debugWriter != nil {
		tw := i.debugWriter
		fmt.Fprintf(tw, "PIN block encode operation finished\n")
		if isTruncated {
			fmt.Fprintf(tw, "The pin is truncated on the right as 12 digits\n")
		}
		fmt.Fprintf(tw, "%s\n", strings.Repeat("*", 36))
		fmt.Fprintf(tw, "PIN\t: %s\n", pin)
		if pad == "" {
			fmt.Fprintf(tw, "PAD\t: N/A\n")
		} else {
			fmt.Fprintf(tw, "PAD\t: %s\n", pad)
		}
		fmt.Fprintf(tw, "Format\t: %s\n", i.format)
		fmt.Fprintf(tw, "%s\n", strings.Repeat("-", 36))
		fmt.Fprintf(tw, "Formatted PIN block\t: %s\n", strings.ToUpper(pinBlock))
		tw.Write([]byte("\n"))
	}

	return strings.ToUpper(pinBlock), nil
}

func (i *visa3Object) Decode(pinBlock, account string) (string, error) {
	if len(pinBlock) != 16 {
		return "", fmt.Errorf("pin block must be 16 characters")
	}

	index := strings.Index(pinBlock, delimiter)
	if index < 4 || index > 12 {
		return "", fmt.Errorf("unable to parse pin block")
	}

	pin := pinBlock[:index]
	// write decode information
	if i.debugWriter != nil {
		tw := i.debugWriter
		fmt.Fprintf(tw, "PIN block decode operation finished\n")
		fmt.Fprintf(tw, "%s\n", strings.Repeat("*", 36))
		fmt.Fprintf(tw, "Formatted PIN block\t: %s\n", strings.ToUpper(pinBlock))
		pad := pinBlock[index+1:]
		fmt.Fprintf(tw, "PAD\t: %s\n", pad)
		fmt.Fprintf(tw, "Format\t: %s\n", i.format)
		fmt.Fprintf(tw, "%s\n", strings.Repeat("-", 36))
		fmt.Fprintf(tw, "Decoded PIN\t: %s\n", pin)
		tw.Write([]byte("\n"))
	}

	return pin, nil
}
