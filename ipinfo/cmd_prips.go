package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/complete/v3"
	"github.com/ipinfo/complete/v3/predict"
	"github.com/spf13/pflag"
)

var completionsPrips = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpPrips() {
	fmt.Printf(
		`Usage: %s prips [<opts>] <ip-range | cidr>

Description:
  Accepts CIDRs (e.g. 8.8.8.0/24) and IP ranges (e.g. 8.8.8.0-8.8.8.255).

  # List all IPs in a CIDR.
  $ %[1]s prips 8.8.8.0/24

  # List all IPs in multiple CIDRs.
  $ %[1]s prips 8.8.8.0/24 8.8.2.0/24 8.8.1.0/24

  # List all IPs in an IP range.
  $ %[1]s prips 8.8.8.0-8.8.8.255

  # List all IPs in multiple CIDRs and IP ranges.
  $ %[1]s prips 1.1.1.0/30 8.8.8.0-8.8.8.255 2.2.2.0/30 7.7.7.0,7.7.7.10

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

	return lib.OutputIPsFrom(args, true, true)
}
