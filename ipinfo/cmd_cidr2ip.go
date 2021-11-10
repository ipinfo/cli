package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsCIDR2IP = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpCIDR2IP() {
	fmt.Printf(
		`Usage: %s cidr2ip [<opts>] <cidrs | filepath>

Description:
  Accepts CIDRs and file paths to files containing CIDRs, converting
  them all to individual IPs within those ranges.

  $ %[1]s cidr2ip 8.8.8.0/24

  

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdCIDR2IP() error {
	f := lib.CmdCIDR2IPFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdCIDR2IP(f, pflag.Args()[1:], printHelpCIDR2IP)
}
