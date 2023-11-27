package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolIsLinkLocalMulticast = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-q":      predict.Nothing,
		"--quiet": predict.Nothing,
	},
}

func printHelpToolIsLinkLocalMulticast() {
	fmt.Printf(
		`Usage: %s tool is_link_local_multicast [<opts>] <cidr | ip | ip-range | filepath>

Description:
  Checks if the input is a link local multicast address.
  Inputs can be IPs, IP ranges, CIDRs, or filepath to a file

Examples:
  $ %[1]s tool is_link_local_multicast 224.0.0.0
  $ %[1]s tool is_link_local_multicast 169.200.0.0
  $ %[1]s tool is_link_local_multicast ff02::2
  $ %[1]s tool is_link_local_multicast fe80::

  # Check CIDR.
  $ %[1]s tool is_link_local_multicast 224.0.0.0/32
  $ %[1]s tool is_link_local_multicast 169.200.0.0/32
  $ %[1]s tool is_link_local_multicast ff02::1/64
  $ %[1]s tool is_link_local_multicast fe80::1/64

  # Check IP range.
  $ %[1]s tool is_link_local_multicast 224.0.0.1-224.255.255.255
  $ %[1]s tool is_link_local_multicast 169.254.0.1-169.254.255.0
  $ %[1]s tool is_link_local_multicast ff02::1-ff02::ffff
  $ %[1]s tool is_link_local_multicast fe80::1-fe80::ffff

  # Check for file.
  $ %[1]s tool is_link_local_multicast /path/to/file.txt 

  # Check entries from stdin.
  $ cat /path/to/file1.txt | %[1]s tool is_link_local_multicast

Options:
  --help, -h
    show help.
  --quiet, -q
    quiet mode; suppress additional output.
`, progBase)
}

func cmdToolIsLinkLocalMulticast() (err error) {
	f := lib.CmdToolIsLinkLocalMulticastFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIsLinkLocalMulticast(f, pflag.Args()[2:], printHelpToolIsLinkLocalMulticast)
}
