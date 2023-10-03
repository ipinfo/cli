package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolPrefixAddr = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpToolPrefixAddr() {
	fmt.Printf(
		`Usage: %s tool prefix addr <cidr>

Description:
	returns prefix's IP address.

Examples:
  # CIDR Valid Examples.
  $ %[1]s tool prefix addr 192.168.0.0/16
  $ %[1]s tool prefix addr 10.0.0.0/8
  $ %[1]s tool prefix addr 2001:0db8:1234::/48
  $ %[1]s tool prefix addr 2606:2800:220:1::/64

  # CIDR Invalid Examples.
  $ %[1]s tool prefix addr 192.168.0.0/40
  $ %[1]s tool prefix addr 2001:0db8:1234::/129

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdToolPrefixAddr() (err error) {
	f := lib.CmdToolPrefixAddrFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolPrefixAddr(f, pflag.Args()[3:], printHelpToolPrefixAddr)
}
