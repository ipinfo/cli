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
	fmt.Printf(`
%s tool is_link_local_unicast [<opts>] <cidr | ip | ip-range | filepath>

Description: Checks the provided address is a Link Local Unicast
Inputs can be IPs, IP ranges, CIDRs, or filepath to a file

Examples
  #IPv4/IPv6 Address.
  $ %[1]s tool is_link_local_unicast 169.254.0.0 | fe80::
  $ %[1]s tool is_link_local_unicast 224.200.0.0 | 2000::

  #IPv4/IPv6 Address Range
  $ %[1]s tool is_link_local_unicast 169.254.0.1-169.254.255.255 | fe80::-fe80::ffff
  $ %[1]s tool is_link_local_unicast 168.254.0.1-169.254.255.255 | 2000::-2000::ffff

  #Check CIDR
  $ %[1]s tool is_link_local_unicast 169.0.0.1/32 | fe80::1/64
  $ %[1]s tool is_link_local_unicast 127.0.0.1/32 | 2000::1/64

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
