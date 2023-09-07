package main

import (
	"fmt"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

// cmdToolUnmap is the handler for the "unmap" command.
var completionsToolUnmap = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

// printHelpToolUnmap prints the help message for the "unmap" command.
func printHelpToolUnmap() {
	fmt.Printf(
		`Usage: %s tool unmap [<opts>] <ip>

Description:
  Unmap returns an IP with any IPv4-mapped IPv6 address prefix removed.

  That is, if the IP is an IPv6 address wrapping an IPv4 address, it returns the
  wrapped IPv4 address. Otherwise it returns the IP unmodified.

Examples:
  %[1]s tool unmap "::ffff:8.8.8.8"
  %[1]s tool unmap "192.180.32.1"
  %[1]s tool unmap "::ffff:192.168.1.1"

Options:
  General:
    --help, -h
      show help.
`, progBase)
}

// cmdToolUnmap is the handler for the "unmap" command.
func cmdToolUnmap() error {
	f := lib.CmdToolUnmapFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolUnmap(f, pflag.Args()[2:], printHelpToolUnmap)
}
