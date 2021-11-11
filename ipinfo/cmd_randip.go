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
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
		"-n":     predict.Nothing,
		"-count": predict.Nothing,
		"-t":     predict.Nothing,
		"-type":  predict.Nothing,
	},
}

func printHelpRandIP() {
	fmt.Printf(
		`Usage: %s randip [<opts>]

Description:
  Generates random IP/IPs.
  Generates 1 IPv4 without providing any arguments. 

  $ %[1]s randip -t ipv4 -n 5

Options:
  --help, -h
    show help.
  --count, -n 
    number of IPs to generate
  --type, -t
    type of IP to generate IPv4/IPv6
`, progBase)
}

func cmdRandIP() error {
	f := lib.CmdRandIPFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdRandIP(f, pflag.Args(), printHelpRandIP)
}
