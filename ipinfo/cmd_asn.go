package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func cmdASN(c *cli.Context) error {
	asn := c.String("asn")

	data, err := ii.GetASNDetails(asn)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", data)

	return nil
}
