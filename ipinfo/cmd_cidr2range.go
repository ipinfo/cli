package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsCIDR2Range = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpCIDR2Range() {
	fmt.Printf(
		`Usage: %s cidr2range [<opts>] <cidr | filepath>

Description:
  Accepts CIDRs and file paths to files containing CIDRs, converting them all
  to IP ranges.

  If a file is input, it is assumed that the CIDR to convert is the first entry
  of each line. Other data is allowed and copied transparently.

  The IP range output is of the form "<start>-<end>" where "<start>" comes
  before or is equal to "<end>" in numeric value.

Examples:
  # Get the range for CIDR 1.1.1.0/30.
  $ %[1]s cidr2range 1.1.1.0/30

  # Convert CIDR entries to IP ranges in 2 files.
  $ %[1]s cidr2range /path/to/file1.txt /path/to/file2.txt

  # Convert CIDR entries to IP ranges from stdin.
  $ cat /path/to/file1.txt | %[1]s cidr2range

  # Convert CIDR entries to IP ranges from stdin and a file.
  $ cat /path/to/file1.txt | %[1]s cidr2range /path/to/file2.txt

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdCIDR2Range() (err error) {
	f := lib.CmdCIDR2RangeFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdCIDR2Range(f, pflag.Args()[1:], printHelpCIDR2Range)
}
