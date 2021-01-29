package main

import (
	"github.com/urfave/cli/v2"
)

func cmdMyIP(c *cli.Context) error {
	data, err := ii.GetIPInfo(nil)
	if err != nil {
		return err
	}

	if c.Bool("json") {
		return outputJSON(data)
	}

	outputFriendlyCore(data)

	return nil
}
