package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolPrev = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-q":      predict.Nothing,
		"--quiet": predict.Nothing,
	},
}

func printHelpToolPrev() {
	fmt.Printf(
		`Usage: %s tool prev [<opts>] <ip | filepath>

Description:
  Finds the previous IP address for given input efficiently 
  Input must be IPs.

Examples:
  # Find the previous IP for the given inputs 
  $ %[1]s tool prev 1.1.1.0

  # Find prev IP from stdin.
  $ cat /path/to/file1.txt | %[1]s tool prev 

Options:
  --help, -h
    show help.
  --quiet, -q
    quiet mode; suppress additional output.
`, progBase)
}

func cmdToolPrev() (err error) {
	f := lib.CmdToolPrevFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolPrev(f, pflag.Args()[2:], printHelpToolPrev)
}
