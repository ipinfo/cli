package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolPrefixMasked = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpToolPrefixMasked() {
	fmt.Printf(
		`Usage: %s tool prefix masked <cidr>

Description:
  Returns canonical form of a prefix, masking off non-high bits, and returns the zero if invalid.

Examples:
  # CIDR Valid Examples.
  $ %[1]s tool prefix masked 192.168.0.0/16
  $ %[1]s tool prefix masked 10.0.0.0/8
  $ %[1]s tool prefix masked 2001:0db8:1234::/48
  $ %[1]s tool prefix masked 2606:2800:220:1::/64

  # CIDR Invalid Examples.
  $ %[1]s tool prefix masked 192.168.0.0/40
  $ %[1]s tool prefix masked 2001:0db8:1234::/129

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdToolPrefixMasked() (err error) {
	f := lib.CmdToolPrefixMaskedFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolPrefixMasked(f, pflag.Args()[3:], printHelpToolPrefixMasked)
}
