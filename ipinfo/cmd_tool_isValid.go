package main

import (
	"fmt"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolIsValid = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

// printHelpToolIsValid prints the help message for the "is_valid" command.
func printHelpToolIsValid() {
	fmt.Printf(
		`Usage: %s tool is_valid <ip>

Description:
  Reports whether an IP is valid.

Examples:
  %[1]s is_valid "190.87.89.1"
  %[1]s is_valid "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
  %[1]s is_valid "::"
  %[1]s is_valid "0"
  %[1]s is_valid ""

Options:
  General:
    --help, -h
      show help.
`, progBase)
}

// cmdToolIsValid is the handler for the "is_valid" command.
func cmdToolIsValid() error {
	f := lib.CmdToolIsValidFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIsValid(f, pflag.Args()[2:], printHelpToolIsValid)
}
