package main

import (
	"github.com/urfave/cli/v2"
)

func cmdASN(c *cli.Context) error {
	asn := c.String("asn")

	data, err := ii.GetASNDetails(asn)
	if err != nil {
		return err
	}

	return outputJSON(data)
}
