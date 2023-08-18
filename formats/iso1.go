package formats

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type iso1Object struct {
	Filler string

	version     string
	format      string
	debugWriter io.Writer
}

func (i *iso1Object) getVersion() string {
	if i.version == iso1Version || i.version == iso2Version {
		return i.version
	}
	return iso1Version
}

// Padding returns padding pattern
func (i *iso1Object) padding(pin string) (string, error) {
	if len(pin) < 4 {
		return "", fmt.Errorf("pin length must be between 4 and 12 digits")
	}

	length := 14 - len(pin)
	if length < 1 {
		return "", nil
	}

	if i.Filler == "" {
		// ISO2
		//  fill is 0xF instead of random digits
		return randomLetters(length, hexDigits)
	}
	return strings.Repeat(i.Filler, length), nil
}

// SetDebugWriter will set writer for getting output message of encoding and decoding logic
func (i *iso1Object) SetDebugWriter(writer io.Writer) {
	i.debugWriter = tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
}

// Encode returns the ISO1 PIN block for the given PIN
//
//	The `ISO-1` PIN block format is equivalent to an `ECI-4` PIN block format
//	and is recommended for usage where no PAN data is available.
func (i *iso1Object) Encode(pin, account string) (string, error) {
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

	// A PIN that is longer than 12 digits is truncated on the right.
	// The first nibble (which identifies the block format) has the value 1.
	pinBlock := fmt.Sprintf("%s%X%s%s", i.getVersion(), len(pin), pin, pad)

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

func (i *iso1Object) Decode(pinBlock, account string) (string, error) {
	if len(pinBlock) != 16 {
		return "", fmt.Errorf("pin block must be 16 characters")
	}

	var pinLength int
	var decodedBlock string

	_, err := fmt.Sscanf(pinBlock, i.getVersion()+"%01X%s", &pinLength, &decodedBlock)
	if err != nil {
		return "", fmt.Errorf("unable to parse pin block")
	}

	if len(decodedBlock) < pinLength {
		return "", fmt.Errorf("parsed pin length is incorrect")
	}

	// write decode information
	if i.debugWriter != nil {
		tw := i.debugWriter
		fmt.Fprintf(tw, "PIN block decode operation finished\n")
		fmt.Fprintf(tw, "%s\n", strings.Repeat("*", 36))
		fmt.Fprintf(tw, "Formatted PIN block\t: %s\n", strings.ToUpper(decodedBlock))
		var pad string
		if len(decodedBlock) > pinLength {
			pad = decodedBlock[pinLength:]
		}
		if pad == "" {
			fmt.Fprintf(tw, "PAD\t: N/A\n")
		} else {
			fmt.Fprintf(tw, "PAD\t: %s\n", pad)
		}
		fmt.Fprintf(tw, "Format\t: %s\n", i.format)
		fmt.Fprintf(tw, "%s\n", strings.Repeat("-", 36))
		fmt.Fprintf(tw, "Decoded PIN\t: %s\n", decodedBlock[:pinLength])
		tw.Write([]byte("\n"))
	}

	return decodedBlock[:pinLength], nil
}
