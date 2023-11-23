package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolIsMulticast = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-q":      predict.Nothing,
		"--quiet": predict.Nothing,
	},
}

func printHelpToolIsMultiCast() {
	fmt.Printf(
		`Usage: %s tool is_multicast [<opts>] <cidr | ip | ip-range | filepath>

Description:
  Checks if the input is a multicast address.
  Inputs can be IPs, IP ranges, CIDRs, or filepath to a file

Examples:
  $ %[1]s tool is_multicast 239.0.0.0
  $ %[1]s tool is_multicast ff00::
  $ %[1]s tool is_multicast 127.0.0.0
  $ %[1]s tool is_multicast ::1

  # Check CIDR.
  $ %[1]s tool is_multicast 239.0.0.0/32
  $ %[1]s tool is_multicast 139.0.0.0/32
  $ %[1]s tool is_multicast ff00::1/64
  $ %[1]s tool is_multicast ::1/64

  # Check IP range.
  $ %[1]s tool is_multicast 239.0.0.0-239.255.255.1
  $ %[1]s tool is_multicast 240.0.0.0-240.255.255.1
  $ %[1]s tool is_multicast ff00::1-ff00::ffff
  $ %[1]s tool is_multicast ::1-::ffff

  # Check for file.
  $ %[1]s tool is_multicast /path/to/file.txt 

  # Check entries from stdin.
  $ cat /path/to/file1.txt | %[1]s tool is_multicast

Options:
  --help, -h
    show help.
  --quiet, -q
    quiet mode; suppress additional output.
`, progBase)
}

func cmdToolIsMultiCast() (err error) {
	f := lib.CmdToolIsMulticastFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIsMulticast(f, pflag.Args()[2:], printHelpToolIsMultiCast)
}
