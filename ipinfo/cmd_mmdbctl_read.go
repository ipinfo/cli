package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	mmdbLib "github.com/ipinfo/mmdbctl/lib"
	"github.com/spf13/pflag"
)

var predictReadFmts = []string{
	"json",
	"json-compact",
	"json-pretty",
	"tsv",
	"csv",
}

var completionsRead = &complete.Command{
	Flags: map[string]complete.Predictor{
		"--nocolor": predict.Nothing,
		"-h":        predict.Nothing,
		"--help":    predict.Nothing,
		"-f":        predict.Set(predictReadFmts),
		"--format":  predict.Set(predictReadFmts),
	},
}

func printHelpRead() {
	fmt.Printf(
		`Usage: %s mmdbctl read [<opts>] <ip | ip-range | cidr | filepath> <mmdb>

Options:
  General:
    --nocolor
      disable colored output.
    --help, -h
      show help.

  Format:
    -f <format>, --format <format>
      the output format.
      can be "json", "json-compact", "json-pretty", "tsv" or "csv".
      note that "json" is short for "json-compact".
      default: json.
`, progBase)
}

func cmdRead() error {
	f := mmdbLib.CmdReadFlags{}
	f.Init()
	pflag.Parse()
	if pflag.NArg() <= 2 {
		f.Help = true
	}

	return mmdbLib.CmdRead(f, pflag.Args()[2:], printHelpRead)
}
