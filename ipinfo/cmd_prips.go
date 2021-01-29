package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func cmdPrips(c *cli.Context) error {
	args := c.Args().Slice()

	// require args.
	if len(args) == 0 {
		return cli.ShowCommandHelp(c, "prips")
	}

	// ensure we only have CIDRs or IPs, but not both.
	hasCIDR := false
	for _, arg := range args {
		if isCIDR(arg) {
			hasCIDR = true
		} else if !isIP(arg) {
			return errNotIP
		} else if hasCIDR {
			return errCannotMixCIDRAndIPs
		}
	}

	// output CIDRs if that's what we have.
	if hasCIDR {
		for _, arg := range args {
			ips, err := cidrToIPs(arg)
			if err != nil {
				return err
			}

			for _, ip := range ips {
				fmt.Println(ip)
			}
		}
		return nil
	}

	// IP range input requires 2 IPs.
	if len(args) != 2 {
		return errIPRangeRequiresTwoIPs
	}

	// now we definitely have 2 IPs only.
	ips, err := ipRangeToIPs(args[0], args[1])
	if err != nil {
		return err
	}

	for _, ip := range ips {
		fmt.Println(ip)
	}

	return nil
}
