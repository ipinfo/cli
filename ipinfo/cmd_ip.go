package main

import (
	"net"

	"github.com/urfave/cli/v2"
)

func cmdIP(c *cli.Context) error {
	ipStr := c.String("ip")
	ip := net.ParseIP(ipStr)

	data, err := ii.GetIPInfo(ip)
	if err != nil {
		return err
	}

	if c.Bool("json") {
		return outputJSON(data)
	}

	outputFriendlyCore(data)

	return nil
}
