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
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpToolIsV6() {
	fmt.Printf(
		`Usage: %s tool is_v6 <ip>

Description:
  check if IPs are IPv6.

Examples:
  $ %[1]s tool is_v6 "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
  $ %[1]s tool is_v6 "2001:0db8:85a3::8a2e:0370:7334"
  $ %[1]s tool is_v6 "::7334"

Options:
  General:
    --help, -h
      show help.
`, progBase)
}

func cmdToolIsV6() error {
	f := lib.CmdToolIsV6Flags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIsV6(f, pflag.Args()[2:], printHelpToolIsV6)
}
