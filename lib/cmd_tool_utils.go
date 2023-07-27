package lib

import (
	"bufio"
	"io"
	"strings"
)

// Function type for the callback used in scanrdr
type CIDRProcessor func(string) error

// Function to read lines from the io.Reader and process each line as a CIDR using the provided callback function
func scanrdr(r io.Reader, processCIDR CIDRProcessor) error {
	buf := bufio.NewReader(r)
	for {
		d, err := buf.ReadString('\n')
		if err == io.EOF {
			if len(d) == 0 {
				break
			}
		} else if err != nil {
			return err
		}

		sepIdx := strings.IndexAny(d, "\n")
		if sepIdx == -1 {
			sepIdx = len(d)
		}

		cidrStr := d[:sepIdx]
		if err := processCIDR(cidrStr); err != nil {
			return err
		}
	}
	return nil
}
