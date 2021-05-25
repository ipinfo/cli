package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/spf13/pflag"
)

func printHelpRange2CIDR() {
	fmt.Printf(
		`Usage: %s range2cidr [<opts>]

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdRange2CIDR() (err error) {
	f := lib.CmdRange2CIDRFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdRange2CIDR(f, pflag.Args()[1:], printHelpRange2CIDR)
}
