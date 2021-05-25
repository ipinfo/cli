package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/spf13/pflag"
)

func printHelpCIDR2Range() {
	fmt.Printf(
		`Usage: %s cidr2range [<opts>]

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdCIDR2Range() (err error) {
	f := lib.CmdCIDR2RangeFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdCIDR2Range(f, pflag.Args()[1:], printHelpCIDR2Range)
}
