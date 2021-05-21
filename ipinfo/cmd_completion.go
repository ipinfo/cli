package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/install"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsCompletion = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpCompletion() {
	fmt.Printf(
		`Usage: %s completion [<opts>]

Description:
  Install the code needed to allow auto-completion for various shells.

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdCompletion() error {
	var fHelp bool

	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpCompletion()
		return nil
	}

	return install.Install(progBase)
}
