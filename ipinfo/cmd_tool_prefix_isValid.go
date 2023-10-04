package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolPrefixIsValid = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpToolPrefixIsValid() {
	fmt.Printf(
		`Usage: %s tool prefix is_valid <cidr>

Description:
  Reports whether a prefix is valid.

Examples:
  # CIDR Valid Examples.
  $ %[1]s tool prefix is_valid 192.168.0.0/16
  $ %[1]s tool prefix is_valid 10.0.0.0/8
  $ %[1]s tool prefix is_valid 2001:0db8:1234::/48
  $ %[1]s tool prefix is_valid 2606:2800:220:1::/64

  # CIDR Invalid Examples.
  $ %[1]s tool prefix is_valid 192.168.0.0/40
  $ %[1]s tool prefix is_valid 2001:0db8:1234::/129

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdToolPrefixIsValid() (err error) {
	f := lib.CmdToolPrefixIsValidFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolPrefixIsValid(f, pflag.Args()[3:], printHelpToolPrefixIsValid)
}
