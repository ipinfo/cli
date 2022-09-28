package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsSplitCIDR = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpSplitCIDR() {
	fmt.Printf(
		`Usage: %s splitcidr <cidr> <split>

Description:
  splits a larger CIDR into smaller CIDRs.

  $ %[1]s splitcidr 8.8.8.0/24 25

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdSplitCIDR() error {
	f := lib.CmdSplitCIDRFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdSplitCIDR(f, pflag.Args()[1:], printHelpSplitCIDR)
}
