package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsRange2IP = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpRange2IP() {
	fmt.Printf(
		`Usage: %s range2ip [<opts>] <ip-range | filepath>

Description:
  Accepts IP ranges and file paths to files containing IP ranges, converting
  them all to individual IPs within those ranges.

  $ %[1]s range2ip 8.8.8.0-8.8.8.255

  IP ranges can be of the form "<start><sep><end>" where "<sep>" can be "," or
  "-", and "<start>" and "<end>" can be any 2 IPs; order does not matter.


Options:
  --help, -h
    show help.
`, progBase)
}

func cmdRange2IP() error {
	f := lib.CmdRange2IPFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdRange2IP(f, pflag.Args()[1:], printHelpRange2IP)
}
