package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolLower = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-q":      predict.Nothing,
		"--quiet": predict.Nothing,
	},
}

func printHelpToolLower() {
	fmt.Printf(
		`Usage: %s tool lower [<opts>] <cidr | ip | ip-range | filepath>

Description:
  Calculates the lower IP address (start address of a network) for the given inputs.
  Input can be a mixture of IPs, IP ranges or CIDRs.
  
Examples:
  # Finds lower IP for CIDR.
  $ %[1]s tool lower 192.168.1.0/24

  # Finds lower IP for IP range.
  $ %[1]s tool lower 1.1.1.0-1.1.1.244

  # Finds lower IPs from stdin.
  $ cat /path/to/file.txt | %[1]s tool lower

  # Find lower IPs from file.
  $ %[1]s tool lower /path/to/file1.txt 

Options:
  --help, -h
    show help.
  --quiet, -q
    quiet mode; suppress additional output.
`, progBase)
}

func cmdToolLower() (err error) {
	f := lib.CmdToolLowerFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolLower(f, pflag.Args()[2:], printHelpToolLower)
}
