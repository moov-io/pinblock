package formats

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type eciObject struct {
	format      string
	version     string
	debugWriter io.Writer
}

func (i *eciObject) getVersion() string {
	if i.version == eci2Version || i.version == eci3Version || i.version == visa2Version {
		return i.version
	}
	return eci2Version
}

// Padding returns padding pattern
func (i *eciObject) padding(pin string) (string, error) {
	if len(pin) < 4 {
		return "", fmt.Errorf("pin length must be between 4 and 12 digits")
	}

	length := 16 - len(pin)
	t := hexLetters
	if i.getVersion() == visa2Version {
		t = hexDigits
	}
	return randomLetters(length, t)
}

// SetDebugWriter will set writer for getting output message of encoding and decoding logic
func (i *eciObject) SetDebugWriter(writer io.Writer) {
	i.debugWriter = tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
}

// Encode returns the OEM-1 PIN block for the given PIN
func (i *eciObject) Encode(pin, account string) (string, error) {
	isTruncated := false

	if len(pin) < 4 || len(pin) > 12 {
		return "", fmt.Errorf("pin length must be between 4 and 12 digits")
	}

	var rawPin string
	if i.getVersion() == eci2Version {
		if len(pin) > 4 {
			pin = pin[:4]
			isTruncated = true
		}
		rawPin = pin
	} else {
		if len(pin) > 6 {
			pin = pin[:6]
			isTruncated = true
		}
		rawPin = pin
		pin = fmt.Sprintf("%d%s", len(pin), pin)
		pin = pin + strings.Repeat("0", 7-len(pin))
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
		fmt.Fprintf(tw, "PIN\t: %s\n", rawPin)
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

func (i *eciObject) Decode(pinBlock, account string) (string, error) {
	if len(pinBlock) != 16 {
		return "", fmt.Errorf("pin block must be 16 characters")
	}

	pinLength := 4
	if i.getVersion() == eci2Version {
		pinLength = 4
	} else {
		_, err := fmt.Sscanf(pinBlock, "%01d%s", &pinLength, &pinBlock)
		if err != nil {
			return "", fmt.Errorf("unable to parse pin block")
		}
		if pinLength > 6 || pinLength < 4 {
			return "", fmt.Errorf("unable to parse pin block")
		}
	}

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
