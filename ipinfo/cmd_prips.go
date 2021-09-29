package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
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
		`Usage: %s prips [<opts>] <ip | ip-range | cidr | file>

Description:
  Accepts CIDRs (e.g. 8.8.8.0/24) and IP ranges (e.g. 8.8.8.0-8.8.8.255).

Examples:
  # List all IPs in a CIDR.
  $ %[1]s prips 8.8.8.0/24

  # List all IPs in multiple CIDRs.
  $ %[1]s prips 8.8.8.0/24 8.8.2.0/24 8.8.1.0/24

  # List all IPs in an IP range.
  $ %[1]s prips 8.8.8.0-8.8.8.255

  # List all IPs in multiple CIDRs and IP ranges.
  $ %[1]s prips 1.1.1.0/30 8.8.8.0-8.8.8.255 2.2.2.0/30 7.7.7.0,7.7.7.10

  # List all IPs from stdin input (newline-separated).
  $ echo '1.1.1.0/30\n8.8.8.0-8.8.8.255\n7.7.7.0,7.7.7.10' | %[1]s prips

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdPrips() error {
	f := lib.CmdPripsFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdPrips(f, pflag.Args()[1:], printHelpPrips)
}
