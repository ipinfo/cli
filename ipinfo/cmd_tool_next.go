package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolNext = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-q":      predict.Nothing,
		"--quiet": predict.Nothing,
	},
}

func printHelpToolNext() {
	fmt.Printf(
		`Usage: %s tool next [<opts>] <ip | filepath>

Description:
  Finds the next IP address for the given input IP.
  Inputs must be IPs.

Examples:
  # Find the next IP for the given inputs 
  $ %[1]s tool next 1.1.1.0

  # Find next IP from stdin.
  $ cat /path/to/file1.txt | %[1]s tool next 

Options:
  --help, -h
    show help.
  --quiet, -q
    quiet mode; suppress additional output.
`, progBase)
}

func cmdToolNext() (err error) {
	f := lib.CmdToolNextFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolNext(f, pflag.Args()[2:], printHelpToolNext)
}
