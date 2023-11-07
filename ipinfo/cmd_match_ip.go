package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/spf13/pflag"
)

func printHelpMatchIP() {
	fmt.Printf(
		`Usage: %s matchip --list <file> --overlapping-with <file>

Options:
  General:
    --filter, -f
      file containing a list of IP, CIDR, and/or Ranges for filtering.
    --overlap-check, -o
      file containing a list of CIDR and/or Ranges to check for overlap.
    --help
      show help.
`, progBase)
}

func cmdMatchIP() error {
	f := lib.CmdMatchIPFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdMatchIP(f, pflag.Args()[1:], printHelpMatchIP)
}
