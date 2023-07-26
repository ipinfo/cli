package main

import (
	"fmt"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsN2IP6 = &complete.Command{
	Flags: map[string]complete.Predictor{
		"--nocolor": predict.Nothing,
		"-h":        predict.Nothing,
		"--help":    predict.Nothing,
	},
}

// printHelpN2IP6 prints the help message for the "n2ip6" command.
func printHelpN2IP6() {
	fmt.Printf(
		`Usage: %s n2ip6 [<opts>] <expr>

Example:
  %[1]s n2ip6 "190.87.89.1"
  %[1]s n2ip6 "2001:0db8:85a3:0000:0000:8a2e:0370:7334
  %[1]s n2ip6 "2001:0db8:85a3::8a2e:0370:7334
  %[1]s n2ip6 "::7334
  %[1]s n2ip6 "7334::""
	

Options:
  General:
    --nocolor
      disable colored output.
    --help, -h
      show help.
`, progBase)
}

// cmdN2IP6 is the handler for the "n2ip6" command.
func cmdN2IP6() error {
	f := lib.CmdN2IP6Flags{}
	f.Init()
	pflag.Parse()
	if pflag.NArg() <= 1 && pflag.NFlag() == 0 {
		f.Help = true
	}

	return lib.CmdN2IP6(f, pflag.Args()[1:], printHelpN2IP6)
}
