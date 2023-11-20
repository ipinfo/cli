package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolIsGlobalUnicast = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-q":      predict.Nothing,
		"--quiet": predict.Nothing,
	},
}

func printHelpToolIsGlobalUnicast() {
	fmt.Printf(`
%s tool is_global_unicast [<opts>] <cidr | ip | ip-range | filepath>

Description: Checks the provided address is a Global Unicast Address
Inputs can be IPs, IP ranges, CIDRs, or filepath to a file

Examples
  #IPv4/IPv6 Address.
  $ %[1]s tool is_global_unicast 10.255.0.0 | 2000::1
  $ %[1]s tool is_global_unicast 255.255.255.255 | ff00::1

  #IPv4/IPv6 Address Range
  $ %[1]s tool is_global_unicast 10.0.0.1-10.8.95.6 | 2000::1-2000::ffff
  $ %[1]s tool is_global_unicast 0.0.0.0-0.255.95.6 | ff00::1-ff00::ffff

  #Check CIDR
  $ %[1]s tool is_global_unicast 10.255.0.1/32 | 2000::1/64
  $ %[1]s tool is_global_unicast 154.0.0.1/32 | ff00::2540/64

  # Check for file.
  $ %[1]s tool is_global_unicast /path/to/file.txt 

  # Check entries from stdin.
  $ cat /path/to/file1.txt | %[1]s tool is_global_unicast

  Options:
  --help, -h
    show help.
  --quiet, -q
    quiet mode; suppress additional output.
`, progBase)
}

func cmdToolisGlobalUnicast() (err error) {
	f := lib.CmdToolIsGlobalUnicastFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIsGlobalUnicast(f, pflag.Args()[2:], printHelpToolIsGlobalUnicast)
}
