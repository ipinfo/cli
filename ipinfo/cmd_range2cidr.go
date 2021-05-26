package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsRange2CIDR = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpRange2CIDR() {
	fmt.Printf(
		`Usage: %s range2cidr [<opts>] <ip-range | filepath>

Description:
  Accepts IP ranges and file paths to files containing IP ranges, converting
  them all to CIDRs (and multiple CIDRs if required).

  If a file is input, it is assumed that the IP range to convert is the first
  entry of each line. Other data is allowed and copied transparently.

  If multiple CIDRs are needed to represent an IP range on a line with other
  data, the data is copied per CIDR required. For example:

    in[0]: "1.1.1.0,1.1.1.2,other-data"
    out[0]: "1.1.1.0/31,other-data"
    out[1]: "1.1.1.2/32,other-data"

  IP ranges can of the form "<start><sep><end>" where "<sep>" can be "," or
  "-", and "<start>" and "<end>" can be any 2 IPs; order does not matter, but
  the resulting CIDRs are printed in the order they cover the range.

Examples:
  # Get all CIDRs for range 1.1.1.0-1.1.1.2.
  $ %[1]s range2cidr 1.1.1.0-1.1.1.2
  $ %[1]s range2cidr 1.1.1.0,1.1.1.2

  # Convert all range entries to CIDRs in 2 files.
  $ %[1]s range2cidr /path/to/file1.txt /path/to/file2.txt

  # Convert all range entries to CIDRs from stdin.
  $ cat /path/to/file1.txt | %[1]s range2cidr

  # Convert all range entries to CIDRs from stdin and a file.
  $ cat /path/to/file1.txt | %[1]s range2cidr /path/to/file2.txt

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdRange2CIDR() (err error) {
	f := lib.CmdRange2CIDRFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdRange2CIDR(f, pflag.Args()[1:], printHelpRange2CIDR)
}
