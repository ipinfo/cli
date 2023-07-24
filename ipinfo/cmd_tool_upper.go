// cmd_tool_upper.go

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
  Input can be IPs, IP ranges, CIDRs, and/or filepath to a file containing any of these.
  Works for both IPv4 and IPv6.

  If input contains CIDRs, it calculates the upper IP address for each CIDR.
  If input contains single IPs, it calculates the upper IP address for each IP.

  IP range can be of format <start-ip><SEP><end-ip>, where <SEP> can either
  be a ',' or a '-'.

Examples:
  # Calculate upper IP for a CIDR.
  $ %[1]s tool upper 192.168.1.0/24

  # Calculate upper IPs for CIDRs.
  $ %[1]s tool upper 192.168.1.0/24 10.0.0.0/16

  # Calculate upper IPs for a range and a CIDR.
  $ %[1]s tool upper 192.168.1.10-192.168.1.20 10.0.0.0/16

  # Calculate upper IPs from a file.
  $ %[1]s tool upper /path/to/file.txt

  # Calculate upper IPs from stdin.
  $ cat /path/to/file.txt | %[1]s tool upper

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
