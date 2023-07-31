package main

import (
	"fmt"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsIP2n = &complete.Command{
	Flags: map[string]complete.Predictor{
		"--nocolor": predict.Nothing,
		"-h":        predict.Nothing,
		"--help":    predict.Nothing,
	},
}

// printHelpIp2n prints the help message for the "ip2n" command.
func printHelpIp2n() {
	fmt.Printf(
		`Usage: %s ip2n <ip>

Example:
  %[1]s ip2n "190.87.89.1"
  %[1]s ip2n "2001:0db8:85a3:0000:0000:8a2e:0370:7334
  %[1]s ip2n "2001:0db8:85a3::8a2e:0370:7334
  %[1]s ip2n "::7334
  %[1]s ip2n "7334::""
	

Options:
  General:
    --nocolor
      disable colored output.
    --help, -h
      show help.
`, progBase)
}

// cmdIP2n is the handler for the "ip2n" command.
func cmdIP2n() error {
	f := lib.CmdIP2nFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdIP2n(f, pflag.Args()[1:], printHelpIp2n)
}
