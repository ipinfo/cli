package main

import (
	"fmt"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolIsV4 = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpToolIsV4() {
	fmt.Printf(
		`Usage: %s tool is_v4 <ip>

Description:
  check if IPs are IPv4.

Examples:
  $ %[1]s tool is_v4 "192.168.1.1"
  $ %[1]s tool is_v4 "1.1.1.1"

Options:
  General:
    --help, -h
      show help.
`, progBase)
}

func cmdToolIsV4() error {
	f := lib.CmdToolIsV4Flags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIsV4(f, pflag.Args()[2:], printHelpToolIsV4)
}
