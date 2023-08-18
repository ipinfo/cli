package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolUpper = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-q":      predict.Nothing,
		"--quiet": predict.Nothing,
	},
}

func printHelpToolUpper() {
	fmt.Printf(
		`Usage: %s tool upper [<opts>] <cidr | ip | ip-range | filepath>

Description:
  Calculates the upper IP address (end address of a network) for the given inputs.
  Input can be a mixture of Ips, IP ranges or CIDRs.

Examples:
  # Finds upper IP for CIDR.
  $ %[1]s tool upper 192.168.1.0/24

  # Finds upper IP for IP range.
  $ %[1]s tool upper 1.1.1.0-1.1.1.244

  # Finds upper IPs from stdin.
  $ cat /path/to/file.txt | %[1]s tool upper

  # Find upper IPs from file.
  $ %[1]s tool upper /path/to/file1.txt 

Options:
  --help, -h
    Show help.
  --quiet, -q
    Quiet mode; suppress additional output.
`, progBase)
}

func cmdToolUpper() (err error) {
	f := lib.CmdToolUpperFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolUpper(f, pflag.Args()[2:], printHelpToolUpper)
}
