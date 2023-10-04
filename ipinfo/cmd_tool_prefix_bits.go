package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolPrefixBits = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpToolPrefixBits() {
	fmt.Printf(
		`Usage: %s tool prefix bits <cidr>

Description:
  Returns the length of a prefix and reports -1 if invalid.

Examples:
  # CIDR Valid Examples.
  $ %[1]s tool prefix bits 192.168.0.0/16
  $ %[1]s tool prefix bits 10.0.0.0/8
  $ %[1]s tool prefix bits 2001:0db8:1234::/48
  $ %[1]s tool prefix bits 2606:2800:220:1::/64

  # CIDR Invalid Examples.
  $ %[1]s tool prefix bits 192.168.0.0/40
  $ %[1]s tool prefix bits 2001:0db8:1234::/129

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdToolPrefixBits() (err error) {
	f := lib.CmdToolPrefixBitsFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolPrefixBits(f, pflag.Args()[3:], printHelpToolPrefixBits)
}
