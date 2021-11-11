package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsRandIP = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-n":      predict.Nothing,
		"--count": predict.Nothing,
		"-4":      predict.Nothing,
		"--ipv4":  predict.Nothing,
		"-6":      predict.Nothing,
		"--ipv6":  predict.Nothing,
	},
}

func printHelpRandIP() {
	fmt.Printf(
		`Usage: %s randip [<opts>]

Description:
  Generates random IP/IPs.
  By default, generates 1 random IPv4 address, but can be configured to generate
  any number of a combination of IPv4/IPv6 addresses.

  $ %[1]s randip -t ipv4 -n 5

Options:
  --help, -h
    show help.
  --count, -n 
    number of IPs to generate.
  --type, -t
    type of IP to generate IPv4/IPv6
`, progBase)
}

func cmdRandIP() error {
	f := lib.CmdRandIPFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdRandIP(f, pflag.Args()[1:], printHelpRandIP)
}
