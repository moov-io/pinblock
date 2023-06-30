package formats

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type oemObject struct {
	format      string
	debugWriter io.Writer
}

// Padding returns padding pattern
func (i *oemObject) padding(pin string) (string, error) {
	if len(pin) < 4 {
		return "", fmt.Errorf("pin length must be between 4 and 12 digits")
	}

	length := 16 - len(pin)
	if length < 1 {
		return "", nil
	}

	hexLetters := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'A', 'B', 'C', 'D', 'E', 'F'}
	var exclusiveLetter byte

	for _, l := range hexLetters {
		if !strings.Contains(pin, string(l)) {
			exclusiveLetter = l
			break
		}
	}

	filler := string(exclusiveLetter)
	return strings.Repeat(filler, length), nil
}

// SetDebugWriter will set writer for getting output message of encoding and decoding logic
func (i *oemObject) SetDebugWriter(writer io.Writer) {
	i.debugWriter = tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
}

// Encode returns the OEM-1 PIN block for the given PIN
func (i *oemObject) Encode(pin string) (string, error) {
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

func (i *oemObject) Decode(pinBlock string) (string, error) {
	if len(pinBlock) != 16 {
		return "", fmt.Errorf("pin block must be 16 characters")
	}

	// getting pin length
	pinLength := len(strings.ReplaceAll(pinBlock, pinBlock[len(pinBlock)-1:], ""))

	// write decode information
	if i.debugWriter != nil {
		tw := i.debugWriter
		fmt.Fprintf(tw, "PIN block decode operation finished\n")
		fmt.Fprintf(tw, "%s\n", strings.Repeat("*", 36))
		fmt.Fprintf(tw, "Formatted PIN block\t: %s\n", strings.ToUpper(pinBlock))
		var pad string
		if len(pinBlock) > pinLength {
			pad = pinBlock[pinLength:]
		}
		if pad == "" {
			fmt.Fprintf(tw, "PAD\t: N/A\n")
		} else {
			fmt.Fprintf(tw, "PAD\t: %s\n", pad)
		}
		fmt.Fprintf(tw, "Format\t: %s\n", i.format)
		fmt.Fprintf(tw, "%s\n", strings.Repeat("-", 36))
		fmt.Fprintf(tw, "Decoded PIN\t: %s\n", pinBlock[:pinLength])
		tw.Write([]byte("\n"))
	}

	return pinBlock[:pinLength], nil
}
