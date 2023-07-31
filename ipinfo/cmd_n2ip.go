package main

import (
	"fmt"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

// cmdN2IP is the handler for the "n2ip" command.
var completionsN2IP = &complete.Command{
	Flags: map[string]complete.Predictor{
		"--nocolor": predict.Nothing,
		"-h":        predict.Nothing,
		"--help":    predict.Nothing,
		"-6":        predict.Set(predictReadFmts),
		"--ipv6":    predict.Set(predictReadFmts),
	},
}

// printHelpN2IP prints the help message for the "n2ip" command.
func printHelpN2IP() {
	fmt.Printf(
		`Usage: %s n2ip [<opts>] <expr>

Example:
  %[1]s n2ip "2*2828-1"
  %[1]s n2ip "190.87.89.1*2"
  %[1]s n2ip "2001:0db8:85a3:0000:0000:8a2e:0370:7334*6"

Options:
  General:
    --nocolor
      disable colored output.
    --help, -h
      show help.
    --ipv6, -6
      force conversion to IPv6 address
`, progBase)
}

// cmdN2IP is the handler for the "n2ip" command.
func cmdN2IP() error {
	f := lib.CmdN2IPFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdN2IP(f, pflag.Args()[1:], printHelpN2IP)
}
