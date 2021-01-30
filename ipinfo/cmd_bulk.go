package main

import (
	"fmt"
	"net"

	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/urfave/cli/v2"
)

func cmdBulk(c *cli.Context) (err error) {
	var ips []net.IP

	args := c.Args().Slice()

	// check for stdin, implied or explicit.
	if len(args) == 0 || (len(args) == 1 && args[0] == "-") {
		ips = ipsFromStdin()

		goto lookup
	}

	// check for IP range.
	if isIP(args[0]) {
		if len(args) != 2 {
			return errIPRangeRequiresTwoIPs
		}
		if !isIP(args[1]) {
			return errNotIP
		}

		ips, err = ipsFromRange(args[0], args[1])
		if err != nil {
			return err
		}

		goto lookup
	}

	// check for all CIDRs.
	if isCIDR(args[0]) {
		for _, arg := range args[1:] {
			if !isCIDR(arg) {
				return errNotCIDR
			}
		}

		ips, err = ipsFromCIDRs(args)
		if err != nil {
			return err
		}

		goto lookup
	}

	// check for all filepaths.
	if fileExists(args[0]) {
		for _, arg := range args[1:] {
			if !fileExists(arg) {
				return errNotFile
			}
		}

		ips, err = ipsFromFiles(args)
		if err != nil {
			return err
		}

		goto lookup
	}

lookup:

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
