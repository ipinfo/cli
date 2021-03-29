package main

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/ipinfo/cli/lib"
)

func printHelpPrips() {
	fmt.Printf(
		`Usage: %s prips [<opts>] <cidrs or ip-range>

Description:
  Accepts CIDRs (e.g. 8.8.8.0/24) or an IP range (e.g. 8.8.8.0 8.8.8.255).

  # List all IPs in a CIDR.
  $ %[1]s prips 8.8.8.0/24

  # List all IPs in multiple CIDRs.
  $ %[1]s prips 8.8.8.0/24 8.8.2.0/24 8.8.1.0/24

  # List all IPs in an IP range.
  $ %[1]s prips 8.8.8.0 8.8.8.255

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdPrips() error {
	var fHelp bool

	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpPrips()
		return nil
	}

	args := pflag.Args()[1:]

	// require args.
	if len(args) == 0 {
		printHelpPrips()
		return nil
	}

	// ensure we only have CIDRs or IPs, but not both.
	hasCIDR := false
	for _, arg := range args {
		if lib.IsCIDR(arg) {
			hasCIDR = true
		} else if !lib.IsIP(arg) {
			return lib.ErrNotIP
		} else if hasCIDR {
			return lib.ErrCannotMixCIDRAndIPs
		}
	}

	// output CIDRs if that's what we have.
	if hasCIDR {
		for _, arg := range args {
			if err := lib.OutputIPsFromCIDR(arg); err != nil {
				return err
			}
		}
		return nil
	}

	// IP range input requires 2 IPs.
	if len(args) != 2 {
		return lib.ErrIPRangeRequiresTwoIPs
	}

	// now we definitely have 2 IPs only.
	if err := lib.OutputIPsFromRange(args[0], args[1]); err != nil {
		return err
	}

	return nil
}
