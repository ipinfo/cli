package main

import (
	"net"

	"github.com/urfave/cli/v2"
)

func cmdIP(c *cli.Context) error {
	var ip net.IP

	args := c.Args()
	ipStr := args.First()
	if ip = net.ParseIP(ipStr); ip == nil {
		return errNotIP
	}

	// TODO

	return nil
}
