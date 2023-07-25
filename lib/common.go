package lib

import (
	"bufio"
	"io"
	"strings"
)

type CIDRProcessor func(string) error

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
