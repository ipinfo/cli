package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolAggregate = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-q":      predict.Nothing,
		"--quiet": predict.Nothing,
	},
}

func printHelpToolAggregate() {
	fmt.Printf(
		`Usage: %s aggregate [<opts>] <cidr | filepath>

Description:
  Accepts IPs, IP ranges, and CIDRs, aggregating them efficiently.
  Input can be IPs, IP ranges, CIDRs, and/or a filepath containing
  any of these. Works for both IPv4 and IPv6.

  If input contains single IPs, it tries to merge them into the input CIDRs,
  otherwise they are printed to the output as it is.

  IP range can be of format <start-ip><SEP><end-ip>, where <SEP> can either
  be a ',' or a '-'.

Examples:
  # Aggregate two CIDRs.
  $ %[1]s aggregate 1.1.1.0/30 1.1.1.0/28

  # Aggregate enteries from 2 files.
  $ %[1]s aggregate /path/to/file1.txt /path/to/file2.txt

  # Aggregate enteries from stdin.
  $ cat /path/to/file1.txt | %[1]s aggregate

  # Aggregate enteries from stdin and a file.
  $ cat /path/to/file1.txt | %[1]s aggregate /path/to/file2.txt

Options:
  --help, -h
    show help.
  --quiet, -q
	quiet mode; suppress additional output.
`, progBase)
}

func cmdToolAggregate() (err error) {
	f := lib.CmdToolAggregateFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolAggregate(f, pflag.Args()[2:], printHelpToolAggregate)
}
