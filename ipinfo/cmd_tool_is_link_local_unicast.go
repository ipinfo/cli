package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolIsLinkLocalUnicast = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-q":      predict.Nothing,
		"--quiet": predict.Nothing,
	},
}

func printHelpToolIsLinkLocalUnicast() {
	fmt.Printf(
		`Usage: %s tool is_link_local_unicast [<opts>] <cidr | ip | ip-range | filepath>

Description:
  Checks if the input is a link local unicast address.
  Inputs can be IPs, IP ranges, CIDRs, or filepath to a file

Examples:
  $ %[1]s tool is_link_local_unicast 169.254.0.0
  $ %[1]s tool is_link_local_unicast 127.0.0.0
  $ %[1]s tool is_link_local_unicast fe80::1
  $ %[1]s tool is_link_local_unicast ::1

  # Check CIDR.
  $ %[1]s tool is_link_local_unicast 169.254.0.0/32
  $ %[1]s tool is_link_local_unicast 139.0.0.0/32
  $ %[1]s tool is_link_local_unicast fe80::1/64
  $ %[1]s tool is_link_local_unicast ::1/64

  # Check IP range.
  $ %[1]s tool is_link_local_unicast 169.254.0.0-169.254.255.1
  $ %[1]s tool is_link_local_unicast 240.0.0.0-240.255.255.1
  $ %[1]s tool is_link_local_unicast fe80::1-feb0::1
  $ %[1]s tool is_link_local_unicast ::1-::ffff

  # Check for file.
  $ %[1]s tool is_link_local_unicast /path/to/file.txt 

  # Check entries from stdin.
  $ cat /path/to/file1.txt | %[1]s tool is_link_local_unicast

Options:
  --help, -h
    show help.
  --quiet, -q
    quiet mode; suppress additional output.
`, progBase)
}

func cmdToolIsLinkLocalUnicast() (err error) {
	f := lib.CmdToolIsLinkLocalUnicastFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIsLinkLocalUnicast(f, pflag.Args()[2:], printHelpToolIsLinkLocalUnicast)
}
