package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolIsUnspecified = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-q":      predict.Nothing,
		"--quiet": predict.Nothing,
	},
}

func printHelpToolIsUnspecified() {
	fmt.Printf(`
%s tool is_unspecified [<opts>] <ip | filepath>

Description: Checks the provided address is an Unspecified Address
Inputs can be any IPv4 or IPv6 Address or filepath to a file

Examples
  # IPv4/IPv6 Address.
  $ %[1]s tool is_unspecified 0.0.0.0 | ::
  $ %[1]s tool is_unspecified 124.198.16.8 | fe80::2

  # Check for file.
  $ %[1]s tool is_unspecified /path/to/file.txt 

  # Check entries from stdin.
  $ cat /path/to/file1.txt | %[1]s tool is_unspecified

  Options:
  --help, -h
    show help.
  --quiet, -q
    quiet mode; suppress additional output.
`, progBase)
}

func cmdToolIsUnspecified() (err error) {
	f := lib.CmdToolIsUnspecifiedFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIsUnspecified(f, pflag.Args()[2:], printHelpToolIsUnspecified)
}
