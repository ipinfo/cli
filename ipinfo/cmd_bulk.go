package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/urfave/cli/v2"
)

func cmdBulk(c *cli.Context) error {
	args := c.Args().Slice()

	// check for stdin, implied or explicit.
	if len(args) == 0 || (len(args) == 1 && args[0] == "-") {
		ips := inputIPsFromStdin()
		if len(ips) == 0 {
			fmt.Println("no input ips")
			return nil
		}

		data, err := ii.GetIPInfoBatch(ips, ipinfo.BatchReqOpts{})
		if err != nil {
			return err
		}

		if c.Bool("csv") {
			return outputCSVBatchCore(data)
		}

		return outputJSON(data)
	}

	// check for IP range.
	if isIP(args[0]) {
		if len(args) != 2 {
			return errIPRangeRequiresTwoIPs
		}
		if !isIP(args[1]) {
			return errNotIP
		}

		// TODO
		return nil
	}

	// check for all CIDRs.
	if isCIDR(args[0]) {
		for _, arg := range args[1:] {
			if !isCIDR(arg) {
				return errNotCIDR
			}
		}

		// TODO
		return nil
	}

	// check for all filepaths.
	if fileExists(args[0]) {
		for _, arg := range args[1:] {
			if !fileExists(arg) {
				return errNotFile
			}
		}

		// TODO
		return nil
	}

	return errInvalidInput
}
