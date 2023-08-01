package main

import (
	"fmt"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

// cmdToolN2IP is the handler for the "n2ip" command.
var completionsToolN2IP = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
		"-6":     predict.Set(predictReadFmts),
		"--ipv6": predict.Set(predictReadFmts),
	},
}

// printHelpToolN2IP prints the help message for the "n2ip" command.
func printHelpToolN2IP() {
	fmt.Printf(
		`Usage: %s n2ip tool [<opts>] <number>

Description:
  Converts a given numeric representation to its corresponding IPv4 or IPv6 address, and can also evaluate a mathematical expression for conversion.

Examples:
  %[1]s n2ip "4294967295 + 87"
  %[1]s n2ip "4294967295" --ipv6
  %[1]s n2ip -6 "201523715"
  %[1]s n2ip "51922968585348276285304963292200960"
  %[1]s n2ip "a:: - 4294967295"

Options:
  General:
    --help, -h
      show help.
    --ipv6, -6
      force conversion to IPv6 address
`, progBase)
}

// cmdToolN2IP is the handler for the "n2ip" command.
func cmdToolN2IP() error {
	f := lib.CmdToolN2IPFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolN2IP(f, pflag.Args()[2:], printHelpToolN2IP)
}
