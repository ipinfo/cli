package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolIsV6 = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-q":      predict.Nothing,
		"--quiet": predict.Nothing,
	},
}

func printHelpToolIsV6() {
	fmt.Printf(
		`Usage: %s tool is_v6 [<opts>] <cidr | ip | ip-range | filepath>

Description:
  Checks if the input is an IPv6 address.
  Inputs can be IPs, IP ranges, CIDRs, or filepath to a file

Examples:
  # Check CIDR.
  $ %[1]s tool is_v6 2001:db8::/32

  # Check IP range.
  $ %[1]s tool is_v6 2001:db8::1-2001:db8::10

  # Check for file.
  $ %[1]s tool is_v6 /path/to/file1.txt 

  # Check entries from stdin.
  $ cat /path/to/file1.txt | %[1]s tool is_v6

Options:
  --help, -h
    show help.
  --quiet, -q
    quiet mode; suppress additional output.
`, progBase)
}

func cmdToolIsV6() (err error) {
	f := lib.CmdToolIsV6Flags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIsV6(f, pflag.Args()[2:], printHelpToolIsV6)
}
