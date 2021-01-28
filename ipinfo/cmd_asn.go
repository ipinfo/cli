package main

import (
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

func cmdAsn(c *cli.Context) error {
	args := c.Args()
	asn := args.First()

	// check length.
	if len(asn) < 3 {
		return errNotASN
	}

	// uppercase if necessary.
	asn = strings.ToUpper(asn)

	// ensure "AS" prefix.
	if !strings.HasPrefix(asn, "AS") {
		return errNotASN
	}

	// ensure number suffix.
	asnNumStr := asn[2:]
	if _, err := strconv.Atoi(asnNumStr); err != nil {
		return errNotASN
	}

	// TODO

	return nil
}
