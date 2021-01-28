package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func cmdMyIP(c *cli.Context) error {
	data, err := ii.GetIPInfo(nil)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", data)

	return nil
}
