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
  Accepts IPs, IP ranges, and CIDRs, aggregating them efficiently.
  Input can be IPs, IP ranges, CIDRs, and/or filepath to a file
  containing any of these. Works for both IPv4 and IPv6.

  If input contains single IPs, it tries to merge them into the input CIDRs,
  otherwise they are printed to the output as they are.

  IP range can be of format <start-ip><SEP><end-ip>, where <SEP> can either
  be a ',' or a '-'.

Examples:
  # Lower two CIDRs.
  $ %[1]s tool lower 1.1.1.0/30 1.1.1.0/28

  # Lower IP range and CIDR.
  $ %[1]s tool lower 1.1.1.0-1.1.1.244 1.1.1.0/28

  # Lower enteries from 2 files.
  $ %[1]s tool lower /path/to/file1.txt /path/to/file2.txt

  # Lower enteries from stdin.
  $ cat /path/to/file1.txt | %[1]s tool lower

  # Lower enteries from stdin and a file.
  $ cat /path/to/file1.txt | %[1]s tool lower /path/to/file2.txt

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
