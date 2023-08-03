package main

import (
	"fmt"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolIP2n = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

// printHelpToolIp2n prints the help message for the "ip2n" command.
func printHelpToolIp2n() {
	fmt.Printf(
		`Usage: %s tool ip2n <ip>

Description:
  Converts an IPv4 or IPv6 address to its decimal representation.

Examples:
  %[1]s ip2n "190.87.89.1"
  %[1]s ip2n "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
  %[1]s ip2n "2001:0db8:85a3::8a2e:0370:7334"
  %[1]s ip2n "::7334"
  %[1]s ip2n "7334::"

Options:
  General:
    --help, -h
      show help.
`, progBase)
}

// cmdToolIP2n is the handler for the "ip2n" command.
func cmdToolIP2n() error {
	f := lib.CmdToolIP2nFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIP2n(f, pflag.Args()[2:], printHelpToolIp2n)
}
