package pinblock

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/moov-io/pinblock/formats"
)

const (
	DescribeEncode = "encode"
	DescribeDecode = "decode"
)

func Describe(iso formats.ISO, w io.Writer, method, pin, account string) error {

	if method == DescribeEncode {
		return describeEncode(iso, w, pin, account)
	} else if method == DescribeDecode {
		return describeDecode(iso, w, pin, account)
	}

	return fmt.Errorf("describe does not support a %s format", method)
}

func describeEncode(iso formats.ISO, w io.Writer, pin, account string) error {

	result, err := iso.Encode(pin, account)
	if err != nil {
		return fmt.Errorf("pin block encode operation failed")
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "PIN block encode operation finished\n")
	fmt.Fprintf(tw, "%s\n", strings.Repeat("*", 36))
	fmt.Fprintf(tw, "PAN\t: %s\n", account)
	fmt.Fprintf(tw, "PIN\t: %s\n", pin)
	if pad, _ := iso.Padding(pin); len(pad) == 0 {
		fmt.Fprintf(tw, "PAD\t: N/A\n")
	} else {
		fmt.Fprintf(tw, "PAD\t: %s\n", pad)
	}
	fmt.Fprintf(tw, "Format\t: %s\n", iso.Format())
	fmt.Fprintf(tw, "%s\n", strings.Repeat("-", 36))
	fmt.Fprintf(tw, "Clear PIN block\t: %s\n\n", result)

	tw.Flush()

	return nil
}

func describeDecode(iso formats.ISO, w io.Writer, pinBlock, account string) error {

	result, err := iso.Decode(pinBlock, account)
	if err != nil {
		return fmt.Errorf("pin block decode operation failed")
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "PIN block decode operation finished\n")
	fmt.Fprintf(tw, "%s\n", strings.Repeat("*", 36))
	fmt.Fprintf(tw, "PIN block\t: %s\n", pinBlock)
	fmt.Fprintf(tw, "PAN block\t: %s\n", account)
	if pad, _ := iso.Padding(result); len(pad) == 0 {
		fmt.Fprintf(tw, "PAD\t: N/A\n")
	} else {
		fmt.Fprintf(tw, "PAD\t: %s\n", pad)
	}
	fmt.Fprintf(tw, "Format\t: %s\n", iso.Format())
	fmt.Fprintf(tw, "%s\n", strings.Repeat("-", 36))
	fmt.Fprintf(tw, "Decoded PIN\t: %s\n\n", result)

	tw.Flush()

	return nil
}
