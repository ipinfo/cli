package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsLogout = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpLogout() {
	fmt.Printf(
		`Usage: %s logout [<opts>]

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdLogout() error {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpLogout()
		return nil
	}

	// delete but don't return an error; just log it.
	if gConfig.Token == "" {
		fmt.Println("not logged in")
		return nil
	}

	gConfig.Token = ""
	if err := SaveConfig(gConfig); err != nil {
		return err
	}

	fmt.Println("logged out")

	return nil
}
