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
		`Usage: %s tool lower [<opts>] <cidr>

Description:
  Calculates the lower IP address (start address of a network) for the given inputs.
  Input can be IP, IP range or CIDR.
  
  If input contains CIDRs, it calculates the lower IP address for each CIDR.
  If input contains IP ranges, it calculates the lower IP address for each Ip range.
  
Examples:
  # Calculate lower IP for IP, IP range and CIDR.
  $ %[1]s tool lower 192.168.1.0/24
  
  # Calculate lower IPs for IPs, IP ranges and CIDRs.
  $ %[1]s tool lower 192.168.1.0/24 10.0.0.0/16
  
  # Calculate lower IPs from stdin.
  $ cat /path/to/file.txt | %[1]s tool lower

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
