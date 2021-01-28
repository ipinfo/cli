package main

import (
	"fmt"
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

	fmt.Printf("%v\n", data)

	return nil
}
