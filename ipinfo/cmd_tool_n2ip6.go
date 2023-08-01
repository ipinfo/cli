package main

import (
	"fmt"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolN2IP6 = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

// printHelpToolN2IP6 prints the help message for the "n2ip6" command.
func printHelpToolN2IP6() {
	fmt.Printf(
		`Usage: %s tool n2ip6 [<opts>] <number>

Description:
  Converts a given numeric representation to its corresponding IPv6 address, and can also evaluate a mathematical expression for conversion.

Examples:
  %[1]s n2ip "4294967295 + 87"
  %[1]s n2ip "4294967295"
  %[1]s n2ip "201523715"
  %[1]s n2ip "51922968585348276285304963292200960"
  %[1]s n2ip "a:: - 4294967295"

Options:
  General:
    --help, -h
      show help.
`, progBase)
}

// cmdToolN2IP6 is the handler for the "n2ip6" command.
func cmdToolN2IP6() error {
	f := lib.CmdToolN2IP6Flags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolN2IP6(f, pflag.Args()[2:], printHelpToolN2IP6)
}
